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

package opcua

import (
	"context"
	"encoding/binary"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type SubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle
	plcSubscriber       *Subscriber
	connection          *Connection
	subscriptionRequest apiModel.PlcSubscriptionRequest
	tagNames            []string
	consumers           []apiModel.PlcSubscriptionEventConsumer
	subscriptionId      uint32
	cycleTime           time.Duration
	revisedCycleTime    time.Duration
	clientHandles       atomic.Uint32

	destroy      atomic.Bool
	subscriberWg sync.WaitGroup
	complete     bool

	log zerolog.Logger
}

func NewSubscriptionHandle(log zerolog.Logger, subscriber *Subscriber, connection *Connection, subscriptionRequest apiModel.PlcSubscriptionRequest, subscriptionId uint32, cycleTime time.Duration) *SubscriptionHandle {
	s := &SubscriptionHandle{
		plcSubscriber:       subscriber,
		connection:          connection,
		subscriptionRequest: subscriptionRequest,
		tagNames:            subscriptionRequest.GetTagNames(),
		subscriptionId:      subscriptionId,
		cycleTime:           cycleTime,
		log:                 log,
	}
	s.clientHandles.Store(1)
	s.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(subscriber, s)
	_, err := s.onSubscribeCreateMonitoredItemsRequest()
	if err != nil {
		subscriber.onDisconnect()
	}
	s.startSubscriber()
	return s
}

func (h *SubscriptionHandle) onSubscribeCreateMonitoredItemsRequest() (readWriteModel.CreateMonitoredItemsResponse, error) {
	requestList := make([]readWriteModel.MonitoredItemCreateRequest, len(h.tagNames))

	for _, tagName := range h.tagNames {
		tagDefaultPlcSubscription := h.subscriptionRequest.GetTag(tagName)

		idNode, err := generateNodeId(tagDefaultPlcSubscription.(Tag))
		if err != nil {
			return nil, errors.Wrap(err, "error genertating node id")
		}

		readValueId := readWriteModel.NewReadValueId(
			idNode,
			0xD,
			NULL_STRING,
			readWriteModel.NewQualifiedName(0, NULL_STRING))

		var monitoringMode readWriteModel.MonitoringMode
		switch tagDefaultPlcSubscription.GetPlcSubscriptionType() {
		case apiModel.SubscriptionCyclic:
			monitoringMode = readWriteModel.MonitoringMode_monitoringModeSampling
		case apiModel.SubscriptionChangeOfState:
			monitoringMode = readWriteModel.MonitoringMode_monitoringModeReporting
		case apiModel.SubscriptionEvent:
			monitoringMode = readWriteModel.MonitoringMode_monitoringModeReporting
		default:
			monitoringMode = readWriteModel.MonitoringMode_monitoringModeReporting
		}

		clientHandle := h.clientHandles.Add(1) - 1

		parameters := readWriteModel.NewMonitoringParameters(
			clientHandle,
			float64(h.cycleTime),  // sampling interval
			NULL_EXTENSION_OBJECT, // filter, null means use default
			1,                     // queue size
			true,                  // discard oldest
		)

		request := readWriteModel.NewMonitoredItemCreateRequest(
			readValueId, monitoringMode, parameters)

		requestList = append(requestList, request)
	}

	requestHeader := readWriteModel.NewRequestHeader(h.connection.channel.getAuthenticationToken(),
		h.connection.channel.getCurrentDateTime(),
		h.connection.channel.getRequestHandle(),
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT)

	createMonitoredItemsRequest := readWriteModel.NewCreateMonitoredItemsRequest(
		requestHeader,
		h.subscriptionId,
		readWriteModel.TimestampsToReturn_timestampsToReturnBoth,
		requestList,
	)

	identifier := createMonitoredItemsRequest.GetExtensionId()
	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewRootExtensionObject(
		expandedNodeId,
		createMonitoredItemsRequest,
		identifier,
	)

	ctx, cancel := context.WithTimeout(context.Background(), REQUEST_TIMEOUT)
	defer cancel()
	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		return nil, errors.Wrapf(err, "Unable to serialise the ReadRequest")
	}

	responseChan := make(chan readWriteModel.CreateMonitoredItemsResponse, 100) // TODO: bit oversized to not block anything. Discards errors
	errorChan := make(chan error, 100)                                          // TODO: bit oversized to not block anything. Discards errors
	consumer := func(opcuaResponse []byte) {
		unknownExtensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			errorChan <- errors.Wrapf(err, "Unable to read the reply")
			return
		}
		var responseMessage readWriteModel.CreateMonitoredItemsResponse
		switch unknownExtensionObject := unknownExtensionObject.(type) {
		case readWriteModel.CreateMonitoredItemsResponse:
			responseMessage = unknownExtensionObject
		case readWriteModel.ServiceFault:
			serviceFault := unknownExtensionObject
			header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
			errorChan <- errors.Errorf("Subscription ServiceFault returned from server with error code,  '%s'", header.GetServiceResult())
			h.plcSubscriber.onDisconnect()
			return
		default:
			errorChan <- errors.Errorf("Unexpected type %T received", unknownExtensionObject)
			h.plcSubscriber.onDisconnect()
			return
		}

		array := make([]readWriteModel.MonitoredItemCreateResult, len(responseMessage.GetResults()))
		for i, definition := range responseMessage.GetResults() {
			array[i] = definition.(readWriteModel.MonitoredItemCreateResult)
		}
		for index, arrayLength := 0, len(array); index < arrayLength; index++ {
			result := array[index]
			if code, ok := readWriteModel.OpcuaStatusCodeByValue(result.GetStatusCode().GetStatusCode()); !ok || code != readWriteModel.OpcuaStatusCode_Good {
				h.log.Error().Str("tag", h.tagNames[index]).Msg("Invalid Tag, subscription created without this tag")
			} else {
				h.log.Debug().Str("tag", h.tagNames[index]).Msg("Tag was added to the subscription")
			}
		}
		responseChan <- responseMessage
	}

	errorDispatcher := func(err error) {
		errorChan <- errors.Wrap(err, "error received")
	}

	h.connection.channel.submit(ctx, h.connection.messageCodec, errorDispatcher, consumer, buffer)

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errorChan:
		return nil, errors.Wrap(err, "error received")
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), "context ended")
	}
}

/**
 * startSubscriber Main subscriber loop. For subscription, we still need to send a request the server on every cycle.
 * Which includes a request for an update of the previously agreed upon list of tags.
 * The server will respond at most once every cycle.
 */
func (h *SubscriptionHandle) startSubscriber() {
	h.log.Trace().Msg("Starting Subscription")

	h.subscriberWg.Add(1)
	go func() {
		defer h.subscriberWg.Done()

		var outstandingAcknowledgements []readWriteModel.SubscriptionAcknowledgement
		var outstandingRequests []uint32
		for !h.destroy.Load() {

			requestHandle := h.connection.channel.getRequestHandle()

			//If we are waiting on a response and haven't received one, just wait until we do. A keep alive will be sent out eventually
			if len(outstandingRequests) <= 1 {
				requestHeader := readWriteModel.NewRequestHeader(h.connection.channel.getAuthenticationToken(),
					h.connection.channel.getCurrentDateTime(),
					requestHandle,
					0,
					NULL_STRING,
					uint32(h.revisedCycleTime*10),
					NULL_EXTENSION_OBJECT)

				//Make a copy of the outstanding requests, so it isn't modified while we are putting the ack list together.
				acks := slices.Clone(outstandingAcknowledgements)
				ackLength := len(acks)
				if ackLength == 0 {
					ackLength = -1
				}
				{ // golang version of remove all
					tmpOutstandingAcknowledgements := map[readWriteModel.SubscriptionAcknowledgement]bool{}
					for _, acknowledgement := range outstandingAcknowledgements {
						tmpOutstandingAcknowledgements[acknowledgement] = true
					}
					for _, ack := range acks {
						delete(tmpOutstandingAcknowledgements, ack)
					}
					outstandingAcknowledgements = make([]readWriteModel.SubscriptionAcknowledgement, len(tmpOutstandingAcknowledgements))
					count := 0
					for ack := range tmpOutstandingAcknowledgements {
						outstandingAcknowledgements[count] = ack
						count++
					}
				}

				publishRequest := readWriteModel.NewPublishRequest(
					requestHeader,
					acks,
				)

				identifier := publishRequest.GetExtensionId()
				extExpandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
					false, //Server Index Specified
					readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
					nil,
					nil)

				extObject := readWriteModel.NewRootExtensionObject(
					extExpandedNodeId,
					publishRequest,
					identifier,
				)

				ctx := context.Background()

				buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
				if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
					h.log.Error().Err(err).Msg("Unable to serialise the ReadRequest")
					continue
				}

				consumer := func(opcuaResponse []byte) {
					var responseMessage readWriteModel.PublishResponse
					var serviceFault readWriteModel.ServiceFault
					unknownExtensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
					if err != nil {
						h.log.Error().Err(err).Msg("Unable to parse the returned Subscription response")
						h.plcSubscriber.onDisconnect()
						return
					}
					switch unknownExtensionObject := unknownExtensionObject.(type) {
					case readWriteModel.PublishResponse:
						responseMessage = unknownExtensionObject
					case readWriteModel.ServiceFault:
						serviceFault = unknownExtensionObject
						header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
						h.log.Debug().
							Stringer("serviceResult", header.GetServiceResult()).
							Msg("Subscription ServiceFault returned from server with error code, ignoring as it is probably just a result of a Delete Subscription Request")
						//h.plcSubscriber.onDisconnect(context);
						return
					}
					if serviceFault == nil {
						handle := responseMessage.GetResponseHeader().(readWriteModel.ResponseHeader).GetRequestHandle()
						index := 0
						for i, request := range outstandingRequests {
							if request == handle {
								index = i
							}
						}
						outstandingRequests = append(outstandingRequests[:index], outstandingRequests[index+1:]...)

						for _, availableSequenceNumber := range responseMessage.GetAvailableSequenceNumbers() {
							outstandingAcknowledgements = append(outstandingAcknowledgements, readWriteModel.NewSubscriptionAcknowledgement(h.subscriptionId, availableSequenceNumber))
						}

						for _, notificationMessage := range responseMessage.GetNotificationMessage().(readWriteModel.NotificationMessage).GetNotificationData() {
							notification := notificationMessage.GetBody()
							if notification, ok := notification.(readWriteModel.DataChangeNotification); ok {
								h.log.Trace().Msg("Found a Data Change notification")
								items := notification.GetMonitoredItems()
								monitoredNoticiations := make([]readWriteModel.MonitoredItemNotification, len(items))
								for i, item := range items {
									monitoredNoticiations[i] = item.(readWriteModel.MonitoredItemNotification)
								}
								h.onSubscriptionValue(monitoredNoticiations)
							} else {
								h.log.Warn().Msg("Unsupported Notification type")
							}
						}
					}
				}

				errorDispatcher := func(err error) {
					h.log.Error().Err(err).Msg("error received")
					h.plcSubscriber.onDisconnect()
				}

				h.connection.channel.submit(ctx, h.connection.messageCodec, errorDispatcher, consumer, buffer)
			}
			//Put the subscriber loop to sleep for the rest of the cycle.
			time.Sleep(h.revisedCycleTime)
		}
		//Wait for any outstanding responses to arrive, using the request timeout length
		//sleep(this.revisedCycleTime * 10);
		h.complete = true
	}()
}

// stopSubscriber stops the subscriber either on disconnect or on error
func (h *SubscriptionHandle) stopSubscriber() {
	h.destroy.Store(true)

	requestHandle := h.connection.channel.getRequestHandle()

	requestHeader := readWriteModel.NewRequestHeader(
		h.connection.channel.getAuthenticationToken(),
		h.connection.channel.getCurrentDateTime(),
		requestHandle,
		0,
		NULL_STRING,
		uint32(h.revisedCycleTime*10),
		NULL_EXTENSION_OBJECT,
	)

	subscriptions := []uint32{h.subscriptionId}
	deleteSubscriptionrequest := readWriteModel.NewDeleteSubscriptionsRequest(requestHeader,
		subscriptions,
	)

	identifier := deleteSubscriptionrequest.GetExtensionId()
	extExpandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil,
	)

	extObject := readWriteModel.NewRootExtensionObject(
		extExpandedNodeId,
		deleteSubscriptionrequest,
		identifier,
	)

	ctx := context.Background()

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		h.log.Error().Err(err).Msg("Unable to serialise the ReadRequest")
		return
	}

	consumer := func(opcuaResponse []byte) {
		var responseMessage readWriteModel.DeleteSubscriptionsResponse
		unknownExtensionObject, err := readWriteModel.ExtensionObjectParseWithBuffer[readWriteModel.ExtensionObject](ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			h.log.Error().Err(err).Msg("Unable to parse the returned Subscription response")
			h.plcSubscriber.onDisconnect()
			return
		}
		switch unknownExtensionObject := unknownExtensionObject.(type) {
		case readWriteModel.DeleteSubscriptionsResponse:
			responseMessage = unknownExtensionObject
		case readWriteModel.ServiceFault:
			serviceFault := unknownExtensionObject
			header := serviceFault.GetResponseHeader().(readWriteModel.ResponseHeader)
			h.log.Debug().
				Stringer("serviceResult", header.GetServiceResult()).
				Msg("Subscription ServiceFault returned from server with error code, ignoring as it is probably just a result of a Delete Subscription Request")
			return
		}
		h.log.Debug().Stringer("responseMessage", responseMessage).Msg("Received response")
	}

	errorDispatcher := func(err error) {
		h.log.Error().Err(err).Msg("error received")
		h.plcSubscriber.onDisconnect()
	}

	h.connection.channel.submit(ctx, h.connection.messageCodec, errorDispatcher, consumer, buffer)
}

/**
 * onSubscriptionValue Receive the returned values from the OPCUA server and format it so that it can be received by the PLC4X client.
 *
 * @param values - array of data values to be sent to the client.
 */
func (h *SubscriptionHandle) onSubscriptionValue(values []readWriteModel.MonitoredItemNotification) {
	var tagNameList []string
	var dataValues []readWriteModel.DataValue
	for _, value := range values {
		tagNameList = append(tagNameList, h.tagNames[value.GetClientHandle()-1])
		dataValues = append(dataValues, value.GetValue())
	}
	_, responseCodes, responseValues := readResponse(h.log, nil, tagNameList, dataValues)
	// TODO: big oof. Where do we get all those nil values from?
	event := spiModel.NewDefaultPlcSubscriptionEvent(nil, nil, nil, nil, responseCodes, responseValues)

	for _, consumer := range h.consumers {
		consumer(event)
	}
}
