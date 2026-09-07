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
package org.apache.plc4x.java.opcua.config;

import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Field;
import java.security.cert.X509Certificate;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.mock;

class OpcuaConfigurationTest {

    @Test
    void defaultAccessorsReturnEmptyState() {
        OpcuaConfiguration cfg = new OpcuaConfiguration();
        // Exercise every public read-side accessor; values are null/false/0
        // before the Configuration framework wires fields in. The point here
        // is the surface staying stable, not the value semantics.
        assertThat(cfg.getProtocolCode()).isNull();
        assertThat(cfg.getTransportCode()).isNull();
        assertThat(cfg.getTransportConfig()).isNull();
        assertThat(cfg.isDiscovery()).isFalse();
        assertThat(cfg.getUsername()).isNull();
        assertThat(cfg.getPassword()).isNull();
        assertThat(cfg.getSecurityPolicy()).isNull();
        assertThat(cfg.getMessageSecurity()).isNull();
        assertThat(cfg.getKeyStoreFile()).isNull();
        assertThat(cfg.getKeyStoreType()).isNull();
        assertThat(cfg.getKeyStorePassword()).isNull();
        assertThat(cfg.getTrustStoreFile()).isNull();
        assertThat(cfg.getTrustStoreType()).isNull();
        assertThat(cfg.getTrustStorePassword()).isNull();
        assertThat(cfg.getEncodingLimits()).isNull();
        assertThat(cfg.getServerCertificate()).isNull();
        assertThat(cfg.getEndpointHost()).isNull();
        assertThat(cfg.getEndpointPort()).isNull();
        assertThat(cfg.getChannelLifetime()).isZero();
        assertThat(cfg.getSessionTimeout()).isZero();
        assertThat(cfg.getRequestTimeout()).isZero();
        assertThat(cfg.getNegotiationTimeout()).isZero();
        assertThat(cfg.getSubscriptionQueueSize()).isZero();
        assertThat(cfg.getBrowseMaxReferencesPerNode()).isZero();
        assertThat(cfg.getBrowseMaxTotalNodes()).isZero();
        assertThat(cfg.getBrowseMaxDepth()).isZero();
    }

    @Test
    void subscriptionQueueSizeDefaultsToOneAndParsesFromConfig() {
        // An empty configuration must apply the @LongDefaultValue(1),
        // and an explicit value must be parsed under the "subscription-queue-size" parameter name.
        ConfigurationFactory factory = new ConfigurationFactory();
        assertThat(factory.createConfiguration(OpcuaConfiguration.class, "")
            .getSubscriptionQueueSize()).isEqualTo(1L);
        assertThat(factory.createConfiguration(OpcuaConfiguration.class, "subscription-queue-size=10")
            .getSubscriptionQueueSize()).isEqualTo(10L);
    }

    @Test
    void browseBoundsAreConfigurable() {
        // A browse walks whatever tree the server describes, and how large that is cannot be
        // known before walking it - so each of the three ways it can run away has a knob.
        OpcuaConfiguration cfg = new OpcuaConfiguration();
        cfg.setBrowseMaxReferencesPerNode(10);
        cfg.setBrowseMaxTotalNodes(20);
        cfg.setBrowseMaxDepth(3);
        assertThat(cfg.getBrowseMaxReferencesPerNode()).isEqualTo(10);
        assertThat(cfg.getBrowseMaxTotalNodes()).isEqualTo(20);
        assertThat(cfg.getBrowseMaxDepth()).isEqualTo(3);
    }

    @Test
    void setServerCertificateRoundtripsAndToStringRedactsPassword() throws Exception {
        OpcuaConfiguration cfg = new OpcuaConfiguration();
        X509Certificate cert = mock(X509Certificate.class);
        cfg.setServerCertificate(cert);
        assertThat(cfg.getServerCertificate()).isSameAs(cert);

        // Inject a username + password via reflection (no setters; production
        // wiring goes through the Configuration framework). The toString must
        // not leak the password value.
        setField(cfg, "username", "alice");
        setField(cfg, "password", "supersecret");
        String s = cfg.toString();
        assertThat(s).contains("alice").doesNotContain("supersecret");
    }

    private static void setField(Object target, String name, Object value) throws Exception {
        Field f = target.getClass().getDeclaredField(name);
        f.setAccessible(true);
        f.set(target, value);
    }
}
