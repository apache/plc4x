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

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcStruct;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcBrowseResponse implements PlcBrowseResponse, Serializable {

    private final PlcBrowseRequest request;

    private final Map<String, PlcResponseCode> responseCodes;

    private final Map<String, List<PlcBrowseItem>> values;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcBrowseResponse(@JsonProperty("request") PlcBrowseRequest request, @JsonProperty("responseCodes") Map<String, PlcResponseCode> responseCodes, @JsonProperty("values") Map<String, List<PlcBrowseItem>> values) {
        this.request = request;
        this.responseCodes = responseCodes;
        this.values = values;
    }

    @Override
    public PlcBrowseRequest getRequest() {
        return request;
    }

    public LinkedHashSet<String> getQueryNames() {
        return request.getQueryNames();
    }

    @Override
    public PlcResponseCode getResponseCode(String queryName) {
        return responseCodes.get(queryName);
    }

    @Override
    public List<PlcBrowseItem> getValues(String queryName) {
        return values.get(queryName);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcBrowseResponse");
        writeBuffer.popContext("PlcBrowseResponse");
    }

}
