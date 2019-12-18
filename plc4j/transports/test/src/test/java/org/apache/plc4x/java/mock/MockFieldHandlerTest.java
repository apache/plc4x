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

package org.apache.plc4x.java.mock;

import org.apache.plc4x.java.api.model.PlcField;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class MockFieldHandlerTest implements WithAssertions {

    MockFieldHandler SUT = new MockFieldHandler();

    @Mock
    PlcField plcField;

    @Test
    void createField() {
        assertThat(SUT.createField("")).isNotNull();
    }

    @Test
    void encodeBoolean() {
        assertThat(SUT.encodeBoolean(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeByte() {
        assertThat(SUT.encodeByte(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeShort() {
        assertThat(SUT.encodeShort(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeInteger() {
        assertThat(SUT.encodeInteger(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeBigInteger() {
        assertThat(SUT.encodeBigInteger(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeLong() {
        assertThat(SUT.encodeLong(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeFloat() {
        assertThat(SUT.encodeFloat(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeBigDecimal() {
        assertThat(SUT.encodeBigDecimal(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeDouble() {
        assertThat(SUT.encodeDouble(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeString() {
        assertThat(SUT.encodeString(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeTime() {
        assertThat(SUT.encodeTime(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeDate() {
        assertThat(SUT.encodeDate(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeDateTime() {
        assertThat(SUT.encodeDateTime(plcField, new Object[0])).isNotNull();
    }

    @Test
    void encodeByteArray() {
        assertThat(SUT.encodeByteArray(plcField, new Object[0])).isNotNull();
    }
}