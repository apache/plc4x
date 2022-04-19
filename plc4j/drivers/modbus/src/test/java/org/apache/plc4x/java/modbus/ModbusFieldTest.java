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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.modbus.base.field.ModbusFieldHoldingRegister;
import org.apache.plc4x.java.modbus.base.field.ModbusIdentificationRegister;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldInputRegister;
import org.apache.plc4x.java.modbus.base.field.ModbusExtendedRegister;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldDiscreteInput;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldCoil;
import org.apache.plc4x.java.modbus.readwrite.ModbusDeviceInformationLevel;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class ModbusFieldTest {

    @Test
    public void testHolding_INT_ARRAY_RANGE() {
        for (int i = 1; i < 125; i++) {
          final ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("400001:INT[" + i + "]");
          Assertions.assertEquals(i, holdingregister.getNumberOfElements());
        }
    }

    @Test
    public void testInput_INT_ARRAY_RANGE() {
        for (int i = 1; i < 125; i++) {
          final ModbusFieldInputRegister inputregister = ModbusFieldInputRegister.of("300001:INT[" + i + "]");
          Assertions.assertEquals(i, inputregister.getNumberOfElements());
        }
    }

    @Test
    public void testExtended_INT_ARRAY_RANGE() {
        for (int i = 1; i < 125; i++) {
          final ModbusExtendedRegister extendedRegister = ModbusExtendedRegister.of("600001:INT[" + i + "]");
          Assertions.assertEquals(i, extendedRegister.getNumberOfElements());
        }
    }

    @Test
    public void testCoil_INT_ARRAY_RANGE() {
        for (int i = 1; i < 2000; i++) {
          final ModbusFieldCoil coil = ModbusFieldCoil.of("000001:BOOL[" + i + "]");
          Assertions.assertEquals(i, coil.getNumberOfElements());
        }
    }

    @Test
    public void testDiscreteInput_INT_ARRAY_RANGE() {
        for (int i = 1; i < 2000; i++) {
          final ModbusFieldDiscreteInput discreteInput = ModbusFieldDiscreteInput.of("100001:BOOL[" + i + "]");
          Assertions.assertEquals(i, discreteInput.getNumberOfElements());
        }
    }

    @Test
    public void testIdentificationField() {
        ModbusDeviceInformationLevel level = ModbusDeviceInformationLevel.EXTENDED;
        short objectId = 10;

        ModbusIdentificationRegister identification = ModbusIdentificationRegister.of("identification:EXTENDED:10");
        Assertions.assertEquals(level, identification.getLevel());
        Assertions.assertEquals(objectId, identification.getObjectId());

        identification = ModbusIdentificationRegister.of("identification:0x03:10");
        Assertions.assertEquals(level, identification.getLevel());
        Assertions.assertEquals(objectId, identification.getObjectId());

        identification = ModbusIdentificationRegister.of("identification:0x3:0xA");
        Assertions.assertEquals(level, identification.getLevel());
        Assertions.assertEquals(objectId, identification.getObjectId());
    }
}
