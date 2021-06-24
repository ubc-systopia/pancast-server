package routes

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"pancast-server/models"
	server_utils "pancast-server/server-utils"
	"pancast-server/types"
	"time"
)

type RegistrationParameters struct {
	ServerKey  string
	DeviceID   uint32
	Clock      uint32
	Secret     []byte
	OTPs       []string
	LocationID string
}

func RegisterController(deviceType int64, keyLoc string, db *sql.DB) (RegistrationParameters, error) {
	var output RegistrationParameters
	// temporary placeholder for location
	tempBeaconLocation := "LOCATION"
	output.LocationID = tempBeaconLocation

	// temporary placeholder for OTPs
	output.OTPs = []string{}

	// get public key
	key, err := ioutil.ReadFile(keyLoc)
	if err != nil {
		return RegistrationParameters{}, err
	}
	output.ServerKey = string(key)

	// compute current time
	output.Clock = GetCurrentMinuteStamp()

	// using the AES-256 standard, where keys have 32 bytes
	aesKey, err := GenerateRandomByteString(32)
	if err != nil {
		return RegistrationParameters{}, err
	}
	output.Secret = aesKey

	// compute available device ID
	id, err := models.GetLowestAvailableDeviceID(db)
	if err != nil {
		return RegistrationParameters{}, err
	}
	output.DeviceID = id

	err = deviceDatabaseHandler(deviceType, output, db)
	if err != nil {
		return RegistrationParameters{}, err
	}
	return output, nil
}

func deviceDatabaseHandler(dType int64, params RegistrationParameters, db *sql.DB) error {
	deviceData := types.RegistrationData{
		DeviceID:  params.DeviceID,
		Secret:    params.Secret,
		Clock:     params.Clock,
	}
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = models.CreateDevice(deviceData, tx, ctx)
	if err != nil {
		return err
	}
	switch dType {
	case server_utils.DONGLE:
		dongleData := types.DongleRegistrationData{
			DeviceData: deviceData,
			OTPs:       params.OTPs,
		}
		err = models.CreateDongle(dongleData, tx, ctx)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	case server_utils.BEACON:
		beaconData := types.BeaconRegistrationData{
			DeviceData: deviceData,
			LocationID: params.LocationID,
		}
		err = models.CreateBeacon(beaconData, tx, ctx)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	default:
		return errors.New("bad device type")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentMinuteStamp() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Minute))
}

func GenerateRandomByteString(n int) ([]byte, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return key, err
	}
	return key, nil
}

func (r *RegistrationParameters) ConvertToJSONString() (string, error) {
	output, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(output), nil
}
