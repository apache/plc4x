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
package plc4go

import (
	"errors"
	"net/url"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/bacnetip"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/knxnetip"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/modbus"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/s7"
)

// This is the main entry point for PLC4Go applications
type PlcDriverManager interface {
	// List the names of all drivers registered in the system
	ListDriverNames() []string
	// Get access to a driver instance for a given driver-name
	GetDriver(driverName string) (PlcDriver, error)

	// Get a connection to a remote PLC for a given plc4x connection-string
	GetConnectedConnection(connectionString string) <-chan PlcConnectionConnectResult
}

func NewPlcDriverManager() PlcDriverManager {
	return plcDriverManger{
		drivers: map[string]PlcDriver{
			"bacnetip": bacnetip.NewBacnetIpDriver(),
			"knxnetip": knxnetip.NewKnxNetIpDriver(),
			"modbus":   modbus.NewModbusDriver(),
			"s7":       s7.NewS7Driver(),
		},
	}
}

type plcDriverManger struct {
	drivers map[string]PlcDriver
}

func (m plcDriverManger) ListDriverNames() []string {
	driverNames := make([]string, len(m.drivers))
	for driverName := range m.drivers {
		driverNames = append(driverNames, driverName)
	}
	return driverNames
}

func (m plcDriverManger) GetDriver(driverName string) (PlcDriver, error) {
	if m.drivers[driverName] == nil {
		return nil, errors.New("Couldn't find driver " + driverName)
	}
	return m.drivers[driverName], nil
}

func (m plcDriverManger) GetConnectedConnection(connectionString string) <-chan PlcConnectionConnectResult {
	connectionUrl, err := url.Parse(connectionString)
	if err != nil {
		ch := make(chan PlcConnectionConnectResult)
		ch <- NewPlcConnectionConnectResult(nil, errors.New("Error parsing connection string: "+err.Error()))
		return ch
	}
	driverName := connectionUrl.Scheme
	driver, err := m.GetDriver(driverName)
	if err != nil {
		ch := make(chan PlcConnectionConnectResult)
		ch <- NewPlcConnectionConnectResult(nil, errors.New("Error getting driver for connection string: "+err.Error()))
		return ch
	}
	connection, err := driver.GetConnection(connectionString)
	if err != nil {
		ch := make(chan PlcConnectionConnectResult)
		ch <- NewPlcConnectionConnectResult(nil, errors.New("Error connecting for connection string: "+err.Error()))
		return ch
	}
	ch := make(chan PlcConnectionConnectResult)
	go func() {
		errConsumer := connection.Connect()
		connectionErr := <-errConsumer
		ch <- connectionErr
	}()
	return ch
}
