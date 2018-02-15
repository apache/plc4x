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

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.java.api.messages.items.RequestItem;
import org.junit.Before;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.Mockito.mock;

public class PlcRequestTest {

    private List<RequestItem> requestItems;

    private PlcRequest SUT;

    @Before
    public void setUp() {
        requestItems = new ArrayList<>();
        SUT = new PlcRequest<RequestItem>(requestItems) {
        };
    }

    @Test
    public void addItem() {
        SUT.addItem(mock(RequestItem.class));
    }

    @Test
    public void getRequestItems() {
        assertThat(SUT.getRequestItems()).isEmpty();
    }

    @Test
    public void getRequestItem() {
        assertThat(SUT.getRequestItem()).isEqualTo(Optional.empty());
        requestItems.add(mock(RequestItem.class));
        assertThat(SUT.getRequestItem().isPresent()).isTrue();
        requestItems.add(mock(RequestItem.class));
        assertThatThrownBy(() -> SUT.getRequestItem()).isInstanceOf(IllegalStateException.class);
    }

    @Test
    public void setRequestItem() {
        SUT.setRequestItem(mock(RequestItem.class));
        requestItems.add(mock(RequestItem.class));
        assertThatThrownBy(() -> SUT.setRequestItem(mock(RequestItem.class))).isInstanceOf(IllegalStateException.class);
    }

    @Test
    public void getNumberOfItems() {
        assertThat(SUT.getNumberOfItems()).isEqualTo(0);
    }

    @Test
    public void isMultiValue() {
        assertThat(SUT.isMultiValue()).isFalse();
    }

    @Test
    public void isEmpty() {
        assertThat(SUT.isEmpty()).isTrue();
    }

}