package cronjobs

/*
	This file contains all the functions which need to be scheduled periodically.
	As of now, it only contains the cuckoo filter (re-)initialization procedure.
*/

import (
	"database/sql"
	"log"
	"os"
//	"fmt"
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

func GenerateEphemeralIDList(db *sql.DB, params DiffprivParameters, mode []string) ([][]byte, int) {
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
				return nil, 0
			}
			ephIDs = append(ephIDs, ephID)
		}
	}
	// generating a known amount of entries
	if server_utils.StringSliceContains(mode, "CUCKOO_FIXED_ITEMS") {
		fixedCount := 10 // 95% of 512
		for i := 0; i < fixedCount; i++ {
			dummy, err := server_utils.GenerateRandomByteString(15)
			if err != nil {
				log.Println(err)
				break
			}
			ephIDs = append(ephIDs, dummy)
		}
		length += fixedCount
		log.Printf("(DEV) Fixed number of ephemeral IDs generated: %d", fixedCount)
	}

	// sample random variable from Laplacian distribution, and create dummy ephemeral IDs
	if !server_utils.StringSliceContains(mode, "CUCKOO_NODIFF") {
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

	if server_utils.StringSliceContains(mode, "CUCKOO_LOG_EPHIDS") {
		f, _ := os.Create("ephid_list.txt")
		_, _ = f.WriteString(server_utils.ByteSlicesToHexString(ephIDs))
	}
	return server_utils.ShuffleByteArray(ephIDs), length
}

func CreateNewFilter(ephIDs [][]byte, length int) (*cuckoo.Filter, error) {

	numBuckets := length / 4 // division by 4 because each bucket can hold 4 ephemeral IDs
	numBuckets = server_utils.NextPowerOfTwo(numBuckets)

	// tries to create a filter, doubling in size if not possible, and ultimately terminating
	// once the filter becomes too big to be feasibly transferred.
	filter, err := server_utils.AllocateFilter(numBuckets, ephIDs)
	if err != nil {
		log.Println(err)
		filter, _ = cuckoo.CreateFilter(4) // dummy var
		return filter, nil
	}
	return filter, nil
}


func CreateChunkedFilters(ephIDs [][]byte, length int) ([]*cuckoo.Filter, error) {
	var chunks []*cuckoo.Filter
	if length == 0 {
		// no cuckoo filter will be created at all
		return chunks, nil
	}
	for length > 0 {
		// create a chunk
		currentEphIDInChunk := server_utils.MAX_CHUNK_EPHID_COUNT
		if currentEphIDInChunk > len(ephIDs) {
			currentEphIDInChunk = len(ephIDs)
		}
		for currentEphIDInChunk != 0 {
			// creates a filter
			filter, err := cuckoo.CreateFilter(server_utils.MAX_CHUNK_EPHID_COUNT / cuckoo.BUCKET_SIZE)
			if err != nil {
				return nil, err
			}
			// tries to insert currentEphIDInChunk ephIDs
			success := true
			for _, ephID := range ephIDs[:currentEphIDInChunk] {
				result := filter.Insert(ephID)
				if !result {
					success = false
				}
			}
			if success {
				// we were able to fit in this many ephemeral IDs in the filter,
				// take out the first currentEphIDInChunk ephIDs and loop again
				ephIDs = ephIDs[currentEphIDInChunk:]
				chunks = append(chunks, filter)
				break
			}
			// failed, so we try again with 1 less ephID
			currentEphIDInChunk -= 1
		}
		length -= currentEphIDInChunk
	}
	log.Printf("%d chunks created\r\n", len(chunks))
	return chunks, nil
}
