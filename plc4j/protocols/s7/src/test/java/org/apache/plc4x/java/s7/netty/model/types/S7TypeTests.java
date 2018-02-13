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

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

public class S7TypeTests {

    @Test
    @Category(FastTests.class)
    public void dataTransportErrorCode() {
        DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.INVALID_ADDRESS;

        assertThat(DataTransportErrorCode.valueOf((byte)5)).isEqualTo(DataTransportErrorCode.INVALID_ADDRESS).withFailMessage("5 incorrectly mapped");
        assertThat(dataTransportErrorCode.getCode()).isEqualTo((byte)5).withFailMessage("code is not 5");
    }

    @Test
    @Category(FastTests.class)
    public void unknownDataTransportErrorCode() {
        assertThat(DataTransportErrorCode.valueOf((byte)0xFE)).isNull();

        DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.INVALID_ADDRESS;
        assertThat(dataTransportErrorCode.getCode()).isEqualTo((byte) 0x05).withFailMessage("code is not 0x05");
    }

    @Test
    @Category(FastTests.class)
    public void dataTransportSize() {
        DataTransportSize dataTransportSize = DataTransportSize.DINTEGER;

        assertThat(DataTransportSize.valueOf((byte)6)).isEqualTo(DataTransportSize.DINTEGER).withFailMessage("6 incorrectly mapped");
        assertThat(dataTransportSize.getCode()).isEqualTo((byte)6).withFailMessage("code is not 6");
        assertThat(dataTransportSize.isSizeInBits()).isFalse().withFailMessage("Unexpected return from bit size");
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

        assertThat(ParameterType.valueOf((byte)0x1E)).isEqualTo(ParameterType.UPLOAD).withFailMessage("0x1E incorrectly mapped");
        assertThat(parameterType.getCode()).isEqualTo((byte)0x1E).withFailMessage("code is not 0x1E");
    }

    @Test
    @Category(FastTests.class)
    public void unknownParameterType() {
        ParameterType parameterType = ParameterType.UPLOAD;

        assertThat(ParameterType.valueOf((byte)0xFF)).isNull();;
    }

    @Test
    @Category(FastTests.class)
    public void specificationType() {
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;

        assertThat(SpecificationType.valueOf((byte)0x12)).isEqualTo(SpecificationType.VARIABLE_SPECIFICATION).withFailMessage("0x12 incorrectly mapped");
        assertThat(specificationType.getCode()).isEqualTo((byte)0x12).withFailMessage("code is not 0x12");
    }

    @Test
    @Category(FastTests.class)
    public void transportSize() {
        TransportSize transportSize = TransportSize.TIME;

        assertThat(TransportSize.valueOf((byte)0x0B)).isEqualTo(TransportSize.TIME).withFailMessage("0x0B incorrectly mapped");
        assertThat(transportSize.getCode()).isEqualTo((byte)0x0B).withFailMessage("code is not 0x0B");
    }

    @Test
    @Category(FastTests.class)
    public void variableAddressingMode() {
        VariableAddressingMode variableAddressingMode = VariableAddressingMode.ALARM_ACK;

        assertThat(VariableAddressingMode.valueOf((byte)0x19)).isEqualTo(VariableAddressingMode.ALARM_ACK).withFailMessage("0x19 incorrectly mapped");
        assertThat(variableAddressingMode.getCode()).isEqualTo((byte)0x19).withFailMessage("code is not 0x19");
    }

    @Test
    @Category(FastTests.class)
    public void memoryAccess() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;

        assertThat(MemoryArea.valueOf((byte)0x84)).isEqualTo(MemoryArea.DATA_BLOCKS).withFailMessage("0x84 incorrectly mapped");
        assertThat(memoryArea.getCode()).isEqualTo((byte)0x84).withFailMessage("code is not 0x84");
    }

    @Test
    @Category(FastTests.class)
    public void unknownMemoryAccess() {
        assertThat(MemoryArea.valueOf((byte)0xFF)).isNull();
     }

}