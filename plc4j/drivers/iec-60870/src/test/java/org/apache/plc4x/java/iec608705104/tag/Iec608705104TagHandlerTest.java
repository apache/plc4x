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

import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104TagHandler;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertNull;

class Iec608705104TagHandlerTest {

    @Test
    void parseTagYieldsPlaceholderTagUntilGrammarLands() {
        Iec608705104TagHandler handler = new Iec608705104TagHandler();
        Iec608705104Tag tag = assertInstanceOf(Iec608705104Tag.class, handler.parseTag("1/0/2"));
        assertEquals(0, tag.getAdsuAddress());
        assertEquals(0, tag.getObjectAddress());
    }

    @Test
    void parseQueryIsNotSupported() {
        // PlcBrowseRequest doesn't apply to push-only IEC-60870 — the handler
        // returns null to signal "no query syntax".
        assertNull(new Iec608705104TagHandler().parseQuery("anything"));
    }

}
