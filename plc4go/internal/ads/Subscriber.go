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

package ads

import (
	"context"
	"github.com/apache/plc4x/plc4go/spi/options"
	"runtime/debug"
	"time"

	dirverModel "github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

func (m *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(m.GetPlcTagHandler(), m.GetPlcValueHandler(), m)
}

func (m *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	// TODO: Implement this ...
	return nil
}

func (m *Connection) Subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	defaultSubscriptionRequest := subscriptionRequest.(*spiModel.DefaultPlcSubscriptionRequest)

	// Subscription requests are unfortunately now supported as multi-item requests,
	// so we have to send one subscription request for every tag.
	// Prepare the subscription requests, by ensuring any symbolic tags are correctly resolved
	// and splitting them up into single-tag requests.
	subSubscriptionRequests := map[string]apiModel.PlcSubscriptionRequest{}
	subResults := map[string]apiModel.PlcSubscriptionRequestResult{}
	// Iterate over all requests and add the result-channels to the list
	for _, tagName := range subscriptionRequest.GetTagNames() {
		tag := subscriptionRequest.GetTag(tagName)
		var directTag dirverModel.DirectPlcTag
		switch tag.(type) {
		case dirverModel.DirectPlcTag:
			directTag = tag.(dirverModel.DirectPlcTag)
		case dirverModel.SymbolicPlcTag:
			symbolicTag := tag.(dirverModel.SymbolicPlcTag)
			directTagPtr, err := m.driverContext.getDirectTagForSymbolTag(symbolicTag)
			if err != nil {
				subResults[tagName] = spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Wrap(err, "error resolving symbolic tag"))
				continue
			}
			directTag = *directTagPtr
		default:
			subResults[tagName] = spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.New("invalid tag type"))
			continue
		}

		subscriptionType := defaultSubscriptionRequest.GetType(tagName)
		interval := defaultSubscriptionRequest.GetInterval(tagName)
		preRegisteredConsumers := defaultSubscriptionRequest.GetPreRegisteredConsumers(tagName)
		subSubscriptionRequests[tagName] = spiModel.NewDefaultPlcSubscriptionRequest(m,
			[]string{tagName},
			map[string]apiModel.PlcTag{tagName: directTag},
			map[string]spiModel.SubscriptionType{tagName: subscriptionType},
			map[string]time.Duration{tagName: interval},
			map[string][]apiModel.PlcSubscriptionEventConsumer{tagName: preRegisteredConsumers})
	}

	// If this is a single item request, we can take a shortcut.
	if len(subSubscriptionRequests) == 1 {
		tagName := subscriptionRequest.GetTagNames()[0]
		subSubscriptionRequest := subSubscriptionRequests[tagName]
		return m.subscribe(ctx, subSubscriptionRequest)
	}

	// Create a sub-result-channel map and prepare execute the subscriptions.
	subResultChannels := map[string]<-chan apiModel.PlcSubscriptionRequestResult{}
	for tagName, subSubscriptionRequest := range subSubscriptionRequests {
		subResultChannels[tagName] = m.subscribe(ctx, subSubscriptionRequest)
	}

	// Create a new result-channel, which completes as soon as all sub-result-channels have returned
	globalResultChannel := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			}
		}()
		// Iterate over all sub-results
		for _, subResultChannel := range subResultChannels {
			select {
			case <-ctx.Done():
				globalResultChannel <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, ctx.Err())
				return
			case subResult := <-subResultChannel:
				// These are all single value requests ... so it's safe to assume this shortcut.
				tagName := subResult.GetRequest().GetTagNames()[0]
				subResults[tagName] = subResult
			}
		}
		// As soon as all are done, process the results
		result := m.processSubscriptionResponses(ctx, subscriptionRequest, subResults)
		// Return the final result
		globalResultChannel <- result
	}()

	return globalResultChannel
}

func (m *Connection) subscribe(ctx context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest) <-chan apiModel.PlcSubscriptionRequestResult {
	responseChan := make(chan apiModel.PlcSubscriptionRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				responseChan <- spiModel.NewDefaultPlcSubscriptionRequestResult(subscriptionRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		// At this point we are sure to only have single item direct tag requests.
		tagName := subscriptionRequest.GetTagNames()[0]
		directTag := subscriptionRequest.GetTag(tagName).(dirverModel.DirectPlcTag)
		if directTag.DataType == nil {
			directTag.DataType = m.driverContext.dataTypeTable[directTag.ValueType.String()]
		}

		response, err := m.ExecuteAdsAddDeviceNotificationRequest(ctx, directTag.IndexGroup, directTag.IndexOffset, directTag.DataType.GetSize(), model.AdsTransMode_ON_CHANGE, 0, 0)
		if err != nil {
			responseChan <- spiModel.NewDefaultPlcSubscriptionRequestResult(
				subscriptionRequest,
				nil,
				err,
			)
		}
		// Create a new subscription handle.
		subscriptionHandle := dirverModel.NewAdsSubscriptionHandle(m, tagName, directTag, options.WithCustomLogger(m.log))
		responseChan <- spiModel.NewDefaultPlcSubscriptionRequestResult(
			subscriptionRequest,
			spiModel.NewDefaultPlcSubscriptionResponse(
				subscriptionRequest,
				map[string]apiModel.PlcResponseCode{tagName: apiModel.PlcResponseCode_OK},
				map[string]apiModel.PlcSubscriptionHandle{tagName: subscriptionHandle},
				options.WithCustomLogger(m.log),
			),
			nil,
		)
		// Store it together with the returned ADS handle.
		m.subscriptions[response.GetNotificationHandle()] = subscriptionHandle
	}()
	return responseChan
}

func (m *Connection) processSubscriptionResponses(_ context.Context, subscriptionRequest apiModel.PlcSubscriptionRequest, subscriptionResults map[string]apiModel.PlcSubscriptionRequestResult) apiModel.PlcSubscriptionRequestResult {
	if len(subscriptionResults) == 1 {
		m.log.Debug().Msg("We got only one response, no merging required")
		for tagName := range subscriptionResults {
			return subscriptionResults[tagName]
		}
	}

	m.log.Trace().Msg("Merging requests")
	responseCodes := map[string]apiModel.PlcResponseCode{}
	subscriptionHandles := map[string]apiModel.PlcSubscriptionHandle{}
	var err error = nil
	for _, subscriptionResult := range subscriptionResults {
		if subscriptionResult.GetErr() != nil {
			m.log.Debug().Err(subscriptionResult.GetErr()).Msgf("Error during subscription")
			if err == nil {
				// Lazy initialization of multi error
				err = utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{subscriptionResult.GetErr()}}
			} else {
				multiError := err.(utils.MultiError)
				multiError.Errors = append(multiError.Errors, subscriptionResult.GetErr())
			}
		} else if subscriptionResult.GetResponse() != nil {
			if len(subscriptionResult.GetResponse().GetRequest().GetTagNames()) > 1 {
				m.log.Error().Int("numberOfTags", len(subscriptionResult.GetResponse().GetRequest().GetTagNames())).Msg("We should only get 1")
			}
			for _, tagName := range subscriptionResult.GetResponse().GetRequest().GetTagNames() {
				handle, err := subscriptionResult.GetResponse().GetSubscriptionHandle(tagName)
				if err != nil {
					responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
				} else {
					responseCodes[tagName] = subscriptionResult.GetResponse().GetResponseCode(tagName)
					subscriptionHandles[tagName] = handle
				}
			}
		}
	}
	return spiModel.NewDefaultPlcSubscriptionRequestResult(
		subscriptionRequest,
		spiModel.NewDefaultPlcSubscriptionResponse(
			subscriptionRequest,
			responseCodes,
			subscriptionHandles,
			options.WithCustomLogger(m.log),
		),
		err,
	)
}

func (m *Connection) handleIncomingDeviceNotificationRequest(deviceNotificationRequest model.AdsDeviceNotificationRequest) {
	for _, stampHeader := range deviceNotificationRequest.GetAdsStampHeaders() {
		for _, sample := range stampHeader.GetAdsNotificationSamples() {
			notificationHandle := sample.GetNotificationHandle()
			subscriptionHandler, ok := m.subscriptions[notificationHandle]
			if !ok {
				continue
			}
			adsSubscriptionHandler, ok := subscriptionHandler.(*dirverModel.AdsSubscriptionHandle)
			if !ok {
				continue
			}

			// If no one is consuming, we don't need to waste effort in parsing.
			if adsSubscriptionHandler.GetNumConsumers() == 0 {
				continue
			}

			// Create a readBuffer containing the sample data
			readBuffer := utils.NewReadBufferByteBased(sample.GetData())

			// Parse the data according to theÏ type data stored in the tag
			directTag := adsSubscriptionHandler.GetDirectTag()
			plcValue, err := m.parsePlcValue(directTag.DataType, directTag.DataType.GetArrayInfo(), readBuffer)
			if err != nil {
				continue
			}

			// Publish the parsed value to the subscribers.
			adsSubscriptionHandler.PublishPlcValue(plcValue)
		}
	}
}

func (m *Connection) Unsubscribe(ctx context.Context, unsubscriptionRequest apiModel.PlcUnsubscriptionRequest) <-chan apiModel.PlcUnsubscriptionRequestResult {
	return nil
}

func (m *Connection) Register(consumer apiModel.PlcSubscriptionEventConsumer, handles []apiModel.PlcSubscriptionHandle) apiModel.PlcConsumerRegistration {
	return nil
}

func (m *Connection) Unregister(registration apiModel.PlcConsumerRegistration) {
}
