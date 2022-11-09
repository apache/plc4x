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
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.assertj.core.api.Condition;
import org.eclipse.milo.examples.server.ExampleServer;
import org.junit.jupiter.api.*;

import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.INET_ADDRESS_PATTERN;
import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.URI_PATTERN;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.fail;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.ExecutionException;

@Disabled("Currently seems to block")
public class OpcuaPlcDriverTest {

    @BeforeAll
    static void setUp() {
        assumeTrue(() -> {
            String OS = System.getProperty("os.name").toLowerCase();
            if (OS.contains("nix")
                || OS.contains("nux")
                || OS.contains("aix")) {
                return false;
            }

            return true;
        }, "somehow opcua doesn't run properly on linux");
    }

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
    private String miloLocalAddress = "127.0.0.1:12686/milo";
    //Tcp pattern of OPC UA
    private String opcPattern = "opcua:tcp://";

    private String paramSectionDivider = "?";
    private String paramDivider = "&";

    private String tcpConnectionAddress = opcPattern + miloLocalAddress;

    private List<String> connectionStringValidSet = List.of(tcpConnectionAddress);
    private List<String> connectionStringCorruptedSet = List.of();

    private String discoveryValidParamTrue = "discovery=true";
    private String discoveryValidParamFalse = "discovery=false";
    private String discoveryCorruptedParamWrongValueNum = "discovery=1";
    private String discoveryCorruptedParamWronName = "diskovery=false";

    List<String> discoveryParamValidSet = List.of(discoveryValidParamTrue, discoveryValidParamFalse);
    List<String> discoveryParamCorruptedSet = List.of(discoveryCorruptedParamWrongValueNum, discoveryCorruptedParamWronName);

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

    @Test
    public void connectionNoParams() {
        connectionStringValidSet.forEach(connectionString -> {
            try {
                PlcConnection opcuaConnection = new PlcDriverManager().getConnection(connectionString);
                Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                assertThat(opcuaConnection).is(is_connected);
                opcuaConnection.close();
                assertThat(opcuaConnection).isNot(is_connected);
            } catch (PlcConnectionException e) {
                fail("Exception during connectionNoParams while connecting Test EXCEPTION: " + e.getMessage());
            } catch (Exception e) {
                fail("Exception during connectionNoParams while closing Test EXCEPTION: " + e.getMessage());
            }
        });
    }

    @Test
    public void connectionWithDiscoveryParam() {
        connectionStringValidSet.forEach(connectionAddress -> {
            discoveryParamValidSet.forEach(discoveryParam -> {
                String connectionString = connectionAddress + paramSectionDivider + discoveryParam;
                try {
                    PlcConnection opcuaConnection = new PlcDriverManager().getConnection(connectionString);
                    Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                    assertThat(opcuaConnection).is(is_connected);
                    opcuaConnection.close();
                    assertThat(opcuaConnection).isNot(is_connected);
                } catch (PlcConnectionException e) {
                    fail("Exception during connectionWithDiscoveryParam while connecting Test EXCEPTION: " + e.getMessage());
                } catch (Exception e) {
                    fail("Exception during connectionWithDiscoveryParam while closing Test EXCEPTION: " + e.getMessage());
                }
            });
        });
    }

    @Test
    public void readVariables() {
        try {
            PlcConnection opcuaConnection = new PlcDriverManager().getConnection(tcpConnectionAddress);
            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
            assertThat(opcuaConnection).is(is_connected);

            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder();
            builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Byte", BYTE_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Double", DOUBLE_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Float", FLOAT_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Int16", INT16_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Int32", INT32_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Int64", INT64_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("Integer", INTEGER_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("SByte", SBYTE_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("String", STRING_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("UInt16", UINT16_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("UInt32", UINT32_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("UInt64", UINT64_IDENTIFIER_READ_WRITE);
            builder.addTagAddress("UInteger", UINTEGER_IDENTIFIER_READ_WRITE);

            builder.addTagAddress("BoolArray", BOOL_ARRAY_IDENTIFIER);
            //builder.addTagAddress("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER);
            builder.addTagAddress("ByteArray", BYTE_ARRAY_IDENTIFIER);
            builder.addTagAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER);
            builder.addTagAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER);
            builder.addTagAddress("Int16Array", INT16_ARRAY_IDENTIFIER);
            builder.addTagAddress("Int32Array", INT32_ARRAY_IDENTIFIER);
            builder.addTagAddress("Int64Array", INT64_ARRAY_IDENTIFIER);
            builder.addTagAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER);
            builder.addTagAddress("StringArray", STRING_ARRAY_IDENTIFIER);
            builder.addTagAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER);
            builder.addTagAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER);
            builder.addTagAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER);

            builder.addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE);

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
            assert !opcuaConnection.isConnected();
        } catch (Exception e) {
            fail("Exception during readVariables Test EXCEPTION: " + e.getMessage());
        }
    }

    @Test
    public void writeVariables() throws Exception {
        PlcConnection opcuaConnection = new PlcDriverManager().getConnection(tcpConnectionAddress);
        Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
        assertThat(opcuaConnection).is(is_connected);

        PlcWriteRequest.Builder builder = opcuaConnection.writeRequestBuilder();
        builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE, true);
        builder.addTagAddress("Byte", BYTE_IDENTIFIER_READ_WRITE + ";BYTE", (short) 3);
        builder.addTagAddress("Double", DOUBLE_IDENTIFIER_READ_WRITE, 0.5d);
        builder.addTagAddress("Float", FLOAT_IDENTIFIER_READ_WRITE, 0.5f);
        //builder.addTagAddress("Int16", INT16_IDENTIFIER_READ_WRITE + "", (short) 1);
        builder.addTagAddress("Int32", INT32_IDENTIFIER_READ_WRITE, 42);
        builder.addTagAddress("Int64", INT64_IDENTIFIER_READ_WRITE, 42L);
        builder.addTagAddress("Integer", INTEGER_IDENTIFIER_READ_WRITE, 42);
        builder.addTagAddress("SByte", SBYTE_IDENTIFIER_READ_WRITE + ";SINT", -127);
        builder.addTagAddress("String", STRING_IDENTIFIER_READ_WRITE, "Helllo Toddy!");
        builder.addTagAddress("UInt16", UINT16_IDENTIFIER_READ_WRITE + ";UINT", 65535);
        builder.addTagAddress("UInt32", UINT32_IDENTIFIER_READ_WRITE + ";UDINT", 101010101L);
        builder.addTagAddress("UInt64", UINT64_IDENTIFIER_READ_WRITE + ";ULINT", new BigInteger("1337"));
        builder.addTagAddress("UInteger", UINTEGER_IDENTIFIER_READ_WRITE + ";UDINT", 102020202L);


        builder.addTagAddress("BooleanArray", BOOL_ARRAY_IDENTIFIER, new Boolean[]{true, true, true, true, true});
        builder.addTagAddress("ByteArray", BYTE_ARRAY_IDENTIFIER + ";BYTE", new Short[]{1, 100, 100, 255, 123});
        builder.addTagAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER, new Double[]{1.0, 2.0, 3.0, 4.0, 5.0});
        builder.addTagAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER, new Float[]{1.0F, 2.0F, 3.0F, 4.0F, 5.0F});
        builder.addTagAddress("Int16Array", INT16_ARRAY_IDENTIFIER, new Short[]{1, 2, 3, 4, 5});
        builder.addTagAddress("Int32Array", INT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5});
        builder.addTagAddress("Int64Array", INT64_ARRAY_IDENTIFIER, new Long[]{1L, 2L, 3L, 4L, 5L});
        builder.addTagAddress("IntegerArray", INT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5});
        builder.addTagAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER, new Byte[]{1, 2, 3, 4, 5});
        builder.addTagAddress("StringArray", STRING_ARRAY_IDENTIFIER, new String[]{"1", "2", "3", "4", "5"});
        builder.addTagAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER + ";UINT", new Short[]{1, 2, 3, 4, 5});
        builder.addTagAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER + ";UDINT", new Integer[]{1, 2, 3, 4, 5});
        builder.addTagAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER + ";ULINT", new Long[]{1L, 2L, 3L, 4L, 5L});

        builder.addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE, "11");

        PlcWriteRequest request = builder.build();
        PlcWriteResponse response = request.execute().get();

        assertThat(response.getResponseCode("Bool")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Byte")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Double")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Float")).isEqualTo(PlcResponseCode.OK);
        //assertThat(response.getResponseCode("Int16")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Int32")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Int64")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Integer")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("SByte")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInt16")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInt32")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInt64")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInteger")).isEqualTo(PlcResponseCode.OK);

        assertThat(response.getResponseCode("BooleanArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("ByteArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("DoubleArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("FloatArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Int16Array")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Int32Array")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("Int64Array")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("IntegerArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("SByteArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("StringArray")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInt16Array")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInt32Array")).isEqualTo(PlcResponseCode.OK);
        assertThat(response.getResponseCode("UInt64Array")).isEqualTo(PlcResponseCode.OK);

        assertThat(response.getResponseCode("DoesNotExists")).isEqualTo(PlcResponseCode.NOT_FOUND);

        opcuaConnection.close();
        assert !opcuaConnection.isConnected();
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
            private PlcConnection connection;

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

                } catch (ExecutionException executionException) {
                    executionException.printStackTrace();
                } catch (InterruptedException interruptedException) {
                    interruptedException.printStackTrace();
                }

            }
        }

        class WriteWorker extends Thread {
            private PlcConnection connection;

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
                } catch (ExecutionException executionException) {
                    executionException.printStackTrace();
                } catch (InterruptedException interruptedException) {
                    interruptedException.printStackTrace();
                }
            }
        }


        try {
            PlcConnection opcuaConnection = new PlcDriverManager().getConnection(tcpConnectionAddress);
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

}
