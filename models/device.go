package models

import (
	"context"
	"database/sql"
	"errors"
	server_utils "pancast-server/server-utils"
	"pancast-server/types"
	"log"
	"strconv"
)

type Device struct {
	DeviceID    int
	SecretKey   []byte
	ClockInit   int
	ClockOffset int
}

func CreateDevice(r types.RegistrationData, tx *sql.Tx, ctx context.Context) error {
	query := "INSERT INTO device VALUES (?,?,?,?)"
	_, err := tx.ExecContext(ctx, query, r.DeviceID, r.Secret, r.Clock, 0)
	if err != nil {
		return err
	}
	return nil
}

func GetNextDeviceID(db *sql.DB, dType int64) (uint32, error) {
  var freeArr []uint32
  var qstring string

  if (dType == server_utils.BEACON) {
    qstring = "SELECT device_id FROM beacon ORDER BY device_id;"
  } else {
    qstring = "SELECT device_id FROM dongle ORDER BY device_id;"
  }

  rows, err := db.Query(qstring)
  if (err != nil) {
    return 0, err
  }

  for rows.Next() {
    var candidate uint32
    err = rows.Scan(&candidate)
    if (err != nil) {
      return 0, err
    }

    freeArr = append(freeArr, candidate)
  }

  if (dType == server_utils.BEACON && len(freeArr) >= server_utils.MAX_BEACON_IDS) {
    return 0, nil
  }

  var newElem uint32
  if (len(freeArr) <= 0) {
    newElem = 0
  } else {
    lastElem := freeArr[len(freeArr)-1]
    newElem = lastElem + 1
  }
  log.Println("dType: " + strconv.Itoa(int(dType)) + ", new ID: " + strconv.FormatInt(int64(newElem), 16))

  return newElem, nil
}

func GetLowestAvailableDeviceID(db *sql.DB, dType int64) (uint32, error) {
	var freeArr []uint32
	// assume sorted
	rows, err := db.Query("SELECT device_id FROM device ORDER BY device_id;")
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var candidate uint32
		err = rows.Scan(&candidate)
		if err != nil {
			return 0, err
		}
		freeArr = append(freeArr, candidate)
	}
	var id uint32
	if dType == server_utils.BEACON {
		if len(freeArr) < server_utils.MAX_BEACON_IDS {
			id = (server_utils.BROADCAST_SERVICE_UUID << 16) + uint32(len(freeArr))
		} else {
			err = errors.New("Max number of beacons registered")
		}
	} else {
		if len(freeArr) < 1e10 {
			id = uint32(len(freeArr))
		} else {
			err = errors.New("Max number of devices registered")
		}
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
