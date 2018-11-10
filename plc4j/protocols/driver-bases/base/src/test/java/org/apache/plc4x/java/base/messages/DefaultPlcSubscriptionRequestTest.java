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
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class DefaultPlcSubscriptionRequestTest {

    private PlcSubscriber subscriber;
    private PlcField fooField;
    private SubscriptionPlcField fooSubscriptionField;
    private DefaultPlcSubscriptionRequest SUT;

    @BeforeEach
    void setUp() {
        subscriber = mock(PlcSubscriber.class);
        fooField = mock(PlcField.class);
        fooSubscriptionField = mock(SubscriptionPlcField.class);
        when(fooSubscriptionField.getPlcField()).thenReturn(fooField);
        SUT = new DefaultPlcSubscriptionRequest(
            subscriber, new LinkedHashMap<>(Collections.singletonMap("foo", fooSubscriptionField)));
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
        PlcField bar = SUT.getField("bar");
        assertThat(bar, nullValue());
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
    void getSubscriptionFields() {
        List<SubscriptionPlcField> fields = SUT.getSubscriptionFields();
        assertThat(fields, notNullValue());
        assertThat(fields.size(), equalTo(1));
        SubscriptionPlcField field = fields.iterator().next();
        assertThat(field, equalTo(fooSubscriptionField));
    }

    @Test
    void getSubscriptionPlcFieldMap() {
        LinkedHashMap<String, SubscriptionPlcField> map = SUT.getSubscriptionPlcFieldMap();
        assertThat(map, notNullValue());
        assertThat(map.size(), equalTo(1));
        assertThat(map.get("foo"), equalTo(fooSubscriptionField));
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
    void getNamedSubscriptionFields() {
        List<Pair<String, SubscriptionPlcField>> namedFields = SUT.getNamedSubscriptionFields();
        assertThat(namedFields, notNullValue());
        assertThat(namedFields.size(), equalTo(1));
        Pair<String, SubscriptionPlcField> entry = namedFields.iterator().next();
        assertThat(entry.getKey(), equalTo("foo"));
        assertThat(entry.getValue(), equalTo(fooSubscriptionField));
    }

    @Test
    void getSubscriber() {
        assertThat(SUT.getSubscriber(), equalTo(subscriber));
    }

    @Test
    void builder() {
        PlcFieldHandler fieldHandler = mock(PlcFieldHandler.class);
        when(fieldHandler.createField(anyString())).thenReturn(mock(PlcField.class));
        DefaultPlcSubscriptionRequest.Builder builder = new DefaultPlcSubscriptionRequest.Builder(subscriber, fieldHandler);
        builder.addChangeOfStateField("state", "bar");
        builder.addCyclicField("cyclic", "bar", Duration.of(3, ChronoUnit.SECONDS));
        builder.addEventField("event", "bar");
        assertThrows(PlcRuntimeException.class, () -> builder.addEventField("event", "bar"));
        PlcSubscriptionRequest subscriptionRequest = builder.build();
        assertThat(subscriptionRequest, notNullValue());
        assertThat(subscriptionRequest.getNumberOfFields(), equalTo(3));
        assertThat(subscriptionRequest.getField("state"), notNullValue());
        assertThat(subscriptionRequest.getField("cyclic"), notNullValue());
        assertThat(subscriptionRequest.getField("event"), notNullValue());
    }

}