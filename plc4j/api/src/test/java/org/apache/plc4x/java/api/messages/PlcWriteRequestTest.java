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
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.messages.items.WriteRequestItem;
import org.apache.plc4x.java.api.model.Address;
import org.junit.Before;
import org.junit.Test;

import java.util.Collections;

import static org.junit.Assert.fail;
import static org.mockito.Mockito.mock;

public class PlcWriteRequestTest {

    Address dummyAddress;

    @Before
    public void setUp() {
        dummyAddress = mock(Address.class);
    }

    @Test
    public void constructor() {
        new PlcWriteRequest();
        new PlcWriteRequest(new WriteRequestItem<>(String.class, dummyAddress, ""));
        new PlcWriteRequest(String.class, dummyAddress);
        new PlcWriteRequest(Collections.singletonList(new WriteRequestItem<>(String.class, dummyAddress)));
    }

    @Test
    public void builder() {
        { // empty
            try {
                PlcWriteRequest.builder().build();
                fail("An empty build should fail.");
            } catch (IllegalStateException e) {
                // expected
            }
        }
        { // one item implicit type
            PlcWriteRequest.builder()
                .addItem(dummyAddress, "")
                .build();
        }
        { // one item
            PlcWriteRequest.builder()
                .addItem(String.class, dummyAddress)
                .build();
        }
        { // one item prebuild
            PlcWriteRequest.builder()
                .addItem(new WriteRequestItem<>(String.class, dummyAddress))
                .build();
        }
        { // two different item
            PlcWriteRequest.builder()
                .addItem(String.class, dummyAddress)
                .addItem(Byte.class, dummyAddress)
                .build();
        }
        { // two different item typeSafe
            try {
                PlcWriteRequest.builder()
                    .addItem(String.class, dummyAddress)
                    .addItem(Byte.class, dummyAddress)
                    .build(String.class);
                fail("Mixed types build should fail.");
            } catch (IllegalStateException e) {
                // expected
            }
        }
        { // two different item typeSafe
            try {
                PlcWriteRequest.builder()
                    .addItem(String.class, dummyAddress)
                    .addItem(Byte.class, dummyAddress)
                    .build(Byte.class);
                fail("Mismatch of types should have failed.");
            } catch (ClassCastException e) {
                // expected
            }
        }
        { // two equal item typeSafe
            PlcWriteRequest.builder()
                .addItem(Byte.class, dummyAddress)
                .addItem(Byte.class, dummyAddress)
                .build(Byte.class);
        }
    }

}