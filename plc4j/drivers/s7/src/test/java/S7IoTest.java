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

import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.utils.ascii.AsciiBox;
import org.junit.jupiter.api.Test;
import org.skyscreamer.jsonassert.JSONAssert;
import org.skyscreamer.jsonassert.JSONCompareMode;

import java.io.ByteArrayInputStream;
import java.util.Collections;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class S7IoTest {

    @Test
    void TestS7MessageBytes() throws Exception {
        String wantBoxStringSerialized = ""
            + "╔═TPKTPacket═══════════════════════════════════════════════════════════════════════════════════════════════════════╗\n"
            + "║╔═protocolId╗╔═reserved╗╔═len═════╗                                                                               ║\n"
            + "║║  0x03 3   ║║ 0x00 0  ║║0x001d 29║                                                                               ║\n"
            + "║╚═══════════╝╚═════════╝╚═════════╝                                                                               ║\n"
            + "║╔═payload════════════════════════════════════════════════════════════════════════════════════════════════════════╗║\n"
            + "║║╔═COTPPacket═══════════════════════════════════════════════════════════════════════════════════════════════════╗║║\n"
            + "║║║╔═headerLength╗╔═tpduCode╗╔═COTPPacketData═════╗╔═parameters═════════════════════════════════════════════════╗║║║\n"
            + "║║║║   0x05 5    ║║0xf0 240 ║║╔═eot════╗╔═tpduRef╗║║╔═COTPParameter════════════════════════════════════════════╗║║║║\n"
            + "║║║╚═════════════╝╚═════════╝║║b0 false║║ 0xd 13 ║║║║╔═parameterType╗╔═parameterLength╗╔═COTPParameterTpduSize╗║║║║║\n"
            + "║║║                          ║╚════════╝╚════════╝║║║║   0xc0 192   ║║     0x01 1     ║║╔═tpduSize══════════╗ ║║║║║║\n"
            + "║║║                          ╚════════════════════╝║║╚══════════════╝╚════════════════╝║║╔═COTPTpduSize════╗║ ║║║║║║\n"
            + "║║║                                                ║║                                  ║║║0x0c 12 SIZE_4096║║ ║║║║║║\n"
            + "║║║                                                ║║                                  ║║╚═════════════════╝║ ║║║║║║\n"
            + "║║║                                                ║║                                  ║╚═══════════════════╝ ║║║║║║\n"
            + "║║║                                                ║║                                  ╚══════════════════════╝║║║║║\n"
            + "║║║                                                ║╚══════════════════════════════════════════════════════════╝║║║║\n"
            + "║║║                                                ╚════════════════════════════════════════════════════════════╝║║║\n"
            + "║║║╔═payload══════════════════════════════════════════════════════════════════════════════════╗                  ║║║\n"
            + "║║║║╔═S7Message══════════════════════════════════════════════════════════════════════════════╗║                  ║║║\n"
            + "║║║║║╔═protocolId╗╔═messageType╗╔═reserved╗╔═tpduReference╗╔═parameterLength╗╔═payloadLength╗║║                  ║║║\n"
            + "║║║║║║  0x32 50  ║║   0x03 3   ║║0x0000 0 ║║  0x000b 11   ║║    0x0002 2    ║║   0x0005 5   ║║║                  ║║║\n"
            + "║║║║║╚═══════════╝╚════════════╝╚═════════╝╚══════════════╝╚════════════════╝╚══════════════╝║║                  ║║║\n"
            + "║║║║║╔═S7MessageResponseData═══╗╔═parameter═════════════════════════════════════╗            ║║                  ║║║\n"
            + "║║║║║║╔═errorClass╗╔═errorCode╗║║╔═S7Parameter═════════════════════════════════╗║            ║║                  ║║║\n"
            + "║║║║║║║  0x00 0   ║║  0x00 0  ║║║║╔═parameterType╗╔═S7ParameterReadVarResponse╗║║            ║║                  ║║║\n"
            + "║║║║║║╚═══════════╝╚══════════╝║║║║    0x04 4    ║║        ╔═numItems╗        ║║║            ║║                  ║║║\n"
            + "║║║║║╚═════════════════════════╝║║╚══════════════╝║        ║ 0x01 1  ║        ║║║            ║║                  ║║║\n"
            + "║║║║║                           ║║                ║        ╚═════════╝        ║║║            ║║                  ║║║\n"
            + "║║║║║                           ║║                ╚═══════════════════════════╝║║            ║║                  ║║║\n"
            + "║║║║║                           ║╚═════════════════════════════════════════════╝║            ║║                  ║║║\n"
            + "║║║║║                           ╚═══════════════════════════════════════════════╝            ║║                  ║║║\n"
            + "║║║║║╔═payload══════════════════════════════════════════════════════════════╗                ║║                  ║║║\n"
            + "║║║║║║╔═S7Payload══════════════════════════════════════════════════════════╗║                ║║                  ║║║\n"
            + "║║║║║║║╔═S7PayloadReadVarResponse═════════════════════════════════════════╗║║                ║║                  ║║║\n"
            + "║║║║║║║║╔═items══════════════════════════════════════════════════════════╗║║║                ║║                  ║║║\n"
            + "║║║║║║║║║╔═S7VarPayloadDataItem═════════════════════════════════════════╗║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║╔═returnCode══════════════╗╔═transportSize══════╗╔═dataLength╗║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║║╔═DataTransportErrorCode╗║║╔═DataTransportSize╗║║ 0x0001 1  ║║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║║║      0xff 255 OK      ║║║║    0x03 3 BIT    ║║╚═══════════╝║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║║╚═══════════════════════╝║║╚══════════════════╝║             ║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║╚═════════════════════════╝╚════════════════════╝             ║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║╔═data═══════════════════════════════════════╗╔═padding╗      ║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║║0|01                            '.         '║║        ║      ║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║║╚════════════════════════════════════════════╝╚════════╝      ║║║║║                ║║                  ║║║\n"
            + "║║║║║║║║║╚══════════════════════════════════════════════════════════════╝║║║║                ║║                  ║║║\n"
            + "║║║║║║║║╚════════════════════════════════════════════════════════════════╝║║║                ║║                  ║║║\n"
            + "║║║║║║║╚══════════════════════════════════════════════════════════════════╝║║                ║║                  ║║║\n"
            + "║║║║║║╚════════════════════════════════════════════════════════════════════╝║                ║║                  ║║║\n"
            + "║║║║║╚══════════════════════════════════════════════════════════════════════╝                ║║                  ║║║\n"
            + "║║║║╚════════════════════════════════════════════════════════════════════════════════════════╝║                  ║║║\n"
            + "║║║╚══════════════════════════════════════════════════════════════════════════════════════════╝                  ║║║\n"
            + "║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║║\n"
            + "║╚════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝║\n"
            + "╚══════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝";
        String wantBoxStringSerializedCompact = ""
            + "╔═TPKTPacket═════════════════════════════════════════════════════════════════════════════════════╗\n"
            + "║╔═protocolId╗╔═reserved╗╔═len═════╗                                                             ║\n"
            + "║║  0x03 3   ║║ 0x00 0  ║║0x001d 29║                                                             ║\n"
            + "║╚═══════════╝╚═════════╝╚═════════╝                                                             ║\n"
            + "║╔═payload/COTPPacket═══════════════════════════════════════════════════════════════════════════╗║\n"
            + "║║╔═headerLength╗╔═tpduCode╗╔═COTPPacketData═════╗                                              ║║\n"
            + "║║║   0x05 5    ║║0xf0 240 ║║╔═eot════╗╔═tpduRef╗║                                              ║║\n"
            + "║║╚═════════════╝╚═════════╝║║b0 false║║ 0xd 13 ║║                                              ║║\n"
            + "║║                          ║╚════════╝╚════════╝║                                              ║║\n"
            + "║║                          ╚════════════════════╝                                              ║║\n"
            + "║║╔═parameters/COTPParameter════════════════════════════════════════════════════════════╗       ║║\n"
            + "║║║╔═parameterType╗╔═parameterLength╗╔═COTPParameterTpduSize/tpduSize/COTPTpduSize═════╗║       ║║\n"
            + "║║║║   0xc0 192   ║║     0x01 1     ║║                0x0c 12 SIZE_4096                ║║       ║║\n"
            + "║║║╚══════════════╝╚════════════════╝╚═════════════════════════════════════════════════╝║       ║║\n"
            + "║║╚═════════════════════════════════════════════════════════════════════════════════════╝       ║║\n"
            + "║║╔═payload/S7Message══════════════════════════════════════════════════════════════════════════╗║║\n"
            + "║║║╔═protocolId╗╔═messageType╗╔═reserved╗╔═tpduReference╗╔═parameterLength╗╔═payloadLength╗    ║║║\n"
            + "║║║║  0x32 50  ║║   0x03 3   ║║0x0000 0 ║║  0x000b 11   ║║    0x0002 2    ║║   0x0005 5   ║    ║║║\n"
            + "║║║╚═══════════╝╚════════════╝╚═════════╝╚══════════════╝╚════════════════╝╚══════════════╝    ║║║\n"
            + "║║║╔═S7MessageResponseData═══╗╔═parameter/S7Parameter═══════════════════════════════════╗      ║║║\n"
            + "║║║║╔═errorClass╗╔═errorCode╗║║╔═parameterType╗╔═S7ParameterReadVarResponse/numItems═══╗║      ║║║\n"
            + "║║║║║  0x00 0   ║║  0x00 0  ║║║║    0x04 4    ║║                0x01 1                 ║║      ║║║\n"
            + "║║║║╚═══════════╝╚══════════╝║║╚══════════════╝╚═══════════════════════════════════════╝║      ║║║\n"
            + "║║║╚═════════════════════════╝╚═════════════════════════════════════════════════════════╝      ║║║\n"
            + "║║║╔═payload/S7Payload/S7PayloadReadVarResponse/items/S7VarPayloadDataItem════════════════════╗║║║\n"
            + "║║║║╔═returnCode/DataTransportErrorCode════════════╗╔═transportSize/DataTransportSize════════╗║║║║\n"
            + "║║║║║                 0xff 255 OK                  ║║               0x03 3 BIT               ║║║║║\n"
            + "║║║║╚══════════════════════════════════════════════╝╚════════════════════════════════════════╝║║║║\n"
            + "║║║║╔═dataLength╗╔═data═══════════════════════════════════════╗                               ║║║║\n"
            + "║║║║║ 0x0001 1  ║║0|01                            '.         '║                               ║║║║\n"
            + "║║║║╚═══════════╝╚════════════════════════════════════════════╝                               ║║║║\n"
            + "║║║╚══════════════════════════════════════════════════════════════════════════════════════════╝║║║\n"
            + "║║╚════════════════════════════════════════════════════════════════════════════════════════════╝║║\n"
            + "║╚══════════════════════════════════════════════════════════════════════════════════════════════╝║\n"
            + "╚════════════════════════════════════════════════════════════════════════════════════════════════╝";
        String wantXml =
            "<TPKTPacket>\n" +
                "  <protocolId dataType=\"uint\" bitLength=\"8\">3</protocolId>\n" +
                "  <reserved dataType=\"uint\" bitLength=\"8\">0</reserved>\n" +
                "  <len dataType=\"uint\" bitLength=\"16\">29</len>\n" +
                "  <payload>\n" +
                "    <COTPPacket>\n" +
                "      <headerLength dataType=\"uint\" bitLength=\"8\">5</headerLength>\n" +
                "      <tpduCode dataType=\"uint\" bitLength=\"8\">240</tpduCode>\n" +
                "      <COTPPacketData>\n" +
                "        <eot dataType=\"bit\" bitLength=\"1\">false</eot>\n" +
                "        <tpduRef dataType=\"uint\" bitLength=\"7\">13</tpduRef>\n" +
                "      </COTPPacketData>\n" +
                "      <parameters isList=\"true\">\n" +
                "        <COTPParameter>\n" +
                "          <parameterType dataType=\"uint\" bitLength=\"8\">192</parameterType>\n" +
                "          <parameterLength dataType=\"uint\" bitLength=\"8\">1</parameterLength>\n" +
                "          <COTPParameterTpduSize>\n" +
                "            <tpduSize>\n" +
                "              <COTPTpduSize dataType=\"uint\" bitLength=\"8\" stringRepresentation=\"SIZE_4096\">12</COTPTpduSize>\n" +
                "            </tpduSize>\n" +
                "          </COTPParameterTpduSize>\n" +
                "        </COTPParameter>\n" +
                "      </parameters>\n" +
                "      <payload>\n" +
                "        <S7Message>\n" +
                "          <protocolId dataType=\"uint\" bitLength=\"8\">50</protocolId>\n" +
                "          <messageType dataType=\"uint\" bitLength=\"8\">3</messageType>\n" +
                "          <reserved dataType=\"uint\" bitLength=\"16\">0</reserved>\n" +
                "          <tpduReference dataType=\"uint\" bitLength=\"16\">11</tpduReference>\n" +
                "          <parameterLength dataType=\"uint\" bitLength=\"16\">2</parameterLength>\n" +
                "          <payloadLength dataType=\"uint\" bitLength=\"16\">5</payloadLength>\n" +
                "          <S7MessageResponseData>\n" +
                "            <errorClass dataType=\"uint\" bitLength=\"8\">0</errorClass>\n" +
                "            <errorCode dataType=\"uint\" bitLength=\"8\">0</errorCode>\n" +
                "          </S7MessageResponseData>\n" +
                "          <parameter>\n" +
                "            <S7Parameter>\n" +
                "              <parameterType dataType=\"uint\" bitLength=\"8\">4</parameterType>\n" +
                "              <S7ParameterReadVarResponse>\n" +
                "                <numItems dataType=\"uint\" bitLength=\"8\">1</numItems>\n" +
                "              </S7ParameterReadVarResponse>\n" +
                "            </S7Parameter>\n" +
                "          </parameter>\n" +
                "          <payload>\n" +
                "            <S7Payload>\n" +
                "              <S7PayloadReadVarResponse>\n" +
                "                <items isList=\"true\">\n" +
                "                  <S7VarPayloadDataItem>\n" +
                "                    <returnCode>\n" +
                "                      <DataTransportErrorCode dataType=\"uint\" bitLength=\"8\" stringRepresentation=\"OK\">255</DataTransportErrorCode>\n" +
                "                    </returnCode>\n" +
                "                    <transportSize>\n" +
                "                      <DataTransportSize dataType=\"uint\" bitLength=\"8\" stringRepresentation=\"BIT\">3</DataTransportSize>\n" +
                "                    </transportSize>\n" +
                "                    <dataLength dataType=\"uint\" bitLength=\"16\">1</dataLength>\n" +
                "                    <data dataType=\"byte\" bitLength=\"8\">0x01</data>\n" +
                "                    <padding isList=\"true\">\n" +
                "                    </padding>\n" +
                "                  </S7VarPayloadDataItem>\n" +
                "                </items>\n" +
                "              </S7PayloadReadVarResponse>\n" +
                "            </S7Payload>\n" +
                "          </payload>\n" +
                "        </S7Message>\n" +
                "      </payload>\n" +
                "    </COTPPacket>\n" +
                "  </payload>\n" +
                "</TPKTPacket>";
        String wantJson = "{\n" +
            "  \"TPKTPacket\": {\n" +
            "    \"len\": 29,\n" +
            "    \"len__plc4x_bitLength\": 16,\n" +
            "    \"len__plc4x_dataType\": \"uint\",\n" +
            "    \"payload\": {\n" +
            "      \"COTPPacket\": {\n" +
            "        \"COTPPacketData\": {\n" +
            "          \"eot\": false,\n" +
            "          \"eot__plc4x_bitLength\": 1,\n" +
            "          \"eot__plc4x_dataType\": \"bit\",\n" +
            "          \"tpduRef\": 13,\n" +
            "          \"tpduRef__plc4x_bitLength\": 7,\n" +
            "          \"tpduRef__plc4x_dataType\": \"uint\"\n" +
            "        },\n" +
            "        \"payload\": {\n" +
            "          \"S7Message\": {\n" +
            "            \"S7MessageResponseData\": {\n" +
            "              \"errorClass\": 0,\n" +
            "              \"errorClass__plc4x_bitLength\": 8,\n" +
            "              \"errorClass__plc4x_dataType\": \"uint\",\n" +
            "              \"errorCode\": 0,\n" +
            "              \"errorCode__plc4x_bitLength\": 8,\n" +
            "              \"errorCode__plc4x_dataType\": \"uint\"\n" +
            "            },\n" +
            "            \"parameter\": {\n" +
            "              \"S7Parameter\": {\n" +
            "                \"S7ParameterReadVarResponse\": {\n" +
            "                  \"numItems\": 1,\n" +
            "                  \"numItems__plc4x_bitLength\": 8,\n" +
            "                  \"numItems__plc4x_dataType\": \"uint\"\n" +
            "                },\n" +
            "                \"parameterType\": 4,\n" +
            "                \"parameterType__plc4x_bitLength\": 8,\n" +
            "                \"parameterType__plc4x_dataType\": \"uint\"\n" +
            "              }\n" +
            "            },\n" +
            "            \"payload\": {\n" +
            "              \"S7Payload\": {\n" +
            "                \"S7PayloadReadVarResponse\": {\n" +
            "                  \"items\": [\n" +
            "                    {\n" +
            "                      \"S7VarPayloadDataItem\": {\n" +
            "                        \"data\": \"0x01\",\n" +
            "                        \"dataLength\": 1,\n" +
            "                        \"dataLength__plc4x_bitLength\": 16,\n" +
            "                        \"dataLength__plc4x_dataType\": \"uint\",\n" +
            "                        \"data__plc4x_bitLength\": 8,\n" +
            "                        \"data__plc4x_dataType\": \"byte\",\n" +
            "                        \"returnCode\": {\n" +
            "                          \"DataTransportErrorCode\": 255,\n" +
            "                          \"DataTransportErrorCode__plc4x_bitLength\": 8,\n" +
            "                          \"DataTransportErrorCode__plc4x_dataType\": \"uint\",\n" +
            "                          \"DataTransportErrorCode__plc4x_stringRepresentation\": \"OK\"\n" +
            "                        },\n" +
            "                        \"transportSize\": {\n" +
            "                          \"DataTransportSize\": 3,\n" +
            "                          \"DataTransportSize__plc4x_bitLength\": 8,\n" +
            "                          \"DataTransportSize__plc4x_dataType\": \"uint\",\n" +
            "                          \"DataTransportSize__plc4x_stringRepresentation\": \"BIT\"\n" +
            "                        }\n" +
            "                      }\n" +
            "                    }\n" +
            "                  ]\n" +
            "                }\n" +
            "              }\n" +
            "            },\n" +
            "            \"messageType\": 3,\n" +
            "            \"messageType__plc4x_bitLength\": 8,\n" +
            "            \"messageType__plc4x_dataType\": \"uint\",\n" +
            "            \"parameterLength\": 2,\n" +
            "            \"parameterLength__plc4x_bitLength\": 16,\n" +
            "            \"parameterLength__plc4x_dataType\": \"uint\",\n" +
            "            \"payloadLength\": 5,\n" +
            "            \"payloadLength__plc4x_bitLength\": 16,\n" +
            "            \"payloadLength__plc4x_dataType\": \"uint\",\n" +
            "            \"protocolId\": 50,\n" +
            "            \"protocolId__plc4x_bitLength\": 8,\n" +
            "            \"protocolId__plc4x_dataType\": \"uint\",\n" +
            "            \"reserved\": 0,\n" +
            "            \"reserved__plc4x_bitLength\": 16,\n" +
            "            \"reserved__plc4x_dataType\": \"uint\",\n" +
            "            \"tpduReference\": 11,\n" +
            "            \"tpduReference__plc4x_bitLength\": 16,\n" +
            "            \"tpduReference__plc4x_dataType\": \"uint\"\n" +
            "          }\n" +
            "        },\n" +
            "        \"headerLength\": 5,\n" +
            "        \"headerLength__plc4x_bitLength\": 8,\n" +
            "        \"headerLength__plc4x_dataType\": \"uint\",\n" +
            "        \"parameters\": [\n" +
            "          {\n" +
            "            \"COTPParameter\": {\n" +
            "              \"COTPParameterTpduSize\": {\n" +
            "                \"tpduSize\": {\n" +
            "                  \"COTPTpduSize\": 12,\n" +
            "                  \"COTPTpduSize__plc4x_bitLength\": 8,\n" +
            "                  \"COTPTpduSize__plc4x_dataType\": \"uint\",\n" +
            "                  \"COTPTpduSize__plc4x_stringRepresentation\": \"SIZE_4096\"\n" +
            "                }\n" +
            "              },\n" +
            "              \"parameterLength\": 1,\n" +
            "              \"parameterLength__plc4x_bitLength\": 8,\n" +
            "              \"parameterLength__plc4x_dataType\": \"uint\",\n" +
            "              \"parameterType\": 192,\n" +
            "              \"parameterType__plc4x_bitLength\": 8,\n" +
            "              \"parameterType__plc4x_dataType\": \"uint\"\n" +
            "            }\n" +
            "          }\n" +
            "        ],\n" +
            "        \"tpduCode\": 240,\n" +
            "        \"tpduCode__plc4x_bitLength\": 8,\n" +
            "        \"tpduCode__plc4x_dataType\": \"uint\"\n" +
            "      }\n" +
            "    },\n" +
            "    \"protocolId\": 3,\n" +
            "    \"protocolId__plc4x_bitLength\": 8,\n" +
            "    \"protocolId__plc4x_dataType\": \"uint\",\n" +
            "    \"reserved\": 0,\n" +
            "    \"reserved__plc4x_bitLength\": 8,\n" +
            "    \"reserved__plc4x_dataType\": \"uint\"\n" +
            "  }\n" +
            "}";


        TPKTPacket tpktPacket = new TPKTPacket(
            new COTPPacketData(
                Collections.singletonList(new COTPParameterTpduSize(COTPTpduSize.SIZE_4096)),
                new S7MessageResponseData(
                    11,
                    new S7ParameterReadVarResponse((short) 1),
                    new S7PayloadReadVarResponse(
                        Collections.singletonList(
                            new S7VarPayloadDataItem(
                                DataTransportErrorCode.OK,
                                DataTransportSize.BIT,
                                new byte[]{0x1}
                            )
                        )
                    ),
                    (short) 0,
                    (short) 0
                ),
                false,
                (short) 13
            )
        );
        // To string
        {
            // TODO: implement me
            tpktPacket.toString();
        }

        // To box
        {
            WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased();
            tpktPacket.serialize(writeBufferBoxBased);
            AsciiBox gotBox = writeBufferBoxBased.getBox();
            assertEquals(wantBoxStringSerialized, gotBox.toString());
        }

        // To box compact
        {
            WriteBufferBoxBased writeBufferBoxBased = new WriteBufferBoxBased(true, true);
            tpktPacket.serialize(writeBufferBoxBased);
            AsciiBox gotBox = writeBufferBoxBased.getBox();
            assertEquals(wantBoxStringSerializedCompact, gotBox.toString());
        }

        // Xml
        {
            WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
            tpktPacket.serialize(writeBufferXmlBased);
            String gotXml = writeBufferXmlBased.getXmlString();
            assertEquals(wantXml, gotXml);
            ReadBufferXmlBased readBufferXmlBased = new ReadBufferXmlBased(new ByteArrayInputStream(gotXml.getBytes()));
            TPKTPacket reReadTpktPacket = TPKTPacket.staticParse(readBufferXmlBased);
            assertThat(reReadTpktPacket).usingRecursiveComparison().isEqualTo(tpktPacket);
        }
        // json
        {
            WriteBufferJsonBased writeBufferJsonBased = new WriteBufferJsonBased();
            tpktPacket.serialize(writeBufferJsonBased);
            String gotJson = writeBufferJsonBased.getJsonString();
            JSONAssert.assertEquals(wantJson, gotJson, JSONCompareMode.LENIENT);
            ReadBufferJsonBased readBufferXmlBased = new ReadBufferJsonBased(new ByteArrayInputStream(gotJson.getBytes()));
            TPKTPacket reReadTpktPacket = TPKTPacket.staticParse(readBufferXmlBased);
            assertThat(reReadTpktPacket).usingRecursiveComparison().isEqualTo(tpktPacket);
        }
    }
}
