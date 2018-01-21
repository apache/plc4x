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
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;

import java.util.List;
import java.util.Optional;

public class TypeSafePlcReadResponse<T> extends PlcReadResponse {

    public TypeSafePlcReadResponse(PlcReadResponse plcReadResponse) {
        super(plcReadResponse.getRequest(), plcReadResponse.getResponseItems());
        List<? extends ReadResponseItem<?>> responseItems = plcReadResponse.getResponseItems();
        Class type = null;
        for (ReadResponseItem<?> responseItem : responseItems) {
            if (!responseItem.getValues().isEmpty()) {
                type = responseItem.getValues().get(0).getClass();
                break;
            }
        }
        if (type != null) {
            for (ReadResponseItem<?> responseItem : responseItems) {
                checkList(responseItem.getValues(), type);
            }
        }
    }

    public TypeSafePlcReadResponse(TypeSafePlcReadRequest<T> request, ReadResponseItem<T> responseItem) {
        super(request, responseItem);
        checkList(responseItem.getValues(), request.getDataType());
    }

    public TypeSafePlcReadResponse(TypeSafePlcReadRequest<T> request, List<ReadResponseItem<T>> responseItems) {
        super(request, responseItems);
        for (ReadResponseItem<T> responseItem : responseItems) {
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
    public List<? extends ReadResponseItem<T>> getResponseItems() {
        return (List<? extends ReadResponseItem<T>>) super.getResponseItems();
    }

    @Override
    @SuppressWarnings("unchecked")
    public Optional<ReadResponseItem<T>> getResponseItem() {
        return (Optional<ReadResponseItem<T>>) super.getResponseItem();
    }

    private static void checkList(List<?> list, Class<?> type) {
        for (Object o : list) {
            if (!type.isAssignableFrom(o.getClass())) {
                throw new IllegalArgumentException("Unexpected data type " + o.getClass() + " on readRequestItem. Expected " + type);
            }
        }
    }
}
