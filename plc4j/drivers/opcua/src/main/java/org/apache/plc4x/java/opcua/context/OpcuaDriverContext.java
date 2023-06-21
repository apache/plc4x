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

import java.util.Optional;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.nio.file.FileSystems;
import java.security.KeyPair;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.Security;
import java.security.cert.X509Certificate;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class OpcuaDriverContext implements DriverContext, HasConfiguration<OpcuaConfiguration> {

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
    private Boolean isEncrypted = false;
    private PascalByteString thumbprint;
    private byte[] senderCertificate;
    private CertificateKeyPair certificateKeyPair;

    public void openKeyStore(OpcuaConfiguration configuration) throws Exception {
        this.isEncrypted = true;
        String certDirectory = configuration.getCertDirectory();
        File securityTempDir = new File(certDirectory, "security");
        if (!securityTempDir.exists() && !securityTempDir.mkdirs()) {
            throw new PlcRuntimeException("Unable to create directory please confirm folder permissions on " + certDirectory);
        }
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        File serverKeyStore = securityTempDir.toPath().resolve(configuration.getKeyStoreFile()).toFile();

        File pkiDir = FileSystems.getDefault().getPath(certDirectory).resolve("pki").toFile();
        char[] password = configuration.getKeyStorePassword().toCharArray();
        if (!serverKeyStore.exists()) {
            certificateKeyPair = CertificateGenerator.generateCertificate();
            LOGGER.info("Creating new KeyStore at {}", serverKeyStore);
            keyStore.load(null, password);
            keyStore.setKeyEntry("plc4x-certificate-alias", certificateKeyPair.getKeyPair().getPrivate(), password, new X509Certificate[]{certificateKeyPair.getCertificate()});
            keyStore.store(new FileOutputStream(serverKeyStore), password);
        } else {
            LOGGER.info("Loading KeyStore at {}", serverKeyStore);
            keyStore.load(new FileInputStream(serverKeyStore), password);
            String alias = keyStore.aliases().nextElement();
            KeyPair kp = new KeyPair(keyStore.getCertificate(alias).getPublicKey(),
                (PrivateKey) keyStore.getKey(alias, password));
            certificateKeyPair = new CertificateKeyPair(kp, (X509Certificate) keyStore.getCertificate(alias));
        }
    }

    public String getHost() {
        return host;
    }

    public void setHost(String host) {
        this.host = host;
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
    
    public void setTransportEndpoint(String transportEndpoint) {
        this.transportEndpoint = transportEndpoint;
    }

    public Boolean getEncrypted() {
        return isEncrypted;
    }

    public PascalByteString getThumbprint() {
        return thumbprint;
    }

    public void setThumbprint(PascalByteString thumbprint) {
        this.thumbprint = thumbprint;
    }

    public CertificateKeyPair getCertificateKeyPair() {
        return certificateKeyPair;
    }

    @Override
    public void setConfiguration(OpcuaConfiguration configuration) {
        Matcher matcher = getMatcher(configuration);
        code = matcher.group("transportCode");
        host = matcher.group("transportHost");
        port = matcher.group("transportPort");
        transportEndpoint = matcher.group("transportEndpoint");

        String portAddition = port != null ? ":" + port : "";
        endpoint = "opc." + code + "://" + host + portAddition + transportEndpoint;


        if (configuration.getSecurityPolicy() != null && !(configuration.getSecurityPolicy() == SecurityPolicy.NONE)) {
            try {
                openKeyStore(configuration);
            } catch (Exception e) {
                throw new PlcRuntimeException("Unable to open keystore, please confirm you have the correct permissions");
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

}
