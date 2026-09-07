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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class S7BrowseServiceTest {

    @Test
    void staticAreas_includesMarkersInputsAndOutputs() {
        List<PlcBrowseItem> items = S7BrowseService.staticAreas(false);
        assertEquals(3, items.size());
        assertEquals("%M", items.get(0).getName());
        assertEquals("%I", items.get(1).getName());
        assertEquals("%Q", items.get(2).getName());
        for (PlcBrowseItem item : items) {
            assertTrue(item.isReadable());
            assertTrue(item.isWritable());
        }
    }

    @Test
    void staticAreas_propagatesSubscribableFlag() {
        for (PlcBrowseItem item : S7BrowseService.staticAreas(true)) {
            assertTrue(item.isSubscribable());
        }
        for (PlcBrowseItem item : S7BrowseService.staticAreas(false)) {
            assertFalse(item.isSubscribable());
        }
    }

    @Test
    void buildListBlocksOfTypeRequest_setsCpuFunctionDiscriminators() {
        S7Message msg = S7BrowseService.buildListBlocksOfTypeRequest(7, S7BrowseService.BlockType.DB);
        assertInstanceOf(S7MessageUserData.class, msg);
        S7ParameterUserDataItemCPUFunctions cpu = (S7ParameterUserDataItemCPUFunctions)
            ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
        assertEquals(0x03, cpu.getCpuFunctionGroup());
        assertEquals(0x04, cpu.getCpuFunctionType());
        assertEquals(0x02, cpu.getCpuSubfunction());
        assertEquals(7, msg.getTpduReference());
        S7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest req =
            (S7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        assertEquals(0x3041, req.getBlockType());
    }

    @Test
    void parseListBlocksOfTypeResponse_decodesBlockNumbers() {
        // 4 bytes per entry: blockNumber:uint16 BE + flags:uint8 + lang:uint8.
        byte[] payload = new byte[] {
            0x00, 0x01, 0x0A, 0x07,    // DB1
            0x00, 0x0A, 0x0A, 0x07,    // DB10
            0x00, (byte) 0xFF, 0x0A, 0x07  // DB255
        };
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, payload.length, payload),
            /*errorCode=*/null);
        List<Integer> numbers = S7BrowseService.parseListBlocksOfTypeResponse(response);
        assertEquals(List.of(1, 10, 255), numbers);
    }

    @Test
    void parseListBlocksOfTypeResponse_emptyPayloadReturnsEmptyList() {
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0, new byte[0]),
            /*errorCode=*/null);
        assertEquals(Collections.emptyList(), S7BrowseService.parseListBlocksOfTypeResponse(response));
    }

    @Test
    void parseListBlocksOfTypeResponse_errorCodeReturnsEmptyList() {
        byte[] payload = new byte[] { 0x00, 0x01, 0x0A, 0x07 };
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, payload.length, payload),
            /*errorCode=*/0xD401);
        assertEquals(Collections.emptyList(), S7BrowseService.parseListBlocksOfTypeResponse(response));
    }

    @Test
    void dataBlocksFromNumbers_yieldsBrowseItemsWithDbAddresses() {
        List<PlcBrowseItem> items = S7BrowseService.dataBlocksFromNumbers(List.of(1, 10), true);
        assertEquals(2, items.size());
        assertEquals("DB1", items.get(0).getName());
        assertEquals("DB10", items.get(1).getName());
        assertTrue(items.get(0).isSubscribable());
        S7Tag tag = (S7Tag) items.get(0).getTag();
        assertEquals(MemoryArea.DATA_BLOCKS, tag.getMemoryArea());
        assertEquals(1, tag.getBlockNumber());
    }

    private static S7Message wrapInUserData(S7PayloadUserDataItem payloadItem, Integer errorCode) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x03, (short) 0x02,
            (short) 0x00, null, null, errorCode);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.singletonList(payloadItem)));
    }
}
