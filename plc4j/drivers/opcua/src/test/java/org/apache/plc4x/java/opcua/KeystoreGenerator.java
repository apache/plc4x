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

package org.apache.plc4x.java.opcua;

import java.io.IOException;
import java.io.OutputStream;
import java.security.GeneralSecurityException;
import java.security.Key;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.cert.CertificateEncodingException;
import java.security.cert.X509Certificate;
import java.util.Base64;
import java.util.Set;
import org.eclipse.milo.opcua.stack.core.util.SelfSignedCertificateBuilder;
import org.eclipse.milo.opcua.stack.core.util.SelfSignedCertificateGenerator;

/**
 * Utility to generate server certificate - based on Eclipse Milo stuff.
 */
public class KeystoreGenerator {

    private final String password;
    private final KeyStore keyStore;
    private final X509Certificate certificate;

    public KeystoreGenerator(String password, int length, String applicationUri) {
        this(password, length, applicationUri, "server-ai", "Milo Server");
    }

    public KeystoreGenerator(String password, int length, String applicationUri, String entryName, String commonName) {
        this.password = password;
        try {
            this.keyStore = generate(password, length, applicationUri, entryName, commonName);

            Key serverPrivateKey = keyStore.getKey(entryName, password.toCharArray());
            if (serverPrivateKey instanceof PrivateKey) {
                this.certificate = (X509Certificate) keyStore.getCertificate(entryName);
            } else {
                throw new IllegalStateException("Unexpected keystore entry, expected private key");
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private KeyStore generate(String password, int length, String applicationUri, String entryName, String commonName) throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, password.toCharArray());
        KeyPair keyPair = SelfSignedCertificateGenerator.generateRsaKeyPair(length);
        SelfSignedCertificateBuilder builder = (new SelfSignedCertificateBuilder(
            keyPair)).setCommonName(commonName)
            .setOrganization("Apache Software Foundation")
            .setOrganizationalUnit("PLC4X")
            .setLocalityName("PLC4J")
            .setStateName("CA")
            .setCountryCode("US")
            .setApplicationUri(applicationUri);

        Set<String> hostnames = Set.of("127.0.0.1");

        for (String hostname : hostnames) {
            if (hostname.startsWith("\\d+\\.")) {
                builder.addIpAddress(hostname);
            } else {
                builder.addDnsName(hostname);
            }
        }

        X509Certificate certificate = builder.build();
        keyStore.setKeyEntry(entryName, keyPair.getPrivate(), password.toCharArray(), new X509Certificate[]{ certificate });
        return keyStore;
    }

    public KeyStore getKeyStore() {
        return keyStore;
    }

    public X509Certificate getCertificate() {
        return certificate;
    }

    public void writeKeyStoreTo(OutputStream stream) throws IOException, GeneralSecurityException {
        keyStore.store(stream, password.toCharArray());
        stream.flush();
    }

    public void writeCertificateTo(OutputStream stream) throws IOException, CertificateEncodingException {
        String data = "-----BEGIN CERTIFICATE-----\n" +
            Base64.getMimeEncoder(64, "\n".getBytes()).encodeToString(certificate.getEncoded()) + "\n" +
            "-----END CERTIFICATE-----";

        stream.write(data.getBytes());
        stream.flush();
    }

}
