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

import io.vavr.control.Try;
import java.nio.ByteBuffer;
import java.security.GeneralSecurityException;
import java.security.PrivateKey;
import java.security.Security;
import java.util.List;
import java.util.function.Supplier;
import javax.crypto.Cipher;
import org.apache.plc4x.java.opcua.protocol.chunk.Chunk;
import org.apache.plc4x.java.opcua.protocol.chunk.ChunkFactory;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenResponse;
import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.apache.plc4x.java.opcua.readwrite.SignatureData;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class EncryptionHandler {

    private final Logger logger = LoggerFactory.getLogger(EncryptionHandler.class);

    static {
        // Required for SecurityPolicy.Aes128_Sha128_RsaPss
        Security.addProvider(new BouncyCastleProvider());
    }

    private final Conversation conversation;

    private final SymmetricEncryptionHandler symmetricEncryptionHandler;
    private final AsymmetricEncryptionHandler asymmetricEncryptionHandler;

    public EncryptionHandler(Conversation conversation, PrivateKey senderPrivateKey) {
        this.conversation = conversation;
        this.symmetricEncryptionHandler = new SymmetricEncryptionHandler(conversation, conversation.getSecurityPolicy());
        this.asymmetricEncryptionHandler = new AsymmetricEncryptionHandler(conversation, conversation.getSecurityPolicy(), senderPrivateKey);
    }

    public List<MessagePDU> encodeMessage(MessagePDU message, Supplier<Integer> sequenceSupplier) {
        OpcuaProtocolLimits limits = conversation.getLimits();
        logger.debug("Encoding Message with Security policy {} and encoding limits {}", conversation.getSecurityPolicy(), limits);

        if (message instanceof OpcuaOpenRequest || message instanceof OpcuaOpenResponse) {
            Chunk chunk = new ChunkFactory().create(true, conversation.isSymmetricEncryptionEnabled(), conversation.isSymmetricSigningEnabled(),
                conversation.getSecurityPolicy(), limits,
                conversation.getLocalCertificate(), conversation.getRemoteCertificate()
            );
            return asymmetricEncryptionHandler.encodeMessage(chunk, message, sequenceSupplier);
        }

        Chunk chunk = new ChunkFactory().create(false, conversation.isSymmetricEncryptionEnabled(), conversation.isSymmetricSigningEnabled(),
            conversation.getSecurityPolicy(), limits,
            conversation.getLocalCertificate(), conversation.getRemoteCertificate()
        );
        return symmetricEncryptionHandler.encodeMessage(chunk, message, sequenceSupplier);
    }

    public MessagePDU decodeMessage(MessagePDU message) {
        OpcuaProtocolLimits limits = conversation.getLimits();
        logger.debug("Decoding Message with Security policy {} and encoding limits {}", conversation.getSecurityPolicy(), limits);

        if (message instanceof OpcuaOpenResponse || message instanceof OpcuaOpenRequest) {
            Chunk chunk = new ChunkFactory().create(true, conversation.isSymmetricEncryptionEnabled(), conversation.isSymmetricSigningEnabled(),
                conversation.getSecurityPolicy(), limits,
                conversation.getRemoteCertificate(), conversation.getLocalCertificate()
            );
            return asymmetricEncryptionHandler.decodeMessage(chunk, message);
        }
        Chunk chunk = new ChunkFactory().create(false, conversation.isSymmetricEncryptionEnabled(), conversation.isSymmetricSigningEnabled(),
            conversation.getSecurityPolicy(), limits,
            conversation.getRemoteCertificate(), conversation.getLocalCertificate()
        );
        return symmetricEncryptionHandler.decodeMessage(chunk, message);
    }

    public SignatureData createClientSignature() throws GeneralSecurityException {
        SecurityPolicy securityPolicy = conversation.getSecurityPolicy();
        byte[] lastServerNonce = conversation.getRemoteNonce();
        byte[] cert = Try.of(() -> conversation.getRemoteCertificate().getEncoded()).getOrElse(new byte[0]);
        byte[] bytes = ByteBuffer.allocate(cert.length + lastServerNonce.length).put(cert).put(lastServerNonce).array();
        byte[] signed = asymmetricEncryptionHandler.sign(bytes);
        return new SignatureData(new PascalString(securityPolicy.getAsymmetricSignatureAlgorithm().getUri()), new PascalByteString(signed.length, signed));
    }

    public byte[] encryptPassword(byte[] data) {
        try {
            Cipher cipher = Cipher.getInstance("RSA/ECB/OAEPWithSHA-1AndMGF1Padding");
            cipher.init(Cipher.ENCRYPT_MODE, this.conversation.getRemoteCertificate().getPublicKey());
            return cipher.doFinal(data);
        } catch (Exception e) {
            logger.error("Unable to encrypt Data", e);
            return null;
        }
    }

}
