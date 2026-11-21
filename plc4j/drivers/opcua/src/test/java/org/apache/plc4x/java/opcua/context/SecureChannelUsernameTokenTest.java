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

import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.EnumSource;
import org.junit.jupiter.params.provider.MethodSource;
import org.junit.jupiter.params.provider.NullAndEmptySource;
import org.junit.jupiter.params.provider.ValueSource;

import javax.crypto.Cipher;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.cert.X509Certificate;
import java.util.HexFormat;
import java.util.List;
import java.util.Map.Entry;
import java.util.stream.Stream;

import static java.util.Map.entry;
import static org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter.LITTLE_ENDIAN;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

/** Password and algorithm encoding for UserNameIdentityToken (GH-2791, OPC UA Part 4). */
class SecureChannelUsernameTokenTest {

    private static final PascalString NULL_STRING = new PascalString(null);
    private static final byte[] NONCE = HexFormat.of().parseHex(
        "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f");
    private static Entry<PrivateKey, X509Certificate> server;

    private Conversation conversation;
    private OpcuaDriverContext driverContext;
    private OpcuaConfiguration configuration;

    @BeforeAll
    static void generateServerCertificate() {
        server = TestCertificateGenerator.generate(2048, "CN=server", 3600);
    }

    @BeforeEach
    void setUp() {
        conversation = mock(Conversation.class);
        when(conversation.getSecurityPolicy()).thenReturn(SecurityPolicy.NONE);
        when(conversation.getRemoteNonce()).thenReturn(NONCE);
        when(conversation.getRemoteCertificate()).thenReturn(server.getValue());
        when(conversation.createRequestHeader()).thenReturn(new RequestHeader(
            new NodeId(new NodeIdTwoByte((short) 0)), 0L, 0L, 0L, NULL_STRING, 0L, null));
        driverContext = mock(OpcuaDriverContext.class);
        when(driverContext.getEndpoint()).thenReturn("opc.tcp://localhost:4840");
        configuration = spy(new OpcuaConfiguration());
    }

    static Stream<Arguments> unencryptedPasswords() {
        return Stream.of((String) null, "", SecurityPolicy.NONE.getSecurityPolicyUri())
            .flatMap(uri -> Stream.of(
                Arguments.of(uri, null, new byte[0]),
                Arguments.of(uri, "", new byte[0]),
                Arguments.of(uri, "secret", new byte[]{115, 101, 99, 114, 101, 116}),
                Arguments.of(uri, "päss密碼🔒", HexFormat.of().parseHex("70c3a47373e5af86e7a2bcf09f9492")),
                Arguments.of(uri, "a".repeat(256), "a".repeat(256).getBytes(StandardCharsets.UTF_8))));
    }

    @ParameterizedTest
    @MethodSource("unencryptedPasswords")
    void sendsOnlyUtf8PasswordAndNullAlgorithm(String tokenPolicyUri, String password, byte[] expected) throws Exception {
        configuration.setAllowInsecureCredentials(true);
        UserNameIdentityToken token = usernameToken(password, SecurityPolicy.NONE, tokenPolicyUri);

        assertArrayEquals(expected, token.getPassword().getStringValue());
        assertNull(token.getEncryptionAlgorithm().getStringValue());
        ByteBuffer wire = serialize(token);
        assertArrayEquals("user-policy".getBytes(StandardCharsets.UTF_8), readByteString(wire));
        assertArrayEquals("user".getBytes(StandardCharsets.UTF_8), readByteString(wire));
        assertArrayEquals(expected, readByteString(wire), "ByteString length must count only password bytes");
        assertEquals(-1, wire.getInt(), "null algorithm must be encoded as a UA String length of -1");
        assertFalse(wire.hasRemaining(), "no nonce, inner length or algorithm bytes may follow");
        verify(conversation, never()).encryptPassword(any(), any());
        verify(conversation, never()).getRemoteNonce();
    }

    @Test
    void encodesTheUnencryptedAlgorithmAsANullUaString() throws Exception {
        configuration.setAllowInsecureCredentials(true);
        ByteBuffer wire = serialize(usernameToken("secret", SecurityPolicy.NONE, null));
        readByteString(wire); // policyId
        readByteString(wire); // userName
        readByteString(wire); // password

        assertEquals(-1, wire.getInt(), "use null rather than empty for compatibility with asyncua");
        assertFalse(wire.hasRemaining());
    }

    @Test
    void needsNeitherNonceNorCertificateForAnUnencryptedToken() {
        configuration.setAllowInsecureCredentials(true);
        when(conversation.getRemoteNonce()).thenReturn(null);
        when(conversation.getRemoteCertificate()).thenReturn(null);

        assertArrayEquals(new byte[]{112}, usernameToken("p", SecurityPolicy.NONE, null)
            .getPassword().getStringValue());
    }

    @ParameterizedTest
    @EnumSource(value = SecurityPolicy.class, names = "NONE", mode = EnumSource.Mode.EXCLUDE)
    void encryptsLengthPasswordAndNonceWithTheTokenPolicy(SecurityPolicy tokenPolicy) throws Exception {
        // Real RSA encryption and independent decryption, even though the channel policy is None.
        when(conversation.encryptPassword(any(), any())).thenAnswer(invocation ->
            new EncryptionHandler(conversation, server.getKey())
                .encryptPassword(invocation.getArgument(0), invocation.getArgument(1)));

        for (String password : new String[]{null, "", "secret", "päss密碼🔒", "a".repeat(64)}) {
            UserNameIdentityToken token = usernameToken(password, SecurityPolicy.NONE, tokenPolicy.getSecurityPolicyUri());
            byte[] expected = password == null ? new byte[0] : password.getBytes(StandardCharsets.UTF_8);
            Cipher decrypt = tokenPolicy.getAsymmetricEncryptionAlgorithm().getCipher();
            decrypt.init(Cipher.DECRYPT_MODE, server.getKey());
            ByteBuffer plaintext = ByteBuffer.wrap(decrypt.doFinal(token.getPassword().getStringValue()))
                .order(ByteOrder.LITTLE_ENDIAN);
            assertEquals(expected.length + NONCE.length, plaintext.getInt());
            byte[] actualPassword = new byte[expected.length];
            plaintext.get(actualPassword);
            assertArrayEquals(expected, actualPassword);
            byte[] actualNonce = new byte[NONCE.length];
            plaintext.get(actualNonce);
            assertArrayEquals(NONCE, actualNonce);
            assertFalse(plaintext.hasRemaining());

            String algorithm = tokenPolicy.getAsymmetricEncryptionAlgorithm().getUri();
            assertEquals(algorithm, token.getEncryptionAlgorithm().getStringValue());
            ByteBuffer wire = serialize(token);
            readByteString(wire); // policyId
            readByteString(wire); // userName
            assertArrayEquals(token.getPassword().getStringValue(), readByteString(wire));
            assertArrayEquals(algorithm.getBytes(StandardCharsets.UTF_8), readByteString(wire));
            assertFalse(wire.hasRemaining());
        }
        verify(conversation, times(5)).encryptPassword(any(), eq(tokenPolicy));
    }

    @ParameterizedTest
    @NullAndEmptySource
    void inheritsTheEndpointPolicyWhenTheTokenPolicyIsUnspecified(String tokenPolicyUri) {
        when(conversation.encryptPassword(any(), any())).thenReturn(new byte[]{42});
        UserNameIdentityToken token = usernameToken("secret", SecurityPolicy.Basic256Sha256, tokenPolicyUri);

        assertArrayEquals(new byte[]{42}, token.getPassword().getStringValue());
        verify(conversation).encryptPassword(any(), eq(SecurityPolicy.Basic256Sha256));
    }

    @Test
    void doesNotFallBackToPlaintextWhenPasswordEncryptionFails() {
        configuration.setAllowInsecureCredentials(true);
        when(conversation.encryptPassword(any(), any())).thenThrow(new PlcRuntimeException("encryption failed"));

        PlcRuntimeException exception = assertThrows(PlcRuntimeException.class,
            () -> usernameToken("secret", SecurityPolicy.NONE, SecurityPolicy.Basic256Sha256.getSecurityPolicyUri()));

        assertEquals("encryption failed", exception.getMessage());
        verify(conversation).encryptPassword(any(), eq(SecurityPolicy.Basic256Sha256));
    }

    @Test
    void honorsExplicitNoneTokenPolicyOnAnEncryptedChannel() throws Exception {
        when(conversation.getSecurityPolicy()).thenReturn(SecurityPolicy.Basic256Sha256);
        when(driverContext.getCertificateKeyPair()).thenReturn(new CertificateKeyPair(
            new KeyPair(server.getValue().getPublicKey(), server.getKey()), server.getValue()));
        when(driverContext.getThumbprint()).thenReturn(new PascalByteString(0, new byte[0]));
        when(conversation.createClientSignature()).thenReturn(new SignatureData(NULL_STRING, new PascalByteString(0, new byte[0])));
        configuration.setAllowInsecureCredentials(true);

        UserNameIdentityToken token = usernameToken("p", SecurityPolicy.Basic256Sha256,
            SecurityPolicy.NONE.getSecurityPolicyUri());

        assertArrayEquals(new byte[]{112}, token.getPassword().getStringValue());
        assertNull(token.getEncryptionAlgorithm().getStringValue());
        verify(conversation, never()).encryptPassword(any(), any());
    }

    @ParameterizedTest
    @NullAndEmptySource
    @ValueSource(strings = "http://opcfoundation.org/UA/SecurityPolicy#None")
    void refusesUnencryptedCredentialsUnlessExplicitlyAllowed(String tokenPolicyUri) {
        PlcRuntimeException exception = assertThrows(PlcRuntimeException.class,
            () -> usernameToken("secret", SecurityPolicy.NONE, tokenPolicyUri));

        assertTrue(exception.getMessage().contains("allow-insecure-credentials"));
        verify(conversation, never()).encryptPassword(any(), any());
    }

    @Test
    void keepsAnonymousAuthenticationIndependentOfCredentialsOptIn() {
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration, null);
        UserTokenPolicy policy = new UserTokenPolicy(new PascalString("anonymous-policy"),
            UserTokenType.userTokenTypeAnonymous, NULL_STRING, NULL_STRING, NULL_STRING);
        BinaryExtensionObjectWithMask identity = assertInstanceOf(BinaryExtensionObjectWithMask.class,
            channel.createActivateSessionRequest(endpoint(SecurityPolicy.NONE, policy)).getUserIdentityToken());

        assertEquals("anonymous-policy", assertInstanceOf(AnonymousIdentityToken.class, identity.getBody())
            .getPolicyId().getStringValue());
        verify(conversation, never()).encryptPassword(any(), any());
    }

    private UserNameIdentityToken usernameToken(String password, SecurityPolicy endpointPolicy, String tokenPolicyUri) {
        // The authentication object rejects null passwords; an omitted configuration password is allowed.
        when(configuration.getUsername()).thenReturn("user");
        when(configuration.getPassword()).thenReturn(password);
        SecureChannel channel = new SecureChannel(conversation, driverContext, configuration,
            password == null ? null : new PlcUsernamePasswordAuthentication("user", password));
        UserTokenPolicy policy = new UserTokenPolicy(new PascalString("user-policy"),
            UserTokenType.userTokenTypeUserName, NULL_STRING, NULL_STRING, new PascalString(tokenPolicyUri));
        BinaryExtensionObjectWithMask identity = assertInstanceOf(BinaryExtensionObjectWithMask.class,
            channel.createActivateSessionRequest(endpoint(endpointPolicy, policy)).getUserIdentityToken());
        return assertInstanceOf(UserNameIdentityToken.class, identity.getBody());
    }

    private Entry<EndpointDescription, UserTokenPolicy> endpoint(SecurityPolicy policy, UserTokenPolicy tokenPolicy) {
        return entry(new EndpointDescription(new PascalString("opc.tcp://localhost:4840"),
            mock(ApplicationDescription.class), new PascalByteString(0, new byte[0]),
            policy == SecurityPolicy.NONE ? MessageSecurityMode.messageSecurityModeNone
                : MessageSecurityMode.messageSecurityModeSignAndEncrypt,
            new PascalString(policy.getSecurityPolicyUri()), List.of(tokenPolicy), NULL_STRING, (short) 0), tokenPolicy);
    }

    private static ByteBuffer serialize(UserNameIdentityToken token) throws Exception {
        byte[] bytes = new byte[token.getLengthInBytes()];
        token.serialize(new WriteBufferByteBased(bytes, LITTLE_ENDIAN));
        return ByteBuffer.wrap(bytes).order(ByteOrder.LITTLE_ENDIAN);
    }

    private static byte[] readByteString(ByteBuffer wire) {
        byte[] bytes = new byte[wire.getInt()];
        wire.get(bytes);
        return bytes;
    }
}
