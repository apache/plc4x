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
    "encoding/json"
    "errors"
    "fmt"
    "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
    "github.com/apache/plc4x/plc4go/pkg/plc4go"
    apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "net/url"
    "strconv"
)

type ModbusDriver struct {
    fieldHandler spi.PlcFieldHandler
    plc4go.PlcDriver
}

func NewModbusDriver() *ModbusDriver {
    return &ModbusDriver{
        fieldHandler: NewFieldHandler(),
    }
}

func (m ModbusDriver) GetProtocolCode() string {
    return "modbus"
}

func (m ModbusDriver) GetProtocolName() string {
    return "Modbus"
}

func (m ModbusDriver) GetDefaultTransport() string {
    return "tcp"
}

func (m ModbusDriver) CheckQuery(query string) error {
    _, err := m.fieldHandler.ParseQuery(query)
    return err
}

func (m ModbusDriver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
    // Get an the transport specified in the url
    transport, ok := transports[transportUrl.Scheme]
    if !ok {
        ch := make(chan plc4go.PlcConnectionConnectResult)
        ch <- plc4go.NewPlcConnectionConnectResult(nil, errors.New("couldn't find transport for given transport url "+transportUrl.String()))
        return ch
    }
    // Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
    options["defaultTcpPort"] = []string{"502"}
    // Have the transport create a new transport-instance.
    transportInstance, err := transport.CreateTransportInstance(transportUrl, options)
    if err != nil {
        ch := make(chan plc4go.PlcConnectionConnectResult)
        ch <- plc4go.NewPlcConnectionConnectResult(nil, errors.New("couldn't initialize transport configuration for given transport url "+transportUrl.String()))
        return ch
    }

    // Create a new codec for taking care of encoding/decoding of messages
    defaultChanel := make(chan interface{})
    go func() {
        for {
            msg := <-defaultChanel
            adu := msg.(model.ModbusTcpADU)
            serialized, err := json.Marshal(adu)
            if err != nil {
                fmt.Printf("got error serializing adu: %s\n", err.Error())
            } else {
                fmt.Printf("got message in the default handler %s\n", serialized)
            }
        }
    }()
    codec := NewModbusMessageCodec(transportInstance, nil)

    // If a unit-identifier was provided in the connection string use this, otherwise use the default of 1
    unitIdentifier := uint8(1)
    if value, ok := options["unit-identifier"]; ok {
        var intValue int
        intValue, err = strconv.Atoi(value[0])
        if err == nil {
            unitIdentifier = uint8(intValue)
        }
    }

    // Create the new connection
    connection := NewModbusConnection(unitIdentifier, codec, options, m.fieldHandler)
    return connection.Connect()
}

func (m ModbusDriver) SupportsDiscovery() bool {
    return false
}

func (m ModbusDriver) Discover(callback func(event apiModel.PlcDiscoveryEvent)) error {
    panic("implement me")
}
