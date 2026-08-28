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
package org.apache.plc4x.java.eip.base.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.eip.base.tag.EipTag.MemberElement;
import org.apache.plc4x.java.eip.base.tag.EipTag.PathElement;
import org.apache.plc4x.java.eip.base.tag.EipTag.SymbolElement;
import org.apache.plc4x.java.eip.readwrite.CIPDataTypeCode;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Field;
import java.lang.reflect.Modifier;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * The CIP path a tag address describes. Each member of a structured address is its own path
 * segment on the wire - a Logix controller walks them one at a time - so an address has to be
 * decomposed rather than sent as one symbol whose name happens to contain dots.
 *
 * <p>The decomposition happens once, when the tag is built, so that building a request does not
 * re-parse the address for every read.
 */
class EipTagPathTest {

    @Test
    void plainTagIsASingleSymbol() {
        assertEquals(List.of(new SymbolElement("myTag")), path("myTag:DINT"));
    }

    /** A leading '%' is not part of the symbol name. */
    @Test
    void percentPrefixIsNotPartOfTheSymbol() {
        assertEquals(List.of(new SymbolElement("N40")), path("%N40:DINT"));
    }

    @Test
    void arrayElementIsASymbolFollowedByItsMember() {
        assertEquals(
            List.of(new SymbolElement("hurz_DINT_ARR"), new MemberElement((short) 3)),
            path("hurz_DINT_ARR[3]:DINT"));
    }

    /**
     * The case a single-symbol encoding gets wrong: every member of a structured address is its
     * own segment, so 'myStruct.member' must not be sent as one symbol named "myStruct.member".
     */
    @Test
    void structuredAddressIsOneSymbolPerMember() {
        assertEquals(
            List.of(new SymbolElement("myStruct"), new SymbolElement("member")),
            path("myStruct.member:REAL"));

        assertEquals(
            List.of(new SymbolElement("Program"), new SymbolElement("MainProgram"), new SymbolElement("MyTag")),
            path("Program.MainProgram.MyTag:DINT"));
    }

    @Test
    void structuredAddressCanCarryAnIndex() {
        assertEquals(
            List.of(new SymbolElement("a"), new SymbolElement("b"), new MemberElement((short) 2)),
            path("a.b[2]:DINT"));
    }

    /**
     * The decomposition keeps the members in the order they were written, including an index
     * that is not the last element. The address pattern only accepts a trailing index, so
     * of() never produces this - but a tag built directly still describes a valid CIP path.
     */
    @Test
    void indexBeforeAFurtherMemberKeepsItsPosition() {
        assertNull(EipTag.of("a[2].b:DINT"), "precondition: the address pattern rejects this form");

        assertEquals(
            List.of(new SymbolElement("a"), new MemberElement((short) 2), new SymbolElement("b")),
            new EipTag("a[2].b", CIPDataTypeCode.DINT).getPathElements());
    }

    /** Empty brackets name no element and are rejected rather than quietly ignored. */
    @Test
    void emptyBracketsAreRejected() {
        assertNull(EipTag.of("a[]:DINT"));
    }

    @Test
    void pathIsUnmodifiable() {
        List<PathElement> elements = EipTag.of("a.b:DINT").getPathElements();
        assertThrows(UnsupportedOperationException.class, () -> elements.add(new SymbolElement("c")));
    }

    // --- the tag itself ---

    /** getTag() keeps returning the address as written; the decomposition is separate. */
    @Test
    void getTagIsUnchangedByTheDecomposition() {
        assertEquals("hurz_DINT_ARR[3]", EipTag.of("hurz_DINT_ARR[3]:DINT").getTag());
        assertEquals("a.b", EipTag.of("a.b:DINT").getTag());
    }

    /**
     * Every other driver's tag is immutable; this one was the exception. Element counts below one
     * are also normalised here rather than at every use site.
     */
    @Test
    void tagIsImmutable() {
        for (Field field : EipTag.class.getDeclaredFields()) {
            if (!Modifier.isStatic(field.getModifiers())) {
                assertTrue(Modifier.isFinal(field.getModifiers()), () -> field.getName() + " should be final");
            }
        }
        assertEquals(0, java.util.Arrays.stream(EipTag.class.getMethods())
            .filter(m -> m.getName().startsWith("set"))
            .count(), "EipTag should expose no setters");
    }

    @Test
    void elementCountIsNeverBelowOne() {
        assertEquals(1, new EipTag("a", CIPDataTypeCode.DINT, 0).getElementNb());
        assertEquals(1, new EipTag("a", CIPDataTypeCode.DINT, -5).getElementNb());
        assertEquals(1, EipTag.of("a:DINT").getElementNb());
        assertEquals(8, EipTag.of("a[0..7]:DINT").getElementNb());
    }

    /**
     * An index rides in a MemberID, whose instance field the mspec declares as uint 8. A larger
     * index is rejected while the address is parsed rather than failing later when the request
     * is serialized - and never silently truncated to a different element.
     */
    @Test
    void indexBeyondAMemberIdIsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> EipTag.of("a[256]:DINT"));
        assertThrows(PlcInvalidTagException.class, () -> EipTag.of("a[300]:DINT"));
        assertThrows(PlcInvalidTagException.class, () -> EipTag.of("a[99999]:DINT"));

        assertEquals(
            List.of(new SymbolElement("a"), new MemberElement((short) 255)),
            path("a[255]:DINT"));
    }

    // --- what a consumer sees ---

    /**
     * getArrayInfo() tells a consumer whether it received a scalar or a list. A bare index
     * selects one element, so it reports empty; a range reports its dimensions even when it
     * spans one element. Reading that one element still walks a member path, which is a
     * separate concern.
     */
    @Test
    void arrayInfoDescribesTheValueNotTheFetch() {
        assertTrue(EipTag.of("myArray[4]:DINT").getArrayInfo().isEmpty(), "a bare index is a scalar");
        assertTrue(EipTag.of("myTag:DINT").getArrayInfo().isEmpty(), "no selection is a scalar");

        assertEquals(1, EipTag.of("myArray[4..4]:DINT").getArrayInfo().size(), "a range is an array");
        assertEquals(8, EipTag.of("myArray[0..7]:DINT").getArrayInfo().get(0).getSize());

        // The single element still needs its member segment on the wire.
        assertEquals(
            List.of(new SymbolElement("myArray"), new MemberElement((short) 4)),
            EipTag.of("myArray[4]:DINT").getPathElements());
    }

    private static List<PathElement> path(String address) {
        EipTag tag = EipTag.of(address);
        assertNotNull(tag, address);
        return tag.getPathElements();
    }
}
