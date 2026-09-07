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
package org.apache.plc4x.java.iec608705104.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.params.provider.ValueSource;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class Iec608705104TagTest {

    @Nested
    class Syntax {

        /**
         * Every accepted form, with the ASDU / information object address the
         * octets fold into. Octets are written least significant first, in the
         * same order they appear on the wire.
         */
        @ParameterizedTest
        @CsvSource({
            // address,          asdu,  ioa
            "0/0,                0,     0",
            "1/2,                1,     2",
            "65535/16777215,     65535, 16777215",
            // <low>/<high> ASDU address: 1 + 2 * 256
            "1/2/3,              513,   3",
            "255/255/0,          65535, 0",
            // <low>.<middle>.<high> information object address: 2 + 3 * 256 + 4 * 65536
            "1/2.3.4,            1,     262914",
            "1/0.0.255,          1,     16711680",
            // both fields octet-wise
            "1/2/3.4.5,          513,   328707",
        })
        void parsesEveryAcceptedForm(String address, int expectedAdsuAddress, int expectedObjectAddress) {
            Iec608705104Tag tag = Iec608705104Tag.of(address);
            assertEquals(expectedAdsuAddress, tag.getAdsuAddress());
            assertEquals(expectedObjectAddress, tag.getObjectAddress());
            assertFalse(tag.isWildcarded());
            assertTrue(Iec608705104Tag.isValidAddress(address));
        }

        @Test
        void leadingZerosAndSurroundingWhitespaceAreTolerated() {
            Iec608705104Tag tag = Iec608705104Tag.of("  007/002.03.004  ");
            assertEquals(7, tag.getAdsuAddress());
            assertEquals(2 + 3 * 256 + 4 * 65536, tag.getObjectAddress());
            assertEquals("7/2.3.4", tag.getAddressString());
        }

        @ParameterizedTest
        @ValueSource(strings = {
            "000007/2",         // one zero too many for a whole ASDU address
            "7/0002.3.4",       // ... and for an octet
            "000065535/0",
            "0/000016777215",   // ... and for a whole information object address
        })
        void paddingBeyondTheWidthOfTheMaximumIsRejected(String address) {
            // The digits of a component are capped at the width of its maximum
            // (3 for an octet, 5 for an ASDU address, 8 for an information
            // object address). The Go driver caps them identically, so both
            // languages accept the same set of addresses.
            assertFalse(Iec608705104Tag.isValidAddress(address));
            assertThrows(PlcInvalidTagException.class, () -> Iec608705104Tag.of(address));
        }

        @Test
        void equivalentNotationsDescribeTheSameTag() {
            // 1/2 as octets is 513, 3.4.5 as octets is 328707.
            assertEquals(Iec608705104Tag.of("513/328707"), Iec608705104Tag.of("1/2/3.4.5"));
            assertEquals(Iec608705104Tag.of("513/328707").hashCode(), Iec608705104Tag.of("1/2/3.4.5").hashCode());
            assertEquals(new Iec608705104Tag(513, 3), Iec608705104Tag.of("1/2/3"));
        }

        @Test
        void aWildcardIsNotTheSameAsAZero() {
            assertNotEquals(Iec608705104Tag.of("1/*"), Iec608705104Tag.of("1/0"));
            assertNotEquals(Iec608705104Tag.of("*/1"), Iec608705104Tag.of("0/1"));
            assertNotEquals(Iec608705104Tag.of("*/*/1"), Iec608705104Tag.of("0/0/1"));
        }

        @Test
        void tagsAreOnlyEqualToOtherTags() {
            Iec608705104Tag tag = Iec608705104Tag.of("1/2");
            assertEquals(tag, tag);
            assertNotEquals(tag, "1/2");
            assertNotEquals(tag, null);
        }
    }

    @Nested
    class Wildcards {

        @Test
        void fullWildcardMatchesEverything() {
            Iec608705104Tag tag = Iec608705104Tag.of("*/*");
            assertTrue(tag.isWildcarded());
            assertTrue(tag.isAdsuAddressWildcarded());
            assertTrue(tag.isObjectAddressWildcarded());
            assertEquals(Iec608705104Tag.WILDCARD_ADDRESS, tag.getAdsuAddress());
            assertEquals(Iec608705104Tag.WILDCARD_ADDRESS, tag.getObjectAddress());
            assertEquals(0, tag.getAdsuAddressMask());
            assertEquals(0, tag.getObjectAddressMask());
            assertTrue(tag.matches(0, 0));
            assertTrue(tag.matches(65535, 16777215));
            assertTrue(tag.matches(1234, 5678));
        }

        @Test
        void wildcardOnTheObjectAddressPinsTheStation() {
            Iec608705104Tag tag = Iec608705104Tag.of("1/*");
            assertFalse(tag.isAdsuAddressWildcarded());
            assertTrue(tag.isObjectAddressWildcarded());
            assertEquals(1, tag.getAdsuAddress());
            assertEquals(Iec608705104Tag.WILDCARD_ADDRESS, tag.getObjectAddress());
            assertTrue(tag.matches(1, 0));
            assertTrue(tag.matches(1, 16777215));
            assertFalse(tag.matches(2, 0));
        }

        @Test
        void wildcardOnTheAsduAddressPinsTheObject() {
            Iec608705104Tag tag = Iec608705104Tag.of("*/3.4.5");
            assertTrue(tag.isAdsuAddressWildcarded());
            assertFalse(tag.isObjectAddressWildcarded());
            assertEquals(328707, tag.getObjectAddress());
            assertTrue(tag.matches(0, 328707));
            assertTrue(tag.matches(65535, 328707));
            assertFalse(tag.matches(0, 328708));
        }

        @Test
        void singleAsduOctetCanBeWildcarded() {
            // Only the high octet is pinned: every ASDU address in 512..767.
            Iec608705104Tag tag = Iec608705104Tag.of("*/2/7");
            assertTrue(tag.isAdsuAddressWildcarded());
            assertEquals(0xFF00, tag.getAdsuAddressMask());
            assertTrue(tag.matches(512, 7));
            assertTrue(tag.matches(767, 7));
            assertFalse(tag.matches(768, 7));
            assertFalse(tag.matches(511, 7));
            assertFalse(tag.matches(512, 8));
        }

        @Test
        void singleObjectOctetCanBeWildcarded() {
            // low = 3, high = 5, middle octet free.
            Iec608705104Tag tag = Iec608705104Tag.of("1/2/3.*.5");
            assertEquals(0xFF00FF, tag.getObjectAddressMask());
            assertTrue(tag.matches(513, 3 + 0 * 256 + 5 * 65536));
            assertTrue(tag.matches(513, 3 + 255 * 256 + 5 * 65536));
            assertFalse(tag.matches(513, 4 + 0 * 256 + 5 * 65536));
            assertFalse(tag.matches(513, 3 + 0 * 256 + 6 * 65536));
            assertFalse(tag.matches(514, 3 + 0 * 256 + 5 * 65536));
        }

        @Test
        void fullySpecifiedTagMatchesOnlyItsOwnObject() {
            Iec608705104Tag tag = Iec608705104Tag.of("1/2");
            assertTrue(tag.matches(1, 2));
            assertFalse(tag.matches(1, 3));
            assertFalse(tag.matches(2, 2));
        }
    }

    @Nested
    class RoundTrip {

        @ParameterizedTest
        @ValueSource(strings = {
            "0/0", "1/2", "65535/16777215",
            "1/2/3", "1/2.3.4", "1/2/3.4.5",
            "*/*", "1/*", "*/2", "*/*/3", "1/*/2.3.4", "1/2/3.*.5", "*/*/*.*.*"
        })
        void addressStringReparsesToAnEqualTag(String address) {
            Iec608705104Tag tag = Iec608705104Tag.of(address);
            assertEquals(address, tag.getAddressString());
            Iec608705104Tag reparsed = Iec608705104Tag.of(tag.getAddressString());
            assertEquals(tag, reparsed);
            assertEquals(tag.getAddressString(), reparsed.getAddressString());
            assertEquals(tag.getAdsuAddressMask(), reparsed.getAdsuAddressMask());
            assertEquals(tag.getObjectAddressMask(), reparsed.getObjectAddressMask());
        }

        @Test
        void programmaticallyBuiltTagAlsoRoundTrips() {
            Iec608705104Tag tag = new Iec608705104Tag(513, 328707);
            assertEquals("513/328707", tag.getAddressString());
            assertEquals(tag, Iec608705104Tag.of(tag.getAddressString()));
        }
    }

    @Nested
    class Rejection {

        @ParameterizedTest
        @ValueSource(strings = {
            "",                 // empty
            "1",                // no information object address
            "1/",               // missing information object address
            "/2",               // missing ASDU address
            "1/2/3/4",          // one field too many
            "1/2.3",            // two octets is not a valid information object address
            "1.2/3",            // dots don't apply to the ASDU address
            "1/2.3.4.5",        // one octet too many
            "abc/1",            // not a number
            "1/abc",
            "-1/2",             // no signs
            "1/-2",
            "1 / 2",            // no inner whitespace
            "**/1",             // wildcard is a whole component
            "1*/2",
            "1//2",             // empty component
            "1/2..3",
            "0x1/2",            // decimal only
            "1/2/3.4.5.6",
        })
        void malformedAddressesAreRejected(String address) {
            assertFalse(Iec608705104Tag.isValidAddress(address), address + " must not be considered valid");
            assertFalse(Iec608705104Tag.isSyntacticallyValid(address), address + " must not even be well-formed");
            assertThrows(PlcInvalidTagException.class, () -> Iec608705104Tag.of(address));
        }

        @Test
        void nullIsRejected() {
            assertFalse(Iec608705104Tag.isValidAddress(null));
            assertFalse(Iec608705104Tag.isSyntacticallyValid(null));
            assertThrows(PlcInvalidTagException.class, () -> Iec608705104Tag.of(null));
        }

        @ParameterizedTest
        @ValueSource(strings = {
            "65536/0",          // ASDU address too large
            "0/16777216",       // information object address too large
            "256/0/0",          // ASDU octet too large
            "0/256/0",
            "1/256.0.0",        // information object octet too large
            "1/0.256.0",
            "1/0.0.256",
        })
        void outOfRangeAddressesAreRejected(String address) {
            // These are well-formed, so the syntax-only check lets them
            // through, but isValidAddress() must not: it guards of().
            assertTrue(Iec608705104Tag.isSyntacticallyValid(address));
            assertFalse(Iec608705104Tag.isValidAddress(address));
            assertThrows(PlcInvalidTagException.class, () -> Iec608705104Tag.of(address));
        }

        @ParameterizedTest
        @CsvSource({"-1, 0", "65536, 0", "0, -1", "0, 16777216"})
        void constructorRejectsOutOfRangeAddresses(int adsuAddress, int objectAddress) {
            assertThrows(PlcInvalidTagException.class, () -> new Iec608705104Tag(adsuAddress, objectAddress));
        }
    }

    @Nested
    class PlcTagContract {

        @Test
        void addressCarriesNoTypeInformation() {
            // The datatype follows from the type identification of the ASDU
            // that delivers the object, not from the address.
            assertEquals(PlcValueType.NULL, Iec608705104Tag.of("1/2").getPlcValueType());
        }

        @Test
        void arrayInfoFallsThroughToDefault() {
            assertNotNull(Iec608705104Tag.of("0/0").getArrayInfo());
        }

        @Test
        void toStringShowsTheAddress() {
            assertEquals("Iec608705104Tag{address='1/2/3.4.5'}", Iec608705104Tag.of("1/2/3.4.5").toString());
            assertEquals("Iec608705104Tag{address='3/7'}", new Iec608705104Tag(3, 7).toString());
        }
    }

}
