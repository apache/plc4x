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
package org.apache.plc4x.java.profinet.config;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.BooleanDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

public class ProfinetConfiguration implements Configuration, RawSocketTransportConfiguration {

    @Override
    public boolean getSupportVlans() {
        return RawSocketTransportConfiguration.super.getSupportVlans();
    }

    @Override
    public int getDefaultPort() {
        return 34964;
    }

    @Override
    public Integer getProtocolId() {
        return RawSocketTransportConfiguration.super.getProtocolId();
    }

    @Override
    public PacketHandler getPcapPacketHandler() {
        return null;
    }

    @ConfigurationParameter("devices")
    @StringDefaultValue("")
    private String devices;

    public HashMap<MacAddress, ProfinetDevice> configuredDevices = new HashMap<>();

    public void setDevices(String sDevices) throws DecoderException {
        // TODO:- Add support for passing in configured devices.
        MacAddress macAddress = new MacAddress(Hex.decodeHex("005056c00001"));
        configuredDevices.put(macAddress, new ProfinetDevice(macAddress));
    }

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
