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

// _Client is an interface used for documentation
type _Client interface {
	fmt.Stringer
	Request(args Args, kwargs KWArgs) error
	Confirmation(args Args, kwargs KWArgs) error
	_setClientPeer(server _Server)
	getClientId() *int
}

// Client is an "abstract" struct which is used in another struct as delegate
type Client struct {
	clientID   *int
	clientPeer _Server

	log zerolog.Logger
}

func NewClient(localLog zerolog.Logger, rootStruct _Client, opts ...func(*Client)) (*Client, error) {
	c := &Client{
		log: localLog,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.clientID != nil {
		cid := *c.clientID
		if _, ok := clientMap[cid]; ok {
			return nil, errors.Errorf("already a client %d", cid)
		}
		clientMap[cid] = c

		// automatically bind
		if server, ok := serverMap[cid]; ok {
			if server.serverPeer != nil {
				return nil, errors.Errorf("server %d already bound", cid)
			}

			// Note: we need to pass the requirements (which should contain c as delegate) here
			if err := Bind(localLog, rootStruct, server); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return c, nil
}

func WithClientCid(cid int) func(*Client) {
	return func(c *Client) {
		c.clientID = &cid
	}
}

func (c *Client) Request(args Args, kwargs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")

	if c.clientPeer == nil {
		return errors.Errorf("unbound client: %s", c)
	}
	return c.clientPeer.Indication(args, kwargs)
}

func (c *Client) Confirmation(args Args, kwargs KWArgs) error {
	panic("this should be implemented by outer struct")
}

func (c *Client) _setClientPeer(server _Server) {
	c.clientPeer = server
}

func (c *Client) getClientId() *int {
	return c.clientID
}

func (c *Client) String() string {
	clientPeer := ""
	if c.clientPeer != nil {
		clientPeer = fmt.Sprintf(", clientPeerId: %d", c.clientPeer.getServerId())
	}
	return fmt.Sprintf("Client(cid:%d%s)", c.clientID, clientPeer)
}
