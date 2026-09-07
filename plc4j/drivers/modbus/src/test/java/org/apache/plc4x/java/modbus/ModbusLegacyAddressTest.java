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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.modbus.base.tag.ModbusTag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Addresses written before the array notation was unified must not parse.
 *
 * <p>The brackets moved from after the type to before it, which is what turns an otherwise
 * silent change of meaning into a failure: {@code [4]} used to mean four elements and now means
 * the fifth, so an unmodified address that still parsed would quietly return different data.
 */
class ModbusLegacyAddressTest {

    @Test
    void legacyForm0IsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> ModbusTag.of("holding-register:1:INT[4]"));
    }

    @Test
    void legacyForm1IsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> ModbusTag.of("coil:1:BOOL[8]"));
    }

    @Test
    void theReplacementFormParses() {
        assertNotNull(ModbusTag.of("holding-register:1[0..3]:INT"));
    }

    /** The message has to hand an upgrading user the address to write, not just a regex. */
    @Test
    void theErrorNamesTheReplacementAddress() {
        PlcInvalidTagException thrown =
            assertThrows(PlcInvalidTagException.class, () -> ModbusTag.of("holding-register:1:INT[4]"));
        assertTrue(thrown.getMessage().contains("holding-register:1[0..3]:INT"), thrown::getMessage);
    }
}
