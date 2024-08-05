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

import java.io.ByteArrayOutputStream;
import java.util.concurrent.CountDownLatch;
import java.util.stream.Stream;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.MiloTestContainer;
import org.apache.plc4x.java.opcua.OpcuaPlcDriverTest;
import org.apache.plc4x.java.opcua.KeystoreGenerator;
import org.junit.jupiter.api.*;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.Arguments;
import org.junit.jupiter.params.provider.MethodSource;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;
import org.testcontainers.containers.output.Slf4jLogConsumer;
import org.testcontainers.images.builder.Transferable;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

// ! For some odd reason does this test not work on VMs running in Parallels.
// cdutz: I have done way more than my fair share on tracking down this issue and am simply giving up on it.
// I tracked it down into the core of Milo several times now, but got lost in there.
// It's not a big issue as the GitHub runners and the Apache Jenkins still run the test.
@Testcontainers(disabledWithoutDocker = true)
public class OpcuaSubscriptionHandleTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpcuaPlcDriverTest.class);

    @Container
    public final MiloTestContainer milo = new MiloTestContainer()
        .withLogConsumer(new Slf4jLogConsumer(LOGGER));

    // Address of local milo server
    private static final String miloLocalAddress = "%s:%d/milo";
    //Tcp pattern of OPC UA
    private static final String opcPattern = "opcua:tcp://";

    private final String paramSectionDivider = "?";
    private final String paramDivider = "&";

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

    private PlcConnection opcuaConnection;

    // ! If this test fails, see comment at the top of the class before investigating.
    @BeforeEach
    public void setup() throws Exception {
        // When switching JDK versions from a newer to an older version,
        // this can cause the server to not start correctly.
        // Deleting the directory makes sure the key-store is initialized correctly.
        String tcpConnectionAddress = String.format(opcPattern + miloLocalAddress, milo.getHost(), milo.getMappedPort(12686)) + "?endpoint-port=12686";
        //Connect
        opcuaConnection = new DefaultPlcDriverManager().getConnection(tcpConnectionAddress);
        assertThat(opcuaConnection).extracting(PlcConnection::isConnected).isEqualTo(true);
    }

    @AfterEach
    public void tearDown() throws Exception {
        // Close Connection
        opcuaConnection.close();
        assertThat(opcuaConnection).extracting(PlcConnection::isConnected).isEqualTo(false);
    }

    // ! If this test fails, see comment at the top of the class before investigating.
    @Test
    public void subscribeDoesNotExists() throws Exception {
        String tag = "DoesNotExists";
        String identifier = DOES_NOT_EXIST_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {} test", tag);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateTagAddress(tag, identifier);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get(1000, TimeUnit.MILLISECONDS);
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(tag);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            //This should never be called,
            fail("Received subscription response whereas error was expected");
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @Test
    public void subscribeEvent() throws Exception {
        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addEventTagAddress("ev1", "ns=0;i=2253");
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get(1000, TimeUnit.MILLISECONDS);
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle("ev1");

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            System.out.println("Subscription handle " + plcSubscriptionEvent);
            assert plcSubscriptionEvent.getResponseCode("ev1").equals(PlcResponseCode.OK);
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    // ! If this test fails, see comment at the top of the class before investigating.
    @Test
    public void subscribeMultiple() throws Exception {
        String tag1 = "UInteger";
        String identifier1 = UINTEGER_IDENTIFIER_READ_WRITE;
        String tag2 = "Integer";
        String identifier2 = INTEGER_IDENTIFIER_READ_WRITE;
        LOGGER.info("Starting subscription {}  and {} test", tag1, tag2);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateTagAddress(tag1, identifier1);
        builder.addChangeOfStateTagAddress(tag2, identifier2);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get(1000, TimeUnit.MILLISECONDS);
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(tag1);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(tag1).equals(PlcResponseCode.OK);
            assert plcSubscriptionEvent.getResponseCode(tag2).equals(PlcResponseCode.OK);
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    // ! If this test fails, see comment at the top of the class before investigating.
    @Test
    public void subscribeMultipleWithOneMissing() throws Exception {
        String tag1 = "UInteger";
        String identifier1 = UINTEGER_IDENTIFIER_READ_WRITE;
        String tag2 = "Integer";
        String identifier2 = UINTEGER_IDENTIFIER_READ_WRITE + "_MISSING_GONE";
        LOGGER.info("Starting subscription {}  and {} test", tag1, tag2);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateTagAddress(tag1, identifier1);
        builder.addChangeOfStateTagAddress(tag2, identifier2);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get(10000, TimeUnit.MILLISECONDS);
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(tag1);

        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            assert plcSubscriptionEvent.getResponseCode(tag1).equals(PlcResponseCode.OK);
            assert plcSubscriptionEvent.getResponseCode(tag2).equals(PlcResponseCode.NOT_FOUND);
        });

        //Wait for value to be returned from server
        Thread.sleep(1200);

        subscriptionHandle.stopSubscriber();
    }

    @ParameterizedTest
    @MethodSource("getTags")
    public void subscribeTest(String tag, Class<?> type) throws Exception {
        LOGGER.info("Starting subscription {} test", tag);

        // Create Subscription
        PlcSubscriptionRequest.Builder builder = opcuaConnection.subscriptionRequestBuilder();
        builder.addChangeOfStateTagAddress(tag, tag);
        PlcSubscriptionRequest request = builder.build();

        // Get result of creating subscription
        PlcSubscriptionResponse response = request.execute().get(1000, TimeUnit.MILLISECONDS);
        final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) response.getSubscriptionHandle(tag);

        CountDownLatch latch = new CountDownLatch(1);
        // Create handler for returned value
        subscriptionHandle.register(plcSubscriptionEvent -> {
            Object value = plcSubscriptionEvent.getObject(tag);
            LOGGER.info("Received a response from {} test {} ({})", tag, plcSubscriptionEvent.getPlcValue(tag).toString(), value.getClass());
            assertEquals(PlcResponseCode.OK, plcSubscriptionEvent.getResponseCode(tag));
            assertNotNull(value);
            assertTrue(type.isInstance(value));
            latch.countDown();
        });

        assertTrue(latch.await(1200, TimeUnit.MILLISECONDS));
        subscriptionHandle.stopSubscriber();
    }

    private static Stream<Arguments> getTags() {
        return Stream.of(
            Arguments.of(BOOL_IDENTIFIER_READ_WRITE, Boolean.class),
            Arguments.of(BYTE_IDENTIFIER_READ_WRITE, Short.class),
            Arguments.of(DOUBLE_IDENTIFIER_READ_WRITE, Double.class),
            Arguments.of(FLOAT_IDENTIFIER_READ_WRITE, Float.class),
            Arguments.of(INT16_IDENTIFIER_READ_WRITE, Short.class),
            Arguments.of(INT32_IDENTIFIER_READ_WRITE, Integer.class),
            Arguments.of(INT64_IDENTIFIER_READ_WRITE, Long.class),
            Arguments.of(INTEGER_IDENTIFIER_READ_WRITE, Integer.class),
            Arguments.of(SBYTE_IDENTIFIER_READ_WRITE, byte[].class),
            Arguments.of(STRING_IDENTIFIER_READ_WRITE, String.class),
            Arguments.of(UINT16_IDENTIFIER_READ_WRITE, Integer.class),
            Arguments.of(UINT32_IDENTIFIER_READ_WRITE, Long.class),
            Arguments.of(UINT64_IDENTIFIER_READ_WRITE, Long.class),
            Arguments.of(UINTEGER_IDENTIFIER_READ_WRITE, Long.class)
        );
    }

}
