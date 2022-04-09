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
package org.apache.plc4x.java.s7.readwrite.configuration;

import org.apache.plc4x.java.s7.readwrite.S7Driver;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.tcp.TcpTransportConfiguration;

public class S7Configuration implements Configuration, TcpTransportConfiguration {

    @ConfigurationParameter("local-rack")
    @IntDefaultValue(1)
    public int localRack = 1;

    @ConfigurationParameter("local-slot")
    @IntDefaultValue(1)
    public int localSlot = 1;
    
    @ConfigurationParameter("local-tsap")
    @IntDefaultValue(0)
    public int localTsap = 0;

    @ConfigurationParameter("remote-rack")
    @IntDefaultValue(0)
    public int remoteRack = 0;

    @ConfigurationParameter("remote-slot")
    @IntDefaultValue(0)
    public int remoteSlot = 0;

    @ConfigurationParameter("remote-tsap")
    @IntDefaultValue(0)
    public int remoteTsap = 0;

    @ConfigurationParameter("pdu-size")
    @IntDefaultValue(1024)
    public int pduSize = 1024;

    @ConfigurationParameter("max-amq-caller")
    @IntDefaultValue(8)
    public int maxAmqCaller = 8;

    @ConfigurationParameter("max-amq-callee")
    @IntDefaultValue(8)
    public int maxAmqCallee = 8;

    @ConfigurationParameter("controller-type")
    public String controllerType;

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

    /**
     * Per default port for the S7 protocol is 102.
     * @return 102
     */
    @Override
    public int getDefaultPort() {
        return S7Driver.ISO_ON_TCP_PORT;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            "local-rack=" + localRack +
            ", local-slot=" + localSlot +
            ", local-tsap=" + localTsap +
            ", remote-rack=" + remoteRack +
            ", remote-slot=" + remoteSlot +
            ", remote-tsap=" + remoteTsap +
            ", pduSize=" + pduSize +
            ", maxAmqCaller=" + maxAmqCaller +
            ", maxAmqCallee=" + maxAmqCallee +
            ", controllerType='" + controllerType + '\'' +
            '}';
    }

}
