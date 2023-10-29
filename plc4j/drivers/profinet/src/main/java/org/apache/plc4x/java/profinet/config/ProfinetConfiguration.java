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
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.*;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationProvider;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetConfiguration implements Configuration, TransportConfigurationProvider {

    @Required
    @ConfigurationParameter
    @ParameterConverter(ProfinetDeviceConvertor.class)
    protected ProfinetDevices devices;

    @Required
    @ConfigurationParameter("gsddirectory")
    @ParameterConverter(ProfinetGsdFileConvertor.class)
    protected static GsdFileMap gsdFiles;

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

    @ComplexConfigurationParameter(prefix = "raw", defaultOverrides = {
        @ComplexConfigurationParameterDefaultOverride(name = "resolve-mac-address", value = "true")
    }, requiredOverrides = {})
    private ProfinetRawSocketTransportConfiguration rawSocketTransportConfiguration;

    public ProfinetRawSocketTransportConfiguration getRawSocketTransportConfiguration() {
        return rawSocketTransportConfiguration;
    }

    public void setRawSocketTransportConfiguration(ProfinetRawSocketTransportConfiguration rawSocketTransportConfiguration) {
        this.rawSocketTransportConfiguration = rawSocketTransportConfiguration;
    }

    public static class ProfinetDeviceConvertor implements ConfigurationParameterConverter<ProfinetDevices> {

        public static final String DEVICE_STRING = "((?<devicename>[\\w- ]+)[, ]+(?<deviceaccess>[\\w ]+)[, ]+\\((?<submodules>[\\w, ]*)\\)[, ]*(?<ipaddress>\\d+\\.\\d+\\.\\d+\\.\\d+)?)";
        public static final String DEVICE_ARRAY_STRING = "^\\[(?:(\\[" + DEVICE_STRING + "{1}\\])[, ]?)+\\]";
        public static final Pattern DEVICE_NAME_ARRAY_PATTERN = Pattern.compile(DEVICE_ARRAY_STRING);
        public static final Pattern DEVICE_PARAMETERS = Pattern.compile(DEVICE_STRING);

        @Override
        public Class<ProfinetDevices> getType() {
            return ProfinetDevices.class;
        }

        @Override
        public ProfinetDevices convert(String value) {

            // Split up the connection string into its individual segments.
            value = value.toUpperCase();
            Matcher matcher = DEVICE_NAME_ARRAY_PATTERN.matcher(value);

            if (!matcher.matches()) {
                throw new RuntimeException("Profinet Device Array is not in the correct format " + value + ".");
            }

            Map<String, ConfigurationProfinetDevice> devices = new HashMap<>();
            String[] deviceParameters  = value.substring(1, value.length() - 1).split("[\\[\\]]");
            for (String deviceParameter : deviceParameters) {
                if (deviceParameter.length() > 7) {
                    matcher = DEVICE_PARAMETERS.matcher(deviceParameter);
                    if (matcher.matches()) {
                        devices.put(matcher.group("devicename"),
                            new ConfigurationProfinetDevice(matcher.group("devicename"),
                                               matcher.group("deviceaccess"),
                                               matcher.group("submodules"),
                                               (vendorId, deviceId) -> gsdFiles.getGsdFiles().get("0x" + vendorId + "-0x" + deviceId)
                            )
                        );
                        if (matcher.group("ipaddress") != null) {
                            devices.get(matcher.group("devicename")).setIpAddress(matcher.group("ipaddress"));
                        }
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
            if(value.startsWith("~")) {
                String homeDirectory = System.getProperty("user.home");
                value = value.replaceAll("~", homeDirectory);
            }
            HashMap<String, ProfinetISO15745Profile> gsdFiles = new HashMap<>();
            DirectoryStream<Path> stream = null;
            try {
                stream = Files.newDirectoryStream(Paths.get(value));
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
            } finally {
                try {
                    if (stream != null) {
                        stream.close();
                    }
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
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
    public TransportConfiguration getTransportConfiguration(String transportCode) {
        switch (transportCode) {
            case "raw":
                return rawSocketTransportConfiguration;
        }
        return null;
    }

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
