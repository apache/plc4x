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
	"math"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// plcValueToApplicationTag maps a PlcValue to a BACnet ApplicationTag, the
// inverse of appTagToPlcValue. The caller may pass a hint via the property
// identifier so we can prefer enumerated/character-string encoding when the
// PlcValue type alone is ambiguous (e.g. uint32 → Enumerated vs UnsignedInteger
// for properties like StatusFlags / Units).
//
// Returns an error if the PlcValue type isn't representable in BACnet.
func plcValueToApplicationTag(v apiValues.PlcValue, hint encodingHint) (model.BACnetApplicationTag, error) {
	if v == nil {
		return model.CreateBACnetApplicationTagNull(), nil
	}
	switch v.GetPlcValueType() {
	case apiValues.NULL:
		return model.CreateBACnetApplicationTagNull(), nil
	case apiValues.BOOL:
		// Binary objects' PRESENT_VALUE is Enumerated (INACTIVE=0, ACTIVE=1)
		// on the wire — sending a Boolean application tag gets rejected with
		// REJECT(INVALID_TAG). When the caller hints Enumerated, encode the
		// bool as 0/1 Enumerated.
		if hint == hintEnumerated {
			var v32 uint32
			if v.GetBool() {
				v32 = 1
			}
			return model.CreateBACnetApplicationTagEnumerated(v32), nil
		}
		return model.CreateBACnetApplicationTagBoolean(v.GetBool()), nil
	case apiValues.BYTE, apiValues.USINT, apiValues.UINT, apiValues.UDINT, apiValues.ULINT, apiValues.WORD, apiValues.DWORD, apiValues.LWORD:
		u := v.GetUint64()
		if hint == hintEnumerated {
			if u > math.MaxUint32 {
				return nil, errors.Errorf("%d does not fit a BACnet Enumerated, which this driver models as 32 bits", u)
			}
			return model.CreateBACnetApplicationTagEnumerated(uint32(u)), nil
		}
		return model.CreateBACnetApplicationTagUnsignedInteger(uint(u)), nil
	case apiValues.SINT, apiValues.INT, apiValues.DINT, apiValues.LINT:
		return model.CreateBACnetApplicationTagSignedInteger(int(v.GetInt64())), nil
	case apiValues.REAL:
		return model.CreateBACnetApplicationTagReal(v.GetFloat32()), nil
	case apiValues.LREAL:
		return model.CreateBACnetApplicationTagDouble(v.GetFloat64()), nil
	case apiValues.STRING, apiValues.CHAR, apiValues.WSTRING, apiValues.WCHAR:
		return model.CreateBACnetApplicationTagCharacterString(model.BACnetCharacterEncoding_ISO_10646, v.GetString()), nil
	case apiValues.RAW_BYTE_ARRAY:
		return model.CreateBACnetApplicationTagOctetString(v.GetRaw()), nil
	case apiValues.DATE:
		return dateToApplicationTag(v.GetDate())
	case apiValues.TIME_OF_DAY:
		return timeOfDayToApplicationTag(v.GetTime()), nil
	default:
		return nil, errors.Errorf("PlcValue type %v cannot be encoded as a BACnet ApplicationTag", v.GetPlcValueType())
	}
}

// dateToApplicationTag maps a calendar day to a BACnet Date (ASHRAE 135 clause
// 20.2.12), the inverse of bacnetDateToPlcValue for a concrete date. The day is
// taken from the value's own zone, as a PlcDATE carries no instant. The year octet
// counts from 1900 and 255 means "any year", so only 1900..2154 can be sent; any
// other year is an error instead of wrapping into a different one. The day of week
// is numbered 1=Monday .. 7=Sunday.
func dateToApplicationTag(date time.Time) (model.BACnetApplicationTag, error) {
	year := date.Year()
	if year < 1900 || year > 1900+0xFE {
		return nil, errors.Errorf("year %d of date %s is outside the BACnet date range 1900..2154", year, date.Format(time.DateOnly))
	}
	dayOfWeek := uint8(date.Weekday())
	if date.Weekday() == time.Sunday {
		dayOfWeek = 7
	}
	return model.CreateBACnetApplicationTagDate(uint16(year), uint8(date.Month()), uint8(date.Day()), dayOfWeek), nil
}

// timeOfDayToApplicationTag maps a time of day to a BACnet Time (clause 20.2.13),
// the inverse of bacnetTimeToPlcValue for a concrete time. The wire format holds
// hundredths of a second; a finer part is truncated, not rounded, so 23:59:59.999
// stays on the same day.
func timeOfDayToApplicationTag(timeOfDay time.Time) model.BACnetApplicationTag {
	return model.CreateBACnetApplicationTagTime(
		uint8(timeOfDay.Hour()),
		uint8(timeOfDay.Minute()),
		uint8(timeOfDay.Second()),
		uint8(timeOfDay.Nanosecond()/int(10*time.Millisecond)),
	)
}

// encodingHint nudges plcValueToApplicationTag when the source PlcValue type is
// ambiguous. Today the only hint is enumerated; future additions could include
// bit-string-from-uint or date-from-string.
type encodingHint uint8

const (
	hintNone encodingHint = iota
	hintEnumerated
)

// hintForProperty returns the encoding hint for a BACnet (objectType,
// propertyIdentifier, arrayIndex) triple. A PlcValue carries no BACnet
// datatype, and plc4go reads an Enumerated back as an unsigned integer, so
// without this hint an enumerated value written back goes out as Unsigned (or
// a bool as Boolean) and devices reject it with INVALID_DATA_TYPE or
// INVALID_TAG. The datatypes come from ASHRAE 135, Clause 12 (object types)
// and Clause 21 (BACnetPropertyIdentifier datatypes).
//
// Array index 0 is the array size (Unsigned) and gets no hint. Proprietary
// properties and object types are not in the tables and keep today's
// type-based encoding.
func hintForProperty(objectType, propertyId uint32, arrayIndex *uint) encodingHint {
	if arrayIndex != nil && *arrayIndex == 0 {
		return hintNone
	}
	prop := model.BACnetPropertyIdentifier(propertyId)
	if _, ok := enumeratedProperties[prop]; ok {
		return hintEnumerated
	}
	if objectTypes, ok := enumeratedPropertiesByObjectType[prop]; ok {
		if _, ok := objectTypes[model.BACnetObjectType(objectType)]; ok {
			return hintEnumerated
		}
	}
	return hintNone
}

// enumeratedProperties are standard properties whose datatype is Enumerated
// (or a list or array of Enumerated) on every object type that has them.
var enumeratedProperties = map[model.BACnetPropertyIdentifier]struct{}{
	model.BACnetPropertyIdentifier_ACCEPTED_MODES:                   {}, // BACnetLIST of BACnetLifeSafetyMode
	model.BACnetPropertyIdentifier_ACCESS_ALARM_EVENTS:              {}, // BACnetLIST of BACnetAccessEvent
	model.BACnetPropertyIdentifier_ACCESS_EVENT:                     {},
	model.BACnetPropertyIdentifier_ACCESS_TRANSACTION_EVENTS:        {}, // BACnetLIST of BACnetAccessEvent
	model.BACnetPropertyIdentifier_ALARM_VALUE:                      {}, // BACnetBinaryPV (BI, BV)
	model.BACnetPropertyIdentifier_AUTHENTICATION_STATUS:            {},
	model.BACnetPropertyIdentifier_AUTHORIZATION_EXEMPTIONS:         {}, // BACnetLIST of BACnetAuthorizationExemption
	model.BACnetPropertyIdentifier_AUTHORIZATION_MODE:               {},
	model.BACnetPropertyIdentifier_BACKUP_AND_RESTORE_STATE:         {},
	model.BACnetPropertyIdentifier_BACNET_IP_MODE:                   {},
	model.BACnetPropertyIdentifier_BACNET_IPV6_MODE:                 {},
	model.BACnetPropertyIdentifier_BASE_DEVICE_SECURITY_POLICY:      {}, // BACnetSecurityLevel
	model.BACnetPropertyIdentifier_CAR_ASSIGNED_DIRECTION:           {},
	model.BACnetPropertyIdentifier_CAR_DOOR_COMMAND:                 {}, // BACnetARRAY of BACnetLiftCarDoorCommand
	model.BACnetPropertyIdentifier_CAR_DOOR_STATUS:                  {}, // BACnetARRAY of BACnetDoorStatus
	model.BACnetPropertyIdentifier_CAR_DRIVE_STATUS:                 {},
	model.BACnetPropertyIdentifier_CAR_LOAD_UNITS:                   {}, // BACnetEngineeringUnits
	model.BACnetPropertyIdentifier_CAR_MODE:                         {},
	model.BACnetPropertyIdentifier_CAR_MOVING_DIRECTION:             {},
	model.BACnetPropertyIdentifier_COMMAND:                          {}, // BACnetNetworkPortCommand
	model.BACnetPropertyIdentifier_CONTROLLED_VARIABLE_UNITS:        {},
	model.BACnetPropertyIdentifier_CREDENTIAL_DISABLE:               {},
	model.BACnetPropertyIdentifier_CREDENTIAL_STATUS:                {}, // BACnetBinaryPV
	model.BACnetPropertyIdentifier_DEFAULT_SUBORDINATE_RELATIONSHIP: {},
	model.BACnetPropertyIdentifier_DERIVATIVE_CONSTANT_UNITS:        {},
	model.BACnetPropertyIdentifier_DOOR_ALARM_STATE:                 {},
	model.BACnetPropertyIdentifier_DOOR_STATUS:                      {},
	model.BACnetPropertyIdentifier_ESCALATOR_MODE:                   {},
	model.BACnetPropertyIdentifier_EVENT_STATE:                      {},
	model.BACnetPropertyIdentifier_EVENT_TYPE:                       {},
	model.BACnetPropertyIdentifier_FAILED_ATTEMPT_EVENTS:            {}, // BACnetLIST of BACnetAccessEvent
	model.BACnetPropertyIdentifier_FAULT_SIGNALS:                    {}, // BACnetLIST of BACnetLiftFault / BACnetEscalatorFault
	model.BACnetPropertyIdentifier_FAULT_TYPE:                       {},
	model.BACnetPropertyIdentifier_FILE_ACCESS_METHOD:               {},
	model.BACnetPropertyIdentifier_GROUP_MODE:                       {},
	model.BACnetPropertyIdentifier_IN_PROGRESS:                      {},
	model.BACnetPropertyIdentifier_INTEGRAL_CONSTANT_UNITS:          {},
	model.BACnetPropertyIdentifier_LAST_ACCESS_EVENT:                {},
	model.BACnetPropertyIdentifier_LAST_RESTART_REASON:              {},
	model.BACnetPropertyIdentifier_LAST_STATE_CHANGE:                {}, // BACnetTimerTransition
	model.BACnetPropertyIdentifier_LIFE_SAFETY_ALARM_VALUES:         {}, // BACnetLIST of BACnetLifeSafetyState
	model.BACnetPropertyIdentifier_LOCK_STATUS:                      {},
	model.BACnetPropertyIdentifier_LOGGING_TYPE:                     {},
	model.BACnetPropertyIdentifier_MASKED_ALARM_VALUES:              {}, // BACnetLIST of BACnetDoorAlarmState
	model.BACnetPropertyIdentifier_MODE:                             {}, // BACnetLifeSafetyMode
	model.BACnetPropertyIdentifier_NETWORK_NUMBER_QUALITY:           {},
	model.BACnetPropertyIdentifier_NETWORK_TYPE:                     {},
	model.BACnetPropertyIdentifier_NODE_TYPE:                        {},
	model.BACnetPropertyIdentifier_NOTIFY_TYPE:                      {},
	model.BACnetPropertyIdentifier_OBJECT_TYPE:                      {},
	model.BACnetPropertyIdentifier_OCCUPANCY_STATE:                  {},
	model.BACnetPropertyIdentifier_OPERATION_DIRECTION:              {},
	model.BACnetPropertyIdentifier_OPERATION_EXPECTED:               {},
	model.BACnetPropertyIdentifier_OUTPUT_UNITS:                     {},
	model.BACnetPropertyIdentifier_PASSBACK_MODE:                    {},
	model.BACnetPropertyIdentifier_POLARITY:                         {},
	model.BACnetPropertyIdentifier_PROGRAM_CHANGE:                   {},
	model.BACnetPropertyIdentifier_PROGRAM_STATE:                    {},
	model.BACnetPropertyIdentifier_PROPERTY_LIST:                    {}, // BACnetARRAY of BACnetPropertyIdentifier
	model.BACnetPropertyIdentifier_PROPORTIONAL_CONSTANT_UNITS:      {},
	model.BACnetPropertyIdentifier_PROTOCOL_LEVEL:                   {},
	model.BACnetPropertyIdentifier_REASON_FOR_DISABLE:               {}, // BACnetLIST of BACnetAccessCredentialDisableReason
	model.BACnetPropertyIdentifier_REASON_FOR_HALT:                  {},
	model.BACnetPropertyIdentifier_RELIABILITY:                      {},
	model.BACnetPropertyIdentifier_SECURED_STATUS:                   {},
	model.BACnetPropertyIdentifier_SEGMENTATION_SUPPORTED:           {},
	model.BACnetPropertyIdentifier_SILENCED:                         {},
	model.BACnetPropertyIdentifier_SUBORDINATE_NODE_TYPES:           {}, // BACnetARRAY of BACnetNodeType
	model.BACnetPropertyIdentifier_SUBORDINATE_RELATIONSHIPS:        {}, // BACnetARRAY of BACnetRelationship
	model.BACnetPropertyIdentifier_SYSTEM_STATUS:                    {},
	model.BACnetPropertyIdentifier_TIMER_STATE:                      {},
	model.BACnetPropertyIdentifier_TRANSITION:                       {}, // BACnetLightingTransition
	model.BACnetPropertyIdentifier_UNITS:                            {},
	model.BACnetPropertyIdentifier_USER_TYPE:                        {},
	model.BACnetPropertyIdentifier_VT_CLASSES_SUPPORTED:             {}, // BACnetLIST of BACnetVTClass
	model.BACnetPropertyIdentifier_WRITE_STATUS:                     {},
}

// enumeratedPropertiesByObjectType are standard properties whose datatype
// depends on the object type; they are Enumerated only on the listed ones.
// Present_Value is Unsigned on multi-state objects and REAL on analog ones,
// for example.
var enumeratedPropertiesByObjectType = map[model.BACnetPropertyIdentifier]map[model.BACnetObjectType]struct{}{
	model.BACnetPropertyIdentifier_PRESENT_VALUE: {
		model.BACnetObjectType_ACCESS_DOOR:            {}, // BACnetDoorValue
		model.BACnetObjectType_BINARY_INPUT:           {}, // BACnetBinaryPV
		model.BACnetObjectType_BINARY_LIGHTING_OUTPUT: {}, // BACnetBinaryLightingPV
		model.BACnetObjectType_BINARY_OUTPUT:          {},
		model.BACnetObjectType_BINARY_VALUE:           {},
		model.BACnetObjectType_LIFE_SAFETY_POINT:      {}, // BACnetLifeSafetyState
		model.BACnetObjectType_LIFE_SAFETY_ZONE:       {},
		model.BACnetObjectType_LOAD_CONTROL:           {}, // BACnetShedState
	},
	model.BACnetPropertyIdentifier_RELINQUISH_DEFAULT: {
		model.BACnetObjectType_ACCESS_DOOR:            {},
		model.BACnetObjectType_BINARY_LIGHTING_OUTPUT: {},
		model.BACnetObjectType_BINARY_OUTPUT:          {},
		model.BACnetObjectType_BINARY_VALUE:           {},
	},
	model.BACnetPropertyIdentifier_FEEDBACK_VALUE: {
		model.BACnetObjectType_BINARY_LIGHTING_OUTPUT: {},
		model.BACnetObjectType_BINARY_OUTPUT:          {},
	},
	model.BACnetPropertyIdentifier_INTERFACE_VALUE: { // BACnetOptionalBinaryPV
		model.BACnetObjectType_BINARY_INPUT:  {},
		model.BACnetObjectType_BINARY_OUTPUT: {},
	},
	model.BACnetPropertyIdentifier_TRACKING_VALUE: {
		model.BACnetObjectType_LIFE_SAFETY_POINT: {},
		model.BACnetObjectType_LIFE_SAFETY_ZONE:  {},
	},
	model.BACnetPropertyIdentifier_ACTION: { // the Command object's Action is a BACnetARRAY of BACnetActionList
		model.BACnetObjectType_LOOP: {},
	},
	model.BACnetPropertyIdentifier_ALARM_VALUES: { // BACnetLIST of the object's state enumeration
		model.BACnetObjectType_ACCESS_DOOR:       {},
		model.BACnetObjectType_ACCESS_ZONE:       {},
		model.BACnetObjectType_LIFE_SAFETY_POINT: {},
		model.BACnetObjectType_LIFE_SAFETY_ZONE:  {},
		model.BACnetObjectType_TIMER:             {},
	},
	model.BACnetPropertyIdentifier_MAINTENANCE_REQUIRED: { // BACnetMaintenance; BOOLEAN on a Life Safety Zone
		model.BACnetObjectType_ACCESS_DOOR:       {},
		model.BACnetObjectType_LIFE_SAFETY_POINT: {},
	},
	model.BACnetPropertyIdentifier_FAULT_VALUES: {
		model.BACnetObjectType_ACCESS_DOOR:       {},
		model.BACnetObjectType_LIFE_SAFETY_POINT: {},
		model.BACnetObjectType_LIFE_SAFETY_ZONE:  {},
	},
}
