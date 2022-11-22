/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.canopen.transport;

import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.MessageInput;
import org.apache.plc4x.java.transport.can.FrameData;

import java.util.function.Function;

/**
 * Adapter from {@link CANOpenFrame} to transport friendly form of frame represented by {@link FrameData}.
 *
 * Since {@link CANOpenFrame} is generated, hence it can be serialized. Yet, to save bandwidth with
 * {@link org.apache.plc4x.java.canopen.CANTestTransport} we just do bare minimum to get it working with
 * object to object mapping.
 */
public class CANOpenFrameDataAdapter implements Function<CANOpenFrame, FrameData> {

    @Override
    public FrameData apply(CANOpenFrame frame) {
        return new FrameData() {
            @Override
            public int getNodeId() {
                return frame.getService().getMin() + frame.getNodeId();
            }

            @Override
            public <T extends Message> T read(MessageInput<T> input, Object... args) {
                return (T) frame.getPayload();
            }

            @Override
            public int getDataLength() {
                return frame.getPayload().getLengthInBytes();
            }

            @Override
            public byte[] getData() {
                return new byte[0];
            }
        };
    }

}
