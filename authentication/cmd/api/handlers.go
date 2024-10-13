package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	//log authentication
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.logRequest("Authentication logs", fmt.Sprintf("Failed login attempt from %s", user.Email))
		app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}
	err = app.logRequest("Authentication logs", fmt.Sprintf("Successful login attempt from %s", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged In Successfully with user %v", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.Marshal(entry)
	//log service
	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewReader(jsonData))
	if err != nil {
		log.Print("Error creating the authentication logs")
		return err
	}
	client := http.Client{}
	_, err = client.Do(request)
	return err
}
