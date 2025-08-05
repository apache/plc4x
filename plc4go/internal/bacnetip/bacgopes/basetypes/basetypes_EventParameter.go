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
