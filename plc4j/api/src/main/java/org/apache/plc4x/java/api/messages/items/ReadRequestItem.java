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

import org.apache.plc4x.java.api.model.Address;

import java.util.Objects;

/**
 * Encapsulats one {@link RequestItem} that could be read multiple times, i.e., from the given Address on
 * {@link ReadRequestItem#size} number of Items with equal datatype are read.
 *
 * Thus,
 * <pre>
 *     new ReadRequestItem(Int.class, adress, 5)
 * </pre>
 * basically reads 5 consecutive integers starting at the given {@link Address}.
 *
 * @param <T> Generic Type of expected Datatype.
 */
public class ReadRequestItem<T> extends RequestItem<T> {

    private final int size;

    public ReadRequestItem(Class<T> datatype, Address address) {
        super(datatype, address);
        this.size = 1;
    }

    public ReadRequestItem(Class<T> datatype, Address address, int size) {
        super(datatype, address);
        this.size = size;
    }

    public int getSize() {
        return size;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof ReadRequestItem)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        ReadRequestItem<?> that = (ReadRequestItem<?>) o;
        // TODO 01.18.18 jf: we should also call the comparison of super at least otherwise this can lead to unwanted behavior.
        // Perhaps we should generate a UUID or something for each ReadRequest to have a good equality comparison
        return size == that.size;
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), size);
    }
}
