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
package org.apache.plc4x.java.transport.rawsocket;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.api.configuration.PlcTransportConfiguration;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.configuration.HasConfiguration;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.transport.Transport;
import org.apache.plc4x.java.utils.rawsockets.netty.address.RawSocketAddress;

import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class RawSocketTransport implements Transport, HasConfiguration<RawSocketTransportConfiguration> {

    private static final Pattern TRANSPORT_RAW_SOCKET_IP_PATTERN = Pattern.compile(
        "^((?<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?<hostname>[a-zA-Z0-9.\\-]+))(:(?<port>[0-9]{1,5}))?");

    private static final Pattern TRANSPORT_RAW_SOCKET_MAC_PATTERN = Pattern.compile(
        "^(?<macAddress>[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2})");

    public static final String TRANSPORT_CODE = "raw";

    private RawSocketTransportConfiguration configuration;

    @Override
    public String getTransportCode() {
        return TRANSPORT_CODE;
    }

    @Override
    public String getTransportName() {
        return "Raw Ethernet Transport";
    }

    @Override
    public void setConfiguration(RawSocketTransportConfiguration configuration) {
        this.configuration = configuration;
    }

    @Override
    public ChannelFactory createChannelFactory(String transportConfig) {
        final Matcher macMatcher = TRANSPORT_RAW_SOCKET_MAC_PATTERN.matcher(transportConfig);
        if(macMatcher.matches()) {
            String macAddressString = macMatcher.group("macAddress");
            try {
                byte[] macAddress = Hex.decodeHex(macAddressString.replace(":", ""));
                // Create the fully qualified remote socket address which we should connect to.
                RawSocketAddress address = new RawSocketAddress(macAddress);
                return new RawSocketChannelFactory(address);
            } catch (DecoderException e) {
                throw new RuntimeException(e);
            }
        }

        final Matcher ipMatcher = TRANSPORT_RAW_SOCKET_IP_PATTERN.matcher(transportConfig);
        if(!ipMatcher.matches()) {
            throw new PlcRuntimeException("Invalid url for Raw socket transport");
        }

        String ip = ipMatcher.group("ip");
        String hostname = ipMatcher.group("hostname");
        String portString = ipMatcher.group("port");

        // If the port wasn't specified, try to get a default port from the configuration.
        int port;
        if(portString != null) {
            port = Integer.parseInt(portString);
        } else if ((configuration != null) &&
            (configuration.getDefaultPort() != RawSocketTransportConfiguration.NO_DEFAULT_PORT)) {
            port = configuration.getDefaultPort();
        } else {
            throw new PlcRuntimeException("No port defined");
        }

        // Create the fully qualified remote socket address which we should connect to.
        SocketAddress address = new InetSocketAddress((ip == null) ? hostname : ip, port);

        RawSocketChannelFactory rawSocketChannelFactory = new RawSocketChannelFactory(address);
        if(configuration != null) {
            rawSocketChannelFactory.setConfiguration(configuration);
        }
        return rawSocketChannelFactory;
    }

    @Override
    public Class<? extends PlcTransportConfiguration> getTransportConfigType() {
        return DefaultRawSocketTransportConfiguration.class;
    }

}
