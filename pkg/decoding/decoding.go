package decoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	//"github.com/golang/gddo/httputil/header"
)

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	value := r.Header.Get("Content-Type")
	if value != "" {
		if value != "application/json" {
			fmt.Println(r.Header.Get("Content-Type"))
			return fmt.Errorf("Content-Type header is not application/json. Error: %v", http.StatusUnsupportedMediaType)
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return fmt.Errorf("%s. Error: %v", msg, http.StatusBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return fmt.Errorf("%s. Error: %v", msg, http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return fmt.Errorf("%s. Error: %v", msg, http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return fmt.Errorf("%s. Error: %v", msg, http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return fmt.Errorf("%s. Error: %v", msg, http.StatusBadRequest)

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return fmt.Errorf("%s. Error: %v", msg, http.StatusRequestEntityTooLarge)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		return fmt.Errorf("%s. Error: %v", msg, http.StatusBadRequest)
	}

	return nil
}
