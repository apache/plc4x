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

package org.apache.plc4x.java.s7.model;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.java.s7.netty.model.types.TransportSize;
import org.apache.plc4x.test.FastTests;
import org.junit.experimental.categories.Category;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.util.stream.Stream;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.*;

class S7FieldTests {

    private static Stream<Arguments> validFieldQueries() {
        return Stream.of(
            Arguments.of("%I0.1:BOOL",          TransportSize.BOOL,  MemoryArea.INPUTS,      0,  0,  1),
            Arguments.of("%ID64:REAL",          TransportSize.REAL,  MemoryArea.INPUTS,      0,  64, 0),
            Arguments.of("%Q0.4:BOOL",          TransportSize.BOOL,  MemoryArea.OUTPUTS,     0,  0,  4),
            Arguments.of("%DB1.DBX38.1:BOOL",   TransportSize.BOOL,  MemoryArea.DATA_BLOCKS, 1,  38, 1)/*,
            // Not quite sure about how Data Block addresses look like, in my TIA portal they all have the prefix "DB".
            Arguments.of("%DB3.DX4.1:BOOL",     S7DataType.BOOL,  MemoryArea.DATA_BLOCKS, 3,  4,  1),
            Arguments.of("%DB3.DB4:INT",        S7DataType.INT,   MemoryArea.DATA_BLOCKS, 3,  4,  0),
            Arguments.of("%DB3.DB4:UINT",       S7DataType.UINT,  MemoryArea.DATA_BLOCKS, 3,  4,  0),
            Arguments.of("%DB3.DW4:REAL",       S7DataType.REAL,  MemoryArea.DATA_BLOCKS, 3,  4,  0)*/
        );
    }

    private static Stream<Arguments> invalidFieldQueries() {
        return Stream.of(
            Arguments.of("%I0:BOOL"),
            Arguments.of("%IW64:REAL"),
            Arguments.of("%DB1.DBX38:BOOL")
        );
    }

    @ParameterizedTest
    @Category(FastTests.class)
    @MethodSource("validFieldQueries")
    void testValidFieldQueryParsing(String fieldQuery, TransportSize expectedClientType, MemoryArea expectedMemoryArea,
                                    int expectedMemoryBlockNumber, int expectedByteOffset, int expectedBitOffset) {
        S7Field field = S7Field.of(fieldQuery);
        assertThat(field, notNullValue());
        assertThat(field.getDataType(), equalTo(expectedClientType));
        assertThat(field.getMemoryArea(), equalTo(expectedMemoryArea));
        assertThat(field.getBlockNumber(), equalTo(expectedMemoryBlockNumber));
        assertThat(field.getByteOffset(), equalTo(expectedByteOffset));
        assertThat(field.getBitOffset(), equalTo((short) expectedBitOffset));
    }

    @ParameterizedTest
    @Category(FastTests.class)
    @MethodSource("invalidFieldQueries")
    void testInvalidFieldQueryParsing(String fieldQuery) {
        try {
            S7Field.of(fieldQuery);
            fail("Should have thrown an exception");
        } catch (PlcRuntimeException e) {
            // This was expected.
        }
    }

    @Test
    void getDefaultJavaType() {
        final PlcField field = S7Field.of("%DB1.DBX38.1:BOOL");

        assertThat(field.getDefaultJavaType(), equalTo(Boolean.class));
    }

    @Test
    void checkGreedyNumFieldsParsing() {
        S7Field field = S7Field.of("%DB56.DBB100:SINT[25]");

        assertEquals(25, field.getNumElements());
    }
}