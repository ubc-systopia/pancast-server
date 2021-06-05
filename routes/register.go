package routes

import (
	"database/sql"
)

func RegisterController(deviceType int, db *sql.DB) {
	//mockDevice := server.BLE_NETWORKED
	// what type of device are we dealing with? mocking a networked beacon for now
	// TODO: Compute a secret key to give to a beacon

	// insert device into database

	// handle registration logic
	// tons of crypto
	// maybe some database access
}
