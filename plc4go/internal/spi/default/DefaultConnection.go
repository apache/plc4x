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
	"github.com/apache/plc4x/plc4go/internal/spi/options"
	"time"

	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DefaultConnectionRequirements defines the required at a implementing connection when using DefaultConnection
// additional options can be set using the functions returning WithOption (e.g. WithDefaultTtl, WithPlcFieldHandler...)
type DefaultConnectionRequirements interface {
	// GetConnection should return the implementing connection when using DefaultConnection
	GetConnection() plc4go.PlcConnection
	// GetMessageCodec should return the spi.MessageCodec in use
	GetMessageCodec() spi.MessageCodec
}

// DefaultConnection should be used as an embedded struct. All defined methods here have default implementations
type DefaultConnection interface {
	plc4go.PlcConnection
	spi.TransportInstanceExposer
	spi.HandlerExposer
	SetConnected(connected bool)
	GetTtl() time.Duration
}

// NewDefaultConnection is the factory for a DefaultConnection
func NewDefaultConnection(requirements DefaultConnectionRequirements, options ...options.WithOption) DefaultConnection {
	return buildDefaultConnection(requirements, options...)
}

// WithDefaultTtl ttl is time.Second * 10 by default
func WithDefaultTtl(defaultTtl time.Duration) options.WithOption {
	return withDefaultTtl{defaultTtl: defaultTtl}
}

func WithPlcFieldHandler(plcFieldHandler spi.PlcFieldHandler) options.WithOption {
	return withPlcFieldHandler{plcFieldHandler: plcFieldHandler}
}

func WithPlcValueHandler(plcValueHandler spi.PlcValueHandler) options.WithOption {
	return withPlcValueHandler{plcValueHandler: plcValueHandler}
}

// DefaultConnectionMetadata implements the model.PlcConnectionMetadata interface
type DefaultConnectionMetadata struct {
	ConnectionAttributes map[string]string
	ProvidesReading      bool
	ProvidesWriting      bool
	ProvidesSubscribing  bool
	ProvidesBrowsing     bool
}

type DefaultPlcConnectionConnectResult interface {
	plc4go.PlcConnectionConnectResult
}

func NewDefaultPlcConnectionConnectResult(connection plc4go.PlcConnection, err error) DefaultPlcConnectionConnectResult {
	return &plcConnectionConnectResult{
		connection: connection,
		err:        err,
	}
}

type DefaultPlcConnectionCloseResult interface {
	plc4go.PlcConnectionCloseResult
	GetTraces() []spi.TraceEntry
}

func NewDefaultPlcConnectionCloseResult(connection plc4go.PlcConnection, err error) plc4go.PlcConnectionCloseResult {
	return &plcConnectionCloseResult{
		connection: connection,
		err:        err,
		traces:     nil,
	}
}

func NewDefaultPlcConnectionCloseResultWithTraces(connection plc4go.PlcConnection, err error, traces []spi.TraceEntry) plc4go.PlcConnectionCloseResult {
	return &plcConnectionCloseResult{
		connection: connection,
		err:        err,
		traces:     traces,
	}
}

type DefaultPlcConnectionPingResult interface {
	plc4go.PlcConnectionPingResult
}

func NewDefaultPlcConnectionPingResult(err error) plc4go.PlcConnectionPingResult {
	return &plcConnectionPingResult{
		err: err,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withDefaultTtl struct {
	options.Option
	// defaultTtl the time to live after a close
	defaultTtl time.Duration
}

type withPlcFieldHandler struct {
	options.Option
	plcFieldHandler spi.PlcFieldHandler
}

type withPlcValueHandler struct {
	options.Option
	plcValueHandler spi.PlcValueHandler
}

type defaultConnection struct {
	DefaultConnectionRequirements
	// defaultTtl the time to live after a close
	defaultTtl time.Duration
	// connected indicates if a connection is connected
	connected    bool
	fieldHandler spi.PlcFieldHandler
	valueHandler spi.PlcValueHandler
}

func buildDefaultConnection(requirements DefaultConnectionRequirements, options ...options.WithOption) DefaultConnection {
	defaultTtl := time.Second * 10
	var fieldHandler spi.PlcFieldHandler
	var valueHandler spi.PlcValueHandler

	for _, option := range options {
		switch option.(type) {
		case withDefaultTtl:
			defaultTtl = option.(withDefaultTtl).defaultTtl
		case withPlcFieldHandler:
			fieldHandler = option.(withPlcFieldHandler).plcFieldHandler
		case withPlcValueHandler:
			valueHandler = option.(withPlcValueHandler).plcValueHandler
		}
	}

	return &defaultConnection{
		requirements,
		defaultTtl,
		false,
		fieldHandler,
		valueHandler,
	}
}

type plcConnectionConnectResult struct {
	connection plc4go.PlcConnection
	err        error
}

func (d *plcConnectionConnectResult) GetConnection() plc4go.PlcConnection {
	return d.connection
}

func (d *plcConnectionConnectResult) GetErr() error {
	return d.err
}

type plcConnectionCloseResult struct {
	connection plc4go.PlcConnection
	err        error
	traces     []spi.TraceEntry
}

func (d *plcConnectionCloseResult) GetConnection() plc4go.PlcConnection {
	return d.connection
}

func (d *plcConnectionCloseResult) GetErr() error {
	return d.err
}

func (d plcConnectionCloseResult) GetTraces() []spi.TraceEntry {
	return d.traces
}

type plcConnectionPingResult struct {
	err error
}

func (d *plcConnectionPingResult) GetErr() error {
	return d.err
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (d *defaultConnection) SetConnected(connected bool) {
	d.connected = connected
}

func (d *defaultConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult)
	go func() {
		err := d.GetMessageCodec().Connect()
		d.SetConnected(true)
		connection := d.GetConnection()
		ch <- NewDefaultPlcConnectionConnectResult(connection, err)
	}()
	return ch
}

func (d *defaultConnection) BlockingClose() {
	log.Trace().Msg("blocking close connection")
	closeResults := d.GetConnection().Close()
	timeout := time.NewTimer(d.GetTtl())
	d.SetConnected(false)
	select {
	case <-closeResults:
		if !timeout.Stop() {
			<-timeout.C
		}
		return
	case <-timeout.C:
		timeout.Stop()
		return
	}
}

func (d *defaultConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	log.Trace().Msg("close connection")
	err := d.GetTransportInstance().Close()
	d.SetConnected(false)
	ch := make(chan plc4go.PlcConnectionCloseResult)
	go func() {
		ch <- NewDefaultPlcConnectionCloseResult(d.GetConnection(), err)
	}()
	return ch
}

func (d *defaultConnection) IsConnected() bool {
	// TODO: should we check here if the transport is connected?
	return d.connected
}

func (d *defaultConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	ch := make(chan plc4go.PlcConnectionPingResult)
	go func() {
		if d.GetConnection().IsConnected() {
			ch <- NewDefaultPlcConnectionPingResult(nil)
		} else {
			ch <- NewDefaultPlcConnectionPingResult(errors.New("not connected"))
		}
	}()
	return ch
}

func (d *defaultConnection) GetTtl() time.Duration {
	return d.defaultTtl
}

func (d *defaultConnection) GetMetadata() model.PlcConnectionMetadata {
	return DefaultConnectionMetadata{
		ConnectionAttributes: nil,
		ProvidesReading:      false,
		ProvidesWriting:      false,
		ProvidesSubscribing:  false,
		ProvidesBrowsing:     false,
	}
}

func (d *defaultConnection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	panic("not implemented")
}

func (d *defaultConnection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	panic("not implemented")
}

func (d *defaultConnection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	panic("not implemented")
}

func (d *defaultConnection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	panic("not implemented")
}

func (d *defaultConnection) BrowseRequestBuilder() model.PlcBrowseRequestBuilder {
	panic("not implemented")
}

func (d *defaultConnection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := d.GetMessageCodec().(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (d *defaultConnection) GetPlcFieldHandler() spi.PlcFieldHandler {
	return d.fieldHandler
}

func (d *defaultConnection) GetPlcValueHandler() spi.PlcValueHandler {
	return d.valueHandler
}

func (m DefaultConnectionMetadata) GetConnectionAttributes() map[string]string {
	return m.ConnectionAttributes
}

func (m DefaultConnectionMetadata) CanRead() bool {
	return m.ProvidesReading
}

func (m DefaultConnectionMetadata) CanWrite() bool {
	return m.ProvidesWriting
}

func (m DefaultConnectionMetadata) CanSubscribe() bool {
	return m.ProvidesSubscribing
}

func (m DefaultConnectionMetadata) CanBrowse() bool {
	return m.ProvidesBrowsing
}
