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

package modbus

const (
	// fileRecordLength is how many registers one extended-register file holds. The extended
	// register area is addressed flat, but FC 0x14/0x15 address it as a set of files, so an
	// address has to be mapped onto a file and an offset within it (plc4j
	// ModbusTcpConnection.FC_EXTENDED_REGISTERS_FILE_RECORD_LENGTH).
	fileRecordLength = 10000
	// extendedRegisterReferenceType is the only reference type FC 0x14 and FC 0x15 define; every
	// item carries it (plc4j ModbusTcpConnection.getReadRequestPdu passes a literal 6).
	extendedRegisterReferenceType = uint8(6)
)

// extendedRegisterGroup is the part of an extended-register request that falls into one file: the
// file it addresses, the offset of the first register within that file and how many registers it
// covers.
type extendedRegisterGroup struct {
	fileNumber   uint16
	recordNumber uint16
	// lengthWords is the number of registers, not bytes; a register is two bytes wide.
	lengthWords uint16
}

// splitExtendedRegister maps a flat extended-register address and a length onto the file layout FC
// 0x14 and FC 0x15 use. The first file starts at file number 1, and a request that would run past
// the end of a file continues at offset 0 of the next one (plc4j
// ModbusTcpConnection.getReadRequestPdu / getWriteRequestPdu).
//
// Unlike plc4j, which only ever emits a second group, this keeps splitting until the whole request
// is covered: with an address of 65535 registers and a length in the same order of magnitude a
// request spans up to seven files, and plc4j's single extra group would claim more registers than
// a file holds.
func splitExtendedRegister(address uint16, lengthWords uint16) []extendedRegisterGroup {
	fileNumber := address/fileRecordLength + 1
	recordNumber := address % fileRecordLength

	var groups []extendedRegisterGroup
	for remaining := lengthWords; remaining > 0; {
		inThisFile := min(uint16(fileRecordLength)-recordNumber, remaining)
		groups = append(groups, extendedRegisterGroup{
			fileNumber:   fileNumber,
			recordNumber: recordNumber,
			lengthWords:  inThisFile,
		})
		remaining -= inThisFile
		fileNumber++
		recordNumber = 0
	}
	return groups
}
