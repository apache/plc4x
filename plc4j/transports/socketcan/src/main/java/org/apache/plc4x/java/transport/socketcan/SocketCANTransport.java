/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.transport.socketcan;

import io.netty.buffer.ByteBuf;

import java.util.function.Function;
import java.util.function.ToIntFunction;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.socketcan.readwrite.SocketCANFrame;
import org.apache.plc4x.java.socketcan.readwrite.io.SocketCANFrameIO;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.connection.ChannelFactory;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;
import org.apache.plc4x.java.transport.can.CANTransport;
import org.apache.plc4x.java.transport.can.FrameData;
import org.apache.plc4x.java.transport.socketcan.netty.address.SocketCANAddress;

/**
 * CAN specific transport which rely on JavaCAN library.
 */
public class SocketCANTransport implements CANTransport<SocketCANFrame> {

    @Override
    public String getTransportCode() {
        return "socketcan";
    }

    @Override
    public String getTransportName() {
        return "SocketCAN Transport (based on JavaCAN)";
    }

    @Override
    public ChannelFactory createChannelFactory(String transportConfig) {
        SocketCANAddress address = new SocketCANAddress(transportConfig);
        return new SocketCANChannelFactory(address);
    }

    @Override
    public ToIntFunction<ByteBuf> getEstimator() {
        return (buff) -> 16;
    }

    /*@Override
    public MessageIO<SocketCANFrame, SocketCANFrame> getMessageIO(Configuration cfg) {
        return new SocketCANFrameIO();
    }*/

    @Override
    public Class<SocketCANFrame> getMessageType() {
        return SocketCANFrame.class;
    }

    @Override
    public CANFrameBuilder<SocketCANFrame> getTransportFrameBuilder() {
        return new SocketCANFrameBuilder();
    }

    @Override
    public Function<SocketCANFrame, FrameData> adapter() {
        return new Function<SocketCANFrame, FrameData>() {
            @Override
            public FrameData apply(SocketCANFrame frame) {
                return new FrameData() {
                    @Override
                    public int getNodeId() {
                        return frame.getIdentifier();
                    }

                    @Override
                    public int getDataLength() {
                        return frame.getData().length;
                    }

                    @Override
                    public byte[] getData() {
                        return frame.getData();
                    }

                    @Override
                    public <T extends Message> T read(MessageIO<T, T> serializer, Object... args) {
                        try {
                            return serializer.parse(new ReadBufferByteBased(frame.getData(), ByteOrder.LITTLE_ENDIAN), args);
                        } catch (ParseException e) {
                            throw new PlcRuntimeException(e);
                        }
                    }
                };
            }
        };
    }
}
