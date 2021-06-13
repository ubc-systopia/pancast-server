package main

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
	srv, _, done := server.StartServer(config.GetServerURL(conf))
	<-done
	fmt.Println("\nShutting down...")
	err = srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
