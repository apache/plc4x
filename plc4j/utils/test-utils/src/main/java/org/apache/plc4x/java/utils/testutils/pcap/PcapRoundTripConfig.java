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

import java.util.*;

/**
 * Immutable configuration for a PCAP round-trip test.
 *
 * <p>Specifies which pcap file to read, how to extract protocol messages from it,
 * and how to parse/serialize them for round-trip verification.
 *
 * <p>Use {@link #builder()} to create instances.
 */
public class PcapRoundTripConfig {

    private final String pcapResource;
    private final String rootType;
    private final String packageName;
    private final String byteOrder;
    private final TransportType transportType;
    private final int protocolPort;
    private final FramingSpec framingSpec;
    private final List<ParserArgSpec> argSpecs;
    private final Set<String> ignoredTestCases;

    private PcapRoundTripConfig(Builder builder) {
        this.pcapResource = Objects.requireNonNull(builder.pcapResource, "pcapResource is required");
        this.rootType = Objects.requireNonNull(builder.rootType, "rootType is required");
        this.packageName = Objects.requireNonNull(builder.packageName, "packageName is required");
        this.byteOrder = Objects.requireNonNull(builder.byteOrder, "byteOrder is required");
        this.transportType = Objects.requireNonNull(builder.transportType, "transportType is required");
        this.protocolPort = builder.protocolPort;
        this.framingSpec = builder.framingSpec;
        this.argSpecs = Collections.unmodifiableList(new ArrayList<>(builder.argSpecs));
        this.ignoredTestCases = Collections.unmodifiableSet(new HashSet<>(builder.ignoredTestCases));

        // TCP transport requires a framing spec for message boundary detection
        if (transportType == TransportType.TCP && framingSpec == null) {
            throw new IllegalStateException("TCP transport requires a framingSpec");
        }
    }

    /** Classpath resource path to the pcap file, e.g. {@code "/protocols/eip/enip_test.pcap"}. */
    public String pcapResource() { return pcapResource; }

    /** Root message type name, e.g. {@code "EipPacket"}. */
    public String rootType() { return rootType; }

    /** Fully qualified package of the generated readwrite classes. */
    public String packageName() { return packageName; }

    /** Byte order for the protocol: {@code "BIG_ENDIAN"} or {@code "LITTLE_ENDIAN"}. */
    public String byteOrder() { return byteOrder; }

    /** Transport type determining how messages are extracted from the pcap. */
    public TransportType transportType() { return transportType; }

    /** Protocol port (TCP/UDP) or EtherType (ETHERNET) used for packet filtering and direction inference. */
    public int protocolPort() { return protocolPort; }

    /** TCP message framing specification; null for UDP and ETHERNET transports. */
    public FramingSpec framingSpec() { return framingSpec; }

    /** Ordered list of parser argument specifications. */
    public List<ParserArgSpec> argSpecs() { return argSpecs; }

    /** Test case names to skip (marked as assumptions, not failures). */
    public Set<String> ignoredTestCases() { return ignoredTestCases; }

    /**
     * Creates a new builder for constructing a {@link PcapRoundTripConfig}.
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Builder for {@link PcapRoundTripConfig}.
     */
    public static class Builder {
        private String pcapResource;
        private String rootType;
        private String packageName;
        private String byteOrder = "BIG_ENDIAN";
        private TransportType transportType;
        private int protocolPort;
        private FramingSpec framingSpec;
        private final List<ParserArgSpec> argSpecs = new ArrayList<>();
        private final Set<String> ignoredTestCases = new HashSet<>();

        private Builder() {}

        /** Sets the classpath resource path to the pcap file. */
        public Builder pcapResource(String pcapResource) {
            this.pcapResource = pcapResource;
            return this;
        }

        /** Sets the root message type name used for parsing. */
        public Builder rootType(String rootType) {
            this.rootType = rootType;
            return this;
        }

        /** Sets the fully qualified package name of the generated readwrite classes. */
        public Builder packageName(String packageName) {
            this.packageName = packageName;
            return this;
        }

        /** Sets the byte order. Defaults to {@code "BIG_ENDIAN"} if not called. */
        public Builder byteOrder(String byteOrder) {
            this.byteOrder = byteOrder;
            return this;
        }

        /** Sets the transport type (TCP, UDP, or ETHERNET). */
        public Builder transportType(TransportType transportType) {
            this.transportType = transportType;
            return this;
        }

        /** Sets the protocol port for TCP/UDP filtering, or EtherType for ETHERNET. */
        public Builder protocolPort(int protocolPort) {
            this.protocolPort = protocolPort;
            return this;
        }

        /** Sets the TCP framing specification for message boundary detection. */
        public Builder framingSpec(FramingSpec framingSpec) {
            this.framingSpec = framingSpec;
            return this;
        }

        /**
         * Convenience method to set the TCP framing specification with individual parameters.
         *
         * @param lengthFieldOffset byte offset to the length field
         * @param lengthFieldSize   size of the length field in bytes
         * @param bigEndian         whether the length field is big-endian
         * @param lengthAdjustment  added to length field value to get total message size
         */
        public Builder framingSpec(int lengthFieldOffset, int lengthFieldSize, boolean bigEndian, int lengthAdjustment) {
            this.framingSpec = new FramingSpec(lengthFieldOffset, lengthFieldSize, bigEndian, lengthAdjustment);
            return this;
        }

        /** Adds a fixed-value parser argument. */
        public Builder fixedArg(String name, String value) {
            this.argSpecs.add(new ParserArgSpec.FixedArg(name, value));
            return this;
        }

        /** Adds a direction-dependent boolean parser argument (true for responses, false for requests). */
        public Builder directionDependentArg(String name) {
            this.argSpecs.add(new ParserArgSpec.DirectionDependentArg(name));
            return this;
        }

        /** Adds a test case name to the ignored set. */
        public Builder ignoreTestCase(String testCaseName) {
            this.ignoredTestCases.add(testCaseName);
            return this;
        }

        /**
         * Builds the immutable configuration.
         *
         * @throws NullPointerException  if required fields are null
         * @throws IllegalStateException if TCP transport has no framing spec
         */
        public PcapRoundTripConfig build() {
            return new PcapRoundTripConfig(this);
        }
    }
}
