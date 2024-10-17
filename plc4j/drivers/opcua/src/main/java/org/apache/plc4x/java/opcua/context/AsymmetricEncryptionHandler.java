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
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.Signature;
import java.security.SignatureException;
import javax.crypto.Cipher;
import org.apache.plc4x.java.opcua.protocol.chunk.Chunk;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

public class AsymmetricEncryptionHandler extends BaseEncryptionHandler {

    private final PrivateKey senderPrivateKey;

    public AsymmetricEncryptionHandler(Conversation conversation, SecurityPolicy securityPolicy, PrivateKey senderPrivateKey) {
        super(conversation, securityPolicy);
        this.senderPrivateKey = senderPrivateKey;
    }

    protected void verify(WriteBufferByteBased buffer, Chunk chunk, int messageLength) throws Exception {
        int signatureStart = messageLength - chunk.getSignatureSize();
        byte[] message = buffer.getBytes(0, signatureStart);
        byte[] signatureData = buffer.getBytes(signatureStart, signatureStart + chunk.getSignatureSize());

        Signature signature = securityPolicy.getAsymmetricSignatureAlgorithm().getSignature();
        signature.initVerify(conversation.getRemoteCertificate().getPublicKey());
        signature.update(message);
        if (signature.verify(signatureData)) {
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

        Cipher cipher = securityPolicy.getAsymmetricEncryptionAlgorithm().getCipher();
        cipher.init(Cipher.DECRYPT_MODE, senderPrivateKey);

        int bodyLength = 0;
        for (int block = 0; block < blockCount; block++) {
            int pos = block * chunk.getCipherTextBlockSize();

            bodyLength += cipher.doFinal(encrypted, pos, chunk.getCipherTextBlockSize(), plainText, bodyLength);
        }

        chunkBuffer.setPos(bodyStart);
        byte[] decrypted = new byte[bodyLength];
        System.arraycopy(plainText, 0, decrypted, 0, bodyLength);
        chunkBuffer.writeByteArray("payload", decrypted);
        return bodyLength;
    }

    protected void encrypt(WriteBufferByteBased buffer, int securityHeaderSize, int plainTextBlockSize, int cipherTextBlockSize, int blockCount) throws Exception {
        int bodyStart = 12 + securityHeaderSize;
        byte[] copy = buffer.getBytes(bodyStart, bodyStart + (plainTextBlockSize * blockCount));
        byte[] encrypted = new byte[cipherTextBlockSize * blockCount];

        // copy of bytes from sequence header over payload, padding bytes and signature
        Cipher cipher = securityPolicy.getAsymmetricEncryptionAlgorithm().getCipher();
        cipher.init(Cipher.ENCRYPT_MODE, conversation.getRemoteCertificate().getPublicKey());

        for (int block = 0; block < blockCount; block++) {
            int pos = block * plainTextBlockSize;
            int target = block * cipherTextBlockSize;

            cipher.doFinal(copy, pos, plainTextBlockSize, encrypted, target);
        }

        buffer.setPos(bodyStart);
        buffer.writeByteArray("encrypted", encrypted);
    }

    public byte[] sign(byte[] contentsToSign) throws GeneralSecurityException {
        Signature signature = securityPolicy.getAsymmetricSignatureAlgorithm().getSignature();
        signature.initSign(senderPrivateKey);
        signature.update(contentsToSign);
        return signature.sign();
    }

}
