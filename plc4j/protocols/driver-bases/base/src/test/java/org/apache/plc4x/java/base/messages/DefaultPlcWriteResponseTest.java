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

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collection;
import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.Map;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class DefaultPlcWriteResponseTest {

    @Mock
    private InternalPlcWriteRequest request;
    private DefaultPlcWriteResponse SUT;

    @BeforeEach
    void setUp() {
        SUT = new DefaultPlcWriteResponse(request, Collections.singletonMap("foo", PlcResponseCode.OK));
    }

    @Test
    void getValues() {
        Map<String, PlcResponseCode> values = SUT.getValues();
        assertThat(values, notNullValue());
        assertThat(values.size(), equalTo(1));
        assertThat(values.keySet().iterator().next(), equalTo("foo"));
        assertThat(values.values().iterator().next(), equalTo(PlcResponseCode.OK));
    }

    @Test
    void getRequest() {
        InternalPlcWriteRequest actRequest = SUT.getRequest();
        assertThat(actRequest, equalTo(request));
    }

    @Test
    void getFieldNames() {
        when(request.getFieldNames()).thenReturn(new LinkedHashSet<>(Collections.singletonList("foo")));
        Collection<String> fieldNames = SUT.getFieldNames();
        assertThat(fieldNames, notNullValue());
        assertThat(fieldNames.size(), equalTo(1));
        assertThat(fieldNames.iterator().next(), equalTo("foo"));
    }

    @Test
    void getField() {
        when(request.getField("foo")).thenReturn(mock(PlcField.class));
        PlcField field = SUT.getField("foo");
        assertThat(field, notNullValue());
    }

    @Test
    void getResponseCode() {
        PlcResponseCode responseCode = SUT.getResponseCode("foo");
        assertThat(responseCode, notNullValue());
        assertThat(responseCode, equalTo(PlcResponseCode.OK));
    }

}