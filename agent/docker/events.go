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
package docker

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

func (DockerContainerEngine) Subscribe(stopped chan string) {
	ctx := context.Background()
	cli, cxnErr := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if cxnErr != nil {
		panic(cxnErr)
	}
	defer cli.Close()

	for {
		events, err := cli.Events(ctx, types.EventsOptions{})
		process(events, err, stopped)
	}
}

func process(events <-chan events.Message, err <-chan error, stopped chan string) {
	for {
		select {
		case event := <-events:
			if event.Action == "die" {
				stopped <- event.ID
			}
		case each := <-err:
			//Once the stream has been completely read an io.EOF error will be sent over the error channel
			if each != io.EOF {
				log.Printf("Error processing docker event: %s", each)
			}
			return
		}
	}
}