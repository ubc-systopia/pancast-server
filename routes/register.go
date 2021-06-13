package routes

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
)

type RegistrationParameters struct {
	ServerKey  string
	DeviceID   uint64
	Clock      uint64
	Secret     string
	OTPs       []string
	LocationID string
}

func RegisterController(deviceType int64, keyLoc string, db *sql.DB) (RegistrationParameters, error) {
	var output RegistrationParameters
	key, err := ioutil.ReadFile(keyLoc)
	if err != nil {
		return RegistrationParameters{}, err
	}
	output.ServerKey = string(key)

	// TODO: Compute current clock offset

	// TODO: Create secret AES key

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
