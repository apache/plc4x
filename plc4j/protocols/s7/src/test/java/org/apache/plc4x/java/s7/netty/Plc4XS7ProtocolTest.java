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
package org.apache.plc4x.java.s7.netty;

import io.netty.channel.ChannelFuture;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelOutboundHandlerAdapter;
import io.netty.channel.ChannelPromise;
import io.netty.channel.embedded.EmbeddedChannel;
import io.netty.handler.codec.EncoderException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.s7.netty.model.messages.S7Message;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.java.s7.netty.util.S7PlcFieldHandler;
import org.junit.Before;
import org.junit.Test;

import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.hamcrest.Matchers.*;
import static org.hamcrest.collection.IsCollectionWithSize.hasSize;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class Plc4XS7ProtocolTest {

    private EmbeddedChannel SUT;
    private PlcReadRequest.Builder readRequestBuilder =
        new DefaultPlcReadRequest.Builder(mock(PlcReader.class), new S7PlcFieldHandler());
    private PlcWriteRequest.Builder writeRequestBuilder =
        new DefaultPlcWriteRequest.Builder(mock(PlcWriter.class), new S7PlcFieldHandler());
    private CompletableFuture<S7Message> writeFuture;

    @Before
    public void setUp() {
        writeFuture = new CompletableFuture<>();
        SUT = new EmbeddedChannel(new ChannelOutboundHandlerAdapter() {
            @Override
            public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) {
                if(msg instanceof S7Message) {
                    writeFuture.complete((S7Message) msg);
                    promise.setSuccess();
                } else {
                    promise.setFailure(new PlcProtocolException(
                        "Got message of type " + msg.getClass().getSimpleName()));
                }
            }
        }, new Plc4XS7Protocol());
    }

    @Test
    @SuppressWarnings("unchecked")
    public void testInvalidFieldType() {
        CompletableFuture<InternalPlcReadRequest> future = new CompletableFuture<>();
        DefaultPlcReadRequest readRequest = mock(DefaultPlcReadRequest.class);
        when(readRequest.getFieldNames()).thenReturn(new LinkedHashSet<>(Collections.singleton("foo")));
        when(readRequest.getField("foo")).thenReturn(mock(PlcField.class));
        PlcRequestContainer container = new PlcRequestContainer(readRequest, future);
        ChannelFuture channelFuture = SUT.writeOneOutbound(container);
        assertThat("The promise should have been set to 'success'", channelFuture.isSuccess(), equalTo(false));

        Throwable exception = channelFuture.cause();
        assertThat("An exception should have been thrown", exception, notNullValue());
        assertThat(exception, instanceOf(EncoderException.class));
        EncoderException encoderException = (EncoderException) exception;
        assertThat(encoderException.getCause(), instanceOf(PlcProtocolException.class));
    }

    @Test
    @SuppressWarnings("unchecked")
    public void testSimpleReadVarRequest() throws InterruptedException, ExecutionException, TimeoutException {
        CompletableFuture<InternalPlcReadRequest> future = new CompletableFuture<>();
        PlcRequestContainer container = new PlcRequestContainer(
            (DefaultPlcReadRequest) readRequestBuilder.addItem("foo", "%Q0:BYTE").build(), future);
        ChannelFuture channelFuture = SUT.writeOneOutbound(container);
        assertThat("The promise should have been set to 'success'", channelFuture.isSuccess(), equalTo(true));

        Throwable exception = channelFuture.cause();
        assertThat("No exception should have been thrown", exception, nullValue());

        S7Message writtenMessage = writeFuture.get(100, TimeUnit.MILLISECONDS);
        assertThat("The protocol layer should have output something", writtenMessage, notNullValue());
        assertThat("The protocol layer should have output something", writtenMessage, instanceOf(S7RequestMessage.class));

        assertThat("The message should have one parameter", writtenMessage.getParameters(), hasSize(1));
        assertThat("The message should have no payload", writtenMessage.getPayloads(), empty());
        assertThat("The request container should be assigned parent to the write message",
            writtenMessage.getParent(), equalTo(container));

        assertThat(writtenMessage.getParameters().get(0), instanceOf(VarParameter.class));
        VarParameter varParameter = (VarParameter) writtenMessage.getParameters().get(0);
        assertThat(varParameter.getItems(), hasSize(1));

        assertThat(varParameter.getItems().get(0), instanceOf(S7AnyVarParameterItem.class));
        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameter.getItems().get(0);
        assertThat(s7AnyVarParameterItem.getSpecificationType(), equalTo(SpecificationType.VARIABLE_SPECIFICATION));
        assertThat(s7AnyVarParameterItem.getMemoryArea(), equalTo(MemoryArea.OUTPUTS));
        assertThat(s7AnyVarParameterItem.getDataBlockNumber(), equalTo(0));
        assertThat(s7AnyVarParameterItem.getByteOffset(), equalTo(0));
        assertThat(s7AnyVarParameterItem.getBitOffset(), equalTo((byte) 0));
        assertThat(s7AnyVarParameterItem.getNumElements(), equalTo(1));
        assertThat(s7AnyVarParameterItem.getDataType(), equalTo(TransportSize.BYTE));
    }

    @Test
    @SuppressWarnings("unchecked")
    public void testSimpleWriteVarRequest() throws InterruptedException, ExecutionException, TimeoutException {
        CompletableFuture<InternalPlcWriteRequest> future = new CompletableFuture<>();
        PlcRequestContainer container = new PlcRequestContainer(
            (DefaultPlcWriteRequest) writeRequestBuilder.addItem("foo", "%Q0:BYTE", (byte) 0x42).build(), future);
        ChannelFuture channelFuture = SUT.writeOneOutbound(container);
        assertThat("The promise should have been set to 'success'", channelFuture.isSuccess(), equalTo(true));

        Throwable exception = channelFuture.cause();
        assertThat("No exception should have been thrown", exception, nullValue());

        S7Message writtenMessage = writeFuture.get(100, TimeUnit.MILLISECONDS);
        assertThat("The protocol layer should have output something", writtenMessage, notNullValue());
        assertThat("The protocol layer should have output something", writtenMessage, instanceOf(S7RequestMessage.class));

        assertThat("The message should have one parameter", writtenMessage.getParameters(), hasSize(1));
        assertThat("The message should have one payload", writtenMessage.getPayloads(), hasSize(1));
        assertThat("The request container should be assigned parent to the write message",
            writtenMessage.getParent(), equalTo(container));

        assertThat(writtenMessage.getParameters().get(0), instanceOf(VarParameter.class));
        VarParameter varParameter = (VarParameter) writtenMessage.getParameters().get(0);
        assertThat(varParameter.getItems(), hasSize(1));
        assertThat(varParameter.getItems().get(0), instanceOf(S7AnyVarParameterItem.class));
        S7AnyVarParameterItem s7AnyVarParameterItem = (S7AnyVarParameterItem) varParameter.getItems().get(0);
        assertThat(s7AnyVarParameterItem.getSpecificationType(), equalTo(SpecificationType.VARIABLE_SPECIFICATION));
        assertThat(s7AnyVarParameterItem.getMemoryArea(), equalTo(MemoryArea.OUTPUTS));
        assertThat(s7AnyVarParameterItem.getDataBlockNumber(), equalTo(0));
        assertThat(s7AnyVarParameterItem.getByteOffset(), equalTo(0));
        assertThat(s7AnyVarParameterItem.getBitOffset(), equalTo((byte) 0));
        assertThat(s7AnyVarParameterItem.getNumElements(), equalTo(1));
        assertThat(s7AnyVarParameterItem.getDataType(), equalTo(TransportSize.BYTE));

        assertThat(writtenMessage.getPayloads().get(0), instanceOf(VarPayload.class));
        VarPayload varPayload = (VarPayload) writtenMessage.getPayloads().get(0);
        assertThat(varPayload.getItems(), hasSize(1));
        assertThat(varPayload.getItems().get(0), instanceOf(VarPayloadItem.class));
        VarPayloadItem varPayloadItem = varPayload.getItems().get(0);
        assertThat(varPayloadItem.getReturnCode(), equalTo(DataTransportErrorCode.RESERVED));
        assertThat(varPayloadItem.getDataTransportSize(), equalTo(DataTransportSize.BYTE_WORD_DWORD));
        assertThat(varPayloadItem.getData(), notNullValue());
        assertThat(varPayloadItem.getData().length, equalTo(1));
        assertThat(varPayloadItem.getData()[0], equalTo((byte) 0x42));
    }

}
