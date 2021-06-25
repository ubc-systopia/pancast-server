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
	// nothing complicated here, just dumping the filter as a byte array

	// new changeee
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(cf.Buckets)))
	data := cf.Encode()
	payload := append(length, data...)
	return payload

	//param := UpdateReturnParameters{
	//	Length: int64(len(cf.Buckets)),
	//	Filter: cf.Encode(),
	//}
	//jsonData, err := json.Marshal(param)
	//if err != nil {
	//	log.Println(err)
	//	return []byte{}
	//}
	//return jsonData
	//return cf.Encode()


}
