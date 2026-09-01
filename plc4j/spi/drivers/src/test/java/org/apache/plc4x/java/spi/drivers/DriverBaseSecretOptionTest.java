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
package org.apache.plc4x.java.spi.drivers;

import org.apache.plc4x.java.api.metadata.Option;
import org.apache.plc4x.java.api.metadata.OptionMetadata;
import org.apache.plc4x.java.api.types.OptionType;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.Required;
import org.apache.plc4x.java.spi.config.annotations.Secret;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.List;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The {@link Secret} marking on a configuration field has to survive the trip through
 * {@code DriverBase}'s reflective option extraction and come out as {@link Option#isSecret()}.
 * That is the whole point of the annotation: it is declared once beside the parameter, and every
 * consumer of the driver metadata reads it from there rather than re-deciding by name.
 *
 * @see ConnectionStringRedactorTest for the redaction that the same marking drives
 */
class DriverBaseSecretOptionTest {

    /** Marked on a superclass, to prove the hierarchy walk carries the flag too. */
    static class BaseConfig implements Configuration {
        @Secret
        @ConfigurationParameter("inherited-token")
        public String inheritedToken;

        @ConfigurationParameter("inherited-host")
        public String inheritedHost;
    }

    static class DriverConfig extends BaseConfig {
        @Secret
        @Required
        @Description("the password")
        @ConfigurationParameter("password")
        public String password;

        @ConfigurationParameter("username")
        public String username;

        @Secret
        @ConfigurationParameter("psk-key")
        public String pskKey;

        // Deliberately unmarked: the identity says *which* key was refused, and hiding it costs
        // the operator the diagnosis while protecting nothing.
        @ConfigurationParameter("psk-identity")
        public String pskIdentity;
    }

    /** Minimal {@link DriverBase} that exists only to expose {@link DriverConfig} as metadata. */
    static final class StubDriver extends DriverBase {
        @Override public String getProtocolCode() { return "stub"; }
        @Override public String getProtocolName() { return "Stub"; }
        @Override public Optional<String> getDefaultTransportCode() { return Optional.empty(); }
        @Override public List<String> getSupportedTransportCodes() { return List.of(); }
        @Override protected Class<? extends Configuration> getConfigurationClass() { return DriverConfig.class; }
        @Override protected ConnectionBase<?> getConnection(Configuration c, TransportInstance<?> t, AuditLog a) {
            throw new UnsupportedOperationException("test stub does not connect");
        }
    }

    private static List<Option> options() {
        OptionMetadata metadata = new StubDriver().getMetadata()
            .getProtocolConfigurationOptionMetadata()
            .orElseThrow(() -> new AssertionError("driver must expose protocol configuration options"));
        return metadata.getOptions();
    }

    private static Option option(String key) {
        return options().stream()
            .filter(o -> key.equals(o.getKey()))
            .findFirst()
            .orElseThrow(() -> new AssertionError("no option named " + key + " in " + options()));
    }

    @Test
    @DisplayName("a marked field becomes a secret option")
    void markedFieldIsSecret() {
        assertTrue(option("password").isSecret());
    }

    @Test
    @DisplayName("an unmarked field does not")
    void unmarkedFieldIsNotSecret() {
        assertFalse(option("username").isSecret());
    }

    @Test
    @DisplayName("a marked pre-shared key is secret, its identity is not")
    void marksThePreSharedKeyButNotItsIdentity() {
        // The pair that motivates declaring this at the field: no name-based rule separates them.
        assertTrue(option("psk-key").isSecret());
        assertFalse(option("psk-identity").isSecret());
    }

    @Test
    @DisplayName("a marking on an inherited field survives the hierarchy walk")
    void inheritedMarkingIsCarried() {
        assertTrue(option("inherited-token").isSecret());
        assertFalse(option("inherited-host").isSecret());
    }

    @Test
    @DisplayName("marking a field changes nothing else about its option")
    void markingLeavesTheOtherAttributesAlone() {
        // The new annotation is read alongside the existing ones, not instead of them.
        Option password = option("password");
        assertEquals(OptionType.STRING, password.getType());
        assertEquals("the password", password.getDescription());
        assertTrue(password.isRequired());
    }

    @Test
    @DisplayName("every declared parameter is still extracted")
    void extractionStillReturnsEveryParameter() {
        assertEquals(6, options().size(), "options were: " + options());
    }
}
