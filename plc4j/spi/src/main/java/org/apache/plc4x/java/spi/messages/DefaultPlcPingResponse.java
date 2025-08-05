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

package org.apache.plc4x.java.spi.messages;

import org.apache.plc4x.java.api.messages.PlcPingRequest;
import org.apache.plc4x.java.api.messages.PlcPingResponse;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

public class DefaultPlcPingResponse implements PlcPingResponse, Serializable {

    private final PlcPingRequest request;
    private final PlcResponseCode responseCode;

    public DefaultPlcPingResponse(PlcPingRequest request, PlcResponseCode responseCode) {
        this.request = request;
        this.responseCode = responseCode;
    }

    @Override
    public PlcPingRequest getRequest() {
        return request;
    }

    @Override
    public PlcResponseCode getResponseCode() {
        return responseCode;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcPingResponse");

        writeBuffer.pushContext("request");
        if(request instanceof Serializable) {
            ((Serializable) request).serialize(writeBuffer);
        }
        writeBuffer.popContext("request");

        writeBuffer.popContext("PlcPingResponse");
    }

}
