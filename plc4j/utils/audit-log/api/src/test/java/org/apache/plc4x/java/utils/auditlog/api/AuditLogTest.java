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

import org.apache.plc4x.java.utils.auditlog.api.config.AuditLogConfiguration;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for AuditLog API
 */
class AuditLogTest {

    // ==================== Builder Pattern Tests ====================

    @Test
    void testBuilderWithNoConfigReturnsNoOp() {
        AuditLog auditLog = AuditLog.builder().build();

        assertNotNull(auditLog);
        assertNotNull(auditLog.getConfig());
        assertFalse(auditLog.isEnabled());

        // Should not throw when calling methods
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message"));
        assertDoesNotThrow(auditLog::close);
    }

    @Test
    void testBuilderWithNullDriverTestsuiteFileReturnsNoOp() {
        AuditLog auditLog = AuditLog.builder()
            .withAuditLogFile(null)
            .build();

        assertNotNull(auditLog);
        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testBuilderWithEmptyDriverTestsuiteFileReturnsNoOp() {
        AuditLog auditLog = AuditLog.builder()
            .withAuditLogFile("")
            .build();

        assertNotNull(auditLog);
        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testBuilderWithWhitespaceDriverTestsuiteFileReturnsNoOp() {
        AuditLog auditLog = AuditLog.builder()
            .withAuditLogFile("   ")
            .build();

        assertNotNull(auditLog);
        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testBuilderWithDriverTestsuiteFile() {
        AuditLog auditLog = AuditLog.builder()
            .withAuditLogFile("/tmp/test-audit.log")
            .build();

        assertNotNull(auditLog);
        // Config says enabled, but we get NoOp since impl is not on classpath
        assertTrue(auditLog.isEnabled());

        // Should not throw when calling methods
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message"));
        assertDoesNotThrow(auditLog::close);
    }

    @Test
    void testBuilderWithConfiguration() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/tmp/test.log";

        AuditLog auditLog = AuditLog.builder()
            .withConfiguration(config)
            .build();

        assertNotNull(auditLog);
        assertTrue(auditLog.isEnabled());
        assertEquals("/tmp/test.log", auditLog.getConfig().auditLogFile);
    }

    @Test
    void testBuilderWithNullConfiguration() {
        AuditLog auditLog = AuditLog.builder()
            .withConfiguration(null)
            .build();

        assertNotNull(auditLog);
        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testBuilderChaining() {
        // Test that builder methods return the builder for chaining
        AuditLog.Builder builder = AuditLog.builder();
        assertSame(builder, builder.withAuditLogFile("/tmp/test.log"));
        assertSame(builder, builder.withConfiguration(null));
    }

    @Test
    void testBuilderOverridesValues() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/original/path.log";

        AuditLog auditLog = AuditLog.builder()
            .withConfiguration(config)
            .withAuditLogFile("/new/path.log")
            .build();

        assertNotNull(auditLog);
        assertEquals("/new/path.log", auditLog.getConfig().auditLogFile);
    }

    @Test
    void testBuilderMultipleBuilds() {
        AuditLog.Builder builder = AuditLog.builder()
            .withAuditLogFile("/tmp/test.log");

        AuditLog auditLog1 = builder.build();
        AuditLog auditLog2 = builder.build();

        assertNotNull(auditLog1);
        assertNotNull(auditLog2);
        assertNotSame(auditLog1, auditLog2);
    }

    // ==================== Common Functionality Tests ====================

    @Test
    void testNoOpImplementationDoesNotFail() {
        AuditLog auditLog = AuditLog.builder().build();

        // Test all event types
        for (AuditLogEventType eventType : AuditLogEventType.values()) {
            assertDoesNotThrow(() -> auditLog.write(eventType, "test message"));
        }

        assertDoesNotThrow(auditLog::close);
    }

    @Test
    void testGetConfigReturnsCorrectConfig() {
        AuditLog auditLog = AuditLog.builder()
            .withAuditLogFile("/tmp/test.log")
            .build();

        assertNotNull(auditLog.getConfig());
        assertEquals("/tmp/test.log", auditLog.getConfig().auditLogFile);
    }

    @Test
    void testWriteWithNullValues() {
        AuditLog auditLog = AuditLog.builder().build();

        // Should handle null values gracefully
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, null));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.API_REQUEST, null));
    }

    @Test
    void testWriteWithEmptyStrings() {
        AuditLog auditLog = AuditLog.builder().build();

        // Should handle empty strings gracefully
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, ""));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.API_RESPONSE, ""));
    }

    @Test
    void testMultipleCallsToClose() {
        AuditLog auditLog = AuditLog.builder().build();

        // Should handle multiple calls to close without issues
        assertDoesNotThrow(auditLog::close);
        assertDoesNotThrow(auditLog::close);
        assertDoesNotThrow(auditLog::close);
    }

    @Test
    void testWriteAfterClose() {
        AuditLog auditLog = AuditLog.builder().build();

        auditLog.close();

        // Should still not throw after close (no-op implementation)
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message"));
    }

    @Test
    void testWriteWithDataObject() {
        AuditLog auditLog = AuditLog.builder().build();

        // Should handle data object gracefully
        Object data = new Object();
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.CONNECT, "message", data));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message", null));
    }

    @Test
    void testIsEnabledReturnsFalseForNoOp() {
        AuditLog auditLog = AuditLog.builder().build();

        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testAllEventTypesCanBeWritten() {
        AuditLog auditLog = AuditLog.builder().build();

        // Test all event types can be written without errors
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "system message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.CONFIG, "config message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.CONNECT, "connect message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.OUTGOING_BYTES, "outgoing bytes"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.OUTGOING_MESSAGE, "outgoing message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.INCOMING_BYTES, "incoming bytes"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.INCOMING_MESSAGE, "incoming message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.API_REQUEST, "request message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.API_RESPONSE, "response message"));
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.CLOSE, "close message"));
    }

    @Test
    void testWriteWithDataObjectNullData() {
        AuditLog auditLog = AuditLog.builder().build();

        // Test with null data - should convert to "null" string
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message", null));
    }

    @Test
    void testWriteWithDataObjectNonNullData() {
        AuditLog auditLog = AuditLog.builder().build();

        // Test with non-null data - should call toString()
        Object data = new Object() {
            @Override
            public String toString() {
                return "custom-data";
            }
        };
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message", data));
    }

    @Test
    void testNoOpAuditLogClose() {
        AuditLog auditLog = AuditLog.builder().build();

        // Test that close can be called multiple times without issue
        auditLog.close();
        auditLog.close();

        // Can still write after close (no-op)
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message"));
    }

    @Test
    void testNoOpAuditLogWriteWithAllEventTypes() {
        AuditLog auditLog = AuditLog.builder().build();

        // Write with data object for each event type
        for (AuditLogEventType eventType : AuditLogEventType.values()) {
            assertDoesNotThrow(() -> auditLog.write(eventType, "message", new Object()));
        }
    }

    @Test
    void testBuilderWithSource() {
        AuditLog auditLog = AuditLog.builder()
            .withSource("test-connection")
            .build();

        assertNotNull(auditLog);
        assertEquals("test-connection", auditLog.getSource());
    }

    @Test
    void testBuilderWithNullSource() {
        AuditLog auditLog = AuditLog.builder()
            .withSource(null)
            .build();

        assertNotNull(auditLog);
        assertEquals("unknown", auditLog.getSource());
    }

    @Test
    void testBuilderWithSourceAndConfig() {
        AuditLog auditLog = AuditLog.builder()
            .withSource("my-source")
            .withAuditLogFile("/tmp/test.log")
            .build();

        assertNotNull(auditLog);
        assertEquals("my-source", auditLog.getSource());
        assertTrue(auditLog.isEnabled());
    }

    @Test
    void testGetSourceDefaultValue() {
        AuditLog auditLog = AuditLog.builder().build();

        assertEquals("unknown", auditLog.getSource());
    }
}
