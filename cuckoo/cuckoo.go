package cuckoo

import (
	"errors"
	"log"
	"math"
	"math/rand"
)

type Filter struct {
	bucketMask uint
	buckets []Bucket
}

type Bucket struct {
	fp []Fingerprint
}

type Fingerprint uint32 // (FINGERPRINT_SIZE * 8 bits)

func CreateFilter(numBuckets int) (*Filter, error) {
	if numBuckets == 0 {
		return nil, errors.New("cannot have zero buckets")
	}
	if !IsPowerOfTwo(numBuckets) {
		return nil, errors.New("buckets is not a power of 2")
	}
	var tempFilter []Bucket
	for i := 0; i < numBuckets; i++ {
		var bkt Bucket
		for j := 0; j < BUCKET_SIZE; j++ {
			var fp Fingerprint
			fp = Fingerprint(0)
			bkt.fp = append(bkt.fp, fp)
		}
		tempFilter = append(tempFilter, bkt)
	}
	return &Filter{buckets: tempFilter, bucketMask: uint(numBuckets - 1)}, nil
}

func (cf *Filter) Insert(item []byte) bool {
	index, fp := GetIndexAndFingerprint(item, cf.bucketMask)
	if cf.buckets[index].insert(fp) {
		//log.Println(cf.buckets)
		return true
	}
	secondIndex := GetAltIndex(fp, index, cf.bucketMask)
	if cf.buckets[secondIndex].insert(fp) {
		//log.Println(cf.buckets)
		return true
	}
	// now start evicting elements from buckets
	switch rand.Intn(2) {
	case 0:
		return cf.reinsert(fp, index)
	case 1:
		return cf.reinsert(fp, secondIndex)
	default:
		log.Println("invocation to reinsert: should not have been called")
		return cf.reinsert(fp, index) // we are NOT supposed to hit this line of code
	}
}

func (b *Bucket) insert(fp Fingerprint) bool {
	for i, el := range b.fp {
		if el == 0 {
			b.fp[i] = fp
			return true
		}
	}
	return false
}

func (cf *Filter) reinsert(fp Fingerprint, alternateIndex uint) bool {
	for i := 0; i < MAX_EVICTIONS; i++ { // maximum num of evictions before filter is 'full'
		randomElementPosition := rand.Intn(BUCKET_SIZE)
		toBeInserted := fp
		randomElementEvicted := cf.buckets[alternateIndex].fp[randomElementPosition]
		cf.buckets[alternateIndex].fp[randomElementPosition] = toBeInserted
		randomElementAltIndex := GetAltIndex(randomElementEvicted, alternateIndex, cf.bucketMask)
		if cf.buckets[randomElementAltIndex].insert(randomElementEvicted) {
			return true
		}
		// keep trying to evict elements
	}
	return false // mission failed, shuffling elements around didn't work
}

func (cf *Filter) Lookup(item []byte) bool {
	index, fp := GetIndexAndFingerprint(item, cf.bucketMask)
	secondIndex := GetAltIndex(fp, index, cf.bucketMask)
	return cf.buckets[index].lookup(fp) || cf.buckets[secondIndex].lookup(fp)
}

func (b *Bucket) lookup(fp Fingerprint) bool {
	for _, el := range b.fp {
		if fp == el {
			return true
		}
	}
	return false
}

func (cf *Filter) Delete(item []byte) bool {
	index, fp := GetIndexAndFingerprint(item, cf.bucketMask)
	if cf.buckets[index].delete(fp) {
		return true
	}
	secondIndex := GetAltIndex(fp, index, cf.bucketMask)
	if cf.buckets[secondIndex].delete(fp) {
		return true
	}
	return false
}

func (b *Bucket) delete(fp Fingerprint) bool {
	for i, el := range b.fp {
		if el == fp {
			b.fp[i] = 0
			return true
		}
	}
	return false
}

// encodes bytes in big endian order, and truncates the first 5 bits
// (presumably all 0's since we set the maximum value of fingerprints
// to be 2^27 - 1
func (cf *Filter) Encode() []byte {
	ArrayLen := int(math.Ceil(float64(len(cf.buckets) * BUCKET_SIZE * FINGERPRINT_BITS) / float64(8)))
	output := make([]byte, ArrayLen)
	byteOffset := 0
	bitOffset := 0
	remainderBits := int(FINGERPRINT_BITS - 8 * math.Floor(FINGERPRINT_BITS / 8))
	for _, bucket := range cf.buckets {
		for _, fingerprint := range bucket.fp {
			arrayElements := FingerprintToByteArray(fingerprint)
			arrayElements = ShiftByteArray(arrayElements)
			for i, elToInsert := range arrayElements {
				if i != BUCKET_SIZE - 1 {
					output = WriteByteToPositionAndBitOffset(output, elToInsert, byteOffset, bitOffset)
					byteOffset++
				} else {
					output = WriteBitsToPositionAndBitOffset(output, elToInsert, byteOffset, bitOffset, remainderBits)
					bitOffset += remainderBits
					if bitOffset >= 8 {
						byteOffset++
						bitOffset -= 8
					}
				}
			}
		}
	}
	return output
}

// associated decoding function
func Decode(filter []byte) (*Filter, error) {
	numBuckets := (8 * len(filter)) / (BUCKET_SIZE * FINGERPRINT_BITS)
	if math.Floor(float64(numBuckets)) != math.Ceil(float64(numBuckets)) {
		return &Filter{buckets: nil}, errors.New("filter decode alignment issues")
	}
	if !IsPowerOfTwo(numBuckets) {
		return &Filter{}, errors.New("num of buckets not a power of 2")
	}
	cf := &Filter{
		buckets: make([]Bucket, numBuckets),
		bucketMask: uint(numBuckets - 1),
	}
	byteOffset := 0
	bitOffset := 0
	remainderBits := int(FINGERPRINT_BITS - 8 * math.Floor(FINGERPRINT_BITS / 8))
	numBytesToRead := int(math.Ceil(float64(FINGERPRINT_BITS) / float64(8)))
	for bucketIndex, _ := range cf.buckets {
		cf.buckets[bucketIndex].fp = make([]Fingerprint, BUCKET_SIZE)
		for positionWithinBucket, _ := range cf.buckets[bucketIndex].fp {
			fp := ReadFingerprintFromPositionAndBitOffset(filter, byteOffset, bitOffset, FINGERPRINT_BITS)
			cf.buckets[bucketIndex].fp[positionWithinBucket] = fp
			byteOffset += numBytesToRead
			bitOffset += remainderBits
			if bitOffset >= 8 {
				bitOffset -= 8
			} else {
				byteOffset -= 1
			}
		}
	}
	return cf, nil
}