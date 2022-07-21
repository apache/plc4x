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

import java.nio.charset.StandardCharsets;
import java.util.Arrays;

public class StaticHelper {

    public static Checksum readAndValidateChecksum(ReadBuffer readBuffer, Message message, boolean srchk) throws ParseException {
        if (!srchk) {
            return null;
        }
        byte checksum = readBytesFromHex("chksum", readBuffer, 2, false)[0];
        try {
            byte actualChecksum = getChecksum(message);
            if (checksum != actualChecksum) {
                throw new ParseException(String.format("Expected checksum 0x%x doesn't match actual checksum 0x%x", checksum, actualChecksum));
            }
        } catch (SerializationException e) {
            throw new ParseException("Unable to calculate checksum", e);
        }
        return new Checksum(checksum);
    }

    public static void calculateChecksum(WriteBuffer writeBuffer, Message message, boolean srchk) throws SerializationException {
        if (!srchk) {
            // Nothing to do when srchck is disabled
            return;
        }
        writeToHex("chksum", writeBuffer, new byte[]{getChecksum(message)});
    }

    private static byte getChecksum(Message message) throws SerializationException {
        byte checksum = 0x0;
        WriteBufferByteBased checksumWriteBuffer = new WriteBufferByteBased(message.getLengthInBytes());
        message.serialize(checksumWriteBuffer);
        for (byte aByte : checksumWriteBuffer.getBytes()) {
            checksum += aByte;
        }
        checksum = (byte) ~checksum;
        checksum++;
        return checksum;
    }

    public static void writeCBusCommand(WriteBuffer writeBuffer, CBusCommand cbusCommand) throws SerializationException {
        writeToHex("cbusCommand", writeBuffer, cbusCommand);
    }

    public static CBusCommand readCBusCommand(ReadBuffer readBuffer, int payloadLength, CBusOptions cBusOptions, boolean srchk) throws ParseException {
        byte[] rawBytes = readBytesFromHex("cbusCommand", readBuffer, payloadLength, srchk);
        return CBusCommand.staticParse(new ReadBufferByteBased(rawBytes), cBusOptions);
    }

    public static void writeEncodedReply(WriteBuffer writeBuffer, EncodedReply encodedReply) throws SerializationException {
        writeToHex("encodedReply", writeBuffer, encodedReply);
    }

    public static EncodedReply readEncodedReply(ReadBuffer readBuffer, int payloadLength, CBusOptions cBusOptions, RequestContext requestContext, boolean srchk) throws ParseException {
        byte[] rawBytes = readBytesFromHex("encodedReply", readBuffer, payloadLength, srchk);
        return EncodedReply.staticParse(new ReadBufferByteBased(rawBytes), cBusOptions, requestContext);
    }

    public static void writeCALData(WriteBuffer writeBuffer, CALData calData) throws SerializationException {
        writeToHex("calData", writeBuffer, calData);
    }

    public static CALData readCALData(ReadBuffer readBuffer, Integer payloadLength) throws ParseException {
        byte[] rawBytes = readBytesFromHex("calData", readBuffer, payloadLength, false);
        return CALData.staticParse(new ReadBufferByteBased(rawBytes), (RequestContext) null);
    }

    private static byte[] readBytesFromHex(String logicalName, ReadBuffer readBuffer, int payloadLength, boolean srchk) throws ParseException {
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
        if (srchk) {
            // We need to reset the last to hex bytes
            readBuffer.reset(readBuffer.getPos() - 2);
            hexBytes = Arrays.copyOf(hexBytes, hexBytes.length - 2);
        }
        byte[] rawBytes;
        try {
            rawBytes = Hex.decodeHex(new String(hexBytes));
        } catch (DecoderException e) {
            throw new ParseException("error getting hex", e);
        }
        return rawBytes;
    }

    private static void writeToHex(String logicalName, WriteBuffer writeBuffer, Message message) throws SerializationException {
        // TODO: maybe we use a writeBuffer hex based
        WriteBufferByteBased payloadWriteBuffer = new WriteBufferByteBased(message.getLengthInBytes() * 2);
        message.serialize(payloadWriteBuffer);
        writeToHex(logicalName, writeBuffer, payloadWriteBuffer.getBytes());
    }

    private static void writeToHex(String logicalName, WriteBuffer writeBuffer, byte[] bytes) throws SerializationException {
        byte[] hexBytes = Hex.encodeHexString(bytes, false).getBytes(StandardCharsets.UTF_8);
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

    public static boolean knowsMeteringCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return MeteringCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsTriggerControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return TriggerControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsEnableControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return EnableControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsTemperatureBroadcastCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return TemperatureBroadcastCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsAccessControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return AccessControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsMediaTransportControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return MediaTransportControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsClockAndTimekeepingCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return ClockAndTimekeepingCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsTelephonyCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return TelephonyCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsAirConditioningCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return AirConditioningCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsMeasurementCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return MeasurementCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }

    public static boolean knowsErrorReportingCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPos();
        try {
            return ErrorReportingCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (ParseException ignore) {
            return false;
        } finally {
            readBuffer.reset(oldPos);
        }
    }
}
