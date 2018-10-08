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
package org.apache.plc4x.java.base.messages;

import java.util.Objects;
import java.util.concurrent.CompletableFuture;

/**
 * Helper mostly used internally to split up big requests into multiple sub-requests.
 *
 * @param <T> type of request.
 * @param <R> type of response.
 */
public class PlcRequestContainer<T extends InternalPlcRequest, R extends InternalPlcResponse> implements PlcProtocolMessage {

    private final T request;
    private final CompletableFuture<R> responseFuture;

    public PlcRequestContainer(T request, CompletableFuture<R> responseFuture) {
        Objects.requireNonNull(request, "Request must not be null");
        Objects.requireNonNull(responseFuture, "Response future must not be null");
        this.request = request;
        this.responseFuture = responseFuture;
    }

    public T getRequest() {
        return request;
    }

    public CompletableFuture<R> getResponseFuture() {
        return responseFuture;
    }

    /**
     * {@link PlcRequestContainer} objects don't have parents.
     *
     * @return null
     */
    @Override
    public PlcProtocolMessage getParent() {
        return null;
    }

    @Override
    public String toString() {
        return "PlcRequestContainer{" +
            "request=" + request +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof PlcRequestContainer)) {
            return false;
        }
        PlcRequestContainer<?, ?> that = (PlcRequestContainer<?, ?>) o;
        return Objects.equals(request, that.request) &&
            Objects.equals(responseFuture, that.responseFuture);
    }

    @Override
    public int hashCode() {
        return Objects.hash(request, responseFuture);
    }

}
