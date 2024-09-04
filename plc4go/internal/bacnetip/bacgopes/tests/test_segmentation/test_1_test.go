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

package test_segmentation

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructors"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/tests"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// This struct turns off the deferred startup function call that broadcasts I-Am-Router-To-Network and Network-Number-Is
//
//	messages.
type _NetworkServiceElement struct {
	*bacgopes.NetworkServiceElement
}

func new_NetworkServiceElement(localLog zerolog.Logger) (*_NetworkServiceElement, error) {
	n := &_NetworkServiceElement{}
	var err error
	n.NetworkServiceElement, err = bacgopes.NewNetworkServiceElement(localLog, bacgopes.WithNetworkServiceElementStartupDisabled(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	return n, nil
}

type ApplicationNetwork struct {
	*tests.StateMachineGroup

	trafficLog *tests.TrafficLog
	vlan       *bacgopes.Network
	sniffer    *SnifferNode
	td         *ApplicationStateMachine
	iut        *ApplicationStateMachine

	log zerolog.Logger
}

func NewApplicationNetwork(localLog zerolog.Logger, tdDeviceObject, iutDeviceObject *bacgopes.LocalDeviceObject) (*ApplicationNetwork, error) {
	a := &ApplicationNetwork{
		log: localLog,
	}
	a.StateMachineGroup = tests.NewStateMachineGroup(localLog)

	// Reset the time machine
	tests.ResetTimeMachine(time.Time{})

	// Create a traffic log
	a.trafficLog = new(tests.TrafficLog)

	// make a little LAN
	a.vlan = bacgopes.NewNetwork(a.log, bacgopes.WithNetworkBroadcastAddress(bacgopes.NewLocalBroadcast(nil)), bacgopes.WithNetworkTrafficLogger(a.trafficLog))

	// sniffer
	var err error
	a.sniffer, err = NewSnifferNode(a.log, a.vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating sniffer node")
	}

	// test device
	a.td, err = NewApplicationStateMachine(a.log, tdDeviceObject, a.vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application state machine")
	}
	a.Append(a.td)

	// implementation under test
	a.iut, err = NewApplicationStateMachine(a.log, iutDeviceObject, a.vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application state machine")
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

func (a *ApplicationNetwork) _debug(format string, args bacgopes.Args) {
	a.log.Debug().Msgf(format, args)
}

//go:generate plc4xGenerator -type=SnifferNode -suffix=_test
type SnifferNode struct {
	bacgopes.Client

	name    string
	address *bacgopes.Address
	node    *bacgopes.Node

	log zerolog.Logger
}

func NewSnifferNode(localLog zerolog.Logger, vlan *bacgopes.Network) (*SnifferNode, error) {
	s := &SnifferNode{
		name: "sniffer",
		log:  localLog,
	}
	s.address, _ = bacgopes.NewAddress(localLog)
	var err error
	s.Client, err = bacgopes.NewClient(localLog, s)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}

	// create a promiscuous node, added to the network
	s.node, err = bacgopes.NewNode(localLog, s.address, bacgopes.WithNodeLan(vlan), bacgopes.WithNodePromiscuous(true))
	if err != nil {
		return nil, errors.Wrap(err, "error creating node")
	}
	s.log.Debug().Stringer("node", s.node).Msg("node")

	// bind the node
	err = bacgopes.Bind(s.log, s, s.node)
	if err != nil {
		return nil, errors.Wrap(err, "error binding node")
	}
	return s, nil
}

func (s *SnifferNode) Request(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("request")
	return errors.New("sniffers don't request")
}

func (s *SnifferNode) Confirmation(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	s.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("confirmation")
	pdu := args.Get0PDU()

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

//go:generate plc4xGenerator -type=ApplicationStateMachine -suffix=_test
type ApplicationStateMachine struct {
	*bacgopes.Application
	tests.StateMachineContract

	address *bacgopes.Address
	asap    *bacgopes.ApplicationServiceAccessPoint
	smap    *bacgopes.StateMachineAccessPoint
	nsap    *bacgopes.NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	node    *bacgopes.Node

	confirmedPrivateResult any

	log zerolog.Logger
}

func NewApplicationStateMachine(localLog zerolog.Logger, localDevice *bacgopes.LocalDeviceObject, vlan *bacgopes.Network, opts ...func(*ApplicationStateMachine)) (*ApplicationStateMachine, error) {
	a := &ApplicationStateMachine{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}
	// build and address and save it
	_, instance := bacgopes.ObjectIdentifierStringToTuple(localDevice.ObjectIdentifier)
	var err error
	a.address, err = bacgopes.NewAddress(a.log, instance)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	a.log.Debug().Stringer("address", a.address).Msg("address")

	// continue with initialization
	a.Application, err = bacgopes.NewApplication(a.log, localDevice)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application io controller")
	}
	var init func()
	a.StateMachineContract, init = tests.NewStateMachine(a.log, a, tests.WithStateMachineName(localDevice.ObjectName))
	init()

	// include a application decoder
	a.asap, err = bacgopes.NewApplicationServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point so it
	// can know if it should support segmentation

	// the segmentation state machines need access to the same device
	// information cache as the application
	deviceInfoCache := a.GetDeviceInfoCache()
	a.smap, err = bacgopes.NewStateMachineAccessPoint(a.log, localDevice, bacgopes.WithStateMachineAccessPointDeviceInfoCache(deviceInfoCache))
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// a network service access point will be needed
	a.nsap, err = bacgopes.NewNetworkServiceAccessPoint(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	a.nse, err = new_NetworkServiceElement(a.log)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service element")
	}
	err = bacgopes.Bind(a.log, a.nse, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding network service element")
	}

	// bind the top layers
	err = bacgopes.Bind(a.log, a, a.asap, a.smap, a.nsap)
	if err != nil {
		return nil, errors.Wrap(err, "error binding top layers")
	}

	// create a node, added to the network
	a.node, err = bacgopes.NewNode(a.log, a.address, bacgopes.WithNodeLan(vlan))
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

func (a *ApplicationStateMachine) Send(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Send")

	// send the apdu down the stack
	return a.Request(args, kwargs)
}

func (a *ApplicationStateMachine) Indication(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Indication")

	// let the state machine know the request was received
	err := a.Receive(args, bacgopes.NoKWArgs)
	if err != nil {
		return errors.Wrap(err, "error receiving indication")
	}

	// allow the application to process it
	return a.Application.Indication(args, kwargs)
}

func (a *ApplicationStateMachine) Confirmation(args bacgopes.Args, kwargs bacgopes.KWArgs) error {
	a.log.Debug().Stringer("args", args).Stringer("kwargs", kwargs).Msg("Confirmation")

	// forward the confirmation to the state machine
	return a.Receive(args, kwargs)
}

// TODO
func (a *ApplicationStateMachine) doConfirmedPrivateTransferRequest(_ struct{}) {
	// TODO: who calls that?
}

func SegmentationTest(t *testing.T, prefix string, cLen, sLen int) {
	t.Skip("to many things missing here... TODO: finish me") // TODO: finish me
	tests.ExclusiveGlobalTimeMachine(t)
	testingLogger := testutils.ProduceTestingLogger(t)

	// client device object
	octets206 := model.MaxApduLengthAccepted_NUM_OCTETS_206
	segmentation := model.BACnetSegmentation_SEGMENTED_BOTH
	maxSegmentsAccepted := model.MaxSegmentsAccepted_NUM_SEGMENTS_04
	tdDeviceObject := &bacgopes.LocalDeviceObject{
		ObjectName:                "td",
		ObjectIdentifier:          "device:10",
		MaximumApduLengthAccepted: &octets206,
		SegmentationSupported:     &segmentation,
		MaxSegmentsAccepted:       &maxSegmentsAccepted,
		VendorIdentifier:          999,
	}

	// server device object
	maxSegmentsAccepted = model.MaxSegmentsAccepted_NUM_SEGMENTS_64
	iutDeviceObject := &bacgopes.LocalDeviceObject{
		ObjectName:                "td",
		ObjectIdentifier:          "device:10",
		MaximumApduLengthAccepted: &octets206,
		SegmentationSupported:     &segmentation,
		MaxSegmentsAccepted:       &maxSegmentsAccepted,
		VendorIdentifier:          999,
	}

	// create a network
	anet, err := NewApplicationNetwork(testingLogger, tdDeviceObject, iutDeviceObject)
	require.NoError(t, err)

	// tell the device info cache of the client about the server
	if false {
		// TODO: what
		//anet.td.GetDeviceInfoCache().GetDeviceInfo(anet.iut.address.String())
		// # update the rest of the values
		//        iut_device_info.maxApduLengthAccepted = iut_device_object.maxApduLengthAccepted
		//        iut_device_info.segmentationSupported = iut_device_object.segmentationSupported
		//        iut_device_info.vendorID = iut_device_object.vendorIdentifier
		//        iut_device_info.maxSegmentsAccepted = iut_device_object.maxSegmentsAccepted
		//        iut_device_info.maxNpduLength = s_ndpu_len
	}

	// tell the device info cache of the server device about the client
	if false {
		// TODO: what
		//anet.iut.GetDeviceInfoCache().GetDeviceInfo(anet.td.address.String())
		//   # update the rest of the values
		//        td_device_info.maxApduLengthAccepted = td_device_object.maxApduLengthAccepted
		//        td_device_info.segmentationSupported = td_device_object.segmentationSupported
		//        td_device_info.vendorID = td_device_object.vendorIdentifier
		//        td_device_info.maxSegmentsAccepted = td_device_object.maxSegmentsAccepted
		//        td_device_info.maxNpduLength = c_ndpu_len
	}

	// build a request string
	var requestString string
	if cLen != 0 {
		requestString = utils.RandomString(cLen)
	}

	if sLen != 0 {
		anet.iut.confirmedPrivateResult = Any(CharacterString(utils.RandomString(sLen)))
	}

	var trq model.BACnetServiceAckConfirmedPrivateTransfer
	// send the request, get it acked
	anet.td.GetStartState().Doc(prefix+"-0").
		Send(bacgopes.NewPDU(ConfirmedPrivateTransferRequest(bacgopes.NewKWArgs(
			"vendorId", 999, "serviceNumber", 1,
			"serviceParameters", requestString,
			"destination", anet.iut.address,
		)), nil), nil).Doc(prefix+"-1").
		Receive(bacgopes.NewArgs(trq), bacgopes.NoKWArgs).Doc(prefix + "-2").
		Success("")

	// no IUT application layer matching
	anet.iut.GetStartState().Success("")

	// run the group
	err = anet.Run(9)
	assert.NoError(t, err)
}

func Test1(t *testing.T) {
	SegmentationTest(t, "7-1", 0, 0)
}

func Test2(t *testing.T) {
	SegmentationTest(t, "7-2", 10, 0)
}

func Test3(t *testing.T) {
	SegmentationTest(t, "7-3", 100, 0)
}

func Test4(t *testing.T) {
	SegmentationTest(t, "7-4", 200, 0)
}

func Test5(t *testing.T) {
	SegmentationTest(t, "7-5", 0, 10)
}

func Test6(t *testing.T) {
	SegmentationTest(t, "7-6", 0, 200)
}

func Test7(t *testing.T) {
	SegmentationTest(t, "7-7", 300, 0)
}

func Test8(t *testing.T) {
	SegmentationTest(t, "7-8", 300, 300)
}

func Test9(t *testing.T) {
	SegmentationTest(t, "7-9", 600, 600)
}
