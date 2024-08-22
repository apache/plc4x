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

	"github.com/apache/plc4x/plc4go/internal/bacnetip"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/tests"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	trafficLog *tests.TrafficLog
	vlan       *bacnetip.Network
	sniffer    *SnifferNode
	td         *ApplicationStateMachine
	iut        *ApplicationStateMachine

	log zerolog.Logger
}

func NewApplicationNetwork(localLog zerolog.Logger, applicationNetworkRequirements ApplicationNetworkRequirements, tdDeviceObject, iutDeviceObject *bacnetip.LocalDeviceObject) (*ApplicationNetwork, error) {
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
	a.vlan = bacnetip.NewNetwork(a.log, bacnetip.WithNetworkBroadcastAddress(bacnetip.NewLocalBroadcast(nil)), bacnetip.WithNetworkTrafficLogger(a.trafficLog))

	// sniffer
	var err error
	a.sniffer, err = NewSnifferNode(a.log, a.vlan)
	if err != nil {
		return nil, errors.Wrap(err, "error creating sniffer node")
	}

	// test device
	a.td, err = NewApplicationStateMachine(a.log, tdDeviceObject, a.vlan, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application state machine")
	}
	a.Append(a.td)

	// implementation under test
	a.iut, err = NewApplicationStateMachine(a.log, iutDeviceObject, a.vlan, a)
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

type ApplicationStateMachineRequirements interface {
	Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error
}

type ApplicationStateMachine struct {
	ApplicationStateMachineRequirements
	*bacnetip.Application
	tests.StateMachine

	address *bacnetip.Address
	asap    *bacnetip.ApplicationServiceAccessPoint
	smap    *bacnetip.StateMachineAccessPoint
	nsap    *bacnetip.NetworkServiceAccessPoint
	nse     *_NetworkServiceElement
	node    *bacnetip.Node

	confirmedPrivateResult string

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
	a.Application, err = bacnetip.NewApplication(a.log, localDevice)
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

	// send the apdu down the stack
	return a.Request(args, kwargs)
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
	return a.Receive(args, kwargs)
}

// TODO
func (a *ApplicationStateMachine) doConfirmedPrivateTransferRequest(_ struct{}) {
	// TODO: who calls that?
}

type TODOWHATTODOWITHTHAT struct {
}

func (T TODOWHATTODOWITHTHAT) Indication(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	//TODO implement me
	panic("implement me")
}

func (T TODOWHATTODOWITHTHAT) Confirmation(args bacnetip.Args, kwargs bacnetip.KWArgs) error {
	//TODO implement me
	panic("implement me")
}

func SegmentationTest(t *testing.T, prefix string, cLen, sLen int) {
	tests.LockGlobalTimeMachine(t)
	testingLogger := testutils.ProduceTestingLogger(t)
	tests.NewGlobalTimeMachine(testingLogger)

	// client device object
	octets206 := model.MaxApduLengthAccepted_NUM_OCTETS_206
	segmentation := model.BACnetSegmentation_SEGMENTED_BOTH
	maxSegmentsAccepted := model.MaxSegmentsAccepted_NUM_SEGMENTS_04
	tdDeviceObject := &bacnetip.LocalDeviceObject{
		ObjectName:                "td",
		ObjectIdentifier:          "device:10",
		MaximumApduLengthAccepted: &octets206,
		SegmentationSupported:     &segmentation,
		MaxSegmentsAccepted:       &maxSegmentsAccepted,
		VendorIdentifier:          999,
	}

	// server device object
	maxSegmentsAccepted = model.MaxSegmentsAccepted_NUM_SEGMENTS_64
	iutDeviceObject := &bacnetip.LocalDeviceObject{
		ObjectName:                "td",
		ObjectIdentifier:          "device:10",
		MaximumApduLengthAccepted: &octets206,
		SegmentationSupported:     &segmentation,
		MaxSegmentsAccepted:       &maxSegmentsAccepted,
		VendorIdentifier:          999,
	}

	// create a network
	anet, err := NewApplicationNetwork(testingLogger, new(TODOWHATTODOWITHTHAT), tdDeviceObject, iutDeviceObject)
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
		anet.iut.confirmedPrivateResult = utils.RandomString(sLen)
	}

	wrapRequestString := func() (
		numberOfDataElements model.BACnetApplicationTagUnsignedInteger,
		data []model.BACnetConstructedDataElement,
		openingTag model.BACnetOpeningTag,
		peekedTagHeader model.BACnetTagHeader,
		closingTag model.BACnetClosingTag,
		tagNumber uint8,
		arrayIndexArgument model.BACnetTagPayloadUnsignedInteger,
	) {
		numberOfDataElements = model.CreateBACnetApplicationTagUnsignedInteger(1)
		elementArgCreator := func() (
			peekedTagHeader model.BACnetTagHeader,
			applicationTag model.BACnetApplicationTag,
			contextTag model.BACnetContextTag,
			constructedData model.BACnetConstructedData,
			objectTypeArgument model.BACnetObjectType,
			propertyIdentifierArgument model.BACnetPropertyIdentifier,
			arrayIndexArgument model.BACnetTagPayloadUnsignedInteger,
		) {
			peekedTagHeader = model.CreateBACnetTagHeaderBalanced(false, 0, 0)
			applicationTag = model.CreateBACnetApplicationTagCharacterString(model.BACnetCharacterEncoding_ISO_8859_1, requestString)
			return
		}
		data = []model.BACnetConstructedDataElement{
			model.NewBACnetConstructedDataElement(elementArgCreator()),
		}
		openingTag = model.CreateBACnetOpeningTag(2)
		peekedTagHeader = model.CreateBACnetTagHeaderBalanced(false, 2, 2)
		closingTag = model.CreateBACnetClosingTag(2)
		return
	}

	cpt := model.NewBACnetConfirmedServiceRequestConfirmedPrivateTransfer(
		model.CreateBACnetVendorIdContextTagged(0, 999),
		model.CreateBACnetContextTagUnsignedInteger(1, 1),
		model.NewBACnetConstructedDataUnspecified(wrapRequestString()),
		0,
	)

	apdu := model.NewAPDUConfirmedRequest(
		false,
		false,
		true,
		maxSegmentsAccepted,
		model.MaxApduLengthAccepted_NUM_OCTETS_1024,
		0,
		nil,
		nil,
		cpt,
		nil,
		nil,
		0)

	rq := bacnetip.NewPDU(apdu)

	var trq model.BACnetServiceAckConfirmedPrivateTransfer
	// send the request, get it acked
	anet.td.GetStartState().Doc(prefix+"-0").
		Send(rq, nil).Doc(prefix+"-1").
		Receive(bacnetip.NewArgs(trq), bacnetip.NoKWArgs).Doc(prefix + "-2").
		Success("")

	// no IUT application layer matching
	anet.iut.GetStartState().Success("")

	// run the group
	err = anet.Run(9)
	assert.NoError(t, err)
}

func Test1(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-1", 0, 0)
}

func Test2(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-2", 10, 0)
}

func Test3(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-3", 100, 0)
}

func Test4(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-4", 200, 0)
}

func Test5(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-5", 0, 10)
}

func Test6(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-6", 0, 200)
}

func Test7(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-7", 300, 0)
}

func Test8(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-8", 300, 300)
}

func Test9(t *testing.T) {
	t.Skip("not ready yet") // TODO: figure out why it is failing
	SegmentationTest(t, "7-9", 600, 600)
}
