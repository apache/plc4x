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

package org.apache.plc4x.java.api.types;

import org.apache.plc4x.java.api.exceptions.PlcIncompatibleDatatypeException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValues;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class PlcValueTest {

    @Nested
    class Complex {

        @Test
        void complexTestCase_isComplex() {
            PlcValue value = PlcValues.of("Entry 1", PlcValues.of(
                PlcValues.of(true),
                PlcValues.of("Pimmel"),
                PlcValues.of(false),
                PlcValues.of("Arsch"),
                PlcValues.of(1278391)
            ));

            System.out.println(value);

            assertThrows(PlcIncompatibleDatatypeException.class, value::getBoolean);
            assertTrue(value.getValue("Entry 1").getIndex(0).getBoolean());
        }
    }

}