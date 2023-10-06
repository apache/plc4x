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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

public class PlcStruct extends PlcValueAdapter {

    private final Map<String, PlcValue> map;

    @Override
    public Object getObject() {
        return map;
    }

    public PlcStruct(Map<String, PlcValue> map) {
        this.map = Collections.unmodifiableMap(map);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.Struct;
    }

    @Override
    public int getLength() {
        return 1;
    }

    @Override
    public boolean isStruct() {
        return true;
    }

    @Override
    public Set<String> getKeys() {
        return map.keySet();
    }

    @Override
    public boolean hasKey(String key) {
        return map.containsKey(key);
    }

    @Override
    public PlcValue getValue(String key) {
        return map.get(key);
    }

    @Override
    public Map<String, ? extends PlcValue> getStruct() {
        return map;
    }

    @Override
    public String toString() {
        return "{" + map.entrySet().stream().map(entry -> String.format("\"%s\": %s", entry.getKey(), entry.getValue())).collect(Collectors.joining(",")) + "}";
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcStruct");
        for (Map.Entry<String, PlcValue> entry : map.entrySet()) {
            String tagName = entry.getKey();
            writeBuffer.pushContext(tagName);
            PlcValue tagValue = entry.getValue();
            if (!(tagValue instanceof Serializable)) {
                throw new PlcRuntimeException("Error serializing. List item doesn't implement XmlSerializable");
            }
            ((Serializable) tagValue).serialize(writeBuffer);
            writeBuffer.pushContext(tagName);
        }
        writeBuffer.popContext("PlcStruct");
    }

}
