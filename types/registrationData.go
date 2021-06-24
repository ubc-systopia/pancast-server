package types

import "fmt"

type IRegistrationData interface {
	CreateSpecificDevice() bool
}

type RegistrationData struct {
	DeviceID uint32
	Secret   []byte
	Clock    uint32
}

func ConcatRegistrationData(r RegistrationData) string {
	return fmt.Sprintf("(%d, '%s', %d, %d)",
		r.DeviceID, r.Secret, r.Clock, 0)
}
