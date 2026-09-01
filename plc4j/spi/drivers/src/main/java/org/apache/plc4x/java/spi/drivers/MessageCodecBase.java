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
package org.apache.plc4x.java.spi.drivers;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.function.Consumer;

/**
 * Abstract base class for protocol message codecs.
 * <p>
 * Provides common functionality for encoding and decoding protocol messages:
 * <ul>
 *   <li>Sending messages via transport</li>
 *   <li>Receiving and parsing messages from transport</li>
 *   <li>Transport lifecycle management (isOpen, close)</li>
 * </ul>
 *
 * @param <M> the message type this codec handles (must extend Message)
 */
public abstract class MessageCodecBase<M extends Message> {

    protected final Logger logger;
    protected final String protocolName;
    protected final TransportInstance<?> transportInstance;
    protected final Consumer<M> messageHandler;

    protected MessageCodecBase(String protocolName, TransportInstance<?> transportInstance, Consumer<M> messageHandler) {
        this.logger = LoggerFactory.getLogger(getClass());
        this.protocolName = protocolName;
        this.transportInstance = transportInstance;
        this.messageHandler = messageHandler;
    }

    /**
     * Returned from {@link #calculateTotalMessageSize} to say the header does frame a message, but
     * not enough of it has arrived to say how long it is. Nothing is consumed and the receive cycle
     * ends; the codec is asked again when more bytes turn up.
     */
    protected static final int NEED_MORE_DATA = -1;

    /**
     * Returned from {@link #calculateTotalMessageSize} to say the bytes at the front of the buffer
     * are not the start of a message at all - a magic number that does not match, a function code
     * that means nothing. Waiting cannot fix that, because no byte that arrives later changes the
     * bytes already there, so the front byte is dropped and framing tried again from the next one.
     */
    protected static final int DESYNCHRONIZED = -2;

    protected abstract int getMinimumHeaderSize();

    /**
     * Works out how long the message at the front of the buffer is.
     *
     * @return the message length in bytes, or {@link #NEED_MORE_DATA} to be asked again once more
     * has arrived, or {@link #DESYNCHRONIZED} to have the front byte dropped and framing retried.
     * A length below {@link #getMinimumHeaderSize()} is treated as {@link #DESYNCHRONIZED}, since
     * a message cannot be shorter than the header it was read from.
     */
    protected abstract int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException;

    protected abstract M parseMessage(ReadBufferByteBased readBuffer) throws BufferException;

    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size]);
    }

    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data);
    }

    /**
     * Drops the byte at the front of the buffer so the next framing attempt starts from a
     * different place. One byte rather than a guess at the next boundary, because a guess needs
     * knowledge of the protocol that lives in the codec, and because one byte is enough to
     * guarantee the receive cycle always moves.
     */
    private void resynchronize() throws TransportException {
        transportInstance.read(1);
    }

    public void send(M message) throws MessageCodecException {
        try {
            int messageSize = message.getLengthInBytes();
            WriteBufferByteBased writeBuffer = createWriteBuffer(messageSize);
            message.serialize(writeBuffer);
            byte[] messageBytes = writeBuffer.getBytes();

            if (logger.isTraceEnabled()) {
                logger.trace("Sending {} message: {} bytes", protocolName, messageBytes.length);
            }

            transportInstance.write(messageBytes);
        } catch (BufferException e) {
            throw new MessageCodecException("Failed to serialize " + protocolName + " message", e);
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to send " + protocolName + " message", e);
        }
    }

    public void processIncomingData() throws MessageCodecException {
        try {
            while (true) {
                int minimumHeaderSize = getMinimumHeaderSize();

                int availableBytes = transportInstance.getNumBytesAvailable();
                if (availableBytes < minimumHeaderSize) {
                    return;
                }

                byte[] header = transportInstance.peekReadableBytes(minimumHeaderSize);
                int totalMessageSize;
                try {
                    totalMessageSize = calculateTotalMessageSize(header, availableBytes);
                } catch (MessageCodecException e) {
                    // A codec that cannot frame these bytes at all is describing the same
                    // situation as DESYNCHRONIZED, and losing the connection over it would let
                    // anyone who can put a byte on the wire end the conversation.
                    logger.warn("Could not frame incoming {} data, skipping a byte to resynchronise",
                        protocolName, e);
                    resynchronize();
                    continue;
                }

                if (totalMessageSize == NEED_MORE_DATA) {
                    return;
                }

                if (totalMessageSize <= DESYNCHRONIZED || totalMessageSize < minimumHeaderSize) {
                    // Either the codec said so, or it named a length no message could have. Both
                    // have to cost a byte: returning without consuming one would read the same
                    // header again on the next cycle, and the cycle after that, with the transport
                    // still reporting itself open and nothing ever moving.
                    logger.warn("Skipping a byte of {} data that does not begin a message (framed as {})",
                        protocolName, totalMessageSize);
                    resynchronize();
                    continue;
                }

                int receiveCapacity = transportInstance.getReceiveCapacity();
                if (totalMessageSize > receiveCapacity) {
                    // Waiting for this would be waiting forever: the transport cannot hold that
                    // many bytes at once, so they will never all be here together however long we
                    // wait. A length we can never satisfy is not a message we are part way through.
                    logger.warn("Skipping a byte of {} data framed as {} bytes, more than the {} "
                        + "the transport can hold", protocolName, totalMessageSize, receiveCapacity);
                    resynchronize();
                    continue;
                }

                if (availableBytes < totalMessageSize) {
                    if (logger.isTraceEnabled()) {
                        logger.trace("Waiting for complete message: have {} bytes, need {} bytes",
                            availableBytes, totalMessageSize);
                    }
                    return;
                }

                byte[] messageBytes = transportInstance.read(totalMessageSize);

                if (logger.isTraceEnabled()) {
                    logger.trace("Received {} message: {} bytes", protocolName, messageBytes.length);
                }

                ReadBufferByteBased readBuffer = createReadBuffer(messageBytes);
                M message;
                try {
                    message = parseMessage(readBuffer);
                } catch (RuntimeException e) {
                    // Adversarial input can trip an unchecked failure inside a generated parser.
                    // Report it as the parse failure it is, rather than letting it escape past
                    // the caller's error handling and out of whatever thread we are running on.
                    throw new MessageCodecException("Failed to parse " + protocolName + " message", e);
                }
                try {
                    messageHandler.accept(message);
                } catch (RuntimeException e) {
                    // The handler is driver code reacting to a message the peer chose to send. Name
                    // it for what it is rather than letting it travel on as whatever it happened to
                    // be, so a caller can tell handling a message from receiving one.
                    throw new MessageCodecException(
                        "Failed to handle " + protocolName + " message", e);
                }
            }
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to receive " + protocolName + " message", e);
        } catch (BufferException e) {
            throw new MessageCodecException("Failed to parse " + protocolName + " message", e);
        }
    }

    public boolean isOpen() {
        return transportInstance.isOpen();
    }

    public void close() throws MessageCodecException {
        try {
            transportInstance.close();
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to close transport", e);
        }
    }

    protected TransportInstance<?> getTransportInstance() {
        return transportInstance;
    }

}
