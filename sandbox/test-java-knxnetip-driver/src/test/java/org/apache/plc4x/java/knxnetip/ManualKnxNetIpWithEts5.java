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

import io.netty.channel.ChannelHandlerContext;
import org.apache.commons.lang3.builder.ToStringStyle;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.connection.NettyPlcConnection;
import org.apache.plc4x.java.base.connection.UdpSocketChannelFactory;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.ets5.passive.KNXGroupAddress;
import org.apache.plc4x.java.ets5.passive.KnxDatapoint;
import org.apache.plc4x.java.ets5.passive.io.KNXGroupAddressIO;
import org.apache.plc4x.java.ets5.passive.io.KnxDatapointIO;
import org.apache.plc4x.java.knxnetip.connection.KnxNetIpConnection;
import org.apache.plc4x.java.knxnetip.ets5.Ets5Parser;
import org.apache.plc4x.java.knxnetip.ets5.model.Ets5Model;
import org.apache.plc4x.java.knxnetip.ets5.model.GroupAddress;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.apache.plc4x.java.utils.ReadBuffer;

import java.io.File;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class ManualKnxNetIpWithEts5 {

    public static void main(String[] args) throws Exception {
        Ets5Model ets5Model = new Ets5Parser().parse(new File("/Users/christofer.dutz/Projects/Apache/PLC4X-Documents/KNX/Stettiner Str. 13/StettinerStr-Soll-Ist-Temperatur.knxproj"));;
        final byte groupAddressType = ets5Model.getGroupAddressType();

        InetAddress inetAddress = InetAddress.getByName("192.168.42.11");
        ChannelFactory channelFactory = new UdpSocketChannelFactory(inetAddress, KnxNetIpConnection.KNXNET_IP_PORT);

        NettyPlcConnection connection = new KnxNetIpConnection(channelFactory, "",
            new PlcMessageToMessageCodec<KNXNetIPMessage, PlcRequestContainer>() {
                @Override
                protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {
                    // Ignore for now ...
                }

                @Override
                protected void decode(ChannelHandlerContext ctx, KNXNetIPMessage packet, List<Object> out) throws Exception {
                    if(packet instanceof TunnelingRequest) {
                        TunnelingRequest request = (TunnelingRequest) packet;
                        CEMI cemiPayload = request.getCemi();
                        if(cemiPayload instanceof CEMIBusmonInd) {
                            CEMIBusmonInd cemiBusmonInd = (CEMIBusmonInd) cemiPayload;
                            if(cemiBusmonInd.getCemiFrame() instanceof CEMIFrameData) {
                                CEMIFrameData cemiDataFrame = (CEMIFrameData) cemiBusmonInd.getCemiFrame();

                                // The first byte is actually just 6 bit long, but we'll treat it as a full one.
                                // So here we create a byte array containing the first and all the following bytes.
                                byte[] payload = new byte[1 + cemiDataFrame.getData().length];
                                payload[0] = cemiDataFrame.getDataFirstByte();
                                System.arraycopy(cemiDataFrame.getData(), 0, payload, 1, cemiDataFrame.getData().length);

                                final KNXAddress sourceAddress = cemiDataFrame.getSourceAddress();
                                final byte[] destinationGroupAddress = cemiDataFrame.getDestinationAddress();

                                ReadBuffer addressReadBuffer = new ReadBuffer(destinationGroupAddress);
                                // Decode the group address depending on the project settings.
                                KNXGroupAddress destinationAddress =
                                    KNXGroupAddressIO.parse(addressReadBuffer, ets5Model.getGroupAddressType());
                                final GroupAddress groupAddress = ets5Model.getGroupAddresses().get(destinationAddress);

                                ReadBuffer rawDataReader = new ReadBuffer(payload);

                                final KnxDatapoint datapoint = KnxDatapointIO.parse(rawDataReader, groupAddress.getType().getMainType(), groupAddress.getType().getSubType());
                                final String jsonDatapoint = datapoint.toString(ToStringStyle.JSON_STYLE);

                                System.out.println(destinationAddress.toString() + " " + jsonDatapoint);
                            } else {
                                System.out.println(packet);
                            }
                        } else {
                            System.out.println(packet);
                        }
                    } else {
                        System.out.println(packet);
                    }
                }
            });

        connection.connect();
        TimeUnit.SECONDS.sleep(3000);
        connection.close();
    }

}
