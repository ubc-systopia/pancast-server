package models

import (
	"context"
	"database/sql"
	"pancast-server/types"
)

func CreateBeacon(b types.BeaconRegistrationData, tx *sql.Tx, ctx context.Context) error {
	query := "INSERT INTO beacon VALUES (?,?)"
	_, err := tx.ExecContext(ctx, query, b.DeviceData.DeviceID, b.LocationID)
	if err != nil {
		return err
	}
	return nil
}
