package roleselector

import (
	"encoding/json"
	"log"
	"net/http"
)

type StoppedEvent struct {
	RoleId		string `json:"roleId"`
	ServiceId	string `json:"serviceId"`
	ImageId		string `json:"imageId"`
}

func MonitorStopped(stoppedChan chan StoppedEvent, mux *http.ServeMux) {
	mux.HandleFunc("/stopped", func(w http.ResponseWriter, r *http.Request) {
		var stopped StoppedEvent
		if err := json.NewDecoder(r.Body).Decode(&stopped); err != nil {
			log.Printf("Could not parse body. %s", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			stoppedChan <- stopped
		}
	})
}