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
package org.apache.plc4x.java.umas;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.umas.readwrite.ModbusTcpADU;

import java.util.function.Consumer;

/**
 * MessageCodec for the UMAS wire format (UMAS tunneled inside Modbus/TCP).
 *
 * <p>Every frame starts with the 7-byte MBAP header:
 * <pre>
 *   transactionId(2) + protocolId(2) + length(2) + unitId(1)
 * </pre>
 * The {@code length} field at offset 4 (big-endian, unsigned 16-bit) counts
 * the unitId byte plus the entire PDU that follows, so the total frame size
 * is {@code length + 6}. Same logic as the legacy
 * {@code UmasDriver.ByteLengthEstimator}.</p>
 *
 * <p>UMAS responses use a generic function key (0xFE) and need the original
 * request's function key to discriminate the response subtype during parse.
 * We peek the transaction ID off the MBAP header and look up the tracked
 * function key in {@link UmasFunctionKeyTracker}.</p>
 */
public class UmasMessageCodec extends MessageCodecBase<ModbusTcpADU> {

    private static final int HEADER_SIZE = 6;
    private static final int LENGTH_OFFSET = 4;
    private static final int LENGTH_OVERHEAD = 6;

    /** The connection's tracker, not a shared one: a transaction id only means something within
     * the connection that issued it. */
    private final UmasFunctionKeyTracker functionKeyTracker;

    public UmasMessageCodec(TransportInstance<?> transportInstance, Consumer<ModbusTcpADU> messageHandler,
                            UmasFunctionKeyTracker functionKeyTracker) {
        super("UMAS", transportInstance, messageHandler);
        this.functionKeyTracker = functionKeyTracker;
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        int length = ((header[LENGTH_OFFSET] & 0xFF) << 8) | (header[LENGTH_OFFSET + 1] & 0xFF);
        return length + LENGTH_OVERHEAD;
    }

    @Override
    protected ModbusTcpADU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        // Peek the transaction ID (first 2 bytes, big-endian) so we can pull
        // the matching UMAS function key out of the tracker for response
        // discrimination — without advancing the read position.
        byte[] all = readBuffer.getBytes();
        int transactionId = ((all[0] & 0xFF) << 8) | (all[1] & 0xFF);
        short fk = functionKeyTracker.consumeFunctionKey(transactionId);
        return ModbusTcpADU.staticParse(readBuffer, fk);
    }

}
