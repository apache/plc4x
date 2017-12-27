package org.apache.plc4x.java.s7.netty.model.payloads;

import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportErrorCode;
import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class S7PayloadTests {

    @Test
    @Tag("fast")
    void varPayloadItem() {
        DataTransportErrorCode returnCode = DataTransportErrorCode.NOT_FOUND;
        DataTransportSize dataTransportSize = DataTransportSize.INTEGER;
        byte[] data = {(byte)0xFF};

        VarPayloadItem  varPayloadItem = new VarPayloadItem(returnCode, dataTransportSize, data);
        assertTrue(varPayloadItem.getReturnCode() == DataTransportErrorCode.NOT_FOUND, "Unexpected data transport error code");
        assertTrue(varPayloadItem.getDataTransportSize() == DataTransportSize.INTEGER, "Unexpected data transport size");
        assertTrue(varPayloadItem.getData()[0] == (byte) 0xFF, "Unexpected user data");
    }

    @Test
    @Tag("fast")
    void varPayload() {
        ParameterType parameterType = ParameterType.DOWNLOAD_ENDED;
        ArrayList<VarPayloadItem> payloadItems = new ArrayList<>();
        byte[] data = {(byte)0xFF};
        
        payloadItems.add(new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, data));

        VarPayload  varPayload = new VarPayload(parameterType, payloadItems);
        assertTrue(varPayload.getType() == ParameterType.DOWNLOAD_ENDED, "Unexpected parameter type");
        assertTrue(varPayload.getPayloadItems().containsAll(payloadItems) , "Unexpected pay load items");
    }

}