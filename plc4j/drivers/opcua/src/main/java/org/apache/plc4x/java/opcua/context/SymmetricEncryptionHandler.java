package org.apache.plc4x.java.opcua.context;

import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaAPU;
import org.apache.plc4x.java.opcua.readwrite.OpcuaMessageResponse;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.EncryptionAlgorithm;
import org.apache.plc4x.java.opcua.security.SecurityPolicy.MacSignatureAlgorithm;
import org.apache.plc4x.java.opcua.security.SymmetricKeys;
import org.apache.plc4x.java.spi.generation.*;

import javax.crypto.Cipher;
import javax.crypto.Mac;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.nio.ByteBuffer;
import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;


public class SymmetricEncryptionHandler {
    private static final int SECURE_MESSAGE_HEADER_SIZE = 12;
    private static final int SEQUENCE_HEADER_SIZE = 8;
    private static final int SYMMETRIC_SECURITY_HEADER_SIZE = 4;

    private final SecurityPolicy policy;


    private SymmetricKeys keys = null;

    public SymmetricEncryptionHandler(SecurityPolicy policy) {
        this.policy = policy;
    }

    /**
     * Docs: https://reference.opcfoundation.org/Core/Part6/v104/docs/6.7
     *
     * @param pdu
     * @param message
     * @param clientNonce
     * @param serverNonce
     * @return
     */
    public ReadBuffer encodeMessage(MessagePDU pdu, byte[] message, byte[] clientNonce, byte[] serverNonce) {
        int unencryptedLength = pdu.getLengthInBytes();
        int messageLength = message.length;

        int beforeBodyLength = unencryptedLength - messageLength; // message header, security header, sequence header

        int cipherTextBlockSize = 16;
        int plainTextBlockSize = 16;
        int signatureSize = policy.getSymmetricSignatureAlgorithm().getSymmetricSignatureSize();


        int maxChunkSize = 8196;
        int paddingOverhead = 1;


        int securityHeaderSize = SYMMETRIC_SECURITY_HEADER_SIZE;
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

        SymmetricKeys symmetricKeys = getSymmetricKeys(clientNonce, serverNonce);

        try {
            WriteBufferByteBased buf = new WriteBufferByteBased(frameSize, ByteOrder.LITTLE_ENDIAN);
            OpcuaAPU opcuaAPU = new OpcuaAPU(pdu);
            opcuaAPU.serialize(buf);

            writePadding(paddingSize, buf);
            updateFrameSize(frameSize, buf);

            byte[] sign = sign(buf.getBytes(), symmetricKeys.getClientKeys());
            buf.writeByteArray(sign);

            buf.setPos(SECURE_MESSAGE_HEADER_SIZE + securityHeaderSize);

            byte[] encrypted = encrypt(securityHeaderSize, frameSize, buf, symmetricKeys);
            buf.writeByteArray(encrypted);

            return new ReadBufferByteBased(buf.getBytes(), ByteOrder.LITTLE_ENDIAN);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public OpcuaAPU decodeMessage(OpcuaAPU pdu, byte[] clientNonce, byte[] serverNonce) {
        MessagePDU message = pdu.getMessage();

        OpcuaMessageResponse a = (OpcuaMessageResponse) message;


        int cipherTextBlockSize = 16; // different for aes256

        byte[] textMessage = a.getMessage();


        int blockCount = (SEQUENCE_HEADER_SIZE + textMessage.length) / cipherTextBlockSize;
        int plainTextBufferSize = cipherTextBlockSize * blockCount;


        try {
            WriteBufferByteBased buf = new WriteBufferByteBased(pdu.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            pdu.serialize(buf);

            EncryptionAlgorithm transformation = policy.getSymmetricEncryptionAlgorithm();
            SymmetricKeys symmetricKeys = getSymmetricKeys(clientNonce, serverNonce);
            Cipher cipher = getCipher(symmetricKeys.getServerKeys(), transformation, Cipher.DECRYPT_MODE);

            ByteBuffer buffer = ByteBuffer.allocate(plainTextBufferSize);
            byte[] bytes = buf.getBytes(pdu.getLengthInBytes() - plainTextBufferSize, pdu.getLengthInBytes());
            ByteBuffer originalMessage = ByteBuffer.wrap(bytes);


            cipher.doFinal(originalMessage, buffer);

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

    private byte[] encrypt(int securityHeaderSize, int frameSize, WriteBufferByteBased buf, SymmetricKeys symmetricKeys) throws Exception {
        ByteBuffer buffer = ByteBuffer.allocate(frameSize - buf.getPos());
        ByteBuffer originalMessage = ByteBuffer.wrap(buf.getBytes(SECURE_MESSAGE_HEADER_SIZE + securityHeaderSize, frameSize));


        EncryptionAlgorithm transformation = policy.getSymmetricEncryptionAlgorithm();
        Cipher cipher = getCipher(symmetricKeys.getClientKeys(), transformation, Cipher.ENCRYPT_MODE);


        cipher.doFinal(originalMessage, buffer);

        return buffer.array();
    }

    private static Cipher getCipher(SymmetricKeys.Keys symmetricKeys, EncryptionAlgorithm transformation, int mode) throws NoSuchAlgorithmException, NoSuchPaddingException, InvalidKeyException, InvalidAlgorithmParameterException {
        Cipher cipher = transformation.getCipher();

        SecretKeySpec keySpec = new SecretKeySpec(symmetricKeys.getEncryptionKey(), "AES");
        IvParameterSpec ivSpec = new IvParameterSpec(symmetricKeys.getInitializationVector());

        cipher.init(mode, keySpec, ivSpec);
        return cipher;
    }

    private static void updateFrameSize(int frameSize, WriteBufferByteBased buf) throws SerializationException {
        int initPosition = buf.getPos();
        buf.setPos(4);
        buf.writeInt(32, frameSize);
        buf.setPos(initPosition);
    }


    public byte[] sign(byte[] data, SymmetricKeys.Keys symmetricKeys) {
        try {
            MacSignatureAlgorithm algorithm = policy.getSymmetricSignatureAlgorithm();
            Mac signature = algorithm.getSignature();
            signature.init(new SecretKeySpec(symmetricKeys.getSignatureKey(), algorithm.getName()));
            signature.update(data);
            return signature.doFinal();
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

    private SymmetricKeys getSymmetricKeys(byte[] clientNonce, byte[] serverNonce) {
        if (keys == null) {
            keys = SymmetricKeys.generateKeyPair(clientNonce, serverNonce, policy.getSymmetricSignatureAlgorithm());
        }
        return keys;
    }
}
