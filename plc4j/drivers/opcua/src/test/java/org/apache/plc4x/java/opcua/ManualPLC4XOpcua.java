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

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.protocol.OpcuaSubscriptionHandle;
import org.eclipse.milo.examples.server.ExampleServer;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Optional;

/**
 * This class serves only as a manual entry point for ad-hoc tests of the OPC UA PLC4J driver.
 * <p>
 * <p>
 * The current version is tested against a public server, which is to be replaced later by a separate instance of the Milo framework.
 * Afterwards the code represented here will be used as an example for the introduction page.
 * <p>
 */
public class ManualPLC4XOpcua {
    private static final String BOOL_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Boolean";
    private static final String BYTE_STRING_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/ByteString";
    private static final String BYTE_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Byte";
    private static final String DOUBLE_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Double";
    private static final String FLOAT_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Float";
    private static final String INT16_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Int16";
    private static final String INT32_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Int32";
    private static final String INT64_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Int64";
    private static final String INTEGER_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/Integer";
    private static final String SBYTE_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/SByte";
    private static final String STRING_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/String";
    private static final String UINT16_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/UInt16";
    private static final String UINT32_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/UInt32";
    private static final String UINT64_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/UInt64";
    private static final String UINTEGER_IDENTIFIER = "ns=2;s=HelloWorld/ScalarTypes/UInteger";

    //Arrays
    private static final String BOOL_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/BooleanArray";
    private static final String BYTE_STRING_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/ByteStringArray";
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
    private static final String UINTEGER_ARRAY_IDENTIFIER = "ns=2;s=HelloWorld/ArrayTypes/UIntegerArray";

    //Don't exists
    private static final String DOES_NOT_EXIST_IDENTIFIER = "ns=2;i=12512623";

    public static void main(String... args) throws Exception {
        // When switching JDK versions from a newer to an older version,
        // this can cause the server to not start correctly.
        // Deleting the directory makes sure the key-store is initialized correctly.
        Path securityBaseDir = Paths.get(System.getProperty("java.io.tmpdir"), "server", "security");
        try {
            Files.delete(securityBaseDir);
        } catch (Exception e) {
            // Ignore this ...
        }

        ExampleServer testServer = new ExampleServer();
        testServer.startup().get();
        try (PlcConnection opcuaConnection = new DefaultPlcDriverManager().getConnection("opcua:tcp://127.0.0.1:12686/milo?discovery=false")) {
            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder().
                addTagAddress("Bool", BOOL_IDENTIFIER).
                addTagAddress("ByteString", BYTE_STRING_IDENTIFIER).
                addTagAddress("Byte", BYTE_IDENTIFIER).
                addTagAddress("Double", DOUBLE_IDENTIFIER).
                addTagAddress("Float", FLOAT_IDENTIFIER).
                addTagAddress("Int16", INT16_IDENTIFIER).
                addTagAddress("Int32", INT32_IDENTIFIER).
                addTagAddress("Int64", INT64_IDENTIFIER).
                addTagAddress("Integer", INTEGER_IDENTIFIER).
                addTagAddress("SByte", SBYTE_IDENTIFIER).
                addTagAddress("String", STRING_IDENTIFIER).
                addTagAddress("UInt16", UINT16_IDENTIFIER).
                addTagAddress("UInt32", UINT32_IDENTIFIER).
                addTagAddress("UInt64", UINT64_IDENTIFIER).
                addTagAddress("UInteger", UINTEGER_IDENTIFIER).

                addTagAddress("BoolArray", BOOL_ARRAY_IDENTIFIER).
                addTagAddress("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER).
                addTagAddress("ByteArray", BYTE_ARRAY_IDENTIFIER).
                addTagAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER).
                addTagAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER).
                addTagAddress("Int16Array", INT16_ARRAY_IDENTIFIER).
                addTagAddress("Int32Array", INT32_ARRAY_IDENTIFIER).
                addTagAddress("Int64Array", INT64_ARRAY_IDENTIFIER).
                addTagAddress("IntegerArray", INTEGER_ARRAY_IDENTIFIER).
                addTagAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER).
                addTagAddress("StringArray", STRING_ARRAY_IDENTIFIER).
                addTagAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER).
                addTagAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER).
                addTagAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER).
                addTagAddress("UIntegerArray", UINTEGER_ARRAY_IDENTIFIER).

                addTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER);

            PlcReadRequest request = builder.build();


            PlcReadResponse response = request.execute().get();

            //Collection coll = response.getAllStrings("String");

            PlcWriteRequest.Builder wBuilder = opcuaConnection.writeRequestBuilder().
                addTagAddress("w-Bool", BOOL_IDENTIFIER, true).

                //addTagAddress("w-ByteString", BYTE_STRING_IDENTIFIER, "TEST".getBytes()).
                    addTagAddress("w-Byte", BYTE_IDENTIFIER, (byte) 1).

                addTagAddress("w-Double", DOUBLE_IDENTIFIER, (double) 0.25).
                addTagAddress("w-Float", FLOAT_IDENTIFIER, (float) 0.25).
                addTagAddress("w-INT16", INT16_IDENTIFIER, 12).
                addTagAddress("w-Int32", INT32_IDENTIFIER, (int) 314).
                addTagAddress("w-Int64", INT64_IDENTIFIER, (long) 123125).
                addTagAddress("w-Integer", INTEGER_IDENTIFIER, (int) 314).
                addTagAddress("w-SByte", SBYTE_IDENTIFIER, (byte) 1).
                addTagAddress("w-String", STRING_IDENTIFIER, "TEST").
                addTagAddress("w-UInt16", UINT16_IDENTIFIER, new BigInteger("12")).
                addTagAddress("w-UInt32", UINT32_IDENTIFIER, new BigInteger("123")).
                addTagAddress("w-UInt64", UINT64_IDENTIFIER, new BigInteger("1245152")).
                addTagAddress("w-UInteger", UINTEGER_IDENTIFIER, new BigInteger("1245152"));
            PlcWriteRequest writeRequest = wBuilder.build();
            PlcWriteResponse wResponse = writeRequest.execute().get();

            // Create Subscription
            PlcSubscriptionRequest.Builder sBuilder = opcuaConnection.subscriptionRequestBuilder().
                addChangeOfStateTagAddress("Bool", BOOL_IDENTIFIER).
                addChangeOfStateTagAddress("ByteString", BYTE_STRING_IDENTIFIER).
                addChangeOfStateTagAddress("Byte", BYTE_IDENTIFIER).
                addChangeOfStateTagAddress("Double", DOUBLE_IDENTIFIER).
                addChangeOfStateTagAddress("Float", FLOAT_IDENTIFIER).
                addChangeOfStateTagAddress("Int16", INT16_IDENTIFIER).
                addChangeOfStateTagAddress("Int32", INT32_IDENTIFIER).
                addChangeOfStateTagAddress("Int64", INT64_IDENTIFIER).
                addChangeOfStateTagAddress("Integer", INTEGER_IDENTIFIER).
                addChangeOfStateTagAddress("SByte", SBYTE_IDENTIFIER).
                addChangeOfStateTagAddress("String", STRING_IDENTIFIER).
                addChangeOfStateTagAddress("UInt16", UINT16_IDENTIFIER).
                addChangeOfStateTagAddress("UInt32", UINT32_IDENTIFIER).
                addChangeOfStateTagAddress("UInt64", UINT64_IDENTIFIER).
                addChangeOfStateTagAddress("UInteger", UINTEGER_IDENTIFIER).

                addChangeOfStateTagAddress("BoolArray", BOOL_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("ByteArray", BYTE_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("Int16Array", INT16_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("Int32Array", INT32_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("Int64Array", INT64_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("IntegerArray", INTEGER_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("StringArray", STRING_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER).
                addChangeOfStateTagAddress("UIntegerArray", UINTEGER_ARRAY_IDENTIFIER).

                addChangeOfStateTagAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER);
            PlcSubscriptionRequest subscriptionRequest = sBuilder.build();

            // Get result of creating subscription
            PlcSubscriptionResponse sResponse = subscriptionRequest.execute().get();
            final OpcuaSubscriptionHandle subscriptionHandle = (OpcuaSubscriptionHandle) sResponse.getSubscriptionHandle("Bool");

            // Create handler for returned value
            subscriptionHandle.register(plcSubscriptionEvent -> {
                assert plcSubscriptionEvent.getResponseCode("Bool").equals(PlcResponseCode.OK);
            });

            //Wait for value to be returned from server
            Thread.sleep(1200);

            subscriptionHandle.stopSubscriber();

            Thread.sleep(20000);
        }
    }
}
