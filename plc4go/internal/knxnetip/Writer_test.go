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
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const testCommunicationChannelId = 42

// writerCodec is a spi.MessageCodec test double which lets a test decide how the gateway
// reacts to an outgoing tunneling-request. Only the methods the write-path uses are
// implemented, anything else would panic on the (nil) embedded interface, which is exactly
// what a test double should do.
type writerCodec struct {
	spi.MessageCodec

	mutex        sync.Mutex
	sentRequests []spi.Message
	expectations []writerCodecExpectation

	// respond is invoked for every SendRequest. Returning an error makes the codec report a
	// send-failure, otherwise the test drives the (optional) replies itself.
	respond func(codec *writerCodec, message spi.Message, handleMessage spi.HandleMessage, handleError spi.HandleError) error
}

type writerCodecExpectation struct {
	accepts       spi.AcceptsMessage
	handleMessage spi.HandleMessage
	handleError   spi.HandleError
}

func (c *writerCodec) Expect(_ context.Context, _ string, accepts spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.expectations = append(c.expectations, writerCodecExpectation{accepts, handleMessage, handleError})
}

func (c *writerCodec) SendRequest(_ context.Context, _ string, message spi.Message, accepts spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	c.mutex.Lock()
	c.sentRequests = append(c.sentRequests, message)
	respond := c.respond
	c.mutex.Unlock()
	if respond == nil {
		return nil
	}
	return respond(c, message, func(message spi.Message) error {
		if !accepts(message) {
			return errors.New("the request expectation didn't accept the message")
		}
		return handleMessage(message)
	}, handleError)
}

// deliver hands the given message to the first registered expectation which accepts it,
// the way the real codec dispatches an incoming message.
func (c *writerCodec) deliver(message spi.Message) bool {
	c.mutex.Lock()
	expectations := slices.Clone(c.expectations)
	c.mutex.Unlock()
	for _, expectation := range expectations {
		if expectation.accepts(message) {
			_ = expectation.handleMessage(message)
			return true
		}
	}
	return false
}

func (c *writerCodec) getSentRequests() []spi.Message {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return slices.Clone(c.sentRequests)
}

// newWriterConnection builds a connection which is wired up far enough for the write-path:
// it has a client address, a communication-channel-id and a codec answering as instructed.
func newWriterConnection(t *testing.T, connectionOptions map[string][]string, respond func(codec *writerCodec, message spi.Message, handleMessage spi.HandleMessage, handleError spi.HandleError) error) (*Connection, *writerCodec) {
	t.Helper()
	codec := &writerCodec{respond: respond}
	connection := &Connection{
		messageCodec: codec,
		options:      connectionOptions,
		tagHandler:   NewTagHandler(),
		valueHandler: NewValueHandler(),
		metadata:     &ConnectionMetadata{},
		valueCache:   map[uint16][]byte{},
		log:          testutils.ProduceTestingLogger(t),
	}
	connection.CommunicationChannelId = testCommunicationChannelId
	connection.ClientKnxAddress = driverModel.NewKnxAddress(1, 1, 1)
	connection.SequenceCounter = -1
	t.Cleanup(connection.wg.Wait)
	return connection, codec
}

// confirmingGateway acts like a real gateway: it acknowledges the tunneling-request and
// then echoes the frame back as an LDataCon.
func confirmingGateway(errorFlag bool) func(*writerCodec, spi.Message, spi.HandleMessage, spi.HandleError) error {
	return func(codec *writerCodec, message spi.Message, handleMessage spi.HandleMessage, _ spi.HandleError) error {
		tunnelingRequest := message.(driverModel.TunnelingRequest)
		sequenceCounter := tunnelingRequest.GetTunnelingRequestDataBlock().GetSequenceCounter()
		if err := handleMessage(driverModel.NewTunnelingResponse(
			driverModel.NewTunnelingResponseDataBlock(testCommunicationChannelId, sequenceCounter, driverModel.Status_NO_ERROR))); err != nil {
			return err
		}
		frame := tunnelingRequest.GetCemi().(driverModel.LDataReq).GetDataFrame().(driverModel.LDataExtended)
		codec.deliver(driverModel.NewTunnelingRequest(
			driverModel.NewTunnelingRequestDataBlock(testCommunicationChannelId, sequenceCounter),
			driverModel.NewLDataCon(0, nil, driverModel.NewLDataExtended(
				frame.GetFrameType(), frame.GetNotRepeated(), frame.GetPriority(),
				frame.GetAcknowledgeRequested(), errorFlag, frame.GetGroupAddress(),
				frame.GetHopCount(), frame.GetExtendedFrameFormat(),
				frame.GetSourceAddress(), frame.GetDestinationAddress(), frame.GetApdu()))))
		return nil
	}
}

// executeWrite runs a write in the background so a regression shows up as a failing test
// instead of a hanging test-run.
func executeWrite(t *testing.T, ctx context.Context, writeRequest apiModel.PlcWriteRequest) apiModel.PlcWriteRequestResult {
	t.Helper()
	results := writeRequest.Execute(ctx)
	select {
	case result := <-results:
		require.NotNil(t, result, "the write result must never be nil")
		return result
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the write never completed the result channel")
		return nil
	}
}

func Test_Writer_Write_success(t *testing.T) {
	connection, codec := newWriterConnection(t, nil, confirmingGateway(false))

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	require.NotNil(t, result.GetResponse())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("switch"))

	// Verify the frame we put on the wire really is a GroupValueWrite for 1/2/3 carrying
	// the value in the 6 bit small-payload.
	sentRequests := codec.getSentRequests()
	require.Len(t, sentRequests, 1)
	tunnelingRequest := sentRequests[0].(driverModel.TunnelingRequest)
	assert.Equal(t, uint8(testCommunicationChannelId), tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId())
	frame := tunnelingRequest.GetCemi().(driverModel.LDataReq).GetDataFrame().(driverModel.LDataExtended)
	assert.True(t, frame.GetGroupAddress())
	assert.Equal(t, connection.ClientKnxAddress, frame.GetSourceAddress())
	assert.Equal(t, []byte{0x0A, 0x03}, frame.GetDestinationAddress())
	groupValueWrite := frame.GetApdu().(driverModel.ApduDataContainer).GetDataApdu().(driverModel.ApduDataGroupValueWrite)
	assert.Equal(t, int8(1), groupValueWrite.GetDataFirstByte())
	assert.Empty(t, groupValueWrite.GetData())
}

// Test_Writer_Write_largeDatapointType makes sure datapoint-types which don't fit into the
// 6 bit small-payload are sent as trailing data bytes instead.
func Test_Writer_Write_largeDatapointType(t *testing.T) {
	connection, codec := newWriterConnection(t, nil, confirmingGateway(false))

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("brightness", "1/2/4:DPT_Scaling", uint8(0x80))
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("brightness"))

	sentRequests := codec.getSentRequests()
	require.Len(t, sentRequests, 1)
	frame := sentRequests[0].(driverModel.TunnelingRequest).GetCemi().(driverModel.LDataReq).GetDataFrame().(driverModel.LDataExtended)
	groupValueWrite := frame.GetApdu().(driverModel.ApduDataContainer).GetDataApdu().(driverModel.ApduDataGroupValueWrite)
	assert.Equal(t, int8(0), groupValueWrite.GetDataFirstByte())
	assert.Equal(t, []byte{0x80}, groupValueWrite.GetData())
}

// Test_Writer_Write_timeout is the regression test for the write which never completed its
// result channel: a silent gateway has to end up as a REQUEST_TIMEOUT response code.
func Test_Writer_Write_timeout(t *testing.T) {
	connection, codec := newWriterConnection(t, map[string][]string{"request-timeout": {"100"}},
		func(*writerCodec, spi.Message, spi.HandleMessage, spi.HandleError) error {
			// The gateway swallows the request without ever answering.
			return nil
		})

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	require.NotNil(t, result.GetResponse())
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("switch"))
	assert.Len(t, codec.getSentRequests(), 1)
}

// Test_Writer_Write_canceledContext drives the same "always complete" guarantee through a
// context which is already done.
func Test_Writer_Write_canceledContext(t *testing.T) {
	connection, _ := newWriterConnection(t, nil,
		func(*writerCodec, spi.Message, spi.HandleMessage, spi.HandleError) error {
			return nil
		})

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	result := executeWrite(t, ctx, writeRequest)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("switch"))
}

// Test_Writer_Write_multipleTags makes sure a multi-tag request doesn't hang either and
// that every tag gets its own frame and its own response code.
func Test_Writer_Write_multipleTags(t *testing.T) {
	connection, codec := newWriterConnection(t, nil, confirmingGateway(false))

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("first", "1/2/3:DPT_Switch", true)
	writeRequestBuilder.AddTagAddress("second", "1/2/4:DPT_Switch", false)
	// A pattern addresses more than one device, which a write must refuse.
	writeRequestBuilder.AddTagAddress("pattern", "1/2/*:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	require.NotNil(t, result.GetResponse())
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("first"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("second"))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, result.GetResponse().GetResponseCode("pattern"))
	// Only the two resolvable tags made it onto the bus.
	assert.Len(t, codec.getSentRequests(), 2)
}

// Test_Writer_Write_gatewayReportsError makes sure a negative confirmation isn't reported
// as a successful write.
func Test_Writer_Write_gatewayReportsError(t *testing.T) {
	connection, _ := newWriterConnection(t, nil, confirmingGateway(true))

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("switch"))
}

// Test_Writer_Write_negativeAck covers a gateway which refuses the tunneling-request itself.
func Test_Writer_Write_negativeAck(t *testing.T) {
	connection, _ := newWriterConnection(t, nil,
		func(_ *writerCodec, message spi.Message, handleMessage spi.HandleMessage, _ spi.HandleError) error {
			sequenceCounter := message.(driverModel.TunnelingRequest).GetTunnelingRequestDataBlock().GetSequenceCounter()
			return handleMessage(driverModel.NewTunnelingResponse(driverModel.NewTunnelingResponseDataBlock(
				testCommunicationChannelId, sequenceCounter, driverModel.Status_INVALID_CONNECTION_ID)))
		})

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("switch"))
}

// Test_Writer_Write_sendFails covers the codec refusing to send at all.
func Test_Writer_Write_sendFails(t *testing.T) {
	connection, _ := newWriterConnection(t, nil,
		func(*writerCodec, spi.Message, spi.HandleMessage, spi.HandleError) error {
			return errors.New("transport is gone")
		})

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("switch"))
}

// Test_Writer_Write_roundTripsThroughTheReadPath pins that what the write-path puts onto
// the bus is exactly what the read-path decodes again. The read-path used to skip the first
// payload byte, which ate the value of a small datapoint-type and shifted a bigger one.
func Test_Writer_Write_roundTripsThroughTheReadPath(t *testing.T) {
	tests := []struct {
		name       string
		tagAddress string
		value      any
		assertion  func(t *testing.T, value apiValues.PlcValue)
	}{
		{
			name:       "a datapoint-type which fits into the embedded data bits",
			tagAddress: "1/2/3:DPT_Switch",
			value:      true,
			assertion: func(t *testing.T, value apiValues.PlcValue) {
				assert.True(t, value.GetBool())
			},
		},
		{
			name:       "a datapoint-type which needs its own data byte",
			tagAddress: "1/2/3:DPT_Scaling",
			value:      uint8(0x80),
			assertion: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint8(0x80), value.GetUint8())
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			connection, codec := newWriterConnection(t, nil, confirmingGateway(false))

			writeRequestBuilder := connection.WriteRequestBuilder()
			writeRequestBuilder.AddTagAddress("tag", test.tagAddress, test.value)
			writeRequest, err := writeRequestBuilder.Build()
			require.NoError(t, err)

			result := executeWrite(t, t.Context(), writeRequest)
			require.NoError(t, result.GetErr())
			require.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("tag"))

			// Reassemble the group-value payload the way the connection does for an
			// incoming telegram and feed it back through the read-path.
			sentRequests := codec.getSentRequests()
			require.Len(t, sentRequests, 1)
			frame := sentRequests[0].(driverModel.TunnelingRequest).GetCemi().(driverModel.LDataReq).GetDataFrame().(driverModel.LDataExtended)
			groupValueWrite := frame.GetApdu().(driverModel.ApduDataContainer).GetDataApdu().(driverModel.ApduDataGroupValueWrite)
			payload := append([]byte{byte(groupValueWrite.GetDataFirstByte())}, groupValueWrite.GetData()...)

			connection.handleValueCacheUpdate(t.Context(), frame.GetDestinationAddress(), payload)

			readRequestBuilder := connection.ReadRequestBuilder()
			readRequestBuilder.AddTagAddress("tag", test.tagAddress)
			readRequest, err := readRequestBuilder.Build()
			require.NoError(t, err)
			readResult := <-readRequest.Execute(t.Context())
			require.NoError(t, readResult.GetErr())
			require.Equal(t, apiModel.PlcResponseCode_OK, readResult.GetResponse().GetResponseCode("tag"))
			test.assertion(t, readResult.GetResponse().GetValue("tag"))
		})
	}
}

// Test_Writer_Write_ignoresForeignConfirmations pins that the confirmation expectation
// really only matches the frame this write sent: a confirmation echoing a different payload
// (e.g. a concurrent write of another client on the same tunnel) must not satisfy it.
func Test_Writer_Write_ignoresForeignConfirmations(t *testing.T) {
	foreignConfirmation := func(codec *writerCodec, message spi.Message, handleMessage spi.HandleMessage, _ spi.HandleError) error {
		tunnelingRequest := message.(driverModel.TunnelingRequest)
		sequenceCounter := tunnelingRequest.GetTunnelingRequestDataBlock().GetSequenceCounter()
		if err := handleMessage(driverModel.NewTunnelingResponse(
			driverModel.NewTunnelingResponseDataBlock(testCommunicationChannelId, sequenceCounter, driverModel.Status_NO_ERROR))); err != nil {
			return err
		}
		frame := tunnelingRequest.GetCemi().(driverModel.LDataReq).GetDataFrame().(driverModel.LDataExtended)
		// Same group address, but somebody else wrote a different value to it.
		assert.False(t, codec.deliver(driverModel.NewTunnelingRequest(
			driverModel.NewTunnelingRequestDataBlock(testCommunicationChannelId, sequenceCounter),
			driverModel.NewLDataCon(0, nil, driverModel.NewLDataExtended(
				frame.GetFrameType(), frame.GetNotRepeated(), frame.GetPriority(),
				frame.GetAcknowledgeRequested(), false, frame.GetGroupAddress(),
				frame.GetHopCount(), frame.GetExtendedFrameFormat(),
				driverModel.NewKnxAddress(1, 1, 2), frame.GetDestinationAddress(),
				driverModel.NewApduDataContainer(false, 0,
					driverModel.NewApduDataGroupValueWrite(0, nil)))))),
			"a confirmation of somebody else's frame must not be accepted")
		return nil
	}
	connection, _ := newWriterConnection(t, map[string][]string{"request-timeout": {"200"}}, foreignConfirmation)

	writeRequestBuilder := connection.WriteRequestBuilder()
	writeRequestBuilder.AddTagAddress("switch", "1/2/3:DPT_Switch", true)
	writeRequest, err := writeRequestBuilder.Build()
	require.NoError(t, err)

	result := executeWrite(t, t.Context(), writeRequest)
	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("switch"))
}

func Test_SerializeGroupValue(t *testing.T) {
	switchType := driverModel.KnxDatapointType_DPT_Switch
	scalingType := driverModel.KnxDatapointType_DPT_Scaling
	unknownType := driverModel.KnxDatapointType_DPT_UNKNOWN

	t.Run("no value", func(t *testing.T) {
		_, _, err := SerializeGroupValue(nil, &switchType)
		assert.Error(t, err)
	})
	t.Run("no datapoint-type", func(t *testing.T) {
		_, _, err := SerializeGroupValue(spiValues.NewPlcBOOL(true), nil)
		assert.Error(t, err)
	})
	t.Run("unknown datapoint-type", func(t *testing.T) {
		_, _, err := SerializeGroupValue(spiValues.NewPlcBOOL(true), &unknownType)
		assert.Error(t, err)
	})
	t.Run("small payload is packed into the first byte", func(t *testing.T) {
		dataFirstByte, data, err := SerializeGroupValue(spiValues.NewPlcBOOL(true), &switchType)
		require.NoError(t, err)
		assert.Equal(t, int8(1), dataFirstByte)
		assert.Empty(t, data)
	})
	t.Run("bigger payloads are sent as data bytes", func(t *testing.T) {
		dataFirstByte, data, err := SerializeGroupValue(spiValues.NewPlcUSINT(0x80), &scalingType)
		require.NoError(t, err)
		assert.Equal(t, int8(0), dataFirstByte)
		assert.Equal(t, []byte{0x80}, data)
	})
}
