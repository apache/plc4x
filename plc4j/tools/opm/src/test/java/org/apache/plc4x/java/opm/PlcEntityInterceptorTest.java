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
package org.apache.plc4x.java.opm;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Answers;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.mockito.stubbing.Answer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.Callable;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyLong;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
public class PlcEntityInterceptorTest implements WithAssertions {

    private static final Logger LOGGER = LoggerFactory.getLogger(PlcEntityInterceptorTest.class);

    DefaultPlcDriverManager driverManager;

    MockConnection connection;

    PlcEntityManager entityManager;

    @Mock
    MockDevice mockDevice;

    @BeforeEach
    void setUp() throws Exception {
        driverManager = new DefaultPlcDriverManager();
        connection = (MockConnection) driverManager.getConnection("mock:test");
        connection.setDevice(mockDevice);
        entityManager = new PlcEntityManager(driverManager);
    }

    @Test
    public void getPlcReadResponse_catchesInterruptedException_rethrows() throws InterruptedException {
        AtomicBoolean exceptionWasThrown = new AtomicBoolean(false);
        // Run in different Thread
        Thread thread = new Thread(() -> {
            try {
                runGetPlcResponseWIthException(invocation -> {
                    throw new InterruptedException();
                });
            } catch (InterruptedException | ExecutionException | TimeoutException e) {
                LOGGER.warn("Fetched exception", e);
                Thread.currentThread().interrupt();
            } catch (OPMException e) {
                exceptionWasThrown.set(true);
            }
        });
        thread.start();
        thread.join();
        assertTrue(exceptionWasThrown.get());
    }

    @Test
    public void getPlcReadResponse_catchesExecutionException_rethrows() {
        assertThatThrownBy(() -> runGetPlcResponseWIthException(invocation -> {
            throw new ExecutionException(new Exception());
        }))
            .isInstanceOf(OPMException.class);
    }

    @Test
    public void getPlcReadResponse_timeoutOnGet_rethrows() {
        PlcReadRequest request = mock(PlcReadRequest.class);
        when(request.execute()).thenReturn(new CompletableFuture<>());

        assertThatThrownBy(() -> PlcEntityInterceptor.getPlcReadResponse(request))
            .isInstanceOf(OPMException.class);
    }

    @Test
    public void getTyped_notOkResponse_throws() {
        DefaultPlcReadResponse response = new DefaultPlcReadResponse(null,
            Collections.singletonMap("tag", new ResponseItem<>(PlcResponseCode.NOT_FOUND, null)));
        assertThatThrownBy(() -> PlcEntityInterceptor.getTyped(Long.class, response, "tag"))
            .isInstanceOf(PlcRuntimeException.class)
            .hasMessage("Unable to read specified tag 'tag', response code was 'NOT_FOUND'");
    }

    @Test
    public void getterWithNoField() throws OPMException {
        BadEntity entity = entityManager.connect(BadEntity.class, "mock:test");

        assertThatThrownBy(entity::getTag1)
            .isInstanceOf(OPMException.class)
            .hasMessage("Unable to identify tag with name 'tag1' for call to 'getTag1'");
    }

    @Nested
    class Misc {

        @Mock
        Callable callable;

        @Mock(answer = Answers.RETURNS_DEEP_STUBS)
        PlcConnectionManager connectionManager;

        class MiscEntity {

            @PlcTag("asd")
            private String ok2;

            public void getTest(String a) {
            }

            public String getOk() {
                return "";
            }

            public String getOk2() {
                return ok2;
            }

            public void setOk2(String ok) {
            }

            public void setOkOk(String ok, String ok2) {
            }

            public void someNotSetterMethod(String arg) {
            }

            public void something() {

            }
        }

        @Test
        void missingCases() throws Exception {
            when(callable.call()).then(invocation -> {
                throw new PlcRuntimeException("broken");
            });
            Map<String, Instant> lastFetched = new HashMap<>();
            Map<String, Instant> lastWritten = new HashMap<>();
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptGetter(null, MiscEntity.class.getDeclaredMethod("something"), callable, null, null, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Exception during forwarding call");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptGetter(null, MiscEntity.class.getDeclaredMethod("getTest", String.class), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Only getter with no arguments are supported");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptGetter(null, MiscEntity.class.getDeclaredMethod("getOk"), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessageMatching("Unable to identify tag with name .*");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptGetter(null, MiscEntity.class.getDeclaredMethod("getOk2"), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Problem during processing");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptGetter(null, MiscEntity.class.getDeclaredMethod("getOk2"), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Problem during processing")
                .hasStackTraceContaining(" Unable to read specified tag 'org.apache.plc4x.java.opm.PlcEntityInterceptorTest$Misc$MiscEntity.ok2', response code was 'Mock for PlcResponseCode");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptSetter(null, MiscEntity.class.getDeclaredMethod("setOk2", String.class), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Problem during processing")
                .hasStackTraceContaining(" Unable to read specified tag 'org.apache.plc4x.java.opm.PlcEntityInterceptorTest$Misc$MiscEntity.ok2', response code was 'Mock for PlcResponseCode");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptSetter(null, MiscEntity.class.getDeclaredMethod("setOkOk", String.class, String.class), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Only setter with one arguments are supported");
            assertThatThrownBy(() -> PlcEntityInterceptor.interceptSetter(null, MiscEntity.class.getDeclaredMethod("someNotSetterMethod", String.class), callable, null, connectionManager, null, lastFetched, lastWritten))
                .isInstanceOf(OPMException.class)
                .hasMessage("Unable to forward invocation someNotSetterMethod on connected PlcEntity");
        }
    }

    private void runGetPlcResponseWIthException(Answer a) throws InterruptedException, ExecutionException, TimeoutException, OPMException {
        PlcReadRequest request = mock(PlcReadRequest.class);
        CompletableFuture future = mock(CompletableFuture.class);
        when(future.get(anyLong(), any())).then(a);
        when(request.execute()).thenReturn(future);

        PlcEntityInterceptor.getPlcReadResponse(request);
    }

    @PlcEntity
    public static class BadEntity {

        public BadEntity() {
            // For OPM
        }

        // Getter with no tag
        public String getTag1() {
            return "";
        }

    }

}