//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package modbus

import (
    "errors"
    "fmt"
    "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "time"
)

type ModbusExpectation struct {
    expiration     time.Time
    acceptsMessage spi.AcceptsMessage
    handleMessage  spi.HandleMessage
}

type ModbusMessageCodec struct {
    expectationCounter            int32
	transportInstance             transports.TransportInstance
	defaultIncomingMessageChannel chan interface{}
	expectations                  []ModbusExpectation
}

func NewModbusMessageCodec(transportInstance transports.TransportInstance, defaultIncomingMessageChannel chan interface{}) *ModbusMessageCodec {
	codec := &ModbusMessageCodec{
        expectationCounter:            1,
		transportInstance:             transportInstance,
		defaultIncomingMessageChannel: defaultIncomingMessageChannel,
		expectations:                  []ModbusExpectation{},
	}
	// Start a worker that handles processing of responses
	go work(codec)
	return codec
}

func (m *ModbusMessageCodec) Connect() error {
	return m.transportInstance.Connect()
}

func (m *ModbusMessageCodec) Disconnect() error {
	return m.transportInstance.Close()
}

func (m *ModbusMessageCodec) Send(message interface{}) error {
	// Cast the message to the correct type of struct
	adu := model.CastModbusTcpADU(message)
	// Serialize the request
	wb := utils.NewWriteBuffer()
	err := adu.Serialize(*wb)
	if err != nil {
		return errors.New("error serializing request " + err.Error())
	}

	// Send it to the PLC
	err = m.transportInstance.Write(wb.GetBytes())
	if err != nil {
		return errors.New("error sending request " + err.Error())
	}
	return nil
}

func (m *ModbusMessageCodec) Expect(acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, ttl time.Duration) error {
    expectation := ModbusExpectation{
        expiration:     time.Now().Add(ttl),
        acceptsMessage: acceptsMessage,
        handleMessage:  handleMessage,
    }
    m.expectations = append(m.expectations, expectation)
    return nil
}

func (m *ModbusMessageCodec) SendRequest(message interface{}, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, ttl time.Duration) error {
    // Send the actual message
    err := m.Send(message)
    if err != nil {
        return err
    }
    return m.Expect(acceptsMessage, handleMessage, ttl)
}

func (m *ModbusMessageCodec) GetDefaultIncomingMessageChannel() chan interface{} {
    return m.defaultIncomingMessageChannel
}

func (m *ModbusMessageCodec) receive() (interface{}, error) {
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.transportInstance.GetNumReadableBytes(); (err == nil) && (num >= 6) {
		data, err := m.transportInstance.PeekReadableBytes(6)
		if err != nil {
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		packetSize := (uint32(data[4]) << 8) + uint32(data[5]) + 6
		if num >= packetSize {
			data, err = m.transportInstance.Read(packetSize)
			if err != nil {
				// TODO: Possibly clean up ...
				return nil, nil
			}
			rb := utils.NewReadBuffer(data)
			adu, err := model.ModbusTcpADUParse(rb, true)
			if err != nil {
				// TODO: Possibly clean up ...
				return nil, nil
			}
			return adu, nil
		}
	}
	return nil, nil
}

func work(m *ModbusMessageCodec) {
	// Start an endless loop
	// TODO: Provide some means to terminate this ...
    for {
        if len(m.expectations) > 0 {
            message, err := m.receive()
            if err != nil {
                fmt.Printf("got an error reading from transport %s", err.Error())
            } else if message != nil {
                now := time.Now()
                messageHandled := false
                // Go through all expectations
                for index, expectation := range m.expectations {
                    // Check if this expectation has expired.
                    if now.After(expectation.expiration) {
                        // Remove this expectation from the list.
                        m.expectations = append(m.expectations[:index], m.expectations[index+1:]...)
                        break
                    }

                    // Check if the current message matches the expectations
                    // If it does, let it handle the message.
                    if accepts := expectation.acceptsMessage(message); accepts {
                        err = expectation.handleMessage(message)
                        if err == nil {
                            messageHandled = true
                            // Remove this expectation from the list.
                            m.expectations = append(m.expectations[:index], m.expectations[index+1:]...)
                        }
                        break
                    }
                }

                // If the message has not been handled and a default handler is provided, call this ...
                if !messageHandled {
                    if m.defaultIncomingMessageChannel != nil {
                        m.defaultIncomingMessageChannel <- message
                    } else {
                        fmt.Printf("No handler registered for handling message %s", message)
                    }
                }
            }
        } else {
            // Sleep for 10ms
            time.Sleep(10 * time.Millisecond)
        }
    }
}

func (m ModbusMessageCodec) GetTransportInstance() transports.TransportInstance {
	return m.transportInstance
}
