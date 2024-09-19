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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

type GroupObject struct {
	Object
	objectType           string // TODO: migrateme
	properties           []Property
	_object_supports_cov bool
}

func NewGroupObject(arg Arg) (*GroupObject, error) {
	o := &GroupObject{
		objectType: "group",
		properties: []Property{
			NewReadableProperty("listOfGroupMembers", ListOfP(NewReadAccessSpecification)),
			NewReadableProperty("presentValue", ListOfP(NewReadAccessResult)),
		},
	}
	// TODO: @register_object_type
	return o, nil
}
