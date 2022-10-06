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

package org.apache.plc4x.java.spi.utils.hex;

import org.apache.commons.lang3.StringUtils;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import org.apache.commons.lang3.tuple.Pair;

import static org.apache.plc4x.java.spi.utils.hex.Hex.*;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class HexTest {

    @BeforeEach
    void setUp() {
        DebugHex = true;
    }

    @Nested
    class Dump {
        @Test
        void testDump() {
            byte[] data = "1234567890abcdefghijklmnopqrstuvwxyz".getBytes();
            String want = "" +
                "00|31 32 33 34 35 36 37 38 39 30 '1234567890'\n" +
                "10|61 62 63 64 65 66 67 68 69 6a 'abcdefghij'\n" +
                "20|6b 6c 6d 6e 6f 70 71 72 73 74 'klmnopqrst'\n" +
                "30|75 76 77 78 79 7a             'uvwxyz    '";
            assertEquals(want, dump(data));
        }

        @Test
        void testBiggerDump() {
            byte[] data = StringUtils.repeat("Lorem ipsum", 90).getBytes();
            String want = "" +
                "000|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "010|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "020|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "030|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "040|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "050|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "060|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "070|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "080|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "090|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "100|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "110|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "120|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "130|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "140|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "150|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "160|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "170|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "180|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "190|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "200|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "210|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "220|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "230|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "240|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "250|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "260|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "270|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "280|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "290|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "300|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "310|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "320|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "330|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "340|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "350|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "360|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "370|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "380|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "390|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "400|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "410|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "420|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "430|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "440|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "450|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "460|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "470|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "480|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "490|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "500|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "510|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "520|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "530|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "540|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "550|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "560|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "570|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "580|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "590|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "600|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "610|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "620|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "630|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "640|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "650|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "660|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "670|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "680|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "690|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "700|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "710|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "720|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "730|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "740|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "750|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "760|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "770|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "780|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "790|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "800|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "810|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "820|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "830|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "840|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "850|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "860|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "870|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'\n" +
                "880|4c 6f 72 65 6d 20 69 70 73 75 'Lorem ipsu'\n" +
                "890|6d 4c 6f 72 65 6d 20 69 70 73 'mLorem ips'\n" +
                "900|75 6d 4c 6f 72 65 6d 20 69 70 'umLorem ip'\n" +
                "910|73 75 6d 4c 6f 72 65 6d 20 69 'sumLorem i'\n" +
                "920|70 73 75 6d 4c 6f 72 65 6d 20 'psumLorem '\n" +
                "930|69 70 73 75 6d 4c 6f 72 65 6d 'ipsumLorem'\n" +
                "940|20 69 70 73 75 6d 4c 6f 72 65 ' ipsumLore'\n" +
                "950|6d 20 69 70 73 75 6d 4c 6f 72 'm ipsumLor'\n" +
                "960|65 6d 20 69 70 73 75 6d 4c 6f 'em ipsumLo'\n" +
                "970|72 65 6d 20 69 70 73 75 6d 4c 'rem ipsumL'\n" +
                "980|6f 72 65 6d 20 69 70 73 75 6d 'orem ipsum'";
            assertEquals(want, dump(data));
        }

        @Test
        void noSize() {
            byte[] data = "a".getBytes();
            String want = "" +
                "0|61                            'a         '";
            assertEquals(want, dump(data));
        }
    }

    @Nested
    class FixedWidth {
        @Test
        void testDump() {
            byte[] data = "1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab".getBytes(StandardCharsets.ISO_8859_1);
            String want = "" +
                "000|31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 '1234567890abcdefghijklmnop'\n" +
                "026|71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 'qrstuvwxyz.1234567890abcde'\n" +
                "052|66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 61 31 32 'fghijklmnopqrstuvwxyz.aa12'\n" +
                "078|33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 '34567890abcdefghijklmnopqr'\n" +
                "104|73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 'stuvwxyz.1234567890abcdefg'\n" +
                "130|68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 61 62          'hijklmnopqrstuvwxyz.aab   '";
            assertEquals(want, dump(data, 110));
        }

        @Test
        void mimimumSize() {
            byte[] data = "1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aa1234567890abcdefghijklmnopqrstuvwxyz\3231234567890abcdefghijklmnopqrstuvwxyz\323aab".getBytes(StandardCharsets.ISO_8859_1);
            String want = "" +
                "000|31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 '1234567890abcdefghijklmnop'\n" +
                "026|71 72 73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 'qrstuvwxyz.1234567890abcde'\n" +
                "052|66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 61 31 32 'fghijklmnopqrstuvwxyz.aa12'\n" +
                "078|33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 68 69 6a 6b 6c 6d 6e 6f 70 71 72 '34567890abcdefghijklmnopqr'\n" +
                "104|73 74 75 76 77 78 79 7a d3 31 32 33 34 35 36 37 38 39 30 61 62 63 64 65 66 67 'stuvwxyz.1234567890abcdefg'\n" +
                "130|68 69 6a 6b 6c 6d 6e 6f 70 71 72 73 74 75 76 77 78 79 7a d3 61 61 62          'hijklmnopqrstuvwxyz.aab   '";
            assertEquals(want, dump(data, 110));
        }
    }

    @Nested
    class MaskString {
        @Test
        void printableChars() {
            byte[] data = "1234567890abcdefghijklmnopqrstuvwxyz".getBytes(StandardCharsets.ISO_8859_1);
            String want = "1234567890abcdefghijklmnopqrstuvwxyz";
            assertEquals(want, maskString(data));
        }

        @Test
        void unPrintableChars() {
            byte[] data = "\3231234567890abcdefghijklmnopqrstuvwxyz.\323".getBytes(StandardCharsets.ISO_8859_1);
            String want = ".1234567890abcdefghijklmnopqrstuvwxyz..";
            assertEquals(want, maskString(data));
        }
    }

    @Nested
    class MinMax {
        @Test
        void nullData() {
            assertEquals("", dump(null, Integer.MIN_VALUE));
        }

        @Test
        void emptyData() {
            assertEquals("", dump(new byte[0], Integer.MIN_VALUE));
        }

        @Test
        void minus1Data() {
            assertEquals("0|01 '.'", dump(new byte[]{0x1}, -1));
        }
    }

    @Nested
    class CalculateBytesPerRowAndIndexWidth {
        @Test
        void OneByteMinIntWidth() {
            assertEquals(Pair.of(1, 1), calculateBytesPerRowAndIndexWidth(1, Integer.MIN_VALUE));
        }

        @Test
        void TenByteMinIntWidth() {
            assertEquals(Pair.of(1, 2), calculateBytesPerRowAndIndexWidth(10, Integer.MIN_VALUE));
        }

        @Test
        void HundredByteMinIntWidth() {
            assertEquals(Pair.of(1, 3), calculateBytesPerRowAndIndexWidth(100, Integer.MIN_VALUE));
        }

        @Test
        void HundredByteTwelveWidth() {
            assertEquals(Pair.of(1, 3), calculateBytesPerRowAndIndexWidth(100, 12));
        }

        @Test
        void HundredFiftyThreeByteHundredThirtySixWidth() {
            assertEquals(Pair.of(32, 3), calculateBytesPerRowAndIndexWidth(153, 136));
        }

        @Test
        void HundredFiftyThreeByteCalculatedWidth() {
            int quoteRune = 1;
            int numberOfBytes = 153;
            int indexWidth = 3;
            int charRepresentation = 1;
            // 000 AF FE AF FE ..... '....*'
            int calculated = indexWidth + blankWidth + (numberOfBytes * byteWidth) + quoteRune + (numberOfBytes * charRepresentation) + quoteRune;
            assertEquals(Pair.of(153, 3), calculateBytesPerRowAndIndexWidth(153, calculated));
        }
    }
}
