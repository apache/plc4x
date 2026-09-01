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

package s7

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

func TestPickCyclicInterval(t *testing.T) {
	tests := []struct {
		requested time.Duration
		base      readWriteModel.TimeBase
		factor    uint8
	}{
		{0, readWriteModel.TimeBase_B1SEC, 1},
		{100 * time.Millisecond, readWriteModel.TimeBase_B01SEC, 1},
		{500 * time.Millisecond, readWriteModel.TimeBase_B01SEC, 5},
		{time.Second, readWriteModel.TimeBase_B1SEC, 1},
		{2500 * time.Millisecond, readWriteModel.TimeBase_B1SEC, 3},
		{30 * time.Second, readWriteModel.TimeBase_B10SEC, 3},
		{255 * time.Second, readWriteModel.TimeBase_B1SEC, 255},
		{2550 * time.Second, readWriteModel.TimeBase_B10SEC, 255},
	}
	for _, test := range tests {
		t.Run(test.requested.String(), func(t *testing.T) {
			interval := pickCyclicInterval(test.requested)
			assert.Equal(t, test.base, interval.base, "base")
			assert.Equal(t, test.factor, interval.factor, "factor")
		})
	}
}

// roundTrip serializes a request and parses it back, proving the declared lengths are
// consistent with the actual wire representation.
func roundTrip(t *testing.T, packet readWriteModel.TPKTPacket) readWriteModel.S7MessageUserData {
	t.Helper()
	data, err := packet.Serialize()
	require.NoError(t, err)
	parsed, err := readWriteModel.TPKTPacketParse(t.Context(), data)
	require.NoError(t, err)
	cotpPacketData := parsed.GetPayload().(readWriteModel.COTPPacketData)
	return cotpPacketData.GetPayload().(readWriteModel.S7MessageUserData)
}

func TestBuildMsgSubscriptionRequest(t *testing.T) {
	message := roundTrip(t, buildMsgSubscriptionRequest(42, readWriteModel.AlarmStateType_ALARM_S_INITIATE))
	payload := message.GetPayload().(readWriteModel.S7PayloadUserData).GetItems()[0]
	request := payload.(readWriteModel.S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest)
	assert.Equal(t, uint8(0x80), request.GetSubscription())
	assert.Equal(t, msgSubscriptionMagicKey, request.GetMagicKey())
	require.NotNil(t, request.GetAlarmtype())
	assert.Equal(t, readWriteModel.AlarmStateType_ALARM_S_INITIATE, *request.GetAlarmtype())
}

func TestBuildCyclicSubscribeRequest(t *testing.T) {
	tag := NewTag(readWriteModel.MemoryArea_DATA_BLOCKS, 5, 10, 0, 2, readWriteModel.TransportSize_INT).(plcTag)
	message := roundTrip(t, buildCyclicSubscribeRequest(43, []PlcTag{tag}, cyclicInterval{readWriteModel.TimeBase_B1SEC, 2}))
	payload := message.GetPayload().(readWriteModel.S7PayloadUserData).GetItems()[0]
	request := payload.(readWriteModel.S7PayloadUserDataItemCyclicServicesSubscribeRequest)
	assert.Equal(t, uint16(1), request.GetItemsCount())
	assert.Equal(t, readWriteModel.TimeBase_B1SEC, request.GetTimeBase())
	assert.Equal(t, uint8(2), request.GetTimeFactor())
	item := request.GetItem()[0].(readWriteModel.CycServiceItemAnyType)
	assert.Equal(t, readWriteModel.TransportSize_INT, item.GetTransportSize())
	assert.Equal(t, uint16(2), item.GetLength())
	assert.Equal(t, uint16(5), item.GetDbNumber())
	assert.Equal(t, readWriteModel.MemoryArea_DATA_BLOCKS, item.GetMemoryArea())
	assert.Equal(t, uint32(10<<3), item.GetAddress())
}

func TestBuildAlarmQueryRequest(t *testing.T) {
	message := roundTrip(t, buildAlarmQueryRequest(44, readWriteModel.QueryType_ALARM_S))
	payload := message.GetPayload().(readWriteModel.S7PayloadUserData).GetItems()[0]
	request := payload.(readWriteModel.S7PayloadUserDataItemCpuFunctionAlarmQueryRequest)
	assert.Equal(t, readWriteModel.SyntaxIdType_ALARM_QUERYREQSET, request.GetSyntaxId())
	assert.Equal(t, readWriteModel.QueryType_ALARM_S, request.GetQueryType())
	assert.Equal(t, readWriteModel.AlarmType_ALARM_S, request.GetAlarmType())
}

func TestParseCyclicSubscribeResponse(t *testing.T) {
	makeResponse := func(payloadItem readWriteModel.S7PayloadUserDataItem, sequenceNumber uint8) readWriteModel.S7Message {
		return readWriteModel.NewS7MessageUserData(
			45,
			readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{
				readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
					0x12, 0x08, 0x02, 0x01, sequenceNumber,
					new(uint8(0)), new(uint8(1)), new(uint16(0)),
				),
			}),
			readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{payloadItem}),
		)
	}

	t.Run("empty response carries jobId", func(t *testing.T) {
		jobId, err := parseCyclicSubscribeResponse(makeResponse(
			readWriteModel.NewS7PayloadUserDataItemCyclicServicesSubscribeEmptyResponse(
				readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_NULL, 0),
			7,
		))
		require.NoError(t, err)
		assert.Equal(t, uint8(7), jobId)
	})
	t.Run("value-carrying response", func(t *testing.T) {
		jobId, err := parseCyclicSubscribeResponse(makeResponse(
			readWriteModel.NewS7PayloadUserDataItemCyclicServicesSubscribeResponse(
				readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_OCTET_STRING, 8,
				1,
				[]readWriteModel.AssociatedValueType{
					readWriteModel.NewAssociatedValueType(readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_BYTE_WORD_DWORD, 2, []uint8{0x00, 0x04}),
				},
			),
			9,
		))
		require.NoError(t, err)
		assert.Equal(t, uint8(9), jobId)
	})
}

func TestExtractCyclicPushItems(t *testing.T) {
	message := readWriteModel.NewS7MessageUserData(
		0,
		readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{
			readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
				0x12, 0x00, 0x02, 0x01, 3,
				new(uint8(0)), new(uint8(1)), nil,
			),
		}),
		readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{
			readWriteModel.NewS7PayloadUserDataItemCyclicServicesPush(
				readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_OCTET_STRING, 12,
				2,
				[]readWriteModel.AssociatedValueType{
					readWriteModel.NewAssociatedValueType(readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_BYTE_WORD_DWORD, 2, []uint8{0x00, 0x04}),
					readWriteModel.NewAssociatedValueType(readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_BYTE_WORD_DWORD, 1, []uint8{0xFF}),
				},
			),
		}),
	)
	items := extractCyclicPushItems(message)
	require.Len(t, items, 2)
	assert.Equal(t, []byte{0x00, 0x04}, items[0])
	assert.Equal(t, []byte{0xFF}, items[1])

	group, functionType, subfunction, ok := userDataPushKey(message)
	require.True(t, ok)
	assert.Equal(t, uint8(0x02), group)
	assert.Equal(t, uint8(0x00), functionType)
	assert.Equal(t, uint8(0x01), subfunction)
	jobId, ok := userDataSequenceNumber(message)
	require.True(t, ok)
	assert.Equal(t, uint8(3), jobId)
}

func TestParseAlarmIndication(t *testing.T) {
	state := func(mask uint8) readWriteModel.State {
		return readWriteModel.NewState(
			mask&0x80 != 0, mask&0x40 != 0, mask&0x20 != 0, mask&0x10 != 0,
			mask&0x08 != 0, mask&0x04 != 0, mask&0x02 != 0, mask&0x01 != 0,
		)
	}
	message := readWriteModel.NewS7MessageUserData(
		0,
		readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{
			readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
				0x12, 0x00, 0x04, 0x12, 0,
				new(uint8(0)), new(uint8(1)), nil,
			),
		}),
		readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{
			readWriteModel.NewS7PayloadAlarmS(
				readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_OCTET_STRING, 0,
				readWriteModel.NewAlarmMessagePushType(
					readWriteModel.NewDateAndTime(24, 8, 18, 10, 30, 0, 123, 1),
					0x11,
					1,
					[]readWriteModel.AlarmMessageObjectPushType{
						readWriteModel.NewAlarmMessageObjectPushType(
							0x04, readWriteModel.SyntaxIdType_ALARM_INDSET, 0,
							0x12345678,
							state(0x01), state(0x00), state(0x00), state(0x01),
							nil,
						),
					},
				),
			),
		}),
	)
	value := parseAlarmIndication(message)
	require.NotNil(t, value)
	require.True(t, value.IsStruct())
	structValue := value.GetStruct()
	assert.Equal(t, uint32(0x12345678), structValue["eventId"].GetUint32())
	assert.Equal(t, uint32(0x01), structValue["eventState"].GetUint32())
	assert.Equal(t, uint32(0x01), structValue["ackStateComing"].GetUint32())
	assert.Equal(t, "2024-08-18T10:30:00.123", structValue["plcTimestamp"].GetString())
	assert.Contains(t, structValue, "receivedAt")
}
