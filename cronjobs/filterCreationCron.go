package cronjobs

import (
	cryptorand "crypto/rand"
	"database/sql"
	"log"
	"pancast-server/cuckoo"
	"pancast-server/models"
	"pancast-server/server-utils"
)

type DiffprivParameters struct {
	Mean        int64
	Sensitivity float64
	Epsilon     float64
	Delta       float64
}

func CreateNewFilter(cf *cuckoo.Filter, db *sql.DB, params DiffprivParameters) {
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
	junkCount := server_utils.SampleLaplacianDistribution(params.Mean, params.Sensitivity, params.Epsilon, params.Delta)
	for i := int64(0); i < junkCount; i++ {
		dummy := make([]byte, 15)
		_, err := cryptorand.Read(dummy)
		if err != nil {
			log.Println(err)
			break
		}
		ephIDs = append(ephIDs, dummy)
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
