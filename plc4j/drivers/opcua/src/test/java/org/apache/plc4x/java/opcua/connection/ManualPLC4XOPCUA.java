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
 * @author Matthias Milan Stlrljic
 * Created by Matthias Milan Stlrljic on 10.05.2019
 */
package org.apache.plc4x.java.opcua.connection;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.base.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.apache.plc4x.java.opcua.protocol.OpcuaField;
import org.apache.plc4x.java.opcua.protocol.model.OpcuaPlcFieldHandler;

import java.math.BigInteger;
import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Collection;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.function.Consumer;

public class ManualPLC4XOPCUA {
    public static void main(String args[]){




        OPCUATcpPlcConnection opcuaConnection = null;
        OpcuaPlcFieldHandler fieldH = new OpcuaPlcFieldHandler();
        PlcField field = fieldH.createField("ns=2;i=10855");
        try {
            opcuaConnection = (OPCUATcpPlcConnection)
                new PlcDriverManager().getConnection("opcua:tcp://opcua.demo-this.com:51210/UA/SampleServer");

        } catch (PlcConnectionException e) {
            e.printStackTrace();
        }
        try {
            PlcReadRequest.Builder builder = opcuaConnection.readRequestBuilder();
            //builder.addItem("String", "ns=2;i=10855");
            builder.addItem("Bool", "ns=2;i=10844");
            builder.addItem("ByteString", "ns=2;i=10858");
            builder.addItem("Byte", "ns=2;i=10846");
            builder.addItem("Double", "ns=2;i=10854");
            builder.addItem("Float", "ns=2;i=10853");
            builder.addItem("Int16", "ns=2;i=10847");
            builder.addItem("Int32", "ns=2;i=10849");
            builder.addItem("Int64", "ns=2;i=10851");
            builder.addItem("Integer", "ns=2;i=10869");
            builder.addItem("SByte", "ns=2;i=10845");
            builder.addItem("String", "ns=2;i=10855");
            builder.addItem("UInt16", "ns=2;i=10848");
            builder.addItem("UInt32", "ns=2;i=10850");
            builder.addItem("UInt64", "ns=2;i=10852");
            builder.addItem("UInteger", "ns=2;i=10870");

            builder.addItem("DoesNotExists", "ns=2;i=12512623");

            PlcReadRequest request = builder.build();
            PlcReadResponse response = opcuaConnection.read(request).get();
            Collection coll = response.getAllStrings("String");

            PlcWriteRequest.Builder wBuilder = opcuaConnection.writeRequestBuilder();
            wBuilder.addItem("w-Bool", "ns=2;i=11012", true);
            /*
            wBuilder.addItem("w-ByteString", "ns=2;i=10858", "TEST".getBytes());
            wBuilder.addItem("w-Byte", "ns=2;i=10846", (byte)1);
            wBuilder.addItem("w-Double", "ns=2;i=10854", (double)0.25);
            wBuilder.addItem("w-Float", "ns=2;i=10853", (float)0.25);
            wBuilder.addItem("w-Int16", "ns=2;i=10847", (short)12);
            wBuilder.addItem("w-Int32", "ns=2;i=10849", (int)314);
            wBuilder.addItem("w-Int64", "ns=2;i=10851", (long)123125);
            wBuilder.addItem("w-Integer", "ns=2;i=10869", (int)314);
            wBuilder.addItem("w-SByte", "ns=2;i=10845", (short)23);
            wBuilder.addItem("w-String", "ns=2;i=10855", "TEST");
            wBuilder.addItem("w-UInt16", "ns=2;i=10848", (int)222);
            wBuilder.addItem("w-UInt32", "ns=2;i=10850", (long)21412);
            wBuilder.addItem("w-UInt64", "ns=2;i=10852", new BigInteger("1245152"));
            wBuilder.addItem("w-UInteger", "ns=2;i=10870", new BigInteger("1245152"));
*/
            PlcWriteRequest writeRequest = wBuilder.build();
            PlcWriteResponse wResponse = opcuaConnection.write(writeRequest).get();

            PlcSubscriptionResponse subResp = opcuaConnection.subscribe(new DefaultPlcSubscriptionRequest(
                opcuaConnection,
                new LinkedHashMap<>(
                    Collections.singletonMap("field1",
                        new SubscriptionPlcField(PlcSubscriptionType.CHANGE_OF_STATE, OpcuaField.of("ns=2;i=10855"), Duration.of(1, ChronoUnit.SECONDS)))
                )
            )).get();

            Consumer<PlcSubscriptionEvent> consumer = plcSubscriptionEvent -> System.out.println(plcSubscriptionEvent.toString());
            PlcConsumerRegistration registration = opcuaConnection.register(consumer, subResp.getSubscriptionHandles());
            Thread.sleep(7000);
            registration.unregister();
            Thread.sleep(200000);
            opcuaConnection.close();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
