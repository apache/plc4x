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

import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.s7.configuration.S7CotpTransportConfiguration;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.lang.reflect.Method;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.*;

class S7DriverTest {

    private final S7Driver driver = new S7Driver();

    @Test
    void protocolMetadata() {
        assertEquals("s7", driver.getProtocolCode());
        assertEquals("Siemens S7 (Basic)", driver.getProtocolName());
        assertEquals("cotp", driver.getDefaultTransportCode().orElseThrow());
        assertEquals(java.util.Collections.singletonList("cotp"), driver.getSupportedTransportCodes());
    }

    @Test
    void defaultPorts() {
        assertEquals(Set.of(S7Driver.ISO_ON_TCP_PORT), driver.defaultPorts("cotp"));
        assertEquals(Set.of(S7Driver.ISO_ON_TCP_PORT), driver.defaultPorts("COTP"));
        assertTrue(driver.defaultPorts("tcp").isEmpty());
    }

    @Test
    void capabilities() throws Exception {
        for (String name : new String[]{"canPing", "canRead", "canWrite"}) {
            Method m = S7Driver.class.getDeclaredMethod(name);
            m.setAccessible(true);
            assertEquals(Boolean.TRUE, m.invoke(driver), name + " should return true");
        }
    }

    @Test
    void configurationClasses() throws Exception {
        Method getConfig = S7Driver.class.getDeclaredMethod("getConfigurationClass");
        getConfig.setAccessible(true);
        assertEquals(S7Configuration.class, getConfig.invoke(driver));

        Method getTransport = S7Driver.class.getDeclaredMethod("getTransportConfigurationClass", Transport.class);
        getTransport.setAccessible(true);

        Transport<?> cotp = Mockito.mock(Transport.class);
        Mockito.when(cotp.getTransportCode()).thenReturn("cotp");
        assertEquals(S7CotpTransportConfiguration.class, getTransport.invoke(driver, cotp));
    }

    @Test
    void getTransportConfigurationClassFallsThroughForUnknownTransport() throws Exception {
        Method m = S7Driver.class.getDeclaredMethod("getTransportConfigurationClass", Transport.class);
        m.setAccessible(true);
        Transport<?> tcp = Mockito.mock(Transport.class);
        Mockito.when(tcp.getTransportCode()).thenReturn("tcp");
        Object result = m.invoke(driver, tcp);
        assertNotEquals(S7CotpTransportConfiguration.class, result);
    }

    @Test
    void prepareTagDelegatesToHandler() {
        S7Tag tag = driver.prepareTag("%DB1.DBW0:INT");
        assertNotNull(tag);
        assertEquals(1, tag.getBlockNumber());
        assertEquals(0, tag.getByteOffset());
    }

    /**
     * Regression: {@code ConnectionBase} calls {@link S7CotpConnection#getMaxConcurrentRequests()}
     * from its own constructor — before our subclass field initializers run. An earlier
     * implementation read from {@code driverContext} which is still null at that point,
     * blowing up with NPE on every connection attempt.
     */
    @Test
    void connectionConstructsWithoutNpe() {
        S7Configuration cfg = new S7Configuration();
        cfg.setMaxAmqCallee(4);
        var transport = Mockito.mock(org.apache.plc4x.java.spi.transports.api.TransportInstance.class);
        var auditLog = Mockito.mock(org.apache.plc4x.java.utils.auditlog.api.AuditLog.class);
        S7CotpConnection conn = new S7CotpConnection(cfg, transport, auditLog);
        assertNotNull(conn);
    }
}
