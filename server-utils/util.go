package server_utils

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"pancast-server/cuckoo"
	"time"
)

func Write(res http.ResponseWriter, output string) {
	_, err := fmt.Fprint(res, output)
	if err != nil {
		log.Fatal(err)
	}
}

func NextPowerOfTwo(num int) int {
	logVal := math.Log2(float64(num))
	exponent := math.Ceil(logVal)
	return int(math.Pow(2, exponent))
}

func AllocateFilter(initSize int, ephIDs [][]byte) (*cuckoo.Filter, error) {
	if initSize > int(math.Pow(2, EXPONENT_TOO_LARGE)) {
		return nil, errors.New("filter has grown too large")
	}
	filter, err := cuckoo.CreateFilter(initSize)
	if err != nil {
		return nil, err
	}
	for _, ephID := range ephIDs {
		result := filter.Insert(ephID)
		if !result {
			return AllocateFilter(initSize*2, ephIDs)
		}
	}
	// success
	return filter, nil
}

func GetCurrentMinuteStamp() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Minute))
}
