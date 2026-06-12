/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.opcua.protocol.chunk;

import java.nio.ByteBuffer;
import org.apache.plc4x.java.opcua.readwrite.BinaryPayload;
import org.apache.plc4x.java.opcua.readwrite.ExtensiblePayload;
import org.apache.plc4x.java.opcua.readwrite.ExtensionObject;
import org.apache.plc4x.java.opcua.readwrite.MessagePDU;
import org.apache.plc4x.java.opcua.readwrite.Payload;
import org.apache.plc4x.java.opcua.readwrite.RootExtensionObject;
import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;

public class PayloadConverter {

    /**
     * Full set of buffer options the OPC UA binary encoding relies on. Just the
     * byte-order isn't enough — the typed field writers/readers query the buffer
     * for the integer/float/string encodings too, and fail without them.
     */
    public static final WithOption[] LITTLE_ENDIAN = {
        WithByteBasedOption.WithByteOrder("LITTLE_ENDIAN"),
        WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
        WithOption.WithSignedIntegerEncoding("twos-complement"),
        WithOption.WithFloatEncoding("IEEE754"),
        WithOption.WithStringEncoding("UTF8")
    };

    public static BinaryPayload toBinary(Payload payload) throws BufferException {
        if (payload instanceof BinaryPayload) {
            return (BinaryPayload) payload;
        }
        return toBinary((ExtensiblePayload) payload);
    }

    public static BinaryPayload toBinary(ExtensiblePayload extensible) throws BufferException {
        ExtensionObject payload = extensible.getPayload();
        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[payload.getLengthInBytes()], LITTLE_ENDIAN);
        payload.serialize(buffer);
        return new BinaryPayload(extensible.getSequenceHeader(), buffer.getBytes());
    }

    public static ExtensiblePayload toExtensible(BinaryPayload binary) throws BufferException {
        byte[] payload = binary.getPayload();
        ReadBufferByteBased buffer = new ReadBufferByteBased(payload, LITTLE_ENDIAN);
        RootExtensionObject extensionObject = (RootExtensionObject) RootExtensionObject.staticParse(buffer, false);
        return new ExtensiblePayload(binary.getSequenceHeader(), extensionObject);
    }

    public static byte[] toStream(Payload payload) throws BufferException {
        return serialize(payload);
    }

    public static byte[] toStream(MessagePDU apdu) throws BufferException {
        return serialize(apdu);
    }

    private static byte[] serialize(Message message) throws BufferException {
        WriteBufferByteBased buffer = new WriteBufferByteBased(new byte[message.getLengthInBytes()], LITTLE_ENDIAN);
        message.serialize(buffer);
        return buffer.getBytes();
    }

    public static Payload fromStream(byte[] payload, boolean binary) throws BufferException {
        ReadBufferByteBased buffer = new ReadBufferByteBased(payload, LITTLE_ENDIAN);
        return Payload.staticParse(buffer, binary, (long) (binary ? payload.length - 8 : -1));
    }

    public static MessagePDU fromStream(ByteBuffer chunkBuffer, boolean response) throws BufferException {
        ReadBufferByteBased buffer = new ReadBufferByteBased(chunkBuffer.array(), LITTLE_ENDIAN);
        return MessagePDU.staticParse(buffer, response, true);
    }

    public static MessagePDU pduFromStream(byte[] message, boolean response) throws BufferException {
        ReadBufferByteBased buffer = new ReadBufferByteBased(message, LITTLE_ENDIAN);
        return MessagePDU.staticParse(buffer, response, true);
    }
}
