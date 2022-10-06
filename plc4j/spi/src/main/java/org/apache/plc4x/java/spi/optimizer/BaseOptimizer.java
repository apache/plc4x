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

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;

public abstract class BaseOptimizer {

    protected List<PlcRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        return Collections.singletonList(readRequest);
    }

    protected PlcReadResponse processReadResponses(PlcReadRequest readRequest, Map<PlcRequest, Either<PlcResponse, Exception>> readResponses) {
        Map<String, ResponseItem<PlcValue>> fields = new HashMap<>();
        for (Map.Entry<PlcRequest, Either<PlcResponse, Exception>> requestsEntries : readResponses.entrySet()) {
            PlcReadRequest curRequest = (PlcReadRequest) requestsEntries.getKey();
            Either<PlcResponse, Exception> readResponse = requestsEntries.getValue();
            for (String fieldName : curRequest.getFieldNames()) {
                if (readResponse.isLeft()) {
                    PlcReadResponse subReadResponse = (PlcReadResponse) readResponse.getLeft();
                    PlcResponseCode responseCode = subReadResponse.getResponseCode(fieldName);
                    PlcValue value = (responseCode == PlcResponseCode.OK) ?
                        subReadResponse.getAsPlcValue().getValue(fieldName) : null;
                    fields.put(fieldName, new ResponseItem<>(responseCode, value));
                } else {
                    fields.put(fieldName, new ResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                }
            }
        }
        return new DefaultPlcReadResponse(readRequest, fields);
    }

    protected List<PlcRequest> processWriteRequest(PlcWriteRequest writeRequest, DriverContext driverContext) {
        return Collections.singletonList(writeRequest);
    }

    protected PlcWriteResponse processWriteResponses(PlcWriteRequest writeRequest,
                                                     Map<PlcRequest, Either<PlcResponse, Exception>> writeResponses) {
        Map<String, PlcResponseCode> fields = new HashMap<>();
        for (Map.Entry<PlcRequest, Either<PlcResponse, Exception>> requestsEntries : writeResponses.entrySet()) {
            PlcWriteRequest subWriteRequest = (PlcWriteRequest) requestsEntries.getKey();
            Either<PlcResponse, Exception> writeResponse = requestsEntries.getValue();
            for (String fieldName : subWriteRequest.getFieldNames()) {
                if (writeResponse.isLeft()) {
                    PlcWriteResponse subWriteResponse = (PlcWriteResponse) writeResponse.getLeft();
                    fields.put(fieldName, subWriteResponse.getResponseCode(fieldName));
                } else {
                    fields.put(fieldName, PlcResponseCode.INTERNAL_ERROR);
                }
            }
        }
        return new DefaultPlcWriteResponse(writeRequest, fields);
    }

    protected List<PlcRequest> processSubscriptionRequest(PlcSubscriptionRequest subscriptionRequest,
                                                          DriverContext driverContext) {
        return Collections.singletonList(subscriptionRequest);
    }

    protected PlcSubscriptionResponse processSubscriptionResponses(PlcSubscriptionRequest subscriptionRequest,
                                                                   Map<PlcRequest, Either<PlcResponse, Exception>> subscriptionResponses) {
        // TODO: Implement
        return null;
    }

    protected List<PlcRequest> processUnsubscriptionRequest(PlcRequest unsubscriptionRequest,
                                                            DriverContext driverContext) {
        return Collections.singletonList(unsubscriptionRequest);
    }

    protected PlcUnsubscriptionResponse processUnsubscriptionResponses(PlcRequest unsubscriptionRequest,
                                                                       Map<PlcRequest, Either<PlcResponse, Exception>> unsubscriptionResponses) {
        // TODO: Implement
        return null;
    }

    public CompletableFuture<PlcReadResponse> optimizedRead(PlcReadRequest readRequest, Plc4xProtocolBase reader) {
        List<PlcRequest> subRequests = processReadRequest(readRequest, reader.getDriverContext());
        return send(readRequest, subRequests, request -> reader.read((PlcReadRequest) request),
            response -> processReadResponses(readRequest, response));
    }

    public CompletableFuture<PlcWriteResponse> optimizedWrite(PlcWriteRequest writeRequest, Plc4xProtocolBase writer) {
        List<PlcRequest> subRequests = processWriteRequest(writeRequest, writer.getDriverContext());
        return send(writeRequest, subRequests, request -> writer.write((PlcWriteRequest) request),
            response -> processWriteResponses(writeRequest, response));
    }

    public CompletableFuture<PlcSubscriptionResponse> optimizedSubscribe(
            PlcSubscriptionRequest subscriptionRequest, Plc4xProtocolBase subscriber) {
        List<PlcRequest> subRequests = processSubscriptionRequest(subscriptionRequest, subscriber.getDriverContext());
        return send(subscriptionRequest, subRequests, request -> subscriber.subscribe((PlcSubscriptionRequest) request),
            response -> processSubscriptionResponses(subscriptionRequest, response));
    }

    public CompletableFuture<PlcUnsubscriptionResponse> optmizedUnsubscribe(
            PlcUnsubscriptionRequest unsubscriptionRequest, Plc4xProtocolBase subscriber) {
        List<PlcRequest> subRequests = processUnsubscriptionRequest(unsubscriptionRequest, subscriber.getDriverContext());
        return send(unsubscriptionRequest, subRequests, request -> subscriber.unsubscribe((PlcUnsubscriptionRequest) request),
            response -> processUnsubscriptionResponses(unsubscriptionRequest, response));
    }

    private CompletableFuture send(PlcRequest originalRequest,
                                   List<? extends PlcRequest> requests,
                                   Function<PlcRequest, CompletableFuture<PlcResponse>> sender,
                                   Function<Map<PlcRequest, Either<PlcResponse, Exception>>, PlcResponse> responseProcessor) {
        // If this send has only one sub-request and this matches the original one, don't do any special handling
        // and just forward the request to the normal sending method.
        if((requests.size() == 1) && (requests.get(0) == originalRequest)) {
            return sender.apply(requests.get(0));
        }
        // If at least one sub request is requested, split up each field request into a separate sub-request
        // And have the reader process each one independently. After the last sub-request is finished,
        // Merge the results back together.
        else if (!requests.isEmpty()) {
            // Create a new future which will be used to return the aggregated response back to the application.
            CompletableFuture<PlcResponse> parentFuture = new CompletableFuture<>();

            // Create one sub-request for every single field and store the futures in a map.
            Map<PlcRequest, CompletableFuture<PlcResponse>> subFutures = new HashMap<>();
            for (PlcRequest subRequest : requests) {
                subFutures.put(subRequest, sender.apply(subRequest));
            }

            // As soon as all sub-futures are done, merge the individual responses back to one big response.
            CompletableFuture.allOf(subFutures.values().toArray(new CompletableFuture[0])).handle((aVoid, t) -> {
                if (t != null) {
                    parentFuture.completeExceptionally(t);
                }
                Map<PlcRequest, Either<PlcResponse, Exception>> results = new HashMap<>();
                for (Map.Entry<PlcRequest, CompletableFuture<PlcResponse>> subFutureEntry : subFutures.entrySet()) {
                    PlcRequest subRequest = subFutureEntry.getKey();
                    CompletableFuture<PlcResponse> subFuture = subFutureEntry.getValue();
                    try {
                        final PlcResponse subResponse = subFuture.get();
                        results.put(subRequest, Either.left(subResponse));
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        results.put(subRequest, Either.right(new Exception("Something went wrong")));
                    } catch (Exception e) {
                        results.put(subRequest, Either.right(new Exception("Something went wrong")));
                    }
                }
                PlcResponse response = responseProcessor.apply(results);
                parentFuture.complete(response);
                return Void.TYPE;
            }).exceptionally(throwable -> {
                // TODO: If would be cool if we could still process all of the successful ones ...
                parentFuture.completeExceptionally(throwable);
                return null;
            });
            return parentFuture;
        } else {
            return CompletableFuture.completedFuture(responseProcessor.apply(Collections.emptyMap()));
        }
    }

}
