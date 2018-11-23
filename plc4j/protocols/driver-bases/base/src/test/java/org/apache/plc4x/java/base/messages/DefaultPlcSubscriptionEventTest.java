/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.java.base.messages;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultByteFieldItem;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.HashMap;
import java.util.Map;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;

class DefaultPlcSubscriptionEventTest {

    private DefaultPlcSubscriptionEvent SUT;

    @BeforeEach
    void setUp() {
        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();
        fields.put("foo", new ImmutablePair<>(PlcResponseCode.OK, new DefaultByteFieldItem((byte) 0x42)));
        SUT = new DefaultPlcSubscriptionEvent(Instant.now(), fields);
    }

    @Test
    void getFieldNames() {
        assertThat(SUT.getFieldNames(), notNullValue());
        assertThat(SUT.getFieldNames().size(), equalTo(1));
        assertThat(SUT.getFieldNames().iterator().next(), equalTo("foo"));
    }

    @Test
    void getField() {
        assertThrows(UnsupportedOperationException.class, () -> SUT.getField("foo"));
    }

    @Test
    void getTimestamp() {
        Instant timestamp = SUT.getTimestamp();
        assertThat(timestamp, notNullValue());
    }

}