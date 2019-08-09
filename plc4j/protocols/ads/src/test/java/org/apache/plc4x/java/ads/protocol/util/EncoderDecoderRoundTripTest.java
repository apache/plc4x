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
import org.apache.plc4x.java.base.messages.items.*;
import org.apache.plc4x.java.base.util.HexUtil;
import org.assertj.core.api.WithAssertions;
import org.junit.Ignore;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.util.Arrays;
import java.util.Set;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class EncoderDecoderRoundTripTest implements WithAssertions {

    private static Set<AdsDataType> untestedDataTypes = Stream.of(
        AdsDataType.STRING,
        AdsDataType.ARRAY,
        AdsDataType.POINTER,
        AdsDataType.ENUM,
        AdsDataType.STRUCT,
        AdsDataType.ALIAS,
        AdsDataType.SUB_RANGE_DATA_TYPE,
        AdsDataType.UNKNOWN
    ).collect(Collectors.toSet());

    @ParameterizedTest
    @MethodSource("generateData")
    @Ignore("Disabled for now as downgrading the Netty version broke the API")
    public void testMin(AdsDataType adsDataType, Number min, Number unused) throws Exception {
        byte[] bytes = LittleEndianEncoder.encodeData(adsDataType, min);
        assertThat(bytes).isNotEmpty();
        assertThat(bytes).hasSize(adsDataType.getTargetByteSize());
        BaseDefaultFieldItem baseDefaultFieldItem = LittleEndianDecoder.decodeData(adsDataType, bytes);

        if (baseDefaultFieldItem instanceof DefaultBooleanFieldItem) {
            assertThat(baseDefaultFieldItem.getBoolean(0)).isEqualTo(false);
        } else if (baseDefaultFieldItem instanceof DefaultLocalTimeFieldItem) {
            assertThat(baseDefaultFieldItem.getObject(0)).isEqualTo(LocalTime.ofNanoOfDay(TimeUnit.MILLISECONDS.toNanos((Long) min)));
        } else if (baseDefaultFieldItem instanceof DefaultLocalDateFieldItem) {
            assertThat(baseDefaultFieldItem.getObject(0)).isEqualTo(LocalDate.ofEpochDay(TimeUnit.SECONDS.toDays((Long) min)));
        } else if (baseDefaultFieldItem instanceof DefaultLocalDateTimeFieldItem) {
            assertThat(baseDefaultFieldItem.getObject(0)).isEqualTo(LocalDateTime.ofEpochSecond((Long) min, 0, ZoneOffset.UTC));
        } else {
            assertThat(baseDefaultFieldItem.getObject(0)).as("Min of %s using %s is equals %s\n%s", baseDefaultFieldItem, adsDataType, min, HexUtil.toHex(bytes)).isEqualTo(min);
        }
    }

    @ParameterizedTest
    @MethodSource("generateData")
    @Ignore("Disabled for now as downgrading the Netty version broke the API")
    public void testMax(AdsDataType adsDataType, Number unused, Number max) throws Exception {
        byte[] bytes = LittleEndianEncoder.encodeData(adsDataType, max);
        assertThat(bytes).isNotEmpty();
        assertThat(bytes).hasSize(adsDataType.getTargetByteSize());
        BaseDefaultFieldItem baseDefaultFieldItem = LittleEndianDecoder.decodeData(adsDataType, bytes);

        if (baseDefaultFieldItem instanceof DefaultBooleanFieldItem) {
            assertThat(baseDefaultFieldItem.getBoolean(0)).isEqualTo(true);
        } else if (baseDefaultFieldItem instanceof DefaultLocalTimeFieldItem) {
            assertThat(baseDefaultFieldItem.getObject(0)).isEqualTo(LocalTime.ofNanoOfDay(TimeUnit.MILLISECONDS.toNanos((Long) max)));
        } else if (baseDefaultFieldItem instanceof DefaultLocalDateFieldItem) {
            assertThat(baseDefaultFieldItem.getObject(0)).isEqualTo(LocalDate.ofEpochDay(TimeUnit.SECONDS.toDays((Long) max)));
        } else if (baseDefaultFieldItem instanceof DefaultLocalDateTimeFieldItem) {
            assertThat(baseDefaultFieldItem.getObject(0)).isEqualTo(LocalDateTime.ofEpochSecond((Long) max, 0, ZoneOffset.UTC));
        } else {
            assertThat(baseDefaultFieldItem.getObject(0)).as("Max of %s using %s is equals %s\n%s", baseDefaultFieldItem, adsDataType, max, HexUtil.toHex(bytes)).isEqualTo(max);
        }
    }

    public static Stream<Arguments> generateData() {
        return Arrays.stream(AdsDataType.values())
            .filter(adsDataType -> !untestedDataTypes.contains(adsDataType))
            .map(adsDataType -> Arguments.of(adsDataType, adsDataType.getLowerBound(), adsDataType.getUpperBound()));
    }
}
