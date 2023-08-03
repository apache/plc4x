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
package org.apache.plc4x.java.opcua;

import io.vavr.collection.List;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.assertj.core.api.Condition;
import org.eclipse.milo.examples.server.ExampleServer;
import org.junit.jupiter.api.*;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.concurrent.ExecutionException;
import java.util.stream.Stream;

import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.INET_ADDRESS_PATTERN;
import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.URI_PATTERN;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.fail;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

public class OpcuaPlcDriverTest {

    // Read only variables of milo example server of version 3.6
    private static final String BOOL_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Boolean";
    private static final String BYTE_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Byte";
    private static final String DOUBLE_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Double";
    private static final String FLOAT_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Float";
    private static final String INT16_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Int16";
    private static final String INT32_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Int32";
    private static final String INT64_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Int64";
    private static final String INTEGER_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Integer";
    private static final String SBYTE_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/SByte";
    private static final String STRING_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/String";
    private static final String UINT16_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/UInt16";
    private static final String UINT32_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/UInt32";
    private static final String UINT64_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/UInt64";
    private static final String UINTEGER_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/UInteger";
    private static final String DOES_NOT_EXIST_IDENTIFIER_READ_WRITE = "ns=2;i=12512623";

    // At the moment not used in PLC4X or in the OPC UA driver
    private static final String BYTE_STRING_IDENTIFIER_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/ByteString";
    private static final String DATE_TIME_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/DateTime";
    private static final String DURATION_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Duration";
    private static final String GUID_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Guid";
    private static final String LOCALIZED_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/LocalizedText";
    private static final String NODE_ID_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/NodeId";
    private static final String QUALIFIED_NAM_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/QualifiedName";
    private static final String UTC_TIME_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/UtcTime";
    private static final String VARIANT_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Variant";
    private static final String XML_ELEMENT_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/XmlElement";

    //Arrays
    private static final String BOOL_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/BooleanArray";
    //private static final String BYTE_STRING_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/ByteStringArray";
    private static final String BYTE_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/ByteArray";
    private static final String DOUBLE_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/DoubleArray";
    private static final String FLOAT_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/FloatArray";
    private static final String INT16_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/Int16Array";
    private static final String INT32_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/Int32Array";
    private static final String INT64_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/Int64Array";
    private static final String INTEGER_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/IntegerArray";
    private static final String SBYTE_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/SByteArray";
    private static final String STRING_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/StringArray";
    private static final String UINT16_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/UInt16Array";
    private static final String UINT32_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/UInt32Array";
    private static final String UINT64_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/UInt64Array";
    private static final String DATE_TIME_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/DateTimeArray";

    // Address of local milo server
    private final String miloLocalAddress = "127.0.0.1:12686/milo";
    //Tcp pattern of OPC UA
    private final String opcPattern = "opcua:tcp://";

    private final String paramSectionDivider = "?";
    private final String paramDivider = "&";

    private final String tcpConnectionAddress = opcPattern + miloLocalAddress;

    private final List<String> connectionStringValidSet = List.of(tcpConnectionAddress);
    private final List<String> connectionStringCorruptedSet = List.of();

    private final String discoveryValidParamTrue = "discovery=true";
    private final String discoveryValidParamFalse = "discovery=false";
    private final String discoveryCorruptedParamWrongValueNum = "discovery=1";
    private final String discoveryCorruptedParamWrongName = "diskovery=false";

    final List<String> discoveryParamValidSet = List.of(discoveryValidParamTrue, discoveryValidParamFalse);
    List<String> discoveryParamCorruptedSet = List.of(discoveryCorruptedParamWrongValueNum, discoveryCorruptedParamWrongName);

    private static ExampleServer exampleServer;

    @BeforeAll
    public static void setup() throws Exception {
        // When switching JDK versions from a newer to an older version,
        // this can cause the server to not start correctly.
        // Deleting the directory makes sure the key-store is initialized correctly.
        Path securityBaseDir = Paths.get(System.getProperty("java.io.tmpdir"), "server", "security");
        try {
            Files.delete(securityBaseDir);
        } catch (Exception e) {
            // Ignore this ...
        }

        exampleServer = new ExampleServer();
        exampleServer.startup().get();
    }

    @AfterAll
    public static void tearDown() throws Exception {
        if (exampleServer != null) {
            exampleServer.shutdown().get();
        }
    }

    @Nested
    class ConnectionRelated {
        @TestFactory
        Stream<DynamicNode> connectionNoParams() {
            return connectionStringValidSet.toStream()
                .map(connectionString -> DynamicTest.dynamicTest(connectionString, () -> {
                    PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
                    Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                    assertThat(opcuaConnection).is(is_connected);
                    opcuaConnection.close();
                    assertThat(opcuaConnection).isNot(is_connected);
                }))
                .map(DynamicNode.class::cast)
                .toJavaStream();
        }

        @TestFactory
        Stream<DynamicNode> connectionWithDiscoveryParam() throws Exception {
            return connectionStringValidSet.toStream()
                .map(connectionAddress -> DynamicContainer.dynamicContainer(connectionAddress, () ->
                    discoveryParamValidSet.toStream().map(discoveryParam -> DynamicTest.dynamicTest(discoveryParam, () -> {
                            String connectionString = connectionAddress + paramSectionDivider + discoveryParam;
                            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
                            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                            assertThat(opcuaConnection).is(is_connected);
                            opcuaConnection.close();
                            assertThat(opcuaConnection).isNot(is_connected);
                        }))
                        .map(DynamicNode.class::cast)
                        .iterator()))
                .map(DynamicNode.class::cast)
                .toJavaStream();
        }
    }

    @Nested
    class readWrite {

        @BeforeEach
        public void testForLinux() {
            checkForLinux(); // TODO: seems those do not run on linux
        }

        @Test
        public void readVariables() throws Exception {
            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
            assertThat(opcuaConnection).is(is_connected);

            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder()
                .addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE)
                .addTagAddress("Byte", BYTE_IDENTIFIER_READ_WRITE)
                .addTagAddress("Double", DOUBLE_IDENTIFIER_READ_WRITE)
                .addTagAddress("Float", FLOAT_IDENTIFIER_READ_WRITE)
                .addTagAddress("Int16", INT16_IDENTIFIER_READ_WRITE)
                .addTagAddress("Int32", INT32_IDENTIFIER_READ_WRITE)
                .addTagAddress("Int64", INT64_IDENTIFIER_READ_WRITE)
                .addTagAddress("Integer", INTEGER_IDENTIFIER_READ_WRITE)
                .addTagAddress("SByte", SBYTE_IDENTIFIER_READ_WRITE)
                .addTagAddress("String", STRING_IDENTIFIER_READ_WRITE)
                .addTagAddress("UInt16", UINT16_IDENTIFIER_READ_WRITE)
                .addTagAddress("UInt32", UINT32_IDENTIFIER_READ_WRITE)
                .addTagAddress("UInt64", UINT64_IDENTIFIER_READ_WRITE)
                .addTagAddress("UInteger", UINTEGER_IDENTIFIER_READ_WRITE)

                .addTagAddress("BoolArray", BOOL_ARRAY_IDENTIFIER)
                //.addTagAddress("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER);
                .addTagAddress("ByteArray", BYTE_ARRAY_IDENTIFIER)
                .addTagAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER)
                .addTagAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER)
                .addTagAddress("Int16Array", INT16_ARRAY_IDENTIFIER)
                .addTagAddress("Int32Array", INT32_ARRAY_IDENTIFIER)
                .addTagAddress("Int64Array", INT64_ARRAY_IDENTIFIER)
                .addTagAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER)
                .addTagAddress("StringArray", STRING_ARRAY_IDENTIFIER)
                .addTagAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER)
                .addTagAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER)
                .addTagAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER)

                .addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE);

            PlcReadRequest request = builder.build();
            PlcReadResponse response = request.execute().get();
            assertThat(response.getResponseCode("Bool")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Byte")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Double")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Float")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Int16")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Int32")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Int64")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Integer")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("SByte")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInt16")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInt32")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInt64")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInteger")).isEqualTo(PlcResponseCode.OK);

            assertThat(response.getResponseCode("BoolArray")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("ByteArray")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("DoubleArray")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("FloatArray")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Int16Array")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Int32Array")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("Int64Array")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("SByteArray")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("StringArray")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInt16Array")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInt32Array")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("UInt64Array")).isEqualTo(PlcResponseCode.OK);

            assertThat(response.getResponseCode("DoesNotExists")).isEqualTo(PlcResponseCode.NOT_FOUND);

            opcuaConnection.close();
            assertThat(opcuaConnection.isConnected()).isFalse();
        }

        @Test
        public void writeVariables() throws Exception {

            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
            assertThat(opcuaConnection).is(is_connected);

            PlcWriteRequest.Builder builder = opcuaConnection.writeRequestBuilder()
                .addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE, true)
                .addTagAddress("Byte", BYTE_IDENTIFIER_READ_WRITE + ";BYTE", (short) 3)
                .addTagAddress("Double", DOUBLE_IDENTIFIER_READ_WRITE, 0.5d)
                .addTagAddress("Float", FLOAT_IDENTIFIER_READ_WRITE, 0.5f)
                //.addTagAddress("Int16", INT16_IDENTIFIER_READ_WRITE + "", (short) 1)
                .addTagAddress("Int32", INT32_IDENTIFIER_READ_WRITE, 42)
                .addTagAddress("Int64", INT64_IDENTIFIER_READ_WRITE, 42L)
                .addTagAddress("Integer", INTEGER_IDENTIFIER_READ_WRITE, 42)
                .addTagAddress("SByte", SBYTE_IDENTIFIER_READ_WRITE + ";SINT", -127)
                .addTagAddress("String", STRING_IDENTIFIER_READ_WRITE, "Helllo Toddy!")
                .addTagAddress("UInt16", UINT16_IDENTIFIER_READ_WRITE + ";UINT", 65535)
                .addTagAddress("UInt32", UINT32_IDENTIFIER_READ_WRITE + ";UDINT", 101010101L)
                .addTagAddress("UInt64", UINT64_IDENTIFIER_READ_WRITE + ";ULINT", new BigInteger("1337"))
                .addTagAddress("UInteger", UINTEGER_IDENTIFIER_READ_WRITE + ";UDINT", 102020202L)


                .addTagAddress("BooleanArray", BOOL_ARRAY_IDENTIFIER, (Object[]) new Boolean[]{true, true, true, true, true})
                .addTagAddress("ByteArray", BYTE_ARRAY_IDENTIFIER + ";BYTE", (Object[]) new Short[]{1, 100, 100, 255, 123})
                .addTagAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER, (Object[]) new Double[]{1.0, 2.0, 3.0, 4.0, 5.0})
                .addTagAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER, (Object[]) new Float[]{1.0F, 2.0F, 3.0F, 4.0F, 5.0F})
                .addTagAddress("Int16Array", INT16_ARRAY_IDENTIFIER, (Object[]) new Short[]{1, 2, 3, 4, 5})
                .addTagAddress("Int32Array", INT32_ARRAY_IDENTIFIER, (Object[]) new Integer[]{1, 2, 3, 4, 5})
                .addTagAddress("Int64Array", INT64_ARRAY_IDENTIFIER, (Object[]) new Long[]{1L, 2L, 3L, 4L, 5L})
                .addTagAddress("IntegerArray", INT32_ARRAY_IDENTIFIER, (Object[]) new Integer[]{1, 2, 3, 4, 5})
                .addTagAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER, (Object[]) new Byte[]{1, 2, 3, 4, 5})
                .addTagAddress("StringArray", STRING_ARRAY_IDENTIFIER, (Object[]) new String[]{"1", "2", "3", "4", "5"})
                .addTagAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER + ";UINT", (Object[]) new Short[]{1, 2, 3, 4, 5})
                .addTagAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER + ";UDINT", (Object[]) new Integer[]{1, 2, 3, 4, 5})
                .addTagAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER + ";ULINT", (Object[]) new Long[]{1L, 2L, 3L, 4L, 5L})

                .addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE, "11");

            PlcWriteRequest request = builder.build();
            PlcWriteResponse response = request.execute().get();

            List.of(
                "Bool",
                "Byte",
                "Double",
                "Float",
                //"Int16", // TODO: why is htat disabled???
                "Int32",
                "Int64",
                "Integer",
                "SByte",
                "String",
                "UInt16",
                "UInt32",
                "UInt64",
                "UInteger",
                "BooleanArray",
                "ByteArray",
                "DoubleArray",
                "FloatArray",
                "Int16Array",
                "Int32Array",
                "Int64Array",
                "IntegerArray",
                "SByteArray",
                "StringArray",
                "UInt16Array",
                "UInt32Array",
                "UInt64Array"
            ).forEach(s -> {
                assertThat(response.getResponseCode(s)).withFailMessage(s + "is not ok").isEqualTo(PlcResponseCode.OK);
            });
            assertThat(response.getResponseCode("DoesNotExists")).isEqualTo(PlcResponseCode.OK);

            opcuaConnection.close();
            assert !opcuaConnection.isConnected();
        }

    }


    @Test
    public void testOpcuaAddressPattern() {
        assertThat(":tcp://localhost").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://localhost:3131").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://www.google.de").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://www.google.de:443").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://127.0.0.1").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://127.0.0.1:251").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://254.254.254.254:1337").matches(INET_ADDRESS_PATTERN);
        assertThat(":tcp://254.254.254.254").matches(INET_ADDRESS_PATTERN);


        assertThat("opcua:tcp://localhost").matches(URI_PATTERN);
        assertThat("opcua:tcp://localhost:3131").matches(URI_PATTERN);
        assertThat("opcua:tcp://www.google.de").matches(URI_PATTERN);
        assertThat("opcua:tcp://www.google.de:443").matches(URI_PATTERN);
        assertThat("opcua:tcp://127.0.0.1").matches(URI_PATTERN);
        assertThat("opcua:tcp://127.0.0.1:251").matches(URI_PATTERN);
        assertThat("opcua:tcp://254.254.254.254:1337").matches(URI_PATTERN);
        assertThat("opcua:tcp://254.254.254.254").matches(URI_PATTERN);

        assertThat("opcua:tcp://127.0.0.1?discovery=false").matches(URI_PATTERN);
        assertThat("opcua:tcp://opcua.demo-this.com:51210/UA/SampleServer?discovery=false").matches(URI_PATTERN);
    }

    /*
        Test added to test the syncronized Trnasactionhandler.
        The test originally failed one out of every 5 or so.
     */

    public void multipleThreads() {
        class ReadWorker extends Thread {
            private final PlcConnection connection;

            public ReadWorker(PlcConnection opcuaConnection) {
                this.connection = opcuaConnection;
            }

            @Override
            public void run() {
                try {
                    PlcReadRequest.Builder read_builder = connection.readRequestBuilder();
                    read_builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE);
                    PlcReadRequest read_request = read_builder.build();

                    for (int i = 0; i < 100; i++) {
                        PlcReadResponse read_response = read_request.execute().get();
                        assertThat(read_response.getResponseCode("Bool")).isEqualTo(PlcResponseCode.OK);
                    }

                } catch (ExecutionException | InterruptedException executionException) {
                    executionException.printStackTrace();
                }

            }
        }

        class WriteWorker extends Thread {
            private final PlcConnection connection;

            public WriteWorker(PlcConnection opcuaConnection) {
                this.connection = opcuaConnection;
            }

            @Override
            public void run() {
                try {
                    PlcWriteRequest.Builder write_builder = connection.writeRequestBuilder();
                    write_builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE, true);
                    PlcWriteRequest write_request = write_builder.build();

                    for (int i = 0; i < 100; i++) {
                        PlcWriteResponse write_response = write_request.execute().get();
                        assertThat(write_response.getResponseCode("Bool")).isEqualTo(PlcResponseCode.OK);
                    }
                } catch (ExecutionException | InterruptedException executionException) {
                    executionException.printStackTrace();
                }
            }
        }


        try {
            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
            assertThat(opcuaConnection).is(is_connected);

            ReadWorker read_worker = new ReadWorker(opcuaConnection);
            WriteWorker write_worker = new WriteWorker(opcuaConnection);
            read_worker.start();
            write_worker.start();

            read_worker.join();
            write_worker.join();

            opcuaConnection.close();
            assert !opcuaConnection.isConnected();
        } catch (Exception e) {
            fail("Exception during readVariables Test EXCEPTION: " + e.getMessage());
        }
    }

    private final String[] validTCPOPC = {
        "localhost",
        "127.0.0.1",
        "254.254.254.254"
    };
    private final int[] validPorts = {
        1337,
        42,
        1,
        24152
    };
    private final String[] nDiscoveryParams = {
        "discovery=false"
    };


    @TestFactory
    Stream<DynamicNode> testConnectionStringPattern() throws Exception {
        return Arrays.stream(validTCPOPC)
            .map(address -> DynamicContainer.dynamicContainer("Address: " + address, () -> Arrays.stream(validPorts)
                    .mapToObj(port -> DynamicTest.dynamicTest("Port: " + port, () -> {
                            assertThat("opcua:tcp://" + address + ":555?discovery=true").matches(URI_PATTERN);
                            assertThat("opcua:tcp://" + address + ":555?discovery=True").matches(URI_PATTERN);
                            assertThat("opcua:tcp://" + address + ":555?discovery=TRUE").matches(URI_PATTERN);
                            assertThat("opcua:tcp://" + address + ":555?Discovery=True").matches(URI_PATTERN);
                            //No Port Specified
                            assertThat("opcua:tcp://" + address + "?discovery=True").matches(URI_PATTERN);
                            //No Transport Specified
                            assertThat("opcua://" + address + ":647?discovery=True").matches(URI_PATTERN);
                            //No Params Specified
                            assertThat("opcua:tcp://" + address + ":111").matches(URI_PATTERN);
                            //No Transport and Params Specified
                            assertThat("opcua://" + address + ":754").matches(URI_PATTERN);
                        })
                    )
                    .map(DynamicNode.class::cast)
                    .iterator()
                )
            );
    }

    static void checkForLinux() {
        assumeTrue(() -> {
            String OS = System.getProperty("os.name").toLowerCase();
            return !OS.contains("nix")
                && !OS.contains("nux")
                && !OS.contains("aix");
        }, "somehow opcua doesn't run properly on linux");
    }
}
