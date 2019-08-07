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

package org.apache.plc4x.protocol.df1;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

import java.util.concurrent.TimeUnit;

/**
 * TODO write comment
 *
 * @author julian
 * Created by julian on 2019-08-07
 */
public class EndToEndTest {

    @org.junit.Test
    public void helloDf1() {
        try (PlcConnection plcConnection = new PlcDriverManager().getConnection("df1:serial:///COM4")) {
            PlcReadRequest request = plcConnection.readRequestBuilder()
                .addItem("erstes", "17:INTEGER")
                .build();

            PlcReadResponse response = request.execute().get(1, TimeUnit.SECONDS);

            System.out.println(request);
        } catch (PlcConnectionException e) {
            e.printStackTrace();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
