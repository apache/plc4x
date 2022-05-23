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

package utils

import (
	"github.com/rs/zerolog"
	"io"
)

type TransportLogger struct {
	source io.ReadWriteCloser
	log    zerolog.Logger
}

type Option func(*TransportLogger)

func NewTransportLogger(source io.ReadWriteCloser, options ...Option) *TransportLogger {
	transportLogger := &TransportLogger{
		source: source,
	}
	for _, option := range options {
		option(transportLogger)
	}
	return transportLogger
}

func WithLogger(log zerolog.Logger) Option {
	return func(logger *TransportLogger) {
		logger.log = log
	}
}

func (t *TransportLogger) Read(p []byte) (int, error) {
	bytesRead, err := t.source.Read(p)
	if bytesRead > 0 {
		t.log.Printf("Read: %s", p[:bytesRead])
	}
	return bytesRead, err
}

func (t *TransportLogger) Write(p []byte) (int, error) {
	bytesWritten, err := t.source.Write(p)
	if bytesWritten > 0 {
		t.log.Printf("Write: %s", p[:bytesWritten])
	}
	return bytesWritten, err
}

func (t *TransportLogger) Close() error {
	return t.source.Close()
}
