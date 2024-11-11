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
package colmenacontext

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

)

//go:embed colmena_service_definition.json
var colmenaServiceDefinition embed.FS

func PublishColmenaServiceDefinition(agentId string) {
	colmenaServiceDefinition, err := colmenaServiceDefinition.ReadFile("colmena_service_definition.json")
	if err != nil {
		panic(err)
	}

	zenohUrl := os.Getenv("ZENOH_URL")
	if zenohUrl == "" {
		log.Printf("Not publishing colmena service definition because ZENOH_URL is not set")
		return
	}
	if agentId == "" {
		log.Printf("Not publishing colmena service definition because AGENT_ID is not set")
		return
	}

	url := fmt.Sprintf("%s/%s", zenohUrl, "colmena_service_definitions/colmena")
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(colmenaServiceDefinition))
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
		return
	}
	log.Printf("COLMENA service definition published")
}
