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
package org.apache.plc4x.java.modbus.tcp.discovery;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * Discovery talks to whatever answers the port, so the length field belongs to a stranger. The next
 * thing done with it is to allocate an array that long.
 */
class ModbusPlcDiscovererLengthTest {

    private static byte[] lengthField(int value) {
        return new byte[]{(byte) ((value >> 8) & 0xFF), (byte) (value & 0xFF)};
    }

    @Test
    void aPlausibleLengthIsAccepted() {
        // The shortest real answer: unit id, function code and one byte.
        assertEquals(6 + 3, ModbusPlcDiscoverer.aduLength(lengthField(3)));
        // The longest a Modbus TCP ADU can be.
        assertEquals(260, ModbusPlcDiscoverer.aduLength(lengthField(254)));
    }

    @Test
    void aLengthThatWouldGoNegativeIsTurnedAway() {
        // 0x7FFF plus the header overflows a short, which is how this became a negative array size.
        assertEquals(-1, ModbusPlcDiscoverer.aduLength(lengthField(0x7FFF)));
        // 0xFFFA read signed is -6, which plus the header used to be a length of exactly zero.
        assertEquals(-1, ModbusPlcDiscoverer.aduLength(lengthField(0xFFFA)));
        assertEquals(-1, ModbusPlcDiscoverer.aduLength(lengthField(0xFFFF)));
    }

    @Test
    void aLengthNoModbusAduCouldHaveIsTurnedAway() {
        assertEquals(-1, ModbusPlcDiscoverer.aduLength(lengthField(0)));
        assertEquals(-1, ModbusPlcDiscoverer.aduLength(lengthField(255)));
        assertEquals(-1, ModbusPlcDiscoverer.aduLength(lengthField(4096)));
    }
}
