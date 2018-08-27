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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.messages.items.PlcReadRequestItem;
import org.apache.plc4x.java.api.messages.items.PlcReadResponseItem;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

/**
 * Response to a {@link PlcReadRequest}.
 * Contains the values read from the PLC but untyped.
 * <p>
 * Values are extracted using the {@link PlcReadRequestItem}s that were send in the read request.
 * <p>
 * If only a variables of one type are requested it is better to use
 * {@link org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest} which leads to a
 * {@link org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse}.
 */
public class PlcReadResponse extends PlcResponse<PlcReadRequest, PlcReadResponseItem<?>, PlcReadRequestItem<?>> {

    public PlcReadResponse(PlcReadRequest request, PlcReadResponseItem<?> responseItems) {
        super(request, Collections.singletonList(responseItems));
    }

    public PlcReadResponse(PlcReadRequest request, List<? extends PlcReadResponseItem<?>> responseItems) {
        super(request, responseItems);
    }

    @SuppressWarnings("unchecked")
    public <T> Optional<PlcReadResponseItem<T>> getValue(PlcReadRequestItem<T> item) {
        return (Optional) super.getValue(item);
    }

    @Override
    public String toString() {
        return "PlcReadResponse{} " + super.toString();
    }
}
