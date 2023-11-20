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

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.NoSuchPaddingException;
import java.security.NoSuchAlgorithmException;
import java.security.Signature;
import java.util.Arrays;

public enum SecurityPolicy {
    NONE("http://opcfoundation.org/UA/SecurityPolicy#None",
        new MacSignatureAlgorithm(""),
        new EncryptionAlgorithm(""),
        new SignatureAlgorithm("", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
        new EncryptionAlgorithm(""),
        0, 0, 0, 1, 0
    ),

    Basic128Rsa15("http://opcfoundation.org/UA/SecurityPolicy#Basic128Rsa15",
        new MacSignatureAlgorithm("HmacSHA1"),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA1withRSA", "http://www.w3.org/2000/09/xmldsig#rsa-sha1"),
        new EncryptionAlgorithm("RSA/ECB/PKCS1Padding"),
        20, 16, 16, 16, 16
    ),

    Basic256("http://opcfoundation.org/UA/SecurityPolicy#Basic256",
        new MacSignatureAlgorithm("HmacSHA1"),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA1withRSA", "http://www.w3.org/2000/09/xmldsig#rsa-sha1"),
        new EncryptionAlgorithm("RSA/ECB/OAEPWithSHA-1AndMGF1Padding"),
        20, 24, 32, 16, 32
    ),

    Basic256Sha256("http://opcfoundation.org/UA/SecurityPolicy#Basic256Sha256",
        new MacSignatureAlgorithm("HmacSHA256"),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA256withRSA", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
        new EncryptionAlgorithm("RSA/ECB/OAEPWithSHA-1AndMGF1Padding"),
        32, 32, 32, 16, 32
    ),

    Aes128_Sha256_RsaOaep("http://opcfoundation.org/UA/SecurityPolicy#Aes128_Sha256_RsaOaep",
        new MacSignatureAlgorithm("HmacSHA256"),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA256withRSA", "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"),
        new EncryptionAlgorithm("RSA/ECB/OAEPWithSHA-1AndMGF1Padding"),
        32, 32, 16, 16, 32
    ),

    Aes256_Sha256_RsaPss("http://opcfoundation.org/UA/SecurityPolicy#Aes256_Sha256_RsaPss",
        new MacSignatureAlgorithm("HmacSHA256"),
        new EncryptionAlgorithm("AES/CBC/NoPadding"),
        new SignatureAlgorithm("SHA256withRSA/PSS", "http://opcfoundation.org/UA/security/rsa-pss-sha2-256"),
        new EncryptionAlgorithm("RSA/ECB/OAEPWithSHA256AndMGF1Padding"),
        32, 32, 32, 16, 32
    );


    private final String securityPolicyUri;
    private final MacSignatureAlgorithm symmetricSignatureAlgorithm;


    private final EncryptionAlgorithm symmetricEncryptionAlgorithm;
    private final SignatureAlgorithm asymmetricSignatureAlgorithm;
    private final EncryptionAlgorithm asymmetricEncryptionAlgorithm;
    private final int symmetricSignatureSize;
    private final int signatureKeySize;
    private final int encryptionKeySize;
    private final int encryptionBlockSize;
    private final int nonceLength;

    SecurityPolicy(String securityPolicyUri,
        MacSignatureAlgorithm symmetricSignatureAlgorithm,
        EncryptionAlgorithm symmetricEncryptionAlgorithm,
        SignatureAlgorithm asymmetricSignatureAlgorithm,
        EncryptionAlgorithm asymmetricEncryptionAlgorithm,
        int symmetricSignatureSize,
        int signatureKeySize, int encryptionKeySize,
        int encryptionBlockSize, int nonceLength
    ) {
        this.securityPolicyUri = securityPolicyUri;
        this.symmetricSignatureAlgorithm = symmetricSignatureAlgorithm;
        this.symmetricEncryptionAlgorithm = symmetricEncryptionAlgorithm;
        this.asymmetricSignatureAlgorithm = asymmetricSignatureAlgorithm;
        this.asymmetricEncryptionAlgorithm = asymmetricEncryptionAlgorithm;
        this.symmetricSignatureSize = symmetricSignatureSize;
        this.signatureKeySize = signatureKeySize;
        this.encryptionKeySize = encryptionKeySize;
        this.encryptionBlockSize = encryptionBlockSize;
        this.nonceLength = nonceLength;
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

    public int getSymmetricSignatureSize() {
        return symmetricSignatureSize;
    }

    public int getSignatureKeySize() {
        return signatureKeySize;
    }

    public int getEncryptionKeySize() {
        return encryptionKeySize;
    }

    public int getEncryptionBlockSize() {
        return encryptionBlockSize;
    }

    public int getNonceLength() {
        return nonceLength;
    }

    public static class MacSignatureAlgorithm {


        private final String name;

        MacSignatureAlgorithm(String name) {
            this.name = name;
        }

        public Mac getSignature() throws NoSuchAlgorithmException {
            return Mac.getInstance(name);
        }

        public String getName() {
            return name;
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
