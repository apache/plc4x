/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.spi.optimizer;

import io.vavr.control.Either;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.DefaultPlcWriteResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;

public abstract class BaseOptimizer {

    protected List<PlcReadRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        return Collections.singletonList(readRequest);
    }

    protected PlcReadResponse processReadResponses(PlcReadRequest readRequest, Map<PlcReadRequest, Either<PlcReadResponse, Exception>> readResponses) {
        Map<String, ResponseItem<PlcValue>> tags = new HashMap<>();
        for (Map.Entry<PlcReadRequest, Either<PlcReadResponse, Exception>> requestsEntries : readResponses.entrySet()) {
            PlcReadRequest curRequest = requestsEntries.getKey();
            Either<PlcReadResponse, Exception> readResponse = requestsEntries.getValue();
            for (String tagName : curRequest.getTagNames()) {
                if (readResponse.isLeft()) {
                    PlcReadResponse subReadResponse = readResponse.getLeft();
                    PlcResponseCode responseCode = subReadResponse.getResponseCode(tagName);
                    PlcValue value = subReadResponse.getAsPlcValue().getValue(tagName);
                    tags.put(tagName, new ResponseItem<>(responseCode, value));
                } else {
                    tags.put(tagName, new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest, tags);
    }

    protected List<PlcWriteRequest> processWriteRequest(PlcWriteRequest writeRequest, DriverContext driverContext) {
        return Collections.singletonList(writeRequest);
    }

    protected PlcWriteResponse processWriteResponses(PlcWriteRequest writeRequest,
                                                     Map<PlcWriteRequest, Either<PlcWriteResponse, Exception>> writeResponses) {
        Map<String, PlcResponseCode> tags = new HashMap<>();
        for (Map.Entry<PlcWriteRequest, Either<PlcWriteResponse, Exception>> requestsEntries : writeResponses.entrySet()) {
            PlcWriteRequest subWriteRequest = requestsEntries.getKey();
            Either<PlcWriteResponse, Exception> writeResponse = requestsEntries.getValue();
            for (String tagName : subWriteRequest.getTagNames()) {
                if (writeResponse.isLeft()) {
                    PlcWriteResponse subWriteResponse = writeResponse.getLeft();
                    tags.put(tagName, subWriteResponse.getResponseCode(tagName));
                } else {
                    tags.put(tagName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }
        return new DefaultPlcWriteResponse(writeRequest, tags);
    }

    protected List<PlcSubscriptionRequest> processSubscriptionRequest(PlcSubscriptionRequest subscriptionRequest,
                                                                      DriverContext driverContext) {
        return Collections.singletonList(subscriptionRequest);
    }

    protected PlcSubscriptionResponse processSubscriptionResponses(PlcSubscriptionRequest subscriptionRequest,
                                                                   Map<PlcSubscriptionRequest, Either<PlcSubscriptionResponse, Exception>> subscriptionResponses) {
        // TODO: Implement
        return null;
    }

    protected List<PlcUnsubscriptionRequest> processUnsubscriptionRequest(PlcUnsubscriptionRequest unsubscriptionRequest,
                                                                          DriverContext driverContext) {
        return Collections.singletonList(unsubscriptionRequest);
    }

    protected PlcUnsubscriptionResponse processUnsubscriptionResponses(PlcUnsubscriptionRequest unsubscriptionRequest,
                                                                       Map<PlcUnsubscriptionRequest, Either<PlcUnsubscriptionResponse, Exception>> unsubscriptionResponses) {
        // TODO: Implement
        return null;
    }

    public CompletableFuture<PlcReadResponse> optimizedRead(PlcReadRequest readRequest, Plc4xProtocolBase<?> reader) {
        List<PlcReadRequest> subRequests = processReadRequest(readRequest, reader.getDriverContext());
        return send(readRequest, subRequests, reader::read, response -> processReadResponses(readRequest, response));
    }

    public CompletableFuture<PlcWriteResponse> optimizedWrite(PlcWriteRequest writeRequest, Plc4xProtocolBase<?> writer) {
        List<PlcWriteRequest> subRequests = processWriteRequest(writeRequest, writer.getDriverContext());
        return send(writeRequest, subRequests, writer::write, response -> processWriteResponses(writeRequest, response));
    }

    public CompletableFuture<PlcSubscriptionResponse> optimizedSubscribe(
        PlcSubscriptionRequest subscriptionRequest, Plc4xProtocolBase<?> subscriber) {
        List<PlcSubscriptionRequest> subRequests = processSubscriptionRequest(subscriptionRequest, subscriber.getDriverContext());
        return send(subscriptionRequest, subRequests, subscriber::subscribe, response -> processSubscriptionResponses(subscriptionRequest, response));
    }

    public CompletableFuture<PlcUnsubscriptionResponse> optimizedUnsubscribe(
        PlcUnsubscriptionRequest unsubscriptionRequest, Plc4xProtocolBase<?> subscriber) {
        List<PlcUnsubscriptionRequest> subRequests = processUnsubscriptionRequest(unsubscriptionRequest, subscriber.getDriverContext());
        return send(unsubscriptionRequest, subRequests, subscriber::unsubscribe, response -> processUnsubscriptionResponses(unsubscriptionRequest, response));
    }

    private <REQ extends PlcRequest, RES extends PlcResponse> CompletableFuture<RES> send(
        REQ originalRequest,
        List<REQ> requests,
        Function<REQ, CompletableFuture<RES>> sender,
        Function<Map<REQ, Either<RES, Exception>>, RES> responseProcessor) {
        // If this send has only one sub-request and this matches the original one, don't do any special handling
        // and just forward the request to the normal sending method.
        if ((requests.size() == 1) && (requests.get(0) == originalRequest)) {
            return sender.apply(requests.get(0));
        }
        // If at least one sub request is requested, split up each tag request into a separate sub-request
        // And have the reader process each one independently. After the last sub-request is finished,
        // Merge the results back together.
        if (requests.isEmpty()) {
            return CompletableFuture.completedFuture(responseProcessor.apply(Collections.emptyMap()));
        }

        // Create a new future which will be used to return the aggregated response back to the application.
        CompletableFuture<RES> parentFuture = new CompletableFuture<>();

        // Create one sub-request for every single tag and store the futures in a map.
        Map<REQ, CompletableFuture<RES>> subFutures = new HashMap<>();
        for (REQ subRequest : requests) {
            subFutures.put(subRequest, sender.apply(subRequest));
        }

        // As soon as all sub-futures are done, merge the individual responses back to one big response.
        CompletableFuture.allOf(subFutures.values().toArray(new CompletableFuture[0])).handle((aVoid, t) -> {
            if (t != null) {
                parentFuture.completeExceptionally(t);
            }
            Map<REQ, Either<RES, Exception>> results = new HashMap<>();
            for (Map.Entry<REQ, CompletableFuture<RES>> subFutureEntry : subFutures.entrySet()) {
                REQ subRequest = subFutureEntry.getKey();
                CompletableFuture<RES> subFuture = subFutureEntry.getValue();
                try {
                    final RES subResponse = subFuture.get();
                    results.put(subRequest, Either.left(subResponse));
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    results.put(subRequest, Either.right(new Exception("Something went wrong")));
                } catch (Exception e) {
                    results.put(subRequest, Either.right(new Exception("Something went wrong")));
                }
            }
            RES response = responseProcessor.apply(results);
            parentFuture.complete(response);
            return Void.TYPE;
        }).exceptionally(throwable -> {
            // TODO: If would be cool if we could still process all of the successful ones ...
            parentFuture.completeExceptionally(throwable);
            return null;
        });
        return parentFuture;
    }

}
