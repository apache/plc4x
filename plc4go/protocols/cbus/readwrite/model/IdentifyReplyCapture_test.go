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

	"github.com/apache/plc4x/plc4go/spi/utils"
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
// Manufacturer, Type, FirmwareVersion, Summary, ExtendedDiagnosticSummary, NetworkVoltage and
// OutputUnitSummary and serializes them back.
func TestIdentifyReplyCaptureRoundTrip(t *testing.T) {
	for _, raw := range []string{
		"8900434C495053414C204F\r\n",         // Manufacturer "CLIPSAL "
		"890150435F434E49454421\r\n",         // Type "PC_CNIED"
		"8902352E342E3030202010\r\n",         // FirmwareVersion "5.4.00  "
		"8A03B64DB4B68CDEF1540057\r\n",       // Summary
		"8D04FFFFFF00001897C095BE0005AB\r\n", // ExtendedDiagnosticSummary
		"86049900860732392E385629\r\n",       // NetworkVoltage "29.8V", long form (see TestIdentifyReplyNetworkVoltageCapture)
		"8510020000FF6A\r\n",                 // OutputUnitSummary
	} {
		t.Run(raw[:4], func(t *testing.T) {
			toClient := parseCapturedIdentifyReply(t, raw)
			serialized, err := toClient.Serialize()
			require.NoError(t, err)
			assert.Equal(t, raw, string(serialized))
		})
	}
}

// TestIdentifyReplyEmptyCapture parses replies without data, as units send them for an attribute they
// have nothing for (a REPLY of 1 byte, the attribute alone), and serializes them back. Delays ("810C73")
// and DSIStatus ("81116E") are the same interface's captured replies; the others are constructed the same
// way for the remaining attributes whose layout needs data, and one in the long form.
func TestIdentifyReplyEmptyCapture(t *testing.T) {
	for _, tt := range []struct {
		raw       string
		attribute Attribute
	}{
		{"81007F\r\n", Attribute_Manufacturer},
		{"81017E\r\n", Attribute_Type},
		{"81027D\r\n", Attribute_FirmwareVersion},
		{"81037C\r\n", Attribute_Summary},
		{"81047B\r\n", Attribute_ExtendedDiagnosticSummary},
		{"810778\r\n", Attribute_NetworkVoltage},
		{"810C73\r\n", Attribute_Delays}, // captured
		{"81106F\r\n", Attribute_OutputUnitSummary},
		{"81116E\r\n", Attribute_DSIStatus}, // captured
		// the long form a unit sends in SMART mode with IDMON: from unit $04 through serial interface $99
		{"86049900810C50\r\n", Attribute_Delays},
	} {
		t.Run(tt.attribute.String(), func(t *testing.T) {
			identifyReply := parseCapturedIdentifyReplyCALData(t, tt.raw)
			assert.Equal(t, tt.attribute, identifyReply.GetAttribute())
			assert.IsType(t, NewIdentifyReplyCommandEmpty(), identifyReply.GetIdentifyReplyCommand())
			serialized, err := parseCapturedIdentifyReply(t, tt.raw).Serialize()
			require.NoError(t, err)
			assert.Equal(t, tt.raw, string(serialized))
		})
	}
}

// TestIdentifyReplyDelaysRestrikeOnly parses a Delays reply with one data byte (constructed): no terminal
// delays and the re-strike delay. One byte is the least a Delays reply with data carries, so it must
// not parse as IdentifyReplyCommandEmpty.
func TestIdentifyReplyDelaysRestrikeOnly(t *testing.T) {
	const raw = "820C056D\r\n"
	identifyReply := parseCapturedIdentifyReplyCALData(t, raw)
	delays, ok := identifyReply.GetIdentifyReplyCommand().(IdentifyReplyCommandDelays)
	require.True(t, ok, "%T", identifyReply.GetIdentifyReplyCommand())
	assert.Empty(t, delays.GetTerminalLevels())
	assert.Equal(t, byte(0x05), delays.GetReStrikeDelay())
	serialized, err := parseCapturedIdentifyReply(t, raw).Serialize()
	require.NoError(t, err)
	assert.Equal(t, raw, string(serialized))
}

// TestIdentifyReplyUnknownAttribute does not parse the same interface's reply without data to an identify
// of an attribute the Attribute enum does not know (0x12, captured "81126D") as an identify reply: Go would
// read the unknown attribute as 0 (Manufacturer) and serialize it as such.
func TestIdentifyReplyUnknownAttribute(t *testing.T) {
	_, err := CBusMessageParse[CBusMessage](
		context.Background(),
		[]byte("81126D\r\n"),
		true,
		NewRequestContext(true),
		NewCBusOptions(false, false, false, false, false, false, false, false, true),
	)
	assert.ErrorContains(t, err, "unknown identify attribute")
	// recoverable: inside a confirmation the optional reply then counts as absent
	assert.ErrorIs(t, err, utils.ParseAssertError{})
}

// TestIdentifyReplyEmptyArrayCapture parses the same interface's replies without data to the identify of
// attributes whose reply length is set by the number of the unit's variables (Serial Interface User Guide
// 7.2, "dynamic"): they stay their own type with an empty array, and serialize back.
func TestIdentifyReplyEmptyArrayCapture(t *testing.T) {
	for _, tt := range []struct {
		raw     string
		command IdentifyReplyCommand
	}{
		{"81057A\r\n", NewIdentifyReplyCommandNetworkTerminalLevels(nil)},
		{"810679\r\n", NewIdentifyReplyCommandTerminalLevels(nil)},
		{"810877\r\n", NewIdentifyReplyCommandGAVValuesCurrent(nil)},
		{"810976\r\n", NewIdentifyReplyCommandGAVValuesStored(nil)},
		{"810A75\r\n", NewIdentifyReplyCommandGAVPhysicalAddresses(nil)},
		{"810B74\r\n", NewIdentifyReplyCommandLogicalAssignment(nil)},
		{"810D72\r\n", NewIdentifyReplyCommandMinimumLevels(nil)},
		{"810E71\r\n", NewIdentifyReplyCommandMaximumLevels(nil)},
		{"810F70\r\n", NewIdentifyReplyCommandCurrentSenseLevels(nil)},
	} {
		t.Run(tt.raw[:4], func(t *testing.T) {
			identifyReply := parseCapturedIdentifyReplyCALData(t, tt.raw)
			assert.IsType(t, tt.command, identifyReply.GetIdentifyReplyCommand())
			assert.Zero(t, identifyReply.GetIdentifyReplyCommand().GetLengthInBytes(context.Background()))
			serialized, err := parseCapturedIdentifyReply(t, tt.raw).Serialize()
			require.NoError(t, err)
			assert.Equal(t, tt.raw, string(serialized))
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

// TestIdentifyReplyOutputUnitSummaryLengths parses OutputUnitSummary replies of each length the Serial
// Interface User Guide (7.2) allows and serializes them back. The reply is the unit flags, the GAV store
// enable bytes and the time from the last recovery of mains; its length is dynamic (set by the unit's
// variables), the guide shows the longest, 4 bytes. The 4 byte replies are captured (a C-Bus Ethernet
// interface in BASIC mode, and unit $0D through serial interface $02 in SMART mode); the 3 and 2 byte
// replies are the first capture with one or both GAV store enable bytes left out.
func TestIdentifyReplyOutputUnitSummaryLengths(t *testing.T) {
	gavStoreEnabled := func(b byte) *byte { return &b }
	for _, tt := range []struct {
		name                   string
		raw                    string
		unitFlags              byte
		gavStoreEnabledByte1   *byte
		gavStoreEnabledByte2   *byte
		timeFromLastRecoverSec uint8
	}{
		{"4 bytes, captured", "8510020000FF6A\r\n", 0x02, gavStoreEnabled(0x00), gavStoreEnabled(0x00), 0xFF},
		{"4 bytes, captured long form", "860D0200851008FFFFFFD1\r\n", 0x08, gavStoreEnabled(0xFF), gavStoreEnabled(0xFF), 0xFF},
		{"3 bytes", "84100200FF6B\r\n", 0x02, gavStoreEnabled(0x00), nil, 0xFF},
		{"2 bytes", "831002FF6C\r\n", 0x02, nil, nil, 0xFF},
	} {
		t.Run(tt.name, func(t *testing.T) {
			identifyReply := parseCapturedIdentifyReplyCALData(t, tt.raw)
			assert.Equal(t, Attribute_OutputUnitSummary, identifyReply.GetAttribute())
			summary, ok := identifyReply.GetIdentifyReplyCommand().(IdentifyReplyCommandOutputUnitSummary)
			require.True(t, ok, "%T", identifyReply.GetIdentifyReplyCommand())
			unitFlags, err := summary.GetUnitFlags().Serialize()
			require.NoError(t, err)
			assert.Equal(t, []byte{tt.unitFlags}, unitFlags)
			assert.Equal(t, tt.gavStoreEnabledByte1, summary.GetGavStoreEnabledByte1())
			assert.Equal(t, tt.gavStoreEnabledByte2, summary.GetGavStoreEnabledByte2())
			assert.Equal(t, tt.timeFromLastRecoverSec, summary.GetTimeFromLastRecoverOfMainsInSeconds())
			serialized, err := parseCapturedIdentifyReply(t, tt.raw).Serialize()
			require.NoError(t, err)
			assert.Equal(t, tt.raw, string(serialized))
		})
	}
}
