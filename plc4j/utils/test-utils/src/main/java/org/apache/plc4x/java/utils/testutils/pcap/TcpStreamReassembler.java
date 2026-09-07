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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Protocol-agnostic TCP stream reassembler that extracts individual protocol messages from
 * raw TCP segment payloads.
 *
 * <p>pcap4j does not perform TCP stream reassembly. This class collects TCP payloads per
 * stream direction, concatenates them into per-direction byte buffers, and uses a
 * {@link FramingSpec} to extract individual protocol messages from each buffer.
 *
 * <p>The framing logic reads a length field at a fixed offset in the message header,
 * then computes the total message length as {@code lengthFieldValue + lengthAdjustment}.
 *
 * <p>Direction is inferred from the protocol port: packets sent TO the port are requests,
 * packets sent FROM the port are responses.
 */
public class TcpStreamReassembler {

    private static final Logger LOGGER = LoggerFactory.getLogger(TcpStreamReassembler.class);

    private final FramingSpec framingSpec;
    private final int protocolPort;

    // Keyed by "srcIP:srcPort->dstIP:dstPort", stores accumulated payload bytes
    private final Map<String, ByteArrayOutputStream> streamBuffers = new LinkedHashMap<>();
    // Tracks whether each stream direction is a response (srcPort == protocolPort)
    private final Map<String, Boolean> streamDirections = new LinkedHashMap<>();

    /**
     * Creates a reassembler for the given framing specification and protocol port.
     *
     * @param framingSpec  how to detect message boundaries in the TCP stream
     * @param protocolPort the well-known port of the protocol (used for direction inference)
     */
    public TcpStreamReassembler(FramingSpec framingSpec, int protocolPort) {
        this.framingSpec = framingSpec;
        this.protocolPort = protocolPort;
    }

    /**
     * Adds a TCP segment payload to the appropriate stream buffer.
     *
     * @param srcIp   source IP address
     * @param srcPort source port
     * @param dstIp   destination IP address
     * @param dstPort destination port
     * @param payload the TCP segment payload bytes
     */
    public void addSegment(String srcIp, int srcPort, String dstIp, int dstPort, byte[] payload) {
        if (payload == null || payload.length == 0) {
            return;
        }
        String key = srcIp + ":" + srcPort + "->" + dstIp + ":" + dstPort;
        streamBuffers.computeIfAbsent(key, k -> new ByteArrayOutputStream());
        try {
            streamBuffers.get(key).write(payload);
        } catch (IOException e) {
            // ByteArrayOutputStream.write() never throws IOException
            throw new RuntimeException("Unexpected IOException from ByteArrayOutputStream", e);
        }
        // srcPort == protocolPort means the server is sending → response
        streamDirections.putIfAbsent(key, srcPort == protocolPort);
    }

    /**
     * Extracts all complete protocol messages from the accumulated stream buffers.
     *
     * <p>Messages are returned in the order they appear across all stream directions,
     * interleaved by stream to approximate the original conversation order.
     * Each message is assigned a global sequential index.
     *
     * @return list of framed messages with raw bytes and direction metadata
     */
    public List<FramedMessage> extractMessages() {
        List<FramedMessage> allMessages = new ArrayList<>();
        int globalIndex = 0;

        for (Map.Entry<String, ByteArrayOutputStream> entry : streamBuffers.entrySet()) {
            String streamKey = entry.getKey();
            byte[] streamData = entry.getValue().toByteArray();
            boolean isResponse = streamDirections.getOrDefault(streamKey, false);

            List<byte[]> messages = extractFrames(streamData, streamKey);
            for (byte[] messageBytes : messages) {
                allMessages.add(new FramedMessage(messageBytes, isResponse, globalIndex++));
            }
        }

        return allMessages;
    }

    /**
     * Extracts individual framed messages from a contiguous byte buffer using the configured
     * {@link FramingSpec}.
     *
     * @param data      the concatenated TCP payload bytes for one stream direction
     * @param streamKey identifier for logging (e.g. "10.0.0.1:44818->10.0.0.2:12345")
     * @return list of complete message byte arrays
     */
    private List<byte[]> extractFrames(byte[] data, String streamKey) {
        List<byte[]> frames = new ArrayList<>();
        int offset = 0;
        int minHeaderSize = framingSpec.lengthFieldOffset() + framingSpec.lengthFieldSize();

        while (offset + minHeaderSize <= data.length) {
            // Read the length field value
            int lengthFieldValue = readLengthField(data, offset);
            if (lengthFieldValue < 0) {
                LOGGER.debug("Negative length field value {} at offset {} in stream {}, stopping extraction",
                    lengthFieldValue, offset, streamKey);
                break;
            }

            int totalMessageLength = lengthFieldValue + framingSpec.lengthAdjustment();
            if (totalMessageLength <= 0) {
                LOGGER.debug("Non-positive total message length {} at offset {} in stream {}, stopping extraction",
                    totalMessageLength, offset, streamKey);
                break;
            }

            if (offset + totalMessageLength > data.length) {
                // Incomplete message at end of stream — common when captures end mid-stream
                LOGGER.debug("Incomplete message at offset {} in stream {}: need {} bytes but only {} remain",
                    offset, streamKey, totalMessageLength, data.length - offset);
                break;
            }

            byte[] frame = new byte[totalMessageLength];
            System.arraycopy(data, offset, frame, 0, totalMessageLength);
            frames.add(frame);
            offset += totalMessageLength;
        }

        if (offset < data.length) {
            LOGGER.debug("Discarded {} trailing bytes in stream {} (likely incomplete message at end of capture)",
                data.length - offset, streamKey);
        }

        return frames;
    }

    /**
     * Reads the length field value from the given position in the data buffer.
     *
     * @param data   the byte buffer
     * @param offset the start of the current message
     * @return the length field value as an unsigned integer
     */
    private int readLengthField(byte[] data, int offset) {
        int fieldOffset = offset + framingSpec.lengthFieldOffset();
        ByteBuffer buffer = ByteBuffer.wrap(data, fieldOffset, framingSpec.lengthFieldSize());
        buffer.order(framingSpec.bigEndian() ? ByteOrder.BIG_ENDIAN : ByteOrder.LITTLE_ENDIAN);

        return switch (framingSpec.lengthFieldSize()) {
            case 1 -> buffer.get() & 0xFF;
            case 2 -> buffer.getShort() & 0xFFFF;
            case 4 -> buffer.getInt();
            default -> throw new IllegalStateException("Unsupported length field size: " + framingSpec.lengthFieldSize());
        };
    }
}
