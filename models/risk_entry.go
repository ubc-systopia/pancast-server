package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pancast-server/types"
)

func CreateRiskEntries(input []types.Entry, db *sql.DB) bool {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	query := "INSERT INTO risk_entries VALUES %s;"
	values := types.ConcatEntries(input)
	statement := fmt.Sprintf(query, values)
	_, err = tx.ExecContext(ctx, statement)
	log.Println(err)
	if err != nil {
		_ = tx.Rollback()
		return false
	}
	err = tx.Commit()
	if err != nil {
		return false
	}
	return true
}

func GetRiskEphIDs(db *sql.DB) *sql.Rows {
	rows, err := db.Query("SELECT eph_id FROM risk_entries;")
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
