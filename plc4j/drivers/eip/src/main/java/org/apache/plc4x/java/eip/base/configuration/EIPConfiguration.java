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

import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.Since;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.generation.ByteOrder;

public class EIPConfiguration implements PlcConnectionConfiguration {

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

    @ConfigurationParameter("force-unconnected-operation")
    @BooleanDefaultValue(false)
    @Description("Forces the driver to use unconnected requests.")
    @Since("0.13.0")
    private boolean forceUnconnectedOperation = false;

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

    public ByteOrder getByteOrder() {
        return this.bigEndian ? ByteOrder.BIG_ENDIAN : ByteOrder.LITTLE_ENDIAN;
    }

    public void setByteOrder(ByteOrder byteOrder) {
        this.bigEndian = byteOrder == ByteOrder.BIG_ENDIAN;
    }

    public boolean isForceUnconnectedOperation() {
        return forceUnconnectedOperation;
    }

    public void setForceUnconnectedOperation(boolean forceUnconnectedOperation) {
        this.forceUnconnectedOperation = forceUnconnectedOperation;
    }

}
