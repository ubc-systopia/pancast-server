package routes

import (
	"encoding/json"
	"log"
	"pancast-server/cuckoo"
)

type UpdateReturnParameters struct {
	Length int
	Filter []byte
}

func UpdateController(cf *cuckoo.Filter) []byte {
	// nothing complicated here :)
	param := UpdateReturnParameters{
		Length: len(cf.Buckets),
		Filter: cf.Encode(),
	}
	jsonData, err := json.Marshal(param)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return jsonData
}
