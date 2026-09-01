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
import org.apache.plc4x.java.spi.config.annotations.Secret;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.transport.tcp.config.TcpTransportConfiguration;

/**
 * Configuration for TLS-PSK (Pre-Shared Key) transport connections.
 * Extends TcpTransportConfiguration with PSK-specific options.
 */
public class PskTlsTransportConfiguration extends TcpTransportConfiguration {

    /**
     * PSK identity string for TLS-PSK authentication.
     * Must be used together with psk-key.
     */
    // Deliberately not @Secret: the identity says which key the device refused, which is the one
    // thing an operator needs when a PSK handshake fails. Hiding it costs the diagnosis and
    // protects nothing - the key it names is marked instead.
    @ConfigurationParameter("psk-identity")
    @Description("PSK identity string for TLS-PSK authentication. Must be used together with psk-key.")
    public String pskIdentity;

    /**
     * PSK key as hexadecimal string for TLS-PSK authentication.
     * Must be used together with psk-identity.
     */
    @Secret
    @ConfigurationParameter("psk-key")
    @Description("PSK key as hexadecimal string for TLS-PSK authentication. Must be used together with psk-identity.")
    public String pskKey;

    /**
     * Whether to log TLS session keys to the audit log in SSLKEYLOGFILE format.
     * This allows Wireshark to decrypt captured TLS traffic for debugging.
     * WARNING: This exposes session keys — only enable for debugging.
     */
    @ConfigurationParameter("log-session-keys")
    @BooleanDefaultValue(false)
    @Description("Log TLS session keys to the audit log in SSLKEYLOGFILE format for Wireshark decryption.")
    public boolean logSessionKeys;

    /**
     * Validates the PSK configuration. Throws IllegalArgumentException if:
     * - pskIdentity is null or empty
     * - pskKey is null, empty, has odd length, or contains non-hex characters
     */
    public void validatePskConfiguration() {
        if (pskIdentity == null || pskIdentity.isEmpty()) {
            throw new IllegalArgumentException("psk-identity must not be empty.");
        }

        if (pskKey == null || pskKey.isEmpty()) {
            throw new IllegalArgumentException("psk-key must not be empty.");
        }

        if (pskKey.length() % 2 != 0) {
            throw new IllegalArgumentException(
                "psk-key must have even length (each byte is 2 hex characters). Length: " + pskKey.length());
        }

        if (!pskKey.matches("[0-9a-fA-F]+")) {
            throw new IllegalArgumentException(
                "psk-key must contain only hexadecimal characters (0-9, a-f, A-F).");
        }
    }

    /**
     * Converts the hex-encoded PSK key to a byte array.
     *
     * @return the PSK key as a byte array
     * @throws IllegalStateException if pskKey is null
     */
    public byte[] getPskKeyBytes() {
        if (pskKey == null) {
            throw new IllegalStateException("psk-key is not set.");
        }
        byte[] bytes = new byte[pskKey.length() / 2];
        for (int i = 0; i < bytes.length; i++) {
            bytes[i] = (byte) Integer.parseInt(pskKey.substring(i * 2, i * 2 + 2), 16);
        }
        return bytes;
    }

}
