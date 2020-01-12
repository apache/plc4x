/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.knxnetip.configuration;

import org.apache.plc4x.java.knxnetip.KnxNetIpDriver;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.udp.UdpTransportConfiguration;

public class KnxNetIpConfiguration implements Configuration, UdpTransportConfiguration {

    @ConfigurationParameter("knxproj-file-path")
    public String knxprojFilePath;

    @ConfigurationParameter("group-address-type")
    @IntDefaultValue(3)
    public int groupAddressType = 3;

    public String getKnxprojFilePath() {
        return knxprojFilePath;
    }

    public void setKnxprojFilePath(String knxprojFilePath) {
        this.knxprojFilePath = knxprojFilePath;
    }

    public int getGroupAddressType() {
        return groupAddressType;
    }

    public void setGroupAddressType(int groupAddressType) {
        this.groupAddressType = groupAddressType;
    }

    @Override
    public int getDefaultPort() {
        return KnxNetIpDriver.KNXNET_IP_PORT;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            "knxprojFilePath=" + knxprojFilePath + ", " +
            "groupAddressType=" + groupAddressType +
            '}';
    }

}
