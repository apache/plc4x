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
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;

public class PayloadConverter {

    public static BinaryPayload toBinary(Payload payload) throws SerializationException {
        if (payload instanceof BinaryPayload) {
            return (BinaryPayload) payload;
        }

        return toBinary((ExtensiblePayload) payload);
    }


    public static BinaryPayload toBinary(ExtensiblePayload extensible) throws SerializationException {
        ExtensionObject payload = extensible.getPayload();

        WriteBufferByteBased buffer = new WriteBufferByteBased(payload.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
        payload.serialize(buffer);

        return new BinaryPayload(extensible.getSequenceHeader(), buffer.getBytes());
    }

    public static ExtensiblePayload toExtensible(BinaryPayload binary) throws ParseException {
        byte[] payload = binary.getPayload();

        ReadBufferByteBased buffer = new ReadBufferByteBased(payload, ByteOrder.LITTLE_ENDIAN);
        RootExtensionObject extensionObject = (RootExtensionObject) RootExtensionObject.staticParse(buffer, false);

        return new ExtensiblePayload(binary.getSequenceHeader(), extensionObject);
    }

    public static byte[] toStream(Payload payload) throws SerializationException {
        return serialize(payload);
    }

    public static byte[] toStream(MessagePDU apdu) throws SerializationException {
        return serialize(apdu);
    }

    private static byte[] serialize(Message message) throws SerializationException {
        WriteBufferByteBased buffer = new WriteBufferByteBased(message.getLengthInBytes(), ByteOrder.LITTLE_ENDIAN);
        message.serialize(buffer);

        return buffer.getBytes();
    }

    public static Payload fromStream(byte[] payload, boolean binary) throws ParseException {
        ReadBufferByteBased buffer = new ReadBufferByteBased(payload, ByteOrder.LITTLE_ENDIAN);
        return Payload.staticParse(buffer, binary, (long) (binary ? payload.length - 8 : -1));
    }

    public static MessagePDU fromStream(ByteBuffer chunkBuffer, boolean response) throws ParseException {
        ReadBufferByteBased buffer = new ReadBufferByteBased(chunkBuffer.array(), ByteOrder.LITTLE_ENDIAN);
        return MessagePDU.staticParse(buffer, response, true);
    }

    public static MessagePDU pduFromStream(byte[] message, boolean response) throws ParseException {
        ReadBufferByteBased buffer = new ReadBufferByteBased(message, ByteOrder.LITTLE_ENDIAN);
        return MessagePDU.staticParse(buffer, response, true);
    }
}
