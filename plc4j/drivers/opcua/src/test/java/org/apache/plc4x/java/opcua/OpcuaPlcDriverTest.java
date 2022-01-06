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
import org.apache.plc4x.java.opcuaserver.OPCUAServer;
import org.assertj.core.api.Condition;
import org.eclipse.milo.examples.server.ExampleServer;
import org.eclipse.milo.opcua.sdk.server.OpcUaServer;
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
import java.nio.file.StandardCopyOption;
import java.util.concurrent.CompletableFuture;

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
    private String miloLocalAddress = "127.0.0.1:12686/plc4x";
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

    private static OpcUaServer exampleServer;

    @BeforeAll
    public static void setup() {
        try {
            // When switching JDK versions from a newer to an older version,
            // this can cause the server to not start correctly.
            // Deleting the directory makes sure the key-store is initialized correctly.
            Path baseDir = Paths.get(System.getProperty("java.io.tmpdir"), "server");
            try {
                Files.delete(baseDir);
            } catch (Exception e) {
                // Ignore this ...
            }

            System.out.println("Working Directory = " + System.getProperty("user.dir"));
            Files.copy(Path.of("src/test/resources/config.yml"), baseDir, StandardCopyOption.REPLACE_EXISTING);

            String[] args = {"", "-c", baseDir + "/config.yml", "-t"};
            OPCUAServer serverInit = new OPCUAServer(args);
            exampleServer = serverInit.getServer();
            serverInit.getServer().startup().get();
            CompletableFuture<Void> future = new CompletableFuture<>();
            future.get();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @AfterAll
    public static void tearDown() {
        try {
            exampleServer.shutdown().get();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Test
    public void connectionNoParams(){
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
    public void connectionWithDiscoveryParam(){
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
            builder.addItem("Bool", BOOL_IDENTIFIER_READ_WRITE);
            builder.addItem("Byte", BYTE_IDENTIFIER_READ_WRITE);
            builder.addItem("Double", DOUBLE_IDENTIFIER_READ_WRITE);
            builder.addItem("Float", FLOAT_IDENTIFIER_READ_WRITE);
            builder.addItem("Int16", INT16_IDENTIFIER_READ_WRITE);
            builder.addItem("Int32", INT32_IDENTIFIER_READ_WRITE);
            builder.addItem("Int64", INT64_IDENTIFIER_READ_WRITE);
            builder.addItem("Integer", INTEGER_IDENTIFIER_READ_WRITE);
            builder.addItem("SByte", SBYTE_IDENTIFIER_READ_WRITE);
            builder.addItem("String", STRING_IDENTIFIER_READ_WRITE);
            builder.addItem("UInt16", UINT16_IDENTIFIER_READ_WRITE);
            builder.addItem("UInt32", UINT32_IDENTIFIER_READ_WRITE);
            builder.addItem("UInt64", UINT64_IDENTIFIER_READ_WRITE);
            builder.addItem("UInteger", UINTEGER_IDENTIFIER_READ_WRITE);

            builder.addItem("BoolArray", BOOL_ARRAY_IDENTIFIER);
            //builder.addItem("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER);
            builder.addItem("ByteArray", BYTE_ARRAY_IDENTIFIER);
            builder.addItem("DoubleArray", DOUBLE_ARRAY_IDENTIFIER);
            builder.addItem("FloatArray", FLOAT_ARRAY_IDENTIFIER);
            builder.addItem("Int16Array", INT16_ARRAY_IDENTIFIER);
            builder.addItem("Int32Array", INT32_ARRAY_IDENTIFIER);
            builder.addItem("Int64Array", INT64_ARRAY_IDENTIFIER);
            builder.addItem("SByteArray", SBYTE_ARRAY_IDENTIFIER);
            builder.addItem("StringArray", STRING_ARRAY_IDENTIFIER);
            builder.addItem("UInt16Array", UINT16_ARRAY_IDENTIFIER);
            builder.addItem("UInt32Array", UINT32_ARRAY_IDENTIFIER);
            builder.addItem("UInt64Array", UINT64_ARRAY_IDENTIFIER);

            builder.addItem("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE);

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
        builder.addItem("Bool", BOOL_IDENTIFIER_READ_WRITE, true);
        builder.addItem("Byte", BYTE_IDENTIFIER_READ_WRITE + ";BYTE", (short) 3);
        builder.addItem("Double", DOUBLE_IDENTIFIER_READ_WRITE, 0.5d);
        builder.addItem("Float", FLOAT_IDENTIFIER_READ_WRITE, 0.5f);
        //builder.addItem("Int16", INT16_IDENTIFIER_READ_WRITE + "", (short) 1);
        builder.addItem("Int32", INT32_IDENTIFIER_READ_WRITE, 42);
        builder.addItem("Int64", INT64_IDENTIFIER_READ_WRITE, 42L);
        builder.addItem("Integer", INTEGER_IDENTIFIER_READ_WRITE, 42);
        builder.addItem("SByte", SBYTE_IDENTIFIER_READ_WRITE + ";SINT", -127);
        builder.addItem("String", STRING_IDENTIFIER_READ_WRITE, "Helllo Toddy!");
        builder.addItem("UInt16", UINT16_IDENTIFIER_READ_WRITE + ";UINT", 65535);
        builder.addItem("UInt32", UINT32_IDENTIFIER_READ_WRITE + ";UDINT", 101010101L);
        builder.addItem("UInt64", UINT64_IDENTIFIER_READ_WRITE + ";ULINT", new BigInteger("1337"));
        builder.addItem("UInteger", UINTEGER_IDENTIFIER_READ_WRITE + ";UDINT", 102020202L);


        builder.addItem("BooleanArray", BOOL_ARRAY_IDENTIFIER, new Boolean[] {true, true, true, true, true});
        builder.addItem("ByteArray", BYTE_ARRAY_IDENTIFIER + ";BYTE", new Short[] {1, 100, 100, 255, 123});
        builder.addItem("DoubleArray", DOUBLE_ARRAY_IDENTIFIER, new Double[] {1.0,2.0,3.0,4.0,5.0});
        builder.addItem("FloatArray", FLOAT_ARRAY_IDENTIFIER, new Float[] {1.0F,2.0F,3.0F,4.0F,5.0F});
        builder.addItem("Int16Array", INT16_ARRAY_IDENTIFIER, new Short[] {1,2,3,4,5});
        builder.addItem("Int32Array", INT32_ARRAY_IDENTIFIER, new Integer[] {1,2,3,4,5});
        builder.addItem("Int64Array", INT64_ARRAY_IDENTIFIER, new Long[] {1L,2L,3L,4L,5L});
        builder.addItem("IntegerArray", INT32_ARRAY_IDENTIFIER, new Integer[] {1,2,3,4,5});
        builder.addItem("SByteArray", SBYTE_ARRAY_IDENTIFIER, new Byte[] {1,2,3,4,5});
        builder.addItem("StringArray", STRING_ARRAY_IDENTIFIER, new String[] {"1","2","3","4","5"});
        builder.addItem("UInt16Array", UINT16_ARRAY_IDENTIFIER + ";UINT", new Short[] {1,2,3,4,5});
        builder.addItem("UInt32Array", UINT32_ARRAY_IDENTIFIER + ";UDINT", new Integer[] {1,2,3,4,5});
        builder.addItem("UInt64Array", UINT64_ARRAY_IDENTIFIER + ";ULINT", new Long[] {1L,2L,3L,4L,5L});

        builder.addItem("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE, "11");

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

}
