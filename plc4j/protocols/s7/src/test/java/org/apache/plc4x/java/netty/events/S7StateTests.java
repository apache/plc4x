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

package org.apache.plc4x.java.netty.events;

import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class S7StateTests {

    @Test
    @Category(FastTests.class)
    public void testInitialS7ConnectionEvent() {
        S7ConnectionEvent s7event = new S7ConnectionEvent();

        assertThat(s7event.getState(), equalTo(S7ConnectionState.INITIAL));
    }

    @Test
    @Category(FastTests.class)
    public void testS7ConnectionEvent() {
        S7ConnectionEvent s7event = new S7ConnectionEvent(S7ConnectionState.SETUP_COMPLETE);

        assertThat(s7event.getState(), equalTo(S7ConnectionState.SETUP_COMPLETE) );
    }

}