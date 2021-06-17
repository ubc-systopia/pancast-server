package routes

import (
	"database/sql"
	cuckoo "github.com/panmari/cuckoofilter"
)

func UpdateController(cf *cuckoo.Filter, db *sql.DB) []byte {
	// nothing complicated here :)
	return cf.Encode()
}
