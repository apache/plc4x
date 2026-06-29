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
package org.apache.plc4x.java.s7.configuration;

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.spi.config.Configuration;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;

public class S7Configuration implements Configuration {

    // Note: the COTP addressing parameters (local/remote rack, slot, device-group and the
    // local-tsap/remote-tsap overrides) are owned by S7CotpTransportConfiguration, which is the
    // configuration the COTP transport actually consumes to build the connection request.

    @ConfigurationParameter("pdu-size")
    @IntDefaultValue(1024)
    @Description("Maximum size of an S7 data-packet sent to and received from the remote PLC.")
    protected int pduSize = 1024;

    @ConfigurationParameter("max-amq-caller")
    @IntDefaultValue(8)
    @Description("Maximum number of unconfirmed requests the PLC will accept in parallel.")
    protected int maxAmqCaller = 8;

    @ConfigurationParameter("max-amq-callee")
    @IntDefaultValue(8)
    @Description("Maximum number of unconfirmed responses or requests PLC4X will accept in parallel.")
    protected int maxAmqCallee = 8;

    @ConfigurationParameter("controller-type")
    @StringDefaultValue("ANY")
    @Description("Skip controller-type detection and assume the given type.")
    protected ControllerType controllerType = ControllerType.ANY;

    @ConfigurationParameter("read-timeout")
    @IntDefaultValue(10000)
    @Description("Maximum waiting time (in milliseconds) for a single S7 request/response exchange.")
    protected int readTimeout = 10000;

    @ConfigurationParameter("ha-heartbeat-interval")
    @IntDefaultValue(4000)
    @Description("S7H dual-path only: interval between heartbeat ticks (in milliseconds). "
        + "Each tick pings each inner connection so a standby disruption is detected within "
        + "interval + ha-failover-timeout. Lower values detect faster but generate more "
        + "background traffic. Default 4000 (4s).")
    protected int haHeartbeatInterval = 4000;

    @ConfigurationParameter("ha-failover-timeout")
    @IntDefaultValue(2000)
    @Description("S7H dual-path only: maximum time (in milliseconds) the wrapper waits for "
        + "an operation on the active inner before swapping to the alternate. The same value "
        + "is used as the per-tick ping timeout in the heartbeat. Lower values fail over "
        + "faster but risk swapping on transient slow responses. Default 2000 (2s).")
    protected int haFailoverTimeout = 2000;

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

    public int getHaHeartbeatInterval() { return haHeartbeatInterval; }
    public void setHaHeartbeatInterval(int v) { this.haHeartbeatInterval = v; }

    public int getHaFailoverTimeout() { return haFailoverTimeout; }
    public void setHaFailoverTimeout(int v) { this.haFailoverTimeout = v; }

}
