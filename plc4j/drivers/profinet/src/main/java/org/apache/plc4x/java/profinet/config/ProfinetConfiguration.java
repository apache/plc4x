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
import org.apache.plc4x.java.profinet.device.ProfinetDevice;
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.profinet.readwrite.MacAddress;
import org.apache.plc4x.java.spi.configuration.BaseConfiguration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
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

    private final Logger logger = LoggerFactory.getLogger(ProfinetConfiguration.class);
    public static final Pattern MACADDRESS_ARRAY_PATTERN = Pattern.compile("^\\[(([A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2}:[A-F0-9]{2})(,)?)*\\]");
    public static final Pattern SUB_MODULE_ARRAY_PATTERN = Pattern.compile("(\\[[\\w, ]*\\])");

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

    public LinkedHashMap<String, ProfinetDevice> configuredDevices = new LinkedHashMap<>();

    private final Map<String, ProfinetISO15745Profile> gsdFiles = new HashMap<>();

    public String getDevices() {
        return devices;
    }

    public void setDevices(String sDevices) throws DecoderException, PlcConnectionException {

        // Split up the connection string into its individual segments.
        Matcher matcher = MACADDRESS_ARRAY_PATTERN.matcher(sDevices.toUpperCase());

        if (!matcher.matches()) {
            throw new PlcConnectionException("Profinet Device Array is not in the correct format " + sDevices + ".");
        }

        String[] devices = sDevices.substring(1, sDevices.length() - 1).split("[ ,]");

        for (String device : devices) {
            MacAddress macAddress = new MacAddress(Hex.decodeHex(device.replace(":", "")));
            configuredDevices.put(device.replace(":", "").toUpperCase(), new ProfinetDevice(macAddress, this));
        }
    }

    public void setSubModules() throws DecoderException, PlcConnectionException {

        // Split up the connection string into its individual segments.
        Matcher matcher = SUB_MODULE_ARRAY_PATTERN.matcher(subModules.toUpperCase());
        if (!matcher.matches()) {
            throw new PlcConnectionException("Profinet Submodule Array is not in the correct format " + subModules + ".");
        }
        String[] devices = new String[matcher.groupCount()];
        for (int i = 0; i < matcher.groupCount(); i++) {
            devices[i] = matcher.group(i).replace(" ", "");
        }

        if (matcher.groupCount() != configuredDevices.size()) {
            throw new PlcConnectionException("Configured device array size doesn't match the submodule array size");
        }

        int index = 0;
        for (Map.Entry<String, ProfinetDevice> entry : configuredDevices.entrySet()) {
            entry.getValue().setSubModules(devices[index]);
            index += 1;
        }
    }

    public Map<String, ProfinetISO15745Profile> readGsdFiles() {
        try {
            DirectoryStream<Path> stream = Files.newDirectoryStream(Paths.get(this.gsdDirectory));
            XmlMapper xmlMapper = new XmlMapper();
            for (Path file : stream) {
                try {
                    ProfinetISO15745Profile gsdFile = xmlMapper.readValue(file.toFile(), ProfinetISO15745Profile.class);
                    if (gsdFile.getProfileHeader() != null && gsdFile.getProfileHeader().getProfileIdentification().equals("PROFINET Device Profile") && gsdFile.getProfileHeader().getProfileClassID().equals("Device")) {
                        String id = gsdFile.getProfileBody().getDeviceIdentity().getVendorId() + "-" + gsdFile.getProfileBody().getDeviceIdentity().getDeviceID();
                        logger.debug("Adding GSDML file for {}", gsdFile.getProfileBody().getDeviceIdentity().getVendorName().getValue());
                        this.gsdFiles.put(id, gsdFile);
                    }
                } catch (IOException e) {
                    // Pass - Ignore any files that aren't xml files.
                    logger.debug(String.valueOf(e));
                }
            }
        } catch (IOException e) {
            throw new RuntimeException("GSDML File directory is un-readable");
        }
        return this.gsdFiles;
    }

    public Map<String, ProfinetISO15745Profile> getGsdFiles() {
        return gsdFiles;
    }

    public HashMap<String, ProfinetDevice> getConfiguredDevices() {
        return configuredDevices;
    }

    public void setConfiguredDevices(LinkedHashMap<String, ProfinetDevice> configuredDevices) {
        this.configuredDevices = configuredDevices;
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

    @Override
    public String toString() {
        return "Configuration{" +
            '}';
    }

}
