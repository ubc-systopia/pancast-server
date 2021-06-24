package cronjobs

import (
	"database/sql"
	"log"
	"pancast-server/cuckoo"
	"pancast-server/models"
	"pancast-server/server-utils"
)

func CreateNewFilter(cf *cuckoo.Filter, db *sql.DB) {
	rows := models.GetRiskEphIDs(db)
	length := models.GetNumOfRecentRiskEphIDs(db)
	// division by 4 is because there are 4 entries per bucket
	baseLength := server_utils.NextPowerOfTwo(length) / 4
	var ephIDs [][]byte
	for rows.Next() {
		var ephID []byte
		err := rows.Scan(&ephID)
		if err != nil {
			log.Fatal(err)
		}
		ephIDs = append(ephIDs, ephID)
	}
	if baseLength < 4 {
		baseLength = 4
	}
	filter, err := server_utils.AllocateFilter(baseLength, ephIDs)
	if err != nil {
		log.Println(err)
		cf, _ = cuckoo.CreateFilter(4) // dummy var
		return
	}
	cf = filter
	return
}
