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
package org.apache.plc4x.java.firmata.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;

class FirmataTagHandlerTest {

    private FirmataTagHandler handler;

    @BeforeEach
    void setUp() {
        handler = new FirmataTagHandler();
    }

    @Test
    void parsesDigitalAddress() {
        PlcTag tag = handler.parseTag("digital:5");
        assertInstanceOf(FirmataTagDigital.class, tag);
    }

    @Test
    void parsesAnalogAddress() {
        PlcTag tag = handler.parseTag("analog:0");
        assertInstanceOf(FirmataTagAnalog.class, tag);
    }

    @Test
    void unknownPrefixThrows() {
        assertThrows(PlcInvalidTagException.class, () -> handler.parseTag("uart:0"));
    }

    @Test
    void queryParsingNotSupported() {
        assertThrows(UnsupportedOperationException.class, () -> handler.parseQuery("*"));
    }

}
