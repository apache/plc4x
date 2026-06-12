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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7TagTest {

    @Test
    void parseDataBlockBool() {
        S7Tag tag = S7Tag.of("%DB1.DBX0.0:BOOL");
        assertEquals(TransportSize.BOOL, tag.getDataType());
        assertEquals(MemoryArea.DATA_BLOCKS, tag.getMemoryArea());
        assertEquals(1, tag.getBlockNumber());
        assertEquals(0, tag.getByteOffset());
        assertEquals(0, tag.getBitOffset());
        assertEquals(1, tag.getNumberOfElements());
    }

    @Test
    void parseDataBlockInt() {
        S7Tag tag = S7Tag.of("%DB10.DBW20:INT");
        assertEquals(TransportSize.INT, tag.getDataType());
        assertEquals(10, tag.getBlockNumber());
        assertEquals(20, tag.getByteOffset());
    }

    @Test
    void parseInputBool() {
        S7Tag tag = S7Tag.of("%I0.0:BOOL");
        assertEquals(TransportSize.BOOL, tag.getDataType());
        assertEquals(MemoryArea.INPUTS, tag.getMemoryArea());
        assertEquals(0, tag.getByteOffset());
        assertEquals(0, tag.getBitOffset());
    }

    @Test
    void parseFlagDword() {
        S7Tag tag = S7Tag.of("%MD100:DINT");
        assertEquals(TransportSize.DINT, tag.getDataType());
        assertEquals(MemoryArea.FLAGS_MARKERS, tag.getMemoryArea());
        assertEquals(100, tag.getByteOffset());
    }

    @Test
    void parseDataBlockByteArray() {
        S7Tag tag = S7Tag.of("%DB1.DBB0:BYTE[10]");
        assertEquals(TransportSize.BYTE, tag.getDataType());
        assertEquals(10, tag.getNumberOfElements());
    }

    @Test
    void parseInvalidTagThrows() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("not a tag"));
    }

    @Test
    void boolWithoutBitOffsetThrows() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1.DBX0:BOOL"));
    }

    @Test
    void equalsAndHashCode() {
        S7Tag a = S7Tag.of("%DB1.DBW0:INT");
        S7Tag b = S7Tag.of("%DB1.DBW0:INT");
        S7Tag c = S7Tag.of("%DB1.DBW2:INT");
        assertEquals(a, b);
        assertEquals(a.hashCode(), b.hashCode());
        assertNotEquals(a, c);
    }

    @Test
    void matchesRecognizesValidPatterns() {
        assertTrue(S7Tag.matches("%DB1.DBX0.0:BOOL"));
        assertTrue(S7Tag.matches("%MW0:INT"));
        assertFalse(S7Tag.matches("nonsense"));
    }

    @Test
    void parseDataBlockShortForm() {
        S7Tag tag = S7Tag.of("%DB5:10:INT");
        assertEquals(5, tag.getBlockNumber());
        assertEquals(10, tag.getByteOffset());
        assertEquals(TransportSize.INT, tag.getDataType());
    }

    @Test
    void parseRawByteArrayKeyword() {
        S7Tag tag = S7Tag.of("%DB1.DBB0:RAW_BYTE_ARRAY[16]");
        assertEquals(TransportSize.BYTE, tag.getDataType());
        assertEquals(16, tag.getNumberOfElements());
    }

    @Test
    void plcValueTypeMappings() {
        assertEquals(PlcValueType.INT, S7Tag.of("%MW0:INT").getPlcValueType());
        assertEquals(PlcValueType.BOOL, S7Tag.of("%I0.0:BOOL").getPlcValueType());
        assertEquals(PlcValueType.TIME,
            new S7Tag(TransportSize.S5TIME, MemoryArea.FLAGS_MARKERS, 0, 0, (byte) 0, 1).getPlcValueType());
        assertEquals(PlcValueType.WORD,
            new S7Tag(TransportSize.COUNTER, MemoryArea.COUNTERS, 0, 0, (byte) 0, 1).getPlcValueType());
    }

    @Test
    void arrayInfoNonScalar() {
        S7Tag scalar = S7Tag.of("%MW0:INT");
        assertTrue(scalar.getArrayInfo().isEmpty());
        S7Tag arr = S7Tag.of("%DB1.DBB0:BYTE[10]");
        assertEquals(1, arr.getArrayInfo().size());
        assertEquals(0, arr.getArrayInfo().get(0).getLowerBound());
        assertEquals(9, arr.getArrayInfo().get(0).getUpperBound());
    }

    @Test
    void invalidBlockNumberThrows() {
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB99999.DBW0:INT"));
    }

    @Test
    void mismatchedTransferSizeThrows() {
        // 'DBB' (byte) but data type 'INT' (word) -> mismatch.
        assertThrows(PlcInvalidTagException.class, () -> S7Tag.of("%DB1.DBB0:INT"));
    }

    @Test
    void plcDataTypeReturnsName() {
        assertEquals("INT", S7Tag.of("%MW0:INT").getPlcDataType());
    }

    @Test
    void counterTypeShufflesOffsets() {
        // The constructor splits the byte address into byte/bit for COUNTER.
        S7Tag tag = new S7Tag(TransportSize.COUNTER, MemoryArea.COUNTERS, 0, 9, (byte) 0, 1);
        // 9 -> byte 1, bit 1
        assertEquals(1, tag.getByteOffset());
        assertEquals(1, tag.getBitOffset());
    }

    @Test
    void toStringContainsKeyFields() {
        String s = S7Tag.of("%DB1.DBW0:INT").toString();
        assertTrue(s.contains("INT"));
        assertTrue(s.contains("DATA_BLOCKS"));
    }

}
