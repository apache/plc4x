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
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

// TODO: implement
//
//go:generate plc4xGenerator -type=Collector -prefix=capability_
type Collector struct {
	capabilities []*Capability

	log zerolog.Logger
}

func NewCollector(localLog zerolog.Logger) *Collector {
	return &Collector{log: localLog}
}

func (c *Collector) searchCapability() {
	panic("not implemented") // TODO: implement me
}

// CapabilityFunctions generator yields functions that match the requested capability sorted by z-index.
func (c *Collector) CapabilityFunctions(fn string) []func(args Args, kwargs KWArgs) error {
	c.log.Trace().Msg("CapabilityFunctions")

	// build a list of functions to call
	var fns []func(args Args, kwargs KWArgs) error
	for _, capability := range c.capabilities {
		xfn := capability.getFN(fn)
		c.log.Trace().Stringer("capability", capability).Bool("xfn", xfn != nil).Msg("cap")
		if xfn != nil {
			// TODO: sorting
			fns = append(fns, xfn)
		}
	}

	// sort them by z-index
	// TODO: sorting

	// now yield them in order
	// TODO: what?

	return fns
}

func (c *Collector) AddCapability(cls any) {
	// TODO: implement
	return
}
