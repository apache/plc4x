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
package org.apache.plc4x.java.knxnetip.configuration;

import org.apache.plc4x.java.knxnetip.readwrite.KnxLayer;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.annotations.ConfigurationParameter;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.IntDefaultValue;
import org.apache.plc4x.java.spi.configuration.annotations.defaults.StringDefaultValue;
import org.apache.plc4x.java.spi.configuration.exceptions.ConfigurationException;

public class KnxNetIpConfiguration implements Configuration {

    @ConfigurationParameter("knxproj-file-path")
    public String knxprojFilePath;

    @ConfigurationParameter("knxproj-password")
    public String knxprojPassword;

    @ConfigurationParameter("group-address-num-levels")
    @IntDefaultValue(3)
    public int groupAddressNumLevels = 3;

    @ConfigurationParameter("connection-type")
    @StringDefaultValue("LINK_LAYER")
    public String connectionType = "LINK_LAYER";

    public String getKnxprojFilePath() {
        return knxprojFilePath;
    }

    public void setKnxprojFilePath(String knxprojFilePath) {
        this.knxprojFilePath = knxprojFilePath;
    }

    public String getKnxprojPassword() {
        return knxprojPassword;
    }

    public void setKnxprojPassword(String knxprojPassword) {
        this.knxprojPassword = knxprojPassword;
    }

    public int getGroupAddressNumLevels() {
        return groupAddressNumLevels;
    }

    public void setGroupAddressNumLevels(int groupAddressNumLevels) {
        this.groupAddressNumLevels = groupAddressNumLevels;
    }

    public String getConnectionType() {
        return connectionType;
    }

    public void setConnectionType(String connectionType) {
        // Try to parse the provided value, if it doesn't match any of the constants,
        // throw an error.
        try {
            KnxLayer.valueOf("TUNNEL_" + connectionType.toUpperCase());
        } catch (IllegalArgumentException e) {
            throw new ConfigurationException("Value provided for connection-type invalid.");
        }
        this.connectionType = connectionType.toUpperCase();
    }

    @Override
    public String toString() {
        return "Configuration{" +
            "knxprojFilePath=" + knxprojFilePath + ", " +
            "groupAddressNumLevels=" + groupAddressNumLevels +
            '}';
    }

}
