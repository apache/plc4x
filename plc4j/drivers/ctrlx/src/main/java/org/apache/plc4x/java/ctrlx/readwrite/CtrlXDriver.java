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

package org.apache.plc4x.java.ctrlx.readwrite;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.authentication.PlcUsernamePasswordAuthentication;
import org.apache.plc4x.java.spi.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcDiscoveryRequest;
import org.apache.plc4x.java.ctrlx.readwrite.configuration.CtrlXConfiguration;
import org.apache.plc4x.java.ctrlx.readwrite.connection.CtrlXConnection;
import org.apache.plc4x.java.ctrlx.readwrite.discovery.CtrlXPlcDiscoverer;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.messages.DefaultPlcDiscoveryRequest;

import java.util.regex.Matcher;

public class CtrlXDriver implements PlcDriver {

    @Override
    public String getProtocolCode() {
        return "ctrlx";
    }

    @Override
    public String getProtocolName() {
        return "CtrlX";
    }

    @Override
    public PlcConnection getConnection(String url) throws PlcConnectionException {
        throw new PlcConnectionException("CtrlX connections require authentication.");
    }

    @Override
    public PlcConnection getConnection(String connectionString, PlcAuthentication authentication) throws PlcConnectionException {
        ConfigurationFactory configurationFactory = new ConfigurationFactory();
        // Split up the connection string into its individual segments.
        Matcher matcher = GeneratedDriverBase.URI_PATTERN.matcher(connectionString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-config}(?{parameter-string)?'");
        }
        final String protocolCode = matcher.group("protocolCode");
        String transportCodeMatch = matcher.group("transportCode");
        final String transportCode = (transportCodeMatch != null) ? transportCodeMatch : "https";
        final String transportConfig = matcher.group("transportConfig");
        final String paramString = matcher.group("paramString");

        // Check if the protocol code matches this driver.
        if (!protocolCode.equals(getProtocolCode())) {
            // Actually this shouldn't happen as the DriverManager should have not used this driver in the first place.
            throw new PlcConnectionException(
                "This driver is not suited to handle this connection string");
        }

        // Create the configuration object.
        PlcConnectionConfiguration configuration = configurationFactory
            .createConfiguration(CtrlXConfiguration.class, protocolCode, transportCode, transportConfig, paramString);
        if (configuration == null) {
            throw new PlcConnectionException("Unsupported configuration");
        }

        // CtrlX only supports "https" as transport.
        if(!"https".equals(transportCode)) {
            throw new PlcConnectionException("Only 'https' transport is supported by this driver");
        }

        if((!(authentication instanceof PlcUsernamePasswordAuthentication))) {
            throw new PlcConnectionException("CtrlX connections require username-password authentication");
        }
        PlcUsernamePasswordAuthentication usernamePasswordAuthentication =
            (PlcUsernamePasswordAuthentication) authentication;

        return new CtrlXConnection(
            String.format("%s://%s", transportCode, transportConfig),
            usernamePasswordAuthentication.getUsername(),
            usernamePasswordAuthentication.getPassword());
    }

    @Override
    public PlcDiscoveryRequest.Builder discoveryRequestBuilder() {
        return new DefaultPlcDiscoveryRequest.Builder(new CtrlXPlcDiscoverer());
    }

}
