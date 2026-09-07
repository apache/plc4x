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
package org.apache.plc4x.java.iec608705104;

import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;
import org.apache.plc4x.java.iec608705104.readwrite.APDU;

/**
 * MessageCodec for IEC 60870-5-104 APDUs.
 *
 * <p>Every APDU starts with the fixed {@link APDU#STARTBYTE} {@code 0x68}
 * followed by a 1-byte length that counts the bytes <em>after</em> the
 * length field itself — so the total frame size is {@code length + 2}.</p>
 */
public class Iec60870MessageCodec extends MessageCodecBase<APDU> {

    private static final int MIN_HEADER = 2;

    public Iec60870MessageCodec(TransportInstance<?> transportInstance, Consumer<APDU> messageHandler) {
        super("IEC 60870-5-104", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MIN_HEADER;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        if ((header[0] & 0xFF) != APDU.STARTBYTE) {
            throw new MessageCodecException(
                "Invalid APDU start byte: 0x" + Integer.toHexString(header[0] & 0xFF));
        }
        return (header[1] & 0xFF) + 2;
    }

    @Override
    protected APDU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return APDU.staticParse(readBuffer);
    }

}
