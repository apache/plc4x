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

import java.nio.ByteBuffer;
import java.security.GeneralSecurityException;
import java.security.PrivateKey;
import java.security.Signature;
import java.security.Security;
import java.util.List;
import java.util.function.Supplier;
import java.security.cert.X509Certificate;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
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

    private final SecureChannelState conversation;

    private final SymmetricEncryptionHandler symmetricEncryptionHandler;
    private final AsymmetricEncryptionHandler asymmetricEncryptionHandler;

    public EncryptionHandler(SecureChannelState conversation, PrivateKey senderPrivateKey) {
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
                conversation.getLocalCertificate(), conversation.getRemoteCertificate(),
                conversation.getLocalCertificateChainSize()
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
        byte[] cert;
        try {
            cert = conversation.getRemoteCertificate().getEncoded();
        } catch (Exception e) {
            cert = new byte[0];
        }
        byte[] bytes = ByteBuffer.allocate(cert.length + lastServerNonce.length).put(cert).put(lastServerNonce).array();
        byte[] signed = asymmetricEncryptionHandler.sign(bytes);
        return new SignatureData(new PascalString(securityPolicy.getAsymmetricSignatureAlgorithm().getUri()), new PascalByteString(signed.length, signed));
    }

    /**
     * Proof that the client holds the private key of the certificate it presented in an
     * {@code X509IdentityToken}: a signature over the server certificate followed by the server
     * nonce (OPC UA Part 4, 5.6.3). Same bytes as {@link #createClientSignature()}, but signed
     * with the <em>user's</em> key rather than the application instance key, and with the
     * algorithm of the user token policy - which need not be the one securing the channel.
     *
     * @param userPrivateKey the private key belonging to the user certificate that was sent
     * @param policy         the security policy governing the user token
     */
    public SignatureData createUserTokenSignature(PrivateKey userPrivateKey, SecurityPolicy policy)
        throws GeneralSecurityException {
        byte[] serverNonce = conversation.getRemoteNonce();
        byte[] serverCertificate;
        try {
            serverCertificate = conversation.getRemoteCertificate().getEncoded();
        } catch (Exception e) {
            serverCertificate = new byte[0];
        }
        byte[] bytes = ByteBuffer.allocate(serverCertificate.length + serverNonce.length)
            .put(serverCertificate).put(serverNonce).array();

        Signature signature = policy.getAsymmetricSignatureAlgorithm().getSignature();
        signature.initSign(userPrivateKey);
        signature.update(bytes);
        byte[] signed = signature.sign();
        return new SignatureData(new PascalString(policy.getAsymmetricSignatureAlgorithm().getUri()),
            new PascalByteString(signed.length, signed));
    }

    /**
     * Encrypts a user token password with the server's public key.
     *
     * @param data   the encodeable password (length prefix, password, server nonce)
     * @param policy the security policy of the selected user token policy - it decides the
     *               algorithm, which is not necessarily the one used for the secure channel
     * @return the encrypted password, or null if it could not be encrypted
     */
    public byte[] encryptPassword(byte[] data, SecurityPolicy policy) {
        X509Certificate remoteCertificate = this.conversation.getRemoteCertificate();
        if (remoteCertificate == null) {
            // Without the server's certificate there is no key to encrypt to. Returning null here
            // used to leave the caller measuring the length of nothing, so a missing certificate
            // arrived as a NullPointerException from somewhere unrelated.
            throw new PlcRuntimeException(
                "Cannot encrypt the password: the server sent no certificate to encrypt it to");
        }
        try {
            Cipher cipher = policy.getAsymmetricEncryptionAlgorithm().getCipher();
            cipher.init(Cipher.ENCRYPT_MODE, remoteCertificate.getPublicKey());
            return cipher.doFinal(data);
        } catch (GeneralSecurityException e) {
            // Reported rather than swallowed. The password is about to go on the wire and the only
            // alternative to encrypting it is not sending it, so this cannot end in a null that
            // the caller discovers by dereferencing it.
            throw new PlcRuntimeException(
                "Could not encrypt the password with " + policy + " for the server's certificate", e);
        }
    }

}
