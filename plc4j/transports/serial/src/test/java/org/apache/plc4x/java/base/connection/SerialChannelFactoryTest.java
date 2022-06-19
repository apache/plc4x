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
package org.apache.plc4x.java.base.connection;

import com.fazecast.jSerialComm.SerialPort;
import io.netty.buffer.ByteBuf;
import io.netty.channel.Channel;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInitializer;
import io.netty.handler.codec.ByteToMessageCodec;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.transport.serial.DummyHandler;
import org.apache.plc4x.java.transport.serial.SerialChannel;
import org.apache.plc4x.java.transport.serial.SerialChannelFactory;
import org.apache.plc4x.java.transport.serial.SerialChannelHandler;
import org.apache.plc4x.java.transport.serial.SerialSocketAddress;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.net.UnknownHostException;
import java.util.List;

public class SerialChannelFactoryTest {

    private static final Logger logger = LoggerFactory.getLogger(SerialChannelFactoryTest.class);

    @Test
    public void showAllPorts() {
        System.out.println("-------------------------------------");
        System.out.println(" Starting to Display all Serial Ports");
        System.out.println("-------------------------------------");
        for (SerialPort commPort : SerialPort.getCommPorts()) {
            System.out.println(commPort.getDescriptivePortName());
        }
    }

    @Test
    public void createChannel() throws PlcConnectionException, InterruptedException, UnknownHostException {
        SerialChannelFactory asdf = new SerialChannelFactory(new SerialSocketAddress("TEST-port1", DummyHandler.INSTANCE));
        // final TcpSocketChannelFactory factory = new TcpSocketChannelFactory(InetAddress.getLocalHost(), 5432);
        final Channel channel = asdf.createChannel(new ChannelInitializer<SerialChannel>() {
            @Override
            protected void initChannel(SerialChannel ch) throws Exception {
                ch.pipeline().addLast(new DemoCodec());
            }
        });
        Thread.sleep(100);
        for (int i = 1; i <= 10; i++) {
            Thread.sleep(10);
            DummyHandler.INSTANCE.fireEvent(1);
        }
        Thread.sleep(100);
        channel.close().sync();
    }

    @Test
    @Disabled("Seems to cause problems on Windows, but as it generally only woks on devices with a 'JBLFlip3-SPPDev' it's not much use anyway")
    public void createChannelToSBL() throws PlcConnectionException, InterruptedException, UnknownHostException {
        SerialChannelFactory asdf = new SerialChannelFactory(new SerialSocketAddress("JBLFlip3-SPPDev"));
        Channel channel = null;
        try {
            channel = asdf.createChannel(new ChannelInitializer<SerialChannel>() {
                @Override
                protected void initChannel(SerialChannel ch) throws Exception {
                    ch.pipeline().addLast(new DemoCodec());
                }
            });
        } catch (Exception e) {
            // do nothing
        }
        Thread.sleep(5_000);
        if (channel != null) {
            channel.close().sync();
        }
    }

    private static class DemoCodec extends ByteToMessageCodec<Object> {
        @Override
        protected void encode(ChannelHandlerContext channelHandlerContext, Object o, ByteBuf byteBuf) throws Exception {
            // do nothing here
        }

        @Override
        protected void decode(ChannelHandlerContext channelHandlerContext, ByteBuf byteBuf, List<Object> list) throws Exception {
            byteBuf.markReaderIndex();
            StringBuilder sb = new StringBuilder();
            for (int i = 1; i <= byteBuf.readableBytes(); i++) {
                sb.append(byteBuf.readByte()).append(", ");
            }
            byteBuf.resetReaderIndex();
            logger.debug("We currently have {} readable bytes: {}", byteBuf.readableBytes(), sb.toString());
        }
    }
}