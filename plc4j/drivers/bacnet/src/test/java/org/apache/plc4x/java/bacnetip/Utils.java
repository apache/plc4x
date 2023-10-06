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
package org.apache.plc4x.java.bacnetip;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

import java.util.stream.IntStream;

import static org.junit.jupiter.api.Assertions.assertNotNull;

public class Utils {
    static final boolean DUMP_PACKAGES = false;

    static final int PAYLOAD_START_INDEX = 42;

    static BVLC tryParseBytes(int[] rawBytesAsInts) throws ParseException {
        return tryParseBytes(rawBytesAsInts, PAYLOAD_START_INDEX);
    }

    static BVLC tryParseBytes(int[] rawBytesAsInts, int startIndex) throws ParseException {
        return tryParseBytes(rawBytesAsInts, startIndex, DUMP_PACKAGES);
    }

    static BVLC tryParseBytes(int[] rawBytesAsInts, int startIndex, boolean dumpPackages) throws ParseException {
        var rawBytes = (byte[]) ArrayUtils.toPrimitive(IntStream.of(rawBytesAsInts).boxed().map(Integer::byteValue).toArray(Byte[]::new));
        rawBytes = ArrayUtils.subarray(rawBytes, startIndex, rawBytes.length);
        BVLC bvlc = BVLC.staticParse(new ReadBufferByteBased(rawBytes));
        assertNotNull(bvlc);
        if (dumpPackages) System.out.println(bvlc);
        return bvlc;
    }

    static BVLC tryParseHex(String hex) throws ParseException, DecoderException {
        return tryParseHex(hex, DUMP_PACKAGES);
    }

    static BVLC tryParseHex(String hex, boolean dumpPackages) throws ParseException, DecoderException {
        return tryParseHex(hex, 0, dumpPackages);
    }

    static BVLC tryParseHex(String hex, int startIndex, boolean dumpPackages) throws ParseException, DecoderException {
        byte[] rawBytes = Hex.decodeHex(hex);
        rawBytes = ArrayUtils.subarray(rawBytes, startIndex, rawBytes.length);
        BVLC bvlc = BVLC.staticParse(new ReadBufferByteBased(rawBytes));
        assertNotNull(bvlc);
        if (dumpPackages) System.out.println(bvlc);
        return bvlc;
    }
}
