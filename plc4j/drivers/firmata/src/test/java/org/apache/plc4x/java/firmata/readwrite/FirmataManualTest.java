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

package org.apache.plc4x.java.firmata.readwrite;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;

public class FirmataManualTest {

    public static void main(String[] args) throws Exception {

        // Default for the Firmata Arduino sketch is, 57600 baud, 8 data-bits, no parity and one stop bit (Which is the default for the Serial transport)
        try(PlcConnection connection = PlcDriverManager.getDefault().getConnectionManager().getConnection("firmata:///dev/tty.usbmodem1101")) {
            PlcSubscriptionRequest subscriptionRequest = connection.subscriptionRequestBuilder().addEventTagAddress("analog1", "analog:1").setConsumer(plcSubscriptionEvent -> {
                System.out.println("Incoming Event: " + plcSubscriptionEvent.getPlcValue("analog1").getInteger());
            }).build();
            PlcSubscriptionResponse plcSubscriptionResponse = subscriptionRequest.execute().get();
            System.out.println(plcSubscriptionResponse);
            Thread.sleep(10000);
        }
    }
}
