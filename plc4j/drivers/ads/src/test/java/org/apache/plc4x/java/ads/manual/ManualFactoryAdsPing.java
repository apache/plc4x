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
package org.apache.plc4x.java.ads.manual;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.concurrent.TimeUnit;

public class ManualFactoryAdsPing {

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = PlcDriverManager.getDefault().getConnectionFactory().getConnection("ads:tcp://192.168.23.20:48898?target-ams-port=851&source-ams-port=65534&source-ams-net-id=192.168.23.220.1.1&target-ams-net-id=192.168.23.20.1.1")){
            for (int i = 0; i < 10; i++) {
                PlcPingResponse pingResponse = connection.ping().get(1000, TimeUnit.MILLISECONDS);
                if (pingResponse.getResponseCode() == PlcResponseCode.OK) {
                    System.out.println("Ping Success");
                } else {
                    System.out.println("Ping Failure");
                }
                Thread.sleep(500);
            }
        }
    }

}
