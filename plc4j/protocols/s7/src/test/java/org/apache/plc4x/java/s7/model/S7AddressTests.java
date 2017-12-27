/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package org.apache.plc4x.java.s7.model;

import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7AddressTests {

    @Test
    @Tag("fast")
    void testS7Address() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        S7Address s7Address = new S7Address(memoryArea, (short) 0x100);

        assertTrue(s7Address.getMemoryArea() == MemoryArea.DATA_BLOCKS, "Wrong memory area");
        assertTrue( s7Address.getByteOffset() == 0x100, "Memory area byte offset incorrect");
    }

    @Test
    @Tag("fast")
    void testS7BitAddress() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        S7BitAddress s7Address = new S7BitAddress(memoryArea, (short) 0x50, (byte) 0x4);

        assertTrue(s7Address.getMemoryArea() == MemoryArea.DATA_BLOCKS, "Wrong memory area");
        assertTrue( s7Address.getByteOffset() == 0x50, "Memory area byte offset incorrect");
        assertTrue( s7Address.getBitOffset() == 0x4, "Memory area bit offset incorrect");
    }

    @Test
    @Tag("fast")
    void testS7DatBlockAddress() {
        S7DataBlockAddress s7Address = new S7DataBlockAddress((short) 1, (short) 0x50);

        assertTrue(s7Address.getMemoryArea() == MemoryArea.DATA_BLOCKS, "Wrong memory area");
        assertTrue(s7Address.getDataBlockNumber() == 1, "Memory block number incorrect");
        assertTrue( s7Address.getByteOffset() == 0x50, "Memory area byte offset incorrect");
    }

}