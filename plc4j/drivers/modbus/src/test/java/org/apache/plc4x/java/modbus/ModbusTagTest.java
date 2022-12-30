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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ModbusTagTest {

    @Test
    void testHolding_INT_ARRAY_RANGE() {
        for (int i = 1; i <= 125; i++) {
          final ModbusTagHoldingRegister holdingregister = ModbusTagHoldingRegister.of("400001:INT[" + i + "]");
          assertEquals(i, holdingregister.getNumberOfElements());
        }

        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () -> {
            ModbusTagHoldingRegister.of("400001:INT[126]");
        });
        assertTrue(exception.getMessage().startsWith("quantity may not be larger than "));
    }

    @Test
    void testInput_INT_ARRAY_RANGE() {
        for (int i = 1; i <= 125; i++) {
          final ModbusTagInputRegister inputregister = ModbusTagInputRegister.of("300001:INT[" + i + "]");
          assertEquals(i, inputregister.getNumberOfElements());
        }
        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () -> {
            ModbusTagInputRegister.of("300001:INT[126]");
        });
        assertTrue(exception.getMessage().startsWith("quantity may not be larger than "));
    }

    @Test
    void testExtended_INT_ARRAY_RANGE() {
        for (int i = 1; i <= 125; i++) {
          final ModbusTagExtendedRegister extendedRegister = ModbusTagExtendedRegister.of("600001:INT[" + i + "]");
          assertEquals(i, extendedRegister.getNumberOfElements());
        }
        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () -> {
            ModbusTagExtendedRegister.of("600001:INT[126]");
        });
        assertTrue(exception.getMessage().startsWith("quantity may not be larger than "));
    }

    @Test
    void testCoil_INT_ARRAY_RANGE() {
        for (int i = 1; i <= 2000; i++) {
          final ModbusTagCoil coil = ModbusTagCoil.of("000001:BOOL[" + i + "]");
          assertEquals(i, coil.getNumberOfElements());
        }
        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () -> {
            ModbusTagCoil.of("000001:BOOL[2001]");
        });
        assertTrue(exception.getMessage().startsWith("quantity may not be larger than "));
    }

    @Test
    void testDiscreteInput_INT_ARRAY_RANGE() {
        for (int i = 1; i <= 2000; i++) {
          final ModbusTagDiscreteInput discreteInput = ModbusTagDiscreteInput.of("100001:BOOL[" + i + "]");
          assertEquals(i, discreteInput.getNumberOfElements());
        }
        IllegalArgumentException exception = assertThrows(IllegalArgumentException.class, () -> {
            ModbusTagDiscreteInput.of("100001:BOOL[2001]");
        });
        assertTrue(exception.getMessage().startsWith("quantity may not be larger than "));
    }

}
