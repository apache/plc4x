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

package org.apache.plc4x.java.transports.api;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.when;

class OsgiTransportManagerTest {

    private OsgiTransportManager manager;
    private Transport<?> mockTransport;

    @BeforeEach
    void setUp() {
        manager = new OsgiTransportManager();
        mockTransport = Mockito.mock(Transport.class);
        when(mockTransport.getTransportCode()).thenReturn("tcp");
        when(mockTransport.getTransportName()).thenReturn("TCP Transport");
    }

    @Test
    void testGetTransport_whenTransportNotRegistered() {
        Optional<Transport> result = manager.getTransport("tcp");
        assertFalse(result.isPresent());
    }

    @Test
    void testBindTransport_registersTransport() {
        manager.bindTransport(mockTransport);

        Optional<Transport> result = manager.getTransport("tcp");
        assertTrue(result.isPresent());
        assertEquals(mockTransport, result.get());
    }

    @Test
    void testBindTransport_duplicateTransportIgnored() {
        Transport<?> mockTransport2 = Mockito.mock(Transport.class);
        when(mockTransport2.getTransportCode()).thenReturn("tcp");
        when(mockTransport2.getTransportName()).thenReturn("TCP Transport 2");

        manager.bindTransport(mockTransport);
        manager.bindTransport(mockTransport2);

        Optional<Transport> result = manager.getTransport("tcp");
        assertTrue(result.isPresent());
        assertEquals(mockTransport, result.get());
    }

    @Test
    void testUnbindTransport_removesTransport() {
        manager.bindTransport(mockTransport);

        Optional<Transport> result = manager.getTransport("tcp");
        assertTrue(result.isPresent());

        manager.unbindTransport(mockTransport);

        result = manager.getTransport("tcp");
        assertFalse(result.isPresent());
    }

    @Test
    void testBindMultipleTransports() {
        Transport<?> mockTransport2 = Mockito.mock(Transport.class);
        when(mockTransport2.getTransportCode()).thenReturn("udp");
        when(mockTransport2.getTransportName()).thenReturn("UDP Transport");

        manager.bindTransport(mockTransport);
        manager.bindTransport(mockTransport2);

        Optional<Transport> tcpResult = manager.getTransport("tcp");
        Optional<Transport> udpResult = manager.getTransport("udp");

        assertTrue(tcpResult.isPresent());
        assertTrue(udpResult.isPresent());
        assertEquals(mockTransport, tcpResult.get());
        assertEquals(mockTransport2, udpResult.get());
    }

    @Test
    void testGetTransport_returnsEmptyForUnknownCode() {
        manager.bindTransport(mockTransport);

        Optional<Transport> result = manager.getTransport("unknown");
        assertFalse(result.isPresent());
    }
}
