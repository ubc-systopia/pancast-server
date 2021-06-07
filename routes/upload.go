package routes

import "database/sql"

type UploadParameters struct {
	EphemeralID []byte
	DongleClock uint64
	BeaconClock uint64
	BeaconID    uint64
	LocationID  string
	Type        int
}

func UploadController(input UploadParameters, db *sql.DB) {

}
