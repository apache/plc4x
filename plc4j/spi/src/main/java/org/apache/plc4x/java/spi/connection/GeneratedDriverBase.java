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
package org.apache.plc4x.java.spi.connection;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.*;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.api.value.PlcValueHandler;

import java.util.ServiceLoader;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public abstract class GeneratedDriverBase<BASE_PACKET extends Message> implements PlcDriver {

    public static final String PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE = "PLC4X_FORCE_AWAIT_SETUP_COMPLETE";
    public static final String PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE = "PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE";
    public static final String PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE = "PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE";

    private static final Pattern URI_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9]*))?://(?<transportConfig>[^?]*)(\\?(?<paramString>.*))?");

    protected abstract Class<? extends Configuration> getConfigurationType();

    protected boolean canRead() {
        return false;
    }

    protected boolean canWrite() {
        return false;
    }

    protected boolean canSubscribe() {
        return false;
    }

    protected boolean awaitSetupComplete() {
        return true;
    }

    protected boolean awaitDisconnectComplete() {
        return false;
    }

    protected boolean awaitDiscoverComplete() {
        return false;
    }

    protected BaseOptimizer getOptimizer() {
        return null;
    }

    protected abstract PlcFieldHandler getFieldHandler();

    protected abstract PlcValueHandler getValueHandler();

    protected abstract String getDefaultTransport();

    protected abstract ProtocolStackConfigurer<BASE_PACKET> getStackConfigurer();

    protected ProtocolStackConfigurer<BASE_PACKET> getStackConfigurer(Transport transport) {
        return getStackConfigurer();
    }

    protected void initializePipeline(ChannelFactory channelFactory) {
        // Override in derived drivers.
    }

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        // Split up the connection string into it's individual segments.
        Matcher matcher = URI_PATTERN.matcher(connectionString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-address}(?{parameter-string)?'");
        }
        final String protocolCode = matcher.group("protocolCode");
        final String transportCode = (matcher.group("transportCode") != null) ?
            matcher.group("transportCode") : getDefaultTransport();
        final String transportConfig = matcher.group("transportConfig");
        final String paramString = matcher.group("paramString");

        // Check if the protocol code matches this driver.
        if (!protocolCode.equals(getProtocolCode())) {
            // Actually this shouldn't happen as the DriverManager should have not used this driver in the first place.
            throw new PlcConnectionException(
                "This driver is not suited to handle this connection string");
        }

        // Create the configuration object.
        Configuration configuration = new ConfigurationFactory().createConfiguration(
            getConfigurationType(), paramString);
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
        if(channelFactory == null) {
            throw new PlcConnectionException("Unable to get channel factory from url " + transportConfig);
        }
        configure(configuration, channelFactory);

        // Give drivers the option to customize the channel.
        initializePipeline(channelFactory);

        // Make the "await setup complete" overridable via system property.
        boolean awaitSetupComplete = awaitSetupComplete();
        if(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE) != null) {
            awaitSetupComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE));
        }

        // Make the "await disconnect complete" overridable via system property.
        boolean awaitDisconnectComplete = awaitDisconnectComplete();
        if(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE) != null) {
            awaitDisconnectComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE));
        }

        // Make the "await disconnect complete" overridable via system property.
        boolean awaitDiscoverComplete = awaitDiscoverComplete();
        if(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE) != null) {
            awaitDiscoverComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_DISCOVER_COMPLETE));
        }

        return new DefaultNettyPlcConnection(
            canRead(), canWrite(), canSubscribe(),
            getFieldHandler(),
            getValueHandler(),
            configuration,
            channelFactory,
            awaitSetupComplete,
            awaitDisconnectComplete,
            awaitDiscoverComplete,
            getStackConfigurer(transport),
            getOptimizer());
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Authentication not supported.");
    }


}
