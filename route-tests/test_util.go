package route_tests

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	server "pancast-server"
	"pancast-server/config"
)

func GetServerRoute(route string, conf config.StartParameters) string {
	return "https://" + config.GetServerURL(conf) + route
}

func ReadJSONPayload(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SetupServer() (*http.Server, *server.Env, config.StartParameters) {
	var conf config.StartParameters
	err := config.ReadJSONConfig("../config/test_app_config.json", &conf)
	if err != nil {
		log.Fatal(err)
	}
	srv, env, _ := server.StartServer(conf)
	return srv, env, conf
}

func TeardownServer(srv *http.Server) error {
	err := srv.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}
