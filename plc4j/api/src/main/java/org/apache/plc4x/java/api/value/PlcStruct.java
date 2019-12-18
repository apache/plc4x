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

package org.apache.plc4x.java.api.value;

import java.util.Collections;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

public class PlcStruct extends PlcValueAdapter {

    private final Map<String, PlcValue> map;

    PlcStruct(Map<String, PlcValue> map) {
        this.map = Collections.unmodifiableMap(map);
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

}
