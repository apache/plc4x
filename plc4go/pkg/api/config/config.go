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

package config

import (
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
)

// TraceTransactionManagerWorkers when set to true the transaction manager displays worker states in log
var (
	TraceTransactionManagerWorkers      bool
	TraceTransactionManagerTransactions bool
	TraceDefaultMessageCodecWorker      bool
)

// TraceConnectionCache when set to true the connection cache outputs logs by default
var (
	TraceConnectionCache bool
)

// WithCustomLogger is a global option to supply a custom logger
func WithCustomLogger(logger zerolog.Logger) WithOption {
	return options.WithCustomLogger(logger)
}

// WithPassLoggerToModel enables passing of log to the model
func WithPassLoggerToModel(passLogger bool) WithOption {
	return options.WithPassLoggerToModel(passLogger)
}

// WithOption is a marker interface for options
type WithOption interface {
	options.WithOption
}
