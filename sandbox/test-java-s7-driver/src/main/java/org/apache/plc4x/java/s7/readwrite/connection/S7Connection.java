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
package org.apache.plc4x.java.s7.readwrite.connection;

import io.netty.channel.*;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.connection.GenericNettyPlcConnection;
import org.apache.plc4x.java.spi.connection.NettyPlcConnection;
import org.apache.plc4x.java.spi.events.ConnectEvent;
import org.apache.plc4x.java.spi.events.ConnectedEvent;
import org.apache.plc4x.java.s7.readwrite.TPKTPacket;
import org.apache.plc4x.java.s7.readwrite.protocol.S7ProtocolLogic;
import org.apache.plc4x.java.s7.readwrite.protocol.S7ProtocolMessage;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcFieldHandler;
import org.apache.plc4x.java.spi.Plc4xNettyWrapper;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteRequest;
import org.apache.plc4x.java.spi.messages.PlcReader;
import org.apache.plc4x.java.spi.messages.PlcWriter;
import org.apache.plc4x.java.spi.parser.ConnectionParser;
import org.apache.plc4x.java.tcp.connection.TcpSocketChannelFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.concurrent.CompletableFuture;

public class S7Connection extends GenericNettyPlcConnection<TPKTPacket> implements PlcReader, PlcWriter {

    private static final int ISO_ON_TCP_PORT = 102;

    private static final Logger logger = LoggerFactory.getLogger(S7Connection.class);

    public S7Connection(InetAddress address, String params) {
        this(new TcpSocketChannelFactory(address, ISO_ON_TCP_PORT), params);
    }

    public S7Connection(ChannelFactory channelFactory, String params) {
        super(channelFactory, true, new S7PlcFieldHandler(),
            TPKTPacket.class,
            new S7ProtocolMessage(),
            new S7ProtocolLogic(
                ConnectionParser.parse("a://1.1.1.1/" + params, S7Configuration.class)
            )
        );
    }

}
