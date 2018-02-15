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

import static org.assertj.core.api.Assertions.assertThat;

import org.apache.plc4x.java.s7.netty.model.types.MemoryArea;
import org.apache.plc4x.test.FastTests;
import org.junit.Test;
import org.junit.experimental.categories.Category;

public class S7AddressTests {

    @Test
    @Category(FastTests.class)
    public void testS7Address() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        S7Address s7Address = new S7Address(memoryArea, (short) 0x100);

        assertThat(s7Address.getMemoryArea()).isEqualTo(MemoryArea.DATA_BLOCKS);
        assertThat(s7Address.getByteOffset()).isEqualTo((short) 0x100);
    }

    @Test
    @Category(FastTests.class)
    public void testS7BitAddress() {
        MemoryArea memoryArea = MemoryArea.DATA_BLOCKS;
        S7BitAddress s7Address = new S7BitAddress(memoryArea, (short) 0x50, (byte) 0x4);

        assertThat(s7Address.getMemoryArea()).isEqualTo(MemoryArea.DATA_BLOCKS);
        assertThat(s7Address.getByteOffset()).isEqualTo((short) 0x50);
        assertThat(s7Address.getBitOffset()).isEqualTo((byte) 0x4);
    }

    @Test
    @Category(FastTests.class)
    public void testS7DatBlockAddress() {
        S7DataBlockAddress s7Address = new S7DataBlockAddress((short) 1, (short) 0x50);

        assertThat(s7Address.getMemoryArea()).isEqualTo(MemoryArea.DATA_BLOCKS);
        assertThat(s7Address.getDataBlockNumber()).isEqualTo((short) 1);
        assertThat(s7Address.getByteOffset()).isEqualTo((short) 0x50);
    }

}