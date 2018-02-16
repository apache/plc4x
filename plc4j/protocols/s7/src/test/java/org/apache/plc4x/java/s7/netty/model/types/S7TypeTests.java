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

package org.apache.plc4x.java.s7.netty.model.types;

import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;

public class S7TypeTests {

    @Test
    @Category(FastTests.class)
    public void dataTransportErrorCode() {
        DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.INVALID_ADDRESS;

        assertThat("5 incorrectly mapped", DataTransportErrorCode.valueOf((byte) 5), equalTo(DataTransportErrorCode.INVALID_ADDRESS));
        assertThat("code is not 5", dataTransportErrorCode.getCode(), equalTo((byte) 5));
    }

    @Test
    @Category(FastTests.class)
    public void unknownDataTransportErrorCode() {
        DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.INVALID_ADDRESS;

        assertThat(DataTransportErrorCode.valueOf((byte) 0xFE), nullValue());
        assertThat(DataTransportErrorCode.valueOf((byte)0xFE)).isNull();

        DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.INVALID_ADDRESS;
        assertThat(dataTransportErrorCode.getCode()).isEqualTo((byte) 0x05).withFailMessage("code is not 0x05");
    }

    @Test
    @Category(FastTests.class)
    public void dataTransportSize() {
        DataTransportSize dataTransportSize = DataTransportSize.DINTEGER;

        assertThat("6 incorrectly mapped", DataTransportSize.valueOf((byte) 6), equalTo(DataTransportSize.DINTEGER));
        assertThat(dataTransportSize.getCode(), equalTo((byte) 6));
        assertThat("Unexpected return from bit size", dataTransportSize.isSizeInBits(), is(false));
    }

    @Test
    @Category(FastTests.class)
    public void headerErrorClass() {
        assertThat(HeaderErrorClass.valueOf((byte) 0x87)).isEqualTo(HeaderErrorClass.ACCESS_ERROR).withFailMessage("0x87 incorrectly mapped");

        HeaderErrorClass headerErrorClass = HeaderErrorClass.ACCESS_ERROR;
        assertThat(headerErrorClass.getCode()).isEqualTo((byte) 0x87).withFailMessage("code is not 0x87");
    }

    @Test
    @Category(FastTests.class)
    public void messageType() {
        MessageType messageType = MessageType.ACK;

        assertThat("2 incorrectly mapped", MessageType.valueOf((byte) 2), equalTo(MessageType.ACK));
        assertThat("code is not 2", messageType.getCode(), equalTo((byte) 2));
        assertThat(MessageType.valueOf((byte)2)).isEqualTo(MessageType.ACK).withFailMessage("2 incorrectly mapped");

        MessageType messageType = MessageType.ACK;
        assertThat(messageType.getCode()).isEqualTo((byte)2).withFailMessage("code is not 2");
    }

    @Test
    @Category(FastTests.class)
    public void parameterError() {
        assertThat(ParameterError.valueOf((short) 0x011C)).isEqualTo(ParameterError.PROTOCOL_ERROR).withFailMessage("0x011C incorrectly mapped");;

        ParameterError parameterError = ParameterError.PROTOCOL_ERROR;
        assertThat(parameterError.getCode()).isEqualTo((short) 0x011C).withFailMessage("code is not 0x011C");
    }

    @Test
    @Category(FastTests.class)
    public void parameterType() {
        ParameterType parameterType = ParameterType.UPLOAD;

        assertThat("0x1E incorrectly mapped", ParameterType.valueOf((byte) 0x1E), equalTo(ParameterType.UPLOAD));
        assertThat("code is not 0x1E", parameterType.getCode(), equalTo((byte) 0x1E));
    }

    @Test
    @Category(FastTests.class)
    public void unknownParameterType() {
        ParameterType parameterType = ParameterType.UPLOAD;

        assertThat(ParameterType.valueOf((byte) 0xFF), nullValue());
        ;
    }

    @Test
    @Category(FastTests.class)
    public void specificationType() {
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;

        assertThat("0x12 incorrectly mapped", SpecificationType.valueOf((byte) 0x12), equalTo(SpecificationType.VARIABLE_SPECIFICATION));
        assertThat("code is not 0x12", specificationType.getCode(), equalTo((byte) 0x12));
    }

    @Test
    @Category(FastTests.class)
    public void transportSize() {
        TransportSize transportSize = TransportSize.TIME;

        assertThat("0x0B incorrectly mapped", TransportSize.valueOf((byte) 0x0B), equalTo(TransportSize.TIME));
        assertThat("code is not 0x0B", transportSize.getCode(), equalTo((byte) 0x0B));
    }

    @Test
    @Category(FastTests.class)
    public void variableAddressingMode() {
        VariableAddressingMode variableAddressingMode = VariableAddressingMode.ALARM_ACK;

        assertThat("0x19 incorrectly mapped", VariableAddressingMode.valueOf((byte) 0x19), equalTo(VariableAddressingMode.ALARM_ACK));
        assertThat("code is not 0x19", variableAddressingMode.getCode(), equalTo((byte) 0x19));
    }

    @Test
    @Category(FastTests.class)
    public void memoryAccess() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;

        assertThat("0x84 incorrectly mapped", MemoryArea.valueOf((byte) 0x84), equalTo(MemoryArea.DATA_BLOCKS));
        assertThat("code is not 0x84", memoryArea.getCode(), equalTo((byte) 0x84));
    }

    @Test
    @Category(FastTests.class)
    public void unknownMemoryAccess() {
        assertThat(MemoryArea.valueOf((byte) 0xFF), nullValue());
    }

}