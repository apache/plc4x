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

import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThatCode;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class PinnedCertificateVerifierTest {

    @Test
    void acceptsTheExactPinnedCertificate() throws Exception {
        X509Certificate pinned = TestCertificateGenerator.generate(2048, "CN=server", 3600).getValue();
        PinnedCertificateVerifier verifier = new PinnedCertificateVerifier(pinned);
        assertThatCode(() -> verifier.checkCertificateTrusted(pinned)).doesNotThrowAnyException();
    }

    @Test
    void rejectsADifferentCertificate() throws Exception {
        X509Certificate pinned = TestCertificateGenerator.generate(2048, "CN=server", 3600).getValue();
        X509Certificate impostor = TestCertificateGenerator.generate(2048, "CN=server", 3600).getValue();
        PinnedCertificateVerifier verifier = new PinnedCertificateVerifier(pinned);
        assertThatThrownBy(() -> verifier.checkCertificateTrusted(impostor))
            .isInstanceOf(CertificateException.class);
    }

    @Test
    void rejectsNullCertificate() throws Exception {
        X509Certificate pinned = TestCertificateGenerator.generate(2048, "CN=server", 3600).getValue();
        PinnedCertificateVerifier verifier = new PinnedCertificateVerifier(pinned);
        assertThatThrownBy(() -> verifier.checkCertificateTrusted(null))
            .isInstanceOf(CertificateException.class);
    }
}
