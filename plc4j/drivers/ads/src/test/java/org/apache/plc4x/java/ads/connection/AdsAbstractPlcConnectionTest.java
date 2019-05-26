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
package org.apache.plc4x.java.ads.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.plc4x.java.ads.api.commands.AdsReadWriteResponse;
import org.apache.plc4x.java.ads.api.commands.types.Data;
import org.apache.plc4x.java.ads.api.commands.types.Result;
import org.apache.plc4x.java.ads.api.generic.types.AmsNetId;
import org.apache.plc4x.java.ads.api.generic.types.AmsPort;
import org.apache.plc4x.java.ads.model.DirectAdsField;
import org.apache.plc4x.java.ads.model.SymbolicAdsField;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcFieldRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.messages.*;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Answers;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.jupiter.MockitoExtension;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.Queue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.*;

@SuppressWarnings("unchecked")
@ExtendWith(MockitoExtension.class)
class AdsAbstractPlcConnectionTest implements WithAssertions {

    private static final Logger LOGGER = LoggerFactory.getLogger(AdsAbstractPlcConnectionTest.class);

    private AdsAbstractPlcConnection SUT;

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    private ChannelFactory channelFactory;

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    private Channel channel;

    @BeforeEach
    void setUp() throws Exception {
        SUT = new AdsAbstractPlcConnection(channelFactory, mock(AmsNetId.class), mock(AmsPort.class), mock(AmsNetId.class), mock(AmsPort.class)) {
            @Override
            protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
                return null;
            }
        };

        when(channelFactory.createChannel(any())).thenReturn(channel);

        SUT.connect();
    }

    @Nested
    class Lifecycle {
        @Test
        void lazyConstructor() {
            AdsAbstractPlcConnection constructed = new AdsAbstractPlcConnection(channelFactory, mock(AmsNetId.class), mock(AmsPort.class)) {
                @Override
                protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
                    return null;
                }
            };
            assertEquals(AdsAbstractPlcConnection.generateAMSNetId(), constructed.getSourceAmsNetId());
            assertEquals(AdsAbstractPlcConnection.generateAMSPort(), constructed.getSourceAmsPort());
        }

        @Test
        void close() throws Exception {
            Map fieldMapping = (Map) FieldUtils.getDeclaredField(AdsAbstractPlcConnection.class, "fieldMapping", true).get(SUT);
            fieldMapping.put(mock(SymbolicAdsField.class), mock(DirectAdsField.class));
            SUT.close();
        }
    }


    @Nested
    class Communication {

        @Test
        void read() {
            CompletableFuture<PlcReadResponse> read = SUT.read(mock(InternalPlcReadRequest.class));
            assertNotNull(read);

            simulatePipelineError(() -> SUT.read(mock(InternalPlcReadRequest.class)));
        }

        @Test
        void write() {
            CompletableFuture<PlcWriteResponse> write = SUT.write(mock(InternalPlcWriteRequest.class));
            assertNotNull(write);

            simulatePipelineError(() -> SUT.write(mock(InternalPlcWriteRequest.class)));
        }

        @Test
        void send() {
            CompletableFuture send = SUT.send(mock(InternalPlcProprietaryRequest.class));
            assertNotNull(send);

            simulatePipelineError(() -> SUT.send(mock(InternalPlcProprietaryRequest.class)));
        }

        void simulatePipelineError(FutureProducingTestRunnable futureProducingTestRunnable) {
            ChannelFuture channelFuture = mock(ChannelFuture.class);
            // Simulate error in the pipeline
            when(channelFuture.addListener(any())).thenAnswer(invocation -> {
                Future future = mock(Future.class);
                when(future.isSuccess()).thenReturn(false);
                when(future.cause()).thenReturn(new DummyException());
                GenericFutureListener genericFutureListener = invocation.getArgument(0);
                genericFutureListener.operationComplete(future);
                return mock(ChannelFuture.class);
            });
            when(channel.writeAndFlush(any())).thenReturn(channelFuture);
            assertThrows(DummyException.class, () -> {
                CompletableFuture completableFuture = futureProducingTestRunnable.run();
                try {
                    completableFuture.get(3, TimeUnit.SECONDS);
                    fail("Should have thrown a ExecutionException");
                } catch (ExecutionException e) {
                    if (e.getCause() instanceof DummyException) {
                        throw (DummyException) e.getCause();
                    }
                    throw e;
                }
            });
        }
    }

    @Nested
    class Symbolic {

        @BeforeEach
        void setUp() {
            SUT.clearMapping();
        }

        @Test
        void mapFields() {
            SUT.mapFields(mock(PlcFieldRequest.class));
        }

        @Test
        void mapSingleField() {
            when(channel.writeAndFlush(any(PlcRequestContainer.class))).then(invocation -> {
                PlcRequestContainer plcRequestContainer = invocation.getArgument(0);
                PlcProprietaryResponse plcProprietaryResponse = Mockito.mock(InternalPlcProprietaryResponse.class, RETURNS_DEEP_STUBS);
                AdsReadWriteResponse adsReadWriteResponse = Mockito.mock(AdsReadWriteResponse.class, RETURNS_DEEP_STUBS);
                when(adsReadWriteResponse.getResult()).thenReturn(Result.of(0));
                when(adsReadWriteResponse.getData()).thenReturn(Data.of(new byte[]{1, 2, 3, 4}));
                when(plcProprietaryResponse.getResponse()).thenReturn(adsReadWriteResponse);
                plcRequestContainer.getResponseFuture().complete(plcProprietaryResponse);
                return mock(ChannelFuture.class);
            });

            SUT.mapFields(SymbolicAdsField.of("Main.byByte[0]:BYTE"));
            SUT.mapFields(SymbolicAdsField.of("Main.byByte[0]:BYTE"));
            verify(channel, times(1)).writeAndFlush(any(PlcRequestContainer.class));
            SUT.clearMapping();
            reset(channel);
        }

        @Test
        void mapSingleFieldNegative() {
            when(channel.writeAndFlush(any(PlcRequestContainer.class))).then(invocation -> {
                PlcRequestContainer plcRequestContainer = invocation.getArgument(0);
                PlcProprietaryResponse plcProprietaryResponse = Mockito.mock(InternalPlcProprietaryResponse.class, RETURNS_DEEP_STUBS);
                AdsReadWriteResponse adsReadWriteResponse = Mockito.mock(AdsReadWriteResponse.class, RETURNS_DEEP_STUBS);
                when(adsReadWriteResponse.getResult()).thenReturn(Result.of(1));
                when(plcProprietaryResponse.getResponse()).thenReturn(adsReadWriteResponse);
                plcRequestContainer.getResponseFuture().complete(plcProprietaryResponse);
                return mock(ChannelFuture.class);
            });

            assertThatThrownBy(() -> SUT.mapFields(SymbolicAdsField.of("Main.byByte[0]:INT64")))
                .isInstanceOf(PlcRuntimeException.class)
                .hasMessageMatching("Non error code received .*");
            verify(channel, times(1)).writeAndFlush(any(PlcRequestContainer.class));
            SUT.clearMapping();
            reset(channel);
        }
    }


    @Nested
    class Misc {
        @Test
        void remainingMethods() {
            assertThat(SUT.canRead()).isTrue();
            assertThat(SUT.canWrite()).isTrue();
            assertThat(SUT.readRequestBuilder()).isNotNull();
            assertThat(SUT.writeRequestBuilder()).isNotNull();
        }

        @Test
        void testToString() {
            String s = SUT.toString();
            assertNotNull(s);
        }

        @Test
        void getTargetAmsNetId() {
            AmsNetId targetAmsNetId = SUT.getTargetAmsNetId();
            assertNotNull(targetAmsNetId);
        }

        @Test
        void getTargetAmsPort() {
            AmsPort targetAmsPort = SUT.getTargetAmsPort();
            assertNotNull(targetAmsPort);
        }

        @Test
        void getSourceAmsNetId() {
            AmsNetId sourceAmsNetId = SUT.getSourceAmsNetId();
            assertNotNull(sourceAmsNetId);
        }

        @Test
        void getSourceAmsPort() {
            AmsPort sourceAmsPort = SUT.getSourceAmsPort();
            assertNotNull(sourceAmsPort);
        }

        @Test
        void generateAMSNetId() {
            AmsNetId targetAmsNetId = AdsAbstractPlcConnection.generateAMSNetId();
            assertNotNull(targetAmsNetId);
        }

        @Test
        void generateAMSPort() {
            AmsPort amsPort = AdsAbstractPlcConnection.generateAMSPort();
            assertNotNull(amsPort);
        }

        // TODO: Commented out as it was causing problems with Java 11
        /*@Test
        void getFromFuture() throws Exception {
            runInThread(() -> {
                CompletableFuture completableFuture = mock(CompletableFuture.class, RETURNS_DEEP_STUBS);
                Object fromFuture = SUT.getFromFuture(completableFuture, 1);
                assertNotNull(fromFuture);
            });
            runInThread(() -> {
                CompletableFuture completableFuture = mock(CompletableFuture.class, RETURNS_DEEP_STUBS);
                when(completableFuture.get(anyLong(), any())).thenThrow(InterruptedException.class);
                assertThrows(PlcRuntimeException.class, () -> SUT.getFromFuture(completableFuture, 1));
            });
            runInThread(() -> {
                CompletableFuture completableFuture = mock(CompletableFuture.class, RETURNS_DEEP_STUBS);
                when(completableFuture.get(anyLong(), any())).thenThrow(ExecutionException.class);
                assertThrows(PlcRuntimeException.class, () -> SUT.getFromFuture(completableFuture, 1));
            });
            runInThread(() -> {
                CompletableFuture completableFuture = mock(CompletableFuture.class, RETURNS_DEEP_STUBS);
                when(completableFuture.get(anyLong(), any())).thenThrow(TimeoutException.class);
                assertThrows(PlcRuntimeException.class, () -> SUT.getFromFuture(completableFuture, 1));
            });
            assertFalse("The current Thread should not be interrupted", Thread.currentThread().isInterrupted());
        }*/

        /**
         * Runs tests steps in a dedicated {@link Thread} so a possible {@link InterruptedException} doesn't lead to a
         * interrupt flag being set on the main Thread ({@link Thread#isInterrupted}).
         *
         * @param testRunnable a special {@link Runnable} which adds a {@code throws Exception} to the {@code run} signature.
         * @throws InterruptedException when this {@link Thread} gets interrupted.
         */
        void runInThread(TestRunnable testRunnable) throws InterruptedException {
            Thread thread = new Thread(() -> {
                try {
                    testRunnable.run();
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });
            Queue<Throwable> uncaughtExceptions = new ConcurrentLinkedQueue<>();
            thread.setUncaughtExceptionHandler((t, e) -> uncaughtExceptions.add(e));
            thread.start();
            thread.join();
            if (!uncaughtExceptions.isEmpty()) {
                uncaughtExceptions.forEach(throwable -> LOGGER.error("Assertion Error: Unexpected Exception", throwable));
                throw new AssertionError("Test failures. Check log");
            }
        }
    }

    /**
     * Variant of {@link Runnable} which adds a {@code throws Exception} to the {@code run} signature.
     */
    private interface TestRunnable {
        /**
         * @throws Exception when the test throws a exception.
         * @see Runnable#run()
         */
        void run() throws Exception;
    }

    private static class DummyException extends Exception {

    }

    @FunctionalInterface
    private interface FutureProducingTestRunnable {
        CompletableFuture run() throws Exception;
    }
}