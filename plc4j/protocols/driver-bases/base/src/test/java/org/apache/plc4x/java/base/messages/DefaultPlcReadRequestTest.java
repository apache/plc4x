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

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class DefaultPlcReadRequestTest {

    private PlcReader reader;
    private PlcField fooField;
    private DefaultPlcReadRequest SUT;

    @BeforeEach
    void setUp() {
        reader = mock(PlcReader.class);
        fooField = mock(PlcField.class);
        SUT = new DefaultPlcReadRequest(reader, new LinkedHashMap<>(Collections.singletonMap("foo", fooField)));
    }

    @Test
    void execute() {
        SUT.execute();
    }

    @Test
    void getNumberOfFields() {
        assertThat(SUT.getNumberOfFields(), equalTo(1));
    }

    @Test
    void getFieldNames() {
        LinkedHashSet<String> fieldNames = SUT.getFieldNames();
        assertThat(fieldNames.size(), equalTo(1));
        assertThat(fieldNames.iterator().next(), equalTo("foo"));
    }

    @Test
    void getField() {
        PlcField foo = SUT.getField("foo");
        assertThat(foo, notNullValue());
    }

    @Test
    void getFields() {
        List<PlcField> fields = SUT.getFields();
        assertThat(fields, notNullValue());
        assertThat(fields.size(), equalTo(1));
        PlcField field = fields.iterator().next();
        assertThat(field, equalTo(fooField));
    }

    @Test
    void getNamedFields() {
        List<Pair<String, PlcField>> namedFields = SUT.getNamedFields();
        assertThat(namedFields, notNullValue());
        assertThat(namedFields.size(), equalTo(1));
        Pair<String, PlcField> entry = namedFields.iterator().next();
        assertThat(entry.getKey(), equalTo("foo"));
        assertThat(entry.getValue(), equalTo(fooField));
    }

    @Test
    void getReader() {
        assertThat(SUT.getReader(), equalTo(reader));
    }

    @Test
    void builder() {
        PlcFieldHandler fieldHandler = mock(PlcFieldHandler.class);
        when(fieldHandler.createField(anyString())).thenReturn(mock(PlcField.class));
        DefaultPlcReadRequest.Builder builder = new DefaultPlcReadRequest.Builder(reader, fieldHandler);
        builder.addItem("foo", "bar");
        assertThrows(PlcRuntimeException.class, () -> builder.addItem("foo", "bar"));
        PlcReadRequest readRequest = builder.build();
        assertThat(readRequest, notNullValue());
        assertThat(readRequest.getNumberOfFields(), equalTo(1));
        assertThat(readRequest.getField("foo"), notNullValue());
    }

}