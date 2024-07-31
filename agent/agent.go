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
package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"colmena.bsc.es/agent/device"
	"colmena.bsc.es/agent/docker"
	"colmena.bsc.es/agent/role"

	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(log.Llongfile)
	var local string
	app := &cli.App{
        Name:  "COLMENA",
        Usage: "Start the COLMENA agent",
		Flags: []cli.Flag{
            &cli.StringFlag{
                Name:        "local",
                Value:       "",
                Usage:       "load service description file",
                Destination: &local,
            },
        },
        Action: func(cCtx *cli.Context) error {

			deviceFeature := cCtx.Args().Get(0)
			deviceStrategy := cCtx.Args().Get(1)
			interfc := cCtx.Args().Get(2) 
			deviceName := cCtx.Args().Get(3)
			
            containerEngine := docker.DockerContainerEngine{}
			deviceInfo := device.ParseInformation(deviceFeature, deviceStrategy, deviceName)
			log.Print(deviceInfo)
			deviceInfo.Interfc = interfc
			found := make(chan []role.Role)
			done := handleOsInterrupt()

			go role.Finder(found, local)
			go role.Run(deviceInfo, done, found, containerEngine, role.KpiMatcherFunc(role.UsingDcp))
			waitForever()
			return nil
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func handleOsInterrupt() chan os.Signal {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	return done
}

func waitForever() {
	time.Sleep(1000000000000)
}
