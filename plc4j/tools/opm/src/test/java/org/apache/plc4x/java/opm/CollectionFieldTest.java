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
package org.apache.plc4x.java.opm;

import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.lang.reflect.Field;
import java.util.Collection;
import java.util.Collections;
import java.util.List;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * A tag reading several values has to be converted to the element type of the annotated field -
 * see GH-1947. The values used to be handed out as the raw PlcValues, so
 * {@code List<Integer> l; l.get(0)} threw a ClassCastException, and an array field could not be
 * assigned at all.
 */
class CollectionFieldTest {

    /** The reporter's field: {@code @PlcTag("input-register:1:UINT[2]") List<Integer> registers;} */
    private List<Integer> registers;
    private Collection<Integer> asCollection;
    private Set<Integer> asSet;
    private int[] asPrimitiveArray;
    private Integer[] asBoxedArray;
    private List<String> asStrings;
    private List<?> withoutElementType;
    private java.util.Deque<Integer> unsupportedCollection;

    @Test
    void listOfIntegerYieldsIntegers() throws Exception {
        Object value = convert("registers");

        assertInstanceOf(List.class, value);
        List<?> list = (List<?>) value;
        assertEquals(2, list.size());
        // The point of the issue: the elements are Integers, not PlcUINTs.
        assertEquals(Integer.valueOf(1), list.get(0));
        assertEquals(Integer.valueOf(2), list.get(1));
        for (Object element : list) {
            assertInstanceOf(Integer.class, element);
        }
    }

    @Test
    void collectionAndSetAreAlsoSupported() throws Exception {
        assertInstanceOf(Collection.class, convert("asCollection"));
        assertInstanceOf(Set.class, convert("asSet"));
        assertEquals(2, ((Collection<?>) convert("asCollection")).size());
    }

    /**
     * The reporter's second finding: an array field could not be assigned at all.
     */
    @Test
    void primitiveArrayIsAssignable() throws Exception {
        Object value = convert("asPrimitiveArray");

        assertInstanceOf(int[].class, value);
        assertArrayEquals(new int[]{1, 2}, (int[]) value);
    }

    @Test
    void boxedArrayIsAssignable() throws Exception {
        Object value = convert("asBoxedArray");

        assertInstanceOf(Integer[].class, value);
        assertArrayEquals(new Integer[]{1, 2}, (Integer[]) value);
    }

    @Test
    void elementTypeDrivesTheConversion() throws Exception {
        Object value = convert("asStrings");

        List<?> list = (List<?>) value;
        assertInstanceOf(String.class, list.get(0));
    }

    /**
     * Without a resolvable element type the values come back as plain objects rather than failing.
     */
    @Test
    void unknownElementTypeFallsBackToObjects() throws Exception {
        Object value = convert("withoutElementType");

        assertEquals(2, ((List<?>) value).size());
    }

    @Test
    void unsupportedCollectionTypeIsReported() {
        assertThrows(ClassCastException.class, () -> convert("unsupportedCollection"));
    }

    /**
     * The whole way through OPM, as the reporter used it: a @PlcTag on a List field, read through
     * the entity manager.
     */
    @Test
    void entityWithAListFieldIsPopulated() throws Exception {
        DefaultPlcDriverManager driverManager = new DefaultPlcDriverManager();
        MockConnection connection = (MockConnection) driverManager.getConnection("mock:collections");
        MockDevice mockDevice = Mockito.mock(MockDevice.class);
        Mockito.when(mockDevice.read(Mockito.any())).thenReturn(
            new DefaultPlcResponseItem<>(PlcResponseCode.OK,
                new PlcList(List.of(new PlcUINT(1), new PlcUINT(2)))));
        connection.setDevice(mockDevice);

        RegisterEntity entity = new PlcEntityManager(driverManager)
            .read(RegisterEntity.class, "mock:collections");

        assertEquals(List.of(1, 2), entity.getRegisters());
        assertArrayEquals(new int[]{1, 2}, entity.getRaw());
    }

    @PlcEntity
    public static class RegisterEntity {

        @PlcTag("input-register:1:UINT[2]")
        private List<Integer> registers;

        @PlcTag("input-register:1:UINT[2]")
        private int[] raw;

        public RegisterEntity() {
            // For OPM
        }

        public List<Integer> getRegisters() {
            return registers;
        }

        public int[] getRaw() {
            return raw;
        }
    }

    /** Converts a two-element UINT response into the declared type of the named field. */
    private Object convert(String fieldName) throws NoSuchFieldException {
        Field field = CollectionFieldTest.class.getDeclaredField(fieldName);
        return PlcEntityInterceptor.getTyped(field.getType(), field.getGenericType(), response(), "tag");
    }

    private DefaultPlcReadResponse response() {
        PlcValue value = new PlcList(List.of(new PlcUINT(1), new PlcUINT(2)));
        PlcResponseItem<PlcValue> item = new DefaultPlcResponseItem<>(PlcResponseCode.OK, value);
        return new DefaultPlcReadResponse(null, Collections.singletonMap("tag", item));
    }
}
