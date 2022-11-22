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
package org.apache.plc4x.java.transport.socketcan;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.socketcan.readwrite.SocketCANFrame;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;

public class SocketCANFrameBuilder implements CANFrameBuilder<SocketCANFrame> {

    private int nodeId;
    private byte[] data;

    @Override
    public CANFrameBuilder<SocketCANFrame> withId(int nodeId) {
        this.nodeId = nodeId;
        return this;
    }

    @Override
    public CANFrameBuilder<SocketCANFrame> withData(byte[] data) {
        this.data = data;
        return this;
    }

    @Override
    public SocketCANFrame create() throws PlcRuntimeException {
        return new SocketCANFrame(nodeId, data);
    }

}
