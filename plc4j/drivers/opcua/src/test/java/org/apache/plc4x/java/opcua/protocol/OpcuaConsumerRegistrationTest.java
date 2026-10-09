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
package org.apache.plc4x.java.opcua.protocol;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.opcua.OpcuaConnection;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.LinkedHashSet;
import java.util.List;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Unregistering a consumer used to bounce between the registration and the connection until
 * the stack overflowed - see GH-2775.
 */
public class OpcuaConsumerRegistrationTest {

    private OpcuaConnection connection;
    private OpcuaSubscriptionHandle handle;

    @BeforeEach
    void setUp() {
        connection = new OpcuaConnection(new OpcuaConfiguration(), null, mock(AuditLog.class));
        PlcSubscriptionRequest request = mock(PlcSubscriptionRequest.class);
        when(request.getTagNames()).thenReturn(new LinkedHashSet<>(List.of("value-1")));
        handle = new OpcuaSubscriptionHandle(connection, null, request, 1L, 1000L);
    }

    @Test
    void unregisterOfHandleRegistrationReturns() {
        Consumer<PlcSubscriptionEvent> consumer = event -> { };
        PlcConsumerRegistration registration = handle.register(consumer);

        assertDoesNotThrow(registration::unregister);
    }

    @Test
    void unregisterOfConnectionRegistrationReturns() {
        Consumer<PlcSubscriptionEvent> consumer = event -> { };
        PlcConsumerRegistration registration = connection.registerConsumer(consumer, List.of(handle));

        assertDoesNotThrow(registration::unregister);
    }

}
