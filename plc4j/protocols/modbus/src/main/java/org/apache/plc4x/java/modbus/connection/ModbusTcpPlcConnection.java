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
package org.apache.plc4x.java.modbus.connection;

import com.digitalpetri.modbus.codec.ModbusRequestEncoder;
import com.digitalpetri.modbus.codec.ModbusResponseDecoder;
import com.digitalpetri.modbus.codec.ModbusTcpCodec;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandler;
import io.netty.channel.ChannelInitializer;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.TcpSocketChannelFactory;
import org.apache.plc4x.java.base.protocol.SingleItemToSingleRequestProtocol;
import org.apache.plc4x.java.modbus.netty.Plc4XModbusProtocol;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.InetAddress;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class ModbusTcpPlcConnection extends BaseModbusPlcConnection {

    private static final int MODBUS_TCP_PORT = 502;

    private static final Logger logger = LoggerFactory.getLogger(ModbusTcpPlcConnection.class);

    private ModbusTcpPlcConnection(InetAddress address, String params) {
        this(new TcpSocketChannelFactory(address, MODBUS_TCP_PORT), params);
        logger.info("Configured ModbusTcpPlcConnection with: host-name {}", address.getHostAddress());
    }

    public ModbusTcpPlcConnection(InetAddress address, int port, String params) {
        this(new TcpSocketChannelFactory(address, port), params);
        logger.info("Configured ModbusTcpPlcConnection with: host-name {}", address.getHostAddress());
    }

    public ModbusTcpPlcConnection(ChannelFactory channelFactory, String params) {
        super(channelFactory, params);
    }

    public static ModbusTcpPlcConnection of(InetAddress address, String params) {
        return new ModbusTcpPlcConnection(address, params);
    }

    public static ModbusTcpPlcConnection of(InetAddress address, int port, String params) {
        return new ModbusTcpPlcConnection(address, port, params);
    }

    @Override
    protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
        return new ChannelInitializer() {
            @Override
            protected void initChannel(Channel channel) {
                channel.pipeline().addLast(new ModbusTcpCodec(new ModbusRequestEncoder(), new ModbusResponseDecoder()));
                channel.pipeline().addLast(new Plc4XModbusProtocol());
                channel.pipeline().addLast(new SingleItemToSingleRequestProtocol(ModbusTcpPlcConnection.this, ModbusTcpPlcConnection.this, timer));
            }
        };
    }

    public InetAddress getRemoteAddress() {
        return ((TcpSocketChannelFactory) channelFactory).getAddress();
    }
}
