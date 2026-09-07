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
package org.apache.plc4x.java.iec608705104.messages;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.iec608705104.messages.Iec608705104PlcSubscriptionEvent;
import org.apache.plc4x.java.iec608705104.tag.Iec608705104Tag;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertThrows;

class Iec608705104PlcSubscriptionEventTest {

    @Test
    void getTagReturnsTheBoundPlcTag() {
        // The default subscription event throws when asked for the tag — we
        // override it specifically so listeners can recover the ASDU/IOA
        // mapping without re-parsing the subscription name.
        Iec608705104Tag tag = new Iec608705104Tag(1, 2);
        PlcValue value = PlcBOOL.of(true);
        PlcResponseItem<PlcValue> item = new DefaultPlcResponseItem<>(PlcResponseCode.OK, value);
        Iec608705104PlcSubscriptionEvent event = new Iec608705104PlcSubscriptionEvent(
            Instant.EPOCH, Map.of("a", tag), Map.of("a", item));

        assertSame(tag, event.getTag("a"));
        assertNull(event.getTag("missing"));
        assertEquals(Instant.EPOCH, event.getTimestamp());
    }

    @Test
    void rejectsNullTagsMap() {
        // Without the tag map, getTag() falls back to the superclass which
        // throws — caller bugs should fail fast at construction instead.
        assertThrows(NullPointerException.class,
            () -> new Iec608705104PlcSubscriptionEvent(Instant.EPOCH, null, Map.of()));
    }

}
