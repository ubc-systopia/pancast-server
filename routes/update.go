package routes

import (
	"encoding/binary"
	"pancast-server/cuckoo"
)

type UpdateReturnParameters struct {
	Length int64
	Filter []byte
}

func UpdateController(cf *cuckoo.Filter) []byte {
	if cf == nil {
		return []byte{}
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(cf.Buckets))*cuckoo.FINGERPRINT_BITS*cuckoo.BUCKET_SIZE/8) // add ceil
	data := cf.Encode()
	payload := append(length, data...)
	return payload

}
