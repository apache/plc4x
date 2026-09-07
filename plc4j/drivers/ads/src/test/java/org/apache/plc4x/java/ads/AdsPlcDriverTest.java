/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.ads;

import org.apache.plc4x.java.ads.configuration.AdsConfiguration;
import org.apache.plc4x.java.ads.configuration.AdsTcpTransportConfiguration;
import org.apache.plc4x.java.ads.discovery.AdsPlcDiscoverer;
import org.apache.plc4x.java.ads.readwrite.Constants;
import org.apache.plc4x.java.ads.tag.AdsTag;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.transports.api.Transport;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.lang.reflect.Method;
import java.util.Set;
import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.*;

class AdsPlcDriverTest {

    private final AdsPlcDriver driver = new AdsPlcDriver();

    @Test
    void protocolMetadata() {
        assertEquals("ads", driver.getProtocolCode());
        assertEquals("Beckhoff TwinCat ADS", driver.getProtocolName());
        assertEquals("tcp", driver.getDefaultTransportCode().orElseThrow());
        assertEquals(java.util.Collections.singletonList("tcp"), driver.getSupportedTransportCodes());
    }

    @Test
    void defaultPorts() {
        assertEquals(Set.of(Constants.ADSTCPDEFAULTPORT), driver.defaultPorts("tcp"));
        assertEquals(Set.of(Constants.ADSTCPDEFAULTPORT), driver.defaultPorts("TCP"));
        assertTrue(driver.defaultPorts("udp").isEmpty());
    }

    @Test
    void capabilities() throws Exception {
        // Use reflection to access the protected can* methods.
        for (String name : new String[]{"canDiscover", "canPing", "canRead", "canWrite", "canSubscribe", "canBrowse"}) {
            var m = AdsPlcDriver.class.getDeclaredMethod(name);
            m.setAccessible(true);
            assertEquals(Boolean.TRUE, m.invoke(driver), name + " should return true");
        }
    }

    @Test
    void configurationClasses() throws Exception {
        var getConfigurationClass = AdsPlcDriver.class.getDeclaredMethod("getConfigurationClass");
        getConfigurationClass.setAccessible(true);
        assertEquals(AdsConfiguration.class, getConfigurationClass.invoke(driver));

        var getTransportConfigurationClass = AdsPlcDriver.class.getDeclaredMethod("getTransportConfigurationClass", Transport.class);
        getTransportConfigurationClass.setAccessible(true);

        Transport<?> tcp = Mockito.mock(Transport.class);
        Mockito.when(tcp.getTransportCode()).thenReturn("tcp");
        assertEquals(AdsTcpTransportConfiguration.class, getTransportConfigurationClass.invoke(driver, tcp));
    }

    @Test
    void discoveryRequestBuilderNotNull() {
        assertNotNull(driver.discoveryRequestBuilder());
    }

    @Test
    void prepareTagDelegatesToHandler() {
        AdsTag tag = driver.prepareTag("Main.value");
        assertNotNull(tag);
    }

    @Test
    void getTransportConfigurationClassFallsThroughForUnknownTransport() throws Exception {
        Method m = AdsPlcDriver.class.getDeclaredMethod("getTransportConfigurationClass", Transport.class);
        m.setAccessible(true);
        Transport<?> udp = Mockito.mock(Transport.class);
        Mockito.when(udp.getTransportCode()).thenReturn("udp");
        // Falls through to super.getTransportConfigurationClass — should not return AdsTcpTransportConfiguration.class.
        Object result = m.invoke(driver, udp);
        assertNotEquals(AdsTcpTransportConfiguration.class, result);
    }

    @Test
    void getConnectionBuildsAdsTcpConnection() throws Exception {
        Method m = AdsPlcDriver.class.getDeclaredMethod(
            "getConnection", org.apache.plc4x.java.spi.config.Configuration.class, TransportInstance.class, AuditLog.class);
        m.setAccessible(true);
        ConnectionBase<?> connection = (ConnectionBase<?>) m.invoke(driver,
            new AdsConfiguration(), Mockito.mock(TransportInstance.class), Mockito.mock(AuditLog.class));
        assertNotNull(connection);
    }

    @Test
    void discoverDelegatesToDiscoverer() throws Exception {
        CompletableFuture<PlcDiscoveryResponse> future = new CompletableFuture<>();
        try (var mocked = Mockito.mockConstruction(AdsPlcDiscoverer.class,
            (mock, ctx) -> Mockito.when(mock.discover(Mockito.any())).thenReturn(future))) {
            AdsPlcDriver d = new AdsPlcDriver();
            PlcDiscoveryRequest request = Mockito.mock(PlcDiscoveryRequest.class);
            assertSame(future, d.discover(request));
            Mockito.verify(mocked.constructed().get(0)).discover(request);
        }
    }
}
