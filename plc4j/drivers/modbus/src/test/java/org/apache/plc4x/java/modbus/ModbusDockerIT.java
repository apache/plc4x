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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcNullAuthentication;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.testutils.manual.BasicPlcTest;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.extension.RegisterExtension;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.wait.strategy.Wait;
import org.testcontainers.images.builder.ImageFromDockerfile;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.time.Duration;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration test for all Modbus driver variants (TCP, RTU, ASCII) across
 * all supported transports (TCP, UDP, TLS) using a pymodbus-based test container.
 *
 * <p>The test container runs pymodbus servers pre-populated with known values
 * across all 4 Modbus address spaces: coils, discrete inputs, input registers,
 * and holding registers. Each driver/transport combination gets its own server
 * port inside the container.</p>
 *
 * <p>Tests are organized by driver type first, then by transport type:
 * <ol>
 *   <li>Modbus TCP — via TCP, UDP, TLS</li>
 *   <li>Modbus RTU — via TCP, UDP</li>
 *   <li>Modbus ASCII — via TCP, UDP</li>
 * </ol>
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@Testcontainers(disabledWithoutDocker = true)
public class ModbusDockerIT {

    private String host;
    private int tcpPort;
    private int rtuTcpPort;
    private int asciiTcpPort;
    private int tlsPort;
    private int udpPort;
    private int udpRtuPort;
    private int udpAsciiPort;

    /** Pymodbus server ports inside the container — one per driver/transport combination. */
    private static final int PYMODBUS_TCP_PORT = 5020;
    private static final int PYMODBUS_RTU_TCP_PORT = 5021;
    private static final int PYMODBUS_ASCII_TCP_PORT = 5022;
    private static final int PYMODBUS_TLS_PORT = 5023;
    private static final int PYMODBUS_UDP_PORT = 5025;
    private static final int PYMODBUS_UDP_RTU_PORT = 5026;
    private static final int PYMODBUS_UDP_ASCII_PORT = 5027;

    @SuppressWarnings("resource")
    private final GenericContainer<?> modbusSimulator = new GenericContainer<>(
        new ImageFromDockerfile()
            .withFileFromClasspath("Dockerfile", "modbus/Dockerfile")
            .withFileFromClasspath("modbus_server.py", "modbus/modbus_server.py"))
        .withExposedPorts(PYMODBUS_TCP_PORT, PYMODBUS_RTU_TCP_PORT, PYMODBUS_ASCII_TCP_PORT, PYMODBUS_TLS_PORT)
        .withCreateContainerCmdModifier(cmd -> {
            // Add UDP port exposure — Testcontainers' withExposedPorts only handles TCP.
            var udpExposed = com.github.dockerjava.api.model.ExposedPort.udp(PYMODBUS_UDP_PORT);
            var udpRtuExposed = com.github.dockerjava.api.model.ExposedPort.udp(PYMODBUS_UDP_RTU_PORT);
            var udpAsciiExposed = com.github.dockerjava.api.model.ExposedPort.udp(PYMODBUS_UDP_ASCII_PORT);

            // Add to exposed ports list
            var existing = cmd.getExposedPorts();
            var all = new com.github.dockerjava.api.model.ExposedPort[existing.length + 3];
            System.arraycopy(existing, 0, all, 0, existing.length);
            all[existing.length] = udpExposed;
            all[existing.length + 1] = udpRtuExposed;
            all[existing.length + 2] = udpAsciiExposed;
            cmd.withExposedPorts(all);

            // Add dynamic port bindings for all UDP ports
            var ports = cmd.getHostConfig().getPortBindings();
            if (ports == null) {
                ports = new com.github.dockerjava.api.model.Ports();
            }
            ports.bind(udpExposed, com.github.dockerjava.api.model.Ports.Binding.empty());
            ports.bind(udpRtuExposed, com.github.dockerjava.api.model.Ports.Binding.empty());
            ports.bind(udpAsciiExposed, com.github.dockerjava.api.model.Ports.Binding.empty());
            cmd.getHostConfig().withPortBindings(ports);
        })
        .waitingFor(Wait.forLogMessage(".*Modbus server started.*", 1)
            .withStartupTimeout(Duration.ofSeconds(120)));

    @RegisterExtension
    ContainerArtifactsExtension artifacts = new ContainerArtifactsExtension(modbusSimulator);

    @BeforeAll
    void startContainer() {
        modbusSimulator.start();

        host = modbusSimulator.getHost();
        tcpPort = modbusSimulator.getMappedPort(PYMODBUS_TCP_PORT);
        rtuTcpPort = modbusSimulator.getMappedPort(PYMODBUS_RTU_TCP_PORT);
        asciiTcpPort = modbusSimulator.getMappedPort(PYMODBUS_ASCII_TCP_PORT);
        tlsPort = modbusSimulator.getMappedPort(PYMODBUS_TLS_PORT);

        // Retrieve the dynamically mapped UDP ports from Docker's port bindings
        var bindings = modbusSimulator.getContainerInfo().getNetworkSettings().getPorts().getBindings();
        udpPort = Integer.parseInt(bindings.get(
            com.github.dockerjava.api.model.ExposedPort.udp(PYMODBUS_UDP_PORT))[0].getHostPortSpec());
        udpRtuPort = Integer.parseInt(bindings.get(
            com.github.dockerjava.api.model.ExposedPort.udp(PYMODBUS_UDP_RTU_PORT))[0].getHostPortSpec());
        udpAsciiPort = Integer.parseInt(bindings.get(
            com.github.dockerjava.api.model.ExposedPort.udp(PYMODBUS_UDP_ASCII_PORT))[0].getHostPortSpec());

        System.out.println("Pymodbus servers started:"
            + "\n  TCP:       " + host + ":" + tcpPort
            + "\n  TCP (RTU): " + host + ":" + rtuTcpPort
            + "\n  TCP (ASCII): " + host + ":" + asciiTcpPort
            + "\n  TLS:       " + host + ":" + tlsPort
            + "\n  UDP:       " + host + ":" + udpPort
            + "\n  UDP (RTU): " + host + ":" + udpRtuPort
            + "\n  UDP (ASCII): " + host + ":" + udpAsciiPort);
    }

    @AfterAll
    void stopContainer() {
        if (modbusSimulator != null) {
            modbusSimulator.stop();
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus TCP via TCP — Holding Register Reads (all data types)
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("TCP/TCP: Connection")
    void testTcpTcpConnection() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            assertTrue(conn.isConnected(), "Connection should be established");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read BOOL (TRUE)")
    void testTcpTcpReadBool() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:1:BOOL")
                .build().execute().get();
            assertTrue(resp.getBoolean("value"), "Expected BOOL value to be TRUE");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read BYTE (42)")
    void testTcpTcpReadByte() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:2:BYTE")
                .build().execute().get();
            assertEquals(42, (int) resp.getShort("value"), "Expected BYTE value to be 42");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read WORD (42424)")
    void testTcpTcpReadWord() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:3:WORD")
                .build().execute().get();
            assertEquals(42424, resp.getInteger("value"), "Expected WORD value to be 42424");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read DWORD (4242442424)")
    void testTcpTcpReadDWord() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:4:DWORD")
                .build().execute().get();
            assertEquals(4242442424L, resp.getLong("value"), "Expected DWORD value to be 4242442424");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read LWORD (4242442424242424242)")
    void testTcpTcpReadLWord() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:6:LWORD")
                .build().execute().get();
            assertEquals(4242442424242424242L, resp.getBigInteger("value").longValue());
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read SINT (-42)")
    void testTcpTcpReadSInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:10:SINT")
                .build().execute().get();
            assertEquals(-42, (int) resp.getByte("value"), "Expected SINT value to be -42");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read USINT (42)")
    void testTcpTcpReadUSInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:11:USINT")
                .build().execute().get();
            assertEquals(42, (int) resp.getShort("value"), "Expected USINT value to be 42");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read INT (-2424)")
    void testTcpTcpReadInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:12:INT")
                .build().execute().get();
            assertEquals(-2424, (int) resp.getShort("value"), "Expected INT value to be -2424");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read UINT (42424)")
    void testTcpTcpReadUInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:13:UINT")
                .build().execute().get();
            assertEquals(42424, resp.getInteger("value"), "Expected UINT value to be 42424");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read DINT (-242442424)")
    void testTcpTcpReadDInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:14:DINT")
                .build().execute().get();
            assertEquals(-242442424, resp.getInteger("value"), "Expected DINT value to be -242442424");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read UDINT (4242442424)")
    void testTcpTcpReadUDInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:16:UDINT")
                .build().execute().get();
            assertEquals(4242442424L, resp.getLong("value"), "Expected UDINT value to be 4242442424");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read LINT (-4242442424242424242)")
    void testTcpTcpReadLInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:18:LINT")
                .build().execute().get();
            assertEquals(-4242442424242424242L, resp.getLong("value"));
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read ULINT (4242442424242424242)")
    void testTcpTcpReadULInt() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:22:ULINT")
                .build().execute().get();
            assertEquals(4242442424242424242L, resp.getBigInteger("value").longValue());
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read REAL (3.141593)")
    void testTcpTcpReadReal() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:26:REAL")
                .build().execute().get();
            assertEquals(3.141593f, resp.getFloat("value"), 0.0001f);
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read LREAL (2.71828182846)")
    void testTcpTcpReadLReal() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:28:LREAL")
                .build().execute().get();
            assertEquals(2.71828182846, resp.getDouble("value"), 0.00000001);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus TCP via TCP — Multi-tag, Write/Readback, Coils, Discrete Inputs, Input Registers, Subscriptions
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("TCP/TCP: Read multiple values in single request")
    void testTcpTcpReadMultipleValues() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("byte", "holding-register:2:BYTE")
                .addTagAddress("word", "holding-register:3:WORD")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();

            assertTrue(resp.getBoolean("bool"), "Expected BOOL to be TRUE");
            assertEquals(42, (int) resp.getShort("byte"), "Expected BYTE to be 42");
            assertEquals(42424, resp.getInteger("word"), "Expected WORD to be 42424");
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("TCP/TCP: BasicPlcTest (read/write all types)")
    void testTcpTcpRunManualTest() throws Exception {
        BasicPlcTest test = new BasicPlcTest(String.format("modbus-tcp://%s:%d", host, tcpPort), new PlcNullAuthentication(), true, true, true, true, 100);
        test.addTestCase("holding-register:1:BOOL", new PlcBOOL(true));
        test.addTestCase("holding-register:2:BYTE", new PlcBYTE(42));
        test.addTestCase("holding-register:3:WORD", new PlcWORD(42424));
        test.addTestCase("holding-register:4:DWORD", new PlcDWORD(4242442424L));
        test.addTestCase("holding-register:6:LWORD", new PlcLWORD(4242442424242424242L));
        test.addTestCase("holding-register:10:SINT", new PlcSINT(-42));
        test.addTestCase("holding-register:11:USINT", new PlcUSINT(42));
        test.addTestCase("holding-register:12:INT", new PlcINT(-2424));
        test.addTestCase("holding-register:13:UINT", new PlcUINT(42424));
        test.addTestCase("holding-register:14:DINT", new PlcDINT(-242442424));
        test.addTestCase("holding-register:16:UDINT", new PlcUDINT(4242442424L));
        test.addTestCase("holding-register:18:LINT", new PlcLINT(-4242442424242424242L));
        test.addTestCase("holding-register:22:ULINT", new PlcULINT(4242442424242424242L));
        test.addTestCase("holding-register:26:REAL", new PlcREAL(3.141593F));
        test.addTestCase("holding-register:28:LREAL", new PlcLREAL(2.71828182846D));
        test.run();
    }

    @Test
    @DisplayName("TCP/TCP: Write and read-back BOOL")
    void testTcpTcpWriteReadBackBool() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder()
                .addTagAddress("value", "holding-register:100:BOOL", true)
                .build().execute().get();
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:100:BOOL")
                .build().execute().get();
            assertTrue(resp.getBoolean("value"), "Expected written BOOL TRUE to read back as TRUE");
        }
    }

    @Test
    @DisplayName("TCP/TCP: Write and read-back unsigned integers (BYTE, WORD, DWORD, LWORD)")
    void testTcpTcpWriteReadBackUnsignedIntegers() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder()
                .addTagAddress("value", "holding-register:101:BYTE", (short) 77)
                .build().execute().get();
            var resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:101:BYTE")
                .build().execute().get();
            assertEquals(77, (int) resp.getShort("value"), "Expected BYTE 77 read-back");

            conn.writeRequestBuilder()
                .addTagAddress("value", "holding-register:102:WORD", 54321)
                .build().execute().get();
            resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:102:WORD")
                .build().execute().get();
            assertEquals(54321, (int) resp.getInteger("value"), "Expected WORD 54321 read-back");

            conn.writeRequestBuilder()
                .addTagAddress("value", "holding-register:103:DWORD", 3141592653L)
                .build().execute().get();
            resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:103:DWORD")
                .build().execute().get();
            assertEquals(3141592653L, resp.getLong("value"), "Expected DWORD 3141592653 read-back");

            conn.writeRequestBuilder()
                .addTagAddress("value", "holding-register:105:LWORD", 1234567890123456789L)
                .build().execute().get();
            resp = conn.readRequestBuilder()
                .addTagAddress("value", "holding-register:105:LWORD")
                .build().execute().get();
            assertEquals(1234567890123456789L, resp.getBigInteger("value").longValue());
        }
    }

    @Test
    @DisplayName("TCP/TCP: Write and read-back signed integers")
    void testTcpTcpWriteReadBackSignedIntegers() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:111:SINT", (byte) -100).build().execute().get();
            assertEquals(-100, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:111:SINT").build().execute().get().getByte("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:112:USINT", (short) 200).build().execute().get();
            assertEquals(200, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:112:USINT").build().execute().get().getShort("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:113:INT", (short) -12345).build().execute().get();
            assertEquals(-12345, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:113:INT").build().execute().get().getShort("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:114:UINT", 54321).build().execute().get();
            assertEquals(54321, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:114:UINT").build().execute().get().getInteger("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:115:DINT", -1234567890).build().execute().get();
            assertEquals(-1234567890, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:115:DINT").build().execute().get().getInteger("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:117:UDINT", 3141592653L).build().execute().get();
            assertEquals(3141592653L, conn.readRequestBuilder().addTagAddress("value", "holding-register:117:UDINT").build().execute().get().getLong("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:119:LINT", -1234567890123456789L).build().execute().get();
            assertEquals(-1234567890123456789L, conn.readRequestBuilder().addTagAddress("value", "holding-register:119:LINT").build().execute().get().getLong("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:123:ULINT", 1234567890123456789L).build().execute().get();
            assertEquals(1234567890123456789L, conn.readRequestBuilder().addTagAddress("value", "holding-register:123:ULINT").build().execute().get().getBigInteger("value").longValue());
        }
    }

    @Test
    @DisplayName("TCP/TCP: Write and read-back floats (REAL, LREAL)")
    void testTcpTcpWriteReadBackFloats() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:127:REAL", 2.71828f).build().execute().get();
            assertEquals(2.71828f, conn.readRequestBuilder().addTagAddress("value", "holding-register:127:REAL").build().execute().get().getFloat("value"), 0.0001f);

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:129:LREAL", 1.41421356237).build().execute().get();
            assertEquals(1.41421356237, conn.readRequestBuilder().addTagAddress("value", "holding-register:129:LREAL").build().execute().get().getDouble("value"), 0.00000001);
        }
    }

    @Test
    @DisplayName("TCP/TCP: Write and read-back multiple tags in single request")
    void testTcpTcpWriteReadBackMultipleTags() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder()
                .addTagAddress("bool", "holding-register:133:BOOL", true)
                .addTagAddress("int", "holding-register:134:INT", (short) -9876)
                .addTagAddress("real", "holding-register:135:REAL", 1.618f)
                .build().execute().get();

            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:133:BOOL")
                .addTagAddress("int", "holding-register:134:INT")
                .addTagAddress("real", "holding-register:135:REAL")
                .build().execute().get();

            assertTrue(resp.getBoolean("bool"));
            assertEquals(-9876, (int) resp.getShort("int"));
            assertEquals(1.618f, resp.getFloat("real"), 0.001f);
        }
    }

    @Test
    @DisplayName("TCP/TCP: Overwrite existing value")
    void testTcpTcpOverwriteValue() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:140:INT", (short) 1000).build().execute().get();
            assertEquals(1000, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:140:INT").build().execute().get().getShort("value"));

            conn.writeRequestBuilder().addTagAddress("value", "holding-register:140:INT", (short) 2000).build().execute().get();
            assertEquals(2000, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:140:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read pre-populated coils")
    void testTcpTcpReadCoils() throws Exception {
        boolean[] expected = {true, false, true, true, false, true, false, false, true, false};
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            for (int i = 0; i < expected.length; i++) {
                int address = i + 1;
                var resp = conn.readRequestBuilder()
                    .addTagAddress("value", "coil:" + address)
                    .build().execute().get();
                assertEquals(expected[i], resp.getBoolean("value"),
                    "Expected coil:" + address + " to be " + expected[i]);
            }
        }
    }

    @Test
    @DisplayName("TCP/TCP: Write and read-back coils (state transitions)")
    void testTcpTcpWriteReadBackCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "coil:100", true).build().execute().get();
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:100").build().execute().get().getBoolean("value"));

            conn.writeRequestBuilder().addTagAddress("value", "coil:101", false).build().execute().get();
            assertFalse(conn.readRequestBuilder().addTagAddress("value", "coil:101").build().execute().get().getBoolean("value"));

            conn.writeRequestBuilder().addTagAddress("value", "coil:101", true).build().execute().get();
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:101").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read pre-populated discrete inputs")
    void testTcpTcpReadDiscreteInputs() throws Exception {
        boolean[] expected = {false, true, false, true, true, false, true, true, false, true};
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            for (int i = 0; i < expected.length; i++) {
                int address = i + 1;
                var resp = conn.readRequestBuilder()
                    .addTagAddress("value", "discrete-input:" + address)
                    .build().execute().get();
                assertEquals(expected[i], resp.getBoolean("value"),
                    "Expected discrete-input:" + address + " to be " + expected[i]);
            }
        }
    }

    @Test
    @DisplayName("TCP/TCP: Read pre-populated input registers")
    void testTcpTcpReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager()
            .getConnection(String.format("modbus-tcp://%s:%d", host, tcpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "input-register:1:BOOL").build().execute().get().getBoolean("value"));
            assertEquals(-2424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:2:INT").build().execute().get().getShort("value"));
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
            assertEquals(-242442424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:4:DINT").build().execute().get().getInteger("value"));
            assertEquals(3.141593f, conn.readRequestBuilder().addTagAddress("value", "input-register:6:REAL").build().execute().get().getFloat("value"), 0.0001f);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus TCP via UDP
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("TCP/UDP: Read holding registers")
    void testTcpUdpReadHoldingRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:udp://%s:%d", host, udpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("int", "holding-register:12:INT")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();
            assertTrue(resp.getBoolean("bool"));
            assertEquals(-2424, (int) resp.getShort("int"));
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("TCP/UDP: Write and read-back")
    void testTcpUdpWriteReadBack() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:udp://%s:%d", host, udpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:180:INT", (short) 4444).build().execute().get();
            assertEquals(4444, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:180:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("TCP/UDP: Read coils")
    void testTcpUdpReadCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:udp://%s:%d", host, udpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:1").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("TCP/UDP: Read discrete inputs")
    void testTcpUdpReadDiscreteInputs() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:udp://%s:%d", host, udpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "discrete-input:2").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("TCP/UDP: Read input registers")
    void testTcpUdpReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:udp://%s:%d", host, udpPort))) {
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus TCP via TLS
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("TCP/TLS: Read holding registers")
    void testTcpTlsReadHoldingRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:tls://%s:%d?tls.verify=false", host, tlsPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("int", "holding-register:12:INT")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();
            assertTrue(resp.getBoolean("bool"));
            assertEquals(-2424, (int) resp.getShort("int"));
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("TCP/TLS: Write and read-back")
    void testTcpTlsWriteReadBack() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:tls://%s:%d?tls.verify=false", host, tlsPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:170:INT", (short) -5555).build().execute().get();
            assertEquals(-5555, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:170:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("TCP/TLS: Read coils")
    void testTcpTlsReadCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:tls://%s:%d?tls.verify=false", host, tlsPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:1").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("TCP/TLS: Read discrete inputs")
    void testTcpTlsReadDiscreteInputs() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:tls://%s:%d?tls.verify=false", host, tlsPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "discrete-input:2").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("TCP/TLS: Read input registers")
    void testTcpTlsReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-tcp:tls://%s:%d?tls.verify=false", host, tlsPort))) {
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus RTU via TCP
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("RTU/TCP: Read holding registers")
    void testRtuTcpReadHoldingRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:tcp://%s:%d", host, rtuTcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("int", "holding-register:12:INT")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();
            assertTrue(resp.getBoolean("bool"));
            assertEquals(-2424, (int) resp.getShort("int"));
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("RTU/TCP: Write and read-back")
    void testRtuTcpWriteReadBack() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:tcp://%s:%d", host, rtuTcpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:150:INT", (short) -7777).build().execute().get();
            assertEquals(-7777, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:150:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("RTU/TCP: Read coils")
    void testRtuTcpReadCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:tcp://%s:%d", host, rtuTcpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:1").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("RTU/TCP: Read discrete inputs")
    void testRtuTcpReadDiscreteInputs() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:tcp://%s:%d", host, rtuTcpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "discrete-input:2").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("RTU/TCP: Read input registers")
    void testRtuTcpReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:tcp://%s:%d", host, rtuTcpPort))) {
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus RTU via UDP
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("RTU/UDP: Read holding registers")
    void testRtuUdpReadHoldingRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:udp://%s:%d", host, udpRtuPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("int", "holding-register:12:INT")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();
            assertTrue(resp.getBoolean("bool"));
            assertEquals(-2424, (int) resp.getShort("int"));
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("RTU/UDP: Write and read-back")
    void testRtuUdpWriteReadBack() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:udp://%s:%d", host, udpRtuPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:185:INT", (short) -8888).build().execute().get();
            assertEquals(-8888, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:185:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("RTU/UDP: Read coils")
    void testRtuUdpReadCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:udp://%s:%d", host, udpRtuPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:1").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("RTU/UDP: Read discrete inputs")
    void testRtuUdpReadDiscreteInputs() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:udp://%s:%d", host, udpRtuPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "discrete-input:2").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("RTU/UDP: Read input registers")
    void testRtuUdpReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-rtu:udp://%s:%d", host, udpRtuPort))) {
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus ASCII via TCP
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("ASCII/TCP: Read holding registers")
    void testAsciiTcpReadHoldingRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:tcp://%s:%d", host, asciiTcpPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("int", "holding-register:12:INT")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();
            assertTrue(resp.getBoolean("bool"));
            assertEquals(-2424, (int) resp.getShort("int"));
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("ASCII/TCP: Write and read-back")
    void testAsciiTcpWriteReadBack() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:tcp://%s:%d", host, asciiTcpPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:160:INT", (short) 3333).build().execute().get();
            assertEquals(3333, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:160:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("ASCII/TCP: Read coils")
    void testAsciiTcpReadCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:tcp://%s:%d", host, asciiTcpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:1").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("ASCII/TCP: Read discrete inputs")
    void testAsciiTcpReadDiscreteInputs() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:tcp://%s:%d", host, asciiTcpPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "discrete-input:2").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("ASCII/TCP: Read input registers")
    void testAsciiTcpReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:tcp://%s:%d", host, asciiTcpPort))) {
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Modbus ASCII via UDP
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    @DisplayName("ASCII/UDP: Read holding registers")
    void testAsciiUdpReadHoldingRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:udp://%s:%d", host, udpAsciiPort))) {
            var resp = conn.readRequestBuilder()
                .addTagAddress("bool", "holding-register:1:BOOL")
                .addTagAddress("int", "holding-register:12:INT")
                .addTagAddress("real", "holding-register:26:REAL")
                .build().execute().get();
            assertTrue(resp.getBoolean("bool"));
            assertEquals(-2424, (int) resp.getShort("int"));
            assertEquals(3.141593f, resp.getFloat("real"), 0.0001f);
        }
    }

    @Test
    @DisplayName("ASCII/UDP: Write and read-back")
    void testAsciiUdpWriteReadBack() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:udp://%s:%d", host, udpAsciiPort))) {
            conn.writeRequestBuilder().addTagAddress("value", "holding-register:190:INT", (short) 6666).build().execute().get();
            assertEquals(6666, (int) conn.readRequestBuilder().addTagAddress("value", "holding-register:190:INT").build().execute().get().getShort("value"));
        }
    }

    @Test
    @DisplayName("ASCII/UDP: Read coils")
    void testAsciiUdpReadCoils() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:udp://%s:%d", host, udpAsciiPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "coil:1").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("ASCII/UDP: Read discrete inputs")
    void testAsciiUdpReadDiscreteInputs() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:udp://%s:%d", host, udpAsciiPort))) {
            assertTrue(conn.readRequestBuilder().addTagAddress("value", "discrete-input:2").build().execute().get().getBoolean("value"));
        }
    }

    @Test
    @DisplayName("ASCII/UDP: Read input registers")
    void testAsciiUdpReadInputRegisters() throws Exception {
        try (var conn = new DefaultPlcDriverManager().getConnection(
            String.format("modbus-ascii:udp://%s:%d", host, udpAsciiPort))) {
            assertEquals(42424, (int) conn.readRequestBuilder().addTagAddress("value", "input-register:3:UINT").build().execute().get().getInteger("value"));
        }
    }

}
