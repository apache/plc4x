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
package org.apache.plc4x.java.api.exceptions;

public class PlcProtocolPayloadTooBigException extends PlcProtocolException {

    private String protocolName;
    private int maxSize;
    private int actualSize;
    private transient Object payload;

    public PlcProtocolPayloadTooBigException(String protocolName, int maxSize, int actualSize, Object payload) {
        super("Payload for protocol '" + protocolName + "' with size " + actualSize +
            " exceeded allowed maximum of " + maxSize);
        this.protocolName = protocolName;
        this.maxSize = maxSize;
        this.actualSize = actualSize;
        this.payload = payload;
    }

    public String getProtocolName() {
        return protocolName;
    }

    public int getMaxSize() {
        return maxSize;
    }

    public int getActualSize() {
        return actualSize;
    }

    public Object getPayload() {
        return payload;
    }

}
