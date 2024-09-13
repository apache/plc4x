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
