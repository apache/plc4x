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
import org.apache.plc4x.java.api.messages.items.ResponseItem;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import static org.mockito.Mockito.mock;

class PlcResponseTest {

    private List<ResponseItem> responseItems;

    private PlcResponse SUT;

    @BeforeEach
    void setUp() {
        responseItems = new ArrayList<>();
        SUT = new PlcResponse<PlcRequest, ResponseItem, RequestItem>(mock(PlcRequest.class), responseItems) {
        };
    }

    @Test
    void getRequest() {
        Assertions.assertNotNull(SUT.getRequest());
    }

    @Test
    void getResponseItems() {
        Assertions.assertEquals(0, SUT.getResponseItems().size());
    }

    @Test
    void getResponseItem() {
        Assertions.assertEquals(Optional.empty(), SUT.getResponseItem());
        responseItems.add(mock(ResponseItem.class));
        Assertions.assertTrue(SUT.getResponseItem().isPresent());
        responseItems.add(mock(ResponseItem.class));
        Assertions.assertThrows(IllegalStateException.class, () -> {
            SUT.getResponseItem();
        });
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

    @Test
    void getValue() {
        Assertions.assertEquals(Optional.empty(), SUT.getValue(null));
    }

}