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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

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
