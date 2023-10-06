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

package org.apache.plc4x.java.profinet.device;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.pcap4j.core.PcapAddress;
import org.pcap4j.core.PcapIpV4Address;
import org.pcap4j.core.PcapNetworkInterface;

import java.net.InetAddress;
import java.net.UnknownHostException;

public class ProfinetNetworkInterface implements NetworkInterface {

    private String ipAddress;
    private String subnet;
    private String gateway;

    public ProfinetNetworkInterface(PcapNetworkInterface networkInterface) {
        if (networkInterface.getAddresses().size() == 0) {
            throw new PlcRuntimeException("No Inet4 Address assigned to interface");
        }
        for (PcapAddress address : networkInterface.getAddresses()) {
            if (address instanceof PcapIpV4Address) {
                ipAddress = address.getAddress().toString();
                subnet = address.getNetmask().toString();
                if (address.getDestinationAddress() != null) {
                    gateway = address.getDestinationAddress().toString();
                } else {
                    gateway = "0.0.0.0";
                }
            }
        }

    }

    @Override
    public String getIpAddress() {
        return ipAddress;
    }

    @Override
    public String getSubnet() {
        return subnet;
    }

    @Override
    public String getGateway() {
        return gateway;
    }

    public byte[] getIpAddressAsByteArray() {
        try {
            if (this.ipAddress != null) {
                return InetAddress.getByName(this.ipAddress.replace("/", "")).getAddress();
            }
            return new byte[4];
        } catch (UnknownHostException e) {
            return new byte[4];
        }
    }

    public byte[] getSubnetAsByteArray() {
        try {
            if (this.subnet != null) {
                return InetAddress.getByName(this.subnet).getAddress();
            }
            return new byte[4];
        } catch (UnknownHostException e) {
            return new byte[4];
        }
    }

    public byte[] getGatewayAsByteArray() {
        try {
            if (this.gateway != null) {
                return InetAddress.getByName(this.gateway).getAddress();
            }
            return new byte[4];
        } catch (UnknownHostException e) {
            return new byte[4];
        }
    }

}
