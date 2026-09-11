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
package org.apache.plc4x.java.modbus.base.optimizer;

import org.apache.plc4x.java.modbus.base.ModbusRegisterCodec;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagCoil;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagExtendedRegister;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.opentest4j.AssertionFailedError;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.*;

class ModbusReadOptimizerTest {

    private static Stream<Arguments> coilInputData() {
        return Stream.of(
            // Simple one tag coil test
            Arguments.of(new ModbusTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(1, optimizedReads.size());
                    ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagCoil.class, mergedTag);
                    assertEquals(0, mergedTag.getAddress());
                    assertEquals(1, mergedTag.getNumberOfElements());
                }),

            // In this test, the two adjacent coils will be joined together to one array.
            Arguments.of(new ModbusTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                    new ModbusTagCoil(1, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(1, optimizedReads.size());
                    ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagCoil.class, mergedTag);
                    assertEquals(0, mergedTag.getAddress());
                    assertEquals(2, mergedTag.getNumberOfElements());
                }),

            // In this test, the two coils with a larger gap will be joined together to one array.
            Arguments.of(new ModbusTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                    new ModbusTagCoil(100, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(1, optimizedReads.size());
                    ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagCoil.class, mergedTag);
                    assertEquals(0, mergedTag.getAddress());
                    assertEquals(101, mergedTag.getNumberOfElements());
                }),

            // In this test, the two coils have a too large gap to be read in one block.
            Arguments.of(new ModbusTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                    new ModbusTagCoil(2100, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(2, optimizedReads.size());
                    ModbusTag first = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagCoil.class, first);
                    assertEquals(0, first.getAddress());
                    assertEquals(1, first.getNumberOfElements());
                    ModbusTag second = optimizedReads.get(1).mergedTag;
                    assertInstanceOf(ModbusTagCoil.class, second);
                    assertEquals(2100, second.getAddress());
                    assertEquals(1, second.getNumberOfElements());
                })
        );
    }

    @ParameterizedTest
    @MethodSource("coilInputData")
    void coilTests(ModbusTag[] tags, CheckResult check) {
        processReadRequest(tags, check);
    }

    private static Stream<Arguments> holdingRegisterInputData() {
        return Stream.of(
            // Simple one tag register test
            Arguments.of(new ModbusTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(1, optimizedReads.size());
                    ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagHoldingRegister.class, mergedTag);
                    assertEquals(0, mergedTag.getAddress());
                    assertEquals(1, mergedTag.getNumberOfElements());
                }),

            // Two adjacent registers joined together.
            Arguments.of(new ModbusTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                    new ModbusTagHoldingRegister(1, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(1, optimizedReads.size());
                    ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagHoldingRegister.class, mergedTag);
                    assertEquals(0, mergedTag.getAddress());
                    assertEquals(2, mergedTag.getNumberOfElements());
                }),

            // Two registers with a larger gap joined together.
            Arguments.of(new ModbusTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                    new ModbusTagHoldingRegister(100, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(1, optimizedReads.size());
                    ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagHoldingRegister.class, mergedTag);
                    assertEquals(0, mergedTag.getAddress());
                    assertEquals(101, mergedTag.getNumberOfElements());
                }),

            // Two registers too far apart - split into two requests.
            Arguments.of(new ModbusTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                    new ModbusTagHoldingRegister(2100, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) optimizedReads -> {
                    assertEquals(2, optimizedReads.size());
                    ModbusTag first = optimizedReads.getFirst().mergedTag;
                    assertInstanceOf(ModbusTagHoldingRegister.class, first);
                    assertEquals(0, first.getAddress());
                    assertEquals(1, first.getNumberOfElements());
                    ModbusTag second = optimizedReads.get(1).mergedTag;
                    assertInstanceOf(ModbusTagHoldingRegister.class, second);
                    assertEquals(2100, second.getAddress());
                    assertEquals(1, second.getNumberOfElements());
                })
        );
    }

    @ParameterizedTest
    @MethodSource("holdingRegisterInputData")
    void holdingRegisterTests(ModbusTag[] tags, CheckResult check) {
        processReadRequest(tags, check);
    }

    /**
     * Tags addressing different units may never be merged into a single block read, and the
     * unit-id of a group has to be carried over to the merged tag - otherwise the connection
     * would silently fall back to its default unit-id (see GitHub issue #2686).
     */
    @Test
    void registersWithDifferentUnitIdsAreNotMerged() {
        processReadRequest(new ModbusTag[]{
                new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.singletonMap("unit-id", "2")),
                new ModbusTagHoldingRegister(1, 1, ModbusDataType.INT, Collections.singletonMap("unit-id", "3"))
            },
            optimizedReads -> {
                assertEquals(2, optimizedReads.size());
                // The addresses would have been adjacent, so without the unit-id they'd be one block.
                ModbusTag first = optimizedReads.getFirst().mergedTag;
                assertEquals((short) 2, first.getUnitId());
                assertEquals(0, first.getAddress());
                assertEquals(1, first.getNumberOfElements());
                ModbusTag second = optimizedReads.get(1).mergedTag;
                assertEquals((short) 3, second.getUnitId());
                assertEquals(1, second.getAddress());
                assertEquals(1, second.getNumberOfElements());
            });
    }

    @Test
    void coilsWithDifferentUnitIdsAreNotMerged() {
        processReadRequest(new ModbusTag[]{
                new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.singletonMap("unit-id", "2")),
                new ModbusTagCoil(1, 1, ModbusDataType.BOOL, Collections.singletonMap("unit-id", "3"))
            },
            optimizedReads -> {
                assertEquals(2, optimizedReads.size());
                assertEquals((short) 2, optimizedReads.getFirst().mergedTag.getUnitId());
                assertEquals((short) 3, optimizedReads.get(1).mergedTag.getUnitId());
            });
    }

    /**
     * Tags of the same unit are still merged, and the merged tag keeps the unit-id.
     */
    @Test
    void registersWithSameUnitIdAreMergedKeepingTheUnitId() {
        processReadRequest(new ModbusTag[]{
                new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.singletonMap("unit-id", "7")),
                new ModbusTagHoldingRegister(1, 1, ModbusDataType.INT, Collections.singletonMap("unit-id", "7"))
            },
            optimizedReads -> {
                assertEquals(1, optimizedReads.size());
                ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                assertEquals((short) 7, mergedTag.getUnitId());
                assertEquals(0, mergedTag.getAddress());
                assertEquals(2, mergedTag.getNumberOfElements());
            });
    }

    /**
     * Tags with an explicit unit-id must not be mixed with tags that use the connection default.
     */
    @Test
    void tagsWithoutUnitIdAreNotMergedWithTagsHavingOne() {
        processReadRequest(new ModbusTag[]{
                new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                new ModbusTagHoldingRegister(1, 1, ModbusDataType.INT, Collections.singletonMap("unit-id", "2"))
            },
            optimizedReads -> {
                assertEquals(2, optimizedReads.size());
                // Tags without a unit-id are grouped first and keep using the connection default.
                assertNull(optimizedReads.getFirst().mergedTag.getUnitId());
                assertEquals(0, optimizedReads.getFirst().mergedTag.getAddress());
                assertEquals((short) 2, optimizedReads.get(1).mergedTag.getUnitId());
                assertEquals(1, optimizedReads.get(1).mergedTag.getAddress());
            });
    }

    /**
     * A coil array (BOOL[n]) has to yield all n values, not just the first one - see GH-2060.
     * The block read already covers the whole array, so the data is on the wire either way.
     */
    @Test
    void coilArrayReturnsAllElements() {
        // 8 coils starting at address 0, response bits (LSB first): 1,0,0,0,1,1,0,1
        Map<String, PlcResponseItem<PlcValue>> response =
            splitSingleRead(new ModbusTagCoil(0, 8, ModbusDataType.BOOL, Collections.emptyMap()),
                new byte[]{(byte) 0b10110001});

        PlcValue value = response.get("tag0").getValue();
        assertInstanceOf(PlcList.class, value);
        assertEquals(8, value.getLength());
        boolean[] expected = {true, false, false, false, true, true, false, true};
        for (int i = 0; i < expected.length; i++) {
            assertEquals(expected[i], value.getIndex(i).getBoolean(), "bit " + i);
        }
    }

    /**
     * An array sitting at an offset inside the block must be read from the right bit position and
     * keep reading into the following byte when it spans a byte boundary.
     */
    @Test
    void coilArrayCrossingAByteBoundary() {
        // A scalar coil at address 0 makes the block start there, so the array at address 6
        // really does start at bit 6 of the first byte and runs into the second one.
        LinkedHashMap<String, ModbusTag> tags = new LinkedHashMap<>();
        tags.put("scalar", new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()));
        tags.put("array", new ModbusTagCoil(6, 4, ModbusDataType.BOOL, Collections.emptyMap()));

        // bits 6,7 of byte 0 = 0,1 ; bits 0,1 of byte 1 = 0,1
        Map<String, PlcResponseItem<PlcValue>> response =
            splitRead(tags, new byte[]{(byte) 0b10000000, (byte) 0b00000010});

        PlcValue value = response.get("array").getValue();
        assertEquals(4, value.getLength());
        assertFalse(value.getIndex(0).getBoolean());
        assertTrue(value.getIndex(1).getBoolean());
        assertFalse(value.getIndex(2).getBoolean());
        assertTrue(value.getIndex(3).getBoolean());
    }

    /**
     * A scalar coil must keep returning a plain PlcBOOL rather than a single-element list.
     */
    @Test
    void singleCoilStillReturnsAScalar() {
        Map<String, PlcResponseItem<PlcValue>> response =
            splitSingleRead(new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                new byte[]{(byte) 0b00000001});

        PlcValue value = response.get("tag0").getValue();
        assertInstanceOf(PlcBOOL.class, value);
        assertTrue(value.getBoolean());
    }

    /**
     * A coil holds a single bit, so anything other than BOOL has to be reported as unsupported
     * rather than silently answered with the first bit.
     */
    @Test
    void nonBoolCoilIsReportedAsUnsupported() {
        Map<String, PlcResponseItem<PlcValue>> response =
            splitSingleRead(new ModbusTagCoil(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                new byte[]{(byte) 0b00000001, (byte) 0b00000000});

        assertEquals(PlcResponseCode.UNSUPPORTED, response.get("tag0").getResponseCode());
        assertNull(response.get("tag0").getValue());
    }

    /**
     * A device answering with fewer coils than requested must not blow up the whole response.
     */
    @Test
    void shortCoilResponseYieldsAnError() {
        Map<String, PlcResponseItem<PlcValue>> response =
            splitSingleRead(new ModbusTagCoil(0, 16, ModbusDataType.BOOL, Collections.emptyMap()),
                new byte[]{(byte) 0b00000001});

        assertEquals(PlcResponseCode.INTERNAL_ERROR, response.get("tag0").getResponseCode());
    }

    /** Selects the byte-swapping big endian order, which a tag may ask for by itself. */
    private static final Map<String, String> BIG_ENDIAN_BYTE_SWAP =
        Collections.singletonMap("byte-order", ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP.name());

    /**
     * A CHAR occupies a whole register but carries one byte, so a slice cut to its payload length
     * cannot be swapped in pairs - the byte would stay where it is, while the unoptimized read,
     * which is handed the whole register, swaps it. Merging a tag into a block must not change the
     * value it decodes to.
     */
    @Test
    void oddLengthScalarIsSwappedTheSameWayAnUnoptimizedReadSwapsIt() throws BufferException {
        ModbusTag tag = new ModbusTagHoldingRegister(1, 1, ModbusDataType.CHAR, BIG_ENDIAN_BYTE_SWAP);
        byte[] registers = {'A', 'B'};

        Map<String, PlcResponseItem<PlcValue>> response = splitSingleRead(tag, registers);

        assertEquals(PlcResponseCode.OK, response.get("tag0").getResponseCode());
        assertEquals("B", unoptimizedRead(tag, registers).getString());
        assertEquals(unoptimizedRead(tag, registers).getString(), response.get("tag0").getValue().getString());
    }

    /**
     * The same for a string of odd declared length: three characters sit in two registers, and the
     * swap has to see both of them.
     */
    @Test
    void oddLengthStringIsSwappedTheSameWayAnUnoptimizedReadSwapsIt() throws BufferException {
        ModbusTag text = new ModbusTagHoldingRegister(1, 1, 3, ModbusDataType.STRING, BIG_ENDIAN_BYTE_SWAP);
        // A second tag behind it, so the block really does cover both of the string's registers
        // rather than the test handing over registers nobody asked for.
        LinkedHashMap<String, ModbusTag> tags = new LinkedHashMap<>();
        tags.put("text", text);
        tags.put("guard", new ModbusTagHoldingRegister(3, 1, ModbusDataType.INT, BIG_ENDIAN_BYTE_SWAP));
        byte[] blockData = {'A', 'B', 'C', 'D', 0x12, 0x34};

        Map<String, PlcResponseItem<PlcValue>> response = splitRead(tags, blockData);

        assertEquals(PlcResponseCode.OK, response.get("text").getResponseCode());
        byte[] stringRegisters = {'A', 'B', 'C', 'D'};
        assertEquals("BAD", unoptimizedRead(text, stringRegisters).getString());
        assertEquals(unoptimizedRead(text, stringRegisters).getString(), response.get("text").getValue().getString());
    }

    /**
     * Decodes the given registers for a single tag the way a read that was never merged into a
     * block does (see the connections' {@code toPlcValue}): the whole registers the device
     * answered with, swapped as one, then parsed.
     */
    private static PlcValue unoptimizedRead(ModbusTag tag, byte[] registers) throws BufferException {
        ModbusByteOrder byteOrder = tag.getByteOrder();
        byte[] data = registers;
        if (byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP || byteOrder == ModbusByteOrder.LITTLE_ENDIAN_BYTE_SWAP) {
            data = byteSwap(data);
        }
        boolean bigEndian = (byteOrder == ModbusByteOrder.BIG_ENDIAN || byteOrder == ModbusByteOrder.BIG_ENDIAN_BYTE_SWAP);
        return ModbusRegisterCodec.parse(new ReadBufferByteBased(data), tag.getDataType(),
            tag.getNumberOfElements(), bigEndian, tag.getStringLength());
    }

    private static byte[] byteSwap(byte[] in) {
        byte[] out = new byte[in.length];
        for (int i = 0; i < out.length - 1; i += 2) {
            out[i] = in[i + 1];
            out[i + 1] = in[i];
        }
        if (in.length % 2 != 0) {
            out[in.length - 1] = in[in.length - 1];
        }
        return out;
    }

    /** The largest PDU the Modbus specification allows, request or response. */
    private static final int MAX_MODBUS_PDU_SIZE = 253;

    /** The number of registers one extended-register file holds. */
    private static final int EXTENDED_REGISTER_FILE_LENGTH = 10000;

    /**
     * An extended-register block is answered over FC 0x14 (Read File Record), whose per-item
     * framing leaves room for fewer registers than FC 0x03 does. A 125-register block is one byte
     * too wide for it, so two tags spanning that many registers must be read separately - even
     * though the plain register areas happily merge them.
     */
    @Test
    void extendedRegistersAreNotMergedWiderThanAReadFileRecordResponseFits() {
        processReadRequest(new ModbusTag[]{
                new ModbusTagExtendedRegister(1, 1, ModbusDataType.INT, Collections.emptyMap()),
                new ModbusTagExtendedRegister(125, 1, ModbusDataType.INT, Collections.emptyMap())
            },
            optimizedReads -> {
                assertEquals(2, optimizedReads.size());
                ModbusTag first = optimizedReads.getFirst().mergedTag;
                assertInstanceOf(ModbusTagExtendedRegister.class, first);
                assertEquals(1, first.getAddress());
                assertEquals(1, first.getNumberOfElements());
                ModbusTag second = optimizedReads.get(1).mergedTag;
                assertInstanceOf(ModbusTagExtendedRegister.class, second);
                assertEquals(125, second.getAddress());
                assertEquals(1, second.getNumberOfElements());
            });

        // The same two addresses in a holding register: the narrower ceiling is the extended
        // registers' own, not a tightening of the limit every area is cut with.
        processReadRequest(new ModbusTag[]{
                new ModbusTagHoldingRegister(1, 1, ModbusDataType.INT, Collections.emptyMap()),
                new ModbusTagHoldingRegister(125, 1, ModbusDataType.INT, Collections.emptyMap())
            },
            optimizedReads -> {
                assertEquals(1, optimizedReads.size());
                assertEquals(125, optimizedReads.getFirst().mergedTag.getNumberOfElements());
            });
    }

    /**
     * Extended registers that do fit stay merged - the ceiling must not cost a round trip it
     * doesn't have to.
     */
    @Test
    void extendedRegistersThatFitAreStillMerged() {
        processReadRequest(new ModbusTag[]{
                new ModbusTagExtendedRegister(1, 1, ModbusDataType.INT, Collections.emptyMap()),
                new ModbusTagExtendedRegister(123, 1, ModbusDataType.INT, Collections.emptyMap())
            },
            optimizedReads -> {
                assertEquals(1, optimizedReads.size());
                ModbusTag mergedTag = optimizedReads.getFirst().mergedTag;
                assertInstanceOf(ModbusTagExtendedRegister.class, mergedTag);
                assertEquals(1, mergedTag.getAddress());
                assertEquals(123, mergedTag.getNumberOfElements());
            });
    }

    /**
     * Whatever the optimizer emits for extended registers has to be answerable: no block may ask
     * for more than an FC 0x14 response can carry back, including one that crosses a file boundary
     * and is therefore answered as two items.
     */
    @Test
    void everyExtendedRegisterBlockFitsTheModbusPduLimit() {
        // Densely addressed tags, so every block is cut as wide as the optimizer will allow, both
        // inside a single file and across the boundary between two of them.
        LinkedHashMap<String, ModbusTag> tags = new LinkedHashMap<>();
        for (int address = 1; address < 250; address++) {
            tags.put("low" + address, new ModbusTagExtendedRegister(address, 1, ModbusDataType.INT, Collections.emptyMap()));
        }
        for (int address = 9900; address < 10200; address++) {
            tags.put("boundary" + address, new ModbusTagExtendedRegister(address, 1, ModbusDataType.INT, Collections.emptyMap()));
        }

        ModbusReadOptimizer optimizer = new ModbusReadOptimizer(2000, 125, ModbusByteOrder.BIG_ENDIAN);
        List<ModbusReadOptimizer.OptimizedRead> optimizedReads = optimizer.optimizeReads(tags);

        assertFalse(optimizedReads.isEmpty());
        boolean sawBlockCrossingAFileBoundary = false;
        for (ModbusReadOptimizer.OptimizedRead optimizedRead : optimizedReads) {
            ModbusTag block = optimizedRead.mergedTag;
            sawBlockCrossingAFileBoundary |= itemsOf(block) > 1;
            int pduSize = readFileRecordResponseSize(block);
            assertTrue(pduSize <= MAX_MODBUS_PDU_SIZE,
                "the block at " + block.getAddress() + " covering " + block.getNumberOfElements()
                    + " registers is answered with a " + pduSize + " byte PDU");
        }
        assertTrue(sawBlockCrossingAFileBoundary, "the two-item shape wasn't exercised");
    }

    /**
     * The size of the FC 0x14 response answering the given extended-register block: function code
     * (1) + response data length (1) + per item [data length (1) + reference type (1) + data].
     */
    private static int readFileRecordResponseSize(ModbusTag block) {
        return 1 + 1 + itemsOf(block) * 2 + block.getNumberOfElements() * 2;
    }

    /**
     * The number of items the response is split into: one per file the block reaches into, the
     * way the connections group an extended-register read.
     */
    private static int itemsOf(ModbusTag block) {
        int firstFile = block.getAddress() / EXTENDED_REGISTER_FILE_LENGTH;
        int lastFile = (block.getAddress() + block.getNumberOfElements() - 1) / EXTENDED_REGISTER_FILE_LENGTH;
        return lastFile - firstFile + 1;
    }

    /**
     * Runs a single tag through the optimizer and splits the given raw response for it.
     */
    private Map<String, PlcResponseItem<PlcValue>> splitSingleRead(ModbusTag tag, byte[] blockData) {
        LinkedHashMap<String, ModbusTag> tagMap = new LinkedHashMap<>();
        tagMap.put("tag0", tag);
        return splitRead(tagMap, blockData);
    }

    /**
     * Runs the given tags through the optimizer and splits the given raw response for the
     * single block read they are expected to produce.
     */
    private Map<String, PlcResponseItem<PlcValue>> splitRead(LinkedHashMap<String, ModbusTag> tags, byte[] blockData) {
        ModbusReadOptimizer optimizer = new ModbusReadOptimizer(2000, 125, ModbusByteOrder.BIG_ENDIAN);
        List<ModbusReadOptimizer.OptimizedRead> optimizedReads = optimizer.optimizeReads(tags);
        assertEquals(1, optimizedReads.size());
        return optimizer.splitResponse(optimizedReads.getFirst(), PlcResponseCode.OK, blockData);
    }

    void processReadRequest(ModbusTag[] tags, CheckResult check) {
        ModbusReadOptimizer optimizer = new ModbusReadOptimizer(2000, 125, ModbusByteOrder.BIG_ENDIAN);
        LinkedHashMap<String, ModbusTag> tagMap = new LinkedHashMap<>();
        int i = 0;
        for (ModbusTag tag : tags) {
            tagMap.put("tag" + i++, tag);
        }
        List<ModbusReadOptimizer.OptimizedRead> optimizedReads = optimizer.optimizeReads(tagMap);
        check.isValid(optimizedReads);
    }

    @FunctionalInterface
    protected interface CheckResult {
        void isValid(List<ModbusReadOptimizer.OptimizedRead> optimizedReads) throws AssertionFailedError;
    }

}
