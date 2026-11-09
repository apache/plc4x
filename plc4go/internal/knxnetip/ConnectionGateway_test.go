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

package knxnetip

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

// Tests of the connection against fakeKnxGateway, over the real codec and udp transport.

// Connecting used to panic on every gateway: each service family of the SEARCH_RESPONSE was
// type-asserted to an interface with GetTypeName, which no generated model type has.
func Test_Connection_Connect_listsTheServiceFamilies(t *testing.T) {
	gateway := newFakeKnxGateway(t)
	connection := connectToFakeGateway(t, t.Context(), gateway)

	assert.Equal(t, []string{"KnxNetIpCore v1", "KnxNetIpTunneling v1"}, connection.metadata.SupportedServices)
	assert.Empty(t, gateway.Violations())
}

// The connection created its codec without the options, so the codec logged to the global
// logger instead of the connection's.
func Test_Connection_codecUsesTheConnectionLogger(t *testing.T) {
	var global bytes.Buffer
	previous := log.Logger
	log.Logger = zerolog.New(&global).Level(zerolog.TraceLevel)
	t.Cleanup(func() { log.Logger = previous })

	gateway := newFakeKnxGateway(t)
	connectToFakeGateway(t, t.Context(), gateway)

	assert.Empty(t, global.String())
}
