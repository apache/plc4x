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
package org.apache.plc4x.java.can.generic.transport;

import java.util.function.Supplier;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;
import org.apache.plc4x.java.transport.can.CANTransport;
import org.apache.plc4x.java.transport.can.FrameData;

/**
 * Generic CAN frame handler turn a wire level message @{@link FrameData} into a wrapper which
 * does not have any specific other than node id and data.
 *
 * Because it is used by generic purpose driver it can not assume any semantics on message role.
 */
public class GenericCANFrameDataHandler implements CANTransport.FrameHandler<Message, GenericFrame> {

    private final Supplier<CANFrameBuilder<Message>> frameBuilder;

    public GenericCANFrameDataHandler(Supplier<CANFrameBuilder<Message>> frameBuilder) {
        this.frameBuilder = frameBuilder;
    }

    @Override
    public GenericFrame fromCAN(FrameData frame) {
        return new GenericFrame(frame.getNodeId(), frame.getData());
    }

    @Override
    public Message toCAN(GenericFrame frame) {
        return frameBuilder.get().withId(frame.getNodeId())
            .withData(frame.getData())
            .create();
    }

}
