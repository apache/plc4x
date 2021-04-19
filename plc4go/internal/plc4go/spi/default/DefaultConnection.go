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

package _default

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
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
func NewDefaultConnection(requirements DefaultConnectionRequirements, options ...WithOption) DefaultConnection {
	return buildDefaultConnection(requirements, options...)
}

// WithDefaultTtl ttl is time.Second * 10 by default
func WithDefaultTtl(defaultTtl time.Duration) WithOption {
	return withDefaultTtl{defaultTtl: defaultTtl}
}

func WithPlcFieldHandler(plcFieldHandler spi.PlcFieldHandler) WithOption {
	return withPlcFieldHandler{plcFieldHandler: plcFieldHandler}
}

func WithPlcValueHandler(plcValueHandler spi.PlcValueHandler) WithOption {
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

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withDefaultTtl struct {
	option
	// defaultTtl the time to live after a close
	defaultTtl time.Duration
}

type withPlcFieldHandler struct {
	option
	plcFieldHandler spi.PlcFieldHandler
}

type withPlcValueHandler struct {
	option
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

func buildDefaultConnection(requirements DefaultConnectionRequirements, options ...WithOption) DefaultConnection {
	defaultTtl := time.Second * 10
	var fieldHandler spi.PlcFieldHandler
	var valueHandler spi.PlcValueHandler

	for _, option := range options {
		if !option.isOption() {
			panic("not a option")
		}
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
		ch <- plc4go.NewPlcConnectionConnectResult(connection, err)
	}()
	return ch
}

func (d *defaultConnection) BlockingClose() {
	log.Trace().Msg("blocking close connection")
	closeResults := d.GetConnection().Close()
	d.SetConnected(false)
	select {
	case <-closeResults:
		return
	case <-time.After(d.GetTtl()):
		return
	}
}

func (d *defaultConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	log.Trace().Msg("close connection")
	d.SetConnected(false)
	ch := make(chan plc4go.PlcConnectionCloseResult)
	go func() {
		ch <- plc4go.NewPlcConnectionCloseResult(d.GetConnection(), nil)
	}()
	return ch
}

func (d *defaultConnection) IsConnected() bool {
	return d.connected
}

func (d *defaultConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	ch := make(chan plc4go.PlcConnectionPingResult)
	go func() {
		if d.GetConnection().IsConnected() {
			ch <- plc4go.NewPlcConnectionPingResult(nil)
		} else {
			ch <- plc4go.NewPlcConnectionPingResult(errors.New("not connected"))
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
