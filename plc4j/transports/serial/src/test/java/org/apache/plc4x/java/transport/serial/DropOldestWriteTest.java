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
package org.apache.plc4x.java.transport.serial;

import org.apache.plc4x.java.spi.transports.api.RingBuffer;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class DropOldestWriteTest {

    @Test
    void writesWithoutDropWhenSpaceSuffices() {
        RingBuffer buffer = new RingBuffer(8);
        int dropped = SerialTransportInstance.writeDroppingOldest(buffer, new byte[]{1, 2, 3}, 0, 3);
        assertEquals(0, dropped);
        assertArrayEquals(new byte[]{1, 2, 3}, buffer.read(3));
    }

    @Test
    void dropsOldestOnOverflow() {
        RingBuffer buffer = new RingBuffer(4);
        SerialTransportInstance.writeDroppingOldest(buffer, new byte[]{'a', 'b'}, 0, 2);
        int dropped = SerialTransportInstance.writeDroppingOldest(buffer, new byte[]{'c', 'd', 'e', 'f'}, 0, 4);
        assertEquals(2, dropped, "the two oldest bytes must be dropped");
        assertArrayEquals(new byte[]{'c', 'd', 'e', 'f'}, buffer.read(4), "newest bytes survive");
    }

    @Test
    void chunkLargerThanCapacityKeepsTail() {
        RingBuffer buffer = new RingBuffer(4);
        byte[] big = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'};
        int dropped = SerialTransportInstance.writeDroppingOldest(buffer, big, 0, big.length);
        assertEquals(6, dropped);
        assertArrayEquals(new byte[]{'6', '7', '8', '9'}, buffer.read(4), "only the tail fits");
    }
}
