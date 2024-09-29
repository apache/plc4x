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

import java.lang.reflect.Field;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.OffsetDateTime;
import java.time.ZoneOffset;
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
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), BigInteger.valueOf(42), new PlcBYTE(42)),
            Arguments.of(new MockTag("mock", PlcValueType.BYTE), BigDecimal.valueOf(42.34), new PlcBYTE(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.WORD), BigInteger.valueOf(42), new PlcWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.WORD), BigDecimal.valueOf(42.34), new PlcWORD(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), BigInteger.valueOf(42), new PlcDWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DWORD), BigDecimal.valueOf(42.34), new PlcDWORD(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), BigInteger.valueOf(42), new PlcLWORD(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LWORD), BigDecimal.valueOf(42.34), new PlcLWORD(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.USINT), BigInteger.valueOf(42), new PlcUSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.USINT), BigDecimal.valueOf(42.34), new PlcUSINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.UINT), BigInteger.valueOf(42), new PlcUINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UINT), BigDecimal.valueOf(42.34), new PlcUINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), BigInteger.valueOf(42), new PlcUDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.UDINT), BigDecimal.valueOf(42.34), new PlcUDINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), BigInteger.valueOf(42), new PlcULINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.ULINT), BigDecimal.valueOf(42.34), new PlcULINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.SINT), BigInteger.valueOf(42), new PlcSINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.SINT), BigDecimal.valueOf(42.34), new PlcSINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.INT), BigInteger.valueOf(42), new PlcINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.INT), BigDecimal.valueOf(42.34), new PlcINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.DINT), BigInteger.valueOf(42), new PlcDINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.DINT), BigDecimal.valueOf(42.34), new PlcDINT(42)),
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
            Arguments.of(new MockTag("mock", PlcValueType.LINT), BigInteger.valueOf(42), new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), BigDecimal.valueOf(42.34), new PlcLINT(42)),
            Arguments.of(new MockTag("mock", PlcValueType.LINT), "42", new PlcLINT(42)),
            // TODO: Add some range tests (values above the max and below the min value)

            // REAL values
            Arguments.of(new MockTag("mock", PlcValueType.REAL), true, new PlcREAL(1.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), false, new PlcREAL(0.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), (byte) 42, new PlcREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), (short) 42, new PlcREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), 42, new PlcREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), (long) 42, new PlcREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), (float) 42.34, new PlcREAL(42.34)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), (double) 42.35, new PlcREAL(42.35)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), BigInteger.valueOf(42), new PlcREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), BigDecimal.valueOf(42.34), new PlcREAL(42.34)),
            Arguments.of(new MockTag("mock", PlcValueType.REAL), "42.1", new PlcREAL(42.1)),
            // TODO: Add some range tests (values above the max and below the min value)

            // LREAL values
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), true, new PlcLREAL(1.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), false, new PlcLREAL(0.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), (byte) 42, new PlcLREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), (short) 42, new PlcLREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), 42, new PlcLREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), (long) 42, new PlcLREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), (float) 42.34, new PlcLREAL(42.34000015258789)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), (double) 42.35, new PlcLREAL(42.35)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), BigInteger.valueOf(42), new PlcLREAL(42.0)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), BigDecimal.valueOf(42.34), new PlcLREAL(42.34)),
            Arguments.of(new MockTag("mock", PlcValueType.LREAL), "42.1", new PlcLREAL(42.1)),
            // TODO: Add some range tests (values above the max and below the min value)

            // CHAR values
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), true, new PlcCHAR('T')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), false, new PlcCHAR('F')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), (byte) 65, new PlcCHAR('A')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), (short) 66, new PlcCHAR('B')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), 67, new PlcCHAR('C')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), (long) 68, new PlcCHAR('D')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), (float) 69.34, new PlcCHAR('E')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), (double) 70.35, new PlcCHAR('F')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), BigInteger.valueOf(71), new PlcCHAR('G')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), BigDecimal.valueOf(72), new PlcCHAR('H')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), 'I', new PlcCHAR('I')),
            Arguments.of(new MockTag("mock", PlcValueType.CHAR), "J", new PlcCHAR('J')),

            // WCHAR values
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), true, new PlcWCHAR('T')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), false, new PlcWCHAR('F')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), (byte) 65, new PlcWCHAR('A')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), (short) 66, new PlcWCHAR('B')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), 67, new PlcWCHAR('C')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), (long) 68, new PlcWCHAR('D')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), (float) 69.34, new PlcWCHAR('E')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), (double) 70.35, new PlcWCHAR('F')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), BigInteger.valueOf(71), new PlcWCHAR('G')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), BigDecimal.valueOf(72), new PlcWCHAR('H')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), 'I', new PlcWCHAR('I')),
            Arguments.of(new MockTag("mock", PlcValueType.WCHAR), "J", new PlcWCHAR('J')),

            // STRING values
            Arguments.of(new MockTag("mock", PlcValueType.STRING), true, new PlcSTRING("true")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), false, new PlcSTRING("false")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), (byte) 42, new PlcSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), (short) 42, new PlcSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), 42, new PlcSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), (long) 42, new PlcSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), (float) 42.34, new PlcSTRING("42.34")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), (double) 42.35, new PlcSTRING("42.35")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), BigInteger.valueOf(42), new PlcSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), BigDecimal.valueOf(42.34), new PlcSTRING("42.34")),
            Arguments.of(new MockTag("mock", PlcValueType.STRING), "Wolf", new PlcSTRING("Wolf")),

            // WSTRING values
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), true, new PlcWSTRING("true")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), false, new PlcWSTRING("false")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), (byte) 42, new PlcWSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), (short) 42, new PlcWSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), 42, new PlcWSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), (long) 42, new PlcWSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), (float) 42.34, new PlcWSTRING("42.34")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), (double) 42.35, new PlcWSTRING("42.35")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), BigInteger.valueOf(42), new PlcWSTRING("42")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), BigDecimal.valueOf(42.34), new PlcWSTRING("42.34")),
            Arguments.of(new MockTag("mock", PlcValueType.WSTRING), "Lamm", new PlcWSTRING("Lamm")),

            // TIME values (Numeric values are interpreted as milliseconds)
            Arguments.of(new MockTag("mock", PlcValueType.TIME), Duration.ofMillis(1234), new PlcTIME(Duration.parse("PT1.234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), (byte) 123, new PlcTIME(Duration.parse("PT0.123S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), (short) 1234, new PlcTIME(Duration.parse("PT1.234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), 1234, new PlcTIME(Duration.parse("PT1.234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), 12345678901L, new PlcTIME(Duration.parse("PT3429H21M18.901S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), (float) 1234.56, new PlcTIME(Duration.parse("PT1.234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), (double) 1234.56, new PlcTIME(Duration.parse("PT1.234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), BigInteger.valueOf(12345678901L), new PlcTIME(Duration.parse("PT3429H21M18.901S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), BigDecimal.valueOf(12345678901L), new PlcTIME(Duration.parse("PT3429H21M18.901S"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME), "PT3429H21M18.901S", new PlcTIME(Duration.parse("PT3429H21M18.901S"))),

            // LTIME values (Numeric values are interpreted as nanoseconds)
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), Duration.ofMillis(1234), new PlcLTIME(Duration.parse("PT1.234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), (byte) 123, new PlcLTIME(Duration.parse("PT0.000000123S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), (short) 1234, new PlcLTIME(Duration.parse("PT0.000001234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), 1234, new PlcLTIME(Duration.parse("PT0.000001234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), 12345678901L, new PlcLTIME(Duration.parse("PT12.345678901S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), (float) 1234.56, new PlcLTIME(Duration.parse("PT0.000001234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), (double) 1234.56, new PlcLTIME(Duration.parse("PT0.000001234S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), BigInteger.valueOf(12345678901L), new PlcLTIME(Duration.parse("PT12.345678901S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), BigDecimal.valueOf(12345678901L), new PlcLTIME(Duration.parse("PT12.345678901S"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME), "PT12.345678901S", new PlcLTIME(Duration.parse("PT12.345678901S"))),

            // DATE values (Numeric values are interpreted as days since epoch)
            Arguments.of(new MockTag("mock", PlcValueType.DATE), LocalDate.ofEpochDay(1234), new PlcDATE(LocalDate.ofEpochDay(1234))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), (byte) 123, new PlcDATE(LocalDate.ofEpochDay(123))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), (short) 1234, new PlcDATE(LocalDate.ofEpochDay(1234))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), 1234, new PlcDATE(LocalDate.ofEpochDay(1234))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), 12345678901L, new PlcDATE(LocalDate.ofEpochDay(12345678901L))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), (float) 1234.56, new PlcDATE(LocalDate.ofEpochDay(1234))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), (double) 1234.56, new PlcDATE(LocalDate.ofEpochDay(1234))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), BigInteger.valueOf(12345678901L), new PlcDATE(LocalDate.ofEpochDay(12345678901L))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), BigDecimal.valueOf(12345678901L), new PlcDATE(LocalDate.ofEpochDay(12345678901L))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE), LocalDate.parse("+33803290-10-08"), new PlcDATE(LocalDate.ofEpochDay(12345678901L))),

            // LDATE values (Numeric values are interpreted as seconds since epoch)
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), LocalDate.ofEpochDay(1234), new PlcLDATE(LocalDate.ofEpochDay(1234))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), (byte) 123, new PlcLDATE(LocalDate.parse("1970-01-01"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), (short) 12345, new PlcLDATE(LocalDate.parse("1970-01-01"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), 1234567890, new PlcLDATE(LocalDate.parse("2009-02-13"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), 12345678901L, new PlcLDATE(LocalDate.parse("2361-03-21"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), (float) 123456.56, new PlcLDATE(LocalDate.parse("1970-01-02"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), (double) 12345678.9, new PlcLDATE(LocalDate.parse("1970-05-23"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), BigInteger.valueOf(12345678901L), new PlcLDATE(LocalDate.parse("2361-03-21"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), BigDecimal.valueOf(12345678901L), new PlcLDATE(LocalDate.parse("2361-03-21"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE), "1978-03-28", new PlcLDATE(LocalDate.parse("1978-03-28"))),

            // TIME_OF_DAY values (Numeric values are interpreted as seconds since midnight)
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), LocalTime.ofSecondOfDay(1234), new PlcTIME_OF_DAY(LocalTime.parse("00:20:34"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), (byte) 123, new PlcTIME_OF_DAY(LocalTime.parse("00:02:03"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), (short) 12345, new PlcTIME_OF_DAY(LocalTime.parse("03:25:45"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), 67890, new PlcTIME_OF_DAY(LocalTime.parse("18:51:30"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), 67890L, new PlcTIME_OF_DAY(LocalTime.parse("18:51:30"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), (float) 67890.56, new PlcTIME_OF_DAY(LocalTime.parse("18:51:30"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), (double) 67890.9, new PlcTIME_OF_DAY(LocalTime.parse("18:51:30"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), BigInteger.valueOf(67890L), new PlcTIME_OF_DAY(LocalTime.parse("18:51:30"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), BigDecimal.valueOf(67890L), new PlcTIME_OF_DAY(LocalTime.parse("18:51:30"))),
            Arguments.of(new MockTag("mock", PlcValueType.TIME_OF_DAY), "12:34:56", new PlcTIME_OF_DAY(LocalTime.parse("12:34:56"))),

            // LTIME_OF_DAY values (Numeric values are interpreted as milliseconds since midnight)
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), LocalTime.ofSecondOfDay(1234), new PlcLTIME_OF_DAY(LocalTime.parse("00:20:34"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), (byte) 123, new PlcLTIME_OF_DAY(LocalTime.parse("00:00:00.123"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), (short) 12345, new PlcLTIME_OF_DAY(LocalTime.parse("00:00:12.345"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), 12345678, new PlcLTIME_OF_DAY(LocalTime.parse("03:25:45.678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), 12345678L, new PlcLTIME_OF_DAY(LocalTime.parse("03:25:45.678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), (float) 123456.56, new PlcLTIME_OF_DAY(LocalTime.parse("00:02:03.456"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), (double) 12345678.9, new PlcLTIME_OF_DAY(LocalTime.parse("03:25:45.678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), BigInteger.valueOf(12345678L), new PlcLTIME_OF_DAY(LocalTime.parse("03:25:45.678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), BigDecimal.valueOf(12345678L), new PlcLTIME_OF_DAY(LocalTime.parse("03:25:45.678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LTIME_OF_DAY), "12:34:56", new PlcLTIME_OF_DAY(LocalTime.parse("12:34:56"))),

            // DATE_AND_TIME values (Numeric values are interpreted as seconds since epoch)
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), LocalDateTime.ofEpochSecond(1234, 0, ZoneOffset.UTC), new PlcDATE_AND_TIME(LocalDateTime.parse("1970-01-01T00:20:34"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), (byte) 123, new PlcDATE_AND_TIME(LocalDateTime.parse("1970-01-01T00:02:03"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), (short) 12345, new PlcDATE_AND_TIME(LocalDateTime.parse("1970-01-01T03:25:45"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), 12345678, new PlcDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), 12345678L, new PlcDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), (float) 123456.56, new PlcDATE_AND_TIME(LocalDateTime.parse("1970-01-02T10:17:36"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), (double) 12345678.9, new PlcDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), BigInteger.valueOf(12345678L), new PlcDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), BigDecimal.valueOf(12345678L), new PlcDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_TIME), "1978-03-28T12:34:56", new PlcDATE_AND_TIME(LocalDateTime.parse("1978-03-28T12:34:56"))),

            // DATE_AND_LTIME values (Numeric values are interpreted as milliseconds since epoch)
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), LocalDateTime.ofEpochSecond(1234, 5678, ZoneOffset.UTC), new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:20:34.000005678"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), (byte) 123, new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:00.000000123"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), (short) 12345, new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:00.000012345"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), 1234567890, new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:01.234567890"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), 12345678901L, new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:12.345678901"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), (float) 123456.56, new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:00.000123456"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), (double) 12345678.9, new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:00.012345678"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), BigInteger.valueOf(12345678901L), new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:12.345678901"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), BigDecimal.valueOf(12345678901L), new PlcDATE_AND_LTIME(LocalDateTime.parse("1970-01-01T00:00:12.345678901"))),
            Arguments.of(new MockTag("mock", PlcValueType.DATE_AND_LTIME), "1978-03-28T01:02:03", new PlcDATE_AND_LTIME(LocalDateTime.parse("1978-03-28T01:02:03"))),

            // LDATE_AND_TIME values
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), LocalDateTime.ofEpochSecond(1234, 5678, ZoneOffset.UTC), new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-01-01T00:20:34.000005678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), (byte) 123, new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-01-01T00:00:00.123"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), (short) 12345, new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-01-01T00:00:12.345"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), 1234567890, new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-01-15T06:56:07.890"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), 12345678901L, new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18.901"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), (float) 123456.56, new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-01-01T00:02:03.456"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), (double) 12345678.9, new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-01-01T03:25:45.678"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), BigInteger.valueOf(12345678901L), new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18.901"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), BigDecimal.valueOf(12345678901L), new PlcLDATE_AND_TIME(LocalDateTime.parse("1970-05-23T21:21:18.901"))),
            Arguments.of(new MockTag("mock", PlcValueType.LDATE_AND_TIME), "1978-03-28T01:02:03", new PlcLDATE_AND_TIME(LocalDateTime.parse("1978-03-28T01:02:03"))),

            // RAW_BYTE_ARRAY values
            Arguments.of(new MockTag("mock", PlcValueType.RAW_BYTE_ARRAY), new byte[]{(byte) 1, (byte) 2, (byte) 3, (byte) 4, (byte) 5, (byte) 6}, new PlcRawByteArray(new byte[]{(byte) 1, (byte) 2, (byte) 3, (byte) 4, (byte) 5, (byte) 6}))
        );
    }

    @ParameterizedTest
    @MethodSource("getSingleElementPlcValues")
    void testSingleElementPlcValues(PlcTag plcTag, Object input, PlcValue expected) {
        PlcValueHandler sut = new DefaultPlcValueHandler();

        PlcValue plcValue = sut.newPlcValue(plcTag, input);

        assertNotNull(plcValue);
        assertEquals(expected, plcValue);
    }

    private static Stream<Arguments> getPlcValueTypesRanges() {
        return Stream.of(
            Arguments.of(PlcBYTE.class),
            Arguments.of(PlcWORD.class),
            Arguments.of(PlcDWORD.class),
            // TODO: Need to find out how to test these ...
            //Arguments.of(PlcLWORD.class),
            Arguments.of(PlcUSINT.class),
            Arguments.of(PlcUINT.class),
            Arguments.of(PlcUDINT.class),
            // TODO: Need to find out how to test these ...
            //Arguments.of(PlcULINT.class),
            Arguments.of(PlcSINT.class),
            Arguments.of(PlcINT.class),
            Arguments.of(PlcDINT.class),
            // TODO: Need to find out how to test these ...
            //Arguments.of(PlcLINT.class),
            //Arguments.of(PlcREAL.class),
            //Arguments.of(PlcLREAL.class),
            Arguments.of(PlcCHAR.class),
            Arguments.of(PlcWCHAR.class)
        );
    }

    @ParameterizedTest
    @MethodSource("getPlcValueTypesRanges")
    void testPlcValueTypesRanges(Class<? extends PlcValue> plcValueType) throws Exception {
        Field minValueField = plcValueType.getField("MIN_VALUE");
        Number minValue = (Number) minValueField.get(null);
        Field maxValueField = plcValueType.getField("MAX_VALUE");
        Number maxValue = (Number) maxValueField.get(null);

        // Set value to one less than the lower bound (Should fail)
        try {
            plcValueType.getConstructor(Long.class).newInstance(minValue.longValue() - 1);
            fail("Set value to one less than the lower bound should have failed");
        } catch (Exception e) {
            // We want an exception here ...
        }

        // Set value to exactly the lower bound (Should succeed)
        try {
            PlcValue plcValue = plcValueType.getConstructor(Long.class).newInstance(minValue.longValue());
            assertEquals(minValue.longValue(), plcValue.getLong());
        } catch (Exception e) {
            fail("Set value to exactly the lower bound should have succeeded");
        }

        // Set value to one more than the lower bound (Should succeed)
        try {
            PlcValue plcValue = plcValueType.getConstructor(Long.class).newInstance(minValue.longValue() + 1);
            assertEquals(minValue.longValue() + 1, plcValue.getLong());
        } catch (Exception e) {
            fail("Set value to one more than the lower bound should have succeeded");
        }

        // Set value to one less than the upper bound (Should succeed)
        try {
            PlcValue plcValue = plcValueType.getConstructor(Long.class).newInstance(maxValue.longValue() - 1);
            assertEquals(maxValue.longValue() - 1, plcValue.getLong());
        } catch (Exception e) {
            fail("Set value to one less than the upper bound should have succeeded");
        }

        // Set value to exactly the upper bound (Should succeed)
        try {
            PlcValue plcValue = plcValueType.getConstructor(Long.class).newInstance(maxValue.longValue());
            assertEquals(maxValue.longValue(), plcValue.getLong());
        } catch (Exception e) {
            fail("Set value to exactly the upper bound should have succeeded");
        }

        // Set value to one more than the upper bound (Should fail)
        try {
            plcValueType.getConstructor(Long.class).newInstance(maxValue.longValue() + 1);
            fail("Set value to one more than the upper bound should have failed");
        } catch (Exception e) {
            // We want an exception here ...
        }
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