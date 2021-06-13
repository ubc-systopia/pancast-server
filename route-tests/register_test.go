package route_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)



// all tests fall under TestRegistration
func TestRegistration(t *testing.T) {
	srv, env, conf := SetupServer()

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

	TeardownServer(srv)
}
