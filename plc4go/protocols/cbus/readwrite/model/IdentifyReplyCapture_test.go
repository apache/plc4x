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

package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIdentifyReplySummaryCapture parses a real C-Bus Ethernet interface's (Type "PC_CNIED") reply
// to an identify of the Summary attribute (request "2103", reply captured on the wire) and
// serializes it back. Per the Serial Interface User Guide (7.2) the reply carries 9 data bytes: the
// part name (6 bytes), the unit service type (1 byte) and the version as 2 hex bytes; this unit's
// version 0x5400 matches its FirmwareVersion reply "5.4.00". Its part name bytes are not ASCII.
func TestIdentifyReplySummaryCapture(t *testing.T) {
	const raw = "8A03B64DB4B68CDEF1540057\r\n"
	toClient := parseCapturedIdentifyReply(t, raw)
	replyOrConfirmation, ok := toClient.GetReply().(ReplyOrConfirmationReply)
	require.True(t, ok, "%T", toClient.GetReply())
	encodedReply, ok := replyOrConfirmation.GetReply().(ReplyEncodedReply)
	require.True(t, ok, "%T", replyOrConfirmation.GetReply())
	calReply, ok := encodedReply.GetEncodedReply().(EncodedReplyCALReply)
	require.True(t, ok, "%T", encodedReply.GetEncodedReply())
	identifyReply, ok := calReply.GetCalReply().GetCalData().(CALDataIdentifyReply)
	require.True(t, ok, "%T", calReply.GetCalReply().GetCalData())
	assert.Equal(t, Attribute_Summary, identifyReply.GetAttribute())
	summary, ok := identifyReply.GetIdentifyReplyCommand().(IdentifyReplyCommandSummary)
	require.True(t, ok, "%T", identifyReply.GetIdentifyReplyCommand())

	assert.EqualValues(t, []byte{0xB6, 0x4D, 0xB4, 0xB6, 0x8C, 0xDE}, summary.GetPartName())
	assert.Equal(t, byte(0xF1), summary.GetUnitServiceType())
	assert.EqualValues(t, 0x5400, summary.GetVersion())

	serialized, err := toClient.Serialize()
	require.NoError(t, err)
	assert.Equal(t, raw, string(serialized))
}

// TestIdentifyReplyNetworkVoltageCapture parses the same interface's reply to an identify of
// NetworkVoltage (request "2107"): per the Serial Interface User Guide (7.2) 2 ASCII bytes of whole
// volts, the decimal point ($2E), 1 ASCII byte of the first decimal place and 'V' ($56), here "29.8V".
// The capture holds the short form "860732392E38564C", taken in BASIC mode; its header $86 (a REPLY of 6
// bytes) is also the marker of a long form CAL reply, which the parser assumes. This test wraps the
// captured CAL data in the long form a unit sends in SMART mode with IDMON (plc4go's default), from unit
// $04 through serial interface $99 (the addressing of the guide's example in 9.2).
func TestIdentifyReplyNetworkVoltageCapture(t *testing.T) {
	const raw = "86049900860732392E385629\r\n"
	identifyReply := parseCapturedIdentifyReplyCALData(t, raw)
	assert.Equal(t, Attribute_NetworkVoltage, identifyReply.GetAttribute())
	networkVoltage, ok := identifyReply.GetIdentifyReplyCommand().(IdentifyReplyCommandNetworkVoltage)
	require.True(t, ok, "%T", identifyReply.GetIdentifyReplyCommand())
	assert.Equal(t, "29", networkVoltage.GetVolts())
	assert.Equal(t, "8", networkVoltage.GetVoltsDecimalPlace())
}

// TestIdentifyReplyCaptureRoundTrip parses the same interface's replies to the identify of
// Manufacturer, Type, FirmwareVersion, Summary, ExtendedDiagnosticSummary and NetworkVoltage and
// serializes them back.
func TestIdentifyReplyCaptureRoundTrip(t *testing.T) {
	for _, raw := range []string{
		"8900434C495053414C204F\r\n",         // Manufacturer "CLIPSAL "
		"890150435F434E49454421\r\n",         // Type "PC_CNIED"
		"8902352E342E3030202010\r\n",         // FirmwareVersion "5.4.00  "
		"8A03B64DB4B68CDEF1540057\r\n",       // Summary
		"8D04FFFFFF00001897C095BE0005AB\r\n", // ExtendedDiagnosticSummary
		"86049900860732392E385629\r\n",       // NetworkVoltage "29.8V", long form (see TestIdentifyReplyNetworkVoltageCapture)
	} {
		t.Run(raw[:4], func(t *testing.T) {
			toClient := parseCapturedIdentifyReply(t, raw)
			serialized, err := toClient.Serialize()
			require.NoError(t, err)
			assert.Equal(t, raw, string(serialized))
		})
	}
}

// parseCapturedIdentifyReplyCALData parses a captured reply to an identify request down to its CAL data.
func parseCapturedIdentifyReplyCALData(t *testing.T, raw string) CALDataIdentifyReply {
	t.Helper()
	toClient := parseCapturedIdentifyReply(t, raw)
	replyOrConfirmation, ok := toClient.GetReply().(ReplyOrConfirmationReply)
	require.True(t, ok, "%T", toClient.GetReply())
	encodedReply, ok := replyOrConfirmation.GetReply().(ReplyEncodedReply)
	require.True(t, ok, "%T", replyOrConfirmation.GetReply())
	calReply, ok := encodedReply.GetEncodedReply().(EncodedReplyCALReply)
	require.True(t, ok, "%T", encodedReply.GetEncodedReply())
	identifyReply, ok := calReply.GetCalReply().GetCalData().(CALDataIdentifyReply)
	require.True(t, ok, "%T", calReply.GetCalReply().GetCalData())
	return identifyReply
}

// parseCapturedIdentifyReply parses a reply to an identify request as the interface sent it (checksum on).
func parseCapturedIdentifyReply(t *testing.T, raw string) CBusMessageToClient {
	t.Helper()
	msg, err := CBusMessageParse[CBusMessage](
		context.Background(),
		[]byte(raw),
		true,
		NewRequestContext(true),
		NewCBusOptions(false, false, false, false, false, false, false, false, true),
	)
	require.NoError(t, err)
	toClient, ok := msg.(CBusMessageToClient)
	require.True(t, ok, "%T", msg)
	return toClient
}
