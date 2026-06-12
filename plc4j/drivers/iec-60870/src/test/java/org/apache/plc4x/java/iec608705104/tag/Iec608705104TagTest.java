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

import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class Iec608705104TagTest {

    @Test
    void carriesAdsuAndObjectAddress() {
        Iec608705104Tag tag = new Iec608705104Tag(10, 42);
        assertEquals(10, tag.getAdsuAddress());
        assertEquals(42, tag.getObjectAddress());
    }

    @Test
    void plcValueTypeIsNullUntilTagSyntaxIsImplemented() {
        // The tag-string grammar (see the commented regex in Iec608705104Tag)
        // hasn't been wired through yet, so the type and address-string are
        // intentionally placeholders. This test pins the current contract so
        // future work can break it on purpose, not by accident.
        Iec608705104Tag tag = new Iec608705104Tag(1, 2);
        assertEquals(PlcValueType.NULL, tag.getPlcValueType());
        assertNull(tag.getAddressString());
    }

    @Test
    void arrayInfoFallsThroughToDefault() {
        assertNotNull(new Iec608705104Tag(0, 0).getArrayInfo());
    }

    @Test
    void toStringIncludesBothAddresses() {
        String s = new Iec608705104Tag(3, 7).toString();
        assertEquals("Iec608705104Tag{adsuAddress=3, objectAddress=7}", s);
    }

}
