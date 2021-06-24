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
	query := "INSERT INTO epi_entries VALUES %s;"
	values := types.ConcatEntries(input)
	statement := fmt.Sprintf(query, values)
	_, err := db.ExecContext(ctx, statement)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
