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

type ElevatorGroupObject struct {
	Object
}

func NewElevatorGroupObject(options ...Option) (*ElevatorGroupObject, error) {
	e := new(ElevatorGroupObject)
	objectType := "elevatorGroup"
	properties := []Property{
		NewReadableProperty("machineRoomID", Vs2P(NewObjectIdentifier)),
		NewReadableProperty("groupID", V2P(NewUnsigned8)),
		NewReadableProperty("groupMembers", ArrayOfPs(NewObjectIdentifier, 0, 0)),
		NewOptionalProperty("groupMode", V2P(NewLiftGroupMode)),
		NewOptionalProperty("landingCalls", ListOfP(NewLandingCallStatus)),
		NewOptionalProperty("landingCallControl", V2P(NewLandingCallStatus)),
	}
	var err error
	e.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, e)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return e, nil
}
