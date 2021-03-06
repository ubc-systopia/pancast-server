package main

/*
	Entry point of application
*/

import (
	"context"
	"fmt"
	"log"
	server "pancast-server"
	"pancast-server/config"
)

func main() {
	var conf config.StartParameters
	err := config.ReadJSONConfig("config/app_config.json", &conf)
	if err != nil {
		log.Fatal(err)
	}
	srv, _, done := server.StartServer(conf)
	<-done // waiting for SIGKILL or interrupt signals
	fmt.Println("\nShutting down...")
	err = srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
