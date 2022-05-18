package org.apache.plc4x.java.bacnetip;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.plc4x.java.bacnetip.readwrite.BVLC;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;

import java.util.stream.IntStream;

import static org.junit.jupiter.api.Assertions.assertNotNull;

public class Utils {
    static final boolean DUMP_PACKAGES = false;

    static int PAYLOAD_START_INDEX = 42;

    static void tryParseBytes(int[] rawBytesAsInts) throws ParseException {
        tryParseBytes(rawBytesAsInts, PAYLOAD_START_INDEX);
    }

    static void tryParseBytes(int[] rawBytesAsInts, int startIndex) throws ParseException {
        tryParseBytes(rawBytesAsInts, PAYLOAD_START_INDEX, DUMP_PACKAGES);
    }

    static void tryParseBytes(int[] rawBytesAsInts, int startIndex, boolean dumpPackages) throws ParseException {
        var rawBytes = (byte[]) ArrayUtils.toPrimitive(IntStream.of(rawBytesAsInts).boxed().map(Integer::byteValue).toArray(Byte[]::new));
        rawBytes = ArrayUtils.subarray(rawBytes, startIndex, rawBytes.length);
        BVLC bvlc = BVLC.staticParse(new ReadBufferByteBased(rawBytes));
        assertNotNull(bvlc);
        if (dumpPackages) System.out.println(bvlc);
    }
}
