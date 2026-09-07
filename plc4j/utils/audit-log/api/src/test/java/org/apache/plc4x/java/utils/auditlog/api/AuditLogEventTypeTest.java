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

package org.apache.plc4x.java.utils.auditlog.api;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for AuditLogEventType enum
 */
class AuditLogEventTypeTest {

    @Test
    void testAllEventTypesExist() {
        // Verify all expected event types are present
        AuditLogEventType[] types = AuditLogEventType.values();

        assertEquals(12, types.length, "Expected exactly 12 event types");

        // Verify specific types exist
        assertNotNull(AuditLogEventType.valueOf("CONFIG"));
        assertNotNull(AuditLogEventType.valueOf("SYSTEM"));
        assertNotNull(AuditLogEventType.valueOf("CONNECT"));
        assertNotNull(AuditLogEventType.valueOf("OUTGOING_BYTES"));
        assertNotNull(AuditLogEventType.valueOf("OUTGOING_MESSAGE"));
        assertNotNull(AuditLogEventType.valueOf("INCOMING_BYTES"));
        assertNotNull(AuditLogEventType.valueOf("INCOMING_MESSAGE"));
        assertNotNull(AuditLogEventType.valueOf("API_REQUEST"));
        assertNotNull(AuditLogEventType.valueOf("API_RESPONSE"));
        assertNotNull(AuditLogEventType.valueOf("API_EVENT"));
        assertNotNull(AuditLogEventType.valueOf("CLOSE"));
        assertNotNull(AuditLogEventType.valueOf("ERROR"));
    }

    @Test
    void testEventTypeNames() {
        assertEquals("CONFIG", AuditLogEventType.CONFIG.name());
        assertEquals("SYSTEM", AuditLogEventType.SYSTEM.name());
        assertEquals("CONNECT", AuditLogEventType.CONNECT.name());
        assertEquals("OUTGOING_BYTES", AuditLogEventType.OUTGOING_BYTES.name());
        assertEquals("OUTGOING_MESSAGE", AuditLogEventType.OUTGOING_MESSAGE.name());
        assertEquals("INCOMING_BYTES", AuditLogEventType.INCOMING_BYTES.name());
        assertEquals("INCOMING_MESSAGE", AuditLogEventType.INCOMING_MESSAGE.name());
        assertEquals("API_REQUEST", AuditLogEventType.API_REQUEST.name());
        assertEquals("API_RESPONSE", AuditLogEventType.API_RESPONSE.name());
        assertEquals("API_EVENT", AuditLogEventType.API_EVENT.name());
        assertEquals("CLOSE", AuditLogEventType.CLOSE.name());
        assertEquals("ERROR", AuditLogEventType.ERROR.name());
    }

    @Test
    void testInvalidEventTypeThrowsException() {
        assertThrows(IllegalArgumentException.class, () -> {
            AuditLogEventType.valueOf("INVALID_TYPE");
        });
    }

    @Test
    void testEventTypeOrdinals() {
        // Verify ordinals are sequential starting from 0
        assertEquals(0, AuditLogEventType.CONFIG.ordinal());
        assertEquals(1, AuditLogEventType.SYSTEM.ordinal());
        assertEquals(2, AuditLogEventType.CONNECT.ordinal());
        assertEquals(3, AuditLogEventType.OUTGOING_BYTES.ordinal());
        assertEquals(4, AuditLogEventType.OUTGOING_MESSAGE.ordinal());
        assertEquals(5, AuditLogEventType.INCOMING_BYTES.ordinal());
        assertEquals(6, AuditLogEventType.INCOMING_MESSAGE.ordinal());
        assertEquals(7, AuditLogEventType.API_REQUEST.ordinal());
        assertEquals(8, AuditLogEventType.API_RESPONSE.ordinal());
        assertEquals(9, AuditLogEventType.API_EVENT.ordinal());
        assertEquals(10, AuditLogEventType.CLOSE.ordinal());
        assertEquals(11, AuditLogEventType.ERROR.ordinal());
    }
}
