package routes

/*
Controller for the /update route.
*/

import (
	"encoding/binary"
	"pancast-server/cuckoo"
)

type UpdateReturnParameters struct {
	Length int64
	Filter []byte
}

func UpdateController(cf *cuckoo.Filter) []byte {
	// divide the number of ephemeral IDs into a number of cuckoo filters that we keep as chunked storage
	if cf == nil {
		return []byte{}
	}
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(cf.Buckets))*cuckoo.FINGERPRINT_BITS*cuckoo.BUCKET_SIZE/8) // add ceil
	data := cf.Encode()
	payload := append(length, data...)
	return payload
}

func UpdateControllerGetCount(chunks []*cuckoo.Filter) byte {
	return byte(len(chunks))
}

func UpdateControllerGetChunk(chunks []*cuckoo.Filter, num int) []byte {
	if len(chunks) <= num || num < 0 || chunks[num] == nil {
		return []byte{}
	}
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(chunks[num].Buckets))*cuckoo.FINGERPRINT_BITS*cuckoo.BUCKET_SIZE/8) // add ceil
	data := chunks[num].Encode()
	payload := append(length, data...)
	return payload
}
