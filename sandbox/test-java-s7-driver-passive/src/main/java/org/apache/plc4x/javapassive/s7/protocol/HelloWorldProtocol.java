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
package org.apache.plc4x.javapassive.s7.protocol;

import io.netty.channel.ChannelHandlerContext;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.s7.passive.S7Message;
import org.apache.plc4x.java.s7.passive.TPKTPacket;

import java.util.List;

public class HelloWorldProtocol extends PlcMessageToMessageCodec<TPKTPacket, PlcRequestContainer> {

    @Override
    protected void encode(ChannelHandlerContext channelHandlerContext, PlcRequestContainer plcRequestContainer, List<Object> list) throws Exception {
        System.out.println(plcRequestContainer);
    }

    @Override
    protected void decode(ChannelHandlerContext channelHandlerContext, TPKTPacket tpktPacket, List<Object> list) throws Exception {
        StringBuilder sb = new StringBuilder();
        sb.append("[").append(tpktPacket.getPayload().getClass().getName());
        if(tpktPacket.getPayload().getPayload() != null) {
            S7Message s7Message = tpktPacket.getPayload().getPayload();
            sb.append(" \n  [").append(s7Message.getClass()).append("]\n");
        }
        sb.append("]");
        System.out.println(sb.toString());
    }

}
