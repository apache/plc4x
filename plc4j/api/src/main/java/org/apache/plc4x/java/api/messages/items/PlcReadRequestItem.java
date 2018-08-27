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

import java.util.Objects;

/**
 * Encapsulats one {@link RequestItem} that could be read multiple times, i.e., from the given PlcField on
 * {@link PlcReadRequestItem#size} number of Items with equal datatype are read.
 *
 * Thus,
 * <pre>
 *     new PlcReadRequestItem(Int.class, adress, 5)
 * </pre>
 * basically reads 5 consecutive integers starting at the given {@link PlcField}.
 *
 * @param <T> Generic Type of expected Datatype.
 */
public class PlcReadRequestItem<T> extends RequestItem<T> {

    private final int size;

    public PlcReadRequestItem(Class<T> datatype, PlcField field) {
        super(datatype, field);
        this.size = 1;
    }

    public PlcReadRequestItem(Class<T> datatype, PlcField field, int size) {
        super(datatype, field);
        this.size = size;
    }

    public int getSize() {
        return size;
    }

    @Override
    public String toString() {
        return "PlcReadRequestItem{" +
            "size=" + size +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcReadRequestItem)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        PlcReadRequestItem<?> that = (PlcReadRequestItem<?>) o;
        return size == that.size;
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), size);
    }
}
