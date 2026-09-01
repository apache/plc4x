/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.eip.base;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.junit.jupiter.api.*;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.wait.strategy.Wait;
import org.testcontainers.images.builder.ImageFromDockerfile;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration test for the EtherNet/IP driver against a cpppo-based simulator.
 *
 * <p>cpppo (<a href="https://github.com/pjkundert/cpppo">github.com/pjkundert/cpppo</a>)
 * is the de-facto open-source EIP/CIP simulator: it exposes named tags over CIP
 * with the same semantics a Logix controller uses, so it covers the read/write
 * paths the driver actually executes.</p>
 *
 * <p>The container pre-populates the following tags:
 * <ul>
 *   <li>{@code hurz_BOOL : BOOL}</li>
 *   <li>{@code hurz_SINT : SINT}</li>
 *   <li>{@code hurz_INT  : INT}</li>
 *   <li>{@code hurz_DINT : DINT}</li>
 *   <li>{@code hurz_REAL : REAL}</li>
 *   <li>{@code hurz_DINT_ARR : DINT[10]}</li>
 * </ul>
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@Testcontainers(disabledWithoutDocker = true)
public class EipDockerIT {

    private static final int CPPPO_PORT = 44818;

    @SuppressWarnings("resource")
    private final GenericContainer<?> eipSimulator = new GenericContainer<>(
        new ImageFromDockerfile()
            .withFileFromClasspath("Dockerfile", "eip/Dockerfile"))
        .withExposedPorts(CPPPO_PORT)
        // cpppo doesn't print a startup banner (its --print flag only logs on
        // each served request), so wait on the TCP port instead of a log line.
        .waitingFor(Wait.forListeningPort()
            .withStartupTimeout(Duration.ofSeconds(120)));

    private String connectionUrl;

    @BeforeAll
    void startContainer() {
        eipSimulator.start();
        connectionUrl = "eip:tcp://" + eipSimulator.getHost() + ":"
            + eipSimulator.getMappedPort(CPPPO_PORT)
            + "?big-endian=false&force-unconnected-operation=true";
    }

    @AfterAll
    void stopContainer() {
        eipSimulator.stop();
    }

    /**
     * Smoke-test for the most common single-tag round-trip: write a DINT, then
     * read it back and verify the value survived the wire encoding.
     *
     * <p>What this exercises end-to-end:
     * <ol>
     *   <li>Connection setup — {@code ListServices}, {@code RegisterSession},
     *       sender-context correlation, and (because {@code force-unconnected-operation=true})
     *       skipping the Logix-only attribute probe.</li>
     *   <li>Outbound write — {@code CipWriteRequest} wrapped in a
     *       {@code CipUnconnectedRequest} inside {@code CipRRData}, including
     *       ANSI-extended-symbol path encoding for the tag name.</li>
     *   <li>Inbound parsing — {@code CipWriteResponse} round-tripped through
     *       the message codec back to a {@code PlcResponseCode}.</li>
     *   <li>The same machinery for the corresponding {@code CipReadRequest /
     *       Response}, including the little-endian DINT body extraction.</li>
     * </ol>
     *
     * <p>The write-before-read pattern makes the assertion deterministic even
     * if a previous test (or a future change to cpppo's defaults) left a
     * different value in the tag.
     */
    @Test
    void writeReadDint() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(connectionUrl)) {
            // Write a known value first so the read below is deterministic
            // regardless of cpppo's default initial state.
            PlcWriteRequest writeRequest = connection.writeRequestBuilder()
                .addTagAddress("dint", "hurz_DINT:DINT", 4711)
                .build();
            PlcWriteResponse writeResponse = writeRequest.execute().get();
            assertEquals(PlcResponseCode.OK, writeResponse.getResponseCode("dint"));

            // Read it straight back and confirm the value is preserved on
            // the wire (catches encoding/decoding mismatches for DINT).
            PlcReadRequest readRequest = connection.readRequestBuilder()
                .addTagAddress("dint", "hurz_DINT:DINT")
                .build();
            PlcReadResponse readResponse = readRequest.execute().get();
            assertEquals(PlcResponseCode.OK, readResponse.getResponseCode("dint"));
            assertEquals(4711, readResponse.getInteger("dint"));
        }
    }

    /**
     * Verifies multi-tag write / multi-tag read across every basic CIP scalar
     * type the driver currently supports: BOOL, SINT, INT, DINT, REAL.
     *
     * <p>Each tag is written to and then read back in a *single* request
     * containing all five tags. That covers two driver code paths the single-
     * tag test doesn't:
     * <ol>
     *   <li>The fan-out in {@code writeWithoutMessageRouter} /
     *       {@code readWithoutMessageRouter} that chains one CIP request per
     *       tag through the request throttle.</li>
     *   <li>The per-type encoding in {@code encodeValue(...)} and decoding in
     *       {@code parsePlcValue(...)} — sign-handling for the integer types
     *       and the IEEE-754 / little-endian byte order for REAL.</li>
     * </ol>
     *
     * <p>The values are chosen to be non-trivial (negative for signed types,
     * a value that uses the full DINT range, a float with multiple
     * significant digits) so a byte-swap or sign-extension bug would surface
     * as a clear mismatch rather than coincidentally pass.
     */
    @Test
    void writeReadAllBasicTypes() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(connectionUrl)) {
            // Build one request that touches every basic scalar type at once.
            PlcWriteRequest writeRequest = connection.writeRequestBuilder()
                .addTagAddress("b",    "hurz_BOOL:BOOL", true)
                .addTagAddress("s8",   "hurz_SINT:SINT", (byte) -42)
                .addTagAddress("s16",  "hurz_INT:INT",   (short) -2424)
                .addTagAddress("s32",  "hurz_DINT:DINT", -242442424)
                .addTagAddress("f32",  "hurz_REAL:REAL", 3.14159f)
                .build();
            PlcWriteResponse writeResponse = writeRequest.execute().get();
            // Every tag must have committed; if any one of them failed the
            // tag-by-tag fan-out logic isn't doing its job.
            writeResponse.getTagNames().forEach(tn ->
                assertEquals(PlcResponseCode.OK, writeResponse.getResponseCode(tn), tn));

            // Now read the same five tags back, again in a single request.
            PlcReadRequest readRequest = connection.readRequestBuilder()
                .addTagAddress("b",   "hurz_BOOL:BOOL")
                .addTagAddress("s8",  "hurz_SINT:SINT")
                .addTagAddress("s16", "hurz_INT:INT")
                .addTagAddress("s32", "hurz_DINT:DINT")
                .addTagAddress("f32", "hurz_REAL:REAL")
                .build();
            PlcReadResponse readResponse = readRequest.execute().get();
            readResponse.getTagNames().forEach(tn ->
                assertEquals(PlcResponseCode.OK, readResponse.getResponseCode(tn), tn));

            // Value-level assertions per type. Any byte-order or
            // sign-extension regression in encodeValue / parsePlcValue
            // surfaces here.
            assertTrue(readResponse.getBoolean("b"));
            assertEquals((byte) -42, readResponse.getByte("s8"));
            assertEquals((short) -2424, readResponse.getShort("s16"));
            assertEquals(-242442424, readResponse.getInteger("s32"));
            assertEquals(3.14159f, readResponse.getFloat("f32"), 1e-5);
        }
    }

    /**
     * Exercises array-element addressing — i.e. CIP "member ID" path segments
     * appended to the ANSI extended-symbol segment.
     *
     * <p>The address {@code hurz_DINT_ARR[3]:DINT} forces
     * {@code EipTcpConnection.toAnsi(...)} to emit:
     * <pre>
     *   AnsiExtendedSymbolSegment("hurz_DINT_ARR")  +  LogicalSegment(MemberID(3))
     * </pre>
     * which is a different code path from a plain scalar tag — and is the
     * path Rockwell controllers expect for any structured-/array-member
     * access. cpppo accepts the same encoding, so a round-trip on a single
     * array element validates that the path-segment serialisation matches
     * what real Logix devices expect.
     *
     * <p>The {@code 0xCAFEBABE} payload is just a recognisable bit pattern
     * — easy to spot in pcap traces if this test fails.
     */
    @Test
    void writeReadArrayElement() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(connectionUrl)) {
            // Write a single element of a DINT[10] array. The "[3]" suffix
            // becomes a MemberID logical segment in the CIP path.
            PlcWriteRequest writeRequest = connection.writeRequestBuilder()
                .addTagAddress("arrElem", "hurz_DINT_ARR[3]:DINT", 0xCAFEBABE)
                .build();
            assertEquals(PlcResponseCode.OK,
                writeRequest.execute().get().getResponseCode("arrElem"));

            // Read the same element back. If member-id encoding were wrong
            // the read would either fail outright or return a value from a
            // different array slot — both caught by the assertion below.
            PlcReadRequest readRequest = connection.readRequestBuilder()
                .addTagAddress("arrElem", "hurz_DINT_ARR[3]:DINT")
                .build();
            PlcReadResponse readResponse = readRequest.execute().get();
            assertEquals(PlcResponseCode.OK, readResponse.getResponseCode("arrElem"));
            assertEquals(0xCAFEBABE, readResponse.getInteger("arrElem"));
        }
    }

    /**
     * Reading several array elements at once - see GH-1008. The driver used to ask the device
     * for a single element while decoding as many as the tag declared, so this returned either
     * a wrong value or an IndexOutOfBoundsException.
     *
     * <p>The element count has to reach the wire in the {@code CipReadRequest}, which is what
     * separates this from {@link #writeReadArrayElement()} - that one reads a single element
     * and passes either way.
     */
    @Test
    void readMultipleArrayElements() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(connectionUrl)) {
            // Seed the first four elements with recognisable values.
            PlcWriteRequest.Builder writeBuilder = connection.writeRequestBuilder();
            for (int i = 0; i < 4; i++) {
                writeBuilder.addTagAddress("e" + i, "hurz_DINT_ARR[" + i + "]:DINT", 0x1000 + i);
            }
            PlcWriteResponse writeResponse = writeBuilder.build().execute().get();
            for (int i = 0; i < 4; i++) {
                assertEquals(PlcResponseCode.OK, writeResponse.getResponseCode("e" + i));
            }

            // Read all four back with a single array read.
            PlcReadRequest readRequest = connection.readRequestBuilder()
                .addTagAddress("arr", "hurz_DINT_ARR[0..3]:DINT")
                .build();
            PlcReadResponse readResponse = readRequest.execute().get();

            assertEquals(PlcResponseCode.OK, readResponse.getResponseCode("arr"));
            PlcValue value = readResponse.getPlcValue("arr");
            assertEquals(4, value.getLength());
            for (int i = 0; i < 4; i++) {
                assertEquals(0x1000 + i, value.getIndex(i).getInt(), "element " + i);
            }
        }
    }

}
