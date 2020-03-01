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
package org.apache.plc4x.java.bacnetip;

import org.apache.plc4x.java.bacnetip.connection.PassiveBacNetIpPlcConnection;
import org.apache.plc4x.java.bacnetip.protocol.HelloWorldProtocol;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.connection.PcapChannelFactory;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketAddress;
import org.apache.plc4x.java.utils.pcapsockets.netty.PcapSocketChannelConfig;
import org.apache.plc4x.java.utils.pcapsockets.netty.UdpIpPacketHandler;

import java.io.File;

public class PassiveBacNetIpDriverManual {

    public static void main(String[] args) throws Exception {
        NettyPlcConnection connection = new PassiveBacNetIpPlcConnection(new PcapChannelFactory(
            //new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/BacNET/Captures/Merck/BACnetWhoIsRouterToNetwork.pcapng"), null,
            new File("/Users/christofer.dutz/Downloads/20190906_udp.pcapng"), null,
            PassiveBacNetIpDriver.BACNET_IP_PORT, PcapSocketAddress.ALL_PROTOCOLS,
            PcapSocketChannelConfig.SPEED_FAST_DOUBLE, new UdpIpPacketHandler()), "", new HelloWorldProtocol());
        connection.connect();
    }

}
