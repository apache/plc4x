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

import org.apache.plc4x.java.bacnetip.readwrite.BACnetPropertyIdentifier;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithAdditionalStringRepresentation;

public class StaticHelper {

    public static Object readPropertyIdentifier(ReadBuffer readBuffer, Long actualLength) throws ParseException {
        int bitsToRead = (int) (actualLength * 8);
        long readUnsignedLong = readBuffer.readUnsignedLong(bitsToRead);
        return BACnetPropertyIdentifier.enumForValue(readUnsignedLong);
    }

    public static void writePropertyIdentifier(WriteBuffer writeBuffer, BACnetPropertyIdentifier value) throws SerializationException {
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
        writeBuffer.writeUnsignedLong("", bitsToWrite, valueValue, WithAdditionalStringRepresentation(value.name()));
    }
}
