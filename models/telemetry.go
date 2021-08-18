package models

import "database/sql"

func CreateTelemetryEntry(time string, path string, numEntries int, timestamp uint32, db *sql.DB) {
	query := "INSERT INTO telemetry (time_taken, route, num_ids_submitted, time_stamp) VALUES (?, ?, ?, ?);"
	_, _ = db.Exec(query, time, path, numEntries, timestamp)
}
