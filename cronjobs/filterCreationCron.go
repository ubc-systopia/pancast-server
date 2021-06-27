package cronjobs

import (
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

func CreateNewFilter(db *sql.DB, params DiffprivParameters) (*cuckoo.Filter, error) {
	// get data to be broadcast
	rows := models.GetRiskEphIDs(db)
	length := models.GetNumOfRecentRiskEphIDs(db)
	// division by 4 is because there are 4 entries per bucket, and therefore
	// the filter can hold 4 * baseLength entries
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
	// sample random variable from Laplacian distribution, and create dummy ephemeral IDs
	junkCount := server_utils.SampleLaplacianDistribution(params.Mean, params.Sensitivity, params.Epsilon, params.Delta)
	for i := int64(0); i < junkCount; i++ {
		dummy, err := server_utils.GenerateRandomByteString(15)
		if err != nil {
			log.Println(err)
			break
		}
		ephIDs = append(ephIDs, dummy)
	}
	if baseLength < 4 {
		baseLength = 4
	}
	log.Printf("Number of ephemeral IDs generated: %d", junkCount)
	// tries to create a filter, doubling in size if not possible, and ultimately terminating
	// once the filter becomes too big to be feasibly transferred.
	filter, err := server_utils.AllocateFilter(baseLength, server_utils.ShuffleByteArray(ephIDs))
	if err != nil {
		log.Println(err)
		filter, _ = cuckoo.CreateFilter(4) // dummy var
		return filter, err
	}
	return filter, err
}
