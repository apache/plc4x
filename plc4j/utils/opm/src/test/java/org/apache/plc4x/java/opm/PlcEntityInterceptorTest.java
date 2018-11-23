/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.opm;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.DefaultPlcReadResponse;
import org.junit.Test;
import org.mockito.Mockito;
import org.mockito.stubbing.Answer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Collections;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.Mockito.when;

public class PlcEntityInterceptorTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcEntityInterceptorTest.class);

    private void runGetPlcResponseWIthException(Answer a) throws InterruptedException, ExecutionException, TimeoutException, OPMException {
        PlcReadRequest request = Mockito.mock(PlcReadRequest.class);
        CompletableFuture future = Mockito.mock(CompletableFuture.class);
        when(future.get(anyLong(), any())).then(a);
        when(request.execute()).thenReturn(future);

        PlcEntityInterceptor.getPlcReadResponse(request);
    }

    @Test
    public void getPlcReadResponse_catchesInterruptedException_rethrows() throws InterruptedException {
        AtomicBoolean exceptionWasThrown = new AtomicBoolean(false);
        // Run in different Thread
        Thread thread = new Thread(new Runnable() {
            @Override
            public void run() {
                try {
                    runGetPlcResponseWIthException(invocation -> {
                        throw new InterruptedException();
                    });
                } catch (InterruptedException | ExecutionException | TimeoutException e) {
                    LOGGER.warn("Fetched exception", e);
                } catch (OPMException e) {
                    exceptionWasThrown.set(true);
                }
            }
        });
        thread.start();
        thread.join();
        assertTrue(exceptionWasThrown.get());
    }

    @Test(expected = OPMException.class)
    public void getPlcReadResponse_catchesExecutionException_rethrows() throws OPMException, InterruptedException, ExecutionException, TimeoutException {
        runGetPlcResponseWIthException(invocation -> {
            throw new ExecutionException(new Exception());
        });
    }

    @Test(expected = OPMException.class)
    public void getPlcReadResponse_timeoutOnGet_rethrows() throws OPMException {
        PlcReadRequest request = Mockito.mock(PlcReadRequest.class);
        CompletableFuture future = new CompletableFuture<>();
        when(request.execute()).thenReturn(future);

        PlcEntityInterceptor.getPlcReadResponse(request);
    }

    @Test
    public void getTyped_notOkResponse_throws() {
        DefaultPlcReadResponse response = new DefaultPlcReadResponse(null, Collections.singletonMap("field", Pair.of(PlcResponseCode.NOT_FOUND, null)));
        String message = null;
        try {
            PlcEntityInterceptor.getTyped(Long.class, response, "field");
        } catch (PlcRuntimeException e) {
            message = e.getMessage();
        }
        assertEquals("Unable to read specified field 'field', response code was 'NOT_FOUND'", message);
    }

    @Test
    public void getterWithNoField() throws OPMException {
        PlcEntityManager entityManager = new PlcEntityManager();
        BadEntity entity = entityManager.connect(BadEntity.class, "test:test");

        String message = null;
        try {
            entity.getField1();
        } catch (Exception e) {
            message = e.getMessage();
        }
        assertEquals("Unable to identify field with name 'field1' for call to 'getField1'", message);
    }

    @PlcEntity
    public static class BadEntity {

        public BadEntity() {
            // For OPM
        }

        // Getter with no field
        @PlcField("field1")
        public String getField1() {
            return "";
        }

    }

}