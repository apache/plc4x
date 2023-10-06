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
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class ModbusTagDataTypeTest {

    @Test
    public void testHolding_DataType() {
        //Datatype, Length in Words
        String[][] datatypes = {{"INT","1"},
                                {"UINT","1"},
                                {"REAL","2"}};
        for (String[] datatype : datatypes) {
            ModbusTagHoldingRegister holdingregister = ModbusTagHoldingRegister.of("holding-register:1:" + datatype[0]);
            Assertions.assertEquals(datatype[0], holdingregister.getDataType().name());
            Assertions.assertEquals(1, holdingregister.getNumberOfElements());
            Assertions.assertEquals(Integer.parseInt(datatype[1]) * 2, holdingregister.getLengthBytes());
            Assertions.assertEquals(Integer.parseInt(datatype[1]), holdingregister.getLengthWords());
        }
    }

    @Test
    public void testInput_DataType() {
        //Datatype, Length in Words
        String[][] datatypes = {{"INT","1"},
                                {"UINT","1"},
                                {"REAL","2"}};
        for (String[] datatype : datatypes) {
            ModbusTagInputRegister inputregister = ModbusTagInputRegister.of("input-register:1:" + datatype[0]);
            Assertions.assertEquals(datatype[0], inputregister.getDataType().name());
            Assertions.assertEquals(1, inputregister.getNumberOfElements());
            Assertions.assertEquals(Integer.parseInt(datatype[1]) * 2, inputregister.getLengthBytes());
            Assertions.assertEquals(Integer.parseInt(datatype[1]), inputregister.getLengthWords());
        }
    }

    @Test
    public void testExtended_DataType() {
        //Datatype, Length in Words
        String[][] datatypes = {{"INT","1"},
                                {"UINT","1"},
                                {"DINT","2"},
                                {"REAL","2"}};
        for (String[] datatype : datatypes) {
            ModbusTagExtendedRegister extendedregister = ModbusTagExtendedRegister.of("extended-register:1:" + datatype[0]);
            Assertions.assertEquals(datatype[0], extendedregister.getDataType().name());
            Assertions.assertEquals(1, extendedregister.getNumberOfElements());
            Assertions.assertEquals(Integer.parseInt(datatype[1]) * 2, extendedregister.getLengthBytes());
            Assertions.assertEquals(Integer.parseInt(datatype[1]), extendedregister.getLengthWords());
        }
    }

    @Test
    public void testCoil_DataType() {
        //Datatype, Length in Bytes
        String[][] datatypes = {{"BOOL","1"}};
        for (String[] datatype : datatypes) {
            ModbusTagCoil coil = ModbusTagCoil.of("coil:1:" + datatype[0]);
            Assertions.assertEquals(datatype[0], coil.getDataType().name());
            Assertions.assertEquals(1, coil.getNumberOfElements());
        }
    }

    @Test
    public void testDiscreteInput_DataType() {
        //Datatype, Length in Bytes
        String[][] datatypes = {{"BOOL","1"}};
        for (String[] datatype : datatypes) {
            ModbusTagDiscreteInput discrete = ModbusTagDiscreteInput.of("discrete-input:1:" + datatype[0]);
            Assertions.assertEquals(datatype[0], discrete.getDataType().name());
            Assertions.assertEquals(1, discrete.getNumberOfElements());
        }
    }
}
