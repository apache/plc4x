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

import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.messages.PlcBrowseRequest;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class HelloPlc4xDiscoverAndBrowse {

    private static final Logger logger = LoggerFactory.getLogger(HelloPlc4xDiscoverAndBrowse.class);

    public static void main(String[] args) throws Exception {
        // Iterate over all installed drivers and execute their browse functionality (If they support it)
        PlcDriverManager driverManager = PlcDriverManager.getDefault();
        PlcConnectionManager connectionManager = driverManager.getConnectionManager();
        for (String protocolCode : driverManager.listDrivers()) {
            PlcDriver driver = driverManager.getDriver(protocolCode);
            if (driver.getMetadata().canDiscover()) {
                logger.info("Performing discovery for {} protocol", driver.getProtocolName());

                PlcDiscoveryRequest discoveryRequest = driver.discoveryRequestBuilder().build();

                discoveryRequest.executeWithHandler(discoveryItem -> {
                    logger.info(" - Found device with connection-url {}", discoveryItem.getConnectionUrl());
                    try (PlcConnection connection = connectionManager.getConnection(discoveryItem.getConnectionUrl())) {
                        if (connection.getMetadata().canBrowse()) {
                            PlcBrowseRequest browseRequest = connection.browseRequestBuilder().build();
                            browseRequest.execute().whenComplete((browseResponse, throwable) -> {
                                if (throwable != null) {
                                    throwable.printStackTrace();
                                } else {
                                    for (String queryName : browseResponse.getQueryNames()) {
                                        for (PlcBrowseItem value : browseResponse.getValues(queryName)) {
                                            outputBrowseItem(value, 0);
                                        }
                                    }
                                }
                            });
                        }
                    } catch (Exception e) {
                        throw new RuntimeException(e);
                    }
                });
            }
        }
    }

    protected static void outputBrowseItem(PlcBrowseItem browseItem, int indent) {
        System.out.printf("%s%s : %s (%s %s %s)%n",
            StringUtils.repeat("   ", Math.max(0, indent)),
            browseItem.getTag().getAddressString(),
            browseItem.getTag().getPlcValueType().name(),
            browseItem.isReadable() ? "R" : " ",
            browseItem.isWritable() ? "W" : " ",
            browseItem.isSubscribable() ? "S" : " ");
        if (!browseItem.getChildren().isEmpty()) {
            for (PlcBrowseItem child : browseItem.getChildren().values()) {
                outputBrowseItem(child, indent + 1);
            }
        }
    }

}
