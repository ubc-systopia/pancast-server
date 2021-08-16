package types

import (
	"encoding/base64"
	"fmt"
)

type Entry struct {
	EphemeralID string // (should be byte array, but for now its a string)
	LocationID  uint64
	DongleClock uint64
	BeaconClock uint64
	BeaconID    uint64
}

func ConcatEntries(input []Entry) string {
	output := ""
	for idx, entry := range input {
		output += createEntries(entry)
		if idx != len(input)-1 { // not last element
			output += ","
		}
	}
	return output
}

func createEntries(input Entry) string {
	decodedEphID, err := base64.StdEncoding.DecodeString(input.EphemeralID)
	if err != nil {
		decodedEphID = []byte(input.EphemeralID)
	}
	return fmt.Sprintf("('%s', '%d', %d, %d, %d)", decodedEphID, input.LocationID,
		input.DongleClock, input.BeaconClock, input.BeaconID)
}
