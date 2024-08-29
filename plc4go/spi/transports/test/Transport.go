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

package test

import (
	"net/url"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

type Transport struct {
	preregisteredInstances map[url.URL]transports.TransportInstance

	log zerolog.Logger
}

func NewTransport(_options ...options.WithOption) *Transport {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Transport{
		preregisteredInstances: map[url.URL]transports.TransportInstance{},
		log:                    customLogger,
	}
}

func (m *Transport) GetTransportCode() string {
	return "test"
}

func (m *Transport) GetTransportName() string {
	return "Test Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	if _, ok := options["failTestTransport"]; ok {
		return nil, errors.New("test transport failed on purpose")
	}
	if preregisteredInstance, ok := m.preregisteredInstances[transportUrl]; ok {
		m.log.Trace().Stringer("transportUrl", &transportUrl).Msg("Returning pre registered instance")
		return preregisteredInstance, nil
	}
	m.log.Trace().Msg("create transport instance")
	return NewTransportInstance(m, _options...), nil
}

func (m *Transport) AddPreregisteredInstances(transportUrl url.URL, preregisteredInstance transports.TransportInstance) error {
	if _, ok := m.preregisteredInstances[transportUrl]; ok {
		return errors.Errorf("registered instance for %v already registered", transportUrl)
	}
	m.preregisteredInstances[transportUrl] = preregisteredInstance
	return nil
}

func (m *Transport) Close() error {
	m.log.Trace().Msg("Closing")
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
