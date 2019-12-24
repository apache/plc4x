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
package org.apache.plc4x.java.knxnetip;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.base.connection.UdpSocketChannelFactory;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpConfiguration;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpFieldHandler;
import org.apache.plc4x.java.knxnetip.protocol.KnxNetIpProtocolLogic;
import org.apache.plc4x.java.knxnetip.protocol.KnxNetIpProtocolMessage;
import org.apache.plc4x.java.knxnetip.readwrite.KNXNetIPMessage;
import org.apache.plc4x.java.spi.connection.DefaultNettyPlcConnection;
import org.apache.plc4x.java.spi.parser.ConnectionParser;

import java.net.*;

public class KnxNetIpDriver implements PlcDriver {

    public static final int KNXNET_IP_PORT = 3671;

    @Override
    public String getProtocolCode() {
        return "knxnet-ip";
    }

    @Override
    public String getProtocolName() {
        return "KNXNet/IP";
    }

    @Override
    public PlcConnection connect(String connectionString) throws PlcConnectionException {
        ConnectionParser parser = new ConnectionParser(getProtocolCode(), connectionString);
        KnxNetIpConfiguration configuration = parser.createConfiguration(KnxNetIpConfiguration.class);
        SocketAddress address = parser.getSocketAddress(KNXNET_IP_PORT);
        return new DefaultNettyPlcConnection<>(new UdpSocketChannelFactory(address), true, new KnxNetIpFieldHandler(),
            KNXNetIPMessage.class,
            new KnxNetIpProtocolMessage(),
            new KnxNetIpProtocolLogic(
                configuration
            )
        );
    }

    @Override
    public PlcConnection connect(String url, PlcAuthentication authentication) throws PlcConnectionException {
        throw new PlcConnectionException("KNXNet/IP connections don't support authentication.");
    }

}
