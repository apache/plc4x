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
package org.apache.plc4x.java.opcua.context;

import java.nio.ByteBuffer;
import java.security.GeneralSecurityException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.function.Supplier;
import org.apache.plc4x.java.opcua.protocol.chunk.Chunk;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.apache.plc4x.java.opcua.readwrite.ChunkType;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;

/**
 * Common (de)chunking + (de)signing/encryption logic shared by
 * {@link AsymmetricEncryptionHandler} (used for the {@code OPN} chunk) and
 * {@link SymmetricEncryptionHandler} (used for all subsequent {@code MSG}/{@code CLO}
 * chunks).
 *
 * <p>Operates on raw {@code byte[]} working arrays. The new SPI's WriteBuffer
 * doesn't expose absolute {@code setPos}/{@code getPos}/{@code getBytes(start,end)}
 * any more, so the chunking code does its surgery on a plain array and only
 * uses {@link org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased}
 * for the very last "serialize this MessagePDU" step (which the
 * {@link PayloadConverter} hides).</p>
 */
abstract class BaseEncryptionHandler {

    protected static final int SECURE_MESSAGE_HEADER_SIZE = 12;
    protected static final int SEQUENCE_HEADER_SIZE = 8;

    protected final SecureChannelState conversation;
    protected final SecurityPolicy securityPolicy;

    public BaseEncryptionHandler(SecureChannelState conversation, SecurityPolicy securityPolicy) {
        this.conversation = conversation;
        this.securityPolicy = securityPolicy;
    }

    /**
     * Splits {@code message} into one or more wire chunks, applying signing and
     * encryption per the policy held by {@code chunk}. The returned list always
     * has the last entry marked {@link ChunkType#FINAL}.
     */
    public final List<MessagePDU> encodeMessage(Chunk chunk, MessagePDU message, Supplier<Integer> sequenceSupplier) {

        try {
            ByteBuffer messageBuffer = ByteBuffer.wrap(PayloadConverter.toStream(message));
            int sequenceStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();

            // Pull the three header blocks once; they get re-emitted on every chunk
            // (with the chunk type byte and sequenceId rewritten per-chunk).
            byte[] messageHeader = new byte[SECURE_MESSAGE_HEADER_SIZE];
            messageBuffer.get(messageHeader);
            byte[] securityHeader = new byte[chunk.getSecurityHeaderSize()];
            messageBuffer.get(securityHeader);
            byte[] sequenceHeader = new byte[SEQUENCE_HEADER_SIZE];
            messageBuffer.get(sequenceHeader);

            ByteBuffer bodyBuffer = messageBuffer.slice();
            List<MessagePDU> messages = new ArrayList<>();
            boolean first = true;
            while (bodyBuffer.hasRemaining()) {
                int bodySize = Math.min(bodyBuffer.remaining(), chunk.getMaxBodySize());
                int paddingSize = 0;
                if (chunk.isEncrypted()) {
                    int plainTextSize = SEQUENCE_HEADER_SIZE + bodySize + chunk.getPaddingOverhead() + chunk.getSignatureSize();
                    int gap = plainTextSize % chunk.getPlainTextBlockSize();
                    paddingSize = gap > 0 ? chunk.getPlainTextBlockSize() - gap : 0;
                }

                int plainTextContentSize = SEQUENCE_HEADER_SIZE + bodySize + chunk.getSignatureSize() + paddingSize + chunk.getPaddingOverhead();
                if (chunk.isEncrypted()) {
                    assert ((plainTextContentSize % chunk.getPlainTextBlockSize()) == 0);
                }

                int chunkSize = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize() + (plainTextContentSize / chunk.getPlainTextBlockSize()) * chunk.getCipherTextBlockSize();

                // Build the chunk in a plain byte array so we can address any
                // offset directly (writing chunk type / total-length / sequence
                // id after the fact).
                byte[] chunkBytes = new byte[chunkSize];
                int pos = 0;
                System.arraycopy(messageHeader, 0, chunkBytes, pos, messageHeader.length);
                pos += messageHeader.length;
                System.arraycopy(securityHeader, 0, chunkBytes, pos, securityHeader.length);
                pos += securityHeader.length;
                System.arraycopy(sequenceHeader, 0, chunkBytes, pos, sequenceHeader.length);
                pos += sequenceHeader.length;

                // Overwrite the total-length field (offsets 4..7, LE).
                writeUint32LE(chunkBytes, 4, chunkSize);

                ChunkType chunkType = bodyBuffer.remaining() - bodySize > 0 ? ChunkType.CONTINUE : ChunkType.FINAL;
                // Chunk type lives at byte 3.
                if (chunkType != ChunkType.FINAL) {
                    chunkBytes[3] = (byte) chunkType.getValue().charAt(0);
                }
                // Subsequent chunks get a fresh sequence id (sequence header lives
                // at the start of the body section).
                if (!first) {
                    writeUint32LE(chunkBytes, sequenceStart, sequenceSupplier.get());
                }
                first = false;

                // Copy the next slice of the body.
                bodyBuffer.get(chunkBytes, pos, bodySize);
                pos += bodySize;

                if (chunk.isEncrypted()) {
                    // Plain-text-padding byte (and optional MSB padding-overhead
                    // byte when the cipher block is bigger than 256 bytes).
                    int totalPaddingBytes = paddingSize + chunk.getPaddingOverhead();
                    for (int index = 0; index < totalPaddingBytes; index++) {
                        chunkBytes[pos + index] = (byte) paddingSize;
                    }
                    if (chunk.getPaddingOverhead() > 1) {
                        // Last byte of the padding region is the MSB of the size.
                        chunkBytes[pos + totalPaddingBytes - 1] = (byte) ((paddingSize >> 8) & 0xFF);
                    }
                    pos += totalPaddingBytes;
                }

                if (chunk.isSigned()) {
                    byte[] signatureData = sign(Arrays.copyOfRange(chunkBytes, 0, pos));
                    System.arraycopy(signatureData, 0, chunkBytes, pos, signatureData.length);
                    pos += signatureData.length;
                }
                if (chunk.isEncrypted()) {
                    encryptInPlace(chunkBytes, chunk.getSecurityHeaderSize(), chunk.getPlainTextBlockSize(),
                        chunk.getCipherTextBlockSize(), plainTextContentSize / chunk.getPlainTextBlockSize()
                    );
                }

                MessagePDU chunkedMessage = PayloadConverter.pduFromStream(chunkBytes, message.getResponse());
                messages.add(chunkedMessage);
            }
            return messages;
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    /**
     * Decrypts + verifies a single received chunk and returns it with the
     * security/padding/signature stripped off. The returned MessagePDU's
     * total-length header field has been rewritten to match the new size.
     */
    public final MessagePDU decodeMessage(Chunk chunk, MessagePDU message) {
        try {
            if (!chunk.isEncrypted() && !chunk.isSigned()) {
                return message;
            }

            int messageLength = message.getLengthInBytes();
            byte[] chunkBytes = PayloadConverter.toStream(message);

            int bodySize = messageLength - chunk.getSecurityHeaderSize() - SECURE_MESSAGE_HEADER_SIZE;
            if (chunk.isEncrypted()) {
                bodySize = decryptInPlace(chunkBytes, chunk, messageLength);
            }

            if (chunk.isSigned()) {
                verify(chunkBytes, chunk, messageLength);
            }

            int encryptionOverhead = getEncryptionOverhead(chunk, messageLength);
            int paddingSize = getPaddingSize(chunkBytes, chunk, messageLength);

            int payloadStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();
            int payloadEnd = payloadStart + bodySize - paddingSize - chunk.getSignatureSize() - chunk.getPaddingOverhead();
            int expectedPaddingSize = messageLength - payloadEnd - chunk.getSignatureSize() - encryptionOverhead - chunk.getPaddingOverhead();

            if (paddingSize != expectedPaddingSize) {
                throw new IllegalArgumentException("Malformed data detected - expected padding size do not match");
            }

            if (chunk.isEncrypted()) {
                byte paddingByte = (byte) (paddingSize & 0xff);
                for (int index = 0; index < expectedPaddingSize; index++) {
                    if (chunkBytes[payloadEnd + index] != paddingByte) {
                        throw new IllegalArgumentException("Malformed padding byte at index " + index);
                    }
                }
            }

            int overhead = paddingSize + chunk.getSignatureSize() + chunk.getPaddingOverhead() + encryptionOverhead;
            writeUint32LE(chunkBytes, 4, messageLength - overhead);

            return PayloadConverter.pduFromStream(chunkBytes, message.getResponse());
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    /** Write a little-endian uint32 at the given byte offset. */
    private static void writeUint32LE(byte[] buf, int offset, long value) {
        buf[offset]     = (byte) (value & 0xFF);
        buf[offset + 1] = (byte) ((value >> 8) & 0xFF);
        buf[offset + 2] = (byte) ((value >> 16) & 0xFF);
        buf[offset + 3] = (byte) ((value >> 24) & 0xFF);
    }

    private int getEncryptionOverhead(Chunk chunk, int messageLength) {
        if (!chunk.isEncrypted()) {
            return 0;
        }

        int bodyStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();
        int bodySize = messageLength - bodyStart;
        int blockCount = bodySize / chunk.getCipherTextBlockSize();
        // bytes we "lost" after payload got decrypted
        return (chunk.getCipherTextBlockSize() * blockCount) - (chunk.getPlainTextBlockSize() * blockCount);
    }

    private short getPaddingSize(byte[] chunkBytes, Chunk chunk, int messageLength) {
        if (!chunk.isEncrypted()) {
            return 0;
        }

        int bodyStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();
        int bodySize = messageLength - bodyStart;
        int blockCount = bodySize / chunk.getCipherTextBlockSize();
        int encryptionOverhead = (chunk.getCipherTextBlockSize() * blockCount) - (chunk.getPlainTextBlockSize() * blockCount);

        int paddingEnd = messageLength - chunk.getSignatureSize() - encryptionOverhead - chunk.getPaddingOverhead();
        // The chunk's padding-overhead is 1 byte for cipher block <= 256, 2 bytes otherwise
        // (low byte at paddingEnd, high byte at paddingEnd+1).
        if (chunk.getPaddingOverhead() > 1) {
            int paddingSize = ((chunkBytes[paddingEnd + 1] & 0xFF) << 8) | (chunkBytes[paddingEnd] & 0xFF);
            return (short) (paddingSize & 0xFFFF);
        }
        return (short) (chunkBytes[paddingEnd] & 0xFF);
    }

    /** Subclass crypto: in-place verify. {@code chunkBytes} is the full received frame. */
    protected abstract void verify(byte[] chunkBytes, Chunk chunk, int messageLength) throws Exception;

    /**
     * Subclass crypto: in-place decrypt. Returns the number of bytes of body
     * material that decrypted to (i.e. the post-decrypt body length).
     */
    protected abstract int decryptInPlace(byte[] chunkBytes, Chunk chunk, int messageLength) throws Exception;

    /**
     * Subclass crypto: in-place encrypt. Encrypts {@code blockCount} blocks of
     * the plaintext region starting immediately after the security header.
     */
    protected abstract void encryptInPlace(byte[] chunkBytes, int securityHeaderSize, int plainTextBlockSize,
                                           int cipherTextBlockSize, int blockCount) throws Exception;

    protected abstract byte[] sign(byte[] contentsToSign) throws GeneralSecurityException;

}
