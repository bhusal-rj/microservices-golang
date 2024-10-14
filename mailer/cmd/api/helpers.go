package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	// Limit the size of the maximum json request data
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	//difference between unmarshal and decoder is that unmarshall requires the data to be in the form of json
	//but the decoder also accepts the stream of json response.
	// Here the decoder is used for parsing the json value.
	//decodes the one json value at one time
	desc := json.NewDecoder(r.Body)

	//decode the data to get the original data in the data. Decoder reads the data from the stream
	err := desc.Decode(data)

	if err != nil {
		return err
	}

	//attemps to decode again. In case of no data the decoder throws an error
	err = desc.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only single JSON value")
	}
	return nil
}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	fmt.Println(err.Error(), "There is an error")
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}
	return app.writeJSON(w, statusCode, payload)
}
