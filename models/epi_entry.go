package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pancast-server/types"
)

func CreateEpiEntries(input []types.Entry, db *sql.DB) bool {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	query := "INSERT INTO epi_entries VALUES %s;"
	values := types.ConcatEntries(input)
	statement := fmt.Sprintf(query, values)
	_, err = tx.ExecContext(ctx, statement)
	if err != nil {
		_ = tx.Rollback()
		log.Println(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

