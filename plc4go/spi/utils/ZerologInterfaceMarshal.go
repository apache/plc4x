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

package utils

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// PLCMessageFormat defines the format of the PLCMessage that is logged
type PLCMessageFormat uint8

const (
	// PLCMessageAsJSON defines the format of the PLCMessage as JSON
	PLCMessageAsJSON PLCMessageFormat = iota
	// PLCMessageAsString defines the format of the PLCMessage as a string
	PLCMessageAsString
	// PLCMessageAsXML defines the format of the PLCMessage as XML
	PLCMessageAsXML
	// PLCMessageAsByte defines the format of the PLCMessage as a byte array
	PLCMessageAsByte
)

// ZerologInterfacePLCMessageFormat defines the format of the PLCMessage that is logged by the zerolog interface marshal function
var ZerologInterfacePLCMessageFormat = PLCMessageAsJSON

// ZerologDefaultInterfaceMarshalFunc is the default marshal function used by zerolog
var ZerologDefaultInterfaceMarshalFunc = zerolog.InterfaceMarshalFunc

// ZerologMessageInterfaceMarshalFunc is the marshal function used by zerolog to serialize PLCMessages
// (Serializable values) in the format set by ZerologInterfacePLCMessageFormat; anything else goes to
// ZerologDefaultInterfaceMarshalFunc. To use it set zerolog.InterfaceMarshalFunc =
// ZerologMessageInterfaceMarshalFunc, e.g. in an init function.
var ZerologMessageInterfaceMarshalFunc = func(v any) ([]byte, error) {
	if plcMessage, ok := v.(Serializable); ok {
		switch ZerologInterfacePLCMessageFormat {
		case PLCMessageAsJSON:
			jsonWriteBuffer := NewJsonWriteBuffer(
				WithJsonWriteBufferDefaultIdent(false),
			)
			if err := plcMessage.SerializeWithWriteBuffer(context.Background(), jsonWriteBuffer); err != nil {
				return nil, errors.Wrap(err, "error serializing PLCMessage")
			}
			jsonString, err := jsonWriteBuffer.GetJsonString()
			if err != nil {
				return nil, errors.Wrapf(err, "error getting JSON string from PLCMessage")
			}
			b := []byte(jsonString)
			if len(b) > 0 {
				// Remove trailing \n which is added by Encode.
				return b[:len(b)-1], nil
			}
			return b, nil
		case PLCMessageAsString:
			boxWriteBuffer := NewWriteBufferBoxBased(
				WithWriteBufferBoxBasedMergeSingleBoxes(),
				WithWriteBufferBoxBasedOmitEmptyBoxes(),
				WithWriteBufferBoxBasedDesiredWidth(200),
				WithWriteBufferBoxBasedPrintPosLengthFooter(),
			)
			if err := plcMessage.SerializeWithWriteBuffer(context.Background(), boxWriteBuffer); err != nil {
				return nil, errors.Wrap(err, "error serializing PLCMessage")
			}
			boxString := boxWriteBuffer.GetBox().String()
			return json.Marshal(boxString)
		case PLCMessageAsXML:
			xmlWriteBuffer := NewXmlWriteBuffer(
				WithXmlWriteBufferDefaultIdent(false),
			)
			if err := plcMessage.SerializeWithWriteBuffer(context.Background(), xmlWriteBuffer); err != nil {
				return nil, errors.Wrap(err, "error serializing PLCMessage")
			}
			xmlString := xmlWriteBuffer.GetXmlString()
			return json.Marshal(xmlString)
		case PLCMessageAsByte:
			byteWriteBuffer := NewWriteBufferByteBased()
			if err := plcMessage.SerializeWithWriteBuffer(context.Background(), byteWriteBuffer); err != nil {
				return nil, errors.Wrap(err, "error serializing PLCMessage")
			}
			bytes := byteWriteBuffer.GetBytes()
			return json.Marshal(bytes)
		}
	}
	return ZerologDefaultInterfaceMarshalFunc(v)
}
