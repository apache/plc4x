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

import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.Test;

import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Checks a real, CA-issued key store against the loading logic. Certificates from a public CA
 * differ from the self-signed ones the rest of the suite generates: they carry a chain, they can
 * hold extra entries, and their extensions come from whatever the CA issues rather than from our
 * own generator - which is the shape reported as failing in GH-2013 and GH-2196.
 * <p>
 * Such a key store contains private key material, so it cannot live in the repository. Point this
 * test at one instead; without the properties it skips.
 *
 * <pre>
 * export PLC4X_OPCUA_TEST_KEYSTORE_PASSWORD='secret'
 * mvn test -Dtest=RealKeyStoreManualTest \
 *   -Dplc4x.opcua.test.keystore=/absolute/path/to/identity.p12 \
 *   [-Dplc4x.opcua.test.keystore-type=pkcs12]
 * </pre>
 * <p>
 * The password may also be given as {@code -Dplc4x.opcua.test.keystore-password}, but the
 * environment variable is preferable: a password on the command line lands in the shell history
 * and in the process list, and unquoted special characters in it are mangled by the shell before
 * Maven ever sees them.
 */
class RealKeyStoreManualTest {

    private static final String KEY_STORE_PROPERTY = "plc4x.opcua.test.keystore";
    private static final String PASSWORD_PROPERTY = "plc4x.opcua.test.keystore-password";
    private static final String PASSWORD_ENV = "PLC4X_OPCUA_TEST_KEYSTORE_PASSWORD";
    private static final String TYPE_PROPERTY = "plc4x.opcua.test.keystore-type";

    @Test
    void loadsTheClientIdentityFromARealKeyStore() throws Exception {
        String keyStoreFile = System.getProperty(KEY_STORE_PROPERTY);
        Assumptions.assumeTrue(keyStoreFile != null && !keyStoreFile.isEmpty(),
            "set -D" + KEY_STORE_PROPERTY + " to run this against a real key store");

        OpcuaConfiguration configuration = new OpcuaConfiguration();
        set(configuration, "keyStoreFile", keyStoreFile);
        set(configuration, "keyStoreType", System.getProperty(TYPE_PROPERTY, "pkcs12"));
        set(configuration, "keyStorePassword", password());

        OpcuaDriverContext context = new OpcuaDriverContext();
        context.openKeyStore(configuration);
        CertificateKeyPair identity = context.getCertificateKeyPair();

        assertNotNull(identity.getPrivateKey(), "no private key was selected");
        assertTrue(identity.getCertificate().getPublicKey() instanceof RSAPublicKey,
            "OPC UA requires an RSA key, this key store holds a "
                + identity.getCertificate().getPublicKey().getAlgorithm() + " key");

        List<X509Certificate> chain = identity.getCertificateChain();
        assertEquals(identity.getCertificate(), chain.get(0), "the client certificate has to come first");

        // The bytes that would go into the SenderCertificate field.
        int expected = 0;
        for (X509Certificate certificate : chain) {
            expected += certificate.getEncoded().length;
        }
        assertEquals(expected, identity.getEncodedCertificateChain().length);

        // Deliberately printed, not asserted: the interesting part is what a real CA issues.
        System.out.println("key store entry: " + identity.getCertificate().getSubjectX500Principal());
        System.out.println("issuer         : " + identity.getCertificate().getIssuerX500Principal());
        System.out.println("key            : "
            + ((RSAPublicKey) identity.getCertificate().getPublicKey()).getModulus().bitLength() + " bit "
            + identity.getCertificate().getPublicKey().getAlgorithm());
        System.out.println("signature      : " + identity.getCertificate().getSigAlgName());
        System.out.println("chain length   : " + chain.size()
            + " (" + identity.getEncodedCertificateChain().length + " bytes on the wire)");
        for (X509Certificate certificate : chain) {
            System.out.println("  - " + certificate.getSubjectX500Principal()
                + " signed by " + certificate.getIssuerX500Principal());
        }
    }

    /**
     * The environment variable wins: it survives special characters that a shell would otherwise
     * expand, and keeps the password out of the shell history and the process list.
     */
    private static String password() {
        String fromEnvironment = System.getenv(PASSWORD_ENV);
        if (fromEnvironment != null && !fromEnvironment.isEmpty()) {
            return fromEnvironment;
        }
        return System.getProperty(PASSWORD_PROPERTY, "");
    }

    private static void set(OpcuaConfiguration configuration, String fieldName, Object value) throws Exception {
        java.lang.reflect.Field field = OpcuaConfiguration.class.getDeclaredField(fieldName);
        field.setAccessible(true);
        field.set(configuration, value);
    }
}
