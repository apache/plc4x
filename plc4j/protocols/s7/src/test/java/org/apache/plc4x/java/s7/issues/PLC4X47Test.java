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

package org.apache.plc4x.java.s7.issues;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.s7.connection.S7PlcConnection;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;


public class PLC4X47Test {

    @Test
    @Disabled
    public void testLargeRequest() throws Exception {
        /*TestChannelFactory channelFactory = new TestChannelFactory();
        S7PlcConnection connection = new S7PlcConnection( channelFactory, 1, 1, "");
        connection.connect();
        EmbeddedChannel channel = channelFactory.getChannel();*/
        S7PlcConnection connection = (S7PlcConnection) new PlcDriverManager().getConnection("s7://10.10.64.20/1/1");

        PlcReadRequest.Builder builder = connection.readRequestBuilder().get();
        for (int i = 1; i <= 30; i++) {
            // just the first byte of each db
            builder.addItem("field-" + i, "%DB3.DB" + i + ":SINT");
        }
        PlcReadRequest readRequest = builder.build();
        PlcReadResponse readResponse = connection.read(readRequest).get();
        System.out.println(readResponse.getFieldNames().size());
    }

}
