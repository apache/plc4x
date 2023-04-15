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
import org.pcap4j.core.PcapNetworkInterface;

public class ProfinetNetworkInterface implements NetworkInterface {

    private String ipAddress;
    private String subnet;
    private String gateway;

    public ProfinetNetworkInterface(PcapNetworkInterface networkInterface) {
        if (networkInterface.getAddresses().size() == 0) {
            throw new PlcRuntimeException("No Inet4 Address assigned to interface");
        }
        for (PcapAddress address : networkInterface.getAddresses()) {
            ipAddress = address.getAddress().toString();
            subnet = address.getNetmask().toString();
            gateway = address.getDestinationAddress().toString();
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
        return null;
    }

    public byte[] getSubnetAsByteArray() {
        return null;
    }

    public byte[] getGatewayAsByteArray() {
        return null;
    }

}
