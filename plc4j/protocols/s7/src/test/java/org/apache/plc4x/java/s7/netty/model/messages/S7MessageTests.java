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

import static org.assertj.core.api.Assertions.assertThat;

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

public class S7MessageTests {

    @Test
    @Category(FastTests.class)
    public void setupCommunictionsRequestMessage() {
        short tpduReference = 1;
        short maxAmqCaller = 4;
        short maxAmqCallee = 8;
        short pduLength = 128;

        SetupCommunicationRequestMessage setupMessage = new SetupCommunicationRequestMessage(tpduReference, maxAmqCaller, maxAmqCallee, pduLength);

        assertThat(setupMessage.getTpduReference()).isEqualTo(tpduReference).withFailMessage("Unexpected tpdu value");
        assertThat(setupMessage.getMessageType()).isEqualTo(MessageType.JOB).withFailMessage("Unexpected message type");
    }

    @Test
    @Category(FastTests.class)
    public void s7RequestMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = null;
        ArrayList<S7Payload> s7Payloads = null;

        S7RequestMessage message = new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads);

        assertThat(message.getTpduReference()).isEqualTo(tpduReference).withFailMessage("Unexpected tpdu value");
        assertThat(message.getMessageType()).isEqualTo(MessageType.USER_DATA).withFailMessage("Unexpected message type");
        assertThat(message.getPayloads()).isNull();
        assertThat(message.getParameters()).isNull();
    }

    @Test
    @Category(FastTests.class)
    public void s7ResponseMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = null;
        ArrayList<S7Payload> s7Payloads = null;
        byte errorClass = 0x1;
        byte errorCode = 0x23;

        S7ResponseMessage message = new S7ResponseMessage(messageType, tpduReference, s7Parameters, s7Payloads, errorClass, errorCode);

        assertThat(message.getTpduReference()).isEqualTo(tpduReference).withFailMessage("Unexpected tpdu value");
        assertThat(message.getMessageType()).isEqualTo(MessageType.USER_DATA).withFailMessage("Unexpected message type");
        assertThat(message.getErrorClass()).isEqualTo((byte) 0x1).withFailMessage("Unexpected error class");
        assertThat(message.getErrorCode()).isEqualTo((byte) 0x23).withFailMessage("Unexpected error code");
        assertThat(message.getPayloads()).isNull();
        assertThat(message.getParameters()).isNull();
    }

    @Test
    @Category(FastTests.class)
    public void s7MessageParameters() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = new ArrayList<>();
        ArrayList<S7Payload> s7Payloads = new ArrayList<>();
        ParameterType parameterType = ParameterType.READ_VAR;
        ArrayList<VarParameterItem> parameterItems = new ArrayList<>();
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        TransportSize transportSize = TransportSize.INT;
        short numElements = 1;
        byte dataBlock = (byte) 0x1;
        byte byteOffset = (byte) 0x10;
        byte bitOffset = (byte) 0x0;

        parameterItems.add(new S7AnyVarParameterItem(specificationType, memoryArea, transportSize, numElements, dataBlock, byteOffset, bitOffset));

        VarParameter varParameter = new VarParameter(parameterType, parameterItems);

        s7Parameters.add(varParameter);

        S7RequestMessage message = new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads);

        assertThat(message.getTpduReference()).isEqualTo(tpduReference).withFailMessage("Unexpected tpdu value");
        assertThat(message.getMessageType()).isEqualTo(MessageType.USER_DATA).withFailMessage("Unexpected message type");
        assertThat(message.getParameters()).hasSize(1).withFailMessage("Unexpected number of parameters");
        assertThat(message.getParameters()).containsAll(s7Parameters).withFailMessage("Unexpected parameters");
        assertThat(message.getParameter(VarParameter.class).get()).isEqualTo(varParameter).withFailMessage("Parameter missing");
        assertThat(message.getParameter(CpuServicesParameter.class).isPresent()).isFalse().withFailMessage("Contains unexpected parameter");
        assertThat(message.getPayloads()).isEmpty();
    }

    @Test
    @Category(FastTests.class)
    public void s7MessagePayload() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = new ArrayList<>();
        ArrayList<S7Payload> s7Payloads = new ArrayList<>();
        ParameterType parameterType = ParameterType.WRITE_VAR;
        ArrayList<VarPayloadItem> payloadItems = new ArrayList<>();
        byte[] data = {(byte)0x79};
        VarPayload varPayload;

        payloadItems.add(new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BIT, data));
        varPayload = new VarPayload(parameterType, payloadItems);
        s7Payloads.add(varPayload);

        S7RequestMessage message = new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads);

        assertThat(message.getTpduReference()).isEqualTo(tpduReference).withFailMessage("Unexpected tpdu value");
        assertThat(message.getMessageType()).isEqualTo(MessageType.USER_DATA).withFailMessage("Unexpected message type");
        assertThat(message.getPayloads()).hasSize(1).withFailMessage("Unexpected number of payloads");
        assertThat(message.getPayloads()).containsAll(s7Payloads).withFailMessage("Unexpected payloads");
        assertThat(message.getPayload(VarPayload.class).get()).isEqualTo(varPayload).withFailMessage("Payload missing");
        assertThat(message.getPayload(VarParameter.class).isPresent()).isFalse().withFailMessage("Contains unexpected payload"); // No other parameter classes
        assertThat(message.getParameters()).isEmpty();
    }

    @Test
    @Category(FastTests.class)
    public void s7AnyVarParameterItem() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = new ArrayList<>();
        ArrayList<S7Payload> s7Payloads = new ArrayList<>();
        ParameterType parameterType = ParameterType.READ_VAR;
        ArrayList<VarParameterItem> parameterItems = new ArrayList<>();
        SpecificationType specificationType = SpecificationType.VARIABLE_SPECIFICATION;
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        TransportSize transportSize = TransportSize.INT;
        short numElements = 1;
        byte dataBlock = (byte) 0x1;
        byte byteOffset = (byte) 0x10;
        byte bitOffset = (byte) 0x0;

        S7AnyVarParameterItem parameterItem = new S7AnyVarParameterItem(specificationType, memoryArea, transportSize, numElements, dataBlock, byteOffset, bitOffset);

        assertThat(parameterItem.getSpecificationType()).isEqualTo(specificationType).withFailMessage("Unexpected specification type");
        assertThat(parameterItem.getMemoryArea()).isEqualTo(MemoryArea.DATA_BLOCKS).withFailMessage("Unexpected memory area");
        assertThat(parameterItem.getTransportSize()).isEqualTo(transportSize).withFailMessage("Unexpected transport size");
        assertThat(parameterItem.getNumElements()).isEqualTo(numElements).withFailMessage("Unexpected number elements");
        assertThat(parameterItem.getDataBlockNumber()).isEqualTo(dataBlock).withFailMessage("Unexpected data block");
        assertThat(parameterItem.getByteOffset()).isEqualTo(byteOffset).withFailMessage("Unexpected byte offset");
        assertThat(parameterItem.getBitOffset()).isEqualTo(bitOffset).withFailMessage("Unexpected bit offset");
        assertThat(parameterItem.getAddressingMode()).isEqualTo(VariableAddressingMode.S7ANY).withFailMessage("Unexpected adressing mode");
    }

}