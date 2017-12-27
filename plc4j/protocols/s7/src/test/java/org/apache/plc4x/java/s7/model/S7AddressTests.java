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