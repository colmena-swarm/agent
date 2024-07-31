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
	"sync"

	"colmena.bsc.es/agent/zenohclient"
)

type KpiMatcher interface { match(kpi KPI) bool }
type KpiMatcherFunc func(kpi KPI) bool
func (kpiMatcherFunc KpiMatcherFunc) match(kpi KPI) bool { return kpiMatcherFunc(kpi) }
var UsingDcp = func(kpi KPI) bool { return zenohclient.MetricsMet(convertKpi(kpi)) }

func kpisMet(kpis []KPI, kpiMatcher KpiMatcher) bool {
	var wg sync.WaitGroup
	result := make(chan bool, len(kpis))
	for _, kpi := range kpis {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result <- kpiMatcher.match(kpi)
		}()
	}
	wg.Wait()
	close(result)
	for each := range result {
		if !each {
			return false
		}
	}
	return true
}
