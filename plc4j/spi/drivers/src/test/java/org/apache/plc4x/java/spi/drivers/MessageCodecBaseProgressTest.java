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

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.MessageCodecBaseTest.FakeTransport;
import org.apache.plc4x.java.spi.drivers.MessageCodecBaseTest.TestMessage;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Whoever is on the other end decides what the length field says. A receive cycle therefore has to
 * end in one of three states - a message delivered, genuinely missing bytes waited for, or a
 * failure reported - and never in "read the same header again, unchanged, forever".
 */
class MessageCodecBaseProgressTest {

    /** A codec whose framing decision is whatever the test wants it to be. */
    private static class SizeCodec extends MessageCodecBase<TestMessage> {
        private final java.util.function.IntUnaryOperator sizeFromHeader;

        SizeCodec(FakeTransport transport, List<TestMessage> received,
                  java.util.function.IntUnaryOperator sizeFromHeader) {
            super("TEST", transport, received::add);
            this.sizeFromHeader = sizeFromHeader;
        }

        @Override protected int getMinimumHeaderSize() { return 2; }

        @Override protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
            return sizeFromHeader.applyAsInt(header[0] & 0xFF);
        }

        @Override protected TestMessage parseMessage(ReadBufferByteBased readBuffer) throws BufferException {
            return new TestMessage(readBuffer.readUnsignedShort(8));
        }

        @Override protected ReadBufferByteBased createReadBuffer(byte[] data) {
            return new ReadBufferByteBased(data,
                WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        }
    }

    /** Runs enough receive cycles to tell forward progress from standing still. */
    private static int bytesConsumedOverRepeatedCycles(FakeTransport transport, MessageCodecBase<?> codec) {
        int before = transport.getNumBytesAvailable();
        for (int cycle = 0; cycle < 32; cycle++) {
            try {
                codec.processIncomingData();
            } catch (MessageCodecException e) {
                // Reporting a failure is a fine outcome; standing still is not.
            }
        }
        return before - transport.getNumBytesAvailable();
    }

    @Test
    void aDesynchronisedCodecDoesNotStallTheReceiveCycle() throws Exception {
        FakeTransport transport = new FakeTransport();
        transport.feed((byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04);
        SizeCodec codec = new SizeCodec(transport, new ArrayList<>(),
            h -> MessageCodecBase.DESYNCHRONIZED);

        int consumed = bytesConsumedOverRepeatedCycles(transport, codec);
        assertTrue(consumed > 0,
            "a codec that says it is out of step must cost the byte, not repeat forever");
    }

    @Test
    void needingMoreDataStillWaitsForIt() throws Exception {
        FakeTransport transport = new FakeTransport();
        transport.feed((byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04);
        SizeCodec codec = new SizeCodec(transport, new ArrayList<>(),
            h -> MessageCodecBase.NEED_MORE_DATA);

        // Modbus-ASCII and C-Bus scan for a CR/LF terminator and say this until they find one, so
        // consuming here would eat the front of a frame that had simply not finished arriving.
        assertEquals(0, bytesConsumedOverRepeatedCycles(transport, codec),
            "waiting for the rest of a frame must not consume the part that arrived");
    }

    @Test
    void aSizeBelowTheHeaderDoesNotStallTheReceiveCycle() throws Exception {
        FakeTransport transport = new FakeTransport();
        transport.feed((byte) 0x00, (byte) 0x00, (byte) 0x00, (byte) 0x00);
        // The ADS "00 00 FA FF FF FF" shape: a length that cannot even cover its own header.
        SizeCodec codec = new SizeCodec(transport, new ArrayList<>(), h -> 0);

        int consumed = bytesConsumedOverRepeatedCycles(transport, codec);
        assertTrue(consumed > 0,
            "a size below the minimum header size must cost a byte, not repeat forever");
    }

    @Test
    void aFramingFailureDoesNotStallTheReceiveCycle() throws Exception {
        FakeTransport transport = new FakeTransport();
        transport.feed((byte) 0xFF, (byte) 0xFF, (byte) 0xFF, (byte) 0xFF);
        MessageCodecBase<TestMessage> codec = new MessageCodecBase<>("TEST", transport, m -> { }) {
            @Override protected int getMinimumHeaderSize() { return 2; }
            @Override protected int calculateTotalMessageSize(byte[] header, int availableBytes)
                throws MessageCodecException {
                throw new MessageCodecException("cannot frame this");
            }
            @Override protected TestMessage parseMessage(ReadBufferByteBased readBuffer) {
                return new TestMessage(0);
            }
        };

        int consumed = 0;
        int before = transport.getNumBytesAvailable();
        for (int cycle = 0; cycle < 32; cycle++) {
            try {
                codec.processIncomingData();
            } catch (MessageCodecException e) {
                // Reporting it is fine; silently making no progress is not.
            }
        }
        consumed = before - transport.getNumBytesAvailable();
        assertTrue(consumed > 0 || !transport.isOpen(),
            "a codec that cannot frame the data must resync or report the transport closed");
    }

    @Test
    void resyncFindsTheNextGoodFrame() throws Exception {
        FakeTransport transport = new FakeTransport();
        // Two bytes of rubbish that frame as size 1 (below the 2-byte header), then a good
        // 2-byte frame carrying the value 7.
        transport.feed((byte) 0x01, (byte) 0x01, (byte) 0x02, (byte) 0x07);
        List<TestMessage> received = new ArrayList<>();
        SizeCodec codec = new SizeCodec(transport, received, h -> h);

        for (int cycle = 0; cycle < 32; cycle++) {
            codec.processIncomingData();
        }
        assertEquals(0, transport.getNumBytesAvailable(),
            "resynchronising must work through the rubbish rather than sit on it");
    }

    @Test
    void aValidFrameIsStillDelivered() throws Exception {
        FakeTransport transport = new FakeTransport();
        transport.feed((byte) 0x02, (byte) 0x09);
        List<TestMessage> received = new ArrayList<>();
        SizeCodec codec = new SizeCodec(transport, received, h -> h);

        codec.processIncomingData();
        assertEquals(1, received.size(), "a well-framed message must still arrive");
        assertEquals(0, transport.getNumBytesAvailable());
    }

    /**
     * A length larger than the transport can ever hold at once cannot be waited for: those bytes
     * will never be present together however long we wait. Telling that apart from an ordinary
     * partial message is what the capacity is for.
     */
    @Test
    void aLengthLargerThanTheTransportCanHoldDoesNotStall() throws Exception {
        FakeTransport transport = new FakeTransport() {
            @Override
            public int getReceiveCapacity() {
                return 64;
            }
        };
        transport.feed((byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04);
        SizeCodec codec = new SizeCodec(transport, new ArrayList<>(), h -> 100_000);

        assertTrue(bytesConsumedOverRepeatedCycles(transport, codec) > 0,
            "a length the transport can never satisfy must cost a byte, not be waited on forever");
    }

    /**
     * A handler is driver code reacting to a message the peer chose to send, so what it throws
     * belongs to handling rather than to the transport underneath.
     */
    @Test
    void aFailingHandlerIsReportedAsACodecFailure() {
        FakeTransport transport = new FakeTransport();
        transport.feed((byte) 0x02, (byte) 0x09);
        MessageCodecBase<TestMessage> codec = new MessageCodecBase<>("TEST", transport, m -> {
            throw new IllegalStateException("the driver could not use this");
        }) {
            @Override protected int getMinimumHeaderSize() { return 2; }
            @Override protected int calculateTotalMessageSize(byte[] header, int availableBytes) {
                return header[0] & 0xFF;
            }
            @Override protected TestMessage parseMessage(ReadBufferByteBased readBuffer) {
                return new TestMessage(1);
            }
        };

        MessageCodecException e = assertThrows(MessageCodecException.class, codec::processIncomingData);
        assertInstanceOf(IllegalStateException.class, e.getCause());
        assertTrue(e.getMessage().contains("handle"),
            "the failure should say it was handling rather than receiving, but was: " + e.getMessage());
    }
}
