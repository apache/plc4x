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
import java.util.List;
import java.util.function.Supplier;
import org.apache.plc4x.java.opcua.protocol.chunk.Chunk;
import org.apache.plc4x.java.opcua.protocol.chunk.PayloadConverter;
import org.apache.plc4x.java.opcua.readwrite.ChunkType;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

abstract class BaseEncryptionHandler {

    protected static final int SECURE_MESSAGE_HEADER_SIZE = 12;
    protected static final int SEQUENCE_HEADER_SIZE = 8;

    protected final Conversation conversation;
    protected final SecurityPolicy securityPolicy;

    public BaseEncryptionHandler(Conversation conversation, SecurityPolicy securityPolicy) {
        this.conversation = conversation;
        this.securityPolicy = securityPolicy;
    }
    public final List<MessagePDU> encodeMessage(Chunk chunk, MessagePDU message, Supplier<Integer> sequenceSupplier) {

        try {
            ByteBuffer messageBuffer = ByteBuffer.wrap(PayloadConverter.toStream(message));
            int sequenceStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();

            // processed parts of frame
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

                WriteBufferByteBased chunkBuffer = new WriteBufferByteBased(chunkSize, ByteOrder.LITTLE_ENDIAN);
                chunkBuffer.writeByteArray("messageHeader", messageHeader);
                chunkBuffer.writeByteArray("securityHeader", securityHeader);
                chunkBuffer.writeByteArray("sequenceHeader", sequenceHeader);
                updateFrameSize(chunkBuffer, chunkSize);
                ChunkType chunkType = bodyBuffer.remaining() - bodySize > 0 ? ChunkType.CONTINUE : ChunkType.FINAL;
                updateFrame(first, chunkBuffer, chunk, chunkType, sequenceSupplier); // populate headers
                first = false;

                byte[] chunkContents = new byte[bodySize];
                bodyBuffer.get(chunkContents);
                // copy part of message not larger than body size into chunk buffer
                chunkBuffer.writeByteArray("payload", chunkContents);

                if (chunk.isEncrypted()) {
                    for (int index = 0, limit = paddingSize + chunk.getPaddingOverhead(); index < limit; index++) {
                        chunkBuffer.writeByte("padding", (byte) paddingSize);
                    }
                    if (chunk.getPaddingOverhead() > 1) {
                        // override extra padding byte with MSB of padding size
                        chunkBuffer.setPos(bodySize + paddingSize + chunk.getPaddingOverhead());
                        chunkBuffer.writeByte("paddingMSB", (byte) ((paddingSize >> 8) & 0xFF));
                    }
                }

                if (chunk.isSigned()) {
                    byte[] signatureData = sign(chunkBuffer.getBytes(0, chunkBuffer.getPos()));
                    chunkBuffer.writeByteArray("signature", signatureData);
                }
                if (chunk.isEncrypted()) {
                    encrypt(chunkBuffer, chunk.getSecurityHeaderSize(), chunk.getPlainTextBlockSize(),
                        chunk.getCipherTextBlockSize(), plainTextContentSize / chunk.getPlainTextBlockSize()
                    );
                }

                MessagePDU chunkedMessage = PayloadConverter.pduFromStream(chunkBuffer.getBytes(), message.getResponse());
                messages.add(chunkedMessage);
            }
            return messages;
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public final MessagePDU decodeMessage(Chunk chunk, MessagePDU message) {
        try {
            if (!chunk.isEncrypted() && !chunk.isSigned()) {
                return message;
            }

            int messageLength = message.getLengthInBytes();
            WriteBufferByteBased chunkBuffer = new WriteBufferByteBased(messageLength, ByteOrder.LITTLE_ENDIAN);
            message.serialize(chunkBuffer);

            int bodySize = messageLength - chunk.getSecurityHeaderSize() - SECURE_MESSAGE_HEADER_SIZE;
            if (chunk.isEncrypted()) {
                bodySize = decrypt(chunkBuffer, chunk, messageLength);
            }

            if (chunk.isSigned()) {
                verify(chunkBuffer, chunk, messageLength);
            }

            int encryptionOverhead = getEncryptionOverhead(chunk, messageLength);
            int paddingSize = getPaddingSize(chunkBuffer, chunk, messageLength);

            int payloadStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();
            int payloadEnd = payloadStart + bodySize - paddingSize - chunk.getSignatureSize() - chunk.getPaddingOverhead();
            int expectedPaddingSize = messageLength - payloadEnd - chunk.getSignatureSize() - encryptionOverhead - chunk.getPaddingOverhead();

            if (paddingSize != expectedPaddingSize) {
                throw new IllegalArgumentException("Malformed data detected - expected padding size do not match");
            }

            if (chunk.isEncrypted()) {
                byte[] paddingBytes = chunkBuffer.getBytes(payloadEnd, payloadEnd + expectedPaddingSize);
                byte paddingByte = (byte) (paddingSize & 0xff);
                for (int index = 0; index < paddingBytes.length; index++) {
                    if (paddingBytes[index] != paddingByte) {
                        throw new IllegalArgumentException("Malformed padding byte at index " + index);
                    }
                }
            }

            int overhead = paddingSize + chunk.getSignatureSize() + chunk.getPaddingOverhead() + encryptionOverhead;
            updateFrameSize(chunkBuffer, messageLength - overhead);

            return PayloadConverter.pduFromStream(chunkBuffer.getBytes(), message.getResponse());
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private void updateFrame(boolean first, WriteBufferByteBased messageBuffer, Chunk chunk, ChunkType chunkType, Supplier<Integer> sequenceSupplier) throws SerializationException {
        int payloadStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();
        if (chunkType != ChunkType.FINAL) {
            messageBuffer.setPos(3);
            messageBuffer.writeString("chunkType", 8, chunkType.getValue());
        }

        if (!first) {
            messageBuffer.setPos(payloadStart);
            messageBuffer.writeUnsignedLong("sequenceId", 32, sequenceSupplier.get());
        }

        // leave buffer at beginning of message body
        messageBuffer.setPos(payloadStart + 8);
    }

    private void updateFrameSize(WriteBufferByteBased messageBuffer, long frameSize) throws SerializationException {
        int position = messageBuffer.getPos();
        try {
            messageBuffer.setPos(4);
            messageBuffer.writeUnsignedLong("totalLength", 32, frameSize);
        } finally {
            messageBuffer.setPos(position);
        }
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

    private short getPaddingSize(WriteBufferByteBased chunkBuffer, Chunk chunk, int messageLength) {
        if (!chunk.isEncrypted()) {
            return 0;
        }

        int bodyStart = SECURE_MESSAGE_HEADER_SIZE + chunk.getSecurityHeaderSize();
        int bodySize = messageLength - bodyStart;
        int blockCount = bodySize / chunk.getCipherTextBlockSize();
        // bytes we "lost" after payload got decrypted
        int encryptionOverhead = (chunk.getCipherTextBlockSize() * blockCount) - (chunk.getPlainTextBlockSize() * blockCount);

        int paddingEnd = messageLength - chunk.getSignatureSize() - encryptionOverhead - chunk.getPaddingOverhead();
        byte[] padding = chunkBuffer.getBytes(paddingEnd, paddingEnd + chunk.getPaddingOverhead());
        if (padding.length > 2) { // cipher block size exceeds 256 bytes
            int paddingSize = ((padding[1] & 0xFF) << 8) | (padding[0] & 0xFF);
            return (short) (paddingSize & 0xFFFF);
        }
        return (short) (padding[0] & 0xFF);
    }

    protected abstract void verify(WriteBufferByteBased buffer, Chunk chunk, int messageLength) throws Exception;

    protected abstract int decrypt(WriteBufferByteBased chunkBuffer, Chunk chunk, int messageLength) throws Exception;

    protected abstract void encrypt(WriteBufferByteBased buffer, int securityHeaderSize, int plainTextBlockSize, int cipherTextBlockSize, int blockCount) throws Exception;

    protected abstract byte[] sign(byte[] contentsToSign) throws GeneralSecurityException;

}
