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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRenderTest is a lazy test of Default* structs without proper initialization
func TestRenderTest(t *testing.T) {
	suts := []interface {
		fmt.Stringer
		utils.Serializable
	}{
		&DefaultArrayInfo{},
		&DefaultPlcBrowseEvent{},
		&DefaultPlcBrowseItem{},
		&DefaultPlcBrowseRequest{},
		&DefaultPlcBrowseRequestResult{},
		&DefaultPlcBrowseResponse{},
		&DefaultPlcBrowseResponseItem{},
		&DefaultPlcConsumerRegistration{},
		&DefaultPlcDiscoveryItem{},
		&DefaultPlcReadRequest{},
		&DefaultPlcReadRequestResult{},
		&DefaultPlcReadResponse{},
		&DefaultPlcSubscriptionEvent{},
		&DefaultPlcSubscriptionEventItem{},
		&DefaultPlcSubscriptionHandle{},
		&DefaultPlcSubscriptionRequest{},
		&DefaultPlcSubscriptionRequestResult{},
		&DefaultPlcSubscriptionResponse{},
		&DefaultPlcSubscriptionResponseItem{},
		&DefaultPlcTagRequest{},
		//&DefaultPlcUnsubscriptionRequest{}, //TODO: empty file
		&DefaultPlcUnsubscriptionRequestResult{},
		//&DefaultPlcUnsubscriptionResponse{}, //TODO: empty file
		&DefaultPlcWriteRequest{},
		&DefaultPlcWriteRequestResult{},
		&DefaultPlcWriteResponse{},
	}
	for _, sut := range suts {
		t.Run(fmt.Sprintf("%T", sut), func(t *testing.T) {
			t.Run("String", func(t *testing.T) {
				assert.NotEmptyf(t, sut.String(), "string should at least return type informations")
			})
			t.Run("Get*/IsÜ*", func(t *testing.T) {
				valueOf := reflect.ValueOf(sut)
				for i := 0; i < valueOf.NumMethod(); i++ {
					method := valueOf.Method(i)
					methodName := valueOf.Type().Method(i).Name
					if strings.HasPrefix(methodName, "Get") || strings.HasPrefix(methodName, "Is") {
						t.Run(methodName, func(t *testing.T) {
							if na := method.Type().NumIn(); na != 0 {
								t.Skipf("skipping because to many argument: %d", na)
							}
							method.Call(nil)
						})
					}
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				serialize, err := sut.Serialize()
				assert.NoError(t, err)
				_ = serialize
			})
		})
	}
}
