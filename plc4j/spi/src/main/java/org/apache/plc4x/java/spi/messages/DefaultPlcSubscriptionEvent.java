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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;

import java.time.Instant;
import java.util.Collection;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcSubscriptionEvent extends DefaultPlcReadResponse implements PlcSubscriptionEvent {

    public final Instant timestamp;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcSubscriptionEvent(@JsonProperty("timestamp") Instant timestamp,
                                       @JsonProperty("fields") Map<String, ResponseItem<PlcValue>> fields) {
        super(null, fields);
        this.timestamp = timestamp;
    }

    @Override
    @JsonIgnore
    public Collection<String> getFieldNames() {
        return getValues().keySet();
    }

    @Override
    @JsonIgnore
    public PlcField getField(String name) {
        throw new UnsupportedOperationException("getField('" + name + "') not supported on " + this.getClass());
    }

    @Override
    public Instant getTimestamp() {
        return timestamp;
    }

}
