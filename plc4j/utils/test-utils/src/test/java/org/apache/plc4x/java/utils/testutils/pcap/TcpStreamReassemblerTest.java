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
package org.apache.plc4x.java.utils.testutils.pcap;

import org.junit.jupiter.api.Test;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class TcpStreamReassemblerTest {

    // Simple framing: 2-byte big-endian length at offset 0, length = total message size
    private static final FramingSpec SIMPLE_FRAMING = new FramingSpec(0, 2, true, 0);

    // EIP-like framing: 2-byte LE length at offset 2, total = length + 24 byte header
    private static final FramingSpec EIP_FRAMING = new FramingSpec(2, 2, false, 24);

    // Modbus TCP: 2-byte BE length at offset 4, total = length + 6 (MBAP header)
    private static final FramingSpec MODBUS_FRAMING = new FramingSpec(4, 2, true, 6);

    @Test
    void testSingleMessageExtraction() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        // Create a message: 2-byte length header (0x0005 = 5 bytes total) + 3 payload bytes
        byte[] segment = new byte[]{0x00, 0x05, 0x01, 0x02, 0x03};
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, segment);

        List<FramedMessage> messages = reassembler.extractMessages();

        assertEquals(1, messages.size());
        assertArrayEquals(segment, messages.getFirst().rawBytes());
        assertFalse(messages.getFirst().isResponse(), "dst==protocolPort should be request");
        assertEquals(0, messages.getFirst().index());
    }

    @Test
    void testResponseDirectionDetection() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        byte[] segment = new byte[]{0x00, 0x05, 0x01, 0x02, 0x03};
        // Source port matches protocol port → response
        reassembler.addSegment("10.0.0.2", 8080, "10.0.0.1", 12345, segment);

        List<FramedMessage> messages = reassembler.extractMessages();

        assertEquals(1, messages.size());
        assertTrue(messages.getFirst().isResponse(), "src==protocolPort should be response");
    }

    @Test
    void testMultipleMessagesInSingleSegment() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        // Two messages concatenated: [0x0003, 0xAA] + [0x0004, 0xBB, 0xCC]
        byte[] segment = new byte[]{
            0x00, 0x03, (byte) 0xAA,       // msg1: length=3, 1 byte payload
            0x00, 0x04, (byte) 0xBB, (byte) 0xCC  // msg2: length=4, 2 bytes payload
        };
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, segment);

        List<FramedMessage> messages = reassembler.extractMessages();

        assertEquals(2, messages.size());
        assertArrayEquals(new byte[]{0x00, 0x03, (byte) 0xAA}, messages.get(0).rawBytes());
        assertArrayEquals(new byte[]{0x00, 0x04, (byte) 0xBB, (byte) 0xCC}, messages.get(1).rawBytes());
    }

    @Test
    void testMessageSpanningMultipleSegments() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        // Message of length 5, split across two segments
        byte[] segment1 = new byte[]{0x00, 0x05, 0x01};  // header + partial payload
        byte[] segment2 = new byte[]{0x02, 0x03};          // rest of payload

        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, segment1);
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, segment2);

        List<FramedMessage> messages = reassembler.extractMessages();

        assertEquals(1, messages.size());
        assertArrayEquals(new byte[]{0x00, 0x05, 0x01, 0x02, 0x03}, messages.getFirst().rawBytes());
    }

    @Test
    void testIncompleteMessageAtEndOfStream() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        // One complete message followed by an incomplete one
        byte[] segment = new byte[]{
            0x00, 0x03, (byte) 0xAA,  // msg1: complete (3 bytes)
            0x00, 0x05, (byte) 0xBB   // msg2: incomplete (need 5 bytes, only have 3)
        };
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, segment);

        List<FramedMessage> messages = reassembler.extractMessages();

        // Only the complete message should be extracted
        assertEquals(1, messages.size());
        assertArrayEquals(new byte[]{0x00, 0x03, (byte) 0xAA}, messages.getFirst().rawBytes());
    }

    @Test
    void testEmptyPayloadIsIgnored() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, new byte[0]);
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, null);

        List<FramedMessage> messages = reassembler.extractMessages();
        assertTrue(messages.isEmpty());
    }

    @Test
    void testBidirectionalStreams() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        // Request: client → server
        byte[] request = new byte[]{0x00, 0x04, 0x01, 0x02};
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, request);

        // Response: server → client
        byte[] response = new byte[]{0x00, 0x05, 0x03, 0x04, 0x05};
        reassembler.addSegment("10.0.0.2", 8080, "10.0.0.1", 12345, response);

        List<FramedMessage> messages = reassembler.extractMessages();

        assertEquals(2, messages.size());
        // First stream direction's messages come first
        assertFalse(messages.get(0).isResponse());
        assertTrue(messages.get(1).isResponse());
        assertEquals(0, messages.get(0).index());
        assertEquals(1, messages.get(1).index());
    }

    @Test
    void testEipLikeFraming() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(EIP_FRAMING, 44818);

        // EIP packet: 2 bytes command + 2 bytes LE length (0x0002 = 2 payload bytes) + 20 bytes header + 2 bytes data
        // Total = 2 + 24 = 26 bytes... wait, let me recalculate
        // EIP framing: length at offset 2, 2-byte LE, adjustment = 24
        // So total message = length_value + 24
        // If length_value = 2, total = 26

        byte[] packet = new byte[26];
        // Set length field at offset 2 (little-endian, value = 2)
        packet[2] = 0x02;
        packet[3] = 0x00;
        // Fill rest with recognizable data
        for (int i = 4; i < 26; i++) {
            packet[i] = (byte) (i & 0xFF);
        }

        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 44818, packet);

        List<FramedMessage> messages = reassembler.extractMessages();
        assertEquals(1, messages.size());
        assertEquals(26, messages.getFirst().rawBytes().length);
    }

    @Test
    void testModbusLikeFraming() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(MODBUS_FRAMING, 502);

        // Modbus TCP MBAP header: 2 txId + 2 protoId + 2 length (BE) + length bytes
        // Framing: offset=4, size=2, BE, adjustment=6
        // length_value = 5 → total = 5 + 6 = 11 bytes
        byte[] packet = new byte[11];
        packet[0] = 0x00; packet[1] = 0x01;  // transaction ID
        packet[2] = 0x00; packet[3] = 0x00;  // protocol ID
        packet[4] = 0x00; packet[5] = 0x05;  // length = 5
        packet[6] = 0x01;                     // unit ID
        packet[7] = 0x03;                     // function code
        packet[8] = 0x02;                     // byte count
        packet[9] = 0x00; packet[10] = (byte) 0xFF;  // data

        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 502, packet);

        List<FramedMessage> messages = reassembler.extractMessages();
        assertEquals(1, messages.size());
        assertArrayEquals(packet, messages.getFirst().rawBytes());
    }

    @Test
    void testExtractFromEmptyReassembler() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);
        List<FramedMessage> messages = reassembler.extractMessages();
        assertTrue(messages.isEmpty());
    }

    @Test
    void testGlobalIndexing() {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(SIMPLE_FRAMING, 8080);

        // Two messages from client
        byte[] twoMessages = new byte[]{
            0x00, 0x03, 0x01,  // msg 0
            0x00, 0x03, 0x02   // msg 1
        };
        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, twoMessages);

        // One message from server
        byte[] oneMessage = new byte[]{0x00, 0x04, 0x03, 0x04};
        reassembler.addSegment("10.0.0.2", 8080, "10.0.0.1", 12345, oneMessage);

        List<FramedMessage> messages = reassembler.extractMessages();

        assertEquals(3, messages.size());
        assertEquals(0, messages.get(0).index());
        assertEquals(1, messages.get(1).index());
        assertEquals(2, messages.get(2).index());
    }

    @Test
    void testFramingSpecValidation() {
        assertThrows(IllegalArgumentException.class, () ->
            new FramingSpec(-1, 2, true, 0), "Negative offset should be rejected");

        assertThrows(IllegalArgumentException.class, () ->
            new FramingSpec(0, 3, true, 0), "Unsupported length field size should be rejected");

        // Valid sizes should not throw
        assertDoesNotThrow(() -> new FramingSpec(0, 1, true, 0));
        assertDoesNotThrow(() -> new FramingSpec(0, 2, true, 0));
        assertDoesNotThrow(() -> new FramingSpec(0, 4, true, 0));
    }

    @Test
    void testFourByteLengthField() {
        // Framing with 4-byte length field at offset 4, little-endian, adjustment = 0
        FramingSpec spec = new FramingSpec(4, 4, false, 0);
        TcpStreamReassembler reassembler = new TcpStreamReassembler(spec, 4840);

        // Build a message: 4 header bytes + 4-byte LE length (value=12 → total=12) + 4 data bytes
        byte[] packet = new byte[12];
        // Length field at offset 4 (little-endian, value = 12)
        ByteBuffer.wrap(packet, 4, 4).order(ByteOrder.LITTLE_ENDIAN).putInt(12);

        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 4840, packet);

        List<FramedMessage> messages = reassembler.extractMessages();
        assertEquals(1, messages.size());
        assertEquals(12, messages.getFirst().rawBytes().length);
    }

    @Test
    void testOneByteLength() {
        FramingSpec spec = new FramingSpec(0, 1, true, 0);
        TcpStreamReassembler reassembler = new TcpStreamReassembler(spec, 8080);

        // Message: 1-byte length (5) + 4 data bytes = 5 total
        byte[] packet = new byte[]{0x05, 0x01, 0x02, 0x03, 0x04};

        reassembler.addSegment("10.0.0.1", 12345, "10.0.0.2", 8080, packet);

        List<FramedMessage> messages = reassembler.extractMessages();
        assertEquals(1, messages.size());
        assertArrayEquals(packet, messages.getFirst().rawBytes());
    }
}
