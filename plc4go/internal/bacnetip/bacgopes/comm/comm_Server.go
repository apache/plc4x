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

package comm

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Server interface {
	ServerContract
	ServerRequirements
}

// ServerContract provides a set of functions which can be overwritten by a sub struct
type ServerContract interface {
	fmt.Stringer
	utils.Serializable
	Response(args Args, kwArgs KWArgs) error
	_setServerPeer(serverPeer Client)
	HasServerPeer() bool
	GetServerId() *int
}

// ServerRequirements provides a set of functions which must be overwritten by a sub struct
type ServerRequirements interface {
	Indication(args Args, kwArgs KWArgs) error
}

// Server is an "abstract" struct which is used in another struct as delegate
//
//go:generate go tool plc4xGenerator -type=server -prefix=comm_
type server struct {
	serverID   *int
	serverPeer Client `asPtr:"true"`

	// args
	argServerRequirements ServerRequirements `ignore:"true"`

	log zerolog.Logger
}

func NewServer(localLog zerolog.Logger, options ...Option) (ServerContract, error) {
	s := &server{
		log: localLog,
	}
	ApplyAppliers(options, s)
	if _debug != nil {
		_debug("__init__ sid=%v", s.serverID)
	}
	if s.serverID != nil {
		sid := *s.serverID
		if _, ok := serverMap[sid]; ok {
			return nil, errors.Errorf("already a server %d", sid)
		}
		serverMap[sid] = s

		// automatically bind
		if c, ok := clientMap[sid]; ok {
			if c.clientPeer != nil {
				return nil, errors.Errorf("client %d already bound", sid)
			}

			// Note: we need to pass the requirements (which should contain s as delegate) here
			if err := Bind(localLog, c, s.argServerRequirements); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func WithServerSID(sid int, requirements ServerRequirements) GenericApplier[*server] {
	return WrapGenericApplier(func(s *server) {
		s.serverID = &sid
		s.argServerRequirements = requirements
	})
}

func (s *server) Response(args Args, kwArgs KWArgs) error {
	if _debug != nil {
		_debug("response %r %r", args, kwArgs)
	}
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Response")

	if s.serverPeer == nil {
		return errors.New("unbound server")
	}
	return s.serverPeer.Confirmation(args, kwArgs)
}

func (s *server) _setServerPeer(serverPeer Client) {
	s.serverPeer = serverPeer
}

func (s *server) HasServerPeer() bool {
	return s.serverPeer != nil
}

func (s *server) GetServerId() *int {
	return s.serverID
}
