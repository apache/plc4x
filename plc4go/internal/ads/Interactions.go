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
	"fmt"
	"time"

	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

func (m *Connection) ExecuteAdsReadDeviceInfoRequest(ctx context.Context) (model.AdsReadDeviceInfoResponse, error) {
	responseChannel := make(chan model.AdsReadDeviceInfoResponse)
	go func() {
		request := m.NewAdsReadDeviceInfoRequest()
		if err := m.messageCodec.SendRequest(
			ctx,
			request,
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == request.GetUserdata().GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				response := amsTcpPacket.GetUserdata().(model.AdsReadDeviceInfoResponse)
				responseChannel <- response
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second); err != nil {
			close(responseChannel)
		} else {
			close(responseChannel)
		}
	}()
	response, err := ReadWithTimeout(ctx, responseChannel)
	if err != nil {
		return nil, fmt.Errorf("error reading device info: %v", err)
	}
	return response, nil
}

func (m *Connection) ExecuteAdsReadRequest(ctx context.Context, indexGroup uint32, indexOffset uint32, length uint32) (model.AdsReadResponse, error) {
	responseChannel := make(chan model.AdsReadResponse)
	go func() {
		request := m.NewAdsReadRequest(indexGroup, indexOffset, length)
		if err := m.messageCodec.SendRequest(
			ctx,
			request,
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == request.GetUserdata().GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				response := amsTcpPacket.GetUserdata().(model.AdsReadResponse)
				responseChannel <- response
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second*5); err != nil {
			close(responseChannel)
		} else {
			//			close(responseChannel)
		}
	}()
	response, err := ReadWithTimeout(ctx, responseChannel)
	if err != nil {
		return nil, fmt.Errorf("error reading: %v", err)
	}
	return response, nil
}

func (m *Connection) ExecuteAdsWriteRequest(ctx context.Context, indexGroup uint32, indexOffset uint32, data []byte) (model.AdsWriteResponse, error) {
	responseChannel := make(chan model.AdsWriteResponse)
	go func() {
		request := m.NewAdsWriteRequest(indexGroup, indexOffset, data)
		if err := m.messageCodec.SendRequest(
			ctx,
			request,
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == request.GetUserdata().GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				response := amsTcpPacket.GetUserdata().(model.AdsWriteResponse)
				responseChannel <- response
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second); err != nil {
			close(responseChannel)
		} else {
			close(responseChannel)
		}
	}()
	response, err := ReadWithTimeout(ctx, responseChannel)
	if err != nil {
		return nil, fmt.Errorf("error writing: %v", err)
	}
	return response, nil
}

func (m *Connection) ExecuteAdsReadWriteRequest(ctx context.Context, indexGroup uint32, indexOffset uint32, readLength uint32, items []model.AdsMultiRequestItem, writeData []byte) (model.AdsReadWriteResponse, error) {
	responseChannel := make(chan model.AdsReadWriteResponse)
	go func() {
		request := m.NewAdsReadWriteRequest(indexGroup, indexOffset, readLength, items, writeData)
		if err := m.messageCodec.SendRequest(
			ctx,
			request,
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == request.GetUserdata().GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				response := amsTcpPacket.GetUserdata().(model.AdsReadWriteResponse)
				responseChannel <- response
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second); err != nil {
			close(responseChannel)
		} else {
			close(responseChannel)
		}
	}()
	response, err := ReadWithTimeout(ctx, responseChannel)
	if err != nil {
		return nil, fmt.Errorf("error writing: %v", err)
	}
	return response, nil
}

func (m *Connection) ExecuteAdsAddDeviceNotificationRequest(ctx context.Context, indexGroup uint32, indexOffset uint32, length uint32, transmissionMode model.AdsTransMode, maxDelay uint32, cycleTime uint32) (model.AdsAddDeviceNotificationResponse, error) {
	responseChannel := make(chan model.AdsAddDeviceNotificationResponse)
	go func() {
		request := m.NewAdsAddDeviceNotificationRequest(indexGroup, indexOffset, length, transmissionMode, maxDelay, cycleTime)
		if err := m.messageCodec.SendRequest(
			ctx,
			request,
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == request.GetUserdata().GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				response := amsTcpPacket.GetUserdata().(model.AdsAddDeviceNotificationResponse)
				responseChannel <- response
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second); err != nil {
			close(responseChannel)
		} else {
			close(responseChannel)
		}
	}()
	response, err := ReadWithTimeout(ctx, responseChannel)
	if err != nil {
		return nil, fmt.Errorf("error writing: %v", err)
	}
	return response, nil
}

func (m *Connection) ExecuteAdsDeleteDeviceNotificationRequest(ctx context.Context, notificationHandle uint32) (model.AdsDeleteDeviceNotificationResponse, error) {
	responseChannel := make(chan model.AdsDeleteDeviceNotificationResponse)
	go func() {
		request := m.NewAdsDeleteDeviceNotificationRequest(notificationHandle)
		if err := m.messageCodec.SendRequest(
			ctx,
			request,
			func(message spi.Message) bool {
				amsTcpPacket, ok := message.(model.AmsTCPPacket)
				if !ok {
					return false
				}
				return amsTcpPacket.GetUserdata().GetInvokeId() == request.GetUserdata().GetInvokeId()
			},
			func(message spi.Message) error {
				amsTcpPacket := message.(model.AmsTCPPacket)
				response := amsTcpPacket.GetUserdata().(model.AdsDeleteDeviceNotificationResponse)
				responseChannel <- response
				return nil
			},
			func(err error) error {
				return nil
			},
			time.Second); err != nil {
			close(responseChannel)
		} else {
			close(responseChannel)
		}
	}()
	response, err := ReadWithTimeout(ctx, responseChannel)
	if err != nil {
		return nil, fmt.Errorf("error writing: %v", err)
	}
	return response, nil
}

func ReadWithTimeout[T spi.Message](ctx context.Context, ch <-chan T) (T, error) {
	timeout, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	select {
	case m := <-ch:
		return m, nil
	case <-timeout.Done():
		var t T
		return t, fmt.Errorf("timeout")
	}
}
