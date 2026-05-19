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
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// DefaultConnectionRequirements defines the required at a implementing connection when using DefaultConnection
// additional options can be set using the functions returning WithOption (e.g. WithDefaultTtl, WithPlcTagHandler...)
type DefaultConnectionRequirements interface {
	// IsConnected should return the implementing connection check when using DefaultConnection
	IsConnected() bool
	// GetMessageCodec should return the spi.MessageCodec in use
	GetMessageCodec() spi.MessageCodec
}

// DefaultConnection should be used as an embedded struct. All defined methods here have default implementations
type DefaultConnection interface {
	utils.Serializable
	plc4go.PlcConnection
	spi.TransportInstanceExposer
	spi.HandlerExposer
	SetConnected(connected bool)
}

// NewDefaultConnection is the factory for a DefaultConnection
func NewDefaultConnection(requirements DefaultConnectionRequirements, options ...options.WithOption) DefaultConnection {
	return buildDefaultConnection(requirements, options...)
}

func WithPlcTagHandler(tagHandler spi.PlcTagHandler) options.WithOption {
	return withPlcTagHandler{plcTagHandler: tagHandler}
}

func WithPlcValueHandler(plcValueHandler spi.PlcValueHandler) options.WithOption {
	return withPlcValueHandler{plcValueHandler: plcValueHandler}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withPlcTagHandler struct {
	options.Option
	plcTagHandler spi.PlcTagHandler
}

type withPlcValueHandler struct {
	options.Option
	plcValueHandler spi.PlcValueHandler
}

//go:generate go tool plc4xGenerator -type=defaultConnection
type defaultConnection struct {
	DefaultConnectionRequirements `ignore:"true"`
	// connected indicates if a connection is connected
	connected    atomic.Bool
	tagHandler   spi.PlcTagHandler
	valueHandler spi.PlcValueHandler

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func buildDefaultConnection(requirements DefaultConnectionRequirements, _options ...options.WithOption) DefaultConnection {
	var tagHandler spi.PlcTagHandler
	var valueHandler spi.PlcValueHandler

	for _, option := range _options {
		switch option.(type) {
		case withPlcTagHandler:
			tagHandler = option.(withPlcTagHandler).plcTagHandler
		case withPlcValueHandler:
			valueHandler = option.(withPlcValueHandler).plcValueHandler
		}
	}

	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &defaultConnection{
		DefaultConnectionRequirements: requirements,
		tagHandler:                    tagHandler,
		valueHandler:                  valueHandler,

		log: customLogger,
	}
}

func (d *defaultConnection) SetConnected(connected bool) {
	d.log.Trace().Bool("connected", connected).Msg("set connected")
	d.connected.Store(connected)
}

func (d *defaultConnection) Connect(ctx context.Context) error {
	d.log.Trace().Msg("Connecting")
	err := d.GetMessageCodec().Connect(ctx)
	d.SetConnected(true)
	return err
}

func (d *defaultConnection) Close() error {
	d.log.Trace().Msg("close connection")
	if messageCodec := d.GetMessageCodec(); messageCodec != nil {
		d.log.Trace().Msg("disconnecting message codec")
		if err := messageCodec.Disconnect(); err != nil {
			d.log.Warn().Err(err).Msg("Error disconnecting message code")
		} else {
			d.log.Trace().Msg("message codec disconnected")
		}
	}
	var err error
	if transportInstance := d.GetTransportInstance(); transportInstance != nil {
		d.log.Trace().Msg("closing transport instance")
		if err = transportInstance.Close(); err != nil {
			d.log.Warn().Err(err).Msg("Error disconnecting transport instance")
		} else {
			d.log.Trace().Msg("transport instance closed")
		}
	}
	d.SetConnected(false)
	return err
}

func (d *defaultConnection) IsConnected() bool {
	// TODO: should we check here if the transport is connected?
	return d.connected.Load()
}

func (d *defaultConnection) Ping(_ context.Context) error {
	if !d.DefaultConnectionRequirements.IsConnected() {
		return errors.New("not connected")
	}
	return nil
}

func (d *defaultConnection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &DefaultConnectionMetadata{
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

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
