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
package org.apache.plc4x.java.modbus.optimizer;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.Plc4xProtocolBase;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.spi.optimizer.BaseOptimizer;

import java.util.Collections;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class ModbusOptimizer extends BaseOptimizer {


    @Override
    public CompletableFuture<PlcReadResponse> optimizedRead(PlcReadRequest readRequest, Plc4xProtocolBase reader) {
        // If just one field is requested, just forward this to the reader.
        if(readRequest.getNumberOfFields() == 1) {
            return reader.read(readRequest);
        }
        // If more than one filed is requested, split up  each field request into a separate sub-request
        // And have the reader process each one independently. After the last sub-request is finished,
        // Merge the results back together.
        else if (readRequest.getNumberOfFields() > 1) {
            // Create a new future which will be used to return the aggregated response back to the application.
            CompletableFuture<PlcReadResponse> parentFuture = new CompletableFuture<>();

            // Create one sub-request for every single field and store the futures in a map.
            Map<String, CompletableFuture<PlcReadResponse>> subFutures = new HashMap<>();
            for (String fieldName : readRequest.getFieldNames()) {
                PlcField field = readRequest.getField(fieldName);
                PlcReadRequest subRequest = new DefaultPlcReadRequest(
                    ((DefaultPlcReadRequest) readRequest).getReader(),
                    new LinkedHashMap<>(Collections.singletonMap(fieldName, field)));
                subFutures.put(fieldName, reader.read(subRequest));
            }

            // As soon as all sub-futures are done, merge the indivdual responses back to one
            // big response.
            CompletableFuture.allOf(subFutures.values().toArray(new CompletableFuture[0])).thenApply(aVoid -> {
                Map<String, Pair<PlcResponseCode, PlcValue>> fields = new HashMap<>();
                for (Map.Entry<String, CompletableFuture<PlcReadResponse>> fieldEntry : subFutures.entrySet()) {
                    String fieldName = fieldEntry.getKey();
                    CompletableFuture<PlcReadResponse> subFuture = fieldEntry.getValue();
                    try {
                        final PlcReadResponse subResponse = subFuture.get();
                        if (subFuture.isDone()) {
                            fields.put(fieldName, Pair.of(PlcResponseCode.OK,
                                subResponse.getAsPlcValue().getValue(fieldName)));
                        } else {
                            fields.put(fieldName, Pair.of(subResponse.getResponseCode(fieldName), null));
                        }
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        fields.put(fieldName, Pair.of(PlcResponseCode.INTERNAL_ERROR, null));
                    } catch (ExecutionException e) {
                        fields.put(fieldName, Pair.of(PlcResponseCode.INTERNAL_ERROR, null));
                    }
                }
                PlcReadResponse readResponse = new DefaultPlcReadResponse((InternalPlcReadRequest) readRequest, fields);
                parentFuture.complete(readResponse);
                return Void.TYPE;
            });
            return parentFuture;
        }
        return CompletableFuture.completedFuture(
            new DefaultPlcReadResponse((InternalPlcReadRequest) readRequest, Collections.emptyMap()));
    }

}
