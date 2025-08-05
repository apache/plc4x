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
)

// Bind a list of clients and servers together, top down
func Bind(localLog zerolog.Logger, args ...any) error {
	if _debug != nil {
		_debug("bind %r", Args(args))
	}
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
		if _debug != nil {
			_debug("    - client: %r", left)
		}
		leftStringer, _ := left.(fmt.Stringer)
		localLog.Debug().Stringer("left", leftStringer).Type("leftType", left).Msg("left pair element")
		right := args[i+1]
		if _debug != nil {
			_debug("    - server: %r", right)
		}
		rightStringer, _ := right.(fmt.Stringer)
		localLog.Debug().Stringer("right", rightStringer).Type("rightType", right).Msg("right pair element")

		// make sure we're binding clients and servers
		clientCast, okClient := left.(Client)
		serverCast, okServer := right.(Server)
		elementServiceCast, okElementService := left.(ApplicationServiceElement)
		serviceAccessPointCast, okServiceAccessPoint := right.(ServiceAccessPoint)
		if okClient && okServer {
			localLog.Trace().Msg("linking client-server")
			clientCast._setClientPeer(serverCast)
			serverCast._setServerPeer(clientCast)
		} else if okElementService && okServiceAccessPoint { // we could be binding application clients and servers
			localLog.Trace().Msg("linking service-elements")
			elementServiceCast._setElementService(serviceAccessPointCast)
			serviceAccessPointCast._setServiceElement(elementServiceCast)
		} else {
			localLog.Debug().
				Bool("okClient", okClient).
				Bool("okServer", okServer).
				Bool("okElementService", okElementService).
				Bool("okServiceAccessPoint", okServiceAccessPoint).
				Msg("cast states")
			return errors.New("Bind() requires a client and a server")
		}
		if _debug != nil {
			_debug("    - bound")
		}
		localLog.Trace().Msg("bound")
	}
	return nil
}
