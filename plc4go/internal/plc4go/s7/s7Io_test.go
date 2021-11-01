/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package s7

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"reflect"
	"strings"
	"testing"
)

func TestS7MessageBytes(t *testing.T) {
	type debuggable interface {
		utils.Serializable
		fmt.Stringer
	}
	type args struct {
		debuggable debuggable
	}
	tests := []struct {
		name                        string
		args                        args
		wantStringSerialized        string
		wantStringSerializedCompact string
		wantStringXml               string
		wantStringJson              string
		wantDump                    string
	}{
		{
			name: "TPKT Packet with Read var response data",
			args: args{
				debuggable: model.NewTPKTPacket(
					model.NewCOTPPacketData(
						false,
						13,
						[]*model.COTPParameter{model.NewCOTPParameterTpduSize(model.COTPTpduSize_SIZE_4096)},
						model.NewS7MessageResponseData(
							0,
							0,
							11,
							model.NewS7ParameterReadVarResponse(1),
							model.NewS7PayloadReadVarResponse(
								[]*model.S7VarPayloadDataItem{
									model.NewS7VarPayloadDataItem(
										model.DataTransportErrorCode_OK,
										model.DataTransportSize_BIT,
										[]byte{1},
									),
								},
							),
						),
					),
				),
			},
			wantStringSerialized: `
╔═TPKTPacket═══════════════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═protocolId╗╔═reserved╗╔═len═════╗                                                                               ║
║║  0x03 3   ║║ 0x00 0  ║║0x001e 30║                                                                               ║
║╚═══════════╝╚═════════╝╚═════════╝                                                                               ║
║╔═payload════════════════════════════════════════════════════════════════════════════════════════════════════════╗║
║║╔═COTPPacket═══════════════════════════════════════════════════════════════════════════════════════════════════╗║║
║║║╔═headerLength╗╔═tpduCode╗╔═COTPPacketData═════╗╔═parameters═════════════════════════════════════════════════╗║║║
║║║║   0x05 5    ║║0xf0 240 ║║╔═eot════╗╔═tpduRef╗║║╔═COTPParameter════════════════════════════════════════════╗║║║║
║║║╚═════════════╝╚═════════╝║║b0 false║║ 0xd 13 ║║║║╔═parameterType╗╔═parameterLength╗╔═COTPParameterTpduSize╗║║║║║
║║║                          ║╚════════╝╚════════╝║║║║   0xc0 192   ║║     0x01 1     ║║╔═tpduSize══════════╗ ║║║║║║
║║║                          ╚════════════════════╝║║╚══════════════╝╚════════════════╝║║╔═COTPTpduSize════╗║ ║║║║║║
║║║                                                ║║                                  ║║║0x0c 12 SIZE_4096║║ ║║║║║║
║║║                                                ║║                                  ║║╚═════════════════╝║ ║║║║║║
║║║                                                ║║                                  ║╚═══════════════════╝ ║║║║║║
║║║                                                ║║                                  ╚══════════════════════╝║║║║║
║║║                                                ║╚══════════════════════════════════════════════════════════╝║║║║
║║║                                                ╚════════════════════════════════════════════════════════════╝║║║
║║║╔═payload══════════════════════════════════════════════════════════════════════════════════╗                  ║║║
║║║║╔═S7Message══════════════════════════════════════════════════════════════════════════════╗║                  ║║║
║║║║║╔═protocolId╗╔═messageType╗╔═reserved╗╔═tpduReference╗╔═parameterLength╗╔═payloadLength╗║║                  ║║║
║║║║║║  0x32 50  ║║   0x03 3   ║║0x0000 0 ║║  0x000b 11   ║║    0x0002 2    ║║   0x0006 6   ║║║                  ║║║
║║║║║╚═══════════╝╚════════════╝╚═════════╝╚══════════════╝╚════════════════╝╚══════════════╝║║                  ║║║
║║║║║╔═S7MessageResponseData═══╗╔═parameter═════════════════════════════════════╗            ║║                  ║║║
║║║║║║╔═errorClass╗╔═errorCode╗║║╔═S7Parameter═════════════════════════════════╗║            ║║                  ║║║
║║║║║║║  0x00 0   ║║  0x00 0  ║║║║╔═parameterType╗╔═S7ParameterReadVarResponse╗║║            ║║                  ║║║
║║║║║║╚═══════════╝╚══════════╝║║║║    0x04 4    ║║        ╔═numItems╗        ║║║            ║║                  ║║║
║║║║║╚═════════════════════════╝║║╚══════════════╝║        ║ 0x01 1  ║        ║║║            ║║                  ║║║
║║║║║                           ║║                ║        ╚═════════╝        ║║║            ║║                  ║║║
║║║║║                           ║║                ╚═══════════════════════════╝║║            ║║                  ║║║
║║║║║                           ║╚═════════════════════════════════════════════╝║            ║║                  ║║║
║║║║║                           ╚═══════════════════════════════════════════════╝            ║║                  ║║║
║║║║║╔═payload══════════════════════════════════════════════════════════════╗                ║║                  ║║║
║║║║║║╔═S7Payload══════════════════════════════════════════════════════════╗║                ║║                  ║║║
║║║║║║║╔═S7PayloadReadVarResponse═════════════════════════════════════════╗║║                ║║                  ║║║
║║║║║║║║╔═items══════════════════════════════════════════════════════════╗║║║                ║║                  ║║║
║║║║║║║║║╔═S7VarPayloadDataItem═════════════════════════════════════════╗║║║║                ║║                  ║║║
║║║║║║║║║║╔═returnCode══════════════╗╔═transportSize══════╗╔═dataLength╗║║║║║                ║║                  ║║║
║║║║║║║║║║║╔═DataTransportErrorCode╗║║╔═DataTransportSize╗║║ 0x0001 1  ║║║║║║                ║║                  ║║║
║║║║║║║║║║║║      0xff 255 OK      ║║║║    0x03 3 BIT    ║║╚═══════════╝║║║║║                ║║                  ║║║
║║║║║║║║║║║╚═══════════════════════╝║║╚══════════════════╝║             ║║║║║                ║║                  ║║║
║║║║║║║║║║╚═════════════════════════╝╚════════════════════╝             ║║║║║                ║║                  ║║║
║║║║║║║║║║╔═data═══════════════════════════════════════╗╔═padding╗      ║║║║║                ║║                  ║║║
║║║║║║║║║║║0|01                            '.         '║║╔══════╗║      ║║║║║                ║║                  ║║║
║║║║║║║║║║╚════════════════════════════════════════════╝║║0x00 0║║      ║║║║║                ║║                  ║║║
║║║║║║║║║║                                              ║╚══════╝║      ║║║║║                ║║                  ║║║
║║║║║║║║║║                                              ╚════════╝      ║║║║║                ║║                  ║║║
║║║║║║║║║╚══════════════════════════════════════════════════════════════╝║║║║                ║║                  ║║║
║║║║║║║║╚════════════════════════════════════════════════════════════════╝║║║                ║║                  ║║║
║║║║║║║╚══════════════════════════════════════════════════════════════════╝║║                ║║                  ║║║
║║║║║║╚════════════════════════════════════════════════════════════════════╝║                ║║                  ║║║
║║║║║╚══════════════════════════════════════════════════════════════════════╝                ║║                  ║║║
║║║║╚════════════════════════════════════════════════════════════════════════════════════════╝║                  ║║║
║║║╚══════════════════════════════════════════════════════════════════════════════════════════╝                  ║║║
║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║
║╚════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
			wantStringSerializedCompact: `
╔═TPKTPacket═════════════════════════════════════════════════════════════════════════════════════╗
║╔═protocolId╗╔═reserved╗╔═len═════╗                                                             ║
║║  0x03 3   ║║ 0x00 0  ║║0x001e 30║                                                             ║
║╚═══════════╝╚═════════╝╚═════════╝                                                             ║
║╔═payload/COTPPacket═══════════════════════════════════════════════════════════════════════════╗║
║║╔═headerLength╗╔═tpduCode╗╔═COTPPacketData═════╗                                              ║║
║║║   0x05 5    ║║0xf0 240 ║║╔═eot════╗╔═tpduRef╗║                                              ║║
║║╚═════════════╝╚═════════╝║║b0 false║║ 0xd 13 ║║                                              ║║
║║                          ║╚════════╝╚════════╝║                                              ║║
║║                          ╚════════════════════╝                                              ║║
║║╔═parameters/COTPParameter════════════════════════════════════════════════════════════╗       ║║
║║║╔═parameterType╗╔═parameterLength╗╔═COTPParameterTpduSize/tpduSize/COTPTpduSize═════╗║       ║║
║║║║   0xc0 192   ║║     0x01 1     ║║                0x0c 12 SIZE_4096                ║║       ║║
║║║╚══════════════╝╚════════════════╝╚═════════════════════════════════════════════════╝║       ║║
║║╚═════════════════════════════════════════════════════════════════════════════════════╝       ║║
║║╔═payload/S7Message══════════════════════════════════════════════════════════════════════════╗║║
║║║╔═protocolId╗╔═messageType╗╔═reserved╗╔═tpduReference╗╔═parameterLength╗╔═payloadLength╗    ║║║
║║║║  0x32 50  ║║   0x03 3   ║║0x0000 0 ║║  0x000b 11   ║║    0x0002 2    ║║   0x0006 6   ║    ║║║
║║║╚═══════════╝╚════════════╝╚═════════╝╚══════════════╝╚════════════════╝╚══════════════╝    ║║║
║║║╔═S7MessageResponseData═══╗╔═parameter/S7Parameter═══════════════════════════════════╗      ║║║
║║║║╔═errorClass╗╔═errorCode╗║║╔═parameterType╗╔═S7ParameterReadVarResponse/numItems═══╗║      ║║║
║║║║║  0x00 0   ║║  0x00 0  ║║║║    0x04 4    ║║                0x01 1                 ║║      ║║║
║║║║╚═══════════╝╚══════════╝║║╚══════════════╝╚═══════════════════════════════════════╝║      ║║║
║║║╚═════════════════════════╝╚═════════════════════════════════════════════════════════╝      ║║║
║║║╔═payload/S7Payload/S7PayloadReadVarResponse/items/S7VarPayloadDataItem════════════════════╗║║║
║║║║╔═returnCode/DataTransportErrorCode════════════╗╔═transportSize/DataTransportSize════════╗║║║║
║║║║║                 0xff 255 OK                  ║║               0x03 3 BIT               ║║║║║
║║║║╚══════════════════════════════════════════════╝╚════════════════════════════════════════╝║║║║
║║║║╔═dataLength╗╔═data═══════════════════════════════════════╗╔═padding/╗                    ║║║║
║║║║║ 0x0001 1  ║║0|01                            '.         '║║ 0x00 0  ║                    ║║║║
║║║║╚═══════════╝╚════════════════════════════════════════════╝╚═════════╝                    ║║║║
║║║╚══════════════════════════════════════════════════════════════════════════════════════════╝║║║
║║╚════════════════════════════════════════════════════════════════════════════════════════════╝║║
║╚══════════════════════════════════════════════════════════════════════════════════════════════╝║
╚════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
			wantStringXml: `
<TPKTPacket>
  <protocolId dataType="uint" bitLength="8">3</protocolId>
  <reserved dataType="uint" bitLength="8">0</reserved>
  <len dataType="uint" bitLength="16">30</len>
  <payload>
    <COTPPacket>
      <headerLength dataType="uint" bitLength="8">5</headerLength>
      <tpduCode dataType="uint" bitLength="8">240</tpduCode>
      <COTPPacketData>
        <eot dataType="bit" bitLength="1">false</eot>
        <tpduRef dataType="uint" bitLength="7">13</tpduRef>
      </COTPPacketData>
      <parameters isList="true">
        <COTPParameter>
          <parameterType dataType="uint" bitLength="8">192</parameterType>
          <parameterLength dataType="uint" bitLength="8">1</parameterLength>
          <COTPParameterTpduSize>
            <tpduSize>
              <COTPTpduSize dataType="int" bitLength="8" stringRepresentation="SIZE_4096">12</COTPTpduSize>
            </tpduSize>
          </COTPParameterTpduSize>
        </COTPParameter>
      </parameters>
      <payload>
        <S7Message>
          <protocolId dataType="uint" bitLength="8">50</protocolId>
          <messageType dataType="uint" bitLength="8">3</messageType>
          <reserved dataType="uint" bitLength="16">0</reserved>
          <tpduReference dataType="uint" bitLength="16">11</tpduReference>
          <parameterLength dataType="uint" bitLength="16">2</parameterLength>
          <payloadLength dataType="uint" bitLength="16">6</payloadLength>
          <S7MessageResponseData>
            <errorClass dataType="uint" bitLength="8">0</errorClass>
            <errorCode dataType="uint" bitLength="8">0</errorCode>
          </S7MessageResponseData>
          <parameter>
            <S7Parameter>
              <parameterType dataType="uint" bitLength="8">4</parameterType>
              <S7ParameterReadVarResponse>
                <numItems dataType="uint" bitLength="8">1</numItems>
              </S7ParameterReadVarResponse>
            </S7Parameter>
          </parameter>
          <payload>
            <S7Payload>
              <S7PayloadReadVarResponse>
                <items isList="true">
                  <S7VarPayloadDataItem>
                    <returnCode>
                      <DataTransportErrorCode dataType="uint" bitLength="8" stringRepresentation="OK">255</DataTransportErrorCode>
                    </returnCode>
                    <transportSize>
                      <DataTransportSize dataType="uint" bitLength="8" stringRepresentation="BIT">3</DataTransportSize>
                    </transportSize>
                    <dataLength dataType="uint" bitLength="16">1</dataLength>
                    <data dataType="byte" bitLength="8">0x01</data>
                    <padding isList="true">
                      <value dataType="uint" bitLength="8">0</value>
                    </padding>
                  </S7VarPayloadDataItem>
                </items>
              </S7PayloadReadVarResponse>
            </S7Payload>
          </payload>
        </S7Message>
      </payload>
    </COTPPacket>
  </payload>
</TPKTPacket>
`,
			wantStringJson: `
{
  "TPKTPacket": {
    "len": 30,
    "len__plc4x_bitLength": 16,
    "len__plc4x_dataType": "uint",
    "payload": {
      "COTPPacket": {
        "COTPPacketData": {
          "eot": false,
          "eot__plc4x_bitLength": 1,
          "eot__plc4x_dataType": "bit",
          "tpduRef": 13,
          "tpduRef__plc4x_bitLength": 7,
          "tpduRef__plc4x_dataType": "uint"
        },
        "headerLength": 5,
        "headerLength__plc4x_bitLength": 8,
        "headerLength__plc4x_dataType": "uint",
        "parameters": [
          {
            "COTPParameter": {
              "COTPParameterTpduSize": {
                "tpduSize": {
                  "COTPTpduSize": 12,
                  "COTPTpduSize__plc4x_bitLength": 8,
                  "COTPTpduSize__plc4x_dataType": "int",
                  "COTPTpduSize__plc4x_stringRepresentation": "SIZE_4096"
                }
              },
              "parameterLength": 1,
              "parameterLength__plc4x_bitLength": 8,
              "parameterLength__plc4x_dataType": "uint",
              "parameterType": 192,
              "parameterType__plc4x_bitLength": 8,
              "parameterType__plc4x_dataType": "uint"
            }
          }
        ],
        "payload": {
          "S7Message": {
            "S7MessageResponseData": {
              "errorClass": 0,
              "errorClass__plc4x_bitLength": 8,
              "errorClass__plc4x_dataType": "uint",
              "errorCode": 0,
              "errorCode__plc4x_bitLength": 8,
              "errorCode__plc4x_dataType": "uint"
            },
            "messageType": 3,
            "messageType__plc4x_bitLength": 8,
            "messageType__plc4x_dataType": "uint",
            "parameter": {
              "S7Parameter": {
                "S7ParameterReadVarResponse": {
                  "numItems": 1,
                  "numItems__plc4x_bitLength": 8,
                  "numItems__plc4x_dataType": "uint"
                },
                "parameterType": 4,
                "parameterType__plc4x_bitLength": 8,
                "parameterType__plc4x_dataType": "uint"
              }
            },
            "parameterLength": 2,
            "parameterLength__plc4x_bitLength": 16,
            "parameterLength__plc4x_dataType": "uint",
            "payload": {
              "S7Payload": {
                "S7PayloadReadVarResponse": {
                  "items": [
                    {
                      "S7VarPayloadDataItem": {
                        "data": "0x01",
                        "dataLength": 1,
                        "dataLength__plc4x_bitLength": 16,
                        "dataLength__plc4x_dataType": "uint",
                        "data__plc4x_bitLength": 8,
                        "data__plc4x_dataType": "byte",
                        "padding": [
                          {
                            "value": 0,
                            "value__plc4x_bitLength": 8,
                            "value__plc4x_dataType": "uint"
                          }
                        ],
                        "returnCode": {
                          "DataTransportErrorCode": 255,
                          "DataTransportErrorCode__plc4x_bitLength": 8,
                          "DataTransportErrorCode__plc4x_dataType": "uint",
                          "DataTransportErrorCode__plc4x_stringRepresentation": "OK"
                        },
                        "transportSize": {
                          "DataTransportSize": 3,
                          "DataTransportSize__plc4x_bitLength": 8,
                          "DataTransportSize__plc4x_dataType": "uint",
                          "DataTransportSize__plc4x_stringRepresentation": "BIT"
                        }
                      }
                    }
                  ]
                }
              }
            },
            "payloadLength": 6,
            "payloadLength__plc4x_bitLength": 16,
            "payloadLength__plc4x_dataType": "uint",
            "protocolId": 50,
            "protocolId__plc4x_bitLength": 8,
            "protocolId__plc4x_dataType": "uint",
            "reserved": 0,
            "reserved__plc4x_bitLength": 16,
            "reserved__plc4x_dataType": "uint",
            "tpduReference": 11,
            "tpduReference__plc4x_bitLength": 16,
            "tpduReference__plc4x_dataType": "uint"
          }
        },
        "tpduCode": 240,
        "tpduCode__plc4x_bitLength": 8,
        "tpduCode__plc4x_dataType": "uint"
      }
    },
    "protocolId": 3,
    "protocolId__plc4x_bitLength": 8,
    "protocolId__plc4x_dataType": "uint",
    "reserved": 0,
    "reserved__plc4x_bitLength": 8,
    "reserved__plc4x_dataType": "uint"
  }
}
`,
			wantDump: `
00|03 00 00 1e 05 f0 0d c0 01 0c '..........'
10|32 03 00 00 00 0b 00 02 00 06 '2.........'
20|00 00 04 01 ff 03 00 01 01 00 '..........'
`,
		},
		{
			name: "TPKT Packet with write var request data",
			args: args{
				debuggable: model.NewTPKTPacket(
					model.NewCOTPPacketData(
						false,
						13,
						[]*model.COTPParameter{model.NewCOTPParameterTpduSize(model.COTPTpduSize_SIZE_4096)},
						model.NewS7MessageRequest(
							13,
							model.NewS7ParameterWriteVarRequest([]*model.S7VarRequestParameterItem{
								model.NewS7VarRequestParameterItemAddress(model.NewS7AddressAny(
									model.TransportSize_BYTE,
									64,
									13,
									model.MemoryArea_INPUTS,
									0,
									0,
								)),
							}),
							model.NewS7PayloadWriteVarRequest([]*model.S7VarPayloadDataItem{
								model.NewS7VarPayloadDataItem(
									model.DataTransportErrorCode_OK,
									model.DataTransportSize_BYTE_WORD_DWORD,
									[]byte{
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
										0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE, 0xAF, 0xFE,
									},
								),
							}),
						),
					),
				),
			},
			wantStringSerialized: `
╔═TPKTPacket═══════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═protocolId╗╔═reserved╗╔═len══════╗                                                                                  ║
║║  0x03 3   ║║ 0x00 0  ║║0x0066 102║                                                                                  ║
║╚═══════════╝╚═════════╝╚══════════╝                                                                                  ║
║╔═payload════════════════════════════════════════════════════════════════════════════════════════════════════════════╗║
║║╔═COTPPacket═══════════════════════════════════════════════════════════════════════════════════════════════════════╗║║
║║║╔═headerLength╗╔═tpduCode╗╔═COTPPacketData═════╗╔═parameters═════════════════════════════════════════════════╗    ║║║
║║║║   0x05 5    ║║0xf0 240 ║║╔═eot════╗╔═tpduRef╗║║╔═COTPParameter════════════════════════════════════════════╗║    ║║║
║║║╚═════════════╝╚═════════╝║║b0 false║║ 0xd 13 ║║║║╔═parameterType╗╔═parameterLength╗╔═COTPParameterTpduSize╗║║    ║║║
║║║                          ║╚════════╝╚════════╝║║║║   0xc0 192   ║║     0x01 1     ║║╔═tpduSize══════════╗ ║║║    ║║║
║║║                          ╚════════════════════╝║║╚══════════════╝╚════════════════╝║║╔═COTPTpduSize════╗║ ║║║    ║║║
║║║                                                ║║                                  ║║║0x0c 12 SIZE_4096║║ ║║║    ║║║
║║║                                                ║║                                  ║║╚═════════════════╝║ ║║║    ║║║
║║║                                                ║║                                  ║╚═══════════════════╝ ║║║    ║║║
║║║                                                ║║                                  ╚══════════════════════╝║║    ║║║
║║║                                                ║╚══════════════════════════════════════════════════════════╝║    ║║║
║║║                                                ╚════════════════════════════════════════════════════════════╝    ║║║
║║║╔═payload════════════════════════════════════════════════════════════════════════════════════════════════════════╗║║║
║║║║╔═S7Message════════════════════════════════════════════════════════════════════════════════════════════════════╗║║║║
║║║║║╔═protocolId╗╔═messageType╗╔═reserved╗╔═tpduReference╗╔═parameterLength╗╔═payloadLength╗╔═S7MessageRequest╗   ║║║║║
║║║║║║  0x32 50  ║║   0x01 1   ║║0x0000 0 ║║  0x000d 13   ║║   0x000e 14    ║║  0x0044 68   ║║                 ║   ║║║║║
║║║║║╚═══════════╝╚════════════╝╚═════════╝╚══════════════╝╚════════════════╝╚══════════════╝╚═════════════════╝   ║║║║║
║║║║║╔═parameter══════════════════════════════════════════════════════════════════════════════════════════════════╗║║║║║
║║║║║║╔═S7Parameter══════════════════════════════════════════════════════════════════════════════════════════════╗║║║║║║
║║║║║║║╔═parameterType╗                                                                                          ║║║║║║║
║║║║║║║║    0x05 5    ║                                                                                          ║║║║║║║
║║║║║║║╚══════════════╝                                                                                          ║║║║║║║
║║║║║║║╔═S7ParameterWriteVarRequest═════════════════════════════════════════════════════════════════════════════╗║║║║║║║
║║║║║║║║╔═numItems╗                                                                                             ║║║║║║║║
║║║║║║║║║ 0x01 1  ║                                                                                             ║║║║║║║║
║║║║║║║║╚═════════╝                                                                                             ║║║║║║║║
║║║║║║║║╔═items════════════════════════════════════════════════════════════════════════════════════════════════╗║║║║║║║║
║║║║║║║║║╔═S7VarRequestParameterItem══════════════════════════════════════════════════════════════════════════╗║║║║║║║║║
║║║║║║║║║║╔═itemType╗                                                                                         ║║║║║║║║║║
║║║║║║║║║║║ 0x12 18 ║                                                                                         ║║║║║║║║║║
║║║║║║║║║║╚═════════╝                                                                                         ║║║║║║║║║║
║║║║║║║║║║╔═S7VarRequestParameterItemAddress═════════════════════════════════════════════════════════════════╗║║║║║║║║║║
║║║║║║║║║║║╔═itemLength╗                                                                                     ║║║║║║║║║║║
║║║║║║║║║║║║  0x0a 10  ║                                                                                     ║║║║║║║║║║║
║║║║║║║║║║║╚═══════════╝                                                                                     ║║║║║║║║║║║
║║║║║║║║║║║╔═address════════════════════════════════════════════════════════════════════════════════════════╗║║║║║║║║║║║
║║║║║║║║║║║║╔═S7Address════════════════════════════════════════════════════════════════════════════════════╗║║║║║║║║║║║║
║║║║║║║║║║║║║╔═addressType╗                                                                                ║║║║║║║║║║║║║
║║║║║║║║║║║║║║  0x10 16   ║                                                                                ║║║║║║║║║║║║║
║║║║║║║║║║║║║╚════════════╝                                                                                ║║║║║║║║║║║║║
║║║║║║║║║║║║║╔═S7AddressAny═══════════════════════════════════════════════════════════════════════════════╗║║║║║║║║║║║║║
║║║║║║║║║║║║║║╔═transportSize══╗╔═numberOfElements╗╔═dbNumber╗╔═area════════════╗╔═reserved╗╔═byteAddress╗║║║║║║║║║║║║║║
║║║║║║║║║║║║║║║╔═TransportSize╗║║    0x0040 64    ║║0x000d 13║║╔═MemoryArea════╗║║  0x0 0  ║║  0x0000 0  ║║║║║║║║║║║║║║║
║║║║║║║║║║║║║║║║ 0x02 2 BYTE  ║║╚═════════════════╝╚═════════╝║║0x81 129 INPUTS║║╚═════════╝╚════════════╝║║║║║║║║║║║║║║
║║║║║║║║║║║║║║║╚══════════════╝║                              ║╚═══════════════╝║                         ║║║║║║║║║║║║║║
║║║║║║║║║║║║║║╚════════════════╝                              ╚═════════════════╝                         ║║║║║║║║║║║║║║
║║║║║║║║║║║║║║╔═bitAddress╗                                                                               ║║║║║║║║║║║║║║
║║║║║║║║║║║║║║║   0x0 0   ║                                                                               ║║║║║║║║║║║║║║
║║║║║║║║║║║║║║╚═══════════╝                                                                               ║║║║║║║║║║║║║║
║║║║║║║║║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║║║║║║
║║║║║║║║║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║║║║║
║║║║║║║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║║║║
║║║║║║║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║║║
║║║║║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║║
║║║║║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║
║║║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║
║║║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║
║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║
║║║║║╔═payload══════════════════════════════════════════════════════════════════╗                                  ║║║║║
║║║║║║╔═S7Payload══════════════════════════════════════════════════════════════╗║                                  ║║║║║
║║║║║║║╔═S7PayloadWriteVarRequest═════════════════════════════════════════════╗║║                                  ║║║║║
║║║║║║║║╔═items══════════════════════════════════════════════════════════════╗║║║                                  ║║║║║
║║║║║║║║║╔═S7VarPayloadDataItem═════════════════════════════════════════════╗║║║║                                  ║║║║║
║║║║║║║║║║╔═returnCode══════════════╗╔═transportSize══════════╗╔═dataLength╗║║║║║                                  ║║║║║
║║║║║║║║║║║╔═DataTransportErrorCode╗║║╔═DataTransportSize════╗║║0x0200 512 ║║║║║║                                  ║║║║║
║║║║║║║║║║║║      0xff 255 OK      ║║║║0x04 4 BYTE_WORD_DWORD║║╚═══════════╝║║║║║                                  ║║║║║
║║║║║║║║║║║╚═══════════════════════╝║║╚══════════════════════╝║             ║║║║║                                  ║║║║║
║║║║║║║║║║╚═════════════════════════╝╚════════════════════════╝             ║║║║║                                  ║║║║║
║║║║║║║║║║╔═data════════════════════════════════════════╗╔═padding╗         ║║║║║                                  ║║║║║
║║║║║║║║║║║00|af fe af fe af fe af fe af fe '..........'║║        ║         ║║║║║                                  ║║║║║
║║║║║║║║║║║10|af fe af fe af fe af fe af fe '..........'║╚════════╝         ║║║║║                                  ║║║║║
║║║║║║║║║║║20|af fe af fe af fe af fe af fe '..........'║                   ║║║║║                                  ║║║║║
║║║║║║║║║║║30|af fe af fe af fe af fe af fe '..........'║                   ║║║║║                                  ║║║║║
║║║║║║║║║║║40|af fe af fe af fe af fe af fe '..........'║                   ║║║║║                                  ║║║║║
║║║║║║║║║║║50|af fe af fe af fe af fe af fe '..........'║                   ║║║║║                                  ║║║║║
║║║║║║║║║║║60|af fe af fe                   '....      '║                   ║║║║║                                  ║║║║║
║║║║║║║║║║╚═════════════════════════════════════════════╝                   ║║║║║                                  ║║║║║
║║║║║║║║║╚══════════════════════════════════════════════════════════════════╝║║║║                                  ║║║║║
║║║║║║║║╚════════════════════════════════════════════════════════════════════╝║║║                                  ║║║║║
║║║║║║║╚══════════════════════════════════════════════════════════════════════╝║║                                  ║║║║║
║║║║║║╚════════════════════════════════════════════════════════════════════════╝║                                  ║║║║║
║║║║║╚══════════════════════════════════════════════════════════════════════════╝                                  ║║║║║
║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║
║║║╚════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║
║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║
║╚════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
			wantStringSerializedCompact: `
╔═TPKTPacket═════════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═protocolId╗╔═reserved╗╔═len══════╗                                                                        ║
║║  0x03 3   ║║ 0x00 0  ║║0x0066 102║                                                                        ║
║╚═══════════╝╚═════════╝╚══════════╝                                                                        ║
║╔═payload/COTPPacket═══════════════════════════════════════════════════════════════════════════════════════╗║
║║╔═headerLength╗╔═tpduCode╗╔═COTPPacketData═════╗                                                          ║║
║║║   0x05 5    ║║0xf0 240 ║║╔═eot════╗╔═tpduRef╗║                                                          ║║
║║╚═════════════╝╚═════════╝║║b0 false║║ 0xd 13 ║║                                                          ║║
║║                          ║╚════════╝╚════════╝║                                                          ║║
║║                          ╚════════════════════╝                                                          ║║
║║╔═parameters/COTPParameter════════════════════════════════════════════════════════════╗                   ║║
║║║╔═parameterType╗╔═parameterLength╗╔═COTPParameterTpduSize/tpduSize/COTPTpduSize═════╗║                   ║║
║║║║   0xc0 192   ║║     0x01 1     ║║                0x0c 12 SIZE_4096                ║║                   ║║
║║║╚══════════════╝╚════════════════╝╚═════════════════════════════════════════════════╝║                   ║║
║║╚═════════════════════════════════════════════════════════════════════════════════════╝                   ║║
║║╔═payload/S7Message══════════════════════════════════════════════════════════════════════════════════════╗║║
║║║╔═protocolId╗╔═messageType╗╔═reserved╗╔═tpduReference╗╔═parameterLength╗╔═payloadLength╗                ║║║
║║║║  0x32 50  ║║   0x01 1   ║║0x0000 0 ║║  0x000d 13   ║║   0x000e 14    ║║  0x0044 68   ║                ║║║
║║║╚═══════════╝╚════════════╝╚═════════╝╚══════════════╝╚════════════════╝╚══════════════╝                ║║║
║║║╔═parameter/S7Parameter════════════════════════════════════════════════════════════════════════════════╗║║║
║║║║╔═parameterType╗                                                                                      ║║║║
║║║║║    0x05 5    ║                                                                                      ║║║║
║║║║╚══════════════╝                                                                                      ║║║║
║║║║╔═S7ParameterWriteVarRequest═════════════════════════════════════════════════════════════════════════╗║║║║
║║║║║╔═numItems╗                                                                                         ║║║║║
║║║║║║ 0x01 1  ║                                                                                         ║║║║║
║║║║║╚═════════╝                                                                                         ║║║║║
║║║║║╔═items/S7VarRequestParameterItem══════════════════════════════════════════════════════════════════╗║║║║║
║║║║║║╔═itemType╗                                                                                       ║║║║║║
║║║║║║║ 0x12 18 ║                                                                                       ║║║║║║
║║║║║║╚═════════╝                                                                                       ║║║║║║
║║║║║║╔═S7VarRequestParameterItemAddress═══════════════════════════════════════════════════════════════╗║║║║║║
║║║║║║║╔═itemLength╗                                                                                   ║║║║║║║
║║║║║║║║  0x0a 10  ║                                                                                   ║║║║║║║
║║║║║║║╚═══════════╝                                                                                   ║║║║║║║
║║║║║║║╔═address/S7Address════════════════════════════════════════════════════════════════════════════╗║║║║║║║
║║║║║║║║╔═addressType╗                                                                                ║║║║║║║║
║║║║║║║║║  0x10 16   ║                                                                                ║║║║║║║║
║║║║║║║║╚════════════╝                                                                                ║║║║║║║║
║║║║║║║║╔═S7AddressAny═══════════════════════════════════════════════════════════════════════════════╗║║║║║║║║
║║║║║║║║║╔═transportSize/TransportSize═══╗╔═numberOfElements╗╔═dbNumber╗╔═area/MemoryArea╗╔═reserved╗║║║║║║║║║
║║║║║║║║║║          0x02 2 BYTE          ║║    0x0040 64    ║║0x000d 13║║0x81 129 INPUTS ║║  0x0 0  ║║║║║║║║║║
║║║║║║║║║╚═══════════════════════════════╝╚═════════════════╝╚═════════╝╚════════════════╝╚═════════╝║║║║║║║║║
║║║║║║║║║╔═byteAddress╗╔═bitAddress╗                                                                 ║║║║║║║║║
║║║║║║║║║║  0x0000 0  ║║   0x0 0   ║                                                                 ║║║║║║║║║
║║║║║║║║║╚════════════╝╚═══════════╝                                                                 ║║║║║║║║║
║║║║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║║
║║║║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║║
║║║║║║╚════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║║
║║║║║╚══════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║║
║║║║╚════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║║
║║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════╝║║║
║║║╔═payload/S7Payload/S7PayloadWriteVarRequest/items/S7VarPayloadDataItem═════════════════════════╗       ║║║
║║║║╔═returnCode/DataTransportErrorCode════════════╗╔═transportSize/DataTransportSize╗╔═dataLength╗║       ║║║
║║║║║                 0xff 255 OK                  ║║     0x04 4 BYTE_WORD_DWORD     ║║0x0200 512 ║║       ║║║
║║║║╚══════════════════════════════════════════════╝╚════════════════════════════════╝╚═══════════╝║       ║║║
║║║║╔═data════════════════════════════════════════╗                                                ║       ║║║
║║║║║00|af fe af fe af fe af fe af fe '..........'║                                                ║       ║║║
║║║║║10|af fe af fe af fe af fe af fe '..........'║                                                ║       ║║║
║║║║║20|af fe af fe af fe af fe af fe '..........'║                                                ║       ║║║
║║║║║30|af fe af fe af fe af fe af fe '..........'║                                                ║       ║║║
║║║║║40|af fe af fe af fe af fe af fe '..........'║                                                ║       ║║║
║║║║║50|af fe af fe af fe af fe af fe '..........'║                                                ║       ║║║
║║║║║60|af fe af fe                   '....      '║                                                ║       ║║║
║║║║╚═════════════════════════════════════════════╝                                                ║       ║║║
║║║╚═══════════════════════════════════════════════════════════════════════════════════════════════╝       ║║║
║║╚════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║
║╚══════════════════════════════════════════════════════════════════════════════════════════════════════════╝║
╚════════════════════════════════════════════════════════════════════════════════════════════════════════════╝
`,
			wantStringXml: `
<TPKTPacket>
  <protocolId dataType="uint" bitLength="8">3</protocolId>
  <reserved dataType="uint" bitLength="8">0</reserved>
  <len dataType="uint" bitLength="16">102</len>
  <payload>
    <COTPPacket>
      <headerLength dataType="uint" bitLength="8">5</headerLength>
      <tpduCode dataType="uint" bitLength="8">240</tpduCode>
      <COTPPacketData>
        <eot dataType="bit" bitLength="1">false</eot>
        <tpduRef dataType="uint" bitLength="7">13</tpduRef>
      </COTPPacketData>
      <parameters isList="true">
        <COTPParameter>
          <parameterType dataType="uint" bitLength="8">192</parameterType>
          <parameterLength dataType="uint" bitLength="8">1</parameterLength>
          <COTPParameterTpduSize>
            <tpduSize>
              <COTPTpduSize dataType="int" bitLength="8" stringRepresentation="SIZE_4096">12</COTPTpduSize>
            </tpduSize>
          </COTPParameterTpduSize>
        </COTPParameter>
      </parameters>
      <payload>
        <S7Message>
          <protocolId dataType="uint" bitLength="8">50</protocolId>
          <messageType dataType="uint" bitLength="8">1</messageType>
          <reserved dataType="uint" bitLength="16">0</reserved>
          <tpduReference dataType="uint" bitLength="16">13</tpduReference>
          <parameterLength dataType="uint" bitLength="16">14</parameterLength>
          <payloadLength dataType="uint" bitLength="16">68</payloadLength>
          <S7MessageRequest></S7MessageRequest>
          <parameter>
            <S7Parameter>
              <parameterType dataType="uint" bitLength="8">5</parameterType>
              <S7ParameterWriteVarRequest>
                <numItems dataType="uint" bitLength="8">1</numItems>
                <items isList="true">
                  <S7VarRequestParameterItem>
                    <itemType dataType="uint" bitLength="8">18</itemType>
                    <S7VarRequestParameterItemAddress>
                      <itemLength dataType="uint" bitLength="8">10</itemLength>
                      <address>
                        <S7Address>
                          <addressType dataType="uint" bitLength="8">16</addressType>
                          <S7AddressAny>
                            <transportSize>
                              <TransportSize dataType="uint" bitLength="8" stringRepresentation="BYTE">2</TransportSize>
                            </transportSize>
                            <numberOfElements dataType="uint" bitLength="16">64</numberOfElements>
                            <dbNumber dataType="uint" bitLength="16">13</dbNumber>
                            <area>
                              <MemoryArea dataType="uint" bitLength="8" stringRepresentation="INPUTS">129</MemoryArea>
                            </area>
                            <reserved dataType="uint" bitLength="5">0</reserved>
                            <byteAddress dataType="uint" bitLength="16">0</byteAddress>
                            <bitAddress dataType="uint" bitLength="3">0</bitAddress>
                          </S7AddressAny>
                        </S7Address>
                      </address>
                    </S7VarRequestParameterItemAddress>
                  </S7VarRequestParameterItem>
                </items>
              </S7ParameterWriteVarRequest>
            </S7Parameter>
          </parameter>
          <payload>
            <S7Payload>
              <S7PayloadWriteVarRequest>
                <items isList="true">
                  <S7VarPayloadDataItem>
                    <returnCode>
                      <DataTransportErrorCode dataType="uint" bitLength="8" stringRepresentation="OK">255</DataTransportErrorCode>
                    </returnCode>
                    <transportSize>
                      <DataTransportSize dataType="uint" bitLength="8" stringRepresentation="BYTE_WORD_DWORD">4</DataTransportSize>
                    </transportSize>
                    <dataLength dataType="uint" bitLength="16">512</dataLength>
                    <data dataType="byte" bitLength="512">0xaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffe</data>
                    <padding isList="true"></padding>
                  </S7VarPayloadDataItem>
                </items>
              </S7PayloadWriteVarRequest>
            </S7Payload>
          </payload>
        </S7Message>
      </payload>
    </COTPPacket>
  </payload>
</TPKTPacket>
`,
			wantStringJson: `
{
  "TPKTPacket": {
    "len": 102,
    "len__plc4x_bitLength": 16,
    "len__plc4x_dataType": "uint",
    "payload": {
      "COTPPacket": {
        "COTPPacketData": {
          "eot": false,
          "eot__plc4x_bitLength": 1,
          "eot__plc4x_dataType": "bit",
          "tpduRef": 13,
          "tpduRef__plc4x_bitLength": 7,
          "tpduRef__plc4x_dataType": "uint"
        },
        "headerLength": 5,
        "headerLength__plc4x_bitLength": 8,
        "headerLength__plc4x_dataType": "uint",
        "parameters": [
          {
            "COTPParameter": {
              "COTPParameterTpduSize": {
                "tpduSize": {
                  "COTPTpduSize": 12,
                  "COTPTpduSize__plc4x_bitLength": 8,
                  "COTPTpduSize__plc4x_dataType": "int",
                  "COTPTpduSize__plc4x_stringRepresentation": "SIZE_4096"
                }
              },
              "parameterLength": 1,
              "parameterLength__plc4x_bitLength": 8,
              "parameterLength__plc4x_dataType": "uint",
              "parameterType": 192,
              "parameterType__plc4x_bitLength": 8,
              "parameterType__plc4x_dataType": "uint"
            }
          }
        ],
        "payload": {
          "S7Message": {
            "S7MessageRequest": {},
            "messageType": 1,
            "messageType__plc4x_bitLength": 8,
            "messageType__plc4x_dataType": "uint",
            "parameter": {
              "S7Parameter": {
                "S7ParameterWriteVarRequest": {
                  "items": [
                    {
                      "S7VarRequestParameterItem": {
                        "S7VarRequestParameterItemAddress": {
                          "address": {
                            "S7Address": {
                              "S7AddressAny": {
                                "area": {
                                  "MemoryArea": 129,
                                  "MemoryArea__plc4x_bitLength": 8,
                                  "MemoryArea__plc4x_dataType": "uint",
                                  "MemoryArea__plc4x_stringRepresentation": "INPUTS"
                                },
                                "bitAddress": 0,
                                "bitAddress__plc4x_bitLength": 3,
                                "bitAddress__plc4x_dataType": "uint",
                                "byteAddress": 0,
                                "byteAddress__plc4x_bitLength": 16,
                                "byteAddress__plc4x_dataType": "uint",
                                "dbNumber": 13,
                                "dbNumber__plc4x_bitLength": 16,
                                "dbNumber__plc4x_dataType": "uint",
                                "numberOfElements": 64,
                                "numberOfElements__plc4x_bitLength": 16,
                                "numberOfElements__plc4x_dataType": "uint",
                                "reserved": 0,
                                "reserved__plc4x_bitLength": 5,
                                "reserved__plc4x_dataType": "uint",
                                "transportSize": {
                                  "TransportSize": 2,
                                  "TransportSize__plc4x_bitLength": 8,
                                  "TransportSize__plc4x_dataType": "uint",
                                  "TransportSize__plc4x_stringRepresentation": "BYTE"
                                }
                              },
                              "addressType": 16,
                              "addressType__plc4x_bitLength": 8,
                              "addressType__plc4x_dataType": "uint"
                            }
                          },
                          "itemLength": 10,
                          "itemLength__plc4x_bitLength": 8,
                          "itemLength__plc4x_dataType": "uint"
                        },
                        "itemType": 18,
                        "itemType__plc4x_bitLength": 8,
                        "itemType__plc4x_dataType": "uint"
                      }
                    }
                  ],
                  "numItems": 1,
                  "numItems__plc4x_bitLength": 8,
                  "numItems__plc4x_dataType": "uint"
                },
                "parameterType": 5,
                "parameterType__plc4x_bitLength": 8,
                "parameterType__plc4x_dataType": "uint"
              }
            },
            "parameterLength": 14,
            "parameterLength__plc4x_bitLength": 16,
            "parameterLength__plc4x_dataType": "uint",
            "payload": {
              "S7Payload": {
                "S7PayloadWriteVarRequest": {
                  "items": [
                    {
                      "S7VarPayloadDataItem": {
                        "data": "0xaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffeaffe",
                        "dataLength": 512,
                        "dataLength__plc4x_bitLength": 16,
                        "dataLength__plc4x_dataType": "uint",
                        "data__plc4x_bitLength": 512,
                        "data__plc4x_dataType": "byte",
                        "padding": [],
                        "returnCode": {
                          "DataTransportErrorCode": 255,
                          "DataTransportErrorCode__plc4x_bitLength": 8,
                          "DataTransportErrorCode__plc4x_dataType": "uint",
                          "DataTransportErrorCode__plc4x_stringRepresentation": "OK"
                        },
                        "transportSize": {
                          "DataTransportSize": 4,
                          "DataTransportSize__plc4x_bitLength": 8,
                          "DataTransportSize__plc4x_dataType": "uint",
                          "DataTransportSize__plc4x_stringRepresentation": "BYTE_WORD_DWORD"
                        }
                      }
                    }
                  ]
                }
              }
            },
            "payloadLength": 68,
            "payloadLength__plc4x_bitLength": 16,
            "payloadLength__plc4x_dataType": "uint",
            "protocolId": 50,
            "protocolId__plc4x_bitLength": 8,
            "protocolId__plc4x_dataType": "uint",
            "reserved": 0,
            "reserved__plc4x_bitLength": 16,
            "reserved__plc4x_dataType": "uint",
            "tpduReference": 13,
            "tpduReference__plc4x_bitLength": 16,
            "tpduReference__plc4x_dataType": "uint"
          }
        },
        "tpduCode": 240,
        "tpduCode__plc4x_bitLength": 8,
        "tpduCode__plc4x_dataType": "uint"
      }
    },
    "protocolId": 3,
    "protocolId__plc4x_bitLength": 8,
    "protocolId__plc4x_dataType": "uint",
    "reserved": 0,
    "reserved__plc4x_bitLength": 8,
    "reserved__plc4x_dataType": "uint"
  }
}
`,
			wantDump: `
000|03 00 00 66 05 f0 0d c0 01 0c '...f......'
010|32 01 00 00 00 0d 00 0e 00 44 '2........D'
020|05 01 12 0a 10 02 00 40 00 0d '.......@..'
030|81 00 00 00 ff 04 02 00 af fe '..........'
040|af fe af fe af fe af fe af fe '..........'
050|af fe af fe af fe af fe af fe '..........'
060|af fe af fe af fe af fe af fe '..........'
070|af fe af fe af fe af fe af fe '..........'
080|af fe af fe af fe af fe af fe '..........'
090|af fe af fe af fe af fe af fe '..........'
100|af fe                         '..        '
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("Simple 2 String", func(t *testing.T) {
				tt.wantStringSerializedCompact = strings.Trim(tt.wantStringSerializedCompact, "\n")
				if got := tt.args.debuggable.String(); got != tt.wantStringSerializedCompact {
					t.Errorf("String() = '\n%v\n', want '\n%v\n'", got, tt.wantStringSerializedCompact)
				}
			})
			t.Run("Simple 2 Box", func(t *testing.T) {
				boxWriter := utils.NewBoxedWriteBuffer()
				if err := tt.args.debuggable.Serialize(boxWriter); err != nil {
					t.Error(err)
				}
				tt.wantStringSerialized = strings.Trim(tt.wantStringSerialized, "\n")
				if got := string(boxWriter.GetBox()); got != tt.wantStringSerialized {
					t.Errorf("Serialize Boxed() = '\n%v\n', want '\n%v\n'", got, tt.wantStringSerialized)
				}
			})
			t.Run("Simple 2 Compact Box", func(t *testing.T) {
				boxWriter := utils.NewBoxedWriteBufferWithOptions(true, true)
				if err := tt.args.debuggable.Serialize(boxWriter); err != nil {
					t.Error(err)
				}
				tt.wantStringSerializedCompact = strings.Trim(tt.wantStringSerializedCompact, "\n")
				if got := string(boxWriter.GetBox()); got != tt.wantStringSerializedCompact {
					t.Errorf("Serialize BoxedCompact() = '\n%v\n', want '\n%v\n'", got, tt.wantStringSerializedCompact)
				}
			})
			t.Run("Simple 2 Xml", func(t *testing.T) {
				xmlWriteBuffer := utils.NewXmlWriteBuffer()
				if err := tt.args.debuggable.Serialize(xmlWriteBuffer); err != nil {
					t.Error(err)
				}
				tt.wantStringXml = strings.Trim(tt.wantStringXml, "\n")
				if got := xmlWriteBuffer.GetXmlString(); got != tt.wantStringXml {
					t.Errorf("Serialize Xml() = '\n%v\n', want '\n%v\n'", got, tt.wantStringXml)
				}
			})
			t.Run("Simple 2 Json", func(t *testing.T) {
				jsonWriteBuffer := utils.NewJsonWriteBuffer()
				if err := tt.args.debuggable.Serialize(jsonWriteBuffer); err != nil {
					t.Error(err)
				}
				tt.wantStringJson = strings.Trim(tt.wantStringJson, "\n")
				if got, err := jsonWriteBuffer.GetJsonString(); err != nil || strings.Trim(got, "\n") != tt.wantStringJson {
					if err != nil {
						t.Error(err)
					} else {
						t.Errorf("Serialize Json() = '\n%v\n', want '\n%v\n'", got, tt.wantStringJson)
					}
				}
			})
			t.Run("Simple Binary Serialize", func(t *testing.T) {
				buffer := utils.NewWriteBufferByteBased()
				if err := tt.args.debuggable.Serialize(buffer); err != nil {
					t.Error(err)
				}
				tt.wantDump = strings.Trim(tt.wantDump, "\n")
				if got := utils.Dump(buffer.GetBytes()); !reflect.DeepEqual(got, tt.wantDump) {
					t.Errorf("Serialize() = '\n%v\n', want '\n%v\n'", got, tt.wantDump)
				}
			})
			t.Run("xml roundtip", func(t *testing.T) {
				reader := strings.NewReader(tt.wantStringXml)
				readBuffer := utils.NewXmlReadBuffer(reader)
				if got, err := model.TPKTPacketParse(readBuffer); err != nil || !reflect.DeepEqual(got, tt.args.debuggable) {
					if err != nil {
						t.Error(err)
					} else {
						t.Errorf("Roundtrip(xml) = '\n%v\n', want '\n%v\n'", got, tt.wantDump)
					}
				}
			})
			t.Run("json roundtip", func(t *testing.T) {
				reader := strings.NewReader(tt.wantStringJson)
				readBuffer := utils.NewJsonReadBuffer(reader)
				if got, err := model.TPKTPacketParse(readBuffer); err != nil || !reflect.DeepEqual(got, tt.args.debuggable) {
					if err != nil {
						t.Error(err)
					} else {
						t.Errorf("Roundtrip(json) = '\n%v\n', want '\n%v\n'", got, tt.args.debuggable)
					}
				}
			})
		})
	}
}
