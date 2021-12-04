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
package org.apache.plc4x.java.simulated.field;

import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class SimularedFieldHandlerTest implements WithAssertions {

    SimulatedFieldHandler SUT = new SimulatedFieldHandler();

    @Mock
    SimulatedField plcField;

    @Test
    void createField() {
        assertThat(SUT.createField("STATE/bar:Integer")).isNotNull();
    }

    /*@Test
    void encodeBoolean() {
        assertThat(SUT.encodeBoolean(plcField, new Boolean[0])).isNotNull();
    }

    @Test
    void encodeByte() {
        when(plcField.getPlcDataType()).thenReturn("Byte");
        assertThat(SUT.encodeByte(plcField, new Byte[0])).isNotNull();
    }

    @Test
    void encodeShort() {
        when(plcField.getPlcDataType()).thenReturn("Short");
        assertThat(SUT.encodeShort(plcField, new Short[0])).isNotNull();
    }

    @Test
    void encodeInteger() {
        when(plcField.getPlcDataType()).thenReturn("Integer");
        assertThat(SUT.encodeInteger(plcField, new Integer[0])).isNotNull();
    }

    @Test
    void encodeBigInteger() {
        when(plcField.getPlcDataType()).thenReturn("BigInteger");
        assertThat(SUT.encodeBigInteger(plcField, new BigInteger[0])).isNotNull();
    }

    @Test
    void encodeLong() {
        when(plcField.getPlcDataType()).thenReturn("Long");
        assertThat(SUT.encodeLong(plcField, new Long[0])).isNotNull();
    }

    @Test
    void encodeFloat() {
        when(plcField.getPlcDataType()).thenReturn("Float");
        assertThat(SUT.encodeFloat(plcField, new Float[0])).isNotNull();
    }

    @Test
    void encodeBigDecimal() {
        when(plcField.getPlcDataType()).thenReturn("BigDecimal");
        assertThat(SUT.encodeBigDecimal(plcField, new BigDecimal[0])).isNotNull();
    }

    @Test
    void encodeDouble() {
        when(plcField.getPlcDataType()).thenReturn("Double");
        assertThat(SUT.encodeDouble(plcField, new Double[0])).isNotNull();
    }

    @Test
    void encodeString() {
        when(plcField.getPlcDataType()).thenReturn("String");
        assertThat(SUT.encodeString(plcField, new String[0])).isNotNull();
    }

    @Test
    void encodeTime() {
        when(plcField.getPlcDataType()).thenReturn("LocalTime");
        assertThat(SUT.encodeTime(plcField, new LocalTime[0])).isNotNull();
    }

    @Test
    void encodeDate() {
        when(plcField.getPlcDataType()).thenReturn("LocalDate");
        assertThat(SUT.encodeDate(plcField, new LocalDate[0])).isNotNull();
    }

    @Test
    void encodeDateTime() {
        when(plcField.getPlcDataType()).thenReturn("LocalDateTime");
        assertThat(SUT.encodeDateTime(plcField, new LocalDateTime[0])).isNotNull();
    }*/

}
