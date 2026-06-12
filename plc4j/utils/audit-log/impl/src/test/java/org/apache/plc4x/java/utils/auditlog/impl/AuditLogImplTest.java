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

package org.apache.plc4x.java.utils.auditlog.impl;

import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.apache.plc4x.java.utils.auditlog.api.AuditLogEventType;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for AuditLogImpl
 */
class AuditLogImplTest {

    @TempDir
    Path tempDir;

    private Path auditLogFile;
    private AuditLog auditLog;

    @BeforeEach
    void setUp() {
        auditLogFile = tempDir.resolve("audit-test.log");
    }

    @AfterEach
    void tearDown() {
        if (auditLog != null) {
            auditLog.close();
        }
    }

    @Test
    void testCreateAndWriteSingleEvent() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test-source")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, "test message");
        auditLog.close();

        // Verify file was created and contains the message
        assertTrue(Files.exists(auditLogFile));
        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(1, lines.size());
        assertTrue(lines.get(0).contains("[SYSTEM]"));
        assertTrue(lines.get(0).contains("[test-source]"));
        assertTrue(lines.get(0).contains("test message"));
    }

    @Test
    void testWriteMultipleEvents() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("conn-1")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.CONNECT, "Connection opened");
        auditLog.write(AuditLogEventType.SYSTEM, "System message");
        auditLog.write(AuditLogEventType.CLOSE, "Connection closed");
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(3, lines.size());
        assertTrue(lines.get(0).contains("[CONNECT]"));
        assertTrue(lines.get(1).contains("[SYSTEM]"));
        assertTrue(lines.get(2).contains("[CLOSE]"));
    }

    @Test
    void testAllEventTypesCanBeLogged() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test-source")
            .withAuditLogFile(auditLogFile.toString())
            .build();

        // Write all event types
        for (AuditLogEventType eventType : AuditLogEventType.values()) {
            auditLog.write(eventType, "Test message for " + eventType);
        }
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        // Should have one line for each event type
        assertEquals(AuditLogEventType.values().length, lines.size());

        // Verify each event type appears
        for (AuditLogEventType eventType : AuditLogEventType.values()) {
            boolean found = lines.stream().anyMatch(line -> line.contains("[" + eventType + "]"));
            assertTrue(found, "Event type " + eventType + " should be in the log");
        }
    }

    @Test
    void testTimestampFormat() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, "message");
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        // Check timestamp format (yyyy-MM-dd HH:mm:ss.SSS)
        String line = lines.get(0);
        assertTrue(line.matches("^\\[\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3}].*"),
            "Timestamp should match format [yyyy-MM-dd HH:mm:ss.SSS]");
    }

    @Test
    void testWriteWithSpecialCharacters() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test-source")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM,
            "Message with special chars: \n\t [] {} \"quotes\" 'apostrophes'");
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertFalse(lines.isEmpty());
        String content = String.join("\n", lines);
        assertTrue(content.contains("Message with special chars:"));
    }

    @Test
    void testMultipleInstances() throws IOException {
        Path file1 = tempDir.resolve("audit1.log");
        Path file2 = tempDir.resolve("audit2.log");

        AuditLog log1 = AuditLog.builder()
            .withSource("log1")
            .withAuditLogFile(file1.toString())
            .build();
        AuditLog log2 = AuditLog.builder()
            .withSource("log2")
            .withAuditLogFile(file2.toString())
            .build();

        log1.write(AuditLogEventType.SYSTEM, "message from log1");
        log2.write(AuditLogEventType.API_REQUEST, "message from log2");

        log1.close();
        log2.close();

        // Verify both files exist and have correct content
        List<String> lines1 = Files.readAllLines(file1);
        List<String> lines2 = Files.readAllLines(file2);

        assertEquals(1, lines1.size());
        assertEquals(1, lines2.size());

        assertTrue(lines1.get(0).contains("message from log1"));
        assertTrue(lines2.get(0).contains("message from log2"));
    }

    @Test
    void testWriteAfterClose() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, "before close");
        auditLog.close();

        // Writing after close should not throw, but may not write
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "after close"));
    }

    @Test
    void testEmptySourceAndMessage() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, "");
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(1, lines.size());
        assertTrue(lines.get(0).contains("[SYSTEM]"));
    }

    @Test
    void testLongMessages() throws IOException {
        String longMessage = "x".repeat(10000);

        auditLog = AuditLog.builder()
            .withSource("test")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, longMessage);
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertFalse(lines.isEmpty());
        String content = String.join("\n", lines);
        assertTrue(content.contains("xxx")); // Verify some of the message is there
    }

    @Test
    void testWriteWithJsonObject() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test-source")
            .withAuditLogFile(auditLogFile.toString())
            .build();

        // Create a simple object to serialize
        TestData testData = new TestData("test-name", 42, true);
        auditLog.write(AuditLogEventType.SYSTEM, "Test object", testData);

        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(1, lines.size());
        String line = lines.get(0);
        assertTrue(line.contains("[SYSTEM]"));
        assertTrue(line.contains("[test-source]"));
        assertTrue(line.contains("Test object:"));
        // Verify JSON content
        assertTrue(line.contains("\"name\":\"test-name\""));
        assertTrue(line.contains("\"value\":42"));
        assertTrue(line.contains("\"enabled\":true"));
    }

    @Test
    void testWriteWithNullObject() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("test-source")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, "Null object", null);
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(1, lines.size());
        String line = lines.get(0);
        assertTrue(line.contains("Null object: null"));
    }

    @Test
    void testDisabledWhenNoFileConfigured() {
        auditLog = AuditLog.builder().build();

        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testDisabledWithEmptyFile() {
        auditLog = AuditLog.builder()
            .withAuditLogFile("")
            .build();

        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testDisabledWithWhitespaceOnlyFile() {
        auditLog = AuditLog.builder()
            .withAuditLogFile("   ")
            .build();

        assertFalse(auditLog.isEnabled());
    }

    @Test
    void testEnabledWithValidFile() {
        auditLog = AuditLog.builder()
            .withAuditLogFile(auditLogFile.toString())
            .build();

        assertTrue(auditLog.isEnabled());
    }

    @Test
    void testBuilderWithNoConfig() {
        auditLog = AuditLog.builder().build();

        assertFalse(auditLog.isEnabled());
        // Should not throw when writing to disabled log
        assertDoesNotThrow(() -> auditLog.write(AuditLogEventType.SYSTEM, "message"));
    }

    @Test
    void testPlcMessageEventTypes() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("plc-1")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.OUTGOING_MESSAGE, "Sending data");
        auditLog.write(AuditLogEventType.INCOMING_MESSAGE, "Received data");
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(2, lines.size());
        assertTrue(lines.get(0).contains("[OUTGOING_MESSAGE]"));
        assertTrue(lines.get(1).contains("[INCOMING_MESSAGE]"));
    }

    @Test
    void testApiEventTypes() throws IOException {
        auditLog = AuditLog.builder()
            .withSource("api-1")
            .withAuditLogFile(auditLogFile.toString())
            .build();
        auditLog.write(AuditLogEventType.API_REQUEST, "Read request");
        auditLog.write(AuditLogEventType.API_RESPONSE, "Read response");
        auditLog.close();

        List<String> lines = Files.readAllLines(auditLogFile);

        assertEquals(2, lines.size());
        assertTrue(lines.get(0).contains("[API_REQUEST]"));
        assertTrue(lines.get(1).contains("[API_RESPONSE]"));
    }

    @Test
    void testGetConfig() {
        auditLog = AuditLog.builder()
            .withAuditLogFile(auditLogFile.toString())
            .build();

        assertNotNull(auditLog.getConfig());
        assertEquals(auditLogFile.toString(), auditLog.getConfig().auditLogFile);
    }

    @Test
    void testParentDirectoryCreation() throws IOException {
        // Test that parent directories are created if they don't exist
        Path nestedPath = tempDir.resolve("nested/dir/audit.log");

        auditLog = AuditLog.builder()
            .withSource("test")
            .withAuditLogFile(nestedPath.toString())
            .build();
        auditLog.write(AuditLogEventType.SYSTEM, "message");
        auditLog.close();

        assertTrue(Files.exists(nestedPath));
        List<String> lines = Files.readAllLines(nestedPath);
        assertEquals(1, lines.size());
    }

    // Test data class for JSON serialization
    static class TestData {
        private final String name;
        private final int value;
        private final boolean enabled;

        public TestData(String name, int value, boolean enabled) {
            this.name = name;
            this.value = value;
            this.enabled = enabled;
        }

        public String getName() {
            return name;
        }

        public int getValue() {
            return value;
        }

        public boolean isEnabled() {
            return enabled;
        }
    }
}