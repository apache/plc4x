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

package _default

import (
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"net/url"
)

type DefaultDriver interface {
	fmt.Stringer
	plc4go.PlcDriver
	spi.PlcDiscoverer
	GetPlcFieldHandler() spi.PlcFieldHandler
}

func NewDefaultDriver(protocolCode string, protocolName string, defaultTransport string, plcFieldHandler spi.PlcFieldHandler) DefaultDriver {
	return &defaultDriver{
		protocolCode:     protocolCode,
		protocolName:     protocolName,
		defaultTransport: defaultTransport,
		plcFieldHandler:  plcFieldHandler,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type defaultDriver struct {
	protocolCode     string
	protocolName     string
	defaultTransport string
	plcFieldHandler  spi.PlcFieldHandler
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (d *defaultDriver) GetProtocolCode() string {
	return d.protocolCode
}

func (d *defaultDriver) GetProtocolName() string {
	return d.protocolName
}

func (d *defaultDriver) GetDefaultTransport() string {
	return d.defaultTransport
}

func (d *defaultDriver) CheckQuery(query string) error {
	_, err := d.plcFieldHandler.ParseQuery(query)
	return err
}

func (d *defaultDriver) GetConnection(_ url.URL, _ map[string]transports.Transport, _ map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	panic("implement me")
}

func (d *defaultDriver) SupportsDiscovery() bool {
	return false
}

func (d *defaultDriver) Discover(_ func(event apiModel.PlcDiscoveryEvent), _ ...options.WithDiscoveryOption) error {
	panic("not available")
}

func (d *defaultDriver) DiscoverWithContext(_ context.Context, callback func(event apiModel.PlcDiscoveryEvent), discoveryOptions ...options.WithDiscoveryOption) error {
	panic("not available")
}

func (d *defaultDriver) GetPlcFieldHandler() spi.PlcFieldHandler {
	return d.plcFieldHandler
}

func (d *defaultDriver) String() string {
	return fmt.Sprintf("%s (%s) [%s]", d.protocolName, d.protocolCode, d.defaultTransport)
}
