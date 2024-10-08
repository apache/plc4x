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

type AccessCredentialObject struct {
	Object
}

func NewAccessCredentialObject(options ...Option) (*AccessCredentialObject, error) {
	a := new(AccessCredentialObject)
	objectType := "accessCredential"
	properties := []Property{
		NewWritableProperty("globalIdentifier", V2P(NewUnsigned)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("reliability", V2P(NewReliability)),
		NewReadableProperty("credentialStatus", V2P(NewBinaryPV)),
		NewReadableProperty("reasonForDisable", ListOfP(NewAccessCredentialDisableReason)),
		NewReadableProperty("authenticationFactors", ArrayOfP(NewCredentialAuthenticationFactor, 0, 0)),
		NewReadableProperty("activationTime", V2P(NewDateTime)),
		NewReadableProperty("expirationTime", V2P(NewDateTime)),
		NewReadableProperty("credentialDisable", V2P(NewAccessCredentialDisable)),
		NewOptionalProperty("daysRemaining", V2P(NewInteger)),
		NewOptionalProperty("usesRemaining", V2P(NewInteger)),
		NewOptionalProperty("absenteeLimit", V2P(NewUnsigned)),
		NewOptionalProperty("belongsTo", V2P(NewDeviceObjectReference)),
		NewReadableProperty("assignedAccessRights", ArrayOfP(NewAssignedAccessRights, 0, 0)),
		NewOptionalProperty("lastAccessPoint", V2P(NewDeviceObjectReference)),
		NewOptionalProperty("lastAccessEvent", V2P(NewAccessEvent)),
		NewOptionalProperty("lastUseTime", V2P(NewDateTime)),
		NewOptionalProperty("traceFlag", V2P(NewBoolean)),
		NewOptionalProperty("threatAuthority", V2P(NewAccessThreatLevel)),
		NewOptionalProperty("extendedTimeEnable", V2P(NewBoolean)),
		NewOptionalProperty("authorizationExemptions", ListOfP(NewAuthorizationException)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		//       , NewOptionalProperty("masterExemption", V2P(NewBoolean)
		//       , NewOptionalProperty("passbackExemption", V2P(NewBoolean)
		//       , NewOptionalProperty("occupancyExemption", V2P(NewBoolean)
	}
	var err error
	a.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, a)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return a, nil
}
