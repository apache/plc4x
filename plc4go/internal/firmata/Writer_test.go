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

package firmata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// executeWrite runs a write request and hands back its result, insisting that one arrives.
func executeWrite(t *testing.T, connection *Connection, tagName string, tagAddress string, value any) apiModel.PlcWriteRequestResult {
	t.Helper()
	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress(tagName, tagAddress, value)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)
	select {
	case result := <-writeRequest.Execute(testutils.TestContext(t)):
		return result
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the write never completed the result channel")
		return nil
	}
}

// sentBytes drains everything the connection has put on the wire.
func sentBytes(t *testing.T, transportInstance *test.TransportInstance) []byte {
	t.Helper()
	return transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes())
}

func TestWriter_WriteASingleDigitalPin(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	result := executeWrite(t, connection, "led", "digital:13", true)
	require.NoError(t, result.GetErr())
	require.NotNil(t, result.GetResponse())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("led"))

	// A pin has to be switched to output mode before it can be driven, then the value follows.
	assert.Equal(t, []byte{
		0xF4, 0x0D, 0x01, // set pin mode: pin 13, output
		0xF5, 0x0D, 0x01, // set digital pin value: pin 13, on
	}, sentBytes(t, transportInstance))
}

// The mode of a pin is a property of the board, so it only has to be set once - the second write to
// the same pin is just the value (plc4j FirmataConnection.buildWriteMessages).
func TestWriter_WriteRepeatsNoPinMode(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	require.NoError(t, executeWrite(t, connection, "led", "digital:13", true).GetErr())
	sentBytes(t, transportInstance)

	result := executeWrite(t, connection, "led", "digital:13", false)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("led"))
	assert.Equal(t, []byte{0xF5, 0x0D, 0x00}, sentBytes(t, transportInstance))
}

func TestWriter_WriteARunOfDigitalPins(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	result := executeWrite(t, connection, "bar", "digital:2[0..2]", []bool{true, false, true})
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("bar"))

	assert.Equal(t, []byte{
		0xF4, 0x02, 0x01, 0xF4, 0x03, 0x01, 0xF4, 0x04, 0x01, // all three pins become outputs
		0xF5, 0x02, 0x01, 0xF5, 0x03, 0x00, 0xF5, 0x04, 0x01, // one value per pin, in pin order
	}, sentBytes(t, transportInstance))
}

// A run of pins needs one value per pin. The array info of the tag says so, which lets the value
// handler reject the request before it is ever built - and the writer refuses the same thing again
// for a value which was handed in as a plc value directly.
func TestWriter_WriteRejectsTheWrongNumberOfValues(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("bar", "digital:2[0..2]", []bool{true, false})
	_, err := writeRequestBuilder.Build()
	assert.Error(t, err)
	assert.Empty(t, sentBytes(t, transportInstance), "a rejected write must not reconfigure any pin")

	values, err := boolValues(spiValues.NewPlcList([]apiValues.PlcValue{
		spiValues.NewPlcBOOL(true), spiValues.NewPlcBOOL(false),
	}), 3)
	assert.Error(t, err)
	assert.Nil(t, values)
}

func TestBoolValues(t *testing.T) {
	tests := []struct {
		name     string
		value    apiValues.PlcValue
		quantity uint8
		want     []bool
		wantErr  bool
	}{
		{name: "a single value", value: spiValues.NewPlcBOOL(true), quantity: 1, want: []bool{true}},
		{
			name:     "a list of values",
			value:    spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(true)}),
			quantity: 2,
			want:     []bool{false, true},
		},
		{
			name:     "a list of one for a scalar tag",
			value:    spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcBOOL(true)}),
			quantity: 1,
			want:     []bool{true},
		},
		{name: "no value at all", quantity: 1, wantErr: true},
		{name: "a single value for a run of pins", value: spiValues.NewPlcBOOL(true), quantity: 3, wantErr: true},
		{name: "a value which isn't a bool", value: spiValues.NewPlcSTRING("nonsense"), quantity: 1, wantErr: true},
		{
			name:     "a list holding something which isn't a bool",
			value:    spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcSTRING("nonsense")}),
			quantity: 1,
			wantErr:  true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			values, err := boolValues(testCase.value, testCase.quantity)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, values)
		})
	}
}

// Writing an analog pin needs the extended-analog sysex command, which neither this driver nor
// plc4j's speaks - so it is refused rather than silently dropped.
func TestWriter_WriteRefusesAnalogPins(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	result := executeWrite(t, connection, "dial", "analog:2", int16(5))
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, result.GetResponse().GetResponseCode("dial"))
	assert.Empty(t, sentBytes(t, transportInstance))
}

// A pin can only have one mode, so a pin which is being reported on can't be written to.
func TestWriter_WriteRefusesAPinWhichIsAnInput(t *testing.T) {
	connection, transportInstance := newTestConnection(t)
	_, responseCode := subscribe(t, connection, "button", "digital:13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	sentBytes(t, transportInstance)

	result := executeWrite(t, connection, "led", "digital:13", true)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, result.GetResponse().GetResponseCode("led"))
	assert.Empty(t, sentBytes(t, transportInstance))
}

// A rejected pin in the middle of a run must not leave the pins in front of it reconfigured.
func TestWriter_WriteClaimsARunAllOrNothing(t *testing.T) {
	connection, transportInstance := newTestConnection(t)
	_, responseCode := subscribe(t, connection, "button", "digital:3")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	sentBytes(t, transportInstance)

	result := executeWrite(t, connection, "bar", "digital:2[0..2]", []bool{true, true, true})
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, result.GetResponse().GetResponseCode("bar"))
	assert.Empty(t, sentBytes(t, transportInstance))

	connection.pinMutex.Lock()
	defer connection.pinMutex.Unlock()
	assert.NotContains(t, connection.digitalPins, uint8(2), "pin 2 must not have been claimed")
	assert.Equal(t, readWriteModel.PinMode_PinModeInput, connection.digitalPins[uint8(3)])
	assert.NotContains(t, connection.digitalPins, uint8(4), "pin 4 must not have been claimed")
}

// Every tag of a request is answered on its own: plc4j fails the whole write as soon as one tag
// doesn't work out, which isn't any more atomic and loses the tags that would have gone through.
func TestWriter_WriteAnswersEveryTagOnItsOwn(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("led", "digital:13", true)
	writeRequestBuilder.AddTagAddress("dial", "analog:2", int16(5))
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)
	result := <-writeRequest.Execute(testutils.TestContext(t))

	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("led"))
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, result.GetResponse().GetResponseCode("dial"))
	assert.Equal(t, []byte{0xF4, 0x0D, 0x01, 0xF5, 0x0D, 0x01}, sentBytes(t, transportInstance))
}

// A write whose messages never made it onto the wire must not leave the pins recorded as outputs:
// the set-pin-mode messages are what configures them, so a retry has to send them again. Without the
// rollback the retry would find the pins claimed, send only the value and report a success for a pin
// the board never put into output mode.
func TestWriter_AFailedSendGivesThePinClaimBack(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())
	codec.failSends()

	result := executeWrite(t, connection, "led", "digital:9", true)
	require.NoError(t, result.GetErr())
	require.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("led"))
	func() {
		connection.pinMutex.Lock()
		defer connection.pinMutex.Unlock()
		assert.NotContains(t, connection.digitalPins, uint8(9), "the claim of a write that failed has to be given back")
	}()

	codec.allowSends()
	result = executeWrite(t, connection, "led", "digital:9", true)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("led"))
	assert.Equal(t, []byte{
		0xF4, 0x09, 0x01, // set pin mode: pin 9, output - sent again by the retry
		0xF5, 0x09, 0x01, // set digital pin value: pin 9, high
	}, codec.sentBytesOf(t))
}

// The interesting failure is the partial one: the set-pin-mode makes it onto the wire and the value
// behind it doesn't. The pin is not usably configured either way - nothing was written to it - so
// the claim still has to go back, and the retry re-sends both messages.
func TestWriter_APartialSendGivesThePinClaimBack(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())
	codec.failSendsAfter(1)

	result := executeWrite(t, connection, "led", "digital:9", true)
	require.NoError(t, result.GetErr())
	require.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("led"))
	func() {
		connection.pinMutex.Lock()
		defer connection.pinMutex.Unlock()
		assert.NotContains(t, connection.digitalPins, uint8(9), "a claim whose batch died part-way has to be given back")
	}()

	codec.allowSends()
	result = executeWrite(t, connection, "led", "digital:9", true)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("led"))
	assert.Equal(t, []byte{
		0xF4, 0x09, 0x01, // the set pin mode which did get out the first time
		0xF4, 0x09, 0x01, // and both messages again from the retry
		0xF5, 0x09, 0x01,
	}, codec.sentBytesOf(t))
}
