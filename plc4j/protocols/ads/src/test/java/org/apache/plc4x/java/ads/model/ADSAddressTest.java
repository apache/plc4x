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
package org.apache.plc4x.java.ads.model;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.*;

public class ADSAddressTest {

    @Test
    public void of() {
        ADSAddress address = ADSAddress.of("1/10");
        assertThat(address.getIndexGroup(), is(1L));
        assertThat(address.getIndexOffset(), is(10L));
    }

    @Test(expected  = IllegalArgumentException.class)
    public void stringInAddress() {
        ADSAddress address = ADSAddress.of("group/offset");
    }

    @Test(expected  = IllegalArgumentException.class)
    public void singleNumberAddress() {
        ADSAddress address = ADSAddress.of("10");
    }

    @Test(expected  = IllegalArgumentException.class)
    public void wrongSeperator() {
        ADSAddress address = ADSAddress.of("1:10");
    }

    @Test
    public void getGroupAndOffset() {
        ADSAddress address = ADSAddress.of(2L, 20L);
        assertThat(address.getIndexGroup(), is(2L));
        assertThat(address.getIndexOffset(), is(20L));
    }
}