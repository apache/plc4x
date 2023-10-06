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

package org.apache.plc4x.java.profinet;

import org.apache.plc4x.java.profinet.device.NetworkInterface;

public class DummyNetworkInterface implements NetworkInterface {

    private String ipAddress;
    private String subnet;
    private String gateway;

    public DummyNetworkInterface(String ipAddress, String subnet, String gateway) {
        this.ipAddress = ipAddress;
        this.subnet = subnet;
        this.gateway = gateway;
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
