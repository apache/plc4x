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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.modbus.field.ModbusFieldHoldingRegister;
import org.apache.plc4x.java.modbus.field.ModbusFieldInputRegister;
import org.apache.plc4x.java.modbus.field.ModbusExtendedRegister;
import org.apache.plc4x.java.modbus.field.ModbusFieldDiscreteInput;
import org.apache.plc4x.java.modbus.field.ModbusFieldCoil;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.modbus.field.*;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class ModbusEncodeTest {

    @Test
    public void testEncodeIntegerINT() {
        Integer[] object = {1,-1,32000,-32000,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:INT");
        ModbusFieldHandler handler = new ModbusFieldHandler();
        PlcList list = (PlcList) handler.encodeInteger(holdingregister, object);
        Assertions.assertEquals("[1,-1,32000,-32000,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeIntegerUINT() {
        Integer[] object = {1,65535,10,55000,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:UINT");
        ModbusFieldHandler handler = new ModbusFieldHandler();
        PlcList list = (PlcList) handler.encodeInteger(holdingregister, object);
        Assertions.assertEquals("[1,65535,10,55000,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeIntegerDINT() {
        Integer[] object = {1,655354775,-2147483648,2147483647,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:DINT");
        ModbusFieldHandler handler = new ModbusFieldHandler();
        PlcList list = (PlcList) handler.encodeInteger(holdingregister, object);
        Assertions.assertEquals("[1,655354775,-2147483648,2147483647,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeFloatREAL() {
        Float[] object = {(float) 1.1,(float) 1000.1,(float) 100000.1,(float) 6363.9,(float) 879.873,(float) 6,(float) 7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:REAL");
        ModbusFieldHandler handler = new ModbusFieldHandler();
        PlcList list = (PlcList) handler.encodeFloat(holdingregister, object);
        Assertions.assertEquals("[1.1,1000.1,100000.1,6363.9,879.873,6.0,7.0]", list.toString());
    }

}
