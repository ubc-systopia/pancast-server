package server_utils

import (
	"database/sql"
	"fmt"
	cuckoo "github.com/panmari/cuckoofilter"
	"log"
	"net/http"
	"pancast-server/models"
)

func Write(res http.ResponseWriter, output string) {
	_, err := fmt.Fprint(res, output)
	if err != nil {
		log.Fatal(err)
	}
}

func PopulateCuckooFilter(db *sql.DB) *cuckoo.Filter {
	cf := cuckoo.NewFilter(1000000)
	rows := models.GetRiskEphIDs(db)
	for rows.Next() {
		var ephID []byte
		err := rows.Scan(&ephID)
		if err != nil {
			log.Fatal(err)
		}
		cf.Insert(ephID)
	}
	return cf
}
