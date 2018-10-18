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

package org.apache.plc4x.java.base.protocol;

import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelPromise;
import io.netty.channel.PendingWriteQueue;
import io.netty.util.HashedWheelTimer;
import org.apache.commons.lang3.reflect.FieldUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcFieldRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.*;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SingleItemToSingleRequestProtocolTest implements WithAssertions {

    PlcReader mockReader = null;
    PlcWriter mockWriter = null;
    PlcSubscriber mockSubscriber = null;

    @InjectMocks
    SingleItemToSingleRequestProtocol SUT = new SingleItemToSingleRequestProtocol(
        mockReader,
        mockWriter,
        new HashedWheelTimer(),
        TimeUnit.SECONDS.toMillis(1),
        false
    );

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    ChannelHandlerContext channelHandlerContext;

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    ChannelPromise channelPromise;

    @Mock
    CompletableFuture<InternalPlcResponse> responseCompletableFuture;

    @BeforeEach
    void setUp() throws Exception {
        SUT.channelRegistered(channelHandlerContext);
        when(channelHandlerContext.executor().inEventLoop()).thenReturn(true);
    }

    @AfterEach
    void tearDown() throws Exception {
        SUT.channelUnregistered(channelHandlerContext);
    }

    @Nested
    class Misc {
        @Test
        void channelRegistered() throws Exception {
            SUT.channelRegistered(channelHandlerContext);
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 0),
                entry("deliveredItems", 0L),
                entry("erroredItems", 0L),
                entry("erroredContainers", 0L),
                entry("deliveredContainers", 0L)
            );
        }

        @Test
        void channelUnregistered() throws Exception {
            SUT.channelUnregistered(channelHandlerContext);
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 0),
                entry("deliveredItems", 0L),
                entry("erroredItems", 0L),
                entry("deliveredContainers", 0L),
                entry("erroredContainers", 0L)
            );
        }

        @Test
        void channelInactive() throws Exception {
            SUT.channelInactive(channelHandlerContext);
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 0),
                entry("deliveredItems", 0L),
                entry("erroredItems", 0L),
                entry("deliveredContainers", 0L),
                entry("erroredContainers", 0L)
            );
        }
    }

    @Nested
    class Roundtrip {
        @Captor
        ArgumentCaptor<PlcRequestContainer> plcRequestContainerArgumentCaptor;

        @Test
        void simpleRead() throws Exception {
            // Given
            // we have a simple read
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(mockReader), responseCompletableFuture);
            // When
            // we write this
            SUT.write(channelHandlerContext, msg, channelPromise);
            // And
            // and we simulate that all get responded
            verify(channelHandlerContext, times(5)).write(plcRequestContainerArgumentCaptor.capture(), any());
            List<PlcRequestContainer> capturedDownstreamContainers = plcRequestContainerArgumentCaptor.getAllValues();
            capturedDownstreamContainers.forEach(this::produceReadResponse);
            // Then
            // our complete container should complete normally
            verify(responseCompletableFuture).complete(any());
            // And we should have no memory leak
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 5),
                entry("erroredItems", 0L),
                entry("deliveredItems", 5L),
                entry("deliveredContainers", 1L),
                entry("erroredContainers", 0L)
            );
        }

        @Test
        void partialRead() throws Exception {
            // Given
            // we have a simple read
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(mockReader), responseCompletableFuture);
            // When
            // we write this
            SUT.write(channelHandlerContext, msg, channelPromise);
            // And
            // and we simulate that some one responded
            verify(channelHandlerContext, times(5)).write(plcRequestContainerArgumentCaptor.capture(), any());
            List<PlcRequestContainer> capturedDownstreamContainers = plcRequestContainerArgumentCaptor.getAllValues();
            capturedDownstreamContainers.stream().findFirst().map(this::produceReadResponse);
            // Then
            // We create SUT with 1 seconds timeout
            TimeUnit.SECONDS.sleep(2);
            // our complete container should complete normally
            verify(responseCompletableFuture).completeExceptionally(any());
            // And we should have no memory leak
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 5),
                entry("deliveredItems", 1L),
                entry("erroredItems", 4L),
                entry("deliveredContainers", 0L),
                entry("erroredContainers", 1L)
            );
        }

        @Test
        void partialReadOneErrored() throws Exception {
            // Given
            // we have a simple read
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(mockReader), responseCompletableFuture);
            // When
            // we write this
            SUT.write(channelHandlerContext, msg, channelPromise);
            // And
            // and we simulate that some one responded
            verify(channelHandlerContext, times(5)).write(plcRequestContainerArgumentCaptor.capture(), any());
            List<PlcRequestContainer> capturedDownstreamContainers = plcRequestContainerArgumentCaptor.getAllValues();
            capturedDownstreamContainers.stream()
                .findFirst()
                .map(plcRequestContainer ->
                    plcRequestContainer
                        .getResponseFuture()
                        .completeExceptionally(new RuntimeException("ErrorOccurred"))
                );
            // Then
            // We create SUT with 1 seconds timeout
            TimeUnit.SECONDS.sleep(2);
            // our complete container should complete normally
            verify(responseCompletableFuture).completeExceptionally(any());
            // And we should have no memory leak
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 5),
                entry("deliveredItems", 0L),
                entry("erroredItems", 1L),
                entry("deliveredContainers", 0L),
                entry("erroredContainers", 1L)
            );
        }

        @Test
        void noRead() throws Exception {
            // Given
            // we have a simple read
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(mockReader), responseCompletableFuture);
            // When
            // we write this
            SUT.write(channelHandlerContext, msg, channelPromise);
            // And
            // and we simulate that some one responded
            verify(channelHandlerContext, times(5)).write(any(), any());
            // Then
            // We create SUT with 1 seconds timeout
            TimeUnit.SECONDS.sleep(2);
            // our complete container should complete normally
            verify(responseCompletableFuture).completeExceptionally(any());
            // And we should have no memory leak
            assertThat(SUT.getStatistics()).containsOnly(
                entry("queue", 0),
                entry("sentButUnacknowledgedSubContainer", 0),
                entry("correlationToParentContainer", 0),
                entry("containerCorrelationIdMap", 0),
                entry("responsesToBeDelivered", 0),
                entry("correlationIdGenerator", 5),
                entry("deliveredItems", 0L),
                entry("erroredItems", 5L),
                entry("deliveredContainers", 0L),
                entry("erroredContainers", 1L)
            );
        }

        @SuppressWarnings("unchecked")
        private Void produceReadResponse(PlcRequestContainer plcRequestContainer) {
            InternalPlcReadRequest request = (InternalPlcReadRequest) plcRequestContainer.getRequest();
            String fieldName = request.getFieldNames().iterator().next();
            CompletableFuture responseFuture = plcRequestContainer.getResponseFuture();
            HashMap<String, Pair<PlcResponseCode, FieldItem>> responseFields = new HashMap<>();
            responseFields.put(fieldName, Pair.of(PlcResponseCode.OK, mock(FieldItem.class)));
            responseFuture.complete(new DefaultPlcReadResponse(request, responseFields));
            return null;
        }
    }

    @Nested
    class Decoding {
        @Test
        void tryFinish() throws Exception {
            SUT.tryFinish(1, null, new CompletableFuture<>());
            // TODO: add Assertions.
        }

        @Test
        void errored() throws Exception {
            SUT.errored(1, mock(Throwable.class), new CompletableFuture<>());
            // TODO: add Assertions.
        }
    }

    @Nested
    class Encoding {

        @Captor
        ArgumentCaptor<PlcRequestContainer> plcRequestContainerArgumentCaptor;

        @Test
        void empty() throws Exception {
            // Given
            Object msg = null;
            // When
            SUT.write(channelHandlerContext, msg, channelPromise);
            // Then
            verify(channelHandlerContext, times(1)).write(null, channelPromise);
        }

        @Test
        void read() throws Exception {
            // Given
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(mockReader), responseCompletableFuture);
            // When
            SUT.write(channelHandlerContext, msg, channelPromise);
            // Then
            verify(channelHandlerContext, times(5)).write(plcRequestContainerArgumentCaptor.capture(), any());
            List<PlcRequestContainer> capturedValues = plcRequestContainerArgumentCaptor.getAllValues();
            // We check if every request as exactly one field
            assertThat(capturedValues)
                .allMatch(plcRequestContainer -> plcRequestContainer.getRequest() instanceof SingleItemToSingleRequestProtocol.CorrelatedPlcReadRequest)
                .allMatch(plcRequestContainer -> ((SingleItemToSingleRequestProtocol.CorrelatedPlcReadRequest) plcRequestContainer.getRequest()).getNumberOfFields() == 1);
            // In sum we should see all fields
            List<String> fieldNamesList = capturedValues.stream()
                .map(PlcRequestContainer::getRequest)
                .map(PlcFieldRequest.class::cast)
                .map(PlcFieldRequest::getFieldNames)
                .flatMap(Collection::stream)
                .collect(Collectors.toList());
            // There should be no duplications
            assertThat(fieldNamesList).hasSize(5);
            assertThat(fieldNamesList).containsExactly(
                "readField1",
                "readField2",
                "readField3",
                "readField4",
                "readField5"
            );
        }

        @Test
        void write() throws Exception {
            // Given
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcWriteRequest.build(mockWriter), responseCompletableFuture);
            // When
            SUT.write(channelHandlerContext, msg, channelPromise);
            // Then
            verify(channelHandlerContext, times(5)).write(plcRequestContainerArgumentCaptor.capture(), any());
            List<PlcRequestContainer> capturedValues = plcRequestContainerArgumentCaptor.getAllValues();
            // We check if every request as exactly one field
            assertThat(capturedValues)
                .allMatch(plcRequestContainer -> plcRequestContainer.getRequest() instanceof SingleItemToSingleRequestProtocol.CorrelatedPlcWriteRequest)
                .allMatch(plcRequestContainer -> ((SingleItemToSingleRequestProtocol.CorrelatedPlcWriteRequest) plcRequestContainer.getRequest()).getNumberOfFields() == 1);
            // In sum we should see all fields
            List<String> fieldNamesList = capturedValues.stream()
                .map(PlcRequestContainer::getRequest)
                .map(PlcFieldRequest.class::cast)
                .map(PlcFieldRequest::getFieldNames)
                .flatMap(Collection::stream)
                .collect(Collectors.toList());
            // There should be no duplications
            assertThat(fieldNamesList).hasSize(5);
            assertThat(fieldNamesList).containsExactly(
                "writeField1",
                "writeField2",
                "writeField3",
                "writeField4",
                "writeField5"
            );
        }

        @Test
        void subscribe() throws Exception {
            // TODO: implement once available
        }

        @Test
        void unsubcribe() throws Exception {
            // TODO: implement once available
        }

        @Test
        void trySendingMessages() throws Exception {
            PendingWriteQueue queue = (PendingWriteQueue) FieldUtils.getField(SUT.getClass(), "queue", true).get(SUT);
            assertThat(queue.size()).isLessThanOrEqualTo(0);
            queue.add(mock(PlcRequestContainer.class), channelPromise);
            assertThat(queue.size()).isGreaterThan(0);
            SUT.trySendingMessages(channelHandlerContext);
            assertThat(queue.size()).isLessThanOrEqualTo(0);
        }
    }

    private static class TestDefaultPlcReadRequest extends DefaultPlcReadRequest {
        private TestDefaultPlcReadRequest(PlcReader reader, LinkedHashMap<String, PlcField> fields) {
            super(reader, fields);
        }

        private static TestDefaultPlcReadRequest build(PlcReader reader) {
            LinkedHashMap<String, PlcField> fields = new LinkedHashMap<>();
            IntStream.rangeClosed(1, 5).forEach(i -> fields.put("readField" + i, mock(PlcField.class)));
            return new TestDefaultPlcReadRequest(reader, fields);
        }
    }

    private static class TestDefaultPlcWriteRequest extends DefaultPlcWriteRequest {

        private TestDefaultPlcWriteRequest(PlcWriter writer, LinkedHashMap<String, Pair<PlcField, FieldItem>> fields) {
            super(writer, fields);
        }

        private static TestDefaultPlcWriteRequest build(PlcWriter writer) {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> fields = new LinkedHashMap<>();
            IntStream.rangeClosed(1, 5).forEach(i -> fields.put("writeField" + i, Pair.of(mock(PlcField.class), mock(FieldItem.class))));
            return new TestDefaultPlcWriteRequest(writer, fields);
        }
    }

    private static class TestDefaultPlcSubscriptionRequest extends DefaultPlcSubscriptionRequest {

        private TestDefaultPlcSubscriptionRequest(PlcSubscriber subscriber, LinkedHashMap<String, SubscriptionPlcField> fields) {
            super(subscriber, fields);
        }

        private static TestDefaultPlcSubscriptionRequest build(PlcSubscriber subscriber) {
            // TODO: implement me once available
            return new TestDefaultPlcSubscriptionRequest(subscriber, new LinkedHashMap<>());
        }
    }

    private static class TestDefaultPlcUnsubscriptionRequest extends DefaultPlcUnsubscriptionRequest {

        private TestDefaultPlcUnsubscriptionRequest(PlcSubscriber subscriber, Collection<? extends InternalPlcSubscriptionHandle> internalPlcSubscriptionHandles) {
            super(subscriber, internalPlcSubscriptionHandles);
        }

        private static TestDefaultPlcUnsubscriptionRequest build(PlcSubscriber subscriber) {
            // TODO: implement me once available
            return new TestDefaultPlcUnsubscriptionRequest(subscriber, Collections.emptyList());
        }
    }
}