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

package eip

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFromOptionsDefaults(t *testing.T) {
	configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{})
	require.NoError(t, err)
	assert.Equal(t, int8(1), configuration.backplane)
	assert.Equal(t, int8(0), configuration.slot)
	assert.True(t, configuration.bigEndian)
	assert.False(t, configuration.forceUnconnectedOperation)
	assert.Equal(t, "", configuration.communicationPath)
	assert.Equal(t, uint16(0), configuration.connectionSerialNumber)
}

func TestParseFromOptionsCamelCase(t *testing.T) {
	configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
		"bigEndian":                 {"false"},
		"forceUnconnectedOperation": {"true"},
		"communicationPath":         {"1,4,2,192.168.0.1,1,1"},
		"connectionSerialNumber":    {"4660"},
	})
	require.NoError(t, err)
	assert.False(t, configuration.bigEndian)
	assert.True(t, configuration.forceUnconnectedOperation)
	assert.Equal(t, "1,4,2,192.168.0.1,1,1", configuration.communicationPath)
	assert.Equal(t, uint16(4660), configuration.connectionSerialNumber)
}

func TestParseFromOptionsKebabCase(t *testing.T) {
	configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
		"big-endian":                  {"false"},
		"force-unconnected-operation": {"true"},
		"connection-serial-number":    {"4660"},
	})
	require.NoError(t, err)
	assert.False(t, configuration.bigEndian)
	assert.True(t, configuration.forceUnconnectedOperation)
	assert.Equal(t, uint16(4660), configuration.connectionSerialNumber)
}

func TestParseFromOptionsAliasPrecedence(t *testing.T) {
	// When both spellings of an aliased option are supplied with conflicting
	// values, getFromOptionsAliases resolves keys in the order given - the
	// camelCase spelling is listed first, so it wins.
	configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
		"bigEndian":  {"false"},
		"big-endian": {"true"},
	})
	require.NoError(t, err)
	assert.False(t, configuration.bigEndian)
}

func TestParseFromOptionsBadBool(t *testing.T) {
	_, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
		"bigEndian": {"maybe"},
	})
	require.Error(t, err)
}
