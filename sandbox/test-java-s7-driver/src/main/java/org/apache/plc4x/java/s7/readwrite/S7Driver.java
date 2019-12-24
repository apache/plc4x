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
package org.apache.plc4x.java.s7.readwrite;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.s7.readwrite.connection.S7Configuration;
import org.apache.plc4x.java.s7.readwrite.protocol.S7ProtocolLogic;
import org.apache.plc4x.java.s7.readwrite.protocol.S7ProtocolMessage;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcFieldHandler;
import org.apache.plc4x.java.spi.connection.DefaultNettyPlcConnection;
import org.apache.plc4x.java.spi.parser.ConnectionParser;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.osgi.service.component.annotations.Component;

import java.net.InetSocketAddress;
import java.net.SocketAddress;

@Component(service = PlcDriver.class, immediate = true)
public class S7Driver implements PlcDriver {

    private static final int ISO_ON_TCP_PORT = 102;

    @Override
    public String getProtocolCode() {
        return "s7ng";
    }

    @Override
    public String getProtocolName() {
        return "Siemens S7 (Basic)";
    }

    @Override
    public PlcConnection connect(String url) throws PlcConnectionException {
        ConnectionParser parser = new ConnectionParser(getProtocolCode(), url);
        S7Configuration configuration = parser.createConfiguration(S7Configuration.class);
        SocketAddress address = parser.getSocketAddress(ISO_ON_TCP_PORT);
        return new DefaultNettyPlcConnection<>(new TcpSocketChannelFactory(address), true, new S7PlcFieldHandler(),
            TPKTPacket.class,
            new S7ProtocolMessage(),
            new S7ProtocolLogic(
                configuration
            )
        );
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("Basic S7 connections don't support authentication (NG).");
    }

}
