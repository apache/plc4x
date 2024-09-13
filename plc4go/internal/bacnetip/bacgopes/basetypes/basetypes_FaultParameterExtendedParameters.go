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
