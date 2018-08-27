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

import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Arrays;
import java.util.List;
import java.util.Objects;

/**
 * Response to a {@link PlcReadRequestItem}.
 * Can contain a list of values if the size in {@link PlcReadRequestItem} is larger zero.
 *
 * @param <T>
 */
public class PlcReadResponseItem<T> extends ResponseItem<PlcReadRequestItem<T>> {

    private final List<T> values;

    public PlcReadResponseItem(PlcReadRequestItem<T> requestItem, PlcResponseCode responseCode, List<T> values) {
        super(requestItem, responseCode);
        Objects.requireNonNull(values, "Values must not be null");
        for (T value : values) {
            if (!requestItem.getDatatype().isAssignableFrom(value.getClass())) {
                throw new IllegalArgumentException("Datatype of " + value + " doesn't match required datatype of " + requestItem.getDatatype());
            }
        }
        this.values = values;
    }

    @SafeVarargs
    public PlcReadResponseItem(PlcReadRequestItem<T> requestItem, PlcResponseCode responseCode, T... values) {
        this(requestItem, responseCode, Arrays.asList(values));
    }

    public List<T> getValues() {
        return values;
    }

    @Override
    public String toString() {
        return "PlcReadResponseItem{" +
            "values=" + values +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcReadResponseItem)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        PlcReadResponseItem<?> that = (PlcReadResponseItem<?>) o;
        return Objects.equals(values, that.values);
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), values);
    }
}
