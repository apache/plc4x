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
package org.apache.plc4x.java.s7.readwrite.configuration;

import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

public class S7Configuration implements PlcConnectionConfiguration {

    @ConfigurationParameter("local-rack")
    @IntDefaultValue(1)
    @Description("Rack value for the client (PLC4X device).")
    public int localRack = 1;

    @ConfigurationParameter("local-slot")
    @IntDefaultValue(1)
    @Description("Slot value for the client (PLC4X device).")
    public int localSlot = 1;

    @ConfigurationParameter("local-tsap")
    @IntDefaultValue(0)
    @Description("Local Transport Service Access Point.")
    public int localTsap = 0;

    @ConfigurationParameter("remote-rack")
    @IntDefaultValue(0)
    @Description("Rack value for the remote main CPU (PLC).")
    public int remoteRack = 0;

    @ConfigurationParameter("remote-slot")
    @IntDefaultValue(0)
    @Description("Slot value for the remote main CPU (PLC).")
    public int remoteSlot = 0;

    @ConfigurationParameter("remote-rack2")
    @IntDefaultValue(0)
    @Description("Rack value for the remote secondary CPU (PLC).")
    public int remoteRack2 = 0;

    @ConfigurationParameter("remote-slot2")
    @IntDefaultValue(0)
    @Description("Slot value for the remote secondary CPU (PLC).")
    public int remoteSlot2 = 0;

    @ConfigurationParameter("remote-tsap")
    @IntDefaultValue(0)
    @Description("Remote Transport Service Access Point.")
    public int remoteTsap = 0;

    @ConfigurationParameter("pdu-size")
    @IntDefaultValue(1024)
    @Description("Maximum size of a data-packet sent to and received from the remote PLC. During the connection process both parties will negotiate a maximum size both parties can work with and is equal or smaller than the given value is used. The driver will automatically split up large requests to not exceed this value in a request or expected response.")
    public int pduSize = 1024;

    @ConfigurationParameter("max-amq-caller")
    @IntDefaultValue(8)
    @Description("Maximum number of unconfirmed requests the PLC will accept in parallel before discarding with errors. This parameter also will be negotiated during the connection process and the maximum both parties can work with and is equal or smaller than the given value is used. The driver will automatically take care not exceeding this value while processing requests. Too many requests can cause a growing queue.")
    public int maxAmqCaller = 8;

    @ConfigurationParameter("max-amq-callee")
    @IntDefaultValue(8)
    @Description("Maximum number of unconfirmed responses or requests PLC4X will accept in parallel before discarding with errors. This option is available for completeness and is correctly handled out during the connection process, however it is currently not enforced on PLC4X’s side. So if a PLC would send more messages than agreed upon, these would still be processed.")
    public int maxAmqCallee = 8;

    @ConfigurationParameter("controller-type")
    @Description("As part of the connection process, usually the PLC4X S7 driver would try to identify the remote device. However some devices seem to have problems with this and hang up or cause other problems. In such a case, providing the controller-type will skip the identification process and hereby avoid this type of problem. Possible values are:/n- S7_300\n- S7_400\n- S7_1200\n- S7-1500\n- LOGO")
    public String controllerType;

    @ConfigurationParameter("read-timeout")
    @IntDefaultValue(0)
    @Description("This is the maximum waiting time for reading on the TCP channel. As there is no traffic, it must be assumed that the connection with the interlocutor was lost and it must be restarted. When the channel is closed, the \"fail over\" is carried out in case of having the secondary channel, or it is expected that it will be restored automatically, which is done every 4 seconds.")
    public int readTimeout = 0;

    @ConfigurationParameter("ping")
    @BooleanDefaultValue(false)
    @Description("Time for supervision of TCP channels. If the channel is not active, a safe stop of the EventLoop must be performed, to ensure that no additional tasks are created.")
    public boolean ping = false;

    @ConfigurationParameter("ping-time")
    @IntDefaultValue(0)
    @Description("If your application requires sampling times greater than the set \"read-timeout\" time, it is important that the PING option is activated, this will prevent the TCP channel from being closed unnecessarily.")
    public int pingTime = 0;

    @ConfigurationParameter("retry-time")
    @IntDefaultValue(0)
    @Description("Time value in seconds at which the execution of the PING will be scheduled. Generally set by developer experience, but generally should be the same as (read-timeout / 2).")
    public int retryTime = 0;

    public int getLocalRack() {
        return localRack;
    }

    public void setLocalRack(int localRack) {
        this.localRack = localRack;
    }

    public int getLocalSlot() {
        return localSlot;
    }

    public void setLocalSlot(int localSlot) {
        this.localSlot = localSlot;
    }

    public int getLocalTsap() {
        return localTsap;
    }

    public void setLocalTsap(int localTsap) {
        this.localTsap = localTsap;
    }

    public int getRemoteRack() {
        return remoteRack;
    }

    public void setRemoteRack(int remoteRack) {
        this.remoteRack = remoteRack;
    }

    public int getRemoteSlot() {
        return remoteSlot;
    }

    public void setRemoteSlot(int remoteSlot) {
        this.remoteSlot = remoteSlot;
    }

    public int getRemoteRack2() {
        return remoteRack2;
    }

    public void setRemoteRack2(int remoteRack2) {
        this.remoteRack2 = remoteRack2;
    }

    public int getRemoteSlot2() {
        return remoteSlot2;
    }

    public void setRemoteSlot2(int remoteSlot2) {
        this.remoteSlot2 = remoteSlot2;
    }

    public int getRemoteTsap() {
        return remoteTsap;
    }

    public void setRemoteTsap(int remoteTsap) {
        this.remoteTsap = remoteTsap;
    }

    public int getPduSize() {
        return pduSize;
    }

    public void setPduSize(int pduSize) {
        this.pduSize = pduSize;
    }

    public int getMaxAmqCaller() {
        return maxAmqCaller;
    }

    public void setMaxAmqCaller(int maxAmqCaller) {
        this.maxAmqCaller = maxAmqCaller;
    }

    public int getMaxAmqCallee() {
        return maxAmqCallee;
    }

    public void setMaxAmqCallee(int maxAmqCallee) {
        this.maxAmqCallee = maxAmqCallee;
    }

    public String getControllerType() {
        return controllerType;
    }

    public void setControllerType(String controllerType) {
        this.controllerType = controllerType;
    }

    public int getReadTimeout() {
        return readTimeout;
    }

    public void setReadTimeout(int readTimeOut) {
        this.readTimeout = readTimeOut;
    }

    public boolean getPing() {
        return ping;
    }

    public void setPing(boolean ping) {
        this.ping = ping;
    }

    public int getPingTime() {
        return pingTime;
    }

    public void setPingTime(int pingTime) {
        this.pingTime = pingTime;
    }

    public int getRetryTime() {
        return retryTime;
    }

    public void setRetryTime(int retryTime) {
        this.retryTime = retryTime;
    }
    
    @Override
    public String toString() {
        return "Configuration{" +
            "local-rack=" + localRack +
            ", local-slot=" + localSlot +
            ", local-tsap=" + localTsap +
            ", remote-rack=" + remoteRack +
            ", remote-slot=" + remoteSlot +
            ", remote-rack2=" + remoteRack2 +
            ", remote-slot2=" + remoteSlot2 +
            ", remote-tsap=" + remoteTsap +
            ", pduSize=" + pduSize +
            ", maxAmqCaller=" + maxAmqCaller +
            ", maxAmqCallee=" + maxAmqCallee +
            ", controllerType='" + controllerType +
            ", readTimeOut='" + readTimeout +
            ", ping='" + ping +
            ", pingTime='" + pingTime +
            ", retryTime='" + retryTime +
            '\'' +
            '}';
    }
            
}

