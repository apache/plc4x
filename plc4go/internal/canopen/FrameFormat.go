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

package canopen

import (
	"context"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transports/can"
)

// frameFormat is how frames are laid out on the transport's byte stream.
type frameFormat interface {
	recordSize() uint32
	// decode turns one record into a frame. ok is false for a frame that is valid CAN but nothing
	// this driver handles (SYNC, EMCY, extended or remote frames).
	decode(ctx context.Context, record []byte) (frame readWriteModel.CANOpenFrame, ok bool, err error)
	encode(frame readWriteModel.CANOpenFrame) ([]byte, error)
}

// cobRange is one row of the CANopen predefined connection set: which COB-IDs belong to a service,
// and which COB-ID is node 0 of it.
//
// The table exists because the generated CANOpenService enum can't be used for this: SYNC and EMCY
// share the 4 bit function code 0b0001, so in Go they are the same constant and the enum's Min()
// and Max() can only ever answer for one of them.
type cobRange struct {
	service  readWriteModel.CANOpenService
	min, max uint32
	base     uint32
	// hasPayload is whether the protocol spec has a payload type for the service.
	hasPayload bool
}

var cobRanges = []cobRange{
	{readWriteModel.CANOpenService_NMT, 0x000, 0x000, 0x000, true},
	{readWriteModel.CANOpenService_SYNC, 0x080, 0x080, 0x080, false},
	{readWriteModel.CANOpenService_EMCY, 0x081, 0x0FF, 0x080, false},
	{readWriteModel.CANOpenService_TIME, 0x100, 0x100, 0x100, true},
	{readWriteModel.CANOpenService_TRANSMIT_PDO_1, 0x181, 0x1FF, 0x180, true},
	{readWriteModel.CANOpenService_RECEIVE_PDO_1, 0x201, 0x27F, 0x200, true},
	{readWriteModel.CANOpenService_TRANSMIT_PDO_2, 0x281, 0x2FF, 0x280, true},
	{readWriteModel.CANOpenService_RECEIVE_PDO_2, 0x301, 0x37F, 0x300, true},
	{readWriteModel.CANOpenService_TRANSMIT_PDO_3, 0x381, 0x3FF, 0x380, true},
	{readWriteModel.CANOpenService_RECEIVE_PDO_3, 0x401, 0x47F, 0x400, true},
	{readWriteModel.CANOpenService_TRANSMIT_PDO_4, 0x481, 0x4FF, 0x480, true},
	{readWriteModel.CANOpenService_RECEIVE_PDO_4, 0x501, 0x57F, 0x500, true},
	{readWriteModel.CANOpenService_TRANSMIT_SDO, 0x581, 0x5FF, 0x580, true},
	{readWriteModel.CANOpenService_RECEIVE_SDO, 0x601, 0x67F, 0x600, true},
	{readWriteModel.CANOpenService_HEARTBEAT, 0x701, 0x77F, 0x700, true},
}

// resolveCOBID splits a standard identifier into service and node.
func resolveCOBID(id uint32) (cobRange, uint8, bool) {
	for _, row := range cobRanges {
		if id >= row.min && id <= row.max {
			return row, uint8(id - row.base), true
		}
	}
	return cobRange{}, 0, false
}

// cobID is the identifier a frame of the service goes out under. The services SYNC and EMCY can't
// be told apart by their enum value, and this driver never sends either.
func cobID(service readWriteModel.CANOpenService, node uint8) (uint32, error) {
	if service == readWriteModel.CANOpenService_SYNC {
		return 0, errors.New("sending SYNC or EMCY is not supported")
	}
	for _, row := range cobRanges {
		if row.service == service {
			if row.min == row.max {
				return row.base, nil
			}
			if node == 0 || node > maxNodeID {
				return 0, errors.Errorf("node id %d is not in 1..%d", node, maxNodeID)
			}
			return row.base + uint32(node), nil
		}
	}
	return 0, errors.Errorf("unknown CANopen service %s", service)
}

// rawFrameFormat is the 16 byte record of package can, which the SocketCAN and virtual CAN
// transports deliver: the node and service come from the CAN identifier, the payload is the data.
type rawFrameFormat struct{}

func (rawFrameFormat) recordSize() uint32 {
	return can.FrameSize
}

func (rawFrameFormat) decode(ctx context.Context, record []byte) (readWriteModel.CANOpenFrame, bool, error) {
	canFrame, err := can.Decode(record)
	if err != nil {
		return nil, false, err
	}
	if canFrame.Extended || canFrame.Remote || canFrame.Error {
		// CANopen uses 11 bit identifiers; the rest is somebody else's traffic.
		return nil, false, nil
	}
	row, node, ok := resolveCOBID(canFrame.ID)
	if !ok || !row.hasPayload {
		return nil, false, nil
	}
	data := canFrame.Data
	if row.service.Pdo() && len(data) < can.MaxDataLength {
		// The spec reads a PDO as 8 bytes, but a PDO is as long as its mapping. Padding keeps the
		// mapped bytes at the front, which is where a tag decodes them from.
		data = append(append([]byte(nil), data...), make([]byte, can.MaxDataLength-len(data))...)
	}
	payload, err := readWriteModel.CANOpenPayloadParse[readWriteModel.CANOpenPayload](ctx, data, row.service)
	if err != nil {
		return nil, false, errors.Wrapf(err, "error parsing the %s payload from 0x%03X", row.service, canFrame.ID)
	}
	return readWriteModel.NewCANOpenFrame(node, row.service, payload), true, nil
}

func (rawFrameFormat) encode(frame readWriteModel.CANOpenFrame) ([]byte, error) {
	id, err := cobID(frame.GetService(), frame.GetNodeId())
	if err != nil {
		return nil, err
	}
	data, err := frame.GetPayload().Serialize()
	if err != nil {
		return nil, errors.Wrap(err, "error serializing the payload")
	}
	return can.Frame{ID: id, Data: data}.Encode()
}

// compactFrameFormat is the spec's own CANOpenFrame: node id byte, service nibble, payload padded to
// 8 bytes. No real transport speaks it, but the protocol's driver test suites are written in it, so
// the driver uses it on the test transport.
type compactFrameFormat struct{}

// compactRecordSize is the node id byte, the service byte and the 8 padded payload bytes.
const compactRecordSize = 10

func (compactFrameFormat) recordSize() uint32 {
	return compactRecordSize
}

func (compactFrameFormat) decode(ctx context.Context, record []byte) (readWriteModel.CANOpenFrame, bool, error) {
	frame, err := readWriteModel.CANOpenFrameParse(ctx, record)
	if err != nil {
		return nil, false, err
	}
	return frame, true, nil
}

func (compactFrameFormat) encode(frame readWriteModel.CANOpenFrame) ([]byte, error) {
	return frame.Serialize()
}
