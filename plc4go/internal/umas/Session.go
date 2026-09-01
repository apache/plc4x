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
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
)

const (
	// customTypeIdBase is the type id the project's own types start at: everything below it is a
	// primitive UmasDataType, and the n-th entry of the datatype dictionary is
	// customTypeIdBase + n. plc4j calls this CUSTOM_TYPE_THRESHOLD.
	customTypeIdBase = uint16(0x1A)
	// defaultStringBufferSize is how many bytes a STRING read asks for when nothing else is known
	// about the symbol. plc4j's DEFAULT_STRING_BUFFER_SIZE.
	defaultStringBufferSize = uint16(254)
	// arrayClassId is the value the first byte of a resolved custom type carries when the type is
	// an array rather than a struct.
	arrayClassId = uint8(0x04)
)

// session is everything a connected UMAS session knows: what the PLC said about itself during the
// handshake, and the data dictionary the driver downloaded from it. It is the plc4go counterpart of
// the state plc4j keeps in the fields of UmasConnection (which grew out of the old
// UmasDriverContext).
//
// Every read and write needs the project CRC and the symbol table, and browse needs the type
// tables, so all of it is shared between the reader, the writer and the browser and guarded by one
// lock.
type session struct {
	// transactionIdentifier hands out the Modbus transaction identifier of the next request. It is
	// the only thing tying a response to its request.
	transactionIdentifier atomic.Uint32

	mutex sync.RWMutex

	// hostname, plcModel and firmwareVersion are what the PlcIdent and InitComms responses report.
	// They are informational - nothing in the request path depends on them.
	hostname        string
	plcModel        uint16
	firmwareVersion uint16
	// maxFrameSize starts at the configured value and is replaced by what the PLC reports in the
	// InitComms response.
	maxFrameSize uint16
	// pairingKey is the first byte of every UMAS PDU. plc4j declares a field for it and never
	// assigns one, so it stays zero for the whole life of the connection; the same is true here,
	// which keeps the wire bytes identical. It is kept as a field rather than a constant because
	// the PLC reservation requests (FC 0x10/0x11), which this driver doesn't send, are what would
	// hand out a non-zero one.
	pairingKey uint8
	// hardwareId identifies the PLC's memory layout and is required by every data-dictionary
	// request. It is read out of memory block 0x30.
	hardwareId uint32
	// projectCrc has to accompany every variable read and write. plc4j derives it from memory
	// block 0x30 as the sum of the two 32 bit hashes which follow the block's 9 byte header, a
	// relation it documents as discovered by comparing Schneider OPC UA Server traffic against the
	// raw block - there is no independent specification of it.
	projectCrc uint32

	// symbolTable maps the lower-cased symbol name to its dictionary entry.
	symbolTable map[string]readWriteModel.UmasUnlocatedVariableReference
	// dataTypeSizes is the allocated byte size the datatype dictionary reports for a type id.
	dataTypeSizes map[uint16]uint16
	// symbolSizes is the byte size derived from the symbol layout, see computeSymbolSizes.
	symbolSizes map[string]uint16
	// customTypeNames maps a custom type id to the name the project gave it.
	customTypeNames map[uint16]string
	// customTypeFields are the members of a custom struct type. An array type has an entry with no
	// members, which is what tells "resolved as an array" from "not a custom type at all".
	customTypeFields map[uint16][]readWriteModel.UmasUDTDefinition
	// customTypeElementTypeIds is the element type of a custom array type.
	customTypeElementTypeIds map[uint16]uint16
	// customTypeDimensions are the bounds of a custom array type.
	customTypeDimensions map[uint16][]readWriteModel.UmasArrayDimension
}

func newSession(maxFrameSize uint16) *session {
	s := &session{
		maxFrameSize:             maxFrameSize,
		symbolTable:              map[string]readWriteModel.UmasUnlocatedVariableReference{},
		dataTypeSizes:            map[uint16]uint16{},
		symbolSizes:              map[string]uint16{},
		customTypeNames:          map[uint16]string{},
		customTypeFields:         map[uint16][]readWriteModel.UmasUDTDefinition{},
		customTypeElementTypeIds: map[uint16]uint16{},
		customTypeDimensions:     map[uint16][]readWriteModel.UmasArrayDimension{},
	}
	return s
}

// nextTransactionIdentifier hands out the identifier of the next request. The field on the wire is
// 16 bits wide and zero is left out so an identifier never collides with an uninitialized one.
//
// Deliberate deviation from plc4j, whose UmasConnection.nextTransactionId compares-and-sets on
// id + 1 (an identifier it never handed out) and therefore keeps handing out identifiers above
// 0xFFFF, masking them down to 16 bits only in the value it returns while the tracker is keyed on
// the unmasked one.
func (s *session) nextTransactionIdentifier() uint16 {
	for {
		next := s.transactionIdentifier.Add(1)
		if next <= math.MaxUint16 {
			return uint16(next)
		}
		// Someone has to reset the counter; whoever wins gets 1 and the losers try again.
		if s.transactionIdentifier.CompareAndSwap(next, 1) {
			return 1
		}
	}
}

func (s *session) getPairingKey() uint8 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.pairingKey
}

func (s *session) getHardwareId() uint32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.hardwareId
}

func (s *session) getProjectCrc() uint32 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.projectCrc
}

func (s *session) getMaxFrameSize() uint16 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.maxFrameSize
}

func (s *session) setIdentity(hostname string, plcModel uint16, firmwareVersion uint16) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.hostname = hostname
	s.plcModel = plcModel
	s.firmwareVersion = firmwareVersion
}

func (s *session) setCommsParameters(maxFrameSize uint16, firmwareVersion uint16) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if maxFrameSize >= minMaxFrameSize {
		s.maxFrameSize = maxFrameSize
	}
	s.firmwareVersion = firmwareVersion
}

func (s *session) setProjectIdentity(hardwareId uint32, projectCrc uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.hardwareId = hardwareId
	s.projectCrc = projectCrc
}

// lookupSymbol finds the dictionary entry of a symbol. Names are matched case insensitively, the way
// plc4j does by lower-casing both the table keys and the requested address.
func (s *session) lookupSymbol(symbolicAddress string) (readWriteModel.UmasUnlocatedVariableReference, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	symbol, ok := s.symbolTable[strings.ToLower(symbolicAddress)]
	return symbol, ok
}

// hasSymbols says whether the data dictionary has been downloaded.
func (s *session) hasSymbols() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.symbolTable) > 0
}

// symbols hands out the dictionary entries sorted by name, so a browse answers in a stable order
// rather than in Go's randomized map order.
func (s *session) symbols() []readWriteModel.UmasUnlocatedVariableReference {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	names := make([]string, 0, len(s.symbolTable))
	for name := range s.symbolTable {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]readWriteModel.UmasUnlocatedVariableReference, 0, len(names))
	for _, name := range names {
		result = append(result, s.symbolTable[name])
	}
	return result
}

// stringBufferSize is how many bytes a STRING read of this symbol asks for: what the datatype
// dictionary says about its type, else what the symbol layout suggests, else the default. Ported
// from plc4j's buildReadReference.
//
// The dictionary lookup can only ever miss as things stand, in plc4j as much as here: only a symbol
// whose type id is the primitive STRING takes this path at all, and the dictionary is keyed by custom
// type id (customTypeIdBase and up). It is kept because the three-way fallback is what plc4j spells
// out, and because a project whose strings are declared as a custom type would need it the moment
// such a type is recognised as a string.
func (s *session) stringBufferSize(symbol readWriteModel.UmasUnlocatedVariableReference) uint16 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if size, ok := s.dataTypeSizes[symbol.GetDataType()]; ok && size > 0 {
		return size
	}
	if size, ok := s.symbolSizes[strings.ToLower(symbol.GetValue())]; ok && size > 0 {
		return size
	}
	return defaultStringBufferSize
}

// setDataTypes records what the datatype dictionary said. The n-th entry describes the type with id
// customTypeIdBase + n; its declared byte size matters for STRING and for struct types, whose size
// can't be derived from a UmasDataType.
func (s *session) setDataTypes(references []readWriteModel.UmasDatatypeReference) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.dataTypeSizes = map[uint16]uint16{}
	for i, reference := range references {
		s.dataTypeSizes[customTypeIdBase+uint16(i)] = reference.GetDataSize()
	}
}

// setStructType records a resolved custom struct type.
func (s *session) setStructType(typeId uint16, name string, fields []readWriteModel.UmasUDTDefinition) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.customTypeNames[typeId] = name
	s.customTypeFields[typeId] = fields
}

// setArrayType records a resolved custom array type. The empty field list is what plc4j stores too:
// it marks the id as a known custom type without giving it members.
func (s *session) setArrayType(typeId uint16, name string, elementTypeId uint16, dimensions []readWriteModel.UmasArrayDimension) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.customTypeNames[typeId] = name
	s.customTypeFields[typeId] = nil
	s.customTypeElementTypeIds[typeId] = elementTypeId
	s.customTypeDimensions[typeId] = dimensions
}

// customTypeInfo is what the dictionary knows about one of the project's own types.
type customTypeInfo struct {
	// name is what the project calls the type.
	name string
	// isArray says whether it is an array type; if not it is a struct.
	isArray bool
	// elementTypeId and dimensions describe an array type.
	elementTypeId uint16
	dimensions    []readWriteModel.UmasArrayDimension
	// fields are the members of a struct type.
	fields []readWriteModel.UmasUDTDefinition
}

// customType looks up one of the project's own types. The second return value is false for a type id
// which is a primitive, or one whose definition couldn't be downloaded.
func (s *session) customType(typeId uint16) (customTypeInfo, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	fields, known := s.customTypeFields[typeId]
	if !known {
		return customTypeInfo{}, false
	}
	elementTypeId, isArray := s.customTypeElementTypeIds[typeId]
	return customTypeInfo{
		name:          s.customTypeNames[typeId],
		isArray:       isArray,
		elementTypeId: elementTypeId,
		dimensions:    s.customTypeDimensions[typeId],
		fields:        fields,
	}, true
}

// setSymbols replaces the symbol table and recomputes the derived symbol sizes.
func (s *session) setSymbols(symbols []readWriteModel.UmasUnlocatedVariableReference) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.symbolTable = make(map[string]readWriteModel.UmasUnlocatedVariableReference, len(symbols))
	for _, symbol := range symbols {
		s.symbolTable[strings.ToLower(symbol.GetValue())] = symbol
	}
	s.symbolSizes = computeSymbolSizes(s.symbolTable)
}

// computeSymbolSizes derives how many bytes each symbol occupies from the layout of the symbol
// table: inside one memory block the symbols are laid out back to back, so the distance to the next
// symbol's offset is the size of this one. The last symbol of every block has no successor and
// therefore no derived size. Ported from plc4j's UmasConnection.computeSymbolSizes.
func computeSymbolSizes(symbolTable map[string]readWriteModel.UmasUnlocatedVariableReference) map[string]uint16 {
	type entry struct {
		name   string
		offset uint32
	}
	byBlock := map[uint16][]entry{}
	for name, symbol := range symbolTable {
		byBlock[symbol.GetBlock()] = append(byBlock[symbol.GetBlock()], entry{name: name, offset: symbol.GetOffset()})
	}
	sizes := map[string]uint16{}
	for _, entries := range byBlock {
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].offset != entries[j].offset {
				return entries[i].offset < entries[j].offset
			}
			// Two symbols at the same offset can't be ordered by offset; ordering them by name
			// keeps the result independent of Go's map iteration order.
			return entries[i].name < entries[j].name
		})
		for i := 0; i < len(entries)-1; i++ {
			size := entries[i+1].offset - entries[i].offset
			if size > 0 && size <= math.MaxUint16 {
				sizes[entries[i].name] = uint16(size)
			}
		}
	}
	return sizes
}

func (s *session) String() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return fmt.Sprintf("umas.session{hostname: %s, model: %d, firmware: %d, maxFrameSize: %d, hardwareId: 0x%08X, projectCrc: 0x%08X, symbols: %d}",
		s.hostname, s.plcModel, s.firmwareVersion, s.maxFrameSize, s.hardwareId, s.projectCrc, len(s.symbolTable))
}
