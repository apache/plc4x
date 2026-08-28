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
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func readerFor(options map[string][]string) (*OptionReader, *bytes.Buffer) {
	var logged bytes.Buffer
	return NewOptionReader(zerolog.New(&logged), options), &logged
}

func TestOptionReader_GetReturnsTheValue(t *testing.T) {
	reader, _ := readerFor(map[string][]string{"unit-identifier": {"3"}})
	assert.Equal(t, "3", reader.Get("unit-identifier"))
}

func TestOptionReader_GetReturnsEmptyForAnAbsentOption(t *testing.T) {
	reader, _ := readerFor(map[string][]string{})
	assert.Equal(t, "", reader.Get("unit-identifier"))
}

// The behaviour of the twelve getFromOptions copies this replaces.
func TestOptionReader_GetWarnsWhenAnOptionIsRepeated(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"unit-identifier": {"3", "4"}})
	assert.Equal(t, "3", reader.Get("unit-identifier"), "the first value wins")
	assert.Contains(t, logged.String(), "Option must be unique")
}

// A key the driver read is recognised by definition - that is the whole point of recording
// consumption rather than declaring a list of names.
func TestOptionReader_SaysNothingAboutAnOptionThatWasRead(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"unit-identifier": {"3"}})
	reader.Get("unit-identifier")
	reader.ReportUnknown("modbus-tcp")
	assert.Empty(t, logged.String())
}

// Asking for an option that was not supplied still counts as recognising it: a driver that reads
// "request-timeout-ms" recognises the name whether or not this connection string carried it.
func TestOptionReader_AskingForAnAbsentOptionStillRecognisesIt(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"request-timeout-ms": {"5000"}})
	reader.Get("request-timeout-ms")
	reader.Get("unit-identifier")
	reader.ReportUnknown("modbus-tcp")
	assert.Empty(t, logged.String())
}

func TestOptionReader_ReportsAnOptionNothingRead(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"unit-identifier": {"3"}, "nonsense": {"1"}})
	reader.Get("unit-identifier")
	reader.ReportUnknown("modbus-tcp")

	assert.Contains(t, logged.String(), "nonsense")
	assert.Contains(t, logged.String(), "modbus-tcp", "the report names the driver")
	assert.NotContains(t, logged.String(), "unit-identifier")
}

// The overwhelmingly common cause is a name that is nearly right, so leading with the replacement
// turns a puzzling no-op into a one-line fix.
func TestOptionReader_SuggestsTheOptionAMisspellingWasMeantToBe(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"unit-identifer": {"3"}})
	reader.Get("unit-identifier")
	reader.ReportUnknown("modbus-tcp")

	assert.Contains(t, logged.String(), "unit-identifer")
	assert.Contains(t, logged.String(), "didYouMean")
	assert.Contains(t, logged.String(), "unit-identifier")
}

func TestOptionReader_SuggestsNothingWhenNothingIsClose(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"completely-different": {"1"}})
	reader.Get("unit-identifier")
	reader.ReportUnknown("modbus-tcp")

	assert.Contains(t, logged.String(), "completely-different")
	assert.NotContains(t, logged.String(), "didYouMean")
}

// A transport option belongs to another consumer. Reporting it would warn about every connection
// string that sets a timeout, which would teach operators to ignore the warning.
//
// The names are registered here rather than by importing a transport: the transports import this
// package, so a test in it cannot import them back. Each transport registers its own in an init(),
// beside the code that reads them.
func TestOptionReader_SaysNothingAboutTransportOptions(t *testing.T) {
	RegisterTransportOptions("connect-timeout-ms", "read-timeout-ms", "reuse-port")

	reader, logged := readerFor(map[string][]string{
		"connect-timeout-ms": {"5000"}, "read-timeout-ms": {"1000"}, "reuse-port": {"true"},
	})
	reader.ReportUnknown("modbus-tcp")
	assert.Empty(t, logged.String())
}

// The exemption covers what some transport registered, and nothing else. An option no transport
// reads is the driver's to report - which is what a hand-kept central list of every transport's
// options could not say.
func TestOptionReader_ReportsAnOptionNoTransportRegistered(t *testing.T) {
	RegisterTransportOptions("baud-rate")

	reader, logged := readerFor(map[string][]string{
		"baud-rate": {"9600"}, "baud-rat": {"9600"},
	})
	reader.ReportUnknown("modbus-tcp")

	assert.NotContains(t, logged.String(), `"option":"baud-rate"`)
	assert.Contains(t, logged.String(), `"option":"baud-rat"`)
}

// What a driver parses itself - a nested or prefixed group it handles by hand - is its to claim.
func TestOptionReader_SaysNothingAboutOptionsTheDriverClaims(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"browser.depth": {"2"}})
	reader.Ignore("browser.depth")
	reader.ReportUnknown("ads")
	assert.Empty(t, logged.String())
}

func TestOptionReader_SaysNothingAboutAnEmptyOptionMap(t *testing.T) {
	reader, logged := readerFor(map[string][]string{})
	reader.ReportUnknown("modbus-tcp")
	assert.Empty(t, logged.String())
}

// One line per unknown option, in a stable order, so a connection string with several typos
// produces a report that reads the same way twice.
func TestOptionReader_ReportsDeterministically(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"zebra": {"1"}, "alpha": {"2"}})
	reader.ReportUnknown("modbus-tcp")

	output := logged.String()
	assert.Less(t, bytes.Index([]byte(output), []byte("alpha")), bytes.Index([]byte(output), []byte("zebra")))
}

// bacnet-ip has always matched option names without regard to case; consolidating the twelve
// per-driver lookups must not quietly change that, in either direction.
func TestOptionReader_CaseInsensitiveMatchesRegardlessOfCase(t *testing.T) {
	reader, logged := readerFor(map[string][]string{"localdeviceid": {"42"}})
	reader.CaseInsensitive()

	assert.Equal(t, "42", reader.Get("LocalDeviceId"))
	reader.ReportUnknown("bacnet-ip")
	assert.Empty(t, logged.String(), "the supplied spelling counts as consumed")
}

// The other ten drivers match exactly, as they always did.
func TestOptionReader_MatchesExactlyByDefault(t *testing.T) {
	reader, _ := readerFor(map[string][]string{"localdeviceid": {"42"}})
	assert.Equal(t, "", reader.Get("LocalDeviceId"))
}
