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

type LiftObject struct {
	Object
}

func NewLiftObject(options ...Option) (*LiftObject, error) {
	l := new(LiftObject)
	objectType := "lift"
	properties := []Property{
		NewReadableProperty("trackingValue", V2P(NewReal)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("elevatorGroup", Vs2P(NewObjectIdentifier)),
		NewReadableProperty("groupID", V2P(NewUnsigned8)),
		NewReadableProperty("installationID", V2P(NewUnsigned8)),
		NewOptionalProperty("floorText", ArrayOfP(NewCharacterString, 0, 0)),
		NewOptionalProperty("carDoorText", ArrayOfP(NewCharacterString, 0, 0)),
		NewOptionalProperty("assignedLandingCalls", ArrayOfP(NewAssignedLandingCalls, 0, 0)),
		NewOptionalProperty("makingCarCall", ArrayOfP(NewUnsigned8, 0, 0)),
		NewOptionalProperty("registeredCarCall", ArrayOfP(NewLiftCarCallList, 0, 0)),
		NewOptionalProperty("carPosition", V2P(NewUnsigned8)),
		NewOptionalProperty("carMovingDirection", V2P(NewLiftCarDirection)),
		NewOptionalProperty("carAssignedDirection", V2P(NewLiftCarDirection)),
		NewOptionalProperty("carDoorStatus", ArrayOfP(NewDoorStatus, 0, 0)),
		NewOptionalProperty("carDoorCommand", ArrayOfP(NewLiftCarDoorCommand, 0, 0)),
		NewOptionalProperty("carDoorZone", V2P(NewBoolean)),
		NewOptionalProperty("carMode", V2P(NewLiftCarMode)),
		NewOptionalProperty("carLoad", V2P(NewReal)),
		NewOptionalProperty("carLoadUnits", V2P(NewEngineeringUnits)),
		NewOptionalProperty("nextStoppingFloor", V2P(NewUnsigned)),
		NewReadableProperty("passengerAlarm", V2P(NewBoolean)),
		NewOptionalProperty("timeDelay", V2P(NewUnsigned)),
		NewOptionalProperty("timeDelayNormal", V2P(NewUnsigned)),
		NewOptionalProperty("energyMeter", V2P(NewReal)),
		NewOptionalProperty("energyMeterRef", V2P(NewDeviceObjectReference)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewOptionalProperty("outOfService", V2P(NewBoolean)),
		NewOptionalProperty("carDriveStatus", V2P(NewLiftCarDriveStatus)),
		NewOptionalProperty("faultSignals", ListOfP(NewLiftFault)),
		NewOptionalProperty("landingDoorStatus", ArrayOfP(NewLandingDoorStatus, 0, 0)),
		NewOptionalProperty("higherDeck", Vs2P(NewObjectIdentifier)),
		NewOptionalProperty("lowerDeck", Vs2P(NewObjectIdentifier)),
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventAlgorithmInhibitRef", V2P(NewObjectPropertyReference)),
		NewOptionalProperty("eventAlgorithmInhibit", V2P(NewBoolean)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
	}
	var err error
	l.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, l)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return l, nil
}
