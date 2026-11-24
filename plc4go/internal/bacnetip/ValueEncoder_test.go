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
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// writtenApplicationTag builds the single WriteProperty request for one
// (object type, property, array index) and returns the application tag that
// carries the value.
func writtenApplicationTag(t *testing.T, objType readWriteModel.BACnetObjectType, propId readWriteModel.BACnetPropertyIdentifier, arrayIndex *uint, value apiValues.PlcValue) readWriteModel.BACnetApplicationTag {
	t.Helper()
	tag := makeTag(objType, 1, propId).(*plcTag)
	tag.Properties[0].ArrayIndex = arrayIndex
	got, err := newTestWriter(t).buildServiceRequest(writeRequestFor(t, []writeTagSpec{{name: "p", tag: tag, value: value}}))
	require.NoError(t, err)
	wp, ok := got.(readWriteModel.BACnetConfirmedServiceRequestWriteProperty)
	require.True(t, ok, "expected WriteProperty, got %T", got)
	cd, ok := wp.GetPropertyValue().(readWriteModel.BACnetConstructedDataUnspecified)
	require.True(t, ok, "expected ConstructedDataUnspecified, got %T", wp.GetPropertyValue())
	require.Len(t, cd.GetData(), 1)
	return cd.GetData()[0].GetApplicationTag()
}

func TestBuildServiceRequest_EnumeratedPropertiesGoOutEnumerated(t *testing.T) {
	idx := func(i uint) *uint { return &i }
	cases := []struct {
		name       string
		objType    readWriteModel.BACnetObjectType
		propId     readWriteModel.BACnetPropertyIdentifier
		arrayIndex *uint
	}{
		{"BO Polarity", readWriteModel.BACnetObjectType_BINARY_OUTPUT, readWriteModel.BACnetPropertyIdentifier_POLARITY, nil},
		{"AI Reliability", readWriteModel.BACnetObjectType_ANALOG_INPUT, readWriteModel.BACnetPropertyIdentifier_RELIABILITY, nil},
		{"AV Event_State", readWriteModel.BACnetObjectType_ANALOG_VALUE, readWriteModel.BACnetPropertyIdentifier_EVENT_STATE, nil},
		{"AV Notify_Type", readWriteModel.BACnetObjectType_ANALOG_VALUE, readWriteModel.BACnetPropertyIdentifier_NOTIFY_TYPE, nil},
		{"AV Units", readWriteModel.BACnetObjectType_ANALOG_VALUE, readWriteModel.BACnetPropertyIdentifier_UNITS, nil},
		{"Loop Action", readWriteModel.BACnetObjectType_LOOP, readWriteModel.BACnetPropertyIdentifier_ACTION, nil},
		{"Access Door Present_Value", readWriteModel.BACnetObjectType_ACCESS_DOOR, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil},
		{"Life Safety Point Present_Value", readWriteModel.BACnetObjectType_LIFE_SAFETY_POINT, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil},
		{"Load Control Present_Value", readWriteModel.BACnetObjectType_LOAD_CONTROL, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil},
		{"BO Relinquish_Default", readWriteModel.BACnetObjectType_BINARY_OUTPUT, readWriteModel.BACnetPropertyIdentifier_RELINQUISH_DEFAULT, nil},
		{"BV Alarm_Value", readWriteModel.BACnetObjectType_BINARY_VALUE, readWriteModel.BACnetPropertyIdentifier_ALARM_VALUE, nil},
		{"Life Safety Point Maintenance_Required", readWriteModel.BACnetObjectType_LIFE_SAFETY_POINT, readWriteModel.BACnetPropertyIdentifier_MAINTENANCE_REQUIRED, nil},
		{"BV Present_Value", readWriteModel.BACnetObjectType_BINARY_VALUE, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil},
		{"Lift Car_Door_Status[1]", readWriteModel.BACnetObjectType_LIFT, readWriteModel.BACnetPropertyIdentifier_CAR_DOOR_STATUS, idx(1)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := writtenApplicationTag(t, c.objType, c.propId, c.arrayIndex, spiValues.NewPlcUDINT(1))
			enum, ok := got.(readWriteModel.BACnetApplicationTagEnumerated)
			require.True(t, ok, "expected Enumerated, got %T", got)
			assert.Equal(t, uint32(1), enum.GetActualValue())
			raw, err := got.Serialize()
			require.NoError(t, err)
			assert.Equal(t, []byte{0x91, 0x01}, raw)
		})
	}
}

func TestBuildServiceRequest_NonEnumeratedPropertiesKeepTheirTag(t *testing.T) {
	idx := func(i uint) *uint { return &i }
	cases := []struct {
		name       string
		objType    readWriteModel.BACnetObjectType
		propId     readWriteModel.BACnetPropertyIdentifier
		arrayIndex *uint
		value      apiValues.PlcValue
		want       readWriteModel.BACnetApplicationTag
	}{
		{"MSV Present_Value is Unsigned", readWriteModel.BACnetObjectType_MULTI_STATE_VALUE, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil,
			spiValues.NewPlcUDINT(2), readWriteModel.CreateBACnetApplicationTagUnsignedInteger(2)},
		{"MSO Relinquish_Default is Unsigned", readWriteModel.BACnetObjectType_MULTI_STATE_OUTPUT, readWriteModel.BACnetPropertyIdentifier_RELINQUISH_DEFAULT, nil,
			spiValues.NewPlcUDINT(2), readWriteModel.CreateBACnetApplicationTagUnsignedInteger(2)},
		{"MSO Feedback_Value is Unsigned", readWriteModel.BACnetObjectType_MULTI_STATE_OUTPUT, readWriteModel.BACnetPropertyIdentifier_FEEDBACK_VALUE, nil,
			spiValues.NewPlcUDINT(2), readWriteModel.CreateBACnetApplicationTagUnsignedInteger(2)},
		{"Array size of an enumerated array is Unsigned", readWriteModel.BACnetObjectType_LIFT, readWriteModel.BACnetPropertyIdentifier_CAR_DOOR_STATUS, idx(0),
			spiValues.NewPlcUDINT(2), readWriteModel.CreateBACnetApplicationTagUnsignedInteger(2)},
		{"Life Safety Zone Maintenance_Required is Boolean", readWriteModel.BACnetObjectType_LIFE_SAFETY_ZONE, readWriteModel.BACnetPropertyIdentifier_MAINTENANCE_REQUIRED, nil,
			spiValues.NewPlcBOOL(true), readWriteModel.CreateBACnetApplicationTagBoolean(true)},
		{"Out_Of_Service is Boolean", readWriteModel.BACnetObjectType_BINARY_OUTPUT, readWriteModel.BACnetPropertyIdentifier_OUT_OF_SERVICE, nil,
			spiValues.NewPlcBOOL(true), readWriteModel.CreateBACnetApplicationTagBoolean(true)},
		{"AV Present_Value is Real", readWriteModel.BACnetObjectType_ANALOG_VALUE, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil,
			spiValues.NewPlcREAL(2.5), readWriteModel.CreateBACnetApplicationTagReal(2.5)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := writtenApplicationTag(t, c.objType, c.propId, c.arrayIndex, c.value)
			gotRaw, err := got.Serialize()
			require.NoError(t, err)
			wantRaw, err := c.want.Serialize()
			require.NoError(t, err)
			assert.Equal(t, wantRaw, gotRaw)
		})
	}
}

func TestBuildServiceRequest_ProprietaryPropertyStaysUnsigned(t *testing.T) {
	objType := readWriteModel.BACnetObjectType_ANALOG_VALUE
	proprietary := uint32(1000)
	tag := &plcTag{
		ObjectId:   objectId{ObjectIdType: &objType, ObjectIdInstance: 1},
		Properties: []property{{PropertyIdentifierProprietary: &proprietary}},
	}
	got, err := newTestWriter(t).buildServiceRequest(writeRequestFor(t, []writeTagSpec{{name: "p", tag: tag, value: spiValues.NewPlcUDINT(1)}}))
	require.NoError(t, err)
	cd := got.(readWriteModel.BACnetConfirmedServiceRequestWriteProperty).GetPropertyValue().(readWriteModel.BACnetConstructedDataUnspecified)
	_, ok := cd.GetData()[0].GetApplicationTag().(readWriteModel.BACnetApplicationTagUnsignedInteger)
	assert.True(t, ok, "a proprietary property has no known datatype and stays Unsigned, got %T", cd.GetData()[0].GetApplicationTag())
}

// TestHintForProperty_MatchesTheMspecEnumeratedCases cross-checks the
// hand-written datatype table against the generated BACnetConstructedData
// parser: every (object type, property) in the table must parse an Enumerated
// into a property-specific case, and every property-specific case that takes
// an Enumerated must be in the table unless the standard says otherwise. A
// table entry is also wrong if the mspec has a case specific to that object
// type that takes a Boolean or Unsigned, as for a Life Safety Zone's
// Maintenance_Required.
func TestHintForProperty_MatchesTheMspecEnumeratedCases(t *testing.T) {
	// Opening tag 3, Enumerated 1, closing tag 3.
	payload := []byte{0x3E, 0x91, 0x01, 0x3F}
	otherDatatypes := map[string][]byte{
		"Boolean":  {0x3E, 0x11, 0x3F},
		"Unsigned": {0x3E, 0x21, 0x01, 0x3F},
	}
	notEnumeratedInTheStandard := func(objType readWriteModel.BACnetObjectType, propId readWriteModel.BACnetPropertyIdentifier) bool {
		switch propId {
		case readWriteModel.BACnetPropertyIdentifier_ACKED_TRANSITIONS,
			readWriteModel.BACnetPropertyIdentifier_ACK_REQUIRED:
			return true // BACnetEventTransitionBits, a BIT STRING; the mspec keys them on tag 9
		case readWriteModel.BACnetPropertyIdentifier_SCHEDULE_DEFAULT:
			return true // any datatype
		case readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE:
			return objType == readWriteModel.BACnetObjectType_CHANNEL || objType == readWriteModel.BACnetObjectType_SCHEDULE // any datatype
		case readWriteModel.BACnetPropertyIdentifier_INTERFACE_VALUE:
			// BACnetOptionalUnsigned in the standard; the mspec parses BACnetOptionalBinaryPV.
			return objType == readWriteModel.BACnetObjectType_MULTI_STATE_INPUT || objType == readWriteModel.BACnetObjectType_MULTI_STATE_OUTPUT
		case readWriteModel.BACnetPropertyIdentifier_ALARM_VALUES,
			readWriteModel.BACnetPropertyIdentifier_FAULT_VALUES,
			readWriteModel.BACnetPropertyIdentifier_MAINTENANCE_REQUIRED,
			readWriteModel.BACnetPropertyIdentifier_TRACKING_VALUE:
			return true // the mspec's catch-all case takes any object type; the table lists the enumerated ones
		}
		return false
	}
	for _, objType := range readWriteModel.BACnetObjectTypeValues {
		if objType == readWriteModel.BACnetObjectType_VENDOR_PROPRIETARY_VALUE {
			continue
		}
		for _, propId := range readWriteModel.BACnetPropertyIdentifierValues {
			switch propId {
			case readWriteModel.BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE,
				readWriteModel.BACnetPropertyIdentifier_ALL,
				readWriteModel.BACnetPropertyIdentifier_REQUIRED,
				readWriteModel.BACnetPropertyIdentifier_OPTIONAL:
				continue
			}
			enumeratedCase, err := typedConstructedDataCase(payload, objType, propId)
			mspecTyped := enumeratedCase != ""
			tableEnumerated := hintForProperty(uint32(objType), uint32(propId), nil) == hintEnumerated
			if tableEnumerated {
				assert.True(t, mspecTyped, "%s.%s is Enumerated in the table but the mspec has no case for it (err=%v)", objType, propId, err)
				for name, other := range otherDatatypes {
					otherCase, _ := typedConstructedDataCase(other, objType, propId)
					// A wildcard case also matches an object type that has no case of its own.
					wildcardCase, _ := typedConstructedDataCase(other, readWriteModel.BACnetObjectType_VENDOR_PROPRIETARY_VALUE, propId)
					assert.True(t, otherCase == "" || otherCase == enumeratedCase || otherCase == wildcardCase,
						"%s.%s is Enumerated in the table but the mspec has a separate %s case %s for this object type", objType, propId, name, otherCase)
				}
			} else if mspecTyped && !notEnumeratedInTheStandard(objType, propId) {
				t.Errorf("%s.%s: the mspec parses an Enumerated into a typed case but the table has no Enumerated entry", objType, propId)
			}
		}
	}
}

// typedConstructedDataCase returns the name of the property-specific
// BACnetConstructedData case the generated parser routes the payload into, or
// "" for the generic Unspecified case or a parse error. Some typed cases panic
// in ReadEnumGenericFailing for enum types it does not implement yet; reaching
// that panic still means a typed case was chosen, but its name is unknown.
func typedConstructedDataCase(payload []byte, objType readWriteModel.BACnetObjectType, propId readWriteModel.BACnetPropertyIdentifier) (typedCase string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); !ok || !strings.Contains(msg, "not yet implemented") {
				panic(r)
			}
			typedCase, err = "unimplemented enum", nil
		}
	}()
	parsed, err := readWriteModel.BACnetConstructedDataParse[readWriteModel.BACnetConstructedData](context.Background(), payload, 3, objType, propId, nil)
	if err != nil {
		return "", err
	}
	if _, unspecified := parsed.(readWriteModel.BACnetConstructedDataUnspecified); unspecified {
		return "", nil
	}
	return fmt.Sprintf("%T", parsed), nil
}
