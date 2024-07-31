/*
 *  Copyright 2002-2024 Barcelona Supercomputing Center (www.bsc.es)
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
package role

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type ServiceDefinitionId struct {
	Value string `json:"value"` 
}

type ServiceDefinition struct {
	Id ServiceDefinitionId `json:"id"`
	Kpis json.RawMessage `json:"kpis"`
	DockerRoleDefinitions []Role `json:"dockerRoleDefinitions"`
}

func Finder(found chan []Role, local string) {
	loadFromDisk(found, local)
	httpServer(found)
}

func loadFromDisk(found chan []Role, location string) {
	if location == "" {
		return
	}

	log.Printf("Loading service definition from %s", location)

	jsonData, err := os.ReadFile(location)
    if err != nil {
        log.Fatalf("Error opening file: %s. %s\n", location, err)
    }
	
	serviceDefinition := ServiceDefinition{}
	err = json.Unmarshal([]byte(jsonData), &serviceDefinition)
	if err != nil {
		log.Fatalf("Could not parse service definition: %s", err)
		return
	}
	addServiceId(serviceDefinition.DockerRoleDefinitions, serviceDefinition.Id.Value)
	log.Printf("Parsed service definition from disk: %v\n", serviceDefinition.DockerRoleDefinitions)
	found <- serviceDefinition.DockerRoleDefinitions
}

func httpServer(found chan []Role) {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		serviceDefinition := ServiceDefinition{}
        err := json.NewDecoder(r.Body).Decode(&serviceDefinition)
		if err != nil {
			log.Printf("Could not parse service definition. %s", err)
			return
		}
		addServiceId(serviceDefinition.DockerRoleDefinitions, serviceDefinition.Id.Value)
		log.Printf("Received service definition: %v\n", serviceDefinition.DockerRoleDefinitions)
		found <- serviceDefinition.DockerRoleDefinitions
    })
	
    listener, err := net.Listen("tcp", ":50551")
	if err != nil {
		panic(err)
	}

	log.Printf("Listening for service descriptions on port: %v", listener.Addr().(*net.TCPAddr).Port)
	panic(http.Serve(listener, nil))
}

func addServiceId(roles []Role, serviceId string ) {
	for i := range roles {
		roles[i].ServiceId = strings.ToLower(serviceId)
	}
}
