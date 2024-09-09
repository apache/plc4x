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
)

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
