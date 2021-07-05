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
package org.apache.plc4x.java.canopen.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.regex.Pattern;

/**
 * Generic field type which defines node address and address pattern (index/subindex).
 */
public abstract class CANOpenField implements PlcField, Serializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?:(0[xX](?<indexHex>[0-9a-fA-F]+))|(?<index>\\d+))/(?:(0[xX](?<subIndexHex>[0-9a-fA-F]+))|(?<subIndex>\\d+)):(?<canDataType>\\w+)(\\[(?<numberOfElements>\\d)])?");
    public static final Pattern NODE_PATTERN = Pattern.compile("(?<nodeId>\\d+)");

    private final int nodeId;

    public CANOpenField(int nodeId) {
        this.nodeId = nodeId;
    }

    public int getNodeId() {
        return nodeId;
    }

    public static CANOpenField of(String addressString) throws PlcInvalidFieldException {
        if (CANOpenSDOField.matches(addressString)) {
            return CANOpenSDOField.of(addressString);
        } else if (CANOpenPDOField.matches(addressString)) {
            return CANOpenPDOField.of(addressString);
        } else if (CANOpenNMTField.matches(addressString)) {
            return CANOpenNMTField.of(addressString);
        } else if (CANOpenHeartbeatField.matches(addressString)) {
            return CANOpenHeartbeatField.of(addressString);
        }

        throw new PlcInvalidFieldException("Unable to parse address: " + addressString);
    }

}
