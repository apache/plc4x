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
	"context"
	"encoding/binary"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// The three data-dictionary payloads are parsed by the generated model, the same way plc4j's
// UmasConnection parses them: the payload of the enclosing UmasPDUReadUnlocatedVariableResponse is
// handed to the response type that matches the request, and its records come back already decoded.
//
// The one thing that has to be supplied from here is the byte order. None of these mspec types pins
// a byte order on its fields, so it is the buffer that decides, and UMAS is little endian throughout.
// That rules out the convenience parsers (UmasPDUReadDatatypeNamesResponseParse and friends): they
// build a utils.NewReadBufferByteBased with no options, which is big endian, and every multi-byte
// field of these records - dataType, block, offset, dataSize, unknown1 - would come back
// byte-swapped. plc4j passes LITTLE_ENDIAN explicitly at the equivalent call, and so does
// dictionaryReadBuffer below.

// dictionaryReadBuffer wraps a data-dictionary payload in the little endian buffer its records are
// encoded in.
func dictionaryReadBuffer(block []byte) utils.ReadBufferByteBased {
	return utils.NewReadBufferByteBased(block, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
}

// parseSymbolTable reads the payload of a DD02 request for block 0xFFFF: the project's symbols.
func parseSymbolTable(ctx context.Context, block []byte) ([]readWriteModel.UmasUnlocatedVariableReference, error) {
	response, err := readWriteModel.UmasPDUReadUnlocatedVariableNamesResponseParseWithBuffer(ctx, dictionaryReadBuffer(block))
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing a symbol table out of %d bytes", len(block))
	}
	return response.GetRecords(), nil
}

// parseDatatypeNames reads the payload of a DD03 request: the project's datatype dictionary.
//
// Note that the reserved byte each record carries is stricter here than it was when this was parsed
// by hand: the Go spi's ReadReservedField fails the parse when the byte is not 0x00, where plc4j's
// only logs it. A device that puts something else there loses its whole dictionary rather than one
// record, which is a difference in the shared field reader and not something this driver can decide.
func parseDatatypeNames(ctx context.Context, block []byte) ([]readWriteModel.UmasDatatypeReference, error) {
	response, err := readWriteModel.UmasPDUReadDatatypeNamesResponseParseWithBuffer(ctx, dictionaryReadBuffer(block))
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing a datatype dictionary out of %d bytes", len(block))
	}
	return response.GetRecords(), nil
}

// parseUdtDefinition reads the payload of a DD02 request for a struct type: the type's members.
func parseUdtDefinition(ctx context.Context, block []byte) ([]readWriteModel.UmasUDTDefinition, error) {
	response, err := readWriteModel.UmasPDUReadUmasUDTDefinitionResponseParseWithBuffer(ctx, dictionaryReadBuffer(block))
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing a UDT definition out of %d bytes", len(block))
	}
	return response.GetRecords(), nil
}
