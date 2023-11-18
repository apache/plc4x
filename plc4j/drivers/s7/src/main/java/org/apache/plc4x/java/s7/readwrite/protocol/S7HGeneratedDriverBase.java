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
package org.apache.plc4x.java.s7.readwrite.protocol;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.configuration.S7TcpTransportConfiguration;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.connection.ProtocolStackConfigurer;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationTypeProvider;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ServiceLoader;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

public class S7HGeneratedDriverBase extends GeneratedDriverBase<TPKTPacket> implements TransportConfigurationTypeProvider {

    private static final Logger logger = LoggerFactory.getLogger(S7HGeneratedDriverBase.class);

    private static final Pattern URI_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9]*))?://(?<transportConfig>[^?^/]*)(\\?(?<paramString>.*))?");

    private static final Pattern URI_H_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9]*))?://(?<transportConfig>[^?]*)/(?<transportConfig2>[^?]*)\\?(?<paramString>.*)?");

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        // Split up the connection string into it's individual segments.
        Matcher smatcher = URI_PATTERN.matcher(connectionString);
        Matcher hmatcher = URI_H_PATTERN.matcher(connectionString);
        if (!smatcher.matches() && !hmatcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-address}(?{parameter-string)?'");
        }

        Matcher matcher = (smatcher.matches()) ? smatcher : hmatcher;

        final String protocolCode = matcher.group("protocolCode");
        final String transportCode = (matcher.group("transportCode") != null) ?
            matcher.group("transportCode") : getDefaultTransport();
        final String transportConfig = matcher.group("transportConfig");
        final String transportConfig2 = (hmatcher.matches()) ? matcher.group("transportConfig2") : null;
        final String paramString = matcher.group("paramString");

        // Check if the protocol code matches this driver.
        if (!protocolCode.equals(getProtocolCode())) {
            // Actually this shouldn't happen as the DriverManager should have not used this driver in the first place.
            throw new PlcConnectionException(
                "This driver is not suited to handle this connection string");
        }

        // Create the configuration object.
        Configuration configuration = new ConfigurationFactory().createConfiguration(
            getConfigurationType(), protocolCode, transportCode, transportConfig, paramString);
        if (configuration == null) {
            throw new PlcConnectionException("Unsupported configuration");
        }

        // Try to find a transport in order to create a communication channel.
        Transport transport = null;
        ServiceLoader<Transport> transportLoader = ServiceLoader.load(
            Transport.class, Thread.currentThread().getContextClassLoader());
        for (Transport curTransport : transportLoader) {
            if (curTransport.getTransportCode().equals(transportCode)) {
                transport = curTransport;
                break;
            }
        }
        if (transport == null) {
            throw new PlcConnectionException("Unsupported transport " + transportCode);
        }

        // Inject the configuration into the transport.
        configure(configuration, transport);

        // Create an instance of the communication channel which the driver should use.
        ChannelFactory channelFactory = transport.createChannelFactory(transportConfig);
        if (channelFactory == null) {
            throw new PlcConnectionException("Unable to get channel factory from url " + transportConfig);
        }
        configure(configuration, channelFactory);

        // Create an instance of the communication channel which the driver should use.
        ChannelFactory secondaryChannelFactory = null;
        if (hmatcher.matches()) {
            secondaryChannelFactory = transport.createChannelFactory(transportConfig2);
            if (secondaryChannelFactory == null) {
                logger.info("Unable to get channel factory from url " + transportConfig2);
            }
        }

        if (hmatcher.matches() && (secondaryChannelFactory != null)) {
            configure(configuration, secondaryChannelFactory);
        }

        // Give drivers the option to customize the channel.
        initializePipeline(channelFactory);

        // Give drivers the option to customize the channel.
        if (hmatcher.matches())
            initializePipeline(secondaryChannelFactory);

        // Make the "await setup complete" overridable via system property.
        boolean awaitSetupComplete = awaitSetupComplete();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE) != null) {
            awaitSetupComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE));
        }

        // Make the "await disconnect complete" overridable via system property.
        boolean awaitDisconnectComplete = awaitDisconnectComplete();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE) != null) {
            awaitDisconnectComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE));
        }

        // Make the "await disconnect complete" overridable via system property.
        boolean awaitDiscoverComplete = awaitDiscoverComplete();
        if (System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE) != null) {
            awaitDiscoverComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE));
        }

        return new S7HPlcConnection(
            canPing(),
            canRead(), canWrite(), canSubscribe(), canBrowse(),
            getTagHandler(),
            getValueHandler(),
            configuration,
            channelFactory,
            secondaryChannelFactory,
            false,
            awaitSetupComplete,
            awaitDisconnectComplete,
            awaitDiscoverComplete,
            getStackConfigurer(transport),
            getOptimizer(),
            getAuthentication());
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        throw new UnsupportedOperationException("getConfigurationType, Not supported yet.");
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        throw new UnsupportedOperationException("getTagHandler, Not supported yet.");
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        throw new UnsupportedOperationException("getValueHandler, Not supported yet.");
    }

    @Override
    protected String getDefaultTransport() {
        throw new UnsupportedOperationException("getDefaultTransport, Not supported yet.");
    }

    @Override
    protected ProtocolStackConfigurer<TPKTPacket> getStackConfigurer() {
        throw new UnsupportedOperationException("getStackConfigurer, Not supported yet.");
    }

    @Override
    public String getProtocolCode() {
        throw new UnsupportedOperationException("getProtocolCode, Not supported yet.");
    }

    @Override
    public String getProtocolName() {
        throw new UnsupportedOperationException("getProtocolName, Not supported yet.");
    }

    public PlcAuthentication getAuthentication() {
        return null;
    }

    @Override
    public Class<? extends TransportConfiguration> getTransportConfigurationType(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return S7TcpTransportConfiguration.class;
        }
        return null;
    }

}
