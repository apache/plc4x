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
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.messages.PlcFieldRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.*;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collection;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SingleItemToSingleRequestProtocolTest implements WithAssertions {

    @InjectMocks
    SingleItemToSingleRequestProtocol SUT = new SingleItemToSingleRequestProtocol(TimeUnit.SECONDS.toMillis(1), false);

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
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(), responseCompletableFuture);
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
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(), responseCompletableFuture);
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
        void noRead() throws Exception {
            // Given
            // we have a simple read
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(), responseCompletableFuture);
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
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcReadRequest.build(), responseCompletableFuture);
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
            PlcRequestContainer<?, ?> msg = new PlcRequestContainer<>(TestDefaultPlcWriteRequest.build(), responseCompletableFuture);
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
        void trySendingMessages() throws Exception {
            SUT.trySendingMessages(channelHandlerContext);
            // TODO: add assertions
        }
    }

    private static class TestDefaultPlcReadRequest extends DefaultPlcReadRequest {

        private TestDefaultPlcReadRequest(LinkedHashMap<String, PlcField> fields) {
            super(fields);
        }

        private static TestDefaultPlcReadRequest build() {
            LinkedHashMap<String, PlcField> fields = new LinkedHashMap<>();
            fields.put("readField1", mock(PlcField.class));
            fields.put("readField2", mock(PlcField.class));
            fields.put("readField3", mock(PlcField.class));
            fields.put("readField4", mock(PlcField.class));
            fields.put("readField5", mock(PlcField.class));
            return new TestDefaultPlcReadRequest(fields);
        }
    }

    private static class TestDefaultPlcWriteRequest extends DefaultPlcWriteRequest {

        private TestDefaultPlcWriteRequest(LinkedHashMap<String, Pair<PlcField, FieldItem>> fields) {
            super(fields);
        }

        private static TestDefaultPlcWriteRequest build() {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> fields = new LinkedHashMap<>();
            fields.put("writeField1", Pair.of(mock(PlcField.class), mock(FieldItem.class)));
            fields.put("writeField2", Pair.of(mock(PlcField.class), mock(FieldItem.class)));
            fields.put("writeField3", Pair.of(mock(PlcField.class), mock(FieldItem.class)));
            fields.put("writeField4", Pair.of(mock(PlcField.class), mock(FieldItem.class)));
            fields.put("writeField5", Pair.of(mock(PlcField.class), mock(FieldItem.class)));
            return new TestDefaultPlcWriteRequest(fields);
        }
    }
}