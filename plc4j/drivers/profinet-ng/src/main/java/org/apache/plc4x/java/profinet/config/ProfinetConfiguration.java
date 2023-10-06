/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.profinet.gsdml.ProfinetISO15745Profile;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.Required;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.transport.rawsocket.RawSocketTransportConfiguration;
import org.apache.plc4x.java.utils.pcap.netty.handlers.PacketHandler;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;

public class ProfinetConfiguration implements Configuration, RawSocketTransportConfiguration {

    @Required
    @ConfigurationParameter("gsd-directory")
    @StringDefaultValue("~/.gsd")
    public String gsdDirectory;

    @ConfigurationParameter("dap-id")
    public String dapId;

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

    @Override
    public boolean isResolveMacAccess() {
        return true;
    }

    public ProfinetISO15745Profile getGsdProfile(int vendorId, int deviceId) {
        String value = gsdDirectory;
        if(value.startsWith("~")) {
            String homeDirectory = System.getProperty("user.home");
            value = value.replaceAll("~", homeDirectory);
        }

        DirectoryStream<Path> stream = null;
        try {
            stream = Files.newDirectoryStream(Paths.get(value));
            XmlMapper xmlMapper = new XmlMapper();
            for (Path file : stream) {
                try {
                    ProfinetISO15745Profile gsdFile = xmlMapper.readValue(file.toFile(), ProfinetISO15745Profile.class);
                    if (gsdFile.getProfileHeader() != null && gsdFile.getProfileHeader().getProfileIdentification().equals("PROFINET Device Profile") && gsdFile.getProfileHeader().getProfileClassID().equals("Device")) {
                        int curVendorId = Integer.parseInt(gsdFile.getProfileBody().getDeviceIdentity().getVendorId().substring(2), 16);
                        int curDeviceId = Integer.parseInt(gsdFile.getProfileBody().getDeviceIdentity().getDeviceID().substring(2), 16);
                        // Check, if this is the device we're looking for.
                        if ((curVendorId == vendorId) && (curDeviceId == deviceId)) {
                            return gsdFile;
                        }
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
        return null;
    }

}
