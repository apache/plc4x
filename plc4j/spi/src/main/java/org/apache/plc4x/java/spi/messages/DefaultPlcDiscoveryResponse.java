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
package org.apache.plc4x.java.spi.messages;

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.*;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcDiscoveryResponse implements PlcDiscoveryResponse, Serializable {

    private final PlcDiscoveryRequest request;
    private final PlcResponseCode responseCode;
    private final List<PlcDiscoveryItem> values;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcDiscoveryResponse(@JsonProperty("request") PlcDiscoveryRequest request,
                                       @JsonProperty("responseCode") PlcResponseCode responseCode,
                                       @JsonProperty("values") List<PlcDiscoveryItem> values) {
        this.request = request;
        this.responseCode = responseCode;
        this.values = values;
    }

    @Override
    public PlcDiscoveryRequest getRequest() {
        return request;
    }

    public PlcResponseCode getResponseCode() {
        return responseCode;
    }

    public List<PlcDiscoveryItem> getValues() {
        return values;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext("PlcDiscoveryResponse");

        // TODO: Implement

        writeBuffer.popContext("PlcDiscoveryResponse");
    }

}
