package cuckoo

import (
	"math/rand"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestGetFingerprint(t *testing.T) {
	assertVar := assert.New(t)
	assertVar.Equal(uint32(1), GetFingerprint(0))
	assertVar.Equal(uint32(2), GetFingerprint(1))
	assertVar.Equal(uint32(134217727), GetFingerprint(134217726))
	assertVar.Equal(uint32(1), GetFingerprint(134217727))
	assertVar.Equal(uint32(65782274), GetFingerprint(200000000))
}

func TestGetIndexAndFingerprint(t *testing.T) {
	assertVar := assert.New(t)
	MaxLoops := 1000
	for i := 0; i < MaxLoops; i++ {
		item := make([]byte, 15)
		rand.Read(item)
		index, fp := GetIndexAndFingerprint(item, TEST_NUM_BUCKETS-1)
		assertVar.GreaterOrEqual(index, uint(0))
		assertVar.LessOrEqual(index, uint(TEST_NUM_BUCKETS-1))
		assertVar.Equal(int(fp>>27), 0)
	}
}

func TestGetAltIndex(t *testing.T) {
	assertVar := assert.New(t)
	MaxLoops := 1000
	for i := 0; i < MaxLoops; i++ {
		item := make([]byte, 15)
		rand.Read(item)
		index, fp := GetIndexAndFingerprint(item, TEST_NUM_BUCKETS-1)
		newIndex := GetAltIndex(fp, index, TEST_NUM_BUCKETS-1)
		assertVar.GreaterOrEqual(newIndex, uint(0))
		assertVar.LessOrEqual(newIndex, uint(TEST_NUM_BUCKETS-1))
	}
}

//func TestWriteBitsToPositionAndBitOffset(t *testing.T) {
//	assertVal := assert.New(t)
//	arr := make([]byte, 2)
//	arr = WriteByteToPositionAndBitOffset(arr, byte(170), 0, 0)
//	arr = WriteByteToPositionAndBitOffset(arr, byte(85), 1, 0)
//	el1 := ReadFingerprintFromPositionAndBitOffset(arr, 0, 0, 8)
//	el2 := ReadFingerprintFromPositionAndBitOffset(arr, 1, 0, 8)
//	assertVal.Equal(uint64(el1), uint64(170))
//	assertVal.Equal(uint64(el2), uint64(85))
//}

//func TestWriteHalfway(t *testing.T) {
//	arr := make([]byte, 8)
//	elementsToInsert := FingerprintToByteArray(22180501)
//	shiftedElements := ShiftByteArray(elementsToInsert)
//	arr = WriteByteToPositionAndBitOffset(arr, shiftedElements[0], 0, 4)
//	arr = WriteByteToPositionAndBitOffset(arr, shiftedElements[1], 1, 4)
//	arr = WriteByteToPositionAndBitOffset(arr, shiftedElements[2], 2, 4)
//	arr = WriteBitsToPositionAndBitOffset(arr, shiftedElements[3], 3, 4, 3)
//}
