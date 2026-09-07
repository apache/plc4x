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
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;

class ModbusTagTest {

    private static final Logger LOG = LoggerFactory.getLogger(ModbusTagTest.class);

    private void verifyModbusTag(List<String> tagPatterns,
                                 int allowedMax,
                                 Class<? extends ModbusTag> expectedClass,
                                 int expectedAddressShift) {
        // Ensure all tag patterns compile to the right tag
        for (int i = 1; i <= allowedMax; i++) {
            List<ModbusTag> tags = new ArrayList<>();
            for (String tagPattern : tagPatterns) {
                // The templates spell the selection as an inclusive range, so the last index is
                // one below the element count the loop is asserting.
                final ModbusTag modbusTag = ModbusTag.of(String.format(tagPattern, i - 1));
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
            List.of("coil:1[0..%d]:BOOL", "00001[0..%d]:BOOL", "000001[0..%d]:BOOL", "0x00001[0..%d]:BOOL"),
            2000,
            ModbusTagCoil.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testDiscreteInput_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("discrete-input:1[0..%d]:BOOL", "10001[0..%d]:BOOL", "100001[0..%d]:BOOL", "1x00001[0..%d]:BOOL"),
            2000,
            ModbusTagDiscreteInput.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testHolding_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("holding-register:1[0..%d]:INT", "40001[0..%d]:INT", "400001[0..%d]:INT", "4x00001[0..%d]:INT"),
            125,
            ModbusTagHoldingRegister.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testInput_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("input-register:1[0..%d]:INT", "30001[0..%d]:INT", "300001[0..%d]:INT", "3x00001[0..%d]:INT"),
            125,
            ModbusTagInputRegister.class,
            PROTOCOL_ADDRESS_OFFSET
        );
    }

    @Test
    void testExtended_INT_ARRAY_RANGE() {
        verifyModbusTag(
            List.of("extended-register:1[0..%d]:INT", "60001[0..%d]:INT", "600001[0..%d]:INT", "6x00001[0..%d]:INT"),
            125,
            ModbusTagExtendedRegister.class,
            0 // Addresses for extended memory start at address 0 instead of 1
        );
    }


    /**
     * The quantity is capped at 125 and the address range is checked, but both checks run after
     * the number has been parsed - so a count too wide to be a number never reached them and left
     * as a NumberFormatException instead of the invalid address it is.
     */
    @Test
    void aCountTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThrows(PlcInvalidTagException.class,
            () -> ModbusTagHoldingRegister.of("holding-register:1[0..99999999998]:INT"));
    }

    @Test
    void anAddressTooWideToBeANumberIsAlsoAnInvalidTag() {
        assertThrows(PlcInvalidTagException.class,
            () -> ModbusTagHoldingRegister.of("holding-register:99999999999:INT"));
    }

    @Test
    void aStringLengthTooWideToBeANumberIsAlsoAnInvalidTag() {
        assertThrows(PlcInvalidTagException.class,
            () -> ModbusTagHoldingRegister.of("holding-register:1:STRING(99999999999)"));
    }
}
