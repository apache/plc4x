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

package bvllservice

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

//go:generate go tool plc4xGenerator -type=_MultiplexServer -prefix=bvllservice_
type _MultiplexServer struct {
	ServerContract
	multiplexer *UDPMultiplexer `stringer:"true"`
}

func _New_MultiplexServer(localLog zerolog.Logger, multiplexer *UDPMultiplexer) (*_MultiplexServer, error) {
	m := &_MultiplexServer{
		multiplexer: multiplexer,
	}
	var err error
	m.ServerContract, err = NewServer(localLog) // TODO: do we need to pass sid?
	if err != nil {
		return nil, errors.Wrap(err, "error creating server")
	}
	return m, nil
}

func (m *_MultiplexServer) Indication(args Args, kwArgs KWArgs) error {
	return m.multiplexer.Indication(NA(m, args), NoKWArgs())
}
