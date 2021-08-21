/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"math"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////
// Internal helper functions
///////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *Connection) interceptIncomingMessage(interface{}) {
	m.resetTimeout()
	if m.connectionStateTimer != nil {
		// Reset the timer for sending the ConnectionStateRequest
		m.connectionStateTimer.Reset(60 * time.Second)
	}
}

func (m *Connection) castIpToKnxAddress(ip net.IP) *driverModel.IPAddress {
	return driverModel.NewIPAddress(utils.ByteArrayToInt8Array(ip)[len(ip)-4:])
}

func (m *Connection) handleIncomingTunnelingRequest(tunnelingRequest *driverModel.TunnelingRequest) {
	go func() {
		lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi.Child)
		if lDataInd == nil {
			return
		}
		var destinationAddress []int8
		switch lDataInd.DataFrame.Child.(type) {
		case *driverModel.LDataExtended:
			dataFrame := driverModel.CastLDataExtended(lDataInd.DataFrame)
			destinationAddress = dataFrame.DestinationAddress
			switch dataFrame.Apdu.Child.(type) {
			case *driverModel.ApduDataContainer:
				container := driverModel.CastApduDataContainer(dataFrame.Apdu)
				switch container.DataApdu.Child.(type) {
				case *driverModel.ApduDataGroupValueWrite:
					groupValueWrite := driverModel.CastApduDataGroupValueWrite(container.DataApdu)
					if destinationAddress == nil {
						return
					}
					var payload []int8
					payload = append(payload, groupValueWrite.DataFirstByte)
					payload = append(payload, groupValueWrite.Data...)

					m.handleValueCacheUpdate(destinationAddress, payload)
				default:
					if dataFrame.GroupAddress {
						return
					}
					// If this is an individual address and it is targeted at us, we need to ack that.
					targetAddress := Int8ArrayToKnxAddress(dataFrame.DestinationAddress)
					if *targetAddress == *m.ClientKnxAddress {
						log.Info().Msg("Acknowleding an unhandled data message.")
						_ = m.sendDeviceAck(*dataFrame.SourceAddress, dataFrame.Apdu.Counter, func(err error) {})
					}
				}
			case *driverModel.ApduControlContainer:
				if dataFrame.GroupAddress {
					return
				}
				// If this is an individual address and it is targeted at us, we need to ack that.
				targetAddress := Int8ArrayToKnxAddress(dataFrame.DestinationAddress)
				if *targetAddress == *m.ClientKnxAddress {
					log.Info().Msg("Acknowleding an unhandled contol message.")
					_ = m.sendDeviceAck(*dataFrame.SourceAddress, dataFrame.Apdu.Counter, func(err error) {})
				}
			}
		default:
			log.Info().Msg("Unknown unhandled message.")
		}
	}()
}

func (m *Connection) handleValueCacheUpdate(destinationAddress []int8, payload []int8) {
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
			subscriber.handleValueChange(destinationAddress, payload, changed)
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
	log.Warn().Msg("Reset connection")
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
			log.Debug().Msgf("Subscriber %v already added", subscriber)
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

// TODO: we can replace this with reflect.DeepEqual()
func (m *Connection) sliceEqual(a, b []int8) bool {
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

func KnxAddressToString(knxAddress *driverModel.KnxAddress) string {
	return fmt.Sprintf("%d.%d.%d", knxAddress.MainGroup, knxAddress.MiddleGroup, knxAddress.SubGroup)
}
