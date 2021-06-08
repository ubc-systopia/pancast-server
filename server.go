package server

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pancast-server/database"
	"pancast-server/routes"
	serverutils "pancast-server/server-utils"
	"strconv"
)

type Env struct {
	db *sql.DB
}

func basic(w http.ResponseWriter, req *http.Request) {
	serverutils.Write(w, "Welcome")
}

func StartServer(address string) (*http.Server, chan os.Signal) {
	db := database.InitDatabaseConnection()
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
	return server, done
}

func (env *Env) registerNewDeviceIndex(w http.ResponseWriter, req *http.Request) {
	deviceType, err := strconv.ParseInt(req.FormValue("type"), 10, 32)
	if err != nil || !isValidType(deviceType) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	params, err := routes.RegisterController(deviceType, env.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		output, _ := params.ConvertToJSONString()
		serverutils.Write(w, output)
	}
}

func isValidType(deviceType int64) bool {
	return deviceType == serverutils.DONGLE || deviceType == serverutils.BEACON
}

func (env *Env) uploadRiskEncountersIndex(w http.ResponseWriter, req *http.Request) {
	body, errBody := ioutil.ReadAll(req.Body)
	if errBody != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	input := routes.ConvertStringToUploadParam(body)
	if !isValidDatabase(input.Type) {
		w.WriteHeader(http.StatusBadRequest)
	}
	if input.Type == 0 && !hasPermissionToUploadToRiskDatabase() {
		w.WriteHeader(http.StatusForbidden)
	} else {
		err := routes.UploadController(input, env.db)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
		} else {
			serverutils.Write(w, "Success!")
		}
	}
}

func isValidDatabase(databaseType int64) bool {
	return databaseType == serverutils.RISK || databaseType == serverutils.EPI
}

func hasPermissionToUploadToRiskDatabase() bool {
	// TODO: implement
	return true
}

func (env *Env) updateRiskAssessmentIndex(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
}
