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
package org.apache.plc4x.java.ads.protocol.util;

import org.apache.plc4x.java.ads.model.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.util.Arrays;
import java.util.Calendar;
import java.util.Date;
import java.util.stream.Stream;

import static org.apache.plc4x.java.base.util.Assert.assertByteEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

@Disabled("needs finishing")
public class LittleEndianEncoderTest {

    @ParameterizedTest
    @MethodSource("createAdsDataTypePears")
    public void decodeData(AdsDataType adsDataType, byte[] expectedData, Class<?> clazz, Object[] values) throws Exception {
        assertByteEquals(expectedData, LittleEndianEncoder.encodeData(adsDataType, values));
    }

    @Test
    public void negativeTest() {
        assertThrows(PlcUnsupportedDataTypeException.class, () -> LittleEndianEncoder.encodeData(AdsDataType.UNKNOWN, ""));
    }

    private static Stream<Arguments> createAdsDataTypePears() {
        Calendar calendar1 = Calendar.getInstance();
        //calendar1.set(2003, Calendar.DECEMBER, 23, 13, 3, 0);
        calendar1.setTime(new Date(1072180980436L));

        return Arrays.stream(AdsDataType.values())
            .filter(adsDataType -> adsDataType != AdsDataType.UNKNOWN)
            .flatMap(adsDataType -> Stream.of(
                Arguments.of(new byte[]{0x01, 0x00, 0x01, 0x00}, Boolean.class, true, false, true, false),
                Arguments.of(new byte[]{0x12, 0x03, 0x05, 0x7f}, Byte.class, (byte) 0x12, (byte) 0x03, (byte) 0x05, (byte) 0x7f),
                Arguments.of(new byte[]{0x1, 0x00}, Short.class, (short) 1),
                Arguments.of(new byte[]{0x0e, 0x00, 0x50, 0x00}, Short.class, (short) 14, (short) 80),
                Arguments.of(new byte[]{0x5a, 0x0a, 0x00, 0x00}, Integer.class, 2650),
                Arguments.of(new byte[]{0x5a, 0x0a, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00}, Integer.class, 2650, 80),
                Arguments.of(new byte[]{(byte) 0xc3, (byte) 0xf5, 0x48, 0x40}, Float.class, 3.14f),
                Arguments.of(new byte[]{(byte) 0xc3, (byte) 0xf5, 0x48, 0x40, 0x14, (byte) 0xae, 0x07, 0x40}, Float.class, 3.14f, 2.12f),
                Arguments.of(new byte[]{0x1F, (byte) 0x85, (byte) 0xEB, 0x51, (byte) 0xB8, 0x1E, 0x09, 0x40}, Double.class, 3.14),
                Arguments.of(new byte[]{0x1F, (byte) 0x85, (byte) 0xEB, 0x51, (byte) 0xB8, 0x1E, 0x09, 0x40, (byte) 0xF6, 0x28, 0x5C, (byte) 0x8F, (byte) 0xC2, (byte) 0xF5, 0x00, 0x40}, Double.class, 3.14, 2.12),
                Arguments.of(new byte[]{(byte) 0x40, (byte) 0x79, (byte) 0xFB, (byte) 0xB5, (byte) 0x4C, (byte) 0xC9, (byte) 0xC3, (byte) 0x01}, Calendar.class, calendar1),
                Arguments.of(new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x00}, String.class, "plc4x"),
                Arguments.of(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}, String.class, "HelloWorld!"),
                Arguments.of(new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x00, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}, String.class, "plc4x", "HelloWorld!")
            ).map(arguments -> Arguments.of(adsDataType, arguments.get()[0], arguments.get()[1], arguments.get()[2])));
    }
}