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
	"os"
	"time"

	"colmena.bsc.es/agent/device"
	"colmena.bsc.es/agent/docker"
	"golang.org/x/exp/maps"
)

type KPI struct {
	Key           string
	Threshold     float64
	FromUnit      float64
	FromType	  string
	Comparison	  string
}

func (kpi *KPI) UnmarshalJSON(data []byte) error {
	var parsedKpi string
	if err := json.Unmarshal(data, &parsedKpi); err != nil {
		return err
	}
	*kpi = parseKpi(parsedKpi)
	return nil
}

type Role struct {
	ServiceId			 string
	Id                   string 	`json:"id"`
	ImageId              string 	`json:"imageId"`
	HardwareRequirements []string 	`json:"hardwareRequirements"`
	Kpis                 []KPI 		`json:"kpis"`
}

func Run(device device.Info, done chan os.Signal, found chan []Role, containerEngine docker.ContainerEngine, kpiMatcher KpiMatcher) {
	all := make(map[string]Role)
	running := make(map[string]Role)
	containerIdToRole := make(map[string]Role)
	roleIdToContainerId := make(map[string]string)
	start := make(chan Role)
	stop := make(chan Role)
	stopped := make(chan string)
	ticker := time.NewTicker(5 * time.Second)
	go containerEngine.Subscribe(stopped)
	for {
		select {
		case roles := <-found:
			for _, role := range roles {
				_, contains := running[role.Id]
				if !contains {
					all[role.Id] = role
				}
			}
			go roleMatch(device, all, running, kpiMatcher, start, stop)
		case starting := <-start:
			running[starting.Id] = starting
			containerId, err := containerEngine.RunContainer(starting.ImageId, device.Interfc, device)
			if (err != nil) {
				log.Printf("Could not start roleId: %v, %v", starting.Id, err)
			} else {
				containerIdToRole[containerId] = starting
				roleIdToContainerId[starting.Id] = containerId
				log.Printf("Started roleId: %v, imageId: %v, containerId: %v", starting.Id, starting.ImageId, containerId)
			}
		case stopping := <-stop:
			log.Printf("Stopping roleId: %v", stopping.Id)
			delete(running, stopping.Id)
			containerEngine.StopContainer(roleIdToContainerId[stopping.Id])
			log.Printf("Stopped roleId: %v", stopping.Id)
		case stoppedContainerid := <-stopped:
			stoppedRole, contains := containerIdToRole[stoppedContainerid]
			if contains {
				log.Printf("Container stopped. roleId: %v containerId: %v", stoppedRole.Id, stoppedContainerid)
				delete(running, stoppedRole.Id)
			}
		case <-ticker.C:
			go func() {
				log.Println("Running periodic role match...")
				roleMatch(device, all, running, kpiMatcher, start, stop)
				log.Println("Finished periodic role match")
			}()
		case <-done:
			log.Println("Agent stopping. Stopping all containers")
			for _, containerId := range maps.Keys(containerIdToRole) {
				log.Printf("Stopping %v", containerId)
				containerEngine.StopContainer(containerId)
			}
			log.Println("Goodbye")
			os.Exit(0)
		}
	}
}

func roleMatch(device device.Info, all map[string]Role, running map[string]Role, kpiMatcher KpiMatcher, start chan Role, stop chan Role) {
	toStop := matchstop(device, maps.Values(running), kpiMatcher)
	for _, each := range toStop {
		stop <- each
	}
	toStart := matchstart(device, notRunning(all, running), kpiMatcher)
	for _, each := range toStart {
		start <- each
	}
}

func notRunning(all map[string]Role, running map[string]Role) []Role {
	var notRunning = []Role{}
	for _, role := range all {
		_, contains := running[role.Id]
		if !contains {
			notRunning = append(notRunning, role)
		}
	}
	return notRunning
}
