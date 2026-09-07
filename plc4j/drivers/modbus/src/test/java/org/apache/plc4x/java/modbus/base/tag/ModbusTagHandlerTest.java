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
package org.apache.plc4x.java.modbus.base.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ModbusTagHandlerTest {

    private ModbusTagHandler handler;

    @BeforeEach
    void setUp() {
        handler = new ModbusTagHandler();
    }

    @Test
    void testParseTag_discreteInput() {
        PlcTag tag = handler.parseTag("1x00001:BOOL");
        assertInstanceOf(ModbusTagDiscreteInput.class, tag);
    }

    @Test
    void testParseTag_coil() {
        PlcTag tag = handler.parseTag("0x00001:BOOL");
        assertInstanceOf(ModbusTagCoil.class, tag);
    }

    @Test
    void testParseTag_holdingRegister() {
        PlcTag tag = handler.parseTag("4x00001:INT");
        assertInstanceOf(ModbusTagHoldingRegister.class, tag);
    }

    @Test
    void testParseTag_inputRegister() {
        PlcTag tag = handler.parseTag("3x00001:INT");
        assertInstanceOf(ModbusTagInputRegister.class, tag);
    }

    @Test
    void testParseTag_extendedRegister() {
        PlcTag tag = handler.parseTag("6x00001:INT");
        assertInstanceOf(ModbusTagExtendedRegister.class, tag);
    }

    @Test
    void testParseTag_invalid() {
        assertThrows(PlcInvalidTagException.class, () -> handler.parseTag("invalid-address"));
    }

    @Test
    void testParseQuery_unsupported() {
        assertThrows(UnsupportedOperationException.class, () -> handler.parseQuery("some-query"));
    }
}
