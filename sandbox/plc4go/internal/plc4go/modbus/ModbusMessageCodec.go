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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/modbus/readwrite/model"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/transports"
    "time"
)

type ModbusExpectation struct {
    timestamp time.Time
    check func(interface{}) bool
    responseChannel chan interface{}
}

type ModbusMessageCodec struct {
    transportInstance transports.TransportInstance
    defaultIncomingMessageChannel chan interface{}
    expectations []ModbusExpectation
}

func NewModbusMessageCodec(transportInstance transports.TransportInstance, defaultIncomingMessageChannel chan interface{}) *ModbusMessageCodec {
    codec := &ModbusMessageCodec{
        transportInstance: transportInstance,
        defaultIncomingMessageChannel: defaultIncomingMessageChannel,
        expectations: []ModbusExpectation{},
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
    wb := spi.WriteBufferNew()
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

func (m *ModbusMessageCodec) Receive() (interface{}, error) {
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
            rb := spi.ReadBufferNew(data)
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

func (m *ModbusMessageCodec) Expect(check func(interface{}) bool) chan interface{} {
    responseChanel := make(chan interface{})
    expectation := ModbusExpectation{
        timestamp: time.Now(),
        check: check,
        responseChannel: responseChanel,
    }
    m.expectations = append(m.expectations, expectation)
    return responseChanel
}

func work(codec *ModbusMessageCodec) {
    // Start an endless loop
    // TODO: Provide some means to terminate this ...
    for {
        if len(codec.expectations) > 0 {
            message, err := codec.Receive()
            if err != nil {
                fmt.Errorf("got an error reading from transport %s", err.Error())
            } else if message != nil {
                messageHandled := false
                // Go through all expectations
                for index, expectation := range codec.expectations {
                    // Check if the current message matches the expectations
                    if expectation.check(message) {
                        // Send the decoded message back
                        expectation.responseChannel <- message
                        // Remove this expectation from the list.
                        codec.expectations = append(codec.expectations[:index], codec.expectations[index+1:]...)
                        messageHandled = true
                        break
                    }
                }
                // If the message has not been handled and a default handler is provided, call this ...
                if !messageHandled {
                    if codec.defaultIncomingMessageChannel != nil {
                        codec.defaultIncomingMessageChannel <- message
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


