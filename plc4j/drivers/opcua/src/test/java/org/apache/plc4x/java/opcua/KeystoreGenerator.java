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
import java.math.BigInteger;
import java.security.GeneralSecurityException;
import java.security.Key;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.CertificateEncodingException;
import java.security.cert.X509Certificate;
import java.time.Duration;
import java.time.Instant;
import java.util.Base64;
import java.util.Date;
import java.util.Set;
import org.bouncycastle.asn1.x500.X500NameBuilder;
import org.bouncycastle.asn1.x500.style.BCStyle;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.ExtendedKeyUsage;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.GeneralName;
import org.bouncycastle.asn1.x509.GeneralNames;
import org.bouncycastle.asn1.x509.KeyPurposeId;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;

/**
 * Utility to generate the client / server keystores used by the OPC UA integration test.
 *
 * <p>The certificates mirror what Milo's {@code SelfSignedCertificateBuilder} produced before:
 * a self-signed RSA certificate valid for a year, carrying the OPC UA application URI as a
 * SAN URI entry - the test server refuses to start without one - plus the usual key usages an
 * OPC UA stack expects. It is built with BouncyCastle here so that the test classpath does not
 * need Milo (see {@code src/test/resources/opcua/server/Dockerfile}).
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
        KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
        keyPairGenerator.initialize(length, new SecureRandom());
        KeyPair keyPair = keyPairGenerator.generateKeyPair();

        X500NameBuilder nameBuilder = new X500NameBuilder(BCStyle.INSTANCE)
            .addRDN(BCStyle.CN, commonName)
            .addRDN(BCStyle.O, "Apache Software Foundation")
            .addRDN(BCStyle.OU, "PLC4X")
            .addRDN(BCStyle.L, "PLC4J")
            .addRDN(BCStyle.ST, "CA")
            .addRDN(BCStyle.C, "US");

        Set<String> hostnames = Set.of("127.0.0.1");
        GeneralName[] subjectAltNames = new GeneralName[hostnames.size() + 1];
        subjectAltNames[0] = new GeneralName(GeneralName.uniformResourceIdentifier, applicationUri);
        int index = 1;
        for (String hostname : hostnames) {
            subjectAltNames[index++] = new GeneralName(GeneralName.dNSName, hostname);
        }

        Instant notBefore = Instant.now();
        JcaX509v3CertificateBuilder builder = new JcaX509v3CertificateBuilder(
            nameBuilder.build(),
            new BigInteger(64, new SecureRandom()),
            Date.from(notBefore),
            Date.from(notBefore.plus(Duration.ofDays(365))),
            nameBuilder.build(),
            keyPair.getPublic()
        );
        builder.addExtension(Extension.basicConstraints, false, new BasicConstraints(false));
        builder.addExtension(Extension.keyUsage, false, new KeyUsage(
            KeyUsage.digitalSignature | KeyUsage.nonRepudiation | KeyUsage.keyEncipherment
                | KeyUsage.dataEncipherment | KeyUsage.keyCertSign));
        builder.addExtension(Extension.extendedKeyUsage, false, new ExtendedKeyUsage(
            new KeyPurposeId[]{KeyPurposeId.id_kp_clientAuth, KeyPurposeId.id_kp_serverAuth}));
        builder.addExtension(Extension.subjectAlternativeName, false, new GeneralNames(subjectAltNames));

        X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(
            builder.build(new JcaContentSignerBuilder("SHA256withRSA").build(keyPair.getPrivate())));
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
