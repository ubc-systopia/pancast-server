package models

import (
	"context"
	"database/sql"
	"pancast-server/types"
)

func CreateDongle(d types.DongleRegistrationData, tx *sql.Tx, ctx context.Context) error {
	query := "INSERT INTO dongle VALUES (?)"
	// placeholder: OTP not implemented yet
	_, err := tx.ExecContext(ctx, query, d.DeviceData.DeviceID)
	if err != nil {
		return err
	}
	return nil
}
