package route_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// all tests fall under TestRegistration
func TestUpdate(t *testing.T) {
	srv, env, conf := SetupServer()

	t.Run(
		"update_test_basic",
		func(t *testing.T) {
			req, err := http.NewRequest("GET", GetServerRoute("/update", conf), nil)
			if err != nil {
				t.Fatal(err)
			}
			res := httptest.NewRecorder()
			env.RegisterNewDeviceIndex(res, req)
		},
	)

	err := TeardownServer(srv)
	if err != nil {
		t.Fatal(err)
	}
}