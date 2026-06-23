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

package bacnetip

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// newTestReader returns a Reader that does no I/O — all tests here exercise the
// pure response-decoding path against synthetic APDUs.
func newTestReader(t *testing.T) *Reader {
	t.Helper()
	return &Reader{
		invokeIdGenerator:     &InvokeIdGenerator{},
		maxSegmentsAccepted:   readWriteModel.MaxSegmentsAccepted_NUM_SEGMENTS_16,
		maxApduLengthAccepted: readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1476,
		log:                   zerolog.Nop(),
	}
}

// readRequestForTags wraps the test's tag names in a minimal PlcReadRequest so
// the decoder has something to dispatch against.
func readRequestForTags(tagNames ...string) apiModel.PlcReadRequest {
	tags := map[string]apiModel.PlcTag{}
	for _, n := range tagNames {
		// We use a fully synthetic tag — the decoder only consults names and the
		// returned response codes / values, not the tag's structure.
		tags[n] = nil
	}
	return spiModel.NewDefaultPlcReadRequest(tags, tagNames, nil, nil)
}

// constructedDataFromTag wraps a single ApplicationTag in BACnetConstructedDataUnspecified
// so we can build the BACnetServiceAckReadProperty response fixture.
func constructedDataFromTag(tag readWriteModel.BACnetApplicationTag) readWriteModel.BACnetConstructedData {
	header := readWriteModel.CreateBACnetTagHeaderBalanced(true, 3, 0)
	element := readWriteModel.NewBACnetConstructedDataElement(header, tag, nil, nil)
	return readWriteModel.NewBACnetConstructedDataUnspecified(
		readWriteModel.CreateBACnetOpeningTag(3),
		header,
		readWriteModel.CreateBACnetClosingTag(3),
		nil,
		[]readWriteModel.BACnetConstructedDataElement{element},
	)
}

func buildReadPropertyAck(t *testing.T, objType uint16, instance uint32, propId readWriteModel.BACnetPropertyIdentifier, tag readWriteModel.BACnetApplicationTag) readWriteModel.APDUComplexAck {
	t.Helper()
	values := constructedDataFromTag(tag)
	serviceAck := readWriteModel.NewBACnetServiceAckReadProperty(
		0,
		readWriteModel.CreateBACnetContextTagObjectIdentifier(0, objType, instance),
		readWriteModel.CreateBACnetPropertyIdentifierTagged(1, uint32(propId)),
		nil,
		values,
	)
	return readWriteModel.NewAPDUComplexAck(false, false, 1, nil, nil, serviceAck, nil, nil)
}

func TestToPlc4xReadResponse_RealPresentValue(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("temp")
	ack := buildReadPropertyAck(
		t,
		uint16(readWriteModel.BACnetObjectType_ANALOG_INPUT),
		1,
		readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE,
		readWriteModel.CreateBACnetApplicationTagReal(23.5),
	)

	resp, err := reader.ToPlc4xReadResponse(ack, request)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("temp"))
	val := resp.GetValue("temp")
	require.NotNil(t, val)
	assert.True(t, val.IsFloat32())
	assert.InDelta(t, 23.5, val.GetFloat32(), 0.001)
}

func TestToPlc4xReadResponse_BooleanPresentValue(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("binIn")
	ack := buildReadPropertyAck(
		t,
		uint16(readWriteModel.BACnetObjectType_BINARY_INPUT),
		7,
		readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE,
		readWriteModel.CreateBACnetApplicationTagBoolean(true),
	)
	resp, err := reader.ToPlc4xReadResponse(ack, request)
	require.NoError(t, err)
	val := resp.GetValue("binIn")
	assert.True(t, val.GetBool())
}

func TestToPlc4xReadResponse_CharacterStringObjectName(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("name")
	ack := buildReadPropertyAck(
		t,
		uint16(readWriteModel.BACnetObjectType_DEVICE),
		1234,
		readWriteModel.BACnetPropertyIdentifier_OBJECT_NAME,
		readWriteModel.CreateBACnetApplicationTagCharacterString(readWriteModel.BACnetCharacterEncoding_ISO_10646, "JACE/N4"),
	)
	resp, err := reader.ToPlc4xReadResponse(ack, request)
	require.NoError(t, err)
	assert.Equal(t, "JACE/N4", resp.GetValue("name").GetString())
}

func TestToPlc4xReadResponse_UnsignedInteger(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("count")
	ack := buildReadPropertyAck(
		t,
		uint16(readWriteModel.BACnetObjectType_ACCUMULATOR),
		2,
		readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE,
		readWriteModel.CreateBACnetApplicationTagUnsignedInteger(42),
	)
	resp, err := reader.ToPlc4xReadResponse(ack, request)
	require.NoError(t, err)
	assert.Equal(t, uint64(42), resp.GetValue("count").GetUint64())
}

func TestToPlc4xReadResponse_Enumerated(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("state")
	ack := buildReadPropertyAck(
		t,
		uint16(readWriteModel.BACnetObjectType_BINARY_OUTPUT),
		1,
		readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE,
		readWriteModel.CreateBACnetApplicationTagEnumerated(1), // ACTIVE
	)
	resp, err := reader.ToPlc4xReadResponse(ack, request)
	require.NoError(t, err)
	assert.Equal(t, uint32(1), resp.GetValue("state").GetUint32())
}

func TestToPlc4xReadResponse_Null(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("nullProp")
	ack := buildReadPropertyAck(
		t,
		uint16(readWriteModel.BACnetObjectType_ANALOG_INPUT),
		1,
		readWriteModel.BACnetPropertyIdentifier_OUT_OF_SERVICE,
		readWriteModel.CreateBACnetApplicationTagNull(),
	)
	resp, err := reader.ToPlc4xReadResponse(ack, request)
	require.NoError(t, err)
	val := resp.GetValue("nullProp")
	require.NotNil(t, val)
	assert.Equal(t, apiValues.NULL, val.GetPlcValueType())
}

// ── Error / Reject / Abort paths ───────────────────────────────────────────

func buildErrorAPDU(class readWriteModel.ErrorClass, code readWriteModel.ErrorCode) readWriteModel.APDUError {
	header := readWriteModel.CreateBACnetTagHeaderBalanced(false, 9, 1)
	bacError := readWriteModel.NewError(
		readWriteModel.NewErrorClassTagged(header, class, 0),
		readWriteModel.NewErrorCodeTagged(header, code, 0),
	)
	general := readWriteModel.NewBACnetErrorGeneral(bacError)
	return readWriteModel.NewAPDUError(1, readWriteModel.BACnetConfirmedServiceChoice_READ_PROPERTY, general)
}

func TestToPlc4xReadResponse_UnknownObjectError(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("missing")
	apdu := buildErrorAPDU(readWriteModel.ErrorClass_OBJECT, readWriteModel.ErrorCode_UNKNOWN_OBJECT)
	resp, err := reader.ToPlc4xReadResponse(apdu, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, resp.GetResponseCode("missing"))
}

func TestToPlc4xReadResponse_UnknownPropertyError(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("badProp")
	apdu := buildErrorAPDU(readWriteModel.ErrorClass_PROPERTY, readWriteModel.ErrorCode_UNKNOWN_PROPERTY)
	resp, err := reader.ToPlc4xReadResponse(apdu, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, resp.GetResponseCode("badProp"))
}

func TestToPlc4xReadResponse_WriteAccessDenied(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("ro")
	apdu := buildErrorAPDU(readWriteModel.ErrorClass_PROPERTY, readWriteModel.ErrorCode_WRITE_ACCESS_DENIED)
	resp, err := reader.ToPlc4xReadResponse(apdu, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, resp.GetResponseCode("ro"))
}

func TestToPlc4xReadResponse_AbortSegmentationNotSupported(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("a")
	abort := readWriteModel.NewAPDUAbort(false, 1,
		readWriteModel.NewBACnetAbortReasonTagged(1, readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED, 0))
	resp, err := reader.ToPlc4xReadResponse(abort, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, resp.GetResponseCode("a"))
}

func TestToPlc4xReadResponse_AbortOther(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("a")
	abort := readWriteModel.NewAPDUAbort(false, 1,
		readWriteModel.NewBACnetAbortReasonTagged(1, readWriteModel.BACnetAbortReason_BUFFER_OVERFLOW, 0))
	resp, err := reader.ToPlc4xReadResponse(abort, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, resp.GetResponseCode("a"))
}

func TestToPlc4xReadResponse_Reject(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("r")
	reject := readWriteModel.NewAPDUReject(1,
		readWriteModel.NewBACnetRejectReasonTagged(1, readWriteModel.BACnetRejectReason_UNRECOGNIZED_SERVICE, 0))
	resp, err := reader.ToPlc4xReadResponse(reject, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, resp.GetResponseCode("r"))
}

func TestToPlc4xReadResponse_Segmented_DefensiveFallback(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("big")
	serviceAck := readWriteModel.NewBACnetServiceAckReadProperty(
		0,
		readWriteModel.CreateBACnetContextTagObjectIdentifier(0, uint16(readWriteModel.BACnetObjectType_DEVICE), 1),
		readWriteModel.CreateBACnetPropertyIdentifierTagged(1, uint32(readWriteModel.BACnetPropertyIdentifier_OBJECT_LIST)),
		nil,
		constructedDataFromTag(readWriteModel.CreateBACnetApplicationTagNull()),
	)
	// In the live read flow a segmented APDU is intercepted and reassembled
	// before ToPlc4xReadResponse. Reaching decodeComplexAck with a segmented
	// APDU is unexpected, so it falls back to UNSUPPORTED rather than panicking.
	apdu := readWriteModel.NewAPDUComplexAck(true, true, 1, nil, nil, serviceAck, nil, nil)
	resp, err := reader.ToPlc4xReadResponse(apdu, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, resp.GetResponseCode("big"))
}

// TestReassembledSegments_RoundTrip verifies the core segmentation path: a real
// service ack serialized and split across two APDUComplexAck segments is
// reassembled, reparsed, and decoded back into the original value.
func TestReassembledSegments_RoundTrip(t *testing.T) {
	reader := newTestReader(t)
	request := readRequestForTags("big")

	serviceAck := readWriteModel.NewBACnetServiceAckReadProperty(
		0,
		readWriteModel.CreateBACnetContextTagObjectIdentifier(0, uint16(readWriteModel.BACnetObjectType_ANALOG_INPUT), 1),
		readWriteModel.CreateBACnetPropertyIdentifierTagged(1, uint32(readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE)),
		nil,
		constructedDataFromTag(readWriteModel.CreateBACnetApplicationTagReal(23.5)),
	)
	fullBytes, err := serviceAck.Serialize()
	require.NoError(t, err)
	require.Greater(t, len(fullBytes), 2)

	// Split the serialized service ack across two segments. Segment 0 carries
	// the service-choice byte (as the wire format does); concatenation restores
	// the original bytes.
	split := len(fullBytes) / 2
	seg0 := readWriteModel.NewAPDUComplexAck(true, true, 7, ptrU8(0), ptrU8(1), nil, nil, fullBytes[:split])
	seg1 := readWriteModel.NewAPDUComplexAck(true, false, 7, ptrU8(1), ptrU8(1), nil, nil, fullBytes[split:])

	r := NewInboundReassembler(7, 1)
	ack0, err := r.AcceptSegment(seg0)
	require.NoError(t, err)
	require.NotNil(t, ack0)
	require.False(t, r.Complete())
	_, err = r.AcceptSegment(seg1)
	require.NoError(t, err)
	require.True(t, r.Complete())
	require.Equal(t, fullBytes, r.Bytes())

	parsed, err := readWriteModel.BACnetServiceAckParse[readWriteModel.BACnetServiceAck](context.Background(), r.Bytes(), uint32(len(r.Bytes())))
	require.NoError(t, err)

	resp, err := reader.decodeServiceAck(parsed, request)
	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("big"))
	val := resp.GetValue("big")
	require.NotNil(t, val)
	assert.InDelta(t, 23.5, val.GetFloat32(), 0.001)
}

func ptrU8(v uint8) *uint8 { return &v }

// ── ValueDecoder unit tests ───────────────────────────────────────────────

func TestAppTagToPlcValue_AllPrimitiveTypes(t *testing.T) {
	cases := []struct {
		name string
		tag  readWriteModel.BACnetApplicationTag
		eval func(t *testing.T, v apiValues.PlcValue)
	}{
		{"nil", nil, func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, apiValues.NULL, v.GetPlcValueType()) }},
		{"null tag", readWriteModel.CreateBACnetApplicationTagNull(), func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, apiValues.NULL, v.GetPlcValueType()) }},
		{"bool true", readWriteModel.CreateBACnetApplicationTagBoolean(true), func(t *testing.T, v apiValues.PlcValue) { assert.True(t, v.GetBool()) }},
		{"unsigned 7", readWriteModel.CreateBACnetApplicationTagUnsignedInteger(7), func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, uint64(7), v.GetUint64()) }},
		{"signed -3", readWriteModel.CreateBACnetApplicationTagSignedInteger(-3), func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, int64(-3), v.GetInt64()) }},
		{"real 3.5", readWriteModel.CreateBACnetApplicationTagReal(3.5), func(t *testing.T, v apiValues.PlcValue) { assert.InDelta(t, 3.5, v.GetFloat32(), 1e-6) }},
		{"double 1.25", readWriteModel.CreateBACnetApplicationTagDouble(1.25), func(t *testing.T, v apiValues.PlcValue) { assert.InDelta(t, 1.25, v.GetFloat64(), 1e-9) }},
		{"character string", readWriteModel.CreateBACnetApplicationTagCharacterString(readWriteModel.BACnetCharacterEncoding_ISO_10646, "hi"), func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, "hi", v.GetString()) }},
		{"octet string", readWriteModel.CreateBACnetApplicationTagOctetString([]byte{0xCA, 0xFE}), func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, []byte{0xCA, 0xFE}, v.GetRaw()) }},
		{"enumerated 4", readWriteModel.CreateBACnetApplicationTagEnumerated(4), func(t *testing.T, v apiValues.PlcValue) { assert.Equal(t, uint32(4), v.GetUint32()) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.eval(t, appTagToPlcValue(c.tag))
		})
	}
}

func TestBitsToBytes(t *testing.T) {
	cases := []struct {
		name string
		in   []bool
		want []byte
	}{
		{"empty", nil, []byte{}},
		{"all-zero byte", []bool{false, false, false, false, false, false, false, false}, []byte{0x00}},
		{"all-one byte", []bool{true, true, true, true, true, true, true, true}, []byte{0xFF}},
		{"mixed MSB", []bool{true, false, true, false, false, false, false, false}, []byte{0xA0}},
		{"partial trailing", []bool{true, true, true}, []byte{0xE0}},
		{"9 bits two bytes", []bool{true, false, false, false, false, false, false, false, true}, []byte{0x80, 0x80}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bitsToBytes(c.in)
			if len(c.want) == 0 {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, c.want, got)
			}
		})
	}
}

func TestConstructedDataToPlcValue_NilGivesNull(t *testing.T) {
	v := constructedDataToPlcValue(nil)
	assert.Equal(t, apiValues.NULL, v.GetPlcValueType())
}

func TestConstructedDataToPlcValue_UnspecifiedSingle(t *testing.T) {
	cd := constructedDataFromTag(readWriteModel.CreateBACnetApplicationTagReal(1.5))
	v := constructedDataToPlcValue(cd)
	assert.InDelta(t, 1.5, v.GetFloat32(), 1e-6)
}

func TestMapErrorClassCodeToResponseCode(t *testing.T) {
	cases := []struct {
		class readWriteModel.ErrorClass
		code  readWriteModel.ErrorCode
		want  apiModel.PlcResponseCode
	}{
		{readWriteModel.ErrorClass_OBJECT, readWriteModel.ErrorCode_UNKNOWN_OBJECT, apiModel.PlcResponseCode_NOT_FOUND},
		{readWriteModel.ErrorClass_PROPERTY, readWriteModel.ErrorCode_UNKNOWN_PROPERTY, apiModel.PlcResponseCode_INVALID_ADDRESS},
		{readWriteModel.ErrorClass_PROPERTY, readWriteModel.ErrorCode_WRITE_ACCESS_DENIED, apiModel.PlcResponseCode_ACCESS_DENIED},
		{readWriteModel.ErrorClass_DEVICE, readWriteModel.ErrorCode_OPERATIONAL_PROBLEM, apiModel.PlcResponseCode_REMOTE_ERROR},
		{readWriteModel.ErrorClass_SECURITY, readWriteModel.ErrorCode_OTHER, apiModel.PlcResponseCode_REMOTE_ERROR}, // unmapped → REMOTE_ERROR
	}
	for _, c := range cases {
		got := mapErrorClassCodeToResponseCode(c.class, c.code)
		assert.Equal(t, c.want, got, "class=%v code=%v", c.class, c.code)
	}
}

// Defensive: PlcRawByteArray expose value via GetRaw — keep tied to spi/values
// behavior so future changes there flag here.
var _ = spiValues.NewPlcRawByteArray
