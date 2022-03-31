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
package org.apache.plc4x.java.knxnetip.context;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.knxnetip.configuration.KnxNetIpConfiguration;
import org.apache.plc4x.java.knxnetip.ets5.Ets5Parser;
import org.apache.plc4x.java.knxnetip.ets5.model.Ets5Model;
import org.apache.plc4x.java.knxnetip.readwrite.IPAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxAddress;
import org.apache.plc4x.java.knxnetip.readwrite.KnxLayer;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.io.File;

public class KnxNetIpDriverContext implements DriverContext, HasConfiguration<KnxNetIpConfiguration> {

    private boolean passiveMode = false;
    private KnxAddress gatewayAddress;
    private String gatewayName;
    private IPAddress localIPAddress;
    private int localPort;
    private short communicationChannelId;
    private KnxAddress clientKnxAddress;
    private byte groupAddressType;
    private KnxLayer tunnelConnectionType;
    private Ets5Model ets5Model;

    @Override
    public void setConfiguration(KnxNetIpConfiguration configuration) {
        if (configuration.knxprojFilePath != null) {
            File knxprojFile = new File(configuration.knxprojFilePath);
            if (knxprojFile.exists() && knxprojFile.isFile()) {
                ets5Model = new Ets5Parser().parse(knxprojFile, configuration.knxprojPassword);
                groupAddressType = ets5Model.getGroupAddressType();
            } else {
                throw new PlcRuntimeException(String.format(
                    "File specified with 'knxproj-file-path' does not exist or is not a file: '%s'",
                    configuration.knxprojFilePath));
            }
        } else {
            groupAddressType = (byte) configuration.groupAddressNumLevels;
        }
        tunnelConnectionType = KnxLayer.valueOf("TUNNEL_" + configuration.getConnectionType());
    }

    public boolean isPassiveMode() {
        return passiveMode;
    }

    public void setPassiveMode(boolean passiveMode) {
        this.passiveMode = passiveMode;
    }

    public KnxAddress getGatewayAddress() {
        return gatewayAddress;
    }

    public void setGatewayAddress(KnxAddress gatewayAddress) {
        this.gatewayAddress = gatewayAddress;
    }

    public String getGatewayName() {
        return gatewayName;
    }

    public void setGatewayName(String gatewayName) {
        this.gatewayName = gatewayName;
    }

    public IPAddress getLocalIPAddress() {
        return localIPAddress;
    }

    public void setLocalIPAddress(IPAddress localIPAddress) {
        this.localIPAddress = localIPAddress;
    }

    public int getLocalPort() {
        return localPort;
    }

    public void setLocalPort(int localPort) {
        this.localPort = localPort;
    }

    public short getCommunicationChannelId() {
        return communicationChannelId;
    }

    public void setCommunicationChannelId(short communicationChannelId) {
        this.communicationChannelId = communicationChannelId;
    }

    public KnxAddress getClientKnxAddress() {
        return clientKnxAddress;
    }

    public void setClientKnxAddress(KnxAddress clientKnxAddress) {
        this.clientKnxAddress = clientKnxAddress;
    }

    public byte getGroupAddressType() {
        return groupAddressType;
    }

    public KnxLayer getTunnelConnectionType() {
        return tunnelConnectionType;
    }

    public Ets5Model getEts5Model() {
        return ets5Model;
    }

}
