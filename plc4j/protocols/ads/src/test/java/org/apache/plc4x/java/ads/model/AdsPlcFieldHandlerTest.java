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
package org.apache.plc4x.java.ads.model;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.TestInfo;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.junit.jupiter.params.provider.ValueSource;

import java.lang.reflect.Array;
import java.lang.reflect.InvocationTargetException;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.function.BiFunction;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import java.util.stream.Stream;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.fail;

class AdsPlcFieldHandlerTest {

/*    private static AdsPlcFieldHandler SUT = new AdsPlcFieldHandler();

    private TestInfo testInfo;

    @BeforeEach
    void setUp(TestInfo testInfo) {
        this.testInfo = testInfo;
    }

    @ParameterizedTest
    @ValueSource(strings = {
        "0/0:BOOL",
        "0/1:WORD"
    })
    void createField(String fieldQuery) {
        SUT.createField(fieldQuery);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeOneBitTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "BOOL-BOOLEAN-MIN", "BOOL-BOOLEAN-MAX",
            "BOOL-BYTE-MIN", "BOOL-BYTE-NIL", "BOOL-BYTE-VAL", "BOOL-BYTE-MAX",
            "BOOL-SHORT-MIN", "BOOL-SHORT-NIL", "BOOL-SHORT-VAL", "BOOL-SHORT-MAX",
            "BOOL-INTEGER-MIN", "BOOL-INTEGER-NIL", "BOOL-INTEGER-VAL", "BOOL-INTEGER-MAX",
            "BOOL-LONG-MIN", "BOOL-LONG-NIL", "BOOL-LONG-VAL", "BOOL-LONG-MAX",

            "BYTE-BOOLEAN-MIN", "BYTE-BOOLEAN-MAX",
            "BYTE-BYTE-MIN", "BYTE-BYTE-NIL", "BYTE-BYTE-VAL", "BYTE-BYTE-MAX",
            "BYTE-SHORT-MIN", "BYTE-SHORT-NIL", "BYTE-SHORT-VAL", "BYTE-SHORT-MAX",
            "BYTE-INTEGER-MIN", "BYTE-INTEGER-NIL", "BYTE-INTEGER-VAL", "BYTE-INTEGER-MAX",
            "BYTE-LONG-MIN", "BYTE-LONG-NIL", "BYTE-LONG-VAL", "BYTE-LONG-MAX",

            "WORD-BOOLEAN-MIN", "WORD-BOOLEAN-MAX",
            "WORD-BYTE-MIN", "WORD-BYTE-NIL", "WORD-BYTE-VAL", "WORD-BYTE-MAX",
            "WORD-SHORT-MIN", "WORD-SHORT-NIL", "WORD-SHORT-VAL", "WORD-SHORT-MAX",
            "WORD-INTEGER-MIN", "WORD-INTEGER-NIL", "WORD-INTEGER-VAL", "WORD-INTEGER-MAX",
            "WORD-LONG-MIN", "WORD-LONG-NIL", "WORD-LONG-VAL", "WORD-LONG-MAX",

            "DWORD-BOOLEAN-MIN", "DWORD-BOOLEAN-MAX",
            "DWORD-BYTE-MIN", "DWORD-BYTE-NIL", "DWORD-BYTE-VAL", "DWORD-BYTE-MAX",
            "DWORD-SHORT-MIN", "DWORD-SHORT-NIL", "DWORD-SHORT-VAL", "DWORD-SHORT-MAX",
            "DWORD-INTEGER-MIN", "DWORD-INTEGER-NIL", "DWORD-INTEGER-VAL", "DWORD-INTEGER-MAX",
            "DWORD-LONG-MIN", "DWORD-LONG-NIL", "DWORD-LONG-VAL", "DWORD-LONG-MAX",

            "LWORD-BOOLEAN-MIN", "LWORD-BOOLEAN-MAX",
            "LWORD-BYTE-MIN", "LWORD-BYTE-NIL", "LWORD-BYTE-VAL", "LWORD-BYTE-MAX",
            "LWORD-SHORT-MIN", "LWORD-SHORT-NIL", "LWORD-SHORT-VAL", "LWORD-SHORT-MAX",
            "LWORD-INTEGER-MIN", "LWORD-INTEGER-NIL", "LWORD-INTEGER-VAL", "LWORD-INTEGER-MAX",
            "LWORD-LONG-MIN", "LWORD-LONG-NIL", "LWORD-LONG-VAL", "LWORD-LONG-MAX"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeBoolean);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeOneByteIntegerTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "BYTE-BYTE-NIL", "BYTE-BYTE-VAL", "BYTE-BYTE-MAX",
            "BYTE-SHORT-NIL", "BYTE-SHORT-VAL",
            "BYTE-INTEGER-NIL", "BYTE-INTEGER-VAL",
            "BYTE-LONG-NIL", "BYTE-LONG-VAL",

            "SINT-BYTE-MIN", "SINT-BYTE-NIL", "SINT-BYTE-VAL", "SINT-BYTE-MAX",
            "SINT-SHORT-NIL", "SINT-SHORT-VAL",
            "SINT-INTEGER-NIL", "SINT-INTEGER-VAL",
            "SINT-LONG-NIL", "SINT-LONG-VAL",

            "INT32-BYTE-MIN", "INT32-BYTE-NIL", "INT32-BYTE-VAL", "INT32-BYTE-MAX",
            "INT32-SHORT-MIN", "INT32-SHORT-NIL", "INT32-SHORT-VAL", "INT32-SHORT-MAX",
            "INT32-INTEGER-MIN", "INT32-INTEGER-NIL", "INT32-INTEGER-VAL", "INT32-INTEGER-MAX",
            "INT32-LONG-NIL", "INT32-LONG-VAL",

            "INT64-BYTE-MIN", "INT64-BYTE-NIL", "INT64-BYTE-VAL", "INT64-BYTE-MAX",
            "INT64-SHORT-MIN", "INT64-SHORT-NIL", "INT64-SHORT-VAL", "INT64-SHORT-MAX",
            "INT64-INTEGER-MIN", "INT64-INTEGER-NIL", "INT64-INTEGER-VAL", "INT64-INTEGER-MAX",
            "INT64-LONG-MIN", "INT64-LONG-NIL", "INT64-LONG-VAL", "INT64-LONG-MAX",

            "WORD-BYTE-NIL", "WORD-BYTE-VAL", "WORD-BYTE-MAX",
            "WORD-SHORT-NIL", "WORD-SHORT-VAL", "WORD-SHORT-MAX",
            "WORD-INTEGER-NIL", "WORD-INTEGER-VAL",
            "WORD-LONG-NIL", "WORD-LONG-VAL",

            "DWORD-BYTE-NIL", "DWORD-BYTE-VAL", "DWORD-BYTE-MAX",
            "DWORD-SHORT-NIL", "DWORD-SHORT-VAL", "DWORD-SHORT-MAX",
            "DWORD-INTEGER-NIL", "DWORD-INTEGER-VAL", "DWORD-INTEGER-MAX",
            "DWORD-LONG-NIL", "DWORD-LONG-VAL",

            "INT-BYTE-MIN", "INT-BYTE-NIL", "INT-BYTE-VAL", "INT-BYTE-MAX",
            "INT-SHORT-MIN", "INT-SHORT-NIL", "INT-SHORT-VAL", "INT-SHORT-MAX",
            "INT-INTEGER-NIL", "INT-INTEGER-VAL",
            "INT-LONG-NIL", "INT-LONG-VAL",

            "UINT-BYTE-NIL", "UINT-BYTE-VAL", "UINT-BYTE-MAX",
            "UINT-SHORT-NIL", "UINT-SHORT-VAL", "UINT-SHORT-MAX",
            "UINT-INTEGER-NIL", "UINT-INTEGER-VAL",
            "UINT-LONG-NIL", "UINT-LONG-VAL",

            "DINT-BYTE-MIN", "DINT-BYTE-NIL", "DINT-BYTE-VAL", "DINT-BYTE-MAX",
            "DINT-SHORT-MIN", "DINT-SHORT-NIL", "DINT-SHORT-VAL", "DINT-SHORT-MAX",
            "DINT-INTEGER-MIN", "DINT-INTEGER-NIL", "DINT-INTEGER-VAL", "DINT-INTEGER-MAX",
            "DINT-LONG-NIL", "DINT-LONG-VAL",

            "UDINT-BYTE-NIL", "UDINT-BYTE-VAL", "UDINT-BYTE-MAX",
            "UDINT-SHORT-NIL", "UDINT-SHORT-VAL", "UDINT-SHORT-MAX",
            "UDINT-INTEGER-NIL", "UDINT-INTEGER-VAL", "UDINT-INTEGER-MAX",
            "UDINT-LONG-NIL", "UDINT-LONG-VAL",

            "LINT-BYTE-MIN", "LINT-BYTE-NIL", "LINT-BYTE-VAL", "LINT-BYTE-MAX",
            "LINT-SHORT-MIN", "LINT-SHORT-NIL", "LINT-SHORT-VAL", "LINT-SHORT-MAX",
            "LINT-INTEGER-MIN", "LINT-INTEGER-NIL", "LINT-INTEGER-VAL", "LINT-INTEGER-MAX",
            "LINT-LONG-MIN", "LINT-LONG-NIL", "LINT-LONG-VAL", "LINT-LONG-MAX",

            "ULINT-BYTE-NIL", "ULINT-BYTE-VAL", "ULINT-BYTE-MAX",
            "ULINT-SHORT-NIL", "ULINT-SHORT-VAL", "ULINT-SHORT-MAX",
            "ULINT-INTEGER-NIL", "ULINT-INTEGER-VAL", "ULINT-INTEGER-MAX",
            "ULINT-LONG-NIL", "ULINT-LONG-VAL", "ULINT-LONG-MAX",

            // If it's a unsigned small int field, then any negative number is out of scope.
            "USINT-BYTE-NIL", "USINT-BYTE-VAL", "USINT-BYTE-MAX",
            "USINT-SHORT-NIL", "USINT-SHORT-VAL",
            "USINT-INTEGER-NIL", "USINT-INTEGER-VAL",
            "USINT-LONG-NIL", "USINT-LONG-VAL"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeByte);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeTwoByteIntegerTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "BYTE-BYTE-NIL", "BYTE-BYTE-VAL", "BYTE-BYTE-MAX",
            "BYTE-SHORT-NIL", "BYTE-SHORT-VAL",
            "BYTE-INTEGER-NIL", "BYTE-INTEGER-VAL",
            "BYTE-LONG-NIL", "BYTE-LONG-VAL",

            "WORD-BYTE-NIL", "WORD-BYTE-VAL", "WORD-BYTE-MAX",
            "WORD-SHORT-NIL", "WORD-SHORT-VAL", "WORD-SHORT-MAX",
            "WORD-INTEGER-NIL", "WORD-INTEGER-VAL",
            "WORD-LONG-NIL", "WORD-LONG-VAL",

            "DWORD-BYTE-NIL", "DWORD-BYTE-VAL", "DWORD-BYTE-MAX",
            "DWORD-SHORT-NIL", "DWORD-SHORT-VAL", "DWORD-SHORT-MAX",
            "DWORD-INTEGER-NIL", "DWORD-INTEGER-VAL", "DWORD-INTEGER-MAX",
            "DWORD-LONG-NIL", "DWORD-LONG-VAL",

            "INT-BYTE-MIN", "INT-BYTE-NIL", "INT-BYTE-VAL", "INT-BYTE-MAX",
            "INT-SHORT-MIN", "INT-SHORT-NIL", "INT-SHORT-VAL", "INT-SHORT-MAX",
            "INT-INTEGER-NIL", "INT-INTEGER-VAL",
            "INT-LONG-NIL", "INT-LONG-VAL",

            "SINT-BYTE-MIN", "SINT-BYTE-NIL", "SINT-BYTE-VAL", "SINT-BYTE-MAX",
            "SINT-SHORT-NIL", "SINT-SHORT-VAL",
            "SINT-INTEGER-NIL", "SINT-INTEGER-VAL",
            "SINT-LONG-NIL", "SINT-LONG-VAL",

            "USINT-BYTE-NIL", "USINT-BYTE-VAL", "USINT-BYTE-MAX",
            "USINT-SHORT-NIL", "USINT-SHORT-VAL",
            "USINT-INTEGER-NIL", "USINT-INTEGER-VAL",
            "USINT-LONG-NIL", "USINT-LONG-VAL",

            "DINT-BYTE-MIN", "DINT-BYTE-NIL", "DINT-BYTE-VAL", "DINT-BYTE-MAX",
            "DINT-SHORT-MIN", "DINT-SHORT-NIL", "DINT-SHORT-VAL", "DINT-SHORT-MAX",
            "DINT-INTEGER-MIN", "DINT-INTEGER-NIL", "DINT-INTEGER-VAL", "DINT-INTEGER-MAX",
            "DINT-LONG-NIL", "DINT-LONG-VAL",

            "UDINT-BYTE-NIL", "UDINT-BYTE-VAL", "UDINT-BYTE-MAX",
            "UDINT-SHORT-NIL", "UDINT-SHORT-VAL", "UDINT-SHORT-MAX",
            "UDINT-INTEGER-NIL", "UDINT-INTEGER-VAL", "UDINT-INTEGER-MAX",
            "UDINT-LONG-NIL", "UDINT-LONG-VAL",

            "LINT-BYTE-MIN", "LINT-BYTE-NIL", "LINT-BYTE-VAL", "LINT-BYTE-MAX",
            "LINT-SHORT-MIN", "LINT-SHORT-NIL", "LINT-SHORT-VAL", "LINT-SHORT-MAX",
            "LINT-INTEGER-MIN", "LINT-INTEGER-NIL", "LINT-INTEGER-VAL", "LINT-INTEGER-MAX",
            "LINT-LONG-MIN", "LINT-LONG-NIL", "LINT-LONG-VAL", "LINT-LONG-MAX",

            "ULINT-BYTE-NIL", "ULINT-BYTE-VAL", "ULINT-BYTE-MAX",
            "ULINT-SHORT-NIL", "ULINT-SHORT-VAL", "ULINT-SHORT-MAX",
            "ULINT-INTEGER-NIL", "ULINT-INTEGER-VAL", "ULINT-INTEGER-MAX",
            "ULINT-LONG-NIL", "ULINT-LONG-VAL", "ULINT-LONG-MAX",

            "INT32-BYTE-MIN", "INT32-BYTE-NIL", "INT32-BYTE-VAL", "INT32-BYTE-MAX",
            "INT32-SHORT-MIN", "INT32-SHORT-NIL", "INT32-SHORT-VAL", "INT32-SHORT-MAX",
            "INT32-INTEGER-MIN", "INT32-INTEGER-NIL", "INT32-INTEGER-VAL", "INT32-INTEGER-MAX",
            "INT32-LONG-NIL", "INT32-LONG-VAL",

            "INT64-BYTE-MIN", "INT64-BYTE-NIL", "INT64-BYTE-VAL", "INT64-BYTE-MAX",
            "INT64-SHORT-MIN", "INT64-SHORT-NIL", "INT64-SHORT-VAL", "INT64-SHORT-MAX",
            "INT64-INTEGER-MIN", "INT64-INTEGER-NIL", "INT64-INTEGER-VAL", "INT64-INTEGER-MAX",
            "INT64-LONG-MIN", "INT64-LONG-NIL", "INT64-LONG-VAL", "INT64-LONG-MAX",

            "UINT-BYTE-NIL", "UINT-BYTE-VAL", "UINT-BYTE-MAX",
            "UINT-SHORT-NIL", "UINT-SHORT-VAL", "UINT-SHORT-MAX",
            "UINT-INTEGER-NIL", "UINT-INTEGER-VAL",
            "UINT-LONG-NIL", "UINT-LONG-VAL"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeShort);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeFourByteIntegerTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "BYTE-BYTE-NIL", "BYTE-BYTE-VAL", "BYTE-BYTE-MAX",
            "BYTE-SHORT-NIL", "BYTE-SHORT-VAL",
            "BYTE-INTEGER-NIL", "BYTE-INTEGER-VAL",
            "BYTE-LONG-NIL", "BYTE-LONG-VAL",

            "WORD-BYTE-NIL", "WORD-BYTE-VAL", "WORD-BYTE-MAX",
            "WORD-SHORT-NIL", "WORD-SHORT-VAL", "WORD-SHORT-MAX",
            "WORD-INTEGER-NIL", "WORD-INTEGER-VAL",
            "WORD-LONG-NIL", "WORD-LONG-VAL",

            "SINT-BYTE-MIN", "SINT-BYTE-NIL", "SINT-BYTE-VAL", "SINT-BYTE-MAX",
            "SINT-SHORT-NIL", "SINT-SHORT-VAL",
            "SINT-INTEGER-NIL", "SINT-INTEGER-VAL",
            "SINT-LONG-NIL", "SINT-LONG-VAL",

            "USINT-BYTE-NIL", "USINT-BYTE-VAL", "USINT-BYTE-MAX",
            "USINT-SHORT-NIL", "USINT-SHORT-VAL",
            "USINT-INTEGER-NIL", "USINT-INTEGER-VAL",
            "USINT-LONG-NIL", "USINT-LONG-VAL",

            "INT-BYTE-MIN", "INT-BYTE-NIL", "INT-BYTE-VAL", "INT-BYTE-MAX",
            "INT-SHORT-MIN", "INT-SHORT-VAL", "INT-SHORT-NIL", "INT-SHORT-MAX",
            "INT-INTEGER-NIL", "INT-INTEGER-VAL",
            "INT-LONG-NIL", "INT-LONG-VAL",

            "UINT-BYTE-NIL", "UINT-BYTE-VAL", "UINT-BYTE-MAX",
            "UINT-SHORT-NIL", "UINT-SHORT-VAL", "UINT-SHORT-MAX",
            "UINT-INTEGER-NIL", "UINT-INTEGER-VAL",
            "UINT-LONG-NIL", "UINT-LONG-VAL",

            "LINT-BYTE-MIN", "LINT-BYTE-NIL", "LINT-BYTE-VAL", "LINT-BYTE-MAX",
            "LINT-SHORT-MIN", "LINT-SHORT-NIL", "LINT-SHORT-VAL", "LINT-SHORT-MAX",
            "LINT-INTEGER-MIN", "LINT-INTEGER-NIL", "LINT-INTEGER-VAL", "LINT-INTEGER-MAX",
            "LINT-LONG-MIN", "LINT-LONG-NIL", "LINT-LONG-VAL", "LINT-LONG-MAX",

            "ULINT-BYTE-NIL", "ULINT-BYTE-VAL", "ULINT-BYTE-MAX",
            "ULINT-SHORT-NIL", "ULINT-SHORT-VAL", "ULINT-SHORT-MAX",
            "ULINT-INTEGER-NIL", "ULINT-INTEGER-VAL", "ULINT-INTEGER-MAX",
            "ULINT-LONG-NIL", "ULINT-LONG-VAL", "ULINT-LONG-MAX",

            "DWORD-BYTE-NIL", "DWORD-BYTE-VAL", "DWORD-BYTE-MAX",
            "DWORD-SHORT-NIL", "DWORD-SHORT-VAL", "DWORD-SHORT-MAX",
            "DWORD-INTEGER-NIL", "DWORD-INTEGER-VAL", "DWORD-INTEGER-MAX",
            "DWORD-LONG-NIL", "DWORD-LONG-VAL",

            "DINT-BYTE-MIN", "DINT-BYTE-NIL", "DINT-BYTE-VAL", "DINT-BYTE-MAX",
            "DINT-SHORT-MIN", "DINT-SHORT-NIL", "DINT-SHORT-VAL", "DINT-SHORT-MAX",
            "DINT-INTEGER-MIN", "DINT-INTEGER-NIL", "DINT-INTEGER-VAL", "DINT-INTEGER-MAX",
            "DINT-LONG-NIL", "DINT-LONG-VAL",

            "UDINT-BYTE-NIL", "UDINT-BYTE-VAL", "UDINT-BYTE-MAX",
            "UDINT-SHORT-NIL", "UDINT-SHORT-VAL", "UDINT-SHORT-MAX",
            "UDINT-INTEGER-NIL", "UDINT-INTEGER-VAL", "UDINT-INTEGER-MAX",
            "UDINT-LONG-NIL", "UDINT-LONG-VAL",

            "INT32-BYTE-MIN", "INT32-BYTE-NIL", "INT32-BYTE-VAL", "INT32-BYTE-MAX",
            "INT32-SHORT-MIN", "INT32-SHORT-NIL", "INT32-SHORT-VAL", "INT32-SHORT-MAX",
            "INT32-INTEGER-MIN", "INT32-INTEGER-NIL", "INT32-INTEGER-VAL", "INT32-INTEGER-MAX",
            "INT32-LONG-NIL", "INT32-LONG-VAL",

            "INT64-BYTE-MIN", "INT64-BYTE-NIL", "INT64-BYTE-VAL", "INT64-BYTE-MAX",
            "INT64-SHORT-MIN", "INT64-SHORT-NIL", "INT64-SHORT-VAL", "INT64-SHORT-MAX",
            "INT64-INTEGER-MIN", "INT64-INTEGER-NIL", "INT64-INTEGER-VAL", "INT64-INTEGER-MAX",
            "INT64-LONG-MIN", "INT64-LONG-NIL", "INT64-LONG-VAL", "INT64-LONG-MAX"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeInteger);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeEightByteIntegerTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "BYTE-BYTE-NIL", "BYTE-BYTE-VAL", "BYTE-BYTE-MAX",
            "BYTE-SHORT-NIL", "BYTE-SHORT-VAL",
            "BYTE-INTEGER-NIL", "BYTE-INTEGER-VAL",
            "BYTE-LONG-NIL", "BYTE-LONG-VAL",

            "WORD-BYTE-NIL", "WORD-BYTE-VAL", "WORD-BYTE-MAX",
            "WORD-SHORT-NIL", "WORD-SHORT-VAL", "WORD-SHORT-MAX",
            "WORD-INTEGER-NIL", "WORD-INTEGER-VAL",
            "WORD-LONG-NIL", "WORD-LONG-VAL",

            "DWORD-BYTE-NIL", "DWORD-BYTE-VAL", "DWORD-BYTE-MAX",
            "DWORD-SHORT-NIL", "DWORD-SHORT-VAL", "DWORD-SHORT-MAX",
            "DWORD-INTEGER-NIL", "DWORD-INTEGER-VAL", "DWORD-INTEGER-MAX",
            "DWORD-LONG-NIL", "DWORD-LONG-VAL",

            "SINT-BYTE-MIN", "SINT-BYTE-NIL", "SINT-BYTE-VAL", "SINT-BYTE-MAX",
            "SINT-SHORT-NIL", "SINT-SHORT-VAL",
            "SINT-INTEGER-NIL", "SINT-INTEGER-VAL",
            "SINT-LONG-NIL", "SINT-LONG-VAL",

            "USINT-BYTE-NIL", "USINT-BYTE-VAL", "USINT-BYTE-MAX",
            "USINT-SHORT-NIL", "USINT-SHORT-VAL",
            "USINT-INTEGER-NIL", "USINT-INTEGER-VAL",
            "USINT-LONG-NIL", "USINT-LONG-VAL",

            "INT-BYTE-MIN", "INT-BYTE-NIL", "INT-BYTE-VAL", "INT-BYTE-MAX",
            "INT-SHORT-MIN", "INT-SHORT-NIL", "INT-SHORT-VAL", "INT-SHORT-MAX",
            "INT-INTEGER-NIL", "INT-INTEGER-VAL",
            "INT-LONG-NIL", "INT-LONG-VAL",

            "LINT-BYTE-MIN", "LINT-BYTE-NIL", "LINT-BYTE-VAL", "LINT-BYTE-MAX",
            "LINT-SHORT-MIN", "LINT-SHORT-NIL", "LINT-SHORT-VAL", "LINT-SHORT-MAX",
            "LINT-INTEGER-MIN", "LINT-INTEGER-NIL", "LINT-INTEGER-VAL", "LINT-INTEGER-MAX",
            "LINT-LONG-MIN", "LINT-LONG-NIL", "LINT-LONG-VAL", "LINT-LONG-MAX",

            "UINT-BYTE-NIL", "UINT-BYTE-VAL", "UINT-BYTE-MAX",
            "UINT-SHORT-NIL", "UINT-SHORT-VAL", "UINT-SHORT-MAX",
            "UINT-INTEGER-NIL", "UINT-INTEGER-VAL",
            "UINT-LONG-NIL", "UINT-LONG-VAL",

            "DINT-BYTE-MIN", "DINT-BYTE-NIL", "DINT-BYTE-VAL", "DINT-BYTE-MAX",
            "DINT-SHORT-MIN", "DINT-SHORT-NIL", "DINT-SHORT-VAL", "DINT-SHORT-MAX",
            "DINT-INTEGER-MIN", "DINT-INTEGER-NIL", "DINT-INTEGER-VAL", "DINT-INTEGER-MAX",
            "DINT-LONG-NIL", "DINT-LONG-VAL",

            "UDINT-BYTE-NIL", "UDINT-BYTE-VAL", "UDINT-BYTE-MAX",
            "UDINT-SHORT-NIL", "UDINT-SHORT-VAL", "UDINT-SHORT-MAX",
            "UDINT-INTEGER-NIL", "UDINT-INTEGER-VAL", "UDINT-INTEGER-MAX",
            "UDINT-LONG-NIL", "UDINT-LONG-VAL",

            "INT32-BYTE-MIN", "INT32-BYTE-NIL", "INT32-BYTE-VAL", "INT32-BYTE-MAX",
            "INT32-SHORT-MIN", "INT32-SHORT-NIL", "INT32-SHORT-VAL", "INT32-SHORT-MAX",
            "INT32-INTEGER-MIN", "INT32-INTEGER-NIL", "INT32-INTEGER-VAL", "INT32-INTEGER-MAX",
            "INT32-LONG-NIL", "INT32-LONG-VAL",

            "INT64-BYTE-MIN", "INT64-BYTE-NIL", "INT64-BYTE-VAL", "INT64-BYTE-MAX",
            "INT64-SHORT-MIN", "INT64-SHORT-NIL", "INT64-SHORT-VAL", "INT64-SHORT-MAX",
            "INT64-INTEGER-MIN", "INT64-INTEGER-NIL", "INT64-INTEGER-VAL", "INT64-INTEGER-MAX",
            "INT64-LONG-MIN", "INT64-LONG-NIL", "INT64-LONG-VAL", "INT64-LONG-MAX",

            "ULINT-BYTE-NIL", "ULINT-BYTE-VAL", "ULINT-BYTE-MAX",
            "ULINT-SHORT-NIL", "ULINT-SHORT-VAL", "ULINT-SHORT-MAX",
            "ULINT-INTEGER-NIL", "ULINT-INTEGER-VAL", "ULINT-INTEGER-MAX",
            "ULINT-LONG-NIL", "ULINT-LONG-VAL", "ULINT-LONG-MAX",

            "LWORD-BYTE-MIN", "LWORD-BYTE-NIL", "LWORD-BYTE-VAL", "LWORD-BYTE-MAX",
            "LWORD-SHORT-MIN", "LWORD-SHORT-NIL", "LWORD-SHORT-VAL", "LWORD-SHORT-MAX",
            "LWORD-INTEGER-MIN", "LWORD-INTEGER-NIL", "LWORD-INTEGER-VAL", "LWORD-INTEGER-MAX",
            "LWORD-LONG-MIN", "LWORD-LONG-NIL", "LWORD-LONG-VAL", "LWORD-LONG-MAX"

            // TODO: Somehow test ULWORD too ...
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeLong);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeFourByteFloatTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "REAL-FLOAT-MIN",
            "REAL-FLOAT-NIL",
            "REAL-FLOAT-VAL",
            "REAL-FLOAT-MAX",
            "REAL-DOUBLE-NIL",
            "REAL-DOUBLE-VAL",
            "LREAL-FLOAT-MIN",
            "LREAL-FLOAT-NIL",
            "LREAL-FLOAT-VAL",
            "LREAL-FLOAT-MAX",
            "LREAL-DOUBLE-MIN",
            "LREAL-DOUBLE-NIL",
            "LREAL-DOUBLE-VAL",
            "LREAL-DOUBLE-MAX"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeFloat);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeEightByteFloatTypes(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "REAL-FLOAT-MIN", "REAL-FLOAT-NIL", "REAL-FLOAT-VAL", "REAL-FLOAT-MAX",
            "REAL-DOUBLE-NIL", "REAL-DOUBLE-VAL",

            "LREAL-FLOAT-MIN", "LREAL-FLOAT-NIL", "LREAL-FLOAT-VAL", "LREAL-FLOAT-MAX",
            "LREAL-DOUBLE-MIN", "LREAL-DOUBLE-NIL", "LREAL-DOUBLE-VAL", "LREAL-DOUBLE-MAX"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeDouble);
    }

    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeString(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
            "CHAR-BYTE-MIN", "CHAR-BYTE-NIL", "CHAR-BYTE-VAL", "CHAR-BYTE-MAX",
            "CHAR-SHORT-MIN", "CHAR-SHORT-NIL", "CHAR-SHORT-VAL", "CHAR-SHORT-MAX",
            "CHAR-INTEGER-MIN", "CHAR-INTEGER-NIL", "CHAR-INTEGER-VAL", "CHAR-INTEGER-MAX",
            "CHAR-LONG-MIN", "CHAR-LONG-NIL", "CHAR-LONG-VAL", "CHAR-LONG-MAX",
            "CHAR-STRING-MIN",
            "WCHAR-BYTE-MIN", "WCHAR-BYTE-NIL", "WCHAR-BYTE-VAL", "WCHAR-BYTE-MAX",
            "WCHAR-SHORT-MIN", "WCHAR-SHORT-NIL", "WCHAR-SHORT-VAL", "WCHAR-SHORT-MAX",
            "WCHAR-INTEGER-MIN", "WCHAR-INTEGER-NIL", "WCHAR-INTEGER-VAL", "WCHAR-INTEGER-MAX",
            "WCHAR-LONG-MIN", "WCHAR-LONG-NIL", "WCHAR-LONG-VAL", "WCHAR-LONG-MAX",
            "WCHAR-STRING-MIN",
            "STRING-BYTE-MIN", "STRING-BYTE-NIL", "STRING-BYTE-VAL", "STRING-BYTE-MAX",
            "STRING-SHORT-MIN", "STRING-SHORT-NIL", "STRING-SHORT-VAL", "STRING-SHORT-MAX",
            "STRING-INTEGER-MIN", "STRING-INTEGER-NIL", "STRING-INTEGER-VAL", "STRING-INTEGER-MAX",
            "STRING-LONG-MIN", "STRING-LONG-NIL", "STRING-LONG-VAL", "STRING-LONG-MAX",
            "STRING-STRING-MIN", "STRING-STRING-VAL", "STRING-STRING-MAX",
            "WSTRING-BYTE-MIN", "WSTRING-BYTE-NIL", "WSTRING-BYTE-VAL", "WSTRING-BYTE-MAX",
            "WSTRING-SHORT-MIN", "WSTRING-SHORT-NIL", "WSTRING-SHORT-VAL", "WSTRING-SHORT-MAX",
            "WSTRING-INTEGER-MIN", "WSTRING-INTEGER-NIL", "WSTRING-INTEGER-VAL", "WSTRING-INTEGER-MAX",
            "WSTRING-LONG-MIN", "WSTRING-LONG-NIL", "WSTRING-LONG-VAL", "WSTRING-LONG-MAX",
            "WSTRING-STRING-MIN", "WSTRING-STRING-VAL", "WSTRING-STRING-MAX"
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeString);
    }

    @Disabled("Not implemented yet")
    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeTime(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeTime);
    }

    @Disabled("Not implemented yet")
    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeDate(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeDate);
    }

    @Disabled("Not implemented yet")
    @ParameterizedTest
    @MethodSource("createInputArrays")
    void encodeDateTime(String name, PlcField field, Object[] values) {
        Set<String> expectedSuccess = new HashSet<>(Arrays.asList(
        ));
        encode(name, field, values, expectedSuccess, SUT::encodeDateTime);
    }

    private static Stream<Arguments> createInputArrays() {
        // Generate valid fields for each ads type.
        Map<AdsDataType, PlcField> fields = new HashMap<>();
        Arrays.stream(AdsDataType.values()).forEach(adsDataType -> {
            String fieldQuery = "1/2:" + adsDataType.name();
            fields.put(adsDataType, SUT.createField(fieldQuery));
        });
        // Generate output for each combination of ads and Java type.
        Stream<Arguments> values = null;
        for (AdsDataType adsDataType : AdsDataType.values()) {
            PlcField field = fields.get(adsDataType);
            for (InputTypes javaType : InputTypes.values()) {
                Object[] testValues = javaType.values;

                BaseDefaultFieldItem fieldItem;
                try {
                    fieldItem = javaType.fieldItemType.getDeclaredConstructor(Array.newInstance(testValues[0].getClass(), 0).getClass()).newInstance(new Object[]{testValues});
                } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException e) {
                    throw new PlcRuntimeException("Error initializing field class " + javaType.fieldItemType.getSimpleName(), e);
                }

                Stream<Arguments> curValues;
                // Min, Max
                if (testValues.length == 2) {
                    curValues = Stream.of(
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-MIN", field, createOneElementArray(testValues[0])),
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-MAX", field, createOneElementArray(testValues[1])));
                }
                // Value, Min, Max
                else if (testValues.length == 3) {
                    curValues = Stream.of(
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-MIN", field, createOneElementArray(testValues[1])),
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-VAL", field, createOneElementArray(testValues[0])),
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-MAX", field, createOneElementArray(testValues[2])));
                }
                // Zero, Value, Min, Max
                else if (testValues.length == 4) {
                    curValues = Stream.of(
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-MIN", field, createOneElementArray(testValues[2])),
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-NIL", field, createOneElementArray(testValues[0])),
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-VAL", field, createOneElementArray(testValues[1])),
                        Arguments.of(adsDataType.name() + "-" + javaType.name() + "-MAX", field, createOneElementArray(testValues[3])));
                } else {
                    throw new RuntimeException("Expecting 2, 3 or 4 valued test-input");
                }
                if (values == null) {
                    values = curValues;
                } else {
                    values = Stream.concat(values, curValues);
                }
            }
        }
        return values;
    }

    private static Object[] createOneElementArray(Object value) {
        Class<?> type = value.getClass();
        Object[] array = (Object[]) Array.newInstance(type, 1);
        array[0] = value;
        return array;
    }

    enum InputTypes {
        BOOLEAN(DefaultBooleanFieldItem.class, new Boolean[]{false, true}),
        BYTE(DefaultByteFieldItem.class, new Byte[]{(byte) 0, (byte) 42, Byte.MIN_VALUE, Byte.MAX_VALUE}),
        SHORT(DefaultShortFieldItem.class, new Short[]{(short) 0, (short) 42, Short.MIN_VALUE, Short.MAX_VALUE}),
        INTEGER(DefaultIntegerFieldItem.class, new Integer[]{0, 42, Integer.MIN_VALUE, Integer.MAX_VALUE}),
        LONG(DefaultLongFieldItem.class, new Long[]{(long) 0, (long) 42, Long.MIN_VALUE, Long.MAX_VALUE}),
        FLOAT(DefaultFloatFieldItem.class, new Float[]{(float) 0.0, (float) 42.23, -Float.MAX_VALUE, Float.MAX_VALUE}),
        DOUBLE(DefaultDoubleFieldItem.class, new Double[]{0.0, 42.23, -Double.MAX_VALUE, Double.MAX_VALUE}),
        // Creates an empty sting as min and a 254 char long string as max.
        STRING(DefaultStringFieldItem.class, new String[]{"Hurz", "", IntStream.range(0, 127).mapToObj(i -> "a").collect(Collectors.joining(""))}),
        TIME(DefaultLocalTimeFieldItem.class, new LocalTime[]{LocalTime.now(), LocalTime.MIDNIGHT, LocalTime.MIN, LocalTime.MAX}),
        DATE(DefaultLocalDateFieldItem.class, new LocalDate[]{LocalDate.now(), LocalDate.MIN, LocalDate.MAX}),
        DATETIME(DefaultLocalDateTimeFieldItem.class, new LocalDateTime[]{LocalDateTime.now(), LocalDateTime.MIN, LocalDateTime.MAX});

        private final Class<? extends BaseDefaultFieldItem> fieldItemType;
        private final Object[] values;

        InputTypes(Class<? extends BaseDefaultFieldItem> fieldItemType, Object[] values) {
            this.fieldItemType = fieldItemType;
            this.values = values;
        }
    }

    private void encode(String name, PlcField field, Object[] values, Set<String> expectedSuccess,
                        BiFunction<PlcField, Object[], BaseDefaultFieldItem> encoder) {
        boolean success = expectedSuccess.contains(name);
        try {
            BaseDefaultFieldItem fieldItem = encoder.apply(field, values);
            assertNotNull(fieldItem, "A FieldItem instance should have been returned for testcase " + name);
            if (!success) {
                fail(testInfo.getDisplayName() + ".\nExpected to fail for testcase " + name);
            }
        } catch (Exception e) {
            if (success) {
                fail(testInfo.getDisplayName() + ".\nExpected to succeed for testcase " + name + " got exception " + e.getMessage(), e);
            }
        }
    }
*/
}