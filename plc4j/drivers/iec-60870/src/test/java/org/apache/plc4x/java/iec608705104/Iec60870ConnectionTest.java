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
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.iec608705104.model.Iec608705104SubscriptionHandle;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Covers the subscription dispatch decision. Everything else in
 * {@link Iec60870Connection} needs a live station and is exercised by the
 * (currently disabled) driver testsuite.
 */
class Iec60870ConnectionTest {

    private static PlcConsumerRegistration registrationFor(String... tagAddresses) {
        PlcSubscriber subscriber = mock(PlcSubscriber.class);
        List<PlcSubscriptionHandle> handles = new ArrayList<>();
        for (String tagAddress : tagAddresses) {
            handles.add(new Iec608705104SubscriptionHandle(subscriber, Iec608705104Tag.of(tagAddress)));
        }
        PlcConsumerRegistration registration = mock(PlcConsumerRegistration.class);
        when(registration.getSubscriptionHandles()).thenReturn(handles);
        return registration;
    }

    @Test
    void anExactlyAddressedHandleOnlyMatchesItsOwnObject() {
        PlcConsumerRegistration registration = registrationFor("1/2");
        assertTrue(Iec60870Connection.matchesAnyHandle(registration, 1, 2));
        assertFalse(Iec60870Connection.matchesAnyHandle(registration, 1, 3));
        assertFalse(Iec60870Connection.matchesAnyHandle(registration, 2, 2));
    }

    @Test
    void wildcardHandlesWiden() {
        assertTrue(Iec60870Connection.matchesAnyHandle(registrationFor("*/*"), 4711, 815));
        assertTrue(Iec60870Connection.matchesAnyHandle(registrationFor("1/*"), 1, 815));
        assertFalse(Iec60870Connection.matchesAnyHandle(registrationFor("1/*"), 2, 815));
    }

    @Test
    void anyOfSeveralHandlesIsEnough() {
        PlcConsumerRegistration registration = registrationFor("1/2", "3/*");
        assertTrue(Iec60870Connection.matchesAnyHandle(registration, 1, 2));
        assertTrue(Iec60870Connection.matchesAnyHandle(registration, 3, 99));
        assertFalse(Iec60870Connection.matchesAnyHandle(registration, 2, 2));
    }

    @Test
    void registrationsWithoutIecHandlesNeverMatch() {
        PlcConsumerRegistration empty = mock(PlcConsumerRegistration.class);
        when(empty.getSubscriptionHandles()).thenReturn(List.of());
        assertFalse(Iec60870Connection.matchesAnyHandle(empty, 0, 0));

        PlcConsumerRegistration foreign = mock(PlcConsumerRegistration.class);
        when(foreign.getSubscriptionHandles()).thenReturn(List.of(mock(PlcSubscriptionHandle.class)));
        assertFalse(Iec60870Connection.matchesAnyHandle(foreign, 0, 0));
    }

}
