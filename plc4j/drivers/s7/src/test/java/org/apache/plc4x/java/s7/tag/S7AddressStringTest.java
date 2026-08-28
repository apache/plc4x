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
package org.apache.plc4x.java.s7.tag;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import static org.junit.jupiter.api.Assertions.*;

/**
 * An S7 tag spells itself the way {@link S7Tag#of(String)} reads it back.
 *
 * <p>{@code getAddressString()} used to return {@code null}, so anything carrying a tag as a
 * string - a log line, a browse result, a serialized request - got nothing at all from an S7
 * tag.</p>
 */
class S7AddressStringTest {

    @ParameterizedTest
    @ValueSource(strings = {
        "%M100:INT",
        "%M100[0..9]:INT",
        "%Q0.0:BOOL",
        "%I7.3:BOOL",
        "%DB1.DB20:INT",
        "%DB1.DB20[0..3]:INT",
        "%DB1.DB0.0:BOOL",
        "%DB42.DB28[0..7]:BYTE",
        "%DB1.DB0:STRING(20)",
        "%DB1.DB0[0..2]:STRING(20)",
        "%DB1.DB0:WSTRING(10)",
        "%DB69.DB68:STRING",
    })
    void anAddressRendersAsItselfAndParsesBack(String address) {
        S7Tag tag = S7Tag.of(address);

        assertEquals(address, tag.getAddressString(), "rendered form");
        assertEquals(tag, S7Tag.of(tag.getAddressString()), "re-parsed tag");
    }

    /**
     * A COUNTER address names a counter, which the constructor splits across the byte and bit
     * offsets. Rendering has to put it back together, or the address comes back naming a
     * different counter - 100 stored as byte 12 bit 4 would render as counter 12.
     */
    @Test
    void aCounterAddressSurvivesTheSplitItIsStoredIn() {
        S7Tag counter = S7Tag.of("%DB1.DB100:COUNTER");

        assertEquals("%DB1.DB100:COUNTER", counter.getAddressString());
        assertEquals(counter, S7Tag.of(counter.getAddressString()));
    }

    /**
     * The optional transfer size code only repeats what the type already says, so it is not part
     * of the canonical form. The address it renders to still parses to the same tag.
     */
    @Test
    void theTransferSizeCodeIsNotPartOfTheCanonicalForm() {
        S7Tag tag = S7Tag.of("%DB69.DBX68[0..2]:WSTRING(254)");

        assertEquals("%DB69.DB68[0..2]:WSTRING(254)", tag.getAddressString());
        assertEquals(tag, S7Tag.of(tag.getAddressString()));
    }

    /**
     * A data block is rendered as a data block. Spelling it through the memory area's short name
     * would produce "%D100", which parses back as block 0 of the data-block area - a different
     * address that reads different memory.
     */
    @Test
    void aDataBlockIsNotRenderedThroughItsShortName() {
        S7Tag tag = S7Tag.of("%DB42.DB28:BYTE");

        assertEquals("%DB42.DB28:BYTE", tag.getAddressString());
        assertEquals(42, S7Tag.of(tag.getAddressString()).getBlockNumber());
    }

    /** A declared base is resolved into the address, so what renders back is the resolved one. */
    @Test
    void aDeclaredBaseIsResolvedIntoTheRenderedAddress() {
        S7Tag tag = S7Tag.of("%DB1.DB20[4..7;4]:INT");

        assertEquals("%DB1.DB20[0..3]:INT", tag.getAddressString());
        assertEquals(tag, S7Tag.of(tag.getAddressString()));
    }
}
