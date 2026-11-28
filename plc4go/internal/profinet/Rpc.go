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
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/profinet/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

const (
	// pnioCMPort is the UDP port of the Profinet context manager (PNIO-CM), on device and controller.
	pnioCMPort = 34964
	// rpcRetransmit is how often an unanswered request is sent again.
	rpcRetransmit = time.Second
	// rpcFackWait is how long a fragment waits for the receiver's FACK before the next goes out.
	// A device that sends none just makes each fragment wait this long.
	rpcFackWait = 200 * time.Millisecond
)

// PNIOStatusError is a Profinet error answer: the four PNIOStatus bytes of a response.
type PNIOStatusError struct {
	Operation                                      readWriteModel.DceRpc_Operation
	ErrorCode, ErrorDecode, ErrorCode1, ErrorCode2 uint8
}

func (e *PNIOStatusError) Error() string {
	return fmt.Sprintf("device refused %s: PNIOStatus %02X %02X %02X %02X (error code, decode, code 1, code 2)",
		e.Operation, e.ErrorCode, e.ErrorDecode, e.ErrorCode1, e.ErrorCode2)
}

// rpcHandler answers a request the device sends us (application ready, for one).
type rpcHandler func(operation readWriteModel.DceRpc_Operation, request readWriteModel.PnIoCm_Packet_Req) (readWriteModel.PnIoCm_Packet_Res, error)

// rpcEndpoint is connectionless DCE/RPC over the kernel's UDP, on the PNIO-CM port: it calls the
// device and serves the calls the device makes back. One socket does both, because a device
// answers to the port a request came from and calls the controller on the PNIO-CM port.
type rpcEndpoint struct {
	conn    *net.UDPConn
	handler rpcHandler

	// fragmentSize is the most body one datagram carries; longer packets go out as fragments.
	fragmentSize int

	mutex   sync.Mutex
	pending map[callKey]chan readWriteModel.DceRpc_Packet
	// facks gets the fragment numbers the receiver acknowledges, per call sent as fragments.
	facks map[callKey]chan uint16
	// partial collects incoming fragments per call.
	partial map[callKey]*reassembly
	// answered keeps the last response per incoming activity, to resend for a retransmitted request.
	answered map[[16]byte][][]byte

	wg   sync.WaitGroup
	done chan struct{}

	log zerolog.Logger
}

// listenRPC opens the endpoint. fragmentSize is the most body one datagram carries; 0 is the
// default, which fills the UDP payload devices take.
func listenRPC(local *net.UDPAddr, handler rpcHandler, fragmentSize int, log zerolog.Logger) (*rpcEndpoint, error) {
	if fragmentSize <= 0 {
		fragmentSize = defaultRPCFragmentBody
	}
	conn, err := net.ListenUDP("udp4", local)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening the PNIO-CM port %s", local)
	}
	endpoint := &rpcEndpoint{
		conn:         conn,
		handler:      handler,
		fragmentSize: fragmentSize,
		pending:      map[callKey]chan readWriteModel.DceRpc_Packet{},
		facks:        map[callKey]chan uint16{},
		partial:      map[callKey]*reassembly{},
		answered:     map[[16]byte][][]byte{},
		done:         make(chan struct{}),
		log:          log,
	}
	endpoint.wg.Go(endpoint.receive)
	return endpoint, nil
}

func (e *rpcEndpoint) close() {
	select {
	case <-e.done:
		return
	default:
	}
	close(e.done)
	_ = e.conn.Close()
	e.wg.Wait()
}

// callKey is one call: its activity and its sequence number within it.
type callKey struct {
	activity [16]byte
	sequence uint32
}

func activityKey(uuid readWriteModel.DceRpc_ActivityUuid) [16]byte {
	var key [16]byte
	binary.BigEndian.PutUint32(key[0:], uuid.GetData1())
	binary.BigEndian.PutUint16(key[4:], uuid.GetData2())
	binary.BigEndian.PutUint16(key[6:], uuid.GetData3())
	copy(key[8:], uuid.GetData4())
	return key
}

func newActivity() readWriteModel.DceRpc_ActivityUuid {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = b[6]&0x0F | 0x40 // version 4
	b[8] = b[8]&0x3F | 0x80 // RFC 4122 variant
	return readWriteModel.NewDceRpc_ActivityUuid(binary.BigEndian.Uint32(b[0:]), binary.BigEndian.Uint16(b[4:]), binary.BigEndian.Uint16(b[6:]), b[8:16])
}

// call sends a request to the device and waits for its response, sending it again every second
// while unanswered. A response with a non-zero PNIOStatus comes back as *PNIOStatusError.
//
// All calls of one application relation go out under one activity with rising sequence numbers,
// the way controllers do: a device ties its RPC session to the AR by the activity, and a call on a
// fresh activity reaches it without the AR (p-net then never raises the end of parameterization to
// its application, so it never reports ready).
func (e *rpcEndpoint) call(ctx context.Context, remote *net.UDPAddr, object readWriteModel.DceRpc_ObjectUuid, activity readWriteModel.DceRpc_ActivityUuid, sequence uint32, operation readWriteModel.DceRpc_Operation, request readWriteModel.PnIoCm_Packet_Req) (readWriteModel.PnIoCm_Packet_Res, error) {
	packet := readWriteModel.NewDceRpc_Packet(readWriteModel.DceRpc_PacketType_REQUEST, true, false, false,
		readWriteModel.IntegerEncoding_BIG_ENDIAN, readWriteModel.CharacterEncoding_ASCII, readWriteModel.FloatingPointEncoding_IEEE,
		object, readWriteModel.NewDceRpc_InterfaceUuid_DeviceInterface(), activity, 0, sequence, operation, 0, request)
	data, err := packet.Serialize()
	if err != nil {
		return nil, errors.Wrapf(err, "error serializing the %s request", operation)
	}
	fragments := fragmentPacket(data, e.fragmentSize, true)
	key := callKey{activityKey(activity), sequence}
	answers := make(chan readWriteModel.DceRpc_Packet, 4)
	facks := make(chan uint16, 16)
	e.mutex.Lock()
	e.pending[key] = answers
	if len(fragments) > 1 {
		e.facks[key] = facks
	}
	e.mutex.Unlock()
	defer func() {
		e.mutex.Lock()
		delete(e.pending, key)
		delete(e.facks, key)
		e.mutex.Unlock()
	}()

	// send sends the request, a fragmented one a fragment at a time, each after the receiver
	// acknowledged the one before or rpcFackWait passed.
	send := func() error {
		for i, fragment := range fragments {
			if _, err := e.conn.WriteToUDP(fragment, remote); err != nil {
				return err
			}
			if i == len(fragments)-1 {
				break
			}
			wait := time.NewTimer(rpcFackWait)
		waiting:
			for {
				select {
				case <-ctx.Done():
					wait.Stop()
					return ctx.Err()
				case <-wait.C:
					e.log.Debug().Int("fragment", i).Stringer("operation", operation).Msg("no FACK for a request fragment; sending the next")
					break waiting
				case number := <-facks:
					if int(number) >= i {
						wait.Stop()
						break waiting
					}
				}
			}
		}
		return nil
	}
	if err := send(); err != nil {
		return nil, errors.Wrapf(err, "error sending the %s request to %s", operation, remote)
	}
	retransmit := time.NewTicker(rpcRetransmit)
	defer retransmit.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.Wrapf(ctx.Err(), "no answer from %s to %s", remote, operation)
		case <-e.done:
			return nil, errors.New("the PNIO-CM endpoint is closed")
		case <-retransmit.C:
			if err := send(); err != nil {
				return nil, errors.Wrapf(err, "error sending the %s request again", operation)
			}
		case answer := <-answers:
			switch payload := answer.GetPayload().(type) {
			case readWriteModel.PnIoCm_Packet_Res:
				if payload.GetErrorCode() != 0 || payload.GetErrorDecode() != 0 || payload.GetErrorCode1() != 0 || payload.GetErrorCode2() != 0 {
					return payload, &PNIOStatusError{Operation: operation, ErrorCode: payload.GetErrorCode(), ErrorDecode: payload.GetErrorDecode(), ErrorCode1: payload.GetErrorCode1(), ErrorCode2: payload.GetErrorCode2()}
				}
				return payload, nil
			case readWriteModel.PnIoCm_Packet_Fault:
				return nil, errors.Errorf("%s answered %s with a fault, status 0x%08X", remote, operation, payload.GetStatus())
			case readWriteModel.PnIoCm_Packet_Rej:
				return nil, errors.Errorf("%s rejected %s, status 0x%08X", remote, operation, payload.GetStatus())
			case readWriteModel.PnIoCm_Packet_Working:
				retransmit.Reset(rpcRetransmit) // the device has it and is busy: wait, don't resend yet
			case readWriteModel.PnIoCm_Packet_NoCall:
				if err := send(); err != nil {
					return nil, err
				}
			}
		}
	}
}

func (e *rpcEndpoint) receive() {
	defer func() {
		if err := recover(); err != nil {
			e.log.Error().Str("stack", string(debug.Stack())).Interface("err", err).Msg("PNIO-CM receive loop panic-ed")
		}
	}()
	buffer := make([]byte, 65536)
	for {
		n, from, err := e.conn.ReadFromUDP(buffer)
		if err != nil {
			select {
			case <-e.done:
				return
			default:
			}
			e.log.Debug().Err(err).Msg("PNIO-CM read error")
			continue
		}
		data := append([]byte(nil), buffer[:n]...)
		if h, ok := parseRPCHeader(data); ok {
			if h.packetType == rpcPacketTypeFack {
				e.mutex.Lock()
				facks, waiting := e.facks[h.key()]
				e.mutex.Unlock()
				if waiting {
					select {
					case facks <- h.fragment:
					default:
					}
				}
				continue
			}
			if h.isFragment() {
				whole, complete := e.reassemble(h, data, from)
				if !complete {
					continue
				}
				data = whole
			}
		}
		packet, err := readWriteModel.DceRpc_PacketParse(context.Background(), data)
		if err != nil {
			e.log.Debug().Err(err).Stringer("from", from).Msg("not a DCE/RPC packet the model reads")
			continue
		}
		key := activityKey(packet.GetActivityUuid())
		if packet.GetPacketType() == readWriteModel.DceRpc_PacketType_REQUEST || packet.GetPacketType() == readWriteModel.DceRpc_PacketType_PING {
			e.serve(packet, key, from)
			continue
		}
		e.mutex.Lock()
		answers, ok := e.pending[callKey{key, packet.GetSequenceNumber()}]
		e.mutex.Unlock()
		if ok {
			select {
			case answers <- packet:
			default:
			}
		}
	}
}

// reassemble takes an incoming fragment, acknowledging it when the sender asks for that; complete
// is true with the whole packet once its last fragment has arrived.
func (e *rpcEndpoint) reassemble(h rpcHeader, fragment []byte, from *net.UDPAddr) ([]byte, bool) {
	e.mutex.Lock()
	r, ok := e.partial[h.key()]
	if !ok {
		if len(e.partial) > 64 {
			e.partial = map[callKey]*reassembly{}
		}
		r = &reassembly{}
		e.partial[h.key()] = r
	}
	whole, accepted, done := r.add(h, fragment)
	if done {
		delete(e.partial, h.key())
	}
	e.mutex.Unlock()
	if accepted && !done && h.flags&rpcFlagNoFack == 0 {
		if _, err := e.conn.WriteToUDP(fackFor(fragment), from); err != nil {
			e.log.Debug().Err(err).Msg("error sending a FACK")
		}
	}
	return whole, done
}

// serve answers a call from the device.
func (e *rpcEndpoint) serve(packet readWriteModel.DceRpc_Packet, key [16]byte, from *net.UDPAddr) {
	e.mutex.Lock()
	cached, seen := e.answered[key]
	e.mutex.Unlock()
	if packet.GetPacketType() == readWriteModel.DceRpc_PacketType_PING {
		if seen {
			e.write(cached, from)
		} else {
			e.reply(packet, from, readWriteModel.DceRpc_PacketType_WORKING, readWriteModel.NewPnIoCm_Packet_Working())
		}
		return
	}
	if seen {
		e.write(cached, from)
		return
	}
	request, ok := packet.GetPayload().(readWriteModel.PnIoCm_Packet_Req)
	if !ok || e.handler == nil {
		return
	}
	response, err := e.handler(packet.GetOperation(), request)
	if err != nil {
		e.log.Warn().Err(err).Stringer("operation", packet.GetOperation()).Msg("refusing a call from the device")
		e.reply(packet, from, readWriteModel.DceRpc_PacketType_REJECT, readWriteModel.NewPnIoCm_Packet_Rej(0x1C010003)) // nca_op_rng_error
		return
	}
	data := e.reply(packet, from, readWriteModel.DceRpc_PacketType_RESPONSE, response)
	if data != nil {
		e.mutex.Lock()
		if len(e.answered) > 64 {
			e.answered = map[[16]byte][][]byte{}
		}
		e.answered[key] = data
		e.mutex.Unlock()
	}
}

// write sends datagrams, the fragments of one packet.
func (e *rpcEndpoint) write(datagrams [][]byte, to *net.UDPAddr) {
	for _, datagram := range datagrams {
		if _, err := e.conn.WriteToUDP(datagram, to); err != nil {
			e.log.Debug().Err(err).Msg("error sending a response to the device")
			return
		}
	}
}

// reply answers a call from the device. A response too long for one datagram goes out as
// fragments that ask for no FACK: this runs on the receive loop, which cannot wait for them.
func (e *rpcEndpoint) reply(request readWriteModel.DceRpc_Packet, to *net.UDPAddr, packetType readWriteModel.DceRpc_PacketType, payload readWriteModel.PnIoCm_Packet) [][]byte {
	response := readWriteModel.NewDceRpc_Packet(packetType, false, true, false,
		readWriteModel.IntegerEncoding_BIG_ENDIAN, readWriteModel.CharacterEncoding_ASCII, readWriteModel.FloatingPointEncoding_IEEE,
		request.GetObjectUuid(), request.GetInterfaceUuid(), request.GetActivityUuid(), 0, request.GetSequenceNumber(), request.GetOperation(), 0, payload)
	data, err := response.Serialize()
	if err != nil {
		e.log.Warn().Err(err).Msg("error serializing a response to the device")
		return nil
	}
	datagrams := fragmentPacket(data, e.fragmentSize, false)
	e.write(datagrams, to)
	return datagrams
}
