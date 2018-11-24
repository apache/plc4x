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
package org.apache.plc4x.java.simulated.connection;

import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultLongFieldItem;
import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.junit.Assert.*;

public class TestDeviceTest {

    @Test
    public void random() {
        TestDevice device = new TestDevice("foobar");
        TestField field = TestField.of("RANDOM/foo:INTEGER");

        Optional<BaseDefaultFieldItem> value = device.get(field);

        assertTrue(value.isPresent());
    }

    @Test
    public void read() {
        TestDevice device = new TestDevice("foobar");
        TestField field = TestField.of("STATE/bar:INTEGER");

        Optional<BaseDefaultFieldItem> value = device.get(field);
        assertFalse(value.isPresent());

        device.set(field, new DefaultLongFieldItem((long) 42));
        value = device.get(field);
        assertTrue(value.isPresent());
        BaseDefaultFieldItem fieldItem = value.get();
        assertEquals(42L, (long) fieldItem.getLong(0));
    }

}
