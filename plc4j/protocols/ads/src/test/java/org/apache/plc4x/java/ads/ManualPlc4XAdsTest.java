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
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.*;

import java.util.Collection;
import java.util.concurrent.CompletableFuture;

public class ManualPlc4XAdsTest {

    public static void main(String... args) throws Exception {
        String connectionUrl;
        if (args.length > 0 && "serial".equalsIgnoreCase(args[0])) {
            System.out.println("Using serial");
            connectionUrl = "ads:serial:///dev/ttys003/10.10.64.40.1.1:851/10.10.56.23.1.1:30000";
        } else {
            System.out.println("Using tcp");
            connectionUrl = "ads:tcp://10.10.64.40/10.10.64.40.1.1:851/192.168.113.1.1.1:30000";
        }
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {
            System.out.println("PlcConnection " + plcConnection);

            PlcReadRequest readRequest = plcConnection.readRequestBuilder().get().addItem("station", "Allgemein_S2.Station:BYTE").build();
            CompletableFuture<? extends PlcReadResponse> response = readRequest.execute();
            PlcReadResponse readResponse = response.get();
            System.out.println("Response " + readResponse);
            Collection<Integer> stations = readResponse.getAllIntegers("station");
            stations.forEach(System.out::println);

            PlcSubscriptionRequest subscriptionRequest = plcConnection.subscriptionRequestBuilder().get().addChangeOfStateField("stationChange", "Allgemein_S2.Station:BYTE").build();
            CompletableFuture<? extends PlcSubscriptionResponse> subscribeResponse = subscriptionRequest.execute();
            PlcSubscriptionResponse plcSubscriptionResponse = subscribeResponse.get();

            // TODO: figure out what to do with this
            /*PlcConsumerRegistration plcConsumerRegistration = plcSubscriber.register(System.out::println, plcSubscriptionResponse.getSubscriptionHandles());

            TimeUnit.SECONDS.sleep(5);

            plcSubscriber.unregister(plcConsumerRegistration);
            PlcUnsubscriptionRequest unsubscriptionRequest = plcConnection.unsubscriptionRequestBuilder().get().addHandles(plcSubscriptionResponse.getSubscriptionHandles()).build();
            CompletableFuture<PlcUnsubscriptionResponse> unsubscriptionResponse = plcSubscriber.unsubscribe(unsubscriptionRequest);

            unsubscriptionResponse
                .get(5, TimeUnit.SECONDS);
            System.out.println(unsubscriptionResponse);*/
        }
        System.exit(0);
    }
}
