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
	"fmt"
	"log"
	"net/http"
	"os"

	"colmena.bsc.es/role-selector/types"
)

const SLA_MANAGER_URL = "SLA_MANAGER_URL"
const defaultSlaManagerUrl = "http://localhost:8081"

type KpiRetriever interface {
	Get(serviceId string) ([]types.KPI, error)
}

type KpiRetrieverClient struct{}

func getSlaManagerUrl() string {
	url := os.Getenv(SLA_MANAGER_URL)
	if url == "" {
		log.Printf("SLA_MANAGER_URL is not set, using default: %s", defaultSlaManagerUrl)
		return defaultSlaManagerUrl
	}
	return url
}
func (k KpiRetrieverClient) Get(serviceId string) ([]types.KPI, error) {
	url := fmt.Sprintf("%s/api/v1/kpis/%s", getSlaManagerUrl(), serviceId)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to get KPI: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var response types.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("Failed to decode KPI: %v", err)
		return nil, err
	}

	alerts := response.Response
	var kpis []types.KPI
	//parse kpis to kpi structs
	for _, alert := range alerts {
		kpis = append(kpis, types.ParseKpis(alert.Kpis)...)
	}
	return kpis, nil
}
