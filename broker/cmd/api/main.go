package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "3004"

type Config struct{}

func main() {
	app := Config{}
	log.Printf("Starting the broker service at port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.route(),
	}

	//start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic("There has been an error", err)
	}
}
