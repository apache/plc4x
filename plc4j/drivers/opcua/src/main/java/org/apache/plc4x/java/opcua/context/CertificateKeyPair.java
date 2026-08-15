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

import java.security.GeneralSecurityException;
import java.security.PrivateKey;
import org.bouncycastle.asn1.x509.GeneralName;

import java.io.ByteArrayOutputStream;
import java.security.KeyPair;
import java.security.MessageDigest;
import java.security.cert.X509Certificate;
import java.util.Collection;
import java.util.List;
import java.util.Optional;

public class CertificateKeyPair {

    private final KeyPair keyPair;
    private final X509Certificate certificate;
    private final List<X509Certificate> certificateChain;
    private final byte[] encodedCertificateChain;
    private final byte[] thumbprint;

    public CertificateKeyPair(KeyPair keyPair, X509Certificate certificate) throws GeneralSecurityException {
        this(keyPair, List.of(certificate));
    }

    /**
     * @param certificateChain the client certificate followed by the certificates that signed it,
     *                         as it comes out of a key store. Everything but the first element is
     *                         only relevant for a CA-signed certificate.
     */
    public CertificateKeyPair(KeyPair keyPair, List<X509Certificate> certificateChain) throws GeneralSecurityException {
        if (certificateChain == null || certificateChain.isEmpty()) {
            throw new IllegalArgumentException("A certificate chain needs at least the certificate itself");
        }
        this.keyPair = keyPair;
        this.certificate = certificateChain.get(0);
        this.certificateChain = List.copyOf(certificateChain);
        this.encodedCertificateChain = encode(this.certificateChain);
        MessageDigest messageDigest = MessageDigest.getInstance("SHA-1");
        // The thumbprint identifies the sender's own certificate, never the chain.
        this.thumbprint = messageDigest.digest(this.certificate.getEncoded());
    }

    /**
     * The certificate chain, the client's own certificate first.
     */
    public List<X509Certificate> getCertificateChain() {
        return certificateChain;
    }

    /**
     * The bytes that go into the {@code SenderCertificate} field of the asymmetric security
     * header: the DER encoded client certificate, followed by the DER encoding of each CA
     * certificate that signed it. OPC UA Part 6 allows appending the issuing certificates, and
     * servers that cannot otherwise build a path to their trust anchor need them.
     */
    public byte[] getEncodedCertificateChain() {
        return encodedCertificateChain;
    }

    private static byte[] encode(List<X509Certificate> chain) throws GeneralSecurityException {
        ByteArrayOutputStream encoded = new ByteArrayOutputStream();
        for (X509Certificate certificate : chain) {
            encoded.writeBytes(certificate.getEncoded());
        }
        return encoded.toByteArray();
    }

    public KeyPair getKeyPair() {
        return keyPair;
    }

    public X509Certificate getCertificate() {
        return certificate;
    }

    public PrivateKey getPrivateKey() {
        return keyPair.getPrivate();
    }

    public byte[] getThumbPrint() {
        return thumbprint;
    }

    public Optional<String> getApplicationUri() {
        Collection<List<?>> lists;
        try {
            lists = certificate.getSubjectAlternativeNames();
        } catch (Exception e) {
            return Optional.empty();
        }
        if (lists == null) {
            return Optional.empty();
        }
        return lists.stream()
            .filter(l -> l.size() == 2)
            .filter(name -> name.get(0).equals(GeneralName.uniformResourceIdentifier))
            .map(t -> (String) t.get(1))
            .findAny();
    }
}
