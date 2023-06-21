package org.apache.plc4x.java.opcua.context;

import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaOpenResponse;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.generation.*;

import javax.crypto.Cipher;
import java.nio.ByteBuffer;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Signature;
import java.security.cert.Certificate;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;

public class AsymmetricEncryptionHandler {
    private static final int SECURE_MESSAGE_HEADER_SIZE = 12;
    private static final int SEQUENCE_HEADER_SIZE = 8;

    private final SecurityPolicy policy;

    private final X509Certificate serverCertificate;
    private final X509Certificate clientCertificate;

    private final PrivateKey clientPrivateKey;

    public AsymmetricEncryptionHandler(X509Certificate serverCertificate, X509Certificate clientCertificate, PrivateKey clientPrivateKey, PublicKey clientPublicKey, SecurityPolicy policy) {
        this.serverCertificate = serverCertificate;
        this.clientCertificate = clientCertificate;
        this.clientPrivateKey = clientPrivateKey;
        this.policy = policy;
    }

    /**
     * Docs: https://reference.opcfoundation.org/Core/Part6/v104/docs/6.7
     *
     * @param pdu
     * @param message
     * @return
     */
    public ReadBuffer encodeMessage(MessagePDU pdu, byte[] message) {
        int unencryptedLength = pdu.getLengthInBytes();
        int messageLength = message.length;

        int beforeBodyLength = unencryptedLength - messageLength; // message header, security header, sequence header

        int cipherTextBlockSize = (getAsymmetricKeyLength(serverCertificate) + 7) / 8;
        int plainTextBlockSize = (getAsymmetricKeyLength(serverCertificate) + 7) / 8 - policy.getAsymmetricPlainBlock();
        int signatureSize = (getAsymmetricKeyLength(clientCertificate) + 7) / 8;


        int maxChunkSize = 8196;
        int paddingOverhead = cipherTextBlockSize > 256 ? 2 : 1;


        int securityHeaderSize = beforeBodyLength - SEQUENCE_HEADER_SIZE - SECURE_MESSAGE_HEADER_SIZE;
        int maxCipherTextSize = maxChunkSize - securityHeaderSize;
        int maxCipherTextBlocks = maxCipherTextSize / cipherTextBlockSize;
        int maxPlainTextSize = maxCipherTextBlocks * plainTextBlockSize;
        int maxBodySize = maxPlainTextSize - SEQUENCE_HEADER_SIZE - paddingOverhead - signatureSize;

        int bodySize = Math.min(message.length, maxBodySize);

        int plainTextSize = SEQUENCE_HEADER_SIZE + bodySize + paddingOverhead + signatureSize;
        int remaining = plainTextSize % plainTextBlockSize;
        int paddingSize = remaining > 0 ? plainTextBlockSize - remaining : 0;

        int plainTextContentSize = SEQUENCE_HEADER_SIZE + bodySize +
            signatureSize + paddingSize + paddingOverhead;

        int frameSize = SECURE_MESSAGE_HEADER_SIZE + securityHeaderSize +
            (plainTextContentSize / plainTextBlockSize) * cipherTextBlockSize;

        try {
            WriteBufferByteBased buf = new WriteBufferByteBased(frameSize, ByteOrder.LITTLE_ENDIAN);
            OpcuaAPU opcuaAPU = new OpcuaAPU(pdu);
            opcuaAPU.serialize(buf);

            writePadding(paddingSize, buf);
            updateFrameSize(frameSize, buf);

            byte[] sign = sign(buf.getBytes());
            buf.writeByteArray(sign);

            buf.setPos(SECURE_MESSAGE_HEADER_SIZE + securityHeaderSize);

            int blockCount = (frameSize - buf.getPos()) / plainTextBlockSize;// -> plainTextContentSize / plainTextBlockSize

            byte[] encrypted = encrypt(plainTextBlockSize, securityHeaderSize, frameSize, buf, blockCount);
            buf.writeByteArray(encrypted);

            return new ReadBufferByteBased(buf.getBytes(), ByteOrder.LITTLE_ENDIAN);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public OpcuaAPU decodeMessage(OpcuaAPU pdu) {
        MessagePDU message = pdu.getMessage();

        OpcuaOpenResponse a = (OpcuaOpenResponse) message;


        int cipherTextBlockSize = (getAsymmetricKeyLength(serverCertificate) + 7) / 8;
        int signatureSize = (getAsymmetricKeyLength(clientCertificate) + 7) / 8;

        byte[] textMessage = a.getMessage();


        int blockCount = (SEQUENCE_HEADER_SIZE + textMessage.length) / cipherTextBlockSize;
        int plainTextBufferSize = cipherTextBlockSize * blockCount;


        try {
            WriteBufferByteBased buf = new WriteBufferByteBased(pdu.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            pdu.serialize(buf);

            Cipher cipher = policy.getAsymmetricEncryptionAlgorithm().getCipher();
            cipher.init(Cipher.DECRYPT_MODE, clientPrivateKey);

            ByteBuffer buffer = ByteBuffer.allocate(plainTextBufferSize);
            byte[] bytes = buf.getBytes(pdu.getLengthInBytes() - plainTextBufferSize, pdu.getLengthInBytes());
            ByteBuffer originalMessage = ByteBuffer.wrap(bytes);

            for (int blockNumber = 0; blockNumber < blockCount; blockNumber++) {
                originalMessage.limit(originalMessage.position() + cipherTextBlockSize);
                cipher.doFinal(originalMessage, buffer);
            }
            buffer.flip();

            buf.setPos(pdu.getLengthInBytes() - plainTextBufferSize);
            buf.writeByteArray(buffer.array());


            int frameSize = pdu.getLengthInBytes() - plainTextBufferSize + buffer.limit();

            updateFrameSize(frameSize, buf);

            byte[] decryptedMessage = buf.getBytes(0, frameSize);
            ReadBuffer readBuffer = new ReadBufferByteBased(decryptedMessage, ByteOrder.LITTLE_ENDIAN);
            OpcuaAPU opcuaAPU = OpcuaAPU.staticParse(readBuffer, true);
            return opcuaAPU;
        } catch (Exception e) {
            throw new RuntimeException(e);
        }


    }

    private byte[] encrypt(int plainTextBlockSize, int securityHeaderSize, int frameSize, WriteBufferByteBased buf, int blockCount) throws Exception {
        ByteBuffer buffer = ByteBuffer.allocate(frameSize - buf.getPos());
        ByteBuffer originalMessage = ByteBuffer.wrap(buf.getBytes(SECURE_MESSAGE_HEADER_SIZE + securityHeaderSize, frameSize));


        Cipher cipher = policy.getAsymmetricEncryptionAlgorithm().getCipher();
        cipher.init(Cipher.ENCRYPT_MODE, serverCertificate.getPublicKey());

        for (int block = 0; block < blockCount; block++) {
            int position = block * plainTextBlockSize;
            int limit = (block + 1) * plainTextBlockSize;
            originalMessage.position(position);
            originalMessage.limit(limit);

            cipher.doFinal(originalMessage, buffer);

        }
        return buffer.array();
    }

    private static void updateFrameSize(int frameSize, WriteBufferByteBased buf) throws SerializationException {
        int initPosition = buf.getPos();
        buf.setPos(4);
        buf.writeInt(32, frameSize);
        buf.setPos(initPosition);
    }


    public byte[] sign(byte[] data) {
        try {
            Signature signature = policy.getAsymmetricSignatureAlgorithm().getSignature();
            signature.initSign(clientPrivateKey);
            signature.update(data);
            return signature.sign();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private void writePadding(int paddingSize, WriteBufferByteBased buffer) throws Exception {
        buffer.writeByte((byte) paddingSize);
        for (int i = 0; i < paddingSize; i++) {
            buffer.writeByte((byte) paddingSize);
        }
    }


    static int getAsymmetricKeyLength(Certificate certificate) {
        PublicKey publicKey = certificate != null ?
            certificate.getPublicKey() : null;

        return (publicKey instanceof RSAPublicKey) ?
            ((RSAPublicKey) publicKey).getModulus().bitLength() : 0;
    }

}
