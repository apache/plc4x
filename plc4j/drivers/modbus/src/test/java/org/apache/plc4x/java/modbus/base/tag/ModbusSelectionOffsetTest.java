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
package org.apache.plc4x.java.modbus.base.tag;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A selection offset counts <em>elements</em>; a Modbus address counts <em>registers</em>. The two
 * are the same number only for a type that occupies one register, which is why this went unnoticed:
 * every example used INT.
 *
 * <p>The read length already scales - {@code getLengthWords()} multiplies by the data type's size -
 * so an unscaled offset does not shorten the read, it moves it, silently, to the wrong registers.
 */
class ModbusSelectionOffsetTest {

    @Test
    @DisplayName("a one-register type advances one register per element")
    void oneRegisterPerElement() {
        // 40001 is address 0; the fifth INT is four registers along.
        assertEquals(4, ModbusTagHoldingRegister.of("holding-register:1[4]:INT").getAddress());
    }

    @Test
    @DisplayName("a two-register type advances two registers per element")
    void twoRegistersPerElement() {
        // The fifth DINT begins eight registers along, not four.
        assertEquals(8, ModbusTagHoldingRegister.of("holding-register:1[4]:DINT").getAddress());
    }

    @Test
    @DisplayName("a four-register type advances four registers per element")
    void fourRegistersPerElement() {
        assertEquals(8, ModbusTagHoldingRegister.of("holding-register:1[2]:LINT").getAddress());
    }

    @Test
    @DisplayName("a string advances by its declared length")
    void stringsAdvanceByTheirDeclaredLength() {
        // A STRING(20) occupies ten registers, so the third one begins twenty registers along.
        assertEquals(20, ModbusTagHoldingRegister.of("holding-register:1[2]:STRING(20)").getAddress());
    }

    @Test
    @DisplayName("the same rule applies to input and extended registers")
    void theOtherRegisterAreasScaleToo() {
        assertEquals(8, ModbusTagInputRegister.of("input-register:1[4]:DINT").getAddress());
        // An extended register address is not shifted by the protocol offset the way the others
        // are, so its base stays 1 and the eight-register advance lands on 9.
        assertEquals(9, ModbusTagExtendedRegister.of("extended-register:1[4]:DINT").getAddress());
    }

    @Test
    @DisplayName("a bit area addresses bits, so its offset is not scaled")
    void bitAreasAreNotScaled() {
        // One coil is one address; there is no element size to multiply by.
        assertEquals(4, ModbusTagCoil.of("coil:1[4]:BOOL").getAddress());
        assertEquals(4, ModbusTagDiscreteInput.of("discrete-input:1[4]:BOOL").getAddress());
    }

    @Test
    @DisplayName("a bare index is a scalar; a one-element range is a list of one")
    void aRangeIsAnArrayEvenWhenItSpansOneElement() {
        // The notation's own rule, and what a consumer reads getArrayInfo() to decide. The count
        // cannot express it: both of these select exactly one element.
        assertEquals(0, ModbusTagHoldingRegister.of("holding-register:1[4]:INT").getArrayInfo().size(),
            "a bare index selects one element, which is a scalar");
        assertEquals(1, ModbusTagHoldingRegister.of("holding-register:1[4..4]:INT").getArrayInfo().size(),
            "a range is an array even when it spans one element");
        assertTrue(ModbusTagHoldingRegister.of("holding-register:1[4..4]:INT").getArrayInfo().get(0).isRange());
    }

    @Test
    @DisplayName("an address with no selection at all is a scalar")
    void noSelectionIsAScalar() {
        assertEquals(0, ModbusTagHoldingRegister.of("holding-register:1:INT").getArrayInfo().size());
    }

    @Test
    @DisplayName("a multi-element range still reports its elements")
    void aRangeReportsItsElements() {
        assertEquals(1, ModbusTagHoldingRegister.of("holding-register:1[0..3]:INT").getArrayInfo().size());
        assertEquals(4, ModbusTagHoldingRegister.of("holding-register:1[0..3]:INT").getArrayInfo().get(0).getSize());
    }
}
