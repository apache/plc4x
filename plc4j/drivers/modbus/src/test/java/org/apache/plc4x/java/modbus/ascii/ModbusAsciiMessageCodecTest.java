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
package org.apache.plc4x.java.modbus.ascii;

import org.apache.plc4x.java.modbus.readwrite.ModbusAsciiADU;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Method;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class ModbusAsciiMessageCodecTest {

    @SuppressWarnings("unchecked")
    @Test
    void testConstruction() {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusAsciiADU> handler = mock(Consumer.class);

        ModbusAsciiMessageCodec codec = new ModbusAsciiMessageCodec(transportInstance, handler);
        assertNotNull(codec);
    }

    @SuppressWarnings("unchecked")
    @Test
    void testMinimumHeaderSize() throws Exception {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusAsciiADU> handler = mock(Consumer.class);
        ModbusAsciiMessageCodec codec = new ModbusAsciiMessageCodec(transportInstance, handler);

        Method method = codec.getClass().getDeclaredMethod("getMinimumHeaderSize");
        method.setAccessible(true);
        assertEquals(9, method.invoke(codec));
    }

    @SuppressWarnings("unchecked")
    @Test
    void testCalculateTotalMessageSize_findsFrameEnd() throws Exception {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusAsciiADU> handler = mock(Consumer.class);
        ModbusAsciiMessageCodec codec = new ModbusAsciiMessageCodec(transportInstance, handler);

        Method method = codec.getClass().getDeclaredMethod("calculateTotalMessageSize", byte[].class, int.class);
        method.setAccessible(true);

        // Simulate a complete ASCII frame: ':01030000000AFC\r\n' (17 bytes)
        byte[] frameData = ":01030000000AFC\r\n".getBytes(java.nio.charset.StandardCharsets.US_ASCII);
        when(transportInstance.peekReadableBytes(frameData.length)).thenReturn(frameData);
        assertEquals(frameData.length, method.invoke(codec, new byte[9], frameData.length));

        // Simulate incomplete frame (no CR+LF yet) — should return -1
        byte[] incompleteData = ":01030000000AFC".getBytes(java.nio.charset.StandardCharsets.US_ASCII);
        when(transportInstance.peekReadableBytes(incompleteData.length)).thenReturn(incompleteData);
        assertEquals(-1, method.invoke(codec, new byte[9], incompleteData.length));
    }
}
