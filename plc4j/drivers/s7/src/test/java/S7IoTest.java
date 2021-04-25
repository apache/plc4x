/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at
 
   http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.io.TPKTPacketIO;
import org.apache.plc4x.java.s7.readwrite.types.COTPTpduSize;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportSize;
import org.apache.plc4x.java.spi.generation.ReadBufferJsonBased;
import org.apache.plc4x.java.spi.generation.ReadBufferXmlBased;
import org.apache.plc4x.java.spi.generation.WriteBufferJsonBased;
import org.apache.plc4x.java.spi.generation.WriteBufferXmlBased;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayInputStream;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class S7IoTest {

    @Test
    void TestS7MessageBytes() throws Exception {
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
                "              <COTPTpduSize dataType=\"int\" bitLength=\"8\" stringRepresentation=\"SIZE_4096\">12</COTPTpduSize>\n" +
                "            </tpduSize>\n" +
                "          </COTPParameterTpduSize>\n" +
                "        </COTPParameter>\n" +
                "      </parameters>\n" +
                "      <S7Message>\n" +
                "        <protocolId dataType=\"uint\" bitLength=\"8\">50</protocolId>\n" +
                "        <messageType dataType=\"uint\" bitLength=\"8\">3</messageType>\n" +
                "        <reserved dataType=\"uint\" bitLength=\"16\">0</reserved>\n" +
                "        <tpduReference dataType=\"uint\" bitLength=\"16\">11</tpduReference>\n" +
                "        <parameterLength dataType=\"uint\" bitLength=\"16\">2</parameterLength>\n" +
                "        <payloadLength dataType=\"uint\" bitLength=\"16\">5</payloadLength>\n" +
                "        <S7MessageResponseData>\n" +
                "          <errorClass dataType=\"uint\" bitLength=\"8\">0</errorClass>\n" +
                "          <errorCode dataType=\"uint\" bitLength=\"8\">0</errorCode>\n" +
                "        </S7MessageResponseData>\n" +
                "        <S7Parameter>\n" +
                "          <parameterType dataType=\"uint\" bitLength=\"8\">4</parameterType>\n" +
                "          <S7ParameterReadVarResponse>\n" +
                "            <numItems dataType=\"uint\" bitLength=\"8\">1</numItems>\n" +
                "          </S7ParameterReadVarResponse>\n" +
                "        </S7Parameter>\n" +
                "        <S7Payload>\n" +
                "          <S7PayloadReadVarResponse>\n" +
                "            <items isList=\"true\">\n" +
                "              <S7VarPayloadDataItem>\n" +
                "                <returnCode>\n" +
                "                  <DataTransportErrorCode dataType=\"uint\" bitLength=\"8\" stringRepresentation=\"OK\">255</DataTransportErrorCode>\n" +
                "                </returnCode>\n" +
                "                <transportSize>\n" +
                "                  <DataTransportSize dataType=\"uint\" bitLength=\"8\" stringRepresentation=\"BIT\">3</DataTransportSize>\n" +
                "                </transportSize>\n" +
                "                <dataLength dataType=\"uint\" bitLength=\"16\">1</dataLength>\n" +
                "                <data isList=\"true\">\n" +
                "                  <value dataType=\"int\" bitLength=\"8\">1</value>\n" +
                "                </data>\n" +
                "                <padding isList=\"true\">\n" +
                "                </padding>\n" +
                "              </S7VarPayloadDataItem>\n" +
                "            </items>\n" +
                "          </S7PayloadReadVarResponse>\n" +
                "        </S7Payload>\n" +
                "      </S7Message>\n" +
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
            "        \"S7Message\": {\n" +
            "          \"S7MessageResponseData\": {\n" +
            "            \"errorClass\": 0,\n" +
            "            \"errorClass__plc4x_bitLength\": 8,\n" +
            "            \"errorClass__plc4x_dataType\": \"uint\",\n" +
            "            \"errorCode\": 0,\n" +
            "            \"errorCode__plc4x_bitLength\": 8,\n" +
            "            \"errorCode__plc4x_dataType\": \"uint\"\n" +
            "          },\n" +
            "          \"S7Parameter\": {\n" +
            "            \"S7ParameterReadVarResponse\": {\n" +
            "              \"numItems\": 1,\n" +
            "              \"numItems__plc4x_bitLength\": 8,\n" +
            "              \"numItems__plc4x_dataType\": \"uint\"\n" +
            "            },\n" +
            "            \"parameterType\": 4,\n" +
            "            \"parameterType__plc4x_bitLength\": 8,\n" +
            "            \"parameterType__plc4x_dataType\": \"uint\"\n" +
            "          },\n" +
            "          \"S7Payload\": {\n" +
            "            \"S7PayloadReadVarResponse\": {\n" +
            "              \"items\": [\n" +
            "                {\n" +
            "                  \"S7VarPayloadDataItem\": {\n" +
            "                    \"data\": [\n" +
            "                      {\n" +
            "                        \"value\": 1,\n" +
            "                        \"value__plc4x_bitLength\": 8,\n" +
            "                        \"value__plc4x_dataType\": \"int\"\n" +
            "                      }\n" +
            "                    ],\n" +
            "                    \"dataLength\": 1,\n" +
            "                    \"dataLength__plc4x_bitLength\": 16,\n" +
            "                    \"dataLength__plc4x_dataType\": \"uint\",\n" +
            "                    \"padding\": [],\n" +
            "                    \"returnCode\": {\n" +
            "                      \"DataTransportErrorCode\": 255,\n" +
            "                      \"DataTransportErrorCode__plc4x_bitLength\": 8,\n" +
            "                      \"DataTransportErrorCode__plc4x_dataType\": \"uint\",\n" +
            "                      \"DataTransportErrorCode__plc4x_stringRepresentation\": \"OK\"\n" +
            "                    },\n" +
            "                    \"transportSize\": {\n" +
            "                      \"DataTransportSize\": 3,\n" +
            "                      \"DataTransportSize__plc4x_bitLength\": 8,\n" +
            "                      \"DataTransportSize__plc4x_dataType\": \"uint\",\n" +
            "                      \"DataTransportSize__plc4x_stringRepresentation\": \"BIT\"\n" +
            "                    }\n" +
            "                  }\n" +
            "                }\n" +
            "              ]\n" +
            "            }\n" +
            "          },\n" +
            "          \"messageType\": 3,\n" +
            "          \"messageType__plc4x_bitLength\": 8,\n" +
            "          \"messageType__plc4x_dataType\": \"uint\",\n" +
            "          \"parameterLength\": 2,\n" +
            "          \"parameterLength__plc4x_bitLength\": 16,\n" +
            "          \"parameterLength__plc4x_dataType\": \"uint\",\n" +
            "          \"payloadLength\": 5,\n" +
            "          \"payloadLength__plc4x_bitLength\": 16,\n" +
            "          \"payloadLength__plc4x_dataType\": \"uint\",\n" +
            "          \"protocolId\": 50,\n" +
            "          \"protocolId__plc4x_bitLength\": 8,\n" +
            "          \"protocolId__plc4x_dataType\": \"uint\",\n" +
            "          \"reserved\": 0,\n" +
            "          \"reserved__plc4x_bitLength\": 16,\n" +
            "          \"reserved__plc4x_dataType\": \"uint\",\n" +
            "          \"tpduReference\": 11,\n" +
            "          \"tpduReference__plc4x_bitLength\": 16,\n" +
            "          \"tpduReference__plc4x_dataType\": \"uint\"\n" +
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
            "                  \"COTPTpduSize__plc4x_dataType\": \"int\",\n" +
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
                new COTPParameter[]{new COTPParameterTpduSize(COTPTpduSize.SIZE_4096)},
                new S7MessageResponseData(
                    11,
                    new S7ParameterReadVarResponse((short) 1),
                    new S7PayloadReadVarResponse(
                        new S7VarPayloadDataItem[]{
                            new S7VarPayloadDataItem(
                                DataTransportErrorCode.OK,
                                DataTransportSize.BIT,
                                new byte[]{0x1}
                            )
                        }
                    ),
                    (short) 0,
                    (short) 0
                ),
                false,
                (short) 13
            )
        );
        // Xml
        {
            WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
            TPKTPacketIO.staticSerialize(writeBufferXmlBased, tpktPacket);
            String gotXml = writeBufferXmlBased.getXmlString();
            assertEquals(wantXml, gotXml);
            ReadBufferXmlBased readBufferXmlBased = new ReadBufferXmlBased(new ByteArrayInputStream(gotXml.getBytes()));
            TPKTPacket reReadTpktPacket = TPKTPacketIO.staticParse(readBufferXmlBased);
            assertThat(reReadTpktPacket).usingRecursiveComparison().isEqualTo(tpktPacket);
        }
        // json
        {
            WriteBufferJsonBased writeBufferJsonBased = new WriteBufferJsonBased();
            TPKTPacketIO.staticSerialize(writeBufferJsonBased, tpktPacket);
            String gotJson = writeBufferJsonBased.getJsonString();
            ObjectMapper objectMapper = new ObjectMapper();
            assertEquals(objectMapper.readTree(wantJson), objectMapper.readTree(gotJson));
            ReadBufferJsonBased readBufferXmlBased = new ReadBufferJsonBased(new ByteArrayInputStream(gotJson.getBytes()));
            TPKTPacket reReadTpktPacket = TPKTPacketIO.staticParse(readBufferXmlBased);
            assertThat(reReadTpktPacket).usingRecursiveComparison().isEqualTo(tpktPacket);
        }
    }
}
