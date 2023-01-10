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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

type _PCI struct {
	pduUserData    spi.Message
	pduSource      *Address
	pduDestination *Address
}

func _New_PCI(pduUserData spi.Message, pduSource *Address, pduDestination *Address) *_PCI {
	return &_PCI{pduUserData, pduSource, pduDestination}
}

func (p *_PCI) String() string {
	return fmt.Sprintf("pduUserData:\n%s\n, pduSource: %s, pduDestination: %s", p.pduUserData, p.pduSource, p.pduDestination)
}

// _Client is an interface used for documentation
type _Client interface {
	Request(pdu _PDU) error
	Confirmation(pdu _PDU) error
	_setClientPeer(server _Server)
	getClientId() *int
}

// Client is an "abstract" struct which is used in another struct as delegate
type Client struct {
	clientID   *int
	clientPeer _Server
}

func NewClient(cid *int, rootStruct _Client) (*Client, error) {
	c := &Client{
		clientID: cid,
	}
	if cid != nil {
		if _, ok := clientMap[*cid]; ok {
			return nil, errors.Errorf("already a client %d", *cid)
		}
		clientMap[*cid] = c

		// automatically bind
		if server, ok := serverMap[*cid]; ok {
			if server.serverPeer != nil {
				return nil, errors.Errorf("server %d already bound", *cid)
			}

			// Note: we need to pass the rootStruct (which should contain c as delegate) here
			if err := bind(rootStruct, server); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return c, nil
}

func (c *Client) Request(pdu _PDU) error {
	log.Debug().Msgf("request\n%s", pdu)

	if c.clientPeer == nil {
		return errors.New("unbound client")
	}
	return c.clientPeer.Indication(pdu)
}

func (c *Client) Confirmation(_PDU) error {
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
		clientPeer = fmt.Sprintf(" clientPeerId: %d", c.clientPeer.getServerId())
	}
	return fmt.Sprintf("Client(cid:%d)%s", c.clientID, clientPeer)
}

// _Server is an interface used for documentation
type _Server interface {
	Indication(pdu _PDU) error
	Response(pdu _PDU) error
	_setServerPeer(serverPeer _Client)
	getServerId() *int
}

// Server is an "abstract" struct which is used in another struct as delegate
type Server struct {
	serverID   *int
	serverPeer _Client
}

func NewServer(sid *int, rootStruct _Server) (*Server, error) {
	s := &Server{
		serverID: sid,
	}
	if sid != nil {
		if _, ok := serverMap[*sid]; ok {
			return nil, errors.Errorf("already a server %d", *sid)
		}
		serverMap[*sid] = s

		// automatically bind
		if client, ok := clientMap[*sid]; ok {
			if client.clientPeer != nil {
				return nil, errors.Errorf("client %d already bound", *sid)
			}

			// Note: we need to pass the rootStruct (which should contain s as delegate) here
			if err := bind(client, rootStruct); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func (s *Server) Indication(_PDU) error {
	panic("this should be implemented by outer struct")
}

func (s *Server) Response(pdu _PDU) error {
	log.Debug().Msgf("response\n%s", pdu)

	if s.serverPeer == nil {
		return errors.New("unbound server")
	}
	return s.serverPeer.Confirmation(pdu)
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
		serverPeer = fmt.Sprintf(" serverPeerId: %d", s.serverPeer.getClientId())
	}
	return fmt.Sprintf("Server(cid:%d)%s", s.serverID, serverPeer)
}

// _ServiceAccessPoint is a interface used for documentation
type _ServiceAccessPoint interface {
	SapConfirmation(pdu _PDU) error
	SapRequest(pdu _PDU) error
	SapIndication(pdu _PDU) error
	SapResponse(pdu _PDU) error
	_setServiceElement(serviceElement _ApplicationServiceElement)
}

type ServiceAccessPoint struct {
	serviceID      *int
	serviceElement _ApplicationServiceElement
}

func NewServiceAccessPoint(sapID *int, rootStruct _ServiceAccessPoint) (*ServiceAccessPoint, error) {
	s := &ServiceAccessPoint{
		serviceID: sapID,
	}
	if sapID != nil {
		if _, ok := serviceMap[*sapID]; ok {
			return nil, errors.Errorf("already a server %d", *sapID)
		}
		serviceMap[*sapID] = s

		// automatically bind
		if element, ok := elementMap[*sapID]; ok {
			if element.elementService != nil {
				return nil, errors.Errorf("application service element %d already bound", *sapID)
			}

			// Note: we need to pass the rootStruct (which should contain s as delegate) here
			if err := bind(element, rootStruct); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func (s *ServiceAccessPoint) SapRequest(pdu _PDU) error {
	log.Debug().Msgf("SapRequest(%d)\n%s", s.serviceID, pdu)

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Indication(pdu)
}

func (s *ServiceAccessPoint) SapIndication(_PDU) error {
	panic("this should be implemented by outer struct")
}

func (s *ServiceAccessPoint) SapResponse(pdu _PDU) error {
	log.Debug().Msgf("SapResponse(%d)\n%s", s.serviceID, pdu)

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Confirmation(pdu)
}

func (s *ServiceAccessPoint) SapConfirmation(_PDU) error {
	panic("this should be implemented by outer struct")
}

func (s *ServiceAccessPoint) _setServiceElement(serviceElement _ApplicationServiceElement) {
	s.serviceElement = serviceElement
}

// _ApplicationServiceElement is a interface used for documentation
type _ApplicationServiceElement interface {
	Request(pdu _PDU) error
	Indication(pdu _PDU) error
	Response(pdu _PDU) error
	Confirmation(pdu _PDU) error
	_setElementService(elementService _ServiceAccessPoint)
}

type ApplicationServiceElement struct {
	elementID      *int
	elementService _ServiceAccessPoint
}

func NewApplicationServiceElement(aseID *int, rootStruct _ApplicationServiceElement) (*ApplicationServiceElement, error) {
	a := &ApplicationServiceElement{
		elementID: aseID,
	}

	if aseID != nil {
		if _, ok := elementMap[*aseID]; ok {
			return nil, errors.Errorf("already an application service element %d", *aseID)
		}
		elementMap[*aseID] = a

		// automatically bind
		if service, ok := serviceMap[*aseID]; ok {
			if service.serviceElement != nil {
				return nil, errors.Errorf("service access point %d already bound", *aseID)
			}

			// Note: we need to pass the rootStruct (which should contain a as delegate) here
			if err := bind(rootStruct, service); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return a, nil
}

func (a *ApplicationServiceElement) Request(pdu _PDU) error {
	log.Debug().Msgf("Request\n%s", pdu)

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapIndication(pdu)
}

func (a *ApplicationServiceElement) Indication(_PDU) error {
	panic("this should be implemented by outer struct")
}

func (a *ApplicationServiceElement) Response(pdu _PDU) error {
	log.Debug().Msgf("Response\n%s", pdu)

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapConfirmation(pdu)
}

func (a *ApplicationServiceElement) Confirmation(_PDU) error {
	panic("this should be implemented by outer struct")
}

func (a *ApplicationServiceElement) _setElementService(elementService _ServiceAccessPoint) {
	a.elementService = elementService
}

// bind a list of clients and servers together, top down
func bind(args ...interface{}) error {
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

			if err := bind(client, server); err != nil {
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

			if err := bind(element, service); err != nil {
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
		client := args[i]
		log.Debug().Msgf("client %v", client)
		server := args[i+1]
		log.Debug().Msgf("server %v", server)

		// make sure we're binding clients and servers
		clientCast, okClient := client.(_Client)
		serverCast, okServer := server.(_Server)
		elementServiceCast, okElementService := client.(_ApplicationServiceElement)
		serviceAccessPointCast, okServiceAccessPoint := server.(_ServiceAccessPoint)
		if okClient && okServer {
			clientCast._setClientPeer(serverCast)
			serverCast._setServerPeer(clientCast)
		} else if okElementService && okServiceAccessPoint { // we could be binding application clients and servers
			elementServiceCast._setElementService(serviceAccessPointCast)
			serviceAccessPointCast._setServiceElement(elementServiceCast)
		} else {
			return errors.New("bind() requires a client and a server")
		}
	}
	log.Debug().Msg("bound")
	return nil
}
