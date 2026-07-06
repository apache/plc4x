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
package org.apache.plc4x.java.modbus.rtu;

import org.apache.plc4x.java.modbus.readwrite.DriverType;
import org.apache.plc4x.java.modbus.readwrite.ModbusADU;
import org.apache.plc4x.java.modbus.readwrite.ModbusRtuADU;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.MessageCodecBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;

import java.util.function.Consumer;

/**
 * MessageCodec for Modbus RTU protocol.
 * Handles the encoding and decoding of Modbus RTU ADU (Application Data Unit) messages.
 *
 * RTU framing: 1 byte address + PDU + 2 bytes CRC.
 * Since RTU has no explicit length field in the header, the response frame size is
 * derived from the function code (see {@link #calculateTotalMessageSize}): exception
 * responses (fc &gt;= 0x80) are 5 bytes, read responses (0x01-0x04, 0x14, 0x15, 0x17,
 * including the file-record operations) are 5 + the byte-count carried in the third
 * header byte, and write echoes (0x05, 0x06, 0x0F, 0x10) are a fixed 8 bytes.
 * {@link #processIncomingData()} peeks
 * a candidate frame, only consumes it once the generated parser has validated its CRC,
 * and otherwise resynchronizes byte-wise.
 *
 * Limitation: since the size table above is response-shaped, this codec assumes it is
 * only ever used to decode responses (as constructed) - it cannot correctly frame a
 * stream of raw Modbus RTU requests.
 */
public class ModbusRtuMessageCodec extends MessageCodecBase<ModbusRtuADU> {

    // Minimum RTU message: address (1) + function code (1) + CRC (2) = 4 bytes
    private static final int MODBUS_RTU_MIN_SIZE = 4;
    // address + fc + exception code + CRC
    private static final int EXCEPTION_RESPONSE_SIZE = 5;
    // address + fc + 2 bytes address + 2 bytes value/quantity + CRC
    private static final int WRITE_ECHO_RESPONSE_SIZE = 8;
    // address + fc + third byte (byteCount for reads) — all the sizing
    // table needs to compute the expected response size
    private static final int SIZING_HEADER_SIZE = 3;

    private long resyncSkippedBytes;

    public ModbusRtuMessageCodec(TransportInstance<?> transportInstance, Consumer<ModbusRtuADU> messageHandler) {
        super("Modbus RTU", transportInstance, messageHandler);
    }

    @Override
    protected int getMinimumHeaderSize() {
        return MODBUS_RTU_MIN_SIZE;
    }

    /**
     * Computes the expected RESPONSE frame size from the peeked header
     * (address, function code, third byte). RTU has no length field, so
     * the size follows from the function code:
     * exception (fc >= 0x80) = 5; reads (0x01-0x04, 0x14, 0x15, 0x17 —
     * including the file-record read/write operations, whose responses
     * carry byteCount as the third byte covering the whole item array
     * just like the plain reads) = 5 + byteCount (the third byte);
     * write echoes (0x05, 0x06, 0x0F, 0x10) = 8.
     * Returns -1 for unknown function codes — the caller treats that as
     * desynchronization, not as "wait for more data".
     */
    @Override
    protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
        int functionCode = header[1] & 0xFF;
        if (functionCode >= 0x80) {
            return EXCEPTION_RESPONSE_SIZE;
        }
        switch (functionCode) {
            case 0x01, 0x02, 0x03, 0x04, 0x14, 0x15, 0x17:
                return MODBUS_RTU_MIN_SIZE + 1 + (header[2] & 0xFF);
            case 0x05, 0x06, 0x0F, 0x10:
                return WRITE_ECHO_RESPONSE_SIZE;
            default:
                return -1;
        }
    }

    @Override
    protected ModbusRtuADU parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
        ModbusADU modbusParsed = ModbusADU.staticParse(readBuffer, DriverType.MODBUS_RTU, true);
        if (!(modbusParsed instanceof ModbusRtuADU modbusADU)) {
            throw new BufferException("Parsed message is not a ModbusRtuADU");
        }
        return modbusADU;
    }

    /**
     * RTU-specific receive loop: peek before consume. The base class's
     * loop reads a frame before validating it, which loses bytes whenever
     * the stream carries more (or less) than exactly one frame — the norm
     * on a shared serial line. Here a frame is only consumed after the
     * generated parser (which validates the CRC checksum field) accepted
     * it from a peek; anything unparseable advances by a single byte
     * (byte-wise resynchronization, mirroring the plc4go modbus codec).
     * <p>
     * While resynchronizing (i.e. a prior candidate at this stream position
     * already failed), a header that looks valid but demands more bytes than
     * are available is treated as another failed candidate rather than as a
     * reason to wait: a byte-shifted resync candidate can accidentally look
     * like a legitimate function code with a large byte count, which would
     * otherwise stall the codec forever waiting for bytes that were never
     * going to arrive. Waiting is only correct at a clean (non-resyncing)
     * stream position, where an undersized buffer really does mean "the next
     * frame hasn't fully arrived yet". The corollary trade-off: a GENUINE
     * partial frame sitting directly behind garbage is skipped rather than
     * awaited — that response is lost and recovered by the master's own
     * request timeout and retry, which beats stalling every frame queued
     * behind an undersized garbage candidate that will never complete.
     */
    @Override
    public void processIncomingData() throws MessageCodecException {
        try {
            while (true) {
                int availableBytes = getTransportInstance().getNumBytesAvailable();
                if (availableBytes < MODBUS_RTU_MIN_SIZE) {
                    return;
                }
                byte[] header = getTransportInstance().peekReadableBytes(SIZING_HEADER_SIZE);
                int expectedSize = calculateTotalMessageSize(header, availableBytes);
                if (expectedSize < 0) {
                    skipOneByte("unknown function code 0x" + Integer.toHexString(header[1] & 0xFF));
                    continue;
                }
                if (availableBytes < expectedSize) {
                    if (resyncSkippedBytes > 0) {
                        skipOneByte("candidate frame during resync needs " + expectedSize
                            + " bytes, only " + availableBytes + " available");
                        continue;
                    }
                    return; // never consume a partial frame
                }
                byte[] frame = getTransportInstance().peekReadableBytes(expectedSize);
                ModbusRtuADU message;
                try {
                    message = parseMessage(createReadBuffer(frame));
                } catch (BufferException e) {
                    skipOneByte("frame failed validation: " + e.getMessage());
                    continue;
                }
                getTransportInstance().read(expectedSize); // consume the validated frame
                noteResyncComplete();
                messageHandler.accept(message);
            }
        } catch (TransportException e) {
            throw new MessageCodecException("Failed to receive Modbus RTU message", e);
        }
    }

    private void skipOneByte(String reason) throws TransportException {
        if (resyncSkippedBytes == 0) {
            logger.warn("Modbus RTU stream out of sync ({}), resynchronizing byte-wise", reason);
        }
        getTransportInstance().read(1);
        resyncSkippedBytes++;
    }

    private void noteResyncComplete() {
        if (resyncSkippedBytes > 0) {
            logger.warn("Modbus RTU stream resynchronized after skipping {} bytes", resyncSkippedBytes);
            resyncSkippedBytes = 0;
        }
    }

}
