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
 * COTP local/remote TSAPs from the S7 rack/slot/device-group parameters when these are provided
 * (the {@code local-tsap}/{@code remote-tsap} parameters from the parent {@link CotpTransportConfiguration}
 * still take precedence when set explicitly to a non-zero value).
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
    @StringDefaultValue("PG_OR_PC")
    @Description("Local Device Group.")
    protected DeviceGroup localDeviceGroup = DeviceGroup.PG_OR_PC;

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

    public S7CotpTransportConfiguration() {
        super();
        // Default TSAPs: derive from rack/slot/device-group right away. The parser will overwrite
        // these via setLocalTsap/setRemoteTsap (or via the rack/slot setters below) if the user
        // supplied explicit values.
        this.localTsap = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, 1, 1) & 0xFFFF;
        this.remoteTsap = S7TsapIdEncoder.encodeS7TsapId(DeviceGroup.PG_OR_PC, 0, 0) & 0xFFFF;
    }

    @Override
    public int getDefaultPort() {
        return ISO_ON_TCP_PORT;
    }

    private void recomputeTsaps() {
        this.localTsap = S7TsapIdEncoder.encodeS7TsapId(
            localDeviceGroup, localRack, localSlot) & 0xFFFF;
        this.remoteTsap = S7TsapIdEncoder.encodeS7TsapId(
            remoteDeviceGroup, remoteRack, remoteSlot) & 0xFFFF;
    }

    public void setLocalRack(int v) { this.localRack = v; recomputeTsaps(); }
    public void setLocalSlot(int v) { this.localSlot = v; recomputeTsaps(); }
    public void setLocalDeviceGroup(DeviceGroup v) { this.localDeviceGroup = v; recomputeTsaps(); }
    public void setRemoteRack(int v) { this.remoteRack = v; recomputeTsaps(); }
    public void setRemoteSlot(int v) { this.remoteSlot = v; recomputeTsaps(); }
    public void setRemoteDeviceGroup(DeviceGroup v) { this.remoteDeviceGroup = v; recomputeTsaps(); }

    public int getLocalRack() { return localRack; }
    public int getLocalSlot() { return localSlot; }
    public DeviceGroup getLocalDeviceGroup() { return localDeviceGroup; }
    public int getRemoteRack() { return remoteRack; }
    public int getRemoteSlot() { return remoteSlot; }
    public DeviceGroup getRemoteDeviceGroup() { return remoteDeviceGroup; }

    public void setLocalTsap(int localTsap) {
        if (localTsap != 0) {
            this.localTsap = localTsap;
        }
    }

    public void setRemoteTsap(int remoteTsap) {
        if (remoteTsap != 0) {
            this.remoteTsap = remoteTsap;
        }
    }

}
