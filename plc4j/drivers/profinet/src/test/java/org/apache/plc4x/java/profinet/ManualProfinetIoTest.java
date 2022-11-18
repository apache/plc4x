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

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;

public class ManualProfinetIoTest {

    public static void main(String[] args) throws Exception {
        final PlcConnection connection = new PlcDriverManager().getConnection("profinet://192.168.90.1?gsddirectory=/home/missy/Documents/Profinet/gsd&devices=[00:0c:29:75:25:67]&submodules=[[IDM_30,IDM_32,IDM_31,]]");
        final PlcSubscriptionRequest request = connection.subscriptionRequestBuilder().addChangeOfStateField("Default Float", "I have no idea").build();
        final PlcSubscriptionResponse plcResponse = request.execute().get();
        System.out.println(plcResponse);
    }

}
