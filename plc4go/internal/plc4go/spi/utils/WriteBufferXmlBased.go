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

package utils

import (
	"encoding/xml"
	"fmt"
	"math/big"
	"strings"
)

type WriteBufferXmlBased interface {
	WriteBuffer
	GetXmlString() string
}

//NewXmlWriteBuffer returns a WriteBufferXmlBased which renders all information into xml
func NewXmlWriteBuffer() WriteBufferXmlBased {
	var xmlString strings.Builder
	encoder := xml.NewEncoder(&xmlString)
	encoder.Indent("", "  ")
	return &xmlWriteBuffer{
		xmlString:     &xmlString,
		Encoder:       encoder,
		doRenderLists: true,
		doRenderAttr:  true,
	}
}

//NewConfiguredXmlWriteBuffer returns a WriteBufferXmlBased which renders configured information into xml
func NewConfiguredXmlWriteBuffer(renderLists bool, renderAttr bool) WriteBufferXmlBased {
	var xmlString strings.Builder
	encoder := xml.NewEncoder(&xmlString)
	encoder.Indent("", "  ")
	return &xmlWriteBuffer{
		xmlString:     &xmlString,
		Encoder:       encoder,
		doRenderLists: renderLists,
		doRenderAttr:  renderAttr,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type xmlWriteBuffer struct {
	bufferCommons
	xmlString *strings.Builder
	*xml.Encoder
	doRenderLists bool
	doRenderAttr  bool
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (x *xmlWriteBuffer) PushContext(logicalName string, writerArgs ...WithWriterArgs) error {
	// Pre-emptive flush to avoid overflow when for a long time no context gets popped
	if err := x.Flush(); err != nil {
		return err
	}
	attrs := make([]xml.Attr, 0)
	attrs = x.markAsListIfRequired(writerArgs, attrs)
	return x.EncodeToken(xml.StartElement{Name: xml.Name{Local: x.sanitizeLogicalName(logicalName)}, Attr: attrs})
}

func (x *xmlWriteBuffer) WriteBit(logicalName string, value bool, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwBitKey, 1, writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteByte(logicalName string, value byte, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, fmt.Sprintf("%#02x", value), x.generateAttr(rwByteKey, 8, writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteByteArray(logicalName string, data []byte, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, fmt.Sprintf("%#02x", data), x.generateAttr(rwByteKey, uint(len(data)*8), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwUintKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwUintKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwUintKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwUintKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwIntKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwIntKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwIntKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwIntKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwIntKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwFloatKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwFloatKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error {
	return x.encodeElement(logicalName, value, x.generateAttr(rwFloatKey, uint(bitLength), writerArgs...), writerArgs...)
}

func (x *xmlWriteBuffer) WriteString(logicalName string, bitLength uint8, encoding string, value string, writerArgs ...WithWriterArgs) error {
	attr := x.generateAttr(rwStringKey, uint(bitLength), writerArgs...)
	attr = append(attr, xml.Attr{Name: xml.Name{Local: rwEncodingKey}, Value: encoding})
	return x.encodeElement(logicalName, value, attr, writerArgs...)
}

func (x *xmlWriteBuffer) PopContext(logicalName string, _ ...WithWriterArgs) error {
	if err := x.Encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: x.sanitizeLogicalName(logicalName)}}); err != nil {
		return err
	}
	return x.Encoder.Flush()
}

func (x *xmlWriteBuffer) GetXmlString() string {
	return x.xmlString.String()
}

func (x *xmlWriteBuffer) encodeElement(logicalName string, value interface{}, attr []xml.Attr, _ ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: x.sanitizeLogicalName(logicalName)},
		Attr: attr,
	})
}

func (x *xmlWriteBuffer) generateAttr(dataType string, bitLength uint, writerArgs ...WithWriterArgs) []xml.Attr {
	attrs := make([]xml.Attr, 2)
	if !x.doRenderAttr {
		return attrs
	}
	attrs[0] = xml.Attr{
		Name:  xml.Name{Local: rwDataTypeKey},
		Value: dataType,
	}
	attrs[1] = xml.Attr{
		Name:  xml.Name{Local: rwBitLengthKey},
		Value: fmt.Sprintf("%d", bitLength),
	}
	for _, arg := range writerArgs {
		if !arg.isWriterArgs() {
			panic("not a writer arg")
		}
		switch arg.(type) {
		case withAdditionalStringRepresentation:
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: rwStringRepresentationKey},
				Value: arg.(withAdditionalStringRepresentation).stringRepresentation,
			})
		}
	}
	return attrs
}

func (x *xmlWriteBuffer) markAsListIfRequired(writerArgs []WithWriterArgs, attrs []xml.Attr) []xml.Attr {
	if !x.doRenderLists {
		return attrs
	}
	if x.isToBeRenderedAsList(upcastWriterArgs(writerArgs...)...) {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: rwIsListKey},
			Value: "true",
		})
	}
	return attrs
}
