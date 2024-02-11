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
package org.apache.plc4x.java.examples.hellodiscovery;

import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Set;

public class HelloDiscovery {

    private static final Logger logger = LoggerFactory.getLogger(HelloDiscovery.class);

    /**
     * Example code do demonstrate using PLC4X's discovery API.
     *
     * @param args ignored.
     */
    public static void main(String[] args) throws Exception {
        CliOptions options = CliOptions.fromArgs(args);
        if (options == null) {
            CliOptions.printHelp();
            // Could not parse.
            System.exit(1);
        }

        PlcDriverManager plcDriverManager = PlcDriverManager.getDefault();
        Set<String> driverCodes = plcDriverManager.listProtocolCodes();
        for (String driverCode : driverCodes) {
            PlcDriver driver = plcDriverManager.getDriver(driverCode);

            // Check if this driver supports discovery.
            if(driver.getMetadata().canDiscover()) {
                logger.info("Executing Discovery for Driver: {}", driverCode);
                PlcDiscoveryRequest discoveryRequest = driver.discoveryRequestBuilder().build();
                PlcDiscoveryResponse discoveryResponse = discoveryRequest.executeWithHandler(
                    discoveryItem -> logger.info("Intercepted discovery of device with name: {} with connection url: {}",
                        discoveryItem.getName(), discoveryItem.getConnectionUrl())).get();
                /*if(discoveryResponse.getResponseCode() == PlcResponseCode.OK) {
                    logger.info("Discovery finished successfully:");
                    for (PlcDiscoveryItem discoveryItem : discoveryResponse.getValues()) {
                        logger.info("Found device with name: {} with connection url: {}",
                            discoveryItem.getName(), discoveryItem.getConnectionUrl());
                    }
                }*/
            }
        }
    }

}
