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
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayOutputStream;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.function.Consumer;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class MessageCodecBaseTest {

    /**
     * Trivial wire message: {@code [ totalLength=2 ][ value ]} (two unsigned bytes, big-endian).
     */
    record TestMessage(int value) implements Message {
        @Override public int getLengthInBytes() { return 2; }
        @Override public int getLengthInBits() { return 16; }
        @Override public void serialize(WriteBuffer writeBuffer) throws BufferException {
            writeBuffer.writeUnsignedShort(8, (short) 2);
            writeBuffer.writeUnsignedShort(8, (short) value);
        }
    }

    /** Codec for {@link TestMessage}: 1-byte header carrying the total message length. */
    static class TestCodec extends MessageCodecBase<TestMessage> {
        TestCodec(TransportInstance<?> transport, Consumer<TestMessage> handler) {
            super("TEST", transport, handler);
        }
        @Override protected int getMinimumHeaderSize() { return 1; }
        @Override protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
            return header[0] & 0xFF;
        }
        @Override protected TestMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
            readBuffer.readUnsignedShort(8); // length
            return new TestMessage(readBuffer.readUnsignedShort(8));
        }
        // The default buffers carry no integer encoding; supply the big-endian defaults the
        // generated (un)marshalling expects so writeUnsignedShort/readUnsignedShort work.
        @Override protected WriteBufferByteBased createWriteBuffer(int size) {
            return new WriteBufferByteBased(new byte[size],
                WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        }
        @Override protected ReadBufferByteBased createReadBuffer(byte[] data) {
            return new ReadBufferByteBased(data,
                WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        }
    }

    /** In-memory transport: inbound bytes are staged via {@link #feed}; writes are captured. */
    static class FakeTransport implements TransportInstance<TransportConfiguration> {
        private byte[] inbound = new byte[0];
        private int pos = 0;
        private final ByteArrayOutputStream written = new ByteArrayOutputStream();
        private boolean open = true;
        private boolean closed = false;

        void feed(byte... data) {
            byte[] remaining = Arrays.copyOfRange(inbound, pos, inbound.length);
            byte[] combined = new byte[remaining.length + data.length];
            System.arraycopy(remaining, 0, combined, 0, remaining.length);
            System.arraycopy(data, 0, combined, remaining.length, data.length);
            inbound = combined;
            pos = 0;
        }

        byte[] written() { return written.toByteArray(); }
        boolean isClosed() { return closed; }
        void setOpen(boolean open) { this.open = open; }

        @Override public TransportConfiguration getConfiguration() { return null; }
        @Override public boolean isOpen() { return open; }
        @Override public int getNumBytesAvailable() { return inbound.length - pos; }
        @Override public byte[] peekReadableBytes(int numBytes) { return Arrays.copyOfRange(inbound, pos, pos + numBytes); }
        @Override public byte[] read(int numBytes) {
            byte[] out = Arrays.copyOfRange(inbound, pos, pos + numBytes);
            pos += numBytes;
            return out;
        }
        @Override public void write(byte[] bytes) throws TransportException { written.writeBytes(bytes); }
        @Override public void close() { closed = true; open = false; }
    }

    @Test
    void sendSerializesMessageToTransport() throws Exception {
        FakeTransport transport = new FakeTransport();
        TestCodec codec = new TestCodec(transport, m -> { });

        codec.send(new TestMessage(0x2A));

        assertArrayEquals(new byte[]{0x02, 0x2A}, transport.written());
    }

    @Test
    void processIncomingDataDeliversCompleteMessage() throws Exception {
        FakeTransport transport = new FakeTransport();
        List<TestMessage> received = new ArrayList<>();
        TestCodec codec = new TestCodec(transport, received::add);

        transport.feed((byte) 0x02, (byte) 0x7F);
        codec.processIncomingData();

        assertEquals(1, received.size());
        assertEquals(0x7F, received.get(0).value());
        assertEquals(0, transport.getNumBytesAvailable(), "frame should be fully consumed");
    }

    @Test
    void processIncomingDataWaitsForIncompleteMessage() throws Exception {
        FakeTransport transport = new FakeTransport();
        List<TestMessage> received = new ArrayList<>();
        TestCodec codec = new TestCodec(transport, received::add);

        // header says total length 2, but only 1 byte is available
        transport.feed((byte) 0x02);
        codec.processIncomingData();

        assertTrue(received.isEmpty(), "must not deliver until the full frame arrived");
        assertEquals(1, transport.getNumBytesAvailable(), "partial frame must stay buffered");

        // now the rest arrives
        transport.feed((byte) 0x55);
        codec.processIncomingData();
        assertEquals(1, received.size());
        assertEquals(0x55, received.get(0).value());
    }

    @Test
    void processIncomingDataDeliversMultipleFrames() throws Exception {
        FakeTransport transport = new FakeTransport();
        List<TestMessage> received = new ArrayList<>();
        TestCodec codec = new TestCodec(transport, received::add);

        transport.feed((byte) 0x02, (byte) 0x01, (byte) 0x02, (byte) 0x02);
        codec.processIncomingData();

        assertEquals(2, received.size());
        assertEquals(0x01, received.get(0).value());
        assertEquals(0x02, received.get(1).value());
    }

    @Test
    void processIncomingDataReturnsWhenNoHeaderYet() throws Exception {
        FakeTransport transport = new FakeTransport();
        List<TestMessage> received = new ArrayList<>();
        TestCodec codec = new TestCodec(transport, received::add);

        codec.processIncomingData(); // nothing fed
        assertTrue(received.isEmpty());
    }

    @Test
    void isOpenAndCloseDelegateToTransport() throws Exception {
        FakeTransport transport = new FakeTransport();
        TestCodec codec = new TestCodec(transport, m -> { });

        assertTrue(codec.isOpen());
        transport.setOpen(false);
        assertFalse(codec.isOpen());

        codec.close();
        assertTrue(transport.isClosed());
    }

    @Test
    void sendWrapsTransportFailureAsCodecException() {
        TransportInstance<TransportConfiguration> failing = new FakeTransport() {
            @Override public void write(byte[] bytes) throws TransportException {
                throw new TransportException("boom");
            }
        };
        TestCodec codec = new TestCodec(failing, m -> { });
        assertThrows(MessageCodecException.class, () -> codec.send(new TestMessage(1)));
    }
}
