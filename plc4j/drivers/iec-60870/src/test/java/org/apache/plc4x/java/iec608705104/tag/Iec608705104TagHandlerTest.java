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
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.junit.jupiter.params.provider.ValueSource;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;

class Iec608705104TagHandlerTest {

    private final Iec608705104TagHandler handler = new Iec608705104TagHandler();

    @ParameterizedTest
    @CsvSource({
        "1/2,       1,     2",
        "1/2/3,     513,   3",
        "1/2.3.4,   1,     262914",
        "1/2/3.4.5, 513,   328707",
    })
    void parseTagResolvesTheRealAddresses(String address, int expectedAdsuAddress, int expectedObjectAddress) {
        Iec608705104Tag tag = assertInstanceOf(Iec608705104Tag.class, handler.parseTag(address));
        assertEquals(expectedAdsuAddress, tag.getAdsuAddress());
        assertEquals(expectedObjectAddress, tag.getObjectAddress());
        assertEquals(address, tag.getAddressString());
    }

    @Test
    void parseTagKeepsWildcards() {
        Iec608705104Tag tag = assertInstanceOf(Iec608705104Tag.class, handler.parseTag("1/*"));
        assertEquals(1, tag.getAdsuAddress());
        assertEquals(Iec608705104Tag.WILDCARD_ADDRESS, tag.getObjectAddress());
    }

    @ParameterizedTest
    @ValueSource(strings = {"", "1", "1/2/3/4", "1/2.3", "nope", "65536/0"})
    void parseTagRejectsMalformedAddresses(String address) {
        assertThrows(PlcInvalidTagException.class, () -> handler.parseTag(address));
    }

    @Test
    void parseQueryFailsLoudlyBecauseIec60870CannotBeBrowsed() {
        // Returning null here used to make PlcBrowseRequest builders NPE far
        // away from the cause.
        assertThrows(UnsupportedOperationException.class, () -> handler.parseQuery("1/*"));
    }

}
