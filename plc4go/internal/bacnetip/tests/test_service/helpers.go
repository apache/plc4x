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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// This struct turns off the deferred startup function call that broadcasts I-Am-Router-To-Network and Network-Number-Is
//
//	messages.
type _NetworkServiceElement struct {
	*bacnetip.NetworkServiceElement
}

func new_NetworkServiceElement(localLog zerolog.Logger) (*_NetworkServiceElement, error) {
	n := &_NetworkServiceElement{}
	var err error
	n.NetworkServiceElement, err = bacnetip.NewNetworkServiceElement(localLog, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	return n, nil
}

type ApplicationNetworkRequirements interface {
	Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

type ApplicationNetwork struct {
	ApplicationNetworkRequirements
	*tests.StateMachineGroup

	trafficLog      *tests.TrafficLog
	vlan            *bacnetip.Network
	tdDeviceObject  *bacnetip.LocalDeviceObject
	td              *ApplicationStateMachine
	iutDeviceObject *bacnetip.LocalDeviceObject
	iut             *ApplicationStateMachine

	log zerolog.Logger
}

func NewApplicationNetwork(localLog zerolog.Logger, applicationNetworkRequirements ApplicationNetworkRequirements) (*ApplicationNetwork, error) {
	a := &ApplicationNetwork{
		ApplicationNetworkRequirements: applicationNetworkRequirements,
		log:                            localLog,
	}
	a.StateMachineGroup = tests.NewStateMachineGroup(localLog)

	// Reset the time machine
	tests.ResetTimeMachine(time.Time{})

	// Create a traffic log
	a.trafficLog = new(tests.TrafficLog)

	// make a little LAN
	a.vlan = bacnetip.NewNetwork(localLog, bacnetip.WithNetworkBroadcastAddress(bacnetip.NewLocalBroadcast(nil)), bacnetip.WithNetworkTrafficLogger(a.trafficLog))

	// test device object
	octets1024 := model.MaxApduLengthAccepted_NUM_OCTETS_1024
	segmentation := model.BACnetSegmentation_NO_SEGMENTATION
	a.tdDeviceObject = &bacnetip.LocalDeviceObject{
		ObjectName:                "td",
		ObjectIdentifier:          "device:10",
		MaximumApduLengthAccepted: &octets1024,
		SegmentationSupported:     &segmentation,
		VendorIdentifier:          999,
	}

	// test device
	var err error
	a.td, err = NewApplicationStateMachine(localLog, a.tdDeviceObject, a.vlan, a)
	if err != nil {
		return nil, errors.Wrap(err, "error building application state machine")
	}
	a.Append(a.td)

	// implementation under test device object
	octets1024 = model.MaxApduLengthAccepted_NUM_OCTETS_1024
	segmentation = model.BACnetSegmentation_NO_SEGMENTATION
	a.iutDeviceObject = &bacnetip.LocalDeviceObject{
		ObjectName:                "iut",
		ObjectIdentifier:          "device:20",
		MaximumApduLengthAccepted: &octets1024,
		SegmentationSupported:     &segmentation,
		VendorIdentifier:          999,
	}

	// implementation under test
	a.iut, err = NewApplicationStateMachine(localLog, a.iutDeviceObject, a.vlan, a)
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
	tests.RunTimeMachine(a.log, timeLimit, time.Time{})
	if a.log.Debug().Enabled() {
		a.log.Debug().Msg("time machine finished")
		for _, machine := range a.GetStateMachines() {
			a.log.Debug().Stringer("machine", machine).Msg("machine")
			for _, entry := range machine.GetTransactionLog() {
				a.log.Debug().Str("entry", entry).Msg("transaction log entry")
			}
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

func (a *ApplicationNetwork) _debug(format string, args bacnetip.Args) {
	a.log.Debug().Msgf(format, args)
}

type SnifferNode struct {
	*bacnetip.Client

	name    string
	address *bacnetip.Address
	node    *bacnetip.Node

	log zerolog.Logger
}

func NewSnifferNode(localLog zerolog.Logger, vlan *bacnetip.Network) (*SnifferNode, error) {
	s := &SnifferNode{
		name: "sniffer",
		log:  localLog,
	}
	s.address, _ = bacnetip.NewAddress(localLog)
	var err error
	s.Client, err = bacnetip.NewClient(localLog, s)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}

	// create a promiscuous node, added to the network
	s.node, err = bacnetip.NewNode(localLog, s.address, vlan, bacnetip.WithNodePromiscuous(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	s.log.Debug().Stringer("node", s.node).Msg("node")

	// bind the node
	err = bacnetip.Bind(s.log, s, s.node)
	if err != nil {
		return nil, errors.Wrap(err, "error binding node")
	}
	return s, nil
}

func (s *SnifferNode) Request(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("request")
	return errors.New("sniffers don't request")
}

func (s *SnifferNode) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("confirmation")
	pdu := args.Get0PDU()

	// it's and NPDU
	npdu := pdu.GetMessage().(model.NPDU)

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

type SnifferStateMachine struct {
	*bacnetip.Client
	tests.StateMachine

	name    string
	address *bacnetip.Address
	node    *bacnetip.Node

	log zerolog.Logger
}

func NewSnifferStateMachine(localLog zerolog.Logger, vlan *bacnetip.Network) (*SnifferStateMachine, error) {
	s := &SnifferStateMachine{
		name: "sniffer",
		log:  localLog,
	}
	s.address, _ = bacnetip.NewAddress(localLog)

	// continue with initialization
	var err error
	s.Client, err = bacnetip.NewClient(localLog, s)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	var init func()
	s.StateMachine, init = tests.NewStateMachine(localLog, s)
	init()

	// create a promiscuous node, added to the network
	s.node, err = bacnetip.NewNode(localLog, s.address, vlan, bacnetip.WithNodePromiscuous(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	s.log.Debug().Stringer("node", s.node).Msg("node")

	// bind the node
	err = bacnetip.Bind(s.log, s, s.node)
	if err != nil {
		return nil, errors.Wrap(err, "error binding node")
	}
	return s, nil
}

func (s *SnifferStateMachine) Sends(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("request")
	return errors.New("sniffers don't send")
}

func (s *SnifferStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("confirmation")
	pdu := args.Get0PDU()

	// it's and NPDU
	npdu := pdu.GetMessage().(model.NPDU)

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
	return s.Receive(args, bacnetip.NoKWArgs)
}

func (s *SnifferStateMachine) String() string {
	return "SnifferStateMachine" //TODO
}

type ApplicationStateMachineRequirements interface {
	Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
	Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

type ApplicationStateMachine struct {
	ApplicationStateMachineRequirements
	*bacnetip.ApplicationIOController
	tests.StateMachine

	address *bacnetip.Address
	asap    *bacnetip.ApplicationServiceAccessPoint
	smap    *bacnetip.StateMachineAccessPoint
	nsap    *bacnetip.NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	node    *bacnetip.Node

	log zerolog.Logger
}

func NewApplicationStateMachine(localLog zerolog.Logger, localDevice *bacnetip.LocalDeviceObject, vlan *bacnetip.Network, applicationStateMachineRequirements ApplicationStateMachineRequirements) (*ApplicationStateMachine, error) {
	a := &ApplicationStateMachine{
		ApplicationStateMachineRequirements: applicationStateMachineRequirements,
		log:                                 localLog,
	}

	// build and address and save it
	_, instance := bacnetip.ObjectIdentifierStringToTuple(localDevice.ObjectIdentifier)
	var err error
	a.address, err = bacnetip.NewAddress(a.log, instance)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	a.log.Debug().Stringer("address", a.address).Msg("address")

	// continue with initialization
	a.ApplicationIOController, err = bacnetip.NewApplicationIOController(a.log, localDevice)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application io controller")
	}
	var init func()
	a.StateMachine, init = tests.NewStateMachine(a.log, a, tests.WithStateMachineName(localDevice.ObjectName))
	init()

	// include a application decoder
	a.asap, err = bacnetip.NewApplicationServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation

	// the segmentation state machines need access to the same device
	// information cache as the application
	deviceInfoCache := a.GetDeviceInfoCache()
	a.smap, err = bacnetip.NewStateMachineAccessPoint(a.log, localDevice, bacnetip.WithStateMachineAccessPointDeviceInfoCache(deviceInfoCache))
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// a network service access point will be needed
	a.nsap, err = bacnetip.NewNetworkServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacnetip.Bind(a.log, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding network service element")
	}

	// bind the top layers
	err = bacnetip.Bind(a.log, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding top layers")
	}

	// create a node, added to the network
	a.node, err = bacnetip.NewNode(a.log, a.address, vlan)
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

func (a *ApplicationStateMachine) String() string {
	return "ApplicationStateMachine" //TODO:
}

func (a *ApplicationStateMachine) Send(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")

	pdu := args.Get0PDU()
	// build a IOCB to wrap the request
	iocb, err := bacnetip.NewIOCB(a.log, pdu, nil)
	if err != nil {
		return errors.Wrap(err, "error creating iocb")
	}
	return a.Request(bacnetip.NewArgs(iocb), bacnetip.NoKWArgs)
}

func (a *ApplicationStateMachine) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")

	// let the state machine know the request was received
	err := a.Receive(args, bacnetip.NoKWArgs)
	if err != nil {
		return errors.Wrap(err, "error receiving indication")
	}

	// allow the application to process it
	return a.ApplicationStateMachineRequirements.Indication(args, kwargs)
}

func (a *ApplicationStateMachine) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")

	// forward the confirmation to the state machine
	err := a.Receive(args, bacnetip.NoKWArgs)
	if err != nil {
		return errors.Wrap(err, "error receiving indication")
	}

	// allow the application to process it
	return a.ApplicationStateMachineRequirements.Confirmation(args, kwargs)
}

type COVTestClientServicesRequirements interface {
}

type COVTestClientServices struct {
	COVTestClientServicesRequirements
	*bacnetip.Capability

	log zerolog.Logger
}

func (c *COVTestClientServices) doConfirmedCOVNotificationRequest(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	c.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("doConfirmedCOVNotificationRequest")

	panic("TODO: implement me") // TODO:implement me
	return nil
}

func (c *COVTestClientServices) doUnconfirmedCOVNotificationRequest(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	c.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("doUnconfirmedCOVNotificationRequest")

	panic("TODO: implement me") // TODO:implement me
	return nil
}
