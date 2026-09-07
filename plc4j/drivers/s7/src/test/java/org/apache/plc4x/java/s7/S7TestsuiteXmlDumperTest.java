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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.userdata.S7SzlService;
import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Collections;
import java.util.HexFormat;
import java.util.List;

/**
 * Temporary helper dumping the canonical XML serialization of the messages used by the
 * S7 driver testsuite into target/testsuite-xml/.
 */
class S7TestsuiteXmlDumperTest {

    @Test
    void dumpTestsuiteMessages() throws Exception {
        // Setup communication request/response
        dump("01-setup-request", new S7MessageRequest(10,
            new S7ParameterSetupCommunication(8, 8, 1024), null));
        dump("02-setup-response", new S7MessageResponseData(10,
            new S7ParameterSetupCommunication(3, 3, 240), null, (short) 0, (short) 0));

        // SZL identification request/response
        dump("03-szl-request", S7SzlService.buildRequest(11, S7SzlService.COMPONENT_IDENTIFICATION, 0x0001));
        byte[] szlItems = HexFormat.of().parseHex(
            "001c0001001c0001000136455337203231322d31424433302d3058423020202000012020");
        S7ParameterUserDataItemCPUFunctions szlResponseParameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x01, (short) 0x01,
            (short) 0, (short) 0, 0);
        dump("04-szl-response", new S7MessageUserData(11,
            new S7ParameterUserData(Collections.singletonList(szlResponseParameter)),
            new S7PayloadUserData(Collections.singletonList(
                new S7PayloadUserDataItemCpuFunctionReadSzlResponse(
                    DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, szlItems.length, szlItems)))));

        // Read request/response
        dump("05-read-request", new S7MessageRequest(12,
            new S7ParameterReadVarRequest(List.of(
                new S7VarRequestParameterItemAddress(
                    new S7AddressAny(TransportSize.BOOL, 1, 0, MemoryArea.OUTPUTS, 0, (byte) 0)))),
            null));
        dump("06-read-response", new S7MessageResponseData(12,
            new S7ParameterReadVarResponse((short) 1),
            new S7PayloadReadVarResponse(List.of(
                new S7VarPayloadDataItem(DataTransportErrorCode.OK, DataTransportSize.BIT, new byte[]{0x01}))),
            (short) 0, (short) 0));

        // PUT/GET disabled error response
        dump("07-putget-error-response", new S7MessageResponse(12, null, null, (short) 129, (short) 4));
    }

    private void dump(String name, S7Message message) throws Exception {
        WriteBufferXmlBased writeBuffer = new WriteBufferXmlBased();
        writeBuffer.writeMessage(message);
        write(name, writeBuffer.getXmlString());
    }

    private void write(String name, String xml) throws IOException {
        Path directory = Path.of("target", "testsuite-xml");
        Files.createDirectories(directory);
        Files.writeString(directory.resolve(name + ".xml"), xml, StandardCharsets.UTF_8);
    }
}
