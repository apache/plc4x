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

package org.apache.plc4x.java.tools.eventpump;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.metadata.Metadata;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcBYTE;
import org.apache.plc4x.java.spi.values.PlcDATE;
import org.apache.plc4x.java.spi.values.PlcDATE_AND_TIME;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcLREAL;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcTIME_OF_DAY;
import org.apache.plc4x.java.spi.values.PlcULINT;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.*;

class TransformedPlcReadResponseTest {

    private PlcReadResponse mockDelegate;
    private Map<String, PlcValue> transformedValues;
    private TransformedPlcReadResponse transformedResponse;

    @BeforeEach
    void setUp() {
        mockDelegate = mock(PlcReadResponse.class);
        transformedValues = new HashMap<>();
    }

    @Test
    void testGetRequest() {
        PlcReadRequest mockRequest = mock(PlcReadRequest.class);
        when(mockDelegate.getRequest()).thenReturn(mockRequest);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(mockRequest, transformedResponse.getRequest());
        verify(mockDelegate).getRequest();
    }

    @Test
    void testGetTagNames() {
        Collection<String> tagNames = Arrays.asList("tag1", "tag2");
        when(mockDelegate.getTagNames()).thenReturn(tagNames);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(tagNames, transformedResponse.getTagNames());
        verify(mockDelegate).getTagNames();
    }

    @Test
    void testGetTag() {
        PlcTag mockTag = mock(PlcTag.class);
        when(mockDelegate.getTag("tag1")).thenReturn(mockTag);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(mockTag, transformedResponse.getTag("tag1"));
        verify(mockDelegate).getTag("tag1");
    }

    @Test
    void testGetResponseCode() {
        when(mockDelegate.getResponseCode("tag1")).thenReturn(PlcResponseCode.OK);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(PlcResponseCode.OK, transformedResponse.getResponseCode("tag1"));
        verify(mockDelegate).getResponseCode("tag1");
    }

    @Test
    void testGetTagMetadata() {
        Metadata mockMetadata = mock(Metadata.class);
        when(mockDelegate.getTagMetadata("tag1")).thenReturn(mockMetadata);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(mockMetadata, transformedResponse.getTagMetadata("tag1"));
        verify(mockDelegate).getTagMetadata("tag1");
    }

    @Test
    void testGetPlcValueWithTransformation() {
        PlcValue originalValue = new PlcDINT(10);
        PlcValue transformedValue = new PlcDINT(20);

        when(mockDelegate.getPlcValue("tag1")).thenReturn(originalValue);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(transformedValue, transformedResponse.getPlcValue("tag1"));
        verify(mockDelegate, never()).getPlcValue("tag1");
    }

    @Test
    void testGetPlcValueWithoutTransformation() {
        PlcValue originalValue = new PlcDINT(10);

        when(mockDelegate.getPlcValue("tag1")).thenReturn(originalValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(originalValue, transformedResponse.getPlcValue("tag1"));
        verify(mockDelegate).getPlcValue("tag1");
    }

    @Test
    void testGetAsPlcValue() {
        PlcValue mockValue = mock(PlcValue.class);
        when(mockDelegate.getAsPlcValue()).thenReturn(mockValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(mockValue, transformedResponse.getAsPlcValue());
        verify(mockDelegate).getAsPlcValue();
    }

    @Test
    void testGetNumberOfValues() {
        when(mockDelegate.getNumberOfValues("tag1")).thenReturn(5);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(5, transformedResponse.getNumberOfValues("tag1"));
        verify(mockDelegate).getNumberOfValues("tag1");
    }

    @Test
    void testGetObjectWithTransformation() {
        PlcValue transformedValue = new PlcDINT(42);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(42, transformedResponse.getObject("tag1"));
    }

    @Test
    void testGetObjectWithoutTransformation() {
        when(mockDelegate.getObject("tag1")).thenReturn(42);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(42, transformedResponse.getObject("tag1"));
        verify(mockDelegate).getObject("tag1");
    }

    @Test
    void testGetObjectWithIndex() {
        when(mockDelegate.getObject("tag1", 0)).thenReturn(42);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(42, transformedResponse.getObject("tag1", 0));
        verify(mockDelegate).getObject("tag1", 0);
    }

    @Test
    void testGetAllObjects() {
        Collection<Object> objects = Arrays.asList(1, 2, 3);
        when(mockDelegate.getAllObjects("tag1")).thenReturn(objects);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(objects, transformedResponse.getAllObjects("tag1"));
        verify(mockDelegate).getAllObjects("tag1");
    }

    // Boolean tests
    @Test
    void testIsValidBoolean() {
        when(mockDelegate.isValidBoolean("tag1")).thenReturn(true);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.isValidBoolean("tag1"));
        verify(mockDelegate).isValidBoolean("tag1");
    }

    @Test
    void testIsValidBooleanWithIndex() {
        when(mockDelegate.isValidBoolean("tag1", 0)).thenReturn(true);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.isValidBoolean("tag1", 0));
        verify(mockDelegate).isValidBoolean("tag1", 0);
    }

    @Test
    void testGetBooleanWithTransformation() {
        PlcValue transformedValue = new PlcBOOL(true);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.getBoolean("tag1"));
    }

    @Test
    void testGetBooleanWithoutTransformation() {
        when(mockDelegate.getBoolean("tag1")).thenReturn(true);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.getBoolean("tag1"));
        verify(mockDelegate).getBoolean("tag1");
    }

    @Test
    void testGetBooleanWithIndex() {
        when(mockDelegate.getBoolean("tag1", 0)).thenReturn(true);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.getBoolean("tag1", 0));
        verify(mockDelegate).getBoolean("tag1", 0);
    }

    @Test
    void testGetAllBooleans() {
        Collection<Boolean> booleans = Arrays.asList(true, false);
        when(mockDelegate.getAllBooleans("tag1")).thenReturn(booleans);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(booleans, transformedResponse.getAllBooleans("tag1"));
        verify(mockDelegate).getAllBooleans("tag1");
    }

    // Byte tests
    @Test
    void testIsValidByte() {
        when(mockDelegate.isValidByte("tag1")).thenReturn(true);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.isValidByte("tag1"));
    }

    @Test
    void testIsValidByteWithIndex() {
        when(mockDelegate.isValidByte("tag1", 0)).thenReturn(true);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertTrue(transformedResponse.isValidByte("tag1", 0));
    }

    @Test
    void testGetByteWithTransformation() {
        PlcValue transformedValue = new PlcBYTE((byte) 42);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals((byte) 42, transformedResponse.getByte("tag1"));
    }

    @Test
    void testGetByteWithoutTransformation() {
        when(mockDelegate.getByte("tag1")).thenReturn((byte) 42);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals((byte) 42, transformedResponse.getByte("tag1"));
    }

    @Test
    void testGetByteWithIndex() {
        when(mockDelegate.getByte("tag1", 0)).thenReturn((byte) 42);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals((byte) 42, transformedResponse.getByte("tag1", 0));
    }

    @Test
    void testGetAllBytes() {
        Collection<Byte> bytes = Arrays.asList((byte) 1, (byte) 2);
        when(mockDelegate.getAllBytes("tag1")).thenReturn(bytes);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(bytes, transformedResponse.getAllBytes("tag1"));
    }

    // Short tests
    @Test
    void testGetShortWithTransformation() {
        PlcValue transformedValue = new PlcINT((short) 100);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals((short) 100, transformedResponse.getShort("tag1"));
    }

    @Test
    void testGetAllShorts() {
        Collection<Short> shorts = Arrays.asList((short) 1, (short) 2);
        when(mockDelegate.getAllShorts("tag1")).thenReturn(shorts);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(shorts, transformedResponse.getAllShorts("tag1"));
    }

    // Integer tests
    @Test
    void testGetIntegerWithTransformation() {
        PlcValue transformedValue = new PlcDINT(1000);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(1000, transformedResponse.getInteger("tag1"));
    }

    @Test
    void testGetAllIntegers() {
        Collection<Integer> integers = Arrays.asList(1, 2, 3);
        when(mockDelegate.getAllIntegers("tag1")).thenReturn(integers);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(integers, transformedResponse.getAllIntegers("tag1"));
    }

    // Long tests
    @Test
    void testGetLongWithTransformation() {
        PlcValue transformedValue = new PlcLINT(10000L);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(10000L, transformedResponse.getLong("tag1"));
    }

    @Test
    void testGetAllLongs() {
        Collection<Long> longs = Arrays.asList(1L, 2L, 3L);
        when(mockDelegate.getAllLongs("tag1")).thenReturn(longs);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(longs, transformedResponse.getAllLongs("tag1"));
    }

    // Float tests
    @Test
    void testGetFloatWithTransformation() {
        PlcValue transformedValue = new PlcREAL(3.14f);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(3.14f, transformedResponse.getFloat("tag1"), 0.001);
    }

    @Test
    void testGetAllFloats() {
        Collection<Float> floats = Arrays.asList(1.0f, 2.0f);
        when(mockDelegate.getAllFloats("tag1")).thenReturn(floats);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(floats, transformedResponse.getAllFloats("tag1"));
    }

    // Double tests
    @Test
    void testGetDoubleWithTransformation() {
        PlcValue transformedValue = new PlcLREAL(3.14159);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(3.14159, transformedResponse.getDouble("tag1"), 0.00001);
    }

    @Test
    void testGetAllDoubles() {
        Collection<Double> doubles = Arrays.asList(1.0, 2.0);
        when(mockDelegate.getAllDoubles("tag1")).thenReturn(doubles);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(doubles, transformedResponse.getAllDoubles("tag1"));
    }

    // BigInteger tests
    @Test
    void testGetBigIntegerWithTransformation() {
        PlcValue transformedValue = new PlcULINT(BigInteger.valueOf(123456789));
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(BigInteger.valueOf(123456789), transformedResponse.getBigInteger("tag1"));
    }

    @Test
    void testGetAllBigIntegers() {
        Collection<BigInteger> bigIntegers = Arrays.asList(BigInteger.ONE, BigInteger.TEN);
        when(mockDelegate.getAllBigIntegers("tag1")).thenReturn(bigIntegers);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(bigIntegers, transformedResponse.getAllBigIntegers("tag1"));
    }

    // BigDecimal tests
    @Test
    void testGetBigDecimalWithTransformation() {
        PlcValue transformedValue = new PlcLREAL(3.14159);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        BigDecimal result = transformedResponse.getBigDecimal("tag1");
        assertEquals(3.14159, result.doubleValue(), 0.00001);
    }

    @Test
    void testGetAllBigDecimals() {
        Collection<BigDecimal> bigDecimals = Arrays.asList(BigDecimal.ONE, BigDecimal.TEN);
        when(mockDelegate.getAllBigDecimals("tag1")).thenReturn(bigDecimals);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(bigDecimals, transformedResponse.getAllBigDecimals("tag1"));
    }

    // String tests
    @Test
    void testGetStringWithTransformation() {
        PlcValue transformedValue = new PlcSTRING("transformed");
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals("transformed", transformedResponse.getString("tag1"));
    }

    @Test
    void testGetAllStrings() {
        Collection<String> strings = Arrays.asList("a", "b", "c");
        when(mockDelegate.getAllStrings("tag1")).thenReturn(strings);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(strings, transformedResponse.getAllStrings("tag1"));
    }

    // Time tests
    @Test
    void testGetTimeWithTransformation() {
        LocalTime time = LocalTime.of(12, 30);
        PlcValue transformedValue = new PlcTIME_OF_DAY(time);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(time, transformedResponse.getTime("tag1"));
    }

    @Test
    void testGetAllTimes() {
        Collection<LocalTime> times = Arrays.asList(LocalTime.now(), LocalTime.NOON);
        when(mockDelegate.getAllTimes("tag1")).thenReturn(times);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(times, transformedResponse.getAllTimes("tag1"));
    }

    // Date tests
    @Test
    void testGetDateWithTransformation() {
        LocalDate date = LocalDate.of(2024, 10, 24);
        PlcValue transformedValue = new PlcDATE(date);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(date, transformedResponse.getDate("tag1"));
    }

    @Test
    void testGetAllDates() {
        Collection<LocalDate> dates = Arrays.asList(LocalDate.now(), LocalDate.of(2024, 1, 1));
        when(mockDelegate.getAllDates("tag1")).thenReturn(dates);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(dates, transformedResponse.getAllDates("tag1"));
    }

    // DateTime tests
    @Test
    void testGetDateTimeWithTransformation() {
        LocalDateTime dateTime = LocalDateTime.of(2024, 10, 24, 12, 30);
        PlcValue transformedValue = new PlcDATE_AND_TIME(dateTime);
        transformedValues.put("tag1", transformedValue);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(dateTime, transformedResponse.getDateTime("tag1"));
    }

    @Test
    void testGetAllDateTimes() {
        Collection<LocalDateTime> dateTimes = Arrays.asList(LocalDateTime.now(), LocalDateTime.of(2024, 1, 1, 0, 0));
        when(mockDelegate.getAllDateTimes("tag1")).thenReturn(dateTimes);

        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);

        assertEquals(dateTimes, transformedResponse.getAllDateTimes("tag1"));
    }

    // Validation methods
    @Test
    void testIsValidShort() {
        when(mockDelegate.isValidShort("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidShort("tag1"));
    }

    @Test
    void testIsValidInteger() {
        when(mockDelegate.isValidInteger("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidInteger("tag1"));
    }

    @Test
    void testIsValidBigInteger() {
        when(mockDelegate.isValidBigInteger("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidBigInteger("tag1"));
    }

    @Test
    void testIsValidLong() {
        when(mockDelegate.isValidLong("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidLong("tag1"));
    }

    @Test
    void testIsValidFloat() {
        when(mockDelegate.isValidFloat("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidFloat("tag1"));
    }

    @Test
    void testIsValidDouble() {
        when(mockDelegate.isValidDouble("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidDouble("tag1"));
    }

    @Test
    void testIsValidBigDecimal() {
        when(mockDelegate.isValidBigDecimal("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidBigDecimal("tag1"));
    }

    @Test
    void testIsValidString() {
        when(mockDelegate.isValidString("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidString("tag1"));
    }

    @Test
    void testIsValidTime() {
        when(mockDelegate.isValidTime("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidTime("tag1"));
    }

    @Test
    void testIsValidDate() {
        when(mockDelegate.isValidDate("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidDate("tag1"));
    }

    @Test
    void testIsValidDateTime() {
        when(mockDelegate.isValidDateTime("tag1")).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidDateTime("tag1"));
    }

    // Validation with index methods
    @Test
    void testIsValidShortWithIndex() {
        when(mockDelegate.isValidShort("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidShort("tag1", 0));
    }

    @Test
    void testIsValidIntegerWithIndex() {
        when(mockDelegate.isValidInteger("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidInteger("tag1", 0));
    }

    @Test
    void testIsValidBigIntegerWithIndex() {
        when(mockDelegate.isValidBigInteger("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidBigInteger("tag1", 0));
    }

    @Test
    void testIsValidLongWithIndex() {
        when(mockDelegate.isValidLong("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidLong("tag1", 0));
    }

    @Test
    void testIsValidFloatWithIndex() {
        when(mockDelegate.isValidFloat("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidFloat("tag1", 0));
    }

    @Test
    void testIsValidDoubleWithIndex() {
        when(mockDelegate.isValidDouble("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidDouble("tag1", 0));
    }

    @Test
    void testIsValidBigDecimalWithIndex() {
        when(mockDelegate.isValidBigDecimal("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidBigDecimal("tag1", 0));
    }

    @Test
    void testIsValidStringWithIndex() {
        when(mockDelegate.isValidString("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidString("tag1", 0));
    }

    @Test
    void testIsValidTimeWithIndex() {
        when(mockDelegate.isValidTime("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidTime("tag1", 0));
    }

    @Test
    void testIsValidDateWithIndex() {
        when(mockDelegate.isValidDate("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidDate("tag1", 0));
    }

    @Test
    void testIsValidDateTimeWithIndex() {
        when(mockDelegate.isValidDateTime("tag1", 0)).thenReturn(true);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertTrue(transformedResponse.isValidDateTime("tag1", 0));
    }

    // Get with index methods
    @Test
    void testGetShortWithIndex() {
        when(mockDelegate.getShort("tag1", 0)).thenReturn((short) 100);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals((short) 100, transformedResponse.getShort("tag1", 0));
    }

    @Test
    void testGetIntegerWithIndex() {
        when(mockDelegate.getInteger("tag1", 0)).thenReturn(1000);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(1000, transformedResponse.getInteger("tag1", 0));
    }

    @Test
    void testGetBigIntegerWithIndex() {
        when(mockDelegate.getBigInteger("tag1", 0)).thenReturn(BigInteger.TEN);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(BigInteger.TEN, transformedResponse.getBigInteger("tag1", 0));
    }

    @Test
    void testGetLongWithIndex() {
        when(mockDelegate.getLong("tag1", 0)).thenReturn(10000L);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(10000L, transformedResponse.getLong("tag1", 0));
    }

    @Test
    void testGetFloatWithIndex() {
        when(mockDelegate.getFloat("tag1", 0)).thenReturn(3.14f);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(3.14f, transformedResponse.getFloat("tag1", 0), 0.001);
    }

    @Test
    void testGetDoubleWithIndex() {
        when(mockDelegate.getDouble("tag1", 0)).thenReturn(3.14159);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(3.14159, transformedResponse.getDouble("tag1", 0), 0.00001);
    }

    @Test
    void testGetBigDecimalWithIndex() {
        when(mockDelegate.getBigDecimal("tag1", 0)).thenReturn(BigDecimal.TEN);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(BigDecimal.TEN, transformedResponse.getBigDecimal("tag1", 0));
    }

    @Test
    void testGetStringWithIndex() {
        when(mockDelegate.getString("tag1", 0)).thenReturn("test");
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals("test", transformedResponse.getString("tag1", 0));
    }

    @Test
    void testGetTimeWithIndex() {
        LocalTime time = LocalTime.NOON;
        when(mockDelegate.getTime("tag1", 0)).thenReturn(time);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(time, transformedResponse.getTime("tag1", 0));
    }

    @Test
    void testGetDateWithIndex() {
        LocalDate date = LocalDate.now();
        when(mockDelegate.getDate("tag1", 0)).thenReturn(date);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(date, transformedResponse.getDate("tag1", 0));
    }

    @Test
    void testGetDateTimeWithIndex() {
        LocalDateTime dateTime = LocalDateTime.now();
        when(mockDelegate.getDateTime("tag1", 0)).thenReturn(dateTime);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(dateTime, transformedResponse.getDateTime("tag1", 0));
    }

    // Get without transformation methods
    @Test
    void testGetShortWithoutTransformation() {
        when(mockDelegate.getShort("tag1")).thenReturn((short) 100);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals((short) 100, transformedResponse.getShort("tag1"));
    }

    @Test
    void testGetIntegerWithoutTransformation() {
        when(mockDelegate.getInteger("tag1")).thenReturn(1000);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(1000, transformedResponse.getInteger("tag1"));
    }

    @Test
    void testGetBigIntegerWithoutTransformation() {
        when(mockDelegate.getBigInteger("tag1")).thenReturn(BigInteger.TEN);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(BigInteger.TEN, transformedResponse.getBigInteger("tag1"));
    }

    @Test
    void testGetLongWithoutTransformation() {
        when(mockDelegate.getLong("tag1")).thenReturn(10000L);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(10000L, transformedResponse.getLong("tag1"));
    }

    @Test
    void testGetFloatWithoutTransformation() {
        when(mockDelegate.getFloat("tag1")).thenReturn(3.14f);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(3.14f, transformedResponse.getFloat("tag1"), 0.001);
    }

    @Test
    void testGetDoubleWithoutTransformation() {
        when(mockDelegate.getDouble("tag1")).thenReturn(3.14159);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(3.14159, transformedResponse.getDouble("tag1"), 0.00001);
    }

    @Test
    void testGetBigDecimalWithoutTransformation() {
        when(mockDelegate.getBigDecimal("tag1")).thenReturn(BigDecimal.TEN);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(BigDecimal.TEN, transformedResponse.getBigDecimal("tag1"));
    }

    @Test
    void testGetStringWithoutTransformation() {
        when(mockDelegate.getString("tag1")).thenReturn("test");
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals("test", transformedResponse.getString("tag1"));
    }

    @Test
    void testGetTimeWithoutTransformation() {
        LocalTime time = LocalTime.NOON;
        when(mockDelegate.getTime("tag1")).thenReturn(time);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(time, transformedResponse.getTime("tag1"));
    }

    @Test
    void testGetDateWithoutTransformation() {
        LocalDate date = LocalDate.now();
        when(mockDelegate.getDate("tag1")).thenReturn(date);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(date, transformedResponse.getDate("tag1"));
    }

    @Test
    void testGetDateTimeWithoutTransformation() {
        LocalDateTime dateTime = LocalDateTime.now();
        when(mockDelegate.getDateTime("tag1")).thenReturn(dateTime);
        transformedResponse = new TransformedPlcReadResponse(mockDelegate, transformedValues);
        assertEquals(dateTime, transformedResponse.getDateTime("tag1"));
    }
}
