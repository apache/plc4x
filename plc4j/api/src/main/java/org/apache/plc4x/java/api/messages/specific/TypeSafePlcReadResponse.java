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
package org.apache.plc4x.java.api.messages.specific;

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.PlcReadResponseItem;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

public class TypeSafePlcReadResponse<T> extends PlcReadResponse {

    public TypeSafePlcReadResponse(TypeSafePlcReadRequest<T> request, PlcReadResponseItem<T> responseItem) {
        super(request, responseItem);
        Objects.requireNonNull(request, "Request must not be null");
        checkList(responseItem.getValues(), request.getDataType());
    }

    public TypeSafePlcReadResponse(TypeSafePlcReadRequest<T> request, List<PlcReadResponseItem<T>> responseItems) {
        super(request, responseItems);
        Objects.requireNonNull(responseItems, "Request items on " + request + " must not be null");
        for (PlcReadResponseItem<T> responseItem : responseItems) {
            checkList(responseItem.getValues(), request.getDataType());
        }
    }

    @Override
    @SuppressWarnings("unchecked")
    public TypeSafePlcReadRequest<T> getRequest() {
        return (TypeSafePlcReadRequest<T>) super.getRequest();
    }

    @Override
    @SuppressWarnings("unchecked")
    public List<? extends PlcReadResponseItem<T>> getResponseItems() {
        return (List<? extends PlcReadResponseItem<T>>) super.getResponseItems();
    }

    @Override
    @SuppressWarnings("unchecked")
    public Optional<PlcReadResponseItem<T>> getResponseItem() {
        return (Optional<PlcReadResponseItem<T>>) super.getResponseItem();
    }

    public static <T> TypeSafePlcReadResponse<T> of(PlcReadResponse plcReadResponse, Class<T> clazz) {
        Objects.requireNonNull(plcReadResponse, "PlcReadResponse must not be null");
        Objects.requireNonNull(clazz, "Class must not be null");
        if (plcReadResponse instanceof TypeSafePlcReadResponse) {
            @SuppressWarnings("unchecked")
            TypeSafePlcReadResponse<T> typeSafePlcReadResponse = (TypeSafePlcReadResponse<T>) plcReadResponse;
            Class type = typeSafePlcReadResponse.getRequest().getDataType();
            if (type != clazz) {
                throw new IllegalArgumentException("Expected type " + clazz + " doesn't match found type " + type);
            }
            return typeSafePlcReadResponse;
        }
        @SuppressWarnings("unchecked")
        List<PlcReadResponseItem<T>> responseItems = (List<PlcReadResponseItem<T>>) plcReadResponse.getResponseItems();
        Objects.requireNonNull(responseItems, "Response items on " + plcReadResponse + " must not be null");
        if (plcReadResponse.getRequest() instanceof TypeSafePlcReadRequest) {
            @SuppressWarnings("unchecked")
            TypeSafePlcReadRequest<T> typeSafePlcReadRequest = (TypeSafePlcReadRequest<T>) plcReadResponse.getRequest();
            Class type = typeSafePlcReadRequest.getDataType();
            if (type != clazz) {
                throw new IllegalArgumentException("Expected type " + clazz + " doesn't match found type " + type);
            }
            return new TypeSafePlcReadResponse<>(typeSafePlcReadRequest, responseItems);
        }
        for (PlcReadResponseItem<?> responseItem : responseItems) {
            checkList(responseItem.getValues(), clazz);
        }
        TypeSafePlcReadRequest<T> request = new TypeSafePlcReadRequest<>(clazz, plcReadResponse.getRequest());
        return new TypeSafePlcReadResponse<>(request, responseItems);
    }

    private static void checkList(List<?> list, Class<?> type) {
        Objects.requireNonNull(list, "List must not be null");
        Objects.requireNonNull(type, "Type must not be null");
        for (Object o : list) {
            if (!type.isAssignableFrom(o.getClass())) {
                throw new IllegalArgumentException("Unexpected data type " + o.getClass() + " on readRequestItem. Expected " + type);
            }
        }
    }

    @Override
    public String toString() {
        return "TypeSafePlcReadResponse{} " + super.toString();
    }
}
