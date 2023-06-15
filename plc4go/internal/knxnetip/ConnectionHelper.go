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

package knxnetip

import (
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"math"
	"net"
	"runtime/debug"
	"strconv"
	"sync/atomic"
	"time"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/pkg/errors"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////
// Internal helper functions
///////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *Connection) interceptIncomingMessage(spi.Message) {
	m.resetTimeout()
	if m.connectionStateTimer != nil {
		// Reset the timer for sending the ConnectionStateRequest
		m.connectionStateTimer.Reset(60 * time.Second)
	}
}

func (m *Connection) castIpToKnxAddress(ip net.IP) driverModel.IPAddress {
	return driverModel.NewIPAddress(ip[len(ip)-4:])
}

func (m *Connection) handleIncomingTunnelingRequest(ctx context.Context, tunnelingRequest driverModel.TunnelingRequest) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			}
		}()
		lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
		if !ok {
			return
		}
		var destinationAddress []byte
		switch lDataInd.GetDataFrame().(type) {
		case driverModel.LDataExtendedExactly:
			dataFrame := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			destinationAddress = dataFrame.GetDestinationAddress()
			switch dataFrame.GetApdu().(type) {
			case driverModel.ApduDataContainerExactly:
				container := dataFrame.GetApdu().(driverModel.ApduDataContainer)
				switch container.GetDataApdu().(type) {
				case driverModel.ApduDataGroupValueWriteExactly:
					groupValueWrite := container.GetDataApdu().(driverModel.ApduDataGroupValueWrite)
					if destinationAddress == nil {
						return
					}
					var payload []byte
					payload = append(payload, byte(groupValueWrite.GetDataFirstByte()))
					payload = append(payload, groupValueWrite.GetData()...)

					m.handleValueCacheUpdate(ctx, destinationAddress, payload)
				default:
					if dataFrame.GetGroupAddress() {
						return
					}
					// If this is an individual address, and it is targeted at us, we need to ack that.
					ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
					targetAddress := ByteArrayToKnxAddress(ctxForModel, dataFrame.GetDestinationAddress())
					if targetAddress == m.ClientKnxAddress {
						m.log.Info().Msg("Acknowleding an unhandled data message.")
						_ = m.sendDeviceAck(ctx, dataFrame.GetSourceAddress(), dataFrame.GetApdu().GetCounter(), func(err error) {})
					}
				}
			case driverModel.ApduControlContainerExactly:
				if dataFrame.GetGroupAddress() {
					return
				}
				// If this is an individual address, and it is targeted at us, we need to ack that.
				ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
				targetAddress := ByteArrayToKnxAddress(ctxForModel, dataFrame.GetDestinationAddress())
				if targetAddress == m.ClientKnxAddress {
					m.log.Info().Msg("Acknowleding an unhandled contol message.")
					_ = m.sendDeviceAck(ctx, dataFrame.GetSourceAddress(), dataFrame.GetApdu().GetCounter(), func(err error) {})
				}
			}
		default:
			m.log.Info().Msg("Unknown unhandled message.")
		}
	}()
}

func (m *Connection) handleValueCacheUpdate(ctx context.Context, destinationAddress []byte, payload []byte) {
	addressData := uint16(destinationAddress[0])<<8 | (uint16(destinationAddress[1]) & 0xFF)

	m.valueCacheMutex.RLock()
	val, ok := m.valueCache[addressData]
	m.valueCacheMutex.RUnlock()
	changed := false
	if !ok || !m.sliceEqual(val, payload) {
		m.valueCacheMutex.Lock()
		m.valueCache[addressData] = payload
		m.valueCacheMutex.Unlock()
		changed = true
	}
	if m.subscribers != nil {
		for _, subscriber := range m.subscribers {
			subscriber.handleValueChange(ctx, destinationAddress, payload, changed)
		}
	}
}

func (m *Connection) handleTimeout() {
	// If this is the first timeout in a sequence, start the timer.
	/*	if m.connectionTimeoutTimer == nil {
		m.connectionTimeoutTimer = time.NewTimer(m.connectionTtl)
		go func() {
			<-m.connectionTimeoutTimer.C
			m.resetConnection()
		}()
	}*/
}

func (m *Connection) resetTimeout() {
	if m.connectionTimeoutTimer != nil {
		if !m.connectionTimeoutTimer.Stop() {
			select {
			case <-m.connectionTimeoutTimer.C:
			default:
			}
		}
		m.connectionTimeoutTimer = nil
	}
}

func (m *Connection) resetConnection() {
	m.log.Warn().Msg("Reset connection")
}

func (m *Connection) getGroupAddressNumLevels() uint8 {
	if val, ok := m.options["group-address-num-levels"]; ok {
		groupAddressNumLevels, err := strconv.ParseUint(val[0], 10, 8)
		if err == nil {
			return uint8(groupAddressNumLevels)
		}
	}
	// TODO: document magic number
	return 3
}

func (m *Connection) addSubscriber(subscriber *Subscriber) {
	for _, sub := range m.subscribers {
		if sub == subscriber {
			m.log.Debug().Msgf("Subscriber %v already added", subscriber)
			return
		}
	}
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *Connection) removeSubscriber(subscriber *Subscriber) {
	for i, sub := range m.subscribers {
		if sub == subscriber {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
		}
	}
}

func (m *Connection) sliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func (m *Connection) getLocalAddress() (*net.UDPAddr, error) {
	transportInstanceExposer, ok := m.messageCodec.(spi.TransportInstanceExposer)
	if !ok {
		return nil, errors.New("used transport, is not a TransportInstanceExposer")
	}

	// Prepare a SearchReq
	udpTransportInstance, ok := transportInstanceExposer.GetTransportInstance().(*udp.TransportInstance)
	if !ok {
		return nil, errors.New("used transport, is not a UdpTransportInstance")
	}

	return udpTransportInstance.LocalAddress, nil
}

func (m *Connection) getNewSequenceCounter() uint8 {
	sequenceCounter := atomic.AddInt32(&m.SequenceCounter, 1)
	if sequenceCounter >= math.MaxUint8 {
		atomic.StoreInt32(&m.SequenceCounter, -1)
		sequenceCounter = -1
	}
	return uint8(sequenceCounter)
}

func (m *Connection) getNextCounter(targetAddress driverModel.KnxAddress) uint8 {
	m.Lock()
	defer m.Unlock()

	connection, ok := m.DeviceConnections[targetAddress]
	if !ok {
		return 0
	}
	counter := connection.counter
	connection.counter++
	if connection.counter >= 16 {
		connection.counter = 0
	}
	return counter
}

func KnxAddressToString(knxAddress driverModel.KnxAddress) string {
	return fmt.Sprintf("%d.%d.%d", knxAddress.GetMainGroup(), knxAddress.GetMiddleGroup(), knxAddress.GetSubGroup())
}
