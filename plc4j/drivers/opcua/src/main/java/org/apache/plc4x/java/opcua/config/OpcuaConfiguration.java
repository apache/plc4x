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
package org.apache.plc4x.java.opcua.config;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.opcua.context.CertificateGenerator;
import org.apache.plc4x.java.opcua.context.CertificateKeyPair;
import org.apache.plc4x.java.opcua.protocol.OpcuaProtocolLogic;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.tcp.TcpTransportConfiguration;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.FileSystems;
import java.nio.file.Path;
import java.security.*;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

public class OpcuaConfiguration implements Configuration, TcpTransportConfiguration {

    static {
        // Required for SecurityPolicy.Aes256_Sha256_RsaPss
        Security.addProvider(new BouncyCastleProvider());
    }

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaConfiguration.class);

    private String code;
    private String host;
    private String port;
    private String endpoint;
    private String transportEndpoint;
    private String params;
    private Boolean isEncrypted = false;
    private PascalByteString thumbprint;
    private byte[] senderCertificate;

    @ConfigurationParameter("discovery")
    @BooleanDefaultValue(true)
    private boolean discovery;

    @ConfigurationParameter("username")
    private String username;

    @ConfigurationParameter("password")
    private String password;

    @ConfigurationParameter("securityPolicy")
    @StringDefaultValue("None")
    private String securityPolicy;

    @ConfigurationParameter("keyStoreFile")
    private String keyStoreFile;

    @ConfigurationParameter("certDirectory")
    private String certDirectory;

    @ConfigurationParameter("keyStorePassword")
    private String keyStorePassword;

    private CertificateKeyPair ckp;

    public boolean isDiscovery() {
        return discovery;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public String getCertDirectory() {
        return certDirectory;
    }

    public String getSecurityPolicy() {
        return securityPolicy;
    }

    public String getKeyStoreFile() {
        return keyStoreFile;
    }

    public String getKeyStorePassword() {
        return keyStorePassword;
    }

    public PascalByteString getThumbprint() {
        return thumbprint;
    }

    public CertificateKeyPair getCertificateKeyPair() {
        return ckp;
    }

    public boolean isEncrypted() { return isEncrypted; }

    public void setDiscovery(boolean discovery) {
        this.discovery = discovery;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public void setPassword(String password) {
        this.password = password;
    }

    public void setCertDirectory(String certDirectory) {
        this.certDirectory = certDirectory;
    }

    public void setSecurityPolicy(String securityPolicy) {
        this.securityPolicy = securityPolicy;
    }

    public void setKeyStoreFile(String keyStoreFile) {
        this.keyStoreFile = keyStoreFile;
    }

    public void setKeyStorePassword(String keyStorePassword) {
        this.keyStorePassword = keyStorePassword;
    }

    public void setThumbprint(PascalByteString thumbprint) { this.thumbprint = thumbprint; }

    public String getTransportCode() {
        return code;
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

    public byte[] getSenderCertificate() {
        return this.senderCertificate;
    }

    public void setTransportCode(String code) {
        this.code = code;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public void setPort(String port) {
        this.port = port;
    }

    public void setEndpoint(String endpoint) {
        this.endpoint = endpoint;
    }

    public void setTransportEndpoint(String transportEndpoint) { this.transportEndpoint = transportEndpoint; }

    public void openKeyStore() throws Exception {
        this.isEncrypted = true;
        File securityTempDir = new File(certDirectory, "security");
        if (!securityTempDir.exists() && !securityTempDir.mkdirs()) {
            throw new PlcConnectionException("Unable to create directory please confirm folder permissions on "  + certDirectory);
        }
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        File serverKeyStore = securityTempDir.toPath().resolve(keyStoreFile).toFile();

        File pkiDir = FileSystems.getDefault().getPath(certDirectory).resolve("pki").toFile();
        if (!serverKeyStore.exists()) {
            ckp = CertificateGenerator.generateCertificate();
            LOGGER.info("Creating new KeyStore at {}", serverKeyStore);
            keyStore.load(null, keyStorePassword.toCharArray());
            keyStore.setKeyEntry("plc4x-certificate-alias", ckp.getKeyPair().getPrivate(), keyStorePassword.toCharArray(), new X509Certificate[] { ckp.getCertificate() });
            keyStore.store(new FileOutputStream(serverKeyStore), keyStorePassword.toCharArray());
        } else {
            LOGGER.info("Loading KeyStore at {}", serverKeyStore);
            keyStore.load(new FileInputStream(serverKeyStore), keyStorePassword.toCharArray());
            String alias = keyStore.aliases().nextElement();
            KeyPair kp = new KeyPair(keyStore.getCertificate(alias).getPublicKey(),
                (PrivateKey) keyStore.getKey(alias, keyStorePassword.toCharArray()));
            ckp = new CertificateKeyPair(kp,(X509Certificate) keyStore.getCertificate(alias));
        }
    }

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

    public void setSenderCertificate(byte[] certificate) { this.senderCertificate = certificate; }

}

