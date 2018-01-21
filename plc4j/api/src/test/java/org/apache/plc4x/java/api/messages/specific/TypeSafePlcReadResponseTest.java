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

import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.items.ReadRequestItem;
import org.apache.plc4x.java.api.messages.items.ReadResponseItem;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class TypeSafePlcReadResponseTest {

    ReadResponseItem<String> readResponseItemString;

    @BeforeEach
    void setUp() {
        readResponseItemString = new ReadResponseItem<>(mock(ReadRequestItem.class), null, Arrays.asList("", ""));
    }

    @Test
    void constuctor() {
        new TypeSafePlcReadResponse<>(mock(PlcReadResponse.class));
        PlcReadResponse response = mock(PlcReadResponse.class);
        when(response.getResponseItems()).thenReturn((List) Collections.singletonList(readResponseItemString));
        new TypeSafePlcReadResponse<>(response);
        TypeSafePlcReadRequest mock = mock(TypeSafePlcReadRequest.class);
        when(mock.getDataType()).thenReturn(String.class);
        new TypeSafePlcReadResponse<>(mock, readResponseItemString);
        new TypeSafePlcReadResponse<>(mock, Collections.singletonList(readResponseItemString));
        Assertions.assertThrows(IllegalArgumentException.class, () ->{
            when(mock.getDataType()).thenReturn(Byte.class);
            new TypeSafePlcReadResponse<>(mock, readResponseItemString);
        });
    }

    @Test
    void getRequest() {
        new TypeSafePlcReadResponse<>(mock(PlcReadResponse.class)).getRequest();
    }

    @Test
    void getResponseItems() {
        new TypeSafePlcReadResponse<>(mock(PlcReadResponse.class)).getResponseItems();
    }

    @Test
    void getResponseItem() {
        new TypeSafePlcReadResponse<>(mock(PlcReadResponse.class)).getResponseItem();
    }

}