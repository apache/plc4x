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
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/stretchr/testify/assert"
)

// TestRenderTest is a lazy test of Default* structs without proper initialization
func TestRenderTest(t *testing.T) {
	suts := []interface {
		fmt.Stringer
		utils.Serializable
	}{
		&DefaultArrayInfo{},
		&DefaultPlcBrowseItem{},
		&DefaultPlcBrowseRequest{},
		&DefaultPlcBrowseRequestBuilder{},
		&DefaultPlcBrowseRequestResult{},
		&DefaultPlcBrowseResponse{},
		&DefaultPlcBrowseResponseItem{},
		&DefaultPlcConsumerRegistration{},
		&DefaultPlcDiscoveryItem{},
		&DefaultPlcReadRequest{DefaultPlcTagRequest: NewDefaultPlcTagRequest(nil, nil)},
		&DefaultPlcReadRequestResult{},
		&DefaultPlcReadResponse{},
		&DefaultPlcSubscriptionEvent{},
		&DefaultPlcSubscriptionEventItem{},
		&DefaultPlcSubscriptionHandle{},
		&DefaultPlcSubscriptionRequest{DefaultPlcTagRequest: NewDefaultPlcTagRequest(nil, nil)},
		&DefaultPlcSubscriptionRequestResult{},
		&DefaultPlcSubscriptionResponse{},
		&DefaultPlcSubscriptionResponseItem{},
		&DefaultPlcTagRequest{},
		&DefaultPlcUnsubscriptionRequest{},
		&DefaultPlcUnsubscriptionRequestResult{},
		&DefaultPlcUnsubscriptionResponse{},
		&DefaultPlcWriteRequest{DefaultPlcTagRequest: NewDefaultPlcTagRequest(nil, nil)},
		&DefaultPlcWriteRequestResult{},
		&DefaultPlcWriteResponse{},
		&ResponseItem{},
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

// TODO: ensure mocks are created in test context...
// TestRenderTestCustom test some custom objects
func TestRenderTestCustom(t *testing.T) {
	tests := []struct {
		name string
		sut  interface {
			fmt.Stringer
			utils.Serializable
		}
		extraCall func(t *testing.T, _sut any)
	}{
		{
			sut: NewDefaultPlcBrowseItem(
				NewMockPlcTag(t),
				"some name",
				"some datatype",
				true,
				true,
				true,
				map[string]apiModel.PlcBrowseItem{
					"tagid1": NewMockPlcBrowseItem(t),
					"tagid2": NewMockPlcBrowseItem(t),
				},
				map[string]apiValues.PlcValue{
					"tagid1": spiValues.PlcNull{},
					"tagid2": spiValues.PlcNull{},
					"tagid3": nil,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcBrowseRequest(
				map[string]apiModel.PlcQuery{
					"tagid1": NewMockPlcQuery(t),
					"tagid2": NewMockPlcQuery(t),
				},
				[]string{
					"tagid1",
					"tagid2",
				},
				nil,
			),
			extraCall: func(t *testing.T, _sut any) {
				sut := _sut.(*DefaultPlcBrowseRequest)
				// TODO: add browser calls
				sut.GetQuery("tagid")
			},
		},
		{
			sut: NewDefaultPlcBrowseResponse(
				NewDefaultPlcBrowseRequest(nil, nil, nil),
				map[string][]apiModel.PlcBrowseItem{
					"tagid": nil,
				},
				map[string]apiModel.PlcResponseCode{
					"tagid": 0,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
			extraCall: func(t *testing.T, _sut any) {
				sut := _sut.(*DefaultPlcBrowseResponse)
				// TODO: add assertions
				sut.GetResponseCode("tagid")
				sut.GetQueryResults("tagid")
			},
		},
		{
			sut: NewDefaultPlcBrowseResponseItem(0, nil),
		},
		{
			// TODO: we need a mock here for improvement
			sut: NewDefaultPlcConsumerRegistration(nil, nil, nil).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcDiscoveryItem(
				"something",
				"something",
				url.URL{},
				map[string][]string{
					"something": {"else"},
				},
				"something",
				nil,
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcReadRequest(
				map[string]apiModel.PlcTag{
					"tagid": nil,
				},
				[]string{
					"tagid",
				},
				nil,
				nil,
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcReadResponse(
				NewDefaultPlcReadRequest(nil, nil, nil, nil),
				map[string]apiModel.PlcResponseCode{
					"tagid": 0,
				},
				map[string]apiValues.PlcValue{
					"tagid": nil,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcSubscriptionEvent(
				nil,
				map[string]apiModel.PlcTag{
					"tagid": nil,
				},
				map[string]SubscriptionType{
					"tagid": 0,
				},
				map[string]time.Duration{
					"tagid": 0,
				},
				map[string]apiModel.PlcResponseCode{
					"tagid": 0,
				},
				map[string]apiValues.PlcValue{
					"tagid": nil,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcSubscriptionEventItem(
				0,
				nil,
				0,
				0,
				nil,
			),
		},
		{
			sut: NewDefaultPlcSubscriptionHandle(
				nil,
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcSubscriptionRequest(
				nil,
				[]string{"tagid"},
				map[string]apiModel.PlcTag{
					"tagid": nil,
				},
				map[string]SubscriptionType{
					"tagid": 0,
				},
				map[string]time.Duration{
					"tagid": 0,
				},
				map[string][]apiModel.PlcSubscriptionEventConsumer{
					"tagd": nil,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcSubscriptionResponse(
				NewDefaultPlcSubscriptionRequest(nil, nil, nil, nil, nil, nil),
				map[string]apiModel.PlcResponseCode{
					"tagid": 0,
				},
				map[string]apiModel.PlcSubscriptionHandle{
					"tagid": nil,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcSubscriptionResponseItem(
				0,
				NewDefaultPlcSubscriptionHandle(nil),
			),
		},
		{
			sut: NewDefaultPlcTagRequest(
				map[string]apiModel.PlcTag{
					"tagid": nil,
				},
				[]string{"tagid"},
			),
		},
		{
			sut: NewDefaultPlcUnsubscriptionRequestResult(
				NewDefaultPlcUnsubscriptionRequest(),
				NewDefaultPlcUnsubscriptionResponse(),
				nil,
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcWriteRequest(
				map[string]apiModel.PlcTag{
					"tagid": nil,
				},
				[]string{"tagid"},
				map[string]apiValues.PlcValue{
					"tageid": nil,
				},
				nil,
				nil,
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
		{
			sut: NewDefaultPlcWriteResponse(
				NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
				map[string]apiModel.PlcResponseCode{
					"tagid": 0,
				},
			).(interface { // TODO: workaround
				fmt.Stringer
				utils.Serializable
			}),
		},
	}
	for _, tt := range tests {
		sut := tt.sut
		testName := fmt.Sprintf("%T", sut)
		if tt.name != "" {
			testName += " - " + tt.name
		}
		testName += fmt.Sprintf("%T", sut)
		t.Run(testName, func(t *testing.T) {
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
							// TODO: if it is 1 arg and a string it is probably the get tags something: maybe we hardcode it here
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
			if tt.extraCall != nil {
				t.Run("extra call", func(t *testing.T) {
					tt.extraCall(t, sut)
				})
			}
		})
	}
}
