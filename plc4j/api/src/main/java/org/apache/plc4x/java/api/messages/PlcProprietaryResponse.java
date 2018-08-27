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

import org.apache.plc4x.java.api.messages.items.ResponseItem;
import org.apache.plc4x.java.api.types.PlcResponseCode;

import java.util.Collections;
import java.util.Objects;

/**
 * This response can be used to pass proprietary protocol specific messages from the underlying layers.
 */
public class PlcProprietaryResponse<CUSTOM_RESPONSE> extends PlcResponse<PlcProprietaryRequest, PlcProprietaryResponse.DummyResponseItem, PlcProprietaryRequest.DummyRequestItem> {

    private final CUSTOM_RESPONSE response;

    public PlcProprietaryResponse(PlcProprietaryRequest request, CUSTOM_RESPONSE response) {
        super(request, Collections.emptyList());
        this.response = response;
    }

    public CUSTOM_RESPONSE getResponse() {
        return response;
    }

    protected static class DummyResponseItem extends ResponseItem<PlcProprietaryRequest.DummyRequestItem> {

        public DummyResponseItem() {
            super(new PlcProprietaryRequest.DummyRequestItem(), PlcResponseCode.OK);
        }
    }

    @Override
    public String toString() {
        return "PlcProprietaryResponse{" +
            "response=" + response +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcProprietaryResponse)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        PlcProprietaryResponse<?> that = (PlcProprietaryResponse<?>) o;
        return Objects.equals(response, that.response);
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), response);
    }
}
