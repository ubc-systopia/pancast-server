package cuckoo

import (
	"encoding/binary"
	"github.com/dgryski/go-metro"
	"math"
)

func IsPowerOfTwo(numBuckets int) bool {
	logNum := math.Log2(float64(numBuckets))
	if math.Ceil(logNum) != math.Floor(logNum) {
		return false
	}
	return true
}

func GetAltIndex(fp Fingerprint, index uint, bucketMask uint) uint {
	hash := GetHash(FingerprintToByteArray(fp))
	return uint(uint64(index)^hash) & bucketMask
}

func GetFingerprint(hash uint64) uint32 {
	// Use least significant bits for fingerprint.
	// Range of FP is [1, MAX_VALUE_THAT_27_BITS_CAN_REPRESENT]
	FingerprintMask := uint64(math.Pow(2, FINGERPRINT_BITS)) - 1
	fp := uint32(hash%FingerprintMask + 1)
	return fp
}

func GetHash(item []byte) uint64 {
	return metro.Hash64(item, 1337)
}

//
//// getIndicesAndFingerprint returns the 2 bucket indices and fingerprint to be used
func GetIndexAndFingerprint(item []byte, bucketMask uint) (uint, Fingerprint) {
	hash := GetHash(item)
	fp := GetFingerprint(hash)
	// Use least significant bits for deriving index.
	i1 := uint(hash & uint64(bucketMask))
	return i1, Fingerprint(fp)
}

func FingerprintToByteArray(item Fingerprint) []byte {
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(item))
	return output
}

func ShiftByteArray(item []byte) []byte {
	output := make([]byte, FINGERPRINT_SIZE)
	for i, _ := range output {
		if i != 0 {
			output[i-1] |= item[i] >> (8 - SHIFT_LEEWAY)
		}
		output[i] |= item[i] << SHIFT_LEEWAY
	}
	return output
}

func ByteArrayToFingerprint(item []byte) Fingerprint {
	return Fingerprint(binary.BigEndian.Uint32(item))
}

func UnshiftByteArray(item []byte) []byte {
	output := make([]byte, FINGERPRINT_SIZE)
	for i, _ := range output {
		if i == 0 {
			output[i] |= item[i] >> SHIFT_LEEWAY
		} else {
			output[i] |= item[i-1]<<(8-SHIFT_LEEWAY) | item[i]>>SHIFT_LEEWAY
		}
	}
	return output
}

func WriteByteToPositionAndBitOffset(arr []byte, input byte, byteOffset int, bitOffset int) []byte {
	if bitOffset == 0 {
		arr[byteOffset] = input
	} else {
		prevByteLowerBits := input >> bitOffset
		nextByteUpperBits := input << (8 - bitOffset)
		arr[byteOffset] |= prevByteLowerBits
		arr[byteOffset+1] |= nextByteUpperBits
	}
	return arr
}

func WriteBitsToPositionAndBitOffset(arr []byte, input byte, byteOffset int, bitOffset int, numBits int) []byte {
	if numBits+bitOffset <= 8 {
		arr[byteOffset] |= input >> bitOffset
	} else {
		prevByteLowerBits := input >> bitOffset
		nextByteUpperBits := input << (8 - bitOffset)
		arr[byteOffset] |= prevByteLowerBits
		arr[byteOffset+1] |= nextByteUpperBits
	}
	return arr
}

func ReadFingerprintFromPositionAndBitOffset(arr []byte, byteOffset int, bitOffset int, numBits int) Fingerprint {
	numBytesToRead := int(math.Ceil(float64(numBits) / float64(8)))
	remainingBitsToRead := bitOffset + (numBits % 8)
	if remainingBitsToRead >= 8 {
		remainingBitsToRead -= 8
		numBytesToRead++
	}
	output := make([]byte, FINGERPRINT_SIZE)
	for i := 0; i < numBytesToRead; i++ {
		// needs refactoring
		if numBytesToRead > FINGERPRINT_SIZE && i == numBytesToRead-1 {
			// at the last array slot, only add to i - 1th block
			if remainingBitsToRead != 0 {
				output[i-1] |= (arr[byteOffset+i] >> (8 - bitOffset)) &
					GenerateHighValueBitmask(remainingBitsToRead+8-bitOffset)
			}
		} else if i == 0 {
			// at first array slot, add to first block
			output[i] |= arr[byteOffset+i] << bitOffset
		} else {
			// add to i - 1th and ith slot
			output[i-1] |= arr[byteOffset+i] >> (8 - bitOffset)
			if i == numBytesToRead-1 {
				// read remaining bits
				output[i] |= (arr[byteOffset+i] << bitOffset) & GenerateHighValueBitmask(remainingBitsToRead)
			} else {
				// read till offset
				output[i] |= arr[byteOffset+i] << bitOffset
			}

		}
	}

	unshift := UnshiftByteArray(output) // for big endian representation over the wire
	return ByteArrayToFingerprint(unshift)
}

func GenerateLowValueBitmask(numBits int) byte {
	output := 0
	for i := 0; i < numBits; i++ {
		output += int(math.Pow(2, float64(i)))
	}
	return byte(output)
}

func GenerateHighValueBitmask(numBits int) byte {
	output := 0
	for i := 0; i < numBits; i++ {
		output += int(float64(128) * math.Pow(0.5, float64(i)))
	}
	return byte(output)
}
