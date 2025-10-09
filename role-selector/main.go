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
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf("HTTP server listening on port: %v", listener.Addr().(*net.TCPAddr).Port)

	sdChan := make(chan types.ServiceDescription)
	servicedescription.Endpoint(sdChan, mux)
	alertsChan := make(chan types.Alert)
	sla.AlertEndpoint(alertsChan, mux)
	stoppedRoles := make(chan roleselector.StoppedEvent)
	roleselector.MonitorStopped(stoppedRoles, mux)
	go http.Serve(listener, mux)

	selector := &roleselector.RoleSelector{
		ServiceDescriptionChan: sdChan,
		AlertsChan:             alertsChan,
		Hardware:               os.Getenv("HARDWARE"),
		RoleRunner:             &roleselector.DsmRoleRunner{},
		KpiRetriever:           sla.KpiRetrieverClient{},
	}

	policy := getPolicy(selector)
	log.Printf("Using policy: %v", policy.Name())
	selector.SetPolicy(policy)

	selector.Run(stoppedRoles)

	log.Println("Role selector stopping...")
	close(alertsChan)
	close(sdChan)
	policy.Stop()
}

func getPolicy(rs *roleselector.RoleSelector) policy.Policy {
	policyName := os.Getenv("POLICY")
	if policyName == "lazy" {
		return &policy.LazyPolicy{}
	} else if policyName == "consensus" {
		endpoint := os.Getenv("ENDPOINT")
		p, err := policy.NewConsensusPolicy(endpoint, rs.TriggerDecision)
		if err != nil {
			log.Printf("failed to create consensus policy (%v), falling back to LazyPolicy", err)
			return &policy.LazyPolicy{}
		}
		return p
	}
	return &policy.EagerPolicy{}
}
