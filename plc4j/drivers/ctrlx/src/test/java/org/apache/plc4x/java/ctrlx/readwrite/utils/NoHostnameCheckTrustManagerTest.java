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

package org.apache.plc4x.java.ctrlx.readwrite.utils;

import org.junit.jupiter.api.Test;

import javax.net.ssl.SSLEngine;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509ExtendedTrustManager;
import javax.net.ssl.X509TrustManager;
import java.net.Socket;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertThrows;

class NoHostnameCheckTrustManagerTest {

    /**
     * Records which variant was called. Only the plain (no socket, no engine) variants may be
     * reached through the wrapper: those are the ones without the host check.
     */
    private static final class RecordingTrustManager extends X509ExtendedTrustManager {
        final List<String> calls = new ArrayList<>();
        final X509Certificate[] issuers = new X509Certificate[0];
        boolean reject;

        private void record(String call) throws CertificateException {
            calls.add(call);
            if (reject) {
                throw new CertificateException("untrusted");
            }
        }

        @Override
        public void checkClientTrusted(X509Certificate[] chain, String authType) throws CertificateException {
            record("client");
        }

        @Override
        public void checkServerTrusted(X509Certificate[] chain, String authType) throws CertificateException {
            record("server");
        }

        @Override
        public void checkClientTrusted(X509Certificate[] chain, String authType, Socket socket) throws CertificateException {
            record("client+socket");
        }

        @Override
        public void checkServerTrusted(X509Certificate[] chain, String authType, Socket socket) throws CertificateException {
            record("server+socket");
        }

        @Override
        public void checkClientTrusted(X509Certificate[] chain, String authType, SSLEngine engine) throws CertificateException {
            record("client+engine");
        }

        @Override
        public void checkServerTrusted(X509Certificate[] chain, String authType, SSLEngine engine) throws CertificateException {
            record("server+engine");
        }

        @Override
        public X509Certificate[] getAcceptedIssuers() {
            return issuers;
        }
    }

    @Test
    void everyCheckGoesToTheVariantWithoutHostCheck() throws Exception {
        RecordingTrustManager delegate = new RecordingTrustManager();
        NoHostnameCheckTrustManager manager = new NoHostnameCheckTrustManager(delegate);
        X509Certificate[] chain = new X509Certificate[0];
        manager.checkServerTrusted(chain, "RSA");
        manager.checkServerTrusted(chain, "RSA", (Socket) null);
        manager.checkServerTrusted(chain, "RSA", (SSLEngine) null);
        manager.checkClientTrusted(chain, "RSA");
        manager.checkClientTrusted(chain, "RSA", (Socket) null);
        manager.checkClientTrusted(chain, "RSA", (SSLEngine) null);
        assertEquals(List.of("server", "server", "server", "client", "client", "client"), delegate.calls);
        assertSame(delegate.issuers, manager.getAcceptedIssuers());
    }

    @Test
    void theChainIsStillChecked() {
        RecordingTrustManager delegate = new RecordingTrustManager();
        delegate.reject = true;
        NoHostnameCheckTrustManager manager = new NoHostnameCheckTrustManager(delegate);
        assertThrows(CertificateException.class,
            () -> manager.checkServerTrusted(new X509Certificate[0], "RSA", (SSLEngine) null));
    }

    @Test
    void onlyX509TrustManagersAreWrapped() {
        TrustManager other = new TrustManager() {
        };
        X509TrustManager x509 = new RecordingTrustManager();
        TrustManager[] wrapped = ApiClientFactory.withoutHostnameCheck(new TrustManager[]{x509, other});
        assertInstanceOf(NoHostnameCheckTrustManager.class, wrapped[0]);
        assertSame(other, wrapped[1]);
    }

}
