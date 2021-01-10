/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/
package org.apache.plc4x.java.opcua;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.opcua.protocol.*;
import org.apache.plc4x.java.opcua.config.*;
import org.apache.plc4x.java.opcua.readwrite.*;
import org.apache.plc4x.java.opcua.readwrite.io.*;
import org.apache.plc4x.java.spi.connection.*;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.spi.values.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;
import org.apache.plc4x.java.spi.optimizer.SingleFieldOptimizer;
import io.netty.buffer.ByteBuf;

import java.util.ServiceLoader;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import java.util.function.ToIntFunction;

import static org.apache.plc4x.java.spi.configuration.ConfigurationFactory.configure;

/**
 * Implementation of the OPC UA protocol, based on:
 * - Eclipse Milo (https://github.com/eclipse/milo)
 *
 * Created by Matthias Milan Strljic on 10.05.2019
 */
public class OpcuaPlcDriver extends GeneratedDriverBase<OpcuaAPU> {


    public static final Pattern OPCUA_URI_PARAM_PATTERN = Pattern.compile("(?<param>[(\\?|\\&)([^=]+)\\=([^&]+)]+)?"); //later used for regex filtering of the params


    public static final Pattern INET_ADDRESS_PATTERN = Pattern.compile("(:(?<transportCode>tcp))?://" +
                                                                        "(?<transportHost>[\\w.-]+)(:" +
                                                                        "(?<transportPort>\\d*))?");


    public static final Pattern URI_PATTERN = Pattern.compile("^(?<protocolCode>opcua)" +
                                                                    INET_ADDRESS_PATTERN +
                                                                    "(?<transportEndpoint>[\\w/=]+)" +
                                                                    "(?<paramString>[(\\?|\\&)\\w=]+\\=[\\w&]+)*"
                                                                );


    private static final int requestTimeout = 10000;

    @Override
    public String getProtocolCode() {
        return "opcua";
    }

    @Override
    public String getProtocolName() {
        return "Opcua";
    }

    @Override
    protected Class<? extends Configuration> getConfigurationType() {
        return OpcuaConfiguration.class;
    }

    @Override
    protected String getDefaultTransport() {
        return "tcp";
    }

    /**
     * Modbus doesn't have a login procedure, so there is no need to wait for a login to finish.
     * @return false
     */
    @Override
    protected boolean awaitSetupComplete() {
        return true;
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
    protected BaseOptimizer getOptimizer() {
        return new SingleFieldOptimizer();
    }

    @Override
    protected OpcuaPlcFieldHandler getFieldHandler() {
        return new OpcuaPlcFieldHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new IEC61131ValueHandler();
    }

    @Override
    protected ProtocolStackConfigurer<OpcuaAPU> getStackConfigurer() {
        return SingleProtocolStackConfigurer.builder(OpcuaAPU.class, OpcuaAPUIO.class)
            .withProtocol(OpcuaProtocolLogic.class)
            .withPacketSizeEstimator(ByteLengthEstimator.class)
            // Every incoming message is to be treated as a response.
            .withParserArgs(true)
            .littleEndian()
            .build();
    }



    @Override
    public PlcConnection getConnection(String connectionString) throws PlcConnectionException {
        // Split up the connection string into it's individual segments.
        Matcher matcher = URI_PATTERN.matcher(connectionString);
        if (!matcher.matches()) {
            throw new PlcConnectionException(
                "Connection string doesn't match the format '{protocol-code}:({transport-code})?//{transport-host}(:{transport-port})(/{transport-endpoint})(?{parameter-string)?'");
        }
        final String protocolCode = matcher.group("protocolCode");
        final String transportCode = (matcher.group("transportCode") != null) ?
            matcher.group("transportCode") : getDefaultTransport();
        final String transportHost = matcher.group("transportHost");
        final String transportPort = matcher.group("transportPort");
        final String transportEndpoint = matcher.group("transportEndpoint");
        final String paramString = matcher.group("paramString");

        System.out.println(protocolCode + " - " + transportCode + " - " + transportHost + " - " +  transportPort + " - " +  transportEndpoint + " - " +  paramString);

        // Check if the protocol code matches this driver.
        if(!protocolCode.equals(getProtocolCode())) {
            // Actually this shouldn't happen as the DriverManager should have not used this driver in the first place.
            throw new PlcConnectionException(
                "This driver is not suited to handle this connection string");
        }

        // Create the configuration object.
        OpcuaConfiguration configuration = new OpcuaConfiguration(  transportCode,
                                                                    transportHost,
                                                                    transportPort,
                                                                    transportEndpoint,
                                                                    paramString);

        if(configuration == null) {
            throw new PlcConnectionException("Unsupported configuration");
        }

        // Try to find a transport in order to create a communication channel.
        Transport transport = null;
        ServiceLoader<Transport> transportLoader = ServiceLoader.load(
            Transport.class, Thread.currentThread().getContextClassLoader());
        for (Transport curTransport : transportLoader) {
            if(curTransport.getTransportCode().equals(transportCode)) {
                transport = curTransport;
                break;
            }
        }
        if(transport == null) {
            throw new PlcConnectionException("Unsupported transport " + transportCode);
        }

        // Inject the configuration into the transport.
        configure(configuration, transport);

        // Create an instance of the communication channel which the driver should use.
        ChannelFactory channelFactory = transport.createChannelFactory(transportHost + ":" + transportPort);
        if(channelFactory == null) {
            throw new PlcConnectionException("Unable to get channel factory from url " + transportHost + ":" + transportPort);
        }
        configure(configuration, channelFactory);

        // Give drivers the option to customize the channel.
        initializePipeline(channelFactory);

        // Make the "await setup complete" overridable via system property.
        boolean awaitSetupComplete = awaitSetupComplete();
        if(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE) != null) {
            awaitSetupComplete = Boolean.parseBoolean(System.getProperty(PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE));
        }

        return new DefaultNettyPlcConnection(
            canRead(), canWrite(), canSubscribe(),
            getFieldHandler(),
            getValueHandler(),
            configuration,
            channelFactory,
            awaitSetupComplete,
            getStackConfigurer(),
            getOptimizer());
    }

    @Override
    public PlcConnection getConnection(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Authentication not supported.");
    }

    /** Estimate the Length of a Packet */
    public static class ByteLengthEstimator implements ToIntFunction<ByteBuf> {
        @Override
        public int applyAsInt(ByteBuf byteBuf) {
            System.out.println("Estimate:- ");
            if (byteBuf.readableBytes() >= 8) {
                System.out.println("Estimate:- " + Integer.reverseBytes(byteBuf.getInt(byteBuf.readerIndex() + 4)));
                return Integer.reverseBytes(byteBuf.getInt(byteBuf.readerIndex() + 4));
            }
            return -1;
        }
    }

    @Override
    public OpcuaField prepareField(String query){
        return OpcuaField.of(query);
    }

}
