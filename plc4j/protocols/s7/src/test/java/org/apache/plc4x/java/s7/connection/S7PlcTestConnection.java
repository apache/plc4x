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
package org.apache.plc4x.java.s7.connection;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.embedded.EmbeddedChannel;
import org.apache.plc4x.java.base.connection.TestChannelFactory;
import org.apache.plc4x.java.isotp.netty.model.types.*;
import org.apache.plc4x.java.netty.events.S7ConnectionEvent;
import org.apache.plc4x.java.s7.netty.S7Protocol;
import org.apache.plc4x.java.s7.netty.model.types.MessageType;
import org.apache.plc4x.java.s7.netty.model.types.ParameterType;

public class S7PlcTestConnection extends S7PlcConnection {

    public S7PlcTestConnection(int rack, int slot, String params) {
        super(new TestChannelFactory(), rack, slot, params);
    }

    @Override
    protected void sendChannelCreatedEvent() {
        EmbeddedChannel channel = (EmbeddedChannel) getChannel();

        // Send an event to the pipeline telling the Protocol filters what's going on.
        channel.pipeline().fireUserEventTriggered(new S7ConnectionEvent());

        ByteBuf writtenData = channel.readOutbound();
        byte[] connectionRequest = new byte[writtenData.readableBytes()];
        writtenData.readBytes(connectionRequest);
        // TODO: Check the content of the Iso TP connection request.

        // Send an Iso TP connection response back to the pipeline.
        byte[] connectionConfirm = toByteArray(
            new int[] {
                // ISO on TCP packet
                0x03, 0x00, 0x00, 0x16,
                // ISO TP packet
                0x11,
                TpduCode.CONNECTION_CONFIRM.getCode(),
                0x00, 0x01, 0x00, 0x02,
                ProtocolClass.CLASS_0.getCode(),

                ParameterCode.CALLED_TSAP.getCode(), 0x02, DeviceGroup.PG_OR_PC.getCode(), 0x01,
                ParameterCode.CALLING_TSAP.getCode(), 0x02, DeviceGroup.OTHERS.getCode(), 0x12,
                ParameterCode.TPDU_SIZE.getCode(), 0x01, TpduSize.SIZE_512.getCode()
            });
        channel.writeInbound(Unpooled.wrappedBuffer(connectionConfirm));

        // Read a S7 Setup Communication request.
        writtenData = channel.readOutbound();
        byte[] setupCommunicationRequest = new byte[writtenData.readableBytes()];
        writtenData.readBytes(setupCommunicationRequest);
        // TODO: Check the content of the S7 Setup Communication connection request.

        // Send an S7 Setup Communication response back to the pipeline.
        byte[] setupCommunicationResponse = toByteArray(
            new int[] {
                // ISO on TCP packet
                0x03, 0x00,
                0x00, 0x1B,
                // ISO TP packet
                0x02,
                TpduCode.DATA.getCode(),
                0x80,
                S7Protocol.S7_PROTOCOL_MAGIC_NUMBER,
                MessageType.ACK_DATA.getCode(), 0x00, 0x00,
                0x00, 0x01,
                // Parameter Length
                0x00, 0x08,
                // Data Length
                0x00, 0x00,
                // Error codes
                0x00, 0x00,
                // Parameters:
                ParameterType.SETUP_COMMUNICATION.getCode(), 0x00, 0x00, 0x08, 0x00, 0x08, 0x01, 0x00
            });
        channel.writeInbound(Unpooled.wrappedBuffer(setupCommunicationResponse));
    }

    public static byte[] toByteArray(int[] in) {
        byte[] out = new byte[in.length];
        for(int i = 0; i < in.length; i++) {
            out[i] = (byte) in[i];
        }
        return out;
    }

}
