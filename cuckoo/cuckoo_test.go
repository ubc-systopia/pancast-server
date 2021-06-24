package cuckoo

import (
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"reflect"
	"testing"
)

func TestFilterCreation(t *testing.T) {
	cf, err := CreateFilter(TEST_NUM_BUCKETS)
	if err != nil {
		t.Fatal("filter creation error")
	}
	if len(cf.Buckets) != TEST_NUM_BUCKETS {
		t.Fatal("incorrect number of Buckets")
	}
	for _, bucket := range cf.Buckets {
		if len(bucket.Fp) != BUCKET_SIZE {
			t.Fatal("incorrect bucket size")
		}
		for _, fp := range bucket.Fp {
			if reflect.TypeOf(fp).Kind() != reflect.Uint32 {
				t.Fatal("incorrect type")
			}
			if fp != 0 {
				t.Fatal("element not initialized to 0")
			}
		}
	}
}

func TestCreateFilterZeroBuckets(t *testing.T) {
	_, err := CreateFilter(TEST_ZERO_NUM_BUCKETS)
	if err == nil {
		t.Fatal("zero Buckets not allowed")
	}
}

func TestCreateFilterBadBucketCount(t *testing.T) {
	_, err := CreateFilter(TEST_BAD_NUM_BUCKETS)
	if err == nil {
		t.Fatal("non-power of 2 Buckets not allowed")
	}
}

func TestFilterInsert(t *testing.T) {
	cf, err := CreateFilter(TEST_NUM_BUCKETS)
	if err != nil {
		t.Fatal("filter creation error")
	}
	item := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	result := cf.Insert(item)
	if !result {
		t.Fatal("insertion error")
	}
}

func TestFilterLookup(t *testing.T) {
	cf, err := CreateFilter(TEST_NUM_BUCKETS)
	item := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	badItem := []byte{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}
	if err != nil {
		t.Fatal("filter creation error")
	}
	result := cf.Insert(item)
	if !result {
		t.Fatal("insertion error")
	}
	t.Run("regular lookup", func(t *testing.T) {
		lookupResult := cf.Lookup(item)
		if !lookupResult {
			t.Fatal("could not look up previously inserted item")
		}
	})

	t.Run("bad lookup", func(t *testing.T) {
		lookupResult := cf.Lookup(badItem)
		if lookupResult {
			t.Fatal("item should not have been found")
		}
	})
}

func TestFilterDelete(t *testing.T) {
	cf, err := CreateFilter(TEST_NUM_BUCKETS)
	if err != nil {
		t.Fatal("filter creation error")
	}
	item := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	result := cf.Insert(item)
	if !result {
		t.Fatal("insertion error")
	}
	deleteResult := cf.Delete(item)
	if !deleteResult {
		t.Fatal("could not look up previously inserted item")
	}
	badDeleteResult := cf.Delete(item)
	if badDeleteResult {
		t.Fatal("delete should not have returned true")
	}
}

func TestFilterEncode(t *testing.T) {
	assertVal := assert.New(t)
	cf, err := CreateFilter(TEST_NUM_BUCKETS)
	if err != nil {
		t.Fatal("filter creation error")
	}
	item := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	cf.Insert(item)
	byteArray := cf.Encode()
	assertVal.Equal((TEST_NUM_BUCKETS*BUCKET_SIZE*FINGERPRINT_BITS)/8, len(byteArray))
}

func TestFilterDecode(t *testing.T) {
	assertVal := assert.New(t)
	cf, err := CreateFilter(TEST_NUM_BUCKETS)
	if err != nil {
		t.Fatal("filter creation error")
	}
	item := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	cf.Insert(item)
	byteArray := cf.Encode()
	newFilter, err := Decode(byteArray)
	if err != nil {
		t.Fatal("error decoding byte array")
	}
	assertVal.Equal(TEST_NUM_BUCKETS, len(newFilter.Buckets))
	lookupResult := newFilter.Lookup(item)
	if !lookupResult {
		t.Fatal("item should exist within filter")
	}
}

func TestFilterMassDecode(t *testing.T) {
	NumCases := 300
	cf, err := CreateFilter(TEST_NUM_BUCKETS << 2)
	if err != nil {
		t.Fatal("filter creation error")
	}
	testCases := make([][]byte, NumCases)
	for i, _ := range testCases {
		testCases[i] = make([]byte, 15)
		rand.Read(testCases[i])
		cf.Insert(testCases[i])
	}
	byteArray := cf.Encode()
	newCF, err := Decode(byteArray)
	if err != nil {
		t.Fatal("decode error")
	}
	for _, c := range testCases {
		result := newCF.Lookup(c)
		if !result {
			t.Fatal("item should have been found")
		}
	}
}

func TestFilterBatchInsertAndLookup(t *testing.T) {
	NumCases := 300
	insertFailCount := 0
	lookupFailCount := 0
	erroneousLookupFailCount := 0
	cf, err := CreateFilter(TEST_NUM_BUCKETS << 2)
	if err != nil {
		t.Fatal("filter creation error")
	}
	// create entries to put into filter
	testCases := make([][]byte, NumCases)
	for i, _ := range testCases {
		testCases[i] = make([]byte, 15)
		rand.Read(testCases[i])
	}
	// test insertions
	for _, c := range testCases {
		result := cf.Insert(c)
		if !result {
			insertFailCount++
		}
	}
	// test real lookups
	for _, c := range testCases {
		result := cf.Lookup(c)
		if !result {
			lookupFailCount++
		}
	}
	// test lookups that don't exist
	wrongCases := make([][]byte, NumCases)
	for _, c := range wrongCases {
		c = make([]byte, 15)
		rand.Read(c)
	}
	for _, c := range wrongCases {
		result := cf.Lookup(c)
		if result {
			erroneousLookupFailCount++
		}
	}
	log.Printf("Insert failed %d times\n", insertFailCount)
	log.Printf("Lookup failed %d times\n", lookupFailCount)
	log.Printf("Encountered false positive %d times\n", erroneousLookupFailCount)

}
