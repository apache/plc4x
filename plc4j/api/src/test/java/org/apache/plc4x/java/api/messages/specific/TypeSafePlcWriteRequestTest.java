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
package org.apache.plc4x.java.api.messages.specific;

import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.Collections;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class TypeSafePlcWriteRequestTest {

/*    PlcWriteRequestItem<String> writeRequestItemString;

    @Before
    public void setUp() {
        writeRequestItemString = new PlcWriteRequestItem<>(String.class, mock(PlcField.class));
    }

    @Test(expected = IllegalArgumentException.class)
    public void constuctor() {
        new TypeSafePlcWriteRequest<>(String.class);
        new TypeSafePlcWriteRequest<>(String.class, mock(PlcWriteRequest.class));
        PlcWriteRequest request = mock(PlcWriteRequest.class);
        when(request.getRequestItems()).thenReturn(Collections.singletonList(writeRequestItemString));
        new TypeSafePlcWriteRequest<>(String.class, request);
        new TypeSafePlcWriteRequest<>(String.class, mock(PlcField.class));
        new TypeSafePlcWriteRequest<>(String.class, mock(PlcField.class), "");
        new TypeSafePlcWriteRequest<>(String.class, writeRequestItemString);

        // expects an exception
        new TypeSafePlcWriteRequest<>(Byte.class, request);
    }

    @Test
    public void addField() {
        new TypeSafePlcWriteRequest<>(String.class).addField(writeRequestItemString);
    }

    @Test
    public void getCheckedWriteRequestItems() {
        new TypeSafePlcWriteRequest<>(String.class).getCheckedRequestItems();
    }

    @Test
    public void getRequestItem() {
        new TypeSafePlcWriteRequest<>(String.class).getRequestItem();
    }

    @Test
    public void getDataType() {
        assertThat(new TypeSafePlcWriteRequest<>(String.class).getDataType(), equalTo(String.class));
    }*/

}