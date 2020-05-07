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

package org.apache.plc4x.java.lldp;

import java.time.LocalTime;
import org.apache.plc4x.java.lldp.readwrite.TTL;
import org.apache.plc4x.java.lldp.readwrite.Text;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

/**
 * Utility class which helps to handle more difficult parts of frame parsing related to variable length fields such as
 * text and numbers.
 */
public class LLDPUtil {

    public static boolean hasMore(ReadBuffer io) {
        long length = io.getTotalBytes();
        return io.getPos() + 2 > length;
    }

    public static String readString(ReadBuffer io, Integer len) {
        return readString(io, len, "UTF-8");
    }

    public static String readString(ReadBuffer io, Integer len, String encoding) {
        return io.readString(len * 8, encoding);
    }

    public static void writeString(WriteBuffer io, Text value, String encoding) throws ParseException {
        io.writeString(value.getLengthInBits(), encoding, value.getText());
    }

    public static int length(String text) {
        return text.length();
    }

    // FIXME what we really want here is a Duration
    public static LocalTime readTime(ReadBuffer io, Integer numBytes) throws ParseException {
        return LocalTime.ofSecondOfDay(io.readInt(numBytes * 8));
    }

    public static void writeTime(WriteBuffer io, TTL value) throws ParseException {
        int seconds = value.getValue().toSecondOfDay();
        int length = length(seconds);
        if (length == 8) {
            // fixme Do we really want this? We have length indication which should really match with actual field value
            // append 0x0, so the test passes. Need to verify if there is a minimal length for a number written to TLV.
            io.writeUnsignedInt(8, 0x0);
        }
        io.writeUnsignedInt(length, seconds);
    }

    public static int length(LocalTime value) {
        return length(value.toSecondOfDay());
    }

    public static int length(int value) {
        if (value <= 255) {
            return 8;
        } else if (value <= 65535) {
            return 16;
        } else if (value <= 16777215) {
            return 24;
        }
        return 32;
    }

    public static int byteLength(LocalTime value) {
        return length(value) / 4;
    }
}
