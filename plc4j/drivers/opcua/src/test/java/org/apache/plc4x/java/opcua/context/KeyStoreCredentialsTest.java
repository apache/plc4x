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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.math.BigInteger;
import java.security.GeneralSecurityException;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.Date;
import java.util.Map;
import java.util.Map.Entry;
import javax.security.auth.x500.X500Principal;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.junit.jupiter.api.Test;

/**
 * Reading an identity - a certificate and its private key - out of an in-memory key store. This is
 * how a {@code PlcCertificateAuthentication} handed to {@code getConnection(url, authentication)}
 * becomes the user identity of an OPC UA {@code X509IdentityToken} (GH-1845).
 */
class KeyStoreCredentialsTest {

    private static final char[] PASSWORD = "changeit".toCharArray();
    private static final String SOURCE = "the supplied PlcCertificateAuthentication";
    private static final String PURPOSE = "user certificates";

    @Test
    void usesTheRequestedAlias() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        Entry<PrivateKey, X509Certificate> operator = TestCertificateGenerator.generate(2048, "CN=operator", 3600);
        Entry<PrivateKey, X509Certificate> engineer = TestCertificateGenerator.generate(2048, "CN=engineer", 3600);
        addKeyEntry(keyStore, "operator", operator);
        addKeyEntry(keyStore, "engineer", engineer);

        CertificateKeyPair identity = KeyStoreCredentials.load(keyStore, PASSWORD, "operator", SOURCE, PURPOSE);

        assertEquals(operator.getValue(), identity.getCertificate());
        assertEquals(operator.getKey(), identity.getPrivateKey());
    }

    /**
     * "engineer" sorts before "operator", so a key store whose entries are walked in alias order
     * would hand back the wrong identity if the alias were ignored.
     */
    @Test
    void defaultsToTheFirstEntryHoldingAPrivateKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        keyStore.setCertificateEntry("ca", TestCertificateGenerator.generate(2048, "CN=ca", 3600).getValue());
        Entry<PrivateKey, X509Certificate> user = TestCertificateGenerator.generate(2048, "CN=user", 3600);
        addKeyEntry(keyStore, "user", user);

        CertificateKeyPair identity = KeyStoreCredentials.load(keyStore, PASSWORD, null, SOURCE, PURPOSE);

        assertEquals(user.getValue(), identity.getCertificate());
        assertNotNull(identity.getPrivateKey());
    }

    @Test
    void reportsAnAliasThatDoesNotExist() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        addKeyEntry(keyStore, "user", TestCertificateGenerator.generate(2048, "CN=user", 3600));

        GeneralSecurityException exception = assertThrows(GeneralSecurityException.class,
            () -> KeyStoreCredentials.load(keyStore, PASSWORD, "missing", SOURCE, PURPOSE));

        assertTrue(exception.getMessage().contains("missing"), exception.getMessage());
        // The available aliases are what lets the user fix the mistake.
        assertTrue(exception.getMessage().contains("user"), exception.getMessage());
    }

    @Test
    void reportsAnAliasThatHoldsNoPrivateKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        keyStore.setCertificateEntry("ca", TestCertificateGenerator.generate(2048, "CN=ca", 3600).getValue());

        GeneralSecurityException exception = assertThrows(GeneralSecurityException.class,
            () -> KeyStoreCredentials.load(keyStore, PASSWORD, "ca", SOURCE, PURPOSE));

        assertTrue(exception.getMessage().contains("ca"), exception.getMessage());
        assertTrue(exception.getMessage().contains("private key"), exception.getMessage());
    }

    @Test
    void reportsAKeyStoreWithoutAnyPrivateKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        keyStore.setCertificateEntry("ca", TestCertificateGenerator.generate(2048, "CN=ca", 3600).getValue());

        GeneralSecurityException exception = assertThrows(GeneralSecurityException.class,
            () -> KeyStoreCredentials.load(keyStore, PASSWORD, null, SOURCE, PURPOSE));

        assertTrue(exception.getMessage().contains("no entry with a private key"), exception.getMessage());
        assertTrue(exception.getMessage().contains(SOURCE), exception.getMessage());
    }

    /**
     * OPC UA requires RSA keys; an EC key otherwise fails inside the handshake with an error that
     * says nothing about the cause.
     */
    @Test
    void reportsANonRsaKey() throws Exception {
        KeyStore keyStore = emptyKeyStore();
        addKeyEntry(keyStore, "user", generateEcEntry());

        GeneralSecurityException exception = assertThrows(GeneralSecurityException.class,
            () -> KeyStoreCredentials.load(keyStore, PASSWORD, null, SOURCE, PURPOSE));

        assertTrue(exception.getMessage().contains("EC"), exception.getMessage());
        assertTrue(exception.getMessage().contains("RSA"), exception.getMessage());
        assertTrue(exception.getMessage().contains(PURPOSE), exception.getMessage());
    }

    /**
     * A CA-signed user certificate has to keep its chain, the leaf first - the same rule that
     * applies to the application instance certificate.
     */
    @Test
    void keepsTheCertificateChainWithTheLeafFirst() throws Exception {
        Entry<PrivateKey, X509Certificate> ca = TestCertificateGenerator.generate(2048, "CN=user-ca", 3600);
        Entry<PrivateKey, X509Certificate> user =
            TestCertificateGenerator.generateSignedBy(2048, "CN=user", 3600, ca.getKey(), ca.getValue());
        KeyStore keyStore = emptyKeyStore();
        keyStore.setKeyEntry("user", user.getKey(), PASSWORD,
            new Certificate[]{user.getValue(), ca.getValue()});

        CertificateKeyPair identity = KeyStoreCredentials.load(keyStore, PASSWORD, null, SOURCE, PURPOSE);

        assertEquals(2, identity.getCertificateChain().size());
        assertEquals(user.getValue(), identity.getCertificate());
    }

    private static void addKeyEntry(KeyStore keyStore, String alias, Entry<PrivateKey, X509Certificate> entry)
        throws Exception {
        keyStore.setKeyEntry(alias, entry.getKey(), PASSWORD, new Certificate[]{entry.getValue()});
    }

    private static KeyStore emptyKeyStore() throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, PASSWORD);
        return keyStore;
    }

    private static Entry<PrivateKey, X509Certificate> generateEcEntry() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(256, new SecureRandom());
        KeyPair keyPair = generator.generateKeyPair();

        X500Principal dn = new X500Principal("CN=ec-user");
        X509v3CertificateBuilder builder = new JcaX509v3CertificateBuilder(
            dn, BigInteger.valueOf(1), new Date(System.currentTimeMillis() - 1000),
            new Date(System.currentTimeMillis() + 3_600_000), dn, keyPair.getPublic());
        X509CertificateHolder holder = builder.build(
            new JcaContentSignerBuilder("SHA256withECDSA").build(keyPair.getPrivate()));
        X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(holder);
        return Map.entry(keyPair.getPrivate(), certificate);
    }
}
