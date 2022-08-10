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
package org.apache.plc4x.java.examples.helloplc4x.discoverandbrowse;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

public class HelloPlc4x {

    private static final Logger logger = LoggerFactory.getLogger(HelloPlc4x.class);

    public static void main(String[] args) throws Exception {
        // Iterate over all installed drivers and execute their browse functionality (If they support it)
        PlcDriverManager plcDriverManager = new PlcDriverManager();
        for (String protocolCode : plcDriverManager.listDrivers()) {
            // For some reason modbus is failing on my machine ... investigate
            if(protocolCode.startsWith("modbus")) {
                continue;
            }
            PlcDriver driver = plcDriverManager.getDriver(protocolCode);
            if(driver.getMetadata().canDiscover()) {
                logger.info("Performing discovery for {} protocol", driver.getProtocolName());

                PlcDiscoveryRequest discoveryRequest = driver.discoveryRequestBuilder().build();

                PlcDiscoveryResponse discoveryResponse = discoveryRequest.execute().get(60, TimeUnit.SECONDS);
                if(discoveryResponse.getResponseCode() == PlcResponseCode.OK) {
                    for (PlcDiscoveryItem value : discoveryResponse.getValues()) {
                        logger.info(" - Found device with connection-url {}", value.getConnectionUrl());
                    }
                }
            }
        }
    }

}
