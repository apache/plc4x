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

package bacgopes

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:generate go run ../../../tools/plc4xgenerator/gen.go -type=SieveQueue -prefix=iocb_
type SieveQueue struct {
	*IOQController
	requestFn func(apdu PDU)
	address   *Address

	log zerolog.Logger `ignore:"true"`
}

func NewSieveQueue(localLog zerolog.Logger, fn func(apdu PDU), address *Address) (*SieveQueue, error) {
	s := &SieveQueue{}
	var err error
	s.IOQController, err = NewIOQController(localLog, address.String(), s)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a IOQController")
	}

	// Save a reference to the request function
	s.requestFn = fn
	s.address = address
	return s, nil
}

func (s *SieveQueue) ProcessIO(iocb _IOCB) error {
	s.log.Debug().Stringer("iocb", iocb).Msg("ProcessIO")

	// this is now an active request
	if err := s.ActiveIO(iocb); err != nil {
		return errors.Wrap(err, "error on active io")
	}

	// send the request
	s.requestFn(iocb.getRequest())
	return nil
}
