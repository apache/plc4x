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

package simulated

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Connection struct {
	device       *Device
	tagHandler   spi.PlcTagHandler
	valueHandler spi.PlcValueHandler
	options      map[string][]string
	connected    bool
	connectionId string
	tracer       tracer.Tracer
	invalidated  atomic.Bool

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewConnection(device *Device, tagHandler spi.PlcTagHandler, valueHandler spi.PlcValueHandler, connectionOptions map[string][]string, _options ...options.WithOption) *Connection {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	connectionId := utils.GenerateId(4)
	connection := &Connection{
		device:       device,
		tagHandler:   tagHandler,
		valueHandler: valueHandler,
		options:      connectionOptions,
		connected:    false,
		connectionId: connectionId,

		log: customLogger.With().Str("connectionId", connectionId).Logger(),
	}
	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	return connection
}

func (c *Connection) GetConnectionId() string {
	return c.connectionId
}

func (c *Connection) IsTraceEnabled() bool {
	return c.tracer != nil
}

func (c *Connection) GetTracer() tracer.Tracer {
	return c.tracer
}

func (c *Connection) Connect(_ context.Context) error {
	// Check if the connection was already connected
	if c.connected {
		if c.tracer != nil {
			c.tracer.AddTrace("connect", "error: already connected")
		}
		// Return an error to the user.
		return errors.New("already connected")
	}
	var txId string
	if c.tracer != nil {
		txId = c.tracer.AddTransactionalStartTrace("connect", "started")
	}
	if delayString, ok := c.options["connectionDelay"]; ok {
		// This is the length of the array, not the string
		if len(delayString) == 1 {
			delay, err := strconv.Atoi(delayString[0])
			if err == nil {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}
	// If we want the connection to fail, do so, otherwise return the connection.
	if errorString, ok := c.options["connectionError"]; ok {
		// If the ping operation should fail with an error, do so.
		if len(errorString) == 1 {
			return errors.New(errorString[0])
		}
		if c.tracer != nil {
			c.tracer.AddTransactionalTrace(txId, "connect", "error: "+errorString[0])
		}
	} else {
		// Mark the connection as "connected"
		c.connected = true
		c.invalidated.Store(false)
		if c.tracer != nil {
			c.tracer.AddTransactionalTrace(txId, "connect", "success")
		}
	}
	return nil
}

func (c *Connection) Close() error {
	ctx := context.TODO()
	ctx, cancelFunc := utils.WithNamedTimeout(ctx, "connection close timeout", 5*time.Second)
	defer cancelFunc()

	// Check if the connection is connected.
	if !c.connected {
		if c.invalidated.Load() {
			return nil
		}
		if c.tracer != nil {
			c.tracer.AddTrace("close", "error: not connected")
		}
		// Return an error to the user.
		return errors.New("not connected")
	}
	var txId string
	if c.tracer != nil {
		txId = c.tracer.AddTransactionalStartTrace("close", "started")
	}
	// If a delay was configured, wait for the pre-configured time.
	if delayString, ok := c.options["closingDelay"]; ok {
		// This is the length of the array, not the string
		if len(delayString) == 1 {
			delay, err := strconv.Atoi(delayString[0])
			if err == nil {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
	}
	// Mark the connection as "disconnected".
	c.connected = false
	if c.tracer != nil {
		c.tracer.AddTransactionalTrace(txId, "close", "success")
	}
	return nil
}

func (c *Connection) IsConnected() bool {
	return c.connected && !c.IsInvalidated()
}

func (c *Connection) Ping(ctx context.Context) error {
	if c.IsInvalidated() {
		return errors.New("connection has been invalidated")
	}
	// Check if the connection is connected
	if !c.connected {
		if c.tracer != nil {
			c.tracer.AddTrace("ping", "error: not connected")
		}
		return errors.New("not connected")
	}
	var txId string
	if c.tracer != nil {
		txId = c.tracer.AddTransactionalStartTrace("ping", "started")
	}
	if delayString, ok := c.options["pingDelay"]; ok {
		// This is the length of the array, not the string
		if len(delayString) == 1 {
			delay, err := strconv.Atoi(delayString[0])
			if err != nil {
				return errors.Wrapf(err, "invalid delay '%s'", delayString[0])
			}
			timer := time.NewTimer(time.Duration(delay) * time.Millisecond)
			c.log.Info().Msgf("Ping delay of %d ms", delay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	if errorString, ok := c.options["pingError"]; ok {
		// If the ping operation should fail with an error, do so.
		if c.tracer != nil {
			c.tracer.AddTransactionalTrace(txId, "ping", "error: "+errorString[0])
		}
		return errors.New(errorString[0])
	} else {
		// Otherwise, give a positive response.
		if c.tracer != nil {
			c.tracer.AddTransactionalTrace(txId, "ping", "success")
		}
		return nil
	}
}

func (c *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return &_default.DefaultConnectionMetadata{
		ConnectionAttributes: map[string]string{
			"connectionDelay": "Delay applied when connecting",
			"closingDelay":    "Delay applied when closing the connection",
			"pingDelay":       "Delay applied when executing a ping operation",
			"readDelay":       "Delay applied when executing a read operation",
			"writeDelay":      "Delay applied when executing a write operation",
		},
		ProvidesReading: true,
		ProvidesWriting: true,
		// SubscriptionRequestBuilder below does hand out a builder, but the Subscriber behind it
		// is a stub whose Subscribe/Unsubscribe always answer "Not Implemented", so advertising
		// the capability would be a lie. Flip this to true only together with a Subscriber that
		// actually delivers events. Browsing has no implementation at all.
		ProvidesSubscribing: false,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(c.tagHandler, NewReader(c.device, c.options, c.tracer))
}

func (c *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.tagHandler, c.valueHandler, NewWriter(c.device, c.options, c.tracer))
}

func (c *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(c.tagHandler, c.valueHandler, NewSubscriber(c.device, c.options, c.tracer))
}

func (c *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	panic("not provided by simulated connection")
}

func (c *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	panic("not provided by simulated connection")
}

func (c *Connection) String() string {
	return "simulatedConnection"
}

func (c *Connection) Invalidate() {
	if c.invalidated.Swap(true) {
		return
	}
	c.log.Debug().Msg("invalidating connection")
	if err := c.Close(); err != nil {
		c.log.Warn().Err(err).Msg("error closing invalidated connection")
	}
}

func (c *Connection) IsInvalidated() bool {
	return c.invalidated.Load()
}
