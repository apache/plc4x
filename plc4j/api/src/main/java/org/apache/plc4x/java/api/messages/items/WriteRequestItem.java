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

import java.lang.reflect.Array;
import java.util.Optional;

public class WriteRequestItem<T> {

    private final Class<T> datatype;

    private final Address address;

    private final T[] values;

    private WriteResponseItem<T> responseItem;

    @SuppressWarnings("unchecked")
    public WriteRequestItem(Class<T> datatype, Address address, T value) {
        this.datatype = datatype;
        this.address = address;
        this.values = (T[]) Array.newInstance(datatype, 1);
        this.values[0] = value;
        responseItem = null;
    }

    public WriteRequestItem(Class<T> datatype, Address address, T[] values) {
        this.datatype = datatype;
        this.address = address;
        this.values = values;
        responseItem = null;
    }

    public Class<T> getDatatype() {
        return datatype;
    }

    public Address getAddress() {
        return address;
    }

    public T[] getValues() {
        return values;
    }

    public Optional<WriteResponseItem<T>> getResponseItem() {
        return Optional.ofNullable(responseItem);
    }

    protected void setResponseItem(WriteResponseItem<T> responseItem) {
        this.responseItem = responseItem;
    }
}
