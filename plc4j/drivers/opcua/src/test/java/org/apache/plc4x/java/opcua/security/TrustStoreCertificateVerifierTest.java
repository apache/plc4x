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
package org.apache.plc4x.java.opcua.security;

import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.junit.jupiter.api.Test;

import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map.Entry;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * The certificate being judged here belongs to the server we are talking to, and it usually arrives
 * with the issuers that vouch for it. Judging the leaf on its own turns away a chain that is
 * perfectly trustworthy, and asking the client-side question of a server certificate answers about
 * the wrong extended key usage.
 */
class TrustStoreCertificateVerifierTest {

    private static KeyStore storeContaining(X509Certificate... anchors) throws Exception {
        KeyStore trustStore = KeyStore.getInstance(KeyStore.getDefaultType());
        trustStore.load(null, null);
        for (int i = 0; i < anchors.length; i++) {
            trustStore.setCertificateEntry("anchor" + i, anchors[i]);
        }
        return trustStore;
    }

    @Test
    void aCertificateSignedByTheAnchorIsTrusted() throws Exception {
        Entry<PrivateKey, X509Certificate> root =
            TestCertificateGenerator.generateCa(2048, "CN=Test Root", 3600, null, null);
        Entry<PrivateKey, X509Certificate> server = TestCertificateGenerator.generateSignedBy(
            2048, "CN=Test Server", 3600, root.getKey(), root.getValue());

        TrustStoreCertificateVerifier verifier =
            new TrustStoreCertificateVerifier(storeContaining(root.getValue()));

        assertDoesNotThrow(() -> verifier.checkCertificateChainTrusted(
            Collections.singletonList(server.getValue())));
    }

    /**
     * The case the leaf-only check could not answer: the anchor signed an intermediate, and the
     * intermediate signed the server. Both arrive from the server, and both are needed.
     */
    @Test
    void aChainThroughAnIntermediateIsTrustedWhenTheIssuersAreGiven() throws Exception {
        Entry<PrivateKey, X509Certificate> root =
            TestCertificateGenerator.generateCa(2048, "CN=Test Root", 3600, null, null);
        Entry<PrivateKey, X509Certificate> intermediate = TestCertificateGenerator.generateCa(
            2048, "CN=Test Intermediate", 3600, root.getKey(), root.getValue());
        Entry<PrivateKey, X509Certificate> server = TestCertificateGenerator.generateSignedBy(
            2048, "CN=Test Server", 3600, intermediate.getKey(), intermediate.getValue());

        TrustStoreCertificateVerifier verifier =
            new TrustStoreCertificateVerifier(storeContaining(root.getValue()));

        List<X509Certificate> chain =
            Arrays.asList(server.getValue(), intermediate.getValue());
        assertDoesNotThrow(() -> verifier.checkCertificateChainTrusted(chain),
            "a three-level chain the server presented in full must validate");
    }

    @Test
    void aCertificateFromSomewhereElseIsNotTrusted() throws Exception {
        Entry<PrivateKey, X509Certificate> ourRoot =
            TestCertificateGenerator.generateCa(2048, "CN=Our Root", 3600, null, null);
        Entry<PrivateKey, X509Certificate> theirRoot =
            TestCertificateGenerator.generateCa(2048, "CN=Somebody Else", 3600, null, null);
        Entry<PrivateKey, X509Certificate> theirServer = TestCertificateGenerator.generateSignedBy(
            2048, "CN=Their Server", 3600, theirRoot.getKey(), theirRoot.getValue());

        TrustStoreCertificateVerifier verifier =
            new TrustStoreCertificateVerifier(storeContaining(ourRoot.getValue()));

        assertThrows(CertificateException.class, () -> verifier.checkCertificateChainTrusted(
            Collections.singletonList(theirServer.getValue())));
    }

    @Test
    void nothingToJudgeIsNotTheSameAsTrusted() throws Exception {
        Entry<PrivateKey, X509Certificate> root =
            TestCertificateGenerator.generateCa(2048, "CN=Test Root", 3600, null, null);
        TrustStoreCertificateVerifier verifier =
            new TrustStoreCertificateVerifier(storeContaining(root.getValue()));

        assertThrows(CertificateException.class,
            () -> verifier.checkCertificateChainTrusted(Collections.emptyList()));
        assertThrows(CertificateException.class,
            () -> verifier.checkCertificateChainTrusted(null));
    }
}
