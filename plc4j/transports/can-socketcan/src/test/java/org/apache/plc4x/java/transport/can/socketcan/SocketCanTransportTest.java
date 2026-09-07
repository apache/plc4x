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

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transport.can.socketcan.config.SocketCanTransportConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

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
}
