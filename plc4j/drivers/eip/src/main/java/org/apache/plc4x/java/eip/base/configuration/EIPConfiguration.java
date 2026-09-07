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
package org.apache.plc4x.java.eip.base.configuration;

import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.Since;
import org.apache.plc4x.java.spi.config.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;

public class EIPConfiguration implements Configuration {

    @ConfigurationParameter
    @IntDefaultValue(1)
    @Description("Without using routing information the backplane defaults to 1. This is overridden if communicationPath is provided.")
    private int backplane = 1;

    @ConfigurationParameter
    @IntDefaultValue(0)
    @Description("The slot within the backplane the CPU is located.")
    private int slot = 0;

    @ConfigurationParameter("big-endian")
    @BooleanDefaultValue(true)
    @Description("Configure if the connection should be set to transport data in Big-Endian format, or not.")
    private boolean bigEndian = true;

    @ConfigurationParameter("connection-serial-number")
    @IntDefaultValue(0)
    @Description("Connection serial number to use in Forward_Open. CIP wants this unique per " +
        "connection, so the default of 0 means 'pick a random one per connection'. Set it " +
        "explicitly only when the exchange has to be reproducible, e.g. in recorded tests.")
    @Since("1.0.0")
    private int connectionSerialNumber = 0;

    @ConfigurationParameter("force-unconnected-operation")
    @BooleanDefaultValue(false)
    @Description("Forces the driver to use unconnected requests.")
    @Since("0.13.0")
    private boolean forceUnconnectedOperation = false;

    @ConfigurationParameter("request-timeout-ms")
    @IntDefaultValue(10_000)
    @Description("Default timeout for all types of requests.")
    private int requestTimeout;

    @ConfigurationParameter("communication-path")
    @Description("The communication path allows for connection routing across multiple backplanes. " +
        "It uses a common format found in Logix controllers.\n" +
        "It consists of pairs of values, each pair begins with either 1 (Backplane) or 2 (Ethernet), " +
        "followed by a slot in the case of a backplane address, or if using Ethernet an ip address. " +
        "e.g. [1,4,2,192.168.0.1,1,1] - Routes to the 4th slot in the first rack, which is an Ethernet " +
        "module, it then connects to the address 192.168.0.1, then finds the module in slot 1.")
    private String communicationPath;

    public int getBackplane() {
        return backplane;
    }

    public void setBackplane(int backplane) {
        this.backplane = backplane;
    }

    public int getSlot() {
        return slot;
    }

    public void setSlot(int slot) {
        this.slot = slot;
    }

    public int getConnectionSerialNumber() {
        return connectionSerialNumber;
    }

    public boolean isBigEndian() {
        return bigEndian;
    }

    public void setBigEndian(boolean bigEndian) {
        this.bigEndian = bigEndian;
    }

    public boolean isForceUnconnectedOperation() {
        return forceUnconnectedOperation;
    }

    public void setForceUnconnectedOperation(boolean forceUnconnectedOperation) {
        this.forceUnconnectedOperation = forceUnconnectedOperation;
    }

    public int getRequestTimeout() {
        return requestTimeout;
    }

    public void setRequestTimeout(int requestTimeout) {
        this.requestTimeout = requestTimeout;
    }

    public String getCommunicationPath() {
        return communicationPath;
    }

    public void setCommunicationPath(String communicationPath) {
        this.communicationPath = communicationPath;
    }

}
