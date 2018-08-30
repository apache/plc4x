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

import org.junit.Before;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import static org.hamcrest.collection.IsEmptyCollection.empty;
import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.core.IsNull.notNullValue;
import static org.junit.Assert.*;
import static org.mockito.Mockito.mock;

public class PlcResponseTest {

/*    private PlcResponse<PlcRequest, ResponseItem, RequestItem> SUT;

    private List<ResponseItem> responseItems;

    private PlcRequest plcRequest;

    @Before
    public void setUp() {
        responseItems = new ArrayList<>();
        plcRequest = mock(PlcRequest.class);
        SUT = new PlcResponse<PlcRequest, ResponseItem, RequestItem>(plcRequest, responseItems) {
        };
    }

    @Test
    public void getRequest() {
        assertThat(SUT.getRequest(), notNullValue());
    }

    @Test
    public void getResponseItems() {
        assertThat(SUT.getResponseItems(), empty());
    }

    @Test
    public void getResponseItem() {
        assertThat(SUT.getResponseItem(), equalTo(Optional.empty()));
        responseItems.add(mock(ResponseItem.class));
        assertThat(SUT.getResponseItem().isPresent(), is(true));
        responseItems.add(mock(ResponseItem.class));
        try {
            SUT.getResponseItem();
            fail("PlcResponse.getResponseItem() should fail if contains multiple items.");
        } catch (IllegalStateException e) {
            // expected
        }
    }

    @Test
    public void getNumberOfItems() {
        assertThat(SUT.getNumberOfItems(), equalTo(0));
    }

    @Test
    public void isMultiValue() {
        assertThat(SUT.isMultiValue(), is(false));
    }

    @Test
    public void isEmpty() {
        assertThat(SUT.isEmpty(), is(true));
    }

    @Test
    public void getValue() {
        assertThat(SUT.getValue(null), equalTo(Optional.empty()));
    }

    @Test
    public void equalsAndHashCode() {
        PlcResponse other = new PlcResponse<PlcRequest, ResponseItem, RequestItem>(plcRequest, responseItems) {
        };
        assertThat(SUT.hashCode(), equalTo(other.hashCode()));
        assertThat(SUT.equals(other), equalTo(true));
        assertThat(SUT.equals(new Object()), equalTo(false));
        assertEquals(SUT, SUT);
    }*/

}