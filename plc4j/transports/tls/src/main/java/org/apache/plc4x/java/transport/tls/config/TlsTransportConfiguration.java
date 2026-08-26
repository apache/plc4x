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
package org.apache.plc4x.java.transport.tls.config;

import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;

/**
 * Configuration for TLS transport connections.
 * Extends TcpTransportConfiguration with TLS-specific options.
 */
public class TlsTransportConfiguration extends TcpTransportConfiguration {

    /**
     * Whether to verify server certificates against trusted CAs.
     * Set to false for self-signed certificates in development/testing environments.
     * WARNING: Disabling verification still provides encryption but does NOT protect against MITM attacks.
     */
    @ConfigurationParameter("verify-ssl")
    @BooleanDefaultValue(true)
    public boolean verifySsl;

    /**
     * Whether to accept a certificate that is valid but was issued for some other name than the
     * host we connected to.
     *
     * <p>Chain validation says the certificate was issued by someone trusted. It says nothing about
     * who it was issued <em>to</em> - so without this check, a certificate trusted for any host is
     * accepted for every host, which is what a machine in the middle needs.</p>
     */
    @ConfigurationParameter("ignore-common-name")
    @BooleanDefaultValue(false)
    @Description("Accept a server certificate issued for a different host than the one connected to")
    public boolean ignoreCommonName;

    /**
     * A key store holding the certificates this connection should trust, instead of the public
     * certificate authorities the JVM ships with.
     *
     * <p>This is what makes verification usable for a device carrying its own certificate: without
     * it the only way past a private CA is to turn verification off entirely, which is why
     * {@code verify-ssl=false} tends to end up in configurations and stay there.</p>
     */
    @ConfigurationParameter("trust-store-file")
    @Description("Key store of certificates to trust, instead of the JVM's public authorities")
    public String trustStoreFile;

    @ConfigurationParameter("trust-store-password")
    @Description("Password of the trust store named by trust-store-file")
    public String trustStorePassword;

    @ConfigurationParameter("trust-store-type")
    @StringDefaultValue("PKCS12")
    @Description("Type of the trust store named by trust-store-file")
    public String trustStoreType;

    /**
     * TLS protocol version to use. If not set, defaults to TLS 1.3 with fallback to TLS 1.2.
     * Valid values: "TLSv1.2", "TLSv1.3"
     * Some protocols (like Secure ADS) require a specific TLS version.
     */
    @ConfigurationParameter("tls-version")
    @Description("TLS protocol version (e.g., 'TLSv1.2', 'TLSv1.3'). If not set, uses TLS 1.3 with fallback to TLS 1.2.")
    public String tlsVersion;

    /**
     * Returns whether SSL certificate verification is enabled.
     * When true, server certificates must be signed by a trusted CA.
     * When false, any server certificate is accepted (including self-signed).
     *
     * @return true if SSL verification is enabled, false otherwise
     */
    public boolean isVerifySsl() {
        return verifySsl;
    }

    /**
     * Returns the configured TLS version, or null if using default behavior.
     *
     * @return the TLS version string (e.g., "TLSv1.2") or null for default
     */
    public String getTlsVersion() {
        return tlsVersion;
    }

    /**
     * Path to a PKCS12 or JKS keystore containing the client certificate and private key.
     * Used for mutual TLS authentication where the server requires client certificates.
     */
    @ConfigurationParameter("keystore")
    @Description("Path to keystore (PKCS12/JKS) containing the client certificate and private key for mutual TLS.")
    public String keystore;

    /**
     * Password for the keystore specified by the keystore parameter.
     */
    @ConfigurationParameter("keystore-password")
    @Description("Password for the client keystore.")
    public String keystorePassword;

    /**
     * Type of the keystore. Defaults to PKCS12 if not specified.
     * Common values: "PKCS12", "JKS"
     */
    @ConfigurationParameter("keystore-type")
    @Description("Keystore type (e.g., 'PKCS12', 'JKS'). Defaults to PKCS12.")
    public String keystoreType;

    /**
     * Whether to log TLS session keys to the audit log in SSLKEYLOGFILE format.
     * This allows Wireshark to decrypt captured TLS traffic for debugging.
     * WARNING: This exposes session keys — only enable for debugging.
     */
    @ConfigurationParameter("log-session-keys")
    @BooleanDefaultValue(false)
    @Description("Log TLS session keys to the audit log in SSLKEYLOGFILE format for Wireshark decryption.")
    public boolean logSessionKeys;

}
