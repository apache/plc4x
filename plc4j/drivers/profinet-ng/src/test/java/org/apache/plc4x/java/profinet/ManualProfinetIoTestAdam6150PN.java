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
import org.apache.plc4x.java.api.messages.*;

import java.time.Duration;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class ManualProfinetIoTestAdam6150PN {

    public static void main(String[] args) throws Exception {
        // Advantech Adam 6150PN (8 channel digital input, 7 chanel digital output)
        // WireShark filter: "eth.addr == 74:fe:48:82:4a:7c"
        try(PlcConnection connection =  new DefaultPlcDriverManager().getConnection("profinet:raw://192.168.24.42")) {
            // List up all resources found in the remote device.
            PlcBrowseRequest browseRequest = connection.browseRequestBuilder().addQuery("all", "*").build();
            PlcBrowseResponse plcBrowseResponse = browseRequest.execute().get();
            for (String queryName : plcBrowseResponse.getQueryNames()) {
                List<PlcBrowseItem> values = plcBrowseResponse.getValues(queryName);
                for (PlcBrowseItem value : values) {
                    System.out.println(value.getName() + ": " + value.getTag().getAddressString());
                }
            }
        }
    }

}
