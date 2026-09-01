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

import org.apache.http.config.Registry;
import org.apache.http.config.RegistryBuilder;
import org.apache.http.conn.socket.ConnectionSocketFactory;
import org.apache.http.conn.ssl.DefaultHostnameVerifier;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.impl.conn.BasicHttpClientConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration;
import org.apache.plc4x.java.ctrlx.readwrite.rest.datalayer.ApiClient;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.TrustManagerFactory;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.security.GeneralSecurityException;
import java.security.KeyStore;
import java.security.cert.Certificate;
import java.security.cert.CertificateFactory;

public class ApiClientFactory {

    private static final Logger LOGGER = LoggerFactory.getLogger(ApiClientFactory.class);

    public static ApiClient getApiClient(String baseUrl, String username, String password,
                                         CtrlXConfiguration configuration) throws PlcConnectionException {
        try {
            SSLContext sslContext = SSLContext.getInstance("TLS");
            HostnameVerifier hostnameVerifier;

            if (configuration.isAllowFactoryDefaultCertificate()) {
                // The certificate this trusts ships in the jar, so it identifies nobody: anything
                // holding it and its key is trusted, for any host. Kept for a device still on its
                // factory settings, and loud about it.
                LOGGER.warn("allow-factory-default-certificate is set: trusting the certificate "
                    + "shipped with this driver, for any host. Anything presenting it is trusted, "
                    + "and the credentials for this connection travel over that channel");
                Certificate certificate = CertificateFactory.getInstance("X.509").generateCertificate(
                    ApiClientFactory.class.getClassLoader().getResourceAsStream("certs/webserver_cert.pem"));
                sslContext.init(null, trustManagersFor(singleCertificateStore(certificate)), null);
                hostnameVerifier = NoopHostnameVerifier.INSTANCE;
            } else {
                sslContext.init(null, trustManagersFrom(configuration), null);
                hostnameVerifier = configuration.isIgnoreCommonName()
                    ? NoopHostnameVerifier.INSTANCE
                    : new DefaultHostnameVerifier();
                if (configuration.isIgnoreCommonName()) {
                    LOGGER.warn("ignore-common-name is set: a certificate issued for any other "
                        + "host will be accepted for this one");
                }
            }

            SSLConnectionSocketFactory sslConnectionSocketFactory =
                new SSLConnectionSocketFactory(sslContext, hostnameVerifier);
            Registry<ConnectionSocketFactory> socketFactoryRegistry =
                RegistryBuilder.<ConnectionSocketFactory>create()
                    .register("https", sslConnectionSocketFactory)
                    .build();

            // Create a new HTTP client using the ssl context.
            BasicHttpClientConnectionManager connectionManager = new BasicHttpClientConnectionManager(socketFactoryRegistry);
            CloseableHttpClient httpClient = HttpClients.custom().setSSLSocketFactory(sslConnectionSocketFactory).setConnectionManager(connectionManager).build();

            // Log in and get a token, that we'll use in further interactions.
            String bearerToken = OauthHelper.getAccessToken(httpClient, baseUrl, username, password);

            // Create a new API client.
            ApiClient apiClient = new ApiClient(httpClient);
            apiClient.setBasePath(baseUrl + "/automation/api/v2");
            apiClient.setAccessToken(bearerToken);

            return apiClient;
        } catch (GeneralSecurityException e) {
            throw new PlcConnectionException("Error setting up the TLS trust for this connection", e);
        } catch (IOException e) {
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
        // null lets the context use the platform trust managers.
        return null;
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
