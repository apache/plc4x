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

import org.apache.plc4x.java.s7.readwrite.DeviceGroup;
import org.apache.plc4x.java.s7.utils.S7TsapIdEncoder;
import org.apache.plc4x.java.spi.config.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.config.annotations.Description;
import org.apache.plc4x.java.spi.config.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.config.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.cotp.config.CotpTransportConfiguration;

/**
 * COTP transport configuration for the S7 driver. Pins the ISO-on-TCP port (102) and derives the
 * COTP local/remote TSAPs from the S7 rack/slot/device-group parameters. The
 * {@code local-tsap}/{@code remote-tsap} parameters inherited from the parent
 * {@link CotpTransportConfiguration} take precedence when set explicitly to a non-zero value.
 *
 * <p>The TSAP is always a derived value - it is never stored. {@link #getLocalTsap()} and
 * {@link #getRemoteTsap()} compute it on demand from the rack/slot/device-group parameters
 * (or return the explicit override). This avoids any drift between the TSAP and the
 * rack/slot inputs it is derived from, which is important because the SPI configuration
 * parser injects field values directly via reflection and never invokes setters.
 */
public class S7CotpTransportConfiguration extends CotpTransportConfiguration {

    public static final int ISO_ON_TCP_PORT = 102;

    @ConfigurationParameter("local-rack")
    @IntDefaultValue(1)
    @Description("Rack value for the client (PLC4X device).")
    protected int localRack = 1;

    @ConfigurationParameter("local-slot")
    @IntDefaultValue(1)
    @Description("Slot value for the client (PLC4X device).")
    protected int localSlot = 1;

    @ConfigurationParameter("local-device-group")
    @StringDefaultValue("OTHERS")
    @Description("Local Device Group.")
    protected DeviceGroup localDeviceGroup = DeviceGroup.OTHERS;

    @ConfigurationParameter("remote-rack")
    @IntDefaultValue(0)
    @Description("Rack value for the remote main CPU (PLC).")
    protected int remoteRack = 0;

    @ConfigurationParameter("remote-slot")
    @IntDefaultValue(0)
    @Description("Slot value for the remote main CPU (PLC).")
    protected int remoteSlot = 0;

    @ConfigurationParameter("remote-device-group")
    @StringDefaultValue("PG_OR_PC")
    @Description("Remote Device Group.")
    protected DeviceGroup remoteDeviceGroup = DeviceGroup.PG_OR_PC;

    @Override
    public int getDefaultPort() {
        return ISO_ON_TCP_PORT;
    }

    @Override
    public int getLocalTsap() {
        return localTsap != 0
            ? localTsap
            : S7TsapIdEncoder.encodeS7TsapId(localDeviceGroup, localRack, localSlot) & 0xFFFF;
    }

    @Override
    public int getRemoteTsap() {
        return remoteTsap != 0
            ? remoteTsap
            : S7TsapIdEncoder.encodeS7TsapId(remoteDeviceGroup, remoteRack, remoteSlot) & 0xFFFF;
    }

    public int getLocalRack() { return localRack; }
    public int getLocalSlot() { return localSlot; }
    public DeviceGroup getLocalDeviceGroup() { return localDeviceGroup; }
    public int getRemoteRack() { return remoteRack; }
    public int getRemoteSlot() { return remoteSlot; }
    public DeviceGroup getRemoteDeviceGroup() { return remoteDeviceGroup; }

}
