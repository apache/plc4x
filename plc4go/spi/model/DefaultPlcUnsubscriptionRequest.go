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
	"context"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

var _ apiModel.PlcUnsubscriptionRequestBuilder = &DefaultPlcUnsubscriptionRequestBuilder{}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcUnsubscriptionRequestBuilder
type DefaultPlcUnsubscriptionRequestBuilder struct {
	subscriptionHandles []apiModel.PlcSubscriptionHandle
}

func NewDefaultPlcUnsubscriptionRequestBuilder() *DefaultPlcUnsubscriptionRequestBuilder {
	return &DefaultPlcUnsubscriptionRequestBuilder{}
}

func (d *DefaultPlcUnsubscriptionRequestBuilder) AddHandles(subscriptionHandles ...apiModel.PlcSubscriptionHandle) apiModel.PlcUnsubscriptionRequestBuilder {
	subscriptionHandles = append(subscriptionHandles, subscriptionHandles...)
	return d
}

func (d *DefaultPlcUnsubscriptionRequestBuilder) Build() (apiModel.PlcUnsubscriptionRequest, error) {
	return NewDefaultPlcUnsubscriptionRequest(d.subscriptionHandles), nil
}

var _ apiModel.PlcUnsubscriptionRequest = &DefaultPlcUnsubscriptionRequest{}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcUnsubscriptionRequest
type DefaultPlcUnsubscriptionRequest struct {
	subscriptionHandles []apiModel.PlcSubscriptionHandle
}

func NewDefaultPlcUnsubscriptionRequest(subscriptionHandles []apiModel.PlcSubscriptionHandle) *DefaultPlcUnsubscriptionRequest {
	return &DefaultPlcUnsubscriptionRequest{
		subscriptionHandles: subscriptionHandles,
	}
}

func (d *DefaultPlcUnsubscriptionRequest) Execute() <-chan apiModel.PlcUnsubscriptionRequestResult {
	return d.ExecuteWithContext(context.Background())
}

func (d *DefaultPlcUnsubscriptionRequest) ExecuteWithContext(ctx context.Context) <-chan apiModel.PlcUnsubscriptionRequestResult {
	results := make(chan apiModel.PlcUnsubscriptionRequestResult, 1)
	go func() {
		var collectedErrors []error
		for _, handle := range d.subscriptionHandles {
			select {
			case unsubscribe := <-handle.(*DefaultPlcSubscriptionHandle).plcSubscriber.Unsubscribe(ctx, d):
				if err := unsubscribe.GetErr(); err != nil {
					collectedErrors = append(collectedErrors, err)
					continue
				}
			case <-ctx.Done():
				collectedErrors = append(collectedErrors, ctx.Err())
			}
		}
		var err error
		if len(collectedErrors) > 0 {
			err = utils.MultiError{MainError: errors.New("error unsubscribing from all"), Errors: collectedErrors}
		}
		results <- NewDefaultPlcUnsubscriptionRequestResult(d, NewDefaultPlcUnsubscriptionResponse(d), err)
	}()
	return results
}

func (d *DefaultPlcUnsubscriptionRequest) GetSubscriptionHandles() []apiModel.PlcSubscriptionHandle {
	return d.subscriptionHandles
}

func (d *DefaultPlcUnsubscriptionRequest) IsAPlcMessage() bool {
	return true
}
