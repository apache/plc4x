/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
*/
package org.apache.plc4x.java.api.messages.items;

import org.apache.plc4x.java.api.model.PlcField;

import java.util.*;

public class PlcWriteRequestItem<T> extends RequestItem<T> {

    private final List<T> values;

    @SafeVarargs
    public PlcWriteRequestItem(Class<T> dataType, PlcField field, T... values) {
        super(dataType, field);
        Objects.requireNonNull(values, "Values must not be null");
        List<T> checkedList = Collections.checkedList(new ArrayList<>(), dataType);
        checkedList.addAll(Arrays.asList(values));
        this.values = checkedList;
    }

    public List<T> getValues() {
        return values;
    }

    public int getSize() {
        return values.size();
    }

    @Override
    public String toString() {
        return "PlcWriteRequestItem{" +
            "values=" + values +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcWriteRequestItem)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        PlcWriteRequestItem<?> that = (PlcWriteRequestItem<?>) o;
        return Objects.equals(values, that.values);
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), values);
    }
}
