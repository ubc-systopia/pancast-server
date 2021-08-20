package server

/*
	Server code. Starts the server, sets handlers for HTTPS routes, listens on a part and sleeps.
*/

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"math/rand"
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
	cfChunks       []*cuckoo.Filter
	mode           []string
	certificateLoc string
	privateKeyLoc  string
	publicKeyLoc   string
	privacyParams  cronjobs.DiffprivParameters
}

func basic(w http.ResponseWriter, req *http.Request) {
	serverutils.Write(w, "Goodbye")
	w.WriteHeader(http.StatusBadRequest)
	return
}

func StartServer(conf config.StartParameters) (*http.Server, *Env, chan os.Signal) {
	// initialization
	rand.Seed(time.Now().UnixNano()) // ssssssecretsssss
	db := database.InitDatabaseConnection()
	serverURL := config.GetServerURL(conf)
	mean, _ := strconv.Atoi(os.Getenv("MEAN"))
	sens, _ := strconv.ParseFloat(os.Getenv("SENS"), 64)
	epsilon, _ := strconv.ParseFloat(os.Getenv("EPSILON"), 64)
	delta, _ := strconv.ParseFloat(os.Getenv("DELTA"), 64)
	mode := config.GetApplicationMode(conf)

	env := &Env{
		db:             db,
		cf:             nil,
		cfChunks:       nil,
		mode:           mode,
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
	ephIDs, length := cronjobs.GenerateEphemeralIDList(env.db, env.privacyParams, env.mode)
	newFilter, err := cronjobs.CreateNewFilter(ephIDs, length) // create filter on startup for now
	if err != nil {
		log.Fatal(err)
	}
	env.cf = newFilter
	chunks, err := cronjobs.CreateChunkedFilters(ephIDs, length)
	if err != nil {
		log.Fatal(err)
	}
	env.cfChunks = chunks

	// initialize routes
	mux := http.NewServeMux()
	basicHandler := http.HandlerFunc(basic)
	registerHandler := http.HandlerFunc(env.RegisterNewDeviceIndex)
	uploadHandler := http.HandlerFunc(env.UploadRiskEncountersIndex)
	updateHandler := http.HandlerFunc(env.UpdateRiskAssessmentIndex)
	updateCountHandler := http.HandlerFunc(env.UpdateRiskAssessmentGetCountIndex)
	mux.Handle("/", env.TelemetryWrapper(basicHandler))
	mux.Handle("/register", env.TelemetryWrapper(registerHandler))
	mux.Handle("/upload", env.TelemetryWrapper(uploadHandler))
	mux.Handle("/update", env.TelemetryWrapper(updateHandler))
	mux.Handle("/update/count", updateCountHandler)

	// initialize cron job
	c := cron.New()
	_, err = c.AddFunc("@midnight", func() { // starts from the moment this is invoked
		ephIDs, length := cronjobs.GenerateEphemeralIDList(env.db, env.privacyParams, env.mode)
		newFilter, err = cronjobs.CreateNewFilter(ephIDs, length)
		if err != nil {
			log.Println("error updating cuckoo filter")
		}
		env.cf = newFilter
		newChunks, err := cronjobs.CreateChunkedFilters(ephIDs, length)
		if err != nil {
			log.Fatal(err)
		}
		env.cfChunks = newChunks
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
	deviceLocation := req.FormValue("location")
	params, err := routes.RegisterController(deviceType, deviceLocation, env.publicKeyLoc, env.db)
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
	rawNum := req.URL.Query().Get("chunk")
	if rawNum != "" {
		num, err := strconv.Atoi(rawNum)
		if err != nil {
			log.Println(err)
			return
		} else {
			ba := routes.UpdateControllerGetChunk(env.cfChunks, num)
			_, err := w.Write(ba)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		ba := routes.UpdateController(env.cf)
		_, err := w.Write(ba)
		if err != nil {
			log.Println(err)
		}
	}
}

func (env *Env) UpdateRiskAssessmentGetCountIndex(w http.ResponseWriter, req *http.Request) {
	count := routes.UpdateControllerGetCount(env.cfChunks)
	_, err := w.Write(count)
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
