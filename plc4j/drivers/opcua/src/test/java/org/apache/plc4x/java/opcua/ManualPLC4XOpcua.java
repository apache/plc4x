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

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.opcua.field.OpcuaPlcFieldHandler;
import org.apache.plc4x.java.opcua.protocol.OpcuaSubscriptionHandle;
import org.eclipse.milo.examples.server.ExampleServer;

import java.math.BigInteger;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

/**
 * This class serves only as a manual entry point for ad-hoc tests of the OPC UA PLC4J driver.
 * <p>
 * <p>
 * The current version is tested against a public server, which is to be replaced later by a separate instance of the Milo framework.
 * Afterwards the code represented here will be used as an example for the introduction page.
 * <p>
 *
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

    public static void main(String args[]) {
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

            ExampleServer testServer = new ExampleServer();
            testServer.startup().get();
        } catch (Exception e) {
            throw new PlcRuntimeException(e);
        }
        PlcConnection opcuaConnection = null;
        OpcuaPlcFieldHandler fieldH = new OpcuaPlcFieldHandler();
        PlcField field = fieldH.parseField(BOOL_IDENTIFIER);
        try {
            opcuaConnection = new PlcDriverManager().getConnection("opcua:tcp://127.0.0.1:12686/milo?discovery=false");

        } catch (PlcConnectionException e) {
            throw new PlcRuntimeException(e);
        }
        try {
            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder();
            builder.addFieldAddress("Bool", BOOL_IDENTIFIER);
            builder.addFieldAddress("ByteString", BYTE_STRING_IDENTIFIER);
            builder.addFieldAddress("Byte", BYTE_IDENTIFIER);
            builder.addFieldAddress("Double", DOUBLE_IDENTIFIER);
            builder.addFieldAddress("Float", FLOAT_IDENTIFIER);
            builder.addFieldAddress("Int16", INT16_IDENTIFIER);
            builder.addFieldAddress("Int32", INT32_IDENTIFIER);
            builder.addFieldAddress("Int64", INT64_IDENTIFIER);
            builder.addFieldAddress("Integer", INTEGER_IDENTIFIER);
            builder.addFieldAddress("SByte", SBYTE_IDENTIFIER);
            builder.addFieldAddress("String", STRING_IDENTIFIER);
            builder.addFieldAddress("UInt16", UINT16_IDENTIFIER);
            builder.addFieldAddress("UInt32", UINT32_IDENTIFIER);
            builder.addFieldAddress("UInt64", UINT64_IDENTIFIER);
            builder.addFieldAddress("UInteger", UINTEGER_IDENTIFIER);

            builder.addFieldAddress("BoolArray", BOOL_ARRAY_IDENTIFIER);
            builder.addFieldAddress("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER);
            builder.addFieldAddress("ByteArray", BYTE_ARRAY_IDENTIFIER);
            builder.addFieldAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER);
            builder.addFieldAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER);
            builder.addFieldAddress("Int16Array", INT16_ARRAY_IDENTIFIER);
            builder.addFieldAddress("Int32Array", INT32_ARRAY_IDENTIFIER);
            builder.addFieldAddress("Int64Array", INT64_ARRAY_IDENTIFIER);
            builder.addFieldAddress("IntegerArray", INTEGER_ARRAY_IDENTIFIER);
            builder.addFieldAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER);
            builder.addFieldAddress("StringArray", STRING_ARRAY_IDENTIFIER);
            builder.addFieldAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER);
            builder.addFieldAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER);
            builder.addFieldAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER);
            builder.addFieldAddress("UIntegerArray", UINTEGER_ARRAY_IDENTIFIER);

            builder.addFieldAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER);

            PlcReadRequest request = builder.build();


            PlcReadResponse response = request.execute().get();

            //Collection coll = response.getAllStrings("String");

            PlcWriteRequest.Builder wBuilder = opcuaConnection.writeRequestBuilder();
            wBuilder.addFieldAddress("w-Bool", BOOL_IDENTIFIER, true);
            //wBuilder.addField("w-ByteString", BYTE_STRING_IDENTIFIER, "TEST".getBytes());
            wBuilder.addFieldAddress("w-Byte", BYTE_IDENTIFIER, (byte) 1);
            wBuilder.addFieldAddress("w-Double", DOUBLE_IDENTIFIER, (double) 0.25);
            wBuilder.addFieldAddress("w-Float", FLOAT_IDENTIFIER, (float) 0.25);
            wBuilder.addFieldAddress("w-INT16", INT16_IDENTIFIER,  12);
            wBuilder.addFieldAddress("w-Int32", INT32_IDENTIFIER, (int) 314);
            wBuilder.addFieldAddress("w-Int64", INT64_IDENTIFIER, (long) 123125);
            wBuilder.addFieldAddress("w-Integer", INTEGER_IDENTIFIER, (int) 314);
            wBuilder.addFieldAddress("w-SByte", SBYTE_IDENTIFIER, (byte) 1);
            wBuilder.addFieldAddress("w-String", STRING_IDENTIFIER, "TEST");
            wBuilder.addFieldAddress("w-UInt16", UINT16_IDENTIFIER, new BigInteger("12"));
            wBuilder.addFieldAddress("w-UInt32", UINT32_IDENTIFIER, new BigInteger("123"));
            wBuilder.addFieldAddress("w-UInt64", UINT64_IDENTIFIER, new BigInteger("1245152"));
            wBuilder.addFieldAddress("w-UInteger", UINTEGER_IDENTIFIER, new BigInteger("1245152"));
            PlcWriteRequest writeRequest = wBuilder.build();
            PlcWriteResponse wResponse = writeRequest.execute().get();

            // Create Subscription
            PlcSubscriptionRequest.Builder sBuilder = opcuaConnection.subscriptionRequestBuilder();
            sBuilder.addChangeOfStateFieldAddress("Bool", BOOL_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("ByteString", BYTE_STRING_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Byte", BYTE_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Double", DOUBLE_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Float", FLOAT_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Int16", INT16_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Int32", INT32_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Int64", INT64_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Integer", INTEGER_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("SByte", SBYTE_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("String", STRING_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInt16", UINT16_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInt32", UINT32_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInt64", UINT64_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInteger", UINTEGER_IDENTIFIER);

            sBuilder.addChangeOfStateFieldAddress("BoolArray", BOOL_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("ByteStringArray", BYTE_STRING_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("ByteArray", BYTE_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("DoubleArray", DOUBLE_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("FloatArray", FLOAT_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Int16Array", INT16_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Int32Array", INT32_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("Int64Array", INT64_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("IntegerArray", INTEGER_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("SByteArray", SBYTE_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("StringArray", STRING_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInt16Array", UINT16_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInt32Array", UINT32_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UInt64Array", UINT64_ARRAY_IDENTIFIER);
            sBuilder.addChangeOfStateFieldAddress("UIntegerArray", UINTEGER_ARRAY_IDENTIFIER);

            sBuilder.addChangeOfStateFieldAddress("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER);
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
            opcuaConnection.close();

        } catch (Exception e) {
            throw new PlcRuntimeException(e);
        }

    }
}
