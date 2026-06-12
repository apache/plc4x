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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.s7.S7CotpConnection;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseRequest;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;

import java.util.LinkedHashMap;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.util.concurrent.CompletableFuture;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;

/**
 * Integration-shape test for {@link S7CotpConnection#onBrowse}: when the SZL probe at
 * connect time hasn't run (or marked UserData unsupported), browse must complete with a
 * structured failure response — not throw and not return null.
 */
class S7BrowseConnectionTest {

    @Test
    void browseReturnsUnsupportedWhenUserDataServicesAbsent() throws Exception {
        S7Configuration cfg = new S7Configuration();
        cfg.setMaxAmqCallee(4);
        TransportInstance<?> transport = Mockito.mock(TransportInstance.class);
        AuditLog auditLog = Mockito.mock(AuditLog.class);
        S7CotpConnection conn = new S7CotpConnection(cfg, transport, auditLog);
        // userDataServicesSupported defaults to false — same state as a LOGO/non-SZL device.
        assertFalse(conn.getMetadata().isBrowseSupported());

        // S7's tag handler doesn't implement parseQuery yet — this round only adds the
        // browse pipeline, not the query DSL — so we hand-build the request.
        LinkedHashMap<String, PlcQuery> queries = new LinkedHashMap<>();
        queries.put("everything", Mockito.mock(PlcQuery.class));
        PlcBrowseRequest browseRequest = new DefaultPlcBrowseRequest(conn, queries);
        CompletableFuture<PlcBrowseResponse> future = conn.browse(browseRequest);

        PlcBrowseResponse response = future.get();
        assertNotNull(response);
        assertEquals(PlcResponseCode.UNSUPPORTED, response.getResponseCode("everything"));
        assertEquals(0, response.getValues("everything").size());
    }
}
