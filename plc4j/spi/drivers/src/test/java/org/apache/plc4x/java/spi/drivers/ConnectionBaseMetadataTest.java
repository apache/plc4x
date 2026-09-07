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
package org.apache.plc4x.java.spi.drivers;

import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.junit.jupiter.api.Test;

import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The connection metadata is derived from the operations a driver implements. It used to be a fixed
 * "everything is supported", which meant a driver that implements nothing still claimed to support
 * reading, writing, subscribing and browsing.
 */
class ConnectionBaseMetadataTest {

    /**
     * A driver that implements none of the operations - it can connect and nothing else.
     */
    @Test
    void reportsNothingForAConnectionThatImplementsNoOperation() {
        PlcConnectionMetadata metadata = new EmptyConnection().getMetadata();

        assertFalse(metadata.isReadSupported());
        assertFalse(metadata.isWriteSupported());
        assertFalse(metadata.isSubscribeSupported());
        assertFalse(metadata.isBrowseSupported());
    }

    @Test
    void reportsExactlyTheOperationsTheDriverImplements() {
        PlcConnectionMetadata metadata = new ReadWriteConnection().getMetadata();

        assertTrue(metadata.isReadSupported());
        assertTrue(metadata.isWriteSupported());
        assertFalse(metadata.isSubscribeSupported(), "onSubscribe is not overridden");
        assertFalse(metadata.isBrowseSupported(), "onBrowse is not overridden");
    }

    /**
     * An operation provided by an intermediate base class counts as implemented. This is how
     * PollingSubscriptionConnectionBase gives protocols without native subscriptions (Modbus and
     * friends) a working subscribe.
     */
    @Test
    void countsAnOperationInheritedFromAnIntermediateBase() {
        PlcConnectionMetadata metadata = new InheritsSubscribeConnection().getMetadata();

        assertTrue(metadata.isSubscribeSupported(), "onSubscribe comes from the intermediate base");
        assertFalse(metadata.isReadSupported());
    }

    /**
     * Browsing counts when either browse hook is implemented.
     */
    @Test
    void countsBrowseFromTheInterceptorHookAsWell() {
        assertTrue(new BrowseConnection().getMetadata().isBrowseSupported());
    }

    /**
     * A driver whose capabilities depend on the device it connected to overrides this and answers
     * from what it learned during connect. Nothing in the base class may cache across connections
     * and defeat that.
     */
    @Test
    void letsADriverAnswerPerConnection() {
        DeviceDependentConnection limited = new DeviceDependentConnection(false);
        DeviceDependentConnection capable = new DeviceDependentConnection(true);

        assertFalse(limited.getMetadata().isWriteSupported());
        assertTrue(capable.getMetadata().isWriteSupported(),
            "two connections of the same class must be able to answer differently");
    }

    private abstract static class TestConnection extends ConnectionBase<Configuration> {
        TestConnection() {
            super(null, null, null);
        }

        @Override
        public PlcTagHandler getTagHandler() {
            return null;
        }

        @Override
        public PlcValueHandler getValueHandler() {
            return null;
        }
    }

    private static class EmptyConnection extends TestConnection {
    }

    private static class ReadWriteConnection extends TestConnection {
        @Override
        protected CompletableFuture<PlcReadResponse> onRead(PlcReadRequest readRequest) {
            return CompletableFuture.completedFuture(null);
        }

        @Override
        protected CompletableFuture<PlcWriteResponse> onWrite(PlcWriteRequest writeRequest) {
            return CompletableFuture.completedFuture(null);
        }
    }

    private abstract static class SubscribingBase extends TestConnection {
        @Override
        protected CompletableFuture<PlcSubscriptionResponse> onSubscribe(PlcSubscriptionRequest subscriptionRequest) {
            return CompletableFuture.completedFuture(null);
        }
    }

    private static class InheritsSubscribeConnection extends SubscribingBase {
    }

    private static class BrowseConnection extends TestConnection {
        @Override
        protected CompletableFuture<PlcBrowseResponse> onBrowse(PlcBrowseRequest browseRequest) {
            return CompletableFuture.completedFuture(null);
        }
    }

    /**
     * Speaks writing, but whether the device accepts it is only known once connected.
     */
    private static class DeviceDependentConnection extends ReadWriteConnection {
        private final boolean deviceAllowsWriting;

        DeviceDependentConnection(boolean deviceAllowsWriting) {
            this.deviceAllowsWriting = deviceAllowsWriting;
        }

        @Override
        public PlcConnectionMetadata getMetadata() {
            PlcConnectionMetadata derived = super.getMetadata();
            return new PlcConnectionMetadata() {
                @Override
                public boolean isReadSupported() {
                    return derived.isReadSupported();
                }

                @Override
                public boolean isWriteSupported() {
                    return derived.isWriteSupported() && deviceAllowsWriting;
                }

                @Override
                public boolean isSubscribeSupported() {
                    return derived.isSubscribeSupported();
                }

                @Override
                public boolean isBrowseSupported() {
                    return derived.isBrowseSupported();
                }
            };
        }
    }
}
