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

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.security.cert.X509Certificate;

import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.opcua.context.SecureChannel;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.configuration.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.LongDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;

public class OpcuaConfiguration implements PlcConnectionConfiguration {

    @ConfigurationParameter("protocol-code")
    private String protocolCode;

    @ConfigurationParameter("transport-code")
    private String transportCode;

    @ConfigurationParameter("transport-config")
    private String transportConfig;

    @ConfigurationParameter("discovery")
    @BooleanDefaultValue(true)
    @Description("Controls the feature of the discovery endpoint of an OPC UA server which every server\n" +
        "will propagate over an '<address>/discovery' endpoint. The most common issue here is that most servers are not correctly\n" +
        "configured and propagate the wrong external IP or URL address. If that is the case you can disable the discovery by\n" +
        "configuring it with a `false` value.\n" +
        "\n" +
        "The discovery phase is always conducted using `NONE` security policy.")
    private boolean discovery;

    @ConfigurationParameter("username")
    @Description("A username to authenticate to the OPCUA server with.")
    private String username;

    @ConfigurationParameter("password")
    @Description("A password to authenticate to the OPCUA server with.")
    private String password;

    @ConfigurationParameter("security-policy")
    @StringDefaultValue("NONE")
    @Description("The security policy applied to communication channel between driver and OPC UA server.\n" +
        "Default value assumes. Possible options are `NONE`, `Basic128Rsa15`, `Basic256`, `Basic256Sha256`, `Aes128_Sha256_RsaOaep`, `Aes256_Sha256_RsaPss`.")
    private SecurityPolicy securityPolicy;

    @ConfigurationParameter("message-security")
    @StringDefaultValue("SIGN_ENCRYPT")
    @Description("The security policy applied to messages exchanged after handshake phase.\n" +
        "Possible options are `NONE`, `SIGN`, `SIGN_ENCRYPT`.\n" +
        "This option is effective only when `securityPolicy` turns encryption (anything beyond `NONE`).")
    private MessageSecurity messageSecurity;

    @ConfigurationParameter("key-store-file")
    @Description("The Keystore file used to lookup client certificate and its private key.")
    private String keyStoreFile;

    @ConfigurationParameter("key-store-type")
    @StringDefaultValue("pkcs12")
    @Description("Keystore type used to access keystore and private key, defaults to PKCS (for Java 11+).\n" +
        "Possible values are between others `jks`, `pkcs11`, `dks`, `jceks`.")
    private String keyStoreType;

    @ConfigurationParameter("key-store-password")
    @Description("Java keystore password used to access keystore and private key.")
    private String keyStorePassword;

    @ConfigurationParameter("server-certificate-file")
    @Description("Filesystem location where server certificate is located, supported formats are `DER` and `PEM`.")
    private String serverCertificateFile;

    @ConfigurationParameter("trust-store-file")
    @Description("The trust store file used to verify server certificates and its chain.")
    private String trustStoreFile;

    @ConfigurationParameter("trust-store-type")
    @StringDefaultValue("pkcs12")
    @Description("Keystore type used to access keystore and private key, defaults to PKCS (for Java 11+).\n" +
        "Possible values are between others `jks`, `pkcs11`, `dks`, `jceks`.")
    private String trustStoreType;

    @ConfigurationParameter("trust-store-password")
    @Description("Password used to open trust store.")
    private String trustStorePassword;

    // the discovered certificate when discovery is enabled
    private X509Certificate serverCertificate;

    @ConfigurationParameter("channel-lifetime")
    @LongDefaultValue(3600000)
    @Description("Time for which negotiated secure channel, its keys and session remains open. Value in milliseconds, by default 60 minutes.")
    private long channelLifetime;

    @ConfigurationParameter("session-timeout")
    @LongDefaultValue(120000)
    @Description("Expiry time for opened secure session, value in milliseconds. Defaults to 2 minutes.")
    private long sessionTimeout;

    @ConfigurationParameter("negotiation-timeout")
    @LongDefaultValue(60000)
    @Description("Timeout for all negotiation steps prior acceptance of application level operations - this timeout applies to open secure channel, create session and close calls. Defaults to 60 seconds.")
    private long negotiationTimeout;

    @ConfigurationParameter("request-timeout")
    @LongDefaultValue(30000)
    @Description("Timeout for read/write/subscribe calls. Value in milliseconds.")
    private long requestTimeout;

    @ComplexConfigurationParameter(prefix = "encoding", defaultOverrides = {}, requiredOverrides = {})
    @Description("TCP encoding options")
    private Limits limits;

    @ConfigurationParameter("endpoint-host")
    @Description("Endpoint host used to establish secure channel.")
    private String endpointHost;

    @ConfigurationParameter("endpoint-port")
    @Description("Endpoint port used to establish secure channel")
    private Integer endpointPort;

    public String getProtocolCode() {
        return protocolCode;
    }

    public String getTransportCode() {
        return transportCode;
    }

    public String getTransportConfig() {
        return transportConfig;
    }

    public boolean isDiscovery() {
        return discovery;
    }

    public String getUsername() {
        return username;
    }

    public String getPassword() {
        return password;
    }

    public SecurityPolicy getSecurityPolicy() {
        return securityPolicy;
    }

    public MessageSecurity getMessageSecurity() {
        return messageSecurity;
    }

    public String getKeyStoreFile() {
        return keyStoreFile;
    }

    public String getKeyStoreType() {
        return keyStoreType;
    }

    public char[] getKeyStorePassword() {
        return keyStorePassword == null ? null : keyStorePassword.toCharArray();
    }

    public String getTrustStoreFile() {
        return trustStoreFile;
    }

    public String getTrustStoreType() {
        return trustStoreType;
    }

    public char[] getTrustStorePassword() {
        return trustStorePassword == null ? null : trustStorePassword.toCharArray();
    }

    public Limits getEncodingLimits() {
        return limits;
    }

    public X509Certificate getServerCertificate() {
        if (serverCertificate == null && serverCertificateFile != null) {
            // initialize server certificate from configured file
            try {
                byte[] certificateBytes = Files.readAllBytes(Path.of(serverCertificateFile));
                serverCertificate = SecureChannel.getX509Certificate(certificateBytes);
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
        return serverCertificate;
    }

    public void setServerCertificate(X509Certificate serverCertificate) {
        this.serverCertificate = serverCertificate;
    }

    public long getChannelLifetime() {
        return channelLifetime;
    }

    public long getSessionTimeout() {
        return sessionTimeout;
    }

    public long getRequestTimeout() {
        return requestTimeout;
    }

    public long getNegotiationTimeout() {
        return negotiationTimeout;
    }

    public String getEndpointHost() {
        return endpointHost;
    }

    public Integer getEndpointPort() {
        return endpointPort;
    }

    @Override
    public String toString() {
        return "OpcuaConfiguration{" +
            "discovery=" + discovery +
            ", username='" + username + '\'' +
            ", password='" + (password != null ? "******" : null) + '\'' +
            ", securityPolicy='" + securityPolicy + '\'' +
            ", keyStoreFile='" + keyStoreFile + '\'' +
            ", keyStorePassword='" + (keyStorePassword != null ? "******" : null) + '\'' +
            ", limits=" + limits +
            '}';
    }

}

