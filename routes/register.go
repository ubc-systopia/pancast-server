package routes

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"pancast-server/models"
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
	log.Println(err)
	if err != nil {
		return RegistrationParameters{}, err
	}
	output.DeviceID = id

	// TODO: Create device in database
	return output, nil
}

// adapted from https://stackoverflow.com/questions/33031658/getting-rsa-public-key-from-certificate-in-golang
//func getPublicKey() (PublicKey, error) { // no need to do this, just generate a public key using cmd tools
//	certBytes, err := ioutil.ReadFile("pancast.cert")
//	if err != nil {
//		log.Fatal(err)
//		return PublicKey{}, err
//	}
//	certBlock, _ := pem.Decode(certBytes)
//	cert, err := x509.ParseCertificate(certBlock.Bytes)
//	if err != nil {
//		log.Fatal(err)
//		return PublicKey{}, err
//	}
//	certKey := cert.PublicKey.(*rsa.PublicKey)
//	return PublicKey{certKey.N, certKey.E}, nil
//
//}


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

func getPublicKey() (string, error) {
	pubkey, err := ioutil.ReadFile("pancast.pubkey")
	if err != nil {
		return "", err
	}
	return string(pubkey), nil
}

func (r *RegistrationParameters) ConvertToJSONString() (string, error) {
	output, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(output), nil
}
