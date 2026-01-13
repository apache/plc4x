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
import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;

public class StaticHelper {

    public static Checksum readAndValidateChecksum(ReadBuffer readBuffer, Message message, boolean srchk) throws BufferException {
        if (!srchk) {
            return null;
        }
        byte checksum = readBytesFromHex("chksum", readBuffer, false)[0];
        byte actualChecksum = getChecksum(message);
        if (checksum != actualChecksum) {
            throw new BufferException(String.format("Expected checksum 0x%x doesn't match actual checksum 0x%x", checksum, actualChecksum));
        }
        return new Checksum(checksum);
    }

    public static void calculateChecksum(WriteBuffer writeBuffer, Message message, boolean srchk) throws BufferException {
        if (!srchk) {
            // Nothing to do when srchck is disabled
            return;
        }
        writeToHex("chksum", writeBuffer, new byte[]{getChecksum(message)});
    }

    private static byte getChecksum(Message message) throws BufferException {
        byte checksum = 0x0;
        WriteBufferByteBased checksumWriteBuffer = new WriteBufferByteBased(new byte[message.getLengthInBytes()]);
        message.serialize(checksumWriteBuffer);
        for (byte aByte : checksumWriteBuffer.getBytes()) {
            checksum += aByte;
        }
        checksum = (byte) ~checksum;
        checksum++;
        return checksum;
    }

    public static void writeCBusCommand(WriteBuffer writeBuffer, CBusCommand cbusCommand) throws BufferException {
        writeToHex("cbusCommand", writeBuffer, cbusCommand);
    }

    public static CBusCommand readCBusCommand(ReadBuffer readBuffer, CBusOptions cBusOptions, boolean srchk) throws BufferException {
        byte[] rawBytes = readBytesFromHex("cbusCommand", readBuffer, srchk);
        return CBusCommand.staticParse(new ReadBufferByteBased(rawBytes), cBusOptions);
    }

    public static void writeEncodedReply(WriteBuffer writeBuffer, EncodedReply encodedReply) throws BufferException {
        writeToHex("encodedReply", writeBuffer, encodedReply);
    }

    public static EncodedReply readEncodedReply(ReadBuffer readBuffer, CBusOptions cBusOptions, RequestContext requestContext, boolean srchk) throws BufferException {
        byte[] rawBytes = readBytesFromHex("encodedReply", readBuffer, srchk);
        return EncodedReply.staticParse(new ReadBufferByteBased(rawBytes), cBusOptions, requestContext);
    }

    public static void writeCALData(WriteBuffer writeBuffer, CALData calData) throws BufferException {
        writeToHex("calData", writeBuffer, calData);
    }

    public static CALData readCALData(ReadBuffer readBuffer) throws BufferException {
        byte[] rawBytes = readBytesFromHex("calData", readBuffer, false);
        return CALData.staticParse(new ReadBufferByteBased(rawBytes), (RequestContext) null);
    }

    private static byte[] readBytesFromHex(String logicalName, ReadBuffer readBuffer, boolean srchk) throws BufferException {
        int payloadLength = findHexEnd(readBuffer);
        if (payloadLength == 0) {
            throw new BufferException("Length is 0");
        }

        byte[] hexBytes = readBuffer.readBits(payloadLength, WithOption.WithName(logicalName));
        byte lastByte = hexBytes[hexBytes.length - 1];
        if ((lastByte >= 0x67) && (lastByte <= 0x7A)) {
            // We need to reset the alpha
            readBuffer.setPositionInBits(readBuffer.getPositionInBits() - 1);
            hexBytes = Arrays.copyOf(hexBytes, hexBytes.length - 1);
        }
        byte[] rawBytes;
        try {
            rawBytes = Hex.decodeHex(new String(hexBytes));
        } catch (DecoderException e) {
            throw new BufferException("error getting hex", e);
        }
        if (srchk) {
            byte checksum = 0x0;
            for (byte aByte : rawBytes) {
                checksum += aByte;
            }
            if (checksum != 0x0) {
                //throw new ParseException("Checksum validation failed");
            }
            // We need to reset the last to hex bytes
            readBuffer.setPositionInBits(readBuffer.getPositionInBits() - 2);
            rawBytes = Arrays.copyOf(rawBytes, rawBytes.length - 1);
        }
        return rawBytes;
    }

    private static int findHexEnd(ReadBuffer readBuffer) throws BufferException {
        // TODO: find out if there is a smarter way to find the end...
        int oldPos = readBuffer.getPositionInBits();
        int payloadLength = 0;
        while (readBuffer.getRemainingBits() >= 8) {
            char hexByte = (char) readBuffer.readSignedByte(8);
            boolean isHex = hexByte >= 'A' && hexByte <= 'F' || hexByte >= 'a' && hexByte <= 'f';
            boolean isNumber = hexByte >= '0' && hexByte <= '9';
            if (!isHex && !isNumber) {
                break;
            }
            payloadLength++;
        }
        readBuffer.setPositionInBits(oldPos);
        return payloadLength;
    }

    private static void writeToHex(String logicalName, WriteBuffer writeBuffer, Message message) throws BufferException {
        // TODO: maybe we use a writeBuffer hex based
        WriteBufferByteBased payloadWriteBuffer = new WriteBufferByteBased(new byte[message.getLengthInBytes() * 2]);
        message.serialize(payloadWriteBuffer);
        writeToHex(logicalName, writeBuffer, payloadWriteBuffer.getBytes());
    }

    private static void writeToHex(String logicalName, WriteBuffer writeBuffer, byte[] bytes) throws BufferException {
        byte[] hexBytes = Hex.encodeHexString(bytes, false).getBytes(StandardCharsets.UTF_8);
        writeBuffer.writeBits(hexBytes.length * 8, hexBytes, WithOption.WithName(logicalName));
    }

    public static boolean knowsCALCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return CALCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsLightingCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return LightingCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsSecurityCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return SecurityCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsMeteringCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return MeteringCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsTriggerControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return TriggerControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsEnableControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return EnableControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsTemperatureBroadcastCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return TemperatureBroadcastCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsAccessControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return AccessControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsMediaTransportControlCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return MediaTransportControlCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsClockAndTimekeepingCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return ClockAndTimekeepingCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsTelephonyCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return TelephonyCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsAirConditioningCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return AirConditioningCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsMeasurementCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return MeasurementCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }

    public static boolean knowsErrorReportingCommandTypeContainer(ReadBuffer readBuffer) {
        int oldPos = readBuffer.getPositionInBits();
        try {
            return ErrorReportingCommandTypeContainer.isDefined(readBuffer.readUnsignedShort(8));
        } catch (BufferException ignore) {
            return false;
        } finally {
            readBuffer.setPositionInBits(oldPos);
        }
    }
}
