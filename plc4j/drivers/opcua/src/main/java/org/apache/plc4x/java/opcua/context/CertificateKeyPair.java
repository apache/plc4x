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

import io.vavr.control.Try;
import org.bouncycastle.asn1.x509.GeneralName;

import java.security.KeyPair;
import java.security.MessageDigest;
import java.security.cert.X509Certificate;
import java.util.Collection;
import java.util.List;
import java.util.Optional;

public class CertificateKeyPair {

    private final KeyPair keyPair;
    private final X509Certificate certificate;
    private final byte[] thumbprint;

    public CertificateKeyPair(KeyPair keyPair, X509Certificate certificate) throws Exception {
        this.keyPair = keyPair;
        this.certificate = certificate;
        MessageDigest messageDigest = MessageDigest.getInstance("SHA-1");
        this.thumbprint = messageDigest.digest(this.certificate.getEncoded());
    }

    public KeyPair getKeyPair() {
        return keyPair;
    }

    public X509Certificate getCertificate() {
        return certificate;
    }

    public byte[] getThumbPrint() {
        return thumbprint;
    }

    public Optional<String> getApplicationUri() {
        Try<Collection<List<?>>> lists = Try.of(certificate::getSubjectAlternativeNames);
        return lists.toJavaStream()
            .flatMap(Collection::stream)
            .filter(l -> l.size() == 2)
            .filter(name -> name.get(0).equals(GeneralName.uniformResourceIdentifier))
            .map(t -> (String) t.get(1))
            .findAny();
    }
}
