package cuckoo

/*
	An implementation of a Cuckoo filter, based off of https://github.com/seiflotfy/cuckoofilter
	Main addition is the efficient CF export functionality
	Intended recipient for this exported filter are the dongle devices
*/

import (
	"errors"
	"log"
	"math"
	"math/rand"
)

type Filter struct {
	bucketMask uint
	Buckets    []Bucket
}

type Bucket struct {
	Fp []Fingerprint
}

type Fingerprint uint32 // (FINGERPRINT_SIZE * 8 bits)

func CreateFilter(numBuckets int) (*Filter, error) {
	if numBuckets == 0 {
		return nil, errors.New("cannot have zero Buckets")
	}
	if !IsPowerOfTwo(numBuckets) {
		// if not a power of 2, then it introduces a bias in the distribution of hashes
		return nil, errors.New("bucket number not a power of 2")
	}
	var tempFilter []Bucket
	for i := 0; i < numBuckets; i++ {
		var bkt Bucket
		for j := 0; j < BUCKET_SIZE; j++ {
			var fp Fingerprint
			fp = Fingerprint(0)
			bkt.Fp = append(bkt.Fp, fp)
		}
		tempFilter = append(tempFilter, bkt)
	}
	return &Filter{Buckets: tempFilter, bucketMask: uint(numBuckets - 1)}, nil
}

func (cf *Filter) Insert(item []byte) bool {
	index, fp := GetIndexAndFingerprint(item, cf.bucketMask)
	//log.Println("First index")
	//log.Println(index)
	//log.Println("Second index")
	//log.Println(GetAltIndex(fp, index, cf.bucketMask))
	if cf.Buckets[index].insert(fp) {
		return true
	}
	secondIndex := GetAltIndex(fp, index, cf.bucketMask)
	if cf.Buckets[secondIndex].insert(fp) {
		return true
	}
	// now start evicting elements from Buckets
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
	for i, el := range b.Fp {
		if el == 0 {
			b.Fp[i] = fp
			return true
		}
	}
	return false
}

func (cf *Filter) reinsert(fp Fingerprint, alternateIndex uint) bool {
	toInsert := fp
	alternateIndexOfCurrentItem := alternateIndex
	for i := 0; i < MAX_EVICTIONS; i++ { // maximum num of evictions before filter is 'full'
		randomElementPosition := rand.Intn(BUCKET_SIZE)
		randomElementEvicted := cf.Buckets[alternateIndexOfCurrentItem].Fp[randomElementPosition]
		cf.Buckets[alternateIndexOfCurrentItem].Fp[randomElementPosition] = toInsert
		randomElementAltIndex := GetAltIndex(randomElementEvicted, alternateIndexOfCurrentItem, cf.bucketMask)
		if cf.Buckets[randomElementAltIndex].insert(randomElementEvicted) {
			return true
		}
		toInsert = randomElementEvicted
		alternateIndexOfCurrentItem = randomElementAltIndex
		// keep trying to evict elements
	}
	return false // mission failed, shuffling elements around didn't work
}

func (cf *Filter) Lookup(item []byte) bool {
	index, fp := GetIndexAndFingerprint(item, cf.bucketMask)
	secondIndex := GetAltIndex(fp, index, cf.bucketMask)
	return cf.Buckets[index].lookup(fp) || cf.Buckets[secondIndex].lookup(fp)
}

func (b *Bucket) lookup(fp Fingerprint) bool {
	for _, el := range b.Fp {
		if fp == el {
			return true
		}
	}
	return false
}

func (cf *Filter) Delete(item []byte) bool {
	index, fp := GetIndexAndFingerprint(item, cf.bucketMask)
	if cf.Buckets[index].delete(fp) {
		return true
	}
	secondIndex := GetAltIndex(fp, index, cf.bucketMask)
	if cf.Buckets[secondIndex].delete(fp) {
		return true
	}
	return false
}

func (b *Bucket) delete(fp Fingerprint) bool {
	for i, el := range b.Fp {
		if el == fp {
			b.Fp[i] = 0
			return true
		}
	}
	return false
}

// encodes bytes in big endian order, and truncates the first 5 bits
// (presumably all 0's since we set the maximum value of fingerprints
// to be 2^27 - 1
func (cf *Filter) Encode() []byte {
	ArrayLen := int(math.Ceil(float64(len(cf.Buckets)*BUCKET_SIZE*FINGERPRINT_BITS) / float64(8)))
	output := make([]byte, ArrayLen)
	byteOffset := 0
	bitOffset := 0
	remainderBits := int(FINGERPRINT_BITS - 8*math.Floor(FINGERPRINT_BITS/8)) // 8 to constant
	for _, bucket := range cf.Buckets {
		for _, fingerprint := range bucket.Fp {
			arrayElements := FingerprintToByteArrayBigEndian(fingerprint)
			arrayElements = ShiftByteArray(arrayElements)
			for i, elToInsert := range arrayElements {
				if i != BUCKET_SIZE-1 {
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
		return &Filter{Buckets: nil}, errors.New("filter decode alignment issues")
	}
	if !IsPowerOfTwo(numBuckets) {
		return &Filter{}, errors.New("num of Buckets not a power of 2")
	}
	cf := &Filter{
		Buckets:    make([]Bucket, numBuckets),
		bucketMask: uint(numBuckets - 1),
	}
	byteOffset := 0
	bitOffset := 0
	remainderBits := int(FINGERPRINT_BITS - 8*math.Floor(FINGERPRINT_BITS/8))
	numBytesToRead := int(math.Ceil(float64(FINGERPRINT_BITS) / float64(8)))
	for bucketIndex, _ := range cf.Buckets {
		cf.Buckets[bucketIndex].Fp = make([]Fingerprint, BUCKET_SIZE)
		for positionWithinBucket, _ := range cf.Buckets[bucketIndex].Fp {
			fp := ReadFingerprintFromPositionAndBitOffset(filter, byteOffset, bitOffset, FINGERPRINT_BITS)
			cf.Buckets[bucketIndex].Fp[positionWithinBucket] = fp
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
