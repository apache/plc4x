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

package options

import (
	"sort"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

// OptionReader reads connection-string options and remembers which ones were asked for, so the
// ones nothing asked for can be reported.
//
// Twelve drivers each carried their own copy of the lookup, and none of them could say what had
// been left over: an option nobody read was discarded in silence, which is what turns a typo into
// a setting that appears accepted and does nothing. plc4j reports these
// (DriverBase.warnAboutUnknownParameters); this is how plc4go does.
//
// Recording what was read, rather than declaring what is recognised, is deliberate. A declared
// list is a second description of the same thing and drifts from the code that reads the options;
// a key the driver actually read is recognised by definition.
type OptionReader struct {
	log             zerolog.Logger
	options         map[string][]string
	consumed        map[string]bool
	caseInsensitive bool
}

// NewOptionReader wraps the options of one connection.
func NewOptionReader(log zerolog.Logger, options map[string][]string) *OptionReader {
	return &OptionReader{log: log, options: options, consumed: make(map[string]bool, len(options))}
}

// CaseInsensitive matches option names without regard to case, which the bacnet-ip driver has
// always done - its own lookup lower-cased both sides. The other drivers match exactly, and this
// is opt-in so consolidating the lookups did not quietly widen what any of them accepts.
func (r *OptionReader) CaseInsensitive() *OptionReader {
	r.caseInsensitive = true
	return r
}

// Get is the value of an option, or "" when it was not supplied. Asking marks the option as
// consumed, whether or not it was supplied - a driver that asks for an option recognises it.
func (r *OptionReader) Get(key string) string {
	r.consumed[key] = true
	optionValues, ok := r.options[key]
	if !ok && r.caseInsensitive {
		for suppliedKey, values := range r.options {
			if strings.EqualFold(suppliedKey, key) {
				r.consumed[suppliedKey] = true
				optionValues, ok = values, true
				break
			}
		}
	}
	if !ok || len(optionValues) == 0 {
		return ""
	}
	if len(optionValues) > 1 {
		r.log.Warn().Str("key", key).Msg("Option must be unique")
	}
	return optionValues[0]
}

// Ignore marks options as belonging to someone else - the transport, or the driver's own nested
// parsing - so they are not reported as unknown. The transport options every driver's connection
// string may carry are ignored already; this is for what a particular driver adds.
func (r *OptionReader) Ignore(keys ...string) {
	for _, key := range keys {
		r.consumed[key] = true
	}
}

// ReportUnknown logs one warning naming every supplied option that nothing read.
//
// It is a warning and never an error: a stray option must not break a connection that would
// otherwise work, which is the rule plc4j settled on for the same report. The operator is told,
// and decides.
func (r *OptionReader) ReportUnknown(protocolCode string) {
	var unknown []string
	for key := range r.options {
		if r.consumed[key] || isTransportOption(key) {
			continue
		}
		unknown = append(unknown, key)
	}
	if len(unknown) == 0 {
		return
	}
	sort.Strings(unknown)
	for _, key := range unknown {
		event := r.log.Warn().Str("option", key).Str("driver", protocolCode)
		if suggestion := r.suggestionFor(key); suggestion != "" {
			event = event.Str("didYouMean", suggestion)
		}
		event.Msg("Connection string option is not known to this driver and is ignored")
	}
}

// suggestionFor is the option the given unknown one was most likely meant to be, or "" when
// nothing read is close enough to be worth naming. Only options this driver read are candidates,
// so the suggestion cannot point at a name the driver would ignore anyway.
func (r *OptionReader) suggestionFor(unknown string) string {
	// Roughly one edit per four characters, so short names do not match everything. The same
	// budget plc4j's suggestionFor uses.
	budget := len(unknown) / 4
	if budget < 1 {
		budget = 1
	} else if budget > 3 {
		budget = 3
	}
	best, bestDistance := "", budget+1
	for candidate := range r.consumed {
		distance := editDistance(unknown, candidate)
		if distance > budget {
			continue
		}
		if distance < bestDistance || (distance == bestDistance && candidate < best) {
			best, bestDistance = candidate, distance
		}
	}
	return best
}

// transportOptions are read by the transport layer rather than by a driver, so a driver must not
// report them as unknown.
//
// Each transport registers its own names beside the code that reads them, rather than this file
// carrying one list of every transport's options. A central list is a second description of what
// the transports read and drifts from them - the very thing OptionReader exists to avoid on the
// driver side - and it exempts every name from every transport on every connection, so a serial
// option on a TCP connection would pass unremarked.
//
// A transport that is not linked in registers nothing, and its options are then reported as
// unknown. That is the honest answer: a connection whose transport is not present cannot be
// using them. What this still cannot do is subtract what the transport *actually consumed*, the
// way plc4j does at its connect site: a Go driver parses its configuration before the transport
// instance exists, so at report time there is nothing yet to ask.
var (
	transportOptionsMutex sync.RWMutex
	transportOptions      = map[string]bool{}
)

// RegisterTransportOptions records the connection-string options a transport reads. Call it from
// the transport's package initialisation, listing the names that package looks up.
func RegisterTransportOptions(keys ...string) {
	transportOptionsMutex.Lock()
	defer transportOptionsMutex.Unlock()
	for _, key := range keys {
		transportOptions[strings.ToLower(key)] = true
	}
}

// IsTransportOption says whether some linked-in transport reads the given option name. Drivers
// using an OptionReader get this applied for them; it is exported for the OPC UA driver, which
// matches option names against its Configuration's fields by reflection rather than by reading
// them one at a time, and so does its own reporting.
func IsTransportOption(key string) bool {
	return isTransportOption(key)
}

func isTransportOption(key string) bool {
	transportOptionsMutex.RLock()
	defer transportOptionsMutex.RUnlock()
	return transportOptions[strings.ToLower(key)]
}

// editDistance is the Levenshtein distance between two option names.
func editDistance(left, right string) int {
	previous := make([]int, len(right)+1)
	current := make([]int, len(right)+1)
	for j := range previous {
		previous[j] = j
	}
	for i := 1; i <= len(left); i++ {
		current[0] = i
		for j := 1; j <= len(right); j++ {
			substitution := previous[j-1]
			if left[i-1] != right[j-1] {
				substitution++
			}
			current[j] = min(substitution, min(previous[j]+1, current[j-1]+1))
		}
		previous, current = current, previous
	}
	return previous[len(right)]
}
