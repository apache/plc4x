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

package opcua

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An option this driver does not know is reported and ignored, rather than refusing the
// connection.
//
// This driver used to be the only one in PLC4X that failed on an unknown option, so one
// connection string was accepted by plc4j and by every other Go driver, and rejected here. The
// reason it failed is still real - a typo in a security-relevant option falls back to a default -
// which is why the option has to be named in the log rather than passed over in silence.
func TestParseFromOptions_reportsAnUnknownOptionAndCarriesOn(t *testing.T) {
	var logged bytes.Buffer
	log := zerolog.New(&logged)

	configuration, err := ParseFromOptions(log, map[string][]string{
		"securityPolicy": {"None"},
		"scurityPolicy":  {"Basic256"}, // the typo that used to refuse the connection
	})

	require.NoError(t, err, "an unknown option must not refuse the connection")
	assert.Equal(t, "None", configuration.SecurityPolicy, "the options it does know still apply")
	assert.Contains(t, logged.String(), "scurityPolicy", "the ignored option must be named")
}

// A transport-level option belongs to another consumer and is not this driver's to report.
func TestParseFromOptions_saysNothingAboutTransportOptions(t *testing.T) {
	var logged bytes.Buffer
	log := zerolog.New(&logged)

	_, err := ParseFromOptions(log, map[string][]string{"connect-timeout-ms": {"5000"}})

	require.NoError(t, err)
	assert.NotContains(t, logged.String(), "connect-timeout-ms")
}

// A password must never appear in a rendering. The opcua Configuration and SecureChannel render
// themselves through generated code, which used to write the password verbatim - so turning on
// debug logging wrote the PLC password into the log.
//
// The value is planted and then looked for, rather than the field list being asserted: that is
// what makes this hold for a secret added later. A test that checked "password renders as
// <redacted>" would pass while a newly added token leaked.
func TestConfiguration_NoSecretIsRendered(t *testing.T) {
	const sentinel = "hunter2-sentinel-value"

	configuration, err := ParseFromOptions(zerolog.Nop(), map[string][]string{
		"username":         {"operator"},
		"password":         {sentinel},
		"keyStorePassword": {sentinel},
	})
	require.NoError(t, err)
	require.Equal(t, sentinel, configuration.Password, "the value is still available to the driver")

	rendered := configuration.String()
	assert.NotContains(t, rendered, sentinel, "no secret may appear in a rendering")
	assert.Contains(t, rendered, "<redacted>", "and the reader is told a secret is configured")
	assert.Contains(t, rendered, "operator", "the username is not a secret - it says who is connecting")
}
