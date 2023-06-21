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

import static org.apache.plc4x.java.spi.generation.ByteOrder.LITTLE_ENDIAN;

import io.vavr.control.Try;
import java.io.ByteArrayInputStream;
import java.nio.ByteBuffer;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Security;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import javax.crypto.Cipher;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.opcua.protocol.OpcuaProtocolLogic;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenRequest;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenResponse;
import org.apache.plc4x.java.opcua.readwrite.PascalByteString;
import org.apache.plc4x.java.opcua.readwrite.PascalString;
import org.apache.plc4x.java.opcua.readwrite.SignatureData;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class EncryptionHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocolLogic.class);

    static {
        // Required for SecurityPolicy.Aes128_Sha128_RsaPss
        Security.addProvider(new BouncyCastleProvider());
    }


    private X509Certificate serverCertificate;
    private X509Certificate clientCertificate;
    private PrivateKey clientPrivateKey;
    private PublicKey clientPublicKey;
    private final SecurityPolicy securitypolicy;

    private byte[] clientNonce = null;
    private byte[] serverNonce = null;
    private final SymmetricEncryptionHandler symmetricEncryptionHandler;
    private final AsymmetricEncryptionHandler asymmetricEncryptionHandler;

    public EncryptionHandler(CertificateKeyPair ckp, byte[] senderCertificate, SecurityPolicy securityPolicy) {
        if (ckp != null) {
            this.clientPrivateKey = ckp.getKeyPair().getPrivate();
            this.clientPublicKey = ckp.getKeyPair().getPublic();
            this.clientCertificate = ckp.getCertificate();
        }
        if (senderCertificate != null) {
            this.serverCertificate = getCertificateX509(senderCertificate);
        }
        this.securitypolicy = securityPolicy;
        this.symmetricEncryptionHandler = new SymmetricEncryptionHandler(securityPolicy);
        this.asymmetricEncryptionHandler = new AsymmetricEncryptionHandler(serverCertificate, clientCertificate, clientPrivateKey, clientPublicKey, securitypolicy);
    }

    public void setServerCertificate(X509Certificate serverCertificate) {
        this.serverCertificate = serverCertificate;
    }

    public ReadBuffer encodeMessage(MessagePDU pdu, byte[] message) {
        switch (securitypolicy) {
            case NONE:
                return new ReadBufferByteBased(message, LITTLE_ENDIAN);
            case Basic256Sha256:
            case Basic128Rsa15:
                if (pdu instanceof OpcuaOpenRequest) {
                    return asymmetricEncryptionHandler.encodeMessage(pdu, message);
                } else {
                    return symmetricEncryptionHandler.encodeMessage(pdu, message, clientNonce, serverNonce);
                }
            default:
                throw new IllegalStateException("Driver doesn't support security policy: " + securitypolicy);
        }
    }


    public SignatureData createClientSignature(byte[] lastServerNonce) {
        byte[] cert = Try.of(() -> serverCertificate.getEncoded()).getOrElse(new byte[0]);
        byte[] bytes = ByteBuffer.allocate(cert.length+lastServerNonce.length).put(cert).put(lastServerNonce).array();
        byte[] signed = asymmetricEncryptionHandler.sign(bytes);
        return new SignatureData(new PascalString(securitypolicy.getAsymmetricSignatureAlgorithm().getUri()), new PascalByteString(signed.length, signed));
    }

    public OpcuaAPU decodeMessage(OpcuaAPU pdu) {
        LOGGER.info("Decoding Message with Security policy {}", securitypolicy);

        switch (securitypolicy) {
            case NONE:
                return pdu;
            case Basic128Rsa15:
            case Basic256Sha256:
                if (pdu.getMessage() instanceof OpcuaOpenResponse) {
                    return asymmetricEncryptionHandler.decodeMessage(pdu);
                } else {
                    return symmetricEncryptionHandler.decodeMessage(pdu, clientNonce, serverNonce);
                }
            default:
                throw new IllegalStateException("Driver doesn't support security policy: " + securitypolicy);
        }
    }


    public byte[] encryptPassword(byte[] data) {
        try {
            Cipher cipher = Cipher.getInstance("RSA/ECB/OAEPWithSHA-1AndMGF1Padding");
            cipher.init(Cipher.ENCRYPT_MODE, this.serverCertificate.getPublicKey());
            return cipher.doFinal(data);
        } catch (Exception e) {
            LOGGER.error("Unable to encrypt Data", e);
            return null;
        }
    }

    public static X509Certificate getCertificateX509(byte[] senderCertificate) {
        try {
            CertificateFactory factory = CertificateFactory.getInstance("X.509");
            LOGGER.info("Public Key Length {}", senderCertificate.length);
            return (X509Certificate) factory.generateCertificate(new ByteArrayInputStream(senderCertificate));
        } catch (Exception e) {
            LOGGER.error("Unable to get certificate from String {}", senderCertificate);
            return null;
        }
    }

    public void setClientNonce(byte[] clientNonce) {
        this.clientNonce = clientNonce;
    }


    public void setServerNonce(byte[] serverNonce) {
        this.serverNonce = serverNonce;
    }
}
