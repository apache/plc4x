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
package org.apache.plc4x.java.opcua.security;

import java.util.Arrays;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.MacSignatureAlgorithm;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class SymmetricKeys {


    private final Keys clientKeys;
    private final byte[] senderNonce;
    private final Keys serverKeys;
    private final byte[] receiverNonce;


    public SymmetricKeys(Keys clientKeys, byte[] senderNonce, Keys serverKeys, byte[] receiverNonce) {
        this.clientKeys = clientKeys;
        this.senderNonce = senderNonce;
        this.serverKeys = serverKeys;
        this.receiverNonce = receiverNonce;
    }

    public Keys getClientKeys() {
        return clientKeys;
    }

    public byte[] getSenderNonce() {
        return senderNonce;
    }

    public Keys getServerKeys() {
        return serverKeys;
    }

    public byte[] getReceiverNonce() {
        return receiverNonce;
    }

    // make sure that keys are
    public boolean matches(byte[] senderNonce, byte[] receiverNonce) {
        return Arrays.equals(this.senderNonce, senderNonce) && Arrays.equals(this.receiverNonce, receiverNonce);
    }

    public static SymmetricKeys generateKeyPair(byte[] senderNonce, byte[] receiverNonce, SecurityPolicy securityPolicy) {
        int signatureKeySize = securityPolicy.getSignatureKeySize();
        int encryptionKeySize = securityPolicy.getEncryptionKeySize();
        int cipherTextBlockSize = securityPolicy.getEncryptionBlockSize();

        MacSignatureAlgorithm policy = securityPolicy.getSymmetricSignatureAlgorithm();
        byte[] senderSignatureKey = createKey(receiverNonce, senderNonce, 0, signatureKeySize, policy);
        byte[] senderEncryptionKey = createKey(receiverNonce, senderNonce, signatureKeySize, encryptionKeySize, policy);
        byte[] senderInitializationVector = createKey(receiverNonce, senderNonce, signatureKeySize + encryptionKeySize, cipherTextBlockSize, policy);

        byte[] receiverSignatureKey = createKey(senderNonce, receiverNonce, 0, signatureKeySize, policy);
        byte[] receiverEncryptionKey = createKey(senderNonce, receiverNonce, signatureKeySize, encryptionKeySize, policy);
        byte[] receiverInitializationVector = createKey(senderNonce, receiverNonce, signatureKeySize + encryptionKeySize, cipherTextBlockSize, policy);

        return new SymmetricKeys(
            new Keys(senderSignatureKey, senderEncryptionKey, senderInitializationVector), senderNonce,
            new Keys(receiverSignatureKey, receiverEncryptionKey, receiverInitializationVector), receiverNonce
        );
    }

    private static byte[] createKey(byte[] secret, byte[] seed, int offset, int length, MacSignatureAlgorithm macSignatureAlgorithm) {
        try {
            Mac mac = macSignatureAlgorithm.getSignature();

            byte[] tempBytes = hash(macSignatureAlgorithm.getName(), secret, seed, mac, offset + length);
            byte[] key = new byte[length];

            System.arraycopy(tempBytes, offset, key, 0, key.length);

            return key;

        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private static byte[] hash(String transformation,
                               byte[] secret,
                               byte[] seed,
                               Mac mac,
                               int required) throws Exception {

        byte[] out = new byte[required];
        int offset = 0;
        int toCopy;
        byte[] a = seed;
        byte[] tmp;

        while (required > 0) {
            SecretKeySpec key = new SecretKeySpec(secret, transformation);
            mac.init(key);
            mac.update(a);
            a = mac.doFinal();
            mac.reset();
            mac.init(key);
            mac.update(a);
            mac.update(seed);
            tmp = mac.doFinal();
            toCopy = Math.min(required, tmp.length);
            System.arraycopy(tmp, 0, out, offset, toCopy);
            offset += toCopy;
            required -= toCopy;
        }

        return out;
    }


    public static class Keys {
        private final byte[] signatureKey;
        private final byte[] encryptionKey;
        private final byte[] initializationVector;

        public Keys(byte[] signatureKey, byte[] encryptionKey, byte[] initializationVector) {
            this.signatureKey = signatureKey;
            this.encryptionKey = encryptionKey;
            this.initializationVector = initializationVector;
        }

        public byte[] getSignatureKey() {
            return signatureKey;
        }

        public byte[] getEncryptionKey() {
            return encryptionKey;
        }

        public byte[] getInitializationVector() {
            return initializationVector;
        }
    }
}
