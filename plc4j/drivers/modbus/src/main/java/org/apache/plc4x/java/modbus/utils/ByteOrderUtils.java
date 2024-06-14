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
package org.apache.plc4x.java.modbus.utils;

public class ByteOrderUtils {

    /**
     * Byte order is 'A B C D'
     * @param num int number
     * @return int number in big endian format
     */
    public static int toBigEndian(int num){
        return num;
    }

    /**
     * Byte order is 'D C B A'
     * @param num int number
     * @return int number in little endian format
     */
    public static int toLittleEndian(int num){
        return Integer.reverseBytes(num);
    }

    /**
     * Byte order is 'B A D C'
     * @param num int number
     * @return int number in big endian swap format
     */
    public static int toBigEndianWordSwap(int num){
        return (num << 16)|(num >>> 16);
    }

    /**
     * Byte order is 'C D A B'
     * @param num int number
     * @return int number in little endian swap format
     */
    public static int toLittleEndianWordSwap(int num){
        return ((num&0xff00)>>>8)|
                ((num<<8)&0xff00)|
                ((num<<8)&0xff000000)|
                ((num &0xff000000)>>>8);
    }

    /**
     * Byte order is 'A B C D E F G H'
     * @param num long number
     * @return long number in big endian format
     */
    public static long toBigEndian(long num){
        return num;
    }

    /**
     * Byte order is 'H G F E D C B A'
     * @param num long number
     * @return long number in little endian format
     */
    public static long toLittleEndian(long num){
        return Long.reverseBytes(num);
    }

    /**
     * Byte order is 'B A D C F E H G'
     * @param num long number
     * @return long number in big endian format
     */
    public static long toBigEndianWordSwap(long num){
        return (num & 0x00ff00ff00ff00ffL) << 8 | (num >>> 8) & 0x00ff00ff00ff00ffL;
    }

    /**
     * Byte order is 'G H E F C D A B'
     * @param num long number
     * @return long number in little endian format
     */
    public static long toLittleEndianWordSwap(long num){
        return (num & 0xffff000000000000L) >>> 48 |
               (num & 0x0000ffff00000000L) >>> 16 |
               (num & 0x00000000ffff0000L) << 16 |
               (num & 0x000000000000ffffL) << 48;
    }

}
