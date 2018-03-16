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
package org.apache.plc4x.java.ads.protocol;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.ads.api.commands.AdsReadRequest;
import org.apache.plc4x.java.ads.api.commands.AdsReadResponse;
import org.apache.plc4x.java.ads.api.commands.types.*;
import org.apache.plc4x.java.ads.api.generic.AmsPacket;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.api.generic.types.Invoke;
import org.apache.plc4x.java.ads.api.serial.AmsSerialAcknowledgeFrame;
import org.apache.plc4x.java.ads.api.serial.AmsSerialFrame;
import org.apache.plc4x.java.ads.api.serial.types.FragmentNumber;
import org.apache.plc4x.java.ads.api.serial.types.UserData;
import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.junit.Before;
import org.junit.Ignore;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ErrorCollector;

import java.util.ArrayList;
import java.util.Arrays;

import static org.apache.plc4x.java.ads.util.Assert.byteArrayEqualsTo;

/**
 * @see <a href="https://infosys.beckhoff.com/content/1033/tcadsamsserialspec/html/tcamssericalspec,0xsample.htm?id=60692407917020132">example</a>
 */
public class Payload2SerialProtocolExampleConversationTest {

    @Rule
    public ErrorCollector errorCollector = new ErrorCollector();

    private Payload2SerialProtocol SUT;

    @Before
    public void setUp() throws Exception {
        SUT = new Payload2SerialProtocol();
    }

    @Ignore("Either the crc calculation is wrong or the supplied crc are wrong")
    @Test
    public void exampleConversation() throws Exception {
        // 1. Terminal --> PLC : Request of 2 bytre data
        int[] exampleRequestInt = {
            /*Magic Cookie    */    0x01, 0xA5,
            /*Sender          */    0x00,
            /*Empfaenger      */    0x00,
            /*Fragmentnummer  */    0x06,
            /*Datenlaenge     */    0x2C,
            /*NetID Empfaenger*/    0xC0, 0xA8, 0x64, 0xAE, 0x01, 0x01,
            /*Port Nummer     */    0x21, 0x03,
            /*NetID Sender    */    0xC0, 0xA8, 0x64, 0x9C, 0x01, 0x01,
            /*Portnummer      */    0x01, 0x80,
            /*Kommando lesen  */    0x02, 0x00,
            /*Status          */    0x04, 0x00,
            /*Anzahl Datenbyte*/    0x0C, 0x00, 0x00, 0x00,
            /*Fehlercode      */    0x00, 0x00, 0x00, 0x00,
            /*InvokeID        */    0x07, 0x00, 0x00, 0x00,
            /*Index Gruppe    */    0x05, 0xF0, 0x00, 0x00,
            /*Index Offset    */    0x04, 0x00, 0x00, 0x9D,
            /*Anzahl Byte     */    0x02, 0x00, 0x00, 0x00,
            /*Checksumme      */    -24, -110,
        };
        byte[] exampleRequest = ArrayUtils.toPrimitive(Arrays
            .stream(exampleRequestInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));
        AmsPacket amsPacket = AdsReadRequest.of(
            AmsNetId.of(Arrays.copyOfRange(exampleRequest, 6, 12)),
            AmsPort.of(Arrays.copyOfRange(exampleRequest, 12, 14)),
            AmsNetId.of(Arrays.copyOfRange(exampleRequest, 14, 20)),
            AmsPort.of(Arrays.copyOfRange(exampleRequest, 20, 22)),
            Invoke.of(0x07),
            IndexGroup.of(0xF0_05),
            IndexOffset.of(0x9D_00_00_04L),
            Length.of(0x2)
        );
        AmsSerialFrame amsSerialFrame = AmsSerialFrame.of(FragmentNumber.of((byte) 0x06), UserData.of(amsPacket.getBytes()));
        System.out.println(amsSerialFrame.dump());
        System.err.println(ByteValue.of(exampleRequest).dump());
        errorCollector.checkThat(amsSerialFrame.getBytes(), byteArrayEqualsTo(exampleRequest));
        SUT.encode(null, amsPacket.getByteBuf(), new ArrayList<>());

        // PLC --> Terminal : Acknowledge:
        int[] exampleResponsAckInt = {
            /*Magic Cookie   */      0x01, 0x5A,
            /*Sender         */      0x00,
            /*Empfaenger     */      0x00,
            /*Fragmentnummer */      0x06,
            /*Datenlaenge    */      0x00,
            /*Checksumme     */      0x67, 0x5A
        };
        byte[] exampleAckResponse = ArrayUtils.toPrimitive(Arrays
            .stream(exampleResponsAckInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));

        AmsSerialAcknowledgeFrame amsSerialAcknowledgeFrame = AmsSerialAcknowledgeFrame.of(amsSerialFrame.getTransmitterAddress(), amsSerialFrame.getReceiverAddress(), amsSerialFrame.getFragmentNumber());
        System.out.println(amsSerialAcknowledgeFrame.dump());
        System.err.println(ByteValue.of(exampleAckResponse).dump());
        errorCollector.checkThat(amsSerialAcknowledgeFrame.getBytes(), byteArrayEqualsTo(exampleAckResponse));
        SUT.decode(null, amsSerialAcknowledgeFrame.getByteBuf(), new ArrayList<>());

        // PLC sends data:
        int[] exampleResponseInt = {
            /*Magic Cookie     */    0x01, 0xA5,
            /*Sender           */    0x00,
            /*Empfaenger       */    0x00,
            /*Fragmentnummer   */    0xEC,
            /*Anzahl Daten     */    0x2A,
            /*NetID Empfaenger */    0xC0, 0xA8, 0x64, 0x9C, 0x01, 0x01,
            /*Portnummer       */    0x01, 0x80,
            /*NetID Sender     */    0xC0, 0xA8, 0x64, 0xAE, 0x01, 0x01,
            /*Portnummer       */    0x21, 0x03,
            /*Response Lesen   */    0x02, 0x00,
            /*Status           */    0x05, 0x00,
            /*Anzahl Daten     */    0x0A, 0x00, 0x00, 0x00,
            /*Fehlercode       */    0x00, 0x00, 0x00, 0x00,
            /*InvokeID         */    0x07, 0x00, 0x00, 0x00,
            /*Ergebnis         */    0x00, 0x00, 0x00, 0x00,
            /*Anzahl Daten     */    0x02, 0x00, 0x00, 0x00,
            /*Daten            */    0xAF, 0x27,
            /*Checksumme       */    0x04, 0xA9,
        };
        byte[] exampleResponse = ArrayUtils.toPrimitive(Arrays
            .stream(exampleResponseInt)
            .mapToObj(value -> (byte) value)
            .toArray(Byte[]::new));
        AmsPacket amsResponsePacket = AdsReadResponse.of(
            AmsNetId.of(Arrays.copyOfRange(exampleResponse, 6, 12)),
            AmsPort.of(Arrays.copyOfRange(exampleResponse, 12, 14)),
            AmsNetId.of(Arrays.copyOfRange(exampleResponse, 14, 20)),
            AmsPort.of(Arrays.copyOfRange(exampleResponse, 20, 22)),
            Invoke.of(0x07),
            Result.of(0x00),
            Data.of((byte) 0xAF, (byte) 0x27)
        );
        AmsSerialFrame amsResponseSerialFrame = AmsSerialFrame.of(FragmentNumber.of((byte) 0xEC), UserData.of(amsResponsePacket.getBytes()));
        System.out.println(amsResponseSerialFrame.dump());
        System.err.println(ByteValue.of(exampleResponse).dump());
        errorCollector.checkThat(amsResponseSerialFrame.getBytes(), byteArrayEqualsTo(exampleResponse));
    }
}
