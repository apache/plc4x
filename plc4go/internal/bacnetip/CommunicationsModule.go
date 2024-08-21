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
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/spi"

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

type IPCI interface {
	fmt.Stringer
	GetPDUUserData() spi.Message
	GetPDUSource() *Address
	SetPDUSource(source *Address)
	GetPDUDestination() *Address
	SetPDUDestination(*Address)
	Update(pci Arg) error
}

type __PCI struct {
	pduUserData    spi.Message // TODO: should that be PDUUserData rater than spi.Message and do we need another field... lets see...
	pduSource      *Address
	pduDestination *Address
}

var _ IPCI = (*__PCI)(nil)

func new__PCI(pduUserData spi.Message, pduSource *Address, pduDestination *Address) *__PCI {
	return &__PCI{pduUserData, pduSource, pduDestination}
}

func (p *__PCI) GetPDUUserData() spi.Message {
	return p.pduUserData
}

func (p *__PCI) GetPDUSource() *Address {
	return p.pduSource
}

func (p *__PCI) SetPDUSource(source *Address) {
	p.pduSource = source
}

func (p *__PCI) GetPDUDestination() *Address {
	return p.pduDestination
}

func (p *__PCI) SetPDUDestination(destination *Address) {
	p.pduDestination = destination
}

func (p *__PCI) Update(pci Arg) error {
	switch pci := pci.(type) {
	case IPCI:
		p.pduUserData = pci.GetPDUUserData()
		p.pduSource = pci.GetPDUSource()
		p.pduDestination = pci.GetPDUDestination()
		return nil
	default:
		return errors.Errorf("invalid IPCI type %T", pci)
	}
}

func (p *__PCI) deepCopy() *__PCI {
	pduUserData := p.pduUserData // those are immutable so no copy needed
	pduSource := p.pduSource
	if pduSource != nil {
		copyPduSource := *pduSource
		pduSource = &copyPduSource
	}
	pduDestination := p.pduDestination
	if pduDestination != nil {
		copyPduDestination := *pduDestination
		pduDestination = &copyPduDestination
	}
	return &__PCI{pduUserData, pduSource, pduDestination}
}

func (p *__PCI) String() string {
	pduUserDataString := ""
	if p.pduUserData != nil {
		pduUserDataString = p.pduUserData.String()
		if strings.Contains(pduUserDataString, "\n") {
			pduUserDataString = "\n" + pduUserDataString + "\n"
		}
		pduUserDataString = "pduUserData: " + pduUserDataString + " ,"
	}
	return fmt.Sprintf("__PCI{%spduSource: %s, pduDestination: %s}", pduUserDataString, p.pduSource, p.pduDestination)
}

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

			// Note: we need to pass the rootStruct (which should contain c as delegate) here
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

			// Note: we need to pass the rootStruct (which should contain s as delegate) here
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

// _ServiceAccessPoint is an interface used for documentation
type _ServiceAccessPoint interface {
	SapConfirmation(Args, KWArgs) error
	SapRequest(Args, KWArgs) error
	SapIndication(Args, KWArgs) error
	SapResponse(Args, KWArgs) error
	_setServiceElement(serviceElement _ApplicationServiceElement)
}

type ServiceAccessPoint struct {
	serviceID      *int
	serviceElement _ApplicationServiceElement

	log zerolog.Logger
}

func NewServiceAccessPoint(localLog zerolog.Logger, rootStruct _ServiceAccessPoint, opts ...func(point *ServiceAccessPoint)) (*ServiceAccessPoint, error) {
	s := &ServiceAccessPoint{
		log: localLog,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.serviceID != nil {
		sapID := *s.serviceID
		if _, ok := serviceMap[sapID]; ok {
			return nil, errors.Errorf("already a server %d", sapID)
		}
		serviceMap[sapID] = s

		// automatically bind
		if element, ok := elementMap[sapID]; ok {
			if element.elementService != nil {
				return nil, errors.Errorf("application service element %d already bound", sapID)
			}

			// Note: we need to pass the rootStruct (which should contain s as delegate) here
			if err := Bind(localLog, element, rootStruct); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return s, nil
}

func WithServiceAccessPointSapID(sapID int) func(*ServiceAccessPoint) {
	return func(s *ServiceAccessPoint) {
		s.serviceID = &sapID
	}
}

func (s *ServiceAccessPoint) String() string {
	serviceID := "-"
	if s.serviceID != nil {
		serviceID = strconv.Itoa(*s.serviceID)
	}
	return fmt.Sprintf("ServiceAccessPoint(serviceID:%v, serviceElement: %s)", serviceID, s.serviceElement)
}

func (s *ServiceAccessPoint) SapRequest(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Interface("serviceID", s.serviceID).Msg("SapRequest")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Indication(args, kwargs)
}

func (s *ServiceAccessPoint) SapIndication(Args, KWArgs) error {
	panic("this should be implemented by outer struct")
}

func (s *ServiceAccessPoint) SapResponse(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Interface("serviceID", s.serviceID).Msg("SapResponse")

	if s.serviceElement == nil {
		return errors.New("unbound service access point")
	}
	return s.serviceElement.Confirmation(args, kwargs)
}

func (s *ServiceAccessPoint) SapConfirmation(Args, KWArgs) error {
	panic("this should be implemented by outer struct")
}

func (s *ServiceAccessPoint) _setServiceElement(serviceElement _ApplicationServiceElement) {
	s.serviceElement = serviceElement
}

// _ApplicationServiceElement is an interface used for documentation
type _ApplicationServiceElement interface {
	Request(args Args, kwargs KWArgs) error
	Indication(args Args, kwargs KWArgs) error
	Response(args Args, kwargs KWArgs) error
	Confirmation(args Args, kwargs KWArgs) error
	_setElementService(elementService _ServiceAccessPoint)
}

type ApplicationServiceElement struct {
	elementID      *int
	elementService _ServiceAccessPoint

	log zerolog.Logger
}

func NewApplicationServiceElement(localLog zerolog.Logger, rootStruct _ApplicationServiceElement, opts ...func(*ApplicationServiceElement)) (*ApplicationServiceElement, error) {
	a := &ApplicationServiceElement{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}

	if a.elementID != nil {
		aseID := *a.elementID
		if _, ok := elementMap[aseID]; ok {
			return nil, errors.Errorf("already an application service element %d", aseID)
		}
		elementMap[aseID] = a

		// automatically bind
		if service, ok := serviceMap[aseID]; ok {
			if service.serviceElement != nil {
				return nil, errors.Errorf("service access point %d already bound", aseID)
			}

			// Note: we need to pass the rootStruct (which should contain a as delegate) here
			if err := Bind(localLog, rootStruct, service); err != nil {
				return nil, errors.Wrap(err, "error binding")
			}
		}
	}
	return a, nil
}

func WithApplicationServiceElementAseID(aseID int) func(*ApplicationServiceElement) {
	return func(s *ApplicationServiceElement) {
		s.elementID = &aseID
	}
}

func (a *ApplicationServiceElement) Request(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapIndication(args, kwargs)
}

func (a *ApplicationServiceElement) Indication(Args, KWArgs) error {
	panic("this should be implemented by outer struct")
}

func (a *ApplicationServiceElement) Response(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")

	if a.elementService == nil {
		return errors.New("unbound application service element")
	}

	return a.elementService.SapConfirmation(args, kwargs)
}

func (a *ApplicationServiceElement) Confirmation(Args, KWArgs) error {
	panic("this should be implemented by outer struct")
}

func (a *ApplicationServiceElement) _setElementService(elementService _ServiceAccessPoint) {
	a.elementService = elementService
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
		elementServiceCast, okElementService := left.(_ApplicationServiceElement)
		serviceAccessPointCast, okServiceAccessPoint := right.(_ServiceAccessPoint)
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
