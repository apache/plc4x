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
package org.apache.plc4x.java.canopen.transport;

import java.util.function.Supplier;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.canopen.readwrite.utils.StaticHelper;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;
import org.apache.plc4x.java.transport.can.CANTransport.FrameHandler;
import org.apache.plc4x.java.transport.can.FrameData;

/**
 * Dedicated frame handler instance which can create {@link CANOpenFrame} from received {@link FrameData} and dump it back
 * via passed {@link CANFrameBuilder} back to the wire.
 */
public class CANOpenFrameDataHandler implements FrameHandler<Message, CANOpenFrame> {

    private final Supplier<CANFrameBuilder<Message>> builder;

    public CANOpenFrameDataHandler(Supplier<CANFrameBuilder<Message>> builder) {
        this.builder = builder;
    }

    @Override
    public CANOpenFrame fromCAN(FrameData frame) {
        CANOpenService service = StaticHelper.serviceId((short) frame.getNodeId());
        int nodeId = Math.abs(service.getMin() - frame.getNodeId());
        return new CANOpenFrame((short) nodeId, service, frame.read(CANOpenPayload::staticParse, service));
    }

    @Override
    public Message toCAN(CANOpenFrame frame) {
        try {
            CANOpenPayload payload = frame.getPayload();
            WriteBufferByteBased buffer = new WriteBufferByteBased(payload.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
            payload.serialize(buffer);
            return builder.get().withId(frame.getService().getMin() + frame.getNodeId())
                .withData(buffer.getData())
                .create();
        } catch (SerializationException e) {
            throw new PlcRuntimeException(e);
        }
    }

}
