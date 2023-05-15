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
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

type Connection struct {
	device       *Device
	tagHandler   spi.PlcTagHandler
	valueHandler spi.PlcValueHandler
	options      map[string][]string
	connected    bool
	connectionId string
	tracer       *spi.Tracer
}

func NewConnection(device *Device, tagHandler spi.PlcTagHandler, valueHandler spi.PlcValueHandler, options map[string][]string) *Connection {
	connection := &Connection{
		device:       device,
		tagHandler:   tagHandler,
		valueHandler: valueHandler,
		options:      options,
		connected:    false,
		connectionId: utils.GenerateId(4),
	}
	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
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

func (c *Connection) GetTracer() *spi.Tracer {
	return c.tracer
}

func (c *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	return c.ConnectWithContext(context.Background())
}

func (c *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- _default.NewDefaultPlcConnectionCloseResult(nil, errors.Errorf("panic-ed %v", err))
			}
		}()
		// Check if the connection was already connected
		if c.connected {
			if c.tracer != nil {
				c.tracer.AddTrace("connect", "error: already connected")
			}
			// Return an error to the user.
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, errors.New("already connected"))
			return
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
				ch <- _default.NewDefaultPlcConnectionConnectResult(c, errors.New(errorString[0]))
			}
			if c.tracer != nil {
				c.tracer.AddTransactionalTrace(txId, "connect", "error: "+errorString[0])
			}
		} else {
			// Mark the connection as "connected"
			c.connected = true
			if c.tracer != nil {
				c.tracer.AddTransactionalTrace(txId, "connect", "success")
			}
			// Return the connection in a connected state to the user.
			ch <- _default.NewDefaultPlcConnectionConnectResult(c, nil)
		}
	}()
	return ch
}

func (c *Connection) BlockingClose() {
	<-c.Close()
}

func (c *Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	ch := make(chan plc4go.PlcConnectionCloseResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("panic-ed %v", err))
			}
		}()
		// Check if the connection is connected.
		if !c.connected {
			if c.tracer != nil {
				c.tracer.AddTrace("close", "error: not connected")
			}
			// Return an error to the user.
			ch <- _default.NewDefaultPlcConnectionCloseResult(c, errors.New("not connected"))
			return
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
		// Return a new connection to the user.
		ch <- _default.NewDefaultPlcConnectionCloseResult(c, nil)
	}()
	return ch
}

func (c *Connection) IsConnected() bool {
	return c.connected
}

func (c *Connection) Ping() <-chan plc4go.PlcConnectionPingResult {
	ch := make(chan plc4go.PlcConnectionPingResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ch <- _default.NewDefaultPlcConnectionPingResult(errors.Errorf("panic-ed %v", err))
			}
		}()
		// Check if the connection is connected
		if !c.connected {
			if c.tracer != nil {
				c.tracer.AddTrace("ping", "error: not connected")
			}
			// Return an error to the user.
			ch <- _default.NewDefaultPlcConnectionPingResult(errors.New("not connected"))
			return
		}
		var txId string
		if c.tracer != nil {
			txId = c.tracer.AddTransactionalStartTrace("ping", "started")
		}
		if delayString, ok := c.options["pingDelay"]; ok {
			// This is the length of the array, not the string
			if len(delayString) == 1 {
				delay, err := strconv.Atoi(delayString[0])
				if err == nil {
					time.Sleep(time.Duration(delay) * time.Millisecond)
				}
			}
		}
		if errorString, ok := c.options["pingError"]; ok {
			// If the ping operation should fail with an error, do so.
			if len(errorString) == 1 {
				ch <- _default.NewDefaultPlcConnectionPingResult(errors.New(errorString[0]))
			}
			if c.tracer != nil {
				c.tracer.AddTransactionalTrace(txId, "ping", "error: "+errorString[0])
			}
		} else {
			// Otherwise, give a positive response.
			if c.tracer != nil {
				c.tracer.AddTransactionalTrace(txId, "ping", "success")
			}
			ch <- _default.NewDefaultPlcConnectionPingResult(nil)
		}
	}()
	return ch
}

func (c *Connection) GetMetadata() model.PlcConnectionMetadata {
	return _default.DefaultConnectionMetadata{
		ConnectionAttributes: map[string]string{
			"connectionDelay": "Delay applied when connecting",
			"closingDelay":    "Delay applied when closing the connection",
			"pingDelay":       "Delay applied when executing a ping operation",
			"readDelay":       "Delay applied when executing a read operation",
			"writeDelay":      "Delay applied when executing a write operation",
		},
		ProvidesReading:     true,
		ProvidesWriting:     true,
		ProvidesSubscribing: false,
		ProvidesBrowsing:    false,
	}
}

func (c *Connection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	return spiModel.NewDefaultPlcReadRequestBuilder(c.tagHandler, NewReader(c.device, c.options, c.tracer))
}

func (c *Connection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(c.tagHandler, c.valueHandler, NewWriter(c.device, c.options, c.tracer))
}

func (c *Connection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(c.tagHandler, c.valueHandler, NewSubscriber(c.device, c.options, c.tracer))
}

func (c *Connection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	panic("not provided by simulated connection")
}

func (c *Connection) BrowseRequestBuilder() model.PlcBrowseRequestBuilder {
	panic("not provided by simulated connection")
}
