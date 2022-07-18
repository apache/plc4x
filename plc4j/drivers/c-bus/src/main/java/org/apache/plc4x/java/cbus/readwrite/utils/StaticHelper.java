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
import org.apache.plc4x.java.cbus.readwrite.*;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;

public class StaticHelper {

    public static void writeCBusCommand(WriteBuffer writeBuffer, CBusCommand cbusCommand) throws SerializationException {
        writeToHex("cbusCommand", writeBuffer, cbusCommand, cbusCommand.getLengthInBytes());
    }

    public static CBusCommand readCBusCommand(ReadBuffer readBuffer, Integer payloadLength, CBusOptions cBusOptions) throws ParseException {
        byte[] rawBytes = readBytesFromHex("cbusCommand", readBuffer, payloadLength);
        return CBusCommand.staticParse(new ReadBufferByteBased(rawBytes), cBusOptions);
    }

    public static void writeEncodedReply(WriteBuffer writeBuffer, EncodedReply encodedReply) throws SerializationException{
        writeToHex("encodedReply", writeBuffer, encodedReply, encodedReply.getLengthInBytes());
    }

    public static EncodedReply readEncodedReply(ReadBuffer readBuffer, Integer payloadLength, CBusOptions cBusOptions, RequestContext requestContext) throws ParseException {
        byte[] rawBytes = readBytesFromHex("encodedReply", readBuffer, payloadLength);
        return EncodedReply.staticParse(new ReadBufferByteBased(rawBytes), cBusOptions, requestContext);
    }

    public static void writeCALDataOrSetParameter(WriteBuffer writeBuffer, CALDataOrSetParameter calDataOrSetParameter) throws SerializationException {
        writeToHex("calDataOrSetParameter", writeBuffer, calDataOrSetParameter, calDataOrSetParameter.getLengthInBytes());
    }

    public static CALDataOrSetParameter readCALDataOrSetParameter(ReadBuffer readBuffer, Integer payloadLength) throws ParseException {
        byte[] rawBytes = readBytesFromHex("calDataOrSetParameter", readBuffer, payloadLength);
        return CALDataOrSetParameter.staticParse(new ReadBufferByteBased(rawBytes));
    }

    private static byte[] readBytesFromHex(String logicalName, ReadBuffer readBuffer, Integer payloadLength) throws ParseException {
        if (payloadLength == 0) {
            throw new ParseException("Length is 0");
        }
        byte[] hexBytes = readBuffer.readByteArray(logicalName, payloadLength);
        byte lastByte = hexBytes[hexBytes.length - 1];
        if ((lastByte >= 0x67) && (lastByte <= 0x7A)) {
            // We need to reset the alpha
            readBuffer.reset(readBuffer.getPos() - 1);
            hexBytes = Arrays.copyOf(hexBytes, hexBytes.length - 1);
        }
        byte[] rawBytes;
        try {
            rawBytes = Hex.decodeHex(new String(hexBytes));
        } catch (DecoderException e) {
            throw new ParseException("error getting hex", e);
        }
        return rawBytes;
    }

    private static void writeToHex(String logicalName, WriteBuffer writeBuffer, Serializable serializable, int lengthInBytes) throws SerializationException {
        // TODO: maybe we use a writeBuffer hex based
        WriteBufferByteBased payloadWriteBuffer = new WriteBufferByteBased(lengthInBytes * 2);
        serializable.serialize(payloadWriteBuffer);
        byte[] hexBytes = Hex.encodeHexString(payloadWriteBuffer.getBytes(), false).getBytes(StandardCharsets.UTF_8);
        writeBuffer.writeByteArray(logicalName, hexBytes);
    }

    public static boolean knowsCALCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return CALCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsLightingCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return LightingCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsSecurityCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return SecurityCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

}
