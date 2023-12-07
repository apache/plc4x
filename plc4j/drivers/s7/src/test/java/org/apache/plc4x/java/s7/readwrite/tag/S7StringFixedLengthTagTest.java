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

package org.apache.plc4x.java.s7.readwrite.tag;

import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

public class S7StringFixedLengthTagTest {
    /**
     * PLC4X-312 - CAN NOT READ STRING FROM S7 PLC
     */
    @Test
    public void testStringAddress() {
        final S7StringFixedLengthTag s7StringFixedLengthTag = S7StringFixedLengthTag.of("%DB145.DBX38:STRING(8)[1]");
        Assertions.assertEquals(TransportSize.STRING, s7StringFixedLengthTag.getDataType());
        Assertions.assertEquals(1, s7StringFixedLengthTag.getNumberOfElements());
        Assertions.assertEquals(145, s7StringFixedLengthTag.getBlockNumber());
        Assertions.assertEquals(MemoryArea.DATA_BLOCKS, s7StringFixedLengthTag.getMemoryArea());
        Assertions.assertEquals(38, s7StringFixedLengthTag.getByteOffset());
        Assertions.assertEquals(0, s7StringFixedLengthTag.getBitOffset());
        Assertions.assertEquals("S7StringFixedLengthTag", s7StringFixedLengthTag.getClass().getSimpleName());
        Assertions.assertEquals(8, s7StringFixedLengthTag.getStringLength());
    }

}
