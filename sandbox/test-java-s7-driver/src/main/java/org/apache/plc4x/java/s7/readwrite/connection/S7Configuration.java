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

package org.apache.plc4x.java.s7.readwrite.connection;

import org.apache.plc4x.java.spi.parser.ConfigurationParameter;
import org.apache.plc4x.java.spi.parser.IntDefaultValue;
import org.apache.plc4x.java.spi.parser.Required;

public class S7Configuration {

    @ConfigurationParameter("rackId")
    @IntDefaultValue(1)
    public int rack = 1;

    @ConfigurationParameter("slotId")
    @IntDefaultValue(1)
    public int slot = 1;

    @ConfigurationParameter
    @IntDefaultValue(1024)
    public int pduSize = 1024;

    @ConfigurationParameter
    @IntDefaultValue(8)
    public int maxAmqCaller = 8;

    @ConfigurationParameter
    @IntDefaultValue(8)
    public int maxAmqCallee = 8;

    @ConfigurationParameter
    public String controllerType;

    public int getRack() {
        return rack;
    }

    public void setRack(int rack) {
        this.rack = rack;
    }

    public int getSlot() {
        return slot;
    }

    public void setSlot(int slot) {
        this.slot = slot;
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

    @Override
    public String toString() {
        return "Configuration{" +
            "rack=" + rack +
            ", slot=" + slot +
            ", pduSize=" + pduSize +
            ", maxAmqCaller=" + maxAmqCaller +
            ", maxAmqCallee=" + maxAmqCallee +
            ", controllerType='" + controllerType + '\'' +
            '}';
    }

}
