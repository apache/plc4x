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

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageRequest;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlResponse;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.*;

class S7SzlServiceTest {

    @Test
    void buildRequest_producesUserDataMessage() {
        S7Message msg = S7SzlService.buildRequest(42, S7SzlService.MODULE_IDENTIFICATION);
        assertInstanceOf(S7MessageUserData.class, msg);
        assertEquals(42, msg.getTpduReference());
        S7ParameterUserData param = (S7ParameterUserData) msg.getParameter();
        S7ParameterUserDataItemCPUFunctions item =
            (S7ParameterUserDataItemCPUFunctions) param.getItems().get(0);
        // CPU function group=0x04 (CPU services), type=0x04 (request), subfunction=0x01 (ReadSZL).
        assertEquals(0x04, item.getCpuFunctionGroup());
        assertEquals(0x04, item.getCpuFunctionType());
        assertEquals(0x01, item.getCpuSubfunction());
    }

    @Test
    void parseProbeResponse_decodesS71200Article() {
        S7Message response = mockSzlResponse("6ES7 212-1AE40-0XB0      ");
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(response);
        assertEquals(ControllerType.S7_1200, result.controllerType());
        assertTrue(result.articleNumber().startsWith("6ES7 2"), result.articleNumber());
    }

    @Test
    void parseProbeResponse_decodesS71500Article() {
        assertEquals(ControllerType.S7_1500,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 511-1AK00-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_decodesS7300Article() {
        assertEquals(ControllerType.S7_300,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 314-6CH04-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_decodesS7400Article() {
        assertEquals(ControllerType.S7_400,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 416-2XK02-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_rejectsResponseWithoutSiemensArticleNumber() {
        // Response payload decodes as ASCII but doesn't contain a "6ES7"/"6GK" order number
        // or a "CPU NNN" model name. The response is still structurally valid — that's enough
        // to mark UserData services as supported. ControllerType falls back to ANY.
        S7Message resp = mockSzlResponse("VENDORX 999            ");
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(resp);
        assertEquals(ControllerType.ANY, result.controllerType());
        assertEquals("", result.articleNumber());
    }

    @Test
    void parseProbeResponse_decodesCpuModelNameForS7300() {
        S7Message resp = mockSzlResponse("CPU 315-2 PN/DP                                  ");
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(resp);
        assertEquals(ControllerType.S7_300, result.controllerType());
        assertTrue(result.articleNumber().startsWith("CPU 315"), result.articleNumber());
    }

    @Test
    void parseProbeResponse_decodesCpuModelNameForS71200() {
        S7Message resp = mockSzlResponse("CPU 1212C                                        ");
        assertEquals(ControllerType.S7_1200,
            S7SzlService.parseProbeResponse(resp).controllerType());
    }

    @Test
    void parseProbeResponse_decodesCpuModelNameForS71500() {
        S7Message resp = mockSzlResponse("CPU 1511-1 PN                                    ");
        assertEquals(ControllerType.S7_1500,
            S7SzlService.parseProbeResponse(resp).controllerType());
    }

    @Test
    void parseProbeResponse_unknownModelDigitIsAny() {
        assertEquals(ControllerType.ANY,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 999-0AA00-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_acceptsZeroFilledItemAsAny() {
        // Empty/zeroed item is structurally valid — the probe succeeds with an empty
        // article and ANY controllerType. The connection treats that as "UserData works".
        S7Message resp = wrapInUserData(new S7PayloadUserDataItemCpuFunctionReadSzlResponse(
            DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0x10, new byte[16]));
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(resp);
        assertEquals(ControllerType.ANY, result.controllerType());
        assertEquals("", result.articleNumber());
    }

    @Test
    void parseProbeResponse_rejectsMissingSzlItem() {
        S7Message resp = new S7MessageUserData(7,
            new S7ParameterUserData(Collections.emptyList()),
            new S7PayloadUserData(Collections.emptyList()));
        assertThrows(IllegalStateException.class, () -> S7SzlService.parseProbeResponse(resp));
    }

    @Test
    void parseProbeResponse_rejectsNonUserDataMessage() {
        S7Message reqMessage = new S7MessageRequest(1, null, null);
        assertThrows(IllegalArgumentException.class,
            () -> S7SzlService.parseProbeResponse(reqMessage));
    }

    /** Build a 22-byte SZL item: 2 bytes header + 20-byte article-number window. */
    private static S7Message mockSzlResponse(String article) {
        byte[] articleBytes = article.getBytes(StandardCharsets.US_ASCII);
        byte[] szlItem = new byte[22];
        // Bytes 0..1 are the 2-byte index header — opaque to us. Bytes 2..21 hold the MLFB.
        int copyLen = Math.min(articleBytes.length, 20);
        System.arraycopy(articleBytes, 0, szlItem, 2, copyLen);
        // Pad with spaces so the trim() hits real content.
        for (int i = 2 + copyLen; i < 22; i++) {
            szlItem[i] = ' ';
        }
        return wrapInUserData(new S7PayloadUserDataItemCpuFunctionReadSzlResponse(
            DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, szlItem.length, szlItem));
    }

    private static S7Message wrapInUserData(S7PayloadUserDataItem item) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x01,
            (short) 0x00, null, null, null);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Arrays.asList(item)));
    }
}
