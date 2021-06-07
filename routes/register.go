package routes

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
)

type RegistrationParameters struct {
	ServerKey  PublicKey
	DeviceID   uint64
	Clock      uint64
	Secret     string
	OTPs       []string
	LocationID string
}

type PublicKey struct {
	N *big.Int
	E int
}

func RegisterController(deviceType int, db *sql.DB) (RegistrationParameters, error) {
	// 1. The backend's public key
	// 2. The device ID
	// 3. Initial clock value
	// 4. Secret key

	// if it is a beacon, we will give it:
	// 1.

	// obtain public key of server
	// TODO: Compute a secret key to give to a beacon

	// insert device into database

	// handle registration logic
	// tons of crypto
	// maybe some database access
	return RegistrationParameters{}, nil
}

// adapted from https://stackoverflow.com/questions/33031658/getting-rsa-public-key-from-certificate-in-golang
func getPublicKey() (PublicKey, error) {
	certBytes, err := ioutil.ReadFile("pancast.cert")
	if err != nil {
		log.Fatal(err)
		return PublicKey{}, err
	}
	certBlock, _ := pem.Decode(certBytes)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		log.Fatal(err)
		return PublicKey{}, err
	}
	certKey := cert.PublicKey.(*rsa.PublicKey)
	return PublicKey{certKey.N, certKey.E}, nil

}

func (r *RegistrationParameters) ConvertToJSONString() (string, error) {
	output, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return string(output), nil
}
