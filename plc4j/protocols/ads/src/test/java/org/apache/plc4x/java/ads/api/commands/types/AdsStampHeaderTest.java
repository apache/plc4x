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
package org.apache.plc4x.java.ads.api.commands.types;

import org.junit.Test;

import java.util.Arrays;

import static org.apache.plc4x.java.mock.util.Assert.assertByteEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class AdsStampHeaderTest {

    @Test
    public void builder() throws Exception {
        {
            AdsStampHeader adsStampHeader = AdsStampHeader.of(TimeStamp.of(13), Arrays.asList(AdsNotificationSample.of(NotificationHandle.of(13), Data.of("Hello World!")), AdsNotificationSample.of(NotificationHandle.of(14), Data.of("Hello World!"))));
            assertNotNull(adsStampHeader);
            byte[] expected = {
                (byte) 0xD0, 0x7B, 0x40, (byte) 0xD5,
                (byte) 0xDE, (byte) 0xB1, (byte) 0x9D, 0x01,
                0x02, 0x00, 0x00, 0x00,
                0x0D, 0x00, 0x00, 0x00,
                0x0D, 0x00, 0x00, 0x00,
                0x48, 0x65, 0x6C, 0x6C,
                0x6F, 0x20, 0x57, 0x6F,
                0x72, 0x6C, 0x64, 0x21,
                0x00, 0x0E, 0x00, 0x00,
                0x00, 0x0D, 0x00, 0x00,
                0x00, 0x48, 0x65, 0x6C,
                0x6C, 0x6F, 0x20, 0x57,
                0x6F, 0x72, 0x6C, 0x64,
                0x21, 0x00
            };
            assertByteEquals(expected, adsStampHeader.getBytes());
            assertEquals(Samples.of(2), adsStampHeader.getSamples());
            assertEquals(2, adsStampHeader.getAdsNotificationSamples().size());
            assertEquals(Data.of("Hello World!"), adsStampHeader.getAdsNotificationSamples().get(0).getData());
        }
        {
            AdsStampHeader adsStampHeader = AdsStampHeader.of(TimeStamp.of(13), Samples.of(2), Arrays.asList(AdsNotificationSample.of(NotificationHandle.of(13), Data.of("Hello World!")), AdsNotificationSample.of(NotificationHandle.of(14), Data.of("Hello World!"))));
            assertNotNull(adsStampHeader);
            byte[] expected = {
                (byte) 0xD0, 0x7B, 0x40, (byte) 0xD5,
                (byte) 0xDE, (byte) 0xB1, (byte) 0x9D, 0x01,
                0x02, 0x00, 0x00, 0x00,
                0x0D, 0x00, 0x00, 0x00,
                0x0D, 0x00, 0x00, 0x00,
                0x48, 0x65, 0x6C, 0x6C,
                0x6F, 0x20, 0x57, 0x6F,
                0x72, 0x6C, 0x64, 0x21,
                0x00, 0x0E, 0x00, 0x00,
                0x00, 0x0D, 0x00, 0x00,
                0x00, 0x48, 0x65, 0x6C,
                0x6C, 0x6F, 0x20, 0x57,
                0x6F, 0x72, 0x6C, 0x64,
                0x21, 0x00
            };
            assertByteEquals(expected, adsStampHeader.getBytes());
            assertEquals(Samples.of(2), adsStampHeader.getSamples());
            assertEquals(2, adsStampHeader.getAdsNotificationSamples().size());
            assertEquals(Data.of("Hello World!"), adsStampHeader.getAdsNotificationSamples().get(0).getData());
        }

    }

    @Test
    public void sizeCalculation() {
        AdsStampHeader adsStampHeader1 = AdsStampHeader.of(TimeStamp.of(13), Arrays.asList(AdsNotificationSample.of(NotificationHandle.of(13), Data.of("Hello World!")), AdsNotificationSample.of(NotificationHandle.of(14), Data.of("Hello World!"))));
        AdsStampHeader adsStampHeader2 = AdsStampHeader.of(TimeStamp.of(13), Samples.of(2), Arrays.asList(AdsNotificationSample.of(NotificationHandle.of(13), Data.of("Hello World!")), AdsNotificationSample.of(NotificationHandle.of(14), Data.of("Hello World!"))));
        assertEquals(adsStampHeader1, adsStampHeader2);
    }
}