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
	"iter"
	"slices"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Collector interface {
	fmt.Stringer
	utils.Serializable
	CapabilityFunctions(fn string) iter.Seq[GenericFunction]
	SearchCapability(subs ...CollectorOrCapability) []Capability
	AddCapability(cls Capability)
}

// TODO: implement
//
//go:generate go tool plc4xGenerator -type=collector -prefix=capability_
type collector struct {
	capabilities []Capability

	log zerolog.Logger
}

var _ Collector = (*collector)(nil)

func NewCollector(localLog zerolog.Logger, subs ...CollectorOrCapability) (_collector Collector, init func()) {
	c := &collector{log: localLog}
	return c, func() {
		// gather the capabilities
		c.capabilities = c.SearchCapability(subs...)
	}
}

func (c *collector) SearchCapability(subs ...CollectorOrCapability) []Capability {
	if c.log.Debug().Enabled() {
		subsStringers := make([]fmt.Stringer, len(subs))
		for i, sub := range subs {
			subsStringers[i] = sub
		}
		c.log.Debug().Stringers("subs", subsStringers).Msg("SearchCapability")
	}
	var result []Capability
	for _, sub := range subs {
		if sub.IsCollector {
			result = append(result, (*sub.Collector).SearchCapability(sub)...)
		} else if c.IsCapability() {
			result = append(result, *sub.Capability)
		} else {
			panic("impossible")
		}
	}
	return result
}

// CapabilityFunctions generator yields functions that match the requested capability sorted by z-index.
func (c *collector) CapabilityFunctions(fn string) iter.Seq[GenericFunction] {
	c.log.Debug().Str("fn", fn).Msg("CapabilityFunctions")

	// build a list of functions to call
	type fnEntry struct {
		_zindex int
		fn      GenericFunction
	}
	var fns []fnEntry
	for _, _capability := range c.capabilities {
		xfn := _capability.getFN(fn)
		c.log.Trace().Stringer("capability", _capability).Bool("xfn", xfn != nil).Msg("cap")
		if xfn != nil {
			fns = append(fns, fnEntry{_zindex: _capability.getZIndex(), fn: xfn})
		}
	}
	// sort them by z-index
	slices.SortFunc(fns, func(a fnEntry, b fnEntry) int {
		return a._zindex - b._zindex
	})
	c.log.Debug().Interface("fns", fns).Msg("fns")

	return func(yield func(function GenericFunction) bool) {
		// now yield them in order
		for xindx, xfn := range fns {
			c.log.Debug().Int("xindx", xindx).Msg("yield")
			yield(xfn.fn)
		}
	}
}

func (c *collector) AddCapability(cls Capability) {
	c.log.Debug().Stringer("cls", cls).Msg("AddCapability")
	c.capabilities = append(c.capabilities, cls)
	// TODO: not sure what to do here
	return
}

func (c *collector) IsCollector() bool {
	return true
}

func (c *collector) IsCapability() bool {
	return false
}
