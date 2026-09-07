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
package org.apache.plc4x.java.bacnetip;

import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;

import java.util.function.Consumer;

/**
 * Decodes the BVLC envelope used by BACnet/IP:
 * <pre>
 *   byte 0   = 0x81 (BACnet/IP type marker)
 *   byte 1   = BVLC function
 *   byte 2-3 = total BVLC length (big endian, includes header)
 *   ...      = payload (NPDU + APDU)
 * </pre>
 * The total message size lives at offset 2, big-endian, so we need at least
 * 4 header bytes to learn how many bytes to consume.
 */
public class BacNetIpMessageCodec extends MessageCodecBase<BVLC> {

    private static final int HEADER_SIZE = 4;
    private static final int BACNET_TYPE = 0x81;

    public BacNetIpMessageCodec(TransportInstance<?> transportInstance, Consumer<BVLC> messageHandler) {
        super("BACnet/IP", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return HEADER_SIZE;
    }

    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) throws MessageCodecException {
        if ((header[0] & 0xFF) != BACNET_TYPE) {
            // Stream is out of sync — let the framework discard a byte and resync.
            throw new MessageCodecException("Invalid BACnet/IP type marker: 0x"
                + Integer.toHexString(header[0] & 0xFF));
        }
        return ((header[2] & 0xFF) << 8) | (header[3] & 0xFF);
    }

    @Override
    protected BVLC parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        return BVLC.staticParse(readBuffer);
    }

}
