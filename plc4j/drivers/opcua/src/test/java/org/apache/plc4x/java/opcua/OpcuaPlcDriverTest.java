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
import java.time.Duration;
import java.util.concurrent.ConcurrentLinkedDeque;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.authentication.PlcNullAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.spi.values.*;
import org.apache.plc4x.java.utils.testutils.manual.BasicPlcTest;
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
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.assertj.core.api.Assertions.fail;

@Testcontainers(disabledWithoutDocker = true)
public class OpcuaPlcDriverTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaPlcDriverTest.class);

    private static final String APPLICATION_URI = "urn:org.apache:plc4x";
    private static final KeystoreGenerator SERVER_KEY_STORE_GENERATOR = new KeystoreGenerator("password", 2048, APPLICATION_URI);
    private static final KeystoreGenerator CLIENT_KEY_STORE_GENERATOR = new KeystoreGenerator("changeit", 2048, APPLICATION_URI, "plc4x_plus_milo", "client");

    private static final File SECURITY_DIR;
    private static final File CLIENT_KEY_STORE;
    // The server's certificate, exported so the client can pin trust to it. Since the
    // driver now rejects server certificates by default (no permissive fallback), every
    // secured connection needs an explicit trust anchor.
    private static final File SERVER_CERTIFICATE;

    // Provision the server keystore / client material BEFORE the container starts. The
    // container is static (see below) so the Testcontainers extension starts it during
    // its beforeAll callback, which runs before any @BeforeAll method — hence the setup
    // lives in this static initializer rather than in @BeforeAll.
    static {
        try {
            Path tempDirectory = Files.createTempDirectory("plc4x_opcua_client");

            SECURITY_DIR = new File(tempDirectory.toFile().getAbsolutePath(), "server/security");
            File pkiDir = new File(SECURITY_DIR, "pki");
            File trustedCerts = new File(pkiDir, "trusted/certs");
            if (!pkiDir.mkdirs() || !trustedCerts.mkdirs()) {
                throw new IllegalStateException("Could not start test - missing permissions to create temporary files");
            }

            // pre-provisioned server certificate (keystore the Milo server loads on boot)
            try (FileOutputStream bos = new FileOutputStream(new File(pkiDir.getParentFile(), "example-server.pfx"))) {
                SERVER_KEY_STORE_GENERATOR.writeKeyStoreTo(bos);
            }

            // export the server certificate so secured client connections can pin trust to it
            SERVER_CERTIFICATE = Files.createTempFile("plc4x_opcua_server_", ".crt").toAbsolutePath().toFile();
            try (FileOutputStream outputStream = new FileOutputStream(SERVER_CERTIFICATE)) {
                SERVER_KEY_STORE_GENERATOR.writeCertificateTo(outputStream);
            }

            // pre-provisioned client certificate, doing it here because server might be slow in picking up them, and we don't want to wait with our tests
            CLIENT_KEY_STORE = Files.createTempFile("plc4x_opcua_client_", ".p12").toAbsolutePath().toFile();
            try (FileOutputStream outputStream = new FileOutputStream(CLIENT_KEY_STORE)) {
                CLIENT_KEY_STORE_GENERATOR.writeKeyStoreTo(outputStream);
            }
            try (FileOutputStream outputStream = new FileOutputStream(new File(trustedCerts, "plc4x.crt"))) {
                CLIENT_KEY_STORE_GENERATOR.writeCertificateTo(outputStream);
            }
        } catch (Exception e) {
            throw new ExceptionInInitializerError(e);
        }
    }

    // Static so the Milo server boots ONCE for the whole class rather than being torn
    // down and rebuilt before every test method (which is what an instance @Container does).
    @Container
    public static final MiloTestContainer milo = new MiloTestContainer()
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

    @BeforeEach
    public void startUp() throws Exception {
        tcpConnectionAddress = String.format(opcPattern + miloLocalAddress, milo.getHost(), milo.getMappedPort(12686)) + "?endpoint-port=12686";
        connectionStringValidSet = List.of(tcpConnectionAddress);
    }

    /**
     * A tag whose address the tag handler can't parse stays in the request carrying an error code
     * and a {@code null} tag. It has to come back as INVALID_ADDRESS - reading it used to throw a
     * NullPointerException while the driver built the node id (there was a TODO in the read path
     * admitting the check was missing).
     */
    @Test
    void readWithInvalidTagAddressIsReported() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcReadResponse response = connection.readRequestBuilder()
                .addTagAddress("bad", "this-is-not-a-node-id")
                .build().execute().get(30, TimeUnit.SECONDS);

            assertThat(response.getResponseCode("bad")).isEqualTo(PlcResponseCode.INVALID_ADDRESS);
        }
    }

    /**
     * The point of per-tag codes: one bad address must not cost the caller the other tags of the
     * same request, and the good values must not be shifted onto the wrong names - OPC UA maps
     * its results back to tags by position.
     */
    @Test
    void readMixesValidAndInvalidTagAddresses() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcReadResponse response = connection.readRequestBuilder()
                .addTagAddress("bad", "this-is-not-a-node-id")
                .addTagAddress("bool", BOOL_IDENTIFIER_READ_WRITE)
                .addTagAddress("int32", INT32_IDENTIFIER_READ_WRITE)
                .build().execute().get(30, TimeUnit.SECONDS);

            assertThat(response.getResponseCode("bad")).isEqualTo(PlcResponseCode.INVALID_ADDRESS);
            assertThat(response.getResponseCode("bool")).isEqualTo(PlcResponseCode.OK);
            assertThat(response.getResponseCode("int32")).isEqualTo(PlcResponseCode.OK);
            // Values must still belong to the tag that asked for them.
            assertThat(response.getPlcValue("bool").isBoolean()).isTrue();
            assertThat(response.getPlcValue("int32").isLong() || response.getPlcValue("int32").isInteger()).isTrue();
        }
    }

    @Test
    void writeWithInvalidTagAddressIsReported() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcWriteResponse response = connection.writeRequestBuilder()
                .addTagAddress("bad", "this-is-not-a-node-id", 42)
                .build().execute().get(30, TimeUnit.SECONDS);

            assertThat(response.getResponseCode("bad")).isEqualTo(PlcResponseCode.INVALID_ADDRESS);
        }
    }

    @Test
    void writeMixesValidAndInvalidTagAddresses() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcWriteResponse response = connection.writeRequestBuilder()
                .addTagAddress("bad", "this-is-not-a-node-id", 42)
                .addTagAddress("bool", BOOL_IDENTIFIER_READ_WRITE, true)
                .build().execute().get(30, TimeUnit.SECONDS);

            assertThat(response.getResponseCode("bad")).isEqualTo(PlcResponseCode.INVALID_ADDRESS);
            assertThat(response.getResponseCode("bool")).isEqualTo(PlcResponseCode.OK);
        }
    }

    @Test
    void browseWildcardDiscoversWholeAddressSpace() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            // "**" is the "browse everything" wildcard: it starts at the Objects folder and
            // recurses the whole sub-tree.
            org.apache.plc4x.java.api.messages.PlcBrowseResponse response = connection.browseRequestBuilder()
                .addQuery("all", "**")
                .build().execute().get(60, TimeUnit.SECONDS);
            assertThat(response.getResponseCode("all")).isEqualTo(PlcResponseCode.OK);

            java.util.List<org.apache.plc4x.java.api.messages.PlcBrowseItem> top = response.getValues("all");
            // The example server exposes the 'HelloWorld' folder under Objects.
            assertThat(top.stream().map(org.apache.plc4x.java.api.messages.PlcBrowseItem::getName)).contains("HelloWorld");

            java.util.List<org.apache.plc4x.java.api.messages.PlcBrowseItem> all = new ArrayList<>();
            collectBrowseItems(top, all);
            assertThat(all.size()).isGreaterThan(50);

            // Phase 2: server-side type resolution applied across the whole address space — many
            // variables carry a resolved data type, and array variables expose their dimensions.
            assertThat(all).anyMatch(i -> i.getOptions().containsKey("data-type"));
            assertThat(all).anyMatch(i -> !i.getArrayInformation().isEmpty());
        }
    }

    @Test
    void browseDiscoversAddressSpace() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            assertThat(connection.getMetadata().isBrowseSupported()).isTrue();

            org.apache.plc4x.java.api.messages.PlcBrowseResponse response = connection.browseRequestBuilder()
                .addQuery("hw", "ns=2;s=HelloWorld")
                .build().execute().get(30, TimeUnit.SECONDS);
            assertThat(response.getResponseCode("hw")).isEqualTo(PlcResponseCode.OK);

            java.util.List<org.apache.plc4x.java.api.messages.PlcBrowseItem> top = response.getValues("hw");
            assertThat(top).isNotEmpty();

            // 'ScalarTypes' is an Object (container) node: discovered, has children, not readable.
            org.apache.plc4x.java.api.messages.PlcBrowseItem scalarTypes = top.stream()
                .filter(i -> "ScalarTypes".equals(i.getName())).findFirst().orElse(null);
            assertThat(scalarTypes).as("ScalarTypes folder").isNotNull();
            assertThat(scalarTypes.isReadable()).isFalse();
            assertThat(scalarTypes.getChildren()).isNotEmpty();

            // A concrete variable underneath it is readable and writable, with a usable address.
            org.apache.plc4x.java.api.messages.PlcBrowseItem boolVar = scalarTypes.getChildren().values().stream()
                .filter(i -> "Boolean".equals(i.getName())).findFirst().orElse(null);
            assertThat(boolVar).as("ScalarTypes/Boolean variable").isNotNull();
            assertThat(boolVar.isReadable()).isTrue();
            assertThat(boolVar.isWritable()).isTrue();
            assertThat(boolVar.getTag().getAddressString()).contains("ns=2;s=HelloWorld/ScalarTypes/Boolean");

            // Phase 2: the variable's data type is resolved from the server (DataType attribute),
            // so the tag is correctly typed (address gains a ";BOOL" suffix) without a manual hint,
            // and the resolved type is surfaced as a browse option. A scalar has no array info.
            assertThat(boolVar.getTag().getAddressString()).endsWith(";BOOL");
            assertThat(((OpcuaTag) boolVar.getTag()).getDataType())
                .isEqualTo(org.apache.plc4x.java.opcua.readwrite.OpcuaDataType.BOOL);
            assertThat(boolVar.getOptions()).containsKey("data-type");
            assertThat(boolVar.getArrayInformation()).isEmpty();

            // Full-subtree recursion reached deeply-nested property nodes (AnalogValue -> EURange).
            java.util.List<org.apache.plc4x.java.api.messages.PlcBrowseItem> all = new ArrayList<>();
            collectBrowseItems(top, all);
            assertThat(all).anyMatch(i -> "EURange".equals(i.getName()));
            assertThat(all.size()).isGreaterThan(40);
        }
    }

    @Test
    void plc4xTestNamespaceExposesScalarsArraysAndMatrices() throws Exception {
        // The custom Plc4xTestNamespace (ns=3;s=Test/...) mirrors the ADS manual-test structure so
        // the driver can be exercised against the same matrix of addressing variants. This verifies
        // the namespace is present and that scalars, whole arrays and multi-dimensional matrices
        // read back correctly.
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcReadResponse response = connection.readRequestBuilder()
                .addTagAddress("bool", "ns=3;s=Test/Scalar/Bool")
                .addTagAddress("dint", "ns=3;s=Test/Scalar/DInt")
                .addTagAddress("string", "ns=3;s=Test/Scalar/String")
                .addTagAddress("intArray", "ns=3;s=Test/Array/Int")
                .addTagAddress("intMatrix", "ns=3;s=Test/Matrix/Int_2x3")
                .build().execute().get(30, TimeUnit.SECONDS);

            for (String tag : new String[]{"bool", "dint", "string", "intArray", "intMatrix"}) {
                assertThat(response.getResponseCode(tag)).describedAs("read of '%s'", tag)
                    .isEqualTo(PlcResponseCode.OK);
            }

            // Scalars
            assertThat(response.getBoolean("bool")).isTrue();
            assertThat(response.getInteger("dint")).isEqualTo(-12345678);
            assertThat(response.getString("string")).isEqualTo("Hello PLC4X");

            // Whole 1-D array -> PlcList of 5 Int16 values.
            org.apache.plc4x.java.api.value.PlcValue intArray = response.getPlcValue("intArray");
            assertThat(intArray.isList()).isTrue();
            assertThat(intArray.getList()).hasSize(5);
            assertThat(intArray.getList().get(0).getShort()).isEqualTo((short) -3);
            assertThat(intArray.getList().get(4).getShort()).isEqualTo((short) 3);

            // Whole 2-D matrix Int16[2][3] -> nested PlcList (2 rows x 3 cols).
            org.apache.plc4x.java.api.value.PlcValue matrix = response.getPlcValue("intMatrix");
            assertThat(matrix.isList()).isTrue();
            assertThat(matrix.getList()).hasSize(2);
            assertThat(matrix.getList().get(0).getList().get(0).getShort()).isEqualTo((short) 10);
            assertThat(matrix.getList().get(1).getList().get(2).getShort()).isEqualTo((short) -12);
        }
    }

    @Test
    void readsAndWritesArrayElementsViaIndexRange() throws Exception {
        // Array/Int = {-3, -1, 0, 1, 3}; Matrix/Int_2x3 = {{10,11,12},{-10,-11,-12}}.
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcReadResponse read = connection.readRequestBuilder()
                .addTagAddress("elem", "ns=3;s=Test/Array/Int[3]")        // single element -> scalar
                .addTagAddress("elemBase1", "ns=3;s=Test/Array/Int[4;1]") // same element, 1-based
                .addTagAddress("slice", "ns=3;s=Test/Array/Int[1..3]")    // inclusive range -> list
                .addTagAddress("cell", "ns=3;s=Test/Matrix/Int_2x3[1][2]") // matrix element -> scalar
                .build().execute().get(30, TimeUnit.SECONDS);

            for (String tag : new String[]{"elem", "elemBase1", "slice", "cell"}) {
                assertThat(read.getResponseCode(tag)).describedAs("read of '%s'", tag)
                    .isEqualTo(PlcResponseCode.OK);
            }
            // [3] -> index 3 (0-based) -> value 1.
            assertThat(read.getShort("elem")).isEqualTo((short) 1);
            // [4;1] -> base 1 -> index 3 -> value 1 (same element).
            assertThat(read.getShort("elemBase1")).isEqualTo((short) 1);
            // [1..3] -> indices 1,2,3 -> {-1, 0, 1}.
            assertThat(read.getPlcValue("slice").getList()).hasSize(3);
            assertThat(read.getPlcValue("slice").getList().get(0).getShort()).isEqualTo((short) -1);
            assertThat(read.getPlcValue("slice").getList().get(2).getShort()).isEqualTo((short) 1);
            // Matrix [1][2] -> row 1, col 2 -> -12.
            assertThat(read.getShort("cell")).isEqualTo((short) -12);

            // Indexed write: overwrite a single element of Array/DInt and read it back via IndexRange.
            connection.writeRequestBuilder()
                .addTagAddress("w", "ns=3;s=Test/Array/DInt[2]", 424242)
                .build().execute().get(30, TimeUnit.SECONDS);
            PlcReadResponse back = connection.readRequestBuilder()
                .addTagAddress("w", "ns=3;s=Test/Array/DInt[2]")
                .build().execute().get(30, TimeUnit.SECONDS);
            assertThat(back.getResponseCode("w")).isEqualTo(PlcResponseCode.OK);
            assertThat(back.getInteger("w")).isEqualTo(424242);

            // Restore the original value (0) so the shared server stays pristine for other tests.
            connection.writeRequestBuilder()
                .addTagAddress("w", "ns=3;s=Test/Array/DInt[2]", 0)
                .build().execute().get(30, TimeUnit.SECONDS);
        }
    }

    @Test
    void comprehensiveAddressingVariants() throws Exception {
        // Mirrors the ADS manual test (ManualFactoryAdsDriverTestTC3) against the Plc4xTestNamespace:
        // the same matrix of addressing variants — scalars, whole arrays, array elements/slices and
        // multi-dimensional matrices — driven through the generic BasicPlcTest harness.

        // (A) Read + write coverage for the verified-writable set: scalars, 1-D arrays, 1-D index.
        BasicPlcTest readWrite = new BasicPlcTest(tcpConnectionAddress, new PlcNullAuthentication(),
            true, true, true, true, 3);
        // Scalars
        readWrite.addTestCase("ns=3;s=Test/Scalar/Bool", new PlcBOOL(true));
        readWrite.addTestCase("ns=3;s=Test/Scalar/SInt", new PlcSINT(-12));
        readWrite.addTestCase("ns=3;s=Test/Scalar/USInt", new PlcUSINT(250));
        readWrite.addTestCase("ns=3;s=Test/Scalar/Int", new PlcINT(-1234));
        readWrite.addTestCase("ns=3;s=Test/Scalar/UInt", new PlcUINT(54321));
        readWrite.addTestCase("ns=3;s=Test/Scalar/DInt", new PlcDINT(-12345678));
        readWrite.addTestCase("ns=3;s=Test/Scalar/UDInt", new PlcUDINT(305419896L));
        readWrite.addTestCase("ns=3;s=Test/Scalar/LInt", new PlcLINT(-9223372036854770000L));
        readWrite.addTestCase("ns=3;s=Test/Scalar/ULInt", new PlcULINT(new BigInteger("18446744073709551000")));
        readWrite.addTestCase("ns=3;s=Test/Scalar/Real", new PlcREAL(3.14159f));
        readWrite.addTestCase("ns=3;s=Test/Scalar/LReal", new PlcLREAL(2.718281828459045d));
        readWrite.addTestCase("ns=3;s=Test/Scalar/String", new PlcSTRING("Hello PLC4X"));
        // Whole 1-D arrays
        readWrite.addTestCase("ns=3;s=Test/Array/Bool", new PlcList(List.of(
            new PlcBOOL(true), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(true),
            new PlcBOOL(false), new PlcBOOL(false), new PlcBOOL(true), new PlcBOOL(false))));
        readWrite.addTestCase("ns=3;s=Test/Array/Int", new PlcList(List.of(
            new PlcINT(-3), new PlcINT(-1), new PlcINT(0), new PlcINT(1), new PlcINT(3))));
        readWrite.addTestCase("ns=3;s=Test/Array/UInt", new PlcList(List.of(
            new PlcUINT(1), new PlcUINT(10), new PlcUINT(100), new PlcUINT(1000), new PlcUINT(10000))));
        readWrite.addTestCase("ns=3;s=Test/Array/DInt", new PlcList(List.of(
            new PlcDINT(-1000), new PlcDINT(-500), new PlcDINT(0), new PlcDINT(1000000), new PlcDINT(2000000))));
        readWrite.addTestCase("ns=3;s=Test/Array/LReal", new PlcList(List.of(
            new PlcLREAL(1.5), new PlcLREAL(-2.0), new PlcLREAL(0.125))));
        readWrite.addTestCase("ns=3;s=Test/Array/String", new PlcList(List.of(
            new PlcSTRING("alpha"), new PlcSTRING("beta"), new PlcSTRING("gamma"))));
        // Array element / slice via IndexRange (single element -> scalar, range -> list)
        readWrite.addTestCase("ns=3;s=Test/Array/Int[3]", new PlcINT(1));
        readWrite.addTestCase("ns=3;s=Test/Array/Int[1..3]", new PlcList(List.of(
            new PlcINT(-1), new PlcINT(0), new PlcINT(1))));
        readWrite.addTestCase("ns=3;s=Test/Array/DInt[0]", new PlcDINT(-1000));
        readWrite.addTestCase("ns=3;s=Test/Array/DInt[4]", new PlcDINT(2000000));
        readWrite.run();

        // (B) Read-only coverage for multi-dimensional matrices (whole + element access). The read
        // decode (nested PlcList / element selection) is the interesting part here; multidimensional
        // writes are exercised elsewhere.
        BasicPlcTest readOnly = new BasicPlcTest(tcpConnectionAddress, new PlcNullAuthentication(),
            true, false, true, true, 2);
        readOnly.addTestCase("ns=3;s=Test/Matrix/Int_2x3", new PlcList(List.of(
            new PlcList(List.of(new PlcINT(10), new PlcINT(11), new PlcINT(12))),
            new PlcList(List.of(new PlcINT(-10), new PlcINT(-11), new PlcINT(-12))))));
        readOnly.addTestCase("ns=3;s=Test/Matrix/Real_3x2", new PlcList(List.of(
            new PlcList(List.of(new PlcREAL(1.0f), new PlcREAL(1.5f))),
            new PlcList(List.of(new PlcREAL(2.0f), new PlcREAL(2.5f))),
            new PlcList(List.of(new PlcREAL(3.0f), new PlcREAL(3.5f))))));
        readOnly.addTestCase("ns=3;s=Test/Matrix/UInt_2x2x2", new PlcList(List.of(
            new PlcList(List.of(
                new PlcList(List.of(new PlcUINT(1), new PlcUINT(2))),
                new PlcList(List.of(new PlcUINT(3), new PlcUINT(4))))),
            new PlcList(List.of(
                new PlcList(List.of(new PlcUINT(5), new PlcUINT(6))),
                new PlcList(List.of(new PlcUINT(7), new PlcUINT(8))))))));
        // Matrix element via IndexRange -> scalar.
        readOnly.addTestCase("ns=3;s=Test/Matrix/Int_2x3[0][0]", new PlcINT(10));
        readOnly.addTestCase("ns=3;s=Test/Matrix/Int_2x3[1][2]", new PlcINT(-12));
        readOnly.addTestCase("ns=3;s=Test/Matrix/UInt_2x2x2[1][1][1]", new PlcUINT(8));
        readOnly.run();
    }

    @Test
    void writesCustomStructRoundTrip() throws Exception {
        // A struct is one node whose value is a PlcStruct (OPC UA has no per-field node addressing).
        // Write a PlcStruct, then read it back and confirm every field round-tripped (5d encode).
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            java.util.Map<String, org.apache.plc4x.java.api.value.PlcValue> simple = new java.util.LinkedHashMap<>();
            simple.put("Foo", new PlcSTRING("written-struct"));
            simple.put("Bar", new PlcDINT(987654));
            simple.put("Baz", new PlcBOOL(false));
            simple.put("Qux", new PlcLREAL(-1.25d));

            PlcWriteResponse write = connection.writeRequestBuilder()
                .addTagAddress("s", "ns=3;s=Test/Struct/Simple", new PlcStruct(simple))
                .build().execute().get(30, TimeUnit.SECONDS);
            assertThat(write.getResponseCode("s")).isEqualTo(PlcResponseCode.OK);

            PlcReadResponse read = connection.readRequestBuilder()
                .addTagAddress("s", "ns=3;s=Test/Struct/Simple")
                .build().execute().get(30, TimeUnit.SECONDS);
            org.apache.plc4x.java.api.value.PlcValue value = read.getPlcValue("s");
            assertThat(value.isStruct()).isTrue();
            assertThat(value.getStruct().get("Foo").getString()).isEqualTo("written-struct");
            assertThat(value.getStruct().get("Bar").getInt()).isEqualTo(987654);
            assertThat(value.getStruct().get("Baz").getBoolean()).isFalse();
            assertThat(value.getStruct().get("Qux").getDouble()).isEqualTo(-1.25d);
        }
    }

    @Test
    void customStructIsReadableWithoutCrashing() throws Exception {
        // 5a: a custom (user-defined) struct value is captured as a raw ExtensionObject body instead
        // of crashing the parser (it threw BufferException before). The value round-trips at least as
        // a placeholder even when the field layout can't be resolved on this server.
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcReadResponse response = connection.readRequestBuilder()
                .addTagAddress("s", "ns=3;s=Test/Struct/Simple")
                .build().execute().get(30, TimeUnit.SECONDS);
            assertThat(response.getResponseCode("s")).isEqualTo(PlcResponseCode.OK);
        }
    }

    @Test
    void readsCustomStructAsPlcStruct() throws Exception {
        // Plc4xTestStruct instance in the namespace: foo="hello-struct", bar=12345, baz=true, qux=2.5.
        // Exercises the whole struct read path: 5a (raw body), the field-layout resolution (modern
        // DataTypeDefinition attribute, or the legacy type dictionary against Milo) and 5c (decode).
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            PlcReadResponse response = connection.readRequestBuilder()
                .addTagAddress("s", "ns=3;s=Test/Struct/Simple")
                .build().execute().get(30, TimeUnit.SECONDS);
            assertThat(response.getResponseCode("s")).isEqualTo(PlcResponseCode.OK);

            org.apache.plc4x.java.api.value.PlcValue value = response.getPlcValue("s");
            assertThat(value.isStruct()).isTrue();
            assertThat(value.getStruct().get("Foo").getString()).isEqualTo("hello-struct");
            assertThat(value.getStruct().get("Bar").getInt()).isEqualTo(12345);
            assertThat(value.getStruct().get("Baz").getBoolean()).isTrue();
            assertThat(value.getStruct().get("Qux").getDouble()).isEqualTo(2.5d);
        }
    }

    @Test
    void resolvesAndCachesNodeTypes() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            // getConnection hands back the driver connection directly, so we can reach the
            // per-session type cache to observe Phase-3 behaviour.
            OpcuaConnection opcua = (OpcuaConnection) connection;
            assertThat(opcua.nodeTypeCacheSize()).isZero();

            // A lazy single-node lookup misses the cache, reads the server, and caches the result.
            OpcuaConnection.NodeAttributes boolAttrs = opcua
                .resolveNodeAttributes(OpcuaTag.of("ns=2;s=HelloWorld/ScalarTypes/Boolean"))
                .get(30, TimeUnit.SECONDS);
            assertThat(boolAttrs.dataType())
                .isEqualTo(org.apache.plc4x.java.opcua.readwrite.OpcuaDataType.BOOL);
            assertThat(boolAttrs.readable()).isTrue();
            assertThat(opcua.nodeTypeCacheSize()).isEqualTo(1);

            // Looking the same node up again is served from the cache — no new entry, no read.
            opcua.resolveNodeAttributes(OpcuaTag.of("ns=2;s=HelloWorld/ScalarTypes/Boolean"))
                .get(30, TimeUnit.SECONDS);
            assertThat(opcua.nodeTypeCacheSize()).isEqualTo(1);

            // A different node is resolved (Int32 -> DINT) and added to the cache.
            OpcuaConnection.NodeAttributes intAttrs = opcua
                .resolveNodeAttributes(OpcuaTag.of("ns=2;s=HelloWorld/ScalarTypes/Int32"))
                .get(30, TimeUnit.SECONDS);
            assertThat(intAttrs.dataType())
                .isEqualTo(org.apache.plc4x.java.opcua.readwrite.OpcuaDataType.DINT);
            assertThat(opcua.nodeTypeCacheSize()).isEqualTo(2);

            // Browsing shares the same cache: the two already-resolved children are reused and the
            // remaining ScalarTypes children get resolved and cached, growing the cache further.
            connection.browseRequestBuilder().addQuery("scalars", "ns=2;s=HelloWorld/ScalarTypes")
                .build().execute().get(30, TimeUnit.SECONDS);
            assertThat(opcua.nodeTypeCacheSize()).isGreaterThan(2);
        }
    }

    @Test
    void writeWithoutTypeSuffixUsesServerType() throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress)) {
            // None of these addresses carry a ";TYPE" suffix. Phase 4 resolves each node's real
            // OPC UA type from the server, so the writes use the correct Variant type instead of a
            // lossy Java-value guess. UInt16 is the telling case: 65535 arrives as a Java Integer,
            // which value-inference would map to Int32 and the server would reject — the resolved
            // UInt16 type makes it succeed.
            PlcWriteResponse write = connection.writeRequestBuilder()
                .addTagAddress("bool", BOOL_IDENTIFIER_READ_WRITE, true)
                .addTagAddress("int16", INT16_IDENTIFIER_READ_WRITE, (short) -12345)
                .addTagAddress("uint16", UINT16_IDENTIFIER_READ_WRITE, 65535)
                .addTagAddress("int32", INT32_IDENTIFIER_READ_WRITE, 1234567)
                .addTagAddress("float", FLOAT_IDENTIFIER_READ_WRITE, 3.5f)
                .addTagAddress("string", STRING_IDENTIFIER_READ_WRITE, "phase4")
                .build().execute().get(30, TimeUnit.SECONDS);

            for (String tag : new String[]{"bool", "int16", "uint16", "int32", "float", "string"}) {
                assertThat(write.getResponseCode(tag))
                    .describedAs("write of suffix-less tag '%s'", tag)
                    .isEqualTo(PlcResponseCode.OK);
            }

            // Read the values back (also suffix-less) and confirm they round-tripped correctly.
            PlcReadResponse read = connection.readRequestBuilder()
                .addTagAddress("bool", BOOL_IDENTIFIER_READ_WRITE)
                .addTagAddress("int16", INT16_IDENTIFIER_READ_WRITE)
                .addTagAddress("uint16", UINT16_IDENTIFIER_READ_WRITE)
                .addTagAddress("int32", INT32_IDENTIFIER_READ_WRITE)
                .addTagAddress("float", FLOAT_IDENTIFIER_READ_WRITE)
                .addTagAddress("string", STRING_IDENTIFIER_READ_WRITE)
                .build().execute().get(30, TimeUnit.SECONDS);

            assertThat(read.getBoolean("bool")).isTrue();
            assertThat(read.getShort("int16")).isEqualTo((short) -12345);
            assertThat(read.getInteger("uint16")).isEqualTo(65535);
            assertThat(read.getInteger("int32")).isEqualTo(1234567);
            assertThat(read.getFloat("float")).isEqualTo(3.5f);
            assertThat(read.getString("string")).isEqualTo("phase4");
        }
    }

    private static void collectBrowseItems(java.util.List<org.apache.plc4x.java.api.messages.PlcBrowseItem> items,
                                           java.util.List<org.apache.plc4x.java.api.messages.PlcBrowseItem> out) {
        for (org.apache.plc4x.java.api.messages.PlcBrowseItem item : items) {
            out.add(item);
            collectBrowseItems(new ArrayList<>(item.getChildren().values()), out);
        }
    }

    @Nested
    class SmokeTest {
        @Test
        public void manyReconnectionsWithSingleSubscription() throws Exception {
            PlcDriverManager driverManager = new DefaultPlcDriverManager();
            PlcConnectionManager connectionManager = driverManager.getConnectionManager();

            for (int i = 0; i < 3; i++) {
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
        /**
         * A CYCLIC subscription has to report the value on every interval, even when it never
         * changes - see GH-1102. OPC UA itself only notifies on change, so a static value used to
         * produce exactly one event (the monitored item's initial value) and nothing after that.
         */
        @Test
        public void cyclicSubscriptionReportsRepeatedly() throws Exception {
            PlcConnectionManager connectionManager = new DefaultPlcDriverManager().getConnectionManager();

            try (PlcConnection connection = connectionManager.getConnection(tcpConnectionAddress)) {
                ConcurrentLinkedDeque<PlcSubscriptionEvent> events = new ConcurrentLinkedDeque<>();

                PlcSubscriptionRequest request = connection.subscriptionRequestBuilder()
                    // A value nobody is writing to - without the fix this reports once and stops.
                    .addCyclicTagAddress("static", STRING_IDENTIFIER_READ_WRITE, Duration.ofMillis(500))
                    .build();

                PlcSubscriptionResponse response = request.execute().get(60, TimeUnit.SECONDS);
                assertThat(response.getResponseCode("static")).isEqualTo(PlcResponseCode.OK);
                response.getSubscriptionHandles().forEach(handle -> handle.register(events::add));

                // Five cycles of head room for three expected reports.
                for (int i = 0; i < 50 && events.size() < 3; i++) {
                    Thread.sleep(100);
                }

                assertThat(events.size())
                    .withFailMessage("expected repeated cyclic reports, got %d event(s)", events.size())
                    .isGreaterThanOrEqualTo(3);
                assertThat(events.getLast().getPlcValue("static")).isNotNull();

                connection.unsubscriptionRequestBuilder()
                    .addHandles(response.getSubscriptionHandles())
                    .build()
                    .execute();
            }
        }

        /**
         * A subscription covering several tags is served by a single handle, so
         * getSubscriptionHandles() has to report exactly one - see GH-1896. Reporting it once per
         * tag makes the documented "register a consumer on every handle" loop attach N consumers
         * and deliver every event N times.
         */
        @Test
        public void multiTagSubscriptionReportsOneHandle() throws Exception {
            PlcConnectionManager connectionManager = new DefaultPlcDriverManager().getConnectionManager();

            try (PlcConnection connection = connectionManager.getConnection(tcpConnectionAddress)) {
                PlcSubscriptionRequest request = connection.subscriptionRequestBuilder()
                    .addChangeOfStateTag("first", OpcuaTag.of(INTEGER_IDENTIFIER_READ_WRITE))
                    .addChangeOfStateTag("second", OpcuaTag.of(BOOL_IDENTIFIER_READ_WRITE))
                    .addChangeOfStateTag("third", OpcuaTag.of(INT32_IDENTIFIER_READ_WRITE))
                    .build();

                PlcSubscriptionResponse response = request.execute().get(60, TimeUnit.SECONDS);
                assertThat(response.getResponseCode("first")).isEqualTo(PlcResponseCode.OK);
                assertThat(response.getResponseCode("second")).isEqualTo(PlcResponseCode.OK);
                assertThat(response.getResponseCode("third")).isEqualTo(PlcResponseCode.OK);

                assertThat(response.getSubscriptionHandles()).hasSize(1);

                connection.unsubscriptionRequestBuilder()
                    .addHandles(response.getSubscriptionHandles())
                    .build()
                    .execute();
            }
        }

        @Test
        public void manySubscriptionsOnSingleConnection() throws Exception {
            final int numberOfSubscriptions = 3;

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

            String options = params(
                entry("discovery", "false"),
                entry("server-certificate-file", SERVER_CERTIFICATE.toString().replace("\\", "/")),
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

        @Test
        void securedConnectionWithoutTrustAnchorIsRejected() {
            // No trust-store-file and no server-certificate-file: the driver must fail
            // closed rather than blindly trusting whatever certificate the server presents.
            String options = params(
                entry("discovery", "false"),
                entry("key-store-file", CLIENT_KEY_STORE.toString().replace("\\", "/")),
                entry("key-store-password", "changeit"),
                entry("key-store-type", "pkcs12"),
                entry("security-policy", SecurityPolicy.Basic256Sha256.name()),
                entry("message-security", MessageSecurity.SIGN.name())
            );
            assertThatThrownBy(() ->
                new DefaultPlcDriverManager().getConnection(tcpConnectionAddress + PARAM_DIVIDER + options))
                .isInstanceOf(Exception.class);
        }

        @Test
        void securedConnectionWithInsecureVerificationConnects() throws Exception {
            // With insecure-certificate-verification the driver must use the permissive verifier
            // (it wins over pinning), so trust is not checked. The server certificate is still
            // supplied because an encrypted policy needs the server's public key to encrypt the
            // OpenSecureChannel; only its trust verification is bypassed here.
            String options = params(
                entry("discovery", "false"),
                entry("key-store-file", CLIENT_KEY_STORE.toString().replace("\\", "/")),
                entry("key-store-password", "changeit"),
                entry("key-store-type", "pkcs12"),
                entry("server-certificate-file", SERVER_CERTIFICATE.toString().replace("\\", "/")),
                entry("insecure-certificate-verification", "true"),
                entry("security-policy", SecurityPolicy.Basic256Sha256.name()),
                entry("message-security", MessageSecurity.SIGN.name())
            );
            try (PlcConnection opcuaConnection =
                     new DefaultPlcDriverManager().getConnection(tcpConnectionAddress + PARAM_DIVIDER + options)) {
                assertThat(opcuaConnection.isConnected()).isTrue();

                PlcReadResponse response = opcuaConnection.readRequestBuilder()
                    .addTagAddress("String", STRING_IDENTIFIER_READ_WRITE)
                    .build().execute().get();

                assertThat(response.getResponseCode("String")).isEqualTo(PlcResponseCode.OK);
            }
        }
    }

    @Nested
    class readWrite {
        // Addresses no longer need a ";TYPE" suffix: the driver resolves each node's data type from
        // the server (Phase 4) so writes use the correct OPC UA type without a hint — including the
        // 'UInteger' node, whose abstract UInteger type is resolved to an unsigned type sized to the
        // written value.
        Map<String, Entry<String, Object>> tags = Map.ofEntries(
            entry("Bool", entry(BOOL_IDENTIFIER_READ_WRITE, true)),
            entry("Byte", entry(BYTE_IDENTIFIER_READ_WRITE, (short) 3)),
            entry("Double", entry(DOUBLE_IDENTIFIER_READ_WRITE, 0.5d)),
            entry("Float", entry(FLOAT_IDENTIFIER_READ_WRITE, 0.5f)),
            entry("Int16", entry(INT16_IDENTIFIER_READ_WRITE, 1)),
            entry("Int32", entry(INT32_IDENTIFIER_READ_WRITE, 42)),
            entry("Int64", entry(INT64_IDENTIFIER_READ_WRITE, 42L)),
            entry("Integer", entry(INTEGER_IDENTIFIER_READ_WRITE, -127)),
            //entry("SByte", entry(SBYTE_IDENTIFIER_READ_WRITE, )),
            entry("String", entry(STRING_IDENTIFIER_READ_WRITE, "Hello Toddy!")),
            entry("UInt16", entry(UINT16_IDENTIFIER_READ_WRITE, 65535)),
            entry("UInt32", entry(UINT32_IDENTIFIER_READ_WRITE, 101010101L)),
            entry("UInt64", entry(UINT64_IDENTIFIER_READ_WRITE, new BigInteger("1337"))),
            entry("UInteger", entry(UINTEGER_IDENTIFIER_READ_WRITE, 102020202L)),
            entry("BooleanArray", entry(BOOL_ARRAY_IDENTIFIER, new boolean[]{true, true, true, true, true})),
            // entry("ByteStringArray", entry(BYTE_STRING_ARRAY_IDENTIFIER, null)),
            entry("ByteArray", entry(BYTE_ARRAY_IDENTIFIER, new Short[]{1, 100, 100, 255, 123})),
            entry("DoubleArray", entry(DOUBLE_ARRAY_IDENTIFIER, new Double[]{1.0, 2.0, 3.0, 4.0, 5.0})),
            entry("FloatArray", entry(FLOAT_ARRAY_IDENTIFIER, new Float[]{1.0F, 2.0F, 3.0F, 4.0F, 5.0F})),
            entry("Int16Array", entry(INT16_ARRAY_IDENTIFIER, new Short[]{1, 2, 3, 4, 5})),
            entry("Int32Array", entry(INT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5})),
            entry("Int64Array", entry(INT64_ARRAY_IDENTIFIER, new Long[]{1L, 2L, 3L, 4L, 5L})),
            entry("IntegerArray", entry(INT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5})),
            entry("SByteArray", entry(SBYTE_ARRAY_IDENTIFIER, new Byte[]{1, 2, 3, 4, 5})),
            entry("StringArray", entry(STRING_ARRAY_IDENTIFIER, new String[]{"1", "2", "3", "4", "5"})),
            entry("UInt16Array", entry(UINT16_ARRAY_IDENTIFIER, new Short[]{1, 2, 3, 4, 5})),
            entry("UInt32Array", entry(UINT32_ARRAY_IDENTIFIER, new Integer[]{1, 2, 3, 4, 5})),
            entry("UInt64Array", entry(UINT64_ARRAY_IDENTIFIER, new Long[]{1L, 2L, 3L, 4L, 5L})),
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
            PlcReadResponse response = request.execute().get(30, TimeUnit.SECONDS);

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
            PlcWriteResponse response = request.execute().get(30, TimeUnit.SECONDS);

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
            private final CountDownLatch latch;

            public ReadWorker(PlcConnection opcuaConnection, CountDownLatch latch) {
                this.connection = opcuaConnection;
                this.latch = latch;
            }

            @Override
            public void run() {
                try {
                    PlcReadRequest.Builder read_builder = connection.readRequestBuilder();
                    read_builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE);
                    PlcReadRequest read_request = read_builder.build();

                    for (int i = 0; i < 10; i++) {
                        PlcReadResponse read_response = read_request.execute().get();
                        assertThat(read_response.getResponseCode("Bool")).isEqualTo(PlcResponseCode.OK);
                    }
                } catch (ExecutionException e) {
                    LOGGER.error("run aborted", e);
                } catch (InterruptedException e) {
                    LOGGER.error("thread interrupted", e);
                    Thread.currentThread().interrupt();
                } finally {
                    this.latch.countDown();
                }
            }
        }

        class WriteWorker extends Thread {
            private final PlcConnection connection;
            private final CountDownLatch latch;

            public WriteWorker(PlcConnection opcuaConnection, CountDownLatch latch) {
                this.connection = opcuaConnection;
                this.latch = latch;
            }

            @Override
            public void run() {
                try {
                    PlcWriteRequest.Builder write_builder = connection.writeRequestBuilder();
                    write_builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE, true);
                    PlcWriteRequest write_request = write_builder.build();

                    for (int i = 0; i < 10; i++) {
                        PlcWriteResponse write_response = write_request.execute().get();
                        assertThat(write_response.getResponseCode("Bool")).isEqualTo(PlcResponseCode.OK);
                    }
                } catch (ExecutionException e) {
                    LOGGER.error("run aborted", e);
                } catch (InterruptedException e) {
                    LOGGER.error("thread interrupted", e);
                    Thread.currentThread().interrupt();
                } finally {
                    this.latch.countDown();
                }
            }
        }


        PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
        Condition<PlcConnection> is_connected = new Condition<>(PlcConnection::isConnected, "is connected");
        assertThat(opcuaConnection).is(is_connected);

        CountDownLatch latch = new CountDownLatch(2);
        ReadWorker read_worker = new ReadWorker(opcuaConnection, latch);
        WriteWorker write_worker = new WriteWorker(opcuaConnection, latch);
        read_worker.start();
        write_worker.start();

        latch.await();

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
                    // Pin trust to the server certificate; the driver rejects unknown certs by default.
                    entry("server-certificate-file", SERVER_CERTIFICATE.toString().replace("\\", "/")),
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
