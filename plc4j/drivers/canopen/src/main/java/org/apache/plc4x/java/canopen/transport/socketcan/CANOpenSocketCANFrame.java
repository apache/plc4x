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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;
import org.apache.plc4x.java.canopen.transport.CANOpenFrame;
import org.apache.plc4x.java.canopen.transport.socketcan.io.CANOpenSocketCANFrameIO;
import org.apache.plc4x.java.canopen.readwrite.CANOpenPayload;
import org.apache.plc4x.java.canopen.readwrite.types.CANOpenService;
import org.apache.plc4x.java.spi.generation.MessageIO;

import java.util.Objects;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class CANOpenSocketCANFrame implements CANOpenFrame {

    // Properties.
    private final int nodeId;
    private final CANOpenService service;
    private final CANOpenPayload payload;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public CANOpenSocketCANFrame(@JsonProperty("nodeId") int nodeId, @JsonProperty("service") CANOpenService service, @JsonProperty("payload") CANOpenPayload payload) {
        this.nodeId = nodeId;
        this.service = service;
        this.payload = payload;
    }

    public CANOpenSocketCANFrame(int nodeId, CANOpenPayload payload) {
        this(nodeId, payload.getFunction(), payload);
    }

    @Override
    public int getNodeId() {
        return nodeId;
    }

    @Override
    public CANOpenService getService() {
        return service;
    }

    @Override
    public CANOpenPayload getPayload() {
        return payload;
    }

    @Override
    @JsonIgnore
    public int getLengthInBytes() {
        return getLengthInBits() / 8;
    }

    @Override
    @JsonIgnore
    public int getLengthInBits() {
        int lengthInBits = 0;

        // Simple field (node + service)
        lengthInBits += 32;

        // A virtual field doesn't have any in- or output.

        // A virtual field doesn't have any in- or output.

        // A virtual field doesn't have any in- or output.

        // A virtual field doesn't have any in- or output.

        // Implicit Field (size)
        lengthInBits += 8;

        // Reserved Field (reserved)
        lengthInBits += 8;

        // Reserved Field (reserved)
        lengthInBits += 8;

        // Reserved Field (reserved)
        lengthInBits += 8;

        // Array field
        if (payload != null) {
            lengthInBits += 8 * payload.getLengthInBytes();
        }

        // Padding Field (padding)
        int _timesPadding = (8) - payload.getLengthInBytes();
        while (_timesPadding-- > 0) {
            lengthInBits += 8;
        }

        return lengthInBits;
    }

    @Override
    @JsonIgnore
    public MessageIO<CANOpenFrame, CANOpenFrame> getMessageIO() {
        return new CANOpenSocketCANFrameIO();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof CANOpenSocketCANFrame)) {
            return false;
        }
        CANOpenSocketCANFrame that = (CANOpenSocketCANFrame) o;
        return (getNodeId() == that.getNodeId()) &&
            (getService() == that.getService()) &&
            (getPayload() == that.getPayload());
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            getNodeId(),
            getService(),
            getPayload()
        );
    }

    @Override
    public String toString() {
        return toString(ToStringStyle.SHORT_PREFIX_STYLE);
    }

    public String toString(ToStringStyle style) {
        return new ToStringBuilder(this, style)
            .append("nodeId", getNodeId())
            .append("service", getService())
            .append("payload", getPayload())
            .toString();
    }

}
