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
package org.apache.plc4x.java.opcua.context;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.math.BigInteger;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The certificate the driver generates when no key store is configured is what most servers see on
 * a first connection attempt, so it has to satisfy the checks a server performs. It used to be
 * signed with SHA-1, flagged as a CA and carried an over-long, often negative serial number, which
 * servers with a modern security baseline reject (see GH-2127).
 */
class CertificateGeneratorTest {

    private static X509Certificate certificate;

    @BeforeAll
    static void generate() {
        CertificateKeyPair keyPair = CertificateGenerator.generateCertificate();
        assertNotNull(keyPair, "certificate generation failed");
        certificate = keyPair.getCertificate();
    }

    @Test
    void isSignedWithSha256() {
        assertEquals("SHA256WITHRSA", certificate.getSigAlgName().toUpperCase(),
            "SHA-1 signatures are rejected by the SHA-256 based security policies");
    }

    @Test
    void hasAUsableKey() {
        assertInstanceOfRsa();
        assertEquals(CertificateGenerator.DEFAULT_KEY_SIZE,
            ((RSAPublicKey) certificate.getPublicKey()).getModulus().bitLength());
    }

    /**
     * Servers can demand a larger key than the default; 'generated-key-size' asks for one.
     */
    @Test
    void honoursTheRequestedKeySize() {
        CertificateKeyPair keyPair = CertificateGenerator.generateCertificate(4096);

        assertNotNull(keyPair);
        assertEquals(4096, ((RSAPublicKey) keyPair.getCertificate().getPublicKey()).getModulus().bitLength());
        assertEquals("SHA256WITHRSA", keyPair.getCertificate().getSigAlgName().toUpperCase());
    }

    /**
     * X.509 requires a positive serial number of at most 20 octets.
     */
    @Test
    void hasASpecConformingSerialNumber() {
        BigInteger serial = certificate.getSerialNumber();

        assertEquals(1, serial.signum(), "the serial number has to be positive");
        assertTrue(serial.toByteArray().length <= 20,
            "the serial number must not exceed 20 octets, was " + serial.toByteArray().length);
    }

    /**
     * This is an application instance certificate; it signs nothing but itself.
     */
    @Test
    void isNotMarkedAsACertificateAuthority() {
        assertEquals(-1, certificate.getBasicConstraints(), "the certificate must not claim to be a CA");

        boolean[] keyUsage = certificate.getKeyUsage();
        assertNotNull(keyUsage);
        assertFalse(keyUsage[5], "keyCertSign has no business on an application instance certificate");
        assertTrue(keyUsage[0], "digitalSignature is required");
        assertTrue(keyUsage[2], "keyEncipherment is required");
        assertTrue(keyUsage[3], "dataEncipherment is required");
    }

    /**
     * Servers look certificates up by their subject key identifier.
     */
    @Test
    void carriesASubjectKeyIdentifier() {
        assertNotNull(certificate.getExtensionValue("2.5.29.14"), "subject key identifier is missing");
    }

    @Test
    void carriesTheApplicationUriAsSubjectAlternativeName() throws Exception {
        assertNotNull(certificate.getSubjectAlternativeNames(),
            "the application URI has to be present as a subject alternative name");
    }

    private void assertInstanceOfRsa() {
        assertTrue(certificate.getPublicKey() instanceof RSAPublicKey,
            "expected an RSA key, got " + certificate.getPublicKey().getAlgorithm());
    }
}
