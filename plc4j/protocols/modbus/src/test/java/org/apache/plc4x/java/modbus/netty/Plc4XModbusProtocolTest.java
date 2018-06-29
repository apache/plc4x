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
package org.apache.plc4x.java.modbus.netty;

import org.apache.commons.lang3.reflect.MethodUtils;
import org.junit.Test;

import java.util.Arrays;
import java.util.BitSet;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertSame;

public class Plc4XModbusProtocolTest {

    @Test
    public void coilEncoding() throws Exception {
        byte[] bytes = (byte[]) MethodUtils.invokeMethod(new Plc4XModbusProtocol(), true, "produceCoilValue", Arrays.asList(false, false, true, false, false, false, true, false));
        assertSame(1, bytes.length);
        String s1 = String.format("%8s", Integer.toBinaryString(bytes[0] & 0xFF)).replace(' ', '0');
        assertEquals("00100010", s1);

        // TODO: just for reference we could use a bitset too
        BitSet bitSet = new BitSet(8);
        bitSet.set(1);
        bitSet.set(5);
        byte[] bytes1 = bitSet.toByteArray();
        assertSame(1, bytes1.length);

        String s2 = String.format("%8s", Integer.toBinaryString(bytes1[0] & 0xFF)).replace(' ', '0');
        assertEquals("00100010", s2);
    }
}