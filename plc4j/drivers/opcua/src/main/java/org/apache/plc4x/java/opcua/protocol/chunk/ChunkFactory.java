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

package org.apache.plc4x.java.opcua.protocol.chunk;

import io.vavr.control.Try;
import java.security.PublicKey;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import org.apache.commons.codec.digest.DigestUtils;
import org.apache.plc4x.java.opcua.context.Conversation;
import org.apache.plc4x.java.opcua.readwrite.OpcuaProtocolLimits;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;

public class ChunkFactory {

    public static final int ASYMMETRIC_SECURITY_HEADER_SIZE = 59;
    public static final int SYMMETRIC_SECURITY_HEADER_SIZE = 4;

    public Chunk create(boolean asymmetric, Conversation conversation) {
        return create(asymmetric,
            conversation.isSymmetricEncryptionEnabled(),
            conversation.isSymmetricSigningEnabled(),
            conversation.getSecurityPolicy(),
            conversation.getLimits(),
            conversation.getLocalCertificate(),
            conversation.getRemoteCertificate()
        );
    }

    public Chunk create(boolean asymmetric, boolean encrypted, boolean signed, SecurityPolicy securityPolicy,
        OpcuaProtocolLimits limits, X509Certificate localCertificate, X509Certificate remoteCertificate) {

        if (securityPolicy == SecurityPolicy.NONE) {
            return new Chunk(
                asymmetric ? ASYMMETRIC_SECURITY_HEADER_SIZE : SYMMETRIC_SECURITY_HEADER_SIZE,
                1,
                1,
                securityPolicy.getSymmetricSignatureSize(),
                (int) limits.getSendBufferSize(),
                asymmetric,
                false,
                false
            );
        }

        // asymmetric messages are always signed and encrypted, however non-asymmetric messages
        // exchanged after handshake might have message security mode set to NONE which results
        // in no overhead to communication
        boolean encryption = asymmetric || encrypted;
        boolean signing = asymmetric || signed;

        int localAsymmetricKeyLength = asymmetric ? keySize(localCertificate) : 0;
        int remoteAsymmetricKeyLength = asymmetric ? keySize(remoteCertificate) : 0;
        int localCertificateSize = asymmetric ? certificateBytes(localCertificate).length : 0;
        int serverCertificateThumbprint = asymmetric ? certificateThumbprint(remoteCertificate).length : 0;

        int asymmetricSecurityHeaderSize = (12 + securityPolicy.getSecurityPolicyUri().length() + localCertificateSize + serverCertificateThumbprint);
        int asymmetricCipherTextBlockSize = asymmetric ? (remoteAsymmetricKeyLength + 7) / 8 : 0;
        int plainTextBlockSize = asymmetric ? (remoteAsymmetricKeyLength + 7) / 8 : 0;

        int cipherTextBlockSize = asymmetric ? asymmetricCipherTextBlockSize : (encrypted ? securityPolicy.getEncryptionBlockSize() : 1);

        if (securityPolicy == SecurityPolicy.Basic128Rsa15) {
            // 12 + 56 + 674 + 20
            return new Chunk(
                asymmetric ? asymmetricSecurityHeaderSize : SYMMETRIC_SECURITY_HEADER_SIZE,
                cipherTextBlockSize,
                asymmetric ? plainTextBlockSize - 11 : (encrypted ? securityPolicy.getEncryptionBlockSize() : 1),
                asymmetric ? ((localAsymmetricKeyLength + 7) / 8) : securityPolicy.getSymmetricSignatureSize(),
                (int) limits.getSendBufferSize(),
                asymmetric,
                encryption,
                signing
            );
        } else if (securityPolicy == SecurityPolicy.Basic256) {
            return new Chunk(
                // 12 + 56 + 674 + 20
                asymmetric ? asymmetricSecurityHeaderSize : SYMMETRIC_SECURITY_HEADER_SIZE,
                cipherTextBlockSize,
                asymmetric ? plainTextBlockSize - 42 : (encrypted ? securityPolicy.getEncryptionBlockSize() : 1),
                asymmetric ? ((localAsymmetricKeyLength + 7) / 8) : securityPolicy.getSymmetricSignatureSize(),
                (int) limits.getSendBufferSize(),
                asymmetric,
                encryption,
                signing
            );
        } else if (securityPolicy == SecurityPolicy.Basic256Sha256) {
            return new Chunk(
                asymmetric ? asymmetricSecurityHeaderSize : SYMMETRIC_SECURITY_HEADER_SIZE,
                cipherTextBlockSize,
                asymmetric ? plainTextBlockSize - 42 : (encrypted ? securityPolicy.getEncryptionBlockSize() : 1),
                asymmetric ? ((localAsymmetricKeyLength + 7) / 8) : securityPolicy.getSymmetricSignatureSize(),
                (int) limits.getSendBufferSize(),
                asymmetric,
                encryption,
                signing
            );
        } else if (securityPolicy == SecurityPolicy.Aes128_Sha256_RsaOaep) {
            return new Chunk(
                asymmetric ? asymmetricSecurityHeaderSize : SYMMETRIC_SECURITY_HEADER_SIZE,
                cipherTextBlockSize,
                asymmetric ? plainTextBlockSize - 42 : (encrypted ? securityPolicy.getEncryptionBlockSize() : 1),
                asymmetric ? ((localAsymmetricKeyLength + 7) / 8) : securityPolicy.getSymmetricSignatureSize(),
                (int) limits.getSendBufferSize(),
                asymmetric,
                encryption,
                signing
            );
        } else if (securityPolicy == SecurityPolicy.Aes256_Sha256_RsaPss) {
            return new Chunk(
                asymmetric ? asymmetricSecurityHeaderSize : SYMMETRIC_SECURITY_HEADER_SIZE,
                cipherTextBlockSize,
                asymmetric ? plainTextBlockSize - 66 : (encrypted ? securityPolicy.getEncryptionBlockSize() : 1),
                asymmetric ? ((localAsymmetricKeyLength + 7) / 8) : securityPolicy.getSymmetricSignatureSize(),
                (int) limits.getSendBufferSize(),
                asymmetric,
                encryption,
                signing
            );
        }

        throw new IllegalArgumentException("Unsupported security policy " + securityPolicy.name() + "[" + securityPolicy.getSecurityPolicyUri() + "]");
    }

    private static int keySize(X509Certificate certificate) {
        PublicKey publicKey = certificate != null ? certificate.getPublicKey() : null;

        return (publicKey instanceof RSAPublicKey) ? ((RSAPublicKey) publicKey).getModulus().bitLength() : 0;
    }

    private static byte[] certificateThumbprint(X509Certificate certificate) {
        return DigestUtils.sha1(certificateBytes(certificate));
    }

    private static byte[] certificateBytes(X509Certificate certificate) {
        return Try.of(() -> certificate.getEncoded()).getOrElse(new byte[0]);
    }


}
