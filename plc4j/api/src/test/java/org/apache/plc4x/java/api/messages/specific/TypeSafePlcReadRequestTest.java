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
package org.apache.plc4x.java.api.messages.specific;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.model.Address;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Collections;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class TypeSafePlcReadRequestTest {

    ReadRequestItem<String> readRequestItemString;

    @BeforeEach
    void setUp() {
        readRequestItemString = new ReadRequestItem<>(String.class, null);
    }

    @Test
    void constuctor() {
        new TypeSafePlcReadRequest<>(String.class);
        new TypeSafePlcReadRequest<>(String.class, mock(PlcReadRequest.class));
        PlcReadRequest request = mock(PlcReadRequest.class);
        when(request.getRequestItems()).thenReturn(Collections.singletonList(readRequestItemString));
        new TypeSafePlcReadRequest<>(String.class, request);
        new TypeSafePlcReadRequest<>(String.class, mock(Address.class));
        new TypeSafePlcReadRequest<>(String.class, mock(Address.class), 3);
        new TypeSafePlcReadRequest<>(String.class, readRequestItemString);
        Assertions.assertThrows(IllegalArgumentException.class, () -> {
            new TypeSafePlcReadRequest<>(Byte.class, request);
        });
    }

    @Test
    void addItem() {
        new TypeSafePlcReadRequest<>(String.class).addItem(readRequestItemString);
    }

    @Test
    void getCheckedReadRequestItems() {
        new TypeSafePlcReadRequest<>(String.class).getCheckedReadRequestItems();
    }

    @Test
    void getRequestItem() {
        new TypeSafePlcReadRequest<>(String.class).getRequestItem();
    }

    @Test
    void getDataType() {
        Assertions.assertEquals(String.class, new TypeSafePlcReadRequest<>(String.class).getDataType());
    }

}