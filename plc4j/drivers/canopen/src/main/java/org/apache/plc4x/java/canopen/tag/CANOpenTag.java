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
package org.apache.plc4x.java.canopen.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.regex.Pattern;

/**
 * Generic tag type which defines node address and address pattern (index/subindex).
 */
public abstract class CANOpenTag implements PlcTag, Serializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?:(0[xX](?<indexHex>[0-9a-fA-F]+))|(?<index>\\d+))/(?:(0[xX](?<subIndexHex>[0-9a-fA-F]+))|(?<subIndex>\\d+)):(?<canDataType>\\w+)(\\[(?<numberOfElements>\\d)])?");
    public static final Pattern NODE_PATTERN = Pattern.compile("(?<nodeId>\\d+)");

    private final int nodeId;

    public CANOpenTag(int nodeId) {
        this.nodeId = nodeId;
    }

    public int getNodeId() {
        return nodeId;
    }

    public static CANOpenTag of(String addressString) throws PlcInvalidTagException {
        if (CANOpenSDOTag.matches(addressString)) {
            return CANOpenSDOTag.of(addressString);
        } else if (CANOpenPDOTag.matches(addressString)) {
            return CANOpenPDOTag.of(addressString);
        } else if (CANOpenNMTTag.matches(addressString)) {
            return CANOpenNMTTag.of(addressString);
        } else if (CANOpenHeartbeatTag.matches(addressString)) {
            return CANOpenHeartbeatTag.of(addressString);
        }

        throw new PlcInvalidTagException("Unable to parse address: " + addressString);
    }

}
