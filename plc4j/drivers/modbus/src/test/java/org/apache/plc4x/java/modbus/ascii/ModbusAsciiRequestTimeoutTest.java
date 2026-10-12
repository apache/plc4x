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
package org.apache.plc4x.java.modbus.ascii;

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.modbus.ascii.config.ModbusAsciiConfiguration;
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.transports.api.AsyncTransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.Test;

import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

/**
 * A transaction the device never answers is reported as REQUEST_TIMEOUT for the tags it
 * covered, instead of failing the whole request.
 */
class ModbusAsciiRequestTimeoutTest {

    @Test
    void readThatTimesOutIsReportedAsRequestTimeout() throws Exception {
        AsyncTransportInstance<?> transport = silentTransport();
        ModbusAsciiConnection connection = newConnectedConnection(transport);

        PlcReadResponse response = connection.readRequestBuilder()
            .addTagAddress("value", "4x00001:INT")
            .build()
            .execute()
            .get(2, TimeUnit.SECONDS);

        assertEquals(PlcResponseCode.REQUEST_TIMEOUT, response.getResponseCode("value"));
    }

    @Test
    void writeThatTimesOutIsReportedAsRequestTimeout() throws Exception {
        AsyncTransportInstance<?> transport = silentTransport();
        ModbusAsciiConnection connection = newConnectedConnection(transport);

        PlcWriteResponse response = connection.writeRequestBuilder()
            .addTagAddress("value", "4x00001:INT", 42)
            .build()
            .execute()
            .get(2, TimeUnit.SECONDS);

        assertEquals(PlcResponseCode.REQUEST_TIMEOUT, response.getResponseCode("value"));
        // Whether the unanswered write took effect is unknown, so it is never re-sent.
        verify(transport, times(1)).write(any());
    }

    /**
     * A transport that accepts every frame and never delivers a byte back.
     */
    private static AsyncTransportInstance<?> silentTransport() {
        AsyncTransportInstance<?> transport = mock(AsyncTransportInstance.class);
        when(transport.isOpen()).thenReturn(true);
        return transport;
    }

    private static ModbusAsciiConnection newConnectedConnection(AsyncTransportInstance<?> transport) throws Exception {
        ModbusAsciiConfiguration config = new ModbusAsciiConfiguration();
        config.setRequestTimeout(200);
        config.setDefaultUnitIdentifier(1);
        config.setPingAddress("4x00001:BOOL");
        config.setDefaultPayloadByteOrder(ModbusByteOrder.BIG_ENDIAN);
        config.setMaxCoilsPerRequest(2000);
        config.setMaxRegistersPerRequest(125);

        AuditLog auditLog = mock(AuditLog.class);
        when(auditLog.isEnabled()).thenReturn(false);

        ModbusAsciiConnection connection = new ModbusAsciiConnection(config, transport, auditLog);
        connection.connect();
        return connection;
    }
}
