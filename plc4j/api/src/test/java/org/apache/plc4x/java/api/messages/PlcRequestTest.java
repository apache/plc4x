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
package org.apache.plc4x.java.api.messages;

import java.util.ArrayList;
import java.util.Optional;

import static org.hamcrest.collection.IsEmptyCollection.empty;
import static org.hamcrest.core.Is.is;
import static org.hamcrest.core.IsEqual.equalTo;
import static org.mockito.Mockito.mock;

public class PlcRequestTest {

/*    private List<RequestItem> requestItems;

    private PlcRequest<RequestItem> SUT;

    @Before
    public void setUp() {
        requestItems = new ArrayList<>();
        SUT = new PlcRequest<RequestItem>(requestItems) {
        };
    }

    @Test
    public void addField() {
        SUT.addField(mock(RequestItem.class));
    }

    @Test
    public void getRequestItems() {
        assertThat(SUT.getRequestItems(), empty());
    }

    @Test
    public void getRequestItem() {
        assertThat(SUT.getRequestItem(), equalTo(Optional.empty()));
        requestItems.add(mock(RequestItem.class));
        assertThat(SUT.getRequestItem().isPresent(), is(true));
        requestItems.add(mock(RequestItem.class));
        try {
            SUT.getRequestItem();
            fail("Too many items in PlcRequest should have failed.");
        } catch (IllegalStateException e) {
            // expected
        }
    }

    @Test
    public void setRequestItem() {
        SUT.setRequestItem(mock(RequestItem.class));
        requestItems.add(mock(RequestItem.class));
        try {
            SUT.setRequestItem(mock(RequestItem.class));
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
    public void equalsAndHashCode() {
        PlcRequest other = new PlcRequest<RequestItem>(requestItems) {
        };
        assertThat(SUT.hashCode(), equalTo(other.hashCode()));
        assertThat(SUT.equals(other), equalTo(true));
        assertThat(SUT.equals(new Object()), equalTo(false));
        assertEquals(SUT, SUT);
    }*/
}