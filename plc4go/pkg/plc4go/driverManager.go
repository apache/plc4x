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
    "github.com/apache/plc4x/plc4go/internal/plc4go/transports"
    "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "net/url"
)

// This is the main entry point for PLC4Go applications
type PlcDriverManager interface {
    // Manually register a new driver
    RegisterDriver(driver PlcDriver)
    // List the names of all drivers registered in the system
    ListDriverNames() []string
    // Get access to a driver instance for a given driver-name
    GetDriver(driverName string) (PlcDriver, error)

    // Manually register a new driver
    RegisterTransport(transport transports.Transport)
    // List the names of all drivers registered in the system
    ListTransportNames() []string
    // Get access to a driver instance for a given driver-name
    GetTransport(transportName string, connectionString string, options map[string][]string) (transports.Transport, error)

    // Get a connection to a remote PLC for a given plc4x connection-string
    GetConnection(connectionString string) <-chan PlcConnectionConnectResult

    // Execute all available discovery methods on all available drivers using all transports
    Discover(func(event model.PlcDiscoveryEvent)) error
}

type PlcDriverManger struct {
    drivers    map[string]PlcDriver
    transports map[string]transports.Transport
}

func NewPlcDriverManager() PlcDriverManager {
    return PlcDriverManger{
        drivers:    map[string]PlcDriver{},
        transports: map[string]transports.Transport{},
    }
}

func (m PlcDriverManger) RegisterDriver(driver PlcDriver) {
    // If this driver is already registered, just skip resetting it
    for driverName, _ := range m.drivers {
        if driverName == driver.GetProtocolCode() {
            return
        }
    }
    m.drivers[driver.GetProtocolCode()] = driver
}

func (m PlcDriverManger) ListDriverNames() []string {
    var driverNames []string
    for driverName := range m.drivers {
        driverNames = append(driverNames, driverName)
    }
    return driverNames
}

func (m PlcDriverManger) GetDriver(driverName string) (PlcDriver, error) {
    if val, ok := m.drivers[driverName]; ok {
        return val, nil
    }
    return nil, errors.New("couldn't find driver " + driverName)
}

func (m PlcDriverManger) RegisterTransport(transport transports.Transport) {
    // If this transport is already registered, just skip resetting it
    for transportName, _ := range m.transports {
        if transportName == transport.GetTransportCode() {
            return
        }
    }
    m.transports[transport.GetTransportCode()] = transport
}

func (m PlcDriverManger) ListTransportNames() []string {
    var transportNames []string
    for transportName := range m.transports {
        transportNames = append(transportNames, transportName)
    }
    return transportNames
}

func (m PlcDriverManger) GetTransport(transportName string, connectionString string, options map[string][]string) (transports.Transport, error) {
    if val, ok := m.transports[transportName]; ok {
        return val, nil
    }
    return nil, errors.New("couldn't find transport " + transportName)
}

func (m PlcDriverManger) GetConnection(connectionString string) <-chan PlcConnectionConnectResult {
    // Parse the connection string.
    connectionUrl, err := url.Parse(connectionString)
    if err != nil {
        ch := make(chan PlcConnectionConnectResult)
        ch <- NewPlcConnectionConnectResult(nil, errors.New("error parsing connection string: "+err.Error()))
        return ch
    }

    // The options will be used to configure both the transports as well as the connections/drivers
    configOptions := connectionUrl.Query()

    // Find the driver specified in the url.
    driverName := connectionUrl.Scheme
    driver, err := m.GetDriver(driverName)
    if err != nil {
        ch := make(chan PlcConnectionConnectResult)
        ch <- NewPlcConnectionConnectResult(nil, errors.New("error getting driver for connection string: "+err.Error()))
        return ch
    }

    // If a transport is provided alongside the driver, the URL content is decoded as "opaque" data
    // Then we have to re-parse that to get the transport code as well as the host & port information.
    var transportName string
    var transportConnectionString string
    if len(connectionUrl.Opaque) > 0 {
        connectionUrl, err := url.Parse(connectionUrl.Opaque)
        if err != nil {
            ch := make(chan PlcConnectionConnectResult)
            ch <- NewPlcConnectionConnectResult(nil, errors.New("error parsing connection string: "+err.Error()))
            return ch
        }
        transportName = connectionUrl.Scheme
        transportConnectionString = connectionUrl.Host
    } else {
        // If no transport was provided the driver has to provide a default transport.
        transportName = driver.GetDefaultTransport()
        transportConnectionString = connectionUrl.Host
    }
    // If no transport has been specified explicitly or per default, we have to abort.
    if transportName == "" {
        ch := make(chan PlcConnectionConnectResult)
        ch <- NewPlcConnectionConnectResult(nil, errors.New("no transport specified and no default defined by driver"))
        return ch
    }

    // Assemble a correct transport url
    transportUrl := url.URL{
        Scheme: transportName,
        Host:   transportConnectionString,
    }

    // Create a new connection
    return driver.GetConnection(transportUrl, m.transports, configOptions)
}

func (m PlcDriverManger) Discover(callback func(event model.PlcDiscoveryEvent)) error {
    for _, driver := range m.drivers {
        if driver.SupportsDiscovery() {
            err := driver.Discover(callback)
            if err != nil {
                return errors.New("Error running Discover on driver " + driver.GetProtocolName() +
                    ". Got error: " + err.Error())
            }
        }
    }
    return nil
}
