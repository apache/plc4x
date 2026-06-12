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

package org.apache.plc4x.java.ads.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.List;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;

class AdsSubscriptionHandleTest {

    @Test
    void gettersExposeFields() {
        PlcSubscriber subscriber = Mockito.mock(PlcSubscriber.class);
        AdsSubscriptionHandle handle = new AdsSubscriptionHandle(subscriber, "tag", null, 42L);
        assertEquals("tag", handle.getTagName());
        assertNull(handle.getAdsDataType());
        assertEquals(42L, handle.getNotificationHandle());
    }

    @Test
    void registerDelegatesToSubscriber() {
        PlcSubscriber subscriber = Mockito.mock(PlcSubscriber.class);
        PlcConsumerRegistration registration = Mockito.mock(PlcConsumerRegistration.class);
        Consumer<PlcSubscriptionEvent> consumer = e -> { };
        Mockito.when(subscriber.registerConsumer(Mockito.eq(consumer), Mockito.anyList()))
            .thenReturn(registration);

        AdsSubscriptionHandle handle = new AdsSubscriptionHandle(subscriber, "tag", null, 42L);
        assertSame(registration, handle.register(consumer));
        Mockito.verify(subscriber).registerConsumer(Mockito.eq(consumer), Mockito.eq(List.of(handle)));
    }

    @Test
    void equalsAndHashCode() {
        PlcSubscriber subscriber = Mockito.mock(PlcSubscriber.class);
        AdsSubscriptionHandle a = new AdsSubscriptionHandle(subscriber, "tag", null, 1L);
        AdsSubscriptionHandle b = new AdsSubscriptionHandle(subscriber, "tag", null, 1L);
        AdsSubscriptionHandle c = new AdsSubscriptionHandle(subscriber, "other", null, 1L);
        assertEquals(a, a);
        assertEquals(a, b);
        assertNotEquals(a, c);
        assertNotEquals(a, "x");
        assertEquals(a.hashCode(), b.hashCode());
    }

    @Test
    void toStringContainsFields() {
        AdsSubscriptionHandle handle = new AdsSubscriptionHandle(
            Mockito.mock(PlcSubscriber.class), "tag", null, 99L);
        assertTrue(handle.toString().contains("tagName='tag'"));
        assertTrue(handle.toString().contains("99"));
    }
}
