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
