/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.base.connection;

import org.apache.plc4x.java.spi.parser.ConfigurationParameter;
import org.apache.plc4x.java.spi.parser.DoubleDefaultValue;
import org.apache.plc4x.java.spi.parser.IntDefaultValue;

public class PcapConfiguration {

    @ConfigurationParameter("protocol-id")
    @IntDefaultValue(-1)
    private Integer protocolId;

    @ConfigurationParameter("replay-speed-factor")
    @DoubleDefaultValue(1)
    private float replaySpeedFactor;

    @ConfigurationParameter("packet-handler")
    private String packetHandler;

    public Integer getProtocolId() {
        return protocolId;
    }

    public void setProtocolId(Integer protocolId) {
        this.protocolId = protocolId;
    }

    public float getReplaySpeedFactor() {
        return replaySpeedFactor;
    }

    public void setReplaySpeedFactor(float replaySpeedFactor) {
        this.replaySpeedFactor = replaySpeedFactor;
    }

    public String getPacketHandler() {
        return packetHandler;
    }

    public void setPacketHandler(String packetHandler) {
        this.packetHandler = packetHandler;
    }

    @Override
    public String toString() {
        return "PcapConfiguration{" +
            "protocolId=" + protocolId +
            ", replaySpeedFactor=" + replaySpeedFactor +
            ", packetHandler=" + packetHandler +
            '}';
    }

}
