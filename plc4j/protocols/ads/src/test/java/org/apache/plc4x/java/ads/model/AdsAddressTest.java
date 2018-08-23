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

import org.apache.plc4x.java.api.exceptions.PlcInvalidAddressException;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.assertThat;

public class AdsAddressTest {

    @Test
    public void of() throws Exception {
        AdsAddress address = AdsAddress.of("1/10");
        assertThat(address.getIndexGroup(), is(1L));
        assertThat(address.getIndexOffset(), is(10L));
    }

    @Test
    public void ofHex() throws Exception {
        AdsAddress address = AdsAddress.of("0x1/0xff");
        assertThat(address.getIndexGroup(), is(1L));
        assertThat(address.getIndexOffset(), is(255L));
    }

    @Test(expected = PlcInvalidAddressException.class)
    public void stringInAddress() throws Exception {
        AdsAddress address = AdsAddress.of("group/offset");
    }

    @Test(expected = PlcInvalidAddressException.class)
    public void singleNumberAddress() throws Exception {
        AdsAddress address = AdsAddress.of("10");
    }

    @Test(expected = PlcInvalidAddressException.class)
    public void wrongSeperator() throws Exception {
        AdsAddress address = AdsAddress.of("1:10");
    }

    @Test
    public void getGroupAndOffset() {
        AdsAddress address = AdsAddress.of(2L, 20L);
        assertThat(address.getIndexGroup(), is(2L));
        assertThat(address.getIndexOffset(), is(20L));
    }
}