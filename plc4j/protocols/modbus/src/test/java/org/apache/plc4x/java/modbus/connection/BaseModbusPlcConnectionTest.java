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
package org.apache.plc4x.java.modbus.connection;

import io.netty.channel.Channel;
import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandler;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.GenericFutureListener;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.base.connection.ChannelFactory;
import org.apache.plc4x.java.base.messages.InternalPlcReadRequest;
import org.apache.plc4x.java.base.messages.InternalPlcWriteRequest;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Answers;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import static org.apache.plc4x.java.base.util.Junit5Backport.assertThrows;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.fail;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.class)
@SuppressWarnings("unchecked")
public class BaseModbusPlcConnectionTest {

    private static final Logger LOGGER = LoggerFactory.getLogger(BaseModbusPlcConnectionTest.class);

    private BaseModbusPlcConnection SUT;

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    private ChannelFactory channelFactory;

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    private Channel channel;

    @Before
    public void setUp() throws Exception {
        SUT = new BaseModbusPlcConnection(channelFactory, null) {
            @Override
            protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
                return null;
            }
        };

        when(channelFactory.createChannel(any())).thenReturn(channel);

        SUT.connect();
    }

    @Test
    public void lazyConstructor() {
        new BaseModbusPlcConnection(channelFactory, null) {
            @Override
            protected ChannelHandler getChannelHandler(CompletableFuture<Void> sessionSetupCompleteFuture) {
                return null;
            }
        };
    }

    @Test
    public void read() {
        CompletableFuture<PlcReadResponse> read = SUT.read(mock(InternalPlcReadRequest.class));
        assertNotNull(read);

        simulatePipelineError(() -> SUT.read(mock(InternalPlcReadRequest.class)));
    }

    @Test
    public void write() {
        CompletableFuture<PlcWriteResponse> write = SUT.write(mock(InternalPlcWriteRequest.class));
        assertNotNull(write);

        simulatePipelineError(() -> SUT.write(mock(InternalPlcWriteRequest.class)));
    }

    public void simulatePipelineError(FutureProducingTestRunnable futureProducingTestRunnable) {
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

    @Test
    public void testToString() {
        String s = SUT.toString();
        assertNotNull(s);
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