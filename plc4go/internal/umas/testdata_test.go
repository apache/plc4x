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
	"encoding/binary"
)

// There is no recorded UMAS traffic in this repository, in any language, so the payloads the tests
// feed the driver are built here from the layouts the mspec declares. Every builder writes the
// fields of its type in mspec order, little endian, and terminates the names with a NUL byte the way
// the 'parseTerminatedString' manual field does.
//
// These builders are what makes the dictionary parsing testable: build a payload from the declared
// layout, parse it through the generated model, and the two have to agree.

// symbolRecord is one entry of the project's symbol table.
type symbolRecord struct {
	name     string
	dataType uint16
	block    uint16
	offset   uint32
	flags    uint8
	unknown4 uint8
}

// symbolTablePayload builds the payload of a DD02 response for block 0xFFFF:
//
//	range(1) + nextAddress(2) + unknown1(2) + noOfRecords(2) + records
func symbolTablePayload(records ...symbolRecord) []byte {
	payload := []byte{0x00}
	payload = binary.LittleEndian.AppendUint16(payload, 0) // nextAddress: the whole table fits
	payload = binary.LittleEndian.AppendUint16(payload, 0) // unknown1
	payload = binary.LittleEndian.AppendUint16(payload, uint16(len(records)))
	for _, record := range records {
		payload = binary.LittleEndian.AppendUint16(payload, record.dataType)
		payload = binary.LittleEndian.AppendUint16(payload, record.block)
		payload = binary.LittleEndian.AppendUint32(payload, record.offset)
		payload = append(payload, record.flags, record.unknown4)
		payload = append(payload, record.name...)
		payload = append(payload, 0x00)
	}
	return payload
}

// datatypeRecord is one entry of the project's datatype dictionary.
type datatypeRecord struct {
	name string
	// dataSize is the allocated byte size, which is what a STRING read asks for.
	dataSize uint16
	// classIdentifier being non-zero is what makes the driver fetch the type's definition.
	classIdentifier uint8
	// dataType is the primitive this type is based on.
	dataType uint8
}

// datatypeDictionaryPayload builds the payload of a DD03 response:
//
//	range(1) + nextAddress(2) + unknown1(1) + noOfRecords(2) + records
//
// where a record is dataSize(2) + unknown1(2) + classIdentifier(1) + dataType(1) + a reserved byte +
// the NUL terminated name.
func datatypeDictionaryPayload(records ...datatypeRecord) []byte {
	payload := []byte{0x00}
	payload = binary.LittleEndian.AppendUint16(payload, 0) // nextAddress
	payload = append(payload, 0x00)                        // unknown1
	payload = binary.LittleEndian.AppendUint16(payload, uint16(len(records)))
	for _, record := range records {
		payload = binary.LittleEndian.AppendUint16(payload, record.dataSize)
		payload = binary.LittleEndian.AppendUint16(payload, 0) // unknown1
		payload = append(payload, record.classIdentifier, record.dataType)
		payload = append(payload, 0x00) // the reserved byte the mspec pins to zero
		payload = append(payload, record.name...)
		payload = append(payload, 0x00)
	}
	return payload
}

// udtMember is one member of a struct type.
type udtMember struct {
	name     string
	dataType uint16
	offset   uint16
}

// udtDefinitionPayload builds the payload of a DD02 response for a struct type:
//
//	range(1) + unknown1(4) + noOfRecords(2) + records
//
// where a record is dataType(2) + offset(2) + unknown5(2) + unknown4(2) + the NUL terminated name.
// The first byte has to be something other than 0x04, which marks an array.
func udtDefinitionPayload(members ...udtMember) []byte {
	payload := []byte{0x01}
	payload = binary.LittleEndian.AppendUint32(payload, 0) // unknown1
	payload = binary.LittleEndian.AppendUint16(payload, uint16(len(members)))
	for _, member := range members {
		payload = binary.LittleEndian.AppendUint16(payload, member.dataType)
		payload = binary.LittleEndian.AppendUint16(payload, member.offset)
		payload = binary.LittleEndian.AppendUint16(payload, 0) // unknown5
		payload = binary.LittleEndian.AppendUint16(payload, 0) // unknown4
		payload = append(payload, member.name...)
		payload = append(payload, 0x00)
	}
	return payload
}

// arrayDimension is one dimension of an array type, with inclusive bounds.
type arrayDimension struct {
	startIndex uint32
	upperBound uint32
}

// arrayTypeDefinitionPayload builds the payload of a DD02 response for an array type:
//
//	classId(1) = 0x04 + elementTypeId(2) + numberOfDimensions(1) + dimensions
//
// where a dimension is startIndex(4) + upperBound(4). This is the one dictionary payload whose mspec
// type pins LITTLE_ENDIAN on its own fields, so it needs no byte order from the buffer.
func arrayTypeDefinitionPayload(elementTypeId uint16, dimensions ...arrayDimension) []byte {
	payload := []byte{arrayClassId}
	payload = binary.LittleEndian.AppendUint16(payload, elementTypeId)
	payload = append(payload, uint8(len(dimensions)))
	for _, dimension := range dimensions {
		payload = binary.LittleEndian.AppendUint32(payload, dimension.startIndex)
		payload = binary.LittleEndian.AppendUint32(payload, dimension.upperBound)
	}
	return payload
}

// projectMemoryBlockPayload builds the 17 bytes of memory block 0x30 the handshake reads the project
// identity out of: the 9 byte UmasMemoryBlockBasicInfo header (range(2) + notSure(2) + index(1) +
// hardwareId(4)) followed by the two 32 bit hashes whose sum is the project CRC.
func projectMemoryBlockPayload(hardwareId uint32, firstHash uint32, secondHash uint32) []byte {
	payload := []byte{}
	payload = binary.LittleEndian.AppendUint16(payload, 0) // range
	payload = binary.LittleEndian.AppendUint16(payload, 0) // notSure
	payload = append(payload, 0x00)                        // index
	payload = binary.LittleEndian.AppendUint32(payload, hardwareId)
	payload = binary.LittleEndian.AppendUint32(payload, firstHash)
	payload = binary.LittleEndian.AppendUint32(payload, secondHash)
	return payload
}
