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

package org.apache.plc4x.java.simulated.connection;

import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;

import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
@SuppressWarnings("unchecked")
class TestFieldHandlerTest implements WithAssertions {

    TestFieldHandler SUT = new TestFieldHandler();

    @Mock
    TestField plcField;

    @Test
    void createField() {
        assertThat(SUT.createField("STATE/bar:INTEGER")).isNotNull();
    }

    @Test
    void encodeBoolean() {
        when(plcField.getDataType()).thenReturn((Class) Boolean.class);
        assertThat(SUT.encodeBoolean(plcField, new Boolean[0])).isNotNull();
    }

    @Test
    void encodeByte() {
        when(plcField.getDataType()).thenReturn((Class) Byte.class);
        assertThat(SUT.encodeByte(plcField, new Byte[0])).isNotNull();
    }

    @Test
    void encodeShort() {
        when(plcField.getDataType()).thenReturn((Class) Short.class);
        assertThat(SUT.encodeShort(plcField, new Short[0])).isNotNull();
    }

    @Test
    void encodeInteger() {
        when(plcField.getDataType()).thenReturn((Class) Integer.class);
        assertThat(SUT.encodeInteger(plcField, new Integer[0])).isNotNull();
    }

    @Test
    void encodeBigInteger() {
        when(plcField.getDataType()).thenReturn((Class) BigInteger.class);
        assertThat(SUT.encodeBigInteger(plcField, new BigInteger[0])).isNotNull();
    }

    @Test
    void encodeLong() {
        when(plcField.getDataType()).thenReturn((Class) Long.class);
        assertThat(SUT.encodeLong(plcField, new Long[0])).isNotNull();
    }

    @Test
    void encodeFloat() {
        when(plcField.getDataType()).thenReturn((Class) Float.class);
        assertThat(SUT.encodeFloat(plcField, new Float[0])).isNotNull();
    }

    @Test
    void encodeBigDecimal() {
        when(plcField.getDataType()).thenReturn((Class) BigDecimal.class);
        assertThat(SUT.encodeBigDecimal(plcField, new BigDecimal[0])).isNotNull();
    }

    @Test
    void encodeDouble() {
        when(plcField.getDataType()).thenReturn((Class) Double.class);
        assertThat(SUT.encodeDouble(plcField, new Double[0])).isNotNull();
    }

    @Test
    void encodeString() {
        when(plcField.getDataType()).thenReturn((Class) String.class);
        assertThat(SUT.encodeString(plcField, new String[0])).isNotNull();
    }

    @Test
    void encodeTime() {
        when(plcField.getDataType()).thenReturn((Class) LocalTime.class);
        assertThat(SUT.encodeTime(plcField, new LocalTime[0])).isNotNull();
    }

    @Test
    void encodeDate() {
        when(plcField.getDataType()).thenReturn((Class) LocalDate.class);
        assertThat(SUT.encodeDate(plcField, new LocalDate[0])).isNotNull();
    }

    @Test
    void encodeDateTime() {
        when(plcField.getDataType()).thenReturn((Class) LocalDateTime.class);
        assertThat(SUT.encodeDateTime(plcField, new LocalDateTime[0])).isNotNull();
    }

    @Test
    void encodeByteArray() {
        when(plcField.getDataType()).thenReturn((Class) Byte[].class);
        assertThat(SUT.encodeByteArray(plcField, new Byte[0])).isNotNull();
    }
}