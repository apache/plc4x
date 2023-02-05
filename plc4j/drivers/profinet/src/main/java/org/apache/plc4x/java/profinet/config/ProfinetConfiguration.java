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
import org.apache.plc4x.java.profinet.device.GsdFileMap;
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.device.ProfinetDevices;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.ParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

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

    @Required
    @ConfigurationParameter
    @ParameterConverter(ProfinetDeviceConvertor.class)
    protected ProfinetDevices devices;

    @Required
    @ConfigurationParameter("gsddirectory")
    @ParameterConverter(ProfinetGsdFileConvertor.class)
    static protected GsdFileMap gsdFiles;

    @ConfigurationParameter("sendclockfactor")
    @IntDefaultValue(32)
    private int sendClockFactor;

    @ConfigurationParameter("reductionratio")
    @IntDefaultValue(4)
    private int reductionRatio;

    @ConfigurationParameter("watchdogfactor")
    @IntDefaultValue(50)
    private int watchdogFactor;

    @ConfigurationParameter("dataholdfactor")
    @IntDefaultValue(50)
    private int dataHoldFactor;

    public static class ProfinetDeviceConvertor implements ConfigurationParameterConverter<ProfinetDevices> {

        public static final Pattern DEVICE_NAME_ARRAY_PATTERN = Pattern.compile("^\\[(?:(\\[(?:[\\w-]*){1},(?:[\\w]*){1},\\((?:[\\w]*[, ]?)*\\){1}\\])[, ]?)+\\]");
        public static final Pattern DEVICE_PARAMETERS = Pattern.compile("^(?<devicename>[\\w-]*){1}[, ]+(?<deviceaccess>[\\w]*){1}[, ]+\\((?<submodules>[\\w, ]*)\\)");

        @Override
        public Class<ProfinetDevices> getType() {
            return ProfinetDevices.class;
        }

        @Override
        public ProfinetDevices convert(String value) {

            // Split up the connection string into its individual segments.
            value = value.replaceAll(" ", "").toUpperCase();
            Matcher matcher = DEVICE_NAME_ARRAY_PATTERN.matcher(value);

            if (!matcher.matches()) {
                throw new RuntimeException("Profinet Device Array is not in the correct format " + value + ".");
            }

            Map<String, ProfinetDevice> devices = new HashMap<>();
            String[] deviceParameters  = value.substring(1, value.length() - 1).replaceAll(" ", "").split("[\\[\\]]");
            for (String deviceParameter : deviceParameters) {
                if (deviceParameter.length() > 7) {
                    matcher = DEVICE_PARAMETERS.matcher(deviceParameter);
                    if (matcher.matches()) {
                        devices.put(matcher.group("devicename"),
                            new ProfinetDevice(matcher.group("devicename"),
                                               matcher.group("deviceaccess"),
                                               matcher.group("submodules"),
                                               (vendorId, deviceId) -> gsdFiles.getGsdFiles().get("0x" + vendorId + "-0x" + deviceId)
                            )
                        );
                    }
                }
            }

            return new ProfinetDevices(devices);
        }
    }

    public static class ProfinetGsdFileConvertor implements ConfigurationParameterConverter<GsdFileMap> {

        @Override
        public Class<GsdFileMap> getType() {
            return GsdFileMap.class;
        }

        @Override
        public GsdFileMap convert(String value) {
            HashMap<String, ProfinetISO15745Profile> gsdFiles = new HashMap<>();
            try {
                DirectoryStream<Path> stream = Files.newDirectoryStream(Paths.get(value));
                XmlMapper xmlMapper = new XmlMapper();
                for (Path file : stream) {
                    try {
                        ProfinetISO15745Profile gsdFile = xmlMapper.readValue(file.toFile(), ProfinetISO15745Profile.class);
                        if (gsdFile.getProfileHeader() != null && gsdFile.getProfileHeader().getProfileIdentification().equals("PROFINET Device Profile") && gsdFile.getProfileHeader().getProfileClassID().equals("Device")) {
                            String id = gsdFile.getProfileBody().getDeviceIdentity().getVendorId() + "-" + gsdFile.getProfileBody().getDeviceIdentity().getDeviceID();
                            gsdFiles.put(id, gsdFile);
                        }
                    } catch (IOException ignored) {
                    }
                }
            } catch (IOException e) {
                throw new RuntimeException("GSDML File directory is un-readable");
            }
            return new GsdFileMap(gsdFiles);
        }

    }

    public ProfinetDevices getDevices() {
        return this.devices;
    }

    public int getSendClockFactor() {
        return sendClockFactor;
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

    public static GsdFileMap getGsdFiles() {
        return gsdFiles;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
