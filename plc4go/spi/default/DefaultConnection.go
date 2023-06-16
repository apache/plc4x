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
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
)

// DefaultConnectionRequirements defines the required at a implementing connection when using DefaultConnection
// additional options can be set using the functions returning WithOption (e.g. WithDefaultTtl, WithPlcTagHandler...)
type DefaultConnectionRequirements interface {
	// GetConnection should return the implementing connection when using DefaultConnection
	GetConnection() plc4go.PlcConnection
	// GetMessageCodec should return the spi.MessageCodec in use
	GetMessageCodec() spi.MessageCodec
	// ConnectWithContext is declared here for Connect redirection
	ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult
}

// DefaultConnection should be used as an embedded struct. All defined methods here have default implementations
type DefaultConnection interface {
	utils.Serializable
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

func WithPlcTagHandler(tagHandler spi.PlcTagHandler) options.WithOption {
	return withPlcTagHandler{plcTagHandler: tagHandler}
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
	GetTraces() []tracer.TraceEntry
}

func NewDefaultPlcConnectionCloseResult(connection plc4go.PlcConnection, err error) plc4go.PlcConnectionCloseResult {
	return &plcConnectionCloseResult{
		connection: connection,
		err:        err,
		traces:     nil,
	}
}

func NewDefaultPlcConnectionCloseResultWithTraces(connection plc4go.PlcConnection, err error, traces []tracer.TraceEntry) plc4go.PlcConnectionCloseResult {
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

type withPlcTagHandler struct {
	options.Option
	plcTagHandler spi.PlcTagHandler
}

type withPlcValueHandler struct {
	options.Option
	plcValueHandler spi.PlcValueHandler
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=defaultConnection
type defaultConnection struct {
	DefaultConnectionRequirements `ignore:"true"`
	// defaultTtl the time to live after a close
	defaultTtl time.Duration `stringer:"true"`
	// connected indicates if a connection is connected
	connected    atomic.Bool
	tagHandler   spi.PlcTagHandler
	valueHandler spi.PlcValueHandler

	log zerolog.Logger `ignore:"true"`
}

func buildDefaultConnection(requirements DefaultConnectionRequirements, _options ...options.WithOption) DefaultConnection {
	defaultTtl := 10 * time.Second
	var tagHandler spi.PlcTagHandler
	var valueHandler spi.PlcValueHandler

	for _, option := range _options {
		switch option.(type) {
		case withDefaultTtl:
			defaultTtl = option.(withDefaultTtl).defaultTtl
		case withPlcTagHandler:
			tagHandler = option.(withPlcTagHandler).plcTagHandler
		case withPlcValueHandler:
			valueHandler = option.(withPlcValueHandler).plcValueHandler
		}
	}

	return &defaultConnection{
		DefaultConnectionRequirements: requirements,
		defaultTtl:                    defaultTtl,
		tagHandler:                    tagHandler,
		valueHandler:                  valueHandler,

		log: options.ExtractCustomLogger(_options...),
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
	traces     []tracer.TraceEntry
}

func (d *plcConnectionCloseResult) GetConnection() plc4go.PlcConnection {
	return d.connection
}

func (d *plcConnectionCloseResult) GetErr() error {
	return d.err
}

func (d *plcConnectionCloseResult) GetTraces() []tracer.TraceEntry {
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
	d.log.Trace().Msgf("set connected %t", connected)
	d.connected.Store(connected)
}

func (d *defaultConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	return d.DefaultConnectionRequirements.ConnectWithContext(context.Background())
}

func (d *defaultConnection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	d.log.Trace().Msg("Connecting")
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		err := d.GetMessageCodec().ConnectWithContext(ctx)
		d.SetConnected(true)
		connection := d.GetConnection()
		ch <- NewDefaultPlcConnectionConnectResult(connection, err)
	}()
	return ch
}

func (d *defaultConnection) BlockingClose() {
	d.log.Trace().Msg("blocking close connection")
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
	d.log.Trace().Msg("close connection")
	if messageCodec := d.GetMessageCodec(); messageCodec != nil {
		d.log.Trace().Msgf("disconnecting message codec")
		if err := messageCodec.Disconnect(); err != nil {
			d.log.Warn().Err(err).Msg("Error disconnecting message code")
		} else {
			d.log.Trace().Msg("message codec disconnected")
		}
	}
	var err error
	if transportInstance := d.GetTransportInstance(); transportInstance != nil {
		d.log.Trace().Msg("closing transport instance")
		err = transportInstance.Close()
		d.log.Trace().Err(err).Msg("transport instance closed")
	}
	d.SetConnected(false)
	ch := make(chan plc4go.PlcConnectionCloseResult, 1)
	ch <- NewDefaultPlcConnectionCloseResult(d.GetConnection(), err)
	return ch
}

func (d *defaultConnection) IsConnected() bool {
	// TODO: should we check here if the transport is connected?
	return d.connected.Load()
}

func (d *defaultConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	ch := make(chan plc4go.PlcConnectionPingResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- NewDefaultPlcConnectionPingResult(errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
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

func (d *defaultConnection) GetMetadata() apiModel.PlcConnectionMetadata {
	return DefaultConnectionMetadata{
		ConnectionAttributes: nil,
		ProvidesReading:      false,
		ProvidesWriting:      false,
		ProvidesSubscribing:  false,
		ProvidesBrowsing:     false,
	}
}

func (d *defaultConnection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	panic("not provided by actual connection")
}

func (d *defaultConnection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	panic("not provided by actual connection")
}

func (d *defaultConnection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	panic("not provided by actual connection")
}

func (d *defaultConnection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	panic("not provided by actual connection")
}

func (d *defaultConnection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	panic("not provided by actual connection")
}

func (d *defaultConnection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := d.GetMessageCodec().(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (d *defaultConnection) GetPlcTagHandler() spi.PlcTagHandler {
	return d.tagHandler
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
