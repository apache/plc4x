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

import java.math.BigInteger;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.Date;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.TimeUnit;
import javax.security.auth.x500.X500Principal;
import org.bouncycastle.asn1.x509.BasicConstraints;
import org.bouncycastle.asn1.x509.Extension;
import org.bouncycastle.asn1.x509.KeyUsage;
import org.bouncycastle.cert.CertIOException;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.cert.jcajce.JcaX509CertificateConverter;
import org.bouncycastle.cert.jcajce.JcaX509v3CertificateBuilder;
import org.bouncycastle.operator.OperatorCreationException;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;

public class TestCertificateGenerator {

    /**
     * A random serial number that always encodes to the same number of DER bytes.
     * <p>
     * DER drops redundant leading bytes of an INTEGER, so a plain {@code Random.nextLong()} serial
     * encodes to 8 bytes most of the time but occasionally to 7 or fewer - which changes the
     * certificate's total DER length. Test expectations that contain a certificate's size (the
     * OPC UA asymmetric security header, see ChunkFactoryTest) then fail by one byte every few
     * hundred runs. Forcing the value to exactly 63 bits keeps it random and unique while pinning
     * the encoded length at 8 bytes.
     */
    private static BigInteger serialNumber() {
        return new BigInteger(62, new SecureRandom()).setBit(62);
    }

    /**
     * Generates a certificate that can sign others and act as a trust anchor.
     *
     * <p>Path validation refuses to treat a certificate without basic constraints as a certificate
     * authority, so without this a test cannot build a chain at all - and the chain is the case
     * worth testing, since a leaf on its own tells you nothing about how intermediates are
     * handled.</p>
     *
     * @param issuerKey         the key of the issuer, or null for a self-signed root
     * @param issuerCertificate the issuer's certificate, or null for a self-signed root
     */
    public static Entry<PrivateKey, X509Certificate> generateCa(int keySize, String dn, long validitySec,
        PrivateKey issuerKey, X509Certificate issuerCertificate) {
        try {
            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            keyPairGenerator.initialize(keySize, new SecureRandom());
            KeyPair keyPair = keyPairGenerator.generateKeyPair();

            boolean selfSigned = issuerCertificate == null;
            X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(
                selfSigned ? new X500Principal(dn) : issuerCertificate.getSubjectX500Principal(),
                serialNumber(),
                new Date(),
                new Date(System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(validitySec)),
                new X500Principal(dn),
                keyPair.getPublic()
            );
            certGen.addExtension(Extension.basicConstraints, true, new BasicConstraints(true));
            certGen.addExtension(Extension.keyUsage, true,
                new KeyUsage(KeyUsage.keyCertSign | KeyUsage.cRLSign | KeyUsage.digitalSignature));

            X509CertificateHolder cert = certGen.build(new JcaContentSignerBuilder("SHA256withRSA")
                .build(selfSigned ? keyPair.getPrivate() : issuerKey));

            X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(cert);
            return Map.entry(keyPair.getPrivate(), certificate);
        } catch (CertificateException | CertIOException | NoSuchAlgorithmException
                 | OperatorCreationException e) {
            throw new RuntimeException("Could not initialize test - CA generation failed", e);
        }
    }

    /**
     * Generates a certificate signed by the given issuer, so tests can build a real certificate
     * chain. A self-signed certificate never exercises the code paths that deal with chains.
     */
    public static Entry<PrivateKey, X509Certificate> generateSignedBy(int keySize, String dn, long validitySec,
        PrivateKey issuerKey, X509Certificate issuerCertificate) {
        try {
            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            keyPairGenerator.initialize(keySize, new SecureRandom());
            KeyPair keyPair = keyPairGenerator.generateKeyPair();

            X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(
                issuerCertificate.getSubjectX500Principal(),
                serialNumber(),
                new Date(),
                new Date(System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(validitySec)),
                new X500Principal(dn),
                keyPair.getPublic()
            );
            X509CertificateHolder cert = certGen.build(new JcaContentSignerBuilder("SHA256withRSA")
                .build(issuerKey));

            X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(cert);
            return Map.entry(keyPair.getPrivate(), certificate);
        } catch (CertificateException | NoSuchAlgorithmException | OperatorCreationException e) {
            throw new RuntimeException("Could not initialize test - certificate generation failed", e);
        }
    }

    public static Entry<PrivateKey, X509Certificate> generate(int keySize, String dn, long validitySec) {
        try {
            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            keyPairGenerator.initialize(keySize, new SecureRandom());
            KeyPair keyPair = keyPairGenerator.generateKeyPair();

            X509v3CertificateBuilder certGen = new JcaX509v3CertificateBuilder(
                new X500Principal(dn),
                serialNumber(),
                new Date(),
                new Date(System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(validitySec)),
                new X500Principal(dn),
                keyPair.getPublic()
            );
            X509CertificateHolder cert = certGen.build(new JcaContentSignerBuilder("SHA256withRSA")
                .build(keyPair.getPrivate()));

            X509Certificate certificate = new JcaX509CertificateConverter().getCertificate(cert);
            return Map.entry(keyPair.getPrivate(), certificate);
        } catch (CertificateException | NoSuchAlgorithmException | OperatorCreationException e) {
            throw new RuntimeException("Could not initialize test - certificate generation failed");
        }
    }

}
