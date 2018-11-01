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
package org.apache.plc4x.java.test;

import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultIntegerFieldItem;
import org.junit.Test;
import org.mockito.Mockito;

import java.util.Optional;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

/**
 * Test the Mocking.
 *
 */
public class TestDeviceMockTest {

    @Test
    public void testMocking() {
        TestDevice device = new TestDevice("foobar");
        TestField field = TestField.of("MOCK/foo:INTEGER");

        MockDevice mock = Mockito.mock(MockDevice.class);
        device.setMockDevice(mock);

        // Read Field
        Optional<BaseDefaultFieldItem> value = device.get(field);
        // Set Field
        DefaultIntegerFieldItem setValue = new DefaultIntegerFieldItem(1);
        device.set(field, setValue);

        // Verify
        verify(mock, times(1)).get(eq("foo"), any());
        verify(mock, times(1)).set(eq("foo"), eq(setValue));
    }

}