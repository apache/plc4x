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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.utils.StaticHelper;
import org.apache.plc4x.java.canopen.readwrite.CANOpenFrame;
import org.apache.plc4x.java.canopen.readwrite.CANOpenService;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.transport.can.CANFrameBuilder;

/**
 * A frame builder which simply re-constructs {@link CANOpenFrame} from its fields.
 *
 * This type is used to declare wire format of data on CAN bus.
 */
public class IdentityCANOpenFrameBuilder implements CANFrameBuilder<CANOpenFrame> {

    private short nodeId;
    private CANOpenService service;
    private byte[] data;

    @Override
    public CANFrameBuilder<CANOpenFrame> withId(int nodeId) {
        this.service = StaticHelper.serviceId((short) nodeId);
        this.nodeId = (short) Math.abs(service.getMin() - nodeId);
        return this;
    }

    @Override
    public CANFrameBuilder<CANOpenFrame> withData(byte[] data) {
        this.data = data;
        return this;
    }

    @Override
    public CANOpenFrame create() {
        try {
            return new CANOpenFrame(
                nodeId, service, CANOpenPayload.staticParse(new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN), service)
            );
        } catch (ParseException e) {
            throw new PlcRuntimeException(e);
        }
    }

}
