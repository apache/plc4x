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

type NetworkSecurityObject struct {
	Object
}

func NewNetworkSecurityObject(options ...Option) (*NetworkSecurityObject, error) {
	n := new(NetworkSecurityObject)
	objectType := "networkSecurity"
	properties := []Property{
		NewWritableProperty("baseDeviceSecurityPolicy", V2P(NewSecurityLevel)),
		NewWritableProperty("networkAccessSecurityPolicies", ArrayOfP(NewNetworkSecurityPolicy, 0, 0)),
		NewWritableProperty("securityTimeWindow", V2P(NewUnsigned)),
		NewWritableProperty("packetReorderTime", V2P(NewUnsigned)),
		NewReadableProperty("distributionKeyRevision", V2P(NewUnsigned)),
		NewReadableProperty("keySets", ArrayOfP(NewSecurityKeySet, 0, 0)),
		NewWritableProperty("lastKeyServer", V2P(NewAddressBinding)),
		NewWritableProperty("securityPDUTimeout", V2P(NewUnsigned)),
		NewReadableProperty("updateKeySetTimeout", V2P(NewUnsigned)),
		NewReadableProperty("supportedSecurityAlgorithms", ListOfP(NewUnsigned)),
		NewWritableProperty("doNotHide", V2P(NewBoolean)),
	}
	var err error
	n.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, n)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return n, nil
}
