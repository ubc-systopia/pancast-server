package cronjobs

/*
	This file contains all the functions which need to be scheduled periodically.
	As of now, it only contains the cuckoo filter (re-)initialization procedure.
 */

import (
	"database/sql"
	"log"
	"math"
	"os"
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

func CreateNewFilter(db *sql.DB, params DiffprivParameters, mode []string) (*cuckoo.Filter, error) {
	var ephIDs [][]byte
	length := 0
	// get data to be broadcast

	if server_utils.StringSliceContains(mode, "CUCKOO_USE_DATABASE") {
		rows := models.GetRiskEphIDs(db)
		length += models.GetNumOfRecentRiskEphIDs(db)
		// division by 4 is because there are 4 entries per bucket, and therefore
		// the filter can hold 4 * baseLength entries
		for rows.Next() {
			var ephID []byte
			err := rows.Scan(&ephID)
			if err != nil {
				return nil, err
			}
			ephIDs = append(ephIDs, ephID)
		}
	}

	// generating a known amount of entries
	if server_utils.StringSliceContains(mode, "CUCKOO_FIXED_ITEMS") {
		fixedCount := 500
		for i := 0; i < fixedCount; i++ {
			dummy, err := server_utils.GenerateRandomByteString(15)
			if err != nil {
				log.Println(err)
				break
			}
			ephIDs = append(ephIDs, dummy)
		}
		listOfEphemeralIDs := server_utils.ByteSlicesToHexString(ephIDs)
		f, _ := os.Create("ephid_list.txt")
		_, _ = f.WriteString(listOfEphemeralIDs)
		length += fixedCount
		log.Printf("(DEV) Fixed number of ephemeral IDs generated: %d", fixedCount)
	}



	// sample random variable from Laplacian distribution, and create dummy ephemeral IDs
	if server_utils.StringSliceContains(mode, "CUCKOO_NODIFF") {
		junkCount := server_utils.SampleLaplacianDistribution(params.Mean, params.Sensitivity, params.Epsilon, params.Delta)
		for i := int64(0); i < junkCount; i++ {
			dummy, err := server_utils.GenerateRandomByteString(15)
			if err != nil {
				log.Println(err)
				break
			}
			ephIDs = append(ephIDs, dummy)
		}
		length += int(junkCount)
		log.Printf("Number of ephemeral IDs generated: %d", junkCount)
	}

	projectedSizeInBytes := int(math.Floor(float64(length * cuckoo.FINGERPRINT_BITS) / 8))
	baseLength := server_utils.LowerPowerOfTwo(projectedSizeInBytes)
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
