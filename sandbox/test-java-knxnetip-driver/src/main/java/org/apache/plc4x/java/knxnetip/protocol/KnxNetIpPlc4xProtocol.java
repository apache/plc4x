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
package org.apache.plc4x.java.knxnetip.protocol;

import io.netty.channel.ChannelHandlerContext;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.knxnetip.readwrite.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class KnxNetIpPlc4xProtocol extends PlcMessageToMessageCodec<KNXNetIPMessage, PlcRequestContainer> {

    private static final Logger LOGGER = LoggerFactory.getLogger(KnxNetIpPlc4xProtocol.class);

    @Override
    protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) {

    }

    @Override
    protected void decode(ChannelHandlerContext ctx, KNXNetIPMessage msg, List<Object> out) {
        if(msg instanceof TunnelingRequest) {
            TunnelingRequest tunnelingRequest = (TunnelingRequest) msg;
            CEMIBusmonInd busmonInd = (CEMIBusmonInd) tunnelingRequest.getCemi();
            if (busmonInd.getCemiFrame() instanceof CEMIFrameData) {
                outputStringRepresentation((CEMIFrameData) busmonInd.getCemiFrame());
            }
        }
    }

    private void outputStringRepresentation(CEMIFrameData data) {
        final KNXAddress sourceAddress = data.getSourceAddress();
        final byte[] destinationAddress = data.getDestinationAddress();
        final boolean groupAddress = data.getGroupAddress();
        final byte[] payload = new byte[data.getData().length + 1];
        payload[0] = data.getDataFirstByte();
        System.arraycopy(data.getData(), 0, payload, 1, data.getData().length);
        String payloadString = Hex.encodeHexString(payload);
    }

}
