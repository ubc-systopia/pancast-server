package server

import (
	"fmt"
	"log"
	"net/http"
	serverutils "pancast-server/server-utils"
)

func basic(res http.ResponseWriter, req *http.Request) {
	serverutils.Write(res, "Welcome")
}

func StartServer(address string) {
	http.HandleFunc("/", basic)
	//	http.HandleFunc("/register", nil)
	//	http.HandleFunc("/upload", nil)
	fmt.Println("Listening on address: " + address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
