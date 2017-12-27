package org.apache.plc4x.java.s7.netty.model.types;

import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class S7TypeTests {

    @Test
    @Tag("fast")
    void dataTransportErrorCode() {
        DataTransportErrorCode dataTransportErrorCode = DataTransportErrorCode.INVALID_ADDRESS;

        assertTrue(DataTransportErrorCode.valueOf((byte)5) == DataTransportErrorCode.INVALID_ADDRESS, "5 incorrectly mapped");
        assertTrue(dataTransportErrorCode.getCode() == (byte)5, "code is not 5");
    }

    @Test
    @Tag("fast")
    void dataTransportSize() {
        DataTransportSize dataTransportSize = DataTransportSize.DINTEGER;

        assertTrue(DataTransportSize.valueOf((byte)6) == DataTransportSize.DINTEGER, "6 incorrectly mapped");
        assertTrue(dataTransportSize.getCode() == (byte)6, "code is not 6");
    }

    @Test
    @Tag("fast")
    void headerErrorClass() {
        //TODO missing valueOf
    }

    @Test
    @Tag("fast")
    void messageType() {
        MessageType messageType = MessageType.ACK;

        assertTrue(MessageType.valueOf((byte)2) == MessageType.ACK, "2 incorrectly mapped");
        assertTrue(messageType.getCode() == (byte)2, "code is not 2");
    }

    @Test
    @Tag("fast")
    void parameterError() {
        //TODO missing value of
    }

    @Test
    @Tag("fast")
    void parameterType() {
        ParameterType parameterType = ParameterType.UPLOAD;

        assertTrue(ParameterType.valueOf((byte)0x1E) == ParameterType.UPLOAD, "0x1E incorrectly mapped");
        assertTrue(parameterType.getCode() == (byte)0x1E, "code is not 0x1E");
    }

    @Test
    @Tag("fast")
    void specificationType() {
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;

        assertTrue(SpecificationType.valueOf((byte)0x12) == SpecificationType.VARIABLE_SPECIFICATION, "0x12 incorrectly mapped");
        assertTrue(specificationType.getCode() == (byte)0x12, "code is not 0x12");
    }

    @Test
    @Tag("fast")
    void transportSize() {
        TransportSize transportSize = TransportSize.TIME;

        assertTrue(TransportSize.valueOf((byte)0x0B) == TransportSize.TIME, "0x0B incorrectly mapped");
        assertTrue(transportSize.getCode() == (byte)0x0B, "code is not 0x0B");
    }

    @Test
    @Tag("fast")
    void variableAddressingMode() {
        VariableAddressingMode variableAddressingMode = VariableAddressingMode.ALARM_ACK;

        assertTrue(VariableAddressingMode.valueOf((byte)0x19) == VariableAddressingMode.ALARM_ACK, "0x19 incorrectly mapped");
        assertTrue(variableAddressingMode.getCode() == (byte)0x19, "code is not 0x19");
    }
}