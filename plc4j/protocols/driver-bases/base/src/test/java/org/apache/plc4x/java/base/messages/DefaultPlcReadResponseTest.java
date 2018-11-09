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

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.messages.items.DefaultByteArrayFieldItem;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class DefaultPlcReadResponseTest {

    @Mock
    private InternalPlcReadRequest request;
    private DefaultPlcReadResponse SUT;

    @BeforeEach
    void setUp() {
        Byte[] data = new Byte[] {(byte) 0x42, (byte) 0x23};
        Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();
        fields.put("foo", new ImmutablePair<>(PlcResponseCode.OK, new DefaultByteArrayFieldItem(data)));
        fields.put("bar", new ImmutablePair<>(PlcResponseCode.NOT_FOUND, new DefaultByteArrayFieldItem(data)));
        SUT = new DefaultPlcReadResponse(request, fields);
    }

    @Test
    void getRequest() {
        assertThat(SUT.getRequest(), equalTo(request));
    }

    @Test
    void getField() {
        when(request.getField("foo")).thenReturn(mock(PlcField.class));
        PlcField field = SUT.getField("foo");
        assertThat(field, notNullValue());
    }

    @Test
    void getResponseCodeForNonexistentField() {
        assertThrows(PlcInvalidFieldException.class, () -> SUT.getResponseCode("hurz"));
    }

    @Test
    void isValidByteArray() {
        boolean valid = SUT.isValidByteArray("foo");
        assertThat(valid, equalTo(true));
    }

    @Test
    void checkInvalidField() {
        assertThrows(PlcInvalidFieldException.class, () -> SUT.isValidByteArray("hurz"));
    }

    @Test
    void checkNonOkResponseCode() {
        assertThrows(PlcRuntimeException.class, () -> SUT.isValidByteArray("bar"));
    }

    @Test
    void isValidByteArrayWithIndex() {
        boolean valid = SUT.isValidByteArray("foo", 0);
        assertThat(valid, equalTo(true));
    }

    @Test
    void getByteArray() {
        Byte[] data = SUT.getByteArray("foo");
        assertThat(data, notNullValue());
    }

    @Test
    void getByteArrayWithIndex() {
        Byte[] data = SUT.getByteArray("foo", 0);
        assertThat(data, notNullValue());
    }

    @Test
    void getAllByteArrays() {
        Collection<Byte[]> byteArrays = SUT.getAllByteArrays("foo");
        assertThat(byteArrays, notNullValue());
    }

}