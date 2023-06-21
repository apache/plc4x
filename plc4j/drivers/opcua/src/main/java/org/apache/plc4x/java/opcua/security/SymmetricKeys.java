package org.apache.plc4x.java.opcua.security;

import org.apache.plc4x.java.opcua.security.SecurityPolicy.MacSignatureAlgorithm;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class SymmetricKeys {


    private final Keys clientKeys;
    private final Keys serverKeys;


    public SymmetricKeys(Keys clientKeys, Keys serverKeys) {
        this.clientKeys = clientKeys;
        this.serverKeys = serverKeys;
    }

    public Keys getClientKeys() {
        return clientKeys;
    }

    public Keys getServerKeys() {
        return serverKeys;
    }

    public static SymmetricKeys generateKeyPair(byte[] clientNonce, byte[] serverNonce, MacSignatureAlgorithm policy) {
        int signatureKeySize = policy.getKeySize();
        int encryptionKeySize = policy.getKeySize();
        int cipherTextBlockSize = 16;


        byte[] clientSignatureKey = createKey(serverNonce, clientNonce, 0, signatureKeySize, policy);
        byte[] clientEncryptionKey = createKey(serverNonce, clientNonce, signatureKeySize, encryptionKeySize, policy);
        byte[] clientInitializationVector = createKey(serverNonce, clientNonce, signatureKeySize + encryptionKeySize, cipherTextBlockSize, policy);


        byte[] serverSignatureKey = createKey(clientNonce, serverNonce, 0, signatureKeySize, policy);
        byte[] serverEncryptionKey = createKey(clientNonce, serverNonce, signatureKeySize, encryptionKeySize, policy);
        byte[] serverInitializationVector = createKey(clientNonce, serverNonce, signatureKeySize + encryptionKeySize, cipherTextBlockSize, policy);

        return new SymmetricKeys(
            new Keys(clientSignatureKey, clientEncryptionKey, clientInitializationVector),
            new Keys(serverSignatureKey, serverEncryptionKey, serverInitializationVector)
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
