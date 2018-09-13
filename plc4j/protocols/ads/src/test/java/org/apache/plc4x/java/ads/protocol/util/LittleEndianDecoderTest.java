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
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.util.Arrays;
import java.util.Calendar;
import java.util.Collection;
import java.util.Date;
import java.util.stream.Stream;

import static java.util.Arrays.asList;
import static java.util.Collections.singletonList;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

@Disabled("needs finishing")
public class LittleEndianDecoderTest {

    @ParameterizedTest
    @MethodSource("createAdsDataTypePears")
    public void decodeData(AdsDataType adsDataType, Collection expectedTypes, Class<?> clazz, byte[] adsData) {
        assertEquals(expectedTypes, Arrays.asList(LittleEndianDecoder.decodeData(adsDataType, adsData).getValues()));
    }

    @Test
    public void negativeTest() {
        assertThrows(PlcProtocolException.class, () -> LittleEndianDecoder.decodeData(AdsDataType.STRING, new byte[]{0x01}));
        assertThrows(PlcUnsupportedDataTypeException.class, () -> LittleEndianDecoder.decodeData(AdsDataType.UNKNOWN, new byte[10]));
    }

    private static Stream<Arguments> createAdsDataTypePears() {
        Calendar calendar1 = Calendar.getInstance();
        calendar1.setTime(new Date(-11644473600000L));
        Calendar calendar0x0001 = Calendar.getInstance();
        calendar0x0001.setTime(new Date(-4438714196208L));

        return Arrays.stream(AdsDataType.values())
            .filter(adsDataType -> adsDataType != AdsDataType.UNKNOWN)
            .flatMap(adsDataType -> Stream.of(
                Arguments.of(asList(true, false), Boolean.class, new byte[]{0x1, 0x0}),
                Arguments.of(asList((byte) 0x1, (byte) 0x0), Byte.class, new byte[]{0x1, 0x0}),
                Arguments.of(singletonList((short) 1), Short.class, new byte[]{0x1}),
                Arguments.of(singletonList((short) 256), Short.class, new byte[]{0x0, 0x1}),
                Arguments.of(asList((short) 256, (short) 256), Short.class, new byte[]{0x0, 0x1, 0x0, 0x1}),
                Arguments.of(singletonList(1), Integer.class, new byte[]{0x1}),
                Arguments.of(singletonList(16777216), Integer.class, new byte[]{0x0, 0x0, 0x0, 0x1}),
                Arguments.of(asList(16777216, 16777216), Integer.class, new byte[]{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1}),
                Arguments.of(singletonList(1.4E-45f), Float.class, new byte[]{0x1}),
                Arguments.of(singletonList(2.3509887E-38f), Float.class, new byte[]{0x0, 0x0, 0x0, 0x1}),
                Arguments.of(asList(2.3509887E-38f, 2.3509887E-38f), Float.class, new byte[]{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1}),
                Arguments.of(singletonList(4.9E-324), Double.class, new byte[]{0x1}),
                Arguments.of(singletonList(7.2911220195563975E-304), Double.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}),
                Arguments.of(asList(7.2911220195563975E-304, 7.2911220195563975E-304), Double.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}),
                Arguments.of(singletonList(calendar1), Calendar.class, new byte[]{0x1}),
                Arguments.of(singletonList(calendar0x0001), Calendar.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}),
                Arguments.of(asList(calendar0x0001, calendar0x0001), Calendar.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}),
                Arguments.of(singletonList("plc4x"), String.class, new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x0}),
                Arguments.of(singletonList("plc4xplc4x"), String.class, new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x70, 0x6c, 0x63, 0x34, 0x78, 0x0}),
                Arguments.of(asList("plc4x", "plc4x"), String.class, new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x0, 0x70, 0x6c, 0x63, 0x34, 0x78, 0x0})
            ).map(arguments -> Arguments.of(adsDataType, arguments.get()[0], arguments.get()[1], arguments.get()[2])));
    }
}