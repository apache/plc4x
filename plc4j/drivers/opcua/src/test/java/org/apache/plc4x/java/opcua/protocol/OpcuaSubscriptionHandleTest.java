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
package org.apache.plc4x.java.opcua.protocol;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.OpcuaPlcDriverTest;
import org.apache.plc4x.test.DisableOnParallelsVmFlag;
import org.eclipse.milo.examples.server.ExampleServer;
import org.junit.jupiter.api.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

import static org.junit.jupiter.api.Assumptions.assumeTrue;

@DisableOnParallelsVmFlag
public class OpcuaSubscriptionHandleTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaPlcDriverTest.class);

    private static ExampleServer exampleServer;

    // Address of local milo server
    private static String miloLocalAddress = "127.0.0.1:12686/milo";
    //Tcp pattern of OPC UA
    private static String opcPattern = "opcua:tcp://";

    private String paramSectionDivider = "?";
    private String paramDivider = "&";

    private static String tcpConnectionAddress = opcPattern + miloLocalAddress;

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

    private static PlcConnection opcuaConnection;

    @BeforeEach
    public void before() {
    }

    @AfterEach
    public void after() {

    }

    @BeforeAll
    public static void setup() {
        try {
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
            //Connect
            opcuaConnection = new PlcDriverManager().getConnection(tcpConnectionAddress);
            assert opcuaConnection.isConnected();
        } catch (Exception e) {
            e.printStackTrace();
            try {
                exampleServer.shutdown().get();
            } catch (Exception j) {
                j.printStackTrace();
            }
        }
    }

    @AfterAll
    public static void tearDown() {
        try {
            // Close Connection
            opcuaConnection.close();
            assert !opcuaConnection.isConnected();

            exampleServer.shutdown().get();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Test
    public void subscribeBool() throws Exception {
        String field = "Bool";
        String identifier = BOOL_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeByte() throws Exception {
        String field = "Byte";
        String identifier = BYTE_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeDouble() throws Exception {
        String field = "Double";
        String identifier = DOUBLE_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeFloat() throws Exception {
        String field = "Float";
        String identifier = FLOAT_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeInt16() throws Exception {
        String field = "Int16";
        String identifier = INT16_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeInt32() throws Exception {
        String field = "Int32";
        String identifier = INT32_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeInt64() throws Exception {
        String field = "Int64";
        String identifier = INT64_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeInteger() throws Exception {
        String field = "Integer";
        String identifier = INTEGER_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeSByte() throws Exception {
        String field = "SByte";
        String identifier = SBYTE_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeString() throws Exception {
        String field = "String";
        String identifier = STRING_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeUInt16() throws Exception {
        String field = "Uint16";
        String identifier = UINT16_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeUInt32() throws Exception {
        String field = "UInt32";
        String identifier = UINT32_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeUInt64() throws Exception {
        String field = "UInt64";
        String identifier = UINT64_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeUInteger() throws Exception {
        String field = "UInteger";
        String identifier = UINTEGER_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field).equals(PlcResponseCode.OK);
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeDoesNotExists() throws Exception {
        String field = "DoesNotExists";
        String identifier = DOES_NOT_EXIST_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", field);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            //This should never be called,
            assert false;
            LOGGER.info("Received a response from {} test {}", field, plcSubscriptionEvent.getPlcValue(field).toString());
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeMultiple() throws Exception {
        String field1 = "UInteger";
        String identifier1 = UINTEGER_IDENTIFIER_READ_WRITE;
        String field2 = "Integer";
        String identifier2 = INTEGER_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {}  and {} test", field1, field2);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateFieldAddress(field1, identifier1);
        builder.addChangeOfStateFieldAddress(field2, identifier2);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get();
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(field1);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(field1).equals(PlcResponseCode.OK);
            assert plcSubscriptionEvent.getResponseCode(field2).equals(PlcResponseCode.OK);
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

}
