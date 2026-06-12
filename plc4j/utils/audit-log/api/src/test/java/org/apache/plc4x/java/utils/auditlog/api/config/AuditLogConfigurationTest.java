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

package org.apache.plc4x.java.utils.auditlog.api.config;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for AuditLogConfiguration
 */
class AuditLogConfigurationTest {

    @Test
    void testIsEnabledWhenFileIsSet() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/tmp/audit.log";

        assertTrue(config.isEnabled());
    }

    @Test
    void testIsDisabledWhenFileIsNull() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = null;

        assertFalse(config.isEnabled());
    }

    @Test
    void testIsDisabledWhenFileIsEmpty() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "";

        assertFalse(config.isEnabled());
    }

    @Test
    void testIsDisabledWhenFileIsWhitespace() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "   ";

        assertFalse(config.isEnabled());
    }

    @Test
    void testIsDisabledWhenFileIsTabsAndNewlines() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "\t\n\r";

        assertFalse(config.isEnabled());
    }

    @Test
    void testIsEnabledWithRelativePath() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "audit.log";

        assertTrue(config.isEnabled());
    }

    @Test
    void testIsEnabledWithAbsolutePath() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/var/log/audit.log";

        assertTrue(config.isEnabled());
    }

    @Test
    void testIsEnabledWithNestedPath() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/path/to/nested/dir/audit.log";

        assertTrue(config.isEnabled());
    }

    @Test
    void testDefaultConfiguration() {
        AuditLogConfiguration config = new AuditLogConfiguration();

        // By default, auditLogFile is null
        assertNull(config.auditLogFile);
        assertFalse(config.isEnabled());
    }

    @Test
    void testConfigurationImplementsInterface() {
        AuditLogConfiguration config = new AuditLogConfiguration();

        // Verify it implements the Configuration interface
        assertTrue(config instanceof org.apache.plc4x.java.spi.config.Configuration);
    }

    @Test
    void testFilePathWithSpaces() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/path/with spaces/audit.log";

        assertTrue(config.isEnabled());
    }

    @Test
    void testFilePathWithSpecialCharacters() {
        AuditLogConfiguration config = new AuditLogConfiguration();
        config.auditLogFile = "/path/with-special_chars.123/audit.log";

        assertTrue(config.isEnabled());
    }
}
