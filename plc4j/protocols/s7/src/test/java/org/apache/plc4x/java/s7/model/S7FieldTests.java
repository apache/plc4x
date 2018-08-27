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

package org.apache.plc4x.java.s7.model;

import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

import static org.hamcrest.core.IsEqual.equalTo;
import static org.junit.Assert.assertThat;

public class S7FieldTests {

    @Test
    @Category(FastTests.class)
    public void testS7Field() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        S7Field s7Field = new S7Field(memoryArea, (short) 0x100);

        assertThat(s7Field.getMemoryArea(), equalTo(MemoryArea.DATA_BLOCKS));
        assertThat(s7Field.getByteOffset(), equalTo((short) 0x100));
    }

    @Test
    @Category(FastTests.class)
    public void testS7BitField() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        S7BitField s7Field = new S7BitField(memoryArea, (short) 0x50, (byte) 0x4);

        assertThat(s7Field.getMemoryArea(), equalTo(MemoryArea.DATA_BLOCKS));
        assertThat(s7Field.getByteOffset(), equalTo((short) 0x50));
        assertThat(s7Field.getBitOffset(), equalTo((byte) 0x4));
    }

    @Test
    @Category(FastTests.class)
    public void testS7DatBlockField() {
        S7DataBlockField s7Field = new S7DataBlockField((short) 1, (short) 0x50);

        assertThat(s7Field.getMemoryArea(), equalTo(MemoryArea.DATA_BLOCKS));
        assertThat(s7Field.getDataBlockNumber(), equalTo((short) 1));
        assertThat(s7Field.getByteOffset(), equalTo((short) 0x50));
    }

}