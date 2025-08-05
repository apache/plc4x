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

type Client interface {
	ClientContract
	ClientRequirements
}

// ClientContract provides a set of functions which can be overwritten by a sub struct
type ClientContract interface {
	fmt.Stringer
	utils.Serializable
	GetClientID() *int
	Request(args Args, kwArgs KWArgs) error
	_setClientPeer(server Server)
}

// ClientRequirements provides a set of functions which must be overwritten by a sub struct
type ClientRequirements interface {
	Confirmation(args Args, kwArgs KWArgs) error
}

//go:generate go tool plc4xGenerator -type=client -prefix=comm_
type client struct {
	clientID   *int
	clientPeer Server `asPtr:"true"`

	// args
	argClientRequirements ClientRequirements `ignore:"true"`

	log zerolog.Logger
}

var _ ClientContract = (*client)(nil)

func NewClient(localLog zerolog.Logger, options ...Option) (ClientContract, error) {
	c := &client{
		log: localLog,
	}
	ApplyAppliers(options, c)
	if _debug != nil {
		_debug("__init__ cid=%v", c.clientID)
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
			if err := Bind(localLog, c.argClientRequirements, server); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return c, nil
}

func WithClientCID(cid int, requirements ClientRequirements) GenericApplier[*client] {
	return WrapGenericApplier(func(c *client) {
		c.clientID = &cid
		c.argClientRequirements = requirements
	})
}

func (c *client) Request(args Args, kwArgs KWArgs) error {
	if _debug != nil {
		_debug("request %r %r", args, kwArgs)
	}
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Request")

	if c.clientPeer == nil {
		return errors.Errorf("unbound client: %s", c)
	}
	return c.clientPeer.Indication(args, kwArgs)
}

func (c *client) _setClientPeer(server Server) {
	c.clientPeer = server
}

func (c *client) GetClientID() *int {
	return c.clientID
}
