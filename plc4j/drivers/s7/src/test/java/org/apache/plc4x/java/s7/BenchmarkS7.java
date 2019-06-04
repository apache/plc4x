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

package org.apache.plc4x.java.s7;
import io.netty.buffer.Unpooled;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.embedded.EmbeddedChannel;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.base.PlcMessageToMessageCodec;
import org.apache.plc4x.java.base.messages.PlcRequestContainer;
import org.apache.plc4x.java.isoontcp.protocol.IsoOnTcpProtocol;
import org.apache.plc4x.java.isotp.protocol.IsoTPProtocol;
import org.apache.plc4x.java.isotp.protocol.model.types.TpduSize;
import org.apache.plc4x.java.s7.netty.S7Protocol;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.types.S7ControllerType;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class BenchmarkS7 {

    public static void main(String[] args) throws Exception {
        byte[] data = Hex.decodeHex("0300006702f080320100000001005600000407120a10060001032b84000160120a10020001032b840001a0120a10010001032b840001a9120a10050001032b84000150120a10020001032b84000198120a10040001032b84000140120a10020001032b84000190");
        TimedMessageCodec timedMessageCodec = new TimedMessageCodec();
        EmbeddedChannel channel = new EmbeddedChannel(
            new IsoOnTcpProtocol(),
            new IsoTPProtocol((short) 0, (short) 0, TpduSize.SIZE_512),
            new S7Protocol((short) 1, (short) 1, (short) 512, S7ControllerType.S7_1500, null),
            timedMessageCodec
            );

        long start = System.currentTimeMillis();
        int numRuns = 20000;
        for(int i = 0; i < numRuns; i++) {
            run(data, channel, timedMessageCodec);
        }
        long end = System.currentTimeMillis();
        System.out.println("Parsed " + numRuns + " packets in " + (end - start) + "ms");
        System.out.println("That's " + ((float) (end - start) / numRuns) + "ms per packet");
    }

    public static void run(byte[] data, EmbeddedChannel channel, TimedMessageCodec timedMessageCodec) {
        CompletableFuture<S7Message> future = new CompletableFuture<>();
        timedMessageCodec.setFuture(future);
        channel.writeInbound(Unpooled.wrappedBuffer(data));
        try {
            future.get();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private static class TimedMessageCodec extends PlcMessageToMessageCodec<S7Message, PlcRequestContainer> {

        CompletableFuture<S7Message> future;

        public void setFuture(CompletableFuture<S7Message> future) {
            this.future = future;
        }

        @Override
        protected void encode(ChannelHandlerContext ctx, PlcRequestContainer msg, List<Object> out) throws Exception {

        }

        @Override
        protected void decode(ChannelHandlerContext ctx, S7Message msg, List<Object> out) throws Exception {
            future.complete(msg);
        }
    }

}
