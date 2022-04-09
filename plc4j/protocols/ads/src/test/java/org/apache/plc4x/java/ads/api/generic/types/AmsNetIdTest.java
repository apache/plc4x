/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ads.api.generic.types;

import org.junit.Test;

import static org.hamcrest.Matchers.is;
import static org.junit.Assert.assertThat;

public class AmsNetIdTest {

    @Test
    public void netIdBytes() {
        // note bytes in reverse order
        AmsNetId netid = AmsNetId.of((byte)0x01, (byte)0x02, (byte)0x03, (byte)0x04, (byte)0x05, (byte)0x06);
        assertThat(netid.toString(), is("1.2.3.4.5.6"));
    }

    @Test
    public void netIdString() {
        // note bytes in reverse order
        AmsNetId netid = AmsNetId.of("1.2.3.4.5.6");
        assertThat(netid.toString(), is("1.2.3.4.5.6"));
    }
    
    @Test(expected = IllegalArgumentException.class)
    public void netIdTooShort() {
        // note bytes in reverse order
        AmsNetId netid = AmsNetId.of("1.2.3.4");
    }

    @Test(expected = IllegalArgumentException.class)
    public void netIdStringTooLong() {
        // note bytes in reverse order
        AmsNetId netid = AmsNetId.of("1.2.3.4.5.6.7.8");
    }

    @Test(expected = IllegalArgumentException.class)
    public void netIdStringWrongSeperator() {
        // note bytes in reverse order
        AmsNetId netid = AmsNetId.of("1:2:3:4:5:6");
    }

}