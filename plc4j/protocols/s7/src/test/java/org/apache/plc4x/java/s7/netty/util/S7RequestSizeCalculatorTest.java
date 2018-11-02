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

package org.apache.plc4x.java.s7.netty.util;

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.S7Parameter;
import org.apache.plc4x.java.s7.netty.model.params.SetupCommunicationParameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.params.items.VarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.S7Payload;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class S7RequestSizeCalculatorTest {

    @Test
    void getReadVarRequestMessageSize() {
        S7RequestMessage readVarRequest = new S7RequestMessage(
            MessageType.JOB,
            (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.READ_VAR, Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS, TransportSize.BYTE, 1, (short) 1, (short) 0, (byte) 0)))),
            Collections.emptyList(),
            null);
        short size = S7RequestSizeCalculator.getRequestMessageSize(readVarRequest);
        assertEquals(24, size);
    }

    @Test
    void getWriteVarRequestMessageSize() {
        S7RequestMessage readVarRequest = new S7RequestMessage(
            MessageType.JOB,
            (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS, TransportSize.BYTE, 1, (short) 1, (short) 0, (byte) 0)))),
            Collections.singletonList(
                new VarPayload(ParameterType.WRITE_VAR, Collections.singletonList(
                    new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {(byte) 0x01})))),
            null);
        short size = S7RequestSizeCalculator.getRequestMessageSize(readVarRequest);
        assertEquals(30, size);
    }

    @Test
    void getSetupCommunicationRequestMessageSize() {
        S7RequestMessage readVarRequest = new S7RequestMessage(
            MessageType.JOB,
            (short) 1,
            Collections.singletonList(
                new SetupCommunicationParameter((short) 1, (short) 2, (short) 250)),
            Collections.emptyList(),
            null);
        short size = S7RequestSizeCalculator.getRequestMessageSize(readVarRequest);
        assertEquals(18, size);
    }

    @Test
    void getUnimplementedParameterItemTypeMessageSize() {
        S7Parameter mockParameter = mock(S7Parameter.class);
        when(mockParameter.getType()).thenReturn(ParameterType.UPLOAD);
        S7RequestMessage readVarRequest = new S7RequestMessage(
            MessageType.JOB,
            (short) 1,
            Collections.singletonList(mockParameter),
            Collections.emptyList(),
            null);
        Assertions.assertThrows(NotImplementedException.class,
            () -> S7RequestSizeCalculator.getRequestMessageSize(readVarRequest));
    }

    @Test
    void getUnimplementedVarAddressingModeMessageSize() {
        VarParameterItem mockParameterItem = mock(VarParameterItem.class);
        when(mockParameterItem.getAddressingMode()).thenReturn(VariableAddressingMode.ALARM_QUERYREQ);
        S7RequestMessage readVarRequest = new S7RequestMessage(
            MessageType.JOB,
            (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, Collections.singletonList(mockParameterItem))),
            Collections.emptyList(),
            null);
        Assertions.assertThrows(NotImplementedException.class,
            () -> S7RequestSizeCalculator.getRequestMessageSize(readVarRequest));
    }

    @Test
    void getUnimplementedPayloadTypeMessageSize() {
        S7Payload mockPayload = mock(S7Payload.class);
        when(mockPayload.getType()).thenReturn(ParameterType.UPLOAD);
        S7RequestMessage readVarRequest = new S7RequestMessage(
            MessageType.JOB,
            (short) 1,
            Collections.singletonList(
                new VarParameter(ParameterType.WRITE_VAR, Collections.singletonList(
                    new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS, TransportSize.BYTE, 1, (short) 1, (short) 0, (byte) 0)))),
            Collections.singletonList(mockPayload),
            null);
        Assertions.assertThrows(NotImplementedException.class,
            () -> S7RequestSizeCalculator.getRequestMessageSize(readVarRequest));
    }

    @Test
    void getRequestItemTotalSize() {
        short size = S7RequestSizeCalculator.getRequestItemTotalSize(
            new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS, TransportSize.BYTE, 1, (short) 1, (short) 0, (byte) 0),
            new VarPayloadItem(DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {(byte) 0x01})
        );
        assertEquals(18, size);
    }

}