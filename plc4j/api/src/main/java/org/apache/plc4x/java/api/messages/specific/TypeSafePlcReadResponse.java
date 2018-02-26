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

import java.util.List;
import java.util.Objects;
import java.util.Optional;

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;

public class TypeSafePlcReadResponse<T> extends PlcReadResponse {

    public TypeSafePlcReadResponse(TypeSafePlcReadRequest<T> request, ReadResponseItem<T> responseItem) {
        super(request, responseItem);
        Objects.requireNonNull(request, "Request must not be null");
        checkList(responseItem.getValues(), request.getDataType());
    }

    public TypeSafePlcReadResponse(TypeSafePlcReadRequest<T> request, List<ReadResponseItem<T>> responseItems) {
        super(request, responseItems);
        Objects.requireNonNull(responseItems, "Request items on " + request + " must not be null");
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

    /**
     * Cast or convert a PlcReadResponse to a TypeSafePlcReadReadResponse<T>.
     * 
     * WARNING: this is inherently a non-type-safe operation.  It was introduced
     * to serve the implementation of PlcReader.read(TypeSafePlcReadRequest<T>).
     * Additional use of it is not recommended.  This interface is subject to change.
     * 
     * @param plcReadResponse the response implicitly with items of type T
     * @return TypeSafePlcReadReadResponse<T>
     */
    @SuppressWarnings("unchecked")
    public static <T> TypeSafePlcReadResponse<T> of(PlcReadResponse plcReadResponse) {
        // BUG: there seems to be no checking that the readResponse/items
        // in fact are for type T.
        // I don't even think it's possible to do that with the current 'of()' signature
        // and plcReadResponse content.
        //
        // The only consolation is that currently, 'of()' is only really used in the
        // impl of PlcReader.read(TypeSafePlcReadRequest<T>) and that case guarantees
        // that all items are a T.  plcReadResponse isa TypeSafePlcReadResponse in
        // this case.
        //
        // Maybe we just need to doc that this conversion is unsafe and/or rename
        // the method to "unsafeOf()"? 
        // Or, if there were an AbstractPlcReader<T>, that could internally implement
        // this method without exposing this generally it, the PlcReader interface
        // could remove the default implementation of read(TypeSafePlcReadRequest<T>),
        // and protocol implementations could extend AbstractPlcReader.
        //
        // FWIW, in one case there is some checking that all of the items in a response
        // are at least of the same type.
      
        if (plcReadResponse instanceof TypeSafePlcReadResponse) {
            // Warning: no validation that everything in the response is a T.
            return (TypeSafePlcReadResponse<T>) plcReadResponse;
        }
        if (plcReadResponse.getRequest() instanceof TypeSafePlcReadRequest) {
            // Warning: no validation that everything in the response is a T.
            return new TypeSafePlcReadResponse<T>((TypeSafePlcReadRequest<T>) plcReadResponse.getRequest(), (List<ReadResponseItem<T>>) plcReadResponse.getResponseItems());
        }
        List<? extends ReadResponseItem<?>> responseItems = plcReadResponse.getResponseItems();
        Objects.requireNonNull(responseItems, "Response items on " + plcReadResponse + " must not be null");
        Class<?> type = null;
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
        if (type == null) {
            type = Object.class;
        }
        // Warning: no validation that everything in the response is a T.
        // Verified that everything in the response was the same type, but why bother?
        return new TypeSafePlcReadResponse<T>(new TypeSafePlcReadRequest<T>((Class<T>) type, plcReadResponse.getRequest()), (List<ReadResponseItem<T>>) responseItems);
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
}
