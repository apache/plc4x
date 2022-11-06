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

package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonPanickingStrings(t *testing.T) {
	suts := []fmt.Stringer{
		&DefaultPlcBrowseResponse{},
		&DefaultPlcConsumerRegistration{},
		&DefaultPlcDiscoveryItem{},
		&DefaultPlcReadRequest{},
		&DefaultPlcReadRequestResult{},
		&DefaultPlcReadResponse{},
		&DefaultPlcSubscriptionEvent{},
		&DefaultPlcSubscriptionHandle{},
		&DefaultPlcSubscriptionRequest{},
		&DefaultPlcSubscriptionRequestResult{},
		&DefaultPlcSubscriptionResponse{},
		//&DefaultPlcUnsubscriptionRequest{}, //TODO: empty file
		&DefaultPlcUnsubscriptionRequestResult{},
		//&DefaultPlcUnsubscriptionResponse{}, //TODO: empty file
		&DefaultPlcWriteRequest{},
		&DefaultPlcWriteRequestResult{},
		&DefaultPlcWriteResponse{},
		&DefaultPlcRequest{},
		&DefaultResponse{},
		&DefaultPlcBrowseRequestResult{},
		&DefaultPlcBrowseRequest{},
		&DefaultPlcBrowseQueryResult{},
		&DefaultPlcBrowseEvent{},
	}
	for _, sut := range suts {
		t.Run(fmt.Sprintf("%T", sut), func(t *testing.T) {
			assert.NotEmptyf(t, sut.String(), "string should at least return type informations")
		})
	}
}
