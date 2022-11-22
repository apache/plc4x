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
package org.apache.plc4x.java.canopen.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.canopen.readwrite.CANOpenDataType;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.canopen.tag.CANOpenPDOTag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class CANOpenPDOTagTest {

    @Test
    public void testNodeSyntax() {
        final CANOpenPDOTag canTag = CANOpenPDOTag.of("RECEIVE_PDO_2:20:RECORD");

        assertEquals(20, canTag.getNodeId());
        assertEquals(CANOpenService.RECEIVE_PDO_2, canTag.getService());
        assertEquals(CANOpenDataType.RECORD, canTag.getCanOpenDataType());
    }

    @Test
    public void testInvalidSyntax() {
        assertThrows(PlcInvalidTagException.class, () -> CANOpenPDOTag.of("PDO:"));
    }
}