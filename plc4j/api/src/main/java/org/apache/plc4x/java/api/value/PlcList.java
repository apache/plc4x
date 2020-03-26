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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class PlcList extends PlcValueAdapter {

    private final List<PlcValue> listItems;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public PlcList(@JsonProperty("listItems") List<?> listItems) {
        List<PlcValue> safelist = listItems.stream().map(plcValue -> {
            // to avoid unwrapped list cause of type erasure
            if (plcValue instanceof PlcValue) {
                return (PlcValue) plcValue;
            } else {
                return PlcValues.of(plcValue);
            }
        }).collect(Collectors.toList());
        this.listItems = Collections.unmodifiableList(safelist);
    }

    @Override
    public Object getObject() {
        return listItems;
    }

    @Override
    @JsonIgnore
    public boolean isList() {
        return true;
    }

    @Override
    @JsonIgnore
    public int getLength() {
        return listItems.size();
    }

    @Override
    @JsonIgnore
    public PlcValue getIndex(int i) {
        return listItems.get(i);
    }

    @Override
    @JsonIgnore
    public List<? extends PlcValue> getList() {
        return listItems;
    }

    @Override
    @JsonIgnore
    public String toString() {
        return "[" + listItems.stream().map(PlcValue::toString).collect(Collectors.joining(",")) + "]";
    }

}
