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

import java.util.List;

public class ManualProfinetIoTestWagoPN {

    public static void main(String[] args) throws Exception {
        // WireShark filter: "eth.addr == 00:30:de:61:37:79"
        try(PlcConnection connection =  new DefaultPlcDriverManager().getConnection("profinet:raw://00:30:de:61:37:79?ip-address=192.168.24.51")) {
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
