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
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Driver struct {
	_default.DefaultDriver
}

func NewDriver() *Driver {
	return &Driver{
		DefaultDriver: _default.NewDefaultDriver("knxnet-ip", "KNXNet/IP", "udp", NewFieldHandler()),
	}
}

func (m Driver) CheckQuery(query string) error {
	_, err := m.GetPlcFieldHandler().ParseQuery(query)
	return err
}

func (m Driver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	// Get an the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		}()
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	options["defaultUdpPort"] = []string{"3671"}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(transportUrl, options)
	if err != nil {
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't initialize transport configuration for given transport url %#v", transportUrl))
		}()
		return ch
	}

	// Create the new connection
	connection := NewConnection(transportInstance, options, m.GetPlcFieldHandler())
	log.Trace().Str("transport", transportUrl.String()).Stringer("connection", connection).Msg("created new connection instance, trying to connect now")
	return connection.Connect()
}

func (m Driver) SupportsDiscovery() bool {
	return true
}

func (m Driver) Discover(callback func(event apiModel.PlcDiscoveryEvent), discoveryOptions ...options.WithDiscoveryOption) error {
	return NewDiscoverer().Discover(callback, discoveryOptions...)
}
