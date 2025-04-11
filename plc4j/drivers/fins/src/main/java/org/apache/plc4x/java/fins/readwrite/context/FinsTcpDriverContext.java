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
package org.apache.plc4x.java.fins.readwrite.context;

import org.apache.plc4x.java.fins.readwrite.configuration.FinsTcpConfiguration;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.context.DriverContext;

import java.time.Duration;

public class FinsTcpDriverContext implements DriverContext, HasConfiguration<FinsTcpConfiguration> {

    private boolean passiveMode = false;
    private int callingTsapId;
    private int calledTsapId;
    private int pduSize;
    private int maxAmqCaller;
    private int maxAmqCallee;

    private int calledTsapId2;
    private int readTimeout;
    private boolean ping;
    private int pingTime;
    private int retryTime;

    @Override
    public void setConfiguration(FinsTcpConfiguration configuration) {

        if (configuration.localTsap > 0) {
            this.callingTsapId = configuration.localTsap;
        }
        if (configuration.remoteTsap > 0) {
            this.calledTsapId = configuration.remoteTsap;
        }

        this.maxAmqCaller = configuration.maxAmqCaller;
        this.maxAmqCallee = configuration.maxAmqCallee;

        this.readTimeout = configuration.readTimeout;
        this.ping = configuration.ping;
        this.pingTime = (configuration.pingTime == 0) ? 10 : configuration.pingTime;
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

    public int getReadTimeout() {
        return readTimeout;
    }

    public void setReadTimeout(int readTimeout) {
        this.readTimeout = readTimeout;
    }

    public Duration getReadTimeoutDuration() {
        return Duration.ofMillis(readTimeout);
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

}
