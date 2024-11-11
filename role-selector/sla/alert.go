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

package sla

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"colmena.bsc.es/role-selector/types"
)

func AlertEndpoint(alertsChan chan types.Alert, mux *http.ServeMux) {
	mux.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		var alerts types.Alerts // assume types.Alerts is []types.Alert
		if err := json.Unmarshal(body, &alerts); err != nil {
			// Try unmarshaling a single alert instead
			var singleAlert types.Alert
			if err := json.Unmarshal(body, &singleAlert); err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			alerts = append(alerts, singleAlert)
		}

		for _, alert := range alerts {
			alertsChan <- alert
		}

		w.WriteHeader(http.StatusOK)
	})
}
