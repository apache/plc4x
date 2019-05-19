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
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.base.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.apache.plc4x.java.opcua.connection.OpcuaTcpPlcConnection;
import org.apache.plc4x.java.opcua.protocol.OpcuaField;
import org.apache.plc4x.java.opcua.protocol.OpcuaPlcFieldHandler;

import java.math.BigInteger;
import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Collection;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.function.Consumer;
/**
 * This class serves only as a manual entry point for ad-hoc tests of the OPC UA PLC4J driver.
 *
 *
 * The current version is tested against a public server, which is to be replaced later by a separate instance of the Milo framework.
 * Afterwards the code represented here will be used as an example for the introduction page.
 *
 * TODO: replace current public server with local Milo instance
 *
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
public class ManualPLC4XOpcua {
    private static final String BOOL_IDENTIFIER = "ns=2;i=10844";
    private static final String BYTE_STRING_IDENTIFIER = "ns=2;i=10858";
    private static final String BYTE_IDENTIFIER = "ns=2;i=10846";
    private static final String DOUBLE_IDENTIFIER = "ns=2;i=10854";
    private static final String FLOAT_IDENTIFIER = "ns=2;i=10853";
    private static final String INT16_IDENTIFIER = "ns=2;i=10847";
    private static final String INT32_IDENTIFIER = "ns=2;i=10849";
    private static final String INT64_IDENTIFIER = "ns=2;i=10851";
    private static final String INTEGER_IDENTIFIER = "ns=2;i=10869";
    private static final String SBYTE_IDENTIFIER = "ns=2;i=10845";
    private static final String STRING_IDENTIFIER = "ns=2;i=10855";
    private static final String UINT16_IDENTIFIER = "ns=2;i=10848";
    private static final String UINT32_IDENTIFIER = "ns=2;i=10850";
    private static final String UINT64_IDENTIFIER = "ns=2;i=10852";
    private static final String UINTEGER_IDENTIFIER = "ns=2;i=10870";
    private static final String DOES_NOT_EXIST_IDENTIFIER = "ns=2;i=12512623";


    public static void main(String args[]){


        OpcuaTcpPlcConnection opcuaConnection = null;
        OpcuaPlcFieldHandler fieldH = new OpcuaPlcFieldHandler();
        PlcField field = fieldH.createField("ns=2;i=10855");
        try {
            opcuaConnection = (OpcuaTcpPlcConnection)
                new PlcDriverManager().getConnection("opcua:tcp://opcua.demo-this.com:51210/UA/SampleServer");

        } catch (PlcConnectionException e) {
            e.printStackTrace();
        }
        try {
            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder();
            builder.addItem("Bool", BOOL_IDENTIFIER);
            builder.addItem("ByteString", BYTE_STRING_IDENTIFIER);
            builder.addItem("Byte", BYTE_IDENTIFIER);
            builder.addItem("Double", DOUBLE_IDENTIFIER);
            builder.addItem("Float", FLOAT_IDENTIFIER);
            builder.addItem("Int16", INT16_IDENTIFIER);
            builder.addItem("Int32", INT32_IDENTIFIER);
            builder.addItem("Int64", INT64_IDENTIFIER);
            builder.addItem("Integer", INTEGER_IDENTIFIER);
            builder.addItem("SByte", SBYTE_IDENTIFIER);
            builder.addItem("String", STRING_IDENTIFIER);
            builder.addItem("UInt16", UINT16_IDENTIFIER);
            builder.addItem("UInt32", UINT32_IDENTIFIER);
            builder.addItem("UInt64", UINT64_IDENTIFIER);
            builder.addItem("UInteger", UINTEGER_IDENTIFIER);

            builder.addItem("DoesNotExists", DOES_NOT_EXIST_IDENTIFIER);

            PlcReadRequest request = builder.build();
            PlcReadResponse response = opcuaConnection.read(request).get();
            Collection coll = response.getAllStrings("String");

            PlcWriteRequest.Builder wBuilder = opcuaConnection.writeRequestBuilder();
            wBuilder.addItem("w-Bool", BOOL_IDENTIFIER, true);
            wBuilder.addItem("w-ByteString", BYTE_STRING_IDENTIFIER, "TEST".getBytes());
            wBuilder.addItem("w-Byte", BYTE_IDENTIFIER, (byte)1);
            wBuilder.addItem("w-Double", DOUBLE_IDENTIFIER, (double)0.25);
            wBuilder.addItem("w-Float", FLOAT_IDENTIFIER, (float)0.25);
            wBuilder.addItem("w-INT16", INT16_IDENTIFIER, (short)12);
            wBuilder.addItem("w-Int32", INT32_IDENTIFIER, (int)314);
            wBuilder.addItem("w-Int64", INT64_IDENTIFIER, (long)123125);
            wBuilder.addItem("w-Integer", INTEGER_IDENTIFIER, (int)314);
            wBuilder.addItem("w-SByte", SBYTE_IDENTIFIER, (short)23);
            wBuilder.addItem("w-String", STRING_IDENTIFIER, "TEST");
            wBuilder.addItem("w-UInt16", UINT16_IDENTIFIER, (int)222);
            wBuilder.addItem("w-UInt32", UINT32_IDENTIFIER, (long)21412);
            wBuilder.addItem("w-UInt64", UINT64_IDENTIFIER, new BigInteger("1245152"));
            wBuilder.addItem("w-UInteger", UINTEGER_IDENTIFIER, new BigInteger("1245152"));
            PlcWriteRequest writeRequest = wBuilder.build();
            PlcWriteResponse wResponse = opcuaConnection.write(writeRequest).get();

            PlcSubscriptionResponse subResp = opcuaConnection.subscribe(new DefaultPlcSubscriptionRequest(
                opcuaConnection,
                new LinkedHashMap<>(
                    Collections.singletonMap("field1",
                        new SubscriptionPlcField(PlcSubscriptionType.CHANGE_OF_STATE, OpcuaField.of(STRING_IDENTIFIER), Duration.of(1, ChronoUnit.SECONDS)))
                )
            )).get();

            Consumer<PlcSubscriptionEvent> consumer = plcSubscriptionEvent -> System.out.println(plcSubscriptionEvent.toString());
            PlcConsumerRegistration registration = opcuaConnection.register(consumer, subResp.getSubscriptionHandles());
            Thread.sleep(7000);
            registration.unregister();
            Thread.sleep(20000);
            opcuaConnection.close();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
