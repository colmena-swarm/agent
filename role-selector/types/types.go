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

package types

import (
	"log"
)

type Id struct {
	Value string `json:"value"`
}

type DockerContextDefinition struct {
	Id      string `json:"id"`
	ImageId string `json:"imageId"`
}

type DockerRoleDefinition struct {
	Id                   string           `json:"id"`
	ImageId              string           `json:"imageId"`
	HardwareRequirements []string         `json:"hardwareRequirements"`
	Kpis                 []KpiDescription `json:"kpis"`
}

type ServiceDescription struct {
	ServiceId                Id                        `json:"id"`
	DockerContextDefinitions []DockerContextDefinition `json:"dockerContextDefinitions"`
	Kpis                     []KpiDescription          `json:"kpis"`
	DockerRoleDefinitions    []DockerRoleDefinition    `json:"dockerRoleDefinitions"`
}

type KpiDescription struct {
	Query string `json:"query"`
	Scope string `json:"scope"`
}

type Operator string

const (
	EqualTo         Operator = "=="
	NotEqualTo      Operator = "!="
	LessThan        Operator = "<"
	LessThanOrEq    Operator = "<="
	GreaterThan     Operator = ">"
	GreaterThanOrEq Operator = ">="
)

type Role struct {
	Id        string
	ImageId   string
	IsRunning bool
	Resources []Resource
}

type KPI struct {
	Query          string
	Value          float64
	Threshold      float64
	Operator       Operator
	AssociatedRole string
	Level          string
}

type Resource struct {
	Name  string
	Value int
}

type KPIQuery struct {
	RoleId          string  `json:"roleId"`
	Query           string  `json:"query"`
	Level           string  `json:"level"`
	Value           float64 `json:"value"`
	TotalViolations int     `json:"total_violations"`
}

type Query struct {
	Kpis []KPIQuery `json:"KPIs"`
}

type Alert struct {
	ServiceId string     `json:"serviceId"`
	SlaId     string     `json:"slaId"`
	Kpis      []KPIQuery `json:"KPIs"`
}

// Alerts represents an array of Alert objects
type Alerts []Alert

// KPIQueries represents an array of KPIQuery objects
type KPIQueries []Alert

type Response struct {
	Message  string  `json:"Message"`
	Method   string  `json:"Method"`
	Resp     string  `json:"Resp"`
	Response []Alert `json:"Response"`
}

func ParseKpis(kpis []KPIQuery) []KPI {
	var kpiStructs []KPI
	for _, kpi := range kpis {
		kpiStruct, err := ParseToKpiStruct(kpi)
		if err != nil {
			log.Printf("Failed to Parse KPI.")
		}
		kpiStructs = append(kpiStructs, *kpiStruct)
	}
	return kpiStructs
}

func ParseToKpiStruct(kpi KPIQuery) (*KPI, error) {
	kpiStruct := KPI{
		Query:          kpi.Query,
		Value:          kpi.Value,
		Operator:       EqualTo,
		AssociatedRole: kpi.RoleId,
		Level:          kpi.Level,
	}
	return &kpiStruct, nil
}
