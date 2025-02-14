package common

import (
	"encoding/json"
	"net/http"
	"time"
)

type CommonConfig struct {
	EnableExaminations  bool
	Origin              string
	Location            *time.Location
	SignupExpirationDay int
}

func CustomResponseAPI(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeLen, writeErr := w.Write([]byte(`{"errors": "Internal Server Error"}`))
		if writeErr != nil {
			return writeLen, writeErr
		}
		return writeLen, err
	}
	w.WriteHeader(status)
	return w.Write(b)
}
