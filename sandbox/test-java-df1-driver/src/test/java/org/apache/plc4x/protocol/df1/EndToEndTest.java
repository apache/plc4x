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
package org.apache.plc4x.protocol.df1;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

import java.util.concurrent.TimeUnit;

public class EndToEndTest {

    @Test
    @Disabled("Seems to cause problems on Windows if no COM4 is available")
    public void helloDf1() {
        try (PlcConnection plcConnection = new DefaultPlcDriverManager().getConnection("df1:serial:///COM4")) {
            PlcReadRequest request = plcConnection.readRequestBuilder()
                .addTagAddress("hurz", "5:INTEGER")
                .build();

            PlcReadResponse response = request.execute().get(100, TimeUnit.SECONDS);


            // TODO: get the actual read bytes from the response
            System.out.println(response);
            System.out.println("Response code was " + response.getResponseCode("erster"));
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
