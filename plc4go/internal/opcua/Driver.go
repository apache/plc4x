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

package opcua

import (
	"context"
	"net/url"
	"regexp"
	"strconv"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Driver struct {
	_default.DefaultDriver
	awaitSetupComplete      bool
	awaitDisconnectComplete bool

	uriPattern *regexp.Regexp

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,

		uriPattern: regexp.MustCompile(`^((?P<protocolCode>opcua):)?(?P<transportCode>[a-z0-9]*)?:(//)?(?P<transportHost>[\w.-]+)(:(?P<transportPort>\d*))?(?P<transportEndpoint>[\w/=]*)[\\?]?(?P<paramString>([^\\=]+=[^\\=&]+&?)*)`),

		log:      customLogger,
		_options: _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "opcua", "Opcua", "tcp", NewTagHandler())
	return driver
}

func (m *Driver) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	m.log.Debug().
		Stringer("transportUrl", &transportUrl).
		Int("numberTransports", len(transports)).
		Int("numberDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url")

	// Get the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		m.log.Error().Stringer("transportUrl", &transportUrl).Str("scheme", transportUrl.Scheme).Msg("We couldn't find a transport for scheme")
		return m.reportError(errors.Errorf("couldn't find transport for given transport url %v", transportUrl))
	}

	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultTcpPort"] = []string{strconv.FormatUint(uint64(80), 10)}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(
		transportUrl,
		driverOptions,
		append(m._options, options.WithCustomLogger(m.log))...,
	)
	if err != nil {
		m.log.Error().Err(err).Stringer("transportUrl", &transportUrl).Strs("defaultTcpPort", driverOptions["defaultTcpPort"]).Msg("We couldn't create a transport instance for port")
		return m.reportError(errors.Wrapf(err, "couldn't initialize transport configuration for given transport url %s", transportUrl.String()))
	}

	// Split up the connection string into its individual segments.
	var protocolCode, transportCode, transportHost, transportPort, transportEndpoint, paramString string
	if match := utils.GetSubgroupMatches(m.uriPattern, transportUrl.String()); match != nil {
		protocolCode = match["protocolCode"]
		if protocolCode == "" {
			protocolCode = m.GetProtocolCode()
		}
		transportCode = match["transportCode"]
		if transportCode == "" {
			transportCode = m.GetDefaultTransport()
		}
		transportHost = match["transportHost"]
		transportPort = match["transportPort"]
		transportEndpoint = match["transportEndpoint"]
		paramString = match["paramString"]
		_ = paramString // TODO: not sure if we need that
	} else {
		return m.reportError(errors.Errorf("Connection string %s doesn't match the format %s", &transportUrl, m.uriPattern))
	}

	// Check if the protocol code matches this driver.
	if protocolCode != m.GetProtocolCode() {
		// Actually this shouldn't happen as the DriverManager should have not used this driver in the first place.
		return m.reportError(errors.New("This driver is not suited to handle this connection string"))
	}

	// Create the configuration object.
	configuration, err := ParseFromOptions(m.log, driverOptions)
	if err != nil {
		return m.reportError(errors.Wrap(err, "can't parse options"))
	}
	configuration.host = transportHost
	configuration.port = transportPort
	configuration.transportEndpoint = transportEndpoint
	configuration.endpoint = "opc." + transportCode + "://" + transportHost + ":" + transportPort + "" + transportEndpoint

	if securityPolicy := configuration.securityPolicy; securityPolicy != "" && securityPolicy != "None" {
		m.log.Trace().Str("securityPolicy", securityPolicy).Msg("working with security policy")
		if err := configuration.openKeyStore(); err != nil {
			return m.reportError(errors.Wrap(err, "error opening key store"))
		}
	} else {
		m.log.Trace().Msg("no security policy")
	}

	driverContext := NewDriverContext(configuration)
	driverContext.awaitSetupComplete = m.awaitSetupComplete
	driverContext.awaitDisconnectComplete = m.awaitDisconnectComplete

	codec := NewMessageCodec(
		transportInstance,
		NewSecureChannel(m.log, driverContext, configuration),
		append(m._options, options.WithCustomLogger(m.log))...,
	)
	m.log.Debug().Msgf("working with codec:\n%v", codec)

	// Create the new connection
	connection := NewConnection(
		codec, configuration,
		driverContext,
		m.GetPlcTagHandler(),
		driverOptions,
		append(m._options, options.WithCustomLogger(m.log))...,
	)
	m.log.Debug().Msg("created connection, connecting now")
	return connection.ConnectWithContext(ctx)
}

func (m *Driver) SetAwaitSetupComplete(awaitComplete bool) {
	m.awaitSetupComplete = awaitComplete
}

func (m *Driver) SetAwaitDisconnectComplete(awaitComplete bool) {
	m.awaitDisconnectComplete = awaitComplete
}

func (m *Driver) reportError(err error) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	ch <- _default.NewDefaultPlcConnectionConnectResult(nil, err)
	return ch
}
