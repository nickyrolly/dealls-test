package healthcheck

import (
	"encoding/json"
	"net/http"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s := struct {
		Name   string `json:"name"`
		Status int    `json:"status"`
	}{
		Name:   "Server Up1 Connection Successfully",
		Status: 200,
	}

	response, _ := json.Marshal(s)

	w.Write(response)
}
