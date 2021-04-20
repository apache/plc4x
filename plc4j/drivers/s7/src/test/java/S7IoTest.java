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
import org.apache.plc4x.java.spi.generation.ReadBufferXmlBased;
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
                "      <parameters>\n" +
                "        <COTPParameter>\n" +
                "          <parameterType dataType=\"uint\" bitLength=\"8\">192</parameterType>\n" +
                "          <parameterLength dataType=\"uint\" bitLength=\"8\">1</parameterLength>\n" +
                "          <COTPParameterTpduSize>\n" +
                "            <tpduSize>\n" +
                "              <COTPTpduSize dataType=\"int\" bitLength=\"8\">12</COTPTpduSize>\n" +
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
                "            <items>\n" +
                "              <S7VarPayloadDataItem>\n" +
                "                <returnCode>\n" +
                "                  <DataTransportErrorCode dataType=\"uint\" bitLength=\"8\">255</DataTransportErrorCode>\n" +
                "                </returnCode>\n" +
                "                <transportSize>\n" +
                "                  <DataTransportSize dataType=\"uint\" bitLength=\"8\">3</DataTransportSize>\n" +
                "                </transportSize>\n" +
                "                <dataLength dataType=\"uint\" bitLength=\"16\">1</dataLength>\n" +
                "                <data>\n" +
                "                  <value dataType=\"int\" bitLength=\"8\">1</value>\n" +
                "                </data>\n" +
                "                <padding>\n" +
                "                </padding>\n" +
                "              </S7VarPayloadDataItem>\n" +
                "            </items>\n" +
                "          </S7PayloadReadVarResponse>\n" +
                "        </S7Payload>\n" +
                "      </S7Message>\n" +
                "    </COTPPacket>\n" +
                "  </payload>\n" +
                "</TPKTPacket>\n";


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
        assertEquals(wantXml, gotXml);
        System.out.println(gotXml);
        ReadBufferXmlBased readBufferXmlBased = new ReadBufferXmlBased(new ByteArrayInputStream(gotXml.getBytes()));
        TPKTPacket reReadTpktPacket = TPKTPacketIO.staticParse(readBufferXmlBased);
        assertThat(reReadTpktPacket).usingRecursiveComparison().isEqualTo(tpktPacket);
    }
}
