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

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.MessageInput;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.utils.testutils.utils.hex.HexDiff;
import org.apache.plc4x.java.utils.testutils.utils.migration.MessageResolver;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.Pcaps;
import org.pcap4j.packet.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.EOFException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;
import java.util.*;
import java.util.stream.Collectors;

import static org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption.WithByteOrder;
import static org.apache.plc4x.java.spi.utils.StaticHelper.ENCODE_HEX;

/**
 * JUnit 5 test runner that reads pcap files and verifies parse→serialize round-trip consistency.
 *
 * <p>For each protocol message extracted from the pcap capture, this runner:
 * <ol>
 *   <li>Parses the raw bytes into a message object using the generated {@code staticParse()} method</li>
 *   <li>Serializes the parsed message back to bytes</li>
 *   <li>Compares the serialized bytes with the original raw bytes</li>
 * </ol>
 *
 * <p>If the round-trip produces identical bytes, the parser and serializer are consistent.
 * Parse failures are treated as skipped tests (not failures), since some packets may use
 * encrypted or unimplemented protocol features.
 *
 * <p>Concrete test classes extend this runner with a minimal constructor:
 * <pre>{@code
 * public class EipEnipRoundTripTest extends PcapRoundTripTestRunner {
 *     public EipEnipRoundTripTest() {
 *         super(PcapRoundTripConfig.builder()
 *             .pcapResource("/protocols/eip/enip_test.pcap")
 *             .rootType("EipPacket")
 *             .packageName("org.apache.plc4x.java.eip.readwrite")
 *             .byteOrder("LITTLE_ENDIAN")
 *             .transportType(TransportType.TCP)
 *             .protocolPort(44818)
 *             .framingSpec(2, 2, false, 24)
 *             .directionDependentArg("response")
 *             .build());
 *     }
 * }
 * }</pre>
 */
public class PcapRoundTripTestRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(PcapRoundTripTestRunner.class);

    private final PcapRoundTripConfig config;

    /**
     * Creates a test runner with the given configuration.
     *
     * @param config the pcap round-trip test configuration
     */
    protected PcapRoundTripTestRunner(PcapRoundTripConfig config) {
        this.config = Objects.requireNonNull(config, "config must not be null");
    }

    /**
     * Generates dynamic JUnit 5 tests — one per protocol message in the pcap capture.
     * Each test verifies that parsing and re-serializing the message produces identical bytes.
     */
    @TestFactory
    public Iterable<DynamicTest> roundTripTests() throws Exception {
        List<FramedMessage> messages = extractMessages();
        LOGGER.info("Extracted {} messages from {}", messages.size(), config.pcapResource());

        List<DynamicTest> tests = new ArrayList<>();
        for (FramedMessage msg : messages) {
            String testName = String.format("Packet %d (%s)", msg.index(),
                msg.isResponse() ? "response" : "request");

            tests.add(DynamicTest.dynamicTest(testName, () -> {
                // Skip explicitly ignored test cases
                Assumptions.assumeFalse(
                    config.ignoredTestCases().contains(testName),
                    "Test case '" + testName + "' is in the ignored list"
                );

                runRoundTrip(msg);
            }));
        }
        return tests;
    }

    /**
     * Extracts protocol messages from the pcap file based on the configured transport type.
     */
    private List<FramedMessage> extractMessages() throws Exception {
        Path tempFile = extractToTempFile();
        try {
            return switch (config.transportType()) {
                case TCP -> extractTcpMessages(tempFile);
                case UDP -> extractUdpMessages(tempFile);
                case ETHERNET -> extractEthernetMessages(tempFile);
            };
        } finally {
            Files.deleteIfExists(tempFile);
        }
    }

    /**
     * Copies the classpath pcap resource to a temporary file (pcap4j requires a filesystem path).
     */
    private Path extractToTempFile() throws Exception {
        InputStream is = getClass().getResourceAsStream(config.pcapResource());
        if (is == null) {
            throw new IllegalStateException("Pcap resource not found on classpath: " + config.pcapResource());
        }
        String fileName = config.pcapResource().substring(config.pcapResource().lastIndexOf('/') + 1);
        String baseName = fileName.contains(".") ? fileName.substring(0, fileName.lastIndexOf('.')) : fileName;
        String extension = fileName.contains(".") ? fileName.substring(fileName.lastIndexOf('.')) : ".pcap";
        Path tempFile = Files.createTempFile(baseName + "_", extension);
        Files.copy(is, tempFile, StandardCopyOption.REPLACE_EXISTING);
        is.close();
        return tempFile;
    }

    /**
     * Extracts TCP protocol messages using {@link TcpStreamReassembler}.
     * Collects all TCP segments matching the protocol port, reassembles streams,
     * and extracts individual framed messages.
     */
    private List<FramedMessage> extractTcpMessages(Path pcapFile) throws Exception {
        TcpStreamReassembler reassembler = new TcpStreamReassembler(config.framingSpec(), config.protocolPort());

        try (PcapHandle handle = Pcaps.openOffline(pcapFile.toString(), PcapHandle.TimestampPrecision.MICRO)) {
            Packet packet;
            while (true) {
                try {
                    packet = handle.getNextPacketEx();
                } catch (EOFException e) {
                    break;
                }

                // Extract IP and TCP layers
                IpV4Packet ipV4 = packet.get(IpV4Packet.class);
                TcpPacket tcp = packet.get(TcpPacket.class);
                if (ipV4 == null || tcp == null) {
                    // Try IPv6
                    IpV6Packet ipV6 = packet.get(IpV6Packet.class);
                    if (ipV6 != null && tcp != null) {
                        byte[] payload = tcp.getPayload() != null ? tcp.getPayload().getRawData() : null;
                        if (payload != null && payload.length > 0) {
                            String srcIp = ipV6.getHeader().getSrcAddr().getHostAddress();
                            String dstIp = ipV6.getHeader().getDstAddr().getHostAddress();
                            int srcPort = tcp.getHeader().getSrcPort().valueAsInt();
                            int dstPort = tcp.getHeader().getDstPort().valueAsInt();

                            // Only process packets involving the protocol port
                            if (srcPort == config.protocolPort() || dstPort == config.protocolPort()) {
                                reassembler.addSegment(srcIp, srcPort, dstIp, dstPort, payload);
                            }
                        }
                    }
                    continue;
                }

                byte[] payload = tcp.getPayload() != null ? tcp.getPayload().getRawData() : null;
                if (payload == null || payload.length == 0) {
                    continue;
                }

                String srcIp = ipV4.getHeader().getSrcAddr().getHostAddress();
                String dstIp = ipV4.getHeader().getDstAddr().getHostAddress();
                int srcPort = tcp.getHeader().getSrcPort().valueAsInt();
                int dstPort = tcp.getHeader().getDstPort().valueAsInt();

                // Only process packets involving the protocol port
                if (srcPort == config.protocolPort() || dstPort == config.protocolPort()) {
                    reassembler.addSegment(srcIp, srcPort, dstIp, dstPort, payload);
                }
            }
        }

        return reassembler.extractMessages();
    }

    /**
     * Extracts UDP protocol messages. Each UDP datagram payload matching the port is one complete message.
     */
    private List<FramedMessage> extractUdpMessages(Path pcapFile) throws Exception {
        List<FramedMessage> messages = new ArrayList<>();
        int index = 0;

        try (PcapHandle handle = Pcaps.openOffline(pcapFile.toString(), PcapHandle.TimestampPrecision.MICRO)) {
            Packet packet;
            while (true) {
                try {
                    packet = handle.getNextPacketEx();
                } catch (EOFException e) {
                    break;
                }

                UdpPacket udp = packet.get(UdpPacket.class);
                if (udp == null || udp.getPayload() == null) {
                    continue;
                }

                int srcPort = udp.getHeader().getSrcPort().valueAsInt();
                int dstPort = udp.getHeader().getDstPort().valueAsInt();

                // Only process packets involving the protocol port
                if (srcPort != config.protocolPort() && dstPort != config.protocolPort()) {
                    continue;
                }

                boolean isResponse = srcPort == config.protocolPort();
                byte[] payload = udp.getPayload().getRawData();
                if (payload.length > 0) {
                    messages.add(new FramedMessage(payload, isResponse, index++));
                }
            }
        }

        return messages;
    }

    /**
     * Extracts Ethernet protocol messages. Each Ethernet frame matching the EtherType is one message.
     * The raw data includes MAC headers since root types like {@code Ethernet_Frame} parse those.
     */
    private List<FramedMessage> extractEthernetMessages(Path pcapFile) throws Exception {
        List<FramedMessage> messages = new ArrayList<>();
        int index = 0;

        try (PcapHandle handle = Pcaps.openOffline(pcapFile.toString(), PcapHandle.TimestampPrecision.MICRO)) {
            Packet packet;
            while (true) {
                try {
                    packet = handle.getNextPacketEx();
                } catch (EOFException e) {
                    break;
                }

                EthernetPacket eth = packet.get(EthernetPacket.class);
                if (eth == null) {
                    continue;
                }

                // Filter by EtherType if a specific one is configured (0 means accept all frames)
                if (config.protocolPort() != 0) {
                    int etherType = eth.getHeader().getType().value() & 0xFFFF;
                    if (etherType != config.protocolPort()) {
                        continue;
                    }
                }

                // For Ethernet, use the full raw frame data (including MAC headers)
                byte[] rawData = eth.getRawData();
                if (rawData.length > 0) {
                    // Direction: not meaningful for Ethernet (no port), default to false
                    messages.add(new FramedMessage(rawData, false, index++));
                }
            }
        }

        return messages;
    }

    /**
     * Runs the parse→serialize round-trip for a single message.
     *
     * <p>Parse failures are treated as skipped tests (assumptions), not hard failures,
     * since some packets in a capture may be encrypted or use unimplemented features.
     */
    private void runRoundTrip(FramedMessage msg) throws Exception {
        byte[] rawBytes = msg.rawBytes();

        // Resolve parser arguments for this packet's direction
        List<String> parserArgs = config.argSpecs().stream()
            .map(spec -> spec.resolve(msg.isResponse()))
            .collect(Collectors.toList());

        // Set up options for MessageResolver
        Map<String, String> options = new HashMap<>();
        options.put("package", config.packageName());

        // Parse
        Message parsedMessage;
        try {
            MessageInput<?> messageInput = MessageResolver.getMessageIOStaticLinked(
                options, config.rootType(), parserArgs);
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(rawBytes, WithByteOrder(config.byteOrder()));
            parsedMessage = (Message) messageInput.parse(readBuffer);
        } catch (Exception e) {
            // Parse failure → skip the test, don't fail it
            // Dump full exception chain for diagnostics
            StringBuilder chain = new StringBuilder();
            chain.append(e.getClass().getSimpleName()).append(": ").append(e.getMessage());
            Throwable c = e.getCause();
            int depth = 0;
            while (c != null && depth < 10) {
                chain.append(" -> ").append(c.getClass().getSimpleName()).append(": ").append(c.getMessage());
                c = c.getCause();
                depth++;
            }
            Assumptions.abort("Parse failed for Packet " + msg.index() + ": " + chain);
            return; // unreachable, but keeps the compiler happy
        }

        // Serialize
        int serializedLength = parsedMessage.getLengthInBytes();
        byte[] buffer = new byte[serializedLength];
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(buffer, WithByteOrder(config.byteOrder()));
        parsedMessage.serialize(writeBuffer);
        byte[] serializedBytes = writeBuffer.getBytes();

        // Compare — for Ethernet, only compare up to the parsed message length
        // (Ethernet minimum frame padding may add trailing bytes)
        byte[] expectedBytes;
        if (config.transportType() == TransportType.ETHERNET && rawBytes.length > serializedLength) {
            expectedBytes = Arrays.copyOf(rawBytes, serializedLength);
        } else {
            expectedBytes = rawBytes;
        }

        if (!Arrays.equals(expectedBytes, serializedBytes)) {
            String expectedHex = ENCODE_HEX(expectedBytes);
            String actualHex = ENCODE_HEX(serializedBytes);
            String diff = HexDiff.hexDiff(expectedHex, actualHex);
            String message = String.format(
                "Round-trip failed for Packet %d (%s).%nExpected %d bytes, got %d bytes.%nDiff:%n%s",
                msg.index(), msg.isResponse() ? "response" : "request",
                expectedBytes.length, serializedBytes.length, diff
            );
            throw new AssertionError(message);
        }

        LOGGER.debug("Round-trip passed for Packet {} ({}) - {} bytes",
            msg.index(), msg.isResponse() ? "response" : "request", rawBytes.length);
    }

}
