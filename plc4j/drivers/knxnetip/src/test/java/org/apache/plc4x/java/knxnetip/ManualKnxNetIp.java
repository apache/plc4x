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
package org.apache.plc4x.java.knxnetip;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;

import java.util.concurrent.TimeUnit;

public class ManualKnxNetIp {

    // Addresses:
    // */*/10: Temperature
    // */*/12: Heating
    // */*/60: Primary Window
    // */*/64: Second Window
    // */*/101: Power Line 1

    public static void main(String[] args) throws Exception {
        //final PlcConnection connection = new PlcDriverManager().getConnection("knxnet-ip://192.168.42.11?knxproj-file-path=/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/KNX/Stettiner%20Str.%2013/StettinerStr-Soll-Ist-Temperatur.knxproj");
        final PlcConnection connection = new PlcDriverManager().getConnection("knxnet-ip:pcap:///Users/christofer.dutz/Projects/Apache/PLC4X-Documents/KNX/Recording-01.03.2020-2.pcapng?knxproj-file-path=/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/KNX/Stettiner%20Str.%2013/StettinerStr-Soll-Ist-Temperatur.knxproj");
        // Make sure we hang up correctly when terminating.
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            try {
                connection.close();
            } catch (Exception  e) {
                throw new PlcRuntimeException("Error closing connection", e);
            }
        }));

        // Create a new subscription request.
        // The address and the name is just bogus as we're always returning everything.
        // We will probably refactor the API in the near future.
        final PlcSubscriptionRequest subscriptionRequest = connection.subscriptionRequestBuilder()
            .addEventFieldAddress("knxData", "*/*/*")
            .build();

        // Register the subscription
        // The timeout is also just a bogus value as the data is coming in actively
        // We will probably refactor the API in the near future.
        final PlcSubscriptionResponse subscriptionResponse =
            subscriptionRequest.execute().get(1000, TimeUnit.MILLISECONDS);

        // Register a callback which is called on new data being available.
        final PlcSubscriptionHandle subscriptionHandle = subscriptionResponse.getSubscriptionHandle("knxData");
        subscriptionHandle.register(knxData -> {
            System.out.println(knxData.getTimestamp().toString() + " - " +
                ((DefaultPlcSubscriptionEvent) knxData).getValues().get("knxData"));
        });
    }

}
