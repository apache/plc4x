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
package org.apache.plc4x.java.firmata;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;

import java.util.function.Consumer;
import org.apache.plc4x.java.firmata.readwrite.FirmataMessage;

/**
 * Message codec for the Firmata (MIDI-style) wire format. Each Firmata message
 * starts with a status byte where the high nibble identifies the message type
 * and the low nibble carries channel/pin information. The total message length
 * is therefore derivable from the first byte (plus, for sysex messages, by
 * scanning for the {@code 0xF7} terminator).
 *
 * <p>Logic mirrors the legacy {@code FirmataDriver.ByteLengthEstimator}:
 * <ul>
 *   <li>{@code 0xE0} (Analog IO) / {@code 0x90} (Digital IO): 3 bytes</li>
 *   <li>{@code 0xC0} (Subscribe Analog) / {@code 0xD0} (Subscribe Digital): 2 bytes</li>
 *   <li>{@code 0xF0} (System message): variable; sub-command in the low nibble
 *       further refines the length, with {@code 0xF0 0x00} (sysex) being
 *       framed by a trailing {@code 0xF7}.</li>
 * </ul>
 */
public class FirmataMessageCodec extends MessageCodecBase<FirmataMessage> {

    private static final int MIN_HEADER = 1;
    private static final byte SYSEX_END = (byte) 0xF7;

    public FirmataMessageCodec(TransportInstance<?> transportInstance, Consumer<FirmataMessage> messageHandler) {
        super("Firmata", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MIN_HEADER;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        int first = header[0] & 0xFF;
        int type = first & 0xF0;
        switch (type) {
            case 0xE0:
            case 0x90:
                return 3;
            case 0xC0:
            case 0xD0:
                return 2;
            case 0xF0:
                return calculateSystemMessageSize(first, availableBytes);
            default:
                throw new MessageCodecException("Invalid Firmata packet type: 0x" + Integer.toHexString(first));
        }
    }

    private int calculateSystemMessageSize(int firstByte, int availableBytes) throws MessageCodecException {
        int sub = firstByte & 0x0F;
        switch (sub) {
            case 0x00:
                // Sysex — scan forward for the 0xF7 terminator. We need to peek
                // at all currently-buffered bytes to find it; if the terminator
                // isn't in the buffer yet, return -1 so the codec waits for more.
                if (availableBytes < 2) {
                    return -1;
                }
                try {
                    byte[] all = transportInstance.peekReadableBytes(availableBytes);
                    for (int i = 1; i < all.length; i++) {
                        if (all[i] == SYSEX_END) {
                            return i + 1;
                        }
                    }
                    return -1;
                } catch (TransportException e) {
                    throw new MessageCodecException("Failed to peek at incoming Firmata bytes", e);
                }
            case 0x04:
            case 0x05:
            case 0x09:
                return 3;
            case 0x0F:
                return 1;
            default:
                throw new MessageCodecException("Invalid Firmata sysex sub-command: 0x"
                    + Integer.toHexString(firstByte));
        }
    }

    @Override
    protected FirmataMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return FirmataMessage.staticParse(readBuffer, true);
    }

}
