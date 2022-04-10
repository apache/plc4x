/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package spi

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"time"
)

type TraceEntry struct {
	Timestamp     time.Time
	ConnectionId  string
	TransactionId string
	Operation     string
	Message       string
}

type TracerProvider interface {
	EnableTracer()
	GetTracer() *Tracer
}

type Tracer struct {
	connectionId string
	traceEntries []TraceEntry
}

func NewTracer(connectionId string) *Tracer {
	return &Tracer{
		connectionId: connectionId,
		traceEntries: []TraceEntry{},
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

func (t Tracer) GetTraces() []TraceEntry {
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
	transactionId := utils.GenerateId(4)
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
	for _, trace := range traces {
		matches := true
		if connectionIdFilter != "" {
			if trace.ConnectionId != connectionIdFilter {
				matches = false
			}
		}
		if matches && transactionIdFilter != "" {
			if trace.TransactionId != transactionIdFilter {
				matches = false
			}
		}
		if matches && operationFilter != "" {
			if trace.Operation != operationFilter {
				matches = false
			}
		}
		if matches && messageFilter != "" {
			if trace.Message != messageFilter {
				matches = false
			}
		}
		if matches {
			result = append(result, trace)
		}
	}
	return result
}
