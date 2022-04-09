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

package ads

import (
	_default "github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Driver struct {
	_default.DefaultDriver
}

func NewDriver() plc4go.PlcDriver {
	return &Driver{
		DefaultDriver: _default.NewDefaultDriver("ads", "Beckhoff TwinCat ADS", "tcp", NewFieldHandler()),
	}
}

func (m *Driver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	log.Debug().Stringer("transportUrl", &transportUrl).Msgf("Get connection for transport url with %d transport(s) and %d option(s)", len(transports), len(options))
	// Get an the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't find a transport for scheme %s", transportUrl.Scheme)
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		}()
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	options["defaultTcpPort"] = []string{"48898"}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(transportUrl, options)
	if err != nil {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't create a transport instance for port %#v", options["defaultTcpPort"])
		ch := make(chan plc4go.PlcConnectionConnectResult)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.New("couldn't initialize transport configuration for given transport url "+transportUrl.String()))
		return ch
	}

	// Create a new codec for taking care of encoding/decoding of messages
	codec := NewMessageCodec(transportInstance)
	log.Debug().Msgf("working with codec %#v", codec)

	configuration, err := ParseFromOptions(options)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid options")
		ch := make(chan plc4go.PlcConnectionConnectResult)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "invalid configuration"))
		return ch
	}

	// Create the new connection
	connection, err := NewConnection(codec, configuration, m.GetPlcFieldHandler(), options)
	if err != nil {
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "couldn't create connection"))
		}()
		return ch
	}
	log.Debug().Stringer("connection", connection).Msg("created connection, connecting now")
	return connection.Connect()
}
