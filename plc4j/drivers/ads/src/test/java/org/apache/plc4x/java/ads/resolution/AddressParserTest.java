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
package org.apache.plc4x.java.ads.resolution;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class AddressParserTest {

    @Test
    void scalarSymbolNoIndices() {
        AddressParser.AddressPart p = AddressParser.parse("MAIN.g_b1");
        assertEquals("MAIN.g_b1", p.baseSegment());
        assertTrue(p.arrayIndices().isEmpty());
        assertNull(p.child());
    }

    @Test
    void singleIndexOn1DArray() {
        AddressParser.AddressPart p = AddressParser.parse("MAIN.g_arr[3]");
        assertEquals("MAIN.g_arr", p.baseSegment());
        assertEquals(List.of(3), p.arrayIndices());
        assertNull(p.child());
    }

    @Test
    void multiIndexOnMultiDimArray() {
        AddressParser.AddressPart p = AddressParser.parse("MAIN.g_cube[1][2][3]");
        assertEquals("MAIN.g_cube", p.baseSegment());
        assertEquals(List.of(1, 2, 3), p.arrayIndices());
        assertNull(p.child());
    }

    @Test
    void structFieldChain() {
        AddressParser.AddressPart p = AddressParser.parse("MAIN.g_simple.s8");
        assertEquals("MAIN.g_simple", p.baseSegment());
        assertNotNull(p.child());
        assertEquals("s8", p.child().baseSegment());
        assertTrue(p.child().arrayIndices().isEmpty());
        assertNull(p.child().child());
    }

    @Test
    void deepMixedPath() {
        AddressParser.AddressPart p = AddressParser.parse("MAIN.g_chanGrid[1][1].lut[2][2][2]");
        assertEquals("MAIN.g_chanGrid", p.baseSegment());
        assertEquals(List.of(1, 1), p.arrayIndices());
        AddressParser.AddressPart child = p.child();
        assertNotNull(child);
        assertEquals("lut", child.baseSegment());
        assertEquals(List.of(2, 2, 2), child.arrayIndices());
        assertNull(child.child());
    }

    @Test
    void chainBetweenStructAndArray() {
        AddressParser.AddressPart p = AddressParser.parse("MAIN.g_plant.channels[1].setpoints[4]");
        AddressParser.AddressPart channels = p.child();
        AddressParser.AddressPart setpoints = channels.child();
        assertEquals("MAIN.g_plant", p.baseSegment());
        assertEquals("channels", channels.baseSegment());
        assertEquals(List.of(1), channels.arrayIndices());
        assertEquals("setpoints", setpoints.baseSegment());
        assertEquals(List.of(4), setpoints.arrayIndices());
    }

    @Test
    void noNamespacePrefix() {
        AddressParser.AddressPart p = AddressParser.parse("foo[2]");
        assertEquals("foo", p.baseSegment());
        assertEquals(List.of(2), p.arrayIndices());
    }

    @Test
    void emptyAddressRejected() {
        assertThrows(PlcInvalidTagException.class, () -> AddressParser.parse(""));
        assertThrows(PlcInvalidTagException.class, () -> AddressParser.parse(null));
    }

    @Test
    void unmatchedBracketRejected() {
        assertThrows(PlcInvalidTagException.class, () -> AddressParser.parse("MAIN.g_arr[1"));
    }

    @Test
    void nonNumericIndexRejected() {
        assertThrows(PlcInvalidTagException.class, () -> AddressParser.parse("MAIN.g_arr[x]"));
    }
}
