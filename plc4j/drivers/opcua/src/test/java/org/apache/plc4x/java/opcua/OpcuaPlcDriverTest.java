/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/
package org.apache.plc4x.java.opcua;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.connection.OpcuaTcpPlcConnection;
import org.eclipse.milo.examples.server.ExampleServer;
import org.junit.Assert;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.INET_ADDRESS_PATTERN;
import static org.apache.plc4x.java.opcua.OpcuaPlcDriver.OPCUA_URI_PATTERN;
import static org.apache.plc4x.java.opcua.UtilsTest.assertMatching;

/**
 */
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
    private static final String LOCALICED_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/LocalizedText";
    private static final String NODE_ID_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/NodeId";
    private static final String QUALIFIED_NAM_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/QualifiedName";
    private static final String UTC_TIME_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/UtcTime";
    private static final String VARIANT_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/Variant";
    private static final String XML_ELEMENT_READ_WRITE = "ns=2;s=HelloWorld/ScalarTypes/XmlElement";
    // Address of local milo server
    String miloLocalAddress = "127.0.0.1:12686/milo";
    //Tcp pattern of OPC UA
    String opcPattern = "opcua:tcp://";
    private ExampleServer exampleServer;



    @BeforeEach
    public void before() {
        try {
            exampleServer = new ExampleServer();
            exampleServer.startup().get();
        } catch (Exception e) {

        }
    }

    @AfterEach
    public void after() {
        try {
            exampleServer.shutdown().get();
        } catch (Exception e) {

        }
    }

    @Test
    public void readVariables() {
        try {
            OpcuaTcpPlcConnection opcuaConnection = (OpcuaTcpPlcConnection)
                new PlcDriverManager().getConnection(opcPattern + miloLocalAddress);
            assert opcuaConnection.isConnected();

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

            builder.addItem("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER_READ_WRITE);

            PlcReadRequest request = builder.build();
            PlcReadResponse response = opcuaConnection.read(request).get();
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


        } catch (Exception e) {
            Assert.fail("Exception during readVariables Test EXCEPTION: " + e.getMessage());
        }
    }

    @Test
    public void testOpcuaAddressPattern() {

        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://localhost:3131");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://www.google.de:443");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://127.0.0.1");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://127.0.0.1:251");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://254.254.254.254:1337");
        assertMatching(INET_ADDRESS_PATTERN, "tcp://254.254.254.254");


        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://localhost");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://localhost:3131");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://www.google.de");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://www.google.de:443");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://127.0.0.1");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://127.0.0.1:251");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://254.254.254.254:1337");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://254.254.254.254");

        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://127.0.0.1&discovery=false");
        assertMatching(OPCUA_URI_PATTERN, "opcua:tcp://opcua.demo-this.com:51210/UA/SampleServer&discovery=false");

    }

}
