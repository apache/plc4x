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
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi/tracer"
)

type DefaultPlcConnectionCloseResult interface {
	plc4go.PlcConnectionCloseResult
	GetTraces() []tracer.TraceEntry
}

func NewDefaultPlcConnectionCloseResult(connection plc4go.PlcConnection, err error) plc4go.PlcConnectionCloseResult {
	return &defaultPlcConnectionCloseResult{
		connection: connection,
		err:        err,
		traces:     nil,
	}
}

func NewDefaultPlcConnectionCloseResultWithTraces(connection plc4go.PlcConnection, err error, traces []tracer.TraceEntry) plc4go.PlcConnectionCloseResult {
	return &defaultPlcConnectionCloseResult{
		connection: connection,
		err:        err,
		traces:     traces,
	}
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=defaultPlcConnectionCloseResult
type defaultPlcConnectionCloseResult struct {
	connection plc4go.PlcConnection
	err        error
	traces     []tracer.TraceEntry
}

func (d *defaultPlcConnectionCloseResult) GetConnection() plc4go.PlcConnection {
	return d.connection
}

func (d *defaultPlcConnectionCloseResult) GetErr() error {
	return d.err
}

func (d *defaultPlcConnectionCloseResult) GetTraces() []tracer.TraceEntry {
	return d.traces
}
