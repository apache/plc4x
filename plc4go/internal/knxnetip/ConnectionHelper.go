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
	"math"
	"net"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////
// Internal helper functions
///////////////////////////////////////////////////////////////////////////////////////////////////////

// interceptIncomingMessage is called by the codec for every frame which arrives.
//
// The keepalive ticker is deliberately NOT reset here: the KNXnet/IP heartbeat has to
// run at a fixed interval, independent of the traffic on the bus. Resetting it on every
// incoming frame means that on any installation with group telegrams more often than
// every 60s the ConnectionStateRequest would never be sent at all and the gateway would
// drop the tunnel after 120s.
// (Java: KnxNetIpConnection uses scheduleAtFixedRate with a fixed heartbeat interval)
func (m *Connection) interceptIncomingMessage(spi.Message) {
	m.resetTimeout()
}

func (m *Connection) castIpToKnxAddress(ip net.IP) driverModel.IPAddress {
	return driverModel.NewIPAddress(ip[len(ip)-4:])
}

func (m *Connection) handleIncomingTunnelingRequest(ctx context.Context, tunnelingRequest driverModel.TunnelingRequest) {
	m.wg.Go(func() {
		defer func() {
			if err := recover(); err != nil {
				m.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataInd)
		if !ok {
			return
		}
		var destinationAddress []byte
		switch lDataInd.GetDataFrame().(type) {
		case driverModel.LDataExtended:
			dataFrame := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			destinationAddress = dataFrame.GetDestinationAddress()
			switch dataFrame.GetApdu().(type) {
			case driverModel.ApduDataContainer:
				container := dataFrame.GetApdu().(driverModel.ApduDataContainer)
				switch container.GetDataApdu().(type) {
				case driverModel.ApduDataGroupValueWrite:
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
			case driverModel.ApduControlContainer:
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
	})
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
	for _, subscriber := range m.getSubscribers() {
		subscriber.handleValueChange(ctx, destinationAddress, payload, changed)
	}
}

func (m *Connection) handleTimeout() {
	// If this is the first timeout in a sequence, start the timer.
	/*	if m.connectionTimeoutTimer == nil {
		m.connectionTimeoutTimer = time.NewTimer(m.connectionTtl)
			m.wg.Go(func() {
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

// getTunnelConnectionType evaluates the "connection-type" connection option and
// maps it to the KnxLayer used in the ConnectionRequest.
// (Java: KnxNetIpConfiguration#connectionType)
func (m *Connection) getTunnelConnectionType() driverModel.KnxLayer {
	if val, ok := m.options["connection-type"]; ok && len(val) > 0 {
		switch strings.ToUpper(strings.TrimSpace(val[0])) {
		case "LINK_LAYER", "TUNNEL_LINK_LAYER":
			return driverModel.KnxLayer_TUNNEL_LINK_LAYER
		case "RAW", "TUNNEL_RAW":
			return driverModel.KnxLayer_TUNNEL_RAW
		case "BUSMONITOR", "TUNNEL_BUSMONITOR":
			return driverModel.KnxLayer_TUNNEL_BUSMONITOR
		default:
			m.log.Warn().Str("connection-type", val[0]).Msg("Invalid value for connection-type, falling back to LINK_LAYER")
		}
	}
	return driverModel.KnxLayer_TUNNEL_LINK_LAYER
}

// getRequestTimeout evaluates the "request-timeout-ms" connection option (in
// milliseconds) which limits how long we wait for a gateway reply.
// (Java: KnxNetIpConfiguration#requestTimeout)
func (m *Connection) getRequestTimeout() time.Duration {
	if val, ok := m.options["request-timeout-ms"]; ok && len(val) > 0 {
		requestTimeout, err := strconv.ParseUint(val[0], 10, 32)
		if err == nil && requestTimeout > 0 {
			return time.Duration(requestTimeout) * time.Millisecond
		}
		m.log.Warn().Str("request-timeout-ms", val[0]).Msg("Invalid value for request-timeout-ms, falling back to the default")
	}
	return defaultRequestTimeout
}

func (m *Connection) addSubscriber(subscriber *Subscriber) {
	m.subscribersMutex.Lock()
	defer m.subscribersMutex.Unlock()
	if slices.Contains(m.subscribers, subscriber) {
		m.log.Debug().Interface("subscriber", subscriber).Msg("Subscriber already added")
		return
	}
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *Connection) removeSubscriber(subscriber *Subscriber) {
	m.subscribersMutex.Lock()
	defer m.subscribersMutex.Unlock()
	m.subscribers = slices.DeleteFunc(m.subscribers, func(other *Subscriber) bool {
		return other == subscriber
	})
}

// getSubscribers returns a snapshot of the currently registered subscribers so the
// callers can iterate them without holding the lock while calling into them.
func (m *Connection) getSubscribers() []*Subscriber {
	m.subscribersMutex.RLock()
	defer m.subscribersMutex.RUnlock()
	return slices.Clone(m.subscribers)
}

// knxAddressEqual compares two individual knx addresses by value. The generated model
// types are interfaces backed by pointers, so "==" would compare identities instead.
func knxAddressEqual(a, b driverModel.KnxAddress) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.GetMainGroup() == b.GetMainGroup() &&
		a.GetMiddleGroup() == b.GetMiddleGroup() &&
		a.GetSubGroup() == b.GetSubGroup()
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
