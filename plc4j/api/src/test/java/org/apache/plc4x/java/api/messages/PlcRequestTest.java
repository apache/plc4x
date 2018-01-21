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

import org.apache.plc4x.java.api.messages.items.RequestItem;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import static org.mockito.Mockito.mock;

class PlcRequestTest {

    private List<RequestItem> requestItems;

    private PlcRequest SUT;

    @BeforeEach
    void setUp() {
        requestItems = new ArrayList<>();
        SUT = new PlcRequest<RequestItem>(requestItems) {
        };
    }

    @Test
    void addItem() {
        SUT.addItem(mock(RequestItem.class));
    }

    @Test
    void getRequestItems() {
        Assertions.assertEquals(0, SUT.getRequestItems().size());
    }

    @Test
    void getRequestItem() {
        Assertions.assertEquals(Optional.empty(), SUT.getRequestItem());
        requestItems.add(mock(RequestItem.class));
        Assertions.assertTrue(SUT.getRequestItem().isPresent());
        requestItems.add(mock(RequestItem.class));
        Assertions.assertThrows(IllegalStateException.class, () -> SUT.getRequestItem());
    }

    @Test
    void setRequestItem() {
        SUT.setRequestItem(mock(RequestItem.class));
        requestItems.add(mock(RequestItem.class));
        Assertions.assertThrows(IllegalStateException.class, () -> SUT.setRequestItem(mock(RequestItem.class)));
    }

    @Test
    void getNumberOfItems() {
        Assertions.assertEquals(0, SUT.getNumberOfItems());
    }

    @Test
    void isMultiValue() {
        Assertions.assertFalse(SUT.isMultiValue());
    }

    @Test
    void isEmpty() {
        Assertions.assertTrue(SUT.isEmpty());
    }

}