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
package org.apache.plc4x.java.s7.context;

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.configuration.S7Configuration;

/**
 * Per-connection mutable S7 protocol state. Replaces the SPI1 driver context that the
 * old {@code s7light} stack used; now plain POJO held by {@code S7CotpConnection}.
 */
public class S7DriverContext {

    private int pduSize = 1024;
    private int maxAmqCaller = 8;
    private int maxAmqCallee = 8;
    private ControllerType controllerType = ControllerType.ANY;
    private int readTimeout = 10000;

    /**
     * Whether the device responded successfully to the SZL probe at connect time. Drives
     * the per-device {@code isBrowseSupported} / {@code isSubscribeSupported} metadata
     * flags — and, by extension, whether the connection wires up the S7Comm UserData
     * services on first use. LOGO and other low-end devices that don't speak SZL get
     * {@code false} here and only ever see plain Read/Write var traffic.
     */
    private boolean userDataServicesSupported = false;
    private String articleNumber = "";

    public void applyConfiguration(S7Configuration configuration) {
        this.pduSize = configuration.getPduSize();
        this.maxAmqCaller = configuration.getMaxAmqCaller();
        this.maxAmqCallee = configuration.getMaxAmqCallee();
        this.readTimeout = configuration.getReadTimeout();
        this.controllerType = configuration.getControllerType() == null
            ? ControllerType.ANY : configuration.getControllerType();
        // The Siemens LOGO device only supports a small PDU size.
        if (this.controllerType == ControllerType.LOGO && this.pduSize == 1024) {
            this.pduSize = 480;
        }
    }

    public int getPduSize() { return pduSize; }
    public void setPduSize(int v) { this.pduSize = v; }

    public int getMaxAmqCaller() { return maxAmqCaller; }
    public void setMaxAmqCaller(int v) { this.maxAmqCaller = v; }

    public int getMaxAmqCallee() { return maxAmqCallee; }
    public void setMaxAmqCallee(int v) { this.maxAmqCallee = v; }

    public ControllerType getControllerType() { return controllerType; }
    public void setControllerType(ControllerType v) { this.controllerType = v; }

    public int getReadTimeout() { return readTimeout; }
    public void setReadTimeout(int v) { this.readTimeout = v; }

    public boolean isUserDataServicesSupported() { return userDataServicesSupported; }
    public void setUserDataServicesSupported(boolean v) { this.userDataServicesSupported = v; }

    public String getArticleNumber() { return articleNumber; }
    public void setArticleNumber(String v) { this.articleNumber = v == null ? "" : v; }
}
