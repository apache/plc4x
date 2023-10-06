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

import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.pcap4j.core.PcapHandle;
import org.pcap4j.core.PcapIpV4Address;

import java.util.List;

public class LocalNetworkDevice {

    private final MacAddress macAddress;
    private final List<PcapIpV4Address> ipAddresses;
    private final PcapHandle handle;

    public LocalNetworkDevice(MacAddress macAddress, List<PcapIpV4Address> ipAddresses, PcapHandle handle) {
        this.macAddress = macAddress;
        this.ipAddresses = ipAddresses;
        this.handle = handle;
    }

    public MacAddress getMacAddress() {
        return macAddress;
    }

    public List<PcapIpV4Address> getIpAddresses() {
        return ipAddresses;
    }

    public PcapHandle getHandle() {
        return handle;
    }

}
