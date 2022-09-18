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

package org.apache.plc4x.java.profinet.device;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.messages.PlcDiscoveryItemHandler;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ProfinetDeviceMessageHandler implements PlcDiscoveryItemHandler {

    private HashMap<MacAddress, ProfinetDevice> configuredDevices;

    @Override
    public void handle(PlcDiscoveryItem discoveryItem) {
        try {
            MacAddress macAddress = new MacAddress(Hex.decodeHex(discoveryItem.getOptions().get("MacAddress")));
            if (configuredDevices.containsKey(macAddress)) {
                configuredDevices.get(macAddress).handle(discoveryItem);
            }
        } catch (DecoderException e) {
            throw new RuntimeException(e);
        }
    }

    public void setConfiguredDevices(HashMap<MacAddress, ProfinetDevice> configuredDevices) {
        this.configuredDevices = configuredDevices;
    }
}
