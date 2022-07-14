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
		log.Printf("Number of real ephIDs: %d", length)
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
	log.Printf("filter created, #ephids: %d, len: %d, #buckets: %d",
			len(ephIDs), length, numBuckets)
	return filter, nil
}


func CreateChunkedFilters(ephIDs [][]byte, length int) ([]*cuckoo.Filter, error) {
	var chunks []*cuckoo.Filter
	num_ephids := length
	var chunkId = 0

	if length == 0 {
		// no cuckoo filter will be created at all
		return chunks, nil
	}

	type ephIDStatus struct {
		ephID []byte
		result bool
		idxType int
		idxVal uint
		fp	cuckoo.Fingerprint
	}

	var ephIDStatusArray []ephIDStatus = nil
	var numAttempts = 1

	for length > 0 {
		// create a chunk
		currentEphIDInChunk := server_utils.MAX_CHUNK_EPHID_COUNT
		if currentEphIDInChunk > len(ephIDs) {
			currentEphIDInChunk = len(ephIDs)
		}

//		startEphIDInChunk := currentEphIDInChunk

		for currentEphIDInChunk != 0 {
			// creates a filter
			filter, err := cuckoo.CreateFilter(
					server_utils.MAX_CHUNK_EPHID_COUNT / cuckoo.BUCKET_SIZE)
			if err != nil {
				return nil, err
			}

			// tries to insert currentEphIDInChunk ephIDs
			success := true
			for _, ephID := range ephIDs[:currentEphIDInChunk] {
				result, idxType, idxVal, fp := filter.Insert(ephID)
				if !result {
//					log.Printf("[%d] failed %d", chunkId, numAttempts)
					success = false
					ephIDStatusArray = nil
					break
				}
				ephIDStatusArray = append(ephIDStatusArray, ephIDStatus{ephID, result, idxType, idxVal, fp})
			}

			if success {
				// we were able to fit in this many ephemeral IDs in the filter,
				// take out the first currentEphIDInChunk ephIDs and loop again
				for _, es := range ephIDStatusArray {
					log.Printf("chunk [%d/%d] #attempts: %d ephID %x idx[%d]: %d fp %08x",
							chunkId, len(ephIDStatusArray), numAttempts, es.ephID, es.idxType, es.idxVal, es.fp)
				}
				ephIDStatusArray = nil
				numAttempts = 1

				ephIDs = ephIDs[currentEphIDInChunk:]
				chunks = append(chunks, filter)

				break
			}

			// failed, so we try again with 1 less ephID
			numAttempts += 1
			currentEphIDInChunk -= 1
		}

//		log.Printf("[%d] #ephids left: %d, curr chunk size: %d -> %d", chunkId, len(ephIDs), startEphIDInChunk, currentEphIDInChunk)
		chunkId += 1

		length -= currentEphIDInChunk
	}
	log.Printf("#ephids: %d, bucket size: %d, max chunk size: %d max #chunks: %d, #chunks created: %d\r\n",
			num_ephids, cuckoo.BUCKET_SIZE, server_utils.MAX_CHUNK_EPHID_COUNT,
			server_utils.MAX_CHUNK_EPHID_COUNT/cuckoo.BUCKET_SIZE, len(chunks))

	if (len(chunks) >  server_utils.MAX_CHUNK_EPHID_COUNT/cuckoo.BUCKET_SIZE) {
	  return nil, nil
	}

	return chunks, nil
}
