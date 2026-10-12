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
package org.apache.plc4x.java.modbus.tcp;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.*;

/**
 * A transaction the device never answers is reported as REQUEST_TIMEOUT for the tags it
 * covered, instead of failing the whole request: the tags of the other blocks (reads) or the
 * other tags (writes) keep the results the device did deliver.
 */
class ModbusTcpRequestTimeoutTest {

    private static final int READ_COILS = 0x01;

    private FakeModbusTcpServer server;

    @AfterEach
    void tearDown() throws Exception {
        if (server != null) {
            server.close();
        }
    }

    @Test
    void readBlockThatTimesOutDoesNotFailTheOtherBlocks() throws Exception {
        // Coils and holding registers can't share a block; leave the coil read unanswered.
        server = new FakeModbusTcpServer(request -> request.functionCode() == READ_COILS
            ? List.of()
            : List.of(request.answer(request.unitId(), new byte[]{0x03, 0x02, 0x00, (byte) 250})));
        try (PlcConnection connection = server.connect(1)) {
            PlcReadResponse response = connection.readRequestBuilder()
                .addTagAddress("slow", "coil:1:BOOL")
                .addTagAddress("fast", "holding-register:1:INT")
                .build().execute().get(5, TimeUnit.SECONDS);

            assertEquals(PlcResponseCode.REQUEST_TIMEOUT, response.getResponseCode("slow"));
            assertEquals(PlcResponseCode.OK, response.getResponseCode("fast"));
            assertEquals(250, response.getInteger("fast"));
        }
    }

    @Test
    void writeThatTimesOutDoesNotFailTheOtherTags() throws Exception {
        // Tags are written one at a time; leave the one writing the value 1 unanswered. A write
        // single register PDU is function code, address (2 bytes), value (2 bytes).
        server = new FakeModbusTcpServer(request -> request.pdu()[4] == 1
            ? List.of()
            : List.of(request.answer(request.unitId(), request.pdu())));
        try (PlcConnection connection = server.connect(1)) {
            PlcWriteResponse response = connection.writeRequestBuilder()
                .addTagAddress("slow", "holding-register:1:INT", 1)
                .addTagAddress("fast", "holding-register:2:INT", 2)
                .build().execute().get(5, TimeUnit.SECONDS);

            assertEquals(PlcResponseCode.REQUEST_TIMEOUT, response.getResponseCode("slow"));
            assertEquals(PlcResponseCode.OK, response.getResponseCode("fast"));
            // Whether the unanswered write took effect is unknown, so it is never re-sent.
            assertEquals(2, server.requests.size());
        }
    }
}
