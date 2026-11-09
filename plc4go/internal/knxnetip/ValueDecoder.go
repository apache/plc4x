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

package knxnetip

import (
	"context"
	"encoding/hex"

	api "github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type ValueDecoder struct {
	rb utils.ReadBuffer
}

func NewValueDecoder(rb utils.ReadBuffer) ValueDecoder {
	return ValueDecoder{
		rb: rb,
	}
}

func (m ValueDecoder) Decode(typeName string) api.PlcValue {
	ctx := context.TODO()
	datatype, ok := driverModel.KnxDatapointTypeByName(typeName)
	if !ok {
		return nil
	}
	plcValue, err := driverModel.KnxDatapointParseWithBuffer(ctx, m.rb, datatype)
	if err != nil {
		return nil
	}
	return plcValue
}

// RawGroupValue is the payload of a group value whose datatype the driver doesn't know, as it
// arrives on the bus: the byte carrying the embedded data bits, then the data bytes. It is the
// api.RawPlcValue the api promises for such a value: the application decodes it with
// RawDecodeValue once it knows the datatype, and may do so again after RawReset.
type RawGroupValue struct {
	spiValues.PlcRawByteArray
	readBuffer utils.ReadBufferByteBased
}

var _ api.RawPlcValue = (*RawGroupValue)(nil)

func NewRawGroupValue(payload []byte) *RawGroupValue {
	return &RawGroupValue{
		PlcRawByteArray: spiValues.NewPlcRawByteArray(payload),
		readBuffer:      utils.NewReadBufferByteBased(payload),
	}
}

// RawDecodeValue decodes the next value of the given datatype, a KnxDatapointType name, from the
// payload. It returns nil if the payload holds no such value, and leaves nothing more to decode
// then, so a caller looping on RawHasMore ends.
func (m *RawGroupValue) RawDecodeValue(typeName string) api.PlcValue {
	plcValue := NewValueDecoder(m.readBuffer).Decode(typeName)
	if plcValue == nil {
		m.readBuffer.Reset(uint32(len(m.Values)))
	}
	return plcValue
}

// RawHasMore tells if there is payload left to decode.
func (m *RawGroupValue) RawHasMore() bool {
	return m.readBuffer.HasMore(1)
}

// RawReset starts decoding the payload from its beginning again.
func (m *RawGroupValue) RawReset() {
	m.readBuffer.Reset(0)
}

// GetString returns the payload in hex.
func (m *RawGroupValue) GetString() string {
	return hex.EncodeToString(m.Values)
}
