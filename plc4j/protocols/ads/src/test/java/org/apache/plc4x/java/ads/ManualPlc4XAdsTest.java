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
package org.apache.plc4x.java.ads;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.connection.PlcConnection;
import org.apache.plc4x.java.api.connection.PlcReader;
import org.apache.plc4x.java.api.connection.PlcSubscriber;
import org.apache.plc4x.java.api.messages.PlcNotification;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.model.Address;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

public class ManualPlc4XAdsTest {

    public static void main(String... args) throws Exception {
        String connectionUrl;
        if (args.length > 0 && "serial".equalsIgnoreCase(args[0])) {
            System.out.println("Using serial");
            connectionUrl = "ads:serial:///dev/ttys003/10.10.64.40.1.1:851/10.10.56.23.1.1:30000";
        } else {
            System.out.println("Using tcp");
            connectionUrl = "ads:tcp://10.10.64.40/10.10.64.40.1.1:851/10.10.56.23.1.1:30000";
        }
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {
            System.out.println("PlcConnection " + plcConnection);

            PlcReader reader = plcConnection.getReader().orElseThrow(() -> new RuntimeException("No Reader found"));

            Address address = plcConnection.parseAddress("Allgemein_S2.Station");
            CompletableFuture<TypeSafePlcReadResponse<Integer>> response = reader
                .read(new TypeSafePlcReadRequest<>(Integer.class, address));
            TypeSafePlcReadResponse<Integer> readResponse = response.get();
            System.out.println("Response " + readResponse);
            ReadResponseItem<Integer> responseItem = readResponse.getResponseItem().orElseThrow(() -> new RuntimeException("No Item found"));
            System.out.println("ResponseItem " + responseItem);
            responseItem.getValues().stream().map(integer -> "Value: " + integer).forEach(System.out::println);

            Consumer<PlcNotification<Integer>> notificationConsumer = plcNotification -> System.out.println("Received notification " + plcNotification);
            PlcSubscriber plcSubscriber = plcConnection.getSubscriber().orElseThrow(() -> new RuntimeException("Subscribe not available"));
            plcSubscriber.subscribe(notificationConsumer, address, Integer.class);
            TimeUnit.SECONDS.sleep(5);
            plcSubscriber.unsubscribe(notificationConsumer, address);
        }
        System.exit(0);
    }
}
