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
	query := "INSERT INTO epi_entries (eph_id, location_id, time_dongle, time_beacon, beacon_id) VALUES %s AS new ON DUPLICATE KEY UPDATE time_dongle=new.time_dongle, time_beacon=new.time_beacon;"
	values := types.ConcatEntries(input)
	statement := fmt.Sprintf(query, values)
	_, err := db.ExecContext(ctx, statement)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
