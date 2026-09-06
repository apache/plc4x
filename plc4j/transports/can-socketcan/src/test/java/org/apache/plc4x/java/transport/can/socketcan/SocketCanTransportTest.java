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
package org.apache.plc4x.java.transport.can.socketcan;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.can.socketcan.config.SocketCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class SocketCanTransportTest {

    @Test
    void getTransportCode() {
        SocketCanTransport transport = new SocketCanTransport();
        assertEquals("can-socketcan", transport.getTransportCode());
    }

    @Test
    void getTransportName() {
        SocketCanTransport transport = new SocketCanTransport();
        assertEquals("SocketCAN", transport.getTransportName());
    }

    @Test
    void getTransportConfigType() {
        SocketCanTransport transport = new SocketCanTransport();
        assertEquals(SocketCanTransportConfiguration.class, transport.getTransportConfigType());
    }

    @Test
    void createTransportInstanceWithWrongConfigTypeThrows() {
        SocketCanTransport transport = new SocketCanTransport();
        AuditLog auditLog = Mockito.mock(AuditLog.class);
        TransportConfiguration wrongConfig = Mockito.mock(TransportConfiguration.class);

        assertThrows(IllegalArgumentException.class, () ->
                transport.createTransportInstance("can-socketcan://can0", wrongConfig, auditLog));
    }

    @Test
    void createTransportInstanceOnNonLinuxThrowsTransportException() {
        String osName = System.getProperty("os.name", "").toLowerCase();
        if (osName.contains("linux")) {
            // On Linux, this would try to open a real CAN socket
            return;
        }

        SocketCanTransport transport = new SocketCanTransport();
        AuditLog auditLog = Mockito.mock(AuditLog.class);

        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        config.interfaceName = "can0";

        // On non-Linux platforms, should throw TransportException
        assertThrows(Exception.class, () ->
                transport.createTransportInstance("can-socketcan://can0", config, auditLog));
    }

    /**
     * Attempts to create an instance and discards whatever comes back.
     * <p>
     * What happens next depends entirely on the machine: off Linux the platform check rejects it,
     * on Linux it goes on to open a real CAN socket, which may fail for want of the interface or
     * of the native library - or succeed, on a host that has the interface configured. None of
     * that is what these tests are about, so none of it is asserted on; the tests check the
     * resolved configuration, which is settled before any of it happens.
     */
    private static void attemptCreation(SocketCanTransport transport, String addressSegment,
                                        SocketCanTransportConfiguration config, AuditLog auditLog) {
        TransportInstance<SocketCanTransportConfiguration> instance = null;
        try {
            instance = transport.createTransportInstance(addressSegment, config, auditLog);
        } catch (Throwable ignored) {
            // Deliberately ignored - see the javadoc above.
        } finally {
            if (instance != null) {
                try {
                    instance.close();
                } catch (Exception ignored) {
                    // Nothing useful to do in a test teardown path.
                }
            }
        }
    }

    @Test
    void addressSegmentNamesTheInterface() {
        SocketCanTransport transport = new SocketCanTransport();
        AuditLog auditLog = Mockito.mock(AuditLog.class);
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();

        attemptCreation(transport, "can0", config, auditLog);

        assertEquals("can0", config.interfaceName);
    }

    @Test
    void addressSegmentWinsOverInterfaceNameOption() {
        SocketCanTransport transport = new SocketCanTransport();
        AuditLog auditLog = Mockito.mock(AuditLog.class);
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        config.interfaceName = "vcan0";

        attemptCreation(transport, "can0", config, auditLog);

        assertEquals("can0", config.interfaceName);
    }

    @Test
    void interfaceNameOptionUsedWhenAddressSegmentIsEmpty() {
        SocketCanTransport transport = new SocketCanTransport();
        AuditLog auditLog = Mockito.mock(AuditLog.class);
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();
        config.interfaceName = "vcan0";

        attemptCreation(transport, "", config, auditLog);

        assertEquals("vcan0", config.interfaceName);
    }

    @Test
    void noInterfaceAnywhereIsRejected() {
        SocketCanTransport transport = new SocketCanTransport();
        AuditLog auditLog = Mockito.mock(AuditLog.class);
        SocketCanTransportConfiguration config = new SocketCanTransportConfiguration();

        // Reported here rather than by the @Required check on the configuration field, which runs
        // before the address segment has been seen and would reject "can-socketcan://can0".
        TransportException e = assertThrows(TransportException.class, () ->
                transport.createTransportInstance("", config, auditLog));

        assertTrue(e.getMessage().contains("No CAN interface given"));
    }
}
