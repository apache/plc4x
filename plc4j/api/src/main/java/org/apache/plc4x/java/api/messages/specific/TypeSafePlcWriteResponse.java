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

import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.messages.items.WriteResponseItem;

import java.util.List;
import java.util.Optional;

public class TypeSafePlcWriteResponse<T> extends PlcWriteResponse {

    public TypeSafePlcWriteResponse(TypeSafePlcWriteRequest<T> request, WriteResponseItem<T> responseItem) {
        super(request, responseItem);
    }

    @SuppressWarnings("unchecked")
    public TypeSafePlcWriteResponse(TypeSafePlcWriteRequest<T> request, List<WriteResponseItem<T>> responseItems) {
        super(request, responseItems);
    }

    @Override
    @SuppressWarnings("unchecked")
    public TypeSafePlcWriteRequest<T> getRequest() {
        return (TypeSafePlcWriteRequest<T>) super.getRequest();
    }

    @Override
    @SuppressWarnings("unchecked")
    public List<? extends WriteResponseItem<T>> getResponseItems() {
        return (List<WriteResponseItem<T>>) super.getResponseItems();
    }

    @Override
    @SuppressWarnings("unchecked")
    public Optional<WriteResponseItem<T>> getResponseItem() {
        return (Optional<WriteResponseItem<T>>) super.getResponseItem();
    }

    @SuppressWarnings("unchecked")
    public static TypeSafePlcWriteResponse of(PlcWriteResponse plcWriteResponse) {
        if (plcWriteResponse instanceof TypeSafePlcWriteResponse) {
            return (TypeSafePlcWriteResponse) plcWriteResponse;
        }
        if (plcWriteResponse.getRequest() instanceof TypeSafePlcWriteRequest) {
            return new TypeSafePlcWriteResponse((TypeSafePlcWriteRequest) plcWriteResponse.getRequest(), plcWriteResponse.getResponseItems());
        }
        Class<?> referenceType = null;
        for (WriteResponseItem<?> writeResponseItem : plcWriteResponse.getResponseItems()) {
            Class<?> foundDataType = writeResponseItem.getRequestItem().getDatatype();
            if (referenceType == null) {
                referenceType = foundDataType;
            }
            if (referenceType != foundDataType) {
                throw new IllegalArgumentException("invalid types found " + foundDataType + ". Required " + referenceType);
            }
        }
        if (referenceType == null) {
            referenceType = Object.class;
        }
        return new TypeSafePlcWriteResponse(new TypeSafePlcWriteRequest(referenceType, plcWriteResponse.getRequest()), plcWriteResponse.getResponseItems());
    }
}
