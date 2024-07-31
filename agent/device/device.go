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
package device

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Info struct {
	Strategy string
	Features []string
	Interfc string
	Name string
}

func ParseInformation(feature string, strategy string, deviceName string) Info {
	deviceFeatures := []string{feature}
	if (len(deviceName) == 0) {
		deviceName = Hostname()
	}
	return Info{Strategy: strategy, Features: deviceFeatures, Name: deviceName}
}

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Could not get hostname from OS.")
		fmt.Println(err)
		os.Exit(1)
	}
	return hostname
}

func IpAddress() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().(*net.UDPAddr)
    return localAddr.IP
}