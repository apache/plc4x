/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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

import static org.apache.plc4x.java.mock.util.Assert.assertByteEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class AdsNotificationSampleTest {

    @Test
    public void builder() throws Exception {
        {
            AdsNotificationSample adsNotificationSample = AdsNotificationSample.of(NotificationHandle.of(13), Data.of("Hello World!"));
            assertNotNull(adsNotificationSample);
            byte[] expected = {
                0x0D, 0x00, 0x00, 0x00,
                0x0D, 0x00, 0x00, 0x00,
                0x48, 0x65, 0x6C, 0x6C,
                0x6F, 0x20, 0x57, 0x6F,
                0x72, 0x6C, 0x64, 0x21,
                0x00
            };
            assertByteEquals(expected, adsNotificationSample.getBytes());
            assertEquals(SampleSize.of(13), adsNotificationSample.getSampleSize());
            assertEquals(Data.of("Hello World!"), adsNotificationSample.getData());
        }
        {
            AdsNotificationSample adsNotificationSample = AdsNotificationSample.of(NotificationHandle.of(13), SampleSize.of(13), Data.of("Hello World!"));
            assertNotNull(adsNotificationSample);
            byte[] expected = {
                0x0D, 0x00, 0x00, 0x00,
                0x0D, 0x00, 0x00, 0x00,
                0x48, 0x65, 0x6C, 0x6C,
                0x6F, 0x20, 0x57, 0x6F,
                0x72, 0x6C, 0x64, 0x21,
                0x00
            };
            assertByteEquals(expected, adsNotificationSample.getBytes());
            assertEquals(SampleSize.of(13), adsNotificationSample.getSampleSize());
            assertEquals(Data.of("Hello World!"), adsNotificationSample.getData());
        }
    }

    @Test
    public void sizeCalculation() {
        AdsNotificationSample adsNotificationSample1 = AdsNotificationSample.of(NotificationHandle.of(13), SampleSize.of(13), Data.of("Hello World!"));
        AdsNotificationSample adsNotificationSample2 = AdsNotificationSample.of(NotificationHandle.of(13), Data.of("Hello World!"));
        assertEquals(adsNotificationSample1, adsNotificationSample2);
    }
}