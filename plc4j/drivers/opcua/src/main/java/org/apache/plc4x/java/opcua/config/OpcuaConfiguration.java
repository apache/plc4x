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

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.config.annotations.ComplexConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Secret;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.LongDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;

public class OpcuaConfiguration implements Configuration {

    @ConfigurationParameter("protocol-code")
    private String protocolCode;

    @ConfigurationParameter("transport-code")
    private String transportCode;

    @ConfigurationParameter("transport-config")
    private String transportConfig;

    @ConfigurationParameter("discovery")
    @BooleanDefaultValue(true)
    @Description("""
        Controls the feature of the discovery endpoint of an OPC UA server which every server
        will propagate over an '<address>/discovery' endpoint. The most common issue here is that most servers are not correctly
        configured and propagate the wrong external IP or URL address. If that is the case you can disable the discovery by
        configuring it with a `false` value.
        
        The discovery phase is always conducted using `NONE` security policy.""")
    private boolean discovery;

    @ConfigurationParameter("username")
    @Description("A username to authenticate to the OPCUA server with.")
    private String username;

    @Secret
    @ConfigurationParameter("password")
    @Description("A password to authenticate to the OPCUA server with.")
    private String password;

    @ConfigurationParameter("security-policy")
    @StringDefaultValue("Basic256Sha256")
    @Description("""
        The security policy applied to communication channel between driver and OPC UA server.
        Possible options are `NONE`, `Basic128Rsa15`, `Basic256`, `Basic256Sha256`, `Aes128_Sha256_RsaOaep`, `Aes256_Sha256_RsaPss`.
        `NONE` means the channel is neither signed nor encrypted, so anything on the path can read and
        change what is exchanged; it also leaves the server unauthenticated. A policy that signs and
        encrypts needs a trust anchor for the server's certificate - see `trust-store-file` and
        `server-certificate-file`.""")
    private SecurityPolicy securityPolicy;

    @ConfigurationParameter("message-security")
    @StringDefaultValue("SIGN_ENCRYPT")
    @Description("""
        The security policy applied to messages exchanged after handshake phase.
        Possible options are `NONE`, `SIGN`, `SIGN_ENCRYPT`.
        This option is effective only when `securityPolicy` turns encryption (anything beyond `NONE`).""")
    private MessageSecurity messageSecurity;

    @ConfigurationParameter("tls.keystore")
    @Description("The Keystore file used to lookup client certificate and its private key.")
    private String keyStoreFile;

    @ConfigurationParameter("tls.keystore-type")
    @StringDefaultValue("pkcs12")
    @Description("Keystore type used to access keystore and private key, defaults to PKCS (for Java 11+).\n" +
        "Possible values are between others `jks`, `pkcs11`, `dks`, `jceks`.")
    private String keyStoreType;

    @Secret
    @ConfigurationParameter("tls.keystore-password")
    @Description("Java keystore password used to access keystore and private key.")
    private String keyStorePassword;

    @ConfigurationParameter("generated-key-size")
    @IntDefaultValue(2048)
    @Description("Size in bits of the RSA key of the certificate the driver generates when no `tls.keystore` is configured. It is ignored when a key store is supplied, as the key then comes from that store. Some servers require a minimum size; 4096 is a common requirement.")
    private int generatedKeySize;

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

    @Secret
    @ConfigurationParameter("trust-store-password")
    @Description("Password used to open trust store.")
    private String trustStorePassword;

    @ConfigurationParameter("allow-insecure-credentials")
    @BooleanDefaultValue(false)
    @Description("""
        Allows a username and password to be sent over a channel that neither signs nor encrypts.
        Without this, a connection configured with credentials over an unprotected channel fails rather
        than putting the password on the wire where anything on the path can read it. Setting it warns.""")
    private boolean allowInsecureCredentials;

    @ConfigurationParameter("browse-max-references-per-node")
    @IntDefaultValue(65536)
    @Description("""
        Largest number of references the driver will collect for a single node while browsing.
        A Browse is answered in batches, each batch handing back a continuation point for the next, and
        the driver follows them until the server stops. A server that never stops would otherwise grow
        the collected list without limit. The same number is asked of the server as its per-node maximum,
        so it can stop before the driver has to. Set to 0 for no limit.""")
    private int browseMaxReferencesPerNode;

    @ConfigurationParameter("browse-max-total-nodes")
    @IntDefaultValue(1000000)
    @Description("""
        Largest number of nodes a single browse will expand. A browse walks whatever tree the
        server describes, and the driver has no way to know how large that is before walking it, so this
        bounds a tree that turns out to be unreasonable - or endless, if the server keeps naming nodes it
        has not named before. Set to 0 for no limit.""")
    private int browseMaxTotalNodes;

    @ConfigurationParameter("browse-max-depth")
    @IntDefaultValue(64)
    @Description("""
        How deep a browse will recurse into the node tree. Already-visited nodes are never
        expanded twice, so a reference cycle terminates on its own, but a server naming a fresh node at
        every level describes a tree with no bottom. Set to 0 for no limit.""")
    private int browseMaxDepth;

    // Spelled with the "tls." namespace although this driver has no tls transport: OPC UA
    // negotiates its own secure channel, so the setting is declared here rather than inherited
    // from a transport, and the namespace is written into the name by hand. See the tls
    // transport's "verify", which is the same concept one layer down.
    //
    // This replaced "insecure-certificate-verification", whose sense was the opposite. The
    // default is now to verify, so an address that misses the migration fails against a server
    // whose certificate does not validate, rather than connecting without checking it.
    @ConfigurationParameter("tls.verify")
    @BooleanDefaultValue(true)
    @Description("""
        Verifies the OPC UA server's certificate. Set to false to trust any certificate the server presents.
        Turning it off is UNSAFE: it leaves the connection open to man-in-the-middle attacks and defeats the
        integrity/authenticity guarantees of a signed secure channel. Only do so for local testing. In production,
        establish trust with `tls.trust-store` (chain validation) or `server-certificate-file` (certificate
        pinning) instead.""")
    private boolean verifyServerCertificate;

    // the discovered certificate when discovery is enabled
    private X509Certificate serverCertificate;

    @ConfigurationParameter("channel-lifetime-ms")
    @LongDefaultValue(3600000)
    @Description("Time for which negotiated secure channel, its keys and session remains open. Value in milliseconds, by default 60 minutes.")
    private long channelLifetime;

    @ConfigurationParameter("min-channel-lifetime-ms")
    @LongDefaultValue(5000)
    @Description("Shortest secure-channel lifetime this client will work with, in milliseconds. A server may revise the requested channel-lifetime-ms downwards, and the renewal schedule is derived from whatever it returns - so a very short lifetime means very frequent renewals, on an executor shared by every OPC UA connection in this JVM. A server-supplied lifetime below this value is raised to it and a warning is logged. If a server genuinely needs faster renewal, lower this value to accept it; the default is far below any lifetime a conforming server negotiates.")
    private long minChannelLifetime;

    @ConfigurationParameter("session-timeout-ms")
    @LongDefaultValue(120000)
    @Description("Expiry time for opened secure session, value in milliseconds. Defaults to 2 minutes.")
    private long sessionTimeout;

    @ConfigurationParameter("handshake-timeout-ms")
    @LongDefaultValue(60000)
    @Description("Timeout for all negotiation steps prior acceptance of application level operations - this timeout applies to open secure channel, create session and close calls. Defaults to 60 seconds.")
    private long negotiationTimeout;

    @ConfigurationParameter("request-timeout-ms")
    @LongDefaultValue(30000)
    @Description("Timeout for read/write/subscribe calls. Value in milliseconds.")
    private long requestTimeout;

    @ComplexConfigurationParameter(prefix = "encoding", defaultOverrides = {}, requiredOverrides = {})
    @Description("TCP encoding options")
    private Limits limits;

    @ConfigurationParameter("endpoint-host")
    @Description("Endpoint host used to establish secure channel connection. Used when client made connection to server which advertises different hostname than one used for network connection.")
    private String endpointHost;

    @ConfigurationParameter("endpoint-port")
    @Description("Endpoint port used to establish secure channel. Used when client made connection to server which advertises different port number than one used for network connection.")
    private Integer endpointPort;

    @ConfigurationParameter("subscription-queue-size")
    @LongDefaultValue(1)
    @Description("""
        Server-side queue depth per monitored item for subscriptions. 1 (default) keeps only
        the latest value between publishes; higher values retain intermediate changes for fast
        change-of-state tags, whose sampling rate can exceed the publishing (cycle) interval.""")
    private long subscriptionQueueSize;

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

    public int getGeneratedKeySize() {
        return generatedKeySize;
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

    public boolean isAllowInsecureCredentials() {
        return allowInsecureCredentials;
    }

    public int getBrowseMaxReferencesPerNode() {
        return browseMaxReferencesPerNode;
    }

    public int getBrowseMaxTotalNodes() {
        return browseMaxTotalNodes;
    }

    public int getBrowseMaxDepth() {
        return browseMaxDepth;
    }

    public void setBrowseMaxReferencesPerNode(int browseMaxReferencesPerNode) {
        this.browseMaxReferencesPerNode = browseMaxReferencesPerNode;
    }

    public void setBrowseMaxTotalNodes(int browseMaxTotalNodes) {
        this.browseMaxTotalNodes = browseMaxTotalNodes;
    }

    public void setBrowseMaxDepth(int browseMaxDepth) {
        this.browseMaxDepth = browseMaxDepth;
    }

    public boolean isVerifyServerCertificate() {
        return verifyServerCertificate;
    }

    /**
     * The filesystem path of a user-supplied server certificate to pin trust to,
     * or {@code null} if none was configured. Unlike {@link #getServerCertificate()},
     * this never returns a certificate that was discovered over the (unauthenticated)
     * discovery channel, so it is safe to use as a trust anchor.
     */
    public String getServerCertificateFile() {
        return serverCertificateFile;
    }

    public Limits getEncodingLimits() {
        return limits;
    }

    public X509Certificate getServerCertificate() {
        if (serverCertificate == null && serverCertificateFile != null) {
            try {
                byte[] certificateBytes = Files.readAllBytes(Path.of(serverCertificateFile));
                java.security.cert.CertificateFactory factory =
                    java.security.cert.CertificateFactory.getInstance("X.509");
                serverCertificate = (X509Certificate) factory.generateCertificate(
                    new java.io.ByteArrayInputStream(certificateBytes));
            } catch (IOException | java.security.cert.CertificateException e) {
                throw new RuntimeException(e);
            }
        }
        return serverCertificate;
    }

    /** Visible for tests that cover how a token is built rather than whether sending it is wise. */
    public void setAllowInsecureCredentials(boolean allowInsecureCredentials) {
        this.allowInsecureCredentials = allowInsecureCredentials;
    }

    public void setServerCertificate(X509Certificate serverCertificate) {
        this.serverCertificate = serverCertificate;
    }

    public long getMinChannelLifetime() {
        return minChannelLifetime;
    }

    public void setMinChannelLifetime(long minChannelLifetime) {
        this.minChannelLifetime = minChannelLifetime;
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

    public long getSubscriptionQueueSize() { return subscriptionQueueSize; }

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

