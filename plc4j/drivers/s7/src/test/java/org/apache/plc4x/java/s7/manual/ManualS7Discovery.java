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
package org.apache.plc4x.java.s7.manual;

import org.apache.plc4x.java.api.PlcDriverManager;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryResponse;

/**
 * Manual smoke-test for the S7 driver's PROFINET DCP discovery. Sends an IdentifyAll
 * broadcast on every available network interface, listens for 10 s, and prints the
 * S7Comm-capable devices found.
 *
 * <p>Requires libpcap on the host (and root/admin privileges on most platforms — pcap4j
 * needs raw socket access). On macOS, run with {@code sudo}; on Windows, run as
 * Administrator with Npcap installed.
 */
public class ManualS7Discovery {

    public static void main(String[] args) throws Exception {
        PlcDiscoveryResponse response = PlcDriverManager.getDefault()
            .getDriver("s7")
            .discoveryRequestBuilder()
            .build()
            .execute()
            .get();

        if (response.getValues() == null || response.getValues().isEmpty()) {
            System.out.println("No S7-capable devices discovered.");
            return;
        }
        System.out.printf("Discovered %d S7-capable device(s):%n", response.getValues().size());
        for (PlcDiscoveryItem item : response.getValues()) {
            System.out.printf("  - %s%n", item.getName());
            System.out.printf("    connection-url: %s%n", item.getConnectionUrl());
            item.getAttributes().forEach((k, v) ->
                System.out.printf("    %-16s %s%n", k + ":", v));
        }
    }
}
