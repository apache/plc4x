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
package org.apache.plc4x.java.canopen.socketcan;

import org.apache.plc4x.java.canopen.api.conversation.canopen.CANConversation;
import org.apache.plc4x.java.canopen.transport.CANOpenFrame;
import org.apache.plc4x.java.canopen.transport.CANOpenFrameBuilder;
import org.apache.plc4x.java.canopen.transport.CANOpenFrameBuilderFactory;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.ConversationContext.SendRequestContext;

import java.time.Duration;

public class SocketCANConversation implements CANConversation<CANOpenFrame> {

    private final int nodeId;
    private final ConversationContext<CANOpenFrame> context;
    private final int timeout;
    private final CANOpenFrameBuilderFactory factory;

    public SocketCANConversation(int nodeId, ConversationContext<CANOpenFrame> context, int timeout, CANOpenFrameBuilderFactory factory) {
        this.nodeId = nodeId;
        this.context = context;
        this.timeout = timeout;
        this.factory = factory;
    }

    @Override
    public int getNodeId() {
        return nodeId;
    }

    @Override
    public CANOpenFrameBuilder createBuilder() {
        return factory.createBuilder();
    }

    public SendRequestContext<CANOpenFrame> send(CANOpenFrame frame) {
        return context.sendRequest(frame)
            .expectResponse(CANOpenFrame.class, Duration.ofMillis(timeout));
    }

    public void sendToWire(CANOpenFrame frame) {
        context.sendToWire(frame);
    }

}
