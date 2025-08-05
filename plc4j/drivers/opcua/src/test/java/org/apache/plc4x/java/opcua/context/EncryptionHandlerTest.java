/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.Map.Entry;
import java.util.function.Supplier;
import org.apache.commons.codec.digest.DigestUtils;
import org.apache.plc4x.java.opcua.TestCertificateGenerator;
import org.apache.plc4x.java.opcua.readwrite.BinaryPayload;
import org.apache.plc4x.java.opcua.readwrite.ChunkType;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.MessageSecurityMode;
import org.apache.plc4x.java.opcua.readwrite.OpcuaMessageRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.readwrite.OpenChannelMessageRequest;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.apache.plc4x.java.opcua.readwrite.Payload;
import org.apache.plc4x.java.opcua.readwrite.SecurityHeader;
import org.apache.plc4x.java.opcua.readwrite.SequenceHeader;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.mockito.Mockito;

class EncryptionHandlerTest {

    Supplier<Integer> sequenceSupplier = () -> 0;

    static List<Arguments> signKeyLengths() {
        return List.of(
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN, 1024, 1024),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN, 2048, 1024),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN, 1024, 2048),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN, 2048, 2048)
        );
    }

    static List<Arguments> encryptKeyLengths() {
        return List.of(
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN_ENCRYPT, 1024, 1024),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN_ENCRYPT, 2048, 1024),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN_ENCRYPT, 1024, 2048),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN_ENCRYPT, 2048, 2048),
            Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.SIGN_ENCRYPT, 1024, 1024),
            Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.SIGN_ENCRYPT, 2048, 1024),
            Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.SIGN_ENCRYPT, 1024, 2048),
            Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.SIGN_ENCRYPT, 2048, 2048)
        );
    }

    private Entry<CertificateKeyPair, CertificateKeyPair> initialize(int client, int server) throws Exception {
        Entry<PrivateKey, X509Certificate> clientKeyPair = TestCertificateGenerator.generate(client, "cn=client", 3600);
        Entry<PrivateKey, X509Certificate> serverKeyPair = TestCertificateGenerator.generate(server, "cn=server", 3600);

        X509Certificate clientCertificate = clientKeyPair.getValue();
        PublicKey clientPublicKey = clientCertificate.getPublicKey();
        X509Certificate serverCertificate = serverKeyPair.getValue();

        return entry(
            new CertificateKeyPair(new KeyPair(clientPublicKey, clientKeyPair.getKey()), clientCertificate),
            new CertificateKeyPair(new KeyPair(clientPublicKey, serverKeyPair.getKey()), serverCertificate)
        );
    }

    @ParameterizedTest
    @MethodSource("encryptKeyLengths")
    void testAsymmetricEncryption(SecurityPolicy securityPolicy, MessageSecurity messageSecurityMode, int client, int server) throws Exception {
        Entry<CertificateKeyPair, CertificateKeyPair> keyPairs = initialize(client, server);
        CertificateKeyPair clientKeyPair = keyPairs.getKey();
        CertificateKeyPair serverKeyPair = keyPairs.getValue();

        Conversation conversation = createSecureChannel(clientKeyPair.getCertificate(), serverKeyPair.getCertificate(),
            securityPolicy, messageSecurityMode, true, true
        );

        EncryptionHandler handler = new EncryptionHandler(conversation, clientKeyPair.getKeyPair().getPrivate());

        int[] messageSizes = {128};
        for (int messageSize : messageSizes) {
            byte[] messageBytes = new byte[messageSize];
            for (int i = 0; i < messageBytes.length; i++) {
                messageBytes[i] = (byte) i;
            }

            SecurityHeader securityHeader = new SecurityHeader(0, 1);
            SequenceHeader sequenceHeader = new SequenceHeader(1, 1);
            BinaryPayload payload = new BinaryPayload(sequenceHeader, messageBytes);

            OpcuaOpenRequest request = new OpcuaOpenRequest(ChunkType.FINAL,
                new OpenChannelMessageRequest(
                    (int) securityHeader.getSecureChannelId(),
                    new PascalString(securityPolicy.getSecurityPolicyUri()),
                    stringFromBytes(clientKeyPair.getCertificate().getEncoded()),
                    stringFromBytes(DigestUtils.sha1(serverKeyPair.getCertificate().getEncoded()))
                ),
                payload
            );
            List<MessagePDU> pdus = handler.encodeMessage(
                request, sequenceSupplier
            );
            assertEquals(1, pdus.size());

            // decrypt
            conversation = createSecureChannel(serverKeyPair.getCertificate(), clientKeyPair.getCertificate(), securityPolicy,
                messageSecurityMode, true, true);
            EncryptionHandler decrypter = new EncryptionHandler(conversation, serverKeyPair.getPrivateKey());
            MessagePDU decoded = decrypter.decodeMessage(pdus.get(0));
            assertInstanceOf(OpcuaOpenRequest.class, decoded);
            OpcuaOpenRequest decodedRequest = (OpcuaOpenRequest) decoded;
            SequenceHeader decodedSequenceHeader = decodedRequest.getMessage().getSequenceHeader();
            Payload decodedPayload = decodedRequest.getMessage();
            assertEquals(sequenceHeader.getSequenceNumber(), decodedSequenceHeader.getSequenceNumber());
            assertEquals(sequenceHeader.getRequestId(), decodedSequenceHeader.getRequestId());
            assertArrayEquals(messageBytes, ((BinaryPayload) decodedPayload).getPayload());
        }

    }

    @ParameterizedTest
    @MethodSource("encryptKeyLengths")
    void testAsymmetricEncryptionSign(SecurityPolicy securityPolicy, MessageSecurity messageSecurityMode, int client, int server) throws Exception {
        Entry<CertificateKeyPair, CertificateKeyPair> keyPairs = initialize(client, server);
        CertificateKeyPair clientKeyPair = keyPairs.getKey();
        CertificateKeyPair serverKeyPair = keyPairs.getValue();

        Conversation secureChannel = createSecureChannel(clientKeyPair.getCertificate(), serverKeyPair.getCertificate(),
            securityPolicy, messageSecurityMode, true, true);

        EncryptionHandler handler = new EncryptionHandler(secureChannel, clientKeyPair.getPrivateKey());

        int[] messageSizes = {128};
        for (int messageSize : messageSizes) {
            byte[] messageBytes = new byte[messageSize];
            for (int i = 0; i < messageBytes.length; i++) {
                messageBytes[i] = (byte) i;
            }

            SecurityHeader securityHeader = new SecurityHeader(0, 1);
            SequenceHeader sequenceHeader = new SequenceHeader(1, 1);
            BinaryPayload payload = new BinaryPayload(sequenceHeader, messageBytes);

            OpcuaOpenRequest request = new OpcuaOpenRequest(ChunkType.FINAL,
                new OpenChannelMessageRequest(
                    (int) securityHeader.getSecureChannelId(),
                    new PascalString(securityPolicy.getSecurityPolicyUri()),
                    stringFromBytes(clientKeyPair.getCertificate().getEncoded()),
                    stringFromBytes(DigestUtils.sha1(serverKeyPair.getCertificate().getEncoded()))
                ),
                payload
            );
            List<MessagePDU> pdus = handler.encodeMessage(
                request, sequenceSupplier
            );
            assertEquals(1, pdus.size());

            // decrypt
            secureChannel = createSecureChannel(serverKeyPair.getCertificate(), clientKeyPair.getCertificate(), securityPolicy,
                messageSecurityMode, true, true);
            EncryptionHandler decryptHandler = new EncryptionHandler(secureChannel, serverKeyPair.getPrivateKey());
            MessagePDU decoded = decryptHandler.decodeMessage(pdus.get(0));
            assertInstanceOf(OpcuaOpenRequest.class, decoded);
            OpcuaOpenRequest decodedRequest = (OpcuaOpenRequest) decoded;
            SequenceHeader decodedSequenceHeader = decodedRequest.getMessage().getSequenceHeader();
            Payload decodedPayload = decodedRequest.getMessage();
            assertEquals(sequenceHeader.getSequenceNumber(), decodedSequenceHeader.getSequenceNumber());
            assertEquals(sequenceHeader.getRequestId(), decodedSequenceHeader.getRequestId());
            assertArrayEquals(messageBytes, ((BinaryPayload) decodedPayload).getPayload());
        }

    }

    @ParameterizedTest
    @MethodSource("signKeyLengths")
    void testSymmetricEncryption(SecurityPolicy securityPolicy, MessageSecurity messageSecurityMode, int client, int server) throws Exception {
        Entry<CertificateKeyPair, CertificateKeyPair> keyPairs = initialize(client, server);
        CertificateKeyPair clientKeyPair = keyPairs.getKey();
        CertificateKeyPair serverKeyPair = keyPairs.getValue();

        Conversation secureChannel = createSecureChannel(clientKeyPair.getCertificate(), serverKeyPair.getCertificate(), securityPolicy,
            messageSecurityMode, true, true);

        EncryptionHandler handler = new EncryptionHandler(secureChannel, clientKeyPair.getPrivateKey());

        int[] messageSizes = {128};
        for (int messageSize : messageSizes) {
            byte[] messageBytes = new byte[messageSize];
            for (int i = 0; i < messageBytes.length; i++) {
                messageBytes[i] = (byte) i;
            }

            SecurityHeader securityHeader = new SecurityHeader(0, 1);
            SequenceHeader sequenceHeader = new SequenceHeader(1, 1);
            BinaryPayload payload = new BinaryPayload(sequenceHeader, messageBytes);

            OpcuaMessageRequest request = new OpcuaMessageRequest(ChunkType.FINAL,
                securityHeader,
                payload
            );
            List<MessagePDU> pdus = handler.encodeMessage(
                request, sequenceSupplier
            );
            assertEquals(1, pdus.size());

            // decrypt
            secureChannel = createSecureChannel(serverKeyPair.getCertificate(), clientKeyPair.getCertificate(), securityPolicy,
                messageSecurityMode, true, true);
            EncryptionHandler decryptHandler = new EncryptionHandler(secureChannel, serverKeyPair.getPrivateKey());
            MessagePDU decoded = decryptHandler.decodeMessage(pdus.get(0));
            OpcuaMessageRequest decodedRequest = (OpcuaMessageRequest) decoded;
            SequenceHeader decodedSequenceHeader = decodedRequest.getMessage().getSequenceHeader();
            Payload decodedPayload = decodedRequest.getMessage();
            assertEquals(sequenceHeader.getSequenceNumber(), decodedSequenceHeader.getSequenceNumber());
            assertEquals(sequenceHeader.getRequestId(), decodedSequenceHeader.getRequestId());
            assertArrayEquals(messageBytes, ((BinaryPayload) decodedPayload).getPayload());
        }
    }

    @ParameterizedTest
    @MethodSource("signKeyLengths")
    void testSymmetricEncryptionSign(SecurityPolicy securityPolicy, MessageSecurity messageSecurityMode, int client, int server) throws Exception {
        Entry<CertificateKeyPair, CertificateKeyPair> keyPairs = initialize(client, server);
        CertificateKeyPair clientKeyPair = keyPairs.getKey();
        CertificateKeyPair serverKeyPair = keyPairs.getValue();

        Conversation secureChannel = createSecureChannel(clientKeyPair.getCertificate(), serverKeyPair.getCertificate(), securityPolicy,
            messageSecurityMode, true, true);

        EncryptionHandler handler = new EncryptionHandler(secureChannel, clientKeyPair.getPrivateKey());

        int[] messageSizes = {128};
        for (int messageSize : messageSizes) {
            byte[] messageBytes = new byte[messageSize];
            for (int i = 0; i < messageBytes.length; i++) {
                messageBytes[i] = (byte) i;
            }

            SecurityHeader securityHeader = new SecurityHeader(0, 1);
            SequenceHeader sequenceHeader = new SequenceHeader(1, 1);
            BinaryPayload payload = new BinaryPayload(sequenceHeader, messageBytes);

            OpcuaMessageRequest request = new OpcuaMessageRequest(ChunkType.FINAL,
                securityHeader,
                payload
            );
            List<MessagePDU> pdus = handler.encodeMessage(
                request, sequenceSupplier
            );
            assertEquals(1, pdus.size());

            // decrypt
            secureChannel = createSecureChannel(serverKeyPair.getCertificate(), clientKeyPair.getCertificate(), securityPolicy,
                messageSecurityMode, true, true);
            EncryptionHandler decryptHandler = new EncryptionHandler(secureChannel, serverKeyPair.getPrivateKey());
            MessagePDU decoded = decryptHandler.decodeMessage(pdus.get(0));
            OpcuaMessageRequest decodedRequest = (OpcuaMessageRequest) decoded;
            SequenceHeader decodedSequenceHeader = decodedRequest.getMessage().getSequenceHeader();
            Payload decodedPayload = decodedRequest.getMessage();
            assertEquals(sequenceHeader.getSequenceNumber(), decodedSequenceHeader.getSequenceNumber());
            assertEquals(sequenceHeader.getRequestId(), decodedSequenceHeader.getRequestId());
            assertArrayEquals(messageBytes, ((BinaryPayload) decodedPayload).getPayload());
        }
    }

    private static PascalByteString stringFromBytes(byte[] bytes) {
        return new PascalByteString(bytes.length, bytes);
    }

    private static Conversation createSecureChannel(X509Certificate localCertificate, X509Certificate remoteCertificate, SecurityPolicy securityPolicy,
        MessageSecurity messageSecurity, boolean encrypted, boolean signed) {
        OpcuaProtocolLimits limits = new OpcuaProtocolLimits(8196, 8196, 8196 * 10, 10);
        Conversation conversation = Mockito.mock(Conversation.class);
        when(conversation.getLimits()).thenReturn(limits);
        when(conversation.getLocalCertificate()).thenReturn(localCertificate);
        when(conversation.getRemoteCertificate()).thenReturn(remoteCertificate);
        when(conversation.getSecurityPolicy()).thenReturn(securityPolicy);
        when(conversation.getMessageSecurity()).thenReturn(messageSecurity);
        when(conversation.isSymmetricEncryptionEnabled()).thenReturn(encrypted);
        when(conversation.isSymmetricSigningEnabled()).thenReturn(signed);
        when(conversation.getLocalNonce()).thenReturn(new byte[32]);
        when(conversation.getRemoteNonce()).thenReturn(new byte[32]);
        return conversation;
    }

}