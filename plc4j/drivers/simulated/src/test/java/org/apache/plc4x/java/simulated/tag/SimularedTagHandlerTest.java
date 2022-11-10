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
package org.apache.plc4x.java.simulated.tag;

import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.apache.plc4x.java.simulated.tag.SimulatedTagHandler;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class SimularedTagHandlerTest implements WithAssertions {

    SimulatedTagHandler SUT = new SimulatedTagHandler();

    @Mock
    SimulatedTag plcTag;

    @Test
    void createTag() {
        assertThat(SUT.parseTag("STATE/bar:DINT")).isNotNull();
    }

    /*@Test
    void encodeBoolean() {
        assertThat(SUT.encodeBoolean(plcTag, new Boolean[0])).isNotNull();
    }

    @Test
    void encodeByte() {
        when(plcTag.getPlcDataType()).thenReturn("Byte");
        assertThat(SUT.encodeByte(plcTag, new Byte[0])).isNotNull();
    }

    @Test
    void encodeShort() {
        when(plcTag.getPlcDataType()).thenReturn("Short");
        assertThat(SUT.encodeShort(plcTag, new Short[0])).isNotNull();
    }

    @Test
    void encodeInteger() {
        when(plcTag.getPlcDataType()).thenReturn("Integer");
        assertThat(SUT.encodeInteger(plcTag, new Integer[0])).isNotNull();
    }

    @Test
    void encodeBigInteger() {
        when(plcTag.getPlcDataType()).thenReturn("BigInteger");
        assertThat(SUT.encodeBigInteger(plcTag, new BigInteger[0])).isNotNull();
    }

    @Test
    void encodeLong() {
        when(plcTag.getPlcDataType()).thenReturn("Long");
        assertThat(SUT.encodeLong(plcTag, new Long[0])).isNotNull();
    }

    @Test
    void encodeFloat() {
        when(plcTag.getPlcDataType()).thenReturn("Float");
        assertThat(SUT.encodeFloat(plcTag, new Float[0])).isNotNull();
    }

    @Test
    void encodeBigDecimal() {
        when(plcTag.getPlcDataType()).thenReturn("BigDecimal");
        assertThat(SUT.encodeBigDecimal(plcTag, new BigDecimal[0])).isNotNull();
    }

    @Test
    void encodeDouble() {
        when(plcTag.getPlcDataType()).thenReturn("Double");
        assertThat(SUT.encodeDouble(plcTag, new Double[0])).isNotNull();
    }

    @Test
    void encodeString() {
        when(plcTag.getPlcDataType()).thenReturn("String");
        assertThat(SUT.encodeString(plcTag, new String[0])).isNotNull();
    }

    @Test
    void encodeTime() {
        when(plcTag.getPlcDataType()).thenReturn("LocalTime");
        assertThat(SUT.encodeTime(plcTag, new LocalTime[0])).isNotNull();
    }

    @Test
    void encodeDate() {
        when(plcTag.getPlcDataType()).thenReturn("LocalDate");
        assertThat(SUT.encodeDate(plcTag, new LocalDate[0])).isNotNull();
    }

    @Test
    void encodeDateTime() {
        when(plcTag.getPlcDataType()).thenReturn("LocalDateTime");
        assertThat(SUT.encodeDateTime(plcTag, new LocalDateTime[0])).isNotNull();
    }*/

}
