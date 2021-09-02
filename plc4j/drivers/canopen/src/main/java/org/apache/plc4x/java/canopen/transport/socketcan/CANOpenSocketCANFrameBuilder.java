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
package org.apache.plc4x.java.canopen.transport.socketcan;

import org.apache.plc4x.java.canopen.transport.CANOpenFrame;
import org.apache.plc4x.java.canopen.transport.CANOpenFrameBuilder;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.types.CANOpenService;

public class CANOpenSocketCANFrameBuilder implements CANOpenFrameBuilder {

    private int node;
    private CANOpenPayload payload;
    private CANOpenService service;

    @Override
    public CANOpenFrameBuilder withNodeId(int node) {
        this.node = node;
        return this;
    }

    @Override
    public CANOpenFrameBuilder withService(CANOpenService service) {
        this.service = service;
        return this;
    }

    @Override
    public CANOpenFrameBuilder withPayload(CANOpenPayload payload) {
        this.payload = payload;
        return this;
    }

    @Override
    public CANOpenFrame build() {
        return new CANOpenSocketCANFrame(node, service, payload);
    }

}
