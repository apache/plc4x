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

package umas

import (
	"bytes"
	"encoding/binary"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// The three data-dictionary payloads are parsed here by hand instead of through the generated
// UmasPDUReadUnlocatedVariableNamesResponse / UmasPDUReadDatatypeNamesResponse /
// UmasPDUReadUmasUDTDefinitionResponse parsers.
//
// Why: each of their record types ends in a NUL terminated name, which the mspec spells as
//
//	[manual vstring value 'STATIC_CALL("parseTerminatedString", readBuffer, -1)' ...]
//
// where -1 selects the helper's variable length mode ("read until the terminator"). The Go code
// generator drops the unary minus and emits ParseTerminatedString(ctx, readBuffer, 1), which is
// fixed width mode over a single byte: every name would come back as at most one character and the
// record stream would desynchronize on the very first entry. plc4j's generated code passes -(1) and
// parses these payloads correctly, so this is a Go-side generator defect, not a protocol question.
// Since the generated model may not be edited, the driver reads the records itself and hands them to
// the generated constructors, so everything downstream keeps working against the model types.
//
// The layouts below are the mspec's, read little endian - which is also the byte order plc4j's
// UmasConnection passes to the ReadBufferByteBased it parses these payloads with.

const (
	// unlocatedVariableNamesHeaderSize is range(1) + nextAddress(2) + unknown1(2) + noOfRecords(2).
	unlocatedVariableNamesHeaderSize = 7
	// datatypeNamesHeaderSize is range(1) + nextAddress(2) + unknown1(1) + noOfRecords(2).
	datatypeNamesHeaderSize = 6
	// udtDefinitionHeaderSize is range(1) + unknown1(4) + noOfRecords(2).
	udtDefinitionHeaderSize = 7

	// unlocatedVariableReferenceFixedSize is dataType(2) + block(2) + offset(4) + flags(1) +
	// unknown4(1), everything of a symbol record that comes before its name.
	unlocatedVariableReferenceFixedSize = 10
	// datatypeReferenceFixedSize is dataSize(2) + unknown1(2) + classIdentifier(1) + dataType(1) +
	// one reserved byte the mspec pins to 0x00.
	datatypeReferenceFixedSize = 7
	// udtDefinitionFixedSize is dataType(2) + offset(2) + unknown5(2) + unknown4(2).
	udtDefinitionFixedSize = 8
)

// parseSymbolTable reads the payload of a DD02 request for block 0xFFFF: the project's symbols.
func parseSymbolTable(block []byte) ([]readWriteModel.UmasUnlocatedVariableReference, error) {
	if len(block) < unlocatedVariableNamesHeaderSize {
		return nil, errors.Errorf("A symbol table needs at least %d header bytes, got %d", unlocatedVariableNamesHeaderSize, len(block))
	}
	numberOfRecords := binary.LittleEndian.Uint16(block[5:7])
	rest := block[unlocatedVariableNamesHeaderSize:]
	records := make([]readWriteModel.UmasUnlocatedVariableReference, 0, numberOfRecords)
	for i := uint16(0); i < numberOfRecords; i++ {
		if len(rest) < unlocatedVariableReferenceFixedSize {
			return nil, errors.Errorf("Symbol record %d is truncated: %d bytes left", i, len(rest))
		}
		dataType := binary.LittleEndian.Uint16(rest[0:2])
		blockNumber := binary.LittleEndian.Uint16(rest[2:4])
		offset := binary.LittleEndian.Uint32(rest[4:8])
		flags := rest[8]
		unknown4 := rest[9]
		name, remainder, err := readTerminatedString(rest[unlocatedVariableReferenceFixedSize:])
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading the name of symbol record %d", i)
		}
		records = append(records, readWriteModel.NewUmasUnlocatedVariableReference(
			dataType, blockNumber, offset, flags, unknown4, name))
		rest = remainder
	}
	return records, nil
}

// parseDatatypeNames reads the payload of a DD03 request: the project's datatype dictionary.
func parseDatatypeNames(block []byte) ([]readWriteModel.UmasDatatypeReference, error) {
	if len(block) < datatypeNamesHeaderSize {
		return nil, errors.Errorf("A datatype dictionary needs at least %d header bytes, got %d", datatypeNamesHeaderSize, len(block))
	}
	numberOfRecords := binary.LittleEndian.Uint16(block[4:6])
	rest := block[datatypeNamesHeaderSize:]
	records := make([]readWriteModel.UmasDatatypeReference, 0, numberOfRecords)
	for i := uint16(0); i < numberOfRecords; i++ {
		if len(rest) < datatypeReferenceFixedSize {
			return nil, errors.Errorf("Datatype record %d is truncated: %d bytes left", i, len(rest))
		}
		dataSize := binary.LittleEndian.Uint16(rest[0:2])
		unknown1 := binary.LittleEndian.Uint16(rest[2:4])
		classIdentifier := rest[4]
		dataType := rest[5]
		// rest[6] is the reserved byte the mspec pins to 0x00. A device which puts something else
		// there isn't a reason to give up on the whole dictionary, so it is skipped rather than
		// checked - the same thing the generated reserved-field reader does (it logs and carries on).
		name, remainder, err := readTerminatedString(rest[datatypeReferenceFixedSize:])
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading the name of datatype record %d", i)
		}
		records = append(records, readWriteModel.NewUmasDatatypeReference(
			dataSize, unknown1, classIdentifier, dataType, name))
		rest = remainder
	}
	return records, nil
}

// parseUdtDefinition reads the payload of a DD02 request for a struct type: the type's members.
func parseUdtDefinition(block []byte) ([]readWriteModel.UmasUDTDefinition, error) {
	if len(block) < udtDefinitionHeaderSize {
		return nil, errors.Errorf("A UDT definition needs at least %d header bytes, got %d", udtDefinitionHeaderSize, len(block))
	}
	numberOfRecords := binary.LittleEndian.Uint16(block[5:7])
	rest := block[udtDefinitionHeaderSize:]
	records := make([]readWriteModel.UmasUDTDefinition, 0, numberOfRecords)
	for i := uint16(0); i < numberOfRecords; i++ {
		if len(rest) < udtDefinitionFixedSize {
			return nil, errors.Errorf("UDT member %d is truncated: %d bytes left", i, len(rest))
		}
		dataType := binary.LittleEndian.Uint16(rest[0:2])
		offset := binary.LittleEndian.Uint16(rest[2:4])
		unknown5 := binary.LittleEndian.Uint16(rest[4:6])
		unknown4 := binary.LittleEndian.Uint16(rest[6:8])
		name, remainder, err := readTerminatedString(rest[udtDefinitionFixedSize:])
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading the name of UDT member %d", i)
		}
		records = append(records, readWriteModel.NewUmasUDTDefinition(
			dataType, offset, unknown5, unknown4, name))
		rest = remainder
	}
	return records, nil
}

// readTerminatedString takes the NUL terminated string off the front of data and hands back what
// follows it. A run without a terminator is an error rather than a string, because it means the
// record stream is out of step and every following record would be nonsense.
func readTerminatedString(data []byte) (string, []byte, error) {
	terminator := bytes.IndexByte(data, terminatorByte)
	if terminator < 0 {
		return "", nil, errors.New("no NUL terminator in the remaining bytes")
	}
	return string(data[:terminator]), data[terminator+1:], nil
}

// terminatorByte is the NUL byte a UMAS name ends with.
const terminatorByte = byte(0x00)
