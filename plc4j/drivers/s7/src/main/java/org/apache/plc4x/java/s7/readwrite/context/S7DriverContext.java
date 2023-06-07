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
package org.apache.plc4x.java.s7.readwrite.context;

import org.apache.plc4x.java.s7.readwrite.COTPTpduSize;
import org.apache.plc4x.java.s7.readwrite.DeviceGroup;
import org.apache.plc4x.java.s7.readwrite.configuration.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.types.S7ControllerType;
import org.apache.plc4x.java.s7.readwrite.utils.S7TsapIdEncoder;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;

public class S7DriverContext implements DriverContext, HasConfiguration<S7Configuration> {

    private boolean passiveMode = false;
    private int callingTsapId;
    private int calledTsapId;
    private COTPTpduSize cotpTpduSize;
    private int pduSize;
    private int maxAmqCaller;
    private int maxAmqCallee;
    private S7ControllerType controllerType;
 
    
    private int calledTsapId2; 
    private int readTimeout;
    private boolean ping;
    private int pingTime;   
    private int retryTime;

    @Override
    public void setConfiguration(S7Configuration configuration) {
        this.callingTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.OTHERS,
            configuration.localRack, configuration.localSlot);
        this.calledTsapId = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC,
            configuration.remoteRack, configuration.remoteSlot);
        
        this.calledTsapId2 = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC,
            configuration.remoteRack2, configuration.remoteSlot2);        
     

        if (configuration.localTsap > 0) {
        	this.callingTsapId = configuration.localTsap;
        }
        if (configuration.remoteTsap > 0) {
        	this.calledTsapId = configuration.remoteTsap;
        }
        this.controllerType = configuration.controllerType == null ? S7ControllerType.ANY : S7ControllerType.valueOf(configuration.controllerType);

        // Initialize the parameters with initial version (Will be updated during the login process)
        this.cotpTpduSize = getNearestMatchingTpduSize((short) configuration.getPduSize());

        // The Siemens LOGO device seems to only work with very limited settings,
        // so we're overriding some of the defaults.
        if (this.controllerType == S7ControllerType.LOGO && configuration.pduSize == 1024) {
            configuration.pduSize = 480;
            this.pduSize = 480;
        } else {
            // The PDU size is theoretically not bound by the COTP TPDU size, however having a larger
            // PDU size would make the code extremely complex. But even if the protocol would allow this
            // I have never seen this happen in reality. Making is smaller would unnecessarily limit the
            // size, so we're setting it to the maximum that can be included.
            this.pduSize = cotpTpduSize.getSizeInBytes() - 16;
        }

        this.maxAmqCaller = configuration.maxAmqCaller;
        this.maxAmqCallee = configuration.maxAmqCallee;
        
        this.readTimeout = configuration.readTimeout;
        this.ping = configuration.ping;
        this.pingTime = (configuration.pingTime == 0)?10:configuration.pingTime;
        this.retryTime = configuration.retryTime;
    }

    public boolean isPassiveMode() {
        return passiveMode;
    }

    public void setPassiveMode(boolean passiveMode) {
        this.passiveMode = passiveMode;
    }

    public int getCallingTsapId() {
        return callingTsapId;
    }

    public void setCallingTsapId(int callingTsapId) {
        this.callingTsapId = callingTsapId;
    }

    public int getCalledTsapId() {
        return calledTsapId;
    }

    public void setCalledTsapId(int calledTsapId) {
        this.calledTsapId = calledTsapId;
    }
    
    public int getCalledTsapId2() {
        return calledTsapId2;
    }

    public void setCalledTsapId2(int calledTsapId2) {
        this.calledTsapId2 = calledTsapId2;
    }    

    public COTPTpduSize getCotpTpduSize() {
        return cotpTpduSize;
    }

    public void setCotpTpduSize(COTPTpduSize cotpTpduSize) {
        this.cotpTpduSize = cotpTpduSize;
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

    public S7ControllerType getControllerType() {
        return controllerType;
    }

    public void setControllerType(S7ControllerType controllerType) {
        this.controllerType = controllerType;
    }
    
    public int getReadTimeout() {
        return readTimeout;
    }

    public void setReadTimeout(int readTimeout) {
        this.readTimeout = readTimeout;
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

    /**
     * Iterate over all values until one is found that the given tpdu size will fit.
     *
     * @param tpduSizeParameter requested tpdu size.
     * @return smallest {@link COTPTpduSize} which will fit a given size of tpdu.
     */
    protected COTPTpduSize getNearestMatchingTpduSize(short tpduSizeParameter) {
        for (COTPTpduSize value : COTPTpduSize.values()) {
            if (value.getSizeInBytes() >= tpduSizeParameter) {
                return value;
            }
        }
        return null;
    }

}
