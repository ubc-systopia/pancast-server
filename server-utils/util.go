package server_utils

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

func Write(res http.ResponseWriter, output string) {
	_, err := fmt.Fprint(res, output)
	if err != nil {
		log.Fatal(err)
	}
}

func ConcatPublicKey(N *big.Int, E int) string {
	return N.String() + "|" + strconv.Itoa(E)
}
