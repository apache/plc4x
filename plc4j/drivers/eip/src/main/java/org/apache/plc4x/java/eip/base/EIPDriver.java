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
package org.apache.plc4x.java.eip.base;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.configuration.PlcConnectionConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.eip.base.configuration.EIPConfiguration;
import org.apache.plc4x.java.eip.base.configuration.EipTcpTransportConfiguration;
import org.apache.plc4x.java.eip.base.tag.EipTag;
import org.apache.plc4x.java.eip.base.protocol.EipProtocolLogic;
import org.apache.plc4x.java.eip.base.tag.EipTagHandler;
import org.apache.plc4x.java.eip.readwrite.EipPacket;
import org.apache.plc4x.java.spi.configuration.ConfigurationFactory;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.*;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.transport.Transport;

import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.transport.TransportConfiguration;
import org.apache.plc4x.java.spi.transport.TransportConfigurationTypeProvider;

import java.util.ServiceLoader;
import java.util.function.Consumer;
import java.util.function.ToIntFunction;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

public class EIPDriver extends GeneratedDriverBase<EipPacket> implements TransportConfigurationTypeProvider {
    public static final int PORT = 44818;
    private static final Pattern URI_PATTERN = Pattern.compile(
        "^(?<protocolCode>[a-z0-9\\-]*)(:(?<transportCode>[a-z0-9]*))?://(?<transportConfig>[^?]*)(\\?(?<paramString>.*))?");

    private EIPConfiguration configuration;

    @Override
    public String getProtocolCode() {
        return "eip";
    }

    @Override
    public String getProtocolName() {
        return "EthernetIP";
    }

    @Override
    public Class<? extends PlcConnectionConfiguration> getConfigurationType() {
        return EIPConfiguration.class;
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new EipTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new org.apache.plc4x.java.spi.values.PlcValueHandler();
    }

    /**
     * This protocol doesn't have a disconnect procedure, so there is no need to wait for a login to finish.
     * @return false
     */
    @Override
    protected boolean awaitDisconnectComplete() {
        return true;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
    }

    @Override
    protected boolean canRead() {
        return true;
    }

    @Override
    protected boolean canWrite() {
        return true;
    }

    @Override
    protected ProtocolStackConfigurer<EipPacket> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(EipPacket.class, EipPacket::staticParse)
            .withProtocol(EipProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            .withParserArgs(true)
            .withCorruptPacketRemover(CorruptPackageCleaner.class)
            .byteOrder(this.configuration.getByteOrder())
            .build();
    }

    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        // Split up the connection string into its individual segments.
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
        this.configuration = (EIPConfiguration) new ConfigurationFactory().createConfiguration(
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
        if(channelFactory == null) {
            throw new PlcConnectionException("Unable to get channel factory from url " + transportConfig);
        }
        configure(configuration, channelFactory);

        // Give drivers the option to customize the channel.
        initializePipeline(channelFactory);

        // Make the "fire discover event" overridable via system property.
        boolean fireDiscoverEvent = fireDiscoverEvent();
        if(System.getProperty(PROPERTY_PLC4X_FORCE_FIRE_DISCOVER_EVENT) != null) {
            fireDiscoverEvent = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_FIRE_DISCOVER_EVENT));
        }

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
            canPing(), canRead(), canWrite(), canSubscribe(), canBrowse(),
            getTagHandler(),
            getValueHandler(),
            configuration,
            channelFactory,
            fireDiscoverEvent,
            awaitSetupComplete,
            awaitDisconnectComplete,
            awaitDiscoverComplete,
            getStackConfigurer(transport),
            getOptimizer(),
            null);
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf>, HasConfiguration<EIPConfiguration> {
        private EIPConfiguration configuration;

        @Override
        public void setConfiguration(EIPConfiguration configuration) {
            this.configuration = configuration;
        }

        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            if (byteBuf.readableBytes() >= 4) {
                //Second byte for the size and then add the header size 24
                if (configuration.getByteOrder() == ByteOrder.BIG_ENDIAN) {
                    return byteBuf.getUnsignedShort(byteBuf.readerIndex() + 2) + 24;
                } else {
                    return byteBuf.getUnsignedShortLE(byteBuf.readerIndex() + 2) + 24;
                }
            }
            return -1;
        }
    }

     /**Consumes all Bytes till another Magic Byte is found */
    public static class CorruptPackageCleaner implements Consumer<ByteBuf> {
        @Override
        public void accept(ByteBuf byteBuf) {
            while (byteBuf.getUnsignedByte(0) != 0x00) {
                // Just consume the bytes till the next possible start position.
                byteBuf.readByte();
            }
        }
    }

    @Override
    public EipTag prepareTag(String query){
        return EipTag.of(query);
    }

    @Override
    public Class<? extends TransportConfiguration> getTransportConfigurationType(String transportCode) {
        switch (transportCode) {
            case "tcp":
                return EipTcpTransportConfiguration.class;
        }
        return null;
    }

}
