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

package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collections;
import java.util.List;
import java.util.StringJoiner;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.*;

class DefaultPlcValueHandlerTest {

    private static Stream<Arguments> getSingleElementPlcValues() {
        return Stream.of(
            // BOOL values
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), true, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), false, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (byte) 42, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (byte) 0, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (short) 42, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (short) 0, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), 42, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), 0, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (long) 42, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (long) 0, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (float) 4.2, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (float) 0.0, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (double) 4.2, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), (double) 0.0, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), BigInteger.ONE, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), BigInteger.ZERO, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), BigDecimal.ONE, new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), BigInteger.ZERO, new PlcBOOL(false)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), "true", new PlcBOOL(true)),
            Arguments.of(new MockTag("mock", PlcValueType.BOOL), "false", new PlcBOOL(false)),

            // BYTE values
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), true, new PlcBYTE(1)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), false, new PlcBYTE(0)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), (byte) 42, new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), (short) 42, new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), 42, new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), (long) 42, new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), (float) 42.34, new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), (double) 42.34, new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), BigInteger.TEN, new PlcBYTE(10)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), BigDecimal.TEN, new PlcBYTE(10)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), "42", new PlcBYTE(42)),

            // WORD values
            Arguments.of(new MockTag("mock", PlcValueType.WORD), 2, new PlcWORD(2)),

            // DWORD values
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), 3, new PlcDWORD(3)),

            // LWORD values
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), 4, new PlcLWORD(4)),

            // USINT values
            Arguments.of(new MockTag("mock", PlcValueType.USINT), 5, new PlcUSINT(5)),

            // UINT values
            Arguments.of(new MockTag("mock", PlcValueType.UINT), 6, new PlcUINT(6)),

            // UDINT values
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), 7, new PlcUDINT(7)),

            // ULINT values
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), 8, new PlcULINT(8)),

            // SINT values
            Arguments.of(new MockTag("mock", PlcValueType.SINT), 9, new PlcSINT(9)),

            // INT values
            Arguments.of(new MockTag("mock", PlcValueType.INT), 10, new PlcINT(10)),

            // DINT values
            Arguments.of(new MockTag("mock", PlcValueType.DINT), 11, new PlcDINT(11)),

            // LINT values
            Arguments.of(new MockTag("mock", PlcValueType.LINT), 12, new PlcLINT(12)),

            // REAL values
            Arguments.of(new MockTag("mock", PlcValueType.REAL), 1.3, new PlcREAL(1.3)),

            // LREAL values
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), 1.4, new PlcLREAL(1.4)),

            // CHAR values
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), "A", new PlcCHAR("A")),

            // WCHAR values
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), "B", new PlcWCHAR("B")),

            // STRING values
            Arguments.of(new MockTag("mock", PlcValueType.STRING), "Wolf", new PlcSTRING("Wolf")),

            // WSTRING values
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), "Lamm", new PlcWSTRING("Lamm")),

            // TIME values
            Arguments.of(new MockTag("mock", PlcValueType.TIME), 1, new PlcTIME(1)),

            // LTIME values
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), 1, new PlcLTIME(1)),

            // DATE values
            Arguments.of(new MockTag("mock", PlcValueType.DATE), 1, new PlcDATE(1)),

            // LDATE values
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), 1, new PlcLDATE(LocalDate.of(1978, 3,28))),

            // TIME_OF_DAY values
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), 1, new PlcTIME_OF_DAY(LocalTime.of(9, 11))),

            // LTIME_OF_DAY values
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), 1, new PlcLTIME_OF_DAY(LocalTime.of(1, 2, 3))),

            // DATE_AND_TIME values
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), 1, new PlcDATE_AND_TIME(LocalDateTime.of(1978, 3, 28, 1, 2, 3))),

            // DATE_AND_LTIME values
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), 1, new PlcDATE_AND_LTIME(LocalDateTime.of(1978, 3, 28, 3, 2, 1))),

            // LDATE_AND_TIME values
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), 1, new PlcLDATE_AND_TIME(LocalDateTime.of(2015, 10,21, 7, 28))),

            // RAW_BYTE_ARRAY values
            Arguments.of(new MockTag("mock", PlcValueType.RAW_BYTE_ARRAY), 1, new PlcRawByteArray(new byte[]{(byte) 1, (byte) 2, (byte) 3, (byte) 4, (byte) 5, (byte) 6}))
        );
    }

    @ParameterizedTest
    @MethodSource("getSingleElementPlcValues")
    void newPlcValue(PlcTag plcTag, Object input, PlcValue expected) {
        PlcValueHandler sut = new DefaultPlcValueHandler();

        PlcValue plcValue = sut.newPlcValue(plcTag, input);

        assertNotNull(plcValue);
        assertEquals(expected, plcValue);
    }

    @Test
    void newPlcValues() {
    }

    public static class MockTag implements PlcTag {
        private final String addressString;
        private final PlcValueType plcValueType;
        private final List<ArrayInfo> arrayInfo;

        public MockTag(String addressString, PlcValueType plcValueType) {
            this(addressString, plcValueType, Collections.emptyList());
        }

        public MockTag(String addressString, PlcValueType plcValueType, List<ArrayInfo> arrayInfo) {
            this.addressString = addressString;
            this.plcValueType = plcValueType;
            this.arrayInfo = arrayInfo;
        }

        @Override
        public String getAddressString() {
            return addressString;
        }

        @Override
        public PlcValueType getPlcValueType() {
            return plcValueType;
        }

        @Override
        public List<ArrayInfo> getArrayInfo() {
            return arrayInfo;
        }

        @Override
        public String toString() {
            return new StringJoiner(", ", MockTag.class.getSimpleName() + "[", "]")
                .add("addressString='" + addressString + "'")
                .add("plcValueType=" + plcValueType)
                .add("arrayInfo=" + arrayInfo)
                .toString();
        }
    }

}