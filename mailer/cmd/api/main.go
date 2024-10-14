package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct{}

const webPort = "80"

func main() {
	app := Config{}
	log.Println("Starting the mail server at port", 80)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.route(),
	}

	log.Printf("Mailer service up and running at port %s", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println("There hass been an error creating the mail server", err)
		os.Exit(1)
	}

}
