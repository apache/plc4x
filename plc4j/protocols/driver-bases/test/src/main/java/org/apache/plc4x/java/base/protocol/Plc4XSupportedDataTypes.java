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
package org.apache.plc4x.java.base.protocol;

import org.apache.commons.lang3.ArrayUtils;
import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;

import java.io.Serializable;
import java.util.*;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Plc4XSupportedDataTypes {

    private final static Map<Class, Pair<? extends Serializable, byte[]>> littleEndianMap;
    private final static Map<Class, Pair<? extends Serializable, byte[]>> bigEndianMap;

    static {
        Calendar calenderInstance = Calendar.getInstance();
        calenderInstance.setTime(new Date(283686951976960L));
        littleEndianMap = new HashMap<>();
        littleEndianMap.put(Boolean.class, ImmutablePair.of(Boolean.TRUE, new byte[]{0x01}));
        littleEndianMap.put(Byte.class, ImmutablePair.of(Byte.valueOf("1"), new byte[]{0x1}));
        littleEndianMap.put(Short.class, ImmutablePair.of(Short.valueOf("1"), new byte[]{0x1, 0x0}));
        littleEndianMap.put(Calendar.class, ImmutablePair.of(calenderInstance, new byte[]{0x0, 0x0, 0x0, 0x0, 0x4, 0x3, 0x2, 0x1}));
        littleEndianMap.put(Float.class, ImmutablePair.of(Float.valueOf("1"), new byte[]{0x0, 0x0, (byte) 0x80, 0x3F}));
        littleEndianMap.put(Double.class, ImmutablePair.of(Double.valueOf("1"), new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, (byte) 0xF0, 0x3F}));
        littleEndianMap.put(Integer.class, ImmutablePair.of(Integer.valueOf("1"), new byte[]{0x1, 0x0, 0x0, 0x0}));
        littleEndianMap.put(String.class, ImmutablePair.of(String.valueOf("Hello World!"), new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}));
        bigEndianMap = new HashMap<>();
        littleEndianMap.forEach((clazz, pair) -> {
            Serializable serializable = pair.getLeft();
            byte[] littleEndianBytes = pair.getRight();
            byte[] bigEndianBytes = ArrayUtils.clone(littleEndianBytes);
            ArrayUtils.reverse(bigEndianBytes);
            bigEndianMap.put(clazz, ImmutablePair.of(serializable, bigEndianBytes));
        });
    }

    /**
     * A {@link Stream} of {@link Class}es plc4x can currently support.
     *
     * @return a stream of supported data types.
     */
    public static Stream<Class<? extends Serializable>> streamOfPlc4XSupportedDataTypes() {
        return Stream.of(
            Boolean.class,
            Byte.class,
            Short.class,
            Calendar.class,
            Float.class,
            Double.class,
            Integer.class,
            String.class
        );
    }

    /**
     * A {@link Stream} of instances of {@link Class}es plc4x can currently support with their according byte representation.
     *
     * @return a stream of {@link org.apache.commons.lang3.tuple.Pair}s of instances and their byte values.
     * @see #streamOfPlc4XSupportedDataTypes
     */
    public static Stream<? extends Pair<? extends Serializable, byte[]>> streamOfLittleEndianDataTypePairs() {
        return streamOfLittleEndianDataTypePairs(streamOfPlc4XSupportedDataTypes());
    }

    /**
     * A {@link Stream} of instances of {@link Class}es which are defined by {@code inputStream} can currently support with their according byte representation.
     *
     * @param inputStream a stream of {@link org.apache.commons.lang3.tuple.Pair}s of instances and their byte values.
     * @see #streamOfPlc4XSupportedDataTypes
     */
    public static Stream<? extends Pair<? extends Serializable, byte[]>> streamOfLittleEndianDataTypePairs(Stream<Class<? extends Serializable>> inputStream) {
        return inputStream
            .map(littleEndianMap::get)
            .peek(Objects::requireNonNull);
    }

    /**
     * A {@link Stream} of instances of {@link Class}es plc4x can currently support with their according little endian byte representation.
     *
     * @return a stream of {@link org.apache.commons.lang3.tuple.Pair}s of instances and their byte values.
     * @see #streamOfPlc4XSupportedDataTypes
     */
    public static Stream<? extends Pair<? extends Serializable, byte[]>> streamOfBigEndianDataTypePairs() {
        return streamOfBigEndianDataTypePairs(streamOfPlc4XSupportedDataTypes());
    }

    /**
     * A {@link Stream} of instances of {@link Class}es which are defined by {@code inputStream} can currently support with their according big endian byte representation.
     *
     * @param inputStream a stream of {@link org.apache.commons.lang3.tuple.Pair}s of instances and their byte values.
     * @see #streamOfPlc4XSupportedDataTypes
     */
    public static Stream<? extends Pair<? extends Serializable, byte[]>> streamOfBigEndianDataTypePairs(Stream<Class<? extends Serializable>> inputStream) {
        return inputStream
            .map(bigEndianMap::get)
            .peek(Objects::requireNonNull);
    }
}
