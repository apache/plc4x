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

package modbus

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseConfiguration(t *testing.T, connectionOptions map[string][]string) Configuration {
	t.Helper()
	configuration, err := ParseFromOptions(zerolog.Nop(), connectionOptions)
	require.NoError(t, err)
	return configuration
}

// A connection string that says nothing gets what plc4j's ModbusTcpConfiguration defaults to.
func TestParseFromOptions_defaults(t *testing.T) {
	configuration := parseConfiguration(t, map[string][]string{})

	assert.Equal(t, uint8(1), configuration.unitIdentifier)
	assert.Equal(t, BigEndianOrder, configuration.defaultPayloadByteOrder)
	assert.Equal(t, "4x00001:BOOL", configuration.pingAddress)
	assert.Equal(t, 5*time.Second, configuration.requestTimeout)
	assert.Equal(t, DefaultConfiguration(), configuration)
}

// plc4j spells the option default-unit-identifier, the Go driver has always called it
// unit-identifier. Both work, so that neither an existing connection string nor one copied from
// plc4j breaks.
func TestParseFromOptions_unitIdentifierAndItsAlias(t *testing.T) {
	assert.Equal(t, uint8(9), parseConfiguration(t, map[string][]string{"unit-identifier": {"9"}}).unitIdentifier)
	assert.Equal(t, uint8(9), parseConfiguration(t, map[string][]string{"default-unit-identifier": {"9"}}).unitIdentifier)

	// With both spelled out the plc4j one wins.
	both := parseConfiguration(t, map[string][]string{"unit-identifier": {"9"}, "default-unit-identifier": {"3"}})
	assert.Equal(t, uint8(3), both.unitIdentifier)
}

func TestParseFromOptions_defaultPayloadByteOrder(t *testing.T) {
	for _, name := range []string{"BIG_ENDIAN", "LITTLE_ENDIAN", "BIG_ENDIAN_BYTE_SWAP", "LITTLE_ENDIAN_BYTE_SWAP"} {
		t.Run(name, func(t *testing.T) {
			expected, ok := ByteOrderByName(name)
			require.True(t, ok)
			configuration := parseConfiguration(t, map[string][]string{"default-payload-byte-order": {name}})
			assert.Equal(t, expected, configuration.defaultPayloadByteOrder)
		})
	}
}

func TestParseFromOptions_pingAddress(t *testing.T) {
	configuration := parseConfiguration(t, map[string][]string{"ping-address": {"holding-register:5:INT"}})
	assert.Equal(t, "holding-register:5:INT", configuration.pingAddress)
}

// The request timeout is stated in milliseconds, as it is in plc4j.
func TestParseFromOptions_requestTimeout(t *testing.T) {
	configuration := parseConfiguration(t, map[string][]string{"request-timeout-ms": {"250"}})
	assert.Equal(t, 250*time.Millisecond, configuration.requestTimeout)
}

// An option the driver can't make sense of is an error - silently falling back to a default would
// leave the caller talking to the device in a way they didn't ask for.
func TestParseFromOptions_rejectsBadValues(t *testing.T) {
	for _, test := range []struct {
		name              string
		connectionOptions map[string][]string
	}{
		{"unit identifier beyond a byte", map[string][]string{"unit-identifier": {"256"}}},
		{"unit identifier that isn't a number", map[string][]string{"default-unit-identifier": {"nope"}}},
		{"unknown byte order", map[string][]string{"default-payload-byte-order": {"MIDDLE_ENDIAN"}}},
		{"unparsable ping address", map[string][]string{"ping-address": {"this is not an address"}}},
		{"request timeout that isn't a number", map[string][]string{"request-timeout-ms": {"soon"}}},
		{"request timeout of zero", map[string][]string{"request-timeout-ms": {"0"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseFromOptions(zerolog.Nop(), test.connectionOptions)
			assert.Error(t, err)
		})
	}
}

// The request timeout becomes the deadline the codec derives the lifetime of its expectation from.
func TestWithRequestTimeout(t *testing.T) {
	t.Run("applies the configured timeout", func(t *testing.T) {
		ctx, cancel := withRequestTimeout(context.Background(), time.Minute)
		defer cancel()
		deadline, ok := ctx.Deadline()
		require.True(t, ok)
		assert.WithinDuration(t, time.Now().Add(time.Minute), deadline, 5*time.Second)
	})
	t.Run("keeps a deadline the caller brought", func(t *testing.T) {
		callerDeadline := time.Now().Add(time.Hour)
		callerCtx, cancelCaller := context.WithDeadline(context.Background(), callerDeadline)
		defer cancelCaller()

		ctx, cancel := withRequestTimeout(callerCtx, time.Minute)
		defer cancel()
		deadline, ok := ctx.Deadline()
		require.True(t, ok)
		assert.Equal(t, callerDeadline, deadline)
	})
	t.Run("adds no deadline without a timeout", func(t *testing.T) {
		ctx, cancel := withRequestTimeout(context.Background(), 0)
		defer cancel()
		_, ok := ctx.Deadline()
		assert.False(t, ok)
	})
	t.Run("cancelling releases the request", func(t *testing.T) {
		ctx, cancel := withRequestTimeout(context.Background(), time.Minute)
		cancel()
		assert.Error(t, ctx.Err())
	})
}
