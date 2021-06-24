package types

import "fmt"

type DongleRegistrationData struct {
	DeviceData RegistrationData
	OTPs       []string
}

func ConcatDongleData(d DongleRegistrationData) string {
	// for the time being, ignore OTPs
	return fmt.Sprintf("(%d)",
		d.DeviceData.DeviceID)
}
