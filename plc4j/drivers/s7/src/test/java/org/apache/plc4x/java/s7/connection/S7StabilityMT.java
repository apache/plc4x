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

package org.apache.plc4x.java.s7.connection;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.junit.Test;

/**
 * This is a manual test to prove the stability of the Driver.
 *
 * 13:37:36.862 [nioEventLoopGroup-2-1] DEBUG o.a.plc4x.java.s7.netty.S7Protocol - S7 Message with id 1018 queued
 * 13:37:36.862 [nioEventLoopGroup-2-1] TRACE o.a.p.j.isotp.protocol.IsoTPProtocol - ISO Transport Protocol Message sent
 * 13:37:36.862 [nioEventLoopGroup-2-1] DEBUG o.a.p.j.i.protocol.IsoOnTcpProtocol - ISO on TCP Message sent
 * 13:37:36.862 [nioEventLoopGroup-2-1] DEBUG o.a.plc4x.java.s7.netty.S7Protocol - S7 Message with id 1018 sent
 *
 * @author julian
 * Created by julian on 2019-08-07
 */
public class S7StabilityMT {

    @Test
    public void testFailingConnection() throws InterruptedException {
        // Establish a connection to the plc using the url provided as first argument
        for (int j = 1; j <= 2; j++) {
            try (PlcConnection plcConnection = new PlcDriverManager().getConnection("s7://192.168.167.210/0/0")) {

                for (int i = 0; i <= 10000; i++) {
                    try {
                        PlcReadResponse response = plcConnection.readRequestBuilder()
                            .addItem("field", "%DB400:DBW10:INT")
                            .build()
                            .execute()
                            .get();

                        System.out.println(response.getResponseCode("field"));
                    } catch (Exception e) {
                        // do nothjing
                    }
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
            Thread.sleep(100);
        }

        System.out.println("The loop is finished");
    }
}
