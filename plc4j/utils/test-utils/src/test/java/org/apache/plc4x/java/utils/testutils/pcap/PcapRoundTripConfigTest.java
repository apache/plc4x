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

import static org.junit.jupiter.api.Assertions.*;

class PcapRoundTripConfigTest {

    @Test
    void testBuilderCreatesValidConfig() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/protocols/eip/enip_test.pcap")
            .rootType("EipPacket")
            .packageName("org.apache.plc4x.java.eip.readwrite")
            .byteOrder("LITTLE_ENDIAN")
            .transportType(TransportType.TCP)
            .protocolPort(44818)
            .framingSpec(2, 2, false, 24)
            .directionDependentArg("response")
            .build();

        assertEquals("/protocols/eip/enip_test.pcap", config.pcapResource());
        assertEquals("EipPacket", config.rootType());
        assertEquals("org.apache.plc4x.java.eip.readwrite", config.packageName());
        assertEquals("LITTLE_ENDIAN", config.byteOrder());
        assertEquals(TransportType.TCP, config.transportType());
        assertEquals(44818, config.protocolPort());
        assertNotNull(config.framingSpec());
        assertEquals(1, config.argSpecs().size());
        assertTrue(config.ignoredTestCases().isEmpty());
    }

    @Test
    void testBuilderDefaultByteOrder() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .build();

        assertEquals("BIG_ENDIAN", config.byteOrder());
    }

    @Test
    void testBuilderWithFixedArg() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .fixedArg("driverType", "tcp")
            .build();

        assertEquals(1, config.argSpecs().size());
        ParserArgSpec spec = config.argSpecs().getFirst();
        assertInstanceOf(ParserArgSpec.FixedArg.class, spec);
        assertEquals("driverType", spec.name());
        assertEquals("tcp", spec.resolve(true));
        assertEquals("tcp", spec.resolve(false));
    }

    @Test
    void testBuilderWithDirectionDependentArg() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .directionDependentArg("response")
            .build();

        assertEquals(1, config.argSpecs().size());
        ParserArgSpec spec = config.argSpecs().getFirst();
        assertInstanceOf(ParserArgSpec.DirectionDependentArg.class, spec);
        assertEquals("response", spec.name());
        assertEquals("true", spec.resolve(true));
        assertEquals("false", spec.resolve(false));
    }

    @Test
    void testBuilderWithMultipleArgs() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .directionDependentArg("response")
            .fixedArg("driverType", "tcp")
            .build();

        assertEquals(2, config.argSpecs().size());
        assertInstanceOf(ParserArgSpec.DirectionDependentArg.class, config.argSpecs().get(0));
        assertInstanceOf(ParserArgSpec.FixedArg.class, config.argSpecs().get(1));
    }

    @Test
    void testBuilderWithIgnoredTestCases() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .ignoreTestCase("Packet 5 (response)")
            .ignoreTestCase("Packet 12 (request)")
            .build();

        assertEquals(2, config.ignoredTestCases().size());
        assertTrue(config.ignoredTestCases().contains("Packet 5 (response)"));
        assertTrue(config.ignoredTestCases().contains("Packet 12 (request)"));
    }

    @Test
    void testTcpTransportRequiresFramingSpec() {
        PcapRoundTripConfig.Builder builder = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.TCP)
            .protocolPort(44818);

        assertThrows(IllegalStateException.class, builder::build,
            "TCP transport should require a framingSpec");
    }

    @Test
    void testUdpTransportDoesNotRequireFramingSpec() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(47808)
            .build();

        assertNull(config.framingSpec());
    }

    @Test
    void testEthernetTransportDoesNotRequireFramingSpec() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("Ethernet_Frame")
            .packageName("com.test")
            .transportType(TransportType.ETHERNET)
            .protocolPort(0x8892)
            .build();

        assertNull(config.framingSpec());
    }

    @Test
    void testBuilderRequiresPcapResource() {
        PcapRoundTripConfig.Builder builder = PcapRoundTripConfig.builder()
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502);

        assertThrows(NullPointerException.class, builder::build);
    }

    @Test
    void testBuilderRequiresRootType() {
        PcapRoundTripConfig.Builder builder = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502);

        assertThrows(NullPointerException.class, builder::build);
    }

    @Test
    void testBuilderRequiresPackageName() {
        PcapRoundTripConfig.Builder builder = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .transportType(TransportType.UDP)
            .protocolPort(502);

        assertThrows(NullPointerException.class, builder::build);
    }

    @Test
    void testBuilderRequiresTransportType() {
        PcapRoundTripConfig.Builder builder = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .protocolPort(502);

        assertThrows(NullPointerException.class, builder::build);
    }

    @Test
    void testArgSpecsListIsUnmodifiable() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .directionDependentArg("response")
            .build();

        assertThrows(UnsupportedOperationException.class, () ->
            config.argSpecs().add(new ParserArgSpec.FixedArg("extra", "value")));
    }

    @Test
    void testIgnoredTestCasesSetIsUnmodifiable() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.UDP)
            .protocolPort(502)
            .build();

        assertThrows(UnsupportedOperationException.class, () ->
            config.ignoredTestCases().add("Packet 1"));
    }

    @Test
    void testFramingSpecWithConvenienceMethod() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.TCP)
            .protocolPort(502)
            .framingSpec(4, 2, true, 6)
            .build();

        FramingSpec spec = config.framingSpec();
        assertEquals(4, spec.lengthFieldOffset());
        assertEquals(2, spec.lengthFieldSize());
        assertTrue(spec.bigEndian());
        assertEquals(6, spec.lengthAdjustment());
    }

    @Test
    void testFramingSpecWithRecordConstructor() {
        PcapRoundTripConfig config = PcapRoundTripConfig.builder()
            .pcapResource("/test.pcap")
            .rootType("TestPacket")
            .packageName("com.test")
            .transportType(TransportType.TCP)
            .protocolPort(502)
            .framingSpec(new FramingSpec(4, 2, true, 6))
            .build();

        assertNotNull(config.framingSpec());
    }
}
