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
package org.apache.plc4x.java.ctrlx.readwrite;

import org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The driver used to trust one certificate that ships inside its own jar, and to accept it for any
 * host. That certificate identifies nobody - anything holding it and its key was trusted - and
 * because it was the only anchor, a device carrying a properly issued certificate could not be
 * reached at all. Both of those are now choices with names.
 */
class CtrlXTrustConfigurationTest {

    private static CtrlXConfiguration configFrom(String params) {
        return new ConfigurationFactory().createConfiguration(CtrlXConfiguration.class, params);
    }

    @Test
    void byDefaultNothingUnusualIsTrusted() {
        CtrlXConfiguration configuration = configFrom("");
        assertFalse(configuration.isAllowFactoryDefaultCertificate(),
            "the certificate in the jar must not be trusted unless asked for");
        assertFalse(configuration.isIgnoreCommonName(),
            "a certificate must be checked against the host it was presented for");
        assertNull(configuration.getTrustStoreFile());
        assertNull(configuration.getServerCertificateFile());
    }

    @Test
    void theFactoryCertificateCanStillBeAskedFor() {
        CtrlXConfiguration configuration = configFrom("allow-factory-default-certificate=true");
        assertTrue(configuration.isAllowFactoryDefaultCertificate());
    }

    @Test
    void aDevicesOwnCertificateCanBeNamed() {
        CtrlXConfiguration configuration =
            configFrom("server-certificate-file=/etc/plc4x/device.pem");
        assertEquals("/etc/plc4x/device.pem", configuration.getServerCertificateFile());
    }

    @Test
    void aTrustStoreCanBeNamedTheSameWayAsElsewhere() {
        CtrlXConfiguration configuration = configFrom(
            "tls.trust-store=/etc/plc4x/trust.p12&tls.trust-store-password=secret&tls.trust-store-type=JKS");
        assertEquals("/etc/plc4x/trust.p12", configuration.getTrustStoreFile());
        assertEquals("secret", configuration.getTrustStorePassword());
        assertEquals("JKS", configuration.getTrustStoreType());
    }

    @Test
    void theTrustStoreTypeDefaultsToTheSameAsTheOtherDrivers() {
        CtrlXConfiguration configuration = configFrom("tls.trust-store=/etc/plc4x/trust.p12");
        assertEquals("PKCS12", configuration.getTrustStoreType());
    }
}
