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

import static org.hamcrest.core.Is.is;

import static org.hamcrest.MatcherAssert.assertThat;

import org.apache.plc4x.java.s7.netty.model.params.CpuServicesRequestParameter;
import org.apache.plc4x.java.s7.netty.model.params.SetupCommunicationParameter;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.apache.plc4x.java.s7.types.S7DataType;
import org.junit.Test;

import java.util.Arrays;
import java.util.Collections;

public class S7SizeHelperTest {

    @Test
    public void getParametersLengthTest() {
        assertThat(S7SizeHelper.getParametersLength(null), is((short) 0));
        assertThat(S7SizeHelper.getParametersLength(Collections.emptyList()), is((short) 0));
        assertThat(S7SizeHelper.getParametersLength(Arrays.asList(
            new SetupCommunicationParameter((short) 8, (short)8, (short)250),
            new VarParameter(ParameterType.READ_VAR, Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    S7DataType.BYTE, (short) 1, (short) 0, (short) 0, (byte) 0))))), is((short) 22));
    }

    @Test
    public void getPayloadsLengthTest() {
        assertThat(S7SizeHelper.getPayloadsLength(null), is((short) 0));
        assertThat(S7SizeHelper.getPayloadsLength(Collections.emptyList()), is((short) 0));
        assertThat(S7SizeHelper.getPayloadsLength(Collections.singletonList(
            new VarPayload(ParameterType.READ_VAR, Collections.singletonList(new VarPayloadItem(
                DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[]{
                (byte) 0x01}))))), is((short) 5));
    }

    @Test
    public void getParameterLengthTest() {
        assertThat(S7SizeHelper.getParameterLength(
            new CpuServicesRequestParameter(CpuServicesParameterFunctionGroup.CPU_FUNCTIONS,
                CpuServicesParameterSubFunctionGroup.READ_SSL, (byte) 0)), is((short) 8));
        assertThat(S7SizeHelper.getParameterLength(
            new SetupCommunicationParameter((short) 8, (short)8, (short)250)), is((short) 8));
        assertThat(S7SizeHelper.getParameterLength(
            new VarParameter(ParameterType.READ_VAR, Collections.emptyList())), is((short) 2));
        assertThat(S7SizeHelper.getParameterLength(
            new VarParameter(ParameterType.READ_VAR, Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    S7DataType.BYTE, (short) 1, (short) 0, (short) 0, (byte) 0)))), is((short) 14));
        assertThat(S7SizeHelper.getParameterLength(
            new VarParameter(ParameterType.WRITE_VAR, Collections.emptyList())), is((short) 2));
        assertThat(S7SizeHelper.getParameterLength(
            new VarParameter(ParameterType.WRITE_VAR, Collections.singletonList(
                new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
                    S7DataType.BYTE, (short) 1, (short) 0, (short) 0, (byte) 0)))), is((short) 14));
    }

    @Test
    public void getPayloadLengthFromParameterTest() {
        // One bit is transferred inside one byte (4 byte header and one payload)
        assertThat(S7SizeHelper.getPayloadLength(new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
            S7DataType.BOOL, (short) 1, (short) 0, (short) 0, (byte) 0)), is((short) 5));
        assertThat(S7SizeHelper.getPayloadLength(new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
            S7DataType.BYTE, (short) 1, (short) 0, (short) 0, (byte) 0)), is((short) 5));
        assertThat(S7SizeHelper.getPayloadLength(new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
            S7DataType.BYTE, (short) 42, (short) 0, (short) 0, (byte) 0)), is((short) 46));
        assertThat(S7SizeHelper.getPayloadLength(new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
            S7DataType.DWORD, (short) 1, (short) 0, (short) 0, (byte) 0)), is((short) 8));
        assertThat(S7SizeHelper.getPayloadLength(new S7AnyVarParameterItem(SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS,
            S7DataType.DWORD, (short) 42, (short) 0, (short) 0, (byte) 0)), is((short) 172));
    }

    @Test
    public void getPayloadLengthFromPayloadTest() {
        assertThat(S7SizeHelper.getPayloadLength(new VarPayloadItem(
            DataTransportErrorCode.OK, DataTransportSize.BIT, new byte[] {
                (byte) 0x01})), is((short) 5));
        assertThat(S7SizeHelper.getPayloadLength(new VarPayloadItem(
            DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {
                (byte) 0x01})), is((short) 5));
        assertThat(S7SizeHelper.getPayloadLength(new VarPayloadItem(
            DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, new byte[] {
                (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04, (byte) 0x05})), is((short) 9));
        assertThat(S7SizeHelper.getPayloadLength(new VarPayloadItem(
            DataTransportErrorCode.OK, DataTransportSize.REAL, new byte[] {
                (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04})), is((short) 8));
        assertThat(S7SizeHelper.getPayloadLength(new VarPayloadItem(
            DataTransportErrorCode.OK, DataTransportSize.REAL, new byte[] {
                (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04,
                (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04,
                (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04,
                (byte) 0x01, (byte) 0x02, (byte) 0x03, (byte) 0x04})), is((short) 20));
    }

}
