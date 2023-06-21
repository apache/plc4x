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

package options

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

// WithOption is a marker interface for options supplied by the builders like WithDefaultTtl
type WithOption interface {
	isOption() bool
}

// Option is a marker struct which can be used to mark an option
type Option struct {
}

func (_ Option) isOption() bool {
	return true
}

// WithCustomLogger is a global option to supply a custom logger
func WithCustomLogger(logger zerolog.Logger) WithOption {
	return withCustomLogger{logger: logger}
}

// ExtractCustomLogger can be used to extract the custom logger
func ExtractCustomLogger(options ...WithOption) (customLogger zerolog.Logger, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withCustomLogger:
			customLogger, found = option.logger, true
		}
	}
	return
}

// WithPassLoggerToModel enables passing of log to the model
func WithPassLoggerToModel(passLogger bool) WithOption {
	return withPassLoggerToModel{passLogger: passLogger}
}

// ExtractPassLoggerToModel to extract the flag indicating that model should be passed to Model
func ExtractPassLoggerToModel(options ...WithOption) (passLogger bool, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withPassLoggerToModel:
			passLogger, found = option.passLogger, true
		}
	}
	return
}

// WithReceiveTimeout set's a timeout for a receive-operation (similar to SO_RCVTIMEO)
func WithReceiveTimeout(timeout time.Duration) WithOption {
	return withReceiveTimeout{timeout: timeout}
}

// ExtractReceiveTimeout to extract the receive-timeout for reading operations. Defaults to 10 seconds
func ExtractReceiveTimeout(options ...WithOption) (receiveDuration time.Duration, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withReceiveTimeout:
			receiveDuration, found = option.timeout, true
		}
	}
	return
}

// WithTraceTransactionManagerWorkers enables trace transaction manager workers
func WithTraceTransactionManagerWorkers(traceWorkers bool) WithOption {
	return withTraceTransactionManagerWorkers{traceWorkers: traceWorkers}
}

// ExtractTransactionManagerWorkers to extract the flag indicating to trace transaction manager workers
func ExtractTransactionManagerWorkers(options ...WithOption) (traceWorkers bool, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withTraceTransactionManagerWorkers:
			traceWorkers, found = option.traceWorkers, true
		}
	}
	return
}

// WithTraceTransactionManagerTransactions enables trace transaction manager transactions
func WithTraceTransactionManagerTransactions(traceTransactions bool) WithOption {
	return withTraceTransactionManagerTransactions{traceTransactions: traceTransactions}
}

// ExtractTraceTransactionManagerTransactions to extract the flag indicating to trace transaction manager transactions
func ExtractTraceTransactionManagerTransactions(options ...WithOption) (traceTransactions bool, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withTraceTransactionManagerTransactions:
			traceTransactions, found = option.traceTransactions, true
		}
	}
	return
}

// WithTraceDefaultMessageCodecWorker enables trace default message codec worker
func WithTraceDefaultMessageCodecWorker(traceWorkers bool) WithOption {
	return withTraceDefaultMessageCodecWorker{traceWorkers: traceWorkers}
}

// ExtractTraceDefaultMessageCodecWorker to extract the flag indicating to trace default message codec workers
func ExtractTraceDefaultMessageCodecWorker(options ...WithOption) (traceWorkers bool, found bool) {
	for _, option := range options {
		switch option := option.(type) {
		case withTraceDefaultMessageCodecWorker:
			traceWorkers, found = option.traceWorkers, true
		}
	}
	return
}

// WithExecutorOptionTracerWorkers sets a flag which extends logging for workers
func WithExecutorOptionTracerWorkers(traceWorkers bool) WithOption {
	return &withTracerExecutorWorkersOption{traceWorkers: traceWorkers}
}

// ExtractTracerWorkers returns the value from WithExecutorOptionTracerWorkers
func ExtractTracerWorkers(_options ...WithOption) (traceWorkers bool, found bool) {
	for _, option := range _options {
		switch option := option.(type) {
		case *withTracerExecutorWorkersOption:
			traceWorkers, found = option.traceWorkers, true
		}
	}
	return
}

// GetLoggerContextForModel returns a log context if the WithPassLoggerToModel WithOption is set
func GetLoggerContextForModel(ctx context.Context, log zerolog.Logger, options ...WithOption) context.Context {
	passToModel := false
	for _, option := range options {
		switch option := option.(type) {
		case withPassLoggerToModel:
			passToModel = option.passLogger
		}
	}
	if passToModel {
		return log.WithContext(ctx)
	}
	return ctx
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withCustomLogger struct {
	Option
	logger zerolog.Logger
}

type withPassLoggerToModel struct {
	Option
	passLogger bool
}

type withReceiveTimeout struct {
	Option
	timeout time.Duration
}

type withTraceTransactionManagerWorkers struct {
	Option
	traceWorkers bool
}

type withTraceTransactionManagerTransactions struct {
	Option
	traceTransactions bool
}

type withTraceDefaultMessageCodecWorker struct {
	Option
	traceWorkers bool
}

type withTracerExecutorWorkersOption struct {
	Option
	traceWorkers bool
}

//
//
///////////////////////////////////////
///////////////////////////////////////
