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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.Collection;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertTrue;

class DefaultPlcSubscriptionResponseTest {

    /**
     * Drivers commonly create one handle for the whole request and map it to every tag name.
     * getSubscriptionHandles() must report that handle once, otherwise the usual "register a
     * consumer on every handle" loop delivers each event once per tag - see GH-1896.
     */
    @Test
    void oneHandleSharedByManyTagsIsReportedOnce() {
        PlcSubscriptionHandle handle = new TestHandle();
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        for (int i = 0; i < 5; i++) {
            values.put("tag" + i, new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
        }

        Collection<PlcSubscriptionHandle> handles = response(values).getSubscriptionHandles();

        assertEquals(1, handles.size());
        assertSame(handle, handles.iterator().next());
    }

    /**
     * Distinct handles must all survive, even when they compare equal - two subscriptions are two
     * subscriptions.
     */
    @Test
    void distinctHandlesAreAllReported() {
        PlcSubscriptionHandle first = new AlwaysEqualHandle();
        PlcSubscriptionHandle second = new AlwaysEqualHandle();
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        values.put("a", new DefaultPlcResponseItem<>(PlcResponseCode.OK, first));
        values.put("b", new DefaultPlcResponseItem<>(PlcResponseCode.OK, second));

        Collection<PlcSubscriptionHandle> handles = response(values).getSubscriptionHandles();

        assertEquals(2, handles.size());
        assertTrue(handles.stream().anyMatch(h -> h == first));
        assertTrue(handles.stream().anyMatch(h -> h == second));
    }

    /**
     * A mix: two tags share one handle, a third has its own.
     */
    @Test
    void mixedSharedAndDistinctHandles() {
        PlcSubscriptionHandle shared = new TestHandle();
        PlcSubscriptionHandle own = new TestHandle();
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        values.put("a", new DefaultPlcResponseItem<>(PlcResponseCode.OK, shared));
        values.put("b", new DefaultPlcResponseItem<>(PlcResponseCode.OK, shared));
        values.put("c", new DefaultPlcResponseItem<>(PlcResponseCode.OK, own));

        assertEquals(2, response(values).getSubscriptionHandles().size());
    }

    /**
     * Tags that failed to subscribe carry a null handle and must simply be skipped.
     */
    @Test
    void tagsWithoutAHandleAreSkipped() {
        PlcSubscriptionHandle handle = new TestHandle();
        Map<String, PlcResponseItem<PlcSubscriptionHandle>> values = new LinkedHashMap<>();
        values.put("ok", new DefaultPlcResponseItem<>(PlcResponseCode.OK, handle));
        values.put("bad", new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));

        assertEquals(1, response(values).getSubscriptionHandles().size());
    }

    private DefaultPlcSubscriptionResponse response(Map<String, PlcResponseItem<PlcSubscriptionHandle>> values) {
        PlcSubscriptionRequest request = Mockito.mock(PlcSubscriptionRequest.class);
        Mockito.lenient().when(request.getTagNames()).thenReturn(new LinkedHashSet<>(values.keySet()));
        return new DefaultPlcSubscriptionResponse(request, values);
    }

    private static class TestHandle implements PlcSubscriptionHandle {
        @Override
        public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
            return null;
        }
    }

    /** Compares equal to every other instance, to prove de-duplication is by identity. */
    private static class AlwaysEqualHandle implements PlcSubscriptionHandle {
        @Override
        public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer) {
            return null;
        }

        @Override
        public boolean equals(Object obj) {
            return obj instanceof AlwaysEqualHandle;
        }

        @Override
        public int hashCode() {
            return 1;
        }
    }
}
