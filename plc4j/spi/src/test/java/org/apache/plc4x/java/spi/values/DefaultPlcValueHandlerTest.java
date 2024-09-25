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
            // TODO: Add some range tests (values above the max and below the min value)

            // WORD values
            Arguments.of(new MockTag("mock", PlcValueType.WORD), true, new PlcWORD(1)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), false, new PlcWORD(0)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), (byte) 42, new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), (short) 42, new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), 42, new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), (long) 42, new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), (float) 42.34, new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), (double) 42.34, new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), BigInteger.TEN, new PlcWORD(10)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), BigDecimal.TEN, new PlcWORD(10)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), "42", new PlcWORD(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // DWORD values
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), true, new PlcDWORD(1)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), false, new PlcDWORD(0)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), (byte) 42, new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), (short) 42, new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), 42, new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), (long) 42, new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), (float) 42.34, new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), (double) 42.34, new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), BigInteger.TEN, new PlcDWORD(10)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), BigDecimal.TEN, new PlcDWORD(10)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), "42", new PlcDWORD(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // LWORD values
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), true, new PlcLWORD(1)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), false, new PlcLWORD(0)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), (byte) 42, new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), (short) 42, new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), 42, new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), (long) 42, new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), (float) 42.34, new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), (double) 42.34, new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), BigInteger.TEN, new PlcLWORD(10)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), BigDecimal.TEN, new PlcLWORD(10)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), "42", new PlcLWORD(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // USINT values
            Arguments.of(new MockTag("mock", PlcValueType.USINT), true, new PlcUSINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), false, new PlcUSINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), (byte) 42, new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), (short) 42, new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), 42, new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), (long) 42, new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), (float) 42.34, new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), (double) 42.34, new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), BigInteger.TEN, new PlcUSINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), BigDecimal.TEN, new PlcUSINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), "42", new PlcUSINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // UINT values
            Arguments.of(new MockTag("mock", PlcValueType.UINT), true, new PlcUINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), false, new PlcUINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), (byte) 42, new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), (short) 42, new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), 42, new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), (long) 42, new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), (float) 42.34, new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), (double) 42.34, new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), BigInteger.TEN, new PlcUINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), BigDecimal.TEN, new PlcUINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), "42", new PlcUINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // UDINT values
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), true, new PlcUDINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), false, new PlcUDINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), (byte) 42, new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), (short) 42, new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), 42, new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), (long) 42, new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), (float) 42.34, new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), (double) 42.34, new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), BigInteger.TEN, new PlcUDINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), BigDecimal.TEN, new PlcUDINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), "42", new PlcUDINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // ULINT values
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), true, new PlcULINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), false, new PlcULINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), (byte) 42, new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), (short) 42, new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), 42, new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), (long) 42, new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), (float) 42.34, new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), (double) 42.34, new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), BigInteger.TEN, new PlcULINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), BigDecimal.TEN, new PlcULINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), "42", new PlcULINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // SINT values
            Arguments.of(new MockTag("mock", PlcValueType.SINT), true, new PlcSINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), false, new PlcSINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), (byte) 42, new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), (short) 42, new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), 42, new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), (long) 42, new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), (float) 42.34, new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), (double) 42.34, new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), BigInteger.TEN, new PlcSINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), BigDecimal.TEN, new PlcSINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), "42", new PlcSINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // INT values
            Arguments.of(new MockTag("mock", PlcValueType.INT), true, new PlcINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), false, new PlcINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), (byte) 42, new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), (short) 42, new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), 42, new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), (long) 42, new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), (float) 42.34, new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), (double) 42.34, new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), BigInteger.TEN, new PlcINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), BigDecimal.TEN, new PlcINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), "42", new PlcINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // DINT values
            Arguments.of(new MockTag("mock", PlcValueType.DINT), true, new PlcDINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), false, new PlcDINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), (byte) 42, new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), (short) 42, new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), 42, new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), (long) 42, new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), (float) 42.34, new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), (double) 42.34, new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), BigInteger.TEN, new PlcDINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), BigDecimal.TEN, new PlcDINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), "42", new PlcDINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // LINT values
            Arguments.of(new MockTag("mock", PlcValueType.LINT), true, new PlcLINT(1)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), false, new PlcLINT(0)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), (byte) 42, new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), (short) 42, new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), 42, new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), (long) 42, new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), (float) 42.34, new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), (double) 42.34, new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), BigInteger.TEN, new PlcLINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), BigDecimal.TEN, new PlcLINT(10)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), "42", new PlcLINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // REAL values
            Arguments.of(new MockTag("mock", PlcValueType.REAL), 1.3, new PlcREAL(1.3)),
            // TODO: Add some range tests (values above the max and below the min value)

            // LREAL values
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), 1.4, new PlcLREAL(1.4)),
            // TODO: Add some range tests (values above the max and below the min value)

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