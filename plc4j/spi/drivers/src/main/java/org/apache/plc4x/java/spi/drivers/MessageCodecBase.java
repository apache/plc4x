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

    protected abstract int getMinimumHeaderSize();

    protected abstract int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException;

    protected abstract M parseMessage(ReadBufferByteBased readBuffer) throws BufferException;

    protected WriteBufferByteBased createWriteBuffer(int size) {
        return new WriteBufferByteBased(new byte[size]);
    }

    protected ReadBufferByteBased createReadBuffer(byte[] data) {
        return new ReadBufferByteBased(data);
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
                int totalMessageSize = calculateTotalMessageSize(header, availableBytes);
                if (totalMessageSize < 0) {
                    return;
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
                M message = parseMessage(readBuffer);
                messageHandler.accept(message);
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
