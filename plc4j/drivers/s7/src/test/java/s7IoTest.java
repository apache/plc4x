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

import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.io.TPKTPacketIO;
import org.apache.plc4x.java.s7.readwrite.types.COTPTpduSize;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.types.DataTransportSize;
import org.apache.plc4x.java.spi.generation.WriteBufferXmlBased;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class s7IoTest {

    @Disabled("Till the output is aligned to the one of golang")
    @Test
    void TestS7MessageBytes() throws Exception {
        String wantXml =
            "<TPKTPacket>\n" +
            "    <protocolId bitLength=\"8\" dataType=\"uint8\">3</protocolId>\n" +
            "    <reserved bitLength=\"8\" dataType=\"uint8\">0</reserved>\n" +
            "    <len bitLength=\"16\" dataType=\"uint16\">29</len>\n" +
            "    <COTPPacket>\n" +
            "        <headerLength bitLength=\"8\" dataType=\"uint18\">5</headerLength>\n" +
            "        <tpduCode bitLength=\"8\" dataType=\"uint8\">240</tpduCode>\n" +
            "        <COTPPacketData>\n" +
            "            <eot bitLength=\"1\" dataType=\"bit\">false</eot>\n" +
            "            <tpduRef bitLength=\"7\" dataType=\"uint8\">13</tpduRef>\n" +
            "        </COTPPacketData>\n" +
            "        <parameters>\n" +
            "            <COTPParameter>\n" +
            "                <parameterType bitLength=\"8\" dataType=\"uint8\">192</parameterType>\n" +
            "                <parameterLength bitLength=\"8\" dataType=\"uint8\">1</parameterLength>\n" +
            "                <COTPParameterTpduSize>\n" +
            "                    <COTPTpduSize bitLength=\"8\" dataType=\"uint8\" stringRepresentation=\"SIZE_4096\">12</COTPTpduSize>\n" +
            "                </COTPParameterTpduSize>\n" +
            "            </COTPParameter>\n" +
            "        </parameters>\n" +
            "        <S7Message>\n" +
            "            <protocolId bitLength=\"8\" dataType=\"uint8\">50</protocolId>\n" +
            "            <messageType bitLength=\"8\" dataType=\"uint8\">3</messageType>\n" +
            "            <reserved bitLength=\"16\" dataType=\"uint16\">0</reserved>\n" +
            "            <tpduReference bitLength=\"16\" dataType=\"uint16\">11</tpduReference>\n" +
            "            <parameterLength bitLength=\"16\" dataType=\"uint16\">2</parameterLength>\n" +
            "            <payloadLength bitLength=\"16\" dataType=\"uint16\">5</payloadLength>\n" +
            "            <S7MessageResponseData>\n" +
            "                <errorClass bitLength=\"8\" dataType=\"uint8\">0</errorClass>\n" +
            "                <errorCode bitLength=\"8\" dataType=\"uint8\">0</errorCode>\n" +
            "            </S7MessageResponseData>\n" +
            "            <S7Parameter>\n" +
            "                <parameterType bitLength=\"8\" dataType=\"uint8\">4</parameterType>\n" +
            "                <S7ParameterReadVarResponse>\n" +
            "                    <numItems bitLength=\"8\" dataType=\"uint8\">1</numItems>\n" +
            "                </S7ParameterReadVarResponse>\n" +
            "            </S7Parameter>\n" +
            "            <S7Payload>\n" +
            "                <S7PayloadReadVarResponse>\n" +
            "                    <items>\n" +
            "                        <S7VarPayloadDataItem>\n" +
            "                            <DataTransportErrorCode bitLength=\"8\" dataType=\"uint8\" stringRepresentation=\"OK\">255</DataTransportErrorCode>\n" +
            "                            <DataTransportSize bitLength=\"8\" dataType=\"uint8\" stringRepresentation=\"BIT\">3</DataTransportSize>\n" +
            "                            <dataLength bitLength=\"16\" dataType=\"uint16\">1</dataLength>\n" +
            "                            <data>\n" +
            "                                <value bitLength=\"8\" dataType=\"int8\">1</value>\n" +
            "                            </data>\n" +
            "                            <padding></padding>\n" +
            "                        </S7VarPayloadDataItem>\n" +
            "                    </items>\n" +
            "                </S7PayloadReadVarResponse>\n" +
            "            </S7Payload>\n" +
            "        </S7Message>\n" +
            "    </COTPPacket>\n" +
            "</TPKTPacket>";


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
        WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
        TPKTPacketIO.staticSerialize(writeBufferXmlBased, tpktPacket);
        String gotXml = writeBufferXmlBased.getXmlString();
        assertEquals(wantXml,gotXml);
    }
}
