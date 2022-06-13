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
package org.apache.plc4x.java.ads.api.generic.types;

import org.junit.Test;

import static org.hamcrest.Matchers.*;
import static org.junit.Assert.assertThat;

public class StateTest {

    // Not the best unit tests but here for coverage

    @Test
    public void stateBitFields() {
        int bitMask = State.StateMask.RESPONSE.getMask() | State.StateMask.NO_RETURN.getMask() | State.StateMask.ADS_COMMAND.getMask()
                        | State.StateMask.SYSTEM_COMMAND.getMask() | State.StateMask.HIGH_PRIORITY_COMMAND.getMask() | State.StateMask.TIMESTAMP_ADDED.getMask()
                        | State.StateMask.UDP_COMMAND.getMask() | State.StateMask.INIT_COMMAND.getMask() | State.StateMask.BROADCAST.getMask();
        State state = State.of(bitMask);

        assertThat(state.toString(), not((isEmptyString())));
    }

    @Test
    public void equals() {
        State a = State.of((byte) 0x1, (byte) 0x2);
        State b = State.of((byte) 0x1, (byte) 0x2);
        State c = State.of((byte) 0x1, (byte) 0x4);
        byte array[] = {(byte) 0x1, (byte) 0x2};

        assertThat(a.equals(a), is(true));
        assertThat(a.equals(b), is(true));
        assertThat(a.equals(c), is(false));
        assertThat(a.equals(1), is(false));
        assertThat(a.equals((byte) 1), is(false));
        assertThat(a.equals(array), is(false));
    }

}