package main

import (
	"fmt"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	err := app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		fmt.Println(err)
	}
}
