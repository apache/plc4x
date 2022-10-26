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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.opcua.protocol.OpcuaProtocolLogic;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaMessageResponse;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenResponse;
import org.apache.plc4x.java.spi.generation.*;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.io.ByteArrayInputStream;
import java.security.*;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;


public class EncryptionHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaProtocolLogic.class);

    static {
        // Required for SecurityPolicy.Aes256_Sha256_RsaPss
        Security.addProvider(new BouncyCastleProvider());
    }

    private X509Certificate serverCertificate;
    private X509Certificate clientCertificate;
    private PrivateKey clientPrivateKey;
    private PublicKey clientPublicKey;
    private String securitypolicy;

    public EncryptionHandler(CertificateKeyPair ckp, byte[] senderCertificate, String securityPolicy) {
        if (ckp != null) {
            this.clientPrivateKey = ckp.getKeyPair().getPrivate();
            this.clientPublicKey = ckp.getKeyPair().getPublic();
            this.clientCertificate = ckp.getCertificate();
        }
        if (senderCertificate != null) {
            this.serverCertificate = getCertificateX509(senderCertificate);
        }
        this.securitypolicy = securityPolicy;
    }

    public void setServerCertificate(X509Certificate serverCertificate) {
        this.serverCertificate = serverCertificate;
    }

    public ReadBuffer encodeMessage(MessagePDU pdu, byte[] message) {
        int PREENCRYPTED_BLOCK_LENGTH = 190;
        int unencryptedLength = pdu.getLengthInBytes();
        int openRequestLength = message.length;
        int positionFirstBlock = unencryptedLength - openRequestLength - 8;
        int paddingSize = PREENCRYPTED_BLOCK_LENGTH - ((openRequestLength + 256 + 1 + 8) % PREENCRYPTED_BLOCK_LENGTH);
        int preEncryptedLength = openRequestLength + 256 + 1 + 8 + paddingSize;
        if (preEncryptedLength % PREENCRYPTED_BLOCK_LENGTH != 0) {
            throw new PlcRuntimeException("Pre encrypted block length " + preEncryptedLength + " isn't a multiple of the block size");
        }
        int numberOfBlocks = preEncryptedLength / PREENCRYPTED_BLOCK_LENGTH;
        int encryptedLength = numberOfBlocks * 256 + positionFirstBlock;
        WriteBufferByteBased buf = new WriteBufferByteBased(encryptedLength, ByteOrder.LITTLE_ENDIAN);
        try {
            new OpcuaAPU(pdu, false).serialize(buf);
            byte paddingByte = (byte) paddingSize;
            buf.writeByte(paddingByte);
            for (int i = 0; i < paddingSize; i++) {
                buf.writeByte(paddingByte);
            }
            //Writing Message Length
            int tempPos = buf.getPos();
            buf.setPos(4);
            buf.writeInt(32, encryptedLength);
            buf.setPos(tempPos);
            byte[] signature = sign(getBytes(buf.getBytes(), 0, unencryptedLength + paddingSize + 1));
            //Write the signature to the end of the buffer
            for (byte b : signature) {
                buf.writeByte(b);
            }
            buf.setPos(positionFirstBlock);
            encryptBlock(buf, getBytes(buf.getBytes(), positionFirstBlock, positionFirstBlock + preEncryptedLength));
            return new ReadBufferByteBased(buf.getData(), ByteOrder.LITTLE_ENDIAN);
        } catch (SerializationException e) {
            throw new PlcRuntimeException("Unable to parse apu prior to encrypting");
        }
    }

    public OpcuaAPU decodeMessage(OpcuaAPU pdu) {
        LOGGER.info("Decoding Message with Security policy {}", securitypolicy);
        switch (securitypolicy) {
            case "None":
                return pdu;
            case "Basic256Sha256":
                byte[] message;
                if (pdu.getMessage() instanceof OpcuaOpenResponse) {
                    message = ((OpcuaOpenResponse) pdu.getMessage()).getMessage();
                } else if (pdu.getMessage() instanceof OpcuaMessageResponse) {
                    message = ((OpcuaMessageResponse) pdu.getMessage()).getMessage();
                } else {
                    return pdu;
                }
                try {
                    int encryptedLength = pdu.getLengthInBytes();
                    int encryptedMessageLength = message.length + 8;
                    int headerLength = encryptedLength - encryptedMessageLength;
                    int numberOfBlocks = encryptedMessageLength / 256;
                    WriteBufferByteBased buf = new WriteBufferByteBased(headerLength + numberOfBlocks * 256, ByteOrder.LITTLE_ENDIAN);
                    pdu.serialize(buf);
                    byte[] data = getBytes(buf.getBytes(), headerLength, encryptedLength);
                    buf.setPos(headerLength);
                    decryptBlock(buf, data);
                    int tempPos = buf.getPos();
                    buf.setPos(0);
                    if (!checkSignature(getBytes(buf.getBytes(), 0, tempPos))) {
                        LOGGER.info("Signature verification failed: - {}", getBytes(buf.getBytes(), 0, tempPos - 256));
                    }
                    buf.setPos(4);
                    buf.writeInt(32, tempPos - 256);
                    ReadBuffer readBuffer = new ReadBufferByteBased(getBytes(buf.getBytes(), 0, tempPos - 256), ByteOrder.LITTLE_ENDIAN);
                    return OpcuaAPU.staticParse(readBuffer, true);
                } catch (SerializationException | ParseException e) {
                    LOGGER.error("Unable to Parse encrypted message");
                }
        }
        return pdu;
    }

    public void decryptBlock(WriteBuffer buf, byte[] data) {
        try {
            Cipher cipher = Cipher.getInstance("RSA/ECB/OAEPWithSHA-1AndMGF1Padding");
            cipher.init(Cipher.DECRYPT_MODE, this.clientPrivateKey);

            for (int i = 0; i < data.length; i += 256) {
                byte[] decrypted = cipher.doFinal(data, i, 256);
                for (int j = 0; j < 214; j++) {
                    buf.writeByte(decrypted[j]);
                }
            }
        } catch (Exception e) {
            LOGGER.error("Unable to decrypt Data", e);
        }
    }

    public boolean checkSignature(byte[] data) {
        try {
            Signature signature = Signature.getInstance("SHA256withRSA", "BC");
            signature.initVerify(serverCertificate.getPublicKey());
            signature.update(data);
            return signature.verify(data, 0, data.length - 256);
        } catch (Exception e) {
            e.printStackTrace();
            LOGGER.error("Unable to sign Data");
            return false;
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

    public void encryptBlock(WriteBuffer buf, byte[] data) {
        try {
            Cipher cipher = Cipher.getInstance("RSA/ECB/OAEPWithSHA-1AndMGF1Padding");
            cipher.init(Cipher.ENCRYPT_MODE, this.serverCertificate.getPublicKey());
            for (int i = 0; i < data.length; i += 190) {
                LOGGER.info("Iterate:- {}, Data Length:- {}", i, data.length);
                byte[] encrypted = cipher.doFinal(data, i, 190);
                for (int j = 0; j < 256; j++) {
                    buf.writeByte(encrypted[j]);
                }
            }
        } catch (Exception e) {
            LOGGER.error("Unable to encrypt Data");
            e.printStackTrace();
        }
    }

    public void encryptHmacBlock(WriteBuffer buf, byte[] data) {
        try {
            Mac cipher = Mac.getInstance("HmacSHA256");
            SecretKeySpec keySpec = new SecretKeySpec(getSecretKey(), "HmacSHA256");
            cipher.init(keySpec);
        } catch (Exception e) {
            LOGGER.error("Unable to encrypt Data", e);
        }
    }

    public byte[] getSecretKey() {
        return null;
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

    public byte[] sign(byte[] data) {
        try {
            Signature signature = Signature.getInstance("SHA256withRSA", "BC");
            signature.initSign(this.clientPrivateKey);
            signature.update(data);
            byte[] ss = signature.sign();
            LOGGER.info("----------------Signature Length{}", ss.length);
            return ss;
        } catch (Exception e) {
            e.printStackTrace();
            LOGGER.error("Unable to sign Data");
            return null;
        }
    }

    private byte[] getBytes(byte[] bytes, int startPos, int endPos) {
        int numBytes = endPos - startPos;
        byte[] data = new byte[numBytes];
        System.arraycopy(bytes, startPos, data, 0, numBytes);
        return data;
    }
}
