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

import java.io.File;
import java.io.FileOutputStream;
import java.lang.reflect.Array;
import java.net.URLEncoder;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.security.MessageSecurity;
import org.apache.plc4x.java.opcua.security.SecurityPolicy;
import org.apache.plc4x.java.opcua.tag.OpcuaTag;
import org.assertj.core.api.Condition;
import org.assertj.core.api.SoftAssertions;
import org.junit.jupiter.api.*;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.ExecutionException;
import java.util.stream.Stream;
import org.testcontainers.containers.output.Slf4jLogConsumer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static java.util.Map.entry;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.fail;

@Testcontainers(disabledWithoutDocker = true)
public class OpcuaPlcDriverTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaPlcDriverTest.class);

    private static final String APPLICATION_URI = "urn:org.apache:plc4x";
    private static final KeystoreGenerator SERVER_KEY_STORE_GENERATOR = new KeystoreGenerator("password", 2048, APPLICATION_URI);
    private static final KeystoreGenerator CLIENT_KEY_STORE_GENERATOR = new KeystoreGenerator("changeit", 2048, APPLICATION_URI, "plc4x_plus_milo", "client");

    @Container
    public final MiloTestContainer milo = new MiloTestContainer()
        //.withCreateContainerCmdModifier(cmd -> cmd.withHostName("test-opcua-server"))
        .withLogConsumer(new Slf4jLogConsumer(LOGGER))
        .withFileSystemBind(SECURITY_DIR.getAbsolutePath(), "/tmp/server/security")
        //.withEnv("JAVA_TOOL_OPTIONS", "-agentlib:jdwp=transport=dt_socket,address=*:8000,server=y,suspend=y")
//        .withCommand("java -cp '/opt/milo/*:/opt/milo/' org.eclipse.milo.examples.server.TestMiloServer")
        ;

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
    private static final String DOES_NOT_EXISTS_TAG_NAME = "DoesNotExists"; // tag name

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

    // Address of local milo server, since it comes from test container its hostname and port is not static
    private final String miloLocalAddress = "%s:%d/milo";
    //Tcp pattern of OPC UA
    private final String opcPattern = "opcua:tcp://";

    private static final String PARAM_SECTION_DIVIDER = "?";
    private static final String PARAM_DIVIDER = "&";

    private final String discoveryValidParamTrue = "discovery=true";
    private final String discoveryValidParamFalse = "discovery=false";
    private final String discoveryCorruptedParamWrongValueNum = "discovery=1";
    private final String discoveryCorruptedParamWrongName = "diskovery=false";

    private String tcpConnectionAddress;
    private List<String> connectionStringValidSet;

    final List<String> discoveryParamValidSet = List.of(discoveryValidParamTrue, discoveryValidParamFalse);
    List<String> discoveryParamCorruptedSet = List.of(discoveryCorruptedParamWrongValueNum, discoveryCorruptedParamWrongName);
    private static File SECURITY_DIR;
    private static File CLIENT_KEY_STORE;

    @BeforeAll
    public static void prepare() throws Exception {
        Path tempDirectory = Files.createTempDirectory("plc4x_opcua_client");

        SECURITY_DIR = new File(tempDirectory.toFile().getAbsolutePath(), "server/security");
        File pkiDir = new File(SECURITY_DIR, "pki");
        File trustedCerts = new File(pkiDir, "trusted/certs");
        if (!pkiDir.mkdirs() || !trustedCerts.mkdirs()) {
            throw new IllegalStateException("Could not start test - missing permissions to create temporary files");
        }

        // pre-provisioned server certificate
        try (FileOutputStream bos = new FileOutputStream(new File(pkiDir.getParentFile(), "example-server.pfx"))) {
            SERVER_KEY_STORE_GENERATOR.writeKeyStoreTo(bos);
        }

        // pre-provisioned client certificate, doing it here because server might be slow in picking up them, and we don't want to wait with our tests
        CLIENT_KEY_STORE = Files.createTempFile("plc4x_opcua_client_", ".p12").toAbsolutePath().toFile();
        try (FileOutputStream outputStream = new FileOutputStream(CLIENT_KEY_STORE)) {
            CLIENT_KEY_STORE_GENERATOR.writeKeyStoreTo(outputStream);
        }
        try (FileOutputStream outputStream = new FileOutputStream(new File(trustedCerts, "plc4x.crt"))) {
            CLIENT_KEY_STORE_GENERATOR.writeCertificateTo(outputStream);
        }
    }

    @BeforeEach
    public void startUp() throws Exception {
        tcpConnectionAddress = String.format(opcPattern + miloLocalAddress, milo.getHost(), milo.getMappedPort(12686)) + "?endpoint-port=12686";
        connectionStringValidSet = List.of(tcpConnectionAddress);
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
            return connectionStringValidSet.stream()
                .map(connectionString -> DynamicTest.dynamicTest(connectionString, () -> {
                    PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
                    Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                    assertThat(opcuaConnection).is(is_connected);
                    opcuaConnection.close();
                    assertThat(opcuaConnection).isNot(is_connected);
                }))
                .map(DynamicNode.class::cast);
        }

        @TestFactory
        Stream<DynamicNode> connectionWithDiscoveryParam() throws Exception {
            return connectionStringValidSet.stream()
                .map(connectionAddress -> DynamicContainer.dynamicContainer(connectionAddress, () ->
                    discoveryParamValidSet.stream().map(discoveryParam -> DynamicTest.dynamicTest(discoveryParam, () -> {
                            String connectionString = connectionAddress + PARAM_DIVIDER + discoveryParam;
                            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
                            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                            assertThat(opcuaConnection).is(is_connected);
                            opcuaConnection.close();
                            assertThat(opcuaConnection).isNot(is_connected);
                        }))
                        .map(DynamicNode.class::cast)
                        .iterator()))
                .map(DynamicNode.class::cast);
        }

        @Test
        void connectionWithUrlAuthentication() throws Exception {
            DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
            try (PlcConnection opcuaConnection = driverManager.getConnection(tcpConnectionAddress + "&username=admin&password=password2")) {
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
            try (PlcConnection opcuaConnection = driverManager.getConnection(tcpConnectionAddress + "&username=user&password=password1",
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
        void staticConfig() throws Exception {
            DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();

            File certificateFile = Files.createTempFile("plc4x_opcua_server_", ".crt").toAbsolutePath().toFile();
            try (FileOutputStream outputStream = new FileOutputStream(certificateFile)) {
                SERVER_KEY_STORE_GENERATOR.writeCertificateTo(outputStream);
            }

            String options = params(
                entry("discovery", "false"),
                entry("server-certificate-file", certificateFile.toString().replace("\\", "/")),
                entry("key-store-file", CLIENT_KEY_STORE.toString().replace("\\", "/")), // handle windows paths
                entry("key-store-password", "changeit"),
                entry("key-store-type", "pkcs12"),
                entry("security-policy", SecurityPolicy.Basic256Sha256.name()),
                entry("message-security", MessageSecurity.SIGN.name())
            );
            try (PlcConnection opcuaConnection = driverManager.getConnection(tcpConnectionAddress + PARAM_DIVIDER
                + options)) {
                Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
                assertThat(opcuaConnection).is(is_connected);

                PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder()
                    .addTagAddress("String", STRING_IDENTIFIER_READ_WRITE);

                PlcReadRequest request = builder.build();
                PlcReadResponse response = request.execute().get();

                assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
            }
        }
    }

    @Nested
    class readWrite {
        Map<String, Entry<String, Object>> tags = Map.ofEntries(
            entry("Bool", entry(BOOL_IDENTIFIER_READ_WRITE, true)),
            entry("Byte", entry(BYTE_IDENTIFIER_READ_WRITE + ";BYTE", (short) 3)),
            entry("Double", entry(DOUBLE_IDENTIFIER_READ_WRITE, 0.5d)),
            entry("Float", entry(FLOAT_IDENTIFIER_READ_WRITE, 0.5f)),
            entry("Int16", entry(INT16_IDENTIFIER_READ_WRITE + ";INT", 1)),
            entry("Int32", entry(INT32_IDENTIFIER_READ_WRITE, 42)),
            entry("Int64", entry(INT64_IDENTIFIER_READ_WRITE, 42L)),
            entry("Integer", entry(INTEGER_IDENTIFIER_READ_WRITE, -127)),
            //entry("SByte", entry(SBYTE_IDENTIFIER_READ_WRITE, )),
            entry("String", entry(STRING_IDENTIFIER_READ_WRITE, "Hello Toddy!")),
            entry("UInt16", entry(UINT16_IDENTIFIER_READ_WRITE + ";UINT", 65535)),
            entry("UInt32", entry(UINT32_IDENTIFIER_READ_WRITE + ";UDINT", 101010101L)),
            entry("UInt64", entry(UINT64_IDENTIFIER_READ_WRITE + ";ULINT", new BigInteger("1337"))),
            entry("UInteger", entry(UINTEGER_IDENTIFIER_READ_WRITE + ";UDINT", 102020202L)),
            entry("BooleanArray", entry(BOOL_ARRAY_IDENTIFIER, new boolean[]{true, true, true, true, true})),
            // entry("ByteStringArray", entry(BYTE_STRING_ARRAY_IDENTIFIER, null)),
            entry("ByteArray", entry(BYTE_ARRAY_IDENTIFIER + ";BYTE", new Short[]{1, 100, 100, 255, 123})),
            entry("DoubleArray", entry(DOUBLE_ARRAY_IDENTIFIER, new Double[]{1.0, 2.0, 3.0, 4.0, 5.0})),
            entry("FloatArray", entry(FLOAT_ARRAY_IDENTIFIER, new Float[]{1.0F, 2.0F, 3.0F, 4.0F, 5.0F})),
            entry("Int16Array", entry(INT16_ARRAY_IDENTIFIER, new Short[]{1, 2, 3, 4, 5})),
            entry("Int32Array", entry(INT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5})),
            entry("Int64Array", entry(INT64_ARRAY_IDENTIFIER, new Long[]{1L, 2L, 3L, 4L, 5L})),
            entry("IntegerArray", entry(INT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5})),
            entry("SByteArray", entry(SBYTE_ARRAY_IDENTIFIER, new Byte[]{1, 2, 3, 4, 5})),
            entry("StringArray", entry(STRING_ARRAY_IDENTIFIER, new String[]{"1", "2", "3", "4", "5"})),
            entry("UInt16Array", entry(UINT16_ARRAY_IDENTIFIER + ";UINT", new Short[]{1, 2, 3, 4, 5})),
            entry("UInt32Array", entry(UINT32_ARRAY_IDENTIFIER + ";UDINT", new Integer[]{1, 2, 3, 4, 5})),
            entry("UInt64Array", entry(UINT64_ARRAY_IDENTIFIER + ";ULINT", new Long[]{1L, 2L, 3L, 4L, 5L})),
            entry(DOES_NOT_EXISTS_TAG_NAME, entry(DOES_NOT_EXIST_IDENTIFIER_READ_WRITE, "11"))
        );

        @ParameterizedTest
        @MethodSource("org.apache.plc4x.java.opcua.OpcuaPlcDriverTest#getConnectionSecurityPolicies")
        public void readVariables(SecurityPolicy policy, MessageSecurity messageSecurity) throws Exception {
            String connectionString = getConnectionString(policy, messageSecurity);
            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
            assertThat(opcuaConnection).is(is_connected);

            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder();
            tags.forEach((tagName, tagEntry) -> builder.addTagAddress(tagName, tagEntry.getKey()));
            PlcReadRequest request = builder.build();
            PlcReadResponse response = request.execute().get();

            SoftAssertions softly = new SoftAssertions();
            tags.keySet().forEach(tag -> {
                if (DOES_NOT_EXISTS_TAG_NAME.equals(tag)) {
                    softly.assertThat(response.getResponseCode(tag))
                        .describedAs("Tag %s should not exist and return NOT_FOUND status", tag)
                        .isEqualTo(PlcResponseCode.NOT_FOUND);
                } else {
                    softly.assertThat(response.getResponseCode(tag))
                        .describedAs("Tag %s should exist and return OK status", tag)
                        .isEqualTo(PlcResponseCode.OK);
                }
            });
            softly.assertAll();

            opcuaConnection.close();
            assertThat(opcuaConnection.isConnected()).isFalse();
        }

        @ParameterizedTest
        @MethodSource("org.apache.plc4x.java.opcua.OpcuaPlcDriverTest#getConnectionSecurityPolicies")
        public void writeVariables(SecurityPolicy policy, MessageSecurity messageSecurity) throws Exception {
            String connectionString = getConnectionString(policy, messageSecurity);
            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
            Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
            assertThat(opcuaConnection).is(is_connected);

            PlcWriteRequest.Builder builder = opcuaConnection.writeRequestBuilder();
            tags.forEach((tagName, tagEntry) -> {
                System.out.println("Write tag " + tagName + " " + tagEntry);
                try {
                    Object value = tagEntry.getValue();
                    if (value.getClass().isArray()) {
                        Object[] values = new Object[Array.getLength(value)];
                        for (int index = 0; index < Array.getLength(value); index++) {
                            values[index] = Array.get(value, index);
                        }
                        builder.addTagAddress(tagName, tagEntry.getKey(), values);
                    } else {
                        builder.addTagAddress(tagName, tagEntry.getKey(), value);
                    }
                } catch (PlcUnsupportedDataTypeException e) {
                    fail(e.toString());
                }
            });
            PlcWriteRequest request = builder.build();
            PlcWriteResponse response = request.execute().get();

            SoftAssertions softly = new SoftAssertions();
            tags.keySet().forEach(tag -> {
                if (DOES_NOT_EXISTS_TAG_NAME.equals(tag)) {
                    softly.assertThat(response.getResponseCode(DOES_NOT_EXISTS_TAG_NAME))
                        .describedAs("Tag %s should not exist and return NOT_FOUND status", tag)
                        .isEqualTo(PlcResponseCode.NOT_FOUND);
                } else {
                    softly.assertThat(response.getResponseCode(tag))
                        .describedAs("Tag %s should exist and return OK status", tag)
                        .isEqualTo(PlcResponseCode.OK);
                }
            });
            softly.assertAll();

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

    private String getConnectionString(SecurityPolicy policy, MessageSecurity messageSecurity) throws Exception {
        switch (policy) {
            case NONE:
                return tcpConnectionAddress;

            case Basic256:
            case Basic128Rsa15:
            case Basic256Sha256:
            case Aes128_Sha256_RsaOaep:
            case Aes256_Sha256_RsaPss:
                String connectionParams = params(
                    entry("key-store-file", CLIENT_KEY_STORE.getAbsoluteFile().toString().replace("\\", "/")), // handle windows paths
                    entry("key-store-password", "changeit"),
                    entry("key-store-type", "pkcs12"),
                    entry("security-policy", policy.name()),
                    entry("message-security", messageSecurity.name())
                );

                return tcpConnectionAddress + PARAM_DIVIDER + connectionParams;
            default:
                throw new IllegalStateException();
        }
    }

    private static Stream<Arguments> getConnectionSecurityPolicies() {
        return Stream.of(
            Arguments.of(SecurityPolicy.NONE, MessageSecurity.NONE),
            Arguments.of(SecurityPolicy.NONE, MessageSecurity.SIGN),
            Arguments.of(SecurityPolicy.NONE, MessageSecurity.SIGN_ENCRYPT),
            //Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.NONE),
            Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.SIGN),
            Arguments.of(SecurityPolicy.Basic256Sha256, MessageSecurity.SIGN_ENCRYPT),
            //Arguments.of(SecurityPolicy.Basic256, MessageSecurity.NONE),
            Arguments.of(SecurityPolicy.Basic256, MessageSecurity.SIGN),
            Arguments.of(SecurityPolicy.Basic256, MessageSecurity.SIGN_ENCRYPT),
            //Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.NONE),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN),
            Arguments.of(SecurityPolicy.Basic128Rsa15, MessageSecurity.SIGN_ENCRYPT),
            //Arguments.of(SecurityPolicy.Aes128_Sha256_RsaOaep, MessageSecurity.NONE),
            Arguments.of(SecurityPolicy.Aes128_Sha256_RsaOaep, MessageSecurity.SIGN),
            Arguments.of(SecurityPolicy.Aes128_Sha256_RsaOaep, MessageSecurity.SIGN_ENCRYPT),
            //Arguments.of(SecurityPolicy.Aes256_Sha256_RsaPss, MessageSecurity.NONE),
            Arguments.of(SecurityPolicy.Aes256_Sha256_RsaPss, MessageSecurity.SIGN),
            Arguments.of(SecurityPolicy.Aes256_Sha256_RsaPss, MessageSecurity.SIGN_ENCRYPT)
        );
    }

    private static String params(Entry<String, String> ... entries) {
        return Stream.of(entries)
            .map(entry -> entry.getKey() + "=" + URLEncoder.encode(entry.getValue(), Charset.defaultCharset()))
            .collect(Collectors.joining(PARAM_DIVIDER));
    }
}
