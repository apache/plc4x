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
package org.apache.plc4x.java.modbus.messages.items;

import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.junit.Test;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;

@SuppressWarnings("unchecked")
public class DefaultModbusByteArrayFieldItemTest {
    @Test
    public void convertByteArrayToIntegerTest() {
        BaseDefaultFieldItem fieldItem = getFieldItemForIntegerArray();

        Integer itemInteger = fieldItem.getInteger(1);
        assertEquals(456,itemInteger,0);
    }

    @Test
    public void convertByteArrayToIntegerTestReturnsNull() {
        BaseDefaultFieldItem fieldItem = getFieldItemForIntegerArray();

        Integer itemInteger = fieldItem.getInteger(17);
        assertNull(itemInteger);
    }

    private static BaseDefaultFieldItem getFieldItemForIntegerArray() {
        int sizeIntByteBuffer = 12;
        ByteBuffer byteBuffer = ByteBuffer.allocate(sizeIntByteBuffer);
        byteBuffer.order(ByteOrder.BIG_ENDIAN);
        byteBuffer.putInt(0,123);
        byteBuffer.putInt(4,456);
        byteBuffer.putInt(8,789);
        Byte[][] byteArray = new Byte[sizeIntByteBuffer/2][2];
        int cntByteBuffer=0;
        for (byte b : byteBuffer.array()) {
            int shortIndex=cntByteBuffer/2;
            byteArray[shortIndex][cntByteBuffer%2] = b;
            cntByteBuffer++;
        }

        return new DefaultModbusByteArrayFieldItem(byteArray);
    }


    @Test
    public void convertByteArrayToShortTest() {
        BaseDefaultFieldItem fieldItem = getFieldItemForShortArray();

        Short itemShort = fieldItem.getShort(3);
        assertEquals(1011,itemShort,0);
    }

    @Test
    public void convertByteArrayToShortTestReturnsNull() {
        BaseDefaultFieldItem fieldItem = getFieldItemForShortArray();

        Short itemShort = fieldItem.getShort(7);
        assertNull(itemShort);
    }

    private static BaseDefaultFieldItem getFieldItemForShortArray() {
        int sizeIntByteBuffer = 8;
        ByteBuffer byteBuffer = ByteBuffer.allocate(sizeIntByteBuffer);
        byteBuffer.order(ByteOrder.BIG_ENDIAN);
        byteBuffer.putShort(0,(short)123);
        byteBuffer.putShort(2,(short)456);
        byteBuffer.putShort(4,(short)789);
        byteBuffer.putShort(6,(short)1011);
        Byte[][] byteArray = new Byte[sizeIntByteBuffer/2][2];
        int cntByteBuffer=0;
        for (byte b : byteBuffer.array()) {
            int shortIndex=cntByteBuffer/2;
            byteArray[shortIndex][cntByteBuffer%2] = b;
            cntByteBuffer++;
        }

        return new DefaultModbusByteArrayFieldItem(byteArray);
    }

    @Test
    public void convertByteArrayToLongTest() {
        BaseDefaultFieldItem fieldItem = getFieldItemForLongArray();

        Long itemLong = fieldItem.getLong(1);
        assertEquals(456789123L,itemLong,0);
    }

    @Test
    public void convertByteArrayToLongTestReturnsNull() {
        BaseDefaultFieldItem fieldItem = getFieldItemForLongArray();

        Long itemLong = fieldItem.getLong(4);
        assertNull(itemLong);
    }

    private static BaseDefaultFieldItem getFieldItemForLongArray() {
        int sizeIntByteBuffer = 32;
        ByteBuffer byteBuffer = ByteBuffer.allocate(sizeIntByteBuffer);
        byteBuffer.order(ByteOrder.BIG_ENDIAN);
        byteBuffer.putLong(0,123456789L);
        byteBuffer.putLong(8,456789123L);
        byteBuffer.putLong(16,789123456L);
        byteBuffer.putLong(24,101110111011L);
        Byte[][] byteArray = new Byte[sizeIntByteBuffer/2][2];
        int cntByteBuffer=0;
        for (byte b : byteBuffer.array()) {
            int shortIndex=cntByteBuffer/2;
            byteArray[shortIndex][cntByteBuffer%2] = b;
            cntByteBuffer++;
        }

        return new DefaultModbusByteArrayFieldItem(byteArray);
    }
}