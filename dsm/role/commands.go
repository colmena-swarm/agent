/*
 *  Copyright 2002-2025 Barcelona Supercomputing Center (www.bsc.es)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
// package role
package role

import (
	"encoding/json"
	"log"
	"net/http"

	"colmena.bsc.es/agent/docker"
)

type RoleCommand struct {
	ServiceId string `json:"serviceId"`
	RoleId    string `json:"roleId"`
	ImageId   string `json:"imageId"`
}

type CommandListener struct {
	AgentId         string
	Interfc         string
	ContainerEngine docker.ContainerEngine
}

func (c CommandListener) Endpoints() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		roleCmd, err := parseRoleCommand(r)
		if err != nil {
			log.Printf("Could not parse body. %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Received command to start role. serviceId: %v, roleId: %v, imageId: %v", roleCmd.ServiceId, roleCmd.RoleId, roleCmd.ImageId)
		go c.ContainerEngine.RunContainer(roleCmd.RoleId, roleCmd.ImageId, c.AgentId, c.Interfc)
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		roleCmd, err := parseRoleCommand(r)
		if err != nil {
			log.Printf("Could not parse body. %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Received command to stop role. serviceId: %v, roleId: %v, imageId: %v", roleCmd.ServiceId, roleCmd.RoleId, roleCmd.ImageId)
		go c.ContainerEngine.StopContainer(roleCmd.RoleId)
		w.WriteHeader(http.StatusOK)
	})

	return mux
}

func parseRoleCommand(r *http.Request) (RoleCommand, error) {
	var roleCmd RoleCommand
	if err := json.NewDecoder(r.Body).Decode(&roleCmd); err != nil {
		return RoleCommand{}, err
	}
	return roleCmd, nil
}
