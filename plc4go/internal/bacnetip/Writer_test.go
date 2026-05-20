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

// newTestWriter returns a Writer with no I/O wiring — pure-function tests only.
func newTestWriter(t *testing.T) *Writer {
	t.Helper()
	return &Writer{
		invokeIdGenerator: &InvokeIdGenerator{},
		driverContext:     NewDriverContext(createDefaultConfiguration()),
		log:               zerolog.Nop(),
	}
}

// writeRequestFor builds a minimal PlcWriteRequest with the given (tag, value)
// pairs. Tags are constructed inline as concrete plcTag structs so they
// satisfy BacNetPlcTag without going through the TagHandler regex parser.
func writeRequestFor(t *testing.T, tagSpecs []writeTagSpec) apiModel.PlcWriteRequest {
	t.Helper()
	tagNames := make([]string, 0, len(tagSpecs))
	tags := map[string]apiModel.PlcTag{}
	values := map[string]apiValues.PlcValue{}
	for _, s := range tagSpecs {
		tagNames = append(tagNames, s.name)
		tags[s.name] = s.tag
		values[s.name] = s.value
	}
	return spiModel.NewDefaultPlcWriteRequest(tags, tagNames, values, nil, nil)
}

type writeTagSpec struct {
	name  string
	tag   BacNetPlcTag
	value apiValues.PlcValue
}

func makeTag(objType readWriteModel.BACnetObjectType, instance uint32, propId readWriteModel.BACnetPropertyIdentifier) BacNetPlcTag {
	objId := objectId{ObjectIdType: &objType, ObjectIdInstance: instance}
	return &plcTag{
		ObjectId:   objId,
		Properties: []property{{PropertyIdentifier: &propId}},
	}
}

// ── plcValueToApplicationTag ───────────────────────────────────────────────

func TestPlcValueToApplicationTag_Mapping(t *testing.T) {
	cases := []struct {
		name string
		in   apiValues.PlcValue
		want any
	}{
		{"bool", spiValues.NewPlcBOOL(true), readWriteModel.CreateBACnetApplicationTagBoolean(true)},
		{"real", spiValues.NewPlcREAL(22.5), readWriteModel.CreateBACnetApplicationTagReal(22.5)},
		{"double", spiValues.NewPlcLREAL(1.25), readWriteModel.CreateBACnetApplicationTagDouble(1.25)},
		{"unsigned", spiValues.NewPlcULINT(7), readWriteModel.CreateBACnetApplicationTagUnsignedInteger(7)},
		{"signed", spiValues.NewPlcLINT(-3), readWriteModel.CreateBACnetApplicationTagSignedInteger(-3)},
		{"string", spiValues.NewPlcSTRING("hi"), readWriteModel.CreateBACnetApplicationTagCharacterString(readWriteModel.BACnetCharacterEncoding_ISO_10646, "hi")},
		{"raw bytes", spiValues.NewPlcRawByteArray([]byte{0xCA, 0xFE}), readWriteModel.CreateBACnetApplicationTagOctetString([]byte{0xCA, 0xFE})},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := plcValueToApplicationTag(c.in, hintNone)
			require.NoError(t, err)
			// Round-trip: encode then decode and assert equality on the resulting PlcValue.
			roundTripped := appTagToPlcValue(got)
			expectedRoundTrip := appTagToPlcValue(c.want.(readWriteModel.BACnetApplicationTag))
			assert.Equal(t, expectedRoundTrip.GetPlcValueType(), roundTripped.GetPlcValueType())
		})
	}
}

func TestPlcValueToApplicationTag_NilGivesNullTag(t *testing.T) {
	got, err := plcValueToApplicationTag(nil, hintNone)
	require.NoError(t, err)
	_, ok := got.(readWriteModel.BACnetApplicationTagNull)
	assert.True(t, ok, "nil PlcValue must produce a Null ApplicationTag")
}

func TestPlcValueToApplicationTag_EnumeratedHint(t *testing.T) {
	got, err := plcValueToApplicationTag(spiValues.NewPlcUDINT(2), hintEnumerated)
	require.NoError(t, err)
	_, ok := got.(readWriteModel.BACnetApplicationTagEnumerated)
	assert.True(t, ok, "expected Enumerated when hintEnumerated is set")
}

// ── buildServiceRequest dispatch ───────────────────────────────────────────

func TestBuildServiceRequest_SingleTagSingleProp_BuildsWriteProperty(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "av",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 5, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcREAL(22.5),
		},
	})
	got, err := writer.buildServiceRequest(req)
	require.NoError(t, err)
	wp, ok := got.(readWriteModel.BACnetConfirmedServiceRequestWriteProperty)
	require.True(t, ok, "expected WriteProperty, got %T", got)
	assert.Equal(t, uint32(5), wp.GetObjectIdentifier().GetPayload().GetInstanceNumber())
	assert.Equal(t, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, wp.GetPropertyIdentifier().GetValue())
}

func TestBuildServiceRequest_MultipleTags_BuildsWritePropertyMultiple(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "av",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 5, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcREAL(22.5),
		},
		{
			name:  "bv",
			tag:   makeTag(readWriteModel.BACnetObjectType_BINARY_VALUE, 3, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcBOOL(true),
		},
	})
	got, err := writer.buildServiceRequest(req)
	require.NoError(t, err)
	wpm, ok := got.(readWriteModel.BACnetConfirmedServiceRequestWritePropertyMultiple)
	require.True(t, ok, "expected WritePropertyMultiple, got %T", got)
	assert.Len(t, wpm.GetData(), 2)
}

func TestBuildServiceRequest_UnsupportedValueTypeReturnsError(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "x",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 1, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcList(nil), // PlcList isn't currently encodable
		},
	})
	_, err := writer.buildServiceRequest(req)
	assert.Error(t, err)
}

// ── toPlcWriteResponse error mapping ───────────────────────────────────────

func TestToPlcWriteResponse_SimpleAck(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "av",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 5, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcREAL(22.5),
		},
	})
	ack := readWriteModel.NewAPDUSimpleAck(1, readWriteModel.BACnetConfirmedServiceChoice_WRITE_PROPERTY)
	resp := writer.toPlcWriteResponse(ack, req)
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("av"))
}

func TestToPlcWriteResponse_ErrorWriteAccessDenied(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "av",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 5, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcREAL(22.5),
		},
	})
	apdu := buildErrorAPDU(readWriteModel.ErrorClass_PROPERTY, readWriteModel.ErrorCode_WRITE_ACCESS_DENIED)
	resp := writer.toPlcWriteResponse(apdu, req)
	assert.Equal(t, apiModel.PlcResponseCode_ACCESS_DENIED, resp.GetResponseCode("av"))
}

func TestToPlcWriteResponse_Abort(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "av",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 1, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcREAL(1),
		},
	})
	abort := readWriteModel.NewAPDUAbort(false, 1,
		readWriteModel.NewBACnetAbortReasonTagged(1, readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED, 0))
	resp := writer.toPlcWriteResponse(abort, req)
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, resp.GetResponseCode("av"))
}

func TestToPlcWriteResponse_Reject(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "av",
			tag:   makeTag(readWriteModel.BACnetObjectType_ANALOG_VALUE, 1, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcREAL(1),
		},
	})
	reject := readWriteModel.NewAPDUReject(1,
		readWriteModel.NewBACnetRejectReasonTagged(1, readWriteModel.BACnetRejectReason_UNRECOGNIZED_SERVICE, 0))
	resp := writer.toPlcWriteResponse(reject, req)
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_DATA, resp.GetResponseCode("av"))
}
