package models

import (
	"context"
	"database/sql"
	"pancast-server/types"
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

func GetLowestAvailableDeviceID(db *sql.DB) (uint32, error) {
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
	for idx, el := range freeArr {
		if el > uint32(idx) {
			return uint32(idx), nil
		}
	}
	if len(freeArr) < 1e10 {
		return uint32(len(freeArr)), nil
	}
	return 0, err // too many devices
}
