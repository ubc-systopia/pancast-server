package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	serverutils "pancast-server/server-utils"
	"pancast-server/types"
)

func CreateRiskEntries(input []types.Entry, db *sql.DB) bool {
	query := "INSERT INTO risk_entries (eph_id, location_id, time_beacon, time_dongle, beacon_id) VALUES %s AS new ON DUPLICATE KEY UPDATE time_dongle = new.time_dongle, time_beacon = new.time_beacon;"
	values := types.ConcatEntries(input)
	statement := fmt.Sprintf(query, values)
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Failed to begin transaction")
		return false
	}
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
	currentTime := serverutils.GetCurrentMinuteStamp()
	query := fmt.Sprintf("SELECT eph_id FROM risk_entries AS R, device AS D WHERE "+
		"D.device_id = R.beacon_id AND "+
		"%d - D.clock_init - D.clock_offset - R.time_beacon <= %d",
		currentTime, serverutils.MINUTES_IN_14_DAYS)
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error obtaining risk ephemeral IDs")
	}
	return rows
}

func GetNumOfRecentRiskEphIDs(db *sql.DB) int {
	currentTime := serverutils.GetCurrentMinuteStamp()
	query := fmt.Sprintf("SELECT COUNT(eph_id) FROM risk_entries AS R, device AS D WHERE "+
		"D.device_id = R.beacon_id AND "+
		"%d - D.clock_init - D.clock_offset - R.time_beacon <= %d",
		currentTime, serverutils.MINUTES_IN_14_DAYS)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return 0
	} else {
		var output int
		rows.Next()
		err = rows.Scan(&output)
		if err != nil {
			log.Println(err)
			return 0
		}
		return output
	}
}
