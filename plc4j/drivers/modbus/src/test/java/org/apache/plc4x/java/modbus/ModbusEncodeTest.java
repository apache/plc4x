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

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldHoldingRegister;
import org.apache.plc4x.java.modbus.base.field.ModbusFieldCoil;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import java.math.BigInteger;

public class ModbusEncodeTest {

    @Test
    public void testEncodeBooleanBOOL() {
        Boolean[] object = {true,false,true,false,true,false,true,true,false};
        ModbusFieldCoil coils = ModbusFieldCoil.of("coil:8:BOOL");
        PlcList list = (PlcList) PlcValueHandler.of(coils, object);
        Assertions.assertEquals("[true,false,true,false,true,false,true,true,false]", list.toString());
    }

    @Test
    public void testEncodeIntegerSINT() {
        Integer[] object = {1,-1,127,-128,5,6,7,8};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:8:SINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,-1,127,-128,5,6,7,8]", list.toString());
    }

    @Test
    public void testEncodeIntegerUSINT() {
        Integer[] object = {1,255,0,4,5,6,7,8};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:8:USINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,255,0,4,5,6,7,8]", list.toString());
    }

    @Test
    public void testEncodeIntegerBYTE() {
        Integer[] object = {1,255,0,4,5,6,7,8};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:8:BYTE");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,255,0,4,5,6,7,8]", list.toString());
    }

    @Test
    public void testEncodeIntegerINT() {
        Integer[] object = {1,-1,32000,-32000,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:INT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,-1,32000,-32000,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeIntegerUINT() {
        Integer[] object = {1,65535,10,55000,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:UINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,65535,10,55000,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeIntegerWORD() {
        Integer[] object = {1,65535,10,55000,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:WORD");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,65535,10,55000,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeIntegerDINT() {
        Integer[] object = {1,655354775,-2147483648,2147483647,5,6,7};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:DINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,655354775,-2147483648,2147483647,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeLongUDINT() {
        Long[] object = {1L,655354775L,0L,4294967295L,5L,6L,7L};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:UDINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,655354775,0,4294967295,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeLongDWORD() {
        Long[] object = {1L,655354775L,0L,4294967295L,5L,6L,7L};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:DWORD");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,655354775,0,4294967295,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeLongLINT() {
        Long[] object = {1L,655354775L,-9223372036854775808L,9223372036854775807L,5L,6L,7L};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:LINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,655354775,-9223372036854775808,9223372036854775807,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeBigIntegerULINT() {
        BigInteger[] object = {BigInteger.valueOf(1L),BigInteger.valueOf(655354775L),BigInteger.valueOf(0),new BigInteger("18446744073709551615"),BigInteger.valueOf(5L),BigInteger.valueOf(6L),BigInteger.valueOf(7L)};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:ULINT");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,655354775,0,18446744073709551615,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeBigIntegerLWORD() {
        BigInteger[] object = {BigInteger.valueOf(1L),BigInteger.valueOf(655354775L),BigInteger.valueOf(0),new BigInteger("18446744073709551615"),BigInteger.valueOf(5L),BigInteger.valueOf(6L),BigInteger.valueOf(7L)};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:LWORD");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1,655354775,0,18446744073709551615,5,6,7]", list.toString());
    }

    @Test
    public void testEncodeFloatREAL() {
        Float[] object = {1.1f,1000.1f,100000.1f,3.4028232E38f,-3.4028232E38f,-1f,10384759934840.0f};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:REAL");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        //! When using Java 19 it seems the toString method uses a different precision than the previous versions,
        //! so we need to check differently in this case.
        for (int i = 0; i < list.getLength(); i++) {
            PlcValue plcValue = list.getIndex(i);
            Float referenceValue = object[i];
            Assertions.assertEquals(referenceValue, plcValue.getFloat());
        }
    }

    @Test
    public void testEncodeDoubleLREAL() {
        Double[] object = {1.1,1000.1,100000.1,1.7E308,-1.7E308,-1d,10384759934840.0};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:7:LREAL");
        PlcList list = (PlcList) PlcValueHandler.of(holdingregister, object);
        Assertions.assertEquals("[1.1,1000.1,100000.1,1.7E308,-1.7E308,-1.0,1.038475993484E13]", list.toString());
    }

    /*@Test
    public void testEncodeStringSTRING() {
        String[] object = {"Hello Toddy!"};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:8:STRING");
        PlcList list = (PlcList) handler.encodeString(holdingregister, object);
        Assertions.assertEquals("[H,e,l,l,o, ,T,o,d,d,y,!]", list.toString());
    }

    @Test
    public void testEncodeStringWSTRING() {
        String[] object = {"Hello Toddy!"};
        ModbusFieldHoldingRegister holdingregister = ModbusFieldHoldingRegister.of("holding-register:8:WSTRING");
        PlcList list = (PlcList) handler.encodeString(holdingregister, object);
        Assertions.assertEquals("[H,e,l,l,o, ,T,o,d,d,y,!]", list.toString());
    } */
}
