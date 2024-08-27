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

package bacnetip

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type _MultiplexClient struct {
	*Client
	multiplexer *UDPMultiplexer
}

func _New_MultiplexClient(localLog zerolog.Logger, multiplexer *UDPMultiplexer) (*_MultiplexClient, error) {
	m := &_MultiplexClient{
		multiplexer: multiplexer,
	}
	var err error
	m.Client, err = NewClient(localLog, m)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	return m, nil
}

func (m *_MultiplexClient) Confirmation(args Args, kwargs KWArgs) error {
	return m.multiplexer.Confirmation(NewArgs(m, args), NoKWArgs)
}
