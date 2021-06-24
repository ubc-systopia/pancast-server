package server

import (
	"database/sql"
	"fmt"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pancast-server/config"
	"pancast-server/cronjobs"
	"pancast-server/cuckoo"
	"pancast-server/database"
	"pancast-server/routes"
	serverutils "pancast-server/server-utils"
	"strconv"
)

type Env struct {
	db             *sql.DB
	cf             *cuckoo.Filter
	certificateLoc string
	privateKeyLoc  string
	publicKeyLoc   string
	privacyParams  cronjobs.DiffprivParameters
}



func basic(w http.ResponseWriter, req *http.Request) {
	serverutils.Write(w, "Welcome")
}

func StartServer(conf config.StartParameters) (*http.Server, *Env, chan os.Signal) {
	db := database.InitDatabaseConnection()
	serverURL := config.GetServerURL(conf)
	mean, _ := strconv.Atoi(os.Getenv("MEAN"))
	sens, _ := strconv.ParseFloat(os.Getenv("SENS"), 64)
	epsilon, _ := strconv.ParseFloat(os.Getenv("EPSILON"), 64)
	delta, _ := strconv.ParseFloat(os.Getenv("DELTA"), 64)

	env := &Env{
		db:             db,
		cf:             nil,
		certificateLoc: conf.CertificateLoc,
		privateKeyLoc:  conf.PrivateKeyLoc,
		publicKeyLoc:   conf.PublicKeyLoc,
		privacyParams: cronjobs.DiffprivParameters{
			Mean:        int64(mean),
			Sensitivity: sens,
			Epsilon:     epsilon,
			Delta:       delta,
		},

	}
	cronjobs.CreateNewFilter(env.cf, env.db, env.privacyParams) // create filter on startup for now
	http.HandleFunc("/", basic)
	http.HandleFunc("/register", env.RegisterNewDeviceIndex)
	http.HandleFunc("/upload", env.UploadRiskEncountersIndex)
	http.HandleFunc("/update", env.UpdateRiskAssessmentIndex)
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		cronjobs.CreateNewFilter(env.cf, env.db, env.privacyParams)
	})
	if err != nil {
		log.Println("error creating cron job")
	}
	server := &http.Server{Addr: serverURL}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func() {
		fmt.Println("Listening on address: " + serverURL)
		if err := http.ListenAndServeTLS(serverURL, env.certificateLoc, env.privateKeyLoc, nil); err != nil {
			log.Fatal(err)
		}
	}()
	return server, env, done
}

func (env *Env) RegisterNewDeviceIndex(w http.ResponseWriter, req *http.Request) {
	deviceType, err := strconv.ParseInt(req.FormValue("type"), 10, 32)
	if err != nil || !isValidType(deviceType) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	params, err := routes.RegisterController(deviceType, env.publicKeyLoc, env.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
	} else {
		output, _ := params.ConvertToJSONString()
		serverutils.Write(w, output)
	}
}

func isValidType(deviceType int64) bool {
	return deviceType == serverutils.DONGLE || deviceType == serverutils.BEACON
}

func (env *Env) UploadRiskEncountersIndex(w http.ResponseWriter, req *http.Request) {
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

func (env *Env) UpdateRiskAssessmentIndex(w http.ResponseWriter, req *http.Request) {
	// TODO: implement
	ba := routes.UpdateController(env.cf)
	code, err := w.Write(ba)
	if err != nil {
		log.Println(err)
	}
	log.Println(code)
}
