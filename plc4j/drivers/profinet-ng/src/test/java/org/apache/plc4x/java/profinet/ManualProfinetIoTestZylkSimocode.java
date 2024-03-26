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

package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

public class ManualProfinetIoTestZylkSimocode {

    public static void main(String[] args) throws Exception {
        // WireShark filter: eth.addr == 88:3f:99:03:ef:b0
        try(PlcConnection connection =  new DefaultPlcDriverManager().getConnection("profinet:raw://192.168.54.23")) {
            // Create and execute the subscription request.
            PlcSubscriptionRequest subscriptionRequest = connection.subscriptionRequestBuilder()
                .addCyclicTagAddress("inputs", "1.1.INPUT.0:BYTE[10]", Duration.ofMillis(400))
                .addCyclicTagAddress("output", "1.1.OUTPUT.0:DWORD", Duration.ofMillis(400))
                .build();
            PlcSubscriptionResponse subscriptionResponse = subscriptionRequest.execute().get(100000, TimeUnit.MILLISECONDS);
            System.out.println(subscriptionResponse);
        }
    }

}
