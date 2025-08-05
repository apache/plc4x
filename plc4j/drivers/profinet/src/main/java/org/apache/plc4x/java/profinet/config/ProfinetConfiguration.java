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
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.profinet.device.GsdFileMap;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.spi.configuration.ConfigurationParameterConverter;
import org.apache.plc4x.java.spi.configuration.annotations.*;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProfinetConfiguration implements PlcConnectionConfiguration {

    @Required
    @ConfigurationParameter
    @ParameterConverter(ProfinetDeviceConvertor.class)
    @Description("Allows you to specify the devices you would like to communicate to, their device access\n" +
        "module (Taken from the GSD file) as well as a list of submodules.\n" +
        "\n" +
        "This parameter has the format\n" +
        "\n" +
        "----\n" +
        "[[{device-1},{device-access},({submodule-1},{submodule-2})],[{device-2},{device-access},({submodule-1},{submodule-2})],....]\n" +
        "----\n" +
        "\n" +
        "For each available slot specified in the GSD file a submodule needs to be in the connection string, however it can be left blank e.g.\n" +
        "\n" +
        "----\n" +
        "[[{device},{device-access},({submodule-1},)]]\n" +
        "----\n" +
        "\n" +
        "If there is no submodule configured.")
    protected ProfinetDevices devices;

    @Required
    @ConfigurationParameter("gsddirectory")
    @ParameterConverter(ProfinetGsdFileConvertor.class)
    @Description("The directory that is used to store any GSD files. This is used to look up the GSD for device found.")
    protected static GsdFileMap gsdFiles;

    @ConfigurationParameter("sendclockfactor")
    @IntDefaultValue(32)
    @Description("This is used to scale the frequency in which cyclic packets are sent. Increasing this slows down communication.")
    private int sendClockFactor;

    @ConfigurationParameter("reductionratio")
    @IntDefaultValue(4)
    @Description("Is also used to scale the frequency. The formula to calculate the overall cycle time is Cycle Time = SendClockFactor * Reduction Ratio * 31.23us")
    private int reductionRatio;

    @ConfigurationParameter("watchdogfactor")
    @IntDefaultValue(50)
    @Description("Used to specify the maximum number of cycles that is allowed to be missed by a device. An alarm is generated if this is exceeded")
    private int watchdogFactor;

    @ConfigurationParameter("dataholdfactor")
    @IntDefaultValue(50)
    @Description("Specifies the number of cycles a device will keep its outputs in a non-safe state when it hasn't received a cyclic packet. This must be equal to or be greater than the watchdog factor")
    private int dataHoldFactor;

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
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
