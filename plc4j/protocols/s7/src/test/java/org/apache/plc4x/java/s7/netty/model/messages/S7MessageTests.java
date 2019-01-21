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

package org.apache.plc4x.java.s7.netty.model.messages;

import org.apache.plc4x.java.s7.netty.model.params.CpuServicesParameter;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import java.util.ArrayList;
import java.util.List;

import static org.hamcrest.collection.IsCollectionWithSize.hasSize;
import static org.hamcrest.collection.IsEmptyCollection.empty;
import static org.hamcrest.collection.IsIterableContainingInOrder.contains;
import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;

public class S7MessageTests {

    @Test
    @Category(FastTests.class)
    public void setupCommunicationsRequestMessage() {
        short tpduReference = 1;
        short maxAmqCaller = 4;
        short maxAmqCallee = 8;
        short pduLength = 128;

        SetupCommunicationRequestMessage setupMessage = new SetupCommunicationRequestMessage(
            tpduReference, maxAmqCaller, maxAmqCallee, pduLength, null);

        assertThat("Unexpected tpdu value", setupMessage.getTpduReference(), equalTo(tpduReference));
        assertThat("Unexpected message type", setupMessage.getMessageType(), equalTo(MessageType.JOB));
    }

    @Test
    @Category(FastTests.class)
    public void s7RequestMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;

        S7RequestMessage message = new S7RequestMessage(
            messageType, tpduReference, null, null, null);

        assertThat("Unexpected tpdu value", message.getTpduReference(), equalTo(tpduReference));
        assertThat("Unexpected message type", message.getMessageType(), equalTo(MessageType.USER_DATA));
        assertThat(message.getPayloads(), nullValue());
        assertThat(message.getParameters(), nullValue());
    }

    @Test
    @Category(FastTests.class)
    public void s7ResponseMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        byte errorClass = 0x1;
        byte errorCode = 0x23;

        S7ResponseMessage message = new S7ResponseMessage(
            messageType, tpduReference, null, null, errorClass, errorCode);

        assertThat("Unexpected tpdu value", message.getTpduReference(), equalTo(tpduReference));
        assertThat("Unexpected message type", message.getMessageType(), equalTo(MessageType.USER_DATA));
        assertThat("Unexpected error class", message.getErrorClass(), equalTo((byte) 0x1));
        assertThat("Unexpected error code", message.getErrorCode(), equalTo((byte) 0x23));
        assertThat(message.getPayloads(), nullValue());
        assertThat(message.getParameters(), nullValue());
    }

    @Test
    @Category(FastTests.class)
    public void s7MessageParameters() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        List<S7Parameter> s7Parameters = new ArrayList<>();
        List<S7Payload> s7Payloads = new ArrayList<>();
        ParameterType parameterType = ParameterType.READ_VAR;
        List<VarParameterItem> parameterItems = new ArrayList<>();
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        TransportSize dataType = TransportSize.INT;
        short numElements = 1;
        byte dataBlock = (byte) 0x1;
        byte byteOffset = (byte) 0x10;
        byte bitOffset = (byte) 0x0;

        parameterItems.add(new S7AnyVarParameterItem(specificationType, memoryArea, dataType, numElements, dataBlock, byteOffset, bitOffset));

        VarParameter varParameter = new VarParameter(parameterType, parameterItems);

        s7Parameters.add(varParameter);

        S7RequestMessage message = new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads, null);

        assertThat("Unexpected tpdu value", message.getTpduReference(), equalTo(tpduReference));
        assertThat("Unexpected message type", message.getMessageType(), equalTo(MessageType.USER_DATA));
        assertThat("Unexpected number of parameters", message.getParameters(), hasSize(1));
        assertThat("Unexpected parameters", message.getParameters(), contains(varParameter));
        assertThat("Parameter missing", message.getParameter(VarParameter.class).isPresent(), is(true));
        assertThat("Parameter missing", message.getParameter(VarParameter.class).get(), equalTo(varParameter));
        assertThat("Contains unexpected parameter", message.getParameter(CpuServicesParameter.class).isPresent(), is(false));
        assertThat(message.getPayloads(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void s7MessagePayload() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        List<S7Parameter> s7Parameters = new ArrayList<>();
        List<S7Payload> s7Payloads = new ArrayList<>();
        ParameterType parameterType = ParameterType.WRITE_VAR;
        List<VarPayloadItem> payloadItems = new ArrayList<>();
        byte[] data = {(byte) 0x79};
        VarPayload varPayload;

        payloadItems.add(new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, data));
        varPayload = new VarPayload(parameterType, payloadItems);
        s7Payloads.add(varPayload);

        S7RequestMessage message = new S7RequestMessage(
            messageType, tpduReference, s7Parameters, s7Payloads, null);

        assertThat("Unexpected tpdu value", message.getTpduReference(), equalTo(tpduReference));
        assertThat("Unexpected message type", message.getMessageType(), equalTo(MessageType.USER_DATA));
        assertThat("Unexpected number of payloads", message.getPayloads(), hasSize(1));
        assertThat("Unexpected payloads", message.getPayloads(), contains(varPayload));
        assertThat("Payload missing", message.getPayload(VarPayload.class).isPresent(), is(true));
        assertThat("Payload doesn't match", message.getPayload(VarPayload.class).get(), equalTo(varPayload));
        assertThat("Contains unexpected payload", message.getPayload(VarParameter.class).isPresent(), is(false)); // No other parameter classes
        assertThat(message.getParameters(), empty());
    }

    @Test
    @Category(FastTests.class)
    public void s7AnyVarParameterItem() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        List<S7Parameter> s7Parameters = new ArrayList<>();
        List<S7Payload> s7Payloads = new ArrayList<>();
        ParameterType parameterType = ParameterType.READ_VAR;
        List<VarParameterItem> parameterItems = new ArrayList<>();
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        TransportSize dataType = TransportSize.INT;
        int numElements = 1;
        int dataBlock = 0x1;
        int byteOffset = 0x10;
        byte bitOffset = (byte) 0x0;

        S7AnyVarParameterItem parameterItem = new S7AnyVarParameterItem(
            specificationType, memoryArea, dataType, numElements, dataBlock, byteOffset, bitOffset);

        assertThat("Unexpected specification type", parameterItem.getSpecificationType(), equalTo(specificationType));
        assertThat("Unexpected memory area", parameterItem.getMemoryArea(), equalTo(MemoryArea.DATA_BLOCKS));
        assertThat("Unexpected transport size", parameterItem.getDataType(), equalTo(dataType));
        assertThat("Unexpected number elements", parameterItem.getNumElements(), equalTo(numElements));
        assertThat("Unexpected data block", parameterItem.getDataBlockNumber(), equalTo(dataBlock));
        assertThat("Unexpected byte offset", parameterItem.getByteOffset(), equalTo(byteOffset));
        assertThat("Unexpected bit offset", parameterItem.getBitOffset(), equalTo(bitOffset));
        assertThat("Unexpected adressing mode", parameterItem.getAddressingMode(), equalTo(VariableAddressingMode.S7ANY));
    }

}