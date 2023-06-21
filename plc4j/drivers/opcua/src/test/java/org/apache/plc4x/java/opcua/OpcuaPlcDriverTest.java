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

import io.vavr.Tuple2;
import io.vavr.collection.List;
import java.util.ArrayList;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.assertj.core.api.Condition;
import org.eclipse.milo.examples.server.ExampleServer;
import org.junit.jupiter.api.*;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.EnumSource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.ExecutionException;
import java.util.stream.Stream;

import static java.util.concurrent.Executors.newSingleThreadExecutor;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.fail;

public class OpcuaPlcDriverTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaPlcDriverTest.class);

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

    //Restricted
    public static final String STRING_IDENTIFIER_ONLY_ADMIN_READ_WRITE = "ns=2;s=HelloWorld/OnlyAdminCanRead/String";

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
    class SmokeTest {
        @Test
        public void manyReconnectionsWithSingleSubscription() throws Exception {
            PlcDriverManager driverManager = new DefaultPlcDriverManager();
            PlcConnectionManager connectionManager = driverManager.getConnectionManager();

            for (int i = 0; i < 25; i++) {
                try (PlcConnection connection = connectionManager.getConnection(tcpConnectionAddress)) {

                    PlcSubscriptionRequest request = connection.subscriptionRequestBuilder()
                            .addChangeOfStateTag("Demo", OpcuaTag.of(INTEGER_IDENTIFIER_READ_WRITE))
                            .build();

                    PlcSubscriptionResponse response = request.execute().get(60, TimeUnit.SECONDS);
                    assertThat(response.getResponseCode("Demo")).isEqualTo(PlcResponseCode.OK);

                    connection.unsubscriptionRequestBuilder()
                            .addHandles(response.getSubscriptionHandles())
                            .build()
                            .execute();
                }
            }
        }
        @Test
        public void manySubscriptionsOnSingleConnection() throws Exception {
            final int numberOfSubscriptions = 25;

            PlcDriverManager driverManager = new DefaultPlcDriverManager();
            PlcConnectionManager connectionManager = driverManager.getConnectionManager();

            ArrayList<PlcSubscriptionResponse> plcSubscriptionResponses = new ArrayList<>();
            ConcurrentLinkedDeque<PlcSubscriptionEvent> events = new ConcurrentLinkedDeque<>();

            try (PlcConnection connection = connectionManager.getConnection(tcpConnectionAddress)) {
                for (int i = 0; i < numberOfSubscriptions; i++) {
                    PlcSubscriptionRequest request = connection.subscriptionRequestBuilder()
                            .addChangeOfStateTag("Demo", OpcuaTag.of(INTEGER_IDENTIFIER_READ_WRITE))
                            .build();

                    PlcSubscriptionResponse response = request.execute().get(60, TimeUnit.SECONDS);
                    assertThat(response.getResponseCode("Demo")).isEqualTo(PlcResponseCode.OK);

                    plcSubscriptionResponses.add(response);

                    response.getSubscriptionHandles().forEach(handle -> handle.register(events::add));
                }

                for (int i = 0; i < 60; i++) {
                    if (events.size() == numberOfSubscriptions) {
                        break;
                    }
                    try {
                        Thread.sleep(1000);
                    } catch (InterruptedException e) {
                        throw new RuntimeException(e);
                    }
                }
                assertThat(events.size()).isEqualTo(numberOfSubscriptions);

                for (PlcSubscriptionResponse response : plcSubscriptionResponses) {
                    connection.unsubscriptionRequestBuilder()
                            .addHandles(response.getSubscriptionHandles())
                            .build()
                            .execute();
                }
            }
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

        @Test
        void connectionWithUrlAuthentication() throws Exception {
            DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
            try (PlcConnection opcuaConnection = driverManager.getConnection(tcpConnectionAddress + "?username=admin&password=password2")) {
                Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                assertThat(opcuaConnection).is(is_connected);

                PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder()
                        .addTagAddress("String", STRING_IDENTIFIER_ONLY_ADMIN_READ_WRITE);

                PlcReadRequest request = builder.build();
                PlcReadResponse response = request.execute().get();

                assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
            }
        }

        @Test
        void connectionWithPlcAuthentication() throws Exception {
            DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
            try (PlcConnection opcuaConnection = driverManager.getConnection(tcpConnectionAddress,
                    new PlcUsernamePasswordAuthentication("admin", "password2"))) {
                Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                assertThat(opcuaConnection).is(is_connected);

                PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder()
                        .addTagAddress("String", STRING_IDENTIFIER_ONLY_ADMIN_READ_WRITE);

                PlcReadRequest request = builder.build();
                PlcReadResponse response = request.execute().get();

                assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
            }
        }

        @Test
        void connectionWithPlcAuthenticationOverridesUrlParam() throws Exception {
            DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
            try (PlcConnection opcuaConnection = driverManager.getConnection(tcpConnectionAddress + "?username=user&password=password1",
                    new PlcUsernamePasswordAuthentication("admin", "password2"))) {
                Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                assertThat(opcuaConnection).is(is_connected);

                PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder()
                        .addTagAddress("String", STRING_IDENTIFIER_ONLY_ADMIN_READ_WRITE);

                PlcReadRequest request = builder.build();
                PlcReadResponse response = request.execute().get();

                assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
            }
        }
    }

    @Nested
    class readWrite {


        @ParameterizedTest
        @EnumSource(SecurityPolicy.class)
        public void readVariables(SecurityPolicy policy) throws Exception {
            String connectionString = getConnectionString(policy);
            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
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
            List.of(
                "Bool",
                "Byte",
                "Double",
                "Float",
                "Int16",
                "Int32",
                "Int64",
                "Integer",
                "SByte",
                "String",
                "UInt16",
                "UInt32",
                "UInt64",
                "UInteger",
                "BoolArray",
                "ByteArray",
                "DoubleArray",
                "FloatArray",
                "Int16Array",
                "Int32Array",
                "Int64Array",
                "SByteArray",
                "StringArray",
                "UInt16Array",
                "UInt32Array",
                "UInt64Array"
            ).forEach(tag -> assertThat(response.getResponseCode(tag)).isEqualTo(PlcResponseCode.OK));


            assertThat(response.getResponseCode("DoesNotExists")).isEqualTo(PlcResponseCode.NOT_FOUND);

            opcuaConnection.close();
            assertThat(opcuaConnection.isConnected()).isFalse();
        }

        @ParameterizedTest
        @EnumSource(SecurityPolicy.class)
        public void writeVariables(SecurityPolicy policy) throws Exception {

            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(getConnectionString(policy));
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
            assertThat(response.getResponseCode("DoesNotExists")).isEqualTo(PlcResponseCode.NOT_FOUND);

            opcuaConnection.close();
            assert !opcuaConnection.isConnected();
        }

    }

    /*
        Test added to test the synchronized TransactionHandler. (This was disabled before being enabled again so it might be a candidate for those tests not running properly on different platforms)
     */
    @Test
    public void multipleThreads() throws Exception {
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

                } catch (ExecutionException e) {
                    LOGGER.error("run aborted", e);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    throw new RuntimeException(e);
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
                } catch (ExecutionException e) {
                    LOGGER.error("run aborted", e);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    throw new RuntimeException(e);
                }
            }
        }


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
    }

    private String getConnectionString(SecurityPolicy policy) {
        switch (policy) {
            case NONE:
                return tcpConnectionAddress;
            case Basic128Rsa15:
                fail("Unsupported");
                return null;
            case Basic256Sha256:
                Path securityTempDir = Paths.get(System.getProperty("java.io.tmpdir"), "server");
                String keyStoreFile = securityTempDir.resolve("security").resolve("example-server.pfx").toAbsolutePath().toString();

                String certDirectory = securityTempDir.toAbsolutePath().toString();
                String connectionParams = Stream.of(
                        new Tuple2<>("keyStoreFile", keyStoreFile),
                        new Tuple2<>("certDirectory", certDirectory),
                        new Tuple2<>("keyStorePassword", "password"),
                        new Tuple2<>("securityPolicy", policy)
                    )
                    .map(tuple -> tuple._1() + "=" + tuple._2())
                    .collect(Collectors.joining(paramDivider));


                return tcpConnectionAddress + paramSectionDivider + connectionParams;
            default:
                throw new IllegalStateException();
        }
    }
}
