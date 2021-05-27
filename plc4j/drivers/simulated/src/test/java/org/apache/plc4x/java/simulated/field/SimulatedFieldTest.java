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

package org.apache.plc4x.java.simulated.field;

import nl.jqno.equalsverifier.EqualsVerifier;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.simulated.types.SimulatedFieldType;
import org.junit.jupiter.api.Test;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.jupiter.api.Assertions.assertThrows;

public class SimulatedFieldTest {

    @Test
    public void constructor() {
        assertThat(SimulatedField.matches("RANDOM/test:Integer[2]"), equalTo(true));
        SimulatedField field = SimulatedField.of("RANDOM/test:Integer[2]");
        assertThat(field.getType(), equalTo(SimulatedFieldType.RANDOM));
        assertThat(field.getName(), equalTo("test"));
        assertThat(field.getPlcDataType(), equalTo("DINT"));
        assertThat(field.getNumberOfElements(), equalTo(2));
        assertThat(field.toString(),
            equalTo("TestField{type=RANDOM, name='test', dataType='DINT', numElements=2}"));
    }

    /*[TODO] Add support for Full Java Type Names back in after plc4go changes have merged
    @Test
    public void constructor() {
        assertThat(SimulatedField.matches("RANDOM/test:Int[2]"), equalTo(true));
        SimulatedField field = SimulatedField.of("RANDOM/test:Int[2]");
        assertThat(field.getType(), equalTo(SimulatedFieldType.RANDOM));
        assertThat(field.getName(), equalTo("test"));
        assertThat(field.getPlcDataType(), equalTo("Integer"));
        assertThat(field.getNumberOfElements(), equalTo(2));
        assertThat(field.toString(),
            equalTo("TestField{type=RANDOM, name='test', dataType='Integer', numElements=2}"));
    }*/

    @Test
    public void invalidType() {
        assertThat(SimulatedField.matches("RANDOM/test:Foo"), equalTo(true));
        assertThrows(PlcInvalidFieldException.class, () -> SimulatedField.of("RANDOM/test:Foo"));
    }

    @Test
    public void invalidAddress() {
        assertThat(SimulatedField.matches("Foo"), equalTo(false));
        assertThrows(PlcInvalidFieldException.class, () -> SimulatedField.of("Foo"));
    }

    @Test
    public void equalsTest() {
        EqualsVerifier.forClass(SimulatedField.class).usingGetClass().verify();
    }

}
