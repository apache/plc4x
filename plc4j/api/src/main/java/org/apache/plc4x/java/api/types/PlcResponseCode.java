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
package org.apache.plc4x.java.api.types;

import java.util.HashMap;
import java.util.Map;

public enum PlcResponseCode {

    /**
     * Everything went ok.
     */
    OK((short) 0x01),
    /**
     * The requested ressource could not be found on the target device.
     * (The address was syntactically correct, but the item it addressed simply didn't exist)
     */
    NOT_FOUND((short) 0x02),
    /**
     * The remote device denied access to the requested resource.
     * (Possibly remote access is disabled, or an authentication is required)
     */
    ACCESS_DENIED((short) 0x03),
    /**
     * The resource address was syntactically invalid.
     */
    INVALID_ADDRESS((short) 0x04),
    /**
     * The requested datatype does not exist or was not compatible with the requested ressource.
     */
    INVALID_DATATYPE((short) 0x06),
    INVALID_DATA((short) 0x07),
    /**
     * Something went wrong internally in the driver logic.
     * (This is most probably a PLC4X bug)
     */
    INTERNAL_ERROR((short) 0x08),
    /**
     * The remote device is currently unable to process the request due to overload.
     */
    REMOTE_BUSY((short) 0x09),
    /**
     * Something went wrong on the remote side.
     */
    REMOTE_ERROR((short) 0x0A),
    /**
     * The requested resource uses a feature of the driver which has not been implemented.
     */
    UNSUPPORTED((short) 0x0B),
    /**
     * Indicates a response is pending.
     */
    RESPONSE_PENDING((short) 0x0C),
    /**
     * The request and the connection are both healthy, but the data this
     * tag points at is not yet available — typically because an asynchronous
     * source (cyclic IO, subscription stream, lazy cache) has not yet
     * delivered its first value. The caller should retry shortly; this is
     * not a permanent failure.
     *
     * <p>Distinct from {@link #RESPONSE_PENDING}, which signals a
     * request/response is in flight to a remote endpoint. {@code NOT_READY}
     * is for push-model and cache-backed reads where there is no in-flight
     * request to wait on — only the next periodic delivery.</p>
     */
    NOT_READY((short) 0x0D),
    /**
     * The request shape is valid, but a value carried by it falls outside
     * the range the target accepts — e.g. a numeric value too large for
     * the destination datatype, an address beyond the device's declared
     * extent, or a polling interval below the protocol's wire-cycle floor.
     * Distinct from {@link #INVALID_DATA} (which means "malformed") and
     * {@link #INVALID_ADDRESS} (which means "no such address exists").
     */
    OUT_OF_RANGE((short) 0x0E);

    private static final Map<Short, PlcResponseCode> map;

    static {
        map = new HashMap<>();
        for (PlcResponseCode value : PlcResponseCode.values()) {
            map.put((short) value.getValue(), value);
        }
    }

    private final short value;

    PlcResponseCode(short value) {
        this.value = value;
    }

    public short getValue() {
        return value;
    }

    public static PlcResponseCode enumForValue(short value) {
        return map.get(value);
    }

    public static Boolean isDefined(short value) {
        return map.containsKey(value);
    }
}
