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
package org.apache.plc4x.java.knxnetip.maual;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.value.PlcValue;

/**
 * Manual test for browsing KNX group addresses from an ETS project file.
 * Browse is only available when a knxproj file is configured.
 */
public class ManualKnxNetIpRead {

    public static void main(String[] args) throws Exception {
        try (PlcConnection connection = new DefaultPlcDriverManager().getConnection(
            "knxnet-ip://192.168.42.28?" +
            "knxproj-file-path=huiiiii&" +
            "knxproj-password=lalala")) {

            PlcReadRequest readRequest = connection.readRequestBuilder()
                .addTagAddress("Lade-Leistung Batterie", "1/1/210:DPT14")  // Temperature (2-byte float)
                .addTagAddress("Leistung Hausverbrauch", "1/1/211:DPT14") // Switch (boolean)
                .build();

            PlcReadResponse readResponse = readRequest.execute().get();

            PlcValue chargePower = readResponse.getPlcValue("Lade-Leistung Batterie");
            System.out.println("Charge power: " + chargePower + " W");

            PlcValue housePower = readResponse.getPlcValue("Leistung Hausverbrauch");
            System.out.println("House power: " + housePower + " W");
        }
    }

}
