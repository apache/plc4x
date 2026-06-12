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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.modbus.readwrite.ModbusTcpADU;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Method;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class ModbusTcpMessageCodecTest {

    @SuppressWarnings("unchecked")
    @Test
    void testConstruction() {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusTcpADU> handler = mock(Consumer.class);

        ModbusTcpMessageCodec codec = new ModbusTcpMessageCodec(transportInstance, handler);
        assertNotNull(codec);
    }

    @SuppressWarnings("unchecked")
    @Test
    void testMinimumHeaderSize() throws Exception {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusTcpADU> handler = mock(Consumer.class);
        ModbusTcpMessageCodec codec = new ModbusTcpMessageCodec(transportInstance, handler);

        Method method = codec.getClass().getDeclaredMethod("getMinimumHeaderSize");
        method.setAccessible(true);
        assertEquals(6, method.invoke(codec));
    }

    @SuppressWarnings("unchecked")
    @Test
    void testCalculateTotalMessageSize() throws Exception {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusTcpADU> handler = mock(Consumer.class);
        ModbusTcpMessageCodec codec = new ModbusTcpMessageCodec(transportInstance, handler);

        Method method = codec.getClass().getDeclaredMethod("calculateTotalMessageSize", byte[].class, int.class);
        method.setAccessible(true);

        // Header: txId(2) + protocolId(2) + length(2) = 6 bytes
        // Length field says 5 bytes follow (unitId + PDU)
        byte[] header = new byte[]{0x00, 0x01, 0x00, 0x00, 0x00, 0x05};
        assertEquals(11, method.invoke(codec, header, 11)); // 6 + 5
    }
}
