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
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;

import java.util.Collection;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

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
        // TODO: temporary workaround
        Thread.currentThread().setUncaughtExceptionHandler((t, e) -> {
            System.err.println(t + " - " + e.getMessage());
            e.printStackTrace(System.err);
            System.exit(1);
        });
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection(connectionUrl)) {
            System.out.println("PlcConnection " + plcConnection);

            PlcReadRequest.Builder readRequestBuilder = plcConnection.readRequestBuilder().orElseThrow(RuntimeException::new);
            PlcReadRequest readRequest = readRequestBuilder.addItem("station", "Allgemein_S2.Station:BYTE").build();
            CompletableFuture<? extends PlcReadResponse> response = readRequest.execute();
            PlcReadResponse readResponse = response.get();
            System.out.println("Response " + readResponse);
            Collection<Integer> stations = readResponse.getAllIntegers("station");
            stations.forEach(System.out::println);

            // 2. We build a subscription
            PlcSubscriptionRequest.Builder subscriptionRequestBuilder = plcConnection.subscriptionRequestBuilder().orElseThrow(RuntimeException::new);
            PlcSubscriptionRequest subscriptionRequest = subscriptionRequestBuilder.addChangeOfStateField("stationChange", "Allgemein_S2.Station:BYTE").build();
            PlcSubscriptionResponse plcSubscriptionResponse = subscriptionRequest.execute().get();

            List<PlcConsumerRegistration> plcConsumerRegistrations = plcSubscriptionResponse.getSubscriptionHandles().stream()
                .map(plcSubscriptionHandle -> plcSubscriptionHandle.register(System.out::println))
                .collect(Collectors.toList());

            // Now we wait a bit
            TimeUnit.SECONDS.sleep(5);

            // we unregister the listener
            plcConsumerRegistrations.forEach(PlcConsumerRegistration::unregister);

            // we unsubscribe at the plc
            PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder = plcConnection.unsubscriptionRequestBuilder().orElseThrow(RuntimeException::new);
            PlcUnsubscriptionRequest unsubscriptionRequest = unsubscriptionRequestBuilder.addHandles(plcSubscriptionResponse.getSubscriptionHandles()).build();
            CompletableFuture<PlcUnsubscriptionResponse> unsubscriptionResponse = unsubscriptionRequest.execute();

            unsubscriptionResponse
                .get(5, TimeUnit.SECONDS);
            System.out.println(unsubscriptionResponse);
        }
        System.exit(0);
    }
}
