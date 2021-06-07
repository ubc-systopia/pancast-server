package routes

import (
	"database/sql"
	"os"
	"pancast-server/database"
	"testing"
)

var (
	db *sql.DB
)

func setup() {
	db = database.InitDatabaseConnection()
	return
}

func teardown() {
	db.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestRegistrationDongleSuccess(t *testing.T) {
	_, err := RegisterController(0, db)
	if err != nil {
		t.Fail()
	}
}

func TestRegistrationBeaconSuccess(t *testing.T) {
	_, err := RegisterController(1, db)
	if err != nil {
		t.Fail()
	}
}
