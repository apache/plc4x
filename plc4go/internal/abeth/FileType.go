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

package abeth

import (
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

// FileType is the kind of data file a tag addresses. Ported from plc4j's
// org.apache.plc4x.java.abeth.types.FileType.
//
// The type code is what goes on the wire in the "protected typed logical read" command. Note that
// the last three entries are not AB file types at all: WORD, DWORD and SINGLEBIT are conveniences
// the plc4x driver invented on top of the INTEGER file, which is why they share its type code
// 0x89 and differ only in how the bytes that come back are decoded.
type FileType uint8

const (
	FileTypeStatus FileType = iota
	FileTypeBit
	FileTypeTimer
	FileTypeCounter
	FileTypeControl
	FileTypeInteger
	FileTypeFloat
	FileTypeOutput
	FileTypeInput
	FileTypeString
	FileTypeAscii
	FileTypeBcd
	// FileTypeWord is a 2 byte integer out of an INTEGER file.
	FileTypeWord
	// FileTypeDword is a 4 byte integer out of an INTEGER file.
	FileTypeDword
	// FileTypeSinglebit is a single bit out of an INTEGER file.
	FileTypeSinglebit
)

// fileTypeNames are the names a tag address spells a file type with, matching the plc4j enum
// constant names (the tag handler upper-cases what the user wrote before looking it up).
var fileTypeNames = map[FileType]string{
	FileTypeStatus:    "STATUS",
	FileTypeBit:       "BIT",
	FileTypeTimer:     "TIMER",
	FileTypeCounter:   "COUNTER",
	FileTypeControl:   "CONTROL",
	FileTypeInteger:   "INTEGER",
	FileTypeFloat:     "FLOAT",
	FileTypeOutput:    "OUTPUT",
	FileTypeInput:     "INPUT",
	FileTypeString:    "STRING",
	FileTypeAscii:     "ASCII",
	FileTypeBcd:       "BCD",
	FileTypeWord:      "WORD",
	FileTypeDword:     "DWORD",
	FileTypeSinglebit: "SINGLEBIT",
}

// fileTypeCodes are the type codes the wire format uses. WORD, DWORD and SINGLEBIT deliberately
// carry the INTEGER code: to the PLC they *are* reads of an integer file.
var fileTypeCodes = map[FileType]uint8{
	FileTypeStatus:    0x84,
	FileTypeBit:       0x85,
	FileTypeTimer:     0x86,
	FileTypeCounter:   0x87,
	FileTypeControl:   0x88,
	FileTypeInteger:   0x89,
	FileTypeFloat:     0x8A,
	FileTypeOutput:    0x8B,
	FileTypeInput:     0x8C,
	FileTypeString:    0x8D,
	FileTypeAscii:     0x8E,
	FileTypeBcd:       0x8F,
	FileTypeWord:      0x89,
	FileTypeDword:     0x89,
	FileTypeSinglebit: 0x89,
}

// fileTypeValueTypes are the plc4x value types a read of the file type produces, mirroring plc4j's
// FileType constructor arguments.
var fileTypeValueTypes = map[FileType]apiValues.PlcValueType{
	FileTypeStatus:  apiValues.RAW_BYTE_ARRAY,
	FileTypeBit:     apiValues.BOOL,
	FileTypeTimer:   apiValues.TIME,
	FileTypeCounter: apiValues.RAW_BYTE_ARRAY,
	FileTypeControl: apiValues.RAW_BYTE_ARRAY,
	FileTypeInteger: apiValues.INT,
	FileTypeFloat:   apiValues.REAL,
	FileTypeOutput:  apiValues.RAW_BYTE_ARRAY,
	FileTypeInput:   apiValues.RAW_BYTE_ARRAY,
	FileTypeString:  apiValues.STRING,
	FileTypeAscii:   apiValues.STRING,
	FileTypeBcd:     apiValues.RAW_BYTE_ARRAY,
	FileTypeWord:    apiValues.WORD,
	FileTypeDword:   apiValues.DWORD,
	// Deliberate deviation from plc4j, which declares SINGLEBIT as WORD while its decoder produces
	// a boolean from it - a tag whose advertised value type doesn't match the value it yields.
	FileTypeSinglebit: apiValues.BOOL,
}

// GetTypeCode is the byte the "protected typed logical read" command carries for this file type.
func (f FileType) GetTypeCode() uint8 {
	return fileTypeCodes[f]
}

// GetPlcValueType is the plc4x value type a read of this file type produces.
func (f FileType) GetPlcValueType() apiValues.PlcValueType {
	if valueType, ok := fileTypeValueTypes[f]; ok {
		return valueType
	}
	return apiValues.NULL
}

// hasFixedWidth says whether the file type decides on its own how many bytes a read asks for. The
// three synthetic types do (2, 4 and 2 bytes); every real file type takes its width from the
// [<byteSize>] suffix of the tag address.
func (f FileType) hasFixedWidth() bool {
	switch f {
	case FileTypeWord, FileTypeDword, FileTypeSinglebit:
		return true
	default:
		return false
	}
}

func (f FileType) String() string {
	if name, ok := fileTypeNames[f]; ok {
		return name
	}
	return "UNKNOWN"
}

// FileTypeByName resolves the name a tag address spells. The name has to be upper case, which is
// what the tag handler hands in.
func FileTypeByName(name string) (FileType, bool) {
	for fileType, fileTypeName := range fileTypeNames {
		if fileTypeName == name {
			return fileType, true
		}
	}
	return 0, false
}
