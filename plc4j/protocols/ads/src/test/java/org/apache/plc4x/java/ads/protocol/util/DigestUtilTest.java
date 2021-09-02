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
package org.apache.plc4x.java.ads.protocol.util;

import com.github.snksoft.crc.CRC;
import io.netty.buffer.Unpooled;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.StringUtils;
import org.junit.Ignore;
import org.junit.Test;

import java.nio.ByteBuffer;
import java.util.Arrays;
import java.util.stream.IntStream;
import java.util.stream.Stream;

import static org.junit.Assert.assertEquals;

public class DigestUtilTest {

    @Test
    public void sameCalculation() {
        assertEquals(
            DigestUtil.calculateCrc16(new byte[]{47, 99}),
            DigestUtil.calculateCrc16(() -> Unpooled.wrappedBuffer(new byte[]{47}), () -> Unpooled.wrappedBuffer(new byte[]{99})));
    }

    @Test
    public void simpleUnwrap() {
        assertEquals(
            DigestUtil.calculateCrc16(new byte[]{47}),
            DigestUtil.calculateCrc16(() -> Unpooled.wrappedBuffer(new byte[]{47})));
    }

    @Test
    public void displayValue() {
        Object[] crcs = IntStream.range(0, 256)
            .map(value -> DigestUtil.calculateCrc16(new byte[]{(byte) value}))
            .mapToObj(Integer::toHexString)
            .map(s -> StringUtils.leftPad(s, 4, '0'))
            .toArray(String[]::new);
        for (int i = 0; i < 32; i++) {
            String col1 = "" + (i) + '\t' + crcs[i];
            String col2 = "" + (i + 32) + '\t' + crcs[i + 32];
            String col3 = "" + (i + 64) + '\t' + crcs[i + 64];
            String col4 = "" + (i + 96) + '\t' + crcs[i + 96];
            String col5 = "" + (i + 128) + '\t' + crcs[i + 128];
            String col6 = "" + (i + 160) + '\t' + crcs[i + 160];
            String col7 = "" + (i + 192) + '\t' + crcs[i + 192];
            String col8 = "" + (i + 224) + '\t' + crcs[i + 224];
            System.out.println(col1 + "\t\t" + col2 + "\t\t" + col3 + "\t\t" + col4 + "\t\t" + col5 + "\t\t" + col6 + "\t\t" + col7 + "\t\t" + col8);
        }
    }

    @Ignore("This test is used to find proper crc parameters and can be run manually")
    @Test
    public void findRightAdsCrc() throws Exception {
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
        int crcField1Index = exampleResponse.length - 2;
        assertEquals((byte) 0x04, exampleResponse[crcField1Index]);
        int crcField2Index = exampleResponse.length - 1;
        assertEquals((byte) 0xA9, exampleResponse[crcField2Index]);

        IntStream.range(0, 0xFFFF).forEach(polynomial -> {
                //int polynomial = 0x8005;
                IntStream.of(0x0000, 0xFFFF).forEach(init -> {
                    //int init = 0x0000;
                    Stream.of(true, false).forEach(reflectIn -> {
                        //boolean reflectIn = true;
                        Stream.of(true, false).forEach(reflectOut -> {
                            //boolean reflectOut = true;
                            // TODO: check if we can check for finalxor at all
                            int finalXor = 0x0;
                            CRC.Parameters params = new CRC.Parameters(16, polynomial, init, reflectIn, reflectOut, finalXor);
                            String paramString = ""
                                + "polynomial=0x" + StringUtils.leftPad(Integer.toHexString(polynomial), 4, '0')
                                + ", init=0x" + StringUtils.leftPad(Integer.toHexString(init), 4, '0')
                                + ", reflectIn=" + reflectIn
                                + ", reflectOut=" + reflectOut
                                + ", finalXor=" + StringUtils.leftPad(Integer.toHexString(finalXor), 4, '0') + "";
                            CRC crcCalculator = new CRC(params);
                            long currentCrc = crcCalculator.update(crcCalculator.init(), Arrays.copyOfRange(exampleResponse, 0, crcField1Index));
                            short crc16 = crcCalculator.finalCRC16(currentCrc);
                            ByteBuffer buffer = ByteBuffer.allocate(2);
                            buffer.putShort(crc16);
                            byte[] bytes = buffer.array();
                            byte msb = bytes[0];
                            byte lsb = bytes[1];
                            if (lsb == exampleResponse[crcField1Index] && msb == exampleResponse[crcField2Index]) {
                                System.out.println("We found a LE: " + paramString);
                            }
                            if (lsb == exampleResponse[crcField2Index] && msb == exampleResponse[crcField1Index]) {
                                System.out.println("We found a BE " + paramString);
                            }
                            if (lsb == exampleResponse[crcField1Index] || msb == exampleResponse[crcField2Index]) {
                                //System.out.println("We found a LE one of them");
                            } else if (lsb == exampleResponse[crcField2Index] || msb == exampleResponse[crcField1Index]) {
                                //System.out.println("We found a BE one of them");
                            }
                        });
                    });
                });
            }
        );
    }
}