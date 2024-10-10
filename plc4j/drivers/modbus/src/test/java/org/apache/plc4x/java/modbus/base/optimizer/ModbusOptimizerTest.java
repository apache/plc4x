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

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.modbus.base.context.ModbusContext;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagCoil;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.mockito.Mockito;
import org.opentest4j.AssertionFailedError;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.*;

class ModbusOptimizerTest {

    private static Stream<Arguments> coilInputData() {
        return Stream.of(
            // Simple one tag coil test
            Arguments.of(new PlcTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(1, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagCoil.class, firstTag);
                    ModbusTagCoil coil = (ModbusTagCoil) firstTag;
                    assertEquals(0, coil.getAddress());
                    assertEquals(1, coil.getNumberOfElements());
                    assertEquals(ModbusDataType.BYTE, coil.getDataType());
                }),

            // In this test, the two adjacent coils will be joined together to one array.
            Arguments.of(new PlcTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                    new ModbusTagCoil(1, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(1, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagCoil.class, firstTag);
                    ModbusTagCoil coil = (ModbusTagCoil) firstTag;
                    assertEquals(0, coil.getAddress());
                    assertEquals(2, coil.getNumberOfElements());
                    assertEquals(ModbusDataType.BYTE, coil.getDataType());
                }),

            // In this test, the two coils with a larger gap will be joined together to one array.
            Arguments.of(new PlcTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                    new ModbusTagCoil(100, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(1, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagCoil.class, firstTag);
                    ModbusTagCoil coil = (ModbusTagCoil) firstTag;
                    assertEquals(0, coil.getAddress());
                    assertEquals(101, coil.getNumberOfElements());
                    assertEquals(ModbusDataType.BYTE, coil.getDataType());
                }),

            // In this test, the two coils have a too large gap to be read in one block, therefore the result
            // should be a list with two sub-requests.
            Arguments.of(new PlcTag[]{
                    new ModbusTagCoil(0, 1, ModbusDataType.BOOL, Collections.emptyMap()),
                    new ModbusTagCoil(2100, 1, ModbusDataType.BOOL, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(2, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagCoil.class, firstTag);
                    ModbusTagCoil coil = (ModbusTagCoil) firstTag;
                    assertEquals(0, coil.getAddress());
                    assertEquals(1, coil.getNumberOfElements());
                    assertEquals(ModbusDataType.BYTE, coil.getDataType());
                    PlcReadRequest secondReadRequest = readRequests.get(1);
                    assertEquals(1, secondReadRequest.getNumberOfTags());
                    firstTag = secondReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagCoil.class, firstTag);
                    coil = (ModbusTagCoil) firstTag;
                    assertEquals(2100, coil.getAddress());
                    assertEquals(1, coil.getNumberOfElements());
                    assertEquals(ModbusDataType.BYTE, coil.getDataType());
                })
        );
    }

    @ParameterizedTest
    @MethodSource("coilInputData")
    void coilTests(PlcTag[] tags, CheckResult check) {
        processReadRequest(tags, check);
    }

    private static Stream<Arguments> holdingRegisterInputData() {
        return Stream.of(
            // Simple one tag coil test
            Arguments.of(new PlcTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(1, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagHoldingRegister.class, firstTag);
                    ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) firstTag;
                    assertEquals(0, holdingRegister.getAddress());
                    assertEquals(1, holdingRegister.getNumberOfElements());
                    assertEquals(ModbusDataType.WORD, holdingRegister.getDataType());
                }),

            // In this test, the two adjacent coils will be joined together to one array.
            Arguments.of(new PlcTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                    new ModbusTagHoldingRegister(1, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(1, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagHoldingRegister.class, firstTag);
                    ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) firstTag;
                    assertEquals(0, holdingRegister.getAddress());
                    assertEquals(2, holdingRegister.getNumberOfElements());
                    assertEquals(ModbusDataType.WORD, holdingRegister.getDataType());
                }),

            // In this test, the two coils with a larger gap will be joined together to one array.
            Arguments.of(new PlcTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                    new ModbusTagHoldingRegister(100, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(1, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagHoldingRegister.class, firstTag);
                    ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) firstTag;
                    assertEquals(0, holdingRegister.getAddress());
                    assertEquals(101, holdingRegister.getNumberOfElements());
                    assertEquals(ModbusDataType.WORD, holdingRegister.getDataType());
                }),

            // In this test, the two coils have a too large gap to be read in one block, therefore the result
            // should be a list with two sub-requests.
            Arguments.of(new PlcTag[]{
                    new ModbusTagHoldingRegister(0, 1, ModbusDataType.INT, Collections.emptyMap()),
                    new ModbusTagHoldingRegister(2100, 1, ModbusDataType.INT, Collections.emptyMap())
                },
                (CheckResult) readRequests -> {
                    assertEquals(2, readRequests.size());
                    PlcReadRequest firstReadRequest = readRequests.get(0);
                    assertEquals(1, firstReadRequest.getNumberOfTags());
                    PlcTag firstTag = firstReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagHoldingRegister.class, firstTag);
                    ModbusTagHoldingRegister holdingRegister = (ModbusTagHoldingRegister) firstTag;
                    assertEquals(0, holdingRegister.getAddress());
                    assertEquals(1, holdingRegister.getNumberOfElements());
                    assertEquals(ModbusDataType.WORD, holdingRegister.getDataType());
                    PlcReadRequest secondReadRequest = readRequests.get(1);
                    assertEquals(1, secondReadRequest.getNumberOfTags());
                    firstTag = secondReadRequest.getTags().get(0);
                    assertInstanceOf(ModbusTagHoldingRegister.class, firstTag);
                    holdingRegister = (ModbusTagHoldingRegister) firstTag;
                    assertEquals(2100, holdingRegister.getAddress());
                    assertEquals(1, holdingRegister.getNumberOfElements());
                    assertEquals(ModbusDataType.WORD, holdingRegister.getDataType());
                })
        );
    }

    @ParameterizedTest
    @MethodSource("holdingRegisterInputData")
    void holdingRegisterTests(PlcTag[] tags, CheckResult check) {
        processReadRequest(tags, check);
    }





    void processReadRequest(PlcTag[] tags, CheckResult check) {
        PlcReader reader = Mockito.mock(PlcReader.class);
        ModbusContext driverContext = Mockito.mock(ModbusContext.class);
        Mockito.when(driverContext.getByteOrder()).thenReturn(ModbusByteOrder.BIG_ENDIAN);
        Mockito.when(driverContext.getMaxCoilsPerRequest()).thenReturn(2000);
        Mockito.when(driverContext.getMaxRegistersPerRequest()).thenReturn(125);
        LinkedHashMap<String, PlcTagItem<PlcTag>> tagMap = new LinkedHashMap<>();
        int i = 0;
        for (PlcTag tag : tags) {
            tagMap.put("tag" + i++, new DefaultPlcTagItem<>(tag));
        }
        ModbusOptimizer sut = new ModbusOptimizer();
        List<PlcReadRequest> plcReadRequests = sut.processReadRequest(new DefaultPlcReadRequest(reader, tagMap), driverContext);
        check.isValid(plcReadRequests);
    }

    @FunctionalInterface
    protected interface CheckResult {
        void isValid(List<PlcReadRequest> readRequests) throws AssertionFailedError;
    }

}