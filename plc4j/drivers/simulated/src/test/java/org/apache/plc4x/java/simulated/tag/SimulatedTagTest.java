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
package org.apache.plc4x.java.simulated.tag;

import nl.jqno.equalsverifier.EqualsVerifier;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.apache.plc4x.java.simulated.types.SimulatedTagType;
import org.junit.jupiter.api.Test;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;

public class SimulatedTagTest {

    @Test
    public void constructor() {
        assertThat(SimulatedTag.matches("RANDOM/test:DINT[2]"), equalTo(true));
        SimulatedTag tag = SimulatedTag.of("RANDOM/test:DINT[2]");
        assertThat(tag.getType(), equalTo(SimulatedTagType.RANDOM));
        assertThat(tag.getName(), equalTo("test"));
        assertThat(tag.getPlcValueType().name(), equalTo("DINT"));
        assertThat(tag.getArrayInfo().get(0).GetSize(), equalTo(2));
        assertThat(tag.toString(),
            equalTo("SimulatedTag{type=RANDOM, name='test', dataType='DINT', numElements=2}"));
    }

    /*[TODO] Add support for Full Java Type Names back in after plc4go changes have merged
    @Test
    public void constructor() {
        assertThat(SimulatedTag.matches("RANDOM/test:Int[2]"), equalTo(true));
        SimulatedTag tag = SimulatedTag.of("RANDOM/test:Int[2]");
        assertThat(tag.getType(), equalTo(SimulatedTagType.RANDOM));
        assertThat(tag.getName(), equalTo("test"));
        assertThat(tag.getPlcDataType(), equalTo("Integer"));
        assertThat(tag.getNumberOfElements(), equalTo(2));
        assertThat(tag.toString(),
            equalTo("TestTag{type=RANDOM, name='test', dataType='Integer', numElements=2}"));
    }*/

    @Test
    public void invalidType() {
        assertThat(SimulatedTag.matches("RANDOM/test:Foo"), equalTo(true));
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("RANDOM/test:Foo"));
    }

    @Test
    public void invalidAddress() {
        assertThat(SimulatedTag.matches("Foo"), equalTo(false));
        assertThrows(PlcInvalidTagException.class, () -> SimulatedTag.of("Foo"));
    }

    @Test
    public void equalsTest() {
        EqualsVerifier.forClass(SimulatedTag.class).usingGetClass().verify();
    }

}
