/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package capability

import (
	"fmt"
)

type CollectorOrCapability struct {
	Collector    *Collector
	IsCollector  bool
	Capability   *Capability
	IsCapability bool
}

var (
	_ fmt.Stringer = (*CollectorOrCapability)(nil)
)

func CoCCo(co Collector) CollectorOrCapability {
	return CollectorOrCapability{&co, true, nil, false}
}

func CoCCa(ca Capability) CollectorOrCapability {
	return CollectorOrCapability{nil, false, &ca, true}
}

func (c CollectorOrCapability) String() string {
	switch {
	case c.IsCollector:
		return (*c.Collector).String()
	case c.IsCapability:
		return (*c.Capability).String()
	default:
		panic("impossible")
	}
}
