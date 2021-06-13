package route_tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpload(t *testing.T) {
	srv, env, conf := SetupServer()

	t.Run(
		"upload_test_basic",
		func(t *testing.T) {
			payload, err := ReadJSONPayload("uploadPayloads/uploadBasic.json")
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", GetServerRoute("/upload", conf), strings.NewReader(payload))
			if err != nil {
				t.Fatal(err)
			}
			res := httptest.NewRecorder()
			env.UploadRiskEncountersIndex(res, req)
		},
	)

	t.Run(
		"upload_test_less_basic",
		func(t *testing.T) {

		},
	)

	TeardownServer(srv)
}
