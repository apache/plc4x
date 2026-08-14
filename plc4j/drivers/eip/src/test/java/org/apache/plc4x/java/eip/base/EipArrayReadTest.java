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
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcList;
import org.junit.jupiter.api.Test;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Decoding of multi-element (array) reads - see GH-1008. The driver used to request a single
 * element while decoding {@code elementNb} of them, which either misdecoded the value or threw
 * an IndexOutOfBoundsException out of the response handler.
 */
class EipArrayReadTest {

    @Test
    void arrayOfDintsIsFullyDecoded() {
        EipTag tag = EipTag.of("%N40[0]:DINT:8");
        assertNotNull(tag);
        assertEquals(8, tag.getElementNb());

        byte[] raw = dints(1, 2, 3, 4, 5, 6, 7, 8);
        PlcValue value = EipTcpConnection.parsePlcValue(tag, raw, CIPDataTypeCode.DINT);

        assertInstanceOf(PlcList.class, value);
        assertEquals(8, value.getLength());
        for (int i = 0; i < 8; i++) {
            assertEquals(i + 1, value.getIndex(i).getInt(), "element " + i);
        }
    }

    @Test
    void scalarStillDecodesToAScalar() {
        EipTag tag = EipTag.of("%N40:DINT");
        PlcValue value = EipTcpConnection.parsePlcValue(tag, dints(42), CIPDataTypeCode.DINT);

        assertInstanceOf(PlcDINT.class, value);
        assertEquals(42, value.getInt());
    }

    /**
     * The reporter's crash: a reply carrying fewer elements than the tag declares must not throw.
     */
    @Test
    void shortReplyIsReportedInsteadOfThrowing() {
        EipTag tag = EipTag.of("%N40[0]:DINT:8");

        // Only one element's worth of data for an 8-element tag.
        assertNull(EipTcpConnection.parsePlcValue(tag, dints(1), CIPDataTypeCode.DINT));
    }

    @Test
    void arrayOfIntsIsFullyDecoded() {
        EipTag tag = EipTag.of("%N40[0]:INT:4");

        ByteBuffer buffer = ByteBuffer.allocate(8).order(ByteOrder.LITTLE_ENDIAN);
        for (short s : new short[]{10, 20, 30, 40}) {
            buffer.putShort(s);
        }
        PlcValue value = EipTcpConnection.parsePlcValue(tag, buffer.array(), CIPDataTypeCode.INT);

        assertNotNull(value);
        assertEquals(4, value.getLength());
        assertEquals(10, value.getIndex(0).getInt());
        assertEquals(40, value.getIndex(3).getInt());
    }

    private static byte[] dints(int... values) {
        ByteBuffer buffer = ByteBuffer.allocate(values.length * 4).order(ByteOrder.LITTLE_ENDIAN);
        for (int value : values) {
            buffer.putInt(value);
        }
        return buffer.array();
    }
}
