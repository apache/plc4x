package org.apache.plc4x.java.opcua.security;

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.NoSuchPaddingException;
import java.security.NoSuchAlgorithmException;
import java.security.Signature;
import java.util.Arrays;

public enum SecurityPolicy {
    NONE("http://opcfoundation.org/UA/SecurityPolicy#None",
        new MacSignatureAlgorithm("", 0, 32),
        new EncryptionAlgorithm(""),
        new SignatureAlgorithm("", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
        new EncryptionAlgorithm(""),
        1),

    Basic128Rsa15("http://opcfoundation.org/UA/SecurityPolicy#Basic128Rsa15",
        new MacSignatureAlgorithm("HmacSHA1", 20, 16),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA1withRSA", "http://www.w3.org/2000/09/xmldsig#rsa-sha1"),
        new EncryptionAlgorithm("RSA/ECB/PKCS1Padding"),
        11),
    Basic256Sha256("http://opcfoundation.org/UA/SecurityPolicy#Basic256Sha256",
        new MacSignatureAlgorithm("HmacSHA256", 32, 32),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA256withRSA", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
        new EncryptionAlgorithm("RSA/ECB/OAEPWithSHA-1AndMGF1Padding"),
        42);


    private final String securityPolicyUri;
    private final MacSignatureAlgorithm symmetricSignatureAlgorithm;


    private final EncryptionAlgorithm symmetricEncryptionAlgorithm;
    private final SignatureAlgorithm asymmetricSignatureAlgorithm;
    private final EncryptionAlgorithm asymmetricEncryptionAlgorithm;
    private final int asymmetricPlainBlock;

    SecurityPolicy(String securityPolicyUri,
                   MacSignatureAlgorithm symmetricSignatureAlgorithm,
                   EncryptionAlgorithm symmetricEncryptionAlgorithm,
                   SignatureAlgorithm asymmetricSignatureAlgorithm,
                   EncryptionAlgorithm asymmetricEncryptionAlgorithm,
                   int asymmetricPlainBlock) {
        this.securityPolicyUri = securityPolicyUri;
        this.symmetricSignatureAlgorithm = symmetricSignatureAlgorithm;
        this.symmetricEncryptionAlgorithm = symmetricEncryptionAlgorithm;
        this.asymmetricSignatureAlgorithm = asymmetricSignatureAlgorithm;
        this.asymmetricEncryptionAlgorithm = asymmetricEncryptionAlgorithm;
        this.asymmetricPlainBlock = asymmetricPlainBlock;
    }

    public static SecurityPolicy findByName(String securityPolicy) {
        return Arrays.stream(values())
            .filter(v -> v.name().equalsIgnoreCase(securityPolicy))
            .findAny()
            .orElseThrow();
    }

    public MacSignatureAlgorithm getSymmetricSignatureAlgorithm() {
        return symmetricSignatureAlgorithm;
    }

    public String getSecurityPolicyUri() {
        return securityPolicyUri;
    }

    public SignatureAlgorithm getAsymmetricSignatureAlgorithm() {
        return asymmetricSignatureAlgorithm;
    }

    public EncryptionAlgorithm getAsymmetricEncryptionAlgorithm() {
        return asymmetricEncryptionAlgorithm;
    }

    public EncryptionAlgorithm getSymmetricEncryptionAlgorithm() {
        return symmetricEncryptionAlgorithm;
    }

    public int getAsymmetricPlainBlock() {
        return asymmetricPlainBlock;
    }

    public static class MacSignatureAlgorithm {


        private final String name;
        private final int symmetricSignatureSize;
        private final int keySize;

        MacSignatureAlgorithm(String name, int symmetricSignatureSize, int keySize) {
            this.name = name;
            this.symmetricSignatureSize = symmetricSignatureSize;
            this.keySize = keySize;
        }

        public Mac getSignature() throws NoSuchAlgorithmException {
            return Mac.getInstance(name);
        }

        public String getName() {
            return name;
        }

        public int getSymmetricSignatureSize() {
            return symmetricSignatureSize;
        }

        public int getKeySize() {
            return keySize;
        }
    }


    public static class SignatureAlgorithm {
        private final String name;
        private String uri;

        SignatureAlgorithm(String name, String uri) {
            this.name = name;
            this.uri = uri;
        }

        public Signature getSignature() throws NoSuchAlgorithmException {
            return Signature.getInstance(name);
        }

        public String getUri() {
            return uri;
        }
    }

    public static class EncryptionAlgorithm {
        private final String name;

        EncryptionAlgorithm(String name) {
            this.name = name;
        }

        public Cipher getCipher() throws NoSuchPaddingException, NoSuchAlgorithmException {
            return Cipher.getInstance(name);
        }
    }


}
