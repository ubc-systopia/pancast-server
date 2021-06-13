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
