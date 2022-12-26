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

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.BaseConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetConfiguration extends BaseConfiguration implements RawSocketTransportConfiguration {

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

    @ConfigurationParameter("deviceaccess")
    @StringDefaultValue("")
    private String deviceAccess;

    @ConfigurationParameter("devices")
    @StringDefaultValue("")
    private String devices;

    @ConfigurationParameter("gsddirectory")
    @StringDefaultValue("")
    private String gsdDirectory;

    @ConfigurationParameter("sendclockfactor")
    @IntDefaultValue(32)
    private int sendClockFactor;

    @ConfigurationParameter("submodules")
    @StringDefaultValue("")
    private String subModules;

    @ConfigurationParameter("reductionratio")
    @IntDefaultValue(4)
    private int reductionRatio;

    @ConfigurationParameter("watchdogfactor")
    @IntDefaultValue(50)
    private int watchdogFactor;

    @ConfigurationParameter("dataholdfactor")
    @IntDefaultValue(50)
    private int dataHoldFactor;

    public String getDevices() {
        return devices;
    }

    public String getGsdDirectory() {
        return gsdDirectory;
    }

    public void setGsdDirectory(String gsdDirectory) {
        this.gsdDirectory = gsdDirectory;
    }

    public int getSendClockFactor() {
        return sendClockFactor;
    }

    public String getSubModules() {
        return subModules;
    }

    public int getReductionRatio() {
        return reductionRatio;
    }

    public int getWatchdogFactor() {
        return watchdogFactor;
    }

    public int getDataHoldFactor() {
        return dataHoldFactor;
    }

    public InetAddress getIpAddress() {
        try {
            return InetAddress.getByName(getTransportConfig().split(":")[0]);
        } catch (UnknownHostException e) {
            throw new RuntimeException(e);
        }
    }

    public String getDeviceAccess() {
        return deviceAccess;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
