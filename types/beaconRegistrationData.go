package types

import "fmt"

type BeaconRegistrationData struct {
	DeviceData RegistrationData
	LocationID uint64
}

func ConcatBeaconData(b BeaconRegistrationData) string {
	return fmt.Sprintf("(%d, '%s')",
		b.DeviceData.DeviceID, b.LocationID)
}

func (b BeaconRegistrationData) CreateDevice() {

}
