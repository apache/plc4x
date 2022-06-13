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
package org.apache.plc4x.java.transport.virtualcan;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;
import org.apache.plc4x.java.transport.can.CANTransport;
import org.apache.plc4x.java.transport.can.FrameData;
import org.apache.plc4x.java.transport.test.TestTransport;
import org.apache.plc4x.java.transport.virtualcan.io.VirtualCANFrameIO;

import java.util.function.Function;
import java.util.function.ToIntFunction;

public class VirtualCANTransport extends TestTransport implements CANTransport<VirtualCANFrame> {

    @Override
    public String getTransportCode() {
        return "virtualcan"; // vcan is used by socketcan, using virtualcan to avoid mismatches
    }

    @Override
    public String getTransportName() {
        return "Virtual CAN";
    }

    @Override
    public ToIntFunction<ByteBuf> getEstimator() {
        return new ToIntFunction<ByteBuf>() {
            @Override
            public int applyAsInt(ByteBuf value) {
                return value.getShort(value.readerIndex());
            }
        };
    }

    @Override
    public Class<VirtualCANFrame> getMessageType() {
        return VirtualCANFrame.class;
    }

    @Override
    public CANFrameBuilder<VirtualCANFrame> getTransportFrameBuilder() {
        return new CANFrameBuilder<VirtualCANFrame>() {

            private byte[] data;
            private int nodeId;

            @Override
            public CANFrameBuilder<VirtualCANFrame> withId(int nodeId) {
                this.nodeId = nodeId;
                return this;
            }

            @Override
            public CANFrameBuilder<VirtualCANFrame> withData(byte[] data) {
                this.data = data;
                return this;
            }

            @Override
            public VirtualCANFrame create() {
                return new VirtualCANFrame(nodeId, data);
            }
        };
    }

    @Override
    public Function<VirtualCANFrame, FrameData> adapter() {
        return new Function<VirtualCANFrame, FrameData>() {
            @Override
            public FrameData apply(VirtualCANFrame frame) {
                return new FrameData() {
                    @Override
                    public int getNodeId() {
                        return frame.getNodeId();
                    }

                    @Override
                    public <T extends Message> T read(MessageInput<T> input, Object... args) {
                        try {
                            return input.parse(new ReadBufferByteBased(getData(), ByteOrder.LITTLE_ENDIAN), args);
                        } catch (ParseException e) {
                            throw new PlcRuntimeException(e);
                        }
                    }

                    @Override
                    public int getDataLength() {
                        return frame.getData().length;
                    }

                    @Override
                    public byte[] getData() {
                        return frame.getData();
                    }
                };
            }
        };
    }

    @Override
    public MessageInput<VirtualCANFrame> getMessageInput(Configuration configuration) {
        return new VirtualCANFrameIO();
    }

}
