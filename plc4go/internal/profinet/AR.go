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

package profinet

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/profinet/gsdml"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/profinet/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

const (
	blockVersionHigh = 1
	blockVersionLow  = 0

	// argsMaximum is how much response the controller accepts, the value plc4j and the captured
	// controller send.
	argsMaximum = 16696

	// The startup sequence's control commands (IEC 61158-6-10 ControlCommand).
	controlCommandPrmEnd  = 0x0001
	controlCommandRelease = 0x0004
	controlCommandDone    = 0x0008

	// realIdentificationDataIndex is the record that lists every slot and subslot a device has.
	realIdentificationDataIndex = 0xF000

	// rtUDPPort is the "port" field of AR and IOCR blocks for RT over Ethernet: the EtherType.
	rtUDPPort = 0x8892

	// Cyclic data status: primary, data valid, provider run, no station problem.
	dataStatusRun = 0x35
	// iopsGood and iocsGood are the provider and consumer status for good data.
	iopsGood = 0x80
	iocsGood = 0x80
)

var zeroUUID = make([]byte, 16)

// arSettings are the timing and identity choices for an application relation.
type arSettings struct {
	stationName     string
	sendClockFactor uint16 // in 31.25 µs
	reductionRatio  uint16
	watchdogFactor  uint16
	dataHoldFactor  uint16
	requestTimeout  time.Duration
	readyTimeout    time.Duration
	// rpcFragmentSize is the most RPC body one datagram carries; 0 is the default.
	rpcFragmentSize int
}

// cycle is the time between two cyclic frames.
func (s arSettings) cycle() time.Duration {
	return time.Duration(s.sendClockFactor) * time.Duration(s.reductionRatio) * 31250 * time.Nanosecond
}

// device is the IO device the relation is with.
type device struct {
	name               string
	mac                net.HardwareAddr
	ip                 net.IP
	vendorID, deviceID uint16
}

// applicationRelation is an IO controller's AR with one device: the connection management over
// UDP and the cyclic exchange over raw Ethernet.
type applicationRelation struct {
	settings arSettings
	device   device
	index    *gsdml.Index
	station  station
	local    net.HardwareAddr
	localIP  net.IP

	rpc  *rpcEndpoint
	plan *plan
	// modules are GSDML module ids by slot the controller expects beyond what the device reports.
	modules map[uint16]string

	arUUID        []byte
	sessionKey    uint16
	activity      readWriteModel.DceRpc_ActivityUuid
	sequence      atomic.Uint32
	inputFrameID  uint16
	outputFrameID uint16

	// deviceAlarmReference is the device's end of the alarm CR; alarmHigh and alarmLow are the two
	// alarm channels, and onAlarm gets every alarm the device sends.
	deviceAlarmReference uint16
	alarmHigh, alarmLow  *alarmChannel
	onAlarm              func(Alarm)

	mutex       sync.Mutex
	output      []byte
	input       []byte
	inputStatus byte
	lastInput   time.Time

	ready    chan struct{}
	readyOne sync.Once
	running  atomic.Bool
	lost     chan error
	lostOnce sync.Once

	cancel context.CancelFunc
	wg     sync.WaitGroup

	log zerolog.Logger
}

func newApplicationRelation(settings arSettings, dev device, index *gsdml.Index, st station, localIP net.IP, log zerolog.Logger) *applicationRelation {
	arUUID := make([]byte, 16)
	_, _ = rand.Read(arUUID)
	return &applicationRelation{
		settings:   settings,
		device:     dev,
		index:      index,
		station:    st,
		local:      st.HardwareAddress(),
		localIP:    localIP,
		arUUID:     arUUID,
		sessionKey: 1,
		activity:   newActivity(),
		alarmHigh:  newAlarmChannel(true),
		alarmLow:   newAlarmChannel(false),
		ready:      make(chan struct{}),
		lost:       make(chan error, 1),
		log:        log,
	}
}

func (ar *applicationRelation) object() readWriteModel.DceRpc_ObjectUuid {
	return readWriteModel.NewDceRpc_ObjectUuid(0, 1, ar.device.deviceID, ar.device.vendorID)
}

func (ar *applicationRelation) remote() *net.UDPAddr {
	return &net.UDPAddr{IP: ar.device.ip, Port: pnioCMPort}
}

func (ar *applicationRelation) call(ctx context.Context, operation readWriteModel.DceRpc_Operation, blocks ...readWriteModel.PnIoCm_Block) (readWriteModel.PnIoCm_Packet_Res, error) {
	callCtx, cancel := context.WithTimeout(ctx, ar.settings.requestTimeout)
	defer cancel()
	sequence := ar.sequence.Add(1) - 1
	return ar.rpc.call(callCtx, ar.remote(), ar.object(), ar.activity, sequence, operation, readWriteModel.NewPnIoCm_Packet_Req(argsMaximum, argsMaximum, 0, blocks))
}

// establish runs the startup: read what the device has, connect, parameterize, and wait for the
// device to report it is ready, with cyclic data running from the connect on.
func (ar *applicationRelation) establish(ctx context.Context) error {
	var err error
	if ar.rpc, err = listenRPC(&net.UDPAddr{IP: ar.localIP, Port: pnioCMPort}, ar.serve, ar.settings.rpcFragmentSize, ar.log); err != nil {
		return err
	}
	real, err := ar.readRealIdentification(ctx)
	if err != nil {
		return err
	}
	if real, err = withConfiguredModules(ar.index, real, ar.modules); err != nil {
		return err
	}
	if ar.plan, err = makePlan(ar.index, real); err != nil {
		return err
	}
	if err := ar.connect(ctx); err != nil {
		return err
	}
	ar.startCyclic()
	if err := ar.parameterize(ctx); err != nil {
		return err
	}
	if _, err := ar.call(ctx, readWriteModel.DceRpc_Operation_CONTROL,
		readWriteModel.NewPnIoCm_Control_Request_ParameterEnd(blockVersionHigh, blockVersionLow, readWriteModel.NewUuid(ar.arUUID), ar.sessionKey, controlCommandPrmEnd)); err != nil {
		return errors.Wrap(err, "the device refused the end of parameterization")
	}
	readyCtx, cancel := context.WithTimeout(ctx, ar.settings.readyTimeout)
	defer cancel()
	select {
	case <-ar.ready:
		ar.log.Info().Str("device", ar.device.name).Dur("cycle", ar.settings.cycle()).Msg("Profinet application relation running")
		return nil
	case err := <-ar.lost:
		return err
	case <-readyCtx.Done():
		return errors.Errorf("the device did not report application ready within %s", ar.settings.readyTimeout)
	}
}

func (ar *applicationRelation) readRealIdentification(ctx context.Context) ([]realSubmodule, error) {
	header := readWriteModel.NewIODReadRequestHeader(blockVersionHigh, blockVersionLow, 0, readWriteModel.NewUuid(zeroUUID), 0, 0, 0,
		realIdentificationDataIndex, argsMaximum, readWriteModel.NewUuid(zeroUUID))
	response, err := ar.call(ctx, readWriteModel.DceRpc_Operation_READ_IMPLICIT, header)
	if err != nil {
		return nil, errors.Wrap(err, "error reading the device's RealIdentificationData")
	}
	for _, block := range response.GetBlocks() {
		if real, ok := block.(readWriteModel.PnIoCm_Block_RealIdentificationData); ok {
			return realSubmodulesOf(real), nil
		}
	}
	return nil, errors.New("the device answered the RealIdentificationData read without one")
}

func (ar *applicationRelation) connect(ctx context.Context) error {
	p := ar.plan
	ioCr := func(crType readWriteModel.PnIoCm_IoCrType, reference uint16, length int, frameID uint16, input bool) readWriteModel.PnIoCm_Block {
		return readWriteModel.NewPnIoCm_Block_IoCrReq(blockVersionHigh, blockVersionLow, crType, reference, rtUDPPort,
			false, false, false, false, readWriteModel.PnIoCm_RtClass_RT_CLASS_2, uint16(length), frameID,
			ar.settings.sendClockFactor, ar.settings.reductionRatio, 1, 0, 0xFFFFFFFF,
			ar.settings.watchdogFactor, ar.settings.dataHoldFactor, 0xC000,
			readWriteModel.NewMacAddress(make([]byte, 6)), []readWriteModel.PnIoCm_IoCrBlockReqApi{p.ioCrApi(input)})
	}
	blocks := []readWriteModel.PnIoCm_Block{
		readWriteModel.NewPnIoCm_Block_ArReq(blockVersionHigh, blockVersionLow, readWriteModel.PnIoCm_ArType_IO_CONTROLLER,
			readWriteModel.NewUuid(ar.arUUID), ar.sessionKey, readWriteModel.NewMacAddress(ar.local),
			readWriteModel.NewDceRpc_ObjectUuid(0, 1, 0x0001, 0x0000),
			false, false, false, false, readWriteModel.PnIoCm_CompanionArType_SINGLE_AR, false, true, false,
			readWriteModel.PnIoCm_State_ACTIVE, 600, rtUDPPort, ar.settings.stationName),
		ioCr(readWriteModel.PnIoCm_IoCrType_INPUT_CR, 1, p.InputLength, 0x8001, true),
		ioCr(readWriteModel.PnIoCm_IoCrType_OUTPUT_CR, 2, p.OutputLength, 0x8000, false),
	}
	blocks = append(blocks, p.expectedSubmoduleBlocks()...)
	blocks = append(blocks, readWriteModel.NewPnIoCm_Block_AlarmCrReq(blockVersionHigh, blockVersionLow, readWriteModel.PnIoCm_AlarmCrType_ALARM_CR,
		rtUDPPort, false, false, rtaTimeoutFactor, rtaRetries, controllerAlarmReference, 200, 0xC000, 0xA000))

	response, err := ar.call(ctx, readWriteModel.DceRpc_Operation_CONNECT, blocks...)
	if err != nil {
		return errors.Wrap(err, "the device refused the connect")
	}
	for _, block := range response.GetBlocks() {
		switch b := block.(type) {
		case readWriteModel.PnIoCm_Block_IoCrRes:
			switch b.GetIoCrType() {
			case readWriteModel.PnIoCm_IoCrType_INPUT_CR:
				ar.inputFrameID = b.GetFrameId()
			case readWriteModel.PnIoCm_IoCrType_OUTPUT_CR:
				ar.outputFrameID = b.GetFrameId()
			}
		case readWriteModel.PnIoCm_Block_AlarmCrRes:
			ar.deviceAlarmReference = b.GetLocalAlarmReference()
		case readWriteModel.PnIoCm_Block_ModuleDiff:
			for _, api := range b.GetApis() {
				for _, module := range api.GetModules() {
					for _, submodule := range module.GetSubmodules() {
						ar.log.Warn().Uint16("slot", module.GetSlotNumber()).Uint16("subslot", submodule.GetSubslotNumber()).
							Stringer("moduleState", module.GetModuleState()).Msg("the device reports a submodule differing from what was expected")
					}
				}
			}
		}
	}
	if ar.inputFrameID == 0 || ar.outputFrameID == 0 {
		return errors.New("the connect response names no frame ids for the cyclic data")
	}
	ar.output = make([]byte, p.OutputLength)
	ar.input = make([]byte, p.InputLength)
	for _, s := range p.Submodules {
		if s.OutputIOPS >= 0 {
			ar.output[s.OutputIOPS] = iopsGood
		}
		if s.OutputIOCS >= 0 {
			ar.output[s.OutputIOCS] = iocsGood
		}
	}
	return nil
}

// parameterize writes every submodule's GSDML parameter records with their default values.
func (ar *applicationRelation) parameterize(ctx context.Context) error {
	sequence := uint16(1)
	for _, s := range ar.plan.Submodules {
		for _, record := range s.Resolved.Submodule.Records {
			data, err := record.DefaultData()
			if err != nil {
				return errors.Wrapf(err, "slot %d subslot 0x%04X", s.Slot, s.Subslot)
			}
			header := readWriteModel.NewIODWriteRequestHeader(blockVersionHigh, blockVersionLow, sequence, readWriteModel.NewUuid(ar.arUUID),
				0, s.Slot, s.Subslot, uint16(record.Index), uint32(len(data)), readWriteModel.NewUserData(data))
			sequence++
			if _, err := ar.call(ctx, readWriteModel.DceRpc_Operation_WRITE, header); err != nil {
				return errors.Wrapf(err, "error writing parameter record %d to slot %d subslot 0x%04X", record.Index, s.Slot, s.Subslot)
			}
		}
	}
	return nil
}

// serve answers the device's calls: application ready, after parameterization.
func (ar *applicationRelation) serve(operation readWriteModel.DceRpc_Operation, request readWriteModel.PnIoCm_Packet_Req) (readWriteModel.PnIoCm_Packet_Res, error) {
	if operation != readWriteModel.DceRpc_Operation_CONTROL {
		return nil, errors.Errorf("unexpected %s call from the device", operation)
	}
	for _, block := range request.GetBlocks() {
		ready, ok := block.(readWriteModel.PnIoCm_Control_Request_ApplicationReady)
		if !ok {
			continue
		}
		arUUID := ready.GetArUuid().GetData()
		if string(arUUID) != string(ar.arUUID) {
			return nil, errors.New("application ready for another application relation")
		}
		ar.readyOne.Do(func() { close(ar.ready) })
		return readWriteModel.NewPnIoCm_Packet_Res(0, 0, 0, 0, argsMaximum, 0, []readWriteModel.PnIoCm_Block{
			readWriteModel.NewPnIoCm_Control_Response_ApplicationReady(blockVersionHigh, blockVersionLow, ready.GetArUuid(), ready.GetSessionKey(), controlCommandDone, 0),
		}), nil
	}
	return nil, errors.New("a control call without an application ready block")
}

// startCyclic sends the output frame every cycle and receives the input frames, watching for the
// device to stop sending them.
func (ar *applicationRelation) startCyclic() {
	ctx, cancel := context.WithCancel(context.Background())
	ar.cancel = cancel
	ar.running.Store(true)
	ar.wg.Go(func() { ar.sendCyclic(ctx) })
	ar.wg.Go(func() { ar.receiveCyclic(ctx) })
	ar.wg.Go(func() { ar.sendAlarmAcks(ctx, ar.alarmHigh) })
	ar.wg.Go(func() { ar.sendAlarmAcks(ctx, ar.alarmLow) })
}

func (ar *applicationRelation) fail(err error) {
	ar.lostOnce.Do(func() {
		ar.running.Store(false)
		ar.lost <- err
	})
}

func (ar *applicationRelation) sendCyclic(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			ar.log.Error().Str("stack", string(debug.Stack())).Interface("err", err).Msg("cyclic sender panic-ed")
		}
	}()
	cycle := ar.settings.cycle()
	step := ar.settings.sendClockFactor * ar.settings.reductionRatio
	frame := make([]byte, 14+2+len(ar.output)+4)
	copy(frame[0:], ar.device.mac)
	copy(frame[6:], ar.local)
	binary.BigEndian.PutUint16(frame[12:], etherTypeProfinet)
	binary.BigEndian.PutUint16(frame[14:], ar.outputFrameID)
	tail := 16 + len(ar.output)
	counter := uint16(0)
	ticker := time.NewTicker(cycle)
	defer ticker.Stop()
	for {
		ar.mutex.Lock()
		copy(frame[16:], ar.output)
		ar.mutex.Unlock()
		binary.BigEndian.PutUint16(frame[tail:], counter)
		frame[tail+2] = dataStatusRun
		frame[tail+3] = 0
		counter += step
		record, err := rawsocket.EncodeRecord(frame)
		if err == nil {
			err = ar.station.Write(ctx, record)
		}
		if err != nil && ctx.Err() == nil {
			ar.log.Debug().Err(err).Msg("error sending the cyclic output frame")
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (ar *applicationRelation) receiveCyclic(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			ar.log.Error().Str("stack", string(debug.Stack())).Interface("err", err).Msg("cyclic receiver panic-ed")
		}
	}()
	timeout := ar.settings.cycle() * time.Duration(ar.settings.watchdogFactor)
	// Until the first input frame the device may take a moment to start its provider.
	deadline := time.Now().Add(max(timeout, 2*time.Second) + ar.settings.requestTimeout)
	for ctx.Err() == nil {
		readCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		frame, err := readFrame(readCtx, ar.station)
		cancel()
		if err == nil {
			if ar.acceptInput(frame) {
				deadline = time.Now().Add(timeout)
				continue
			}
			ar.acceptAlarm(ctx, frame)
		}
		if ctx.Err() == nil && time.Now().After(deadline) {
			ar.fail(errors.Errorf("no cyclic data from %s for %s: the application relation is lost", ar.device.name, timeout))
			return
		}
	}
}

// acceptInput takes an input frame of this relation into the input image.
func (ar *applicationRelation) acceptInput(frame []byte) bool {
	if len(frame) < 18 || string(frame[6:12]) != string(ar.device.mac) {
		return false
	}
	offset := 12
	if binary.BigEndian.Uint16(frame[offset:]) == 0x8100 {
		offset += 4
	}
	if binary.BigEndian.Uint16(frame[offset:]) != etherTypeProfinet {
		return false
	}
	offset += 2
	if binary.BigEndian.Uint16(frame[offset:]) != ar.inputFrameID {
		return false
	}
	data := frame[offset+2:]
	if len(data) < len(ar.input)+4 {
		return false
	}
	ar.mutex.Lock()
	copy(ar.input, data[:len(ar.input)])
	ar.inputStatus = data[len(ar.input)+2]
	ar.lastInput = time.Now()
	ar.mutex.Unlock()
	return true
}

// inputImage is a copy of the latest input data and its data status.
func (ar *applicationRelation) inputImage() ([]byte, byte, time.Time) {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return append([]byte(nil), ar.input...), ar.inputStatus, ar.lastInput
}

// updateOutput changes part of the output image, sent from the next cycle on.
func (ar *applicationRelation) updateOutput(offset int, data []byte, mask []byte) error {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	if offset < 0 || offset+len(data) > len(ar.output) {
		return errors.Errorf("output bytes %d..%d are outside the %d byte output frame", offset, offset+len(data), len(ar.output))
	}
	for i := range data {
		if mask == nil {
			ar.output[offset+i] = data[i]
		} else {
			ar.output[offset+i] = ar.output[offset+i]&^mask[i] | data[i]&mask[i]
		}
	}
	return nil
}

func (ar *applicationRelation) outputImage() []byte {
	ar.mutex.Lock()
	defer ar.mutex.Unlock()
	return append([]byte(nil), ar.output...)
}

// release ends a running application relation with the Release call, so the device frees it at
// once instead of waiting for its watchdog. It is best effort: a device that does not answer within
// the request timeout is left to the watchdog.
func (ar *applicationRelation) release() {
	if ar.rpc == nil || !ar.running.Load() {
		return
	}
	response, err := ar.call(context.Background(), readWriteModel.DceRpc_Operation_RELEASE,
		readWriteModel.NewPnIoCm_Control_Request_Release(blockVersionHigh, blockVersionLow, readWriteModel.NewUuid(ar.arUUID), ar.sessionKey, controlCommandRelease, 0))
	if err != nil {
		ar.log.Debug().Err(err).Str("device", ar.device.name).Msg("Profinet release failed; the device's watchdog ends the relation")
		return
	}
	for _, block := range response.GetBlocks() {
		if done, ok := block.(readWriteModel.PnIoCm_Control_Response_Release); ok && done.GetControlCommand() == controlCommandDone {
			ar.log.Debug().Str("device", ar.device.name).Msg("Profinet application relation released")
			return
		}
	}
	ar.log.Debug().Str("device", ar.device.name).Msg("Profinet release answered without a release block")
}

// close releases the relation if it is still running, then stops the cyclic data and the RPC
// endpoint.
func (ar *applicationRelation) close() {
	ar.release()
	if ar.cancel != nil {
		ar.cancel()
	}
	ar.wg.Wait()
	if ar.rpc != nil {
		ar.rpc.close()
	}
	ar.running.Store(false)
}
