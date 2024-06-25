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

import java.security.GeneralSecurityException;
import java.security.MessageDigest;
import java.util.Arrays;
import org.apache.plc4x.java.opcua.protocol.chunk.Chunk;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.EncryptionAlgorithm;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.MacSignatureAlgorithm;
import org.apache.plc4x.java.opcua.security.SymmetricKeys;
import org.apache.plc4x.java.spi.generation.*;

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;


public class SymmetricEncryptionHandler extends BaseEncryptionHandler {

    private SymmetricKeys keys = null;
    private byte[] senderNonce;

    public SymmetricEncryptionHandler(Conversation channel, SecurityPolicy policy) {
        super(channel, policy);
    }

    protected void verify(WriteBufferByteBased buffer, Chunk chunk, int messageLength) throws Exception {
        int signatureStart = messageLength - chunk.getSignatureSize();
        byte[] message = buffer.getBytes(0, signatureStart);
        byte[] signatureData = buffer.getBytes(signatureStart, signatureStart + chunk.getSignatureSize());

        SymmetricKeys symmetricKeys = getSymmetricKeys(conversation.getLocalNonce(), conversation.getRemoteNonce());
        MacSignatureAlgorithm algorithm = securityPolicy.getSymmetricSignatureAlgorithm();
        Mac signature = algorithm.getSignature();
        signature.init(new SecretKeySpec(symmetricKeys.getServerKeys().getSignatureKey(), algorithm.getName()));
        signature.update(message);
        byte[] signatureBytes = signature.doFinal();

        if (!MessageDigest.isEqual(signatureData, signatureBytes)) {
            throw new IllegalArgumentException("Invalid signature");
        }
    }

    protected int decrypt(WriteBufferByteBased chunkBuffer, Chunk chunk, int messageLength) throws Exception {
        int bodyStart = 12 + chunk.getSecurityHeaderSize();

        int bodySize = messageLength - bodyStart;
        int blockCount = bodySize / chunk.getCipherTextBlockSize();
        assert(bodySize % chunk.getCipherTextBlockSize() == 0);

        byte[] encrypted = chunkBuffer.getBytes(bodyStart, bodyStart + bodySize);
        byte[] plainText = new byte[chunk.getCipherTextBlockSize() * blockCount];

        SymmetricKeys symmetricKeys = getSymmetricKeys(conversation.getLocalNonce(), conversation.getRemoteNonce());
        Cipher cipher = getCipher(symmetricKeys.getServerKeys(), securityPolicy.getSymmetricEncryptionAlgorithm(), Cipher.DECRYPT_MODE);

        int bodyLength = cipher.doFinal(encrypted, 0, encrypted.length, plainText, 0);

        chunkBuffer.setPos(bodyStart);
        chunkBuffer.writeByteArray("payload", plainText);
        return bodyLength;
    }

    protected void encrypt(WriteBufferByteBased buffer, int securityHeaderSize, int plainTextBlockSize, int cipherTextBlockSize, int blockCount) throws Exception {
        SymmetricKeys symmetricKeys = getSymmetricKeys(conversation.getLocalNonce(), conversation.getRemoteNonce());

        int bodyStart = 12 + securityHeaderSize;
        byte[] copy = buffer.getBytes(bodyStart, bodyStart + (plainTextBlockSize * blockCount));
        byte[] encrypted = new byte[cipherTextBlockSize * blockCount];

        EncryptionAlgorithm transformation = securityPolicy.getSymmetricEncryptionAlgorithm();
        Cipher cipher = getCipher(symmetricKeys.getClientKeys(), transformation, Cipher.ENCRYPT_MODE);
        cipher.doFinal(copy, 0, copy.length, encrypted, 0);

        buffer.setPos(bodyStart);
        buffer.writeByteArray("encrypted", encrypted);
    }

    protected byte[] sign(byte[] data)throws GeneralSecurityException {
        SymmetricKeys symmetricKeys = getSymmetricKeys(conversation.getLocalNonce(), conversation.getRemoteNonce());
        MacSignatureAlgorithm algorithm = securityPolicy.getSymmetricSignatureAlgorithm();
        Mac signature = algorithm.getSignature();
        signature.init(new SecretKeySpec(symmetricKeys.getClientKeys().getSignatureKey(), algorithm.getName()));
        signature.update(data);
        return signature.doFinal();
    }

    private SymmetricKeys getSymmetricKeys(byte[] senderNonce, byte[] receiverNonce) {
        if (keys == null) {
            this.senderNonce = senderNonce;
            keys = SymmetricKeys.generateKeyPair(senderNonce, receiverNonce, securityPolicy);
        } else if (!Arrays.equals(this.senderNonce, senderNonce)) {
            // sender nonce changed, we have to roll new security keys because security token
            // was just renewed.
            // We do not track receiver nonce, because they change at the same time
            this.senderNonce = senderNonce;
            keys = SymmetricKeys.generateKeyPair(senderNonce, receiverNonce, securityPolicy);
        }
        return keys;
    }

    private static Cipher getCipher(SymmetricKeys.Keys symmetricKeys, EncryptionAlgorithm transformation, int mode) throws NoSuchAlgorithmException, NoSuchPaddingException, InvalidKeyException, InvalidAlgorithmParameterException {
        Cipher cipher = transformation.getCipher();

        SecretKeySpec keySpec = new SecretKeySpec(symmetricKeys.getEncryptionKey(), "AES");
        IvParameterSpec ivSpec = new IvParameterSpec(symmetricKeys.getInitializationVector());

        cipher.init(mode, keySpec, ivSpec);
        return cipher;
    }

}
