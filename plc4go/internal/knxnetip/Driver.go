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
	"github.com/rs/zerolog"
	"net/url"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
)

type Driver struct {
	_default.DefaultDriver

	log zerolog.Logger // TODO: use it
}

func NewDriver(_options ...options.WithOption) *Driver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		log: customLogger,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "knxnet-ip", "KNXNet/IP", "udp", NewTagHandler())
	return driver
}

func (m *Driver) CheckQuery(query string) error {
	_, err := m.GetPlcTagHandler().ParseQuery(query)
	return err
}

func (m *Driver) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	// Get an the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultUdpPort"] = []string{"3671"}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(transportUrl, driverOptions, options.WithCustomLogger(m.log))
	if err != nil {
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't initialize transport configuration for given transport url %#v", transportUrl))
		return ch
	}

	// Create the new connection
	connection := NewConnection(transportInstance, driverOptions, m.GetPlcTagHandler(), options.WithCustomLogger(m.log))
	m.log.Trace().Str("transport", transportUrl.String()).Stringer("connection", connection).Msg("created new connection instance, trying to connect now")
	return connection.ConnectWithContext(ctx)
}

func (m *Driver) SupportsDiscovery() bool {
	return true
}

func (m *Driver) DiscoverWithContext(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return NewDiscoverer().Discover(ctx, callback, discoveryOptions...)
}
