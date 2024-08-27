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

// maps of named clients and servers
var clientMap map[int]*Client

var serverMap map[int]*Server

// maps of named SAPs and ASEs
var serviceMap map[int]*ServiceAccessPoint

var elementMap map[int]*ApplicationServiceElement

func init() {
	clientMap = make(map[int]*Client)
	serverMap = make(map[int]*Server)
	serviceMap = make(map[int]*ServiceAccessPoint)
	elementMap = make(map[int]*ApplicationServiceElement)
}

// Bind a list of clients and servers together, top down
func Bind(localLog zerolog.Logger, args ...any) error {
	// generic bind is pairs of names
	if len(args) == 0 {
		// find unbound clients and bind them
		for cid, client := range clientMap {
			// skip those that are already bound
			if client.clientPeer != nil {
				continue
			}

			server, ok := serverMap[cid]
			if !ok {
				return errors.Errorf("unmatched server %d", cid)
			}

			if server.serverPeer != nil {
				return errors.Errorf("server already bound %d", cid)
			}

			if err := Bind(localLog, client, server); err != nil {
				return errors.Wrap(err, "error binding")
			}
		}

		// see if there are any unbound servers
		for sid, server := range serverMap {
			if server.serverPeer != nil {
				continue
			}

			if _, ok := clientMap[sid]; !ok {
				return errors.Errorf("unmatched client %d", sid)
			} else {
				return errors.Errorf("unknown unbound server %d", sid)
			}
		}

		// find unbound application service elements and bind them
		for eid, element := range elementMap {
			// skip those that are already bound
			if element.elementService != nil {
				continue
			}

			service, ok := serviceMap[eid]
			if !ok {
				return errors.Errorf("unmatched element %d", eid)
			}

			if service.serviceElement == nil {
				return errors.Errorf("element already bound %d", eid)
			}

			if err := Bind(localLog, element, service); err != nil {
				return errors.Wrap(err, "error binding")
			}
		}

		// see if there are any unbound services
		for sid, service := range serviceMap {
			if service.serviceElement != nil {
				continue
			}

			if _, ok := elementMap[sid]; !ok {
				return errors.Errorf("unmatched service %d", sid)
			} else {
				return errors.Errorf("unknown unbound service %d", sid)
			}
		}
	}

	// go through the argument pairs
	for i := 0; i < len(args)-1; i++ {
		left := args[i]
		leftStringer, _ := left.(fmt.Stringer)
		localLog.Debug().Stringer("left", leftStringer).Msg("left pair element")
		right := args[i+1]
		rightStringer, _ := right.(fmt.Stringer)
		localLog.Debug().Stringer("right", rightStringer).Msg("right pair element")

		// make sure we're binding clients and servers
		clientCast, okClient := left.(_Client)
		serverCast, okServer := right.(_Server)
		elementServiceCast, okElementService := left.(ApplicationServiceElementContract)
		serviceAccessPointCast, okServiceAccessPoint := right.(ServiceAccessPointContract)
		if okClient && okServer {
			localLog.Trace().Msg("linking client-server")
			clientCast._setClientPeer(serverCast)
			serverCast._setServerPeer(clientCast)
		} else if okElementService && okServiceAccessPoint { // we could be binding application clients and servers
			localLog.Trace().Msg("linking service-elements")
			elementServiceCast._setElementService(serviceAccessPointCast)
			serviceAccessPointCast._setServiceElement(elementServiceCast)
		} else {
			return errors.New("Bind() requires a client and a server")
		}
	}
	localLog.Debug().Msg("bound")
	return nil
}
