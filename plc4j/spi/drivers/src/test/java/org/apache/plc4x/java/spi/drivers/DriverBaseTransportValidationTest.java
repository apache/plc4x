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

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Verifies the SPI-core check that a driver may only be opened over one of the transports it
 * declares it supports, plus the {@code allow-unsupported-transport} opt-out.
 *
 * <p>The tests exercise the guard purely at the {@code DriverBase.getConnection(...)} boundary,
 * with no network and no registered transport, by distinguishing two failure messages:</p>
 * <ul>
 *   <li><b>Guard rejection</b> — message contains {@value #GUARD_MARKER} (the driver does not
 *       support the requested transport).</li>
 *   <li><b>Pre-existing registered-transport failure</b> — message contains
 *       {@value #REGISTERED_MARKER} (the transport passed the guard but is not a registered
 *       transport in this test JVM).</li>
 * </ul>
 */
class DriverBaseTransportValidationTest {

    /** Substring unique to the driver-supported-transport rejection message. */
    private static final String GUARD_MARKER = "is not supported by driver";
    /** Substring of the pre-existing "transport not registered at all" failure message. */
    private static final String REGISTERED_MARKER = "Unsupported transport";

    /**
     * Configurable {@link DriverBase} stub: protocol code, optional default transport, and the
     * declared supported-transport list are all set per test. It never actually connects.
     */
    static final class StubDriver extends DriverBase {
        private final String protocolCode;
        private final String defaultTransport;          // nullable -> no default
        private final List<String> supportedTransports; // may be empty

        StubDriver(String protocolCode, String defaultTransport, List<String> supportedTransports) {
            this.protocolCode = protocolCode;
            this.defaultTransport = defaultTransport;
            this.supportedTransports = supportedTransports;
        }

        @Override public String getProtocolCode() { return protocolCode; }
        @Override public String getProtocolName() { return protocolCode; }
        @Override public Optional<String> getDefaultTransportCode() { return Optional.ofNullable(defaultTransport); }
        @Override public List<String> getSupportedTransportCodes() { return supportedTransports; }
        @Override protected Class<? extends Configuration> getConfigurationClass() { return Configuration.class; }
        @Override protected ConnectionBase<?> getConnection(Configuration c, TransportInstance<?> t, AuditLog a) {
            throw new UnsupportedOperationException("test stub does not connect");
        }
    }

    /** A driver that supports exactly one transport, "supported", with no default. */
    private static StubDriver singleSupportDriver() {
        return new StubDriver("stub", null, List.of("supported"));
    }

    /** Open the connection, assert it throws, and return the (non-null) exception message. */
    private static String messageFromFailedConnect(DriverBase driver, String connectionString) {
        Throwable t = assertThrows(Throwable.class, () -> driver.getConnection(connectionString));
        String msg = t.getMessage();
        assertTrue(msg != null && !msg.isBlank(), "expected a non-empty failure message");
        return msg;
    }

    @Nested
    class StrictCheck {

        @Test
        void unsupportedTransportIsRejectedWithRequestedAndSupportedNamed() {
            String msg = messageFromFailedConnect(singleSupportDriver(), "stub:other://host");
            assertTrue(msg.contains(GUARD_MARKER), "should be the guard rejection: " + msg);
            assertTrue(msg.contains("other"), "message must name the requested transport: " + msg);
            assertTrue(msg.contains("supported"), "message must list the supported transport(s): " + msg);
        }

        @Test
        void supportedTransportPassesTheGuard() {
            String msg = messageFromFailedConnect(singleSupportDriver(), "stub:supported://host");
            assertFalse(msg.contains(GUARD_MARKER), "guard must NOT reject a supported transport: " + msg);
            assertTrue(msg.contains(REGISTERED_MARKER), "should reach the registered-transport lookup: " + msg);
        }

        @Test
        void supportedButUnregisteredTransportStillFailsPreExisting() {
            StubDriver driver = new StubDriver("stub", null, List.of("ghost"));
            String msg = messageFromFailedConnect(driver, "stub:ghost://host");
            assertFalse(msg.contains(GUARD_MARKER), "guard passed (ghost is declared supported): " + msg);
            assertTrue(msg.contains(REGISTERED_MARKER), "pre-existing unregistered-transport failure: " + msg);
        }
    }

    @Nested
    class OptOut {

        @Test
        void optOutAllowsUnsupportedTransportThroughGuard() {
            String msg = messageFromFailedConnect(singleSupportDriver(),
                "stub:other://host?allow-unsupported-transport=true");
            assertFalse(msg.contains(GUARD_MARKER), "opt-out must skip the guard: " + msg);
        }

        @Test
        void optOutDoesNotBypassRegisteredTransportCheck() {
            String msg = messageFromFailedConnect(singleSupportDriver(),
                "stub:other://host?allow-unsupported-transport=true");
            assertTrue(msg.contains(REGISTERED_MARKER),
                "registered-transport lookup must still run under opt-out: " + msg);
        }

        @Test
        void invalidOptOutValueIsTreatedAsStrict() {
            String msg = messageFromFailedConnect(singleSupportDriver(),
                "stub:other://host?allow-unsupported-transport=notabool");
            assertTrue(msg.contains(GUARD_MARKER),
                "invalid opt-out value must fall back to strict and reject: " + msg);
        }
    }

    @Nested
    class DefaultOnlyDriver {

        /** No explicit supported list, but a default transport "supported". */
        private StubDriver defaultOnlyDriver() {
            return new StubDriver("stub", "supported", List.of());
        }

        @Test
        void defaultTransportExplicitPassesTheGuard() {
            String msg = messageFromFailedConnect(defaultOnlyDriver(), "stub:supported://host");
            assertFalse(msg.contains(GUARD_MARKER), "default transport must pass the guard: " + msg);
            assertTrue(msg.contains(REGISTERED_MARKER), msg);
        }

        @Test
        void defaultTransportOmittedPassesTheGuard() {
            String msg = messageFromFailedConnect(defaultOnlyDriver(), "stub://host");
            assertFalse(msg.contains(GUARD_MARKER), "omitted code -> default must pass the guard: " + msg);
            assertTrue(msg.contains(REGISTERED_MARKER), msg);
        }

        @Test
        void otherTransportRejectedForDefaultOnlyDriver() {
            String msg = messageFromFailedConnect(defaultOnlyDriver(), "stub:other://host");
            assertTrue(msg.contains(GUARD_MARKER), "non-default transport must be rejected: " + msg);
            assertTrue(msg.contains("other"), msg);
        }
    }
}
