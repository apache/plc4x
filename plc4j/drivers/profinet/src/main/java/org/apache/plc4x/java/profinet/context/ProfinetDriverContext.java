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
package org.apache.plc4x.java.profinet.context;

import org.apache.plc4x.java.profinet.config.ProfinetConfiguration;
import org.apache.plc4x.java.profinet.readwrite.DceRpc_ActivityUuid;
import org.apache.plc4x.java.profinet.readwrite.IpAddress;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.net.DatagramSocket;

public class ProfinetDriverContext  implements DriverContext, HasConfiguration<ProfinetConfiguration> {

    private DceRpc_ActivityUuid dceRpc_activityUuid;
    private MacAddress localMacAddress;
    private IpAddress localIpAddress;
    private int localUdpPort;
    private MacAddress remoteMacAddress;
    private IpAddress remoteIpAddress;
    private int remoteUdpPort;
    private int sessionKey;

    private DatagramSocket udpSocket;

    @Override
    public void setConfiguration(ProfinetConfiguration configuration) {

    }

    public DceRpc_ActivityUuid getDceRpcActivityUuid() {
        return dceRpc_activityUuid;
    }

    public void setDceRpcActivityUuid(DceRpc_ActivityUuid dceRpc_activityUuid) {
        this.dceRpc_activityUuid = dceRpc_activityUuid;
    }

    public MacAddress getLocalMacAddress() {
        return localMacAddress;
    }

    public void setLocalMacAddress(MacAddress localMacAddress) {
        this.localMacAddress = localMacAddress;
    }

    public IpAddress getLocalIpAddress() {
        return localIpAddress;
    }

    public void setLocalIpAddress(IpAddress localIpAddress) {
        this.localIpAddress = localIpAddress;
    }

    public int getLocalUdpPort() {
        return localUdpPort;
    }

    public void setLocalUdpPort(int localUdpPort) {
        this.localUdpPort = localUdpPort;
    }

    public MacAddress getRemoteMacAddress() {
        return remoteMacAddress;
    }

    public void setRemoteMacAddress(MacAddress remoteMacAddress) {
        this.remoteMacAddress = remoteMacAddress;
    }

    public IpAddress getRemoteIpAddress() {
        return remoteIpAddress;
    }

    public void setRemoteIpAddress(IpAddress remoteIpAddress) {
        this.remoteIpAddress = remoteIpAddress;
    }

    public int getRemoteUdpPort() {
        return remoteUdpPort;
    }

    public void setRemoteUdpPort(int remoteUdpPort) {
        this.remoteUdpPort = remoteUdpPort;
    }

    public int getSessionKey() {
        return sessionKey;
    }

    public void setSessionKey(int sessionKey) {
        this.sessionKey = sessionKey;
    }

    public DatagramSocket getUdpSocket() {
        return udpSocket;
    }

    public void setUdpSocket(DatagramSocket udpSocket) {
        this.udpSocket = udpSocket;
    }
}
