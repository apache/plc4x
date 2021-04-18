//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package values

import (
	"encoding/xml"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	api "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

type PlcValueDecoder interface {
	Decode(typeName string) api.PlcValue
}

type RawPlcValue struct {
	readBuffer utils.ReadBuffer
	decoder    PlcValueDecoder
	PlcValueAdapter
}

func NewRawPlcValue(readBuffer utils.ReadBuffer, decoder PlcValueDecoder) RawPlcValue {
	return RawPlcValue{
		readBuffer: readBuffer,
		decoder:    decoder,
	}
}

func (m RawPlcValue) GetRaw() []byte {
	return m.readBuffer.GetBytes()
}

func (m RawPlcValue) IsList() bool {
	return true
}

func (m RawPlcValue) GetLength() uint32 {
	return uint32(m.readBuffer.GetTotalBytes())
}

func (m RawPlcValue) GetIndex(i uint32) api.PlcValue {
	return NewPlcUSINT(m.readBuffer.GetBytes()[i])
}

func (m RawPlcValue) GetList() []api.PlcValue {
	var plcValues []api.PlcValue
	for _, value := range m.readBuffer.GetBytes() {
		plcValues = append(plcValues, NewPlcUSINT(value))
	}
	return plcValues
}

func (m RawPlcValue) RawDecodeValue(typeName string) api.PlcValue {
	return m.decoder.Decode(typeName)
}

func (m RawPlcValue) RawHasMore() bool {
	return m.readBuffer.HasMore(1)
}

func (m RawPlcValue) RawReset() {
	m.readBuffer.Reset()
}

func (m RawPlcValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "PlcByteArray"}}); err != nil {
		return err
	}

	// TODO: Implement this ...

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "PlcByteArray"}}); err != nil {
		return err
	}
	return nil
}
