package server

import (
	"bytes"
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
	"pancast-server/models"
	"pancast-server/routes"
	serverutils "pancast-server/server-utils"
	"strconv"
	"time"
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
	// initialization
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

	// initialize filter on startup
	newFilter, err := cronjobs.CreateNewFilter(env.db, env.privacyParams) // create filter on startup for now
	if err != nil {
		log.Fatal(err)
	}
	env.cf = newFilter

	// initialize routes
	mux := http.NewServeMux()
	basicHandler := http.HandlerFunc(basic)
	registerHandler := http.HandlerFunc(env.RegisterNewDeviceIndex)
	uploadHandler := http.HandlerFunc(env.UploadRiskEncountersIndex)
	updateHandler := http.HandlerFunc(env.UpdateRiskAssessmentIndex)
	mux.Handle("/", env.TelemetryWrapper(basicHandler))
	mux.Handle("/register", env.TelemetryWrapper(registerHandler))
	mux.Handle("/upload", env.TelemetryWrapper(uploadHandler))
	mux.Handle("/update", env.TelemetryWrapper(updateHandler))

	// initialize cron job
	c := cron.New()
	_, err = c.AddFunc("@midnight", func() { // starts from the moment this is invoked
		newFilter, err = cronjobs.CreateNewFilter(env.db, env.privacyParams)
		if err != nil {
			log.Println("error updating cuckoo filter")
		}
		env.cf = newFilter
	})
	c.Start()
	if err != nil {
		log.Println("error creating cron job")
	}
	server := &http.Server{Addr: serverURL}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func() {
		fmt.Println("Listening on address: " + serverURL)
		if err := http.ListenAndServeTLS(serverURL, env.certificateLoc, env.privateKeyLoc, mux); err != nil {
			log.Fatal(err)
		}
	}()
	return server, env, done
}

func (env *Env) RegisterNewDeviceIndex(w http.ResponseWriter, req *http.Request) {
	deviceType, err := strconv.ParseInt(req.FormValue("type"), 10, 32)
	if err != nil || !isValidType(deviceType) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	params, err := routes.RegisterController(deviceType, env.publicKeyLoc, env.db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
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
		return
	}
	input, err := routes.ConvertStringToUploadParam(body)
	if err != nil || !isValidDatabase(input.Type) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if input.Type == 0 && !hasPermissionToUploadToRiskDatabase() {
		w.WriteHeader(http.StatusForbidden)
		return
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
	ba := routes.UpdateController(env.cf)
	_, err := w.Write(ba)
	if err != nil {
		log.Println(err)
	}
}

func (env *Env) TelemetryWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, _ := ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		recvTime := time.Now()
		h.ServeHTTP(w, req)
		totalTime := time.Since(recvTime)
		log.Println("Request received")
		log.Println("Time elapsed: " + totalTime.String())
		log.Println("Routed for " + req.URL.Path)
		numEntries := -1
		if req.URL.Path == "/upload" {
			input, err := routes.ConvertStringToUploadParam(body)
			if err == nil {
				numEntries = len(input.Entries)
				log.Println("Number of ephemeral IDs submitted: " + strconv.Itoa(numEntries))
			} else {
				log.Println("Bad request")
			}
		}
		models.CreateTelemetryEntry(totalTime.String(), req.URL.Path, numEntries, serverutils.GetCurrentMinuteStamp(), env.db)
	})
}
