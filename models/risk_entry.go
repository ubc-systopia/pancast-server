package models

import (
	"context"
	"database/sql"
	"log"
)

type RiskEntry struct {
	EphemeralID []byte
	DongleClock uint64
	BeaconClock uint64
	BeaconID    uint64
	LocationID  string
}

func CreateRiskEntry(input RiskEntry, db *sql.DB) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := "INSERT INTO risk_entries VALUES (?, ?, ?, ?, ?)"
	_, err = tx.ExecContext(ctx, query, input.EphemeralID, input.LocationID, input.DongleClock,
		input.BeaconClock, input.BeaconID)
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}