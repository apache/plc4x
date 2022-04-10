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
package org.apache.plc4x.java.canopen.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.canopen.readwrite.CANOpenDataType;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class CANOpenSDOFieldTest {

    @Test
    public void testNodeSyntax() {
        final CANOpenSDOField canField = CANOpenSDOField.of("SDO:20:0x10/0xAA:RECORD");

        assertEquals(20, canField.getNodeId());
        assertEquals(20, canField.getAnswerNodeId());
        assertEquals(0x10, canField.getIndex());
        assertEquals(0xAA, canField.getSubIndex());
        assertEquals(CANOpenDataType.RECORD, canField.getCanOpenDataType());
    }

    @Test
    public void testAnswerNodeSyntax() {
        final CANOpenSDOField canField = CANOpenSDOField.of("SDO:20/22:0x10/0xAA:RECORD");

        assertEquals(20, canField.getNodeId());
        assertEquals(22, canField.getAnswerNodeId());
        assertEquals(0x10, canField.getIndex());
        assertEquals(0xAA, canField.getSubIndex());
        assertEquals(CANOpenDataType.RECORD, canField.getCanOpenDataType());
    }

    @Test
    public void testInvalidSyntax() {
        assertThrows(PlcInvalidFieldException.class, () -> CANOpenSDOField.of("SDO:"));
    }

}