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
package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

func (DockerContainerEngine) Subscribe(ctx context.Context) {
	cli, cxnErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if cxnErr != nil {
		panic(cxnErr)
	}
	defer cli.Close()

	events, err := cli.Events(ctx, types.EventsOptions{})
	process(ctx, events, err)
}

func process(ctx context.Context, events <-chan events.Message, err <-chan error) {
	log.Printf("Processing docker events")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-events:
			if event.Action == "die" {
				roleStopped(event)
			}
		case each := <-err:
			//once the stream has been completely read an io.EOF error will be sent over the error channel
			if each != io.EOF {
				log.Printf("Error processing docker event: %s", each)
			}
			return
		}
	}
}

type StoppedEvent struct {
	RoleId		string `json:"roleId"`
	ServiceId	string `json:"serviceId"`
	ImageId		string `json:"imageId"`
}

func roleStopped(event events.Message) {
	if event.Actor.Attributes["es.bsc.colmena.roleId"] == "" ||
		event.Actor.Attributes["es.bsc.colmena.serviceId"] == "" {
		return
	}

	// send the stopped event to the role selector
	roleSelectorUrl := os.Getenv("ROLE_SELECTOR_URL")
	if roleSelectorUrl == "" {
		roleSelectorUrl = "http://role-selector:5555"
	}

	url := fmt.Sprintf("%s/%s", roleSelectorUrl, "stopped")

	stoppedEvent := StoppedEvent{
		RoleId: event.Actor.Attributes["es.bsc.colmena.roleId"],
		ServiceId: event.Actor.Attributes["es.bsc.colmena.serviceId"],
		ImageId: event.Actor.Attributes["image"],
	}

	jsonData, err := json.Marshal(stoppedEvent)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
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
	log.Printf("Stopped event sent to role selector")
}
