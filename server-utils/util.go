package server_utils

/*
	Utilities file. Contains all sorts of utility functions
*/

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
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

func LowerPowerOfTwo(num int) int {
	logVal := math.Log2(float64(num))
	exponent := math.Ceil(logVal)
	return int(math.Pow(2, exponent))
}

func GetIncrementUnit(num int, granularity int) int {
	return num
}

func NextPowerOfTwo(num int) int {
	exponent := math.Floor(math.Log2(float64(num)))
	return int(math.Pow(2, exponent+1))

}

func AllocateFilter(initNumBuckets int, ephIDs [][]byte) (*cuckoo.Filter, error) {
//	log.Println(initNumBuckets)
	if initNumBuckets > int(math.Pow(2, EXPONENT_TOO_LARGE)) {
		return nil, errors.New("filter has grown too large")
	}
	filter, err := cuckoo.CreateFilter(initNumBuckets)
	if err != nil {
		return nil, err
	}
	for _, ephID := range ephIDs {
		result, _, _, _ := filter.Insert(ephID)
		if !result {
			return AllocateFilter(NextPowerOfTwo(initNumBuckets), ephIDs)
		}
	}
	// success
	return filter, nil
}
func GenerateRandomByteString(n int) ([]byte, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return key, err
	}
	return key, nil
}

func SampleLaplacianDistribution(mean int64, sensitivity float64, epsilon float64, delta float64) int64 {
	lambda := sensitivity / epsilon
	t := math.Ceil(lambda * math.Log((delta-1+math.Exp(sensitivity/lambda))/(2*delta)))
	randomBytes := make([]byte, 8)
	_, err := cryptorand.Read(randomBytes)
	if err != nil {
		log.Println("error generating laplacian random variable")
		return -1
	}
	randomNumberSource := rand.NewSource(binary.LittleEndian.Uint64(randomBytes))
	laplacianInstance := distuv.Laplace{
		Mu:    float64(mean),
		Scale: lambda,
		Src:   randomNumberSource,
	}
	randomVar := math.Floor(laplacianInstance.Rand())
	if t+randomVar < 0 {
		return 0
	} else {
		return int64(t + randomVar)
	}
}

func GetCurrentMinuteStamp() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Minute))
}

func ShuffleByteArray(array [][]byte) [][]byte {
	copiedArray := array
	currentIndex := len(copiedArray)
	var tempVal []byte
	randomIndex := 0
	for currentIndex != 0 {
		randomIndex = rand.Intn(len(copiedArray))
		currentIndex -= 1
		tempVal = copiedArray[currentIndex]
		copiedArray[currentIndex] = copiedArray[randomIndex]
		copiedArray[randomIndex] = tempVal
	}
	return copiedArray
}

func DecodeBase64ToByteArray(base64input string) ([]byte, error) {
	output, err := base64.StdEncoding.DecodeString(base64input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func ByteSlicesToHexString(byteArray [][]byte) string {
	finalString := ""
	for _, ephemeralID := range byteArray {
		finalString = finalString + hex.EncodeToString(ephemeralID) + "\n"
	}
	return finalString
}

func StringSliceContains(s []string, i string) bool {
	for _, item := range s {
		if item == i {
			return true
		}
	}
	return false
}
