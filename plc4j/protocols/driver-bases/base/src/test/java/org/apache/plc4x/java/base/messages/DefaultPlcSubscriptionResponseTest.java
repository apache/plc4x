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
import org.apache.plc4x.java.api.exceptions.PlcNotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.hamcrest.core.IsNull.nullValue;
import static org.junit.Assert.assertThat;
import static org.junit.Assert.fail;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;

class DefaultPlcSubscriptionResponseTest {

    private InternalPlcSubscriptionRequest mockRequest;
    private PlcSubscriptionHandle mockSubscriptionHandle;
    private Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> mockValues;
    private DefaultPlcSubscriptionResponse SUT;

    @BeforeEach
    @SuppressWarnings("unchecked")
    void setUp() {
        mockRequest = mock(InternalPlcSubscriptionRequest.class);
        mockSubscriptionHandle = mock(PlcSubscriptionHandle.class);
        mockValues = new HashMap<>();
        mockValues.put("foo", new ImmutablePair(PlcResponseCode.OK, mockSubscriptionHandle));
        mockValues.put("bar", new ImmutablePair(PlcResponseCode.NOT_FOUND, mock(PlcSubscriptionHandle.class)));
        SUT = new DefaultPlcSubscriptionResponse(mockRequest, mockValues);
    }

    @Test
    void getSubscriptionHandle() {
        PlcSubscriptionHandle foo = SUT.getSubscriptionHandle("foo");
        assertThat(foo, equalTo(mockSubscriptionHandle));
        assertThrows(PlcRuntimeException.class, () -> SUT.getSubscriptionHandle("bar"));
        assertThat(SUT.getSubscriptionHandle("hurz"), nullValue());
    }

    @Test
    void getFieldNames() {
        Collection<String> fieldNames = SUT.getFieldNames();
        assertThat(fieldNames, notNullValue());
        assertThat(fieldNames.size(), equalTo(2));
        for (String fieldName : fieldNames) {
            if(!"foo".equals(fieldName) && !"bar".equals(fieldName)) {
                fail();
            }
        }
    }

    @Test
    void getField() {
        assertThrows(PlcNotImplementedException.class, () -> SUT.getField("foo"));
    }

    @Test
    void getResponseCode() {
        PlcResponseCode responseCode = SUT.getResponseCode("foo");
        assertThat(responseCode, notNullValue());
        assertThat(responseCode, equalTo(PlcResponseCode.OK));

        responseCode = SUT.getResponseCode("hurz");
        assertThat(responseCode, nullValue());
    }

    @Test
    void getRequest() {
        PlcSubscriptionRequest request = SUT.getRequest();
        assertThat(request, equalTo(mockRequest));
    }

    @Test
    void getSubscriptionHandles() {
        Collection<PlcSubscriptionHandle> subscriptionHandles = SUT.getSubscriptionHandles();
        assertThat(subscriptionHandles, notNullValue());
        assertThat(subscriptionHandles.size(), equalTo(2));
        boolean found = false;
        for (PlcSubscriptionHandle subscriptionHandle : subscriptionHandles) {
            if(subscriptionHandle.equals(mockSubscriptionHandle)) {
                found = true;
            }
        }
        assertThat(found, equalTo(true));
    }

    @Test
    void getValues() {
        Map<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> values = SUT.getValues();
        assertThat(values, notNullValue());
        assertThat(values, equalTo(mockValues));
    }

}