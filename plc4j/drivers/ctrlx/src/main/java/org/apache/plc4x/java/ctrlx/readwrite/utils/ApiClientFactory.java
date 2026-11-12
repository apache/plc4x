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

package org.apache.plc4x.java.ctrlx.readwrite.utils;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration;
import org.apache.plc4x.java.ctrlx.readwrite.rest.CtrlXRestClient;
import org.apache.plc4x.java.ctrlx.readwrite.rest.CtrlXRestException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.TrustManagerFactory;
import javax.net.ssl.X509TrustManager;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.http.HttpClient;
import java.security.GeneralSecurityException;
import java.security.KeyStore;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;

public class ApiClientFactory {

    private static final Logger LOGGER = LoggerFactory.getLogger(ApiClientFactory.class);

    private ApiClientFactory() {
    }

    public static CtrlXRestClient getApiClient(String baseUrl, String username, String password,
                                               CtrlXConfiguration configuration) throws PlcConnectionException {
        SSLContext sslContext;
        try {
            sslContext = SSLContext.getInstance("TLS");
            if (configuration.isAllowFactoryDefaultCertificate()) {
                // The certificate this trusts ships in the jar, so it identifies nobody: anything
                // holding it and its key is trusted, for any host. Kept for a device still on its
                // factory settings, and loud about it.
                LOGGER.warn("allow-factory-default-certificate is set: trusting the certificate "
                    + "shipped with this driver, for any host. Anything presenting it is trusted, "
                    + "and the credentials for this connection travel over that channel");
                Certificate certificate = CertificateFactory.getInstance("X.509").generateCertificate(
                    ApiClientFactory.class.getClassLoader().getResourceAsStream("certs/webserver_cert.pem"));
                sslContext.init(null,
                    withoutHostnameCheck(trustManagersFor(singleCertificateStore(certificate))), null);
            } else {
                TrustManager[] trustManagers = trustManagersFrom(configuration);
                if (configuration.isIgnoreCommonName()) {
                    LOGGER.warn("ignore-common-name is set: a certificate issued for any other "
                        + "host will be accepted for this one");
                    trustManagers = withoutHostnameCheck(trustManagers);
                }
                sslContext.init(null, trustManagers, null);
            }
        } catch (GeneralSecurityException | IOException e) {
            throw new PlcConnectionException("Error setting up the TLS trust for this connection", e);
        }

        HttpClient httpClient = HttpClient.newBuilder()
            .sslContext(sslContext)
            // What the Apache client this replaced spoke; the device's REST server is only known
            // to work with that.
            .version(HttpClient.Version.HTTP_1_1)
            .build();
        try {
            // Log in and get a token, that we'll use in further interactions.
            return CtrlXRestClient.login(httpClient, baseUrl, username, password);
        } catch (CtrlXRestException e) {
            httpClient.close();
            throw new PlcConnectionException("Error getting access token", e);
        }
    }

    /**
     * The certificates to trust: what the operator named, or - naming nothing - the authorities the
     * JVM ships with, which is what any other HTTPS client would do.
     */
    private static TrustManager[] trustManagersFrom(CtrlXConfiguration configuration)
            throws GeneralSecurityException, IOException {
        if (configuration.getServerCertificateFile() != null
            && !configuration.getServerCertificateFile().isEmpty()) {
            try (InputStream is = new FileInputStream(configuration.getServerCertificateFile())) {
                Certificate certificate = CertificateFactory.getInstance("X.509").generateCertificate(is);
                return trustManagersFor(singleCertificateStore(certificate));
            }
        }
        if (configuration.getTrustStoreFile() != null && !configuration.getTrustStoreFile().isEmpty()) {
            String type = configuration.getTrustStoreType() != null
                ? configuration.getTrustStoreType() : "PKCS12";
            char[] storePassword = configuration.getTrustStorePassword() != null
                ? configuration.getTrustStorePassword().toCharArray() : null;
            KeyStore trustStore = KeyStore.getInstance(type);
            try (InputStream is = new FileInputStream(configuration.getTrustStoreFile())) {
                trustStore.load(is, storePassword);
            }
            return trustManagersFor(trustStore);
        }
        // A null key store gives the platform's trust managers.
        return trustManagersFor(null);
    }

    /**
     * The same trust, minus the check that the certificate was issued for the host connected to.
     */
    static TrustManager[] withoutHostnameCheck(TrustManager[] trustManagers) {
        TrustManager[] wrapped = new TrustManager[trustManagers.length];
        for (int i = 0; i < trustManagers.length; i++) {
            wrapped[i] = trustManagers[i] instanceof X509TrustManager x509
                ? new NoHostnameCheckTrustManager(x509)
                : trustManagers[i];
        }
        return wrapped;
    }

    private static KeyStore singleCertificateStore(Certificate certificate)
            throws GeneralSecurityException, IOException {
        KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
        keyStore.load(null, null);
        keyStore.setCertificateEntry("server", certificate);
        return keyStore;
    }

    private static TrustManager[] trustManagersFor(KeyStore keyStore) throws GeneralSecurityException {
        TrustManagerFactory trustManagerFactory =
            TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
        trustManagerFactory.init(keyStore);
        return trustManagerFactory.getTrustManagers();
    }
}
