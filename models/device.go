package models

import (
	"context"
	"database/sql"
	"log"
)

type Device struct {
	DeviceID    int
	SecretKey   []byte
	ClockInit   int
	ClockOffset int
}

type Beacon struct {
	Self       Device
	LocationID []byte
}

type Dongle struct {
	Self Device
	OTPs [][]byte
}

func CreateDevice(deviceType int, db *sql.DB) {
	ctx := context.Background()
	_, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_ = "INSERT INTO device VALUES "
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
		if el > uint32(idx) + 1 {
			return uint32(idx) + 1, nil
		}
	}
	return 0, err
}
