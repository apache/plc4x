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

package model

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// CBusOptions is the corresponding interface of CBusOptions
type CBusOptions interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetConnect returns Connect (property field)
	GetConnect() bool
	// GetSmart returns Smart (property field)
	GetSmart() bool
	// GetIdmon returns Idmon (property field)
	GetIdmon() bool
	// GetExstat returns Exstat (property field)
	GetExstat() bool
	// GetMonitor returns Monitor (property field)
	GetMonitor() bool
	// GetMonall returns Monall (property field)
	GetMonall() bool
	// GetPun returns Pun (property field)
	GetPun() bool
	// GetPcn returns Pcn (property field)
	GetPcn() bool
	// GetSrchk returns Srchk (property field)
	GetSrchk() bool
}

// CBusOptionsExactly can be used when we want exactly this type and not a type which fulfills CBusOptions.
// This is useful for switch cases.
type CBusOptionsExactly interface {
	CBusOptions
	isCBusOptions() bool
}

// _CBusOptions is the data-structure of this message
type _CBusOptions struct {
	Connect bool
	Smart   bool
	Idmon   bool
	Exstat  bool
	Monitor bool
	Monall  bool
	Pun     bool
	Pcn     bool
	Srchk   bool
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CBusOptions) GetConnect() bool {
	return m.Connect
}

func (m *_CBusOptions) GetSmart() bool {
	return m.Smart
}

func (m *_CBusOptions) GetIdmon() bool {
	return m.Idmon
}

func (m *_CBusOptions) GetExstat() bool {
	return m.Exstat
}

func (m *_CBusOptions) GetMonitor() bool {
	return m.Monitor
}

func (m *_CBusOptions) GetMonall() bool {
	return m.Monall
}

func (m *_CBusOptions) GetPun() bool {
	return m.Pun
}

func (m *_CBusOptions) GetPcn() bool {
	return m.Pcn
}

func (m *_CBusOptions) GetSrchk() bool {
	return m.Srchk
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCBusOptions factory function for _CBusOptions
func NewCBusOptions(connect bool, smart bool, idmon bool, exstat bool, monitor bool, monall bool, pun bool, pcn bool, srchk bool) *_CBusOptions {
	return &_CBusOptions{Connect: connect, Smart: smart, Idmon: idmon, Exstat: exstat, Monitor: monitor, Monall: monall, Pun: pun, Pcn: pcn, Srchk: srchk}
}

// Deprecated: use the interface for direct cast
func CastCBusOptions(structType any) CBusOptions {
	if casted, ok := structType.(CBusOptions); ok {
		return casted
	}
	if casted, ok := structType.(*CBusOptions); ok {
		return *casted
	}
	return nil
}

func (m *_CBusOptions) GetTypeName() string {
	return "CBusOptions"
}

func (m *_CBusOptions) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (connect)
	lengthInBits += 1

	// Simple field (smart)
	lengthInBits += 1

	// Simple field (idmon)
	lengthInBits += 1

	// Simple field (exstat)
	lengthInBits += 1

	// Simple field (monitor)
	lengthInBits += 1

	// Simple field (monall)
	lengthInBits += 1

	// Simple field (pun)
	lengthInBits += 1

	// Simple field (pcn)
	lengthInBits += 1

	// Simple field (srchk)
	lengthInBits += 1

	return lengthInBits
}

func (m *_CBusOptions) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CBusOptionsParse(ctx context.Context, theBytes []byte) (CBusOptions, error) {
	return CBusOptionsParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func CBusOptionsParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (CBusOptions, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (CBusOptions, error) {
		return CBusOptionsParseWithBuffer(ctx, readBuffer)
	}
}

func CBusOptionsParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (CBusOptions, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CBusOptions"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CBusOptions")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	connect, err := ReadSimpleField(ctx, "connect", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'connect' field"))
	}

	smart, err := ReadSimpleField(ctx, "smart", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'smart' field"))
	}

	idmon, err := ReadSimpleField(ctx, "idmon", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'idmon' field"))
	}

	exstat, err := ReadSimpleField(ctx, "exstat", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'exstat' field"))
	}

	monitor, err := ReadSimpleField(ctx, "monitor", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'monitor' field"))
	}

	monall, err := ReadSimpleField(ctx, "monall", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'monall' field"))
	}

	pun, err := ReadSimpleField(ctx, "pun", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'pun' field"))
	}

	pcn, err := ReadSimpleField(ctx, "pcn", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'pcn' field"))
	}

	srchk, err := ReadSimpleField(ctx, "srchk", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'srchk' field"))
	}

	if closeErr := readBuffer.CloseContext("CBusOptions"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CBusOptions")
	}

	// Create the instance
	return &_CBusOptions{
		Connect: connect,
		Smart:   smart,
		Idmon:   idmon,
		Exstat:  exstat,
		Monitor: monitor,
		Monall:  monall,
		Pun:     pun,
		Pcn:     pcn,
		Srchk:   srchk,
	}, nil
}

func (m *_CBusOptions) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CBusOptions) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("CBusOptions"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for CBusOptions")
	}

	if err := WriteSimpleField[bool](ctx, "connect", m.GetConnect(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'connect' field")
	}

	if err := WriteSimpleField[bool](ctx, "smart", m.GetSmart(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'smart' field")
	}

	if err := WriteSimpleField[bool](ctx, "idmon", m.GetIdmon(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'idmon' field")
	}

	if err := WriteSimpleField[bool](ctx, "exstat", m.GetExstat(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'exstat' field")
	}

	if err := WriteSimpleField[bool](ctx, "monitor", m.GetMonitor(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'monitor' field")
	}

	if err := WriteSimpleField[bool](ctx, "monall", m.GetMonall(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'monall' field")
	}

	if err := WriteSimpleField[bool](ctx, "pun", m.GetPun(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'pun' field")
	}

	if err := WriteSimpleField[bool](ctx, "pcn", m.GetPcn(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'pcn' field")
	}

	if err := WriteSimpleField[bool](ctx, "srchk", m.GetSrchk(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'srchk' field")
	}

	if popErr := writeBuffer.PopContext("CBusOptions"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for CBusOptions")
	}
	return nil
}

func (m *_CBusOptions) isCBusOptions() bool {
	return true
}

func (m *_CBusOptions) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
