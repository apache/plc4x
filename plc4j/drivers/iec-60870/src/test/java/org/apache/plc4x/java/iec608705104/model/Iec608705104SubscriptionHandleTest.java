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
package org.apache.plc4x.java.iec608705104.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.iec608705104.model.Iec608705104SubscriptionHandle;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentCaptor;

import java.util.Collection;
import java.util.List;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

class Iec608705104SubscriptionHandleTest {

    @Test
    void registerForwardsToTheSubscriber() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        Iec608705104Tag tag = new Iec608705104Tag(1, 2);
        Iec608705104SubscriptionHandle handle = new Iec608705104SubscriptionHandle(subscriber, tag);
        @SuppressWarnings("unchecked")
        Consumer<PlcSubscriptionEvent> consumer = mock(Consumer.class);
        PlcConsumerRegistration registration = mock(PlcConsumerRegistration.class);
        when(subscriber.registerConsumer(eq(consumer), any())).thenReturn(registration);

        assertSame(registration, handle.register(consumer));

        @SuppressWarnings("unchecked")
        ArgumentCaptor<Collection<PlcSubscriptionHandle>> handlesArg =
            ArgumentCaptor.forClass(Collection.class);
        verify(subscriber).registerConsumer(eq(consumer), handlesArg.capture());
        assertEquals(List.of(handle), handlesArg.getValue());
    }

    @Test
    void getterExposesTag() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        Iec608705104Tag tag = new Iec608705104Tag(3, 4);
        assertSame(tag, new Iec608705104SubscriptionHandle(subscriber, tag).getTag());
    }

    @Test
    void equalsHashCodeKeyedOnTagOnly() {
        // Two handles that point at the same information object must be
        // treated as equivalent even when they were created by different
        // subscriber instances — otherwise the consumer registry would
        // dedupe by identity and double-deliver events.
        PlcSubscriber a = mock(PlcSubscriber.class);
        PlcSubscriber b = mock(PlcSubscriber.class);
        Iec608705104Tag sharedTag = new Iec608705104Tag(5, 6);

        Iec608705104SubscriptionHandle h1 = new Iec608705104SubscriptionHandle(a, sharedTag);
        Iec608705104SubscriptionHandle h2 = new Iec608705104SubscriptionHandle(b, sharedTag);
        Iec608705104SubscriptionHandle h3 = new Iec608705104SubscriptionHandle(a, new Iec608705104Tag(5, 7));

        assertEquals(h1, h2);
        assertEquals(h1.hashCode(), h2.hashCode());
        assertNotEquals(h1, h3);
        assertEquals(h1, h1);
        assertNotEquals(h1, "not a handle");
    }

    @Test
    void toStringMentionsTag() {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        String s = new Iec608705104SubscriptionHandle(subscriber, new Iec608705104Tag(9, 0)).toString();
        assertTrue(s.contains("Iec608705104Tag"));
    }

}
