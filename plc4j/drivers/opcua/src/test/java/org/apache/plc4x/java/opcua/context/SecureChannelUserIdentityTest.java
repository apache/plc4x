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

import static java.util.Map.entry;
import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

import java.nio.ByteBuffer;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.Signature;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.Map.Entry;
import org.apache.plc4x.java.api.authentication.PlcCertificateAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.ActivateSessionRequest;
import org.apache.plc4x.java.opcua.readwrite.ApplicationDescription;
import org.apache.plc4x.java.opcua.readwrite.BinaryExtensionObjectWithMask;
import org.apache.plc4x.java.opcua.readwrite.EndpointDescription;
import org.apache.plc4x.java.opcua.readwrite.MessageSecurityMode;
import org.apache.plc4x.java.opcua.readwrite.NodeId;
import org.apache.plc4x.java.opcua.readwrite.NodeIdTwoByte;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.apache.plc4x.java.opcua.readwrite.RequestHeader;
import org.apache.plc4x.java.opcua.readwrite.SignatureData;
import org.apache.plc4x.java.opcua.readwrite.UserTokenPolicy;
import org.apache.plc4x.java.opcua.readwrite.UserTokenType;
import org.apache.plc4x.java.opcua.readwrite.X509IdentityToken;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

/**
 * Authenticating with a user certificate - the {@code X509IdentityToken} of OPC UA Part 4 - which
 * the driver reaches through {@code getConnection(url, new PlcCertificateAuthentication(...))}
 * (GH-1845).
 */
class SecureChannelUserIdentityTest {

    private static final char[] PASSWORD = "changeit".toCharArray();
    private static final PascalString NULL_STRING = new PascalString(null);

    private Entry<PrivateKey, X509Certificate> user;
    private Entry<PrivateKey, X509Certificate> server;
    private Conversation conversation;
    private OpcuaDriverContext driverContext;
    private OpcuaConfiguration configuration;

    @BeforeEach
    void setUp() throws Exception {
        user = TestCertificateGenerator.generate(2048, "CN=operator", 3600);
        server = TestCertificateGenerator.generate(2048, "CN=server", 3600);

        conversation = Mockito.mock(Conversation.class);
        // A channel with no security at all: it proves the user token is signed with the user key
        // independently of whatever secures the channel.
        when(conversation.getSecurityPolicy()).thenReturn(SecurityPolicy.NONE);
        when(conversation.getRemoteCertificate()).thenReturn(server.getValue());
        when(conversation.getRemoteNonce()).thenReturn("server-nonce-0123456789abcdef0123".getBytes());
        when(conversation.createRequestHeader()).thenReturn(new RequestHeader(
            new NodeId(new NodeIdTwoByte((short) 0)), 0L, 0L, 0L, NULL_STRING, 0L, null));
        // Delegate to the real signing code, so the assertions look at real signatures.
        when(conversation.createUserTokenSignature(any(), any())).thenAnswer(invocation ->
            new EncryptionHandler(conversation, invocation.getArgument(0))
                .createUserTokenSignature(invocation.getArgument(0), invocation.getArgument(1)));

        driverContext = Mockito.mock(OpcuaDriverContext.class);
        when(driverContext.getEndpoint()).thenReturn("opc.tcp://localhost:4840");
        configuration = new OpcuaConfiguration();
    }

    @Test
    void acceptsCertificateAuthenticationInsteadOfRejectingIt() throws Exception {
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration, certificateAuthentication(null));

        assertInstanceOf(SecureChannel.class, channel);
    }

    @Test
    void sendsAnX509IdentityTokenCarryingTheUserCertificate() throws Exception {
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration, certificateAuthentication(null));

        ActivateSessionRequest request = channel.createActivateSessionRequest(
            certificateEndpoint(SecurityPolicy.Basic256Sha256));

        BinaryExtensionObjectWithMask identityToken =
            assertInstanceOf(BinaryExtensionObjectWithMask.class, request.getUserIdentityToken());
        X509IdentityToken token = assertInstanceOf(X509IdentityToken.class, identityToken.getBody());
        assertEquals("cert-policy", token.getPolicyId().getStringValue());
        assertArrayEquals(user.getValue().getEncoded(), token.getCertificateData().getStringValue(),
            "the token carries the DER encoded user certificate");
    }

    @Test
    void signsTheUserTokenWithTheUserKey() throws Exception {
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration, certificateAuthentication(null));

        ActivateSessionRequest request = channel.createActivateSessionRequest(
            certificateEndpoint(SecurityPolicy.Basic256Sha256));

        SignatureData signature = request.getUserTokenSignature();
        assertEquals(SecurityPolicy.Basic256Sha256.getAsymmetricSignatureAlgorithm().getUri(),
            signature.getAlgorithm().getStringValue());

        byte[] serverCertificate = server.getValue().getEncoded();
        byte[] serverNonce = conversation.getRemoteNonce();
        byte[] signed = ByteBuffer.allocate(serverCertificate.length + serverNonce.length)
            .put(serverCertificate).put(serverNonce).array();
        Signature verifier = SecurityPolicy.Basic256Sha256.getAsymmetricSignatureAlgorithm().getSignature();
        verifier.initVerify(user.getValue().getPublicKey());
        verifier.update(signed);
        assertTrue(verifier.verify(signature.getSignature().getStringValue()));

        assertFalse(java.util.Arrays.equals(
                request.getClientSignature().getSignature().getStringValue(),
                signature.getSignature().getStringValue()),
            "the user token signature is not the application instance signature the driver sent before");
    }

    /**
     * The alias picks the identity out of a key store that holds several.
     */
    @Test
    void usesTheCertificateOfTheRequestedAlias() throws Exception {
        Entry<PrivateKey, X509Certificate> engineer = TestCertificateGenerator.generate(2048, "CN=engineer", 3600);
        KeyStore keyStore = keyStoreWith(entry("operator", user), entry("engineer", engineer));
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration,
            new PlcCertificateAuthentication(keyStore, PASSWORD, "engineer"));

        ActivateSessionRequest request = channel.createActivateSessionRequest(
            certificateEndpoint(SecurityPolicy.Basic256Sha256));

        X509IdentityToken token = (X509IdentityToken)
            ((BinaryExtensionObjectWithMask) request.getUserIdentityToken()).getBody();
        assertArrayEquals(engineer.getValue().getEncoded(), token.getCertificateData().getStringValue());
    }

    /**
     * A user token policy of None leaves no algorithm to sign with, so the server could never
     * verify that the client holds the private key. Saying so beats sending an empty signature and
     * having the server reject the session for reasons it does not explain.
     */
    @Test
    void refusesACertificateTokenPolicyThatCarriesNoSecurity() throws Exception {
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration, certificateAuthentication(null));

        PlcRuntimeException exception = assertThrows(PlcRuntimeException.class,
            () -> channel.createActivateSessionRequest(certificateEndpoint(SecurityPolicy.NONE)));

        assertTrue(exception.getMessage().contains("security policy"), exception.getMessage());
    }

    @Test
    void matchesACertificateTokenPolicyWhenAUserCertificateIsSupplied() {
        assertTrue(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeCertificate), null, true));
        assertFalse(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeAnonymous), null, true));
        assertFalse(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeUserName), null, true));
    }

    /**
     * Without a user certificate nothing changes: username when credentials were given, anonymous
     * otherwise.
     */
    @Test
    void keepsMatchingUsernameAndAnonymousPoliciesWhenNoCertificateIsSupplied() {
        assertTrue(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeUserName), "user", false));
        assertFalse(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeAnonymous), "user", false));
        assertTrue(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeAnonymous), null, false));
        assertFalse(SecureChannel.isUserTokenPolicyCompatible(policy(UserTokenType.userTokenTypeCertificate), null, false));
    }

    @Test
    void stillRejectsAuthenticationTypesItCannotUse() {
        PlcRuntimeException exception = assertThrows(PlcRuntimeException.class,
            () -> new SecureChannel(conversation, driverContext, configuration, new UnsupportedAuthentication()));

        assertTrue(exception.getMessage().contains("supports"), exception.getMessage());
    }

    /**
     * Username authentication still produces a UserNameIdentityToken, and the user token signature
     * stays the application instance signature it has always been.
     */
    @Test
    void leavesUsernameAuthenticationUntouched() throws Exception {
        // This covers how the token is built, not whether sending it is wise: the endpoint below
        // protects nothing, which the driver otherwise refuses for a password.
        configuration.setAllowInsecureCredentials(true);
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration,
            new PlcUsernamePasswordAuthentication("user", "secret"));

        ActivateSessionRequest request = channel.createActivateSessionRequest(
            usernameEndpoint());

        assertInstanceOf(org.apache.plc4x.java.opcua.readwrite.UserNameIdentityToken.class,
            ((BinaryExtensionObjectWithMask) request.getUserIdentityToken()).getBody());
        assertEquals(request.getClientSignature(), request.getUserTokenSignature());
    }

    /**
     * A password sent over a channel that neither signs nor encrypts is readable by anything on the
     * path, and unlike a value it stays useful long after it was read.
     */
    @Test
    void refusesToSendAPasswordOverAnUnprotectedChannel() {
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration,
            new PlcUsernamePasswordAuthentication("user", "secret"));

        PlcRuntimeException thrown = assertThrows(PlcRuntimeException.class,
            () -> channel.createActivateSessionRequest(usernameEndpoint()));
        assertTrue(thrown.getMessage().contains("allow-insecure-credentials"),
            "the refusal should name the way out, but was: " + thrown.getMessage());
    }

    /**
     * The way out has to work, or the message above is telling operators to do something that does
     * not help them.
     */
    @Test
    void sendsThePasswordAnywayWhenExplicitlyAllowed() throws Exception {
        configuration.setAllowInsecureCredentials(true);
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration,
            new PlcUsernamePasswordAuthentication("user", "secret"));

        assertInstanceOf(org.apache.plc4x.java.opcua.readwrite.UserNameIdentityToken.class,
            ((BinaryExtensionObjectWithMask) channel.createActivateSessionRequest(usernameEndpoint())
                .getUserIdentityToken()).getBody());
    }

    private PlcCertificateAuthentication certificateAuthentication(String alias) throws Exception {
        return new PlcCertificateAuthentication(keyStoreWith(entry("operator", user)), PASSWORD, alias);
    }

    @SafeVarargs
    private KeyStore keyStoreWith(Entry<String, Entry<PrivateKey, X509Certificate>>... entries) throws Exception {
        KeyStore keyStore = KeyStore.getInstance("PKCS12");
        keyStore.load(null, PASSWORD);
        for (Entry<String, Entry<PrivateKey, X509Certificate>> element : entries) {
            keyStore.setKeyEntry(element.getKey(), element.getValue().getKey(), PASSWORD,
                new Certificate[]{element.getValue().getValue()});
        }
        return keyStore;
    }

    private static UserTokenPolicy policy(UserTokenType tokenType) {
        return new UserTokenPolicy(new PascalString("policy"), tokenType, NULL_STRING, NULL_STRING, NULL_STRING);
    }

    private Entry<EndpointDescription, UserTokenPolicy> certificateEndpoint(SecurityPolicy tokenPolicy) {
        UserTokenPolicy userTokenPolicy = new UserTokenPolicy(new PascalString("cert-policy"),
            UserTokenType.userTokenTypeCertificate, NULL_STRING, NULL_STRING,
            new PascalString(tokenPolicy.getSecurityPolicyUri()));
        return entry(endpoint(userTokenPolicy), userTokenPolicy);
    }

    private Entry<EndpointDescription, UserTokenPolicy> usernameEndpoint() {
        UserTokenPolicy userTokenPolicy = new UserTokenPolicy(new PascalString("user-policy"),
            UserTokenType.userTokenTypeUserName, NULL_STRING, NULL_STRING,
            new PascalString(SecurityPolicy.NONE.getSecurityPolicyUri()));
        return entry(endpoint(userTokenPolicy), userTokenPolicy);
    }

    private EndpointDescription endpoint(UserTokenPolicy userTokenPolicy) {
        return new EndpointDescription(
            new PascalString("opc.tcp://localhost:4840"),
            Mockito.mock(ApplicationDescription.class),
            new PascalByteString(0, new byte[0]),
            MessageSecurityMode.messageSecurityModeNone,
            new PascalString(SecurityPolicy.NONE.getSecurityPolicyUri()),
            List.of(userTokenPolicy),
            NULL_STRING,
            (short) 0);
    }

    private static final class UnsupportedAuthentication
        implements org.apache.plc4x.java.api.authentication.PlcAuthentication {
    }
}
