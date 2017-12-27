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
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;

import static org.junit.jupiter.api.Assertions.*;

class S7MessageTests {

    @Test
    @Tag("fast")
    void setupCommunictionsRequestMessage() {
        short tpduReference = 1;
        short maxAmqCaller = 4;
        short maxAmqCallee = 8;
        short pduLength = 128;

        SetupCommunicationRequestMessage setupMessage = new SetupCommunicationRequestMessage(tpduReference, maxAmqCaller, maxAmqCallee, pduLength);

        assertTrue(setupMessage.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(setupMessage.getMessageType() == MessageType.JOB, "Unexpected message type");
    }

    @Test
    @Tag("fast")
    void s7RequestMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = null;
        ArrayList<S7Payload> s7Payloads = null;

        S7RequestMessage message = new S7RequestMessage(messageType, tpduReference, s7Parameters, s7Payloads);

        assertTrue(message.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(message.getMessageType() == MessageType.USER_DATA, "Unexpected message type");
        assertTrue(message.getPayloads() == null, "Unexpected payloads");
        assertTrue(message.getParameters() == null, "Unexpected parameters");
    }

    @Test
    @Tag("fast")
    void s7ResponseMessage() {
        MessageType messageType = MessageType.USER_DATA;
        short tpduReference = 1;
        ArrayList<S7Parameter> s7Parameters = null;
        ArrayList<S7Payload> s7Payloads = null;
        byte errorClass = 0x1;
        byte errorCode = 0x23;

        S7ResponseMessage message = new S7ResponseMessage(messageType, tpduReference, s7Parameters, s7Payloads, errorClass, errorCode);

        assertTrue(message.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(message.getMessageType() == MessageType.USER_DATA, "Unexpected message type");
        assertTrue(message.getErrorClass() == 0x1, "Unexpected error class");
        assertTrue(message.getErrorCode() == 0x23, "Unexpected error code");
        assertTrue(message.getPayloads() == null, "Unexpected payloads");
        assertTrue(message.getParameters() == null, "Unexpected parameters");
    }

    @Test
    @Tag("fast")
    void s7MessageParameters() {
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

        assertTrue(message.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(message.getMessageType() == MessageType.USER_DATA, "Unexpected message type");
        assertTrue(message.getParameters().size() == 1, "Unexpected number of parameters");
        assertTrue(message.getParameters().containsAll(s7Parameters), "Unexpected parameters");
        assertTrue(message.getParameter(VarParameter.class).equals(varParameter), "Parameter missing");
        assertTrue(message.getParameter(CpuServicesParameter.class) == null, "Contains unexpected parameter");
        assertTrue(message.getPayloads().size() == 0, "Unexpected number of payloads");
    }

    @Test
    @Tag("fast")
    void s7MessagePayload() {
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

        assertTrue(message.getTpduReference() == tpduReference, "Unexpected tpdu value");
        assertTrue(message.getMessageType() == MessageType.USER_DATA, "Unexpected message type");
        assertTrue(message.getPayloads().size() == 1, "Unexpected number of payloads");
        assertTrue(message.getPayloads().containsAll(s7Payloads), "Unexpected payloads");
        assertTrue(message.getPayload(VarPayload.class).equals(varPayload), "Payload missing");
        assertTrue(message.getPayload(VarParameter.class) == null, "Contains unexpected payload"); // No other parameter classes
        assertTrue(message.getParameters().size() == 0, "Unexpected number of parameters");
    }

    @Test
    @Tag("fast")
    void s7AnyVarParameterItem() {
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

        assertTrue(parameterItem.getSpecificationType() == specificationType, "Unexpected specification type");
        assertTrue(parameterItem.getMemoryArea() == MemoryArea.DATA_BLOCKS, "Unexpected memory area");
        assertTrue(parameterItem.getTransportSize() == transportSize, "Unexpected transport size");
        assertTrue(parameterItem.getNumElements() == numElements, "Unexpected number elements");
        assertTrue(parameterItem.getDataBlockNumber() == dataBlock, "Unexpected data block");
        assertTrue(parameterItem.getByteOffset() == byteOffset, "Unexpected byte offset");
        assertTrue(parameterItem.getBitOffset() == bitOffset, "Unexpected bit offset");
        assertTrue(parameterItem.getAddressingMode() == VariableAddressingMode.S7ANY, "Unexpected adressing mode");
    }

}