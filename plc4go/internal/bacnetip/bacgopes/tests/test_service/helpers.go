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

package test_service

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/app"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/capability"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/iocb"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/netservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/service"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/state_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests/time_machine"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/vlan"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// This struct turns off the deferred startup function call that broadcasts I-Am-Router-To-Network and Network-Number-Is
//
//	messages.
type _NetworkServiceElement struct {
	*NetworkServiceElement
}

func new_NetworkServiceElement(localLog zerolog.Logger, options ...Option) (*_NetworkServiceElement, error) {
	n := &_NetworkServiceElement{}
	ApplyAppliers(options, n)
	optionsForParent := AddLeafTypeIfAbundant(options, n)
	var err error
	n.NetworkServiceElement, err = NewNetworkServiceElement(localLog, Combine(optionsForParent, WithNetworkServiceElementStartupDisabled(true))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	return n, nil
}

type ApplicationNetwork struct {
	*StateMachineGroup

	trafficLog      *TrafficLog
	vlan            *Network
	tdDeviceObject  LocalDeviceObject
	td              *ApplicationStateMachine
	iutDeviceObject LocalDeviceObject
	iut             *ApplicationStateMachine

	log zerolog.Logger
}

func NewApplicationNetwork(localLog zerolog.Logger) (*ApplicationNetwork, error) {
	a := &ApplicationNetwork{
		log: localLog,
	}
	a.StateMachineGroup = NewStateMachineGroup(localLog)

	// Reset the time machine
	ResetTimeMachine(time.Time{})

	// Create a traffic log
	a.trafficLog = new(TrafficLog)

	// make a little LAN
	a.vlan = NewNetwork(localLog, WithNetworkBroadcastAddress(NewLocalBroadcast(nil)), WithNetworkTrafficLogger(a.trafficLog))

	var err error
	// test device object
	octets1024 := model.MaxApduLengthAccepted_NUM_OCTETS_1024
	segmentation := model.BACnetSegmentation_NO_SEGMENTATION
	a.tdDeviceObject, err = NewLocalDeviceObject(
		NoArgs,
		NKW(
			KWObjectName, "td",
			KWObjectIdentifier, "device:10",
			KWMaximumApduLengthAccepted, &octets1024,
			KWSegmentationSupported, &segmentation,
			KWVendorIdentifier, 999,
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating test device")
	}

	// test device
	a.td, err = NewApplicationStateMachine(localLog, a.tdDeviceObject, a.vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error building application state machine")
	}
	a.Append(a.td)

	// implementation under test device object
	octets1024 = model.MaxApduLengthAccepted_NUM_OCTETS_1024
	segmentation = model.BACnetSegmentation_NO_SEGMENTATION
	a.iutDeviceObject, err = NewLocalDeviceObject(
		NoArgs,
		NKW(
			KWObjectName, "iut",
			KWObjectIdentifier, "device:20",
			KWMaximumApduLengthAccepted, &octets1024,
			KWSegmentationSupported, &segmentation,
			KWVendorIdentifier, 999,
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating test device")
	}

	// implementation under test
	a.iut, err = NewApplicationStateMachine(localLog, a.iutDeviceObject, a.vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error building application state machine")
	}
	a.Append(a.iut)

	return a, nil
}

func (a *ApplicationNetwork) Run(timeLimit time.Duration) error {
	if timeLimit == 0 {
		timeLimit = 60 * time.Second
	}
	a.log.Debug().Dur("timeLimit", timeLimit).Msg("run")

	// Run the group
	err := a.StateMachineGroup.Run()
	if err != nil {
		return errors.Wrap(err, "error running state machine group")
	}
	a.log.Trace().Msg("group running")

	// run it for some time
	RunTimeMachine(a.log, timeLimit, time.Time{})
	if a.log.Debug().Enabled() {
		a.log.Debug().Msg("time machine finished")
		for _, machine := range a.GetStateMachines() {
			a.log.Debug().Stringer("machine", machine).Stringers("entries", ToStringers(machine.GetTransactionLog())).Msg("machine")
		}

		a.trafficLog.Dump(a._debug)
	}

	// check for success
	success, failed := a.CheckForSuccess()
	if !success {
		return errors.New("application network did not succeed")
	}
	_ = failed

	return nil
}

func (a *ApplicationNetwork) _debug(format string, args Args) {
	a.log.Debug().Msgf(format, args)
}

func (a *ApplicationNetwork) Close() error {
	defer utils.StopWarn(a.log)()
	if err := a.td.Close(); err != nil {
		a.log.Warn().Err(err).Msg("error closing td")
	}
	if err := a.iut.Close(); err != nil {
		a.log.Warn().Err(err).Msg("error closing iut")
	}
	return nil
}

//go:generate go tool plc4xGenerator -type=SnifferNode
type SnifferNode struct {
	ClientContract

	name    string
	address *Address
	node    *Node

	log zerolog.Logger
}

func NewSnifferNode(localLog zerolog.Logger, vlan *Network) (*SnifferNode, error) {
	s := &SnifferNode{
		name: "sniffer",
		log:  localLog,
	}
	s.address, _ = NewAddress(NoArgs)
	var err error
	s.ClientContract, err = NewClient(localLog)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}

	// create a promiscuous node, added to the network
	s.node, err = NewNode(localLog, s.address, WithNodeLan(vlan), WithNodePromiscuous(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	s.log.Debug().Stringer("node", s.node).Msg("node")

	// bind the node
	err = Bind(s.log, s, s.node)
	if err != nil {
		return nil, errors.Wrap(err, "error binding node")
	}
	return s, nil
}

func (s *SnifferNode) Request(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("request")
	return errors.New("sniffers don't request")
}

func (s *SnifferNode) Confirmation(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("confirmation")
	pdu := GA[PDU](args, 0)

	// it's and NPDU
	npdu := pdu.GetRootMessage().(model.NPDU)

	// filter out network layer traffic if there is any, probably not
	if nlm := npdu.GetNlm(); nlm != nil {
		s.log.Debug().Stringer("nlm", nlm).Msg("network message")
		return nil
	}

	// decode as generic APDU
	apdu := npdu.GetApdu()

	// TODO: not sure what to do here
	/*
	 # "lift" the source and destination address
	        if npdu.npduSADR:
	            apdu.pduSource = npdu.npduSADR
	        else:
	            apdu.pduSource = npdu.pduSource
	        if npdu.npduDADR:
	            apdu.pduDestination = npdu.npduDADR
	        else:
	            apdu.pduDestination = npdu.pduDestination
	*/
	_ = apdu

	// TODO: print etc

	return nil
}

//go:generate go tool plc4xGenerator -type=SnifferStateMachine
type SnifferStateMachine struct {
	ClientContract
	StateMachineContract

	name    string
	address *Address
	node    *Node

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, vlan *Network, options ...Option) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		name: "sniffer",
		log:  localLog,
	}
	ApplyAppliers(options, s)
	optionsForParent := AddLeafTypeIfAbundant(options, s)
	s.address, _ = NewAddress(NoArgs)

	// continue with initialization
	var err error
	s.ClientContract, err = NewClient(localLog, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	var init func()
	s.StateMachineContract, init = NewStateMachine(localLog, s)
	init()

	// create a promiscuous node, added to the network
	s.node, err = NewNode(localLog, s.address, Combine(options, WithNodeLan(vlan), WithNodePromiscuous(true))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	s.log.Debug().Stringer("node", s.node).Msg("node")

	// bind the node
	err = Bind(s.log, s, s.node)
	if err != nil {
		return nil, errors.Wrap(err, "error binding node")
	}
	return s, nil
}

func (s *SnifferStateMachine) Send(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("request")
	return errors.New("sniffers don't send")
}

func (s *SnifferStateMachine) Confirmation(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("confirmation")
	pdu := GA[PDU](args, 0)

	// it's and NPDU
	npdu := pdu.GetRootMessage().(model.NPDU)

	// filter out network layer traffic if there is any, probably not
	if nlm := npdu.GetNlm(); nlm != nil {
		s.log.Debug().Stringer("nlm", nlm).Msg("network message")
		return nil
	}

	// decode as generic APDU
	apdu := npdu.GetApdu()

	// TODO: not sure what to do here
	/*
	 # "lift" the source and destination address
	        if npdu.npduSADR:
	            apdu.pduSource = npdu.npduSADR
	        else:
	            apdu.pduSource = npdu.pduSource
	        if npdu.npduDADR:
	            apdu.pduDestination = npdu.npduDADR
	        else:
	            apdu.pduDestination = npdu.pduDestination
	*/
	_ = apdu

	// TODO: print etc

	// pass to the state machine
	return s.Receive(args, NoKWArgs())
}

//go:generate go tool plc4xGenerator -type=ApplicationStateMachine
type ApplicationStateMachine struct {
	*ApplicationIOController
	StateMachineContract

	address *Address
	asap    *ApplicationServiceAccessPoint
	smap    *StateMachineAccessPoint
	nsap    *NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	node    *Node

	log zerolog.Logger
}

func NewApplicationStateMachine(localLog zerolog.Logger, localDevice LocalDeviceObject, vlan *Network, options ...Option) (*ApplicationStateMachine, error) {
	a := &ApplicationStateMachine{
		log: localLog,
	}
	ApplyAppliers(options, a)
	optionsForParent := AddLeafTypeIfAbundant(options, a)

	// build and address and save it
	_, instance := ObjectIdentifierStringToTuple(localDevice.GetObjectIdentifier())
	var err error
	a.address, err = NewAddress(NA(instance))
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	a.log.Debug().Stringer("address", a.address).Msg("address")

	// continue with initialization
	a.ApplicationIOController, err = NewApplicationIOController(a.log, Combine(optionsForParent, WithApplicationLocalDeviceObject(localDevice))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application io controller")
	}
	var init func()
	a.StateMachineContract, init = NewStateMachine(a.log, a, Combine(optionsForParent, WithStateMachineName(localDevice.GetObjectName()))...)
	init()

	// include a application decoder
	a.asap, err = NewApplicationServiceAccessPoint(a.log, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation

	// the segmentation state machines need access to the same device
	// information cache as the application
	deviceInfoCache := a.GetDeviceInfoCache()
	a.smap, err = NewStateMachineAccessPoint(a.log, localDevice, Combine(options, WithStateMachineAccessPointDeviceInfoCache(deviceInfoCache))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// a network service access point will be needed
	a.nsap, err = NewNetworkServiceAccessPoint(a.log, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(a.log, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = Bind(a.log, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding network service element")
	}

	// bind the top layers
	err = Bind(a.log, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding top layers")
	}

	// create a node, added to the network
	a.node, err = NewNode(a.log, a.address, Combine(options, WithNodeLan(vlan))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}

	// bind the network service to the node, no network number
	err = a.nsap.Bind(a.node, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error binding nsap")
	}
	return a, nil
}

func (a *ApplicationStateMachine) Send(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Send")

	pdu := GA[PDU](args, 0)
	// build a IOCB to wrap the request
	iocb, err := NewIOCB(a.log, pdu, nil)
	if err != nil {
		return errors.Wrap(err, "error creating iocb")
	}
	return a.Request(NA(iocb), NoKWArgs())
}

func (a *ApplicationStateMachine) Indication(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Indication")

	// let the state machine know the request was received
	err := a.Receive(args, NoKWArgs())
	if err != nil {
		return errors.Wrap(err, "error receiving indication")
	}

	// allow the application to process it
	return a.Application.Indication(args, kwArgs)
}

func (a *ApplicationStateMachine) Confirmation(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("Confirmation")

	// forward the confirmation to the state machine
	err := a.Receive(args, NoKWArgs())
	if err != nil {
		return errors.Wrap(err, "error receiving indication")
	}

	// allow the application to process it
	return a.ApplicationIOController.Confirmation(args, kwArgs)
}

type COVTestClientServicesRequirements interface {
}

type COVTestClientServices struct {
	COVTestClientServicesRequirements
	Capability

	log zerolog.Logger
}

func (c *COVTestClientServices) doConfirmedCOVNotificationRequest(args Args, kwArgs KWArgs) error {
	c.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("doConfirmedCOVNotificationRequest")

	panic("TODO: implement me") // TODO:implement me
	return nil
}

func (c *COVTestClientServices) doUnconfirmedCOVNotificationRequest(args Args, kwArgs KWArgs) error {
	c.log.Debug().Stringer("args", args).Stringer("kwArgs", kwArgs).Msg("doUnconfirmedCOVNotificationRequest")

	panic("TODO: implement me") // TODO:implement me
	return nil
}
