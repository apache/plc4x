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
import java.security.PrivateKey;
import java.security.Signature;
import java.util.Arrays;
import javax.crypto.Cipher;
import org.apache.plc4x.java.opcua.protocol.chunk.Chunk;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;

/**
 * Encryption + signing for the {@code OpcuaOpenRequest}/{@code Response} chunk
 * — the only chunks that use asymmetric (RSA) crypto. Subsequent traffic is
 * handled by {@link SymmetricEncryptionHandler}.
 */
public class AsymmetricEncryptionHandler extends BaseEncryptionHandler {

    private final PrivateKey senderPrivateKey;

    public AsymmetricEncryptionHandler(SecureChannelState conversation, SecurityPolicy securityPolicy, PrivateKey senderPrivateKey) {
        super(conversation, securityPolicy);
        this.senderPrivateKey = senderPrivateKey;
    }

    @Override
    protected void verify(byte[] chunkBytes, Chunk chunk, int messageLength) throws Exception {
        int signatureStart = messageLength - chunk.getSignatureSize();
        byte[] message = Arrays.copyOfRange(chunkBytes, 0, signatureStart);
        byte[] signatureData = Arrays.copyOfRange(chunkBytes, signatureStart, signatureStart + chunk.getSignatureSize());

        Signature signature = securityPolicy.getAsymmetricSignatureAlgorithm().getSignature();
        signature.initVerify(conversation.getRemoteCertificate().getPublicKey());
        signature.update(message);
        if (!signature.verify(signatureData)) {
            throw new IllegalArgumentException("Invalid signature");
        }
    }

    @Override
    protected int decryptInPlace(byte[] chunkBytes, Chunk chunk, int messageLength) throws Exception {
        int bodyStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();

        int bodySize = messageLength - bodyStart;
        int blockCount = bodySize / chunk.getCipherTextBlockSize();
        assert (bodySize % chunk.getCipherTextBlockSize() == 0);

        byte[] encrypted = Arrays.copyOfRange(chunkBytes, bodyStart, bodyStart + bodySize);
        byte[] plainText = new byte[chunk.getCipherTextBlockSize() * blockCount];

        Cipher cipher = securityPolicy.getAsymmetricEncryptionAlgorithm().getCipher();
        cipher.init(Cipher.DECRYPT_MODE, senderPrivateKey);

        int bodyLength = 0;
        for (int block = 0; block < blockCount; block++) {
            int pos = block * chunk.getCipherTextBlockSize();
            bodyLength += cipher.doFinal(encrypted, pos, chunk.getCipherTextBlockSize(), plainText, bodyLength);
        }

        // Copy the plaintext back into the chunk buffer at the body offset.
        System.arraycopy(plainText, 0, chunkBytes, bodyStart, bodyLength);
        return bodyLength;
    }

    @Override
    protected void encryptInPlace(byte[] chunkBytes, int securityHeaderSize, int plainTextBlockSize, int cipherTextBlockSize, int blockCount) throws Exception {
        int bodyStart = SECURE_MESSAGE_HEADER_SIZE + securityHeaderSize;
        byte[] copy = Arrays.copyOfRange(chunkBytes, bodyStart, bodyStart + (plainTextBlockSize * blockCount));
        byte[] encrypted = new byte[cipherTextBlockSize * blockCount];

        Cipher cipher = securityPolicy.getAsymmetricEncryptionAlgorithm().getCipher();
        cipher.init(Cipher.ENCRYPT_MODE, conversation.getRemoteCertificate().getPublicKey());

        for (int block = 0; block < blockCount; block++) {
            int pos = block * plainTextBlockSize;
            int target = block * cipherTextBlockSize;
            cipher.doFinal(copy, pos, plainTextBlockSize, encrypted, target);
        }

        System.arraycopy(encrypted, 0, chunkBytes, bodyStart, encrypted.length);
    }

    @Override
    public byte[] sign(byte[] contentsToSign) throws GeneralSecurityException {
        Signature signature = securityPolicy.getAsymmetricSignatureAlgorithm().getSignature();
        signature.initSign(senderPrivateKey);
        signature.update(contentsToSign);
        return signature.sign();
    }

}
