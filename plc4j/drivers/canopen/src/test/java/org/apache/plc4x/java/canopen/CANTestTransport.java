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
package org.apache.plc4x.java.canopen;

import io.netty.buffer.ByteBuf;

import java.util.function.Function;
import java.util.function.ToIntFunction;

import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.io.CANOpenFrameIO;
import org.apache.plc4x.java.canopen.transport.CANOpenFrameDataAdapter;
import org.apache.plc4x.java.canopen.transport.IdentityCANOpenFrameBuilder;
import org.apache.plc4x.java.spi.configuration.Configuration;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.generation.MessageInput;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;
import org.apache.plc4x.java.transport.can.CANTransport;
import org.apache.plc4x.java.transport.can.FrameData;
import org.apache.plc4x.java.transport.test.TestTransport;

public class CANTestTransport extends TestTransport implements CANTransport<CANOpenFrame> {
    @Override
    public ToIntFunction<ByteBuf> getEstimator() {
        // id (1 byte), service (1 byte), data (up to 8 bytes, padded)
        return (buffer) -> 10;
    }

    @Override
    public Class<CANOpenFrame> getMessageType() {
        return CANOpenFrame.class;
    }

    @Override
    public MessageInput<CANOpenFrame> getMessageInput(Configuration configuration) {
        return new CANOpenFrameIO();
    }

    @Override
    public CANFrameBuilder<CANOpenFrame> getTransportFrameBuilder() {
        return new IdentityCANOpenFrameBuilder();
    }

    @Override
    public Function<CANOpenFrame, FrameData> adapter() {
        return new CANOpenFrameDataAdapter();
    }
}
