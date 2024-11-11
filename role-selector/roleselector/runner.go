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

package roleselector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type RoleRunner interface {
	Run(roleId string, serviceId string, imageId string)
	Stop(roleId string, serviceId string, imageId string)
}

const DSM_URL = "DSM_URL"

type DsmRoleRunner struct {
}

func (r *DsmRoleRunner) executeAction(action, roleId, serviceId, imageId string) {
	dsmUrl := os.Getenv(DSM_URL)
	if dsmUrl == "" {
		log.Fatalf("%v is not set so RoleId: %v, serviceId: %v cannot be %sed. Exiting...", DSM_URL, roleId, serviceId, action)
		return
	}

	data := map[string]string{
		"roleId":    roleId,
		"serviceId": serviceId,
		"imageId":   imageId,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	url := fmt.Sprintf("%s/%s", dsmUrl, action)
	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(jsonData)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("request failed with status code: %d", resp.StatusCode)
	}
}

func (r *DsmRoleRunner) Run(roleId, serviceId, imageId string) {
	r.executeAction("start", roleId, serviceId, imageId)
}

func (r *DsmRoleRunner) Stop(roleId, serviceId, imageId string) {
	r.executeAction("stop", roleId, serviceId, imageId)
}
