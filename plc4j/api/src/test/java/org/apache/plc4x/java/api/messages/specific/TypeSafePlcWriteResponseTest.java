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

import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.messages.items.PlcWriteRequestItem;
import org.apache.plc4x.java.api.messages.items.PlcWriteResponseItem;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.Before;
import org.junit.Test;

import java.util.Collections;
import java.util.List;

import static org.mockito.Mockito.*;

public class TypeSafePlcWriteResponseTest {

    PlcWriteResponseItem<String> writeResponseItemString;

    @Before
    public void setUp() {
        writeResponseItemString = new PlcWriteResponseItem<>(mock(PlcWriteRequestItem.class), PlcResponseCode.OK);
    }

    @Test
    public void constuctor() {
        new TypeSafePlcWriteResponse<>(mock(TypeSafePlcWriteRequest.class), writeResponseItemString);
        new TypeSafePlcWriteResponse<>(mock(TypeSafePlcWriteRequest.class), Collections.singletonList(writeResponseItemString));
    }

    @Test
    public void of() {
        {
            TypeSafePlcWriteResponse.of(mock(PlcWriteResponse.class, RETURNS_DEEP_STUBS));
        }
        {
            PlcWriteResponse response = mock(PlcWriteResponse.class, RETURNS_DEEP_STUBS);
            when(response.getRequest()).thenReturn(mock(TypeSafePlcWriteRequest.class, RETURNS_DEEP_STUBS));
            TypeSafePlcWriteResponse.of(response);
        }
        {
            PlcWriteResponse response = mock(PlcWriteResponse.class, RETURNS_DEEP_STUBS);
            when(response.getResponseItems()).thenReturn((List) Collections.singletonList(writeResponseItemString));
            TypeSafePlcWriteResponse.of(response);
        }
    }

    @Test
    public void getRequest() {
        new TypeSafePlcWriteResponse<>(mock(TypeSafePlcWriteRequest.class), Collections.emptyList()).getRequest();
    }

    @Test
    public void getResponseItems() {
        new TypeSafePlcWriteResponse<>(mock(TypeSafePlcWriteRequest.class), Collections.emptyList()).getResponseItems();
    }

    @Test
    public void getResponseItem() {
        new TypeSafePlcWriteResponse<>(mock(TypeSafePlcWriteRequest.class), Collections.emptyList()).getResponseItem();
    }

}