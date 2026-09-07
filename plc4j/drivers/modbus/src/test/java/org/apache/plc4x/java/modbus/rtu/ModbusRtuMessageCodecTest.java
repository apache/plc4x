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

import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

class ModbusRtuMessageCodecTest {

    @SuppressWarnings("unchecked")
    @Test
    void testConstruction() {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusRtuADU> handler = mock(Consumer.class);

        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transportInstance, handler);
        assertNotNull(codec);
    }

    @SuppressWarnings("unchecked")
    @Test
    void testMinimumHeaderSize() throws Exception {
        TransportInstance<?> transportInstance = mock(TransportInstance.class);
        Consumer<ModbusRtuADU> handler = mock(Consumer.class);
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transportInstance, handler);

        Method method = codec.getClass().getDeclaredMethod("getMinimumHeaderSize");
        method.setAccessible(true);
        assertEquals(4, method.invoke(codec));
    }

    /**
     * In-memory transport double: bytes are queued via deliver(...) and the
     * codec consumes them through the real TransportInstance contract.
     * Strict: peeking/reading more than available is a test bug and throws.
     */
    static final class ScriptedTransport implements TransportInstance<TransportConfiguration> {
        private final ByteArrayOutputStream buffer = new ByteArrayOutputStream();
        private int readPosition;
        private boolean open = true;

        void deliver(byte[] bytes) {
            buffer.writeBytes(bytes);
        }

        private byte[] bufferedBytes() {
            return buffer.toByteArray();
        }

        @Override
        public TransportConfiguration getConfiguration() {
            return null;
        }

        @Override
        public boolean isOpen() {
            return open;
        }

        @Override
        public int getNumBytesAvailable() {
            return buffer.size() - readPosition;
        }

        @Override
        public byte[] peekReadableBytes(int numBytes) throws TransportException {
            if (numBytes > getNumBytesAvailable()) {
                throw new TransportException("peek beyond available: " + numBytes);
            }
            byte[] all = bufferedBytes();
            byte[] result = new byte[numBytes];
            System.arraycopy(all, readPosition, result, 0, numBytes);
            return result;
        }

        @Override
        public byte[] read(int numBytes) throws TransportException {
            byte[] result = peekReadableBytes(numBytes);
            readPosition += numBytes;
            return result;
        }

        @Override
        public void write(byte[] bytes) {
            // not exercised by these tests
        }

        @Override
        public void close() {
            open = false;
        }
    }

    private static byte[] wireBytes(ModbusRtuADU adu) throws Exception {
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[adu.getLengthInBytes()]);
        adu.serialize(writeBuffer);
        return writeBuffer.getBytes();
    }

    private static byte[] readResponseFrame(int unitId, byte[] registers) throws Exception {
        return wireBytes(new ModbusRtuADU((short) unitId, new ModbusPDUReadHoldingRegistersResponse(registers)));
    }

    private static byte[] readFileRecordResponseFrame(int unitId, byte[] data) throws Exception {
        List<ModbusPDUReadFileRecordResponseItem> items = new ArrayList<>();
        items.add(new ModbusPDUReadFileRecordResponseItem((short) 6, data));
        return wireBytes(new ModbusRtuADU((short) unitId, new ModbusPDUReadFileRecordResponse(items)));
    }

    private static byte[] concat(byte[]... parts) {
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        for (byte[] part : parts) {
            out.writeBytes(part);
        }
        return out.toByteArray();
    }

    @Test
    void singleCompleteResponseIsDispatched() throws Exception {
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        transport.deliver(readResponseFrame(1, new byte[]{0x11, 0x22}));
        codec.processIncomingData();

        assertEquals(1, received.size());
        assertEquals(1, received.get(0).getAddress());
    }

    @Test
    void twoResponsesInOneDeliveryAreBothDispatched() throws Exception {
        // The shared-port batching case — the bug this change fixes.
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        transport.deliver(concat(
            readResponseFrame(1, new byte[]{0x11, 0x22}),
            readResponseFrame(2, new byte[]{0x33, 0x44, 0x55, 0x66})));
        codec.processIncomingData();

        assertEquals(2, received.size());
        assertEquals(1, received.get(0).getAddress());
        assertEquals(2, received.get(1).getAddress());
    }

    @Test
    void partialFrameWaitsAndCompletesAcrossDeliveries() throws Exception {
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        byte[] frame = readResponseFrame(1, new byte[]{0x11, 0x22, 0x33, 0x44});
        transport.deliver(java.util.Arrays.copyOfRange(frame, 0, 5)); // header + partial data
        codec.processIncomingData();
        assertEquals(0, received.size(), "nothing may be consumed while the frame is incomplete");
        assertEquals(5, transport.getNumBytesAvailable(), "partial bytes must remain buffered");

        transport.deliver(java.util.Arrays.copyOfRange(frame, 5, frame.length));
        codec.processIncomingData();
        assertEquals(1, received.size());
    }

    @Test
    void garbagePrefixIsSkippedToTheNextValidFrame() throws Exception {
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        // 0x00 is not a known function code in byte 2 position for any
        // alignment of this prefix — forces byte-wise resync.
        transport.deliver(concat(new byte[]{0x00, 0x00, 0x00}, readResponseFrame(3, new byte[]{0x77, 0x77})));
        codec.processIncomingData();

        assertEquals(1, received.size());
        assertEquals(3, received.get(0).getAddress());
    }

    @Test
    void corruptCrcFrameIsSkippedAndNextFrameParses() throws Exception {
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        byte[] corrupt = readResponseFrame(1, new byte[]{0x11, 0x22});
        corrupt[corrupt.length - 1] ^= (byte) 0xFF; // break the CRC
        transport.deliver(concat(corrupt, readResponseFrame(2, new byte[]{0x33, 0x44})));
        codec.processIncomingData();

        assertEquals(1, received.size(), "the corrupt frame is skipped, the valid one parses");
        assertEquals(2, received.get(0).getAddress());
    }

    @Test
    void writeEchoAndExceptionResponsesAreFramedCorrectly() throws Exception {
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        byte[] writeEcho = wireBytes(new ModbusRtuADU((short) 4, new ModbusPDUWriteSingleRegisterResponse(0x0010, 0x1234)));
        assertEquals(8, writeEcho.length, "sanity: write echo is the fixed 8-byte frame");
        byte[] exception = wireBytes(new ModbusRtuADU((short) 5, new ModbusPDUError(ModbusErrorCode.ILLEGAL_FUNCTION)));
        assertEquals(5, exception.length, "sanity: exception response is the fixed 5-byte frame");

        transport.deliver(concat(writeEcho, exception));
        codec.processIncomingData();

        assertEquals(2, received.size());
        assertEquals(4, received.get(0).getAddress());
        assertEquals(5, received.get(1).getAddress());
    }

    @Test
    void fileRecordResponseIsFramedCorrectly() throws Exception {
        // fc 0x14 (read file record) responses carry byteCount as the third
        // byte just like the plain reads, but the sizing switch used to
        // omit 0x14/0x15 entirely, so this frame would be discarded byte-wise
        // and the following frame would be the only one dispatched.
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        transport.deliver(concat(
            readFileRecordResponseFrame(1, new byte[]{0x11, 0x22, 0x33, 0x44}),
            readResponseFrame(2, new byte[]{0x33, 0x44})));
        codec.processIncomingData();

        assertEquals(2, received.size(), "fc 0x14 (read file record) response must be framed, not discarded");
        assertEquals(1, received.get(0).getAddress());
        assertEquals(2, received.get(1).getAddress());
    }

    @Test
    void foreignUnitIdFrameIsStillDispatched() throws Exception {
        // The codec frames and dispatches everything on the wire; unit-id
        // filtering is the connection layer's job.
        ScriptedTransport transport = new ScriptedTransport();
        List<ModbusRtuADU> received = new ArrayList<>();
        ModbusRtuMessageCodec codec = new ModbusRtuMessageCodec(transport, received::add);

        transport.deliver(readResponseFrame(42, new byte[]{0x01, 0x02}));
        codec.processIncomingData();

        assertEquals(1, received.size());
        assertEquals(42, received.get(0).getAddress());
    }
}
