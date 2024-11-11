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

package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"colmena.bsc.es/role-selector/policy"
	"colmena.bsc.es/role-selector/roleselector"
	"colmena.bsc.es/role-selector/servicedescription"
	"colmena.bsc.es/role-selector/sla"
	"colmena.bsc.es/role-selector/types"
)

func main() {
	mux := http.NewServeMux()

	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		panic(err)
	}
	log.Printf("HTTP server listening on port: %v", listener.Addr().(*net.TCPAddr).Port)

	sdChan := make(chan types.ServiceDescription)
	servicedescription.Endpoint(sdChan, mux)
	alertsChan := make(chan types.Alert)
	sla.AlertEndpoint(alertsChan, mux)
	go http.Serve(listener, mux)

	policy := getPolicy()
	log.Printf("Using policy: %v", policy.Name())
	selector := &roleselector.RoleSelector{
		ServiceDescriptionChan: sdChan,
		AlertsChan:             alertsChan,
		Hardware:               os.Getenv("HARDWARE"),
		Policy:                 policy,
		RoleRunner:             &roleselector.DsmRoleRunner{},
		KpiRetriever:           sla.KpiRetrieverClient{},
	}
	selector.Run()

	log.Println("Role selector stopping...")
	close(alertsChan)
	close(sdChan)
}

func getPolicy() policy.Policy {
	policyName := os.Getenv("POLICY")
	if policyName == "lazy" {
		return &policy.LazyPolicy{}
	}
	return &policy.EagerPolicy{}
}
