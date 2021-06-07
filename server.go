package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pancast-server/routes"
	serverutils "pancast-server/server-utils"
)

type Env struct {
	db *sql.DB
}

func basic(w http.ResponseWriter, req *http.Request) {
	serverutils.Write(w, "Welcome")
}

func StartServer(address string) {
	db := InitDatabaseConnection()
	defer db.Close()
	env := &Env{db: db}
	http.HandleFunc("/", basic)
	http.HandleFunc("/register", env.registerNewDeviceIndex)
	http.HandleFunc("/upload", env.uploadRiskEncountersIndex)
	http.HandleFunc("/update", env.updateRiskAssessmentIndex)
	server := &http.Server{Addr: address}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func() {
		fmt.Println("Listening on address: " + address)
		if err := http.ListenAndServeTLS(address, "pancast.cert", "pancast.key", nil); err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	fmt.Println("\nShutting down...")
	err := server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}

func (env *Env) registerNewDeviceIndex(w http.ResponseWriter, req *http.Request) {
	params, err := routes.RegisterController(0, env.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		output, _ := params.ConvertToJSONString()
		serverutils.Write(w, output)
	}
}

func (env *Env) uploadRiskEncountersIndex(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

func (env *Env) updateRiskAssessmentIndex(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
}
