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

import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

/**
 * Fail-closed certificate verifier used as the secure default: it rejects every
 * server certificate because no trust anchor has been configured. To establish
 * trust, configure either a {@code trust-store-file} (chain validation) or a
 * {@code server-certificate-file} (certificate pinning). As a last resort,
 * certificate verification can be disabled entirely with
 * {@code insecure-certificate-verification=true}, which is unsafe and leaves the
 * connection open to man-in-the-middle attacks.
 */
public class RejectingCertificateVerifier implements CertificateVerifier {

    @Override
    public void checkCertificateTrusted(X509Certificate certificate) throws CertificateException {
        throw new CertificateException("No trust anchor configured for OPC UA server certificate verification. "
            + "Configure 'trust-store-file' or 'server-certificate-file' to establish trust, or set "
            + "'insecure-certificate-verification=true' to disable verification (unsafe).");
    }

}
