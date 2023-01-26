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
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadResponse;

public class ManualTest {

    public static void main(String[] args) {
        CachedPlcConnectionManager cachedPlcConnectionManager = CachedPlcConnectionManager.getBuilder(new DefaultPlcDriverManager()).build();
        for (int i = 0; i < 100; i++){
            try (PlcConnection connection = cachedPlcConnectionManager.getConnection("ads:tcp://192.168.23.20?sourceAmsNetId=192.168.23.200.1.1&sourceAmsPort=65534&targetAmsNetId=192.168.23.20.1.1&targetAmsPort=851")) {
                PlcReadResponse plcReadResponse = connection.readRequestBuilder().addTagAddress("var", "MAIN.hurz_REAL").build().execute().get();
                System.out.printf("Run %d: Value: %f%n", i, plcReadResponse.getFloat("var"));
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        }
    }

}
