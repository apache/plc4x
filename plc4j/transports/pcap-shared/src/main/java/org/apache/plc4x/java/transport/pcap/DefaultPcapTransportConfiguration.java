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

package org.apache.plc4x.java.transport.pcap;

import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

public abstract class DefaultPcapTransportConfiguration implements PcapTransportConfiguration {

    @ConfigurationParameter("support-vlans")
    @BooleanDefaultValue(false)
    @Description("Should VLan packets be automatically unpacked?")
    private boolean supportVlans;

    @ConfigurationParameter("protocol-id")
    @IntDefaultValue(-1)
    @Description("When provided, filters all packets to let only packets matching this ethernet protocol-id pass.")
    private int protocolId;

    @Override
    public boolean getSupportVlans() {
        return supportVlans;
    }

    public void setSupportVlans(boolean supportVlans) {
        this.supportVlans = supportVlans;
    }

    @Override
    public int getProtocolId() {
        return protocolId;
    }

    public void setProtocolId(int protocolId) {
        this.protocolId = protocolId;
    }

}
