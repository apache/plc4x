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
	"fmt"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

const (
	// maxSegmentedUploadSize bounds what a segmented upload accepts. The size comes from the remote
	// node's initiate response, an unsigned 32 bit field, and is checked before any segment is
	// fetched; object dictionary values are nowhere near this. plc4j uses the same bound.
	maxSegmentedUploadSize = 16 * 1024 * 1024
	// expeditedLimit is the most data an expedited transfer carries in its initiate frame.
	expeditedLimit = 4
	// segmentDataLength is the data one segment carries.
	segmentDataLength = 7
)

// Abort codes this driver sends (CiA 301, 7.2.4.3.17).
const (
	abortToggleBit    = uint32(0x05030000)
	abortTimeout      = uint32(0x05040000)
	abortGeneralError = uint32(0x08000000)
	abortOutOfMemory  = uint32(0x05040005)
)

// abortDescriptions are the CiA 301 abort codes, so an abort says what went wrong.
var abortDescriptions = map[uint32]string{
	0x05030000: "toggle bit not alternated",
	0x05040000: "SDO protocol timed out",
	0x05040001: "client/server command specifier not valid or unknown",
	0x05040002: "invalid block size",
	0x05040003: "invalid sequence number",
	0x05040004: "CRC error",
	0x05040005: "out of memory",
	0x06010000: "unsupported access to an object",
	0x06010001: "attempt to read a write only object",
	0x06010002: "attempt to write a read only object",
	0x06020000: "object does not exist in the object dictionary",
	0x06040041: "object cannot be mapped to the PDO",
	0x06040042: "the number and length of the objects to be mapped would exceed PDO length",
	0x06040043: "general parameter incompatibility",
	0x06040047: "general internal incompatibility in the device",
	0x06060000: "access failed due to a hardware error",
	0x06070010: "data type does not match, length of service parameter does not match",
	0x06070012: "data type does not match, length of service parameter too high",
	0x06070013: "data type does not match, length of service parameter too low",
	0x06090011: "sub-index does not exist",
	0x06090030: "invalid value for parameter",
	0x06090031: "value of parameter written too high",
	0x06090032: "value of parameter written too low",
	0x06090036: "maximum value is less than minimum value",
	0x060A0023: "resource not available: SDO connection",
	0x08000000: "general error",
	0x08000020: "data cannot be transferred or stored to the application",
	0x08000021: "data cannot be transferred or stored to the application because of local control",
	0x08000022: "data cannot be transferred or stored to the application because of the present device state",
	0x08000023: "object dictionary dynamic generation fails or no object dictionary is present",
	0x08000024: "no data available",
}

// AbortError is an SDO transfer the remote node aborted.
type AbortError struct {
	Node     uint8
	Index    uint16
	SubIndex uint8
	Code     uint32
}

func (e *AbortError) Error() string {
	description, ok := abortDescriptions[e.Code]
	if !ok {
		description = "unknown abort code"
	}
	return fmt.Sprintf("node %d aborted the SDO transfer of 0x%04X/%d with 0x%08X (%s)", e.Node, e.Index, e.SubIndex, e.Code, description)
}

// sdoTransfer is one SDO upload or download between this connection and one node.
type sdoTransfer struct {
	codec      *MessageCodec
	node       uint8
	answerNode uint8
	address    readWriteModel.IndexAddress
	// stepTimeout bounds the wait for each answer.
	stepTimeout time.Duration
}

func newSDOTransfer(codec *MessageCodec, tag Tag, stepTimeout time.Duration) sdoTransfer {
	return sdoTransfer{
		codec:       codec,
		node:        tag.Node,
		answerNode:  tag.ResponderNode(),
		address:     readWriteModel.NewIndexAddress(tag.Index, tag.SubIndex),
		stepTimeout: stepTimeout,
	}
}

// exchange sends one request and waits for the answering node's next SDO response. An abort
// answer comes back as an *AbortError.
func (s sdoTransfer) exchange(ctx context.Context, command readWriteModel.SDORequestCommand, request readWriteModel.SDORequest) (readWriteModel.SDOResponse, error) {
	answerNode := s.answerNode
	w := s.codec.expect(func(frame readWriteModel.CANOpenFrame) bool {
		if frame.GetNodeId() != answerNode || frame.GetService() != readWriteModel.CANOpenService_TRANSMIT_SDO {
			return false
		}
		payload, ok := frame.GetPayload().(readWriteModel.CANOpenSDOResponse)
		if !ok {
			return false
		}
		// An initiate answer for another object is not ours: it is a late answer to an earlier
		// transfer, or one meant for another client. Skipping it is what the protocol's driver test
		// suite expects. An abort counts whatever object it names, since it ends whatever transfer
		// the server had open with us.
		switch response := payload.GetResponse().(type) {
		case readWriteModel.SDOInitiateUploadResponse:
			return s.sameAddress(response.GetAddress())
		case readWriteModel.SDOInitiateDownloadResponse:
			return s.sameAddress(response.GetAddress())
		}
		return true
	})
	frame := readWriteModel.NewCANOpenFrame(s.node, readWriteModel.CANOpenService_RECEIVE_SDO, readWriteModel.NewCANOpenSDORequest(command, request))
	stepCtx, cancel := context.WithTimeout(ctx, s.stepTimeout)
	defer cancel()
	if err := s.codec.sendFrame(stepCtx, frame); err != nil {
		w.cancel()
		return nil, err
	}
	answer, err := w.wait(stepCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			// Tell the server, so it doesn't keep the transfer open until its own timeout.
			s.abort(abortTimeout)
			return nil, errors.Wrapf(err, "node %d did not answer the SDO request for 0x%04X/%d", answerNode, s.address.GetIndex(), s.address.GetSubindex())
		}
		return nil, err
	}
	response := answer.GetPayload().(readWriteModel.CANOpenSDOResponse).GetResponse()
	if abort, ok := response.(readWriteModel.SDOAbortResponse); ok {
		return nil, &AbortError{
			Node:     answerNode,
			Index:    s.address.GetIndex(),
			SubIndex: s.address.GetSubindex(),
			Code:     abort.GetAbort().GetCode(),
		}
	}
	return response, nil
}

// abort ends the transfer from this side. Best effort: it is sent when something already went wrong.
func (s sdoTransfer) abort(code uint32) {
	ctx, cancel := context.WithTimeout(context.Background(), s.stepTimeout)
	defer cancel()
	frame := readWriteModel.NewCANOpenFrame(s.node, readWriteModel.CANOpenService_RECEIVE_SDO,
		readWriteModel.NewCANOpenSDORequest(readWriteModel.SDORequestCommand_ABORT,
			readWriteModel.NewSDOAbortRequest(readWriteModel.NewSDOAbort(s.address, code))))
	if err := s.codec.sendFrame(ctx, frame); err != nil {
		s.codec.log.Debug().Err(err).Msg("error sending an SDO abort")
	}
}

func (s sdoTransfer) sameAddress(answered readWriteModel.IndexAddress) bool {
	return answered.GetIndex() == s.address.GetIndex() && answered.GetSubindex() == s.address.GetSubindex()
}

func (s sdoTransfer) checkAddress(answered readWriteModel.IndexAddress) error {
	if !s.sameAddress(answered) {
		s.abort(abortGeneralError)
		return errors.Errorf("node %d answered for 0x%04X/%d instead of 0x%04X/%d", s.answerNode,
			answered.GetIndex(), answered.GetSubindex(), s.address.GetIndex(), s.address.GetSubindex())
	}
	return nil
}

// upload reads the object's bytes: expedited when the server sends them in its initiate answer,
// segmented otherwise. Ported from plc4j's SDOUploadConversation.
func (s sdoTransfer) upload(ctx context.Context) ([]byte, error) {
	response, err := s.exchange(ctx, readWriteModel.SDORequestCommand_INITIATE_UPLOAD, readWriteModel.NewSDOInitiateUploadRequest(s.address))
	if err != nil {
		return nil, err
	}
	initiate, ok := response.(readWriteModel.SDOInitiateUploadResponse)
	if !ok {
		s.abort(abortGeneralError)
		return nil, errors.Errorf("unexpected answer %T to an SDO upload request", response)
	}
	if err := s.checkAddress(initiate.GetAddress()); err != nil {
		return nil, err
	}
	switch payload := initiate.GetPayload().(type) {
	case readWriteModel.SDOInitiateExpeditedUploadResponse:
		return payload.GetData(), nil
	case readWriteModel.SDOInitiateSegmentedUploadResponse:
		size := payload.GetBytes()
		if size > maxSegmentedUploadSize {
			s.abort(abortOutOfMemory)
			return nil, errors.Errorf("node %d announced an implausible SDO upload of %d bytes", s.answerNode, size)
		}
		return s.uploadSegments(ctx, int(size))
	}
	s.abort(abortGeneralError)
	return nil, errors.Errorf("unsupported SDO upload kind %T", initiate.GetPayload())
}

func (s sdoTransfer) uploadSegments(ctx context.Context, size int) ([]byte, error) {
	data := make([]byte, 0, size)
	toggle := false
	for {
		response, err := s.exchange(ctx, readWriteModel.SDORequestCommand_SEGMENT_UPLOAD, readWriteModel.NewSDOSegmentUploadRequest(toggle))
		if err != nil {
			return nil, err
		}
		segment, ok := response.(readWriteModel.SDOSegmentUploadResponse)
		if !ok {
			s.abort(abortGeneralError)
			return nil, errors.Errorf("unexpected answer %T to an SDO segment upload request", response)
		}
		if segment.GetToggle() != toggle {
			s.abort(abortToggleBit)
			return nil, errors.Errorf("node %d did not alternate the toggle bit", s.answerNode)
		}
		// Bytes past the announced size are dropped rather than failing the transfer, the way plc4j
		// does; some servers pad their last segment. The announced size is what the value is.
		if room := size - len(data); room > 0 {
			chunk := segment.GetData()
			if len(chunk) > room {
				s.codec.log.Debug().Int("announced", size).Msg("SDO server sent more than it announced, truncating")
				chunk = chunk[:room]
			}
			data = append(data, chunk...)
		}
		if segment.GetLast() {
			return data, nil
		}
		toggle = !toggle
	}
}

// download writes the object's bytes: expedited for up to 4 bytes, segmented above. Ported from
// plc4j's SDODownloadConversation.
func (s sdoTransfer) download(ctx context.Context, data []byte) error {
	if len(data) <= expeditedLimit {
		response, err := s.exchange(ctx, readWriteModel.SDORequestCommand_INITIATE_DOWNLOAD,
			readWriteModel.NewSDOInitiateDownloadRequest(true, true, s.address, readWriteModel.NewSDOInitiateExpeditedUploadResponse(data)))
		if err != nil {
			return err
		}
		initiate, ok := response.(readWriteModel.SDOInitiateDownloadResponse)
		if !ok {
			s.abort(abortGeneralError)
			return errors.Errorf("unexpected answer %T to an SDO download request", response)
		}
		return s.checkAddress(initiate.GetAddress())
	}

	response, err := s.exchange(ctx, readWriteModel.SDORequestCommand_INITIATE_DOWNLOAD,
		readWriteModel.NewSDOInitiateDownloadRequest(false, true, s.address, readWriteModel.NewSDOInitiateSegmentedUploadResponse(uint32(len(data)))))
	if err != nil {
		return err
	}
	initiate, ok := response.(readWriteModel.SDOInitiateDownloadResponse)
	if !ok {
		s.abort(abortGeneralError)
		return errors.Errorf("unexpected answer %T to an SDO download request", response)
	}
	if err := s.checkAddress(initiate.GetAddress()); err != nil {
		return err
	}
	toggle := false
	for offset := 0; offset < len(data); offset += segmentDataLength {
		end := min(offset+segmentDataLength, len(data))
		response, err := s.exchange(ctx, readWriteModel.SDORequestCommand_SEGMENT_DOWNLOAD,
			readWriteModel.NewSDOSegmentDownloadRequest(toggle, end == len(data), data[offset:end]))
		if err != nil {
			return err
		}
		segment, ok := response.(readWriteModel.SDOSegmentDownloadResponse)
		if !ok {
			s.abort(abortGeneralError)
			return errors.Errorf("unexpected answer %T to an SDO segment download", response)
		}
		if segment.GetToggle() != toggle {
			s.abort(abortToggleBit)
			return errors.Errorf("node %d did not alternate the toggle bit", s.answerNode)
		}
		toggle = !toggle
	}
	return nil
}

// dataItemSize is the size argument the generated DataItem code expects for a value of byteLength
// bytes. The spec sizes strings in bits (and UNICODE_STRING divides by 8 once more), records in
// bytes, and ignores the argument for the fixed size types.
func dataItemSize(dataType readWriteModel.CANOpenDataType, byteLength int) int32 {
	switch dataType {
	case readWriteModel.CANOpenDataType_VISIBLE_STRING, readWriteModel.CANOpenDataType_OCTET_STRING:
		return int32(byteLength * 8)
	case readWriteModel.CANOpenDataType_UNICODE_STRING:
		return int32(byteLength * 64)
	}
	return int32(byteLength)
}

// isVariableLength reports whether a type's size comes from the data rather than from the type.
func isVariableLength(dataType readWriteModel.CANOpenDataType) bool {
	switch dataType {
	case readWriteModel.CANOpenDataType_VISIBLE_STRING, readWriteModel.CANOpenDataType_OCTET_STRING,
		readWriteModel.CANOpenDataType_UNICODE_STRING, readWriteModel.CANOpenDataType_RECORD:
		return true
	}
	return false
}

// decodeValue turns object or PDO bytes into a value. Fixed size types read their own width from
// the front of the data, so a PDO padded to 8 bytes decodes the same as an exact one.
func decodeValue(ctx context.Context, dataType readWriteModel.CANOpenDataType, data []byte) (apiValues.PlcValue, error) {
	if !isVariableLength(dataType) {
		width := (int(dataType.NumBits()) + 7) / 8
		if len(data) < width {
			return nil, errors.Errorf("%s needs %d bytes, got %d", dataType, width, len(data))
		}
	}
	return readWriteModel.DataItemParse(ctx, data, dataType, dataItemSize(dataType, len(data)))
}

// encodeValue turns a value into object or PDO bytes.
func encodeValue(dataType readWriteModel.CANOpenDataType, value apiValues.PlcValue) ([]byte, error) {
	byteLength := 0
	switch dataType {
	case readWriteModel.CANOpenDataType_VISIBLE_STRING, readWriteModel.CANOpenDataType_OCTET_STRING:
		byteLength = len(value.GetString())
	case readWriteModel.CANOpenDataType_UNICODE_STRING:
		byteLength = len(value.GetString())
	case readWriteModel.CANOpenDataType_RECORD:
		byteLength = int(value.GetLength())
	}
	return readWriteModel.DataItemSerialize(value, dataType, dataItemSize(dataType, byteLength))
}
