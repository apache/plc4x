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
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * What the caller receives has to match what {@code getArrayInfo()} promises.
 *
 * <p>The decoder decided from the element count, so a one-element range came back as a scalar
 * while the tag reported it as an array - a consumer reading the tag to decide how to render the
 * value was told one thing and handed another.</p>
 */
class EipDecodedShapeTest {

    private static final byte[] ONE_DINT = {0x2A, 0x00, 0x00, 0x00};

    @Test
    @DisplayName("a bare index yields a scalar")
    void aBareIndexIsAScalar() {
        EipTag tag = EipTag.of("%rate[4]:DINT");
        assertTrue(tag.getArrayInfo().isEmpty(), "the tag says scalar");

        PlcValue value = EipTcpConnection.parsePlcValue(tag, ONE_DINT, CIPDataTypeCode.DINT);
        assertFalse(value.isList(), "and so is the value");
        assertEquals(42, value.getInt());
    }

    @Test
    @DisplayName("a one-element range yields a list of one")
    void aOneElementRangeIsAList() {
        EipTag tag = EipTag.of("%rate[4..4]:DINT");
        assertFalse(tag.getArrayInfo().isEmpty(), "the tag says array");

        PlcValue value = EipTcpConnection.parsePlcValue(tag, ONE_DINT, CIPDataTypeCode.DINT);
        assertTrue(value.isList(), "and so is the value - this is what the count could not express");
        assertEquals(1, value.getList().size());
        assertEquals(42, value.getList().get(0).getInt());
    }

    @Test
    @DisplayName("the value's shape always follows the tag's")
    void theShapesAgree() {
        for (String address : new String[]{"%rate:DINT", "%rate[4]:DINT", "%rate[4..4]:DINT"}) {
            EipTag tag = EipTag.of(address);
            PlcValue value = EipTcpConnection.parsePlcValue(tag, ONE_DINT, CIPDataTypeCode.DINT);
            assertEquals(!tag.getArrayInfo().isEmpty(), value.isList(), address);
        }
    }

    @Test
    @DisplayName("a count the request cannot carry is refused")
    void refusesACountTheRequestCannotCarry() {
        // The CIP element count is 16 bits: 65536 elements narrow to zero on the wire, so the
        // device would be asked for nothing at all.
        assertThrows(org.apache.plc4x.java.api.exceptions.PlcInvalidTagException.class,
            () -> EipTag.of("%arr[0..65535]:DINT"));

        // One below the limit is asked for in full.
        assertEquals(65535, EipTag.of("%arr[0..65534]:DINT").getElementNb());
    }
}
