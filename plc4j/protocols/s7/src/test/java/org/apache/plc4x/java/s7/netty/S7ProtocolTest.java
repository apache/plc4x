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

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.isotp.netty.model.IsoTPMessage;
import org.apache.plc4x.java.isotp.netty.model.tpdus.Tpdu;
import org.apache.plc4x.java.netty.NettyTestBase;
import org.apache.plc4x.java.s7.netty.model.messages.S7RequestMessage;
import org.apache.plc4x.java.s7.netty.model.params.VarParameter;
import org.apache.plc4x.java.s7.netty.model.params.items.S7AnyVarParameterItem;
import org.apache.plc4x.java.s7.netty.model.payloads.VarPayload;
import org.apache.plc4x.java.s7.netty.model.payloads.items.VarPayloadItem;
import org.apache.plc4x.java.s7.netty.model.types.*;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.LinkedList;

import static java.util.Collections.singletonList;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.mock;

public class S7ProtocolTest extends NettyTestBase {

    private S7Protocol SUT;

    @BeforeEach
    void setUp() {
        SUT = new S7Protocol((short) 1, (short) 1, (short) 1);
    }

    @Test
    @Tag("fast")
    public void encode() throws Exception {
        //TODO: finish me
        LinkedList<Object> out = new LinkedList<>();
        SUT.encode(null, new S7RequestMessage(
            MessageType.ACK,
            (short) 1,
            singletonList(new VarParameter(ParameterType.WRITE_VAR, singletonList(new S7AnyVarParameterItem(
                SpecificationType.VARIABLE_SPECIFICATION, MemoryArea.DATA_BLOCKS, TransportSize.BIT, (short) 1, (short) 1, (short) 1, (byte) 1
            )))),
            singletonList(new VarPayload(
                ParameterType.WRITE_VAR,
                singletonList(new VarPayloadItem(
                    DataTransportErrorCode.OK,
                    DataTransportSize.BYTE_WORD_DWORD, new byte[]{0})
                ))
            )), out);
        assertThat(out).hasSize(1);
    }

    @Test
    @Tag("fast")
    public void decode() throws Exception {
        //TODO: finish me
        LinkedList<Object> out = new LinkedList<>();
        ByteBuf buffer = Unpooled.buffer();
        // Magic Number
        buffer.writeByte(0x32);
        buffer.writeByte(MessageType.JOB.getCode());
        // Reserved magic value
        buffer.writeShort(0x0000);
        // tpduReference
        buffer.writeShort(0x0000);
        // headerParametersLength
        buffer.writeShort(0x0000);
        // userDataLength
        buffer.writeShort(0x0000);
        SUT.decode(null, new IsoTPMessage(mock(Tpdu.class), buffer), out);
        assertThat(out).hasSize(1);
    }

}
