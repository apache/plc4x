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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.junit.jupiter.api.Assertions.*;
import org.apache.plc4x.java.s7.userdata.S7SzlService;

import static org.mockito.Mockito.*;

/**
 * Mock-based tests for the S7H dual-path failover wrapper. These pin down the failover
 * state machine without needing real PLCs — they construct {@link S7HCotpConnection} with
 * mocked inner {@link S7CotpConnection}s and exercise the public/protected API surface.
 *
 * <p>Live end-to-end validation lives in the manual tests (run against real hardware);
 * this file catches regressions during refactors.
 */
class S7HCotpConnectionTest {

    @Test
    void metadata_subscribeSupported_isMaskedFalseEvenIfInnerSupportsIt() {
        // The wrapper hard-blocks subscriptions because PLC-side subscription state can't
        // be transparently transferred to the survivor on failover. Even if both inners
        // report subscribe-supported, the wrapper must report false to keep callers from
        // setting up streams that silently break on the next swap.
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        PlcConnectionMetadata innerMeta = mock(PlcConnectionMetadata.class);
        when(primary.getMetadata()).thenReturn(innerMeta);
        when(innerMeta.isReadSupported()).thenReturn(true);
        when(innerMeta.isWriteSupported()).thenReturn(true);
        when(innerMeta.isSubscribeSupported()).thenReturn(true);
        when(innerMeta.isBrowseSupported()).thenReturn(true);

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);

        PlcConnectionMetadata meta = wrapper.getMetadata();
        assertTrue(meta.isReadSupported());
        assertTrue(meta.isWriteSupported());
        assertFalse(meta.isSubscribeSupported(),
            "S7H must mask subscribe-supported to false regardless of inner");
        assertTrue(meta.isBrowseSupported());
    }

    @Test
    void onSubscribe_rejectsWithClearMessage() {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);
        org.apache.plc4x.java.api.messages.PlcSubscriptionRequest request =
            mock(org.apache.plc4x.java.api.messages.PlcSubscriptionRequest.class);

        ExecutionException ex = assertThrows(ExecutionException.class,
            () -> wrapper.onSubscribeForTest(request).get(1, TimeUnit.SECONDS));
        assertInstanceOf(PlcRuntimeException.class, ex.getCause());
        assertTrue(ex.getCause().getMessage().contains("Subscribe not supported"),
            "Error message should explain S7H limitation. Was: " + ex.getCause().getMessage());
    }

    @Test
    void isConnected_returnsTrueIfEitherInnerIsConnected() {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(false);
        when(secondary.isConnected()).thenReturn(true);
        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);

        assertTrue(wrapper.isConnected());

        when(secondary.isConnected()).thenReturn(false);
        assertFalse(wrapper.isConnected());
    }

    @Test
    void onRead_routesToActiveAndFailsOverOnTimeout() throws Exception {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(true);
        when(secondary.isConnected()).thenReturn(true);

        // Primary's read never completes (simulates a hung TCP after cable pull).
        CompletableFuture<PlcReadResponse> hung = new CompletableFuture<>();
        // Secondary returns a real response.
        PlcReadResponse expected = mock(PlcReadResponse.class);
        CompletableFuture<PlcReadResponse> done = CompletableFuture.completedFuture(expected);
        when(primary.read(any())).thenReturn(hung);
        when(secondary.read(any())).thenReturn(done);

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);
        wrapper.markActiveAsPrimaryForTest();

        PlcReadRequest request = mock(PlcReadRequest.class);
        // The wrapper-level timeout is 2s by default; expect the read to complete around
        // that deadline via the secondary, not after the test's own 5s timeout.
        PlcReadResponse actual = wrapper.onReadForTest(request).get(5, TimeUnit.SECONDS);
        assertSame(expected, actual);
        verify(primary, times(1)).read(any());
        verify(secondary, times(1)).read(any());
    }

    @Test
    void onRead_returnsFailureWhenBothInnersAreDown() {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(false);
        when(secondary.isConnected()).thenReturn(false);

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);

        PlcReadRequest request = mock(PlcReadRequest.class);
        ExecutionException ex = assertThrows(ExecutionException.class,
            () -> wrapper.onReadForTest(request).get(1, TimeUnit.SECONDS));
        assertInstanceOf(PlcRuntimeException.class, ex.getCause());
        assertTrue(ex.getCause().getMessage().contains("Neither primary nor secondary"));
    }

    @Test
    void onRead_failedFutureFromActive_swapsAndRetriesOnAlt() throws Exception {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(true);
        when(secondary.isConnected()).thenReturn(true);

        // Primary fails fast; secondary succeeds.
        when(primary.read(any())).thenReturn(
            CompletableFuture.failedFuture(new TimeoutException("kernel-detected drop")));
        PlcReadResponse expected = mock(PlcReadResponse.class);
        when(secondary.read(any())).thenReturn(CompletableFuture.completedFuture(expected));

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);
        wrapper.markActiveAsPrimaryForTest();

        PlcReadResponse actual = wrapper.onReadForTest(mock(PlcReadRequest.class))
            .get(2, TimeUnit.SECONDS);
        assertSame(expected, actual);
        verify(primary).read(any());
        verify(secondary).read(any());
    }

    @Test
    void onRead_singleHealthyInner_doesntTryAlt() throws Exception {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(true);
        when(secondary.isConnected()).thenReturn(false);
        PlcReadResponse expected = mock(PlcReadResponse.class);
        when(primary.read(any())).thenReturn(CompletableFuture.completedFuture(expected));

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);
        wrapper.markActiveAsPrimaryForTest();

        assertSame(expected, wrapper.onReadForTest(mock(PlcReadRequest.class)).get());
        verify(primary).read(any());
        verify(secondary, never()).read(any());
    }

    // ------------------------------------------------------------------------------------
    // Test fixture helpers
    // ------------------------------------------------------------------------------------

    @Test
    void readDeviceIdentification_readsFromTheActiveInner() throws Exception {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(true);
        when(secondary.isConnected()).thenReturn(true);
        when(primary.readDeviceIdentification())
            .thenReturn(CompletableFuture.completedFuture(identificationWithFirmware("V 2.6.0")));

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);
        wrapper.markActiveAsPrimaryForTest();

        assertEquals("V 2.6.0",
            wrapper.readDeviceIdentification().get(5, TimeUnit.SECONDS).firmwareVersion());
        verify(secondary, never()).readDeviceIdentification();
    }

    @Test
    void readDeviceIdentification_failsOverToTheSurvivingInner() throws Exception {
        S7CotpConnection primary = mock(S7CotpConnection.class);
        S7CotpConnection secondary = mock(S7CotpConnection.class);
        when(primary.isConnected()).thenReturn(false);
        when(secondary.isConnected()).thenReturn(true);
        when(secondary.readDeviceIdentification())
            .thenReturn(CompletableFuture.completedFuture(identificationWithFirmware("V 4.2.1")));

        TestableS7HCotpConnection wrapper = newWrapper(primary, secondary);
        // A wrapper only reaches this call after onConnect() has picked an active slot, so
        // the failover under test is "active inner went down", not "never connected".
        wrapper.markActiveAsPrimaryForTest();

        assertEquals("V 4.2.1",
            wrapper.readDeviceIdentification().get(5, TimeUnit.SECONDS).firmwareVersion());
    }

    private static S7SzlService.S7DeviceIdentification identificationWithFirmware(String firmware) {
        return new S7SzlService.S7DeviceIdentification(null, null, firmware,
            null, null, null, null, null, null, null, null, null, null);
    }

    private static TestableS7HCotpConnection newWrapper(S7CotpConnection primary, S7CotpConnection secondary) {
        S7Configuration config = new S7Configuration();
        TransportInstance<?> transport = mock(TransportInstance.class);
        AuditLog auditLog = mock(AuditLog.class);
        return new TestableS7HCotpConnection(config, transport, auditLog,
            () -> primary, () -> secondary);
    }

    /**
     * Subclass that exposes a few internals as public test entry points without forcing us
     * to bend the production class. Specifically: route requests through the protected
     * {@code onRead/onSubscribe} hooks (rather than going through the parent's request
     * builder + throttle) and force the active slot for cases that need to bypass
     * onConnect.
     */
    private static class TestableS7HCotpConnection extends S7HCotpConnection {
        TestableS7HCotpConnection(S7Configuration cfg, TransportInstance<?> t, AuditLog a,
                                  java.util.function.Supplier<S7CotpConnection> p,
                                  java.util.function.Supplier<S7CotpConnection> s) {
            super(cfg, t, a, p, s);
        }
        CompletableFuture<PlcReadResponse> onReadForTest(PlcReadRequest req) {
            return onRead(req);
        }
        CompletableFuture<org.apache.plc4x.java.api.messages.PlcSubscriptionResponse>
                onSubscribeForTest(org.apache.plc4x.java.api.messages.PlcSubscriptionRequest req) {
            return onSubscribe(req);
        }
        void markActiveAsPrimaryForTest() {
            // Reach into the parent to flip the active slot. Without this, the wrapper's
            // chooseSlot fallback to "any healthy inner" still works (the primary is the
            // first one tried), so this mostly serves as documentation that the test is
            // exercising the active-primary branch rather than the alt-fallback branch.
            try {
                java.lang.reflect.Field f = S7HCotpConnection.class.getDeclaredField("activeSlot");
                f.setAccessible(true);
                java.util.concurrent.atomic.AtomicReference<?> ar =
                    (java.util.concurrent.atomic.AtomicReference<?>) f.get(this);
                java.lang.reflect.Field pr = S7HCotpConnection.class.getDeclaredField("primaryRef");
                pr.setAccessible(true);
                @SuppressWarnings({"unchecked", "rawtypes"})
                java.util.concurrent.atomic.AtomicReference primary =
                    (java.util.concurrent.atomic.AtomicReference) pr.get(this);
                @SuppressWarnings({"unchecked", "rawtypes"})
                java.util.concurrent.atomic.AtomicReference activeSlot =
                    (java.util.concurrent.atomic.AtomicReference) ar;
                activeSlot.set(primary);
            } catch (ReflectiveOperationException e) {
                throw new RuntimeException(e);
            }
        }
    }
}
