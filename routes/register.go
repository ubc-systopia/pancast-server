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
	"os"
	"strconv"
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

func RegisterController(deviceType int64, db *sql.DB) (RegistrationParameters, error) {
	// Get server's public key from certificate?
	// Client must authenticate with server that it is legit
	// Not complete yet
	var output RegistrationParameters
	//key, err := getPublicKey()
	//if err != nil {
	//	return RegistrationParameters{}, err
	//}
	n := new(big.Int)
	n, _ = n.SetString(os.Getenv("PUBLIC_N"), 10)
	e, _ := strconv.Atoi(os.Getenv("PUBLIC_E"))
	key := PublicKey{
		N: n,
		E: e,
	}
	output.ServerKey = key

	// TODO: Compute current clock offset

	// TODO: Create secret AES key

	// TODO: Create device in database
	return output, nil
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
