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
	"bufio"
	"encoding/xml"
	"fmt"
	"math/big"
	"strings"
)

type WriteBufferXmlBased interface {
	WriteBuffer
	GetXmlString() string
}

func NewXmlWriteBuffer() WriteBufferXmlBased {
	var xmlString strings.Builder
	encoder := xml.NewEncoder(bufio.NewWriterSize(&xmlString, 1024*16))
	encoder.Indent("", "  ")
	return &xmlWriteBuffer{
		xmlString: &xmlString,
		Encoder:   encoder,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type xmlWriteBuffer struct {
	xmlString *strings.Builder
	*xml.Encoder
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (x *xmlWriteBuffer) PushContext(logicalName string, _ ...WithWriterArgs) error {
	// Pre-emptive flush to avoid overflow when for a long time no context gets popped
	if err := x.Encoder.Flush(); err != nil {
		return err
	}
	return x.Encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: sanitizeLogicalName(logicalName)}})
}

func (x *xmlWriteBuffer) WriteBit(logicalName string, value bool, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("bit", 1, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("uint", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("uint", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("uint", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("uint", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("int", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("int", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("int", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("int", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("int", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("float", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("float", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error {
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: generateAttr("float", bitLength, writerArgs...),
	})
}

func (x *xmlWriteBuffer) WriteString(logicalName string, bitLength uint8, encoding string, value string, writerArgs ...WithWriterArgs) error {
	attr := generateAttr("string", bitLength, writerArgs...)
	attr = append(attr, xml.Attr{Name: xml.Name{Local: "encoding"}, Value: encoding})
	return x.EncodeElement(value, xml.StartElement{
		Name: xml.Name{Local: sanitizeLogicalName(logicalName)},
		Attr: attr,
	})
}

func (x *xmlWriteBuffer) PopContext(logicalName string, _ ...WithWriterArgs) error {
	if err := x.Encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: sanitizeLogicalName(logicalName)}}); err != nil {
		return err
	}
	return x.Encoder.Flush()
}

func (x *xmlWriteBuffer) GetXmlString() string {
	return x.xmlString.String()
}

func generateAttr(dataType string, bitLength uint8, writerArgs ...WithWriterArgs) []xml.Attr {
	attrs := make([]xml.Attr, 2)
	attrs[0] = xml.Attr{
		Name:  xml.Name{Local: "dataType"},
		Value: dataType,
	}
	attrs[1] = xml.Attr{
		Name:  xml.Name{Local: "bitLength"},
		Value: fmt.Sprintf("%d", bitLength),
	}
	for _, arg := range writerArgs {
		if !arg.isWriterArgs() {
			panic("not a writer arg")
		}
		switch arg.(type) {
		case withAdditionalStringRepresentation:
			attrs = append(attrs, xml.Attr{
				Name:  xml.Name{Local: "stringRepresentation"},
				Value: arg.(withAdditionalStringRepresentation).stringRepresentation,
			})
		}
	}
	return attrs
}

func sanitizeLogicalName(logicalName string) string {
	if logicalName == "" {
		return "value"
	}
	return logicalName
}
