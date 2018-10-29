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
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.util.Arrays;
import java.util.stream.Stream;

@Disabled("check if the encoder/decoder can be used like this.")
public class EncoderDecoderRoundTripTest implements WithAssertions {

    @ParameterizedTest
    @MethodSource("generateData")
    public void testMin(AdsDataType adsDataType, Double min, Double max) throws Exception {
        byte[] bytes = LittleEndianEncoder.encodeData(adsDataType, min);
        BaseDefaultFieldItem baseDefaultFieldItem = LittleEndianDecoder.decodeData(adsDataType, bytes);

        assertThat(baseDefaultFieldItem.getDouble(0)).isEqualTo(min);
    }

    @ParameterizedTest
    @MethodSource("generateData")
    public void testMax(AdsDataType adsDataType, Double min, Double max) throws Exception {
        byte[] bytes = LittleEndianEncoder.encodeData(adsDataType, max);
        BaseDefaultFieldItem baseDefaultFieldItem = LittleEndianDecoder.decodeData(adsDataType, bytes);

        assertThat(baseDefaultFieldItem.getDouble(0)).isEqualTo(max);
    }

    public static Stream<Arguments> generateData() {
        return Arrays.stream(AdsDataType.values())
            .map(adsDataType -> Arguments.of(adsDataType, adsDataType.getLowerBound(), adsDataType.getUpperBound()));
    }
}
