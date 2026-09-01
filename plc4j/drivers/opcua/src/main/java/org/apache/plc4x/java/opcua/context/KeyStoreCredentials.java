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

import java.security.GeneralSecurityException;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.PrivateKey;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Enumeration;
import java.util.List;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Reads an identity - a certificate, its chain and the matching private key - out of a key store.
 * <p>
 * Two identities are involved in an OPC UA connection and both come from a key store: the
 * application instance certificate that secures the channel, and the user certificate of an
 * {@code X509IdentityToken}. They are loaded the same way, so the alias selection and the
 * validation live here rather than being duplicated per caller.
 */
final class KeyStoreCredentials {

    private static final Logger LOGGER = LoggerFactory.getLogger(KeyStoreCredentials.class);

    private KeyStoreCredentials() {
    }

    /**
     * @param keyStore the key store to read from
     * @param password the password protecting the key entries
     * @param alias    the entry to use, or {@code null} to take the first entry that holds a
     *                 private key
     * @param source   where the key store came from, used in error messages so the user knows
     *                 which of the two identities is at fault
     * @param purpose  what the identity is used for, named when rejecting a non-RSA key
     */
    static CertificateKeyPair load(KeyStore keyStore, char[] password, String alias, String source, String purpose)
        throws GeneralSecurityException {
        String selected = alias == null ? selectKeyAlias(keyStore, source) : requireAlias(keyStore, alias, source);
        X509Certificate certificate = (X509Certificate) keyStore.getCertificate(selected);
        PrivateKey privateKey = (PrivateKey) keyStore.getKey(selected, password);
        validate(selected, certificate, privateKey, purpose);
        LOGGER.info("Using entry '{}' of the key store: {} bit {} certificate for '{}'", selected,
            ((RSAPublicKey) certificate.getPublicKey()).getModulus().bitLength(),
            certificate.getPublicKey().getAlgorithm(), certificate.getSubjectX500Principal());
        return new CertificateKeyPair(new KeyPair(certificate.getPublicKey(), privateKey),
            certificateChain(keyStore, selected, certificate));
    }

    /**
     * Picks the entry to authenticate with: the first one that actually holds a private key.
     * Taking whatever alias came first meant a key store that also carries the signing CA (as
     * produced by the certificate tutorial) could hand the CA's entry to the server, or an entry
     * with no private key at all.
     */
    static String selectKeyAlias(KeyStore keyStore, String source) throws GeneralSecurityException {
        Enumeration<String> aliases = keyStore.aliases();
        while (aliases.hasMoreElements()) {
            String alias = aliases.nextElement();
            if (keyStore.isKeyEntry(alias)) {
                return alias;
            }
            LOGGER.debug("Skipping key store entry '{}': it holds no private key", alias);
        }
        throw new KeyStoreException("The key store '" + source
            + "' contains no entry with a private key, so it cannot be used as a client identity.");
    }

    /**
     * An explicitly requested alias has to exist and hold a private key. Silently falling back to
     * another entry would authenticate as somebody else than asked for.
     */
    private static String requireAlias(KeyStore keyStore, String alias, String source) throws GeneralSecurityException {
        if (!keyStore.containsAlias(alias)) {
            throw new KeyStoreException("The key store '" + source + "' has no entry '" + alias
                + "'. It contains: " + String.join(", ", aliases(keyStore)) + ".");
        }
        if (!keyStore.isKeyEntry(alias)) {
            throw new KeyStoreException("Entry '" + alias + "' of the key store '" + source
                + "' holds no private key, so it cannot be used as an identity.");
        }
        return alias;
    }

    private static List<String> aliases(KeyStore keyStore) throws GeneralSecurityException {
        return Collections.list(keyStore.aliases());
    }

    /**
     * The certificate chain stored under the given alias, the identity's own certificate first. A
     * CA-signed certificate is only verifiable by a peer that can build a path to its trust
     * anchor, which usually means the issuing certificates have to travel with it - see OPC UA
     * Part 6, SenderCertificate. A self-signed certificate simply yields a chain of one.
     */
    static List<X509Certificate> certificateChain(KeyStore keyStore, String alias, X509Certificate certificate)
        throws GeneralSecurityException {
        Certificate[] chain = keyStore.getCertificateChain(alias);
        if (chain == null || chain.length == 0) {
            return List.of(certificate);
        }
        List<X509Certificate> certificates = new ArrayList<>(chain.length);
        for (Certificate element : chain) {
            if (element instanceof X509Certificate x509) {
                certificates.add(x509);
            } else {
                LOGGER.warn("Ignoring a {} certificate in the chain of key store entry '{}'",
                    element.getType(), alias);
            }
        }
        if (certificates.size() > 1) {
            LOGGER.info("Key store entry '{}' carries a chain of {} certificates", alias, certificates.size());
        }
        return certificates.isEmpty() ? List.of(certificate) : certificates;
    }

    /**
     * OPC UA identities are RSA; anything else is rejected by the server during the handshake,
     * typically with an error that says nothing about the cause. Failing here names the offending
     * key store entry instead.
     */
    static void validate(String alias, X509Certificate certificate, PrivateKey privateKey, String purpose)
        throws GeneralSecurityException {
        if (certificate == null) {
            throw new KeyStoreException("Key store entry '" + alias + "' has no certificate.");
        }
        if (privateKey == null) {
            throw new KeyStoreException("Key store entry '" + alias + "' has no private key."
                + " Check the key store password - a wrong one yields an entry without a key.");
        }
        if (!(certificate.getPublicKey() instanceof RSAPublicKey)) {
            throw new KeyStoreException("Key store entry '" + alias + "' uses a "
                + certificate.getPublicKey().getAlgorithm() + " key. OPC UA requires RSA keys for"
                + " " + purpose + "; regenerate the certificate with an RSA key.");
        }
    }
}
