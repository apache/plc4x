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

package org.apache.plc4x.java.opcua.security;

import java.security.MessageDigest;
import java.security.cert.CertificateEncodingException;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.Objects;

/**
 * Certificate verifier that pins trust to a single, explicitly configured server
 * certificate (from {@code server-certificate-file}). The presented certificate
 * is accepted only when its DER encoding is byte-for-byte identical to the pinned
 * certificate. Comparison uses a constant-time check to avoid leaking information
 * through timing.
 */
public class PinnedCertificateVerifier implements CertificateVerifier {

    private final byte[] pinnedEncoded;

    public PinnedCertificateVerifier(X509Certificate pinnedCertificate) throws CertificateEncodingException {
        this.pinnedEncoded = Objects.requireNonNull(pinnedCertificate, "pinnedCertificate must not be null").getEncoded();
    }

    @Override
    public void checkCertificateTrusted(X509Certificate certificate) throws CertificateException {
        if (certificate == null || !MessageDigest.isEqual(pinnedEncoded, certificate.getEncoded())) {
            throw new CertificateException("Server certificate does not match the pinned 'server-certificate-file'.");
        }
    }

}
