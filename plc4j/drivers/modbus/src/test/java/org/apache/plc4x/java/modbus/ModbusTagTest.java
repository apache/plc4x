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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.modbus.base.tag.*;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

import static org.apache.plc4x.java.modbus.base.tag.ModbusTag.PROTOCOL_ADDRESS_OFFSET;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ModbusTagTest {

    private static final Logger LOG = LoggerFactory.getLogger(ModbusTagTest.class);

    private void verifyModbusTag(List<String> tagPatterns,
                                 int allowedMax,
                                 Class<? extends ModbusTag> expectedClass,
                                 int expectedAddressShift) {
        // Ensure all tagpatterns compile to the right tag
        for (int i = 1; i <= allowedMax; i++) {
            List<ModbusTag> tags = new ArrayList<>();
            for (String tagPattern : tagPatterns) {
                final ModbusTag modbusTag = ModbusTag.of(String.format(tagPattern, i));
                assertTrue(expectedClass.isInstance(modbusTag));
                assertEquals(i, modbusTag.getNumberOfElements());
                tags.add(modbusTag);
            }
            // All forms of defining the tag MUST result in an identical modbus tag
            assertEquals(1, tags.stream().distinct().count());
        }

        for (String tagPattern : tagPatterns) {
            IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () ->
                ModbusTag.of(String.format(tagPattern, allowedMax + 1))
            );
            assertTrue(exception.getMessage().startsWith("quantity may not be larger than "));
        }

        // Ensure the getAddressString yields a parseable and identical tag.
        for (String tagPattern : tagPatterns) {
            String addressString1 = String.format(tagPattern, 42);
            LOG.info("Validating {}", addressString1);
            ModbusTag modbusTag1 = ModbusTag.of(addressString1);
            String addressString2 = modbusTag1.getAddressString();
            ModbusTag modbusTag2 = ModbusTag.of(addressString2);
            assertEquals(modbusTag1, modbusTag2, "From input addressString:  " + addressString1);

            // We know ALL examples below request address '1'
            // So this must return the logical address that was requested.
            assertEquals(1, modbusTag1.getLogicalAddress());
            assertEquals(modbusTag1.getLogicalAddress() - expectedAddressShift, modbusTag1.getAddress());

            assertEquals(1, modbusTag2.getLogicalAddress());
            assertEquals(modbusTag2.getLogicalAddress() - expectedAddressShift, modbusTag2.getAddress());
        }

    }

    @Test
    void testCoil_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("coil:1:BOOL[%d]", "00001:BOOL[%d]", "000001:BOOL[%d]", "0x00001:BOOL[%d]"),
            2000,
            ModbusTagCoil.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testDiscreteInput_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("discrete-input:1:BOOL[%d]", "10001:BOOL[%d]", "100001:BOOL[%d]", "1x00001:BOOL[%d]"),
            2000,
            ModbusTagDiscreteInput.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testHolding_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("holding-register:1:INT[%d]", "40001:INT[%d]", "400001:INT[%d]", "4x00001:INT[%d]"),
            125,
            ModbusTagHoldingRegister.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testInput_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("input-register:1:INT[%d]", "30001:INT[%d]", "300001:INT[%d]", "3x00001:INT[%d]"),
            125,
            ModbusTagInputRegister.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testExtended_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("extended-register:1:INT[%d]", "60001:INT[%d]", "600001:INT[%d]", "6x00001:INT[%d]"),
            125,
            ModbusTagExtendedRegister.class,
            0 // Addresses for extended memory start at address 0 instead of 1
        );
    }

}
