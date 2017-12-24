package org.apache.plc4x.java.isotp.netty.model.types;

import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TypeTests {

    @Test
    @Tag("fast")
    void deviceGroupTest() {
        DeviceGroup deviceGroup;

        deviceGroup = DeviceGroup.valueOf((byte)1);
        assertTrue(deviceGroup == DeviceGroup.PG_OR_PC, "1 incorrectly mapped");
        assertTrue(deviceGroup.getCode() == (byte)1, "code is not 1");

        deviceGroup = DeviceGroup.valueOf((byte)2);
        assertTrue(deviceGroup == DeviceGroup.OS, "2 incorrectly mapped");
        assertTrue(deviceGroup.getCode() == (byte)2, "code is not 2");

        deviceGroup = DeviceGroup.valueOf((byte)3);
        assertTrue(deviceGroup == DeviceGroup.OTHERS, "3 incorrectly mapped");
        assertTrue(deviceGroup.getCode() == (byte)3, "code is not 3");
    }

    @Test
    @Tag("fast")
    void deviceGroupUnknownTest() {
        DeviceGroup deviceGroup = DeviceGroup.valueOf((byte)0x40);

        assertNull(deviceGroup, "expected device group to be null");
    }


    @Test
    @Tag("fast")
    void disconnectReasonTest() {
        DisconnectReason disconnectReason = DisconnectReason.valueOf((byte)3);

        assertTrue(DisconnectReason.valueOf((byte)3) == DisconnectReason.ADDRESS_UNKNOWN, "3 incorrectly mapped");
        assertTrue(disconnectReason.getCode() == (byte)3, "code is not 3");
    }

    @Test
    @Tag("fast")
    void diosconectReasonUnknownTest() {
        DisconnectReason disconnectReason = DisconnectReason.valueOf((byte)4);

        assertNull(disconnectReason, "expected disconnect reason to be null");
    }

    @Test
    @Tag("fast")
    void parameterCodeTest() {
        ParameterCode parameterCode = ParameterCode.valueOf((byte)0xC1);

        assertTrue(ParameterCode.valueOf((byte)0xC1) == ParameterCode.CALLING_TSAP, "0xC1 incorrectly mapped");
        assertTrue(parameterCode.getCode() == (byte)0xC1, "code is not 0xC1");
    }

    @Test
    @Tag("fast")
    void parameterCodeUnknownTest() {
        ParameterCode parameterCode = ParameterCode.valueOf((byte)0x90);

        assertNull(parameterCode, "expected parameter code to be null");
    }

    @Test
    @Tag("fast")
    void protocolClassTest() {
        ProtocolClass protocolClass;

        protocolClass = ProtocolClass.valueOf((byte)0x10);
        assertTrue(protocolClass == ProtocolClass.CLASS_1, "0x10 incorrectly mapped");
        assertTrue(protocolClass.getCode() == (byte)0x10, "code is not 0x10");

        protocolClass = ProtocolClass.valueOf((byte)0x20);
        assertTrue(protocolClass == ProtocolClass.CLASS_2, "0x20 incorrectly mapped");
        assertTrue(protocolClass.getCode() == (byte)0x20, "code is not 0x20");

        protocolClass = ProtocolClass.valueOf((byte)0x30);
        assertTrue(protocolClass == ProtocolClass.CLASS_3, "0x30 incorrectly mapped");
        assertTrue(protocolClass.getCode() == (byte)0x30, "code is not 0x30");

        protocolClass = ProtocolClass.valueOf((byte)0x40);
        assertTrue(protocolClass == ProtocolClass.CLASS_4, "0x40 incorrectly mapped");
        assertTrue(protocolClass.getCode() == (byte)0x40, "code is not 0x40");
    }

    @Test
    @Tag("fast")
    void protocolClassUnknownTest() {
        ProtocolClass protocolClass = ProtocolClass.valueOf((byte)0x50);

        assertNull(protocolClass, "expected protocol class to be null");
    }

    @Test
    @Tag("fast")
    void rejectCauseTest() {
        RejectCause rejectCause = RejectCause.valueOf((byte)0x03);

        assertTrue(RejectCause.valueOf((byte)0x03) == RejectCause.INVALID_PARAMETER_TYPE, "0x03 incorrectly mapped");
        assertTrue(rejectCause.getCode() == (byte)0x03, "code is not 0x03");
    }

    @Test
    @Tag("fast")
    void rejectClauseUnknownTest() {
        RejectCause rejectCause = RejectCause.valueOf((byte)0x90);

        assertNull(rejectCause, "expected reject cause to be null");
    }

    @Test
    @Tag("fast")
    void tpduCodeTest() {
        TpduCode tpduCode = TpduCode.valueOf((byte)0xF0);

        assertTrue(TpduCode.valueOf((byte)0xF0) == TpduCode.DATA, "0xF0 incorrectly mapped");
        assertTrue(tpduCode.getCode() == (byte)0xF0, "code is not 0xF0");
    }

    @Test
    @Tag("fast")
    void tpduCodeUnknownTest() {
        TpduCode tpduCode = TpduCode.valueOf((byte)0x01);

        assertNull(tpduCode, "expected tpdu code to be null");
    }

    @Test
    @Tag("fast")
    void typduSizeTest() {
        TpduSize tpduSize = TpduSize.valueOf((byte)0x07);

        assertTrue(TpduSize.valueOf((byte)0x07) == TpduSize.SIZE_128, "0x07 incorrectly mapped");
        assertTrue(tpduSize.getCode() == (byte)0x07, "code is not 0x07");
    }

    @Test
    @Tag("fast")
    void tpduSizeUnknownTest() {
        TpduSize tpduSize = TpduSize.valueOf((byte)0x06);

        assertNull(tpduSize, "expected tpdu size to be null");
    }


}