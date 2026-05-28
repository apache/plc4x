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
package org.apache.plc4x.java.api.authentication;

import java.security.KeyStore;
import java.util.Arrays;
import java.util.Objects;

/**
 * Authentication implementation for mutual TLS (mTLS) using client certificates.
 * <p>
 * This class provides client certificate authentication for TLS connections that require
 * mutual authentication. The client certificate and private key are provided via a KeyStore.
 * <p>
 * Example usage:
 * <pre>{@code
 * KeyStore keyStore = KeyStore.getInstance("PKCS12");
 * try (FileInputStream fis = new FileInputStream("client-cert.p12")) {
 *     keyStore.load(fis, "keystorePassword".toCharArray());
 * }
 *
 * PlcAuthentication auth = new PlcCertificateAuthentication(
 *     keyStore,
 *     "keystorePassword".toCharArray()
 * );
 *
 * PlcConnection connection = driverManager.getConnection(connectionUrl, auth);
 * }</pre>
 */
public class PlcCertificateAuthentication implements PlcAuthentication {

    private final KeyStore keyStore;
    private final char[] keyStorePassword;
    private final String keyAlias;

    /**
     * Creates a new certificate authentication with the specified KeyStore.
     * Uses the first key entry found in the KeyStore.
     *
     * @param keyStore         the KeyStore containing the client certificate and private key
     * @param keyStorePassword the password to access the KeyStore and key entries
     * @throws NullPointerException if keyStore or keyStorePassword is null
     */
    public PlcCertificateAuthentication(KeyStore keyStore, char[] keyStorePassword) {
        this(keyStore, keyStorePassword, null);
    }

    /**
     * Creates a new certificate authentication with the specified KeyStore and key alias.
     *
     * @param keyStore         the KeyStore containing the client certificate and private key
     * @param keyStorePassword the password to access the KeyStore and key entries
     * @param keyAlias         the alias of the key entry to use, or null to use the first key entry
     * @throws NullPointerException if keyStore or keyStorePassword is null
     */
    public PlcCertificateAuthentication(KeyStore keyStore, char[] keyStorePassword, String keyAlias) {
        Objects.requireNonNull(keyStore, "KeyStore must not be null");
        Objects.requireNonNull(keyStorePassword, "KeyStore password must not be null");
        this.keyStore = keyStore;
        this.keyStorePassword = keyStorePassword.clone();
        this.keyAlias = keyAlias;
    }

    /**
     * Returns the KeyStore containing the client certificate and private key.
     *
     * @return the KeyStore
     */
    public KeyStore getKeyStore() {
        return keyStore;
    }

    /**
     * Returns the password for the KeyStore and key entries.
     * Note: Returns a clone for security reasons.
     *
     * @return a copy of the KeyStore password
     */
    public char[] getKeyStorePassword() {
        return keyStorePassword.clone();
    }

    /**
     * Returns the alias of the key entry to use.
     *
     * @return the key alias, or null if the first key entry should be used
     */
    public String getKeyAlias() {
        return keyAlias;
    }

    @Override
    public final boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcCertificateAuthentication that)) {
            return false;
        }
        return Objects.equals(keyStore, that.keyStore) &&
            Arrays.equals(keyStorePassword, that.keyStorePassword) &&
            Objects.equals(keyAlias, that.keyAlias);
    }

    @Override
    public final int hashCode() {
        int result = Objects.hash(keyStore, keyAlias);
        result = 31 * result + Arrays.hashCode(keyStorePassword);
        return result;
    }

    @Override
    public String toString() {
        return "PlcCertificateAuthentication{" +
            "keyStore=" + keyStore.getType() +
            ", keyAlias='" + (keyAlias != null ? keyAlias : "(first entry)") + '\'' +
            '}';
    }

}
