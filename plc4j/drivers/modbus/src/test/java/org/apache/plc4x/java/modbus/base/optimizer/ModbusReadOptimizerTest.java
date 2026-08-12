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

import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagCoil;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.opentest4j.AssertionFailedError;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
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
