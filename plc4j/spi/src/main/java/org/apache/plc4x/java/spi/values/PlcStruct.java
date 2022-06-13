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
package org.apache.plc4x.java.spi.values;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcStruct extends PlcValueAdapter {

    private final Map<String, PlcValue> map;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcStruct(@JsonProperty("map") Map<String, PlcValue> map) {
        this.map = Collections.unmodifiableMap(map);
    }

    @Override
    @JsonIgnore
    public int getLength() {
        return 1;
    }

    @Override
    @JsonIgnore
    public boolean isStruct() {
        return true;
    }

    @Override
    @JsonIgnore
    public Set<String> getKeys() {
        return map.keySet();
    }

    @Override
    @JsonIgnore
    public boolean hasKey(String key) {
        return map.containsKey(key);
    }

    @Override
    @JsonIgnore
    public PlcValue getValue(String key) {
        return map.get(key);
    }

    @Override
    @JsonIgnore
    public Map<String, ? extends PlcValue> getStruct() {
        return map;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return "{" + map.entrySet().stream().map(entry -> String.format("\"%s\": %s", entry.getKey(), entry.getValue())).collect(Collectors.joining(",")) + "}";
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcStruct");
        for (Map.Entry<String, PlcValue> entry : map.entrySet()) {
            String fieldName = entry.getKey();
            writeBuffer.pushContext(fieldName);
            PlcValue fieldValue = entry.getValue();
            if (!(fieldValue instanceof Serializable)) {
                throw new PlcRuntimeException("Error serializing. List item doesn't implement XmlSerializable");
            }
            ((Serializable) fieldValue).serialize(writeBuffer);
            writeBuffer.pushContext(fieldName);
        }
        writeBuffer.popContext("PlcStruct");
    }

}
