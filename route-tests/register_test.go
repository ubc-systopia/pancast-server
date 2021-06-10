package route_tests

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	server "pancast-server"
	"pancast-server/config"
	"testing"
)

func setupRegister() (*http.Server, *server.Env, config.StartParameters) {
	var conf config.StartParameters
	err := config.ReadJSONConfig("config/app_config.json", &conf)
	if err != nil {
		log.Fatal(err)
	}
	srv, env, _ := server.StartServer(config.GetServerURL(conf))
	return srv, env, conf
}

func teardownRegister(srv *http.Server) {
	_ = srv.Shutdown(context.Background())
}

// all tests fall under TestRegistration
func TestRegistration(t *testing.T) {
	srv, env, conf := setupRegister()

	t.Run(
		"registration_test_basic",
		func(t *testing.T) {
			req, err :=http.NewRequest("POST", GetServerRoute("/register", conf), nil)
			if err != nil {
				t.Fatal(err)
			}
			res := httptest.NewRecorder()
			env.RegisterNewDeviceIndex(res, req)
		},
	)

	t.Run(
		"registration_test_less_basic",
		func(t *testing.T) {

		},
	)

	t.Run(
		"registration_test_least_basic",
		func(t *testing.T) {

		},
	)

	teardownRegister(srv)
}
