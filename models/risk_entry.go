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
	for i := 0; i < len(input); i++ {
		query :=
      `INSERT INTO risk_entries (eph_id, location_id, time_beacon, time_dongle, beacon_id)
       VALUES %s AS new
       ON DUPLICATE KEY UPDATE time_dongle = new.time_dongle, time_beacon = new.time_beacon;`
		log.Printf("PREP QUERY === %d", i)
		log.Printf("%s", input[i])
		values := types.ConcatEntries(input[i:i+1])
		statement := fmt.Sprintf(query, values)
		log.Printf("STMT %d ===\n%s", i, statement)
		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Println("Failed to begin transaction")
			return false
		}

		_, err = tx.ExecContext(ctx, statement)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return false
		}

		err = tx.Commit()
		if err != nil {
			return false
		}
	}

	return true
}

func GetRiskEphIDs(db *sql.DB) *sql.Rows {
	currentTime := serverutils.GetCurrentMinuteStamp()
	query := `SELECT HEX(eph_id) FROM risk_entries AS R, device AS D
	WHERE D.device_id = R.beacon_id AND (%d - D.clock_init - D.clock_offset - R.time_beacon) <= %d`
	statement := fmt.Sprintf(query,	currentTime, serverutils.MINUTES_IN_14_DAYS)
	rows, err := db.Query(statement)
	if err != nil {
		log.Printf("Error obtaining risk ephemeral IDs: %s", err)
		return rows
	}

	rows2, err := db.Query(statement)
	count := 0
	for rows2.Next() {
		count += 1
	}
	log.Printf("QUERY res len: %d ===\n%s", count, statement)

	return rows
}

func GetNumOfRecentRiskEphIDs(db *sql.DB) int {
	currentTime := serverutils.GetCurrentMinuteStamp()
	query := `SELECT COUNT(eph_id) FROM risk_entries AS R, device AS D
	WHERE D.device_id = R.beacon_id AND (%d - D.clock_init - D.clock_offset - R.time_beacon) <= %d`
	statement := fmt.Sprintf(query, currentTime, serverutils.MINUTES_IN_14_DAYS)
	rows, err := db.Query(statement)

	log.Printf("QUERY ===\n%s", statement)
	if err != nil {
		log.Println(err)
		return 0
	}

	var output int
	rows.Next()
	err = rows.Scan(&output)
	if err != nil {
		log.Println(err)
		return 0
	}

	return output
}
