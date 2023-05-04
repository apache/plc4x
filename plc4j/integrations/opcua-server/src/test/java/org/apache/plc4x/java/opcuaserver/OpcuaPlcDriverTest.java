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
package org.apache.plc4x.java.opcuaserver;

import io.vavr.collection.List;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.*;

import java.math.BigInteger;


import org.apache.commons.io.FileUtils;
import java.io.File;

import static org.assertj.core.api.Assertions.fail;

/**
 */
@Disabled("We're getting strange OutOfMemoryErrors from this one")
public class OpcuaPlcDriverTest {
    // Read only variables of milo example server of version 3.6
    private static final String BOOL_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_BOOL";
    private static final String BYTE_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_BYTE";
    private static final String DOUBLE_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_LREAL";
    private static final String FLOAT_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_REAL";
    private static final String INT16_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_INT";
    private static final String INT32_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_DINT";
    private static final String INT64_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_LINT";
    private static final String INTEGER_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_DINT";
    private static final String SBYTE_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_USINT";
    private static final String STRING_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_STRING";
    private static final String UINT16_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_UINT";
    private static final String UINT32_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_UDINT";
    private static final String UINT64_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_ULINT";
    private static final String UINTEGER_IDENTIFIER_READ_WRITE = "ns=1;s=Simulated_UDINT";
    private static final String DOES_NOT_EXIST_IDENTIFIER_READ_WRITE = "ns=1;i=12512623";
    // At the moment not used in PLC4X or in the OPC UA driver
    private static final String BYTE_STRING_IDENTIFIER_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/ByteString";
    private static final String DATE_TIME_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/DateTime";
    private static final String DURATION_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/Duration";
    private static final String GUID_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/Guid";
    private static final String LOCALISED_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/LocalizedText";
    private static final String NODE_ID_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/NodeId";
    private static final String QUALIFIED_NAM_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/QualifiedName";
    private static final String UTC_TIME_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/UtcTime";
    private static final String VARIANT_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/Variant";
    private static final String XML_ELEMENT_READ_WRITE = "ns=1;s=HelloWorld/ScalarTypes/XmlElement";
    // Address of local milo server
    private final String miloLocalAddress = "127.0.0.1:12673/plc4x";
    //Tcp pattern of OPC UA
    private final String opcPattern = "opcua:tcp://";

    private final String paramSectionDivider = "?";
    private String paramDivider = "&";

    private final String tcpConnectionAddress = opcPattern + miloLocalAddress;

    private final List<String> connectionStringValidSet = List.of(tcpConnectionAddress);
    private List<String> connectionStringCorruptedSet = List.of();

    private final String discoveryValidParamTrue = "discovery=true";
    private final String discoveryValidParamFalse = "discovery=false";
    private final String discoveryCorruptedParamWrongValueNum = "discovery=1";
    private final String discoveryCorruptedParamWronName = "diskovery=false";

    final List<String> discoveryParamValidSet = List.of(discoveryValidParamTrue, discoveryValidParamFalse);
    List<String> discoveryParamCorruptedSet = List.of(discoveryCorruptedParamWrongValueNum, discoveryCorruptedParamWronName);

    private static OPCUAServer exampleServer;

    @BeforeAll
    public static void setup() throws Exception{
        FileUtils.deleteDirectory(new File("target/test-tmp"));
        exampleServer = new OPCUAServer(new String[]{"-c", "src/test/resources/config.yml", "-i", "-t"});
        exampleServer.startup().get();
    }

    @AfterAll
    public static void tearDown() {
        try {
            exampleServer.shutdown().get();
        } catch (Exception e) {

        }
    }

    @Test
    public void connectionNoParams(){
        connectionStringValidSet.forEach(connectionAddress -> {
                String connectionString = connectionAddress;
                try {
                    PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
                    assert opcuaConnection.isConnected();
                    opcuaConnection.close();
                    assert !opcuaConnection.isConnected();
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
                    PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(connectionString);
                    assert opcuaConnection.isConnected();
                    opcuaConnection.close();
                    assert !opcuaConnection.isConnected();
                } catch (PlcConnectionException e) {
                    fail("Exception during connectionWithDiscoveryParam while connecting Test EXCEPTION: " + e.getMessage());
                } catch (Exception e) {
                    fail("Exception during connectionWithDiscoveryParam while closing Test EXCEPTION: " + e.getMessage());
                }
            });
        });


    }

    @Test
    public void readVariables() throws Exception{

            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
            assert opcuaConnection.isConnected();

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

            builder.addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE);

            PlcReadRequest request = builder.build();
            PlcReadResponse response = request.execute().get();

            assert response.getResponseCode("Bool").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Byte").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Double").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Float").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Int16").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Int32").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Int64").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Integer").equals(PlcResponseCode.OK);
            assert response.getResponseCode("SByte").equals(PlcResponseCode.OK);
            assert response.getResponseCode("String").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInt16").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInt32").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInt64").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInteger").equals(PlcResponseCode.OK);

            assert response.getResponseCode("DoesNotExists").equals(PlcResponseCode.NOT_FOUND);

            opcuaConnection.close();
            assert !opcuaConnection.isConnected();

    }

    @Test
    public void writeVariables() throws Exception {

            PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
            assert opcuaConnection.isConnected();

            PlcWriteRequest.Builder builder = opcuaConnection.writeRequestBuilder();
            builder.addTagAddress("Bool", BOOL_IDENTIFIER_READ_WRITE, true);
            builder.addTagAddress("Byte", BYTE_IDENTIFIER_READ_WRITE + ";BYTE", 255);
            builder.addTagAddress("Double", DOUBLE_IDENTIFIER_READ_WRITE, 0.5d);
            builder.addTagAddress("Float", FLOAT_IDENTIFIER_READ_WRITE, 0.5f);
            builder.addTagAddress("Int16", INT16_IDENTIFIER_READ_WRITE + ";INT", 32000);
            builder.addTagAddress("Int32", INT32_IDENTIFIER_READ_WRITE, 42);
            builder.addTagAddress("Int64", INT64_IDENTIFIER_READ_WRITE, 42L);
            builder.addTagAddress("Integer", INTEGER_IDENTIFIER_READ_WRITE, 42);
            builder.addTagAddress("SByte", SBYTE_IDENTIFIER_READ_WRITE + ";USINT", 100);
            builder.addTagAddress("String", STRING_IDENTIFIER_READ_WRITE, "Helllo Toddy!");
            builder.addTagAddress("UInt16", UINT16_IDENTIFIER_READ_WRITE + ";UINT", 65535);
            builder.addTagAddress("UInt32", UINT32_IDENTIFIER_READ_WRITE + ";UDINT", 100);
            builder.addTagAddress("UInt64", UINT64_IDENTIFIER_READ_WRITE + ";ULINT", new BigInteger("1337"));
            builder.addTagAddress("UInteger", UINTEGER_IDENTIFIER_READ_WRITE + ";UDINT", 100);

            builder.addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE, "Sad Toddy");

            PlcWriteRequest request = builder.build();
            PlcWriteResponse response = request.execute().get();
            assert response.getResponseCode("Bool").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Byte").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Double").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Float").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Int16").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Int32").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Int64").equals(PlcResponseCode.OK);
            assert response.getResponseCode("Integer").equals(PlcResponseCode.OK);
            assert response.getResponseCode("SByte").equals(PlcResponseCode.OK);
            assert response.getResponseCode("String").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInt16").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInt32").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInt64").equals(PlcResponseCode.OK);
            assert response.getResponseCode("UInteger").equals(PlcResponseCode.OK);

            assert response.getResponseCode("DoesNotExists").equals(PlcResponseCode.NOT_FOUND);

            opcuaConnection.close();
            assert !opcuaConnection.isConnected();
    }

}
