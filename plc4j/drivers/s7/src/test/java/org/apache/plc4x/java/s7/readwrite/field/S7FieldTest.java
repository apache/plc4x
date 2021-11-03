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
package org.apache.plc4x.java.s7.readwrite.field;

import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class S7FieldTest {

    @Test
    public void testPlcProxyAddress() {
        final S7Field s7Field = S7Field.of("10-01-00-01-00-2D-84-00-00-08");
        Assertions.assertEquals(TransportSize.BOOL, s7Field.getDataType());
        Assertions.assertEquals(1, s7Field.getNumberOfElements());
        Assertions.assertEquals(45, s7Field.getBlockNumber());
        Assertions.assertEquals(MemoryArea.DATA_BLOCKS, s7Field.getMemoryArea());
        Assertions.assertEquals(1, s7Field.getByteOffset());
        Assertions.assertEquals(0, s7Field.getBitOffset());
    }

    /**
     * PLC4X-312 - CAN NOT READ STRING FROM S7 PLC
     */
    @Test
    public void testStringAddress() {
        final S7Field s7Field = S7Field.of("%DB145.DBX38:STRING(8)[1]");
        Assertions.assertEquals(TransportSize.STRING, s7Field.getDataType());
        Assertions.assertEquals(1, s7Field.getNumberOfElements());
        Assertions.assertEquals(145, s7Field.getBlockNumber());
        Assertions.assertEquals(MemoryArea.DATA_BLOCKS, s7Field.getMemoryArea());
        Assertions.assertEquals(38, s7Field.getByteOffset());
        Assertions.assertEquals(0, s7Field.getBitOffset());
        Assertions.assertEquals("S7StringField", s7Field.getClass().getSimpleName());
        S7StringField s7StringField = (S7StringField) s7Field;
        Assertions.assertEquals(8, s7StringField.getStringLength());
    }

}
