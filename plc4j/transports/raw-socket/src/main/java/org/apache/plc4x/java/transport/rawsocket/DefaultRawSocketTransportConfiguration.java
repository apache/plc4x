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

package org.apache.plc4x.java.transport.rawsocket;

import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Description;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.transport.pcap.DefaultPcapTransportConfiguration;

public abstract class DefaultRawSocketTransportConfiguration extends DefaultPcapTransportConfiguration implements RawSocketTransportConfiguration {

    @ConfigurationParameter("resolve-mac-address")
    @Description("If set to true, the transport will automatically resolve the MAC address for a given IP address (Allows connecting to a raw-socket device using the devices host-name or ip-address).")
    private boolean resolveMacAddress;

    public DefaultRawSocketTransportConfiguration() {
        resolveMacAddress = false;
    }

    @Override
    public boolean isResolveMacAccess() {
        return resolveMacAddress;
    }

    public void setResolveMacAddress(boolean resolveMacAddress) {
        this.resolveMacAddress = resolveMacAddress;
    }

}
