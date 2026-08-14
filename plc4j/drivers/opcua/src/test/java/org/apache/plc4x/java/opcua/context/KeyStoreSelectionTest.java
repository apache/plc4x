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

import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.FileOutputStream;
import java.io.OutputStream;
import java.math.BigInteger;
import java.nio.file.Path;
import java.security.GeneralSecurityException;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.X509Certificate;
import java.util.Date;
import java.util.Map.Entry;

import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;

import javax.security.auth.x500.X500Principal;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Choosing the client identity out of a key store. The driver used to take whatever alias came
 * first, which breaks for the key stores the certificate tutorial produces - those also contain the
 * signing CA - and gave no useful error when the entry could not be used at all (see the
 * discussion in GH-2013 / GH-2196).
 */
class KeyStoreSelectionTest {

    private static final char[] PASSWORD = "changeit".toCharArray();

    @TempDir
    Path tempDir;

    @Test
    void picksTheEntryThatHasAPrivateKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        // A trusted certificate entry sorting before the key entry - "ca" < "client".
        keyStore.setCertificateEntry("ca", TestCertificateGenerator.generate(2048, "CN=ca", 3600).getValue());
        Entry<PrivateKey, X509Certificate> client = TestCertificateGenerator.generate(2048, "CN=client", 3600);
        keyStore.setKeyEntry("client", client.getKey(), PASSWORD, new java.security.cert.Certificate[]{client.getValue()});

        CertificateKeyPair selected = load(keyStore);

        assertEquals(client.getValue(), selected.getCertificate());
        assertNotNull(selected.getKeyPair().getPrivate());
    }

    @Test
    void reportsAKeyStoreWithoutAnyPrivateKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        keyStore.setCertificateEntry("ca", TestCertificateGenerator.generate(2048, "CN=ca", 3600).getValue());

        GeneralSecurityException exception = assertThrows(GeneralSecurityException.class, () -> load(keyStore));

        assertTrue(exception.getMessage().contains("no entry with a private key"), exception.getMessage());
    }

    /**
     * OPC UA application instance certificates have to be RSA. The certificate tutorial offers
     * RSA, DSA and EC, and picking one of the latter two produced an opaque handshake failure.
     */
    @Test
    void reportsANonRsaKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        Entry<PrivateKey, X509Certificate> ecEntry = generateEcEntry();
        keyStore.setKeyEntry("client", ecEntry.getKey(), PASSWORD,
            new java.security.cert.Certificate[]{ecEntry.getValue()});

        GeneralSecurityException exception = assertThrows(GeneralSecurityException.class, () -> load(keyStore));

        assertTrue(exception.getMessage().contains("EC"), exception.getMessage());
        assertTrue(exception.getMessage().contains("RSA"), exception.getMessage());
    }

    private CertificateKeyPair load(KeyStore keyStore) throws Exception {
        Path file = tempDir.resolve("keystore-" + System.nanoTime() + ".p12");
        try (OutputStream outputStream = new FileOutputStream(file.toFile())) {
            keyStore.store(outputStream, PASSWORD);
        }

        // The configuration has no setters - the framework injects the fields.
        OpcuaConfiguration configuration = new OpcuaConfiguration();
        set(configuration, "keyStoreFile", file.toAbsolutePath().toString());
        set(configuration, "keyStoreType", "pkcs12");
        set(configuration, "keyStorePassword", new String(PASSWORD));

        OpcuaDriverContext context = new OpcuaDriverContext();
        context.openKeyStore(configuration);
        return context.getCertificateKeyPair();
    }

    private static void set(OpcuaConfiguration configuration, String fieldName, Object value) throws Exception {
        java.lang.reflect.Field field = OpcuaConfiguration.class.getDeclaredField(fieldName);
        field.setAccessible(true);
        field.set(configuration, value);
    }

    private KeyStore emptyKeyStore() throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, PASSWORD);
        return keyStore;
    }

    private Entry<PrivateKey, X509Certificate> generateEcEntry() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(256, new SecureRandom());
        KeyPair keyPair = generator.generateKeyPair();

        X500Principal dn = new X500Principal("CN=ec-client");
        X509v3CertificateBuilder builder = new JcaX509v3CertificateBuilder(
            dn, BigInteger.valueOf(1), new Date(System.currentTimeMillis() - 1000),
            new Date(System.currentTimeMillis() + 3_600_000), dn, keyPair.getPublic());
        X509CertificateHolder holder = builder.build(
            new JcaContentSignerBuilder("SHA256withECDSA").build(keyPair.getPrivate()));
        X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(holder);
        return java.util.Map.entry(keyPair.getPrivate(), certificate);
    }
}
