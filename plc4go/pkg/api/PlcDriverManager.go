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

package plc4go

import (
	"context"
	"io"
	"net/url"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/options/converter"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// PlcDriverManager is the main entry point for PLC4Go applications
type PlcDriverManager interface {
	io.Closer
	// RegisterDriver Manually register a new driver
	RegisterDriver(driver PlcDriver)
	// ListDriverNames List the names of all drivers registered in the system
	ListDriverNames() []string
	// GetDriver Get access to a driver instance for a given driver-name
	GetDriver(driverName string) (PlcDriver, error)

	// GetConnection Get a connection to a remote PLC for a given plc4x connection-string
	GetConnection(connectionString string) <-chan PlcConnectionConnectResult

	// Discover Execute all available discovery methods on all available drivers using all transports
	Discover(callback func(event model.PlcDiscoveryItem), discoveryOptions ...WithDiscoveryOption) error
	// DiscoverWithContext Execute all available discovery methods on all available drivers using all transports
	DiscoverWithContext(ctx context.Context, callback func(event model.PlcDiscoveryItem), discoveryOptions ...WithDiscoveryOption) error
}

func NewPlcDriverManager(_options ...config.WithOption) PlcDriverManager {
	localLog := options.ExtractCustomLoggerOrDefaultToGlobal(converter.WithOptionToInternal(_options...)...)
	localLog.Trace().Msg("Creating plc driver manager")
	return &plcDriverManger{
		drivers:    map[string]PlcDriver{},
		transports: map[string]transports.Transport{},

		log: localLog,
	}
}

// WithDiscoveryOptionProtocol sets an option for a protocol
func WithDiscoveryOptionProtocol(protocolName string) WithDiscoveryOption {
	return withDiscoveryOption{options.WithDiscoveryOptionProtocol(protocolName)}
}

// WithDiscoveryOptionTransport sets an option for a transportName
func WithDiscoveryOptionTransport(transportName string) WithDiscoveryOption {
	return withDiscoveryOption{options.WithDiscoveryOptionTransport(transportName)}
}

// WithDiscoveryOptionDeviceName sets an option for a deviceName
func WithDiscoveryOptionDeviceName(deviceName string) WithDiscoveryOption {
	return withDiscoveryOption{options.WithDiscoveryOptionDeviceName(deviceName)}
}

// WithDiscoveryOptionLocalAddress sets an option for a localAddress
func WithDiscoveryOptionLocalAddress(localAddress string) WithDiscoveryOption {
	return withDiscoveryOption{options.WithDiscoveryOptionLocalAddress(localAddress)}
}

// WithDiscoveryOptionRemoteAddress sets an option for a remoteAddress
func WithDiscoveryOptionRemoteAddress(remoteAddress string) WithDiscoveryOption {
	return withDiscoveryOption{options.WithDiscoveryOptionRemoteAddress(remoteAddress)}
}

func WithDiscoveryOptionProtocolSpecific(key string, value any) WithDiscoveryOption {
	return withDiscoveryOption{options.WithDiscoveryOptionProtocolSpecific(key, value)}
}

// WithDiscoveryOption is a marker interface for options regarding discovery
// FIXME: this is to avoid leaks spi in the signature move to spi driver or create interfaces. Can also be done by moving spi in a proper module
type WithDiscoveryOption interface {
	isDiscoveryOption() bool
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

//go:generate plc4xGenerator -type=plcDriverManger
type plcDriverManger struct {
	drivers    map[string]PlcDriver
	transports map[string]transports.Transport

	log zerolog.Logger
}

type withDiscoveryOption struct {
	options.WithDiscoveryOption
}

func (w withDiscoveryOption) isDiscoveryOption() bool {
	return true
}

func convertToInternalOptions(withDiscoveryOptions ...WithDiscoveryOption) []options.WithDiscoveryOption {
	result := make([]options.WithDiscoveryOption, len(withDiscoveryOptions))
	for i, discoveryOption := range withDiscoveryOptions {
		result[i] = discoveryOption.(withDiscoveryOption).WithDiscoveryOption
	}
	return result
}

func (m *plcDriverManger) RegisterDriver(driver PlcDriver) {
	m.log.Debug().Str("protocolName", driver.GetProtocolName()).Msg("Registering driver")
	// If this driver is already registered, just skip resetting it
	for driverName := range m.drivers {
		if driverName == driver.GetProtocolCode() {
			m.log.Warn().Str("protocolName", driver.GetProtocolName()).Msg("Already registered")
			return
		}
	}
	m.drivers[driver.GetProtocolCode()] = driver
	m.log.Info().Str("protocolName", driver.GetProtocolName()).Msg("Driver for protocolName registered")
}

func (m *plcDriverManger) ListDriverNames() []string {
	m.log.Trace().Msg("Listing driver names")
	var driverNames []string
	for driverName := range m.drivers {
		driverNames = append(driverNames, driverName)
	}
	m.log.Trace().Int("nDrivers", len(driverNames)).Msg("Found nDrivers driver(s)")
	return driverNames
}

func (m *plcDriverManger) GetDriver(driverName string) (PlcDriver, error) {
	if val, ok := m.drivers[driverName]; ok {
		return val, nil
	}
	return nil, errors.Errorf("couldn't find driver %s", driverName)
}

func (m *plcDriverManger) RegisterTransport(transport transports.Transport) {
	m.log.Debug().Str("transportName", transport.GetTransportName()).Msg("Registering transport")
	// If this transport is already registered, just skip resetting it
	for transportName := range m.transports {
		if transportName == transport.GetTransportCode() {
			m.log.Warn().Str("transportName", transport.GetTransportName()).Msg("Transport already registered")
			return
		}
	}
	m.transports[transport.GetTransportCode()] = transport
	m.log.Info().Str("transportName", transport.GetTransportName()).Msg("Transport for transportName registered")
}

func (m *plcDriverManger) ListTransportNames() []string {
	m.log.Trace().Msg("Listing transport names")
	var transportNames []string
	for transportName := range m.transports {
		transportNames = append(transportNames, transportName)
	}
	m.log.Trace().Int("nTransports", len(transportNames)).Msg("Found nTransports transports")
	return transportNames
}

func (m *plcDriverManger) GetTransport(transportName string, _ string, _ map[string][]string) (transports.Transport, error) {
	if val, ok := m.transports[transportName]; ok {
		m.log.Debug().Str("transportName", transportName).Msg("Returning transport")
		return val, nil
	}
	return nil, errors.Errorf("couldn't find transport %s", transportName)
}

func (m *plcDriverManger) GetConnection(connectionString string) <-chan PlcConnectionConnectResult {
	m.log.Debug().Str("connectionString", connectionString).Msg("Getting connection for connectionString")
	// Parse the connection string.
	connectionUrl, err := url.Parse(connectionString)
	if err != nil {
		m.log.Error().Err(err).Msg("Error parsing connection")
		ch := make(chan PlcConnectionConnectResult, 1)
		ch <- &plcConnectionConnectResult{err: errors.Wrap(err, "error parsing connection string")}
		return ch
	}
	m.log.Debug().Stringer("connectionUrl", connectionUrl).Msg("parsed connection URL")

	// The options will be used to configure both the transports as well as the connections/drivers
	configOptions := connectionUrl.Query()

	// Find the driver specified in the url.
	driverName := connectionUrl.Scheme
	driver, err := m.GetDriver(driverName)
	if err != nil {
		m.log.Err(err).Str("driverName", driverName).Msg("Couldn't get driver for driverName")
		ch := make(chan PlcConnectionConnectResult, 1)
		ch <- &plcConnectionConnectResult{err: errors.Wrap(err, "error getting driver for connection string")}
		return ch
	}
	m.log.Debug().Stringer("connectionUrl", connectionUrl).Str("protocolName", driver.GetProtocolName()).Msg("got driver protocolName")

	// If a transport is provided alongside the driver, the URL content is decoded as "opaque" data
	// Then we have to re-parse that to get the transport code as well as the host & port information.
	var transportName string
	var transportConnectionString string
	var transportPath string
	if len(connectionUrl.Opaque) > 0 {
		m.log.Trace().Msg("we handling a opaque connectionUrl")
		connectionUrl, err := url.Parse(connectionUrl.Opaque)
		if err != nil {
			m.log.Err(err).Str("connectionUrl.Opaque", connectionUrl.Opaque).Msg("Couldn't get transport due to parsing error")
			ch := make(chan PlcConnectionConnectResult, 1)
			ch <- &plcConnectionConnectResult{err: errors.Wrap(err, "error parsing connection string")}
			return ch
		}
		transportName = connectionUrl.Scheme
		transportConnectionString = connectionUrl.Host
		transportPath = connectionUrl.Path
	} else {
		m.log.Trace().Msg("we handling a non-opaque connectionUrl")
		// If no transport was provided the driver has to provide a default transport.
		transportName = driver.GetDefaultTransport()
		transportConnectionString = connectionUrl.Host
		transportPath = connectionUrl.Path
	}
	m.log.Debug().
		Str("transportName", transportName).
		Str("transportConnectionString", transportConnectionString).
		Msg("got a transport")
	// If no transport has been specified explicitly or per default, we have to abort.
	if transportName == "" {
		m.log.Error().Msg("got a empty transport")
		ch := make(chan PlcConnectionConnectResult, 1)
		ch <- &plcConnectionConnectResult{err: errors.New("no transport specified and no default defined by driver")}
		return ch
	}

	// Assemble a correct transport url
	transportUrl := url.URL{
		Scheme: transportName,
		Host:   transportConnectionString,
		Path:   transportPath,
	}
	m.log.Debug().Stringer("transportUrl", &transportUrl).Msg("Assembled transport url")

	// Create a new connection
	return driver.GetConnection(transportUrl, m.transports, configOptions)
}

func (m *plcDriverManger) Discover(callback func(event model.PlcDiscoveryItem), discoveryOptions ...WithDiscoveryOption) error {
	return m.DiscoverWithContext(context.TODO(), callback, discoveryOptions...)
}

func (m *plcDriverManger) DiscoverWithContext(ctx context.Context, callback func(event model.PlcDiscoveryItem), discoveryOptions ...WithDiscoveryOption) error {
	// Check if we've got at least one option to restrict to certain protocols only.
	// If there is at least one, we only check that protocol, if there are none, all
	// available protocols are checked.
	internalOptions := convertToInternalOptions(discoveryOptions...)
	protocolOptions := options.FilterDiscoveryOptionsProtocol(internalOptions)
	discoveryDrivers := map[string]PlcDriver{}
	if len(protocolOptions) > 0 {
		for _, protocolOption := range protocolOptions {
			if driver, ok := m.drivers[protocolOption.GetProtocolName()]; ok {
				discoveryDrivers[driver.GetProtocolName()] = driver
			}
		}
	} else {
		discoveryDrivers = m.drivers
	}

	// Execute discovery on all selected drivers
	for _, driver := range discoveryDrivers {
		if driver.SupportsDiscovery() {
			err := driver.DiscoverWithContext(ctx, callback, internalOptions...)
			if err != nil {
				return errors.Wrapf(err, "Error running Discover on driver %s", driver.GetProtocolName())
			}
		}
	}
	return nil
}

func (m *plcDriverManger) Close() error {
	m.log.Info().Msg("Shutting down driver manager")
	var aggregatedErrors []error
	for s, driver := range m.drivers {
		m.log.Trace().Str("name", s).Msg("closing driver")
		if err := driver.Close(); err != nil {
			aggregatedErrors = append(aggregatedErrors, err)
		}
	}
	for s, transport := range m.transports {
		m.log.Trace().Str("name", s).Msg("closing transport")
		if err := transport.Close(); err != nil {
			aggregatedErrors = append(aggregatedErrors, err)
		}
	}
	if len(aggregatedErrors) > 0 {
		return utils.MultiError{
			MainError: errors.New("error closing everything"),
			Errors:    aggregatedErrors,
		}
	}
	return nil
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
