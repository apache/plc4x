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
package org.apache.plc4x.java.firmata.model;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.firmata.tag.FirmataTag;
import org.apache.plc4x.java.firmata.tag.FirmataTagDigital;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.junit.jupiter.api.Test;

import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.atomic.AtomicReference;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FirmataSubscriptionHandleTest {

    private final FirmataTag pinTwo = FirmataTagDigital.of("digital:2");
    private final FirmataTag pinThree = FirmataTagDigital.of("digital:3");

    @Test
    void registerDelegatesToSubscriber() {
        AtomicReference<Consumer<PlcSubscriptionEvent>> capturedConsumer = new AtomicReference<>();
        AtomicReference<Collection<PlcSubscriptionHandle>> capturedHandles = new AtomicReference<>();
        PlcConsumerRegistration stubRegistration = new StubRegistration();
        RecordingSubscriber subscriber = new RecordingSubscriber(capturedConsumer, capturedHandles, stubRegistration);

        FirmataSubscriptionHandle handle = new FirmataSubscriptionHandle(subscriber, "button", pinTwo);
        Consumer<PlcSubscriptionEvent> consumer = event -> { /* no-op */ };

        PlcConsumerRegistration result = handle.register(consumer);

        assertSame(stubRegistration, result);
        assertSame(consumer, capturedConsumer.get());
        assertTrue(capturedHandles.get().contains(handle));
        assertEquals(1, capturedHandles.get().size());
    }

    @Test
    void gettersExposeConstructorArguments() {
        PlcSubscriber subscriber = new RecordingSubscriber(null, null, new StubRegistration());
        FirmataSubscriptionHandle handle = new FirmataSubscriptionHandle(subscriber, "button", pinTwo);
        assertEquals("button", handle.getName());
        assertSame(pinTwo, handle.getTag());
    }

    @Test
    void nullArgumentsRejected() {
        PlcSubscriber subscriber = new RecordingSubscriber(null, null, new StubRegistration());
        assertThrows(NullPointerException.class,
            () -> new FirmataSubscriptionHandle(null, "x", pinTwo));
        assertThrows(NullPointerException.class,
            () -> new FirmataSubscriptionHandle(subscriber, null, pinTwo));
        assertThrows(NullPointerException.class,
            () -> new FirmataSubscriptionHandle(subscriber, "x", null));
    }

    @Test
    void equalityAndHashUseNameAndTag() {
        PlcSubscriber subscriber = new RecordingSubscriber(null, null, new StubRegistration());
        FirmataSubscriptionHandle a = new FirmataSubscriptionHandle(subscriber, "button", pinTwo);
        FirmataSubscriptionHandle b = new FirmataSubscriptionHandle(subscriber, "button", pinTwo);
        FirmataSubscriptionHandle differentName = new FirmataSubscriptionHandle(subscriber, "led", pinTwo);
        FirmataSubscriptionHandle differentTag  = new FirmataSubscriptionHandle(subscriber, "button", pinThree);

        assertEquals(a, b);
        assertEquals(a.hashCode(), b.hashCode());
        assertNotEquals(a, differentName);
        assertNotEquals(a, differentTag);
        assertNotEquals(a, "not a handle");
        assertEquals(a, a);
        assertNotEquals(null, a);
        assertTrue(a.toString().contains("button"));
    }

    /** Minimal {@link PlcConsumerRegistration} stand-in for test stubs. */
    private static final class StubRegistration implements PlcConsumerRegistration {
        @Override
        public Integer getConsumerId() {
            return 0;
        }

        @Override
        public List<PlcSubscriptionHandle> getSubscriptionHandles() {
            return Collections.emptyList();
        }

        @Override
        public void unregister() {
            // no-op
        }
    }

    /**
     * {@link PlcSubscriber} that records the arguments it was called with so
     * the test can assert the handle's register() forwarded them correctly.
     */
    private static final class RecordingSubscriber implements PlcSubscriber {
        private final AtomicReference<Consumer<PlcSubscriptionEvent>> capturedConsumer;
        private final AtomicReference<Collection<PlcSubscriptionHandle>> capturedHandles;
        private final PlcConsumerRegistration registrationToReturn;

        RecordingSubscriber(AtomicReference<Consumer<PlcSubscriptionEvent>> capturedConsumer,
                            AtomicReference<Collection<PlcSubscriptionHandle>> capturedHandles,
                            PlcConsumerRegistration registrationToReturn) {
            this.capturedConsumer = capturedConsumer;
            this.capturedHandles = capturedHandles;
            this.registrationToReturn = registrationToReturn;
        }

        @Override
        public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
            throw new UnsupportedOperationException("not under test");
        }

        @Override
        public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
            throw new UnsupportedOperationException("not under test");
        }

        @Override
        public PlcConsumerRegistration registerConsumer(Consumer<PlcSubscriptionEvent> consumer,
                                                        Collection<PlcSubscriptionHandle> handles) {
            if (capturedConsumer != null) {
                capturedConsumer.set(consumer);
            }
            if (capturedHandles != null) {
                capturedHandles.set(handles);
            }
            return registrationToReturn;
        }

        @Override
        public void unregisterConsumer(PlcConsumerRegistration plcConsumerRegistration) {
            // no-op
        }
    }

}
