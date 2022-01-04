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
package org.apache.plc4x.java.bacnetip.readwrite.utils;

import org.apache.plc4x.java.bacnetip.readwrite.BACnetContextTag;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetPropertyIdentifier;
import org.apache.plc4x.java.bacnetip.readwrite.BACnetTag;
import org.apache.plc4x.java.bacnetip.readwrite.io.BACnetTagIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.util.List;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithAdditionalStringRepresentation;

public class StaticHelper {

    public static BACnetPropertyIdentifier readPropertyIdentifier(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        long readUnsignedLong = readBuffer.readUnsignedLong("propertyIdentifier",bitsToRead);
        BACnetPropertyIdentifier baCnetPropertyIdentifier = BACnetPropertyIdentifier.enumForValue(readUnsignedLong);
        if (baCnetPropertyIdentifier == null) {
            return BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE;
        }
        return baCnetPropertyIdentifier;
    }

    public static void writePropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier value) throws SerializationException {
        if (value == null || value ==  BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        long valueValue = value.getValue();
        if (valueValue <= 0xffL) {
            bitsToWrite = 8;
        } else if (valueValue <= 0xffffL) {
            bitsToWrite = 16;
        } else if (valueValue <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("propertyIdentifier", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }

    public static void writeProprietaryPropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier baCnetPropertyIdentifier, long value) throws SerializationException {
        if (baCnetPropertyIdentifier != null && baCnetPropertyIdentifier !=  BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return;
        }
        int bitsToWrite;
        if (value <= 0xffL) {
            bitsToWrite = 8;
        } else if (value <= 0xffffL) {
            bitsToWrite = 16;
        } else if (value <= 0xffffffffL) {
            bitsToWrite = 32;
        } else {
            bitsToWrite = 32;
        }
        writeBuffer.writeUnsignedLong("proprietaryPropertyIdentifier", bitsToWrite, value, WithAdditionalStringRepresentation(BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE.name()));
    }

    public static Long readProprietaryPropertyIdentifier(ReadBuffer readBuffer, BACnetPropertyIdentifier value, Long actualLength) throws ParseException {
        if (value!=null&& value !=  BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE) {
            return 0L;
        }
        // We need to reset our reader to the position we read before
        readBuffer.reset((int) (readBuffer.getPos()-actualLength));
        int bitsToRead = (int) (actualLength * 8);
        return readBuffer.readUnsignedLong("proprietaryPropertyIdentifier",bitsToRead);
    }

    public static boolean openingClosingTerminate(ReadBuffer readBuffer, BACnetContextTag openingTag) {
        if (openingTag == null) {
            // If we don't have an opening tag at all we can terminate here
            return true;
        }
        int oldPos = readBuffer.getPos();
        byte aByte;
        try {
            aByte = readBuffer.readByte();
        } catch (ArrayIndexOutOfBoundsException ignore) {
            return true;
        } catch (ParseException ignore) {
            // TODO: we should rethrow the exception here
            return false;
        }
        readBuffer.reset(oldPos);
        return aByte == 0x3F;
    }

    public static BACnetTag parseTags(ReadBuffer readBuffer) throws ParseException {
        return BACnetTagIO.staticParse(readBuffer);
    }

    public static void writeTags(WriteBuffer writeBuffer, BACnetTag value) throws SerializationException {
        value.serialize(writeBuffer);
    }

    public static int tagsLength(List<BACnetTag> data) {
        return data.stream()
            .map(BACnetTag::getLengthInBytes).mapToInt(Integer::intValue).sum();
    }
}
