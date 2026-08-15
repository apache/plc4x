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

import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Paths;
import java.security.GeneralSecurityException;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.List;
import java.util.Optional;
import org.apache.commons.codec.digest.DigestUtils;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.security.CertificateVerifier;
import org.apache.plc4x.java.opcua.security.PermissiveCertificateVerifier;
import org.apache.plc4x.java.opcua.security.PinnedCertificateVerifier;
import org.apache.plc4x.java.opcua.security.RejectingCertificateVerifier;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.security.TrustStoreCertificateVerifier;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.PrivateKey;
import java.security.Security;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Per-connection OPC UA state. Used to be wired up as the old SPI's
 * {@code DriverContext} + {@code HasConfiguration<OpcuaConfiguration>}; with
 * the new SPI the connection just owns one of these directly.
 */
public class OpcuaDriverContext {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaDriverContext.class);

    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("(:(?<transportCode>[a-z0-9]*))?://" +
        "(?<transportHost>[\\w.-]+)(:" +
        "(?<transportPort>\\d*))?");

    public static final Pattern URI_PATTERN = Pattern.compile("^(?<protocolCode>opcua)" +
        INET_ADDRESS_PATTERN +
        "(?<transportEndpoint>[\\w/=]*)[?]?" +
        "(?<paramString>([^=]+=[^=&]+&?)*)"
    );


    static {
        // Required for SecurityPolicy.Aes256_Sha256_RsaPss
        Security.addProvider(new BouncyCastleProvider());
    }

    private String code;
    private String host;
    private String port;
    private String endpoint;
    private String transportEndpoint;
    private CertificateKeyPair certificateKeyPair;
    private X509Certificate serverCertificate;
    private PascalByteString thumbprint;

    // Secure by default: reject every server certificate until an explicit trust
    // anchor (trust store or pinned server certificate) is configured, or the user
    // explicitly opts out of verification. Anything more permissive by default would
    // leave signed secure channels open to man-in-the-middle attacks.
    private CertificateVerifier certificateVerifier = new RejectingCertificateVerifier();


    public void openKeyStore(OpcuaConfiguration configuration) throws IOException, GeneralSecurityException {
        String serverKeyStore = configuration.getKeyStoreFile();

        if (serverKeyStore == null) {
            int keySize = configuration.getGeneratedKeySize();
            if (keySize < 2048 || keySize % 1024 != 0) {
                throw new IllegalArgumentException("'generated-key-size' has to be a multiple of 1024"
                    + " and at least 2048, was " + keySize);
            }
            LOGGER.info("Client certificate not provided, creating temporary {} bit certificate and private key", keySize);
            certificateKeyPair = CertificateGenerator.generateCertificate(keySize);
        } else {
            LOGGER.info("Loading KeyStore at {}", serverKeyStore);

            KeyStore keyStore = openKeyStore(configuration.getKeyStoreFile(), configuration.getKeyStoreType(), configuration.getKeyStorePassword());
            String alias = selectKeyAlias(keyStore, configuration.getKeyStoreFile());
            X509Certificate certificate = (X509Certificate) keyStore.getCertificate(alias);
            PrivateKey privateKey = (PrivateKey) keyStore.getKey(alias, configuration.getKeyStorePassword());
            validateClientKey(alias, certificate, privateKey);
            LOGGER.info("Using entry '{}' of the key store: {} bit {} certificate for '{}'", alias,
                ((RSAPublicKey) certificate.getPublicKey()).getModulus().bitLength(),
                certificate.getPublicKey().getAlgorithm(), certificate.getSubjectX500Principal());
            KeyPair kp = new KeyPair(certificate.getPublicKey(), privateKey);
            certificateKeyPair = new CertificateKeyPair(kp, certificateChain(keyStore, alias, certificate));
        }

        if (configuration.getServerCertificate() != null) {
            serverCertificate = configuration.getServerCertificate();
            byte[] sha1 = DigestUtils.sha1(serverCertificate.getEncoded());
            thumbprint = new PascalByteString(sha1.length, sha1);
        }

        certificateVerifier = buildCertificateVerifier(configuration);
    }

    /**
     * Picks the key store entry to authenticate with: the first one that actually holds a private
     * key. Taking whatever alias came first meant a key store that also carries the signing CA (as
     * produced by the certificate tutorial) could hand the CA's entry to the server, or an entry
     * with no private key at all.
     */
    private String selectKeyAlias(KeyStore keyStore, String keyStoreFile) throws GeneralSecurityException {
        Enumeration<String> aliases = keyStore.aliases();
        while (aliases.hasMoreElements()) {
            String alias = aliases.nextElement();
            if (keyStore.isKeyEntry(alias)) {
                return alias;
            }
            LOGGER.debug("Skipping key store entry '{}': it holds no private key", alias);
        }
        throw new KeyStoreException("The key store '" + keyStoreFile
            + "' contains no entry with a private key, so it cannot be used as a client identity.");
    }

    /**
     * The certificate chain stored under the given alias, the client certificate first. A
     * CA-signed client certificate is only verifiable by a server that can build a path to its
     * trust anchor, which usually means the issuing certificates have to travel with it - see
     * OPC UA Part 6, SenderCertificate. A self-signed certificate simply yields a chain of one.
     */
    private List<X509Certificate> certificateChain(KeyStore keyStore, String alias, X509Certificate certificate)
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
            LOGGER.info("Key store entry '{}' carries a chain of {} certificates; sending all of them",
                alias, certificates.size());
        }
        return certificates.isEmpty() ? List.of(certificate) : certificates;
    }

    /**
     * OPC UA application instance certificates are RSA; anything else is rejected by the server
     * during the handshake, typically with an error that says nothing about the cause. Failing here
     * names the offending key store entry instead.
     */
    private void validateClientKey(String alias, X509Certificate certificate, PrivateKey privateKey)
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
                + " application instance certificates; regenerate the certificate with an RSA key.");
        }
    }

    /**
     * Selects the server-certificate trust strategy, in order of precedence:
     * <ol>
     *   <li>{@code insecure-certificate-verification=true} &rarr; trust everything (unsafe, opt-in only);</li>
     *   <li>a {@code trust-store-file} &rarr; validate the certificate chain against the trust store;</li>
     *   <li>a {@code server-certificate-file} &rarr; pin trust to that exact certificate;</li>
     *   <li>otherwise &rarr; fail closed and reject, since no trust anchor is available.</li>
     * </ol>
     * Note that the pinned certificate is read from the configured file only — a
     * certificate learned over the unauthenticated discovery channel is never used
     * as a trust anchor.
     */
    private CertificateVerifier buildCertificateVerifier(OpcuaConfiguration configuration)
        throws IOException, GeneralSecurityException {
        if (configuration.isInsecureCertificateVerification()) {
            LOGGER.warn("OPC UA server certificate verification is DISABLED "
                + "('insecure-certificate-verification=true'). The connection is vulnerable to "
                + "man-in-the-middle attacks; do not use this in production.");
            return new PermissiveCertificateVerifier();
        }
        if (configuration.getTrustStoreFile() != null) {
            KeyStore trustStore = openKeyStore(configuration.getTrustStoreFile(), configuration.getTrustStoreType(), configuration.getTrustStorePassword());
            return new TrustStoreCertificateVerifier(trustStore);
        }
        if (configuration.getServerCertificateFile() != null) {
            LOGGER.info("Pinning OPC UA server certificate trust to {}", configuration.getServerCertificateFile());
            return new PinnedCertificateVerifier(configuration.getServerCertificate());
        }
        LOGGER.warn("No OPC UA trust anchor configured ('trust-store-file' or 'server-certificate-file'); "
            + "server certificates will be rejected. Set 'insecure-certificate-verification=true' to bypass "
            + "verification for local testing only.");
        return new RejectingCertificateVerifier();
    }

    public String getHost() {
        return host;
    }

    public String getPort() {
        return port;
    }

    public String getEndpoint() {
        return endpoint;
    }

    public String getTransportEndpoint() {
        return transportEndpoint;
    }

    public X509Certificate getServerCertificate() {
        return serverCertificate;
    }

    public CertificateKeyPair getCertificateKeyPair() {
        return certificateKeyPair;
    }

    /**
     * Initialises the per-connection state from the URL components the new SPI
     * already resolved (the {@code protocol-code}/{@code transport-code}/
     * {@code transport-config} fields on {@link OpcuaConfiguration} are
     * never populated by the framework — they're URL components, not
     * configuration parameters).
     *
     * @param transportCode   resolved transport code ({@code "tcp"}, etc.)
     * @param host            remote host (from the transport's resolved socket address)
     * @param port            remote port as a string, or {@code null}
     * @param transportEndpoint path part of the URL after host:port, may be empty
     */
    public void initialize(String transportCode, String host, String port, String transportEndpoint,
                           OpcuaConfiguration configuration) {
        this.code = transportCode;
        this.host = host;
        this.port = port;
        this.transportEndpoint = transportEndpoint != null ? transportEndpoint : "";

        String portAddition = port != null ? ":" + port : "";
        this.endpoint = "opc." + transportCode + "://" + host + portAddition + this.transportEndpoint;

        if (configuration.getSecurityPolicy() != null && configuration.getSecurityPolicy() != SecurityPolicy.NONE) {
            try {
                openKeyStore(configuration);
            } catch (IOException | GeneralSecurityException e) {
                throw new PlcRuntimeException("Unable to open keystore, please confirm you have the correct permissions", e);
            }
        }
    }

    /**
     * @deprecated since 0.14 — superseded by
     * {@link #initialize(String, String, String, String, OpcuaConfiguration)}. The
     * URL-derived parameters this method reads from {@link OpcuaConfiguration}
     * are not populated by the SPI; kept here so the legacy URL-pattern unit
     * tests keep compiling.
     */
    @Deprecated
    public void setConfiguration(OpcuaConfiguration configuration) {
        Matcher matcher = getMatcher(configuration);
        code = matcher.group("transportCode");
        host = matcher.group("transportHost");
        port = matcher.group("transportPort");
        transportEndpoint = matcher.group("transportEndpoint");

        String portAddition = port != null ? ":" + port : "";
        endpoint = "opc." + code + "://" + host + portAddition + transportEndpoint;

        if (configuration.getSecurityPolicy() != null && configuration.getSecurityPolicy() != SecurityPolicy.NONE) {
            try {
                openKeyStore(configuration);
            } catch (IOException | GeneralSecurityException e) {
                throw new PlcRuntimeException("Unable to open keystore, please confirm you have the correct permissions", e);
            }
        }
    }

    private static Matcher getMatcher(OpcuaConfiguration configuration) {
        String uri = configuration.getProtocolCode() + ":" + configuration.getTransportCode() + "://" + configuration.getTransportConfig();

        // Split up the connection string into it's individual segments.
        Matcher matcher = URI_PATTERN.matcher(uri);
        if (!matcher.matches()) {
            throw new PlcRuntimeException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-host}(:{transport-port})(/{transport-endpoint})(?{parameter-string)?': " + uri);
        }
        return matcher;
    }

    public Optional<String> getApplicationUri() {
        return Optional.ofNullable(certificateKeyPair)
            .flatMap(CertificateKeyPair::getApplicationUri);
    }

    public PascalByteString getThumbprint() {
        return thumbprint;
    }

    public CertificateVerifier getCertificateVerifier() {
        return certificateVerifier;
    }

    private static KeyStore openKeyStore(String keyStoreFile, String keyStoreType, char[] password) throws IOException, GeneralSecurityException {
        File serverKeyStore = null;
        if (keyStoreFile != null) {
            serverKeyStore = Paths.get(keyStoreFile).toFile();
        }
        if (keyStoreFile == null || !serverKeyStore.exists()) {
            throw new FileNotFoundException("Invalid parameter - specified file " + keyStoreFile + " does not exist");
        }

        KeyStore keyStore = KeyStore.getInstance(keyStoreType);
        // The stream has to be closed explicitly - on Windows an open handle keeps the file locked.
        try (InputStream inputStream = new FileInputStream(serverKeyStore)) {
            keyStore.load(inputStream, password);
        }
        return keyStore;
    }

}
