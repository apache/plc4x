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

package service

import "github.com/apache/plc4x/plc4go/internal/bacnetip"

type SubscriptionList struct {
	//TODO: implement me
}

type Subscription struct {
	*bacnetip.OneShotTask
	*bacnetip.DebugContents
	//TODO: implement me
}

type COVDetection struct {
	*DetectionAlgorithm
	//TODO: implement me
}

type GenericCriteria struct {
	*COVDetection
	//TODO: implement me
}

type COVIncrementCriteria struct {
	*COVDetection
	//TODO: implement me
}

type AccessDoorCriteria struct {
	*COVDetection
	//TODO: implement me
}

type AccessPointCriteria struct {
	*COVDetection
	//TODO: implement me
}

type CredentialDataInputCriteria struct {
	*COVDetection
	//TODO: implement me
}

type LoadControlCriteria struct {
	*COVDetection
	//TODO: implement me
}

type PulseConverterCriteria struct {
	*COVIncrementCriteria
	//TODO: implement me
}

var criteriaTypeMap any //TODO: implement me

type ActiveCOVSubscription struct {
	*bacnetip.Property
	//TODO: implement me
}

type ChangeOfValuesServices struct {
	*bacnetip.Capability
	//TODO: implement me
}
