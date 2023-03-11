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

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

public class PlcList extends PlcValueAdapter {

    private final List<PlcValue> listItems;

    public PlcList() {
        listItems = new ArrayList<>();
    }

    public PlcList(List<PlcValue> listItems) {
        List<PlcValue> safelist = listItems.stream().map(plcValue ->
            // to avoid unwrapped list cause of type erasure
            plcValue
        ).collect(Collectors.toList());
        this.listItems = Collections.unmodifiableList(safelist);
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.List;
    }

    public void add(PlcValue value) {
        listItems.add(value);
    }

    @Override
    public Object getObject() {
        return getList();
    }

    @Override
    public boolean isList() {
        return true;
    }

    @Override
    public int getLength() {
        return listItems.size();
    }

    @Override
    public PlcValue getIndex(int i) {
        return listItems.get(i);
    }

    @Override
    public List<PlcValue> getList() {
        return listItems;
    }

    @Override
    public String toString() {
        return "[" + listItems.stream().map(PlcValue::toString).collect(Collectors.joining(",")) + "]";
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcList");
        for (PlcValue listItem : listItems) {
            if (!(listItem instanceof Serializable)) {
                throw new PlcRuntimeException("Error serializing. List item doesn't implement XmlSerializable");
            }
            ((Serializable) listItem).serialize(writeBuffer);
        }
        writeBuffer.popContext("PlcList");
    }

}
