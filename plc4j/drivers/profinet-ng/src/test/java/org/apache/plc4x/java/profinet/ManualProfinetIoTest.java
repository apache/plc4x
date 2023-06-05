/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcBrowseResponse;

import java.util.concurrent.TimeUnit;

public class ManualProfinetIoTest {

    public static void main(String[] args) throws Exception {
        try(PlcConnection connection =  new DefaultPlcDriverManager().getConnection("profinet:raw://192.168.24.31")) {
            PlcBrowseRequest browseRequest = connection.browseRequestBuilder().addQuery("all", "*").build();
            Thread.sleep(10000);
            /*PlcBrowseResponse plcBrowseResponse = browseRequest.execute().get(30000, TimeUnit.MILLISECONDS);
            System.out.println(plcBrowseResponse);*/
        }
    }

}
