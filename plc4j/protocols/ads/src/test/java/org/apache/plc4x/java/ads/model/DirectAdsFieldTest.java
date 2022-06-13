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
package org.apache.plc4x.java.ads.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.junit.Test;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.assertThat;

public class DirectAdsFieldTest {

    @Test
    public void of() {
        DirectAdsField field = DirectAdsField.of("1/10:BYTE[2]");
        assertThat(field.getIndexGroup(), is(1L));
        assertThat(field.getIndexOffset(), is(10L));
        assertThat(field.getAdsDataType(), is(AdsDataType.BYTE));
        assertThat(field.getNumberOfElements(), is(2));
    }

    @Test
    public void ofHex() {
        DirectAdsField field = DirectAdsField.of("0x1/0xff:BYTE[2]");
        assertThat(field.getIndexGroup(), is(1L));
        assertThat(field.getIndexOffset(), is(255L));
        assertThat(field.getAdsDataType(), is(AdsDataType.BYTE));
        assertThat(field.getNumberOfElements(), is(2));
    }

    @Test(expected = PlcInvalidFieldException.class)
    public void stringInField() {
        DirectAdsField field = DirectAdsField.of("group/offset");
    }

    @Test(expected = PlcInvalidFieldException.class)
    public void singleNumberField() {
        DirectAdsField field = DirectAdsField.of("10");
    }

    @Test(expected = PlcInvalidFieldException.class)
    public void wrongSeperator() {
        DirectAdsField field = DirectAdsField.of("1:10");
    }

    @Test
    public void getGroupAndOffset() {
        DirectAdsField field = DirectAdsField.of(2L, 20L, AdsDataType.BYTE, 1);
        assertThat(field.getIndexGroup(), is(2L));
        assertThat(field.getIndexOffset(), is(20L));
        assertThat(field.getAdsDataType(), is(AdsDataType.BYTE));
        assertThat(field.getNumberOfElements(), is(1));
    }

}