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

package tracer

import (
	"fmt"
	"time"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/rs/zerolog"
)

type TraceEntry struct {
	Timestamp     time.Time
	ConnectionId  string
	TransactionId string
	Operation     string
	Message       string
}

type Provider interface {
	EnableTracer()
	GetTracer() *Tracer
}

type Tracer struct {
	connectionId string
	traceEntries []TraceEntry

	log zerolog.Logger
}

func NewTracer(connectionId string, _options ...options.WithOption) *Tracer {
	return &Tracer{
		connectionId: connectionId,
		traceEntries: []TraceEntry{},
		log:          options.ExtractCustomLogger(_options...),
	}
}

func (t *Tracer) GetConnectionId() string {
	return t.connectionId
}

func (t *Tracer) SetConnectionId(connectionId string) {
	t.connectionId = connectionId
}

func (t *Tracer) ResetTraces() {
	t.traceEntries = []TraceEntry{}
}

func (t *Tracer) GetTraces() []TraceEntry {
	return t.traceEntries
}

func (t *Tracer) AddTrace(operation string, message string) {
	t.traceEntries = append(t.traceEntries, TraceEntry{
		Timestamp:     time.Now(),
		ConnectionId:  t.connectionId,
		TransactionId: "",
		Operation:     operation,
		Message:       message,
	})
}

func (t *Tracer) AddTransactionalStartTrace(operation string, message string) string {
	transactionId := utils.GenerateId(t.log, 4)
	t.traceEntries = append(t.traceEntries, TraceEntry{
		Timestamp:     time.Now(),
		ConnectionId:  t.connectionId,
		TransactionId: transactionId,
		Operation:     operation,
		Message:       message,
	})
	return transactionId
}

func (t *Tracer) AddTransactionalTrace(transactionId string, operation string, message string) {
	t.traceEntries = append(t.traceEntries, TraceEntry{
		Timestamp:     time.Now(),
		ConnectionId:  t.connectionId,
		TransactionId: transactionId,
		Operation:     operation,
		Message:       message,
	})
}

func (t *Tracer) FilterTraces(traces []TraceEntry, connectionIdFilter string, transactionIdFilter string, operationFilter string, messageFilter string) []TraceEntry {
	var result []TraceEntry
traceFiltering:
	for _, trace := range traces {
		if connectionIdFilter != "" && trace.ConnectionId != connectionIdFilter {
			continue traceFiltering
		}
		if transactionIdFilter != "" && trace.TransactionId != transactionIdFilter {
			continue traceFiltering
		}
		if operationFilter != "" && trace.Operation != operationFilter {
			continue traceFiltering
		}
		if messageFilter != "" && trace.Message != messageFilter {
			continue traceFiltering
		}
		result = append(result, trace)
	}
	return result
}

func (t *Tracer) String() string {
	return fmt.Sprintf("Tracer for %s", t.connectionId)
}
