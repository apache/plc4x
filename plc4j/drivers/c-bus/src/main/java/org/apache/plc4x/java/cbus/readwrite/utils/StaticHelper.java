/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.cbus.readwrite.utils;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.cbus.readwrite.CALReply;
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.spi.generation.*;

import java.nio.charset.StandardCharsets;

public class StaticHelper {
    public static void writeCBusCommand(WriteBuffer writeBuffer, CBusCommand cbusCommand) throws SerializationException {
        // TODO: maybe we use a writebuffer hex based
        WriteBufferByteBased payloadWriteBuffer = new WriteBufferByteBased(cbusCommand.getLengthInBytes() * 2);
        cbusCommand.serialize(payloadWriteBuffer);
        byte[] hexBytes = Hex.encodeHexString(payloadWriteBuffer.getBytes()).getBytes(StandardCharsets.UTF_8);
        writeBuffer.writeByteArray("cbusCommand", hexBytes);
    }

    public static CBusCommand readCBusCommand(ReadBuffer readBuffer, Integer payloadLength, boolean srcchk) throws ParseException {
        byte[] hexBytes = readBuffer.readByteArray("cbusCommand", payloadLength);
        byte[] rawBytes;
        try {
            rawBytes = Hex.decodeHex(new String(hexBytes));
        } catch (DecoderException e) {
            throw new ParseException("error getting hex", e);
        }
        return CBusCommand.staticParse(new ReadBufferByteBased(rawBytes), srcchk);
    }

    public static void writeCALReply(WriteBuffer writeBuffer, CALReply calReply) throws SerializationException {
        // TODO: maybe we use a writebuffer hex based
        WriteBufferByteBased payloadWriteBuffer = new WriteBufferByteBased(calReply.getLengthInBytes() * 2);
        calReply.serialize(payloadWriteBuffer);
        byte[] hexBytes = Hex.encodeHexString(payloadWriteBuffer.getBytes()).getBytes(StandardCharsets.UTF_8);
        writeBuffer.writeByteArray("calReply", hexBytes);
    }

    public static CALReply readCALReply(ReadBuffer readBuffer, Integer payloadLength) throws ParseException {
        byte[] hexBytes = readBuffer.readByteArray("calReply", payloadLength);
        byte[] rawBytes;
        try {
            rawBytes = Hex.decodeHex(new String(hexBytes));
        } catch (DecoderException e) {
            throw new ParseException("error getting hex", e);
        }
        return CALReply.staticParse(new ReadBufferByteBased(rawBytes));
    }

}
