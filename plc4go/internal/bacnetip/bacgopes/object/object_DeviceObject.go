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

package object

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type DeviceObject struct {
	Object
}

func NewDeviceObject(options ...Option) (*DeviceObject, error) {
	o := new(DeviceObject)
	objectType := "device"
	properties := []Property{
		NewReadableProperty("systemStatus", V2P(NewDeviceStatus)),
		NewReadableProperty("vendorName", V2P(NewCharacterString)),
		NewReadableProperty("vendorIdentifier", V2P(NewUnsigned)),
		NewReadableProperty("modelName", V2P(NewCharacterString)),
		NewReadableProperty("firmwareRevision", V2P(NewCharacterString)),
		NewReadableProperty("applicationSoftwareVersion", V2P(NewCharacterString)),
		NewOptionalProperty("location", V2P(NewCharacterString)),
		NewReadableProperty("protocolVersion", V2P(NewUnsigned)),
		NewReadableProperty("protocolRevision", V2P(NewUnsigned)),
		NewReadableProperty("protocolServicesSupported", V2P(NewServicesSupported)),
		NewReadableProperty("protocolObjectTypesSupported", V2P(NewObjectTypesSupported)),
		NewReadableProperty("objectList", ArrayOfPs(NewObjectIdentifier, 0, 0)),
		NewOptionalProperty("structuredObjectList", ArrayOfPs(NewObjectIdentifier, 0, 0)),
		NewReadableProperty("maxApduLengthAccepted", V2P(NewUnsigned)),
		NewReadableProperty("segmentationSupported", V2P(NewSegmentation)),
		NewOptionalProperty("vtClassesSupported", ListOfP(NewVTClass)),
		NewOptionalProperty("activeVtSessions", ListOfP(NewVTSession)),
		NewOptionalProperty("localTime", Vs2P(NewTime)),
		NewOptionalProperty("localDate", Vs2P(NewDate)),
		NewOptionalProperty("utcOffset", V2P(NewInteger)),
		NewOptionalProperty("daylightSavingsStatus", V2P(NewBoolean)),
		NewOptionalProperty("apduSegmentTimeout", V2P(NewUnsigned)),
		NewReadableProperty("apduTimeout", V2P(NewUnsigned)),
		NewReadableProperty("numberOfApduRetries", V2P(NewUnsigned)),
		NewOptionalProperty("timeSynchronizationRecipients", ListOfP(NewRecipient)),
		NewOptionalProperty("maxMaster", V2P(NewUnsigned)),
		NewOptionalProperty("maxInfoFrames", V2P(NewUnsigned)),
		NewReadableProperty("deviceAddressBinding", ListOfP(NewAddressBinding)),
		NewReadableProperty("databaseRevision", V2P(NewUnsigned)),
		NewOptionalProperty("configurationFiles", ArrayOfPs(NewObjectIdentifier, 0, 0)),
		NewOptionalProperty("lastRestoreTime", V2P(NewTimeStamp)),
		NewOptionalProperty("backupFailureTimeout", V2P(NewUnsigned)),
		NewOptionalProperty("backupPreparationTime", V2P(NewUnsigned)),
		NewOptionalProperty("restorePreparationTime", V2P(NewUnsigned)),
		NewOptionalProperty("restoreCompletionTime", V2P(NewUnsigned)),
		NewOptionalProperty("backupAndRestoreState", V2P(NewBackupState)),
		NewOptionalProperty("activeCovSubscriptions", ListOfP(NewCOVSubscription)),
		NewOptionalProperty("maxSegmentsAccepted", V2P(NewUnsigned)),
		NewOptionalProperty("slaveProxyEnable", ArrayOfP(NewBoolean, 0, 0)),
		NewOptionalProperty("autoSlaveDiscovery", ArrayOfP(NewBoolean, 0, 0)),
		NewOptionalProperty("slaveAddressBinding", ListOfP(NewAddressBinding)),
		NewOptionalProperty("manualSlaveAddressBinding", ListOfP(NewAddressBinding)),
		NewOptionalProperty("lastRestartReason", V2P(NewRestartReason)),
		NewOptionalProperty("timeOfDeviceRestart", V2P(NewTimeStamp)),
		NewOptionalProperty("restartNotificationRecipients", ListOfP(NewRecipient)),
		NewOptionalProperty("utcTimeSynchronizationRecipients", ListOfP(NewRecipient)),
		NewOptionalProperty("timeSynchronizationInterval", V2P(NewUnsigned)),
		NewOptionalProperty("alignIntervals", V2P(NewBoolean)),
		NewOptionalProperty("intervalOffset", V2P(NewUnsigned)),
		NewOptionalProperty("serialNumber", V2P(NewCharacterString)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewOptionalProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		NewOptionalProperty("activeCovMultipleSubscriptions", ListOfP(NewCOVMultipleSubscription)),
		NewOptionalProperty("auditNotificationRecipient", V2P(NewRecipient)),
		NewOptionalProperty("deviceUUID", V2P(NewOctetString)), // size 16,
		NewOptionalProperty("deployedProfileLocation", V2P(NewCharacterString)),
	}
	var err error
	o.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, o)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return o, nil
}
