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

package basetypes

import (
	"fmt"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type vendorRange struct {
	lower, upper int
}

//
//   Bit Strings
//

type AuditOperationFlags struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewAuditOperationFlags(arg Arg) (*AuditOperationFlags, error) {
	s := &AuditOperationFlags{
		vendorRange: vendorRange{32, 63},
		bitNames: map[string]int{
			"read":              0,
			"write":             1,
			"create":            2,
			"delete":            3,
			"lifeSafety":        4,
			"acknowledgeAlarm":  5,
			"deviceDisableComm": 6,
			"deviceEnableComm":  7,
			"deviceReset":       8,
			"deviceBackup":      9,
			"deviceRestore":     10,
			"subscription":      11,
			"notification":      12,
			"auditingFailure":   13,
			"networkChanges":    14,
			"general":           15,
		},
		bitLen: 16,
	}
	panic("implement me")
	return s, nil
}

type DaysOfWeek struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewDaysOfWeek(arg Arg) (*DaysOfWeek, error) {
	s := &DaysOfWeek{
		bitNames: map[string]int{
			"monday":    0,
			"tuesday":   1,
			"wednesday": 2,
			"thursday":  3,
			"friday":    4,
			"saturday":  5,
			"sunday":    6,
		},
		bitLen: 7,
	}
	panic("implement me")
	return s, nil
}

type EventTransitionBits struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewEventTransitionBits(arg Arg) (*EventTransitionBits, error) {
	s := &EventTransitionBits{
		bitNames: map[string]int{
			"toOffnormal": 0,
			"toFault":     1,
			"toNormal":    2,
		},
		bitLen: 3,
	}
	panic("implement me")
	return s, nil
}

type LimitEnable struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewLimitEnable(arg Arg) (*LimitEnable, error) {
	s := &LimitEnable{
		bitNames: map[string]int{
			"lowLimitEnable":  0,
			"highLimitEnable": 1,
		},
		bitLen: 2,
	}
	panic("implement me")
	return s, nil
}

type LogStatus struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewLogStatus(arg Arg) (*LogStatus, error) {
	s := &LogStatus{
		bitNames: map[string]int{
			"logDisabled":    0,
			"bufferPurged":   1,
			"logInterrupted": 2,
		},
		bitLen: 3,
	}
	panic("implement me")
	return s, nil
}

type ObjectTypesSupported struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewObjectTypesSupported(arg Arg) (*ObjectTypesSupported, error) {
	s := &ObjectTypesSupported{
		bitNames: map[string]int{
			"analogInput":           0,
			"analogOutput":          1,
			"analogValue":           2,
			"binaryInput":           3,
			"binaryOutput":          4,
			"binaryValue":           5,
			"calendar":              6,
			"command":               7,
			"device":                8,
			"eventEnrollment":       9,
			"file":                  10,
			"group":                 11,
			"loop":                  12,
			"multiStateInput":       13,
			"multiStateOutput":      14,
			"notificationClass":     15,
			"program":               16,
			"schedule":              17,
			"averaging":             18,
			"multiStateValue":       19,
			"trendLog":              20,
			"lifeSafetyPoint":       21,
			"lifeSafetyZone":        22,
			"accumulator":           23,
			"pulseConverter":        24,
			"eventLog":              25,
			"globalGroup":           26,
			"trendLogMultiple":      27,
			"loadControl":           28,
			"structuredView":        29,
			"accessDoor":            30,
			"accessCredential":      32,
			"accessPoint":           33,
			"accessRights":          34,
			"accessUser":            35,
			"accessZone":            36,
			"credentialDataInput":   37,
			"networkSecurity":       38, // removed revision 22
			"bitstringValue":        39,
			"characterstringValue":  40,
			"datePatternValue":      41,
			"dateValue":             42,
			"datetimePatternValue":  43,
			"datetimeValue":         44,
			"integerValue":          45,
			"largeAnalogValue":      46,
			"octetstringValue":      47,
			"positiveIntegerValue":  48,
			"timePatternValue":      49,
			"timeValue":             50,
			"notificationForwarder": 51,
			"alertEnrollment":       52,
			"channel":               53,
			"lightingOutput":        54,
			"binaryLightingOutput":  55,
			"networkPort":           56,
			"elevatorGroup":         57,
			"escalator":             58,
			"lift":                  59,
			"staging":               60,
			"auditLog":              61,
			"auditReporter":         62,
		},
		bitLen: 63,
	}
	panic("implement me")
	return s, nil
}

type PriorityFilter struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewPriorityFilter(arg Arg) (*PriorityFilter, error) {
	s := &PriorityFilter{
		bitNames: map[string]int{
			"manualLifeSafety":          0,
			"automaticLifeSafety":       1,
			"priority3":                 2,
			"priority4":                 3,
			"criticalEquipmentControls": 4,
			"minimumOnOff":              5,
			"priority7":                 6,
			"manualOperator":            7,
			"priority9":                 8,
			"priority10":                9,
			"priority11":                10,
			"priority12":                11,
			"priority13":                12,
			"priority14":                13,
			"priority15":                14,
			"priority16":                15,
		},
		bitLen: 16,
	}
	panic("implement me")
	return s, nil
}

type ResultFlags struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewResultFlags(arg Arg) (*ResultFlags, error) {
	s := &ResultFlags{
		bitNames: map[string]int{
			"firstItem": 0,
			"lastItem":  1,
			"moreItems": 2,
		},
		bitLen: 3,
	}
	panic("implement me")
	return s, nil
}

type ServicesSupported struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewServicesSupported(arg Arg) (*ServicesSupported, error) {
	s := &ServicesSupported{
		bitNames: map[string]int{
			"acknowledgeAlarm":           0,
			"confirmedCOVNotification":   1,
			"confirmedEventNotification": 2,
			"getAlarmSummary":            3,
			"getEnrollmentSummary":       4,
			"subscribeCOV":               5,
			"atomicReadFile":             6,
			"atomicWriteFile":            7,
			"addListElement":             8,
			"removeListElement":          9,
			"createObject":               10,
			"deleteObject":               11,
			"readProperty":               12,
			// "readPropertyConditional": 13      // removed in version 1 revision 12,
			"readPropertyMultiple":       14,
			"writeProperty":              15,
			"writePropertyMultiple":      16,
			"deviceCommunicationControl": 17,
			"confirmedPrivateTransfer":   18,
			"confirmedTextMessage":       19,
			"reinitializeDevice":         20,
			"vtOpen":                     21,
			"vtClose":                    22,
			"vtData":                     23,
			// "authenticate": 24                 // removed in version 1 revision 11,
			// "requestKey": 25                   // removed in version 1 revision 11,
			"iAm":                                26,
			"iHave":                              27,
			"unconfirmedCOVNotification":         28,
			"unconfirmedEventNotification":       29,
			"unconfirmedPrivateTransfer":         30,
			"unconfirmedTextMessage":             31,
			"timeSynchronization":                32,
			"whoHas":                             33,
			"whoIs":                              34,
			"readRange":                          35,
			"utcTimeSynchronization":             36,
			"lifeSafetyOperation":                37,
			"subscribeCOVProperty":               38,
			"getEventInformation":                39,
			"writeGroup":                         40,
			"subscribeCOVPropertyMultiple":       41,
			"confirmedCOVNotificationMultiple":   42,
			"unconfirmedCOVNotificationMultiple": 43,
			"confirmedAuditNotification":         44,
			"auditLogQuery":                      45,
			"unconfirmedAuditNotification":       46,
			"whoAmI":                             47,
			"youAre":                             48,
		},
		bitLen: 49,
	}
	panic("implement me")
	return s, nil
}

type StatusFlags struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewStatusFlags(arg Arg) (*StatusFlags, error) {
	s := &StatusFlags{
		bitNames: map[string]int{
			"inAlarm":      0,
			"fault":        1,
			"overridden":   2,
			"outOfService": 3,
		},
		bitLen: 4,
	}
	panic("implement me")
	return s, nil
}

//
//   Enumerations
//

type AccessAuthenticationFactorDisable struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessAuthenticationFactorDisable(arg Arg) (*AccessAuthenticationFactorDisable, error) {
	s := &AccessAuthenticationFactorDisable{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"none": 0,
			"disabled":          1,
			"disabledLost":      2,
			"disabledStolen":    3,
			"disabledDamaged":   4,
			"disabledDestroyed": 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessCredentialDisable struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessCredentialDisable(arg Arg) (*AccessCredentialDisable, error) {
	s := &AccessCredentialDisable{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"none": 0,
			"disable":        1,
			"disableManual":  2,
			"disableLockout": 3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessCredentialDisableReason struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessCredentialDisableReason(arg Arg) (*AccessCredentialDisableReason, error) {
	s := &AccessCredentialDisableReason{
		enumerations: map[string]uint64{"disabled": 0,
			"disabledNeedsProvisioning": 1,
			"disabledUnassigned":        2,
			"disabledNotYetActive":      3,
			"disabledExpired":           4,
			"disabledLockout":           5,
			"disabledMaxDays":           6,
			"disabledMaxUses":           7,
			"disabledInactivity":        8,
			"disabledManual":            9,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessEvent struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessEvent(arg Arg) (*AccessEvent, error) {
	s := &AccessEvent{
		vendorRange: vendorRange{512, 65535},
		enumerations: map[string]uint64{"none": 0,
			"granted":                             1,
			"muster":                              2,
			"passbackDetected":                    3,
			"duress":                              4,
			"trace":                               5,
			"lockoutMaxAttempts":                  6,
			"lockoutOther":                        7,
			"lockoutRelinquished":                 8,
			"lockedByHigherPriority":              9,
			"outOfService":                        10,
			"outOfServiceRelinquished":            11,
			"accompanimentBy":                     12,
			"authenticationFactorRead":            13,
			"authorizationDelayed":                14,
			"verificationRequired":                15,
			"deniedDenyAll":                       128,
			"deniedUnknownCredential":             129,
			"deniedAuthenticationUnavailable":     130,
			"deniedAuthenticationFactorTimeout":   131,
			"deniedIncorrectAuthenticationFactor": 132,
			"deniedZoneNoAccessRights":            133,
			"deniedPointNoAccessRights":           134,
			"deniedNoAccessRights":                135,
			"deniedOutOfTimeRange":                136,
			"deniedThreatLevel":                   137,
			"deniedPassback":                      138,
			"deniedUnexpectedLocationUsage":       139,
			"deniedMaxAttempts":                   140,
			"deniedLowerOccupancyLimit":           141,
			"deniedUpperOccupancyLimit":           142,
			"deniedAuthenticationFactorLost":      143,
			"deniedAuthenticationFactorStolen":    144,
			"deniedAuthenticationFactorDamaged":   145,
			"deniedAuthenticationFactorDestroyed": 146,
			"deniedAuthenticationFactorDisabled":  147,
			"deniedAuthenticationFactorError":     148,
			"deniedCredentialUnassigned":          149,
			"deniedCredentialNotProvisioned":      150,
			"deniedCredentialNotYetActive":        151,
			"deniedCredentialExpired":             152,
			"deniedCredentialManualDisable":       153,
			"deniedCredentialLockout":             154,
			"deniedCredentialMaxDays":             155,
			"deniedCredentialMaxUses":             156,
			"deniedCredentialInactivity":          157,
			"deniedCredentialDisabled":            158,
			"deniedNoAccompaniment":               159,
			"deniedIncorrectAccompaniment":        160,
			"deniedLockout":                       161,
			"deniedVerificationFailed":            162,
			"deniedVerificationTimeout":           163,
			"deniedOther":                         164,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessPassbackMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessPassbackMode(arg Arg) (*AccessPassbackMode, error) {
	s := &AccessPassbackMode{
		enumerations: map[string]uint64{"passbackOff": 0,
			"hardPassback": 1,
			"softPassback": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessRuleTimeRangeSpecifier struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessRuleTimeRangeSpecifier(arg Arg) (*AccessRuleTimeRangeSpecifier, error) {
	s := &AccessRuleTimeRangeSpecifier{
		enumerations: map[string]uint64{"specified": 0,
			"always": 1,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessRuleLocationSpecifier struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessRuleLocationSpecifier(arg Arg) (*AccessRuleLocationSpecifier, error) {
	s := &AccessRuleLocationSpecifier{
		enumerations: map[string]uint64{"specified": 0,
			"all": 1,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessUserType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessUserType(arg Arg) (*AccessUserType, error) {
	s := &AccessUserType{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"asset": 0,
			"group":  1,
			"person": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccessZoneOccupancyState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessZoneOccupancyState(arg Arg) (*AccessZoneOccupancyState, error) {
	s := &AccessZoneOccupancyState{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"normal": 0,
			"belowLowerLimit": 1,
			"atLowerLimit":    2,
			"atUpperLimit":    3,
			"aboveUpperLimit": 4,
			"disabled":        5,
			"notSupported":    6,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AccumulatorRecordAccumulatorStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccumulatorRecordAccumulatorStatus(arg Arg) (*AccumulatorRecordAccumulatorStatus, error) {
	s := &AccumulatorRecordAccumulatorStatus{
		enumerations: map[string]uint64{"normal": 0,
			"starting":  1,
			"recovered": 2,
			"abnormal":  3,
			"failed":    4,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type Action struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAction(arg Arg) (*Action, error) {
	s := &Action{
		enumerations: map[string]uint64{"direct": 0,
			"reverse": 1,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AuditLevel struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuditLevel(arg Arg) (*AuditLevel, error) {
	s := &AuditLevel{
		vendorRange: vendorRange{128, 255},
		enumerations: map[string]uint64{"none": 0,
			"auditAll":    1,
			"auditConfig": 2,
			"default":     3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AuditOperation struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuditOperation(arg Arg) (*AuditOperation, error) {
	s := &AuditOperation{
		vendorRange: vendorRange{32, 63},
		enumerations: map[string]uint64{"read": 0,
			"write":             1,
			"create":            2,
			"delete":            3,
			"lifeSafety":        4,
			"acknowledgeAlarm":  5,
			"deviceDisableComm": 6,
			"deviceEnableComm":  7,
			"deviceReset":       8,
			"deviceBackup":      9,
			"deviceRestore":     10,
			"subscription":      11,
			"notification":      12,
			"auditingFailure":   13,
			"networkChanges":    14,
			"general":           15,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AuthenticationFactorType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuthenticationFactorType(arg Arg) (*AuthenticationFactorType, error) {
	s := &AuthenticationFactorType{
		enumerations: map[string]uint64{"undefined": 0,
			"error":              1,
			"custom":             2,
			"simpleNumber16":     3,
			"simpleNumber32":     4,
			"simpleNumber56":     5,
			"simpleAlphaNumeric": 6,
			"abaTrack2":          7,
			"wiegand26":          8,
			"wiegand37":          9,
			"wiegand37facility":  10,
			"facility16card32":   11,
			"facility32card32":   12,
			"fascN":              13,
			"fascNbcd":           14,
			"fascNlarge":         15,
			"fascNlargeBcd":      16,
			"gsa75":              17,
			"chuid":              18,
			"chuidFull":          19,
			"guid":               20,
			"cbeffA":             21,
			"cbeffB":             22,
			"cbeffC":             23,
			"userPassword":       24,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AuthenticationStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuthenticationStatus(arg Arg) (*AuthenticationStatus, error) {
	s := &AuthenticationStatus{
		enumerations: map[string]uint64{"notReady": 0,
			"ready":                          1,
			"disabled":                       2,
			"waitingForAuthenticationFactor": 3,
			"waitingForAccompaniment":        4,
			"waitingForVerification":         5,
			"inProgress":                     6,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AuthorizationException struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuthorizationException(arg Arg) (*AuthorizationException, error) {
	s := &AuthorizationException{
		vendorRange: vendorRange{64, 255},
		enumerations: map[string]uint64{"passback": 0,
			"occupancyCheck":     1,
			"accessRights":       2,
			"lockout":            3,
			"deny":               4,
			"verification":       5,
			"authorizationDelay": 6,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type AuthorizationMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuthorizationMode(arg Arg) (*AuthorizationMode, error) {
	s := &AuthorizationMode{
		vendorRange: vendorRange{64, 65536},
		enumerations: map[string]uint64{"authorize": 0,
			"grantActive":          1,
			"denyAll":              2,
			"verificationRequired": 3,
			"authorizationDelayed": 4,
			"none":                 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type BackupState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewBackupState(arg Arg) (*BackupState, error) {
	s := &BackupState{
		enumerations: map[string]uint64{"idle": 0,
			"preparingForBackup":  1,
			"preparingForRestore": 2,
			"performingABackup":   3,
			"performingARestore":  4,
			"backupFailure":       5,
			"restoreFailure":      6,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type BinaryLightingPV struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewBinaryLightingPV(arg Arg) (*BinaryLightingPV, error) {
	s := &BinaryLightingPV{
		vendorRange: vendorRange{64, 255},
		enumerations: map[string]uint64{"off": 0,
			"on":             1,
			"warn":           2,
			"warnOff":        3,
			"warnRelinquish": 4,
			"stop":           5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type BinaryPV struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewBinaryPV(arg Arg) (*BinaryPV, error) {
	s := &BinaryPV{
		enumerations: map[string]uint64{"inactive": 0,
			"active": 1,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type DeviceStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewDeviceStatus(arg Arg) (*DeviceStatus, error) {
	s := &DeviceStatus{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"operational": 0,
			"operationalReadOnly": 1,
			"downloadRequired":    2,
			"downloadInProgress":  3,
			"nonOperational":      4,
			"backupInProgress":    5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type DoorAlarmState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewDoorAlarmState(arg Arg) (*DoorAlarmState, error) {
	s := &DoorAlarmState{
		vendorRange: vendorRange{256, 65535},
		enumerations: map[string]uint64{"normal": 0,
			"alarm":           1,
			"doorOpenTooLong": 2,
			"forcedOpen":      3,
			"tamper":          4,
			"doorFault":       5,
			"lockDown":        6,
			"freeAccess":      7,
			"egressOpen":      8,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type DoorSecuredStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewDoorSecuredStatus(arg Arg) (*DoorSecuredStatus, error) {
	s := &DoorSecuredStatus{
		enumerations: map[string]uint64{"secured": 0,
			"unsecured": 1,
			"unknown":   2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type DoorStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewDoorStatus(arg Arg) (*DoorStatus, error) {
	s := &DoorStatus{
		enumerations: map[string]uint64{"closed": 0,
			"opened":  1,
			"unknown": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type DoorValue struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewDoorValue(arg Arg) (*DoorValue, error) {
	s := &DoorValue{
		enumerations: map[string]uint64{"lock": 0,
			"unlock":              1,
			"pulseUnlock":         2,
			"extendedPulseUnlock": 3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type EngineeringUnits struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewEngineeringUnits(arg Arg) (*EngineeringUnits, error) {
	s := &EngineeringUnits{
		vendorRange: vendorRange{256, 65535},
		enumerations: map[string]uint64{
			//Acceleration
			"metersPerSecondPerSecond": 166,
			"squareMeters":             0,
			"squareCentimeters":        116,
			"squareFeet":               1,
			"squareInches":             115,
			//Currency
			"currency1":  105,
			"currency2":  106,
			"currency3":  107,
			"currency4":  108,
			"currency5":  109,
			"currency6":  110,
			"currency7":  111,
			"currency8":  112,
			"currency9":  113,
			"currency10": 114,
			//Electrical
			"milliamperes":                2,
			"amperes":                     3,
			"amperesPerMeter":             167,
			"amperesPerSquareMeter":       168,
			"ampereSquareMeters":          169,
			"decibels":                    199,
			"decibelsMillivolt":           200,
			"decibelsVolt":                201,
			"farads":                      170,
			"henrys":                      171,
			"ohms":                        4,
			"ohmMeters":                   172,
			"ohmMeterPerSquareMeter":      237,
			"milliohms":                   145,
			"kilohms":                     122,
			"megohms":                     123,
			"microSiemens":                190,
			"millisiemens":                202,
			"siemens":                     173,
			"siemensPerMeter":             174,
			"teslas":                      175,
			"volts":                       5,
			"millivolts":                  124,
			"kilovolts":                   6,
			"megavolts":                   7,
			"voltAmperes":                 8,
			"kilovoltAmperes":             9,
			"megavoltAmperes":             10,
			"ampereSeconds":               238,
			"ampereSquareHours":           246,
			"voltAmpereHours":             239, //VAh
			"kilovoltAmpereHours":         240, //kVAh
			"megavoltAmpereHours":         241, //MVAh
			"voltAmperesReactive":         11,
			"kilovoltAmperesReactive":     12,
			"megavoltAmperesReactive":     13,
			"voltAmpereHoursReactive":     242, //varh
			"kilovoltAmpereHoursReactive": 243, //kvarh
			"megavoltAmpereHoursReactive": 244, //Mvarh
			"voltsPerDegreeKelvin":        176,
			"voltsPerMeter":               177,
			"voltsSquareHours":            245,
			"degreesPhase":                14,
			"powerFactor":                 15,
			"webers":                      178,
			// Energy
			"joules":                16,
			"kilojoules":            17,
			"kilojoulesPerKilogram": 125,
			"megajoules":            126,
			"joulesPerHours":        247,
			"wattHours":             18,
			"kilowattHours":         19,
			"megawattHours":         146,
			"wattHoursReactive":     203,
			"kilowattHoursReactive": 204,
			"megawattHoursReactive": 205,
			"btus":                  20,
			"kiloBtus":              147,
			"megaBtus":              148,
			"therms":                21,
			"tonHours":              22,
			// Enthalpy
			"joulesPerKilogramDryAir":     23,
			"kilojoulesPerKilogramDryAir": 149,
			"megajoulesPerKilogramDryAir": 150,
			"btusPerPoundDryAir":          24,
			"btusPerPound":                117,
			"joulesPerDegreeKelvin":       127,
			// Entropy
			"kilojoulesPerDegreeKelvin":     151,
			"megajoulesPerDegreeKelvin":     152,
			"joulesPerKilogramDegreeKelvin": 128,
			// Force
			"newton": 153,
			// Frequency
			"cyclesPerHour":                   25,
			"cyclesPerMinute":                 26,
			"hertz":                           27,
			"kilohertz":                       129,
			"megahertz":                       130,
			"perHour":                         131,
			"gramsOfWaterPerKilogramDryAir":   28,
			"percentRelativeHumidity":         29,
			"micrometers":                     194,
			"millimeters":                     30,
			"centimeters":                     118,
			"kilometers":                      193,
			"meters":                          31,
			"inches":                          32,
			"feet":                            33,
			"candelas":                        179,
			"candelasPerSquareMeter":          180,
			"wattsPerSquareFoot":              34,
			"wattsPerSquareMeter":             35,
			"lumens":                          36,
			"luxes":                           37,
			"footCandles":                     38,
			"milligrams":                      196,
			"grams":                           195,
			"kilograms":                       39,
			"poundsMass":                      40,
			"tons":                            41,
			"gramsPerSecond":                  154,
			"gramsPerMinute":                  155,
			"kilogramsPerSecond":              42,
			"kilogramsPerMinute":              43,
			"kilogramsPerHour":                44,
			"poundsMassPerSecond":             119,
			"poundsMassPerMinute":             45,
			"poundsMassPerHour":               46,
			"tonsPerHour":                     156,
			"milliwatts":                      132,
			"watts":                           47,
			"kilowatts":                       48,
			"megawatts":                       49,
			"btusPerHour":                     50,
			"kiloBtusPerHour":                 157,
			"horsepower":                      51,
			"tonsRefrigeration":               52,
			"pascals":                         53,
			"hectopascals":                    133,
			"kilopascals":                     54,
			"pascalSeconds":                   253,
			"millibars":                       134,
			"bars":                            55,
			"poundsForcePerSquareInch":        56,
			"millimetersOfWater":              206,
			"centimetersOfWater":              57,
			"inchesOfWater":                   58,
			"millimetersOfMercury":            59,
			"centimetersOfMercury":            60,
			"inchesOfMercury":                 61,
			"degreesCelsius":                  62,
			"degreesKelvin":                   63,
			"degreesKelvinPerHour":            181,
			"degreesKelvinPerMinute":          182,
			"degreesFahrenheit":               64,
			"degreeDaysCelsius":               65,
			"degreeDaysFahrenheit":            66,
			"deltaDegreesFahrenheit":          120,
			"deltaDegreesKelvin":              121,
			"years":                           67,
			"months":                          68,
			"weeks":                           69,
			"days":                            70,
			"hours":                           71,
			"minutes":                         72,
			"seconds":                         73,
			"hundredthsSeconds":               158,
			"milliseconds":                    159,
			"newtonMeters":                    160,
			"millimetersPerSecond":            161,
			"millimetersPerMinute":            162,
			"metersPerSecond":                 74,
			"metersPerMinute":                 163,
			"metersPerHour":                   164,
			"kilometersPerHour":               75,
			"feetPerSecond":                   76,
			"feetPerMinute":                   77,
			"milesPerHour":                    78,
			"cubicFeet":                       79,
			"cubicFeetPerDay":                 248,
			"cubicMeters":                     80,
			"cubicMetersPerDay":               249,
			"imperialGallons":                 81,
			"milliliters":                     197,
			"liters":                          82,
			"usGallons":                       83,
			"cubicFeetPerSecond":              142,
			"cubicFeetPerMinute":              84,
			"cubicFeetPerHour":                191,
			"cubicMetersPerSecond":            85,
			"cubicMetersPerMinute":            165,
			"cubicMetersPerHour":              135,
			"imperialGallonsPerMinute":        86,
			"millilitersPerSecond":            198,
			"litersPerSecond":                 87,
			"litersPerMinute":                 88,
			"litersPerHour":                   136,
			"usGallonsPerMinute":              89,
			"usGallonsPerHour":                192,
			"degreesAngular":                  90,
			"degreesCelsiusPerHour":           91,
			"degreesCelsiusPerMinute":         92,
			"degreesFahrenheitPerHour":        93,
			"degreesFahrenheitPerMinute":      94,
			"jouleSeconds":                    183,
			"kilogramsPerCubicMeter":          186,
			"kilowattHoursPerSquareMeter":     137,
			"kilowattHoursPerSquareFoot":      138,
			"megajoulesPerSquareMeter":        139,
			"megajoulesPerSquareFoot":         140,
			"noUnits":                         95,
			"newtonSeconds":                   187,
			"newtonsPerMeter":                 188,
			"partsPerMillion":                 96,
			"partsPerBillion":                 97,
			"percent":                         98,
			"percentObscurationPerFoot":       143,
			"percentObscurationPerMeter":      144,
			"percentPerSecond":                99,
			"perMinute":                       100,
			"perSecond":                       101,
			"psiPerDegreeFahrenheit":          102,
			"radians":                         103,
			"radiansPerSecond":                184,
			"revolutionsPerMinute":            104,
			"squareMetersPerNewton":           185,
			"wattsPerMeterPerDegreeKelvin":    189,
			"wattsPerSquareMeterDegreeKelvin": 141,
			"perMille":                        207,
			"gramsPerGram":                    208,
			"kilogramsPerKilogram":            209,
			"gramsPerKilogram":                210,
			"milligramsPerGram":               211,
			"milligramsPerKilogram":           212,
			"gramsPerMilliliter":              213,
			"gramsPerLiter":                   214,
			"milligramsPerLiter":              215,
			"microgramsPerLiter":              216,
			"gramsPerCubicMeter":              217,
			"milligramsPerCubicMeter":         218,
			"microgramsPerCubicMeter":         219,
			"nanogramsPerCubicMeter":          220,
			"gramsPerCubicCentimeter":         221,
			"wattHoursPerCubicMeter":          250,
			"joulesPerCubicMeter":             251,
			"becquerels":                      222,
			"kilobecquerels":                  223,
			"megabecquerels":                  224,
			"gray":                            225,
			"milligray":                       226,
			"microgray":                       227,
			"sieverts":                        228,
			"millisieverts":                   229,
			"microsieverts":                   230,
			"microsievertsPerHour":            231,
			"decibelsA":                       232,
			"nephelometricTurbidityUnit":      233,
			"pH":                              234,
			"gramsPerSquareMeter":             235,
			"minutesPerDegreeKelvin":          236,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ErrorClass struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewErrorClass(arg Arg) (*ErrorClass, error) {
	s := &ErrorClass{
		enumerations: map[string]uint64{"device": 0,
			"object":        1,
			"property":      2,
			"resources":     3,
			"security":      4,
			"services":      5,
			"vt":            6,
			"communication": 7,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ErrorCode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewErrorCode(arg Arg) (*ErrorCode, error) {
	s := &ErrorCode{
		enumerations: map[string]uint64{"abortApduTooLong": 123,
			"abortApplicationExceededReplyTime":  124,
			"abortBufferOverflow":                51,
			"abortInsufficientSecurity":          135,
			"abortInvalidApduInThisState":        52,
			"abortOther":                         56,
			"abortOutOfResources":                125,
			"abortPreemptedByHigherPriorityTask": 53,
			"abortProprietary":                   55,
			"abortSecurityError":                 136,
			"abortSegmentationNotSupported":      54,
			"abortTsmTimeout":                    126,
			"abortWindowSizeOutOfRange":          127,
			"accessDenied":                       85,
			"addressingError":                    115,
			"badDestinationAddress":              86,
			"badDestinationDeviceId":             87,
			"badSignature":                       88,
			"badSourceAddress":                   89,
			"badTimestamp":                       90, //Removed revision 22
			"busy":                               82,
			"bvlcFunctionUnknown":                143,
			"bvlcProprietaryFunctionUnknown":     144,
			"cannotUseKey":                       91, //Removed revision 22
			"cannotVerifyMessageId":              92, // Removed revision 22
			"characterSetNotSupported":           41,
			"communicationDisabled":              83,
			"configurationInProgress":            2,
			"correctKeyRevision":                 93, // Removed revision 22
			"covSubscriptionFailed":              43,
			"datatypeNotSupported":               47,
			"deleteFdtEntryFailed":               120,
			"deviceBusy":                         3,
			"destinationDeviceIdRequired":        94, // Removed revision 22
			"distributeBroadcastFailed":          121,
			"dnsError":                           192,
			"dnsNameResolutionFailed":            190,
			"dnsResolverFailure":                 191,
			"dnsUnavailable":                     189,
			"duplicateEntry":                     137,
			"duplicateMessage":                   95,
			"duplicateName":                      48,
			"duplicateObjectId":                  49,
			"dynamicCreationNotSupported":        4,
			"encryptionNotConfigured":            96,
			"encryptionRequired":                 97,
			"fileAccessDenied":                   5,
			"fileFull":                           128,
			"headerEncodingError":                145,
			"headerNotUnderstood":                146,
			"httpError":                          165,
			"httpNoUpgrade":                      153,
			"httpNotAServer":                     164,
			"httpResourceNotLocal":               154,
			"httpProxyAuthenticationFailed":      155,
			"httpResponseTimeout":                156,
			"httpResponseSyntaxError":            157,
			"httpResponseValueError":             158,
			"httpResponseMissingHeader":          159,
			"httpTemporaryUnavailable":           163,
			"httpUnexpectedResponseCode":         152,
			"httpUpgradeRequired":                161,
			"httpUpgradeError":                   162,
			"httpWebsocketHeaderError":           160,
			"inconsistentConfiguration":          129,
			"inconsistentObjectType":             130,
			"inconsistentParameters":             7,
			"inconsistentSelectionCriterion":     8,
			"incorrectKey":                       98, // Removed revision 22
			"internalError":                      131,
			"invalidArrayIndex":                  42,
			"invalidConfigurationData":           46,
			"invalidDataEncoding":                142,
			"invalidDataType":                    9,
			"invalidEventState":                  73,
			"invalidFileAccessMethod":            10,
			"invalidFileStartPosition":           11,
			"invalidKeyData":                     99, // Removed revision 22
			"invalidParameterDataType":           13,
			"invalidTag":                         57,
			"invalidTimeStamp":                   14,
			"invalidValueInThisState":            138,
			"ipAddressNotReachable":              198,
			"ipError":                            199,
			"keyUpdateInProgress":                100, // Removed revision 22
			"listElementNotFound":                81,
			"listItemNotNumbered":                140,
			"listItemNotTimestamped":             141,
			"logBufferFull":                      75,
			"loggedValuePurged":                  76,
			"malformedMessage":                   101,
			"messageTooLong":                     113,
			"messageIncomplete":                  147,
			"missingRequiredParameter":           16,
			"networkDown":                        58,
			"noAlarmConfigured":                  74,
			"noObjectsOfSpecifiedType":           17,
			"noPropertySpecified":                77,
			"noSpaceForObject":                   18,
			"noSpaceToAddListElement":            19,
			"noSpaceToWriteProperty":             20,
			"noVtSessionsAvailable":              21,
			"nodeDuplicateVmac":                  151,
			"notABacnetScHub":                    148,
			"notConfigured":                      132,
			"notConfiguredForTriggeredLogging":   78,
			"notCovProperty":                     44,
			"notKeyServer":                       102, // Removed revision 22
			"notRouterToDnet":                    110,
			"objectDeletionNotPermitted":         23,
			"objectIdentifierAlreadyExists":      24,
			"other":                              0,
			"operationalProblem":                 25,
			"optionalFunctionalityNotSupported":  45,
			"outOfMemory":                        133,
			"parameterOutOfRange":                80,
			"passwordFailure":                    26,
			"payloadExpected":                    149,
			"propertyIsNotAList":                 22,
			"propertyIsNotAnArray":               50,
			"readAccessDenied":                   27,
			"readBdtFailed":                      117,
			"readFdtFailed":                      119,
			"registerForeignDeviceFailed":        118,
			"rejectBufferOverflow":               59,
			"rejectInconsistentParameters":       60,
			"rejectInvalidParameterDataType":     61,
			"rejectInvalidTag":                   62,
			"rejectMissingRequiredParameter":     63,
			"rejectParameterOutOfRange":          64,
			"rejectTooManyArguments":             65,
			"rejectUndefinedEnumeration":         66,
			"rejectUnrecognizedService":          67,
			"rejectProprietary":                  68,
			"rejectOther":                        69,
			"routerBusy":                         111,
			"securityError":                      114,
			"securityNotConfigured":              103,
			"serviceRequestDenied":               29,
			"sourceSecurityRequired":             104,
			"success":                            84,
			"tcpClosedByLocal":                   195,
			"tcpClosedOther":                     196,
			"tcpConnectTimeout":                  193,
			"tcpConnectionRefused":               194,
			"tcpError":                           197,
			"timeout":                            30,
			"tlsClientAuthenticationFailed":      182,
			"tlsClientCertificateError":          180,
			"tlsClientCertificateExpired":        184,
			"tlsClientCertificateRevoked":        186,
			"tleError":                           188,
			"tlsServerAuthenticationFailed":      183,
			"tlsServerCertificateError":          181,
			"tlsServerCertificateExpired":        185,
			"tlsServerCertificateRevoked":        187,
			"tooManyKeys":                        105, // Removed revision 22
			"unexpectedData":                     150,
			"unknownAuthenticationType":          106,
			"unknownDevice":                      70,
			"unknownFileSize":                    122,
			"unknownKey":                         107, // Removed revision 22
			"unknownKeyRevision":                 108, // Removed revision 22
			"unknownNetworkMessage":              112,
			"unknownObject":                      31,
			"unknownProperty":                    32,
			"unknownSubscription":                79,
			"umknownRoute":                       71,
			"unknownSourceMessage":               109, // Removed revision 22
			"unknownVtClass":                     34,
			"unknownVtSession":                   35,
			"unsupportedObjectType":              36,
			"valueNotInitialized":                72,
			"valueOutOfRange":                    37,
			"valueTooLong":                       134,
			"vtSessionAlreadyClosed":             38,
			"vtSessionTerminationFailure":        39,
			"websocket-close-error":              168,
			"websocket-closed-abnormally":        173,
			"websocket-closed-by-peer":           169,
			"websocket-data-against-policy":      175,
			"websocket-data-inconsistent":        174,
			"websocket-data-not-accepted":        172,
			"websocket-endpoint-leaves":          170,
			"websocket-error":                    179,
			"websocket-extension-missing":        177,
			"websocket-frame-too-long":           176,
			"websocket-protocol-error":           171,
			"websocket-request-unavailable":      178,
			"websocket-scheme-not-supported":     166,
			"websocket-unknown-control-message":  167,
			"writeAccessDenied":                  40,
			"writeBdtFailed":                     116,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type EscalatorMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewEscalatorMode(arg Arg) (*EscalatorMode, error) {
	s := &EscalatorMode{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"stop":         1,
			"up":           2,
			"down":         3,
			"inspection":   4,
			"outOfService": 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type EscalatorOperationDirection struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewEscalatorOperationDirection(arg Arg) (*EscalatorOperationDirection, error) {
	s := &EscalatorOperationDirection{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"stopped":          1,
			"upRatedSpeed":     2,
			"upReducedSpeed":   3,
			"downRatedSpeed":   4,
			"downReducedSpeed": 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type EventState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewEventState(arg Arg) (*EventState, error) {
	s := &EventState{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"normal": 0,
			"fault":           1,
			"offnormal":       2,
			"highLimit":       3,
			"lowLimit":        4,
			"lifeSafetyAlarm": 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type EventType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewEventType(arg Arg) (*EventType, error) {
	s := &EventType{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"changeOfBitstring": 0,
			"changeOfState":  1,
			"changeOfValue":  2,
			"commandFailure": 3,
			"floatingLimit":  4,
			"outOfRange":     5,
			// -- context tag 7 is deprecated
			"changeOfLifeSafety": 8,
			"extended":           9,
			"bufferReady":        10,
			"unsignedRange":      11,
			// -- enumeration value 12 is reserved for future addenda
			"accessEvent":             13,
			"doubleOutOfRange":        14,
			"signedOutOfRange":        15,
			"unsignedOutOfRange":      16,
			"changeOfCharacterstring": 17,
			"changeOfStatusFlags":     18,
			"changeOfReliability":     19,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type FaultType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewFaultType(arg Arg) (*FaultType, error) {
	s := &FaultType{
		enumerations: map[string]uint64{"none": 0,
			"faultCharacterstring": 1,
			"faultExtended":        2,
			"faultLifeSafety":      3,
			"faultState":           4,
			"faultStatusFlags":     5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type FileAccessMethod struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewFileAccessMethod(arg Arg) (*FileAccessMethod, error) {
	s := &FileAccessMethod{
		enumerations: map[string]uint64{"recordAccess": 0,
			"streamAccess": 1,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LiftCarDirection struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLiftCarDirection(arg Arg) (*LiftCarDirection, error) {
	s := &LiftCarDirection{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"none":      1,
			"stopped":   2,
			"up":        3,
			"down":      4,
			"upAndDown": 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LiftCarDoorCommand struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLiftCarDoorCommand(arg Arg) (*LiftCarDoorCommand, error) {
	s := &LiftCarDoorCommand{
		enumerations: map[string]uint64{"none": 0,
			"open":  1,
			"close": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LiftCarDriveStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLiftCarDriveStatus(arg Arg) (*LiftCarDriveStatus, error) {
	s := &LiftCarDriveStatus{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"stationary":      1,
			"braking":         2,
			"accelerate":      3,
			"decelerate":      4,
			"ratedSpeed":      5,
			"singleFloorJump": 6,
			"twoFloorJump":    7,
			"threeFloorJump":  8,
			"multiFloorJump":  9,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LiftCarMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLiftCarMode(arg Arg) (*LiftCarMode, error) {
	s := &LiftCarMode{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"normal":              1,
			"vip":                 2,
			"homing":              3,
			"parking":             4,
			"attendantControl":    5,
			"firefighterControl":  6,
			"emergencyPower":      7,
			"inspection":          8,
			"cabinetRecall":       9,
			"earthquakeOperation": 10,
			"fireOperation":       11,
			"outOfService":        12,
			"occupantEvacuation":  13,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LiftFault struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLiftFault(arg Arg) (*LiftFault, error) {
	s := &LiftFault{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"controllerFault": 0,
			"driveAndMotorFault":           1,
			"governorAndSafetyGearFault":   2,
			"liftShaftDeviceFault":         3,
			"powerSupplyFault":             4,
			"safetyInterlockFault":         5,
			"doorClosingFault":             6,
			"doorOpeningFault":             7,
			"carStoppedOutsideLandingZone": 8,
			"callButtonStuck":              9,
			"startFailure":                 10,
			"controllerSupplyFault":        11,
			"sTestFailure":                 12,
			"runtimeLimitExceeded":         13,
			"positionLost":                 14,
			"driveTemperatureExceeded":     15,
			"loadMeasurementFault":         16,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LiftGroupMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLiftGroupMode(arg Arg) (*LiftGroupMode, error) {
	s := &LiftGroupMode{
		enumerations: map[string]uint64{"unknown": 0,
			"normal":         1,
			"downPeak":       2,
			"twoWay":         3,
			"fourWay":        4,
			"emergencyPower": 5,
			"upPeak":         6,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LifeSafetyMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLifeSafetyMode(arg Arg) (*LifeSafetyMode, error) {
	s := &LifeSafetyMode{
		enumerations: map[string]uint64{"off": 0,
			"on":                       1,
			"test":                     2,
			"manned":                   3,
			"unmanned":                 4,
			"armed":                    5,
			"disarmed":                 6,
			"prearmed":                 7,
			"slow":                     8,
			"fast":                     9,
			"disconnected":             10,
			"enabled":                  11,
			"disabled":                 12,
			"automaticReleaseDisabled": 13,
			"default":                  14,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LifeSafetyOperation struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLifeSafetyOperation(arg Arg) (*LifeSafetyOperation, error) {
	s := &LifeSafetyOperation{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"none": 0,
			"silence":          1,
			"silenceAudible":   2,
			"silenceVisual":    3,
			"reset":            4,
			"resetAlarm":       5,
			"resetFault":       6,
			"unsilence":        7,
			"unsilenceAudible": 8,
			"unsilenceVisual":  9,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LifeSafetyState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLifeSafetyState(arg Arg) (*LifeSafetyState, error) {
	s := &LifeSafetyState{
		enumerations: map[string]uint64{"quiet": 0,
			"preAlarm":        1,
			"alarm":           2,
			"fault":           3,
			"faultPreAlarm":   4,
			"faultAlarm":      5,
			"notReady":        6,
			"active":          7,
			"tamper":          8,
			"testAlarm":       9,
			"testActive":      10,
			"testFault":       11,
			"testFaultAlarm":  12,
			"holdup":          13,
			"duress":          14,
			"tamperAlarm":     15,
			"abnormal":        16,
			"emergencyPower":  17,
			"delayed":         18,
			"blocked":         19,
			"localAlarm":      20,
			"generalAlarm":    21,
			"supervisory":     22,
			"testSupervisory": 23,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LightingInProgress struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLightingInProgress(arg Arg) (*LightingInProgress, error) {
	s := &LightingInProgress{
		enumerations: map[string]uint64{"idle": 0,
			"fadeActive":    1,
			"rampActive":    2,
			"notControlled": 3,
			"other":         4,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LightingOperation struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLightingOperation(arg Arg) (*LightingOperation, error) {
	s := &LightingOperation{
		vendorRange: vendorRange{256, 65535},
		enumerations: map[string]uint64{"none": 0,
			"fadeTo":         1,
			"rampTo":         2,
			"stepUp":         3,
			"stepDown":       4,
			"stepOn":         5,
			"stepOff":        6,
			"warn":           7,
			"warnOff":        8,
			"warnRelinquish": 9,
			"stop":           10,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LightingTransition struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLightingTransition(arg Arg) (*LightingTransition, error) {
	s := &LightingTransition{
		vendorRange: vendorRange{64, 255},
		enumerations: map[string]uint64{"none": 0,
			"fade": 1,
			"ramp": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LockStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLockStatus(arg Arg) (*LockStatus, error) {
	s := &LockStatus{
		enumerations: map[string]uint64{"locked": 0,
			"unlocked": 1,
			"fault":    2,
			"unknown":  3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type LoggingType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLoggingType(arg Arg) (*LoggingType, error) {
	s := &LoggingType{
		vendorRange: vendorRange{64, 255},
		enumerations: map[string]uint64{"polled": 0,
			"cov":       1,
			"triggered": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type Maintenance struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewMaintenance(arg Arg) (*Maintenance, error) {
	s := &Maintenance{
		vendorRange: vendorRange{256, 65535},
		enumerations: map[string]uint64{"none": 0,
			"periodicTest":           1,
			"needServiceOperational": 2,
			"needServiceInoperative": 3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type NodeType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewNodeType(arg Arg) (*NodeType, error) {
	s := &NodeType{
		enumerations: map[string]uint64{"unknown": 0,
			"system":         1,
			"network":        2,
			"device":         3,
			"organizational": 4,
			"area":           5,
			"equipment":      6,
			"point":          7,
			"collection":     8,
			"property":       9,
			"functional":     10,
			"other":          11,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type NotifyType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewNotifyType(arg Arg) (*NotifyType, error) {
	s := &NotifyType{
		enumerations: map[string]uint64{"alarm": 0,
			"event":           1,
			"ackNotification": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type Polarity struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewPolarity(arg Arg) (*Polarity, error) {
	s := &Polarity{
		enumerations: map[string]uint64{"normal": 0,
			"reverse": 1,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ProgramError struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewProgramError(arg Arg) (*ProgramError, error) {
	s := &ProgramError{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"normal": 0,
			"loadFailed": 1,
			"internal":   2,
			"program":    3,
			"other":      4,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ProgramRequest struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewProgramRequest(arg Arg) (*ProgramRequest, error) {
	s := &ProgramRequest{
		enumerations: map[string]uint64{"ready": 0,
			"load":    1,
			"run":     2,
			"halt":    3,
			"restart": 4,
			"unload":  5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ProgramState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewProgramState(arg Arg) (*ProgramState, error) {
	s := &ProgramState{
		enumerations: map[string]uint64{"idle": 0,
			"loading":   1,
			"running":   2,
			"waiting":   3,
			"halted":    4,
			"unloading": 5,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type PropertyIdentifier struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewPropertyIdentifier(arg Arg) (*PropertyIdentifier, error) {
	s := &PropertyIdentifier{
		vendorRange: vendorRange{512, 4194303},
		enumerations: map[string]uint64{"absenteeLimit": 244,
			"acceptedModes":                    175,
			"accessAlarmEvents":                245,
			"accessDoors":                      246,
			"accessEvent":                      247,
			"accessEventAuthenticationFactor":  248,
			"accessEventCredential":            249,
			"accessEventTag":                   322,
			"accessEventTime":                  250,
			"accessTransactionEvents":          251,
			"accompaniment":                    252,
			"accompanimentTime":                253,
			"ackedTransitions":                 0,
			"ackRequired":                      1,
			"action":                           2,
			"actionText":                       3,
			"activationTime":                   254,
			"activeAuthenticationPolicy":       255,
			"activeCovMultipleSubscriptions":   481,
			"activeCovSubscriptions":           152,
			"activeText":                       4,
			"activeVtSessions":                 5,
			"actualShedLevel":                  212,
			"adjustValue":                      176,
			"alarmValue":                       6,
			"alarmValues":                      7,
			"alignIntervals":                   193,
			"all":                              8,
			"allowGroupDelayInhibit":           365,
			"allWritesSuccessful":              9,
			"apduLength":                       399,
			"apduSegmentTimeout":               10,
			"apduTimeout":                      11,
			"applicationSoftwareVersion":       12,
			"archive":                          13,
			"assignedAccessRights":             256,
			"assignedLandingCalls":             447,
			"attemptedSamples":                 124,
			"auditableOperations":              501,
			"auditablePriorityFilter":          500,
			"auditLevel":                       498,
			"auditNotificationRecipient":       499,
			"auditSourceReporter":              497,
			"authenticationFactors":            257,
			"authenticationPolicyList":         258,
			"authenticationPolicyNames":        259,
			"authenticationStatus":             260,
			"authorizationExemptions":          364,
			"authorizationMode":                261,
			"autoSlaveDiscovery":               169,
			"averageValue":                     125,
			"backupAndRestoreState":            338,
			"backupFailureTimeout":             153,
			"backupPreparationTime":            339,
			"bacnetIPGlobalAddress":            407,
			"bacnetIPMode":                     408,
			"bacnetIPMulticastAddress":         409,
			"bacnetIPNATTraversal":             410,
			"bacnetIPUDPPort":                  412,
			"bacnetIPv6Mode":                   435,
			"bacnetIPv6MulticastAddress":       440,
			"bacnetIPv6UDPPort":                438,
			"baseDeviceSecurityPolicy":         327,
			"bbmdAcceptFDRegistrations":        413,
			"bbmdBroadcastDistributionTable":   414,
			"bbmdForeignDeviceTable":           415,
			"belongsTo":                        262,
			"bias":                             14,
			"bitMask":                          342,
			"bitText":                          343,
			"blinkWarnEnable":                  373,
			"bufferSize":                       126,
			"carAssignedDirection":             448,
			"carDoorCommand":                   449,
			"carDoorStatus":                    450,
			"carDoorText":                      451,
			"carDoorZone":                      452,
			"carDriveStatus":                   453,
			"carLoad":                          454,
			"carLoadUnits":                     455,
			"carMode":                          456,
			"carMovingDirection":               457,
			"carPosition":                      458,
			"changeOfStateCount":               15,
			"changeOfStateTime":                16,
			"changesPending":                   416,
			"channelNumber":                    366,
			"clientCovIncrement":               127,
			"command":                          417,
			"commandTimeArray":                 430,
			"configurationFiles":               154,
			"controlGroups":                    367,
			"controlledVariableReference":      19,
			"controlledVariableUnits":          20,
			"controlledVariableValue":          21,
			"count":                            177,
			"countBeforeChange":                178,
			"countChangeTime":                  179,
			"covIncrement":                     22,
			"covPeriod":                        180,
			"covResubscriptionInterval":        128,
			"covuPeriod":                       349,
			"covuRecipients":                   350,
			"credentialDisable":                263,
			"credentials":                      265,
			"credentialsInZone":                266,
			"credentialStatus":                 264,
			"currentCommandPriority":           431,
			"databaseRevision":                 155,
			"dateList":                         23,
			"daylightSavingsStatus":            24,
			"daysRemaining":                    267,
			"deadband":                         25,
			"defaultFadeTime":                  374,
			"defaultPresentValue":              492,
			"defaultRampRate":                  375,
			"defaultStepIncrement":             376,
			"defaultSubordinateRelationship":   490,
			"defaultTimeout":                   393,
			"deleteOnForward":                  502,
			"deployedProfileLocation":          484,
			"derivativeConstant":               26,
			"derivativeConstantUnits":          27,
			"description":                      28,
			"descriptionOfHalt":                29,
			"deviceAddressBinding":             30,
			"deviceType":                       31,
			"deviceUUID":                       507,
			"directReading":                    156,
			"distributionKeyRevision":          328,
			"doNotHide":                        329,
			"doorAlarmState":                   226,
			"doorExtendedPulseTime":            227,
			"doorMembers":                      228,
			"doorOpenTooLongTime":              229,
			"doorPulseTime":                    230,
			"doorStatus":                       231,
			"doorUnlockDelayTime":              232,
			"dutyWindow":                       213,
			"effectivePeriod":                  32,
			"egressActive":                     386,
			"egressTime":                       377,
			"elapsedActiveTime":                33,
			"elevatorGroup":                    459,
			"enable":                           133,
			"energyMeter":                      460,
			"energyMeterRef":                   461,
			"entryPoints":                      268,
			"errorLimit":                       34,
			"escalatorMode":                    462,
			"eventAlgorithmInhibit":            354,
			"eventAlgorithmInhibitRef":         355,
			"eventDetectionEnable":             353,
			"eventEnable":                      35,
			"eventMessageTexts":                351,
			"eventMessageTextsConfig":          352,
			"eventParameters":                  83,
			"eventState":                       36,
			"eventTimeStamps":                  130,
			"eventType":                        37,
			"exceptionSchedule":                38,
			"executionDelay":                   368,
			"exitPoints":                       269,
			"expectedShedLevel":                214,
			"expirationTime":                   270,
			"extendedTimeEnable":               271,
			"failedAttemptEvents":              272,
			"failedAttempts":                   273,
			"failedAttemptsTime":               274,
			"faultHighLimit":                   388,
			"faultLowLimit":                    389,
			"faultParameters":                  358,
			"faultSignals":                     463,
			"faultType":                        359,
			"faultValues":                      39,
			"fdBBMDAddress":                    418,
			"fdSubscriptionLifetime":           419,
			"feedbackValue":                    40,
			"fileAccessMethod":                 41,
			"fileSize":                         42,
			"fileType":                         43,
			"firmwareRevision":                 44,
			"floorNumber":                      506,
			"floorText":                        464,
			"fullDutyBaseline":                 215,
			"globalIdentifier":                 323,
			"groupID":                          465,
			"groupMemberNames":                 346,
			"groupMembers":                     345,
			"groupMode":                        467,
			"higherDeck":                       468,
			"highLimit":                        45,
			"inactiveText":                     46,
			"initialTimeout":                   394,
			"inProcess":                        47,
			"inProgress":                       378,
			"inputReference":                   181,
			"installationID":                   469,
			"instanceOf":                       48,
			"instantaneousPower":               379,
			"integralConstant":                 49,
			"integralConstantUnits":            50,
			"interfaceValue":                   387,
			"intervalOffset":                   195,
			"ipAddress":                        400,
			"ipDefaultGateway":                 401,
			"ipDHCPEnable":                     402,
			"ipDHCPLeaseTime":                  403,
			"ipDHCPLeaseTimeRemaining":         404,
			"ipDHCPServer":                     405,
			"ipDNSServer":                      406,
			"ipSubnetMask":                     411,
			"ipv6Address":                      436,
			"ipv6AutoAddressingEnable":         442,
			"ipv6DefaultGateway":               439,
			"ipv6DHCPLeaseTime":                443,
			"ipv6DHCPLeaseTimeRemaining":       444,
			"ipv6DHCPServer":                   445,
			"ipv6DNSServer":                    441,
			"ipv6PrefixLength":                 437,
			"ipv6ZoneIndex":                    446,
			"issueConfirmedNotifications":      51,
			"isUTC":                            344,
			"keySets":                          330,
			"landingCallControl":               471,
			"landingCalls":                     470,
			"landingDoorStatus":                472,
			"lastAccessEvent":                  275,
			"lastAccessPoint":                  276,
			"lastCommandTime":                  432,
			"lastCredentialAdded":              277,
			"lastCredentialAddedTime":          278,
			"lastCredentialRemoved":            279,
			"lastCredentialRemovedTime":        280,
			"lastKeyServer":                    331,
			"lastNotifyRecord":                 173,
			"lastPriority":                     369,
			"lastRestartReason":                196,
			"lastRestoreTime":                  157,
			"lastStateChange":                  395,
			"lastUseTime":                      281,
			"lifeSafetyAlarmValues":            166,
			"lightingCommand":                  380,
			"lightingCommandDefaultPriority":   381,
			"limitEnable":                      52,
			"limitMonitoringInterval":          182,
			"linkSpeed":                        420,
			"linkSpeedAutonegotiate":           422,
			"linkSpeeds":                       421,
			"listOfGroupMembers":               53,
			"listOfObjectPropertyReferences":   54,
			"listOfSessionKeys":                55,
			"localDate":                        56,
			"localForwardingOnly":              360,
			"localTime":                        57,
			"location":                         58,
			"lockout":                          282,
			"lockoutRelinquishTime":            283,
			"lockStatus":                       233,
			"logBuffer":                        131,
			"logDeviceObjectProperty":          132,
			"loggingObject":                    183,
			"loggingRecord":                    184,
			"loggingType":                      197,
			"logInterval":                      134,
			"lowDiffLimit":                     390,
			"lowerDeck":                        473,
			"lowLimit":                         59,
			"macAddress":                       423,
			"machineRoomID":                    474,
			"maintenanceRequired":              158,
			"makingCarCall":                    475,
			"manipulatedVariableReference":     60,
			"manualSlaveAddressBinding":        170,
			"maskedAlarmValues":                234,
			"masterExemption":                  284,
			"maxActualValue":                   382,
			"maxApduLengthAccepted":            62,
			"maxFailedAttempts":                285,
			"maximumOutput":                    61,
			"maximumSendDelay":                 503,
			"maximumValue":                     135,
			"maximumValueTimestamp":            149,
			"maxInfoFrames":                    63,
			"maxMaster":                        64,
			"maxPresValue":                     65,
			"maxSegmentsAccepted":              167,
			"memberOf":                         159,
			"members":                          286,
			"memberStatusFlags":                347,
			"minActualValue":                   383,
			"minimumOffTime":                   66,
			"minimumOnTime":                    67,
			"minimumOutput":                    68,
			"minimumValue":                     136,
			"minimumValueTimestamp":            150,
			"minPresValue":                     69,
			"mode":                             160,
			"modelName":                        70,
			"modificationDate":                 71,
			"monitoredObjects":                 504,
			"musterPoint":                      287,
			"negativeAccessRules":              288,
			"networkAccessSecurityPolicies":    332,
			"networkInterfaceName":             424,
			"networkNumber":                    425,
			"networkNumberQuality":             426,
			"networkType":                      427,
			"nextStoppingFloor":                476,
			"nodeSubtype":                      207,
			"nodeType":                         208,
			"notificationClass":                17,
			"notificationThreshold":            137,
			"notifyType":                       72,
			"numberOfApduRetries":              73,
			"numberOfAuthenticationPolicies":   289,
			"numberOfStates":                   74,
			"objectIdentifier":                 75,
			"objectList":                       76,
			"objectName":                       77,
			"objectPropertyReference":          78,
			"objectType":                       79,
			"occupancyCount":                   290,
			"occupancyCountAdjust":             291,
			"occupancyCountEnable":             292,
			"occupancyExemption":               293,
			"occupancyLowerLimit":              294,
			"occupancyLowerLimitEnforced":      295,
			"occupancyState":                   296,
			"occupancyUpperLimit":              297,
			"occupancyUpperLimitEnforced":      298,
			"operationDirection":               477,
			"operationExpected":                161,
			"optional":                         80,
			"outOfService":                     81,
			"outputUnits":                      82,
			"packetReorderTime":                333,
			"passbackExemption":                299,
			"passbackMode":                     300,
			"passbackTimeout":                  301,
			"passengerAlarm":                   478,
			"polarity":                         84,
			"portFilter":                       363,
			"positiveAccessRules":              302,
			"power":                            384,
			"powerMode":                        479,
			"prescale":                         185,
			"presentStage":                     493,
			"presentValue":                     85,
			"priority":                         86,
			"priorityArray":                    87,
			"priorityForWriting":               88,
			"processIdentifier":                89,
			"processIdentifierFilter":          361,
			"profileLocation":                  485,
			"profileName":                      168,
			"programChange":                    90,
			"programLocation":                  91,
			"programState":                     92,
			"propertyList":                     371,
			"proportionalConstant":             93,
			"proportionalConstantUnits":        94,
			"protocolLevel":                    482,
			"protocolObjectTypesSupported":     96,
			"protocolRevision":                 139,
			"protocolServicesSupported":        97,
			"protocolVersion":                  98,
			"pulseRate":                        186,
			"readOnly":                         99,
			"reasonForDisable":                 303,
			"reasonForHalt":                    100,
			"recipientList":                    102,
			"recordCount":                      141,
			"recordsSinceNotification":         140,
			"referencePort":                    483,
			"registeredCarCall":                480,
			"reliability":                      103,
			"reliabilityEvaluationInhibit":     357,
			"relinquishDefault":                104,
			"represents":                       491,
			"requestedShedLevel":               218,
			"requestedUpdateInterval":          348,
			"required":                         105,
			"resolution":                       106,
			"restartNotificationRecipients":    202,
			"restoreCompletionTime":            340,
			"restorePreparationTime":           341,
			"routingTable":                     428,
			"scale":                            187,
			"scaleFactor":                      188,
			"scheduleDefault":                  174,
			"securedStatus":                    235,
			"securityPDUTimeout":               334,
			"securityTimeWindow":               335,
			"segmentationSupported":            107,
			"sendNow":                          505,
			"serialNumber":                     372,
			"setpoint":                         108,
			"setpointReference":                109,
			"setting":                          162,
			"shedDuration":                     219,
			"shedLevelDescriptions":            220,
			"shedLevels":                       221,
			"silenced":                         163,
			"slaveAddressBinding":              171,
			"slaveProxyEnable":                 172,
			"stageNames":                       495,
			"stages":                           494,
			"startTime":                        142,
			"stateChangeValues":                396,
			"stateDescription":                 222,
			"stateText":                        110,
			"statusFlags":                      111,
			"stopTime":                         143,
			"stopWhenFull":                     144,
			"strikeCount":                      391,
			"structuredObjectList":             209,
			"subordinateAnnotations":           210,
			"subordinateList":                  211,
			"subordinateNodeTypes":             487,
			"subordinateRelationships":         489,
			"subordinateTags":                  488,
			"subscribedRecipients":             362,
			"supportedFormatClasses":           305,
			"supportedFormats":                 304,
			"supportedSecurityAlgorithms":      336,
			"systemStatus":                     112,
			"tags":                             486,
			"targetReferences":                 496,
			"threatAuthority":                  306,
			"threatLevel":                      307,
			"timeDelay":                        113,
			"timeDelayNormal":                  356,
			"timeOfActiveTimeReset":            114,
			"timeOfDeviceRestart":              203,
			"timeOfStateCountReset":            115,
			"timeOfStrikeCountReset":           392,
			"timerRunning":                     397,
			"timerState":                       398,
			"timeSynchronizationInterval":      204,
			"timeSynchronizationRecipients":    116,
			"totalRecordCount":                 145,
			"traceFlag":                        308,
			"trackingValue":                    164,
			"transactionNotificationClass":     309,
			"transition":                       385,
			"trigger":                          205,
			"units":                            117,
			"updateInterval":                   118,
			"updateKeySetTimeout":              337,
			"updateTime":                       189,
			"userExternalIdentifier":           310,
			"userInformationReference":         311,
			"userName":                         317,
			"userType":                         318,
			"usesRemaining":                    319,
			"utcOffset":                        119,
			"utcTimeSynchronizationRecipients": 206,
			"validSamples":                     146,
			"valueBeforeChange":                190,
			"valueChangeTime":                  192,
			"valueSet":                         191,
			"valueSource":                      433,
			"valueSourceArray":                 434,
			"varianceValue":                    151,
			"vendorIdentifier":                 120,
			"vendorName":                       121,
			"verificationTime":                 326,
			"virtualMACAddressTable":           429,
			"vtClassesSupported":               122,
			"weeklySchedule":                   123,
			"windowInterval":                   147,
			"windowSamples":                    148,
			"writeStatus":                      370,
			"zoneFrom":                         320,
			"zoneMembers":                      165,
			"zoneTo":                           321,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type Relationship struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewRelationship(arg Arg) (*Relationship, error) {
	s := &Relationship{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"default":           1,
			"contains":          2,
			"containedBy":       3,
			"uses":              4,
			"usedBy":            5,
			"commands":          6,
			"commandedBy":       7,
			"adjusts":           8,
			"adjustedBy":        9,
			"ingress":           10,
			"egress":            11,
			"suppliesAir":       12,
			"receivesAir":       13,
			"suppliesHotAir":    14,
			"receivesHotAir":    15,
			"suppliesCoolAir":   16,
			"receivesCoolAir":   17,
			"suppliesPower":     18,
			"receivesPower":     19,
			"suppliesGas":       20,
			"receivesGas":       21,
			"suppliesWater":     22,
			"receivesWater":     23,
			"suppliesHotWater":  24,
			"receivesHotWater":  25,
			"suppliesCoolWater": 26,
			"receivesCoolWater": 27,
			"suppliesSteam":     28,
			"receivesSteam":     29,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type Reliability struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewReliability(arg Arg) (*Reliability, error) {
	s := &Reliability{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"noFaultDetected": 0,
			"noSensor":                      1,
			"overRange":                     2,
			"underRange":                    3,
			"openLoop":                      4,
			"shortedLoop":                   5,
			"noOutput":                      6,
			"unreliableOther":               7,
			"processError":                  8,
			"multiStateFault":               9,
			"configurationError":            10,
			"communicationFailure":          12,
			"memberFault":                   13,
			"monitoredObjectFault":          14,
			"tripped":                       15,
			"lampFailure":                   16,
			"activationFailure":             17,
			"renewDHCPFailure":              18,
			"renewFDRegistrationFailure":    19,
			"restartAutoNegotiationFailure": 20,
			"restartFailure":                21,
			"proprietaryCommandFailure":     22,
			"faultsListed":                  23,
			"referencedObjectFault":         24,
			"multiStateOutOfRange":          25,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type RestartReason struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewRestartReason(arg Arg) (*RestartReason, error) {
	s := &RestartReason{
		vendorRange: vendorRange{64, 255},
		enumerations: map[string]uint64{"unknown": 0,
			"coldstart":          1,
			"warmstart":          2,
			"detectedPowerLost":  3,
			"detectedPoweredOff": 4,
			"hardwareWatchdog":   5,
			"softwareWatchdog":   6,
			"suspended":          7,
			"activateChanges":    8,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type SecurityLevel struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewSecurityLevel(arg Arg) (*SecurityLevel, error) {
	s := &SecurityLevel{
		enumerations: map[string]uint64{"incapable": 0,
			"plain":             1,
			"signed":            2,
			"encrypted":         3,
			"signedEndToEnd":    4,
			"encryptedEndToEnd": 4,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type SecurityPolicy struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewSecurityPolicy(arg Arg) (*SecurityPolicy, error) {
	s := &SecurityPolicy{
		enumerations: map[string]uint64{"plainNonTrusted": 0,
			"plainTrusted":     1,
			"signedTrusted":    2,
			"encryptedTrusted": 3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ShedState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewShedState(arg Arg) (*ShedState, error) {
	s := &ShedState{
		enumerations: map[string]uint64{"shedInactive": 0,
			"shedRequestPending": 1,
			"shedCompliant":      2,
			"shedNonCompliant":   3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type Segmentation struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewSegmentation(arg Arg) (*Segmentation, error) {
	s := &Segmentation{
		enumerations: map[string]uint64{"segmentedBoth": 0,
			"segmentedTransmit": 1,
			"segmentedReceive":  2,
			"noSegmentation":    3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type SilencedState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewSilencedState(arg Arg) (*SilencedState, error) {
	s := &SilencedState{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"unsilenced": 0,
			"audibleSilenced": 1,
			"visibleSilenced": 2,
			"allSilenced":     3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type TimerState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewTimerState(arg Arg) (*TimerState, error) {
	s := &TimerState{
		enumerations: map[string]uint64{"idle": 0,
			"running": 1,
			"expired": 2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type TimerTransition struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewTimerTransition(arg Arg) (*TimerTransition, error) {
	s := &TimerTransition{
		enumerations: map[string]uint64{"none": 0,
			"idleToRunning":    1,
			"runningToIdle":    2,
			"runningToRunning": 3,
			"runningToExpired": 4,
			"forcedToExpired":  5,
			"expiredToIdle":    6,
			"expiredToRunning": 7,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type VTClass struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewVTClass(arg Arg) (*VTClass, error) {
	s := &VTClass{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"defaultTerminal": 0,
			"ansiX364": 1,
			"decVt52":  2,
			"decVt100": 3,
			"decVt220": 4,
			"hp70094":  5,
			"ibm3130":  6,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type WriteStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewWriteStatus(arg Arg) (*WriteStatus, error) {
	s := &WriteStatus{
		enumerations: map[string]uint64{"idle": 0,
			"inProgress": 1,
			"successful": 2,
			"failed":     3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type NetworkType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewNetworkType(arg Arg) (*NetworkType, error) {
	s := &NetworkType{
		enumerations: map[string]uint64{"ethernet": 0,
			"arcnet":  1,
			"mstp":    2,
			"ptp":     3,
			"lontalk": 4,
			"ipv4":    5,
			"zigbee":  6,
			"virtual": 7,
			// "non-bacnet":  8  Removed in Version 1, Revision 18,
			"ipv6":   9,
			"serial": 10,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type ProtocolLevel struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewProtocolLevel(arg Arg) (*ProtocolLevel, error) {
	s := &ProtocolLevel{
		enumerations: map[string]uint64{"physical": 0,
			"protocol":             1,
			"bacnetApplication":    2,
			"nonBacnetApplication": 3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type NetworkNumberQuality struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewNetworkNumberQuality(arg Arg) (*NetworkNumberQuality, error) {
	s := &NetworkNumberQuality{
		enumerations: map[string]uint64{"unknown": 0,
			"learned":           1,
			"learnedConfigured": 2,
			"configured":        3,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type NetworkPortCommand struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewNetworkPortCommand(arg Arg) (*NetworkPortCommand, error) {
	s := &NetworkPortCommand{
		enumerations: map[string]uint64{"idle": 0,
			"discardChanges":         1,
			"renewFdDRegistration":   2,
			"restartSlaveDiscovery":  3,
			"renewDHCP":              4,
			"restartAutonegotiation": 5,
			"disconnect":             6,
			"restartPort":            7,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type IPMode struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewIPMode(arg Arg) (*IPMode, error) {
	s := &IPMode{
		enumerations: map[string]uint64{"normal": 0,
			"foreign": 1,
			"bbmd":    2,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}

type RouterEntryStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewRouterEntryStatus(arg Arg) (*RouterEntryStatus, error) {
	s := &RouterEntryStatus{
		enumerations: map[string]uint64{"available": 0,
			"busy":         1,
			"disconnected": 2,
		},
	}
	panic("implement me enumasdasd")
	return s, nil
}

//
//   Forward Sequences
//

type HostAddress struct {
	*Choice
	choiceElements []Element
}

func NewHostAddress(arg Arg) (*HostAddress, error) {
	s := &HostAddress{
		choiceElements: []Element{
			NewElement("none", V2E(NewNull), WithElementContext(0)),
			NewElement("ipAddress", V2E(NewOctetString), WithElementContext(1)), //  4 octets for B/IP or 16 octets for B/IPv6
			NewElement("name", V2E(NewCharacterString), WithElementContext(2)),  //  Internet host name (see RFC 1123)
		},
	}
	panic("implementchoice")
	return s, nil
}

type HostNPort struct {
	*Sequence
	sequenceElements []Element
}

func NewHostNPort(arg Arg) (*HostNPort, error) {
	s := &HostNPort{
		sequenceElements: []Element{
			NewElement("host", V2E(NewHostAddress), WithElementContext(0)),
			NewElement("port", V2E(NewUnsigned16), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type BDTEntry struct {
	*Sequence
	sequenceElements []Element
}

func NewBDTEntry(arg Arg) (*BDTEntry, error) {
	s := &BDTEntry{
		sequenceElements: []Element{
			NewElement("bbmdAddress", V2E(NewHostNPort), WithElementContext(0)),
			NewElement("broadcastMask", V2E(NewOctetString), WithElementContext(1)), //  shall be present if BACnet/IP, and absent for BACnet/IPv6
		},
	}
	panic("implementchoice")
	return s, nil
}

type FDTEntry struct {
	*Sequence
	sequenceElements []Element
}

func NewFDTEntry(arg Arg) (*FDTEntry, error) {
	s := &FDTEntry{
		sequenceElements: []Element{
			NewElement("bacnetIPAddress", V2E(NewOctetString), WithElementContext(0)),    //  the 6-octet B/IP or 18-octet B/IPv6 address of the registrant
			NewElement("timeToLive", V2E(NewUnsigned16), WithElementContext(1)),          //  time to live in seconds at the time of registration
			NewElement("remainingTimeToLive", V2E(NewUnsigned16), WithElementContext(2)), //  remaining time to live in seconds, incl. grace period
		},
	}
	panic("implementchoice")
	return s, nil
}

type VMACEntry struct {
	*Sequence
	sequenceElements []Element
}

func NewVMACEntry(arg Arg) (*VMACEntry, error) {
	s := &VMACEntry{
		sequenceElements: []Element{
			NewElement("virtualMACAddress", V2E(NewOctetString), WithElementContext(0)), //  maximum size 6 octets
			NewElement("nativeMACAddress", V2E(NewOctetString), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type PropertyReference struct {
	*Sequence
	sequenceElements []Element
}

func NewPropertyReference(arg Arg) (*PropertyReference, error) {
	s := &PropertyReference{
		sequenceElements: []Element{
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(0)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(1), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type RouterEntry struct {
	*Sequence
	sequenceElements []Element
}

func NewRouterEntry(arg Arg) (*RouterEntry, error) {
	s := &RouterEntry{
		sequenceElements: []Element{
			NewElement("networkNumber", V2E(NewUnsigned16), WithElementContext(0)),
			NewElement("macAddress", V2E(NewOctetString), WithElementContext(1)),
			NewElement("status", V2E(NewRouterEntryStatus), WithElementContext(2)), //  Defined Above
			NewElement("performanceIndex", V2E(NewUnsigned8), WithElementContext(3), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NameValue struct {
	*Sequence

	sequenceElements []Element
	name             string
	value            any
}

func NewNameValue(args Args) (*NameValue, error) {
	s := &NameValue{
		sequenceElements: []Element{
			NewElement("name", V2E(NewCharacterString)),
			NewElement("value", Vs2E(NewAnyAtomic), WithElementOptional(true)),
		},
	}
	// default to no value
	s.name = GetOptional[string](args, 0, "")
	s.value = GetOptional[any](args, 1, nil)

	if s.value == nil {
		return s, nil
	}
	switch value := s.value.(type) {
	case IsAtomic:
		s.value = value
	case DateTime:
		s.value = value
	case Tag:
		var err error
		s.value, err = value.AppToObject()
		if err != nil {
			return nil, errors.Wrap(err, "error converting tag")
		}
	default:
		return nil, errors.Errorf("invalid constructor datatype: %T", value)
	}
	return s, nil
}

func (s *NameValue) Encode(arg Arg) error {
	tagList := arg.(*TagList)
	// build a tag and encode the name into it
	tag, err := NewTag(NoArgs)
	if err != nil {
		return errors.Wrap(err, "error converting tag")
	}
	characterString, err := NewCharacterString(s.name)
	if err != nil {
		return errors.Wrap(err, "error creating character string")
	}
	if err := characterString.Encode(tag); err != nil {
		return errors.Wrap(err, "error encoding")
	}
	context, err := tag.AppToContext(0)
	if err != nil {
		return errors.Wrap(err, "error converting tag")
	}
	tagList.Append(context)

	// the value is optional
	if s.value != nil {
		if v, ok := s.value.(*DateTime); ok {
			// has its own encoder
			if err := v.Encode(tagList); err != nil {
				return errors.Wrap(err, "error converting tag")
			}
		} else if e, ok := s.value.(Encoder); ok {
			// atomic values encode into a tag
			tag, err = NewTag(NoArgs)
			if err != nil {
				return errors.Wrap(err, "error creating tag")
			}
			if err := e.Encode(tag); err != nil {
				return errors.Wrap(err, "error converting tag")
			}
			tagList.Append(tag)
		}
	}
	return nil
}

func (s *NameValue) Decode(arg Arg) error {
	tagList := arg.(*TagList)

	// no contents yet
	s.name = ""
	s.value = ""

	// look for the context encoded character string
	tag := tagList.Peek()
	if tag == nil || (tag.GetTagClass() != readWriteModel.TagClass_CONTEXT_SPECIFIC_TAGS) || (tag.GetTagNumber() != 0) {
		return MissingRequiredParameter{Message: fmt.Sprintf("%s is a missing required element of %p", s.name, s)}
	}

	// pop it off and save the value
	tagList.Pop()
	tag, err := tag.ContextToApp(uint(readWriteModel.BACnetDataType_CHARACTER_STRING))
	if err != nil {
		return errors.Wrap(err, "error converting tag")
	}
	characterString, err := NewCharacterString(tag)
	if err != nil {
		return errors.Wrap(err, "error converting tag to string")
	}
	s.name = characterString.GetValue()

	// look for the optional application encoded value
	tag = tagList.Peek()
	if tag != nil && (tag.GetTagClass() == readWriteModel.TagClass_APPLICATION_TAGS) {
		// if it is a date check the next one for a time
		if tag.GetTagNumber() == uint(readWriteModel.BACnetDataType_DATE) && (len(tagList.GetTagList()) >= 2) {
			nextTag := tagList.GetTagList()[1]

			if nextTag.GetTagClass() == readWriteModel.TagClass_APPLICATION_TAGS && (nextTag.GetTagNumber() == uint(readWriteModel.BACnetDataType_TIME)) {
				s.value, err = NewDateTime(NoArgs)
				if err != nil {
					return errors.Wrap(err, "error creating date time")
				}
				if err := s.value.(Decoder).Decode(tagList); err != nil {
					return errors.Wrap(err, "error decoding taglist")
				}
			}

			// just a primitive value
			if s.value == nil {
				tagList.Pop()
				s.value, err = tag.AppToObject()
				if err != nil {
					return errors.Wrap(err, "error converting tag")
				}
			}
		}
	}
	return nil
}

type NameValueCollection struct {
	*Sequence
	sequenceElements []Element
}

func NewNameValueCollection(arg Arg) (*NameValueCollection, error) {
	s := &NameValueCollection{
		sequenceElements: []Element{
			NewElement("members", SequenceOfs(NewNameValue), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DeviceAddress struct {
	*Sequence
	sequenceElements []Element
}

func NewDeviceAddress(arg Arg) (*DeviceAddress, error) {
	s := &DeviceAddress{
		sequenceElements: []Element{
			NewElement("networkNumber", V2E(NewUnsigned)),
			NewElement("macAddress", V2E(NewOctetString)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DeviceObjectPropertyReference struct {
	*Sequence
	sequenceElements []Element
}

func NewDeviceObjectPropertyReference(arg Arg) (*DeviceObjectPropertyReference, error) {
	s := &DeviceObjectPropertyReference{
		sequenceElements: []Element{
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(1)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(2), WithElementOptional(true)),
			NewElement("deviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(3), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DeviceObjectReference struct {
	*Sequence
	sequenceElements []Element
}

func NewDeviceObjectReference(arg Arg) (*DeviceObjectReference, error) {
	s := &DeviceObjectReference{
		sequenceElements: []Element{
			NewElement("deviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0), WithElementOptional(true)),
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DateTime struct {
	*Sequence
	sequenceElements []Element
}

func NewDateTime(arg Arg) (*DateTime, error) {
	s := &DateTime{
		sequenceElements: []Element{
			NewElement("date", Vs2E(NewDate)),
			NewElement("time", Vs2E(NewTime)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DateRange struct {
	*Sequence
	sequenceElements []Element
}

func NewDateRange(arg Arg) (*DateRange, error) {
	s := &DateRange{
		sequenceElements: []Element{
			NewElement("startDate", Vs2E(NewDate)),
			NewElement("endDate", Vs2E(NewDate)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ErrorType struct {
	*Sequence
	sequenceElements []Element
}

func NewErrorType(arg Arg) (*ErrorType, error) {
	s := &ErrorType{
		sequenceElements: []Element{
			NewElement("errorClass", V2E(NewErrorClass)),
			NewElement("errorCode", V2E(NewErrorCode)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LandingCallStatusCommand struct {
	*Choice
	choiceElements []Element
}

func NewLandingCallStatusCommand(arg Arg) (*LandingCallStatusCommand, error) {
	s := &LandingCallStatusCommand{
		choiceElements: []Element{
			NewElement("direction", V2E(NewLiftCarDirection), WithElementContext(1)),
			NewElement("destination", V2E(NewUnsigned8), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LandingCallStatus struct {
	*Sequence
	sequenceElements []Element
}

func NewLandingCallStatus(arg Arg) (*LandingCallStatus, error) {
	s := &LandingCallStatus{
		sequenceElements: []Element{
			NewElement("floorNumber", V2E(NewUnsigned8), WithElementContext(0)),
			NewElement("command", V2E(NewLandingCallStatusCommand)),
			NewElement("floorText", V2E(NewCharacterString), WithElementContext(3), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LandingDoorStatusLandingDoor struct {
	*Sequence
	sequenceElements []Element
}

func NewLandingDoorStatusLandingDoor(arg Arg) (*LandingDoorStatusLandingDoor, error) {
	s := &LandingDoorStatusLandingDoor{
		sequenceElements: []Element{
			NewElement("floorNumber", V2E(NewUnsigned8), WithElementContext(0)),
			NewElement("doorStatus", V2E(NewDoorStatus), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LandingDoorStatus struct {
	*Sequence
	sequenceElements []Element
}

func NewLandingDoorStatus(arg Arg) (*LandingDoorStatus, error) {
	s := &LandingDoorStatus{
		sequenceElements: []Element{
			NewElement("landingDoors", SequenceOf(NewLandingDoorStatusLandingDoor), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LiftCarCallList struct {
	*Sequence
	sequenceElements []Element
}

func NewLiftCarCallList(arg Arg) (*LiftCarCallList, error) {
	s := &LiftCarCallList{
		sequenceElements: []Element{
			NewElement("floorNumbers", SequenceOf(NewUnsigned8), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LightingCommand struct {
	*Sequence
	sequenceElements []Element
}

func NewLightingCommand(arg Arg) (*LightingCommand, error) {
	s := &LightingCommand{
		sequenceElements: []Element{
			NewElement("operation", V2E(NewLightingOperation), WithElementContext(0)),
			NewElement("targetLevel", V2E(NewReal), WithElementContext(1), WithElementOptional(true)),
			NewElement("rampRate", V2E(NewReal), WithElementContext(2), WithElementOptional(true)),
			NewElement("stepIncrement", V2E(NewReal), WithElementContext(3), WithElementOptional(true)),
			NewElement("fadeTime", V2E(NewUnsigned), WithElementContext(4), WithElementOptional(true)),
			NewElement("priority", V2E(NewUnsigned), WithElementContext(5), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ObjectPropertyReference struct {
	*Sequence
	sequenceElements []Element
}

func NewObjectPropertyReference(arg Arg) (*ObjectPropertyReference, error) {
	s := &ObjectPropertyReference{
		sequenceElements: []Element{
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(1)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(2), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type OptionalBinaryPV struct {
	*Choice
	choiceElements []Element
}

func NewOptionalBinaryPV(arg Arg) (*OptionalBinaryPV, error) {
	s := &OptionalBinaryPV{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("binaryPV", V2E(NewBinaryPV)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type OptionalCharacterString struct {
	*Choice
	choiceElements []Element
}

func NewOptionalCharacterString(arg Arg) (*OptionalCharacterString, error) {
	s := &OptionalCharacterString{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("characterstring", V2E(NewCharacterString)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type OptionalPriorityFilter struct {
	*Choice
	choiceElements []Element
}

func NewOptionalPriorityFilter(arg Arg) (*OptionalPriorityFilter, error) {
	s := &OptionalPriorityFilter{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("filter", V2E(NewPriorityFilter)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type OptionalReal struct {
	*Choice
	choiceElements []Element
}

func NewOptionalReal(arg Arg) (*OptionalReal, error) {
	s := &OptionalReal{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("real", V2E(NewReal)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type OptionalUnsigned struct {
	*Choice
	choiceElements []Element
}

func NewOptionalUnsigned(arg Arg) (*OptionalUnsigned, error) {
	s := &OptionalUnsigned{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("unsigned", V2E(NewUnsigned)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ProcessIdSelection struct {
	*Choice
	choiceElements []Element
}

func NewProcessIdSelection(arg Arg) (*ProcessIdSelection, error) {
	s := &ProcessIdSelection{
		choiceElements: []Element{
			NewElement("processIdentifier", V2E(NewUnsigned)),
			NewElement("nullValue", V2E(NewNull)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type PropertyStates struct {
	*Choice
	vendorRange    vendorRange
	choiceElements []Element
}

func NewPropertyStates(arg Arg) (*PropertyStates, error) {
	s := &PropertyStates{
		vendorRange: vendorRange{64, 254},
		choiceElements: []Element{
			NewElement("booleanValue", V2E(NewBoolean), WithElementContext(0)),
			NewElement("binaryValue", V2E(NewBinaryPV), WithElementContext(1)),
			NewElement("eventType", V2E(NewEventType), WithElementContext(2)),
			NewElement("polarity", V2E(NewPolarity), WithElementContext(3)),
			NewElement("programChange", V2E(NewProgramRequest), WithElementContext(4)),
			NewElement("programState", V2E(NewProgramState), WithElementContext(5)),
			NewElement("reasonForHalt", V2E(NewProgramError), WithElementContext(6)),
			NewElement("reliability", V2E(NewReliability), WithElementContext(7)),
			NewElement("state", V2E(NewEventState), WithElementContext(8)),
			NewElement("systemStatus", V2E(NewDeviceStatus), WithElementContext(9)),
			NewElement("units", V2E(NewEngineeringUnits), WithElementContext(10)),
			NewElement("unsignedValue", V2E(NewUnsigned), WithElementContext(11)),
			NewElement("lifeSafetyMode", V2E(NewLifeSafetyMode), WithElementContext(12)),
			NewElement("lifeSafetyState", V2E(NewLifeSafetyState), WithElementContext(13)),
			NewElement("restartReason", V2E(NewRestartReason), WithElementContext(14)),
			NewElement("doorAlarmState", V2E(NewDoorAlarmState), WithElementContext(15)),
			NewElement("action", V2E(NewAction), WithElementContext(16)),
			NewElement("doorSecuredStatus", V2E(NewDoorSecuredStatus), WithElementContext(17)),
			NewElement("doorStatus", V2E(NewDoorStatus), WithElementContext(18)),
			NewElement("doorValue", V2E(NewDoorValue), WithElementContext(19)),
			NewElement("fileAccessMethod", V2E(NewFileAccessMethod), WithElementContext(20)),
			NewElement("lockStatus", V2E(NewLockStatus), WithElementContext(21)),
			NewElement("lifeSafetyOperation", V2E(NewLifeSafetyOperation), WithElementContext(22)),
			NewElement("maintenance", V2E(NewMaintenance), WithElementContext(23)),
			NewElement("nodeType", V2E(NewNodeType), WithElementContext(24)),
			NewElement("notifyType", V2E(NewNotifyType), WithElementContext(25)),
			NewElement("securityLevel", V2E(NewSecurityLevel), WithElementContext(26)),
			NewElement("shedState", V2E(NewShedState), WithElementContext(27)),
			NewElement("silencedState", V2E(NewSilencedState), WithElementContext(28)),
			NewElement("accessEvent", V2E(NewAccessEvent), WithElementContext(30)),
			NewElement("zoneOccupancyState", V2E(NewAccessZoneOccupancyState), WithElementContext(31)),
			NewElement("accessCredentialDisableReason", V2E(NewAccessCredentialDisableReason), WithElementContext(32)),
			NewElement("accessCredentialDisable", V2E(NewAccessCredentialDisable), WithElementContext(33)),
			NewElement("authenticationStatus", V2E(NewAuthenticationStatus), WithElementContext(34)),
			NewElement("backupState", V2E(NewBackupState), WithElementContext(36)),
			NewElement("writeStatus", V2E(NewWriteStatus), WithElementContext(370)),
			NewElement("lightingInProgress", V2E(NewLightingInProgress), WithElementContext(38)),
			NewElement("lightingOperation", V2E(NewLightingOperation), WithElementContext(39)),
			NewElement("lightingTransition", V2E(NewLightingTransition), WithElementContext(40)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type PropertyValue struct {
	*Sequence
	sequenceElements []Element
}

func NewPropertyValue(arg Arg) (*PropertyValue, error) {
	s := &PropertyValue{
		sequenceElements: []Element{
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(0)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(1), WithElementOptional(true)),
			NewElement("value", Vs2E(NewAny), WithElementContext(2)),
			NewElement("priority", V2E(NewUnsigned), WithElementContext(3), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type Recipient struct {
	*Choice
	choiceElements []Element
}

func NewRecipient(arg Arg) (*Recipient, error) {
	s := &Recipient{
		choiceElements: []Element{
			NewElement("device", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("address", V2E(NewDeviceAddress), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type RecipientProcess struct {
	*Sequence
	sequenceElements []Element
}

func NewRecipientProcess(arg Arg) (*RecipientProcess, error) {
	s := &RecipientProcess{
		sequenceElements: []Element{
			NewElement("recipient", V2E(NewRecipient), WithElementContext(0)),
			NewElement("processIdentifier", V2E(NewUnsigned), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type TimerStateChangeValue struct {
	*Choice
	choiceElements []Element
}

func NewTimerStateChangeValue(arg Arg) (*TimerStateChangeValue, error) {
	s := &TimerStateChangeValue{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("boolean", V2E(NewBoolean)),
			NewElement("unsigned", V2E(NewUnsigned)),
			NewElement("integer", V2E(NewInteger)),
			NewElement("real", V2E(NewReal)),
			NewElement("double", V2E(NewDouble)),
			NewElement("octetstring", V2E(NewOctetString)),
			NewElement("characterstring", V2E(NewCharacterString)),
			NewElement("bitstring", Vs2E(NewBitString)),
			NewElement("enumerated", Vs2E(NewEnumerated)),
			NewElement("date", Vs2E(NewDate)),
			NewElement("time", Vs2E(NewTime)),
			NewElement("objectidentifier", Vs2E(NewObjectIdentifier)),
			NewElement("noValue", V2E(NewNull), WithElementContext(0)),
			NewElement("constructedValue", Vs2E(NewAny), WithElementContext(1)),
			NewElement("datetime", V2E(NewDateTime), WithElementContext(2)),
			NewElement("lightingCommand", V2E(NewLightingCommand), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type TimeStamp struct {
	*Choice
	choiceElements []Element
}

func NewTimeStamp(arg Arg) (*TimeStamp, error) {
	s := &TimeStamp{
		choiceElements: []Element{
			NewElement("time", Vs2E(NewTime), WithElementContext(0)),
			NewElement("sequenceNumber", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("dateTime", V2E(NewDateTime), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type TimeValue struct {
	*Sequence
	sequenceElements []Element
}

func NewTimeValue(arg Arg) (*TimeValue, error) {
	s := &TimeValue{
		sequenceElements: []Element{
			NewElement("time", Vs2E(NewTime)),
			NewElement("value", Vs2E(NewAnyAtomic)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type WeekNDay struct {
	*OctetString
}

func NewWeekNDay(arg Arg) (*WeekNDay, error) {
	s := &WeekNDay{}
	panic("implement me")
	return s, nil
}

func (w WeekNDay) String() string {
	value := w.GetValue()
	if len(value) != 3 {
		return "WeekNDay(?): " + w.OctetString.String()
	} else {
		return fmt.Sprintf("WeekNDay(%d, %d, %d)", value[0], value[1], value[2])
	}
}

//
//   Sequences
//

type AccessRule struct {
	*Sequence
	sequenceElements []Element
}

func NewAccessRule(arg Arg) (*AccessRule, error) {
	s := &AccessRule{
		sequenceElements: []Element{
			NewElement("timeRangeSpecifier", V2E(NewAccessRuleTimeRangeSpecifier), WithElementContext(0)),
			NewElement("timeRange", V2E(NewDeviceObjectPropertyReference), WithElementContext(1), WithElementOptional(true)),
			NewElement("locationSpecifier", V2E(NewAccessRuleLocationSpecifier), WithElementContext(2)),
			NewElement("location", V2E(NewDeviceObjectReference), WithElementContext(3), WithElementOptional(true)),
			NewElement("enable", V2E(NewBoolean), WithElementContext(4)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AccessThreatLevel struct {
	*Unsigned
}

func NewAccessThreatLevel(arg Arg) (*AccessThreatLevel, error) {
	s := &AccessThreatLevel{}
	_low_limit := 0
	_ = _low_limit
	_high_limit := 100
	_ = _high_limit
	panic("implement me")
	return s, nil
}

type AccumulatorRecord struct {
	*Sequence
	sequenceElements []Element
}

func NewAccumulatorRecord(arg Arg) (*AccumulatorRecord, error) {
	s := &AccumulatorRecord{
		sequenceElements: []Element{
			NewElement("timestamp", V2E(NewDateTime), WithElementContext(0)),
			NewElement("presentValue", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("accumulatedValue", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("accumulatorStatus", V2E(NewAccumulatorRecordAccumulatorStatus), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ActionCommand struct {
	*Sequence
	sequenceElements []Element
}

func NewActionCommand(arg Arg) (*ActionCommand, error) {
	s := &ActionCommand{
		sequenceElements: []Element{
			NewElement("deviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0), WithElementOptional(true)),
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(1)),
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(2)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(3), WithElementOptional(true)),
			NewElement("propertyValue", Vs2E(NewAny), WithElementContext(4)),
			NewElement("priority", V2E(NewUnsigned), WithElementContext(5), WithElementOptional(true)),
			NewElement("postDelay", V2E(NewUnsigned), WithElementContext(6), WithElementOptional(true)),
			NewElement("quiteOnFailure", V2E(NewBoolean), WithElementContext(7)),
			NewElement("writeSuccessFul", V2E(NewBoolean), WithElementContext(8)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ActionList struct {
	*Sequence
	sequenceElements []Element
}

func NewActionList(arg Arg) (*ActionList, error) {
	s := &ActionList{
		sequenceElements: []Element{
			NewElement("action", SequenceOf(NewActionCommand), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type Address struct {
	*Sequence
	sequenceElements []Element
}

func NewAddress(arg Arg) (*Address, error) {
	s := &Address{
		sequenceElements: []Element{
			NewElement("networkNumber", V2E(NewUnsigned16)),
			NewElement("macAddress", V2E(NewOctetString)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AddressBinding struct {
	*Sequence
	sequenceElements []Element
}

func NewAddressBinding(arg Arg) (*AddressBinding, error) {
	s := &AddressBinding{
		sequenceElements: []Element{
			NewElement("deviceObjectIdentifier", Vs2E(NewObjectIdentifier)),
			NewElement("deviceAddress", V2E(NewDeviceAddress)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AssignedAccessRights struct {
	*Sequence
	sequenceElements []Element
}

func NewAssignedAccessRights(arg Arg) (*AssignedAccessRights, error) {
	s := &AssignedAccessRights{
		sequenceElements: []Element{
			NewElement("assignedAccessRights", V2E(NewDeviceObjectReference), WithElementContext(0)),
			NewElement("enable", V2E(NewBoolean), WithElementContext(1)),
		},
	}
	//serviceChoice: 15,
	panic("implementchoice")
	return s, nil
}

type AssignedLandingCallsLandingCalls struct {
	*Sequence
	sequenceElements []Element
}

func NewAssignedLandingCallsLandingCalls(arg Arg) (*AssignedLandingCallsLandingCalls, error) {
	s := &AssignedLandingCallsLandingCalls{
		sequenceElements: []Element{
			NewElement("floorNumber", V2E(NewUnsigned8), WithElementContext(0)),
			NewElement("direction", V2E(NewLiftCarDirection), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AssignedLandingCalls struct {
	*Sequence
	sequenceElements []Element
}

func NewAssignedLandingCalls(arg Arg) (*AssignedLandingCalls, error) {
	s := &AssignedLandingCalls{
		sequenceElements: []Element{
			NewElement("landingCalls", SequenceOf(NewAssignedLandingCallsLandingCalls), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditNotification struct {
	*Sequence
	sequenceElements []Element
}

func NewAuditNotification(arg Arg) (*AuditNotification, error) {
	s := &AuditNotification{
		sequenceElements: []Element{
			NewElement("sourceTimestamp", V2E(NewTimeStamp), WithElementContext(0), WithElementOptional(true)),
			NewElement("targetTimestamp", V2E(NewTimeStamp), WithElementContext(1), WithElementOptional(true)),
			NewElement("sourceDevice", V2E(NewRecipient), WithElementContext(2)),
			NewElement("sourceObject", Vs2E(NewObjectIdentifier), WithElementContext(3), WithElementOptional(true)),
			NewElement("operation", V2E(NewAuditOperation), WithElementContext(4)),
			NewElement("sourceComment", V2E(NewCharacterString), WithElementContext(5), WithElementOptional(true)),
			NewElement("targetComment", V2E(NewCharacterString), WithElementContext(6), WithElementOptional(true)),
			NewElement("invokeID", V2E(NewUnsigned8), WithElementContext(7), WithElementOptional(true)),
			NewElement("sourceUserID", V2E(NewUnsigned16), WithElementContext(8), WithElementOptional(true)),
			NewElement("sourceUserRole", V2E(NewUnsigned8), WithElementContext(9), WithElementOptional(true)),
			NewElement("targetDevice", V2E(NewRecipient), WithElementContext(10)),
			NewElement("targetObject", Vs2E(NewObjectIdentifier), WithElementContext(11), WithElementOptional(true)),
			NewElement("targetProperty", V2E(NewPropertyReference), WithElementContext(12), WithElementOptional(true)),
			NewElement("targetPriority", V2E(NewUnsigned), WithElementContext(13), WithElementOptional(true)), //  1..16
			NewElement("targetValue", Vs2E(NewAny), WithElementContext(14), WithElementOptional(true)),
			NewElement("currentValue", Vs2E(NewAny), WithElementContext(15), WithElementOptional(true)),
			NewElement("result", V2E(NewErrorType), WithElementContext(16), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditLogRecordLogDatum struct {
	*Choice
	choiceElements []Element
}

func NewAuditLogRecordLogDatum(arg Arg) (*AuditLogRecordLogDatum, error) {
	s := &AuditLogRecordLogDatum{
		choiceElements: []Element{
			NewElement("logStatus", V2E(NewLogStatus), WithElementContext(0)),
			NewElement("auditNotification", V2E(NewAuditNotification), WithElementContext(1)),
			NewElement("timeChange", V2E(NewReal)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditLogRecord struct {
	*Sequence
	sequenceElements []Element
}

func NewAuditLogRecord(arg Arg) (*AuditLogRecord, error) {
	s := &AuditLogRecord{
		sequenceElements: []Element{
			NewElement("timestamp", V2E(NewDateTime), WithElementContext(0)),
			NewElement("logDatum", V2E(NewAuditLogRecordLogDatum), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditLogRecordResult struct {
	*Sequence
	sequenceElements []Element
}

func NewAuditLogRecordResult(arg Arg) (*AuditLogRecordResult, error) {
	s := &AuditLogRecordResult{
		sequenceElements: []Element{
			NewElement("sequenceNumber", V2E(NewUnsigned), WithElementContext(0)), //  Unsigned64
			NewElement("logRecord", V2E(NewAuditLogRecord), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditLogQueryParametersByTarget struct {
	*Sequence
	sequenceElements []Element
}

func NewAuditLogQueryParametersByTarget(arg Arg) (*AuditLogQueryParametersByTarget, error) {
	s := &AuditLogQueryParametersByTarget{
		sequenceElements: []Element{
			NewElement("targetDeviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("targetDeviceAddress", V2E(NewAddress), WithElementContext(1), WithElementOptional(true)),
			NewElement("targetObjectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(2), WithElementOptional(true)),
			NewElement("targetPropertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(3), WithElementOptional(true)),
			NewElement("targetArrayIndex", V2E(NewUnsigned), WithElementContext(4), WithElementOptional(true)),
			NewElement("targetPriority", V2E(NewUnsigned), WithElementContext(5), WithElementOptional(true)),
			NewElement("operations", V2E(NewAuditOperationFlags), WithElementContext(6), WithElementOptional(true)),
			NewElement("successfulActionsOnly", V2E(NewBoolean), WithElementContext(7)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditLogQueryParametersBySource struct {
	*Sequence
	sequenceElements []Element
}

func NewAuditLogQueryParametersBySource(arg Arg) (*AuditLogQueryParametersBySource, error) {
	s := &AuditLogQueryParametersBySource{
		sequenceElements: []Element{
			NewElement("sourceDeviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("sourceDeviceAddress", V2E(NewAddress), WithElementContext(1), WithElementOptional(true)),
			NewElement("sourceObjectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(2), WithElementOptional(true)),
			NewElement("operations", V2E(NewAuditOperationFlags), WithElementContext(3), WithElementOptional(true)),
			NewElement("successfulActionsOnly", V2E(NewBoolean), WithElementContext(4)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuditLogQueryParameters struct {
	*Choice
	choiceElements []Element
}

func NewAuditLogQueryParameters(arg Arg) (*AuditLogQueryParameters, error) {
	s := &AuditLogQueryParameters{
		choiceElements: []Element{
			NewElement("byTarget", V2E(NewAuditLogQueryParametersByTarget), WithElementContext(0)),
			NewElement("bySource", V2E(NewAuditLogQueryParametersBySource), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuthenticationFactor struct {
	*Sequence
	sequenceElements []Element
}

func NewAuthenticationFactor(arg Arg) (*AuthenticationFactor, error) {
	s := &AuthenticationFactor{
		sequenceElements: []Element{
			NewElement("formatType", V2E(NewAuthenticationFactorType), WithElementContext(0)),
			NewElement("formatClass", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("value", V2E(NewOctetString), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuthenticationFactorFormat struct {
	*Sequence
	sequenceElements []Element
}

func NewAuthenticationFactorFormat(arg Arg) (*AuthenticationFactorFormat, error) {
	s := &AuthenticationFactorFormat{
		sequenceElements: []Element{
			NewElement("formatType", V2E(NewAuthenticationFactorType), WithElementContext(0)),
			NewElement("vendorId", V2E(NewUnsigned), WithElementContext(1), WithElementOptional(true)),
			NewElement("vendorFormat", V2E(NewUnsigned), WithElementContext(2), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuthenticationPolicyPolicy struct {
	*Sequence
	sequenceElements []Element
}

func NewAuthenticationPolicyPolicy(arg Arg) (*AuthenticationPolicyPolicy, error) {
	s := &AuthenticationPolicyPolicy{
		sequenceElements: []Element{
			NewElement("credentialDataInput", V2E(NewDeviceObjectReference), WithElementContext(0)),
			NewElement("index", V2E(NewUnsigned), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type AuthenticationPolicy struct {
	*Sequence
	sequenceElements []Element
}

func NewAuthenticationPolicy(arg Arg) (*AuthenticationPolicy, error) {
	s := &AuthenticationPolicy{
		sequenceElements: []Element{
			NewElement("policy", SequenceOf(NewAuthenticationPolicyPolicy), WithElementContext(0)),
			NewElement("orderEnforced", V2E(NewBoolean), WithElementContext(1)),
			NewElement("timeout", V2E(NewUnsigned), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type CalendarEntry struct {
	*Choice
	choiceElements []Element
}

func NewCalendarEntry(arg Arg) (*CalendarEntry, error) {
	s := &CalendarEntry{
		choiceElements: []Element{
			NewElement("date", Vs2E(NewDate), WithElementContext(0)),
			NewElement("dateRange", V2E(NewDateRange), WithElementContext(1)),
			NewElement("weekNDay", V2E(NewWeekNDay), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ChannelValue struct {
	*Choice
	choiceElements []Element
}

func NewChannelValue(arg Arg) (*ChannelValue, error) {
	s := &ChannelValue{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("real", V2E(NewReal)),
			NewElement("enumerated", Vs2E(NewEnumerated)),
			NewElement("unsigned", V2E(NewUnsigned)),
			NewElement("boolean", V2E(NewBoolean)),
			NewElement("integer", V2E(NewInteger)),
			NewElement("double", V2E(NewDouble)),
			NewElement("time", Vs2E(NewTime)),
			NewElement("characterString", V2E(NewCharacterString)),
			NewElement("octetString", V2E(NewOctetString)),
			NewElement("bitString", Vs2E(NewBitString)),
			NewElement("date", Vs2E(NewDate)),
			NewElement("objectidentifier", Vs2E(NewObjectIdentifier)),
			NewElement("lightingCommand", V2E(NewLightingCommand), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ClientCOV struct {
	*Choice
	choiceElements []Element
}

func NewClientCOV(arg Arg) (*ClientCOV, error) {
	s := &ClientCOV{
		choiceElements: []Element{
			NewElement("realIncrement", V2E(NewReal)),
			NewElement("defaultIncrement", V2E(NewNull)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type COVMultipleSubscriptionListOfCOVReference struct {
	*Sequence
	sequenceElements []Element
}

func NewCOVMultipleSubscriptionListOfCOVReference(arg Arg) (*COVMultipleSubscriptionListOfCOVReference, error) {
	s := &COVMultipleSubscriptionListOfCOVReference{
		sequenceElements: []Element{
			NewElement("monitoredProperty", V2E(NewPropertyReference), WithElementContext(0)),
			NewElement("covIncrement", V2E(NewReal), WithElementContext(1), WithElementOptional(true)),
			NewElement("timestamped", V2E(NewBoolean), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type COVMultipleSubscriptionList struct {
	*Sequence
	sequenceElements []Element
}

func NewCOVMultipleSubscriptionList(arg Arg) (*COVMultipleSubscriptionList, error) {
	s := &COVMultipleSubscriptionList{
		sequenceElements: []Element{
			NewElement("monitoredObjectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("listOfCOVReferences", SequenceOf(NewCOVMultipleSubscriptionListOfCOVReference), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type COVMultipleSubscription struct {
	*Sequence
	sequenceElements []Element
}

func NewCOVMultipleSubscription(arg Arg) (*COVMultipleSubscription, error) {
	s := &COVMultipleSubscription{
		sequenceElements: []Element{
			NewElement("recipient", V2E(NewRecipientProcess), WithElementContext(0)),
			NewElement("issueConfirmedNotifications", V2E(NewBoolean), WithElementContext(1)),
			NewElement("timeRemaining", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("maxNotificationDelay", V2E(NewUnsigned), WithElementContext(3)),
			NewElement("listOfCOVSubscriptionSpecifications", SequenceOf(NewCOVMultipleSubscriptionList), WithElementContext(4)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type COVSubscription struct {
	*Sequence
	sequenceElements []Element
}

func NewCOVSubscription(arg Arg) (*COVSubscription, error) {
	s := &COVSubscription{
		sequenceElements: []Element{
			NewElement("recipient", V2E(NewRecipientProcess), WithElementContext(0)),
			NewElement("monitoredPropertyReference", V2E(NewObjectPropertyReference), WithElementContext(1)),
			NewElement("issueConfirmedNotifications", V2E(NewBoolean), WithElementContext(2)),
			NewElement("timeRemaining", V2E(NewUnsigned), WithElementContext(3)),
			NewElement("covIncrement", V2E(NewReal), WithElementContext(4), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type CredentialAuthenticationFactor struct {
	*Sequence
	sequenceElements []Element
}

func NewCredentialAuthenticationFactor(arg Arg) (*CredentialAuthenticationFactor, error) {
	s := &CredentialAuthenticationFactor{
		sequenceElements: []Element{
			NewElement("disable", V2E(NewAccessAuthenticationFactorDisable), WithElementContext(0)),
			NewElement("authenticationFactor", V2E(NewAuthenticationFactor), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DailySchedule struct {
	*Sequence
	sequenceElements []Element
}

func NewDailySchedule(arg Arg) (*DailySchedule, error) {
	s := &DailySchedule{
		sequenceElements: []Element{
			NewElement("daySchedule", SequenceOf(NewTimeValue), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type Destination struct {
	*Sequence
	sequenceElements []Element
}

func NewDestination(arg Arg) (*Destination, error) {
	s := &Destination{
		sequenceElements: []Element{
			NewElement("validDays", V2E(NewDaysOfWeek)),
			NewElement("fromTime", Vs2E(NewTime)),
			NewElement("toTime", Vs2E(NewTime)),
			NewElement("recipient", V2E(NewRecipient)),
			NewElement("processIdentifier", V2E(NewUnsigned)),
			NewElement("issueConfirmedNotifications", V2E(NewBoolean)),
			NewElement("transitions", V2E(NewEventTransitionBits)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type DeviceObjectPropertyValue struct {
	*Sequence
	sequenceElements []Element
}

func NewDeviceObjectPropertyValue(arg Arg) (*DeviceObjectPropertyValue, error) {
	s := &DeviceObjectPropertyValue{
		sequenceElements: []Element{
			NewElement("deviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(1)),
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(2)),
			NewElement("arrayIndex", V2E(NewUnsigned), WithElementContext(3), WithElementOptional(true)),
			NewElement("value", Vs2E(NewAny), WithElementContext(4)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventNotificationSubscription struct {
	*Sequence
	sequenceElements []Element
}

func NewEventNotificationSubscription(arg Arg) (*EventNotificationSubscription, error) {
	s := &EventNotificationSubscription{
		sequenceElements: []Element{
			NewElement("recipient", V2E(NewRecipient), WithElementContext(0)),
			NewElement("processIdentifier", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("issueConfirmedNotifications", V2E(NewBoolean), WithElementContext(2)),
			NewElement("timeRemaining", V2E(NewUnsigned), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfBitstring struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterChangeOfBitstring(arg Arg) (*EventParameterChangeOfBitstring, error) {
	s := &EventParameterChangeOfBitstring{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("bitmask", Vs2E(NewBitString), WithElementContext(1)),
			NewElement("listOfBitstringValues", SequenceOfs(NewBitString), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfState struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterChangeOfState(arg Arg) (*EventParameterChangeOfState, error) {
	s := &EventParameterChangeOfState{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("listOfValues", SequenceOf(NewPropertyStates), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfValueCOVCriteria struct {
	*Choice
	choiceElements []Element
}

func NewEventParameterChangeOfValueCOVCriteria(arg Arg) (*EventParameterChangeOfValueCOVCriteria, error) {
	s := &EventParameterChangeOfValueCOVCriteria{
		choiceElements: []Element{
			NewElement("bitmask", Vs2E(NewBitString), WithElementContext(0)),
			NewElement("referencedPropertyIncrement", V2E(NewReal), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfValue struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterChangeOfValue(arg Arg) (*EventParameterChangeOfValue, error) {
	s := &EventParameterChangeOfValue{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("covCriteria", V2E(NewEventParameterChangeOfValueCOVCriteria), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterCommandFailure struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterCommandFailure(arg Arg) (*EventParameterCommandFailure, error) {
	s := &EventParameterCommandFailure{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("feedbackPropertyReference", V2E(NewDeviceObjectPropertyReference), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterFloatingLimit struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterFloatingLimit(arg Arg) (*EventParameterFloatingLimit, error) {
	s := &EventParameterFloatingLimit{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("setpointReference", V2E(NewDeviceObjectPropertyReference), WithElementContext(1)),
			NewElement("lowDiffLimit", V2E(NewReal), WithElementContext(2)),
			NewElement("highDiffLimit", V2E(NewReal), WithElementContext(3)),
			NewElement("deadband", V2E(NewReal), WithElementContext(4)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterOutOfRange struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterOutOfRange(arg Arg) (*EventParameterOutOfRange, error) {
	s := &EventParameterOutOfRange{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("lowLimit", V2E(NewReal), WithElementContext(1)),
			NewElement("highLimit", V2E(NewReal), WithElementContext(2)),
			NewElement("deadband", V2E(NewReal), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfLifeSafety struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterChangeOfLifeSafety(arg Arg) (*EventParameterChangeOfLifeSafety, error) {
	s := &EventParameterChangeOfLifeSafety{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("listOfLifeSafetyAlarmValues", SequenceOf(NewLifeSafetyState), WithElementContext(1)),
			NewElement("listOfAlarmValues", SequenceOf(NewLifeSafetyState), WithElementContext(2)),
			NewElement("modePropertyReference", V2E(NewDeviceObjectPropertyReference), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterExtendedParameters struct {
	*Choice
	choiceElements []Element
}

func NewEventParameterExtendedParameters(arg Arg) (*EventParameterExtendedParameters, error) {
	s := &EventParameterExtendedParameters{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull), WithElementContext(0)),
			NewElement("real", V2E(NewReal), WithElementContext(1)),
			NewElement("integer", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("boolean", V2E(NewBoolean), WithElementContext(3)),
			NewElement("double", V2E(NewDouble), WithElementContext(4)),
			NewElement("octet", V2E(NewOctetString), WithElementContext(5)),
			NewElement("bitstring", Vs2E(NewBitString), WithElementContext(6)),
			NewElement("enum", Vs2E(NewEnumerated), WithElementContext(7)),
			NewElement("reference", V2E(NewDeviceObjectPropertyReference), WithElementContext(8)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterExtended struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterExtended(arg Arg) (*EventParameterExtended, error) {
	s := &EventParameterExtended{
		sequenceElements: []Element{
			NewElement("vendorId", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("extendedEventType", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("parameters", SequenceOf(NewEventParameterExtendedParameters), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterBufferReady struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterBufferReady(arg Arg) (*EventParameterBufferReady, error) {
	s := &EventParameterBufferReady{
		sequenceElements: []Element{
			NewElement("notificationThreshold", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("previousNotificationCount", V2E(NewUnsigned), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterUnsignedRange struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterUnsignedRange(arg Arg) (*EventParameterUnsignedRange, error) {
	s := &EventParameterUnsignedRange{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("lowLimit", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("highLimit", V2E(NewUnsigned), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterAccessEventAccessEvent struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterAccessEventAccessEvent(arg Arg) (*EventParameterAccessEventAccessEvent, error) {
	s := &EventParameterAccessEventAccessEvent{
		sequenceElements: []Element{
			NewElement("listOfAccessEvents", SequenceOf(NewAccessEvent), WithElementContext(0)),
			NewElement("accessEventTimeReference", V2E(NewDeviceObjectPropertyReference), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterAccessEvent struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterAccessEvent(arg Arg) (*EventParameterAccessEvent, error) {
	s := &EventParameterAccessEvent{
		sequenceElements: []Element{
			NewElement("accessEvent", SequenceOf(NewEventParameterAccessEventAccessEvent), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterDoubleOutOfRange struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterDoubleOutOfRange(arg Arg) (*EventParameterDoubleOutOfRange, error) {
	s := &EventParameterDoubleOutOfRange{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("lowLimit", V2E(NewDouble), WithElementContext(1)),
			NewElement("highLimit", V2E(NewDouble), WithElementContext(2)),
			NewElement("deadband", V2E(NewDouble), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterSignedOutOfRange struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterSignedOutOfRange(arg Arg) (*EventParameterSignedOutOfRange, error) {
	s := &EventParameterSignedOutOfRange{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("lowLimit", V2E(NewInteger), WithElementContext(1)),
			NewElement("highLimit", V2E(NewInteger), WithElementContext(2)),
			NewElement("deadband", V2E(NewUnsigned), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterUnsignedOutOfRange struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterUnsignedOutOfRange(arg Arg) (*EventParameterUnsignedOutOfRange, error) {
	s := &EventParameterUnsignedOutOfRange{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("lowLimit", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("highLimit", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("deadband", V2E(NewUnsigned), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfCharacterString struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterChangeOfCharacterString(arg Arg) (*EventParameterChangeOfCharacterString, error) {
	s := &EventParameterChangeOfCharacterString{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("listOfAlarmValues", SequenceOf(NewCharacterString), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameterChangeOfStatusFlags struct {
	*Sequence
	sequenceElements []Element
}

func NewEventParameterChangeOfStatusFlags(arg Arg) (*EventParameterChangeOfStatusFlags, error) {
	s := &EventParameterChangeOfStatusFlags{
		sequenceElements: []Element{
			NewElement("timeDelay", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("selectedFlags", V2E(NewStatusFlags), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type EventParameter struct {
	*Choice
	choiceElements []Element
}

func NewEventParameter(arg Arg) (*EventParameter, error) {
	s := &EventParameter{
		choiceElements: []Element{
			NewElement("changeOfBitstring", V2E(NewEventParameterChangeOfBitstring), WithElementContext(0)),
			NewElement("changeOfState", V2E(NewEventParameterChangeOfState), WithElementContext(1)),
			NewElement("changeOfValue", V2E(NewEventParameterChangeOfValue), WithElementContext(2)),
			NewElement("commandFailure", V2E(NewEventParameterCommandFailure), WithElementContext(3)),
			NewElement("floatingLimit", V2E(NewEventParameterFloatingLimit), WithElementContext(4)),
			NewElement("outOfRange", V2E(NewEventParameterOutOfRange), WithElementContext(5)),
			NewElement("changeOfLifesafety", V2E(NewEventParameterChangeOfLifeSafety), WithElementContext(8)),
			NewElement("extended", V2E(NewEventParameterExtended), WithElementContext(9)),
			NewElement("bufferReady", V2E(NewEventParameterBufferReady), WithElementContext(10)),
			NewElement("unsignedRange", V2E(NewEventParameterUnsignedRange), WithElementContext(11)),
			NewElement("accessEvent", V2E(NewEventParameterAccessEvent), WithElementContext(13)),
			NewElement("doubleOutOfRange", V2E(NewEventParameterDoubleOutOfRange), WithElementContext(14)),
			NewElement("signedOutOfRange", V2E(NewEventParameterSignedOutOfRange), WithElementContext(15)),
			NewElement("unsignedOutOfRange", V2E(NewEventParameterUnsignedOutOfRange), WithElementContext(16)),
			NewElement("changeOfCharacterstring", V2E(NewEventParameterChangeOfCharacterString), WithElementContext(17)),
			NewElement("changeOfStatusflags", V2E(NewEventParameterChangeOfStatusFlags), WithElementContext(18)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameterCharacterString struct {
	*Sequence
	sequenceElements []Element
}

func NewFaultParameterCharacterString(arg Arg) (*FaultParameterCharacterString, error) {
	s := &FaultParameterCharacterString{
		sequenceElements: []Element{
			NewElement("listOfFaultValues", SequenceOf(NewCharacterString), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameterExtendedParameters struct {
	*Choice
	choiceElements []Element
}

func NewFaultParameterExtendedParameters(arg Arg) (*FaultParameterExtendedParameters, error) {
	s := &FaultParameterExtendedParameters{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("real", V2E(NewReal)),
			NewElement("unsigned", V2E(NewUnsigned)),
			NewElement("boolean", V2E(NewBoolean)),
			NewElement("integer", V2E(NewInteger)),
			NewElement("double", V2E(NewDouble)),
			NewElement("octet", V2E(NewOctetString)),
			NewElement("characterString", V2E(NewCharacterString)),
			NewElement("bitstring", Vs2E(NewBitString)),
			NewElement("enum", Vs2E(NewEnumerated)),
			NewElement("date", Vs2E(NewDate)),
			NewElement("time", Vs2E(NewTime)),
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier)),
			NewElement("reference", V2E(NewDeviceObjectPropertyReference), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameterExtended struct {
	*Sequence
	sequenceElements []Element
}

func NewFaultParameterExtended(arg Arg) (*FaultParameterExtended, error) {
	s := &FaultParameterExtended{
		sequenceElements: []Element{
			NewElement("vendorId", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("extendedFaultType", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("parameters", SequenceOf(NewFaultParameterExtendedParameters), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameterLifeSafety struct {
	*Sequence
	sequenceElements []Element
}

func NewFaultParameterLifeSafety(arg Arg) (*FaultParameterLifeSafety, error) {
	s := &FaultParameterLifeSafety{
		sequenceElements: []Element{
			NewElement("listOfFaultValues", SequenceOf(NewLifeSafetyState), WithElementContext(0)),
			NewElement("modePropertyReference", V2E(NewDeviceObjectPropertyReference), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameterState struct {
	*Sequence
	sequenceElements []Element
}

func NewFaultParameterState(arg Arg) (*FaultParameterState, error) {
	s := &FaultParameterState{
		sequenceElements: []Element{
			NewElement("listOfFaultValues", SequenceOf(NewPropertyStates), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameterStatusFlags struct {
	*Sequence
	sequenceElements []Element
}

func NewFaultParameterStatusFlags(arg Arg) (*FaultParameterStatusFlags, error) {
	s := &FaultParameterStatusFlags{
		sequenceElements: []Element{
			NewElement("statusFlagsReference", V2E(NewDeviceObjectPropertyReference), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type FaultParameter struct {
	*Choice
	choiceElements []Element
}

func NewFaultParameter(arg Arg) (*FaultParameter, error) {
	s := &FaultParameter{
		choiceElements: []Element{
			NewElement("none", V2E(NewNull), WithElementContext(0)),
			NewElement("faultCharacterString", V2E(NewFaultParameterCharacterString), WithElementContext(1)),
			NewElement("faultExtended", V2E(NewFaultParameterExtended), WithElementContext(2)),
			NewElement("faultLifeSafety", V2E(NewFaultParameterLifeSafety), WithElementContext(3)),
			NewElement("faultState", V2E(NewFaultParameterState), WithElementContext(4)),
			NewElement("faultStatusFlags", V2E(NewFaultParameterStatusFlags), WithElementContext(5)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type KeyIdentifier struct {
	*Sequence
	sequenceElements []Element
}

func NewKeyIdentifier(arg Arg) (*KeyIdentifier, error) {
	s := &KeyIdentifier{
		sequenceElements: []Element{
			NewElement("algorithm", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("keyId", V2E(NewUnsigned), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LogDataLogData struct {
	*Choice
	choiceElements []Element
}

func NewLogDataLogData(arg Arg) (*LogDataLogData, error) {
	s := &LogDataLogData{
		choiceElements: []Element{
			NewElement("booleanValue", V2E(NewBoolean), WithElementContext(0)),
			NewElement("realValue", V2E(NewReal), WithElementContext(1)),
			NewElement("enumValue", Vs2E(NewEnumerated), WithElementContext(2)),
			NewElement("unsignedValue", V2E(NewUnsigned), WithElementContext(3)),
			NewElement("signedValue", V2E(NewInteger), WithElementContext(4)),
			NewElement("bitstringValue", Vs2E(NewBitString), WithElementContext(5)),
			NewElement("nullValue", V2E(NewNull), WithElementContext(6)),
			NewElement("failure", V2E(NewErrorType), WithElementContext(7)),
			NewElement("anyValue", Vs2E(NewAny), WithElementContext(8)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LogData struct {
	*Choice
	choiceElements []Element
}

func NewLogData(arg Arg) (*LogData, error) {
	s := &LogData{
		choiceElements: []Element{
			NewElement("logStatus", V2E(NewLogStatus), WithElementContext(0)),
			NewElement("logData", SequenceOf(NewLogDataLogData), WithElementContext(1)),
			NewElement("timeChange", V2E(NewReal), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LogMultipleRecord struct {
	*Sequence
	sequenceElements []Element
}

func NewLogMultipleRecord(arg Arg) (*LogMultipleRecord, error) {
	s := &LogMultipleRecord{
		sequenceElements: []Element{
			NewElement("timestamp", V2E(NewDateTime), WithElementContext(0)),
			NewElement("logData", V2E(NewLogData), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LogRecordLogDatum struct {
	*Choice
	choiceElements []Element
}

func NewLogRecordLogDatum(arg Arg) (*LogRecordLogDatum, error) {
	s := &LogRecordLogDatum{
		choiceElements: []Element{
			NewElement("logStatus", V2E(NewLogStatus), WithElementContext(0)),
			NewElement("booleanValue", V2E(NewBoolean), WithElementContext(1)),
			NewElement("realValue", V2E(NewReal), WithElementContext(2)),
			NewElement("enumValue", Vs2E(NewEnumerated), WithElementContext(3)),
			NewElement("unsignedValue", V2E(NewUnsigned), WithElementContext(4)),
			NewElement("signedValue", V2E(NewInteger), WithElementContext(5)),
			NewElement("bitstringValue", Vs2E(NewBitString), WithElementContext(6)),
			NewElement("nullValue", V2E(NewNull), WithElementContext(7)),
			NewElement("failure", V2E(NewErrorType), WithElementContext(8)),
			NewElement("timeChange", V2E(NewReal), WithElementContext(9)),
			NewElement("anyValue", Vs2E(NewAny), WithElementContext(10)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type LogRecord struct {
	*Sequence
	sequenceElements []Element
}

func NewLogRecord(arg Arg) (*LogRecord, error) {
	s := &LogRecord{
		sequenceElements: []Element{
			NewElement("timestamp", V2E(NewDateTime), WithElementContext(0)),
			NewElement("logDatum", V2E(NewLogRecordLogDatum), WithElementContext(1)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(2), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NetworkSecurityPolicy struct {
	*Sequence
	sequenceElements []Element
}

func NewNetworkSecurityPolicy(arg Arg) (*NetworkSecurityPolicy, error) {
	s := &NetworkSecurityPolicy{
		sequenceElements: []Element{
			NewElement("portId", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("securityLevel", V2E(NewSecurityPolicy), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfBitstring struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfBitstring(arg Arg) (*NotificationParametersChangeOfBitstring, error) {
	s := &NotificationParametersChangeOfBitstring{
		sequenceElements: []Element{
			NewElement("referencedBitstring", Vs2E(NewBitString), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfState struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfState(arg Arg) (*NotificationParametersChangeOfState, error) {
	s := &NotificationParametersChangeOfState{
		sequenceElements: []Element{
			NewElement("newState", V2E(NewPropertyStates), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfValueNewValue struct {
	*Choice
	choiceElements []Element
}

func NewNotificationParametersChangeOfValueNewValue(arg Arg) (*NotificationParametersChangeOfValueNewValue, error) {
	s := &NotificationParametersChangeOfValueNewValue{
		choiceElements: []Element{
			NewElement("changedBits", Vs2E(NewBitString), WithElementContext(0)),
			NewElement("changedValue", V2E(NewReal), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfValue struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfValue(arg Arg) (*NotificationParametersChangeOfValue, error) {
	s := &NotificationParametersChangeOfValue{
		sequenceElements: []Element{
			NewElement("newValue", V2E(NewNotificationParametersChangeOfValueNewValue), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersCommandFailure struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersCommandFailure(arg Arg) (*NotificationParametersCommandFailure, error) {
	s := &NotificationParametersCommandFailure{
		sequenceElements: []Element{
			NewElement("commandValue", Vs2E(NewAny), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("feedbackValue", Vs2E(NewAny), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersFloatingLimit struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersFloatingLimit(arg Arg) (*NotificationParametersFloatingLimit, error) {
	s := &NotificationParametersFloatingLimit{
		sequenceElements: []Element{
			NewElement("referenceValue", V2E(NewReal), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("setpointValue", V2E(NewReal), WithElementContext(2)),
			NewElement("errorLimit", V2E(NewReal), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersOutOfRange struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersOutOfRange(arg Arg) (*NotificationParametersOutOfRange, error) {
	s := &NotificationParametersOutOfRange{
		sequenceElements: []Element{
			NewElement("exceedingValue", V2E(NewReal), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("deadband", V2E(NewReal), WithElementContext(2)),
			NewElement("exceededLimit", V2E(NewReal), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersExtendedParametersType struct {
	*Choice
	choiceElements []Element
}

func NewNotificationParametersExtendedParametersType(arg Arg) (*NotificationParametersExtendedParametersType, error) {
	s := &NotificationParametersExtendedParametersType{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("real", V2E(NewReal)),
			NewElement("integer", V2E(NewUnsigned)),
			NewElement("boolean", V2E(NewBoolean)),
			NewElement("double", V2E(NewDouble)),
			NewElement("octet", V2E(NewOctetString)),
			NewElement("bitstring", Vs2E(NewBitString)),
			NewElement("enum", Vs2E(NewEnumerated)),
			NewElement("propertyValue", V2E(NewDeviceObjectPropertyValue)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersExtended struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersExtended(arg Arg) (*NotificationParametersExtended, error) {
	s := &NotificationParametersExtended{
		sequenceElements: []Element{
			NewElement("vendorId", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("extendedEventType", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("parameters", V2E(NewNotificationParametersExtendedParametersType), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersBufferReady struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersBufferReady(arg Arg) (*NotificationParametersBufferReady, error) {
	s := &NotificationParametersBufferReady{
		sequenceElements: []Element{
			NewElement("bufferProperty", V2E(NewDeviceObjectPropertyReference), WithElementContext(0)),
			NewElement("previousNotification", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("currentNotification", V2E(NewUnsigned), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersUnsignedRange struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersUnsignedRange(arg Arg) (*NotificationParametersUnsignedRange, error) {
	s := &NotificationParametersUnsignedRange{
		sequenceElements: []Element{
			NewElement("exceedingValue", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("exceedingLimit", V2E(NewUnsigned), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersComplexEventType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersComplexEventType(arg Arg) (*NotificationParametersComplexEventType, error) {
	s := &NotificationParametersComplexEventType{
		sequenceElements: []Element{
			NewElement("complexEventType", V2E(NewPropertyValue), WithElementContext(0)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfLifeSafety struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfLifeSafety(arg Arg) (*NotificationParametersChangeOfLifeSafety, error) {
	s := &NotificationParametersChangeOfLifeSafety{
		sequenceElements: []Element{
			NewElement("newState", V2E(NewLifeSafetyState), WithElementContext(0)),
			NewElement("newMode", V2E(NewLifeSafetyMode), WithElementContext(1)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(2)),
			NewElement("operationExpected", V2E(NewLifeSafetyOperation), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersAccessEventType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersAccessEventType(arg Arg) (*NotificationParametersAccessEventType, error) {
	s := &NotificationParametersAccessEventType{
		sequenceElements: []Element{
			NewElement("accessEvent", V2E(NewAccessEvent), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("accessEventTag", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("accessEventTime", V2E(NewTimeStamp), WithElementContext(3)),
			NewElement("accessCredential", V2E(NewDeviceObjectReference), WithElementContext(4)),
			NewElement("authenicationFactor", V2E(NewAuthenticationFactorType), WithElementContext(5), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersDoubleOutOfRangeType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersDoubleOutOfRangeType(arg Arg) (*NotificationParametersDoubleOutOfRangeType, error) {
	s := &NotificationParametersDoubleOutOfRangeType{
		sequenceElements: []Element{
			NewElement("exceedingValue", V2E(NewDouble), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("deadband", V2E(NewDouble), WithElementContext(2)),
			NewElement("exceededLimit", V2E(NewDouble), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersSignedOutOfRangeType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersSignedOutOfRangeType(arg Arg) (*NotificationParametersSignedOutOfRangeType, error) {
	s := &NotificationParametersSignedOutOfRangeType{
		sequenceElements: []Element{
			NewElement("exceedingValue", V2E(NewInteger), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("deadband", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("exceededLimit", V2E(NewInteger), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersUnsignedOutOfRangeType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersUnsignedOutOfRangeType(arg Arg) (*NotificationParametersUnsignedOutOfRangeType, error) {
	s := &NotificationParametersUnsignedOutOfRangeType{
		sequenceElements: []Element{
			NewElement("exceedingValue", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("deadband", V2E(NewUnsigned), WithElementContext(2)),
			NewElement("exceededLimit", V2E(NewUnsigned), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfCharacterStringType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfCharacterStringType(arg Arg) (*NotificationParametersChangeOfCharacterStringType, error) {
	s := &NotificationParametersChangeOfCharacterStringType{
		sequenceElements: []Element{
			NewElement("changedValue", V2E(NewCharacterString), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("alarmValue", V2E(NewCharacterString), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfStatusFlagsType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfStatusFlagsType(arg Arg) (*NotificationParametersChangeOfStatusFlagsType, error) {
	s := &NotificationParametersChangeOfStatusFlagsType{
		sequenceElements: []Element{
			NewElement("presentValue", V2E(NewCharacterString), WithElementContext(0)),
			NewElement("referencedFlags", V2E(NewStatusFlags), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParametersChangeOfReliabilityType struct {
	*Sequence
	sequenceElements []Element
}

func NewNotificationParametersChangeOfReliabilityType(arg Arg) (*NotificationParametersChangeOfReliabilityType, error) {
	s := &NotificationParametersChangeOfReliabilityType{
		sequenceElements: []Element{
			NewElement("reliability", V2E(NewReliability), WithElementContext(0)),
			NewElement("statusFlags", V2E(NewStatusFlags), WithElementContext(1)),
			NewElement("propertyValues", SequenceOf(NewPropertyValue), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type NotificationParameters struct {
	*Choice
	choiceElements []Element
}

func NewNotificationParameters(arg Arg) (*NotificationParameters, error) {
	s := &NotificationParameters{
		choiceElements: []Element{
			NewElement("changeOfBitstring", V2E(NewNotificationParametersChangeOfBitstring), WithElementContext(0)),
			NewElement("changeOfState", V2E(NewNotificationParametersChangeOfState), WithElementContext(1)),
			NewElement("changeOfValue", V2E(NewNotificationParametersChangeOfValue), WithElementContext(2)),
			NewElement("commandFailure", V2E(NewNotificationParametersCommandFailure), WithElementContext(3)),
			NewElement("floatingLimit", V2E(NewNotificationParametersFloatingLimit), WithElementContext(4)),
			NewElement("outOfRange", V2E(NewNotificationParametersOutOfRange), WithElementContext(5)),
			NewElement("complexEventType", V2E(NewNotificationParametersComplexEventType), WithElementContext(6)),
			NewElement("changeOfLifeSafety", V2E(NewNotificationParametersChangeOfLifeSafety), WithElementContext(8)),
			NewElement("extended", V2E(NewNotificationParametersExtended), WithElementContext(9)),
			NewElement("bufferReady", V2E(NewNotificationParametersBufferReady), WithElementContext(10)),
			NewElement("unsignedRange", V2E(NewNotificationParametersUnsignedRange), WithElementContext(11)),
			NewElement("accessEvent", V2E(NewNotificationParametersAccessEventType), WithElementContext(13)),
			NewElement("doubleOutOfRange", V2E(NewNotificationParametersDoubleOutOfRangeType), WithElementContext(14)),
			NewElement("signedOutOfRange", V2E(NewNotificationParametersSignedOutOfRangeType), WithElementContext(15)),
			NewElement("unsignedOutOfRange", V2E(NewNotificationParametersUnsignedOutOfRangeType), WithElementContext(16)),
			NewElement("changeOfCharacterString", V2E(NewNotificationParametersChangeOfCharacterStringType), WithElementContext(17)),
			NewElement("changeOfStatusFlags", V2E(NewNotificationParametersChangeOfStatusFlagsType), WithElementContext(18)),
			NewElement("changeOfReliability", V2E(NewNotificationParametersChangeOfReliabilityType), WithElementContext(19)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ObjectPropertyValue struct {
	*Sequence
	sequenceElements []Element
}

func NewObjectPropertyValue(arg Arg) (*ObjectPropertyValue, error) {
	s := &ObjectPropertyValue{
		sequenceElements: []Element{
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(1)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(2), WithElementOptional(true)),
			NewElement("value", Vs2E(NewAny), WithElementContext(3)),
			NewElement("priority", V2E(NewUnsigned), WithElementContext(4), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ObjectSelector struct {
	*Choice
	choiceElements []Element
}

func NewObjectSelector(arg Arg) (*ObjectSelector, error) {
	s := &ObjectSelector{
		choiceElements: []Element{
			NewElement("none", V2E(NewNull)),
			NewElement("object", Vs2E(NewObjectIdentifier)),
			NewElement("objectType", Vs2E(NewObjectType)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type PortPermission struct {
	*Sequence
	sequenceElements []Element
}

func NewPortPermission(arg Arg) (*PortPermission, error) {
	s := &PortPermission{
		sequenceElements: []Element{
			NewElement("portId", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("enabled", V2E(NewBoolean), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type Prescale struct {
	*Sequence
	sequenceElements []Element
}

func NewPrescale(arg Arg) (*Prescale, error) {
	s := &Prescale{
		sequenceElements: []Element{
			NewElement("multiplier", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("moduloDivide", V2E(NewUnsigned), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type PriorityValue struct {
	*Choice
	choiceElements []Element
}

func NewPriorityValue(arg Arg) (*PriorityValue, error) {
	s := &PriorityValue{
		choiceElements: []Element{
			NewElement("null", V2E(NewNull)),
			NewElement("real", V2E(NewReal)),
			NewElement("enumerated", Vs2E(NewEnumerated)),
			NewElement("unsigned", V2E(NewUnsigned)),
			NewElement("boolean", V2E(NewBoolean)),
			NewElement("integer", V2E(NewInteger)),
			NewElement("double", V2E(NewDouble)),
			NewElement("time", Vs2E(NewTime)),
			NewElement("characterString", V2E(NewCharacterString)),
			NewElement("octetString", V2E(NewOctetString)),
			NewElement("bitString", Vs2E(NewBitString)),
			NewElement("date", Vs2E(NewDate)),
			NewElement("objectidentifier", Vs2E(NewObjectIdentifier)),
			NewElement("constructedValue", Vs2E(NewAny), WithElementContext(0)),
			NewElement("datetime", V2E(NewDateTime), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

// TODO: finish me... dunno what is going on here...
type PriorityArray struct {
}

func NewPriorityArray(arg Arg) (*PriorityArray, error) {
	s := &PriorityArray{}
	ArrayOf(NewPriorityValue, 16, PriorityValue{})
	panic("implementchoice")
	return s, nil
}

type PropertyAccessResultAccessResult struct {
	*Choice
	choiceElements []Element
}

func NewPropertyAccessResultAccessResult(arg Arg) (*PropertyAccessResultAccessResult, error) {
	s := &PropertyAccessResultAccessResult{
		choiceElements: []Element{
			NewElement("propertyValue", Vs2E(NewAny), WithElementContext(4)),
			NewElement("propertyAccessError", V2E(NewErrorType), WithElementContext(5)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type PropertyAccessResult struct {
	*Sequence
	sequenceElements []Element
}

func NewPropertyAccessResult(arg Arg) (*PropertyAccessResult, error) {
	s := &PropertyAccessResult{
		sequenceElements: []Element{
			NewElement("objectIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(0)),
			NewElement("propertyIdentifier", V2E(NewPropertyIdentifier), WithElementContext(1)),
			NewElement("propertyArrayIndex", V2E(NewUnsigned), WithElementContext(2), WithElementOptional(true)),
			NewElement("deviceIdentifier", Vs2E(NewObjectIdentifier), WithElementContext(3), WithElementOptional(true)),
			NewElement("accessResult", V2E(NewPropertyAccessResultAccessResult)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type Scale struct {
	*Choice
	choiceElements []Element
}

func NewScale(arg Arg) (*Scale, error) {
	s := &Scale{
		choiceElements: []Element{
			NewElement("floatScale", V2E(NewReal), WithElementContext(0)),
			NewElement("integerScale", V2E(NewInteger), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type SecurityKeySet struct {
	*Sequence
	sequenceElements []Element
}

func NewSecurityKeySet(arg Arg) (*SecurityKeySet, error) {
	s := &SecurityKeySet{
		sequenceElements: []Element{
			NewElement("keyRevision", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("activationTime", V2E(NewDateTime), WithElementContext(1)),
			NewElement("expirationTime", V2E(NewDateTime), WithElementContext(2)),
			NewElement("keyIds", SequenceOf(NewKeyIdentifier), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ShedLevel struct {
	*Choice
	choiceElements []Element
}

func NewShedLevel(arg Arg) (*ShedLevel, error) {
	s := &ShedLevel{
		choiceElements: []Element{
			NewElement("percent", V2E(NewUnsigned), WithElementContext(0)),
			NewElement("level", V2E(NewUnsigned), WithElementContext(1)),
			NewElement("amount", V2E(NewReal), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type SetpointReference struct {
	*Sequence
	sequenceElements []Element
}

func NewSetpointReference(arg Arg) (*SetpointReference, error) {
	s := &SetpointReference{
		sequenceElements: []Element{
			NewElement("setpointReference", V2E(NewObjectPropertyReference), WithElementContext(0), WithElementOptional(true)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type SpecialEventPeriod struct {
	*Choice
	choiceElements []Element
}

func NewSpecialEventPeriod(arg Arg) (*SpecialEventPeriod, error) {
	s := &SpecialEventPeriod{
		choiceElements: []Element{
			NewElement("calendarEntry", V2E(NewCalendarEntry), WithElementContext(0)),
			NewElement("calendarReference", Vs2E(NewObjectIdentifier), WithElementContext(1)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type SpecialEvent struct {
	*Sequence
	sequenceElements []Element
}

func NewSpecialEvent(arg Arg) (*SpecialEvent, error) {
	s := &SpecialEvent{
		sequenceElements: []Element{
			NewElement("period", V2E(NewSpecialEventPeriod)),
			NewElement("listOfTimeValues", SequenceOf(NewTimeValue), WithElementContext(2)),
			NewElement("eventPriority", V2E(NewUnsigned), WithElementContext(3)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type StageLimitValue struct {
	*Sequence
	sequenceElements []Element
}

func NewStageLimitValue(arg Arg) (*StageLimitValue, error) {
	s := &StageLimitValue{
		sequenceElements: []Element{
			NewElement("limit", V2E(NewReal)),
			NewElement("values", Vs2E(NewBitString)),
			NewElement("deadband", V2E(NewReal)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type ValueSource struct {
	*Choice
	choiceElements []Element
}

func NewValueSource(arg Arg) (*ValueSource, error) {
	s := &ValueSource{
		choiceElements: []Element{
			NewElement("none", V2E(NewNull), WithElementContext(0)),
			NewElement("object", V2E(NewDeviceObjectReference), WithElementContext(1)),
			NewElement("address", V2E(NewAddress), WithElementContext(2)),
		},
	}
	panic("implementchoice")
	return s, nil
}

type VTSession struct {
	*Sequence
	sequenceElements []Element
}

func NewVTSession(arg Arg) (*VTSession, error) {
	s := &VTSession{
		sequenceElements: []Element{
			NewElement("localVtSessionID", V2E(NewUnsigned)),
			NewElement("remoteVtSessionID", V2E(NewUnsigned)),
			NewElement("remoteVtAddress", V2E(NewDeviceAddress)),
		},
	}
	panic("implementchoice")
	return s, nil
}
