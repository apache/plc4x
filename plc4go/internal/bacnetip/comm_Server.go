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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// _Server is an interface used for documentation
type _Server interface {
	fmt.Stringer
	Indication(args Args, kwargs KWArgs) error
	Response(args Args, kwargs KWArgs) error
	_setServerPeer(serverPeer _Client)
	getServerId() *int
}

// Server is an "abstract" struct which is used in another struct as delegate
type Server struct {
	serverID   *int
	serverPeer _Client

	log zerolog.Logger
}

func NewServer(localLog zerolog.Logger, rootStruct _Server, opts ...func(server *Server)) (*Server, error) {
	s := &Server{
		log: localLog,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.serverID != nil {
		sid := *s.serverID
		if _, ok := serverMap[sid]; ok {
			return nil, errors.Errorf("already a server %d", sid)
		}
		serverMap[sid] = s

		// automatically bind
		if client, ok := clientMap[sid]; ok {
			if client.clientPeer != nil {
				return nil, errors.Errorf("client %d already bound", sid)
			}

			// Note: we need to pass the requirements (which should contain s as delegate) here
			if err := Bind(localLog, client, rootStruct); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func WithServerSID(sid int) func(*Server) {
	return func(server *Server) {
		server.serverID = &sid
	}
}

func (s *Server) Indication(Args, KWArgs) error {
	panic("this should be implemented by outer struct")
}

func (s *Server) Response(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")

	if s.serverPeer == nil {
		return errors.New("unbound server")
	}
	return s.serverPeer.Confirmation(args, kwargs)
}

func (s *Server) _setServerPeer(serverPeer _Client) {
	s.serverPeer = serverPeer
}

func (s *Server) getServerId() *int {
	return s.serverID
}

func (s *Server) String() string {
	serverPeer := ""
	if s.serverPeer != nil {
		serverPeer = fmt.Sprintf(", serverPeerId: %d", s.serverPeer.getClientId())
	}
	return fmt.Sprintf("Server(cid:%d%s)", s.serverID, serverPeer)
}
