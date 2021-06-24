package server_utils

import (
	"database/sql"
	cuckoo "github.com/panmari/cuckoofilter"
)

func deleteOutdatedKeys(cf *cuckoo.Filter, db *sql.DB) {
	// TODO: cron job for deleting outdated keys from cuckoo filter
}
